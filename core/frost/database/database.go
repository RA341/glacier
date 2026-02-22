package database

import (
	"embed"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ra341/glacier/shared/database"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

////go:generate go run gorm.io/cli/gorm@latest gen -i .. -o ./generated/queries

//go:embed generated/migrations/*.sql
var migrationDir embed.FS

const migrationPath = "generated/migrations"

const dbName = "frost.db"

func New(basepath string, devMode bool) *gorm.DB {
	fullPath := filepath.Join(basepath, dbName)
	fullPath = fullPath + "?_journal_mode=WAL&_busy_timeout=5000"

	conn, err := database.InitConn(
		"sqlite3",
		fullPath,
		migrationDir,
		migrationPath,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to initialize database")
	}

	sql := sqlite.New(sqlite.Config{
		Conn: conn,
	})

	gormDB, err := database.InitGorm(sql, devMode)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to GORM")
	}

	return gormDB
}
