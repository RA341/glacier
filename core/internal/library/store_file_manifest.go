package library

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type StoreFolderMetadata interface {
	Add(ctx context.Context, gameId int, metadata *FolderManifest) error
	Get(ctx context.Context, gameId int) (FolderManifest, error)
	Delete(ctx context.Context, gameId int) error
	Edit(ctx context.Context, gameId int, metadata *FolderManifest) error
}

type FolderManifest struct {
	gorm.Model

	GameID int  `gorm:"index"`
	Game   Game `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	TotalSize         int64
	AvailableExePaths []string       `gorm:"serializer:json" `
	FileInfo          []FileManifest `gorm:"serializer:json"`
}

type FileManifest struct {
	RelPath  string
	Size     int64
	ModTime  time.Time
	Checksum string
}
