package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wepala/vine-os/core/pericarp/pkg/domain"
	"github.com/wepala/vine-pod/internal/domain/event"
)

func TestResourceCreatedEvent(t *testing.T) {
	t.Run("NewResourceCreatedEvent creates event with correct properties", func(t *testing.T) {
		// Arrange
		resourceID := "resource-123"
		data := `{"@id": "https://example.com/resource1", "@type": "schema:Document"}`
		contentType := "application/ld+json"
		extractedID := "https://example.com/resource1"

		// Act
		event := event.NewResourceCreatedEvent(resourceID, data, contentType, extractedID)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.created", event.EventType())
		assert.Equal(t, resourceID, event.AggregateID())
		assert.Equal(t, 1, event.Version())
		assert.Equal(t, contentType, event.ContentType())
		assert.Equal(t, data, event.Data())
		assert.Equal(t, extractedID, event.ExtractedID())
		assert.WithinDuration(t, time.Now(), event.OccurredAt(), time.Second)
	})

	t.Run("SetVersion updates event version", func(t *testing.T) {
		// Arrange
		event := event.NewResourceCreatedEvent("resource-123", "{}", "application/ld+json", "https://example.com/resource1")

		// Act
		event.SetVersion(5)

		// Assert
		assert.Equal(t, 5, event.Version())
	})

	t.Run("implements domain.Event interface", func(t *testing.T) {
		// Arrange
		event := event.NewResourceCreatedEvent("resource-123", "{}", "application/ld+json", "https://example.com/resource1")

		// Assert
		var _ domain.Event = event
	})

	t.Run("NewResourceCreatedFromJSONLDEvent creates unified event", func(t *testing.T) {
		// Arrange
		resourceID := "resource-123"
		data := `{"@id": "https://example.com/resource1"}`
		extractedID := "https://example.com/resource1"

		// Act
		event := event.NewResourceCreatedFromJSONLDEvent(resourceID, data, extractedID)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.created", event.EventType())
		assert.Equal(t, "application/ld+json", event.ContentType())
	})

	t.Run("NewResourceCreatedFromTurtleEvent creates unified event", func(t *testing.T) {
		// Arrange
		resourceID := "resource-456"
		data := `<https://example.com/resource2> a <http://schema.org/Document> .`
		extractedID := "https://example.com/resource2"

		// Act
		event := event.NewResourceCreatedFromTurtleEvent(resourceID, data, extractedID)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.created", event.EventType())
		assert.Equal(t, "text/turtle", event.ContentType())
	})

	t.Run("NewResourceCreatedFromRDFXMLEvent creates unified event", func(t *testing.T) {
		// Arrange
		resourceID := "resource-789"
		data := `<rdf:RDF><rdf:Description rdf:about="https://example.com/resource3"/></rdf:RDF>`
		extractedID := "https://example.com/resource3"

		// Act
		event := event.NewResourceCreatedFromRDFXMLEvent(resourceID, data, extractedID)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.created", event.EventType())
		assert.Equal(t, "application/rdf+xml", event.ContentType())
	})
}

func TestResourceURIAssignedEvent(t *testing.T) {
	t.Run("NewResourceURIAssignedEvent creates event with correct properties", func(t *testing.T) {
		// Arrange
		resourceID := "resource-abc"
		uri := "https://alice.example.com/notes/note1"

		// Act
		event := event.NewResourceURIAssignedEvent(resourceID, uri)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.uri_assigned", event.EventType())
		assert.Equal(t, resourceID, event.AggregateID())
		assert.Equal(t, 1, event.Version())
		assert.Equal(t, uri, event.URI())
		assert.WithinDuration(t, time.Now(), event.OccurredAt(), time.Second)
	})

	t.Run("SetVersion updates event version", func(t *testing.T) {
		// Arrange
		event := event.NewResourceURIAssignedEvent("resource-abc", "https://example.com/resource")

		// Act
		event.SetVersion(2)

		// Assert
		assert.Equal(t, 2, event.Version())
	})

	t.Run("implements domain.Event interface", func(t *testing.T) {
		// Arrange
		event := event.NewResourceURIAssignedEvent("resource-abc", "https://example.com/resource")

		// Assert
		var _ domain.Event = event
	})
}

