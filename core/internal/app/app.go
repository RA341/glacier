package app

import (
	"fmt"
	"reflect"

	"github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/config/config_manager"
	"github.com/ra341/glacier/internal/database"
	"github.com/ra341/glacier/internal/downloader"
	"github.com/ra341/glacier/internal/indexer"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/metadata"
	"github.com/ra341/glacier/internal/search"
	"github.com/ra341/glacier/pkg/logger"

	"github.com/rs/zerolog/log"
)

func InitMeta(flavour info.FlavourType) {
	info.SetFlavour(flavour)
	info.PrintInfo()
	logger.InitDefault()
}

type App struct {
	Conf          *config.Service
	Library       *library.Service
	DownloadSrv   *downloader.Service
	Search        *search.Service
	ConfigManager *config_manager.Service
}

func NewApp() *App {
	conf := config.New()
	get := conf.Get()
	if get == nil {
		log.Fatal().Msg("config is nil THIS SHOULD NEVER HAPPEN")
		return nil
	}
	logger.InitConsole(get.Logger.Level, get.Logger.Verbose)

	db := database.New(conf.Get().Glacier.ConfigDir, false)

	libDb := library.NewStoreGorm(db)

	confDb := config_manager.NewStore(db)
	configManager := config_manager.New(confDb)

	downSrv := downloader.New(
		configManager.Downloader.LoadService,
		libDb,
		func() *downloader.Config {
			return &conf.Get().Download
		},
	)
	downSrv.StartTracker() // check for previous incomplete downloads
	folderMetaDb := library.NewStoreFolderMetadataGorm(db)
	libSrv := library.New(libDb, folderMetaDb,
		downSrv,
		func() *library.Config {
			return &conf.Get().Library
		},
	)

	metaSrv := metadata.New(configManager.Meta.LoadService)
	indexerSrv := indexer.New(configManager.Indexer.LoadService)

	searchSrv := search.New(metaSrv, indexerSrv)

	a := &App{
		Conf:          conf,
		Library:       libSrv,
		DownloadSrv:   downSrv,
		Search:        searchSrv,
		ConfigManager: configManager,
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
