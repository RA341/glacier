//go:build ignore

package internal

import "gorm.io/cli/gorm/genconfig"

var _ = genconfig.Config{

	//// Map `gen:"name"` tags to helper kinds
	//FieldNameMap: map[string]any{
	//	"json": JSON{}, // use a custom JSON helper where fields are tagged `gen:"json"`
	//},

	IncludeInterfaces: []any{"Store*"},
	// ignore all structs
	IncludeStructs: []any{""},
}
