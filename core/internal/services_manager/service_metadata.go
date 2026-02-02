package services_manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/metadata/providers/igdb"
	metadata "github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/mapsct"
	"github.com/ra341/glacier/pkg/syncmap"
)

type MetadataMap struct {
	initMap map[metadata.ProviderType]ServiceConfigOpts[metadata.Provider]
	metaMap syncmap.Map[string, metadata.Provider]
	store   Store
}

func NewMetadataMap(store Store) ServiceConfigMap[metadata.Provider] {
	return &MetadataMap{
		store: store,
		initMap: map[metadata.ProviderType]ServiceConfigOpts[metadata.Provider]{
			metadata.ProviderIGDB: {InitFn: igdb.New, Config: igdb.Config{}},
		},
	}
}
func (m *MetadataMap) LoadService(id string) (metadata.Provider, error) {
	val, ok := m.metaMap.Load(id)
	if ok {
		return val, nil
	}

	conf, err := m.store.Get(id)
	if err != nil {
		return nil, err
	}

	service, err := m.initService(&conf)
	if err != nil {
		return nil, err
	}

	m.metaMap.Store(id, service)

	return service, nil
}

func (m *MetadataMap) initService(conf *ServiceConfig) (metadata.Provider, error) {
	serviceMap, err := m.loadServiceMap(conf.Flavour)
	if err != nil {
		return nil, err
	}

	return serviceMap.InitFn(conf.Config)
}

func (m *MetadataMap) getServiceSchema(flv string) ([]mapsct.FieldSchema, error) {
	serviceMap, err := m.loadServiceMap(flv)
	if err != nil {
		return nil, err
	}

	return mapsct.GetSchema(serviceMap.Config)
}

func (m *MetadataMap) loadServiceMap(flv string) (ServiceConfigOpts[metadata.Provider], error) {
	concrete, err := metadata.ProviderTypeString(flv)
	if err != nil {
		return ServiceConfigOpts[metadata.Provider]{}, fmt.Errorf("invalid provider name: %w allowed: %s", err, metadata.ProviderTypeStrings())
	}

	val, ok := m.initMap[concrete]
	if !ok {
		return ServiceConfigOpts[metadata.Provider]{}, fmt.Errorf("unsupported downloader : %s", concrete)
	}

	return val, nil
}
