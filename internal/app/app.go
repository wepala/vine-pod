package app

import (
	"context"
	"fmt"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/internal/server"
	"github.com/wepala/vine-pod/pkg/logger"
)

// App represents the main application
type App struct {
	config *config.Config
	logger logger.Logger
	server *server.Server
}

// New creates a new application instance
func New(cfg *config.Config, logger logger.Logger) (*App, error) {
	// Create HTTP server
	srv, err := server.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return &App{
		config: cfg,
		logger: logger,
		server: srv,
	}, nil
}

// Run starts the application
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Starting Vine Pod application", "address", a.config.Address())

	// Start HTTP server
	if err := a.server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the application
func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	// Shutdown HTTP server
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}