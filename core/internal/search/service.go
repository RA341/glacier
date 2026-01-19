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

func (s *Service) GetMetadataResults(providerStr string, query string) ([]metaTypes.Meta, error) {
	if query == "" {
		return nil, nil
	}

	prov, err := metaTypes.ProviderTypeString(providerStr)
	if err != nil {
		return nil, err
	}

	return s.metaSrv.Match(prov, query)
}

func (s *Service) GetIndexerResults(indexerStr string, query string) ([]indexerTypes.Source, error) {
	if query == "" {
		return []indexerTypes.Source{}, nil
	}
	ind, err := indexerTypes.IndexerTypeString(indexerStr)
	if err != nil {
		return nil, err
	}

	return s.indexer.Search(ind, query)
}
