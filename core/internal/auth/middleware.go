package auth

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/shared/api"
	"github.com/rs/zerolog/log"
)

const CookieRefreshToken = "refresh_token"
const CookieSessionToken = "session_token"

func NewMiddleware(srv *Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := checkAuthCookie(srv, r.Cookies(), r.Context(), w)
		if err != nil {
			api.WriteErr(
				w,
				http.StatusUnauthorized,
				connect.CodeUnauthenticated,
				err.Error(),
			)
			return
		}

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

const CtxKeySession = "session-inf"

func checkAuthCookie(srv *Service, cookies []*http.Cookie, ctx context.Context, headerW HeaderWriter) (context.Context, error) {
	sessionToken := ""
	refreshToken := ""
	for _, s := range cookies {
		if s.Name == CookieSessionToken {
			sessionToken = s.Value
		}
		if s.Name == CookieRefreshToken {
			refreshToken = s.Value
		}
	}

	if sessionToken != "" {
		session, err := srv.VerifySession(sessionToken)
		if err == nil {
			return injectSession(ctx, &session), nil
		}
		log.Debug().Err(err).Msg("session is invalid, attempting refresh")
	}

	// session failed/missing, refresh the session
	if refreshToken != "" {
		newSession, newSessionToken, newRefreshToken, err := srv.RefreshSession(refreshToken)
		if err == nil {
			writeSessionCookie(
				headerW,
				&newSession,
				newSessionToken,
				newRefreshToken,
			)

			return injectSession(ctx, &newSession), nil
		}
		log.Warn().Err(err).Msg("refresh attempt failed")
	}

	return ctx, fmt.Errorf("authentication failed no valid session or refresh token")
}

func injectSession(ctx context.Context, s *Session) context.Context {
	ctx = context.WithValue(ctx, CtxKeySession, s)
	ctx = context.WithValue(ctx, user.CtxKeyUser, &s.User)
	return ctx
}

func GetSessionCtx(ctx context.Context) (*Session, error) {
	val := ctx.Value(CtxKeySession)
	if val == nil {
		return nil, fmt.Errorf("session not found in context")
	}

	u, ok := val.(*Session)
	if !ok {
		return nil, fmt.Errorf("could not convert value to session struct")
	}
	return u, nil
}
