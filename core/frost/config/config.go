package config

type Config struct {
	Server Server
}

type Server struct {
	Port    int      `yaml:"port" env:"PORT" default:"9966" help:"Port to listen on"`
	Origins []string `yaml:"origins" env:"ORIGINS" default:"*" help:"Allowed origins in CSV"`
}
