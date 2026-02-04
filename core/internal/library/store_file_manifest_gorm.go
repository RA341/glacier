package library

import (
	"context"
	"fmt"

	download "github.com/ra341/glacier/internal/downloader/types"
	"gorm.io/gorm"
)

type StoreFolderMetadataGorm struct {
	DB *gorm.DB
}

func NewStoreFolderMetadataGorm(gorm *gorm.DB) StoreGameManifest {
	return &StoreFolderMetadataGorm{
		DB: gorm,
	}
}

func (s *StoreFolderMetadataGorm) Q(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

func (s *StoreFolderMetadataGorm) Add(ctx context.Context, gameId int, metadata *FolderManifest) error {
	metadata.GameID = gameId
	err := s.Q(ctx).Save(&metadata).Error
	return err
}

func (s *StoreFolderMetadataGorm) Get(ctx context.Context, gameId int) (FolderManifest, error) {
	var metadata FolderManifest
	err := s.Q(ctx).
		Where("game_id = ?", gameId).
		First(&metadata).
		Error
	return metadata, err
}

func (s *StoreFolderMetadataGorm) Delete(ctx context.Context, gameId int) error {
	return s.Q(ctx).
		Unscoped().
		Where("game_id = ?", gameId).
		Delete(&FolderManifest{}).
		Error
}

func (s *StoreFolderMetadataGorm) ListGamesWithoutManifest(ctx context.Context) ([]int, error) {
	var games []int

	err := s.Q(ctx).
		Model(&Game{}).
		Joins("LEFT JOIN folder_manifests ON folder_manifests.game_id = games.id").
		Where("folder_manifests.id IS NULL").
		Where("games.state = ?", download.Complete).
		Pluck("games.id", &games).
		Error

	return games, err
}

func (s *StoreFolderMetadataGorm) Edit(ctx context.Context, gameId int, metadata *FolderManifest) error {
	return fmt.Errorf("edit manifest is unimplemented")
}
