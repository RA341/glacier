package library

import (
	"context"

	"gorm.io/gorm"
)

type StoreFolderMetadataGorm struct {
	DB *gorm.DB
}

func NewStoreFolderMetadataGorm(gorm *gorm.DB) *StoreFolderMetadataGorm {
	return &StoreFolderMetadataGorm{
		DB: gorm,
	}
}

func (s *StoreFolderMetadataGorm) Q(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}

func (s *StoreFolderMetadataGorm) Add(ctx context.Context, gameId int, metadata *FolderMetadata) error {
	metadata.GameID = gameId
	err := s.Q(ctx).Save(&metadata).Error
	return err
}

func (s *StoreFolderMetadataGorm) Get(ctx context.Context, gameId int) (FolderMetadata, error) {
	var metadata FolderMetadata
	err := s.Q(ctx).Where("game_id = ?", gameId).First(&metadata, gameId).Error
	return metadata, err
}

func (s *StoreFolderMetadataGorm) Delete(ctx context.Context, gameId int) error {
	return s.Q(ctx).
		Unscoped().
		Where("game_id = ?", gameId).
		Delete(&FolderMetadata{}).
		Error
}

func (s *StoreFolderMetadataGorm) Edit(ctx context.Context, gameId int, metadata *FolderMetadata) error {
	//TODO implement me
	panic("implement me")
}
