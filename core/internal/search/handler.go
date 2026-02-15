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

func (h *Handler) GetGameMeta(ctx context.Context, c *connect.Request[v1.GetGameMetaRequest]) (*connect.Response[v1.GetGameMetaResponse], error) {
	meta, err := h.srv.GetGameMeta(ctx, c.Msg.Provider, c.Msg.GameDbId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetGameMetaResponse{
		Meta: meta.ToProto(),
	}), nil
}

func (h *Handler) SearchIndexers(ctx context.Context, req *connect.Request[v1.SearchIndexersRequest]) (*connect.Response[v1.SearchIndexersResponse], error) {
	search, err := h.srv.GetIndexerResults(req.Msg.Q.Indexer, req.Msg.Q.Query)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(search, func(t indexTypes.Source) *v1.GameSource {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.SearchIndexersResponse{
		Results: res,
	}), nil
}

func (h *Handler) SearchMetadata(ctx context.Context, req *connect.Request[v1.SearchMetadataRequest]) (*connect.Response[v1.SearchMetadataResponse], error) {
	search, err := h.srv.GetMetadataResults(ctx, req.Msg.Q.Indexer, req.Msg.Q.Query)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(search, func(t metaTypes.Meta) *v1.GameMetadata {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.SearchMetadataResponse{
		Metadata: res,
	}), nil
}
