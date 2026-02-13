package qbit

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/autobrr/go-qbittorrent"
	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/magnet"
	"github.com/ra341/glacier/pkg/mapsct"
)

type Qbit struct {
	cli *qbittorrent.Client
}

type Config struct {
	Host     string
	User     string
	Password string
}

func New(in map[string]any) (types.Downloader, error) {
	var conf Config

	err := mapsct.ParseMap(&conf, in)
	if err != nil {
		return nil, err
	}

	cli := qbittorrent.NewClient(qbittorrent.Config{
		Host:     conf.Host,
		Username: conf.User,
		Password: conf.Password,
		Timeout:  10,
	})

	err = cli.Login()
	if err != nil {
		return nil, fmt.Errorf("could not login to qbit")
	}

	return &Qbit{
		cli: cli,
	}, nil
}

func (q *Qbit) Download(ctx context.Context, url string, downloadPath string) (downloadID string, err error) {
	magnetLink := url
	if strings.HasPrefix(url, "http") {
		// it's not a magnet link download it first.
		magnetLink, err = q.downloadAndConvertToMagnet(ctx, url, downloadPath)
		if err != nil {
			return "", err
		}
	}

	err = q.cli.AddTorrentFromUrlCtx(ctx, url, map[string]string{
		"savepath": downloadPath,
	})
	if err != nil {
		return "", err
	}

	component, err := magnet.DecodeMagnetURL(magnetLink)
	if err != nil {
		return "", err
	}

	return component.InfoHash, nil
}

func (q *Qbit) downloadAndConvertToMagnet(ctx context.Context, url string, downloadPath string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download torrent file from %s: %w", url, err)
	}
	defer fileutil.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download torrent file from %s: status %s", url, resp.Status)
	}

	return magnet.TorrentFileToMagnet(resp.Body)
}

func (q *Qbit) Cancel(ctx context.Context, downloadId string, removeDownloaded bool) error {
	return q.cli.DeleteTorrentsCtx(ctx, []string{downloadId}, removeDownloaded)
}

func (q *Qbit) Progress(ctx context.Context, download *types.Download) (err error) {
	torrents, err := q.cli.GetTorrentsCtx(ctx, qbittorrent.TorrentFilterOptions{
		Hashes: []string{download.DownloadId},
	})
	if err != nil {
		return err
	}
	if len(torrents) < 1 {
		return fmt.Errorf("torrent not found with hash: %s", download.DownloadId)
	}

	info := torrents[0]

	switch info.State {
	// this is fucking annoying to look at maybe refactor
	case qbittorrent.TorrentStateUploading,
		qbittorrent.TorrentStateStalledUp,
		qbittorrent.TorrentStatePausedUp,
		qbittorrent.TorrentStateQueuedUp,
		qbittorrent.TorrentStateForcedUp,
		qbittorrent.TorrentStateCheckingUp:
		download.State = types.Complete
	case qbittorrent.TorrentStateError,
		qbittorrent.TorrentStateMissingFiles,
		qbittorrent.TorrentStatePausedDl,
		qbittorrent.TorrentStateStoppedDl:
		download.State = types.Error
	default:
		download.State = types.Downloading
	}

	if info.Progress >= 1 {
		download.State = types.Complete
	}

	download.Complete = uint64(info.Downloaded)
	download.Left = uint64(info.AmountLeft)
	download.Progress = fmt.Sprintf("ETA: %v - State: %s", time.Duration(info.ETA)*time.Second, info.State)

	download.IncompletePath = filepath.Join(info.SavePath, info.Name)

	return nil
}
