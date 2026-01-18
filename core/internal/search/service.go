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

func (s *Service) Match(query string) ([]metaTypes.Meta, error) {
	// todo factor out provider
	return s.metaSrv.Match(metaTypes.ProviderIGDB, query)
}

func (s *Service) Search(query string) ([]indexerTypes.IndexerGame, error) {
	if query == "" {
		return nil, nil
	}

	panic("IMPLEMENT ME")
}
