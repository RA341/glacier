package main

import (
	"embed"
	"os"
	"path/filepath"

	frost "github.com/ra341/glacier/frost/app"
	frostConf "github.com/ra341/glacier/frost/config"
	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
	"github.com/ra341/glacier/shared/config"
	"github.com/rs/zerolog/log"
)

//go:embed all:build
var uiDir embed.FS

func init() {
	app.InitMeta(info.FlavourFrost)
}

func main() {
	frostHome, err := os.UserHomeDir()
	if err != nil {
		log.Warn().Err(err).Msg("could not get home directory, using current working directory")
	}

	frostHome = filepath.Join(frostHome, ".frost")
	configDir := filepath.Join(frostHome, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		frost.ShowErr("could not create config dir\n" + err.Error())
		os.Exit(1)
	}

	log.Info().Msgf("Using config dir: %s", frostHome)

	envs := map[string]string{
		"CONFIG_DIR":      configDir,
		"CONFIG_YML_PATH": filepath.Join(configDir, "frost.yml"),
		"LOG_DIR":         filepath.Join(frostHome, "log"),
	}

	config.SetEnvWithMap(frostConf.EnvPrefix, envs)

	frost.NewDesktop(api.WithUIFS(uiDir, "build"))
}
