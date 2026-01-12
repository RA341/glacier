package library

import (
	"context"

	"gorm.io/gorm"
)

//go:generate go run github.com/dmarkham/enumer@latest -type=GameType -output=enum_gametype.go
type GameType int

const (
	// GameTypeUnknown is the default zero-value
	GameTypeUnknown GameType = iota

	// GameTypeInstaller means the file is an installer and must be installed after downloaded
	GameTypeInstaller

	// GameTypeStandalone means the files are ready-to-run after download
	GameTypeStandalone
)

type Game struct {
	gorm.Model

	Name           string
	GameType       GameType
	DownloadedPath string
}

type Store interface {
	// GetByDownload SELECT * FROM @@table WHERE downloaded_path=@path
	GetByDownload(ctx context.Context, path string) (Game, error)
}
