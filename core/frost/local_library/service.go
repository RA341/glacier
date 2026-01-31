package download

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/frost/local_library/download"
	librpc "github.com/ra341/glacier/generated/library/v1"
	"github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/internal/library"
)

type Service struct {
	store      Store
	baseurl    string
	downloader *download.Service
	lib        v1connect.LibraryServiceClient
}

func New(baseurl string, store Store) *Service {
	s := &Service{
		lib:        v1connect.NewLibraryServiceClient(http.DefaultClient, baseurl+"/api/server"),
		store:      store,
		baseurl:    baseurl,
		downloader: download.New(baseurl, store, 50, 100),
	}
	return s
}

func (s *Service) Download(ctx context.Context, gameId int, downloadFolder string) error {
	var ll LocalGame

	game, err := s.lib.GetGame(ctx, connect.NewRequest(&librpc.GetGameRequest{GameId: uint64(gameId)}))
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
