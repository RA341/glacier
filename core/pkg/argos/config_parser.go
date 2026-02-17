package argos

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type FieldProcessor func(field reflect.Value, fieldType reflect.StructField)

func LoadStruct(ptr any, functor FieldProcessor) {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Pointer && field.IsNil() {
			// if it's a struct pointer we want to recurse into
			if field.Type().Elem().Kind() == reflect.Struct {

				field.Set(reflect.New(field.Type().Elem()))
			}
		}

		fVal := field
		if fVal.Kind() == reflect.Pointer {
			fVal = fVal.Elem()
		}
		if fVal.Kind() == reflect.Struct {
			LoadStruct(fVal.Addr().Interface(), functor)
			continue
		}

		functor(field, fieldType)
	}
}

// FieldProcessorTag reads tags and applies them in order
//
//	priority order
//
// 1. 'env' - always overwrites any field value
// 2. isFieldSet
// 3. 'default' - skipped if field is already set
// and applies the value
func FieldProcessorTag(prefixer Prefixer) FieldProcessor {
	return func(field reflect.Value, fieldType reflect.StructField) {

		if field.Kind() == reflect.Map {
			if field.IsNil() {
				newMap := reflect.MakeMap(field.Type())
				field.Set(newMap)
			}
			return
		}

		defaultValue := fieldType.Tag.Get("default")
		if defaultValue == "" {
			log.Fatal().
				Str("field", fieldType.Name).
				Msg("Default value not set for field, all non struct fields must have a default value, if field has no default value mark it with '-'")
		}
		if defaultValue == "-" {
			defaultValue = ""
		}

		if hasEnvTag(fieldType, field, prefixer) {
			// if an env is set then it is always updated regardless of prev value
			return
		}

		if !field.IsZero() {
			// field is set
			return
		}

		setField(fieldType, field, defaultValue)
	}
}

func hasEnvTag(fieldType reflect.StructField, field reflect.Value, prefixer Prefixer) bool {
	envValue := fieldType.Tag.Get("env")
	if envValue == "" {
		return false
	}

	fullEnv := prefixer(envValue)
	valToSet := os.Getenv(fullEnv)
	if valToSet != "" {
		setField(fieldType, field, valToSet)
		return true
	}

	return false
}

func resolveFolder(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to resolve path")
		return path
	}

	err = os.MkdirAll(abs, os.ModePerm)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create folder")
	}

	return abs
}

func setField(fieldType reflect.StructField, field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		if _, ok := fieldType.Tag.Lookup("folder"); ok {
			value = resolveFolder(value)
		}
		field.SetString(value)
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal().Err(err).Str("val", value).Msg("Failed to convert value to int")
		}
		field.SetInt(int64(i))
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			log.Fatal().Err(err).Str("val", value).Msg("Failed to convert value to uint")
		}
		field.SetUint(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatal().Err(err).Str("val", value).Msg("Failed to convert value to bool")
		}
		field.SetBool(b)
		break
	case reflect.Slice:
		elemType := field.Type().Elem()
		// 2. Only proceed if it is a slice of strings
		if elemType.Kind() == reflect.String {
			if value == "" {
				field.Set(reflect.MakeSlice(field.Type(), 0, 0))
			} else {
				split := strings.Split(value, ",")
				field.Set(reflect.ValueOf(split))
			}
			break
		}
		// unsupported slice type
		fallthrough
	default:
		log.Fatal().
			Str("field", fieldType.Name).
			Str("type", field.Type().String()).
			Msg("unsupported field ")
	}
}
