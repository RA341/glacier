package api

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func LoggingMiddleware(next http.Handler) http.Handler {
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
