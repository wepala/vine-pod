package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/wepala/vine-pod/internal/app"
	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/pkg/logger"
	"github.com/wepala/vine-pod/pkg/version"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Print version information
	logger.Info("Starting Vine Pod Solid Server",
		"version", version.Version,
		"commit", version.Commit,
		"build_time", version.BuildTime,
	)

	// Create application
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("Failed to create application", "error", err)
		os.Exit(1)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the application
	go func() {
		if err := application.Run(ctx); err != nil {
			logger.Error("Application failed to start", "error", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", "signal", sig)
	case <-ctx.Done():
		logger.Info("Application context cancelled")
	}

	// Graceful shutdown
	logger.Info("Shutting down gracefully...")
	if err := application.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown gracefully", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}