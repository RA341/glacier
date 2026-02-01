package auth

import (
	"context"
	"net/http"
)

const CookieRefreshToken = "refresh_token"
const CookieSessionToken = "session_token"

const CtxKeyUser = "user"

func NewMiddleware(srv *Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		found := false

		for _, s := range r.Cookies() {
			if s.Name == CookieSessionToken {
				session, err := srv.VerifySession(s.Value)
				if err != nil {
					http.Error(w, "Invalid session cookie", http.StatusUnauthorized)
					return
				}

				value := context.WithValue(
					r.Context(),
					CtxKeyUser,
					session.User,
				)

				r.WithContext(value)
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "session cookie not found", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
