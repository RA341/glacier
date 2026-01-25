package config

type Config struct {
	Server Server
	Files  Files
}

type Server struct {
	Port    int      `yaml:"port" env:"PORT" default:"9966" help:"Port to listen on"`
	Origins []string `yaml:"origins" env:"ORIGINS" default:"*" help:"Allowed origins in CSV"`

	GlacierUrl string `yaml:"glacierUrl" env:"GLACIER_URL" default:"http://localhost:6699" help:"url of glacier server"`
}

type Files struct {
	ConfigDir string `yaml:"configDir" env:"CONFIG_DIR" default:"./config" help:"frost config directory"`
}
