package download

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/time/rate"
)

type Config struct {
	// set automatically by constructor
	httpCli *http.Client
	getBase func() string

	SpeedThrottleInMB  int `yaml:"speedThrottleInMB" env:"SPEED_THROTTLE_IN_MB" default:"0" help:"Set global speed limit for download, Restart is required, 0 is unlimited."`
	MaxConcurrentGames int `yaml:"maxConcurrentGames" env:"MAX_GAMES" default:"1" help:"how many games to download at a time, -1 to set unlimited"`

	// download chunking
	MaxConcurrentFiles      int `yaml:"maxConcurrentFiles" env:"MAX_FILES" default:"50" help:"Maximum number of concurrent files"`
	MaxConcurrentFileChunks int `yaml:"maxConcurrentFileChunks" env:"MAX_CHUNKS" default:"100" help:"Maximum number of chunks in a file to process"`
	ChunkSizeInMB           int `yaml:"chunkSize" env:"CHUNK_SIZE" default:"128" help:"file chunk size in MB"`
}

func (c *Config) UrlBase() string {
	return fmt.Sprintf(
		"%s/library/download",
		strings.TrimRight(c.getBase(), "/"),
	)
}

func (c *Config) GetChunkSize() int64 {
	return int64(c.ChunkSizeInMB) * MB
}

func (c *Config) GetSpeedThrottle() rate.Limit {
	limit := rate.Limit(c.SpeedThrottleInMB * MB)
	return limit
}
