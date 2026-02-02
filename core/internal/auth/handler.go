package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"connectrpc.com/connect"
	frost "github.com/ra341/glacier/frost/app"
	v1 "github.com/ra341/glacier/generated/auth/v1"
	"github.com/ra341/glacier/generated/auth/v1/v1connect"
	"github.com/ra341/glacier/internal/user"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	srv *Service
}

func NewHandler(srv *Service) (string, http.Handler) {
	h := &Handler{srv: srv}
	return v1connect.NewAuthServiceHandler(h)
}

func (h *Handler) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	val := req.Header().Get(frost.FrostHeader)
	var sessionType = Web
	if val == "true" {
		sessionType = Frost
	}

	s, session, refresh, err := h.srv.Login(
		req.Msg.Username,
		req.Msg.Password,
		sessionType,
	)
	if err != nil {
		return nil, err
	}

	resp := connect.NewResponse(&v1.LoginResponse{})

	writeSessionCookie(resp, &s, session, refresh)

	return resp, nil
}

func (h *Handler) Register(ctx context.Context, req *connect.Request[v1.RegisterRequest]) (*connect.Response[v1.RegisterResponse], error) {
	if req.Msg.Password != req.Msg.PasswordVerify {
		return nil, fmt.Errorf("password do not match")
	}

	err := h.srv.Register(req.Msg.Username, req.Msg.Password, user.TechPriest, nil)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.RegisterResponse{}), nil
}

func (h *Handler) Logout(ctx context.Context, req *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error) {
	cookies := req.Header().Get("Cookie")
	err := h.getSessionCookie(ctx, cookies)
	if err != nil {
		log.Warn().Err(err).Msg("logout err")
	}

	// always return success
	return connect.NewResponse(&v1.LogoutResponse{}), nil
}

func (h *Handler) getSessionCookie(ctx context.Context, cookies string) error {
	cookie, err := http.ParseCookie(cookies)
	if err != nil {
		return fmt.Errorf("could not parse cookie %v", err)
	}

	ctx, err = checkAuthCookie(h.srv, cookie, ctx, nil)
	if err != nil {
		return fmt.Errorf("could not validate cookie %v", err)
	}

	sess, err := GetSessionCtx(ctx)
	if err != nil {
		return fmt.Errorf("could not get session %v", err)
	}

	err = h.srv.store.Delete(sess)
	if err != nil {
		return fmt.Errorf("could not delete session %v", err)
	}

	return nil
}

func writeSessionCookie(hw HeaderWriter, s *Session, session, refresh string) {
	if hw == nil {
		// in case of logout
		return
	}

	addCookie(
		hw,
		CookieSessionToken, session,
		s.SessionTokenExpiry,
	)
	addCookie(
		hw,
		CookieRefreshToken, refresh,
		s.RefreshTokenExpiry,
	)
}

type HeaderWriter interface {
	Header() http.Header
}

func addCookie(
	resp HeaderWriter,
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
