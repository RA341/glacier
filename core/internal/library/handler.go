package library

import (
	"context"
	"net/http"
	"time"

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

func (h *Handler) TriggerTracker(ctx context.Context, req *connect.Request[v1.TriggerTrackerRequest]) (*connect.Response[v1.TriggerTrackerResponse], error) {
	h.srv.downloader.TriggerTracker()
	return connect.NewResponse(&v1.TriggerTrackerResponse{}), nil
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

func (h *Handler) List(ctx context.Context, c *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	list, err := h.srv.List(ctx, c.Msg.Query, uint(c.Msg.Offset), uint(c.Msg.Limit))
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t Game) *v1.Game {
		return &v1.Game{
			ID:            uint64(t.ID),
			CreatedAt:     t.CreatedAt.Format(time.RFC3339),
			EditedAt:      t.UpdatedAt.Format(time.RFC3339),
			Meta:          t.Meta.ToProto(),
			GameType:      t.GameType.String(),
			DownloadState: t.Download.ToProto(),
		}
	})

	return connect.NewResponse(&v1.ListResponse{
		GameList: res,
	}), nil
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
