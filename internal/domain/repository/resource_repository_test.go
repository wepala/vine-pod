package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wepala/vine-pod/internal/domain/repository"
	"github.com/wepala/vine-pod/test/fixtures"
)

func TestResourceRepositoryMock(t *testing.T) {
	ctx := fixtures.TestContext()

	// Create a mock repository
	mockRepo := &repository.ResourceRepositoryMock{
		CreateFunc: func(ctx context.Context, resource *repository.Resource) error {
			// Mock successful creation
			return nil
		},
		GetByIDFunc: func(ctx context.Context, id string) (*repository.Resource, error) {
			// Mock returning a test resource
			return &repository.Resource{
				ID:          id,
				URI:         "http://example.com/resource/" + id,
				ContentType: "application/ld+json",
				Data:        `{"@id": "http://example.com/resource/` + id + `"}`,
				CreatedAt:   1234567890,
				UpdatedAt:   1234567890,
			}, nil
		},
		GetByURIFunc: func(ctx context.Context, uri string) (*repository.Resource, error) {
			// Mock returning a resource by URI
			return &repository.Resource{
				ID:          "test-id",
				URI:         uri,
				ContentType: "application/ld+json",
				Data:        `{"@id": "` + uri + `"}`,
				CreatedAt:   1234567890,
				UpdatedAt:   1234567890,
			}, nil
		},
	}

	t.Run("Create resource", func(t *testing.T) {
		resource := &repository.Resource{
			ID:          "test-resource-1",
			URI:         "http://example.com/resource/test-resource-1",
			ContentType: "application/ld+json",
			Data:        `{"@id": "http://example.com/resource/test-resource-1"}`,
			CreatedAt:   1234567890,
			UpdatedAt:   1234567890,
		}

		err := mockRepo.Create(ctx, resource)
		require.NoError(t, err)

		// Verify the mock was called
		assert.Len(t, mockRepo.CreateCalls(), 1)
		assert.Equal(t, resource, mockRepo.CreateCalls()[0].Resource)
	})

	t.Run("Get resource by ID", func(t *testing.T) {
		resourceID := "test-resource-1"

		resource, err := mockRepo.GetByID(ctx, resourceID)
		require.NoError(t, err)
		require.NotNil(t, resource)

		assert.Equal(t, resourceID, resource.ID)
		assert.Equal(t, "http://example.com/resource/"+resourceID, resource.URI)
		assert.Equal(t, "application/ld+json", resource.ContentType)

		// Verify the mock was called
		assert.Len(t, mockRepo.GetByIDCalls(), 1)
		assert.Equal(t, resourceID, mockRepo.GetByIDCalls()[0].ID)
	})

	t.Run("Get resource by URI", func(t *testing.T) {
		resourceURI := "http://example.com/resource/test-resource-1"

		resource, err := mockRepo.GetByURI(ctx, resourceURI)
		require.NoError(t, err)
		require.NotNil(t, resource)

		assert.Equal(t, "test-id", resource.ID)
		assert.Equal(t, resourceURI, resource.URI)
		assert.Equal(t, "application/ld+json", resource.ContentType)

		// Verify the mock was called
		assert.Len(t, mockRepo.GetByURICalls(), 1)
		assert.Equal(t, resourceURI, mockRepo.GetByURICalls()[0].URI)
	})
}

func TestResourceRepositoryTestUtils(t *testing.T) {
	t.Run("Test database utilities", func(t *testing.T) {
		// Test creating an in-memory database
		db := fixtures.TestDB(t)
		defer fixtures.CleanupDB(t, db)

		// Test database connection
		fixtures.AssertDBConnection(t, db)

		// Test that we can run basic SQL
		result := db.Exec("CREATE TABLE test_resources (id TEXT PRIMARY KEY, uri TEXT)")
		require.NoError(t, result.Error)

		// Insert test data
		result = db.Exec("INSERT INTO test_resources (id, uri) VALUES (?, ?)", "test-1", "http://example.com/test-1")
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		// Query test data
		var count int64
		result = db.Model(&struct {
			ID  string `gorm:"column:id"`
			URI string `gorm:"column:uri"`
		}{}).Table("test_resources").Count(&count)
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Test configuration utilities", func(t *testing.T) {
		cfg := fixtures.TestConfig()
		require.NotNil(t, cfg)

		assert.Equal(t, "localhost", cfg.Server.Host)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "debug", cfg.LogLevel)
		assert.Equal(t, "sqlite", cfg.Database.Driver)
		assert.Equal(t, ":memory:", cfg.Database.DSN)
	})

	t.Run("Test logging utilities", func(t *testing.T) {
		logger := fixtures.TestLogger()
		require.NotNil(t, logger)

		// Test that logger works without errors
		logger.Info("test log message")
		logger.Debug("debug message")
		logger.Error("error message")
	})
}
