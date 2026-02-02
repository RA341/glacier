package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

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
	subFS, err := fs.Sub(uiDir, "build")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading frontend directory: %w", err))
	}

	frost.NewTray(frost.WithServerBase(api.WithUIFS(subFS)))
}
