package library

type Config struct {
	GameDir string `yaml:"game" env:"GAME_DIR" default:"./gamestop"`
}
