package app

import "github.com/ra341/glacier/frost/config"

type App struct {
	Conf *config.Service
}

func New() *App {
	conf := config.New()

	return &App{
		Conf: conf,
	}
}
