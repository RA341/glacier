package download

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	hc "github.com/ra341/glacier/frost/http_client"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

const MB = 1024 * 1024

type Service struct {
	Config             config.Provider[Config]
	editStatus         EditStatus
	ActiveDownloads    syncmap.Map[int, *Download]
	DownloadTotalBytes *uint64

	workers      *errgroup.Group
	SpeedLimiter *hc.LimiterService
}

// New
//
// basepath must be: "http://localhost:6699"
func New(baseurl string, config config.Provider[Config], httpCliFactory hc.HttpCliFactory, editStatus EditStatus, ls *hc.LimiterService, downloadSpeedCounter *uint64) *Service {
	transport := &http.Transport{
		// MaxIdleConns is the total connections across all hosts
		MaxIdleConns: 100,

		// MaxIdleConnsPerHost must be >= your worker count.
		// The default is only 2 If you have 30 workers, 28 will
		// constantly create new TCP connections.
		MaxIdleConnsPerHost: 100,

		// Time to keep an idle connection in the pool
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	config().httpCli = httpCliFactory(transport)
	config().base = fmt.Sprintf(
		"%s/library/download",
		strings.TrimRight(baseurl, "/"),
	)

	group := errgroup.Group{}
	group.SetLimit(config().MaxConcurrentGames)

	return &Service{
		Config:             config,
		SpeedLimiter:       ls,
		editStatus:         editStatus,
		DownloadTotalBytes: downloadSpeedCounter,
		workers:            &group,
	}
}

func (d *Service) Download(ctx context.Context, gameId int, downloadPath string, force bool) error {
	if _, ok := d.ActiveDownloads.Load(gameId); ok {
		log.Debug().Msg("download already in progress")
		return nil
	}

	err := d.editStatus(ctx, gameId, &LocalDownload{
		Status:  StatusQueued,
		Started: time.Now(),
		Done:    time.Time{},
	})
	if err != nil {
		return fmt.Errorf("could not update status %w", err)
	}

	download, err := NewDownload(
		d.Config(),
		d.editStatus,
		d.onDone,
		gameId,
		downloadPath,
		force,
	)
	if err != nil {
		return fmt.Errorf("could not start download: %w", err)
	}

	go func() {
		d.workers.Go(func() error {
			download.Start()
			return nil
		})
	}()

	d.ActiveDownloads.Store(gameId, download)

	return nil
}

func (d *Service) Pause(gameId int) error {
	val, ok := d.ActiveDownloads.Load(gameId)
	if !ok {
		return fmt.Errorf("game not found " + strconv.Itoa(gameId))
	}

	val.Pause()
	return nil
}

func (d *Service) Resume(gameId int) error {
	val, ok := d.ActiveDownloads.Load(gameId)
	if !ok {
		return fmt.Errorf("game not found " + strconv.Itoa(gameId))
	}

	val.Resume()
	go func() {
		// add it back to the error group for download limiting
		d.workers.Go(func() error {
			val.Start()
			return nil
		})
	}()

	return nil
}

func (d *Service) Cancel(gameId int) error {
	val, ok := d.ActiveDownloads.LoadAndDelete(gameId)
	if !ok {
		return fmt.Errorf("game not found " + strconv.Itoa(gameId))
	}

	val.Cancel()

	if val.isPaused.Load() {
		// release resources if paused and canceled was called
		// otherwise it should be released by defer in val.Start()
		val.Close()
	}

	log.Info().Msg("download stopped")

	err := os.RemoveAll(val.downloadFolder)
	if err != nil {
		return fmt.Errorf("could not remove download folder: %w", err)
	}

	return nil
}

func (d *Service) onDone(gameId int) {
	d.ActiveDownloads.Delete(gameId)
}
