package config

import (
	"github.com/ra341/glacier/pkg/argos"
	"github.com/ra341/glacier/shared/config"
)

const EnvPrefix = "FROST"
const FrostYml = "frost.yml"
const FrostYmlPathEnv = "FROST_CONFIG_YML_PATH"

func New(printConf bool) *config.Service[Config] {
	defaultPrefixer := argos.WithPrefixer(EnvPrefix)
	return config.New[Config](
		defaultPrefixer,
		FrostYmlPathEnv,
		FrostYml,
	)
}
