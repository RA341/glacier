package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/ra341/glacier/internal/config/config_manager"
	"github.com/ra341/glacier/internal/library"
)

func main() {
	stmts, err := gormschema.
		New("sqlite").
		Load(
			&library.Game{},
			&library.FolderManifest{},
			&config_manager.ServiceConfig{},
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
