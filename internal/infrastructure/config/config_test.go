package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test with default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify default values
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host '0.0.0.0', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}

	if cfg.LogLevel != "info" {
		t.Errorf("Expected default log level 'info', got '%s'", cfg.LogLevel)
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("SOLID_DATA_PATH", "/tmp/test")

	defer func() {
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("SOLID_DATA_PATH")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify environment variable values
	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected host 'localhost', got '%s'", cfg.Server.Host)
	}

	if cfg.Server.Port != 9000 {
		t.Errorf("Expected port 9000, got %d", cfg.Server.Port)
	}

	if cfg.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got '%s'", cfg.LogLevel)
	}

	if cfg.Solid.DataPath != "/tmp/test" {
		t.Errorf("Expected data path '/tmp/test', got '%s'", cfg.Solid.DataPath)
	}
}

func TestAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 3000,
		},
	}

	expected := "127.0.0.1:3000"
	if addr := cfg.Address(); addr != expected {
		t.Errorf("Expected address '%s', got '%s'", expected, addr)
	}
}

func TestGetEnvDuration(t *testing.T) {
	// Test with valid duration
	os.Setenv("TEST_DURATION", "45s")
	defer os.Unsetenv("TEST_DURATION")

	duration := getEnvDuration("TEST_DURATION", "30s")
	expected := 45 * time.Second
	if duration != expected {
		t.Errorf("Expected duration %v, got %v", expected, duration)
	}

	// Test with invalid duration (should return default)
	os.Setenv("TEST_DURATION", "invalid")
	duration = getEnvDuration("TEST_DURATION", "60s")
	expected = 60 * time.Second
	if duration != expected {
		t.Errorf("Expected duration %v, got %v", expected, duration)
	}
}
