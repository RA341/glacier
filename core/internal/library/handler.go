package library

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/library/v1"
	"github.com/ra341/glacier/generated/library/v1/v1connect"
	"github.com/ra341/glacier/internal/metadata/types"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	svc := &Handler{srv: srv}
	return v1connect.NewLibraryServiceHandler(svc)
}

func (h *Handler) ListWithState(ctx context.Context, c *connect.Request[v1.ListWithStateRequest]) (*connect.Response[v1.ListWithStateResponse], error) {
	list, err := h.srv.ListDownloading(ctx, c.Msg.State)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t Game) *v1.Game {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.ListWithStateResponse{
		Game: res,
	}), nil
}

func (h *Handler) List(ctx context.Context, c *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	list, err := h.srv.List(ctx, c.Msg.Query, uint(c.Msg.Offset), uint(c.Msg.Limit))
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t Game) *v1.Game {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.ListResponse{
		GameList: res,
	}), nil
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	err := h.srv.Delete(ctx, uint(c.Msg.GameId))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteResponse{}), nil
}

func (h *Handler) TriggerTracker(ctx context.Context, req *connect.Request[v1.TriggerTrackerRequest]) (*connect.Response[v1.TriggerTrackerResponse], error) {
	h.srv.downloader.TriggerTracker()
	return connect.NewResponse(&v1.TriggerTrackerResponse{}), nil
}

func (h *Handler) Exists(ctx context.Context, req *connect.Request[v1.ExistsRequest]) (*connect.Response[v1.ExistsResponse], error) {
	typeString, err := types.ProviderTypeString(req.Msg.MetadataType)
	if err != nil {
		return nil, err
	}
	
	gameId, err := h.srv.store.Exists(typeString, req.Msg.MetadataGameId)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.ExistsResponse{
		GameId: uint64(gameId),
	}), nil
}

func (h *Handler) GetGame(ctx context.Context, c *connect.Request[v1.GetGameRequest]) (*connect.Response[v1.GetGameResponse], error) {
	get, err := h.srv.Get(ctx, uint(c.Msg.GameId))
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.GetGameResponse{
		Game: get.ToProto(),
	}), nil
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
