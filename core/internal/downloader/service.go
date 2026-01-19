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

	isDownloadTrackerRunning atomic.Bool
	trackerCtxCancelFn       context.CancelFunc
	triggerChan              chan struct{}
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

	downloadId, err := downloader.Download(ctx, download.DownloadUrl, s.conf().IncompletePath)
	if err != nil {
		return err
	}

	err = s.store.UpdateDownloadProgress(ctx, gameId,
		types.Download{
			DownloadId: downloadId,
			State:      types.Downloading,
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
	s.triggerChan = make(chan struct{}, 1)

	go s.downloadTracker(ctx)
}

func (s *Service) TriggerTracker() {
	select {
	case s.triggerChan <- struct{}{}:
		log.Debug().Msg("tracker triggered")
	default:
		log.Debug().Msg("tracker is already triggered")
	}
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
	var done bool

	for {
		select {
		case t := <-timer.C:
			log.Debug().Time("time", t).Msg("checking downloads")
			errTries, done = s.trackDownloader(ctx, errTries)
			if done {
				return
			}
		case <-s.triggerChan:
			errTries, done = s.trackDownloader(ctx, errTries)
			if done {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) trackDownloader(ctx context.Context, errTries int) (int, bool) {
	downloading, err := s.store.ListDownloadState(ctx, types.Downloading)
	if err != nil {
		log.Warn().Err(err).Msg("could not check downloads")
		errTries++
		if errTries > 10 {
			log.Warn().Msg("could not list downloads, something is wrong with the DB")
			return errTries, true
		}
	}
	errTries = 0

	if len(downloading) < 1 {
		log.Warn().Msg("no active downloads found")
		return errTries, true
	}

	for _, dn := range downloading {
		s.checkDownload(ctx, &dn)
	}

	return errTries, false
}

// checks a single download
func (s *Service) checkDownload(ctx context.Context, dn *library.Game) {
	download := dn.Download

	downloader, err := s.cli(download.Client)
	if err != nil {
		dn.Download.State = types.Error
		dn.Download.Progress = fmt.Sprintf("%v", err)
	}

	downId := download.DownloadId
	err = downloader.Progress(ctx, &dn.Download)
	if err != nil {
		dn.Download.State = types.Downloading
		dn.Download.Progress = fmt.Sprintf("unable to get progress: %v", err)
	}

	err = s.store.UpdateDownloadProgress(ctx, dn.ID, dn.Download)
	if err != nil {
		log.Warn().Err(err).Str("download", downId).Msg("failed to update download state")
	}
}
