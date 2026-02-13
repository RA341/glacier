package download

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	hc "github.com/ra341/glacier/frost/http_client"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/rs/zerolog/log"
)

const MB = 1024 * 1024

type Service struct {
	Config             *Config
	editStatus         EditStatus
	ActiveDownloads    syncmap.Map[int, *Download]
	DownloadTotalBytes *uint64
}

// New
//
// basepath must be: "http://localhost:6699"
func New(
	baseurl string,
	config *Config,
	httpCliFactory hc.HttpCliFactory,
	editStatus EditStatus,
	downloadSpeedCounter *uint64,
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

	config.httpCli = httpCliFactory(transport)
	config.base = fmt.Sprintf(
		"%s/library/download",
		strings.TrimRight(baseurl, "/"),
	)

	return &Service{
		Config:             config,
		editStatus:         editStatus,
		DownloadTotalBytes: downloadSpeedCounter,
	}
}

func (d *Service) Download(gameId int, downloadPath string, force bool) error {
	if _, ok := d.ActiveDownloads.Load(gameId); ok {
		log.Debug().Msg("download already in progress")
		return nil
	}

	download, err := NewDownload(
		d.Config,
		d.editStatus,
		d.onDone,
		gameId,
		downloadPath,
		force,
	)
	if err != nil {
		return fmt.Errorf("could not start download: %w", err)
	}

	d.ActiveDownloads.Store(gameId, download)

	return nil
}

func (d *Service) Pause(gameId int) error {
	return fmt.Errorf("not implemented")
}

func (d *Service) Cancel(gameId int) error {
	val, ok := d.ActiveDownloads.LoadAndDelete(gameId)
	if !ok {
		return fmt.Errorf("game not found " + strconv.Itoa(gameId))
	}

	val.Cancel()
	<-val.done // wait for download to stop

	log.Info().Msg("download stopped")

	err := os.RemoveAll(val.downloadFolder)
	if err != nil {
		return fmt.Errorf("could not remove download folder: %w", err)
	}

	return nil
}

func (d *Service) onDone(gameId int) {
	d.ActiveDownloads.Delete(gameId)
}
