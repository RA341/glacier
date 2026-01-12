package app

import (
	"github.com/ra341/glacier/internal/database"
	libraryquery "github.com/ra341/glacier/internal/database/generated/queries/library"
	"github.com/ra341/glacier/internal/library"
)

type App struct {
	Conf    *Config
	Library *library.Service
}

func NewApp() *App {
	db := database.New(".", false)

	sd := libraryquery.Store[library.Game](db)
	libMan := library.NewApp(sd)

	return &App{
		Library: libMan,
	}
}
