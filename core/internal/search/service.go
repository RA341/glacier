package search

import (
	"context"

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

func (s *Service) GetGameMeta(ctx context.Context, metaProvider string, providerDbId string) (*metaTypes.Meta, error) {
	return s.metaSrv.Get(ctx, metaProvider, providerDbId)
}

func (s *Service) GetMetadataResults(ctx context.Context, name string, query string) ([]metaTypes.Meta, error) {
	if query == "" {
		return nil, nil
	}

	return s.metaSrv.Match(ctx, name, query)
}

func (s *Service) GetIndexerResults(name string, query string) ([]indexerTypes.Source, error) {
	if query == "" {
		return []indexerTypes.Source{}, nil
	}

	return s.indexer.Search(name, query)
}
