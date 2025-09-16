package fx

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/internal/server"
	"github.com/wepala/vine-pod/pkg/logger"
)

// ServerModule provides server-related dependencies
var ServerModule = fx.Module("server",
	fx.Provide(NewKratosServer),
	fx.Invoke(RegisterServerLifecycle),
)

// NewKratosServer creates a new Kratos server instance
func NewKratosServer(cfg *config.Config, logger logger.Logger) (*server.SimpleKratosServer, error) {
	return server.NewSimpleKratosServer(cfg, logger)
}

// RegisterServerLifecycle registers server lifecycle hooks with Fx
func RegisterServerLifecycle(
	lc fx.Lifecycle,
	srv *server.SimpleKratosServer,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.Start(ctx); err != nil {
					logger.Error("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}
