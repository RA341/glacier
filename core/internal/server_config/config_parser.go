package server_config

import (
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type FieldProcessor func(field reflect.Value, fieldType reflect.StructField)

func SetDefaultsFromTags(ptr interface{}, functor FieldProcessor) {
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			// if it's a struct pointer we want to recurse into
			if field.Type().Elem().Kind() == reflect.Struct {

				field.Set(reflect.New(field.Type().Elem()))
			}
		}

		fVal := field
		if fVal.Kind() == reflect.Ptr {
			fVal = fVal.Elem()
		}
		if fVal.Kind() == reflect.Struct {
			SetDefaultsFromTags(fVal.Addr().Interface(), functor)
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
				Msg("Default value not set for field, all non struct fields must have a default value")
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

func setField(fieldType reflect.StructField, field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
		break
	case reflect.Int:
		i, _ := strconv.Atoi(value)
		field.SetInt(int64(i))
		break
	case reflect.Bool:
		b, _ := strconv.ParseBool(value)
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
