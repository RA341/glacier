package library

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/internal/library/manifest"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/shared/config"
)

type Downloader interface {
	Add(ctx context.Context, gameId *Game) (err error)
	TriggerTracker()
}

type Service struct {
	config         config.Provider[Config]
	downloader     Downloader
	downloadAssets DownloadAssetTrigger

	store    Store
	manifest *manifest.Service
}

type DownloadAssetTrigger func(ctx context.Context, gameId uint) (err error)

func New(
	store Store,
	fs *manifest.Service,
	downloader Downloader,
	config config.Provider[Config],
	assets DownloadAssetTrigger,
) *Service {
	s := &Service{
		downloader:     downloader,
		config:         config,
		downloadAssets: assets,
		store:          store,
		manifest:       fs,
	}
	return s
}

func (s *Service) GetDownloadManifest(ctx context.Context, gid int, w http.ResponseWriter) error {
	game, err := s.store.GetById(ctx, uint(gid))
	if err != nil {
		return err
	}

	if game.Download.State != types.Complete {
		return fmt.Errorf("game is not complete")
	}

	return s.manifest.GetGameManifest(ctx, gid, game.Download.DownloadPath, w)
}

func (s *Service) FileDownload(ctx context.Context, id int, file string) (*os.File, error) {
	game, err := s.store.GetById(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(game.Download.DownloadPath, file)
	return os.Open(filePath)
}

func (s *Service) Get(ctx context.Context, id uint) (Game, error) {
	return s.store.GetById(ctx, id)
}

func (s *Service) Edit(ctx context.Context, game *Game) error {
	err := s.store.Edit(ctx, game)
	if err != nil {
		return err
	}

	err = s.downloadAssets(ctx, game.ID)
	if err != nil {
		return err
	}

	return err
}

func (s *Service) List(ctx context.Context, query string, offset, limit uint) ([]Game, error) {
	return s.store.List(ctx, query, limit, offset)
}

func (s *Service) ListDownloading(ctx context.Context, state string) ([]Game, error) {
	dState, err := types.DownloadStateString(state)
	if err != nil {
		return nil, err
	}

	return s.store.ListDownloadState(ctx, dState)
}

func (s *Service) Add(ctx context.Context, game *Game) error {
	game.Download.State = types.Queued
	game.Download.DownloadPath = filepath.Join(
		s.config().GameDir,
		filepath.Clean(game.Meta.Name),
	)

	err := s.store.Add(ctx, game)
	if err != nil {
		return err
	}

	err = s.downloader.Add(ctx, game)
	if err != nil {
		game.SetErr(err)
		return s.store.UpdateDownloadProgress(ctx, game.ID, game.Download)
	}

	err = s.downloadAssets(ctx, game.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	err := checkPerms(ctx)
	if err != nil {
		return err
	}

	return s.store.Delete(ctx, id)
}

func (s *Service) ListFiles(ctx context.Context, id uint, downloaded bool, base string) ([]os.DirEntry, error) {
	readPath, err := s.getFolder(ctx, id, downloaded)
	if err != nil {
		return nil, err
	}

	if base != "" {
		readPath = filepath.Join(readPath, base)
	}

	return os.ReadDir(readPath)
}

func (s *Service) DeleteFile(ctx context.Context, id uint, base string, downloaded bool) error {
	readPath, err := s.getFolder(ctx, id, downloaded)
	if err != nil {
		return err
	}

	if base == "" {
		return fmt.Errorf("no file name found")
	}

	readPath = filepath.Join(readPath, base)

	return os.RemoveAll(readPath)
}

func (s *Service) getFolder(ctx context.Context, id uint, downloaded bool) (string, error) {
	game, err := s.store.GetById(ctx, id)
	if err != nil {
		return "", err
	}

	var readPath = game.Download.DownloadPath
	if downloaded {
		readPath = game.Download.IncompletePath
	}
	return readPath, nil
}

func (s *Service) Redownload(ctx context.Context, gameId uint) error {
	game, err := s.store.GetById(ctx, gameId)
	if err != nil {
		return err
	}

	err = s.downloader.Add(ctx, &game)
	if err != nil {
		game.SetErr(err)
		return s.store.UpdateDownloadProgress(ctx, game.ID, game.Download)
	}

	return nil
}

func checkPerms(ctx context.Context) error {
	userInf, err := user.GetUserCtx(ctx)
	if err != nil {
		return err
	}

	if userInf.Role > user.Magos {
		return fmt.Errorf("insufficient permissions to delete")
	}
	return nil
}
