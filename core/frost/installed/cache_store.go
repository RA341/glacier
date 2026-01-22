package download

import "io"

type ChunkState int

const (
	Queued ChunkState = iota
	Complete
	Error
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
}
