package config

type Config struct {
	Server     Server
	Files      Files
	Logger     Logger
	Downloader Downloader
}

type Server struct {
	Port           int      `yaml:"port" env:"PORT" default:"9966" help:"Port to listen on"`
	AllowedOrigins []string `yaml:"origins" env:"ORIGINS" default:"*" help:"Allowed origins in CSV"`

	GlacierUrl string `yaml:"glacierUrl" env:"GLACIER_URL" default:"http://localhost:6699" help:"url of glacier server"`
}

type Downloader struct {
	MaxConcurrentFiles int `yaml:"maxConcurrentFiles" env:"MAX_FILES" default:"50" help:"Maximum number of concurrent files"`
	MaxFileChunks      int `yaml:"maxFileChunks" env:"MAX_CHUNKS" default:"100" help:"Maximum number of chunks in a file to process"`
}

type Files struct {
	ConfigDir string `yaml:"configDir" env:"CONFIG_DIR" default:"./config" help:"frost config directory"`
}

type Logger struct {
	Verbose    bool   `yaml:"verbose" default:"false" env:"LOGGER_VERBOSE" help:"add more info"`
	Level      string `yaml:"level" default:"info" env:"LOGGER_LEVEL" help:"disabled|debug|info|warn|error|fatal"`
	HTTPLogger bool   `yaml:"HTTPLogger" default:"false" env:"LOGGER_HTTP" help:"log api routes"`
}
