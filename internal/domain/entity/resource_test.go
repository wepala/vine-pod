package entity_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wepala/vine-pod/internal/domain/entity"
	"github.com/wepala/vine-pod/internal/domain/event"
)

func TestResource_FromJSONLD(t *testing.T) {
	t.Run("creates resource from valid JSON-LD with @id", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@id": "https://alice.example.com/notes/note1",
			"@type": "schema:Note",
			"schema:text": "My first note"
		}`

		// Act
		resource := entity.NewBasicResource().FromJSONLD(jsonLD)

		// Assert
		assert.NotNil(t, resource)
		assert.Equal(t, "https://alice.example.com/notes/note1", resource.ID())
		assert.False(t, resource.HasErrors())
		assert.True(t, resource.HasUncommittedEvents())

		events := resource.UncommittedEvents()
		assert.Len(t, events, 1)

		createdEvent, ok := events[0].(*event.ResourceCreatedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.created", createdEvent.EventType())
		assert.Equal(t, "https://alice.example.com/notes/note1", createdEvent.ExtractedID())
		assert.Equal(t, jsonLD, createdEvent.Data())
		assert.Equal(t, "application/ld+json", createdEvent.ContentType())
	})

	t.Run("adds error when JSON-LD data is empty", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act
		result := resource.FromJSONLD("")

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "empty JSON-LD data")
	})

	t.Run("adds error when JSON-LD is invalid JSON", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		invalidJSON := `{"@id": "incomplete`

		// Act
		result := resource.FromJSONLD(invalidJSON)

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "JSON-LD validation failed")
	})

	t.Run("adds error when JSON-LD is missing @id", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		jsonLD := `{
			"@type": "schema:Note",
			"schema:text": "Note without @id"
		}`

		// Act
		result := resource.FromJSONLD(jsonLD)

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "JSON-LD validation failed")
	})
}

func TestResource_FromTurtle(t *testing.T) {
	t.Run("creates resource from valid Turtle data", func(t *testing.T) {
		// Arrange
		turtle := `@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix schema: <http://schema.org/> .

<https://alice.example.com/profile/card#me>
    a foaf:Person ;
    foaf:name "Alice Smith" ;
    schema:jobTitle "Software Developer" .`

		// Act
		resource := entity.NewBasicResource().FromTurtle(turtle)

		// Assert
		assert.NotNil(t, resource)
		assert.Equal(t, "https://alice.example.com/profile/card#me", resource.ID())
		assert.False(t, resource.HasErrors())
		assert.True(t, resource.HasUncommittedEvents())

		events := resource.UncommittedEvents()
		assert.Len(t, events, 1)

		createdEvent, ok := events[0].(*event.ResourceCreatedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.created", createdEvent.EventType())
		assert.Equal(t, "https://alice.example.com/profile/card#me", createdEvent.ExtractedID())
		assert.Equal(t, turtle, createdEvent.Data())
		assert.Equal(t, "text/turtle", createdEvent.ContentType())
	})

	t.Run("adds error when Turtle data is empty", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act
		result := resource.FromTurtle("")

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "empty Turtle data")
	})
}

func TestResource_FromRDFXML(t *testing.T) {
	t.Run("creates resource from valid RDF/XML data", func(t *testing.T) {
		// Arrange
		rdfxml := `<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF
  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
  xmlns:foaf="http://xmlns.com/foaf/0.1/">

  <foaf:Person rdf:about="https://alice.example.com/contacts/bob">
    <foaf:name>Bob Johnson</foaf:name>
    <foaf:mbox rdf:resource="mailto:bob@example.com"/>
  </foaf:Person>
</rdf:RDF>`

		// Act
		resource := entity.NewBasicResource().FromRDFXML(rdfxml)

		// Assert
		assert.NotNil(t, resource)
		assert.Equal(t, "https://alice.example.com/contacts/bob", resource.ID())
		assert.False(t, resource.HasErrors())
		assert.True(t, resource.HasUncommittedEvents())

		events := resource.UncommittedEvents()
		assert.Len(t, events, 1)

		createdEvent, ok := events[0].(*event.ResourceCreatedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.created", createdEvent.EventType())
		assert.Equal(t, "https://alice.example.com/contacts/bob", createdEvent.ExtractedID())
		assert.Equal(t, rdfxml, createdEvent.Data())
		assert.Equal(t, "application/rdf+xml", createdEvent.ContentType())
	})

	t.Run("adds error when RDF/XML data is empty", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act
		result := resource.FromRDFXML("")

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "empty RDF/XML data")
	})
}

func TestResource_WithURI(t *testing.T) {
	t.Run("assigns URI to resource", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		uri := "https://alice.example.com/notes/note1"

		// Act
		result := resource.WithURI(uri)

		// Assert
		assert.NotNil(t, result)
		assert.False(t, result.HasErrors())
		assert.True(t, result.HasUncommittedEvents())

		events := result.UncommittedEvents()
		assert.Len(t, events, 1)

		uriEvent, ok := events[0].(*event.ResourceURIAssignedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.uri_assigned", uriEvent.EventType())
		assert.Equal(t, uri, uriEvent.URI())
	})

	t.Run("adds error when URI is empty", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act
		result := resource.WithURI("")

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "invalid URI")
	})
}

