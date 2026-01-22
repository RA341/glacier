package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
	"gorm.io/gorm"
)

type StoreGorm struct {
	db *gorm.DB
}

func NewStoreGorm(db *gorm.DB) Store {
	return &StoreGorm{db: db}
}

// List returns a paginated list of games.
func (s *StoreGorm) List(ctx context.Context, query string, limit int, offset int) ([]LocalGame, error) {
	var games []LocalGame
	err := s.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&games).Error
	return games, err
}

func (s *StoreGorm) ListWithState(ctx context.Context, status download.Status) ([]LocalGame, error) {
	var games []LocalGame
	err := s.db.WithContext(ctx).
		Where("status = ?", status).
		Find(&games).Error
	return games, err
}

func (s *StoreGorm) Add(ctx context.Context, game *LocalGame) error {
	return s.db.WithContext(ctx).Create(game).Error
}

func (s *StoreGorm) Get(ctx context.Context, id int) (LocalGame, error) {
	var game LocalGame
	err := s.db.WithContext(ctx).First(&game, id).Error
	return game, err
}
func (s *StoreGorm) Edit(ctx context.Context, id int, game *LocalGame) error {
	game.ID = uint(id)
	return s.db.WithContext(ctx).Save(game).Error
}

func (s *StoreGorm) Delete(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&LocalGame{}, id).Error
}

func (s *StoreGorm) EditStatus(ctx context.Context, id int, Status download.Status, StatusMessage string) error {
	return s.db.WithContext(ctx).Model(&LocalGame{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         Status,
		"status_message": StatusMessage,
	}).Error
}
