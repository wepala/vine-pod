package server

import (
	"context"
	"net/http"

	zaplog "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"

	"go.uber.org/zap"

	"github.com/wepala/vine-pod/internal/application/service"
	"github.com/wepala/vine-pod/internal/in
	"github.com/wepala/vine-pod/pkg/logger"
)

// SimpleKratosServer represents a simplified Kratos HTTP server
type SimpleKratosServer struct {
	config *config.Config
	logger logger.Logger
	server *kratoshttp.Server
}

// NewSimpleKratosServer creates a new simplified Kratos HTTP server
func NewSimpleKratosServer(cfg *config.Config, logger logger.Logger) (*SimpleKratosServer, error) {
	// Create Kratos logger adapter
	kratosLogger := kratoslog.With(zaplog.NewLogger(logger.GetZapLogger()),
		"service.name", "vine-pod",
		"service.version", "dev",
	)

	// Create services
	healthSvc := service.NewHealthService(cfg, logger)
	versionSvc := service.NewVersionService(cfg, logger)
	solidSvc := service.NewSolidService(cfg, logger)

	// Create Kratos HTTP server with middleware
	srv := kratoshttp.NewServer(
		kratoshttp.Address(cfg.Address()),
		kratoshttp.Middleware(
			recovery.Recovery(),
			logging.Server(kratosLogger),
		),
	)

	// Register routes using standard HTTP handlers
	srv.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := healthSvc.GetHealth(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		if err := versionSvc.GetVersion(r.Context(), w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	srv.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			if err := solidSvc.GetRoot(r.Context(), w, r); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			// Handle Solid protocol requests
			switch r.Method {
			case http.MethodGet:
				if err := solidSvc.GetResource(r.Context(), w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			case http.MethodPost:
				if err := solidSvc.CreateResource(r.Context(), w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			case http.MethodPut:
				if err := solidSvc.UpdateResource(r.Context(), w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			case http.MethodDelete:
				if err := solidSvc.DeleteResource(r.Context(), w, r); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	})

	// CORS will be handled manually in each handler for simplicity

	return &SimpleKratosServer{
		config: cfg,
		logger: logger,
		server: srv,
	}, nil
}

// Start starts the Kratos HTTP server
func (s *SimpleKratosServer) Start(ctx context.Context) error {
	s.logger.Info("Starting Kratos HTTP server", zap.String("address", s.config.Address()))

	go func() {
		if err := s.server.Start(ctx); err != nil {
			s.logger.Error("Kratos HTTP server failed", zap.Error(err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return nil
}

// Shutdown gracefully shuts down the Kratos HTTP server
func (s *SimpleKratosServer) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down Kratos HTTP server")

	if err := s.server.Stop(ctx); err != nil {
		s.logger.Error("Failed to shutdown Kratos server", zap.Error(err))
		return err
	}

	return nil
}
