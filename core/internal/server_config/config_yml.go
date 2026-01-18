package server_config

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ra341/glacier/pkg/fileutil"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

type ConfigYml struct {
	path string
	rw   *sync.RWMutex
}

func NewYml() ConfigYml {
	yml := ConfigYml{
		rw: new(sync.RWMutex),
	}
	yml.loadPath()
	return yml
}

func (cy *ConfigYml) loadPath() {
	defer func() {
		if cy.path == "" {
			log.Fatal().Msg("loadPath did not load a path lol")
		}

		log.Info().Str("path", cy.path).Msg("using config path")
	}()

	if loadPath := os.Getenv(GlacierYmlPathEnv); loadPath != "" {
		if !strings.HasSuffix(loadPath, ".yml") {
			log.Fatal().Str("path", loadPath).Msg("custom config file path must end with .yml")
		}

		abs, err := filepath.Abs(loadPath)
		if err != nil {
			log.Warn().Err(err).Str("path", loadPath).Msg("can't get absolute path for Yml path")
		}

		cy.path = abs
		return
	}

	configPath, err := os.Executable()
	if err != nil {
		log.Fatal().Err(err).Msg("can't get executable path")
	}
	configPath = filepath.Join(filepath.Dir(configPath), GlacierYml)

	cy.path = configPath
}

func (cy *ConfigYml) writeAndLoad(conf *Config) error {
	cy.rw.Lock()
	defer cy.rw.Unlock()

	err := cy.writeYml(conf)
	if err != nil {
		return err
	}

	return cy.writeYml(conf)
}

func (cy *ConfigYml) writeYml(conf *Config) error {
	contents, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return os.WriteFile(cy.path, contents, os.ModePerm)
}

func (cy *ConfigYml) loadYml(conf *Config) error {
	file, err := os.OpenFile(cy.path, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer fileutil.Close(file)

	contents, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(contents, conf)
}

// converts any path to abs path and creates it
func resolvePaths(pathsToResolve []*string) {
	for _, p := range pathsToResolve {
		absPath, err := filepath.Abs(*p)
		if err != nil {
			log.Fatal().Err(err).Str("path", *p).Msg("can't resolve path")
		}
		*p = absPath

		if err = os.MkdirAll(absPath, 0777); err != nil {
			log.Fatal().Err(err).Str("path", *p).Msg("can't create path")
		}
	}
}
