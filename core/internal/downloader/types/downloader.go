package types

import "context"

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=DownloadState -output=enum_download_state.go
type DownloadState int

const (
	// Unknown this should only happen in serialization when an invalid external value is received
	Unknown DownloadState = iota
	Queued
	Downloading
	Complete
	Error
)

// Download todo add download client type
type Download struct {
	Client ClientType
	// DownloadId contains the ID of the download from download client
	DownloadId string
	State      DownloadState
	// Progress progressString contains any message from the client of the download
	Progress string

	DownloadUrl    string
	DownloadPath   string
	IncompletePath string
}

// Downloader generic interface that any downloader must implement
type Downloader interface {
	Download(ctx context.Context, url string, downloadPath string) (downloadID string, err error)
	Cancel(ctx context.Context, downloadId string, removeDownloaded bool) error
	Progress(ctx context.Context, download *Download) (err error)
}
