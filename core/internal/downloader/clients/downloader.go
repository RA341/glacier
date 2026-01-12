package clients

import "context"

//go:generate go run github.com/dmarkham/enumer@latest -type=DownloadState -output=enum_downloadstate.go
type DownloadState int

const (
	Queued DownloadState = iota
	Downloading
	Complete
	Error
)

// Downloader generic interface that any downloader must implement
type Downloader interface {
	Download(ctx context.Context, url string, downloadPath string) (downloadID string, err error)
	Cancel(ctx context.Context, downloadId string, removeDownloaded string) error

	// Progress gets the progress of a download
	//
	// progressString contains the state of the download
	Progress(ctx context.Context, downloadId string) (state DownloadState, progressString string, err error)
}
