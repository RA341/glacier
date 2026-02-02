package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/ra341/glacier/internal/auth"
	"github.com/ra341/glacier/internal/indexer"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/search"
	sm "github.com/ra341/glacier/internal/services_manager"
	"github.com/ra341/glacier/internal/user"
	"github.com/ra341/glacier/shared/api"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*App

	UIFS    fs.FS
	Ctx     context.Context
	uiProxy http.Handler
}

func NewServer(opts ...ServerOpt) {
	var server Server
	ParseOpts(&server, opts...)

	server.App = NewApp()
	config := server.Conf.Get().Server

	router := http.NewServeMux()
	server.RegisterRoutes(router)

	finalMux := WithCors(router, config.AllowedOrigins)

	port := fmt.Sprintf(":%d", config.Port)
	log.Info().Str("port", port).Msg("Starting server...")

	srv := &http.Server{
		Addr: port,
		Handler: h2c.NewHandler(
			finalMux,
			&http2.Server{},
		),
	}

	go func() {
		var err error
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	<-server.Ctx.Done()

	fmt.Println("Context cancelled. Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error occurred while shutting down server: %v\n", err)
	}

	fmt.Println("Server gracefully stopped.")
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

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	apiRouter := http.NewServeMux()
	s.registerApiRoutes(apiRouter)
	api.WithSubRouter(
		mux,
		"/api",
		s.withLogger(apiRouter),
	)

	s.registerUI(mux)
}

func (s *Server) registerUI(mux *http.ServeMux) {
	if s.UIFS != nil {
		log.Info().Msg("using fs UI")
		mux.Handle("/", api.NewSpaHandler(s.UIFS))
		return
	}

	if s.uiProxy != nil {
		log.Info().Msg("using proxy UI")
		mux.Handle("/", s.uiProxy)
		return
	}

	log.Info().Msg("using default UI")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("No UI opt was set when building"))
	})
}

func (s *Server) registerApiRoutes(mux *http.ServeMux) {
	serverRouter := http.NewServeMux()

	protectedRouter := http.NewServeMux()
	s.registerProtectedRoutes(protectedRouter)
	api.WithSubRouter(
		serverRouter,
		"/protected",
		s.withAuth(protectedRouter),
	)

	publicRouter := http.NewServeMux()
	s.registerPublicRoutes(publicRouter)
	api.WithSubRouter(
		serverRouter,
		"/public",
		publicRouter,
	)

	api.WithSubRouter(mux, "/server", serverRouter)
}

func (s *Server) registerProtectedRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("protected fuck off"))
	})

	mux.Handle(search.NewHandler(s.Search))
	mux.Handle(indexer.NewHandler(s.Indexer))

	mux.Handle(library.NewHandler(s.Library))
	api.WithSubRouter(mux,
		"/library/download",
		library.NewHandlerHttp(s.Library),
	)

	mux.Handle(user.NewHandler(s.User))

	adminMiddleware := NewMiddleware(user.AdminMiddleware)
	mux.Handle(adminMiddleware(sm.NewHandler(s.ConfigManager)))
}

type NewHandler func(string, http.Handler) (string, http.Handler)

type Middleware func(http.Handler) http.Handler

func NewMiddleware(m Middleware) NewHandler {
	return func(path string, h http.Handler) (string, http.Handler) {
		return path, m(h)
	}
}

func (s *Server) registerPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("unprotected fuck off"))
	})

	mux.Handle(auth.NewHandler(s.Session))
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// middleware stuff

func (s *Server) withAuth(mux *http.ServeMux) http.Handler {
	if s.Conf.Get().Auth.Disable {
		log.Warn().Msg("CAUTION: AUTHENTICATION IS DISABLED")
		return mux
	}

	return auth.NewMiddleware(
		s.Session,
		mux,
	)
}

func (s *Server) withLogger(protectedRouter *http.ServeMux) http.Handler {
	if !s.Conf.Get().Logger.HTTPLogger {
		return protectedRouter
	}

	log.Info().Msg("using http route logger")
	return api.LoggingMiddleware(protectedRouter)
}
