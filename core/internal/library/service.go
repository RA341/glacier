package library

import (
	"context"
	"fmt"
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
