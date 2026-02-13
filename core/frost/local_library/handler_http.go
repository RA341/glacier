package download

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/ra341/glacier/pkg/fileutil"
	"github.com/ra341/glacier/pkg/ws"

	"github.com/dustin/go-humanize"
	"github.com/gorilla/websocket"
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
	subMux.HandleFunc("/speed", h.LoadDownloadSpeed)

	return subMux
}

func (h *HandlerHttp) LoadDownloadSpeed(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Debug().Msg("download speed tracker stopped")
	}()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer fileutil.Close(conn)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	last := atomic.LoadUint64(h.srv.downloader.DownloadTotalBytes)

	stopChan := make(chan struct{})
	go func() {
		// If the user closes the tab, this goroutine should stop
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				close(stopChan)
				return
			}
		}
	}()

	for {
		select {
		case <-stopChan:
			return
		case <-r.Context().Done():
			return
		case <-ticker.C:
			curr := atomic.LoadUint64(h.srv.downloader.DownloadTotalBytes)

			// Calculate delta (bytes since last second)
			dBytes := curr - last
			last = curr

			payload := struct {
				Speed string `json:"speed"`
				Raw   uint64 `json:"bytes_per_sec"`
			}{
				Speed: fmt.Sprintf("%s/s", humanize.Bytes(dBytes)),
				Raw:   dBytes,
			}

			if err := conn.WriteJSON(payload); err != nil {
				return
			}
		}
	}
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
