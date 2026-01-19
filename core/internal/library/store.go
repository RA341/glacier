package library

import (
	"context"

	download "github.com/ra341/glacier/internal/downloader/types"
	indexer "github.com/ra341/glacier/internal/indexer/types"
	metadata "github.com/ra341/glacier/internal/metadata/types"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	// metadata of the Game
	Meta metadata.Meta `gorm:"embedded"`
	// client where the source is getting downloaded
	Download download.Download `gorm:"embedded"`
	// Source of the Indexer
	Source indexer.Source `gorm:"embedded"`
}

type Store interface {
	Add(ctx context.Context, game *Game) error
	UpdateDownloadProgress(ctx context.Context, id uint, download download.Download) error
	GetById(ctx context.Context, id uint) (Game, error)
	DeleteGame(ctx context.Context, id uint) error
	List(ctx context.Context, limit uint, offset uint) ([]Game, error)
	ListDownloadState(ctx context.Context, state download.DownloadState) ([]Game, error)
}
