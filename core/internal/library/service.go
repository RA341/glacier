package library

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/internal/user"
)

type Downloader interface {
	Add(ctx context.Context, gameId *Game) (err error)
	TriggerTracker()
}

type Service struct {
	config     ConfigLoader
	downloader Downloader

	store    Store
	manifest *ManifestService
}

type ConfigLoader func() *Config

func New(
	store Store,
	fs *ManifestService,
	downloader Downloader,
	config ConfigLoader,
) *Service {

	return &Service{
		downloader: downloader,
		config:     config,
		store:      store,
		manifest:   fs,
	}
}

func (s *Service) Get(ctx context.Context, id uint) (Game, error) {
	return s.store.GetById(ctx, id)
}

func (s *Service) Edit(ctx context.Context, game *Game) error {
	err := checkPerms(ctx)
	if err != nil {
		return err
	}

	return s.store.Edit(ctx, game)
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
