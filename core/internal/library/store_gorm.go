package library

import (
	"context"

	"github.com/ra341/glacier/internal/downloader/types"
	"gorm.io/gorm"
)

type StoreGorm struct {
	gormDB *gorm.DB
}

func NewStoreGorm(gormDB *gorm.DB) *StoreGorm {
	return &StoreGorm{
		gormDB: gormDB,
	}
}

func (s *StoreGorm) Q(ctx context.Context) *gorm.DB {
	return s.gormDB.WithContext(ctx)
}

func (s *StoreGorm) ListDownloadState(ctx context.Context, state types.DownloadState) ([]Game, error) {
	var downloads []Game

	err := s.gormDB.Where("state = ?", state).Find(&downloads).Error

	return downloads, err
}

func (s *StoreGorm) Add(ctx context.Context, game *Game) error {
	return s.Q(ctx).Save(game).Error
}

func (s *StoreGorm) List(ctx context.Context, limit uint, offset uint) ([]Game, error) {
	var games []Game
	err := s.Q(ctx).
		Order("updated_at desc").Offset(int(offset)).
		Limit(int(limit)).Find(&games).
		Error
	return games, err
}

func (s *StoreGorm) UpdateDownloadProgress(ctx context.Context, id uint, download types.Download) error {
	return s.Q(ctx).
		Model(&Game{
			Model: gorm.Model{ID: id},
		}).
		Select(
			`download_id`,
			`state`,
			`progress`,
			`incomplete_path`,
			`left`,
			`complete`,
		).
		UpdateColumns(Game{
			Download: download,
		}).Error
}

func (s *StoreGorm) GetById(ctx context.Context, id uint) (Game, error) {
	var game Game
	err := s.Q(ctx).First(&game, id).Error
	return game, err
}

func (s *StoreGorm) DeleteGame(ctx context.Context, id uint) error {
	return s.Q(ctx).Unscoped().Delete(&Game{}, id).Error
}
