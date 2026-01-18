package database

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

////go:generate go run gorm.io/cli/gorm@latest gen -i .. -o ./generated/queries

//go:embed generated/migrations/*.sql
var migrationDir embed.FS

const migrationPath = "generated/migrations"

const dbName = "glacier.db"

func New(basepath string, devMode bool) *gorm.DB {
	gormDB, err := connect(basepath, devMode)
	if err != nil {
		// Using zerolog for consistency
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}
	return gormDB
}

func connect(basepath string, devMode bool) (*gorm.DB, error) {
	targetPath := filepath.Join(basepath, dbName)
	dbpath, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, fmt.Errorf("unable to get abs path of %s: %w", targetPath, err)
	}

	connectionStr := dbpath + "?_journal_mode=WAL&_busy_timeout=5000"

	sqlDB, err := sql.Open("sqlite3", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open raw sqlite connection: %w", err)
	}

	if err := migrate(sqlDB); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	gormLogLevel := logger.Silent
	if devMode {
		gormLogLevel = logger.Info
	}

	conf := &gorm.Config{
		Logger:      logger.Default.LogMode(gormLogLevel),
		PrepareStmt: true,
	}

	db, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, conf)
	if err != nil {
		return nil, fmt.Errorf("failed to open gorm connection: %w", err)
	}

	log.Info().Str("path", dbpath).Msg("Connected to database")
	return db, nil
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(migrationDir)

	gzlog := GooseZerolog{}
	goose.SetLogger(gzlog)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	log.Info().Msg("Checking for database migrations...")

	if err := goose.Up(db, migrationPath); err != nil {
		return err
	}

	return nil
}

type GooseZerolog struct{}

func (g GooseZerolog) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (g GooseZerolog) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}
