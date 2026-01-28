package indexer

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/indexer/v1"
	"github.com/ra341/glacier/generated/indexer/v1/v1connect"
	"github.com/ra341/glacier/internal/indexer/types"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewIndexerServiceHandler(h)
}

func (h Handler) GetGameType(ctx context.Context, c *connect.Request[v1.GetGameTypeRequest]) (*connect.Response[v1.GetGameTypeResponse], error) {
	res := listutils.ToMap(types.GameTypeStrings(), func(t string) *v1.GameType {
		return &v1.GameType{
			Name: t,
		}
	})

	return connect.NewResponse(&v1.GetGameTypeResponse{
		GameTypes: res,
	}), nil
}
