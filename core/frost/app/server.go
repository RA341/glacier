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
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/pkg/logger"
	"github.com/ra341/glacier/shared/api"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func InitMeta(flavour info.FlavourType) {
	info.SetFlavour(flavour)
	info.PrintInfo()
	logger.InitDefault()
}

func init() {
	InitMeta(info.FlavourFrost)
}

type Server struct {
	app *App

	uiDir string
}

func NewServer(uiDir string) {
	server := &Server{
		app:   New(),
		uiDir: uiDir,
	}
	conf := server.app.Conf.Get().Server

	router := http.NewServeMux()
	server.RegisterRoutes(router)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:      conf.Origins,
		AllowPrivateNetwork: true,
		AllowedMethods:      connectcors.AllowedMethods(),
		AllowedHeaders:      connectcors.AllowedHeaders(),
		ExposedHeaders:      connectcors.ExposedHeaders(),
	})

	finalMux := corsConfig.Handler(router)

	port := fmt.Sprintf(":%d", conf.Port)
	log.Info().Str("port", port).Msg("Starting server...")

	srv := &http.Server{
		Addr: port,
		Handler: h2c.NewHandler(
			finalMux,
			&http2.Server{},
		),
	}

	serverCtx := context.Background() // todo load from config

	go func() {
		var err error
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Error starting server")
		}
	}()

	<-serverCtx.Done()

	fmt.Println("Context cancelled. Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error occurred while shutting down server: %v\n", err)
	}

	fmt.Println("Server gracefully stopped.")
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	s.registerRoutes(mux)

	s.registerUI(mux)
}

func (s *Server) registerUI(mux *http.ServeMux) {
	if s.uiDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(s.uiDir)))
	}

	apiMux := http.NewServeMux()
	s.registerRoutes(apiMux)
	api.WithSubRouter(mux, "/api/frost", apiMux)

	glacierProxy := http.NewServeMux()
	s.registerGlacierProxy(glacierProxy)
	mux.Handle("/api/server/", glacierProxy)
}

func (s *Server) registerGlacierProxy(mux *http.ServeMux) {
	target, err := url.Parse("http://localhost:6699")
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing url")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
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
		return nil
	}

	mux.Handle("/", proxy)
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("Hello from frost"))
	})

	mux.Handle(ll.NewHandler(s.app.LocalLibrarySrv))
}
