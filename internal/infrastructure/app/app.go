package app

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/infrastructure/config"
	"github.com/wepala/vine-pod/internal/infrastructure/server"
	"github.com/wepala/vine-pod/pkg/logger"
)

// App represents the main application
type App struct {
	config       *config.Config
	logger       logger.Logger
	kratosServer *server.SimpleKratosServer
}

// New creates a new application instance
func New(cfg *config.Config, logger logger.Logger) (*App, error) {
	// Create Kratos HTTP server
	srv, err := server.NewSimpleKratosServer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kratos server: %w", err)
	}

	return &App{
		config:       cfg,
		logger:       logger,
		kratosServer: srv,
	}, nil
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting Vine Pod application", zap.String("address", a.config.Address()))

	// Start Kratos HTTP server
	if err := a.kratosServer.Start(ctx); err != nil {
		return fmt.Errorf("failed to start Kratos server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	// Shutdown Kratos HTTP server
	if err := a.kratosServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown Kratos server: %w", err)
	}

	return nil
}
