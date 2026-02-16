package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	ll "github.com/ra341/glacier/frost/local_library"
	"github.com/ra341/glacier/internal/metadata/assets"
)

func main() {
	stmts, err := gormschema.
		New("sqlite").
		Load(
			&ll.LocalGame{},
			&assets.Asset{},
		)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
