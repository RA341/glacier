package config_manager

import (
	"fmt"

	downloaderTypes "github.com/ra341/glacier/internal/downloader/types"
	indexTypes "github.com/ra341/glacier/internal/indexer/types"
	metadata "github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/mapsct"
)

type ServiceHandlers struct {
	init         func(*ServiceConfig) error
	getSupported func() []string
	getSchema    func(string) ([]mapsct.FieldSchema, error)
}

type Service struct {
	Downloader ServiceConfigMap[downloaderTypes.Downloader]
	Indexer    ServiceConfigMap[indexTypes.Indexer]
	Meta       ServiceConfigMap[metadata.Provider]

	store Store

	// holds the map active clients

	// The registry replaces the switch statements
	registry map[ServiceType]ServiceHandlers
}

func New(store Store) *Service {
	s := &Service{
		Downloader: NewDownloaderMap(store),
		Indexer:    NewIndexerMap(store),
		Meta:       NewMetadataMap(store),
		store:      store,
	}

	s.registry = map[ServiceType]ServiceHandlers{
		Metadata: {
			init: func(cfg *ServiceConfig) error {
				_, err := s.Meta.initService(cfg)
				return err
			},
			getSupported: metadata.ProviderTypeStrings,
			getSchema:    s.Meta.getServiceSchema,
		},
		Indexer: {
			init: func(cfg *ServiceConfig) error {
				_, err := s.Indexer.initService(cfg)
				return err
			},
			getSupported: indexTypes.IndexerTypeStrings,
			getSchema:    s.Indexer.getServiceSchema,
		},
		Downloader: {
			init: func(cfg *ServiceConfig) error {
				_, err := s.Downloader.initService(cfg)
				return err
			},
			getSupported: downloaderTypes.ClientTypeStrings,
			getSchema:    s.Downloader.getServiceSchema,
		},
	}

	return s
}

func (s *Service) TestAndSave(cf *ServiceConfig) error {
	err := s.Test(cf)
	if err != nil {
		return err
	}

	return s.store.New(cf)
}

func (s *Service) Test(cfg *ServiceConfig) error {
	if handlers, ok := s.registry[cfg.ServiceType]; ok {
		return handlers.init(cfg)
	}
	return fmt.Errorf("unsupported service type: %v", cfg.ServiceType)
}

func (s *Service) GetSupportedValues(serviceType ServiceType) ([]string, error) {
	if handlers, ok := s.registry[serviceType]; ok {
		return handlers.getSupported(), nil
	}
	return nil, fmt.Errorf("unsupported service type: %v", serviceType)
}

func (s *Service) GetSchema(serviceType ServiceType, flv string) ([]mapsct.FieldSchema, error) {
	if handlers, ok := s.registry[serviceType]; ok {
		return handlers.getSchema(flv)
	}
	return nil, fmt.Errorf("unknown service type: %v", serviceType)
}
