package app

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ra341/glacier/internal/auth"
	"github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/database"
	"github.com/ra341/glacier/internal/downloader"
	"github.com/ra341/glacier/internal/indexer"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/metadata"
	"github.com/ra341/glacier/internal/search"
	"github.com/ra341/glacier/internal/services_manager"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/pkg/logger"

	"github.com/rs/zerolog/log"
)

func InitMeta(flavour info.FlavourType) {
	info.SetFlavour(flavour)
	info.PrintInfo()
	logger.InitDefault()
}

type App struct {
	Conf *config.Service

	Library       *library.Service
	DownloadSrv   *downloader.Service
	Search        *search.Service
	ConfigManager *services_manager.Service
	Indexer       *indexer.Service

	User    *user.Service
	Session *auth.Service
}

func NewApp() *App {
	conf := config.New()
	c := conf.Get()
	if c == nil {
		log.Fatal().Msg("config is nil THIS SHOULD NEVER HAPPEN")
		return nil
	}
	logger.InitConsole(c.Logger.Level, c.Logger.Verbose)

	db := database.New(c.Glacier.ConfigDir, c.Logger.DBLogger)

	libDb := library.NewStoreGorm(db)

	confDb := services_manager.NewStore(db)
	configManager := services_manager.New(confDb)

	manStore := library.NewStoreManifestGorm(db)
	fms := library.NewManifestService(libDb, manStore)

	downSrv := downloader.New(
		configManager.Downloader.LoadService,
		func(id int) {
			fms.GenerateManifest(context.Background(), id)
		},
		libDb,
		func() *downloader.Config {
			return &c.Download
		},
	)
	downSrv.StartTracker() // check for previous incomplete downloads

	libSrv := library.New(libDb, fms,
		downSrv,
		func() *library.Config {
			return &c.Library
		},
	)

	metaSrv := metadata.New(configManager.Meta.LoadService)
	indexerSrv := indexer.New(configManager.Indexer.LoadService)

	searchSrv := search.New(metaSrv, indexerSrv)

	userDb := user.NewStoreGorm(db)
	userSrv := user.NewService(userDb)

	sessionDb := auth.NewStoreGorm(db, c.Auth.MaxConcurrentSessions)
	sessionSrv := auth.New(sessionDb, userSrv, c.Auth.OpenRegistration)

	a := &App{
		Conf:          conf,
		Library:       libSrv,
		DownloadSrv:   downSrv,
		Search:        searchSrv,
		Indexer:       indexerSrv,
		ConfigManager: configManager,
		User:          userSrv,
		Session:       sessionSrv,
	}

	err := a.VerifyServices()
	if err != nil {
		log.Fatal().Err(err).Msg("could not verify services")
	}
	return a
}

func (a *App) VerifyServices() error {
	val := reflect.ValueOf(a).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		// We only care about pointers (services)
		if field.Kind() == reflect.Ptr && field.IsNil() {
			return fmt.Errorf("critical error: service '%s' was not initialized", fieldName)
		}
	}
	return nil
}
