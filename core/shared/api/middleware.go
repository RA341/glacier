package api

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"time"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type Mid func(next http.Handler) http.Handler

func NewMiddleware(handFn Mid) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handFn(next)
	}
}

func SkipIfPath(
	excludedPath []string,
	// this is the next handler in the chain
	// if the path is not equal this will be called otherwise bypassed
	nextHandler func(http.Handler) http.Handler,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		newNext := nextHandler(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Contains(excludedPath, r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			newNext.ServeHTTP(w, r)
		})
	}
}

func WithProxy(target string) http.Handler {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Host
		proxy.ServeHTTP(w, r)
	})
}

func WithCors(router http.Handler, origins []string) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:      origins,
		AllowPrivateNetwork: true,
		AllowedMethods:      connectcors.AllowedMethods(),
		AllowedHeaders:      connectcors.AllowedHeaders(),
		ExposedHeaders:      connectcors.ExposedHeaders(),
	}).Handler(router)
}

func WithHTTPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subLog := log.With().
			Str("url", r.URL.String()).
			Str("method", r.Method).
			Logger()

		if r.Header.Get("Connect-Protocol-Version") != "" {
			subLog.Debug().Str("url", r.URL.String()).
				Msg("connect rpc")
			next.ServeHTTP(w, r)
			return
		}

		// WebSocket request, don't
		// wrap the writer to avoid hijacking errs
		if r.Header.Get("Upgrade") == "websocket" {
			subLog.Debug().Str("url", r.URL.String()).
				Msg("websocket connection")
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()
		//subLog.Debug().Msg("Started request")
		next.ServeHTTP(w, r)

		subLog.Debug().
			//Int("status", wrapped.statusCode).
			Dur("elapsed", time.Since(start)).
			Msg("Completed request")
	})
}
