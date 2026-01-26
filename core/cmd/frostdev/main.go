package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	frost "github.com/ra341/glacier/frost/app"
)

//go:embed all:build
var uiDir embed.FS

func main() {
	subFS, err := fs.Sub(uiDir, "build")
	if err != nil {
		log.Fatal(fmt.Errorf("error loading frontend directory: %w", err))
	}

	frost.NewTray(subFS)
}
