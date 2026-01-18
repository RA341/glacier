package manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/downloader/clients/transmission"
	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/syncmap"
	"github.com/rs/zerolog/log"
)

type Config interface {
	LoadDownloader(cli types.ClientType) (types.ClientConfig, error)
	SetDownloader(cli types.ClientType, conf types.ClientConfig) error
	LoadAllClients() map[string]types.ClientConfig
}

type Service struct {
	cf Config
	// indicates any available downloader
	clients syncmap.Map[types.ClientType, types.Downloader]
}

func New(cf Config) *Service {
	s := &Service{
		cf: cf,
	}
	s.LoadEnabledProviders()
	return s
}

func (s *Service) Get(cli types.ClientType) (types.Downloader, error) {
	value, ok := s.clients.Load(cli)
	if !ok {
		return nil, fmt.Errorf("client not found for %v", cli.String())
	}
	return value, nil
}

func (s *Service) LoadEnabledProviders() {
	providers := s.cf.LoadAllClients()
	for key, config := range providers {
		val, ok := config["enable"]
		if !ok {
			continue
		}

		enabled, ok := val.(bool)
		if !ok || !enabled {
			continue
		}

		clientType, err := types.ClientTypeString(key)
		if err != nil {
			log.Warn().Err(err).
				Str("got", key).
				Strs("allowed", types.ClientTypeStrings()).
				Msg("invalid client type")
			continue
		}

		err = s.LoadClientFromConfig(clientType, config)
		if err != nil {
			log.Warn().Err(err).Str("name", key).Msg("failed to load provider")
			continue
		}
	}
}

func (s *Service) LoadClient(cli types.ClientType) error {
	conf, err := s.cf.LoadDownloader(cli)
	if err != nil {
		return fmt.Errorf("err loading config: %v", err)
	}
	return s.LoadClientFromConfig(cli, conf)
}

func (s *Service) LoadClientFromConfig(cli types.ClientType, conf types.ClientConfig) error {
	var initFn func(config types.ClientConfig) (types.Downloader, error)

	switch cli {
	case types.ClientTransmission:
		initFn = transmission.New
	case types.ClientUnknown:
	default:
		return fmt.Errorf("unknown client: %s", cli.String())
	}

	downloader, err := initFn(conf)
	if err != nil {
		return err
	}

	s.clients.Store(cli, downloader)

	log.Debug().Str("name", cli.String()).Msg("download client loaded")
	return nil
}

func (s *Service) GetActive() []types.ClientType {
	var names []types.ClientType

	s.clients.Range(func(key types.ClientType, value types.Downloader) bool {
		names = append(names, key)
		return true
	})

	return names
}
