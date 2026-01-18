package logger

import (
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogConfig struct {
	Level      string
	Verbose    bool
	Writer     io.Writer
	TimeFormat string
	NoColor    bool
	Caller     bool
}

// parseLevel parses log level string with fallback
func parseLevel(levelStr string) zerolog.Level {
	level, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(levelStr)))
	if err != nil {
		log.Warn().Str("invalid_level", levelStr).Msg("Invalid log level, using info")

		customLevel, err := strconv.Atoi(levelStr)
		if err == nil {
			log.Info().Int("level", customLevel).Msg("Setting custom log level")
			return zerolog.Level(customLevel)
		}

		return zerolog.InfoLevel
	}
	return level
}

// createConsoleWriter creates a configured console writer
func createConsoleWriter(config LogConfig) zerolog.ConsoleWriter {
	//writer := zerolog.ConsoleWriter{
	//	Out:        config.Writer,
	//	TimeFormat: config.TimeFormat,
	//	NoColor:    config.NoColor,
	//}
	//
	//// Customize output format if needed
	//if config.Verbose {
	//	writer.FormatLevel = func(i interface{}) string {
	//		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	//	}
	//}

	return zerolog.ConsoleWriter{
		Out:        config.Writer,
		TimeFormat: config.TimeFormat,
		NoColor:    config.NoColor,
	}
}

// createLogger creates a configured logger
func createLogger(config LogConfig) zerolog.Logger {
	level := parseLevel(config.Level)
	writer := createConsoleWriter(config)

	// Create base logger
	logger := zerolog.New(writer).With().Timestamp()

	// Add caller info if requested
	if config.Caller || config.Verbose {
		logger = logger.Caller()
	}

	return logger.Logger().Level(level)
}

// Init initializes the global logger with configuration
func Init(config LogConfig) {
	logger := createLogger(config)

	// Set global logger
	log.Logger = logger
	zerolog.SetGlobalLevel(parseLevel(config.Level))
}

// DefaultConfig returns sensible defaults
func DefaultConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		Verbose:    false,
		Writer:     os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
		NoColor:    false,
		Caller:     false,
	}
}

// InitDefault initializes logger with default configuration
func InitDefault() {
	Init(DefaultConfig())
}

// InitConsole initializes console logger with specified level and verbose mode
func InitConsole(level string, verbose bool) {
	config := DefaultConfig()
	config.Level = level
	config.Verbose = verbose
	config.Caller = verbose

	Init(config)
}

// TestConfig returns config optimized for testing
func TestConfig() LogConfig {
	return LogConfig{
		Level:      "debug",
		Verbose:    true,
		Writer:     os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    true,
		Caller:     true,
	}
}

// InitForTest initializes logger optimized for testing
func InitForTest() {
	Init(TestConfig())
}

// InitSilent initializes a silent logger (disabled)
func InitSilent() {
	config := DefaultConfig()
	config.Level = "disabled"
	Init(config)
}

// GetLogger returns a new logger instance with the same configuration
func GetLogger() zerolog.Logger {
	return log.Logger
}

// GetLoggerWithFields returns a logger with additional fields
func GetLoggerWithFields(fields map[string]interface{}) zerolog.Logger {
	logger := log.Logger
	for key, value := range fields {
		logger = logger.With().Interface(key, value).Logger()
	}
	return logger
}

// SetLevel changes the global log level at runtime
func SetLevel(level string) {
	parsedLevel := parseLevel(level)
	zerolog.SetGlobalLevel(parsedLevel)
	log.Logger = log.Logger.Level(parsedLevel)

	log.Info().Str("new_level", level).Msg("Log level changed")
}

// IsLevelEnabled checks if a log level is enabled
func IsLevelEnabled(level zerolog.Level) bool {
	return log.Logger.GetLevel() <= level
}

// WithContext returns a logger with context fields
func WithContext(component string) zerolog.Logger {
	return log.With().Str("component", component).Logger()
}

//func ExampleUsage() {
//	// Basic usage
//	InitConsole("debug", true)
//
//	// With custom config
//	config := Config{
//		Level:      "warn",
//		Verbose:    false,
//		Writer:     os.Stdout,
//		TimeFormat: "15:04:05",
//		NoColor:    true,
//		Caller:     true,
//	}
//	Init(config)
//
//	// Component-specific logger
//	apiLogger := WithContext("api")
//	apiLogger.Info().Msg("API server starting")
//
//	// Logger with fields
//	dbLogger := GetLoggerWithFields(map[string]interface{}{
//		"component": "database",
//		"version":   "1.0.0",
//	})
//	dbLogger.Error().Msg("Database connection failed")
//
//	// Runtime level change
//	SetLevel("trace")
//
//	// Check if level is enabled
//	if IsLevelEnabled(zerolog.DebugLevel) {
//		log.Debug().Msg("Debug logging is enabled")
//	}
//}
