package download

import "io"

type Status int

const (
	StatusQueued Status = iota
	StatusMetadata
	StatusDownloading
	StatusError
	StatusComplete
)

type ChunkState int

const (
	ChunkQueued ChunkState = iota
	ChunkComplete
	ChunkError
)

type Chunk struct {
	Start int64      `json:"start"`
	End   int64      `json:"end"`
	State ChunkState `json:"status"`
}

type CacheStore interface {
	io.Closer
	GetFileList() ([]string, error)
	GetChunkLen(file string) (int, error)

	Add(file string, chunk []Chunk) error
	Get(file string) ([]Chunk, bool, error)
	Update(file string, index int, chunk *Chunk) error
	Progress() (progress []FileProgress, err error)
}
