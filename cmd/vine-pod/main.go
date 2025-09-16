package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	fxmodules "github.com/wepala/vine-pod/internal/fx"
	"github.com/wepala/vine-pod/pkg/version"
)

func main() {
	// Create Fx application
	app := fxmodules.NewApp()

	// Create context for application lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the Fx application
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	// Log startup message
	log.Printf("Vine Pod Solid Server started - Version: %s, Commit: %s",
		version.Version, version.Commit)

	// Wait for shutdown signal
	select {
	case sig := <-sigChan:
		log.Printf("Received shutdown signal: %v", sig)
	case <-ctx.Done():
		log.Println("Application context cancelled")
	}

	// Graceful shutdown
	log.Println("Shutting down gracefully...")
	if err := app.Stop(ctx); err != nil {
		log.Fatalf("Failed to shutdown gracefully: %v", err)
	}

	log.Println("Server stopped")
}