package manifest

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type StoreGameManifest interface {
	Add(ctx context.Context, gameId int, metadata *FolderManifest) error
	Get(ctx context.Context, gameId int) (FolderManifest, error)
	Delete(ctx context.Context, gameId int) error
	Edit(ctx context.Context, gameId int, metadata *FolderManifest) error
	ListGamesWithoutManifest(ctx context.Context) ([]Info, error)
}

// Info information required to generate metadata
type Info struct {
	ID           int
	DownloadPath string
}

type FolderManifest struct {
	gorm.Model
	GameID    int `gorm:"column:game_id;uniqueIndex"`
	TotalSize int64
	FileInfo  []FileManifest `gorm:"serializer:json"`
}

type FileManifest struct {
	RelPath  string
	Size     int64
	ModTime  time.Time
	Checksum string
}
