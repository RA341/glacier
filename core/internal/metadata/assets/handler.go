package assets

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/assets/v1"
	"github.com/ra341/glacier/generated/assets/v1/v1connect"
	"github.com/ra341/glacier/internal/user"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}

	authMid := user.RequireRole(user.Omnissiah)

	path, handler := v1connect.NewAssetServiceHandler(h)
	return path, authMid(handler)
}

func (h *Handler) Delete(ctx context.Context, c *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	err := h.srv.delete(ctx, uint(c.Msg.ID))
	if err != nil {
		return nil, err
	}
	return &connect.Response[v1.DeleteResponse]{}, nil
}

func (h *Handler) Edit(ctx context.Context, c *connect.Request[v1.EditRequest]) (*connect.Response[v1.EditResponse], error) {
	var asset Asset
	asset.FromProto(c.Msg.Asset)

	err := h.srv.Edit(ctx, &asset)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.EditResponse{}), nil
}

func (h *Handler) GetTypes(ctx context.Context, c *connect.Request[v1.GetTypesRequest]) (*connect.Response[v1.GetTypesResponse], error) {
	//TODO implement me
	panic("implement me")
}
