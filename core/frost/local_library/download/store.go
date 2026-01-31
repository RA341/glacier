package download

import (
	"io"
	"time"

	v1 "github.com/ra341/glacier/generated/frost_library/v1"
)

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=Status -output=enum_local_download_state.go
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

type Info struct {
	DownloadPath  string
	Status        Status
	StatusMessage string
	Started       time.Time
	Done          time.Time
}

func (g *Info) ToProto() *v1.DownloadInf {
	return &v1.DownloadInf{
		State:        g.Status.String(),
		Message:      g.StatusMessage,
		TimeStarted:  g.Started.Format(time.RFC3339),
		DownloadPath: g.DownloadPath,
	}
}
