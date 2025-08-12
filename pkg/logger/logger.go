package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	waLog "go.mau.fi/whatsmeow/util/log"
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

// WAAdapter adapts our centralized logger to whatsmeow's log interface
type WAAdapter struct {
	module string
}

// NewWALogger creates a new WhatsApp logger adapter
func NewWALogger(module string) waLog.Logger {
	return &WAAdapter{
		module: module,
	}
}

// Errorf logs an error message
func (w *WAAdapter) Errorf(msg string, args ...interface{}) {
	Error().Str("module", w.module).Msgf(msg, args...)
}

// Warnf logs a warning message
func (w *WAAdapter) Warnf(msg string, args ...interface{}) {
	Warn().Str("module", w.module).Msgf(msg, args...)
}

// Infof logs an info message
func (w *WAAdapter) Infof(msg string, args ...interface{}) {
	Info().Str("module", w.module).Msgf(msg, args...)
}

// Debugf logs a debug message
func (w *WAAdapter) Debugf(msg string, args ...interface{}) {
	Debug().Str("module", w.module).Msgf(msg, args...)
}

// Sub creates a sub-logger with additional context
func (w *WAAdapter) Sub(module string) waLog.Logger {
	return &WAAdapter{
		module: fmt.Sprintf("%s/%s", w.module, module),
	}
}
