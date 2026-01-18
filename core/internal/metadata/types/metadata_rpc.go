package types

import (
	"time"

	v1 "github.com/ra341/glacier/generated/search/v1"
)

func (m *Meta) ToProto() *v1.GameMetadata {
	return &v1.GameMetadata{
		ProviderType:  m.ProviderType.String(),
		ID:            m.GameDBID,
		Name:          m.Name,
		Summary:       m.Summary,
		Description:   m.Description,
		URL:           m.URL,
		ThumbnailURL:  m.ThumbnailURL,
		Videos:        m.Videos,
		Platforms:     m.Platforms,
		Genres:        m.Genres,
		Rating:        m.Rating,
		RatingCount:   uint32(m.RatingCount),
		ReleaseDate:   m.ReleaseDate.Format(time.RFC3339),
		ReleaseStatus: m.ReleaseStatus,
		Category:      m.Category,
	}
}

func (m *Meta) FromProto(rpcMeta *v1.GameMetadata) {
	// there shouldn't be an issue parsing, even if there, it's fine to set as default
	parsedDate, _ := time.Parse(time.RFC3339, rpcMeta.ReleaseDate)

	providerType, err := ProviderTypeString(rpcMeta.ProviderType)
	if err != nil {
		providerType = ProviderUnknown
	}

	m.ProviderType = providerType
	m.GameDBID = rpcMeta.ID
	m.Name = rpcMeta.Name
	m.Summary = rpcMeta.Summary
	m.Description = rpcMeta.Description
	m.URL = rpcMeta.URL
	m.ThumbnailURL = rpcMeta.ThumbnailURL
	m.Videos = rpcMeta.Videos
	m.Platforms = rpcMeta.Platforms
	m.Genres = rpcMeta.Genres
	m.Rating = rpcMeta.Rating
	m.RatingCount = uint(rpcMeta.RatingCount)
	m.ReleaseDate = parsedDate
	m.ReleaseStatus = rpcMeta.ReleaseStatus
	m.Category = rpcMeta.Category
}
