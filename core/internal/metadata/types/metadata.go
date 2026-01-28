package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=ProviderType -output=enum_provider_type.go
type ProviderType int

const (
	ProviderUnknown ProviderType = iota
	ProviderIGDB
)

type Provider interface {
	GetMatches(query string) ([]Meta, error)
	// GetFullMetadata id will be whatever ID the provider uses
	GetFullMetadata(id string) (*Meta, error)
}

type Meta struct {
	ProviderType ProviderType `gorm:"uniqueIndex:idx_provider_game"`

	GameDBID string `gorm:"uniqueIndex:idx_provider_game"`

	Name string
	// A short description/blurb of the game.
	ShortDesc string
	// A longer description of the game's plot.
	FullDesc string
	//  The direct link to the game's page on metadata provider.
	URL string
	// todo download the iamge
	ThumbnailURL string   // struct Image {local, remote string}
	Videos       []string `gorm:"serializer:json"`
	Platforms    []string `gorm:"serializer:json"`
	Genres       []string `gorm:"serializer:json"`

	// The average rating from critics/external sites
	Rating      string
	RatingCount uint

	ReleaseDate time.Time
	// The status of the game (e.g., Released, Alpha, Beta, Cancelled).
	ReleaseStatus string
	// Main Game, DLC, Expansion, Remake, Remaster etc
	Category string
}

// StringArray is a custom type for []string that serializes to JSON in SQLite
type StringArray []string

// Value handles saving to the database (Go -> SQL)
func (p *StringArray) Value() (driver.Value, error) {
	if p == nil || len(*p) == 0 {
		return "[]", nil // Store empty JSON array instead of NULL
	}
	b, err := json.Marshal(*p)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan handles loading from the database (SQL -> Go)
func (p *StringArray) Scan(src interface{}) error {
	if src == nil {
		*p = nil
		return nil
	}
	bytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal StringArray: expected []byte, got %T", src)
	}
	return json.Unmarshal(bytes, p)
}
