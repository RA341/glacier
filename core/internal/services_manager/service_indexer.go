package services_manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/indexer/indexers/hydra"
	indexTypes "github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/pkg/mapsct"
	"github.com/ra341/glacier/pkg/syncmap"
)

type IndexerMap struct {
	initMap map[indexTypes.IndexerType]ServiceConfigOpts[indexTypes.Indexer]
	store   Store

	indexerMap syncmap.Map[string, indexTypes.Indexer]
}

func NewIndexerMap(store Store) ServiceConfigMap[indexTypes.Indexer] {
	return &IndexerMap{
		store: store,
		initMap: map[indexTypes.IndexerType]ServiceConfigOpts[indexTypes.Indexer]{
			indexTypes.IndexerHydra: {
				InitFn: hydra.New,
				Config: hydra.Config{},
			},
		},
	}
}

func (s *IndexerMap) LoadService(id string) (indexTypes.Indexer, error) {
	val, ok := s.indexerMap.Load(id)
	if ok {
		return val, nil
	}

	conf, err := s.store.Get(id)
	if err != nil {
		return nil, err
	}

	service, err := s.initService(&conf)
	if err != nil {
		return nil, err
	}

	s.indexerMap.Store(id, service)

	return service, nil
}

func (s *IndexerMap) initService(conf *ServiceConfig) (indexTypes.Indexer, error) {
	val, err := s.loadServiceMap(conf.Flavour)
	if err != nil {
		return nil, err
	}
	return val.InitFn(conf.Config)
}

func (s *IndexerMap) getServiceSchema(flv string) ([]mapsct.FieldSchema, error) {
	val, err := s.loadServiceMap(flv)
	if err != nil {
		return nil, err
	}
	return mapsct.GetSchema(val.Config)
}

func (s *IndexerMap) loadServiceMap(flv string) (ServiceConfigOpts[indexTypes.Indexer], error) {
	concrete, err := indexTypes.IndexerTypeString(flv)
	if err != nil {
		return ServiceConfigOpts[indexTypes.Indexer]{}, fmt.Errorf("invalid provider name: %w allowed: %s", err, indexTypes.IndexerTypeStrings())
	}

	val, ok := s.initMap[concrete]
	if !ok {
		return ServiceConfigOpts[indexTypes.Indexer]{}, fmt.Errorf("unsupported downloader : %s", concrete)
	}

	return val, nil
}
