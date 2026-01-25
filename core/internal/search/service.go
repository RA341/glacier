package search

import (
	"github.com/ra341/glacier/internal/indexer"
	indexerTypes "github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/internal/metadata"
	metaTypes "github.com/ra341/glacier/internal/metadata/types"
)

type Service struct {
	metaSrv *metadata.Service
	indexer *indexer.Service
}

func New(metaSrv *metadata.Service, indexer *indexer.Service) *Service {
	return &Service{
		metaSrv: metaSrv,
		indexer: indexer,
	}
}

func (s *Service) GetMetadataResults(name string, query string) ([]metaTypes.Meta, error) {
	if query == "" {
		return nil, nil
	}

	return s.metaSrv.Match(name, query)
}

func (s *Service) GetIndexerResults(name string, query string) ([]indexerTypes.Source, error) {
	if query == "" {
		return []indexerTypes.Source{}, nil
	}

	return s.indexer.Search(name, query)
}
