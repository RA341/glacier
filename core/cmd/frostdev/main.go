package main

import (
	frost "github.com/ra341/glacier/frost/app"
	frostConf "github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
	"github.com/ra341/glacier/shared/config"
)

func init() {
	app.InitMeta(info.FlavourFrostDevelop)
}

func main() {
	envs := map[string]string{
		"LOGGER_VERBOSE": "true",
		"LOGGER_LEVEL":   "debug",
		"LOGGER_HTTP":    "true",
		"CONFIG_DIR":     "./config",
		//"GLACIER_URL":    "http://192.168.50.123:6699",
		"GLACIER_URL":     "http://localhost:6699",
		"CONFIG_YML_PATH": "./config/frost.yml",
		"START_SILENT":    "true",
	}

	config.SetEnvWithMap(frostConf.EnvPrefix, envs)

	frost.NewDesktop(api.WithUIProxy("http://localhost:5122"))
}
