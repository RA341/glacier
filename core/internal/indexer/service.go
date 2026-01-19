package indexer

import (
	"github.com/ra341/glacier/internal/indexer/manager"
	"github.com/ra341/glacier/internal/indexer/types"
)

type Service struct {
	manager *manager.Service
}

func New(indexer *manager.Service) *Service {
	return &Service{
		manager: indexer,
	}
}

func (s *Service) Search(indexerType types.IndexerType, searchTerm string) ([]types.Source, error) {
	get, err := s.manager.Get(indexerType)
	if err != nil {
		return nil, err
	}
	return get.Search(searchTerm)
}
