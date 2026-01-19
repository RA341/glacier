package types

import "time"

type IndexerConfig = map[string]any

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=GameType -output=enum_game_type.go

// GameType identifies the type of files downloaded
type GameType int

const (
	// GameTypeUnknown is the default zero-value
	GameTypeUnknown GameType = iota

	// GameTypeInstaller means the files must be installed after download
	GameTypeInstaller

	// GameTypeStandalone means the files are ready-to-play after download
	GameTypeStandalone
)

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=IndexerType -output=enum_indexer_type.go
type IndexerType int

const (
	IndexerUnknown IndexerType = iota
	IndexerHydra
)

type Indexer interface {
	Search(query string) ([]Source, error)
	// Close is called for the indexer to clean itself up
	// generally can be empty unless needed
	Close()
}

// Source represents the data retrieved from the indexer
//
// except name and download url all other fields are optional
type Source struct {
	IndexerType IndexerType
	GameType    GameType

	Title       string
	DownloadUrl string

	// optional
	ImageURL   string
	FileSize   string
	CreatedISO string
}

type DownloadInfo struct {
	Title      string    `json:"title"`
	Uris       []string  `json:"uris"`
	UploadDate time.Time `json:"uploadDate"`
	FileSize   string    `json:"fileSize"`
}
