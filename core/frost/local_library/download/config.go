package download

import "net/http"

type Config struct {
	// set automatically by constructor
	httpCli *http.Client
	base    string

	MaxConcurrentFiles      int   `yaml:"maxConcurrentFiles"`
	MaxConcurrentFileChunks int   `yaml:"maxConcurrentFileChunks"`
	ChunkSizeInMB           int64 `yaml:"chunkSize"`
}
