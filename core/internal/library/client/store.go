package client

import "github.com/ra341/glacier/internal/library"

type InstalledGame struct {
	library.Game

	InstalledPath string
}
