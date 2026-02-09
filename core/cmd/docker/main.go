package main

import (
	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
)

func init() {
	app.InitMeta(info.FlavourDocker)
}

func main() {
	app.NewServer(
		api.WithFromPath("./web"),
	)
}
