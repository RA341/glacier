package config_manager

import (
	"fmt"

	"github.com/ra341/glacier/internal/downloader/clients/transmission"
	downloaderTypes "github.com/ra341/glacier/internal/downloader/types"
	"github.com/ra341/glacier/pkg/mapsct"
	"github.com/ra341/glacier/pkg/syncmap"
)

type InitFn[T any] func(config map[string]any) (T, error)

type ServiceConfigOpts[T any] struct {
	InitFn InitFn[T]
	Config any
}

type ServiceConfigMap[T any] interface {
	initService(conf *ServiceConfig) (T, error)
	LoadService(id string) (T, error)
	getServiceSchema(flv string) ([]mapsct.FieldSchema, error)
	loadServiceMap(flv string) (ServiceConfigOpts[T], error)
}

type DownloaderMap struct {
	configMap map[downloaderTypes.ClientType]ServiceConfigOpts[downloaderTypes.Downloader]
	sm        syncmap.Map[string, downloaderTypes.Downloader]
	store     Store
}

func NewDownloaderMap(store Store) ServiceConfigMap[downloaderTypes.Downloader] {
	return &DownloaderMap{
		store: store,
		configMap: map[downloaderTypes.ClientType]ServiceConfigOpts[downloaderTypes.Downloader]{
			downloaderTypes.ClientTransmission: {
				InitFn: transmission.New,
				Config: transmission.Config{},
			},
		},
	}
}

func (d *DownloaderMap) LoadService(id string) (downloaderTypes.Downloader, error) {
	val, ok := d.sm.Load(id)
	if ok {
		return val, nil
	}

	conf, err := d.store.Get(id)
	if err != nil {
		return nil, err
	}

	service, err := d.initService(&conf)
	if err != nil {
		return nil, err
	}

	d.sm.Store(id, service)

	return service, nil
}

func (d *DownloaderMap) initService(conf *ServiceConfig) (downloaderTypes.Downloader, error) {
	val, err := d.loadServiceMap(conf.Flavour)
	if err != nil {
		return nil, err
	}
	return val.InitFn(conf.Config)
}

func (d *DownloaderMap) getServiceSchema(flv string) ([]mapsct.FieldSchema, error) {
	val, err := d.loadServiceMap(flv)
	if err != nil {
		return nil, err
	}
	return mapsct.GetSchema(val.Config)
}

func (d *DownloaderMap) loadServiceMap(flv string) (ServiceConfigOpts[downloaderTypes.Downloader], error) {
	concrete, err := downloaderTypes.ClientTypeString(flv)
	if err != nil {
		return ServiceConfigOpts[downloaderTypes.Downloader]{}, fmt.Errorf("invalid provider name: %w allowed: %s", err, downloaderTypes.ClientTypeStrings())
	}

	val, ok := d.configMap[concrete]
	if !ok {
		return ServiceConfigOpts[downloaderTypes.Downloader]{}, fmt.Errorf("unsupported downloader : %s", concrete)
	}

	return val, nil
}
