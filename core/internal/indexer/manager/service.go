package manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/indexer/indexers/hydra"
	"github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/rs/zerolog/log"
)

type Config interface {
	LoadIndexer(cli types.IndexerType) (types.IndexerConfig, error)
	SetIndexer(cli types.IndexerType, conf types.IndexerConfig) error
	LoadAllIndexers() map[string]types.IndexerConfig
}

type Service struct {
	cf            Config
	activeIndexer syncmap.Map[types.IndexerType, types.Indexer]
}

func New(df Config) *Service {
	s := &Service{
		cf: df,
	}
	s.LoadEnabledIndexers()
	return s
}

func (s *Service) Get(cli types.IndexerType) (types.Indexer, error) {
	value, ok := s.activeIndexer.Load(cli)
	if !ok {
		return nil, fmt.Errorf("client not found for %v", cli.String())
	}
	return value, nil
}

func (s *Service) LoadEnabledIndexers() {
	providers := s.cf.LoadAllIndexers()
	for key, config := range providers {
		val, ok := config["enable"]
		if !ok {
			continue
		}

		enabled, ok := val.(bool)
		if !ok || !enabled {
			continue
		}

		providerType, err := types.IndexerTypeString(key)
		if err != nil {
			log.Warn().Err(err).
				Str("got", key).
				Strs("allowed", types.IndexerTypeStrings()).
				Msg("invalid indexer type")
			continue
		}

		err = s.LoadFromConfig(providerType, config)
		if err != nil {
			log.Warn().Err(err).Str("name", key).Msg("failed to load indexer config")
			continue
		}
	}
}

func (s *Service) LoadProvider(cli types.IndexerType) error {
	conf, err := s.cf.LoadIndexer(cli)
	if err != nil {
		return fmt.Errorf("err loading config: %v", err)
	}

	return s.LoadFromConfig(cli, conf)
}

func (s *Service) LoadFromConfig(cli types.IndexerType, conf types.IndexerConfig) error {
	var initFn func(config types.IndexerConfig) (types.Indexer, error)

	switch cli {
	case types.IndexerHydra:
		initFn = hydra.New
	case types.IndexerUnknown:
	default:
		return fmt.Errorf("unknown indexer: %s", cli.String())
	}

	downloader, err := initFn(conf)
	if err != nil {
		return err
	}

	s.activeIndexer.Store(cli, downloader)

	log.Debug().Str("name", cli.String()).Msg("indexer loaded")
	return nil
}

func (s *Service) GetActive() []types.IndexerType {
	var names []types.IndexerType

	s.activeIndexer.Range(func(key types.IndexerType, value types.Indexer) bool {
		names = append(names, key)
		return true
	})

	return names
}
