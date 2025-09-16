package fx

import (
	"go.uber.org/fx"

	"github.com/wepala/vine-pod/internal/config"
)

// ConfigModule provides configuration-related dependencies
var ConfigModule = fx.Module("config",
	fx.Provide(NewConfig),
)

// NewConfig creates and loads the application configuration
func NewConfig() (*config.Config, error) {
	return config.Load()
}
