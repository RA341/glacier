package argos

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	ColorReset     = "\033[0m"
	ColorBold      = "\033[1m"
	ColorUnderline = "\033[4m"

	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m" // For string values
	ColorYellow  = "\033[33m" // For bools
	ColorBlue    = "\033[34m" // For keys
	ColorMagenta = "\033[35m" // For numbers
	ColorCyan    = "\033[36m"
)

type FieldPrintConfig struct {
	TagName string
	// TagName is the same as the one passed above
	PrintConfig func(TagName string, val *FieldVal)

	// max length of value for this tag
	maxTagLength int
}

func PrintInfo(c interface{}, footer string, opts ...FieldPrintConfig) {
	var tags = make([]string, len(opts))
	for i, opt := range opts {
		tags[i] = opt.TagName
	}

	// start with an empty structPrefix since its at tippity top
	pairs := flattenStruct(
		reflect.ValueOf(c),
		"",
		tags...,
	)

	// Find the length of the longest key for alignment.
	maxKeyLength := 0
	maxValueLength := 0
	//maxHelpLength := 0
	for _, val := range pairs {
		if len(val.Key) > maxKeyLength {
			maxKeyLength = len(val.Key)
		}

		// strip ANSI color codes to get the true visible length of the value.
		cleanValue := ansiRegex.ReplaceAllString(val.Value, "")
		if len(cleanValue) > maxValueLength {
			maxValueLength = len(cleanValue)
		}

		for i, o := range opts {
			// apply any config the client could make
			o.PrintConfig(o.TagName, &val)

			src := val.Tags[o.TagName]
			cleanTag := ansiRegex.ReplaceAllString(src, "")
			if len(cleanTag) > o.maxTagLength {
				o.maxTagLength = len(cleanTag)
				opts[i] = o
			}
		}
	}

	var contentBuilder strings.Builder
	for i, val := range pairs {
		if val.Value == "" {
			// empty imply a map type
			continue
		}

		// padding for the key column
		keyPadding := strings.Repeat(" ", maxKeyLength-len(val.Key))

		// padding for the value column
		cleanValue := ansiRegex.ReplaceAllString(val.Value, "")
		valuePadding := strings.Repeat(" ", maxValueLength-len(cleanValue))

		// Colorize parts for readability
		coloredKey := Colorize(val.Key, ColorBlue+ColorBold)

		// Assemble the line with calculated padding
		// Format: [Key]:[Padding]  [Value][Padding]   [....tags]
		contentBuilder.WriteString(coloredKey)
		contentBuilder.WriteString(":")
		contentBuilder.WriteString(keyPadding)
		contentBuilder.WriteString("  ") // Separator

		contentBuilder.WriteString(val.Value)
		contentBuilder.WriteString(valuePadding)
		contentBuilder.WriteString("  ")

		// tags
		for j, o := range opts {
			tagToPrint := val.Tags[o.TagName]

			cleanTagVal := ansiRegex.ReplaceAllString(tagToPrint, "")
			tagPadding := ""
			if o.maxTagLength > 0 {
				tagPadding = strings.Repeat(" ", o.maxTagLength-len(cleanTagVal))
			}

			contentBuilder.WriteString(tagToPrint)
			// we dont need padding for that last column
			if j < len(opts)-1 {
				contentBuilder.WriteString(tagPadding)
			}
			contentBuilder.WriteString("  ")
		}

		if i < len(pairs)-1 {
			contentBuilder.WriteString("\n")
		}
	}

	contentBuilder.WriteString(fmt.Sprintf("\n\n%s", footer))

	printInBox("Config", contentBuilder.String())
}

func WithUnderLine(value string) string {
	return Colorize(value, ColorRed+ColorUnderline)
}

// FieldVal
//
//	{
//		structKey: {
//			tag1: val1,
//			tag2: val2,
//			tag3: val3,
//		},
//	    structKey2: {
//				tag1: val1,
//				tag2: val2,
//				tag3: val3,
//			},
//		}
type FieldVal struct {
	Key   string
	Value string
	Tags  map[string]string
}

