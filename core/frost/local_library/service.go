package download

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	hc "github.com/ra341/glacier/frost/http_client"
	"github.com/ra341/glacier/frost/local_library/download"
	librpc "github.com/ra341/glacier/generated/library/v1"
	glacier "github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/internal/library"
)

type Service struct {
	store      Store
	baseurl    string
	downloader *download.Service
	lib        glacier.LibraryServiceClient
}

func New(baseurl string, store Store, downloader *download.Service, cli hc.HttpCliFactory) *Service {
	s := &Service{
		lib:        glacier.NewLibraryServiceClient(cli(&http.Transport{}), baseurl),
		store:      store,
		baseurl:    baseurl,
		downloader: downloader,
	}
	return s
}

func (s *Service) Download(ctx context.Context, gameId int, downloadFolder string) error {
	var ll LocalGame

	request := connect.NewRequest(&librpc.GetGameRequest{GameId: uint64(gameId)})
	game, err := s.lib.GetGame(ctx, request)
	if err != nil {
		return fmt.Errorf("could not get game info from server: %w", err)
	}

	var libGame library.Game
	libGame.FromProto(game.Msg.Game)

	ll.GameId = gameId
	ll.Game = libGame
	ll.Download.Started = time.Now()

	err = s.store.Add(ctx, &ll)
	if err != nil {
		return fmt.Errorf("could not add game to DB: %w", err)
	}

	return s.downloader.Download(gameId, downloadFolder)
}

func (s *Service) ListDownloading(ctx context.Context) ([]LocalGame, error) {
	return s.store.ListWithState(
		ctx,
		download.StatusDownloading,
		download.StatusMetadata,
		download.StatusQueued,
	)
}
