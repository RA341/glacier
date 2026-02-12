package app

import (
	"fmt"
	"reflect"

	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/frost/database"
	hc "github.com/ra341/glacier/frost/http_client"
	ll "github.com/ra341/glacier/frost/local_library"
	"github.com/ra341/glacier/frost/local_library/download"
	"github.com/ra341/glacier/frost/secrets"
	"github.com/ra341/glacier/pkg/logger"
	sharedConfig "github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog/log"
)

type App struct {
	Conf            *sharedConfig.Service[config.Config]
	LocalLibrarySrv *ll.Service
	Secret          *secrets.Service
}

func New() *App {
	conf := config.New(true)
	get := conf.Get()

	logger.InitConsole(get.Logger.Level, get.Logger.Verbose)

	db := database.New(get.Files.ConfigDir, false)

	const appName = "dev.radn.glacier.frost"
	secretStore := secrets.NewKeyringStore(appName)
	ss := secrets.NewService(secretStore)

	httpCliFac := hc.NewFrostHttpClientFactory(ss)

	frostProtectedBase := get.Server.GlacierUrl + "/api/server/protected"

	llStore := ll.NewStoreGorm(db)
	downloader := download.New(
		frostProtectedBase,
		&download.Config{
			MaxConcurrentFiles:      5,
			MaxConcurrentFileChunks: 50,
			ChunkSizeInMB:           128 * download.MB,
		}, // todo
		httpCliFac,
		llStore.EditStatus,
	)

	llibSrv := ll.New(
		frostProtectedBase,
		llStore,
		downloader,
		httpCliFac,
	)

	a := &App{
		Conf:            conf,
		LocalLibrarySrv: llibSrv,
		Secret:          ss,
	}
	err := a.VerifyServices()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load services")
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
