package download

import (
	"io"
	"time"

	v1 "github.com/ra341/glacier/generated/frost_library/v1"
)

//go:generate go run github.com/dmarkham/enumer@latest -type=Status -output=enum_local_download_state.go
type Status int

const (
	StatusQueued Status = iota
	StatusMetadata
	StatusDownloading
	StatusError
	StatusComplete
	StatusPaused
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

type LocalDownload struct {
	DownloadPath  string
	Status        Status
	StatusMessage string
	Started       time.Time
	Done          time.Time
}

func (g *LocalDownload) ToProto() *v1.LocalDownload {
	return &v1.LocalDownload{
		DownloadPath:  g.DownloadPath,
		Status:        g.Status.String(),
		StatusMessage: g.StatusMessage,
		Started:       g.Started.Format(time.RFC3339),
		Done:          g.Done.Format(time.RFC3339),
	}
}

func (g *LocalDownload) FromProto(download *v1.LocalDownload) {
	if download == nil {
		return
	}

	g.DownloadPath = download.DownloadPath

	statusString, err := StatusString(download.Status)
	if err != nil {
		statusString = StatusError
	}
	g.Status = statusString

	g.Started, err = time.Parse(time.RFC3339, download.Started)
	if err != nil {
		g.Started = time.Now()
	}

	g.Done, err = time.Parse(time.RFC3339, download.Done)
	if err != nil {
		g.Done = time.Now()
	}

	g.StatusMessage = download.StatusMessage
}
