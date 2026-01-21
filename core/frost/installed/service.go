package download

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ra341/glacier/internal/library"
	"github.com/rs/zerolog/log"
	"resty.dev/v3"
)

type Service struct {
	baseurl    string
	downloader *Downloader
}

func New(baseurl string) *Service {
	return &Service{
		baseurl:    baseurl,
		downloader: NewDownloader(50, 100),
	}
}

func (s *Service) Download(gameId int, downloadFolder string) error {
	var meta library.FolderMetadata

	now := time.Now()

	err := s.downloadMetadata(gameId, &meta)
	if err != nil {
		return err
	}

	elapsed := time.Since(now)
	log.Info().Dur("elapsed", elapsed).Msg("got metadata, starting download...")

	gameDownload := filepath.Join(downloadFolder, strconv.Itoa(gameId))
	err = os.MkdirAll(gameDownload, 0755)
	if err != nil {
		return err
	}

	err = s.downloader.Download(gameDownload, gameId, &meta)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) downloadMetadata(gameId int, meta *library.FolderMetadata) error {
	get, err := resty.New().
		//SetDebug(true).
		SetBaseURL(s.baseurl).R().
		Get("/api/server/library/download/meta/" + strconv.Itoa(gameId))
	if err != nil {
		return err
	}

	if get.IsError() {
		return fmt.Errorf("%s", get.String())
	}

	decoder := gob.NewDecoder(get.Body)
	return decoder.Decode(meta)
}
