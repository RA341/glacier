package config

import "github.com/ra341/glacier/frost/local_library/download"

type Config struct {
	Desktop    Desktop
	Server     Server
	Files      Files
	Logger     Logger
	Downloader download.Config
}

type Desktop struct {
	StartSilent bool `yaml:"startSilent" env:"START_SILENT" default:"false" help:"Disable UI launch on startup"`
}

type Server struct {
	Port           int      `yaml:"port" env:"PORT" default:"9966" help:"Port to listen on"`
	AllowedOrigins []string `yaml:"origins" env:"ORIGINS" default:"*" help:"Allowed origins in CSV"`

	GlacierUrl string `yaml:"glacierUrl" env:"GLACIER_URL" default:"http://localhost:6699" help:"url of glacier server"`
}

type Files struct {
	ConfigDir string `yaml:"configDir" env:"CONFIG_DIR" default:"./config" help:"frost config directory" folder:""`
	LogsDir   string `yaml:"logsDir" env:"LOG_DIR" default:"./logs" help:"log output" folder:""`
}

type Logger struct {
	Verbose    bool   `yaml:"verbose" default:"false" env:"LOGGER_VERBOSE" help:"add more info"`
	Level      string `yaml:"level" default:"info" env:"LOGGER_LEVEL" help:"disabled|debug|info|warn|error|fatal"`
	HTTPLogger bool   `yaml:"HTTPLogger" default:"false" env:"LOGGER_HTTP" help:"log api routes"`
}
