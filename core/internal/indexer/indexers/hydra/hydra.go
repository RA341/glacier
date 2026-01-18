package hydra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/mapsct"

	"github.com/rs/zerolog/log"
	"resty.dev/v3"
)

// uses the hydra launcher sources

type Config struct {
	// stores the downloaded JSON files
	CacheDir string

	// list of the JSON hydra url
	// <source>:<url>
	Sources map[string]string
	// download the JSON again to update with any changes
	UpdateInterval string
	Debug          bool
}

func (c Config) GetInterval() time.Duration {
	duration, err := time.ParseDuration(c.UpdateInterval)
	if err != nil {
		defaultDur := time.Hour * 24
		return defaultDur
	}
	return duration
}

type Hydra struct {
	config Config

	cancel context.CancelFunc
}

func New(config types.IndexerConfig) (types.Indexer, error) {
	raw, err := newRaw(config)
	if err != nil {
		return nil, err
	}
	go raw.startIndexUpdater()

	return raw, nil
}

// useful for testing
func newRaw(config types.IndexerConfig) (*Hydra, error) {
	var conf Config

	err := mapsct.ParseMap(&conf, config)
	if err != nil {
		return nil, err
	}

	abs, err := filepath.Abs(conf.CacheDir)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(abs, 0755)
	if err != nil {
		return nil, err
	}
	conf.CacheDir = abs

	return &Hydra{config: conf}, nil
}

func (h *Hydra) Close() {
	if h.cancel != nil {
		h.cancel()
	}
}

////////////////////////////////////////////////////////////
// search stuff

func (h *Hydra) Search(query string) ([]types.IndexerGame, error) {
	index, err := h.loadIndex()
	if err != nil {
		return nil, err
	}

	var infos []types.IndexerGame
	for _, d := range index.Downloads {
		contains := strings.Contains(strings.ToLower(d.Title), strings.ToLower(query))
		if contains {
			infos = append(infos, types.IndexerGame{
				Title:       d.Title,
				DownloadUrl: d.Uris[0],
				FileSize:    d.FileSize,
				CreatedISO:  d.UploadDate.Format(time.RFC3339),
			})
		}
	}

	return infos, nil
}

type DownloadInfo struct {
	Title      string    `json:"title"`
	Uris       []string  `json:"uris"`
	UploadDate time.Time `json:"uploadDate"`
	FileSize   string    `json:"fileSize"`
}

type JsonResult struct {
	Downloads []DownloadInfo `json:"downloads"`
}

func (h *Hydra) loadIndex() (*JsonResult, error) {
	// todo search all files for now it just gets the first one in the map

	for name, _ := range h.config.Sources {
		path := h.getJsonPath(name)
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		val := &JsonResult{}
		err = json.Unmarshal(contents, val)
		if err != nil {
			return nil, err
		}

		return val, nil
	}

	return &JsonResult{}, nil
}

func (h *Hydra) withPath(name string) string {
	return filepath.Join(h.config.CacheDir, name)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// indexer background updates

// is blocking must be run in a go routine
func (h *Hydra) startIndexUpdater() {
	timer := time.NewTimer(h.config.GetInterval())
	defer timer.Stop()

	// run once at start
	h.updateIndexes()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h.cancel = cancel

	for {
		select {
		case <-timer.C:
			h.updateIndexes()
		case <-ctx.Done():
			log.Info().
				Str("indexer", types.IndexerHydra.String()).
				Msg("stopping index updater")
			return
		}
	}
}

func (h *Hydra) updateIndexes() {
	log.Info().
		Str("indexer", types.IndexerHydra.String()).
		Msg("updating indexes")

	for name, url := range h.config.Sources {
		err := h.downloadIndex(name, url)
		if err != nil {
			log.Warn().Err(err).Str("name", name).Msg("failed to update index")
		}
	}
}

func (h *Hydra) downloadIndex(name string, url string) error {
	etag, err := h.loadEtag(name)
	if err != nil {
		return err
	}

	get, err := resty.New().
		SetHeaders(map[string]string{
			"If-None-Match": etag,
		}).
		//SetDebug(h.config.Debug).
		R().
		Get(url)
	if err != nil {
		return err
	}

	if get.IsError() {
		return fmt.Errorf(
			"could not get indexer %s\n code:%d \n%s",
			name,
			get.StatusCode(),
			get.String(),
		)
	}

	newEtag := get.Header().Get("Etag")
	if get.StatusCode() == http.StatusNotModified {
		log.Debug().Msg("file unmodified")
		return nil
	}

	log.Debug().Str("prev-etag", etag).Str("current-etag", newEtag).Msg("file modified")

	if err = h.WriteJSON(name, get.Body); err != nil {
		return err
	}

	if err = h.writeEtag(name, newEtag); err != nil {
		return err
	}

	return nil
}

func (h *Hydra) WriteJSON(name string, writer io.Reader) error {
	path := h.getJsonPath(name)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil
	}
	defer fileutil.Close(file)

	_, err = io.Copy(file, writer)
	return err
}

func (h *Hydra) getJsonPath(name string) string {
	return h.withPath(fmt.Sprintf("%s.json", name))
}

func (h *Hydra) loadEtag(name string) (string, error) {
	path := h.getEtagPath(name)

	err := h.createFile(path)
	if err != nil {
		return "", err
	}

	readFile, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(readFile), nil
}

func (h *Hydra) writeEtag(name, etag string) error {
	path := h.getEtagPath(name)

	err := h.createFile(path)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(etag), 0644)
}

func (h *Hydra) createFile(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	defer fileutil.Close(f)
	return err
}

func (h *Hydra) getEtagPath(name string) string {
	return h.withPath(fmt.Sprintf("%s.etag", name))
}
