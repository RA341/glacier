package download

import (
	"context"

	"github.com/ra341/glacier/frost/local_library/download"
)

type Service struct {
	store      Store
	baseurl    string
	downloader *download.Service
}

func New(baseurl string, store Store) *Service {
	return &Service{
		store:      store,
		baseurl:    baseurl,
		downloader: download.New(baseurl, store, 50, 100),
	}
}

func (s *Service) Download(gameId int, downloadFolder string) error {
	return s.downloader.Download(gameId, downloadFolder)
}

func (s *Service) ListDownloading(ctx context.Context) (map[int][]download.FileProgress, error) {
	var games = map[int][]download.FileProgress{}

	var err error

	s.downloader.ActiveDownloads.Range(func(key int, value *download.Download) bool {
		var progress []download.FileProgress

		progress, err = value.Progress()
		if err != nil {
			return false
		}

		games[key] = progress
		return true
	})

	if err != nil {
		return nil, err
	}

	return games, nil
}
