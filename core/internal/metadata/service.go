package metadata

import (
	"github.com/ra341/glacier/internal/metadata/manager"
	"github.com/ra341/glacier/internal/metadata/types"
)

type Service struct {
	// manager the providers config
	pm *manager.Service
}

func New(man *manager.Service) *Service {
	return &Service{
		pm: man,
	}
}

func (s *Service) Match(provider types.ProviderType, query string) ([]types.Meta, error) {
	val, err := s.pm.Get(provider)
	if err != nil {
		return nil, err
	}

	return val.GetMatches(query)
}
