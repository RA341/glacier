package config

type Prefixer func(in string) string

const EnvPrefix = "GLACIER"

func DefaultPrefixer() Prefixer {
	return WithPrefixer(EnvPrefix)
}

func WithPrefixer(envPrefix string) Prefixer {
	return func(in string) string {
		return envPrefix + "_" + in
	}
}
