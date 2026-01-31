package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	frost "github.com/ra341/glacier/frost/app"
	"github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/pkg/argos"
)

//go:embed all:build
var uiDir embed.FS

func main() {
	prefixer := argos.WithPrefixer(config.EnvPrefix)
	envs := map[string]string{
		"LOGGER_VERBOSE": "true",
		"LOGGER_LEVEL":   "debug",
		"CONFIG_DIR":     "./config",
		//"GLACIER_URL":    "http://192.168.50.123:6699",
		"GLACIER_URL":     "http://localhost:6699",
		"CONFIG_YML_PATH": "./config/glacier.yml",
	}

	for key, value := range envs {
		err := os.Setenv(prefixer(key), value)
		if err != nil {
			log.Fatalf("could not set %s:%s\nerr:%v", key, value, err)
		}
	}

	subFS, err := fs.Sub(uiDir, "build")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading frontend directory: %w", err))
	}

	frost.NewTray(frost.WithDisableUI(), frost.WithUI(subFS))
}
