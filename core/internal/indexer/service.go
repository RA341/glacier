package indexer

import (
	"github.com/ra341/glacier/internal/indexer/types"
)

type Get func(name string) (types.Indexer, error)

type Service struct {
	get Get
}

func New(get Get) *Service {
	return &Service{
		get: get,
	}
}

func (s *Service) Search(indexerType string, searchTerm string) ([]types.Source, error) {
	get, err := s.get(indexerType)
	if err != nil {
		return nil, err
	}
	return get.Search(searchTerm)
}
