package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
	"github.com/wepala/vine-pod/pkg/version"
)

// Handlers holds all HTTP handlers
type Handlers struct {
	config *config.Config
	logger logger.Logger
}

// New creates a new handlers instance
func New(cfg *config.Config, logger logger.Logger) *Handlers {
	return &Handlers{
		config: cfg,
		logger: logger,
	}
}

// Health handles health check requests
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"service": "vine-pod",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Version handles version information requests
func (h *Handlers) Version(w http.ResponseWriter, r *http.Request) {
	versionInfo := version.Get()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versionInfo)
}

// Root handles root path requests
func (h *Handlers) Root(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "Welcome to Vine Pod - Solid Server",
		"version": version.Version,
		"endpoints": map[string]string{
			"health":  "/health",
			"version": "/version",
			"solid":   "/solid/",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SolidHandler handles Solid protocol requests (placeholder)
func (h *Handlers) SolidHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Solid protocol request received",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("user_agent", r.Header.Get("User-Agent")),
	)

	// This is a placeholder for future Solid protocol implementation
	response := map[string]interface{}{
		"message": "Solid protocol endpoint",
		"method":  r.Method,
		"path":    r.URL.Path,
		"status":  "not_implemented",
		"note":    "This endpoint will be implemented with full Solid protocol support",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(response)
}