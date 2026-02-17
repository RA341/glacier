package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync/atomic"

	"dario.cat/mergo"
	"github.com/ra341/glacier/pkg/argos"

	"github.com/rs/zerolog/log"
)

type Service[T any] struct {
	prefixer argos.Prefixer
	cy       Yml[T]
	conf     atomic.Pointer[T]
}

func New[T any](
	prefixer argos.Prefixer,
	ymlPathEnv, ymlFileName string,
) *Service[T] {
	s := &Service[T]{
		prefixer: prefixer,
		cy:       Yml[T]{},
	}
	s.Init(ymlPathEnv, ymlFileName)
	return s
}

//const GlacierYml = "glacier.yml"
//const GlacierYmlPathEnv = "GLACIER_CONFIG_YML_PATH"

func (s *Service[T]) Init(
	ymlPathEnv string,
	ymlFileName string,
) {
	s.cy = NewYml[T](
		ymlPathEnv,
		ymlFileName,
	)
	err := s.cy.backupCurrent()
	if err != nil {
		log.Fatal().Err(err).Msg("could not backup current config")
	}

	var conf T
	err = s.cy.loadYml(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config file")
	}

	defaultPrefixer := s.prefixer
	rnFn := argos.FieldProcessorTag(defaultPrefixer)
	argos.LoadStruct(&conf, rnFn)

	printConfig(defaultPrefixer, &conf)

	err = s.storeAndLoad(&conf)
	if err != nil {
		log.Fatal().Err(err).Msg("can't init config file")
	}
}

func printConfig[T any](defaultPrefixer argos.Prefixer, conf *T) {
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

func (s *Service[T]) Get() *T {
	return s.conf.Load()
}

func (s *Service[T]) ListFiles(base string) ([]string, error) {
	abs, err := filepath.Abs(base)
	if err != nil {
		return nil, err
	}

	dirs, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	var files = make([]string, 0, len(dirs))
	for _, f := range dirs {
		if !f.IsDir() {
			continue
		}
		files = append(files, filepath.Join(base, f.Name()))
	}

	return files, nil
}

func (s *Service[T]) storeAndLoad(loadCopy *T) error {
	err := s.cy.writeAndLoad(loadCopy)
	if err != nil {
		return err
	}
	s.conf.Store(loadCopy)
	return nil
}

func (s *Service[T]) Set(c *T) error {
	cf := *s.Get()

	err := mergo.Merge(new(cf), c)
	if err != nil {
		return err
	}

	return s.storeAndLoad(c)
}

type FieldVal struct {
	IsStruct bool

	Key       string
	Value     any
	FieldType string

	Help     string
	Env      string
	EnvSet   bool
	Default  string
	IsFolder bool
	IsSecret bool

	Nested map[string]FieldVal
}

func (s *Service[T]) GetSchema() ([]byte, error) {
	var conf T
	conf = *s.Get()

	of := reflect.ValueOf(&conf)
	meta := s.parseConf(of)

	return json.Marshal(meta)
}

func (s *Service[T]) parseConf(v reflect.Value) map[string]FieldVal {
	// If it's a pointer, dereference it.
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	topPairs := make(map[string]FieldVal)
	defaultPrefixer := s.prefixer

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldT := t.Field(i)
		fieldV := v.Field(i)

		if !fieldT.IsExported() {
			continue
		}

		keyName := fieldT.Name
		// full path (e.g., "database.host")
		if fieldV.Kind() == reflect.Struct {
			topPairs[keyName] = FieldVal{
				IsStruct: true,
				Nested:   s.parseConf(fieldV),
			}
			continue
		}

		var actualValue any
		if fieldV.CanInterface() {
			actualValue = fieldV.Interface()
		}

		env := fieldT.Tag.Get("env")

		_, hide := fieldT.Tag.Lookup("hide")
		_, isFolder := fieldT.Tag.Lookup("folder")

		fieldVal := FieldVal{
			Key:       keyName,
			Value:     actualValue,
			FieldType: fieldV.Type().String(),
			Help:      fieldT.Tag.Get("help"),
			Env:       defaultPrefixer(env),
			IsSecret:  hide,
			IsFolder:  isFolder,
			Default:   fieldT.Tag.Get("default"),
			EnvSet:    os.Getenv(defaultPrefixer(env)) != "",
		}

		topPairs[keyName] = fieldVal
	}

	return topPairs
}
