package repository

import (
	"context"

	"github.com/wepala/vine-os/core/pericarp/pkg/domain"
	"github.com/wepala/vine-pod/internal/domain/entity"
)

//go:generate moq -out resource_repository_mock.go . ResourceRepository

// ResourceRepository defines the interface for resource persistence operations
// It works with the Resource domain entity interface and supports event sourcing
type ResourceRepository interface {
	// Save persists a resource entity and its uncommitted events
	Save(ctx context.Context, resource entity.Resource) error

	// GetByID retrieves a resource by its ID and reconstructs it from events
	GetByID(ctx context.Context, id string) (entity.Resource, error)

	// GetByURI retrieves a resource by its URI
	GetByURI(ctx context.Context, uri string) (entity.Resource, error)

	// Delete removes a resource by ID
	Delete(ctx context.Context, id string) error

	// List resources with pagination
	List(ctx context.Context, limit, offset int) ([]entity.Resource, error)

	// FindByContainer retrieves all resources in a container
	FindByContainer(ctx context.Context, containerURI string) ([]entity.Resource, error)

	// LoadEvents retrieves all events for a specific resource aggregate
	LoadEvents(ctx context.Context, aggregateID string) ([]domain.Event, error)

	// LoadEventsFromVersion retrieves events for a resource starting from a specific version
	LoadEventsFromVersion(ctx context.Context, aggregateID string, version int) ([]domain.Event, error)
}
