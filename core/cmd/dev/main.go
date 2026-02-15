package main

import (
	"log"
	"os"

	"github.com/ra341/glacier/internal/app"
	galcierConfig "github.com/ra341/glacier/internal/config"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/pkg/argos"
	"github.com/ra341/glacier/shared/api"
)

func init() {
	app.InitMeta(info.FlavourDevelop)
}

func main() {
	prefixer := argos.WithPrefixer(galcierConfig.EnvPrefix)

	envs := map[string]string{
		"LOGGER_VERBOSE": "true",
		"LOGGER_LEVEL":   "debug",
		"LOGGER_HTTP":    "false",

		"SERVER_PORT":             "6699",
		"AUTH_DISABLE":            "false",
		"AUTH_OPEN_REGISTRATION":  "true",
		"AUTH_MAX_SESSIONS":       "3",
		"AUTH_OIDC_ENABLE":        "false",
		"AUTH_OIDC_ISSUER":        "https://auth.localhost",
		"AUTH_OIDC_CLIENT_ID":     "51668c29-ca37-4bd5-b4b2-dbc9c953ea6d",
		"AUTH_OIDC_CLIENT_SECRET": "xAsPzGKy4MGoExAy0r7kLbBN5Hvr9Pg9",
		"AUTH_OIDC_REDIRECT_URL":  "http://localhost:6699/api/server/public/auth/oidc/callback",

		"USE_YTDLP":    "true",
		"YT_RELAY_URL": "http://localhost:3002",

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

	app.NewServer(api.WithUIProxy("http://localhost:5173"))
}
