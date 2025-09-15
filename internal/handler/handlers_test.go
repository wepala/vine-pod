package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
	"github.com/wepala/vine-pod/pkg/version"
)

func TestHealthHandler(t *testing.T) {
	cfg, _ := config.Load()
	logger := logger.New("info")
	handlers := New(cfg, logger)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handlers.Health(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "vine-pod" {
		t.Errorf("Expected service 'vine-pod', got '%s'", response["service"])
	}
}

func TestVersionHandler(t *testing.T) {
	cfg, _ := config.Load()
	logger := logger.New("info")
	handlers := New(cfg, logger)

	req := httptest.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	handlers.Version(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response version.Info
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check that we get version information
	if response.Version == "" {
		t.Error("Expected version to be set")
	}

	if response.GoVersion == "" {
		t.Error("Expected go_version to be set")
	}

	if response.Platform == "" {
		t.Error("Expected platform to be set")
	}
}

func TestRootHandler(t *testing.T) {
	cfg, _ := config.Load()
	logger := logger.New("info")
	handlers := New(cfg, logger)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handlers.Root(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["message"] == nil {
		t.Error("Expected message field to be present")
	}

	if response["endpoints"] == nil {
		t.Error("Expected endpoints field to be present")
	}
}

func TestSolidHandler(t *testing.T) {
	cfg, _ := config.Load()
	logger := logger.New("info")
	handlers := New(cfg, logger)

	req := httptest.NewRequest("GET", "/solid/test", nil)
	w := httptest.NewRecorder()

	handlers.SolidHandler(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("Expected status code %d, got %d", http.StatusNotImplemented, w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "not_implemented" {
		t.Errorf("Expected status 'not_implemented', got '%v'", response["status"])
	}
}