package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ra341/glacier/internal/auth"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/library/manifest"
	"github.com/ra341/glacier/internal/metadata/assets"
	"github.com/ra341/glacier/internal/services_manager"
	"github.com/ra341/glacier/internal/user"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	dialect := "sqlite"
	if len(os.Args) > 1 {
		dialect = os.Args[1]
	}

	stmts, err := gormschema.
		New(dialect).
		Load(
			&library.Game{},
			&manifest.FolderManifest{},
			&services_manager.ServiceConfig{},
			&user.User{},
			&auth.Session{},
			&assets.Asset{},
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