func TestResource_Update(t *testing.T) {
	t.Run("updates resource with new data", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		// First create the resource with initial data
		jsonLD1 := `{"@id": "https://example.com/resource1", "title": "Initial Title"}`
		resource.FromJSONLD(jsonLD1)
		resource.MarkEventsAsCommitted() // Simulate persistence

		newData := `{"@id": "https://example.com/resource1", "title": "Updated Title"}`
		contentType := "application/ld+json"

		// Act
		result := resource.Update(newData, contentType)

		// Assert
		assert.NotNil(t, result)
		assert.False(t, result.HasErrors())
		assert.True(t, result.HasUncommittedEvents())

		events := result.UncommittedEvents()
		assert.Len(t, events, 1)

		updateEvent, ok := events[0].(*event.ResourceUpdatedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.updated", updateEvent.EventType())
		assert.Equal(t, newData, updateEvent.NewData())
		assert.Equal(t, contentType, updateEvent.ContentType())
	})

	t.Run("adds error when update data is empty", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act
		result := resource.Update("", "application/ld+json")

		// Assert
		assert.True(t, result.HasErrors())
		assert.False(t, result.HasUncommittedEvents())

		errors := result.GetErrors()
		assert.Len(t, errors, 1)
		assert.Contains(t, errors[0].Error(), "empty update data")
	})
}

func TestResource_Delete(t *testing.T) {
	t.Run("deletes resource", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		// First create and assign URI to the resource
		jsonLD := `{"@id": "https://example.com/resource1", "title": "Title"}`
		resource.FromJSONLD(jsonLD).WithURI("https://example.com/resource1")
		resource.MarkEventsAsCommitted() // Simulate persistence

		// Act
		result := resource.Delete()

		// Assert
		assert.NotNil(t, result)
		assert.False(t, result.HasErrors())
		assert.True(t, result.HasUncommittedEvents())

		events := result.UncommittedEvents()
		assert.Len(t, events, 1)

		deleteEvent, ok := events[0].(*event.ResourceDeletedEvent)
		assert.True(t, ok)
		assert.Equal(t, "resource.deleted", deleteEvent.EventType())
	})
}

func TestResource_Metadata(t *testing.T) {
	t.Run("GetContentType returns content type", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		jsonLD := `{"@id": "https://example.com/resource1", "title": "Title"}`
		resource.FromJSONLD(jsonLD)

		// Act
		contentType := resource.GetContentType()

		// Assert
		assert.Equal(t, "application/ld+json", contentType)
	})

	t.Run("GetLastModified returns last modified time", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		jsonLD := `{"@id": "https://example.com/resource1", "title": "Title"}`
		resource.FromJSONLD(jsonLD)

		// Act
		lastModified := resource.GetLastModified()

		// Assert
		assert.WithinDuration(t, time.Now(), lastModified, time.Second)
	})

	t.Run("GetETag returns entity tag", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()
		jsonLD := `{"@id": "https://example.com/resource1", "title": "Title"}`
		resource.FromJSONLD(jsonLD)

		// Act
		etag := resource.GetETag()

		// Assert
		assert.NotEmpty(t, etag)
	})
}

func TestResource_ChainedOperations(t *testing.T) {
	t.Run("can chain multiple operations", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@id": "https://alice.example.com/notes/note1",
			"@type": "schema:Note",
			"schema:text": "My first note"
		}`
		uri := "https://alice.example.com/notes/note1"

		// Act
		resource := entity.NewBasicResource().
			FromJSONLD(jsonLD).
			WithURI(uri)

		// Assert
		assert.NotNil(t, resource)
		assert.Equal(t, "https://alice.example.com/notes/note1", resource.ID())
		assert.False(t, resource.HasErrors())
		assert.True(t, resource.HasUncommittedEvents())

		events := resource.UncommittedEvents()
		assert.Len(t, events, 2) // CreatedFromJSONLD + URIAssigned
	})

	t.Run("stops chain on error", func(t *testing.T) {
		// Act
		resource := entity.NewBasicResource().
			FromJSONLD(""). // This should cause an error
			WithURI("https://example.com/resource")

		// Assert
		assert.True(t, resource.HasErrors())
		// Should not have URI assigned event due to error in chain
		events := resource.UncommittedEvents()
		assert.Len(t, events, 0)
	})
}

func TestResource_ErrorAccumulation(t *testing.T) {
	t.Run("accumulates multiple errors", func(t *testing.T) {
		// Arrange
		resource := entity.NewBasicResource()

		// Act - perform multiple operations that should fail
		result := resource.
			FromJSONLD(""). // Error 1
			WithURI("")     // Error 2

		// Assert
		assert.True(t, result.HasErrors())
		errors := result.GetErrors()
		assert.Len(t, errors, 2)
	})
}

// Test interface compliance
func TestResource_InterfaceCompliance(t *testing.T) {
	t.Run("BasicResource implements Resource interface", func(t *testing.T) {
		var _ entity.Resource = entity.NewBasicResource()
	})
}
