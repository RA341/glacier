package config

import "github.com/ra341/glacier/pkg/argos"

const EnvPrefix = "GLACIER"

func DefaultPrefixer() argos.Prefixer {
	return argos.WithPrefixer(EnvPrefix)
}
