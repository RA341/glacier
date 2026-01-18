package indexer

import (
	"fmt"

	"github.com/ra341/glacier/internal/indexer/types"
)

type ConfigLoader func() *Config

type Config struct {
	Indexers map[string]types.IndexerConfig `yaml:"indexers"`
}

func (c *Config) SetCli(cli types.IndexerType, conf types.IndexerConfig) {
	c.Indexers[cli.String()] = conf
}

func (c *Config) GetCli(cli types.IndexerType) (types.IndexerConfig, error) {
	val, ok := c.Indexers[cli.String()]
	if !ok {
		return types.IndexerConfig{}, fmt.Errorf("indexer not found for %s", cli)
	}
	return val, nil
}
