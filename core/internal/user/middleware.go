package user

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/ra341/glacier/generated/service_config/v1/v1connect"
	"github.com/ra341/glacier/shared/api"
)

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == v1connect.ServiceConfigServiceGetActiveServiceProcedure {
			next.ServeHTTP(w, r)
			return
		}

		user, err := GetUserCtx(r.Context())
		if err != nil {
			api.WriteErr(w,
				http.StatusBadRequest,
				connect.CodeUnauthenticated,
				"user info not found in context",
			)
			return
		}

		if user.Role > Magos {
			api.WriteErr(w,
				http.StatusBadRequest,
				connect.CodeUnauthenticated,
				"insufficient permission to access this resource",
			)
			return
		}

		next.ServeHTTP(w, r)
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