func TestResourceUpdatedEvent(t *testing.T) {
	t.Run("NewResourceUpdatedEvent creates event with correct properties", func(t *testing.T) {
		// Arrange
		resourceID := "resource-def"
		previousData := `{"@id": "https://example.com/resource4", "title": "Old Title"}`
		newData := `{"@id": "https://example.com/resource4", "title": "New Title"}`
		contentType := "application/ld+json"

		// Act
		event := event.NewResourceUpdatedEvent(resourceID, previousData, newData, contentType)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.updated", event.EventType())
		assert.Equal(t, resourceID, event.AggregateID())
		assert.Equal(t, 1, event.Version())
		assert.Equal(t, previousData, event.PreviousData())
		assert.Equal(t, newData, event.NewData())
		assert.Equal(t, contentType, event.ContentType())
		assert.WithinDuration(t, time.Now(), event.OccurredAt(), time.Second)
	})

	t.Run("SetVersion updates event version", func(t *testing.T) {
		// Arrange
		event := event.NewResourceUpdatedEvent("resource-def", "old", "new", "application/ld+json")

		// Act
		event.SetVersion(4)

		// Assert
		assert.Equal(t, 4, event.Version())
	})

	t.Run("implements domain.Event interface", func(t *testing.T) {
		// Arrange
		event := event.NewResourceUpdatedEvent("resource-def", "old", "new", "application/ld+json")

		// Assert
		var _ domain.Event = event
	})
}

func TestResourceDeletedEvent(t *testing.T) {
	t.Run("NewResourceDeletedEvent creates event with correct properties", func(t *testing.T) {
		// Arrange
		resourceID := "resource-ghi"
		uri := "https://alice.example.com/notes/note1"

		// Act
		event := event.NewResourceDeletedEvent(resourceID, uri)

		// Assert
		assert.NotNil(t, event)
		assert.Equal(t, "resource.deleted", event.EventType())
		assert.Equal(t, resourceID, event.AggregateID())
		assert.Equal(t, 1, event.Version())
		assert.Equal(t, uri, event.URI())
		assert.WithinDuration(t, time.Now(), event.OccurredAt(), time.Second)
	})

	t.Run("SetVersion updates event version", func(t *testing.T) {
		// Arrange
		event := event.NewResourceDeletedEvent("resource-ghi", "https://example.com/resource")

		// Act
		event.SetVersion(6)

		// Assert
		assert.Equal(t, 6, event.Version())
	})

	t.Run("implements domain.Event interface", func(t *testing.T) {
		// Arrange
		event := event.NewResourceDeletedEvent("resource-ghi", "https://example.com/resource")

		// Assert
		var _ domain.Event = event
	})
}

// Integration tests to verify all events implement the domain.Event interface
func TestAllEventsImplementDomainEventInterface(t *testing.T) {
	testCases := []struct {
		name  string
		event domain.Event
	}{
		{
			name:  "ResourceCreatedEvent",
			event: event.NewResourceCreatedEvent("id", "{}", "application/ld+json", "https://example.com/resource"),
		},
		{
			name:  "ResourceCreatedFromJSONLDEvent (factory)",
			event: event.NewResourceCreatedFromJSONLDEvent("id", "{}", "https://example.com/resource"),
		},
		{
			name:  "ResourceCreatedFromTurtleEvent (factory)",
			event: event.NewResourceCreatedFromTurtleEvent("id", "", "https://example.com/resource"),
		},
		{
			name:  "ResourceCreatedFromRDFXMLEvent (factory)",
			event: event.NewResourceCreatedFromRDFXMLEvent("id", "", "https://example.com/resource"),
		},
		{
			name:  "ResourceURIAssignedEvent",
			event: event.NewResourceURIAssignedEvent("id", "https://example.com/resource"),
		},
		{
			name:  "ResourceUpdatedEvent",
			event: event.NewResourceUpdatedEvent("id", "old", "new", "application/ld+json"),
		},
		{
			name:  "ResourceDeletedEvent",
			event: event.NewResourceDeletedEvent("id", "https://example.com/resource"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotEmpty(t, tc.event.EventType())
			assert.NotEmpty(t, tc.event.AggregateID())
			assert.Greater(t, tc.event.Version(), 0)
			assert.False(t, tc.event.OccurredAt().IsZero())
		})
	}
}
