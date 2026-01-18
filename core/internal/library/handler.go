package library

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/library/v1"
	"github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	svc := &Handler{srv: srv}
	return v1connect.NewLibraryServiceHandler(svc)
}

func (h *Handler) Add(ctx context.Context, req *connect.Request[v1.AddRequest]) (*connect.Response[v1.AddResponse], error) {
	var game = &Game{}
	game.FromProto(req.Msg.Game)

	err := h.srv.Add(ctx, game)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.AddResponse{}), nil
}

func (h *Handler) GetGameType(ctx context.Context, c *connect.Request[v1.GetGameTypeRequest]) (*connect.Response[v1.GetGameTypeResponse], error) {
	res := listutils.ToMap(GameTypeStrings(), func(t string) *v1.GameType {
		return &v1.GameType{
			Name: t,
		}
	})

	return connect.NewResponse(
		&v1.GetGameTypeResponse{
			GameTypes: res,
		},
	), nil
}
