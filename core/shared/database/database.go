package database

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm(dialector gorm.Dialector, devMode bool) (*gorm.DB, error) {
	gormLogLevel := logger.Silent
	if devMode {
		gormLogLevel = logger.Info
	}

	conf := &gorm.Config{
		Logger:         logger.Default.LogMode(gormLogLevel),
		PrepareStmt:    true,
		TranslateError: true,
	}

	db, err := gorm.Open(dialector, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	log.Info().Msg("Connected to database")
	return db, nil
}

func InitConn(
	dialect string,
	connectionStr string,
	migrationDir fs.FS,
	migrationPath string,
) (*sql.DB, error) {
	sqlDB, err := sql.Open(dialect, connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open raw sqlite connection: %w", err)
	}

	if err := migrate(sqlDB, migrationDir, migrationPath, dialect); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return sqlDB, err
}

func migrate(
	db *sql.DB,
	migrationDir fs.FS,
	migrationPath string,
	dialect string,
) error {
	goose.SetBaseFS(migrationDir)

	gzlog := GooseZerolog{}
	goose.SetLogger(gzlog)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	log.Info().Msg("Checking for database migrations...")

	if err := goose.Up(db, migrationPath); err != nil {
		return err
	}

	return nil
}

type GooseZerolog struct{}

func (g GooseZerolog) Fatalf(format string, v ...any) {
	log.Fatal().Msgf(format, v...)
}

func (g GooseZerolog) Printf(format string, v ...any) {
	log.Info().Msgf(format, v...)
}
