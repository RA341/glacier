package manager

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/indexer/v1"
	"github.com/ra341/glacier/generated/indexer/v1/v1connect"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewIndexerServiceHandler(h)
}

func (h *Handler) GetActiveIndexers(context.Context, *connect.Request[v1.GetActiveIndexersRequest]) (*connect.Response[v1.GetActiveIndexersResponse], error) {
	im := h.srv.cf.LoadAllIndexers()

	var indexers []*v1.Indexer
	for name := range im {
		indexers = append(indexers, &v1.Indexer{
			Name: name,
		})
	}

	return connect.NewResponse(&v1.GetActiveIndexersResponse{
		Indexers: indexers,
	}), nil
}
