package library

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ra341/glacier/internal/downloader/types"
)

type Downloader interface {
	Add(ctx context.Context, gameId uint, download types.Download) (err error)
}

type Service struct {
	downloader Downloader
	store      Store
}

func New(store Store, downloader Downloader) *Service {
	return &Service{
		downloader: downloader,
		store:      store,
	}
}

func (s *Service) List(ctx context.Context, offset, limit uint) ([]Game, error) {
	return s.store.List(ctx, limit, offset)
}

func (s *Service) Add(ctx context.Context, game *Game) error {
	game.Download.State = types.DownloadQueued
	abs, err := filepath.Abs(game.Download.DownloadPath)
	if err != nil {
		return fmt.Errorf("could not resolve absolute path for download: %w", err)
	}
	game.Download.DownloadPath = abs

	err = s.store.Add(ctx, game)
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
