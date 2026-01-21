package download

import (
	"gorm.io/gorm"
)

type InstalledGame struct {
	gorm.Model
	libraryId    uint
	ExePath      string
	DownloadPath string
}
