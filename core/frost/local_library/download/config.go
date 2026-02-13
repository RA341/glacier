package download

import "net/http"

type Config struct {
	// set automatically by constructor
	httpCli *http.Client
	base    string

	MaxConcurrentFiles      int `yaml:"maxConcurrentFiles" env:"MAX_FILES" default:"50" help:"Maximum number of concurrent files"`
	MaxConcurrentFileChunks int `yaml:"maxConcurrentFileChunks" env:"MAX_CHUNKS" default:"100" help:"Maximum number of chunks in a file to process"`
	ChunkSizeInMB           int `yaml:"chunkSize" env:"CHUNK_SIZE" default:"128" help:"file chunk size in MB"`
}

func (c *Config) GetChunkSize() int64 {
	return int64(c.ChunkSizeInMB) * MB
}
