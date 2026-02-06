package config

import (
	"fmt"
	"sync/atomic"

	"github.com/ra341/glacier/pkg/argos"
	"github.com/rs/zerolog/log"
)

// todo this is likely not concurrent safe

type Service struct {
	cy   ConfigYml
	conf atomic.Pointer[Config]
}

func New() *Service {
	s := &Service{}
	s.Init()
	return s
}

const GlacierYml = "glacier.yml"
const GlacierYmlPathEnv = "GLACIER_CONFIG_YML_PATH"

func (s *Service) Init() {
	s.cy = NewYml()
	err := s.cy.backupCurrent()
	if err != nil {
		log.Fatal().Err(err).Msg("could not backup current config")
	}

	var conf Config
	err = s.cy.loadYml(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config file")
	}

	defaultPrefixer := DefaultPrefixer()
	rnFn := argos.FieldProcessorTag(defaultPrefixer)
	argos.LoadStruct(&conf, rnFn)

	pathsToResolve := []*string{
		&conf.Download.IncompletePath,
		&conf.Glacier.ConfigDir,
		&conf.Library.GameDir,
	}
	resolvePaths(pathsToResolve)

	printConfig(defaultPrefixer, &conf)

	err = s.storeAndLoad(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't init config file")
	}
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
	// todo hide
	//redactTag := argos.FieldPrintConfig{
	//	TagName: "hide",
	//	PrintConfig: func(TagName string, val *argos.FieldVal) {
	//		_, ok := val.Tags[TagName]
	//		if ok {
	//			val.Value = argos.Colorize("REDACTED", argos.ColorRed)
	//		}
	//	},
	//}
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

func (s *Service) Get() *Config {
	return s.conf.Load()
}

func (s *Service) storeAndLoad(loadCopy *Config) error {
	err := s.cy.writeAndLoad(loadCopy)
	if err != nil {
		return err
	}
	s.conf.Store(loadCopy)
	return nil
}

//func (s *Service) Set(src *Config) error {
//	newDst := s.loadCopy()
//
//	if err := mergo.Merge(&newDst, src); err != nil {
//		return err
//	}
//
//	s.conf.Store(&newDst)
//
//	return nil
//}

func (s *Service) loadCopy() Config {
	return *s.conf.Load()
}
