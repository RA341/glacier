package mapsct

import (
	"fmt"

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
