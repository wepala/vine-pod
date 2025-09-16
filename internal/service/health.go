package service

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
)

// HealthService handles health check operations
type HealthService struct {
	config *config.Config
	logger logger.Logger
}

// NewHealthService creates a new health service
func NewHealthService(cfg *config.Config, logger logger.Logger) *HealthService {
	return &HealthService{
		config: cfg,
		logger: logger,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// GetHealth handles health check requests
func (s *HealthService) GetHealth(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s.logger.Debug("Health check requested",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	response := HealthResponse{
		Status:  "healthy",
		Service: "vine-pod",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.logger.Error("Failed to encode health response", zap.Error(err))
		return err
	}

	return nil
}
