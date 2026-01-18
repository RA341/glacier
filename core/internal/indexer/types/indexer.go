package types

import "time"

type IndexerConfig = map[string]any

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=IndexerType -output=enum_indexer_type.go
type IndexerType int

const (
	IndexerUnknown IndexerType = iota
	IndexerHydra
)

type Indexer interface {
	Search(query string) ([]IndexerGame, error)
	// Close is called for the indexer to clean itself up
	// generally can be empty unless needed
	Close()
}

// IndexerGame except name and download url all other fields are optional
type IndexerGame struct {
	Title       string
	DownloadUrl string

	// optional
	ImageURL    string
	Description string
	FileSize    string
	CreatedISO  string
}

type DownloadInfo struct {
	Title      string    `json:"title"`
	Uris       []string  `json:"uris"`
	UploadDate time.Time `json:"uploadDate"`
	FileSize   string    `json:"fileSize"`
}
