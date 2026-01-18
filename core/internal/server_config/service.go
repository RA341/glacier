package server_config

import (
	"sync/atomic"

	downloaderTypes "github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/internal/indexer/types"
	metadataTypes "github.com/ra341/glacier/internal/metadata/types"
	"github.com/rs/zerolog/log"
)

// todo likely this is not concurrent safe when reading

type Service struct {
	cy   ConfigYml
	conf atomic.Pointer[Config]
}

func New() *Service {
	s := &Service{}
	s.Init()
	return s
}

const GlacierYml = "glacier.yml"
const GlacierYmlPathEnv = "GLACIER_CONFIG_YML_PATH"

func (s *Service) Init() {
	s.cy = NewYml()

	var conf Config
	err := s.cy.loadYml(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config file")
	}

	defaultPrefixer := DefaultPrefixer()
	rnFn := FieldProcessorTag(defaultPrefixer)
	SetDefaultsFromTags(&conf, rnFn)

	pathsToResolve := []*string{
		&conf.Download.IncompleteDownloadPath,
		&conf.Files.ConfigDir,
		&conf.Library.GameDir,
	}
	resolvePaths(pathsToResolve)

	err = s.storeAndLoad(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't init config file")
	}
}

func (s *Service) Get() *Config {
	return s.conf.Load()
}

func (s *Service) storeAndLoad(loadCopy *Config) error {
	err := s.cy.writeAndLoad(loadCopy)
	if err != nil {
		return err
	}
	s.conf.Store(loadCopy)
	return nil
}

//func (s *Service) Set(src *Config) error {
//	newDst := s.loadCopy()
//
//	if err := mergo.Merge(&newDst, src); err != nil {
//		return err
//	}
//
//	s.conf.Store(&newDst)
//
//	return nil
//}

func (s *Service) loadCopy() Config {
	return *s.conf.Load()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// downloader

func (s *Service) LoadDownloader(cli downloaderTypes.ClientType) (downloaderTypes.ClientConfig, error) {
	return s.conf.Load().Download.GetClient(cli)
}

func (s *Service) SetDownloader(cli downloaderTypes.ClientType, conf downloaderTypes.ClientConfig) error {
	loadCopy := s.loadCopy()
	loadCopy.Download.SetClient(cli, conf)
	return s.storeAndLoad(&loadCopy)
}

func (s *Service) LoadAllClients() map[string]downloaderTypes.ClientConfig {
	return s.conf.Load().Download.Clients
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// providers

func (s *Service) LoadMetaProvider(cli metadataTypes.ProviderType) (metadataTypes.ProviderConfig, error) {
	return s.conf.Load().Metadata.GetCli(cli)
}

func (s *Service) SetMetaProvider(cli metadataTypes.ProviderType, conf metadataTypes.ProviderConfig) error {
	loadCopy := s.loadCopy()
	loadCopy.Metadata.SetCli(cli, conf)

	return s.storeAndLoad(&loadCopy)
}

func (s *Service) LoadAllProviders() map[string]metadataTypes.ProviderConfig {
	return s.conf.Load().Metadata.Providers
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// indexer

func (s *Service) LoadIndexer(cli types.IndexerType) (types.IndexerConfig, error) {
	return s.conf.Load().Indexer.GetCli(cli)
}

func (s *Service) SetIndexer(cli types.IndexerType, conf types.IndexerConfig) error {
	loadCopy := s.loadCopy()
	loadCopy.Indexer.SetCli(cli, conf)

	return s.storeAndLoad(&loadCopy)
}

func (s *Service) LoadAllIndexers() map[string]types.IndexerConfig {
	return s.conf.Load().Indexer.Indexers
}
