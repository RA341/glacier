package downloader

import (
	"time"

	"github.com/rs/zerolog/log"
)

type ConfigLoader func() *Config

type Config struct {
	CheckInterval  string `yaml:"checkInterval" default:"30m" env:"DOWNLOAD_CHECK_TIME" help:"time between checking games status"`
	IncompletePath string `yaml:"incompletePath" default:"./incomplete" env:"INCOMPLETE_DIR" help:"places downloading games here"`
}

func (c *Config) Interval() time.Duration {
	duration, err := time.ParseDuration(c.CheckInterval)
	if err != nil {
		const defaultCheckInterval = 30 * time.Minute
		log.Warn().Err(err).Str("interval", c.CheckInterval).Msg("can't parse check interval")
		duration = defaultCheckInterval
	}

	return duration
}
