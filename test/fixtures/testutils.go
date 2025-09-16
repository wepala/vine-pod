package fixtures

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/wepala/vine-pod/internal/infrastructure/config"
	"github.com/wepala/vine-pod/pkg/logger"
)

// TestConfig creates a test configuration
func TestConfig() *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
		LogLevel: "debug",
		Solid: config.SolidConfig{
			DataPath:   "./test_data",
			EnableCORS: true,
		},
		Database: config.DatabaseConfig{
			Driver:          "sqlite",
			DSN:             ":memory:",
			MaxOpenConns:    5,
			MaxIdleConns:    5,
			ConnMaxLifetime: 300000000000, // 5 minutes in nanoseconds
			ConnMaxIdleTime: 60000000000,  // 1 minute in nanoseconds
		},
	}
}

// TestLogger creates a test logger
func TestLogger() logger.Logger {
	return logger.New("debug")
}

// TestDB creates an in-memory SQLite database for testing
func TestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to create test database")

	// Test the connection
	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB")

	err = sqlDB.Ping()
	require.NoError(t, err, "Failed to ping test database")

	return db
}

// CleanupDB cleans up the test database
func CleanupDB(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB for cleanup")

	err = sqlDB.Close()
	require.NoError(t, err, "Failed to close test database")
}

// TestContext creates a test context with timeout
func TestContext() context.Context {
	return context.Background()
}

// TestResource creates a test resource for testing
func TestResource() map[string]interface{} {
	return map[string]interface{}{
		"@context": "https://www.w3.org/ns/ldp",
		"@id":      "http://example.com/resource/1",
		"@type":    "ldp:Resource",
		"title":    "Test Resource",
		"content":  "This is a test resource for BDD scenarios",
	}
}

// TestContainer creates a test container for testing
func TestContainer() map[string]interface{} {
	return map[string]interface{}{
		"@context": "https://www.w3.org/ns/ldp",
		"@id":      "http://example.com/container/",
		"@type":    "ldp:BasicContainer",
		"title":    "Test Container",
		"contains": []string{},
	}
}

// AssertDBConnection verifies database connection is working
func AssertDBConnection(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB")

	var result int
	err = sqlDB.QueryRow("SELECT 1").Scan(&result)
	require.NoError(t, err, "Database connection test failed")
	require.Equal(t, 1, result, "Database query should return 1")
}

// AssertNoDBConnections verifies no database connections are leaked
func AssertNoDBConnections(t *testing.T, db *gorm.DB) {
	t.Helper()

	sqlDB, err := db.DB()
	require.NoError(t, err, "Failed to get underlying sql.DB")

	stats := sqlDB.Stats()
	require.Equal(t, 0, stats.OpenConnections, "Should have no open database connections")
}
