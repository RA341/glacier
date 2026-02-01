package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/info"
)

func init() {
	app.InitMeta(info.FlavourDevelop)
}

//go:embed all:build
var uifs embed.FS

func main() {
	prefixer := config.DefaultPrefixer()
	envs := map[string]string{
		"LOGGER_VERBOSE": "true",
		"LOGGER_LEVEL":   "debug",
		"LOGGER_HTTP":    "true",

		"SERVER_PORT": "6699",

		"CONFIG_DIR":      "./config",
		"GAME_DIR":        "./gamestop",
		"CONFIG_YML_PATH": "./config/glacier.yml",
	}

	for key, value := range envs {
		err := os.Setenv(prefixer(key), value)
		if err != nil {
			log.Fatalf("could not set %s:%s\nerr:%v", key, value, err)
		}
	}

	subFS, err := fs.Sub(uifs, "build")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading frontend directory: %w", err))
	}

	app.NewServer(app.WithUIFS(subFS))
}
