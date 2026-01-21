package main

import (
	"embed"

	frost "github.com/ra341/glacier/frost/app"
)

//go:embed build
var uiDir embed.FS

func main() {
	frost.NewServer("")
}
