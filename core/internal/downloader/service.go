package downloader

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/internal/library"
	"github.com/rs/zerolog/log"
)

type GetCli func(clientType types.ClientType) (types.Downloader, error)

type Service struct {
	cli   GetCli
	store library.Store
	conf  ConfigLoader

	// these fields are handled internally no need to passed in constructor
	isDownloadTrackerRunning atomic.Bool
	trackerCtxCancelFn       context.CancelFunc
}

func New(cli GetCli, store library.Store, conf ConfigLoader) *Service {
	return &Service{
		cli:   cli,
		store: store,
		conf:  conf,
	}
}

func (s *Service) Add(ctx context.Context, gameId uint, download types.Download) (err error) {
	downloader, err := s.cli(download.Client)
	if err != nil {
		return err
	}

	downloadId, err := downloader.Download(ctx, download.DownloadUrl, s.conf().IncompleteDownloadPath)
	if err != nil {
		return err
	}

	err = s.store.UpdateDownloadProgress(ctx, gameId,
		types.Download{
			DownloadId: downloadId,
			State:      types.DownloadDownloading,
		},
	)
	if err != nil {
		return err
	}

	s.StartTracker()

	return err
}

func (s *Service) Cancel(ctx context.Context, client types.ClientType, downloadID string, removeFiles bool) error {
	downloader, err := s.cli(client)
	if err != nil {
		return err
	}
	return downloader.Cancel(ctx, downloadID, removeFiles)
}

func (s *Service) StartTracker() {
	if !s.isDownloadTrackerRunning.CompareAndSwap(false, true) {
		log.Debug().Msg("download tracker is running")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.trackerCtxCancelFn = cancel

	go s.downloadTracker(ctx)
}

func (s *Service) StopTracker() {
	if s.trackerCtxCancelFn == nil {
		log.Warn().Msg("cancel fn is nil")
		return
	}

	s.trackerCtxCancelFn()
}

func (s *Service) downloadTracker(ctx context.Context) {
	defer func() {
		log.Info().Msg("download tracker has been stopped")
	}()

	s.isDownloadTrackerRunning.Store(true)
	defer s.isDownloadTrackerRunning.Store(false)

	timer := time.NewTimer(s.conf().Interval())
	defer timer.Stop()

	errTries := 0

	for {
		select {
		case t := <-timer.C:
			log.Debug().Time("time", t).Msg("checking downloads")

			downloading, err := s.store.ListDownloadState(ctx, types.DownloadDownloading)
			if err != nil {
				log.Warn().Err(err).Msg("could not check downloads")
				errTries++
				if errTries > 10 {
					log.Warn().Msg("could not list downloads, something is wrong with the DB")
					return
				}
			}
			errTries = 0

			if len(downloading) < 1 {
				log.Warn().Msg("no active downloads found")
				return
			}

			for _, dn := range downloading {
				s.checkDownload(ctx, &dn)
			}

		case <-ctx.Done():
			return
		}
	}
}

// checks a single download
func (s *Service) checkDownload(ctx context.Context, dn *library.Game) {
	download := dn.Download

	downloader, err := s.cli(download.Client)
	if err != nil {
		dn.Download.State = types.DownloadError
		dn.Download.Progress = fmt.Sprintf("%v", err)
	}

	downId := download.DownloadId
	err = downloader.Progress(ctx, &dn.Download)
	if err != nil {
		dn.Download.State = types.DownloadDownloading
		dn.Download.Progress = fmt.Sprintf("unable to get progress: %v", err)
	}

	err = s.store.UpdateDownloadProgress(ctx, dn.ID, dn.Download)
	if err != nil {
		log.Warn().Err(err).Str("download", downId).Msg("failed to update download state")
	}
}
