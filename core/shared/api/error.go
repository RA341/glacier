package api

import (
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
)

// WriteErr writes connect compatible errors in JSON form
//
// will be of form {code: ..., message: ..}
func WriteErr(
	w http.ResponseWriter,
	httpCode int,
	connectCode connect.Code,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	m := fmt.Sprintf(
		`{"code":"%s","message":"%s"}`,
		connectCode.String(),
		message,
	)

	_, err := fmt.Fprint(w, m)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to write response")
	}
}
