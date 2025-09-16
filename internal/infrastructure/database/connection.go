package database

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/wepala/vine-pod/internal/infrastructure/config"
	applogger "github.com/wepala/vine-pod/pkg/logger"
)

// NewGormDB creates a new GORM database connection based on configuration
func NewGormDB(cfg *config.Config, log applogger.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Database.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.Database.DSN)
	case "sqlite":
		// Create data directory if it doesn't exist
		if err := ensureDataDir(cfg.Database.DSN); err != nil {
			return nil, fmt.Errorf("failed to create data directory: %w", err)
		}
		dialector = sqlite.Open(cfg.Database.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}

	// Create GORM logger that integrates with our Zap logger
	gormLogger := NewGormLogger(log)

	// Open database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Database connection established",
		zap.String("driver", cfg.Database.Driver),
		zap.String("dsn", maskDSN(cfg.Database.DSN)),
		zap.Int("max_open_conns", cfg.Database.MaxOpenConns),
		zap.Int("max_idle_conns", cfg.Database.MaxIdleConns),
	)

	return db, nil
}

// ensureDataDir creates the data directory for SQLite if it doesn't exist
func ensureDataDir(dsn string) error {
	dir := filepath.Dir(dsn)
	if dir == "." {
		return nil // Current directory, no need to create
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	return nil
}

// maskDSN masks sensitive information in DSN for logging
func maskDSN(dsn string) string {
	// For SQLite, return as-is
	if !contains(dsn, "://") {
		return dsn
	}

	// For PostgreSQL, mask password
	// This is a simple implementation, could be improved
	return "postgresql://***:***@host:port/dbname"
}

// contains checks if string s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findIndex(s, substr) >= 0
}

// findIndex returns the index of the first occurrence of substr in s, or -1 if not found
func findIndex(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
