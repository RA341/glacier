package library

import (
	"net/http"
	"strconv"
)

type HandlerHttp struct {
	srv *Service
}

func NewHandlerHttp(srv *Service) http.Handler {
	h := &HandlerHttp{
		srv: srv,
	}

	subMux := http.NewServeMux()
	subMux.HandleFunc("GET /meta/{game}", h.getMetadata)
	subMux.HandleFunc("GET /load/{game}/{file}", h.downloadFile)

	return subMux
}

func (h *HandlerHttp) getMetadata(w http.ResponseWriter, r *http.Request) {
	gameId := r.PathValue("game")
	if gameId == "" {
		http.Error(w, "no game id specified", http.StatusBadRequest)
		return
	}
	gid, err := strconv.Atoi(gameId)
	if err != nil {
		http.Error(w, "could not convert id to uint", http.StatusBadRequest)
		return
	}

	err = h.srv.GetDownloadManifest(r.Context(), gid, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerHttp) downloadFile(w http.ResponseWriter, r *http.Request) {
	gameIdStr := r.PathValue("game")
	if gameIdStr == "" {
		http.Error(w, "no game id specified", http.StatusBadRequest)
		return
	}
	gameId, err := strconv.Atoi(gameIdStr)
	if err != nil {
		http.Error(w, "could not convert id to uint", http.StatusBadRequest)
		return
	}

	file := r.PathValue("file")
	if file == "" {
		http.Error(w, "no filename specified", http.StatusBadRequest)
		return
	}

	fullFilePath, err := h.srv.GetDownloadPath(r.Context(), gameId, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, fullFilePath)
}
