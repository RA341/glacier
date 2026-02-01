package auth

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/ra341/glacier/generated/auth/v1"
	"github.com/ra341/glacier/generated/auth/v1/v1connect"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewAuthServiceHandler(h)
}

func (h *Handler) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	typeString, err := SessionTypeString(req.Msg.TokenType)
	if err != nil {
		return nil, err
	}

	s, session, refresh, err := h.srv.Login(req.Msg.Username, req.Msg.Password, typeString)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&v1.LoginResponse{})

	addCookie(
		resp,
		CookieSessionToken, session,
		s.SessionTokenExpiry,
	)
	addCookie(
		resp,
		CookieRefreshToken, refresh,
		s.RefreshTokenExpiry,
	)

	return resp, nil
}

func addCookie[T any](
	resp *connect.Response[T],
	name, val string,
	expiry time.Time,
) {
	cookie := http.Cookie{
		Name:  name,
		Value: val,
		Path:  "/",

		Expires:  expiry,
		HttpOnly: true,
		Secure:   false,
	}

	resp.Header().Add("Set-Cookie", cookie.String())
	return
}

func (h *Handler) Register(ctx context.Context, req *connect.Request[v1.RegisterRequest]) (*connect.Response[v1.RegisterResponse], error) {
	//TODO implement me
	panic("implement me")
}
