package main

import (
	"embed"

	frost "github.com/ra341/glacier/frost/app"
	"github.com/ra341/glacier/internal/app"
	"github.com/ra341/glacier/internal/info"
	"github.com/ra341/glacier/shared/api"
)

//go:embed all:build
var uiDir embed.FS

func init() {
	app.InitMeta(info.FlavourFrost)
}

func main() {
	frost.NewDesktop(api.WithUIFS(uiDir, "build"))
}
