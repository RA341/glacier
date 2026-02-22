package api

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

// ServerBase defines a common server struct to support common server opts
type ServerBase struct {
	Ctx             context.Context
	UIHandler       http.Handler
	RestartHandler  func()
	ShutdownHandler func()
}

func (s *ServerBase) RegisterUI(mux *http.ServeMux, defaultUIHandler http.HandlerFunc) {
	if s.UIHandler != nil {
		mux.Handle("/", s.UIHandler)
		return
	}

	//if s.UIFS != nil {
	//	log.Info().Msg("using FS SPA UI")
	//	mux.Handle("/", NewSpaHandler(s.UIFS))
	//	return
	//}
	//
	//if s.UIProxy != nil {
	//	log.Info().Msg("using proxy UI")
	//	mux.Handle("/", s.UIProxy)
	//	return
	//}

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

func WithCtx(ctx context.Context) ServerOpt {
	return func(o *ServerBase) {
		o.Ctx = ctx
	}
}

func WithUIProxy(target string) ServerOpt {
	u, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(u)

	hand := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Host = u.Host
		proxy.ServeHTTP(w, r)
	})

	return WithUIHandler(hand)
}

func WithUIFS(uifs fs.FS, stripPrefix string) ServerOpt {
	subFS, err := fs.Sub(uifs, stripPrefix)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading frontend directory")
	}

	hand := NewSpaHandler(subFS)
	return WithUIHandler(hand)
}

func WithFromPath(path string) ServerOpt {
	root, err := os.OpenRoot(path)
	if err != nil {
		log.Fatal().Err(err).Str("path", path).Msg("could not load UI from file")
	}

	hand := NewSpaHandler(root.FS())

	return WithUIHandler(hand)
}

func WithUIHandler(ui http.Handler) ServerOpt {
	return func(o *ServerBase) {
		o.UIHandler = ui
	}
}

func WithRestartHandler(restart func()) ServerOpt {
	return func(o *ServerBase) {
		o.RestartHandler = restart
	}
}

func WithShutDownHandler(shutdown func()) ServerOpt {
	return func(o *ServerBase) {
		o.ShutdownHandler = shutdown
	}
}
