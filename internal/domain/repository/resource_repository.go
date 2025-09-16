package repository

import (
	"context"
)

//go:generate moq -out resource_repository_mock.go . ResourceRepository

// ResourceRepository defines the interface for resource persistence operations
type ResourceRepository interface {
	// Create a new resource
	Create(ctx context.Context, resource *Resource) error

	// GetByID retrieves a resource by its ID
	GetByID(ctx context.Context, id string) (*Resource, error)

	// GetByURI retrieves a resource by its URI
	GetByURI(ctx context.Context, uri string) (*Resource, error)

	// Update an existing resource
	Update(ctx context.Context, resource *Resource) error

	// Delete a resource by ID
	Delete(ctx context.Context, id string) error

	// List resources with pagination
	List(ctx context.Context, limit, offset int) ([]*Resource, error)

	// FindByContainer retrieves all resources in a container
	FindByContainer(ctx context.Context, containerURI string) ([]*Resource, error)
}

// Resource represents a domain resource entity
type Resource struct {
	ID          string `json:"id"`
	URI         string `json:"uri"`
	ContentType string `json:"content_type"`
	Data        string `json:"data"`
	Container   string `json:"container,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
