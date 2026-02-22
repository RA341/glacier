package database

import (
	"github.com/ra341/glacier/shared/database"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func New(config *Config, devMode bool) *gorm.DB {
	sqlDB, err := NewSqlConn(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	gormDB, err := database.InitGorm(sqlDB, devMode)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to GORM")
	}

	return gormDB
}
