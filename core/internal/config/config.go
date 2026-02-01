package config

import (
	"github.com/ra341/glacier/internal/downloader"
	"github.com/ra341/glacier/internal/library"
)

type Config struct {
	Glacier  Glacier           `yaml:"glacier"`
	Server   Server            `yaml:"server"`
	Logger   Logger            `yaml:"logger"`
	Auth     Auth              `yaml:"auth"`
	Library  library.Config    `yaml:"library"`
	Download downloader.Config `yaml:"downloader"`
}

type Glacier struct {
	// holds db and the yml file
	ConfigDir string `yaml:"config" env:"CONFIG_DIR" default:"./config" help:"path to the config dir"`
}

type Auth struct {
	MaxConcurrentSessions int  `yaml:"maxConcurrentSessions" env:"MAX_SESSIONS" default:"8" help:"maximum number of logged in sessions per user"`
	Disable               bool `yaml:"disable" env:"AUTH_DISABLE" default:"false" help:"disable authentication"`
}

type Server struct {
	Port           int      `yaml:"port" default:"6699" env:"SERVER_PORT" help:"server port"`
	AllowedOrigins []string `yaml:"allowedOrigins" default:"*" env:"ALLOWED_ORIGINS" help:"allowed origins in CSV"`
}

type Logger struct {
	Verbose    bool   `yaml:"verbose" default:"false" env:"LOGGER_VERBOSE" help:"add more info"`
	Level      string `yaml:"level" default:"info" env:"LOGGER_LEVEL" help:"disabled|debug|info|warn|error|fatal"`
	HTTPLogger bool   `yaml:"HTTPLogger" default:"false" env:"LOGGER_HTTP" help:"log api routes"`
}
