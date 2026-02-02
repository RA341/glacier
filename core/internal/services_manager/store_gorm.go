package services_manager

import (
	"context"

	"gorm.io/gorm"
)

type ServiceConfigManagerGorm struct {
	serviceType string
	db          *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &ServiceConfigManagerGorm{
		//serviceType: serviceType,
		db: db,
	}
}

func (s *ServiceConfigManagerGorm) Q() *gorm.DB {
	return s.db.WithContext(context.Background()).Model(&ServiceConfig{})
}

func (s *ServiceConfigManagerGorm) ListAll(config ServiceType) ([]ServiceConfig, error) {
	var dest []ServiceConfig
	err := s.Q().Where("service_type = ?", config).Order("created_at desc").Find(&dest).Error
	return dest, err
}

func (s *ServiceConfigManagerGorm) ListEnabled(ServiceType ServiceType) ([]ServiceConfig, error) {
	var dest []ServiceConfig
	err := s.Q().Order("created_at desc").
		Where("service_type = ?", ServiceType).
		Where("enabled = ?", true).
		Find(&dest).
		Error

	return dest, err
}

func (s *ServiceConfigManagerGorm) Get(id string) (ServiceConfig, error) {
	var dest ServiceConfig
	err := s.Q().Where("name = ?", id).First(&dest).Error
	return dest, err
}

func (s *ServiceConfigManagerGorm) New(conf *ServiceConfig) error {
	return s.Q().Create(conf).Error
}

func (s *ServiceConfigManagerGorm) Edit(conf *ServiceConfig) error {
	err := s.Q().
		Where("id = ?", conf.ID).
		Select("name", "enabled", "config").
		Updates(conf).Error
	return err
}

func (s *ServiceConfigManagerGorm) Delete(id uint) error {
	return s.Q().Unscoped().Delete(&ServiceConfig{}, id).Error
}
