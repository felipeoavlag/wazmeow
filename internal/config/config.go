package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	WhatsApp WhatsAppConfig
	Log      LogConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	Debug    bool
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host            string
	Port            int
	SSLCertFile     string
	SSLKeyFile      string
	ShutdownTimeout time.Duration
}

// WhatsAppConfig holds WhatsApp configuration
type WhatsAppConfig struct {
	Debug                string
	OSName               string
	MaxSessions          int
	ConnectionTimeout    time.Duration
	QRTimeout            time.Duration
	ReconnectInterval    time.Duration
	MaxReconnectAttempts int
	PoolSize             int
	PoolMaxIdle          int
	PoolMaxLifetime      time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "wazmeow"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "wazmeow"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Debug:    getEnv("DB_DEBUG", "") != "",
		},
		Server: ServerConfig{
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			Port:            getEnvAsInt("SERVER_PORT", 8080),
			SSLCertFile:     getEnv("SSL_CERT_FILE", ""),
			SSLKeyFile:      getEnv("SSL_KEY_FILE", ""),
			ShutdownTimeout: getEnvAsDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		WhatsApp: WhatsAppConfig{
			Debug:                getEnv("WA_DEBUG", ""),
			OSName:               getEnv("WA_OS_NAME", "Mac OS 10"),
			MaxSessions:          getEnvAsInt("WA_MAX_SESSIONS", 100),
			ConnectionTimeout:    getEnvAsDuration("WA_CONNECTION_TIMEOUT", 30*time.Second),
			QRTimeout:            getEnvAsDuration("WA_QR_TIMEOUT", 5*time.Minute),
			ReconnectInterval:    getEnvAsDuration("WA_RECONNECT_INTERVAL", 10*time.Second),
			MaxReconnectAttempts: getEnvAsInt("WA_MAX_RECONNECT_ATTEMPTS", 5),
			PoolSize:             getEnvAsInt("WA_POOL_SIZE", 50),
			PoolMaxIdle:          getEnvAsInt("WA_POOL_MAX_IDLE", 10),
			PoolMaxLifetime:      getEnvAsDuration("WA_POOL_MAX_LIFETIME", time.Hour),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "console"),
		},
	}

	return cfg, nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getEnvAsDuration gets an environment variable as duration with a fallback value
func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}
