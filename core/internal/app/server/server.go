package server

import (
	"context"
	"errors"
	"fmt"

	"net/http"

	"time"

	"github.com/ra341/glacier/internal/app"
	downloadManager "github.com/ra341/glacier/internal/downloader/manager"
	indexerManager "github.com/ra341/glacier/internal/indexer/manager"
	"github.com/ra341/glacier/internal/library"
	metadataManager "github.com/ra341/glacier/internal/metadata/manager"
	"github.com/ra341/glacier/internal/search"
	"github.com/ra341/glacier/shared/api"

	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type Server struct {
	*app.App

	uiDir string
}

func NewServer(uiDir string) {
	server := &Server{
		App:   app.NewApp(),
		uiDir: uiDir,
	}
	config := server.App.Conf.Get().Server

	router := http.NewServeMux()
	server.RegisterRoutes(router)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:      config.AllowedOrigins,
		AllowPrivateNetwork: true,
		AllowedMethods:      connectcors.AllowedMethods(),
		AllowedHeaders:      connectcors.AllowedHeaders(),
		ExposedHeaders:      connectcors.ExposedHeaders(),
	})

	finalMux := corsConfig.Handler(router)

	port := fmt.Sprintf(":%d", config.Port)
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
	apiRouter := http.NewServeMux()
	s.registerRoutes(apiRouter)
	api.WithSubRouter(mux, "/api/server", apiRouter)

	s.registerUI(mux)
}

func (s *Server) registerUI(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(s.uiDir)))
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write([]byte("fuck off"))
	})

	mux.Handle(search.NewHandler(s.Search))
	mux.Handle(library.NewHandler(s.Library))
	api.WithSubRouter(mux,
		"/library/download",
		library.NewHandlerHttp(s.Library),
	)

	mux.Handle(downloadManager.NewHandler(s.DownloadClientManager))
	mux.Handle(metadataManager.NewHandler(s.MetadataManager))
	mux.Handle(indexerManager.NewHandler(s.IndexerManager))
}
