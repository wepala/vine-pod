package fx

import (
	"go.uber.org/fx"
)

// AppModule is the main application module that combines all other modules
var AppModule = fx.Module("app",
	// Core modules
	ConfigModule,
	LoggerModule,

	// Infrastructure modules
	DatabaseModule,
	// Service modules
	ServicesModule,

	// Server module (includes lifecycle management)
	ServerModule,
)

// NewApp creates a new Fx application with all modules
func NewApp() *fx.App {
	return fx.New(
		AppModule,
		// Disable Fx's default logger to use our structured logger
		fx.NopLogger,
	)
}
