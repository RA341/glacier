package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ra341/glacier/pkg/argos"
	"github.com/rs/zerolog/log"
)

type Service struct {
	conf Config
}

func New(printConf bool) *Service {
	s := &Service{}
	s.init(printConf)
	return s
}

const EnvPrefix = "FROST"

func (s *Service) init(printConf bool) {
	var conf Config

	// todo yaml

	defaultPrefixer := argos.WithPrefixer(EnvPrefix)
	rnFn := argos.FieldProcessorTag(defaultPrefixer)
	argos.LoadStruct(&conf, rnFn)

	pathsToResolve := []*string{
		&conf.Files.ConfigDir,
		&conf.Files.LogsDir,
	}
	err := resolvePaths(pathsToResolve)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to resolve paths")
		return
	}

	if printConf {
		printConfig(defaultPrefixer, &conf)
	}

	s.conf = conf
}

func resolvePaths(pathsToResolve []*string) error {
	for _, p := range pathsToResolve {
		absPath, err := filepath.Abs(*p)
		if err != nil {
			return fmt.Errorf("failed to get abs path for %s: %w", *p, err)
		}
		*p = absPath

		if err = os.MkdirAll(absPath, 0777); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Get() *Config {
	return &s.conf
}

func printConfig(defaultPrefixer argos.Prefixer, conf *Config) {
	envDisplay := argos.WithUnderLine("Env:")
	envTag := argos.FieldPrintConfig{
		TagName: "env",
		PrintConfig: func(TagName string, val *argos.FieldVal) {
			v, ok := val.Tags[TagName]
			if ok {
				val.Tags[TagName] = envDisplay + " " +
					argos.Colorize(defaultPrefixer(v), argos.ColorCyan)
			}
		},
	}
	helpTag := argos.FieldPrintConfig{
		TagName: "help",
		PrintConfig: func(TagName string, val *argos.FieldVal) {
			v, ok := val.Tags[TagName]
			if ok {
				val.Tags[TagName] = argos.Colorize(v, argos.ColorYellow)
			}
		},
	}

	ms := argos.Colorize("To modify config, set the respective", argos.ColorMagenta+argos.ColorBold)
	footer := fmt.Sprintf("%s %s", ms, envDisplay)

	argos.PrintInfo(
		conf,
		footer,
		helpTag, envTag,
	)
}
