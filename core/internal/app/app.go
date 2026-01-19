package app

import (
	"fmt"
	"reflect"

	"github.com/ra341/glacier/internal/database"
	"github.com/ra341/glacier/internal/downloader"
	downloadManager "github.com/ra341/glacier/internal/downloader/manager"
	"github.com/ra341/glacier/internal/indexer"
	indexerManager "github.com/ra341/glacier/internal/indexer/manager"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/metadata"
	metadataManager "github.com/ra341/glacier/internal/metadata/manager"
	"github.com/ra341/glacier/internal/search"
	sc "github.com/ra341/glacier/internal/server_config"
	"github.com/ra341/glacier/pkg/logger"

	"github.com/rs/zerolog/log"
)

func init() {
	logger.InitConsole("debug", true)
}

func InitMeta(flavour info.FlavourType) {
	info.SetFlavour(flavour)
	info.PrintInfo()
	logger.InitDefault()
}

type App struct {
	Conf *sc.Service

	Library *library.Service

	DownloadSrv           *downloader.Service
	DownloadClientManager *downloadManager.Service
	IndexerManager        *indexerManager.Service
	Search                *search.Service
	MetadataManager       *metadataManager.Service
}

func NewApp() *App {
	config := sc.New()
	get := config.Get()
	if get == nil {
		log.Fatal().Msg("config is nil THIS SHOULD NEVER HAPPEN")
	}

	db := database.New(config.Get().Glacier.ConfigDir, false)

	libDb := library.NewStoreGorm(db)

	clientMan := downloadManager.New(config)
	downSrv := downloader.New(
		clientMan.Get,
		libDb,
		func() *downloader.Config {
			return &config.Get().Download
		},
	)
	downSrv.StartTracker() // check for previous incomplete downloads

	libSrv := library.New(libDb, downSrv,
		func() *library.Config {
			return &config.Get().Library
		},
	)

	metaMan := metadataManager.New(config)
	metaSrv := metadata.New(metaMan)

	indexerMan := indexerManager.New(config)
	indexerSrv := indexer.New(indexerMan)

	searchSrv := search.New(metaSrv, indexerSrv)

	a := &App{
		Conf:                  config,
		Library:               libSrv,
		IndexerManager:        indexerMan,
		DownloadSrv:           downSrv,
		DownloadClientManager: clientMan,
		MetadataManager:       metaMan,
		Search:                searchSrv,
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
