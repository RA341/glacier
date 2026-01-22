package downloader

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

func (s *Service) Add(ctx context.Context, game *library.Game) (err error) {
	downloader, err := s.cli(game.Download.Client)
	if err != nil {
		return err
	}

	downloadId, err := downloader.Download(
		ctx,
		game.Download.DownloadUrl,
		s.conf().IncompletePath,
	)
	if err != nil {
		return err
	}
	game.Download.State = types.Downloading
	game.Download.DownloadId = downloadId

	err = s.store.UpdateDownloadProgress(ctx, game.ID, game.Download)
	if err != nil {
		return err
	}

	s.StartTracker()
	s.TriggerTracker() // get the initial status

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
	defer func() {
		err := s.store.UpdateDownloadProgress(ctx, dn.ID, dn.Download)
		if err != nil {
			log.Warn().Err(err).Str("download", dn.Download.DownloadId).Msg("failed to update download state")
		}
	}()
	download := dn.Download

	downloader, err := s.cli(download.Client)
	if err != nil {
		dn.Download.State = types.Error
		dn.Download.Progress = fmt.Sprintf("%v", err)
		return
	}

	err = downloader.Progress(ctx, &dn.Download)
	if err != nil {
		dn.Download.State = types.Downloading
		dn.Download.Progress = fmt.Sprintf("unable to get progress: %v", err)
		return
	}

	if dn.Download.State == types.Complete {
		s.completeGameDownload(dn)
		return
	}
}

func (s *Service) completeGameDownload(game *library.Game) {
	err := os.MkdirAll(game.Download.DownloadPath, os.ModePerm)
	if err != nil {
		game.Download.State = types.Error
		game.Download.Progress = fmt.Sprintf("%v", err)
		return
	}

	err = filepath.WalkDir(game.Download.IncompletePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if game.Download.IncompletePath == path {
			return nil
		}

		// relative path from src dir
		relPath, err := filepath.Rel(game.Download.IncompletePath, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(game.Download.DownloadPath, relPath)

		if d.IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}

		parentDir := filepath.Dir(targetPath)
		err = os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = os.Link(path, targetPath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		game.Download.State = types.Error
		game.Download.Progress = fmt.Sprintf("could link file: %v", err)
		return
	}

	log.Debug().Str("title", game.Meta.Name[:24]).Msg("download complete with no errors")
}
