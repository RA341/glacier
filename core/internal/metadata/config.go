package metadata

import (
	"fmt"

	"github.com/ra341/glacier/internal/metadata/types"
)

type ConfigLoader func() *Config

type Config struct {
	// client_name: map[...client]...conf
	Providers map[string]types.ProviderConfig `yaml:"providers"`
}

func (c *Config) GetCli(cli types.ProviderType) (types.ProviderConfig, error) {
	val, ok := c.Providers[cli.String()]
	if !ok {
		return types.ProviderConfig{}, fmt.Errorf("client not found for %s", cli)
	}
	return val, nil
}

func (c *Config) SetCli(cli types.ProviderType, conf types.ProviderConfig) {
	c.Providers[cli.String()] = conf
}
