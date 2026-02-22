//go:build !postgres

package database

import (
	"embed"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ra341/glacier/shared/database"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const sqliteMigration = "generated/migrations/sqlite"

//go:embed generated/migrations/sqlite/*.sql
var migrationDir embed.FS

const dbName = "glacier.db"

type Config struct {
	// maintain the same env as config to place alongside config files
	DatabaseDir string `yaml:"databaseDir" env:"CONFIG_DIR" default:"./config" help:"path to the database dir" folder:""`
}

// NewSqlConn input will be a path to a DB
func NewSqlConn(config *Config) (gorm.Dialector, error) {
	log.Info().Msg("Using SQLite database")

	connStr := filepath.Join(config.DatabaseDir, dbName)
	connStr, err := filepath.Abs(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to get abs path of %s: %w", connStr, err)
	}

	connStr = connStr + "?_journal_mode=WAL&_busy_timeout=5000"

	conn, err := database.InitConn(
		"sqlite3",
		connStr,
		migrationDir,
		sqliteMigration,
	)
	if err != nil {
		return nil, err
	}

	return sqlite.New(sqlite.Config{
		Conn: conn,
	}), err
}
