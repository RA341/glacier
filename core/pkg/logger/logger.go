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
	Writers    []io.Writer
	TimeFormat string
	NoColor    bool
	Caller     bool
}

// parseLevel parses log level string with fallback
func parseLevel(levelStr string) zerolog.Level {
	level, err := zerolog.ParseLevel(strings.ToLower(strings.TrimSpace(levelStr)))
	if err != nil {
		customLevel, err := strconv.Atoi(levelStr)
		if err == nil {
			log.Info().Int("level", customLevel).Msg("Setting custom log level")
			return zerolog.Level(customLevel)
		}

		log.Warn().Str("invalid_level", levelStr).Msg("Invalid log level, using info")
		return zerolog.InfoLevel
	}
	return level
}

// consoleWriter creates a configured console writer
func consoleWriter() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	}
}

// CreateLogger creates a configured logger
func CreateLogger(Level string, Verbose bool, writers ...io.Writer) zerolog.Context {
	level := parseLevel(Level)

	writer := zerolog.MultiLevelWriter(append(writers, consoleWriter())...)

	logger := zerolog.New(writer).Level(level).With().Timestamp()

	if Verbose {
		logger = logger.Caller()
	}

	return logger
}

// Init initializes the global logger with configuration
func Init(config LogConfig) {
	logger := CreateLogger(config.Level, config.Verbose)

	// Set global logger
	log.Logger = logger.Logger()
	zerolog.SetGlobalLevel(parseLevel(config.Level))
}

// DefaultConfig returns sensible defaults
func DefaultConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		Verbose:    false,
		Writers:    []io.Writer{os.Stderr},
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
	config := InitLoggerWithLevel(level, verbose)

	Init(config)
}

func InitLoggerWithLevel(level string, verbose bool) LogConfig {
	config := DefaultConfig()
	config.Level = level
	config.Verbose = verbose
	config.Caller = verbose
	return config
}

// TestConfig returns config optimized for testing
func TestConfig() LogConfig {
	return LogConfig{
		Level:      "debug",
		Verbose:    true,
		Writers:    []io.Writer{os.Stdout},
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
func GetLoggerWithFields(fields map[string]any) zerolog.Logger {
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
