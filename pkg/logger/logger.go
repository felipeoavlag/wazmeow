package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger
func Init(level, format string) {
	// Set log level
	switch strings.ToLower(level) {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log format
	if strings.ToLower(format) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// Add caller information
	log.Logger = log.Logger.With().Caller().Logger()
}

// Info returns an info level logger
func Info() *zerolog.Event {
	return log.Info()
}

// Debug returns a debug level logger
func Debug() *zerolog.Event {
	return log.Debug()
}

// Warn returns a warn level logger
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error returns an error level logger
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal returns a fatal level logger
func Fatal() *zerolog.Event {
	return log.Fatal()
}
