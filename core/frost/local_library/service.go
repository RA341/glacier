package download

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"connectrpc.com/connect"
	hc "github.com/ra341/glacier/frost/http_client"
	"github.com/ra341/glacier/frost/launcher"
	"github.com/ra341/glacier/frost/local_library/download"
	librpc "github.com/ra341/glacier/generated/library/v1"
	glacier "github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/internal/library"
	"github.com/rs/zerolog/log"
)

type Service struct {
	baseurl string
	store   Store
	lib     glacier.LibraryServiceClient

	downloader *download.Service
	launcher   *launcher.Service
}

func New(baseurl string, store Store, downloader *download.Service, cli hc.HttpCliFactory) *Service {
	s := &Service{
		lib:        glacier.NewLibraryServiceClient(cli(&http.Transport{}), baseurl),
		store:      store,
		baseurl:    baseurl,
		downloader: downloader,
		launcher:   launcher.New(),
	}
	s.loadDownloading()

	return s
}

func (s *Service) Running(ctx context.Context, gameID int, exe string) error {
	get, err := s.store.Get(ctx, gameID)
	if err != nil {
		return err
	}

	exe = get.Game.File.Exe
	if exe == "" {
		return fmt.Errorf("game exe not found please")
	}
	fullPath := filepath.Join(get.Download.DownloadPath, exe)

	return s.launcher.Running(ctx, fullPath)
}

func (s *Service) Launch(ctx context.Context, gameID int) error {
	get, err := s.store.Get(ctx, gameID)
	if err != nil {
		return err
	}

	exe := get.Game.File.Exe
	if exe == "" {
		return fmt.Errorf("game exe not found please")
	}
	fullPath := filepath.Join(get.Download.DownloadPath, exe)

	return s.launcher.Launch(fullPath)
}

func (s *Service) Download(
	ctx context.Context,
	gameId int,
	downloadFolder string,
	recheck bool,
	force bool,
) error {
	if recheck {
		game, err := s.store.GetByGameId(ctx, uint64(gameId), false)
		if err != nil {
			return fmt.Errorf("could not find game, did you add it first ?: %w", err)
		}

		return s.downloader.Download(
			ctx,
			gameId,
			game.Download.DownloadPath,
			force,
		)
	}

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

	ll.Download.DownloadPath, err = filepath.Abs(filepath.Join(downloadFolder, libGame.Meta.Name))
	if err != nil {
		return fmt.Errorf("could not get download path: %w", err)
	}

	err = os.MkdirAll(ll.Download.DownloadPath, 0755)
	if err != nil {
		return err
	}

	// todo check for avail space
	// github.com/shirou/gopsutil/v4/disk
	// usage, err := disk.Usage("/")
	// if err != nil {
	//	fmt.Printf("Error: %v\n", err)
	//	return
	//}

	err = s.store.Add(ctx, &ll)
	if err != nil {
		return fmt.Errorf("could not add game to DB: %w", err)
	}

	return s.downloader.Download(
		ctx,
		gameId,
		ll.Download.DownloadPath,
		force,
	)
}

func (s *Service) ListDownloading(ctx context.Context) ([]LocalGame, error) {
	return s.store.ListWithState(
		ctx,
		download.StatusDownloading,
		download.StatusMetadata,
		download.StatusQueued,
	)
}

func (s *Service) Cancel(id int) error {
	var errs []error

	err := s.downloader.Cancel(id)
	if err != nil {
		errs = append(errs, err)
	}

	log.Debug().Msg("removing download from db")
	err = s.store.Delete(context.Background(), id)
	if err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (s *Service) loadDownloading() {
	ctx := context.Background()
	state, err := s.ListDownloading(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("could not list downloading games")
		return
	}

	for _, g := range state {
		err := s.Download(
			ctx,
			int(g.ID),
			g.Download.DownloadPath,
			true,
			false,
		)
		if err != nil {
			log.Warn().Err(err).Msg("could not restart game download")
		}
	}
}

func (s *Service) Delete(ctx context.Context, id uint64) error {
	localGame, err := s.store.Get(ctx, int(id))
	if err != nil {
		return err
	}

	err = os.RemoveAll(localGame.Download.DownloadPath)
	if err != nil {
		return err
	}

	return s.store.Delete(ctx, int(id))
}
