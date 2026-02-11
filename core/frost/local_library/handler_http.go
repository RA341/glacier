package download

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/ws"
	"github.com/rs/zerolog/log"
)

type HandlerHttp struct {
	srv *Service
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin should be more restrictive in production
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewHandlerHttp(srv *Service) http.Handler {
	h := &HandlerHttp{srv: srv}
	subMux := http.NewServeMux()
	subMux.HandleFunc("/running/{gameId}", h.HandlerProcessRunning)

	return subMux
}

func (h *HandlerHttp) HandlerProcessRunning(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.PathValue("gameId")
	if gameIdStr == "" {
		http.Error(w, "missing 'gameID' parameter", http.StatusBadRequest)
		return
	}
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		http.Error(w, "could not convert 'gameID' to int", http.StatusBadRequest)
		return
	}

	exe := r.URL.Query().Get("exe")
	if exe == "" {
		http.Error(w, "missing 'exe' parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn().Err(err).Str("exe", exe).Msg("Failed to upgrade connection")
		return
	}
	defer fileutil.Close(conn)

	err = h.srv.Running(r.Context(), gameId, exe)
	if err != nil {
		ws.WErr(conn, fmt.Errorf("error when process exited"))
		return
	}
}
