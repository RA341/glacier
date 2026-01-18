package server_config

import (
	"github.com/ra341/glacier/internal/downloader"
	"github.com/ra341/glacier/internal/indexer"
	"github.com/ra341/glacier/internal/metadata"
)

type Config struct {
	Server Server `yaml:"server"`
	Logger Logger `yaml:"logger"`
	Files  Files  `yaml:"files"`

	Indexer  indexer.Config    `yaml:"indexer"`
	Download downloader.Config `yaml:"downloader"`
	Metadata metadata.Config   `yaml:"metadata"`
}

type Files struct {
	// holds db and the yml file
	ConfigDir string `yaml:"config" env:"CONFIG_DIR" default:"./config"`
	GameDir   string `yaml:"game" env:"GAME_DIR" default:"./gamestop"`
}

type Server struct {
	Port           int      `yaml:"port" default:"6699" env:"SERVER_PORT"`
	AllowedOrigins []string `yaml:"allowedOrigins" default:"*"`
}

type Logger struct {
	Verbose bool `yaml:"verbose" default:"false"`
	Level   int  `yaml:"level" default:"info"`
}
