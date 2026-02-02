package main

import (
	"log"
	"os"

	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
)

func init() {
	app.InitMeta(info.FlavourDevelop)
}

func main() {
	prefixer := config.DefaultPrefixer()
	envs := map[string]string{
		"LOGGER_VERBOSE": "true",
		"LOGGER_LEVEL":   "debug",
		"LOGGER_HTTP":    "false",

		"SERVER_PORT":            "6699",
		"AUTH_DISABLE":           "false",
		"AUTH_OPEN_REGISTRATION": "true",
		"MAX_SESSIONS":           "3",

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

	devUi := api.WithProxy("http://localhost:5173")

	app.NewServer(api.WithUIProxy(devUi))
}
