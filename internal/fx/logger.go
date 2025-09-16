package fx

import (
	"go.uber.org/fx"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
)

// LoggerModule provides logging-related dependencies
var LoggerModule = fx.Module("logger",
	fx.Provide(NewLogger),
)

// NewLogger creates a new logger instance based on configuration
func NewLogger(cfg *config.Config) logger.Logger {
	return logger.New(cfg.LogLevel)
}
