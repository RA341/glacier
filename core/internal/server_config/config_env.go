package server_config

type Prefixer func(in string) string

const defaultPrefix = "GLACIER"

func DefaultPrefixer() Prefixer {
	return WithPrefixer(defaultPrefix)
}

func WithPrefixer(envPrefix string) Prefixer {
	return func(in string) string {
		return envPrefix + "_" + in
	}
}
