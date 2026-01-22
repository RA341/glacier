package download

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/pkg/syncmap"
)

func NewDownloadClient(maxWorkers, maxConn int) *http.Client {
	transport := &http.Transport{
		// MaxIdleConns is the total connections across all hosts
		MaxIdleConns: maxConn,

		// MaxIdleConnsPerHost must be >= your worker count.
		// The default is only 2 If you have 30 workers, 28 will
		// constantly create new TCP connections.
		MaxIdleConnsPerHost: maxWorkers,

		// Time to keep an idle connection in the pool
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   0,
	}
}

type Downloader struct {
	BaseUrl    string
	httpClient *http.Client

	maxConcurrentFiles      int
	maxConcurrentFileChunks int
	chunkSize               int64

	ActiveDownloads syncmap.Map[int, *Download]
}

func (d *Downloader) getMaxConcurrentFiles() int {
	return d.maxConcurrentFiles
}

func (d *Downloader) getChunkSize() int64 {
	return d.chunkSize
}

func (d *Downloader) getMaxConcurrentFileChunks() int {
	return d.maxConcurrentFileChunks
}

func (d *Downloader) getHttpClient() *http.Client {
	return d.httpClient
}

const MB = 1024 * 1024

func NewDownloader(maxConcurrentFiles, maxConcurrentFileChunks int) *Downloader {
	return &Downloader{
		httpClient: NewDownloadClient(100, 0),
		BaseUrl:    "http://localhost:6699/api/server/library/download",

		maxConcurrentFiles:      maxConcurrentFiles,
		maxConcurrentFileChunks: maxConcurrentFileChunks,
		chunkSize:               MB * 250,
	}
}

func (d *Downloader) ListActiveDownload() {

}

func (d *Downloader) Download(downloadFolder string, gameId int, met *library.FolderMetadata) error {
	download, err := NewDownload(d, d.BaseUrl, gameId, downloadFolder, met)
	if err != nil {
		return fmt.Errorf("could not start download: %w", err)
	}

	d.ActiveDownloads.Store(gameId, download)

	return nil
}
