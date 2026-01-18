package search

import (
	"context"
	"net/http"
	"strings"

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

func (h *Handler) SearchIndexers(ctx context.Context, req *connect.Request[v1.SearchIndexersRequest]) (*connect.Response[v1.SearchIndexersResponse], error) {
	search, err := h.srv.GetIndexerResults(req.Msg.Q.Indexer, req.Msg.Q.Query)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(search, func(t indexTypes.IndexerGame) *v1.GameIndexer {
		return &v1.GameIndexer{
			Title:       t.Title,
			DownloadUrl: t.DownloadUrl,
			ImageURL:    t.ImageURL,
			FileSize:    t.FileSize,
			CreatedISO:  t.CreatedISO,
		}
	})

	return connect.NewResponse(&v1.SearchIndexersResponse{
		Results: res,
	}), nil
}

func (h *Handler) SearchMetadata(ctx context.Context, req *connect.Request[v1.SearchMetadataRequest]) (*connect.Response[v1.SearchMetadataResponse], error) {
	search, err := h.srv.GetMetadataResults(req.Msg.Q.Indexer, req.Msg.Q.Query)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(search, func(t metaTypes.Meta) *v1.GameMetadata {
		t.ThumbnailURL = strings.Replace(t.ThumbnailURL, "t_thumb", "t_cover_big", 1)
		return t.ToProto()
	})

	return connect.NewResponse(&v1.SearchMetadataResponse{
		Metadata: res,
	}), nil
}
