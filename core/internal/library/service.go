package library

import (
	"context"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
)

type Downloader interface {
	Add(ctx context.Context, gameId uint, download types.Download) (err error)
}

type Service struct {
	config     ConfigLoader
	downloader Downloader
	store      Store
}

type ConfigLoader func() *Config

func New(store Store, downloader Downloader, config ConfigLoader) *Service {
	return &Service{
		downloader: downloader,
		store:      store,
		config:     config,
	}
}

func (s *Service) List(ctx context.Context, offset, limit uint) ([]Game, error) {
	return s.store.List(ctx, limit, offset)
}

func (s *Service) Add(ctx context.Context, game *Game) error {
	game.Download.State = types.DownloadQueued
	game.Download.DownloadPath = filepath.Join(
		s.config().GameDir,
		filepath.Clean(game.Meta.Name),
	)

	err := s.store.Add(ctx, game)
	if err != nil {
		return err
	}

	return s.downloader.Add(ctx, game.ID, game.Download)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.store.DeleteGame(ctx, id)
}

func (s *Service) Get(ctx context.Context, id uint) (Game, error) {
	return s.store.GetById(ctx, id)
}
