package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// SpaHandler implements the http.Handler interface and serves a single-page
// application. If a requested file is not found, it serves the 'index.html'
// file, allowing client-side routing to take over.
type SpaHandler struct {
	staticFS   fs.FS
	fileServer http.Handler
}

// NewSpaHandler creates a new handler for serving a single-page application.
// see SpaHandler for info
func NewSpaHandler(staticFS fs.FS) http.Handler {
	return &SpaHandler{
		staticFS:   staticFS,
		fileServer: http.FileServer(http.FS(staticFS)),
	}
}

func (h *SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqPath := path.Clean(r.URL.Path)
	fsPath := strings.TrimPrefix(reqPath, "/")

	if fsPath == "" {
		http.ServeFileFS(w, r, h.staticFS, "index.html")
		return
	}

	info, err := fs.Stat(h.staticFS, fsPath)
	if err != nil {
		// File not found let SvelteKit handle routing
		http.ServeFileFS(w, r, h.staticFS, "index.html")
		return
	}

	// If it's a directory, serve index.html instead of letting
	// the file server redirect or show a listing
	if info.IsDir() {
		http.ServeFileFS(w, r, h.staticFS, "index.html")
		return
	}

	h.fileServer.ServeHTTP(w, r)
}
