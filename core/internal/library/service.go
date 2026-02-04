package library

import (
	"context"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
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
	return s.store.DeleteGame(ctx, id)
}

func (s *Service) Get(ctx context.Context, id uint) (Game, error) {
	return s.store.GetById(ctx, id)
}
