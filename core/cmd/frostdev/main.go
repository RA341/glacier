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

	go func() {
		frost.NewServer(frost.WithUIFS(subFS))
	}()

	err = frost.NewDesktop("ui/ui")
	if err != nil {
		log.Fatal("err starting ui", err)
	}
}
