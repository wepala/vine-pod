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