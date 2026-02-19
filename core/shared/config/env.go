package config

import (
	"log"
	"os"

	"github.com/ra341/glacier/pkg/argos"
)

func SetEnvWithMap(envPrefix string, envs map[string]string) {
	prefixer := argos.WithPrefixer(envPrefix)
	for key, value := range envs {
		err := os.Setenv(prefixer(key), value)
		if err != nil {
			log.Fatalf("could not set %s:%s\nerr:%v", key, value, err)
		}
	}
}
