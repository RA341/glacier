package config

import (
	"fmt"

	"github.com/ra341/glacier/pkg/argos"
)

type Service struct {
	conf Config
}

func New() *Service {
	s := &Service{}
	s.Init()
	return s
}

const EnvPrefix = "FROST"

func (s *Service) Init() {
	var conf Config

	// todo yaml

	defaultPrefixer := argos.WithPrefixer(EnvPrefix)
	rnFn := argos.FieldProcessorTag(defaultPrefixer)
	argos.LoadStruct(&conf, rnFn)

	//pathsToResolve := []*string{
	//	&conf.Download.IncompletePath,
	//	&conf.Glacier.ConfigDir,
	//	&conf.Library.GameDir,
	//}
	//resolvePaths(pathsToResolve)

	printConfig(defaultPrefixer, &conf)

	s.conf = conf
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
