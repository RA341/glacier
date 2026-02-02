package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

// ServerBase defines a common server struct to support common server opts
type ServerBase struct {
	UIFS    fs.FS
	Ctx     context.Context
	UIProxy http.Handler
}

func (s *ServerBase) RegisterUI(mux *http.ServeMux, defaultUIHandler http.HandlerFunc) {
	if s.UIFS != nil {
		log.Info().Msg("using FS SPA UI")
		mux.Handle("/", NewSpaHandler(s.UIFS))
		return
	}

	if s.UIProxy != nil {
		log.Info().Msg("using proxy UI")
		mux.Handle("/", s.UIProxy)
		return
	}

	log.Info().Msg("using default UI")
	mux.HandleFunc("/", defaultUIHandler)
}

func LoadUIFromDir(path string) (fs.FS, error) {
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open UI path: %s : %w", path, err)
	}
	return root.FS(), nil
}

type ServerOpt func(o *ServerBase)

func ParseOpts(srv *ServerBase, opts ...ServerOpt) {
	for _, opt := range opts {
		opt(srv)
	}

	if srv.Ctx == nil {
		srv.Ctx = context.Background()
	}
}

func WithServerBase(sb *ServerBase) ServerOpt {
	return func(o *ServerBase) {
		*o = *sb
	}
}

func WithUIProxy(ui http.Handler) ServerOpt {
	return func(o *ServerBase) {
		o.UIProxy = ui
	}
}

func WithCtx(ctx context.Context) ServerOpt {
	return func(o *ServerBase) {
		o.Ctx = ctx
	}
}

func WithUIFS(uiFs fs.FS) ServerOpt {
	return func(o *ServerBase) {
		o.UIFS = uiFs
	}
}
