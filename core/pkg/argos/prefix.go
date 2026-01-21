package argos

type Prefixer func(in string) string

func WithPrefixer(envPrefix string) Prefixer {
	return func(in string) string {
		return envPrefix + "_" + in
	}
}
