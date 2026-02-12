package config

import (
	"github.com/ra341/glacier/pkg/argos"
	"github.com/ra341/glacier/shared/config"
)

const EnvPrefix = "GLACIER"
const GlacierYml = "glacier.yml"
const GlacierYmlPathEnv = "GLACIER_CONFIG_YML_PATH"

func New() *config.Service[Config] {
	prefixer := argos.WithPrefixer(EnvPrefix)

	return config.New[Config](
		prefixer,
		GlacierYmlPathEnv,
		GlacierYml,
	)
}
