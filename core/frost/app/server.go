package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	ll "github.com/ra341/glacier/frost/local_library"
	"github.com/ra341/glacier/internal/auth"
	"github.com/ra341/glacier/shared/api"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*App
	api.ServerBase
}

func NewServer(opts ...api.ServerOpt) {
	var server Server
	api.ParseOpts(&server.ServerBase, opts...)

	server.App = New()

	conf := server.Conf.Get().Server

	router := http.NewServeMux()
	server.RegisterRoutes(router)

	finalMux := api.WithCors(router, conf.AllowedOrigins)

	port := fmt.Sprintf(":%d", conf.Port)
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
	// todo send a return value to indicate if the tray should also be stopped or not
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	s.registerApiRoutes(mux)
	s.registerUI(mux)
}

func (s *Server) registerUI(mux *http.ServeMux) {
	s.RegisterUI(mux, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("No UI was set when building"))
	})
}

func (s *Server) registerApiRoutes(mux *http.ServeMux) {
	frostMux := http.NewServeMux()
	s.RegisterFrostRoutes(frostMux)
	api.WithSubRouter(
		mux,
		"/api/frost",
		frostMux,
	)

	glacierProxy := http.NewServeMux()
	s.registerGlacierProxy(glacierProxy)
	mux.Handle("/api/server/", glacierProxy)
}

func (s *Server) RegisterFrostRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("A pleasant fuck off from frost"))
	})

	mux.Handle(ll.NewHandler(s.LocalLibrarySrv))
}

func (s *Server) registerGlacierProxy(mux *http.ServeMux) {
	conf := s.App.Conf.Get()
	target, err := url.Parse(conf.Server.GlacierUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing url")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		req.Header.Set(auth.FrostHeader, "true")

		originalDirector(req)
		req.Host = target.Host
	}

	proxy.ModifyResponse = func(resp *http.Response) error {
		// server already has CORS middleware, remove the
		// headers coming from the backend to prevent multiple value errors.
		resp.Header.Del("Access-Control-Allow-Origin")
		resp.Header.Del("Access-Control-Allow-Credentials")
		resp.Header.Del("Access-Control-Allow-Methods")
		resp.Header.Del("Access-Control-Allow-Headers")

		refresh := resp.Header.Get(auth.HeaderFrostRefreshToken)
		session := resp.Header.Get(auth.HeaderFrostSessionToken)
		if refresh != "" && session != "" {
			err := s.Secret.AddSession(session, refresh)
			if err != nil {
				// maybe return error IDK
				log.Warn().Err(err).Msg("Error adding session")
			}
		}

		return nil
	}

	mux.Handle("/", proxy)
}
