package database

import (
	"embed"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/ra341/glacier/shared/database"
	"gorm.io/gorm"
)

////go:generate go run gorm.io/cli/gorm@latest gen -i .. -o ./generated/queries

//go:embed generated/migrations/*.sql
var migrationDir embed.FS

const migrationPath = "generated/migrations"

const dbName = "frost.db"

func New(basepath string, devMode bool) *gorm.DB {
	fullPath := filepath.Join(basepath, dbName)
	return database.New(fullPath, devMode, migrationDir, migrationPath)
}
