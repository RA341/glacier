package download

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	hc "github.com/ra341/glacier/frost/http_client"
	"github.com/ra341/glacier/pkg/syncmap"
)

const MB = 1024 * 1024

type Service struct {
	progress ProgressUpdater

	baseurl    string
	httpClient *http.Client

	maxConcurrentFiles      int
	maxConcurrentFileChunks int
	chunkSize               int64

	ActiveDownloads syncmap.Map[int, *Download]
}

// New
//
// basepath must be: "http://localhost:6699//"
func New(
	basepath string,
	httpCliFactory hc.HttpCliFactory,
	progress ProgressUpdater,
	maxConcurrentFiles, maxConcurrentFileChunks int,
) *Service {
	transport := &http.Transport{
		// MaxIdleConns is the total connections across all hosts
		MaxIdleConns: 100,

		// MaxIdleConnsPerHost must be >= your worker count.
		// The default is only 2 If you have 30 workers, 28 will
		// constantly create new TCP connections.
		MaxIdleConnsPerHost: 100,

		// Time to keep an idle connection in the pool
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &Service{
		httpClient:              httpCliFactory(transport),
		baseurl:                 fmt.Sprintf("%s/library/download", basepath),
		progress:                progress,
		maxConcurrentFiles:      maxConcurrentFiles,
		maxConcurrentFileChunks: maxConcurrentFileChunks,
		chunkSize:               MB * 128,
	}
}

func (d *Service) Download(gameId int, downloadFolder string) error {
	gameDownload := filepath.Join(downloadFolder, strconv.Itoa(gameId))
	err := os.MkdirAll(gameDownload, 0755)
	if err != nil {
		return err
	}

	// todo check for avail space
	download, err := NewDownload(
		d,
		d.onDone,
		d.progress,
		d.baseurl,
		gameDownload,
		gameId,
	)
	if err != nil {
		return fmt.Errorf("could not start download: %w", err)
	}

	d.ActiveDownloads.Store(gameId, download)

	return nil
}
func (d *Service) onDone(gameId int) {
	d.ActiveDownloads.Delete(gameId)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// satisfies Config in downloader.go

func (d *Service) getMaxConcurrentFiles() int {
	return d.maxConcurrentFiles
}

func (d *Service) getChunkSize() int64 {
	return d.chunkSize
}

func (d *Service) getMaxConcurrentFileChunks() int {
	return d.maxConcurrentFileChunks
}

func (d *Service) getHttpClient() *http.Client {
	return d.httpClient
}
