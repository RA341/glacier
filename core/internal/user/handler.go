package user

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/user/v1"
	"github.com/ra341/glacier/generated/user/v1/v1connect"
	"github.com/ra341/glacier/pkg/listutils"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{
		srv: srv,
	}
	return v1connect.NewUserServiceHandler(h)
}

func (h *Handler) Self(ctx context.Context, c *connect.Request[v1.SelfRequest]) (*connect.Response[v1.SelfResponse], error) {
	actionUser, err := GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.SelfResponse{
		User: actionUser.ToProto(),
	}), nil
}

func (h *Handler) List(ctx context.Context, req *connect.Request[v1.ListRequest]) (*connect.Response[v1.ListResponse], error) {
	actionUser, err := GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	list, err := h.srv.List(req.Msg.Query, actionUser)
	if err != nil {
		return nil, err
	}

	res := listutils.ToMap(list, func(t User) *v1.User {
		return t.ToProto()
	})

	return connect.NewResponse(&v1.ListResponse{
		Users: res,
	}), nil
}

func (h *Handler) ListRoles(ctx context.Context, c *connect.Request[v1.ListRolesRequest]) (*connect.Response[v1.ListRolesResponse], error) {
	res := listutils.ToMap(RoleStrings(), func(t string) *v1.Role {
		return &v1.Role{
			Name: t,
		}
	})
	return connect.NewResponse(&v1.ListRolesResponse{Roles: res}), nil
}

func (h *Handler) Delete(ctx context.Context, req *connect.Request[v1.DeleteRequest]) (*connect.Response[v1.DeleteResponse], error) {
	actionBy, err := GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	err = h.srv.Delete(uint(req.Msg.Id), actionBy)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.DeleteResponse{}), nil
}

func (h *Handler) New(ctx context.Context, req *connect.Request[v1.NewRequest]) (*connect.Response[v1.NewResponse], error) {
	actionBy, err := GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	var editUser User
	err = editUser.FromProto(req.Msg.User)
	if err != nil {
		return nil, err
	}

	err = h.srv.New(
		editUser.Username,
		editUser.EncryptedPassword,
		editUser.Role,
		actionBy,
	)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.NewResponse{}), nil
}

func (h *Handler) Edit(ctx context.Context, req *connect.Request[v1.EditRequest]) (*connect.Response[v1.EditResponse], error) {
	actionBy, err := GetUserCtx(ctx)
	if err != nil {
		return nil, err
	}

	var editUser User
	err = editUser.FromProto(req.Msg.User)
	if err != nil {
		return nil, err
	}

	err = h.srv.Edit(&editUser, actionBy)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.EditResponse{}), nil
}
