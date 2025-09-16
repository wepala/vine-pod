package di

import (
	"go.uber.org/fx"

	"github.com/wepala/vine-pod/internal/application/service"
	"github.com/wepala/vine-pod/internal/infrastructure/config"
	"github.com/wepala/vine-pod/pkg/logger"
)

// ServicesModule provides all service dependencies
var ServicesModule = fx.Module("services",
	fx.Provide(
		NewHealthService,
		NewVersionService,
		NewSolidService,
	),
)

// NewHealthService creates a new health service
func NewHealthService(cfg *config.Config, logger logger.Logger) *service.HealthService {
	return service.NewHealthService(cfg, logger)
}

// NewVersionService creates a new version service
func NewVersionService(cfg *config.Config, logger logger.Logger) *service.VersionService {
	return service.NewVersionService(cfg, logger)
}

// NewSolidService creates a new solid service
func NewSolidService(cfg *config.Config, logger logger.Logger) *service.SolidService {
	return service.NewSolidService(cfg, logger)
}
