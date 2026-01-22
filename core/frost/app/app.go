package app

import (
	"os"
	"path/filepath"

	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/frost/database"
	ll "github.com/ra341/glacier/frost/local_library"
	"github.com/rs/zerolog/log"
)

type App struct {
	Conf            *config.Service
	LocalLibrarySrv *ll.Service
}

func New() *App {
	conf := config.New()

	base := "http://localhost:6699"
	configBase := "./config"

	abs, err := filepath.Abs(configBase)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get config path")
	}
	err = os.MkdirAll(abs, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("could create config path")
	}

	db := database.New(configBase, false)
	llStore := ll.NewStoreGorm(db)
	llibSrv := ll.New(base, llStore)

	return &App{
		Conf:            conf,
		LocalLibrarySrv: llibSrv,
	}
}
