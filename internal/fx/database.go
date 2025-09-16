package fx

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/wepala/vine-pod/internal/config"
	"github.com/wepala/vine-pod/internal/database"
	"github.com/wepala/vine-pod/pkg/logger"
)

// DatabaseModule provides database-related dependencies
var DatabaseModule = fx.Module("database",
	fx.Provide(NewDatabase),
	fx.Invoke(RegisterDatabaseLifecycle),
)

// NewDatabase creates a new GORM database connection
func NewDatabase(cfg *config.Config, logger logger.Logger) (*gorm.DB, error) {
	return database.NewGormDB(cfg, logger)
}

// RegisterDatabaseLifecycle registers database lifecycle hooks with Fx
func RegisterDatabaseLifecycle(
	lc fx.Lifecycle,
	db *gorm.DB,
	logger logger.Logger,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Test database connection on startup
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}

			if err := sqlDB.PingContext(ctx); err != nil {
				logger.Error("Database ping failed during startup", zap.Error(err))
				return err
			}

			logger.Info("Database connection verified on startup")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Close database connection on shutdown
			sqlDB, err := db.DB()
			if err != nil {
				logger.Error("Failed to get underlying DB for shutdown", zap.Error(err))
				return err
			}

			if err := sqlDB.Close(); err != nil {
				logger.Error("Failed to close database connection", zap.Error(err))
				return err
			}

			logger.Info("Database connection closed successfully")
			return nil
		},
	})
}
