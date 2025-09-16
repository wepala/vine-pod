package service

import (
	"context"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
	"github.com/wepala/vine-pod/pkg/version"
)

// VersionService handles version information operations
type VersionService struct {
	config *config.Config
	logger logger.Logger
}

// NewVersionService creates a new version service
func NewVersionService(cfg *config.Config, logger logger.Logger) *VersionService {
	return &VersionService{
		config: cfg,
		logger: logger,
	}
}

// GetVersion handles version information requests
func (s *VersionService) GetVersion(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	s.logger.Debug("Version information requested",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
	)

	versionInfo := version.Get()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
		s.logger.Error("Failed to encode version response", zap.Error(err))
		return err
	}

	return nil
}