package search

import (
	"context"
	"net/http"

	v1 "github.com/ra341/glacier/generated/search/v1"
	"github.com/ra341/glacier/generated/search/v1/v1connect"
	indexTypes "github.com/ra341/glacier/internal/indexer/types"
	metaTypes "github.com/ra341/glacier/internal/metadata/types"

	"github.com/ra341/glacier/pkg/listutils"

	"connectrpc.com/connect"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewSearchServiceHandler(h)
}

func (h *Handler) Search(ctx context.Context, req *connect.Request[v1.SearchRequest]) (*connect.Response[v1.SearchResponse], error) {
	search, err := h.srv.Search(req.Msg.Query)
	if err != nil {
		return nil, err
	}

	resu := listutils.ToMap(search, func(t indexTypes.IndexerGame) *v1.GameSearchResult {
		return &v1.GameSearchResult{
			Name:        t.Title,
			DownloadUrl: t.DownloadUrl,
			Size:        t.FileSize,
			UploadDate:  t.CreatedISO,
		}
	})

	return connect.NewResponse(&v1.SearchResponse{
		Results: resu,
	}), nil
}

func (h *Handler) Match(ctx context.Context, req *connect.Request[v1.MatchRequest]) (*connect.Response[v1.MatchResponse], error) {
	search, err := h.srv.Match(req.Msg.Query)
	if err != nil {
		return nil, err
	}

	resu := listutils.ToMap(search, func(t metaTypes.Meta) *v1.GameMetadata {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.MatchResponse{
		Metadata: resu,
	}), nil
}
