package downloader

import (
	"gorm.io/gorm"
)

type Downloads struct {
	gorm.Model

	DownloadedPath string
}

type Store interface {
	// GetByDownload SELECT * FROM @@table WHERE downloaded_path=@path
	//GetByDownload(ctx context.Context, path string) (Game, error)
}
