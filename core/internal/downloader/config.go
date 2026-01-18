package downloader

import (
	"fmt"
	"time"

	"github.com/ra341/glacier/internal/downloader/types"
	"github.com/rs/zerolog/log"
)

type ConfigLoader func() *Config

type Config struct {
	IncompleteDownloadPath string `yaml:"incompleteDownloadPath" default:"./incomplete" env:"INCOMPLETE_DOWNLOAD_DIR"`
	CheckInterval          string `yaml:"checkInterval" default:"30m" env:"DOWNLOAD_CHECK_INTERVAL"`

	// client_name: map[...client]...conf
	Clients map[string]types.ClientConfig `yaml:"clients"`
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

func (c *Config) GetClient(cli types.ClientType) (types.ClientConfig, error) {
	val, ok := c.Clients[cli.String()]
	if !ok {
		return types.ClientConfig{}, fmt.Errorf("client not found for %s", cli)
	}
	return val, nil
}

func (c *Config) SetClient(cli types.ClientType, conf types.ClientConfig) {
	c.Clients[cli.String()] = conf
}
