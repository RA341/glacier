package mapsct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

func ParseMap(output interface{}, raw map[string]any) error {
	mapConfig := &mapstructure.DecoderConfig{
		// error if a field in the struct is NOT in the source map
		ErrorUnset: true,
		Result:     output,
	}

	decoder, err := mapstructure.NewDecoder(mapConfig)
	if err != nil {
		return err
	}

	err = decoder.Decode(raw)
	if err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	return nil
}

type FieldSchema struct {
	Name      string `json:"name"`
	Type      string `json:"type"`                 // string, number, boolean, map, array
	InsertKey string `json:"insert_key"`           // The path for the UI to use
	KeyType   string `json:"key_type,omitempty"`   // e.g. "string" if it's a map
	ValueType string `json:"value_type,omitempty"` // e.g. "string" if it's a map or array
}

func GetSchema(input interface{}) ([]FieldSchema, error) {
	t := reflect.TypeOf(input)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("GetSchema only accepts structs, got %s", t.Kind())
	}

	var schema []FieldSchema

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if !field.IsExported() {
			continue
		}

		name := field.Tag.Get("mapstructure")
		if name == "" || name == "-" {
			name = field.Name
		} else {
			name = strings.Split(name, ",")[0]
		}

		if field.Type.Kind() == reflect.Struct {
			return nil, fmt.Errorf("nested structs are not supported: found struct at field '%s'", field.Name)
		}

		info := FieldSchema{
			Name:      name,
			InsertKey: name,
		}

		switch field.Type.Kind() {
		case reflect.String:
			info.Type = "string"

		case reflect.Bool:
			info.Type = "boolean"

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			info.Type = "number"

		case reflect.Map:
			info.Type = "map"
			info.KeyType = getSimpleType(field.Type.Key())
			info.ValueType = getSimpleType(field.Type.Elem())

			if field.Type.Elem().Kind() == reflect.Struct {
				return nil, fmt.Errorf("nested structs inside maps are not supported: field '%s'", field.Name)
			}

		case reflect.Slice, reflect.Array:
			info.Type = "array"
			info.ValueType = getSimpleType(field.Type.Elem())

			if field.Type.Elem().Kind() == reflect.Struct {
				return nil, fmt.Errorf("nested structs inside slices are not supported: field '%s'", field.Name)
			}

		default:
			info.Type = "unknown"
		}

		schema = append(schema, info)
	}

	return schema, nil
}

func getSimpleType(t reflect.Type) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	default:
		return t.Kind().String()
	}
}
