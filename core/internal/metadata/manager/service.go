package manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/metadata/providers/igdb"
	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/syncmap"

	"github.com/rs/zerolog/log"
)

type Config interface {
	LoadMetaProvider(cli types.ProviderType) (types.ProviderConfig, error)
	SetMetaProvider(cli types.ProviderType, conf types.ProviderConfig) error
	LoadAllProviders() map[string]types.ProviderConfig
}

type Service struct {
	cf        Config
	providers syncmap.Map[types.ProviderType, types.Provider]
}

func New(cf Config) *Service {
	s := &Service{
		cf: cf,
	}
	go s.LoadEnabledProviders()
	return s
}

func (s *Service) Get(cli types.ProviderType) (types.Provider, error) {
	value, ok := s.providers.Load(cli)
	if !ok {
		return nil, fmt.Errorf("client not found for %v", cli.String())
	}
	return value, nil
}

func (s *Service) LoadEnabledProviders() {
	providers := s.cf.LoadAllProviders()
	for key, config := range providers {
		val, ok := config["enable"]
		if !ok {
			continue
		}

		enabled, ok := val.(bool)
		if !ok || !enabled {
			continue
		}

		providerType, err := types.ProviderTypeString(key)
		if err != nil {
			log.Warn().Err(err).
				Str("got", key).
				Strs("allowed", types.ProviderTypeStrings()).
				Msg("invalid provider type")
			continue
		}

		err = s.LoadProviderFromConfig(providerType, config)
		if err != nil {
			log.Warn().Err(err).Str("name", key).Msg("failed to load provider config")
			continue
		}
	}
}

func (s *Service) LoadProvider(cli types.ProviderType) error {
	conf, err := s.cf.LoadMetaProvider(cli)
	if err != nil {
		return fmt.Errorf("err loading config: %v", err)
	}

	return s.LoadProviderFromConfig(cli, conf)
}

func (s *Service) LoadProviderFromConfig(cli types.ProviderType, conf types.ProviderConfig) error {
	var initFn func(config types.ProviderConfig) (types.Provider, error)

	switch cli {
	case types.ProviderIGDB:
		initFn = igdb.New
	case types.ProviderUnknown:
	default:
		return fmt.Errorf("unknown provider: %s", cli.String())
	}

	downloader, err := initFn(conf)
	if err != nil {
		return err
	}

	s.providers.Store(cli, downloader)

	log.Debug().Str("name", cli.String()).Msg("metadata provider loaded")
	return nil
}

func (s *Service) GetActive() []types.ProviderType {
	var names []types.ProviderType

	s.providers.Range(func(key types.ProviderType, value types.Provider) bool {
		names = append(names, key)
		return true
	})

	return names
}
