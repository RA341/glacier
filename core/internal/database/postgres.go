//go:build postgres

package database

import (
	"embed"
	"fmt"

	"github.com/ra341/glacier/shared/database"
	"github.com/rs/zerolog/log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `yaml:"Host" env:"POSTGRES_HOST" default:"-" help:"postgres host"`
	Username string `yaml:"Username" env:"POSTGRES_USER" default:"-" help:"postgres username"`
	Password string `yaml:"Password" env:"POSTGRES_PASS" default:"-" help:"postgres password"`
	Database string `yaml:"Database" env:"POSTGRES_DB" default:"-" help:"postgres database"`
}

//go:embed generated/migrations/postgres/*.sql
var migrationDir embed.FS

const postgresMigration = "generated/migrations/postgres"

// // NewSqlConn input will be a path to a urlExample := "postgres://username:password@localhost:5432/database_name"
func NewSqlConn(config *Config) (gorm.Dialector, error) {
	log.Info().Msg("Using Postgres database")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Database,
	)

	conn, err := database.InitConn(
		"pgx",
		connStr,
		migrationDir,
		postgresMigration,
	)
	if err != nil {
		return nil, err
	}

	dia := postgres.New(postgres.Config{
		Conn: conn,
	})

	return dia, err
}
