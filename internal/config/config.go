package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config representa a configuração completa da aplicação
type Config struct {
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
	Log      LogConfig      `json:"log"`
	CORS     CORSConfig     `json:"cors"`
	Session  SessionConfig  `json:"session"`
	Security SecurityConfig `json:"security"`
	App      AppConfig      `json:"app"`
	Webhook  WebhookConfig  `json:"webhook"`
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
	Debug    bool   `json:"debug"`
}

// ServerConfig configurações do servidor HTTP
type ServerConfig struct {
	Host         string        `json:"host"`
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

// LogConfig configurações de logging
type LogConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// CORSConfig configurações de CORS
type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

// SessionConfig configurações das sessões WhatsApp
type SessionConfig struct {
	MaxSessions    int           `json:"max_sessions"`
	SessionTimeout time.Duration `json:"session_timeout"`
}

// SecurityConfig configurações de segurança
type SecurityConfig struct {
	APIKey            string        `json:"api_key"`
	RateLimitRequests int           `json:"rate_limit_requests"`
	RateLimitWindow   time.Duration `json:"rate_limit_window"`
}

// AppConfig configurações gerais da aplicação
type AppConfig struct {
	Debug       bool   `json:"debug"`
	Environment string `json:"environment"`
}

// WebhookConfig configurações de webhook
type WebhookConfig struct {
	Timeout        time.Duration        `json:"timeout"`
	MaxRetries     int                  `json:"max_retries"`
	RetryDelay     time.Duration        `json:"retry_delay"`
	Workers        int                  `json:"workers"`
	QueueSize      int                  `json:"queue_size"`
	CircuitBreaker CircuitBreakerConfig `json:"circuit_breaker"`
	RateLimit      RateLimitConfig      `json:"rate_limit"`
}

// CircuitBreakerConfig configurações do circuit breaker
type CircuitBreakerConfig struct {
	MaxFailures      int           `json:"max_failures"`
	ResetTimeout     time.Duration `json:"reset_timeout"`
	HalfOpenMaxCalls int           `json:"half_open_max_calls"`
}

// RateLimitConfig configurações de rate limiting
type RateLimitConfig struct {
	RequestsPerSecond int           `json:"requests_per_second"`
	BurstSize         int           `json:"burst_size"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
}

// Load carrega a configuração das variáveis de ambiente
func Load() (*Config, error) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "wazmeow"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			Debug:    parseBool("DB_DEBUG", false),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  parseDuration("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: parseDuration("SERVER_WRITE_TIMEOUT", "30s"),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "INFO"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice("CORS_ALLOWED_ORIGINS", "*"),
			AllowedMethods: parseStringSlice("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
			AllowedHeaders: parseStringSlice("CORS_ALLOWED_HEADERS", "Accept,Authorization,Content-Type,X-CSRF-Token"),
		},
		Session: SessionConfig{
			MaxSessions:    parseInt("MAX_SESSIONS", 100),
			SessionTimeout: parseDuration("SESSION_TIMEOUT", "3600s"),
		},
		Security: SecurityConfig{
			APIKey:            getEnv("API_KEY", ""),
			RateLimitRequests: parseInt("RATE_LIMIT_REQUESTS", 100),
			RateLimitWindow:   parseDuration("RATE_LIMIT_WINDOW", "1m"),
		},
		App: AppConfig{
			Debug:       parseBool("DEBUG", false),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Webhook: WebhookConfig{
			Timeout:    parseDuration("WEBHOOK_TIMEOUT", "30s"),
			MaxRetries: parseInt("WEBHOOK_MAX_RETRIES", 3),
			RetryDelay: parseDuration("WEBHOOK_RETRY_DELAY", "5s"),
			Workers:    parseInt("WEBHOOK_WORKERS", 5),
			QueueSize:  parseInt("WEBHOOK_QUEUE_SIZE", 1000),
			CircuitBreaker: CircuitBreakerConfig{
				MaxFailures:      parseInt("WEBHOOK_CB_MAX_FAILURES", 5),
				ResetTimeout:     parseDuration("WEBHOOK_CB_RESET_TIMEOUT", "60s"),
				HalfOpenMaxCalls: parseInt("WEBHOOK_CB_HALF_OPEN_MAX_CALLS", 3),
			},
			RateLimit: RateLimitConfig{
				RequestsPerSecond: parseInt("WEBHOOK_RATE_LIMIT_RPS", 50),    // Aumentado de 10 para 50
				BurstSize:         parseInt("WEBHOOK_RATE_LIMIT_BURST", 100), // Aumentado de 20 para 100
				CleanupInterval:   parseDuration("WEBHOOK_RATE_LIMIT_CLEANUP", "5m"),
			},
		},
	}

	return config, nil
}

// GetDatabaseURL retorna a URL de conexão com o banco de dados
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetServerAddress retorna o endereço completo do servidor
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// GetServerURL retorna a URL base do servidor para webhooks
func (c *Config) GetServerURL() string {
	// Se o host for 0.0.0.0, usar 127.0.0.1 para URLs de webhook
	host := c.Server.Host
	if host == "0.0.0.0" {
		host = "127.0.0.1"
	}

	// Usar sempre HTTP para webhooks locais na porta 8080
	protocol := "http"

	return fmt.Sprintf("%s://%s:%s", protocol, host, c.Server.Port)
}

// GetSwaggerHost retorna o host para documentação Swagger
// Se SWAGGER_HOST estiver definido, usa ele, senão usa SERVER_HOST
// Se SERVER_HOST for 0.0.0.0, usa localhost para Swagger
func (c *Config) GetSwaggerHost() string {
	// Primeiro verifica se há um SWAGGER_HOST específico configurado
	if swaggerHost := getEnv("SWAGGER_HOST", ""); swaggerHost != "" {
		return swaggerHost
	}

	// Senão, usa o SERVER_HOST
	host := c.Server.Host

	// Se for 0.0.0.0, usar localhost para Swagger (mais apropriado para documentação)
	if host == "0.0.0.0" {
		return fmt.Sprintf("localhost:%s", c.Server.Port)
	}

	return fmt.Sprintf("%s:%s", host, c.Server.Port)
}

// IsProduction verifica se está em ambiente de produção
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment verifica se está em ambiente de desenvolvimento
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// Validate valida a configuração
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST é obrigatório")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER é obrigatório")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME é obrigatório")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT é obrigatório")
	}

	// Validar nível de log
	validLogLevels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	if !contains(validLogLevels, strings.ToUpper(c.Log.Level)) {
		return fmt.Errorf("LOG_LEVEL deve ser um de: %s", strings.Join(validLogLevels, ", "))
	}

	return nil
}

// Funções auxiliares
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseDuration(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	if parsed, err := time.ParseDuration(defaultValue); err == nil {
		return parsed
	}
	return 30 * time.Second
}

func parseStringSlice(key, defaultValue string) []string {
	value := getEnv(key, defaultValue)
	if value == "" {
		return []string{}
	}
	return strings.Split(value, ",")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
