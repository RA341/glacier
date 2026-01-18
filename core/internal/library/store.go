package library

import (
	"context"

	download "github.com/ra341/glacier/internal/downloader/types"
	metadata "github.com/ra341/glacier/internal/metadata/types"

	"gorm.io/gorm"
)

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=GameType -output=enum_gametype.go

// GameType identifies the type of files downloaded
type GameType int

const (
	// GameTypeUnknown is the default zero-value
	GameTypeUnknown GameType = iota

	// GameTypeInstaller means the files must be installed after download
	GameTypeInstaller

	// GameTypeStandalone means the files are ready-to-play after download
	GameTypeStandalone
)

type Game struct {
	gorm.Model
	GameType GameType
	Meta     metadata.Meta     `gorm:"embedded"`
	Download download.Download `gorm:"embedded"`
}

type Store interface {
	Add(ctx context.Context, game *Game) error
	UpdateDownloadProgress(ctx context.Context, id uint, download download.Download) error
	GetById(ctx context.Context, id uint) (Game, error)
	DeleteGame(ctx context.Context, id uint) error
	List(ctx context.Context, limit uint, offset uint) ([]Game, error)
	ListDownloadState(ctx context.Context, state download.DownloadState) ([]Game, error)
}
