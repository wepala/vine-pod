package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/internal/handler"
	"github.com/wepala/vine-pod/internal/middleware"
	"github.com/wepala/vine-pod/pkg/logger"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	logger logger.Logger
	server *http.Server
}

// New creates a new HTTP server
func New(cfg *config.Config, logger logger.Logger) (*Server, error) {
	// Create handlers
	handlers := handler.New(cfg, logger)

	// Setup router with middleware
	mux := http.NewServeMux()

	// Apply middleware
	var handler http.Handler = mux
	handler = middleware.CORS(cfg)(handler)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Recovery(logger)(handler)

	// Register routes
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("GET /version", handlers.Version)
	mux.HandleFunc("GET /{$}", handlers.Root) // Use explicit root pattern

	// Solid protocol routes (placeholder for future implementation)
	mux.HandleFunc("/solid/", handlers.SolidHandler)

	// Create HTTP server
	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		config: cfg,
		logger: logger,
		server: server,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting HTTP server", zap.String("address", s.server.Addr))

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}
