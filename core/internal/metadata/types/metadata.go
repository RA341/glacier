package types

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ra341/glacier/internal/metadata/assets"
)

//go:generate go run github.com/dmarkham/enumer@latest -sql -type=ProviderType -output=enum_provider_type.go
type ProviderType int

const (
	ProviderUnknown ProviderType = iota
	ProviderIGDB
)

// todo
////go:generate go run github.com/dmarkham/enumer@latest -type=Platform -output=enum_platform.go
//type Platform int
//
//const (
//	PlatformUnknown Platform = iota
//	PlatformWindows
//	PlatformLinux
//	PlatformSwitch
//	PlatformAndroid
//	PlatformIOS
//	PlatformMacOS
//)

type Provider interface {
	GetMatches(ctx context.Context, query string) ([]Meta, error)
	// GetFullMetadata id will be whatever ID the provider uses
	GetFullMetadata(ctx context.Context, id string) (*Meta, error)
}

type Meta struct {
	ProviderType ProviderType `gorm:"uniqueIndex:idx_provider_game"`
	// ID assigned buy the metadata provider
	GameDBID string `gorm:"uniqueIndex:idx_provider_game"`

	//  The direct link to the game's page on metadata provider.
	URL  string
	Name string
	// A short description/blurb of the game.
	ShortDesc string
	// A longer description of the game's plot.
	FullDesc string

	Assets []assets.Asset `gorm:"foreignKey:GameID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// The average rating from critics/external sites
	Rating      string
	RatingCount uint

	ReleaseDate time.Time
	// The status of the game (e.g., Released, Alpha, Beta, Cancelled).
	ReleaseStatus string

	// Main Game, DLC, Expansion, Remake, Remaster etc
	Category  string
	Platforms []string `gorm:"serializer:json"`

	// todo create a genre table and auto insert when a unique genre is detected
	Genres []string `gorm:"serializer:json"`
}

func (m *Meta) Thumbnail() assets.Asset {
	for _, asset := range m.Assets {
		if asset.Type == assets.AssetThumbnail {
			return asset
		}
	}
	return assets.Asset{}
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
func (p *StringArray) Scan(src any) error {
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
