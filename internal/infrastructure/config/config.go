package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	LogLevel string
	Solid    SolidConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver          string        // "postgres" or "sqlite"
	DSN             string        // Data Source Name
	MaxOpenConns    int           // Maximum open connections
	MaxIdleConns    int           // Maximum idle connections
	ConnMaxLifetime time.Duration // Connection maximum lifetime
	ConnMaxIdleTime time.Duration // Connection maximum idle time
}

// SolidConfig holds Solid protocol specific configuration
type SolidConfig struct {
	DataPath    string
	AllowOrigin string
	EnableCORS  bool
}

// Load reads configuration from environment variables and returns Config
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvDuration("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getEnvDuration("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getEnvDuration("SERVER_IDLE_TIMEOUT", "60s"),
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			DSN:             getEnv("DB_DSN", "./data/vine-pod.db"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", "5m"),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", "5m"),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Solid: SolidConfig{
			DataPath:    getEnv("SOLID_DATA_PATH", "./data"),
			AllowOrigin: getEnv("SOLID_ALLOW_ORIGIN", "*"),
			EnableCORS:  getEnvBool("SOLID_ENABLE_CORS", true),
		},
	}

	return cfg, nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvDuration(key, defaultValue string) time.Duration {
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

// Address returns the server address string
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
