package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/shared/api"
)

var ErrNotAuthorized = errors.New("insufficient permissions")

var OmniMiddleware = RequireRole(Omnissiah)

func RequireRole(minRole Role) func(http.Handler) http.Handler {
	return api.NewMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := GetUserCtx(r.Context())
			if err != nil {
				api.WriteErr(
					w,
					http.StatusUnauthorized,
					connect.CodeUnauthenticated,
					"user info not found",
				)
				return
			}

			if user.Role > minRole {
				api.WriteErr(
					w,
					http.StatusForbidden,
					connect.CodePermissionDenied,
					"insufficient permission",
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	})
}

const CtxKeyUser = "user"

func GetUserCtx(ctx context.Context) (*User, error) {
	val := ctx.Value(CtxKeyUser)
	if val == nil {
		return nil, fmt.Errorf("user not found in context")
	}

	u, ok := val.(*User)
	if !ok {
		return nil, fmt.Errorf("could not convert value to user struct")
	}
	return u, nil
}
