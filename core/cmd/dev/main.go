package main

import (
	"log"
	"os"

	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/info"
)

func init() {
	app.InitMeta(info.FlavourDevelop)
}

func main() {
	prefixer := config.DefaultPrefixer()
	envs := map[string]string{
		"LOGGER_VERBOSE":  "true",
		"LOGGER_LEVEL":    "debug",
		"SERVER_PORT":     "6699",
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

	//file, err := app.LoadUIFromDir("./web")
	//if err != nil {
	//	log.Fatalf("could not load UI from file:%s\nerr:%v", file, err)
	//	return
	//}
	//app.NewServer(app.WithUIFS(file))

	app.NewServer()
}
