package app

import (
	"os"
	"path/filepath"

	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/frost/database"
	ll "github.com/ra341/glacier/frost/local_library"
	"github.com/ra341/glacier/pkg/logger"
	"github.com/rs/zerolog/log"
)

type App struct {
	Conf            *config.Service
	LocalLibrarySrv *ll.Service
}

func New() *App {
	conf := config.New()
	get := conf.Get()

	logger.InitConsole(get.Logger.Level, get.Logger.Verbose)

	abs, err := filepath.Abs(get.Files.ConfigDir)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config path")
	}
	err = os.MkdirAll(abs, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("could create config path")
	}

	db := database.New(get.Files.ConfigDir, false)
	llStore := ll.NewStoreGorm(db)
	llibSrv := ll.New(get.Server.GlacierUrl, llStore)

	return &App{
		Conf:            conf,
		LocalLibrarySrv: llibSrv,
	}
}
