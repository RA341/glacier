package transmission

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/mapsct"

	"github.com/hekmon/transmissionrpc/v3"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Protocol string
	Host     string
	User     string
	Password string
}

type Client struct {
	cli *transmissionrpc.Client
}

func New(config types.ClientConfig) (types.Downloader, error) {
	var conf Config

	err := mapsct.ParseMap(&conf, config)
	if err != nil {
		return nil, err
	}

	clientStr := ""
	if conf.User != "" && conf.Password != "" {
		clientStr = fmt.Sprintf("%s://%s:%s@%s/transmission/rpc", conf.Protocol, conf.User, conf.Password, conf.Host)
	} else {
		clientStr = fmt.Sprintf("%s://%s/transmission/rpc", conf.Protocol, conf.Host)
	}

	endpoint, err := url.Parse(clientStr)
	if err != nil {
		return nil, err
	}

	tbt, err := transmissionrpc.New(endpoint, nil)
	if err != nil {
		panic(err)
	}

	transmission := &Client{cli: tbt}

	_, _, err = transmission.test()
	if err != nil {
		return nil, fmt.Errorf("transmission client health check failed: %s", err)
	}

	return transmission, nil
}

func (tm *Client) Type() types.ClientType {
	return types.ClientTransmission
}

func (tm *Client) test() (string, string, error) {
	ok, serverVersion, serverMinimumVersion, err := tm.cli.RPCVersion(context.Background())
	if err != nil {
		return "", "", err
	}

	if !ok {
		return "", "", fmt.Errorf("remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
	}

	return "transmission", fmt.Sprintf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
		serverVersion, transmissionrpc.RPCVersion), nil
}

func (tm *Client) Download(ctx context.Context, url string, downloadPath string) (downloadID string, err error) {
	if strings.HasPrefix(url, "magnet:?") {
		log.Debug().Str("url", url[:12]).Msg("Downloading with magnet magnet link")
		return tm.downloadMagnet(ctx, url, downloadPath, "")
	}

	log.Debug().Str("url", url).Msg("Downloading with torrent file")
	return tm.downloadFile(url, downloadID, "")
}

func (tm *Client) downloadFile(torrent, downloadPath, _ string) (string, error) {
	torrentResult, err := tm.cli.TorrentAddFileDownloadDir(context.Background(), torrent, downloadPath)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(*torrentResult.ID, 10), nil
}

func (tm *Client) downloadMagnet(ctx context.Context, magnet, downloadPath, category string) (string, error) {
	torrentResult, err := tm.cli.TorrentAdd(ctx, transmissionrpc.TorrentAddPayload{
		Filename:    &magnet,
		DownloadDir: &downloadPath,
		//Labels:      []string{category},
	})
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(*torrentResult.ID, 10), nil
}

func (tm *Client) Cancel(ctx context.Context, downloadId string, removeDownloaded bool) error {
	return fmt.Errorf(" this function is not implemented you idiot")
}

func (tm *Client) Progress(ctx context.Context, download *types.Download) (err error) {
	var intIds = make([]int64, 1)
	finalId, err := strconv.Atoi(download.DownloadId)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to convert torrent id: %s to int64", download.DownloadId)
		return err
	}
	intIds = append(intIds, int64(finalId))

	infos, err := tm.cli.TorrentGetAllFor(ctx, intIds)
	if err != nil {
		return err
	}
	if len(infos) == 0 {
		return fmt.Errorf("no torrent not found")
	}

	for _, info := range infos {
		download.State = types.Downloading
		if info.Status.String() == "seeding" {
			download.State = types.Complete
		} else if info.Status.String() == "error" ||
			info.Status.String() == "stopped" {
			download.State = types.Error
		}

		download.Progress = fmt.Sprintf(
			"%s\nprogress:%v",
			info.Status.String(),
			info.PercentDone,
		)

		break
	}

	return nil
}
