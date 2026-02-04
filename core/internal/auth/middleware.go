package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/shared/api"
	"github.com/rs/zerolog/log"
)

const CookieRefreshToken = "refresh_token"
const CookieSessionToken = "session_token"

const FrostHeader = "is-frost"
const FrostReqHeader = "req-frost"

func NewMiddleware(srv *Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context
		var err error

		if r.Header.Get(FrostReqHeader) == "true" {
			ctx, err = checkFrostHeaders(srv, r.Header, r.Context(), w)
		} else {
			ctx, err = checkAuthCookie(srv, r.Cookies(), r.Context(), w)
		}

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

func checkFrostHeaders(srv *Service, headers http.Header, ctx context.Context, hw http.ResponseWriter) (retCtx context.Context, err error) {
	session := headers.Get(HeaderFrostSessionToken)
	refresh := headers.Get(HeaderFrostRefreshToken)

	return verifySession(ctx, srv, session, refresh, hw)
}

func checkAuthCookie(srv *Service, cookies []*http.Cookie, ctx context.Context, headerW HeaderWriter) (context.Context, error) {
	session := ""
	refresh := ""
	for _, s := range cookies {
		if s.Name == CookieSessionToken {
			session = s.Value
		}
		if s.Name == CookieRefreshToken {
			refresh = s.Value
		}
	}

	return verifySession(ctx, srv, session, refresh, headerW)
}

func verifySession(ctx context.Context, srv *Service, sessionToken string, refreshToken string, hw HeaderWriter) (context.Context, error) {
	if sessionToken != "" {
		session, err := srv.VerifySession(sessionToken)
		if err == nil {
			return injectSession(ctx, &session), nil
		}

		log.Debug().Err(err).Msg("session is invalid, attempting refresh")
	} else {
		//log.Warn().Msg("empty session token")
	}

	// session failed/missing, refresh the session
	if refreshToken != "" {
		newSession, newSessionToken, newRefreshToken, err := srv.RefreshSession(refreshToken)
		if err == nil {
			log.Debug().
				Str("type", newSession.SessionType.String()).
				Msg("Session refreshed")

			if newSession.SessionType == Web {
				writeWebCookie(
					hw,
					&newSession,
					newSessionToken,
					newRefreshToken,
				)
			}

			if newSession.SessionType == Frost {
				WriteFrostHeader(
					hw,
					newSessionToken,
					newRefreshToken,
				)
			}

			return injectSession(ctx, &newSession), nil
		}

		log.Error().Err(err).Msg("could not refresh session")
	} else {
		//log.Warn().Msg("empty refresh token")
	}

	return ctx, ErrInvalidSession
}

var ErrInvalidSession = errors.New("authentication failed no valid session or refresh token")

const CtxKeySession = "session-inf"

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
