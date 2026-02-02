package services_manager

import (
	"gorm.io/gorm"
)

type Store interface {
	Get(id string) (ServiceConfig, error)
	New(conf *ServiceConfig) error
	Edit(conf *ServiceConfig) error
	Delete(id uint) error
	ListAll(config ServiceType) ([]ServiceConfig, error)
	ListEnabled(ServiceType ServiceType) ([]ServiceConfig, error)
}

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=ServiceType -output=service_type.go
type ServiceType int

const (
	Indexer ServiceType = iota
	Metadata
	Downloader
)

type ServiceConfig struct {
	gorm.Model

	// same name & type & service is not allowed
	// ServiceType overall service type MetadataProvider or DownloadClient
	ServiceType ServiceType `gorm:"index:idx_service_config,unique"`
	Name        string      `gorm:"index:idx_service_config,unique"`

	Enabled bool
	// Flavour sub provider like IGDB or transmission
	Flavour string
	Config  map[string]any `gorm:"serializer:json"`
}
