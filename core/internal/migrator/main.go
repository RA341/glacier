package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/ra341/glacier/internal/auth"
	"github.com/ra341/glacier/internal/library"
	"github.com/ra341/glacier/internal/services_manager"
	"github.com/ra341/glacier/internal/user"
)

func main() {
	stmts, err := gormschema.
		New("sqlite").
		Load(
			&library.Game{},
			&library.FolderManifest{},
			&services_manager.ServiceConfig{},
			&user.User{},
			&auth.Session{},
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
