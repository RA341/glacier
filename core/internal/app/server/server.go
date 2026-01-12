package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ra341/glacier/internal/app"

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

	router := http.NewServeMux()
	server.RegisterRoutes(router)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:      []string{"*"}, // todo load from config
		AllowPrivateNetwork: true,
		AllowedMethods:      connectcors.AllowedMethods(),
		AllowedHeaders:      connectcors.AllowedHeaders(),
		ExposedHeaders:      connectcors.ExposedHeaders(),
	})

	finalMux := corsConfig.Handler(router)

	portNum := 6699 // todo load from config
	port := fmt.Sprintf(":%d", portNum)
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
	s.registerUI(mux)
}

func (s *Server) registerUI(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir(s.uiDir)))
}
