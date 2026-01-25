package metadata

import (
	"github.com/ra341/glacier/internal/metadata/types"
)

type Get func(name string) (types.Provider, error)

type Service struct {
	// manager the providers config
	get Get
}

func New(get Get) *Service {
	return &Service{
		get: get,
	}
}

func (s *Service) Match(provider string, query string) ([]types.Meta, error) {
	val, err := s.get(provider)
	if err != nil {
		return nil, err
	}

	return val.GetMatches(query)
}