// flattenStruct recursively traverses a struct and returns a flat list of KeyValue pairs.
func flattenStruct(v reflect.Value, structPrefix string, tags ...string) []FieldVal {
	// If it's a pointer, dereference it.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var topPairs []FieldVal

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldT := t.Field(i)
		fieldV := v.Field(i)

		keyName := fieldT.Name
		// full path (e.g., "database.host")
		fullKey := keyName
		if structPrefix != "" {
			fullKey = structPrefix + "." + keyName
		}

		if fieldV.Kind() == reflect.Struct {
			nestedPairs := flattenStruct(fieldV, fullKey, tags...)
			topPairs = append(topPairs, nestedPairs...)
			continue
		}

		fieldVal := FieldVal{
			Key:  fullKey,
			Tags: make(map[string]string),
		}
		for _, tag := range tags {
			tagValue, ok := fieldT.Tag.Lookup(tag)
			if !ok {
				continue
			}
			fieldVal.Tags[tag] = tagValue
		}

		switch fieldV.Kind() {
		case reflect.Slice:
			// Format slices as comma-separated strings
			var sliceItems []string
			for j := 0; j < fieldV.Len(); j++ {
				sliceItems = append(sliceItems, formatSimpleValue(fieldV.Index(j)))
			}
			fieldVal.Value = strings.Join(sliceItems, ", ")
		default:
			fieldVal.Value = formatSimpleValue(fieldV)
		}

		topPairs = append(topPairs, fieldVal)
	}

	return topPairs
}

// formatSimpleValue converts a reflect.Value of a simple type to a colored string.
func formatSimpleValue(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return Colorize(v.String(), ColorGreen)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Colorize(strconv.FormatInt(v.Int(), 10), ColorMagenta)
	case reflect.Bool:
		return Colorize(strconv.FormatBool(v.Bool()), ColorYellow)
	default:
		// this indicates that the value is not printable
		return ""
		//if v.IsValid() {
		//	return v.String()
		//}
		//return Colorize("null", "red") // Should not happen with valid structs
	}
}

func Colorize(s, color string) string {
	return color + s + ColorReset
}

const (
	topLeft     = "╭"
	topRight    = "╮"
	bottomLeft  = "╰"
	bottomRight = "╯"
	horizontal  = "─"
	vertical    = "│"
	space       = " "
	hPadding    = 2 // Horizontal padding on each side
	vPadding    = 1 // Vertical padding (empty lines top/bottom)
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)

func printInBox(title, content string) {
	lines := strings.Split(content, "\n")
	maxWidth := utf8.RuneCountInString(title)
	for _, line := range lines {
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		if width := utf8.RuneCountInString(cleanLine); width > maxWidth {
			maxWidth = width
		}
	}
	innerWidth := maxWidth + (hPadding * 2)

	titleBar := buildTitleBar(title, innerWidth, horizontal, space)
	fmt.Println(topLeft + titleBar + topRight)

	printPaddedLine(innerWidth, vertical, space)

	for _, line := range lines {
		cleanLine := ansiRegex.ReplaceAllString(line, "")
		visibleWidth := utf8.RuneCountInString(cleanLine)
		paddingNeeded := innerWidth - visibleWidth

		fmt.Printf("%s%s%s%s%s\n",
			vertical,
			strings.Repeat(space, hPadding),
			line,
			strings.Repeat(space, paddingNeeded-hPadding),
			vertical,
		)
	}

	printPaddedLine(innerWidth, vertical, space)
	fmt.Println(bottomLeft + strings.Repeat(horizontal, innerWidth) + bottomRight)
}

func buildTitleBar(title string, innerWidth int, char, space string) string {
	titleText := space + title + space
	titleWidth := utf8.RuneCountInString(titleText)
	if titleWidth >= innerWidth {
		return titleText[:innerWidth]
	}
	padding := innerWidth - titleWidth
	return strings.Repeat(char, padding/2) + titleText + strings.Repeat(char, padding-(padding/2))
}

func printPaddedLine(innerWidth int, vertical, space string) {
	fmt.Println(vertical + strings.Repeat(space, innerWidth) + vertical)
}
