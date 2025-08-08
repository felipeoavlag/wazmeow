package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// Logger é nossa interface de logging interna
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// zeroLogger implementa a interface do whatsmeow logger usando zerolog
type zeroLogger struct {
	mod string
	zerolog.Logger
}

// NewZeroLogger cria uma nova instância do logger usando zerolog
func NewZeroLogger(level string) *zeroLogger {
	// Configurar nível de log
	logLevel := parseLogLevel(level)
	zerolog.SetGlobalLevel(logLevel)

	// Configurar output com timestamp e cores
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return &zeroLogger{Logger: logger}
}

// parseLogLevel converte string para zerolog.Level
func parseLogLevel(level string) zerolog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "WARN", "WARNING":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

// Implementação da interface waLog.Logger (para compatibilidade com whatsmeow)
func (z *zeroLogger) Warnf(msg string, args ...interface{})  { z.Logger.Warn().Msgf(msg, args...) }
func (z *zeroLogger) Errorf(msg string, args ...interface{}) { z.Logger.Error().Msgf(msg, args...) }
func (z *zeroLogger) Infof(msg string, args ...interface{})  { z.Logger.Info().Msgf(msg, args...) }
func (z *zeroLogger) Debugf(msg string, args ...interface{}) { z.Logger.Debug().Msgf(msg, args...) }

func (z *zeroLogger) Sub(module string) waLog.Logger {
	if z.mod != "" {
		module = fmt.Sprintf("%s/%s", z.mod, module)
	}
	return &zeroLogger{mod: module, Logger: z.Logger.With().Str("sublogger", module).Logger()}
}

// Implementação da nossa interface Logger interna
func (z *zeroLogger) Debug(msg string, args ...interface{}) { z.Logger.Debug().Msgf(msg, args...) }
func (z *zeroLogger) Info(msg string, args ...interface{})  { z.Logger.Info().Msgf(msg, args...) }
func (z *zeroLogger) Warn(msg string, args ...interface{})  { z.Logger.Warn().Msgf(msg, args...) }
func (z *zeroLogger) Error(msg string, args ...interface{}) { z.Logger.Error().Msgf(msg, args...) }
func (z *zeroLogger) Fatal(msg string, args ...interface{}) {
	z.Logger.Error().Msgf(msg, args...)
	os.Exit(1)
}

// InitLogger inicializa o logger padrão usando zerolog
func InitLogger() Logger {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "INFO"
	}
	return NewZeroLogger(level)
}

// NewLogger mantém compatibilidade com código existente
func NewLogger(level string) Logger {
	return NewZeroLogger(level)
}

// GlobalLogger é uma instância global do logger para uso em toda a aplicação
var GlobalLogger *zeroLogger

// InitGlobalLogger inicializa o logger global
func InitGlobalLogger(level string) {
	GlobalLogger = NewZeroLogger(level)
}

// Funções de conveniência para logging global
func Debug(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Debugf(msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Infof(msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Warnf(msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Errorf(msg, args...)
	}
}

func Fatal(msg string, args ...interface{}) {
	if GlobalLogger != nil {
		GlobalLogger.Errorf(msg, args...)
	}
	os.Exit(1)
}

// GetWhatsmeowLogger retorna um logger compatível com whatsmeow
func GetWhatsmeowLogger() waLog.Logger {
	if GlobalLogger == nil {
		// Se o logger global não estiver inicializado, criar um temporário
		return NewZeroLogger("INFO")
	}
	return GlobalLogger
}

// ForWhatsApp retorna um logger compatível com whatsmeow (nome descritivo)
func ForWhatsApp() waLog.Logger {
	return GetWhatsmeowLogger()
}
