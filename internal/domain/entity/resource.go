package entity

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"github.com/wepala/vine-os/core/pericarp/pkg/domain"
	"github.com/wepala/vine-pod/internal/domain/event"
	"github.com/wepala/vine-pod/internal/domain/service"
)

// Resource defines the interface for all LDP resources
type Resource interface {
	// Core pericarp methods
	ID() string
	Version() int
	SequenceNo() int
	UncommittedEvents() []domain.Event
	MarkEventsAsCommitted()
	HasUncommittedEvents() bool
	UncommittedEventCount() int
	AddEvent(event domain.Event)
	HasErrors() bool
	GetErrors() []error
	AddError(err error)

	// Resource construction methods
	FromJSONLD(data string) Resource
	FromTurtle(data string) Resource
	FromRDFXML(data string) Resource
	WithURI(uri string) Resource

	// Resource operations
	Update(data string, contentType string) Resource
	Delete() Resource

	// Resource metadata
	GetContentType() string
	GetLastModified() time.Time
	GetETag() string
}

// BasicResource is the concrete implementation of Resource interface
type BasicResource struct {
	domain.Entity // Embedded pericarp entity (not pointer)

	// Resource-specific state derived from events
	uri          string
	contentType  string
	data         string
	lastModified time.Time
	etag         string
	errors       []error

	// Dependencies
	rdfValidator service.RDFValidationService
}

// NewBasicResource creates a new BasicResource instance
func NewBasicResource() Resource {
	return &BasicResource{
		Entity:       domain.NewEntity(""), // ID will be set when resource is created
		errors:       make([]error, 0),
		lastModified: time.Now(),
		rdfValidator: service.NewStandardRDFValidationService(),
	}
}

// NewBasicResourceWithValidator creates a new BasicResource instance with a custom RDF validator
func NewBasicResourceWithValidator(validator service.RDFValidationService) Resource {
	return &BasicResource{
		Entity:       domain.NewEntity(""), // ID will be set when resource is created
		errors:       make([]error, 0),
		lastModified: time.Now(),
		rdfValidator: validator,
	}
}

// FromJSONLD creates a resource from JSON-LD data
func (r *BasicResource) FromJSONLD(data string) Resource {
	if data == "" {
		r.AddError(errors.New("empty JSON-LD data"))
		return r
	}

	// Don't process if there are already errors
	if r.hasErrorsInChain() {
		return r
	}

	// Use the RDF validation service to validate and extract the resource ID
	resourceID, err := r.rdfValidator.ValidateJSONLD(data)
	if err != nil {
		r.AddError(fmt.Errorf("JSON-LD validation failed: %w", err))
		return r
	}

	// Initialize entity with ID extracted from the resource content if not already set
	if r.ID() == "" {
		r.Entity = domain.NewEntity(resourceID)
	}

	// Create and add the event
	createdEvent := event.NewResourceCreatedEvent(resourceID, data, "application/ld+json", resourceID)
	r.AddEvent(createdEvent)

	// Apply the event to update state
	r.applyResourceCreatedEvent(createdEvent)

	return r
}

// FromTurtle creates a resource from Turtle data
func (r *BasicResource) FromTurtle(data string) Resource {
	if data == "" {
		r.AddError(errors.New("empty Turtle data"))
		return r
	}

	// Don't process if there are already errors
	if r.hasErrorsInChain() {
		return r
	}

	// Use the RDF validation service to validate and extract the resource ID
	resourceID, err := r.rdfValidator.ValidateTurtle(data)
	if err != nil {
		r.AddError(fmt.Errorf("Turtle validation failed: %w", err))
		return r
	}

	// Initialize entity with ID extracted from the resource content if not already set
	if r.ID() == "" {
		r.Entity = domain.NewEntity(resourceID)
	}

	// Create and add the event
	createdEvent := event.NewResourceCreatedEvent(resourceID, data, "text/turtle", resourceID)
	r.AddEvent(createdEvent)

	// Apply the event to update state
	r.applyResourceCreatedEvent(createdEvent)

	return r
}

// FromRDFXML creates a resource from RDF/XML data
func (r *BasicResource) FromRDFXML(data string) Resource {
	if data == "" {
		r.AddError(errors.New("empty RDF/XML data"))
		return r
	}

	// Don't process if there are already errors
	if r.hasErrorsInChain() {
		return r
	}

	// Use the RDF validation service to validate and extract the resource ID
	resourceID, err := r.rdfValidator.ValidateRDFXML(data)
	if err != nil {
		r.AddError(fmt.Errorf("RDF/XML validation failed: %w", err))
		return r
	}

	// Initialize entity with ID extracted from the resource content if not already set
	if r.ID() == "" {
		r.Entity = domain.NewEntity(resourceID)
	}

	// Create and add the event
	createdEvent := event.NewResourceCreatedEvent(resourceID, data, "application/rdf+xml", resourceID)
	r.AddEvent(createdEvent)

	// Apply the event to update state
	r.applyResourceCreatedEvent(createdEvent)

	return r
}

// WithURI assigns a URI to the resource
func (r *BasicResource) WithURI(uri string) Resource {
	if uri == "" {
		r.AddError(errors.New("invalid URI"))
		return r
	}

	// Don't process if there are already errors
	if r.hasErrorsInChain() {
		return r
	}

	// Create and add the event
	uriEvent := event.NewResourceURIAssignedEvent(r.ID(), uri)
	r.AddEvent(uriEvent)

	// Apply the event to update state
	r.applyResourceURIAssignedEvent(uriEvent)

	return r
}

// Update updates the resource with new data
func (r *BasicResource) Update(data string, contentType string) Resource {
	if r.hasErrorsInChain() {
		return r // Don't process if there are already errors
	}

	if data == "" {
		r.AddError(errors.New("empty update data"))
		return r
	}

	// Store previous data for the event
	previousData := r.data

	// Create and add the event
	updateEvent := event.NewResourceUpdatedEvent(r.ID(), previousData, data, contentType)
	r.AddEvent(updateEvent)

	// Apply the event to update state
	r.applyResourceUpdatedEvent(updateEvent)

	return r
}

// Delete marks the resource as deleted
func (r *BasicResource) Delete() Resource {
	if r.hasErrorsInChain() {
		return r // Don't process if there are already errors
	}

	// Create and add the event
	deleteEvent := event.NewResourceDeletedEvent(r.ID(), r.uri)
	r.AddEvent(deleteEvent)

	// Apply the event to update state
	r.applyResourceDeletedEvent(deleteEvent)

	return r
}

// GetContentType returns the content type of the resource
func (r *BasicResource) GetContentType() string {
	return r.contentType
}

// GetLastModified returns the last modified time
func (r *BasicResource) GetLastModified() time.Time {
	return r.lastModified
}

// GetETag returns the entity tag for caching
func (r *BasicResource) GetETag() string {
	if r.etag == "" {
		// Generate ETag based on content and last modified time
		content := fmt.Sprintf("%s-%d", r.data, r.lastModified.Unix())
		r.etag = fmt.Sprintf(`"%x"`, md5.Sum([]byte(content)))
	}
	return r.etag
}

// HasErrors returns true if the resource has accumulated errors
func (r *BasicResource) HasErrors() bool {
	return len(r.errors) > 0
}

// GetErrors returns all accumulated errors
func (r *BasicResource) GetErrors() []error {
	return r.errors
}

// AddError adds an error to the resource
func (r *BasicResource) AddError(err error) {
	r.errors = append(r.errors, err)
}

// Private helper methods

func (r *BasicResource) hasErrorsInChain() bool {
	return r.HasErrors()
}

// Universal event application interface for all resource created events
type ResourceCreatedEventInterface interface {
	Data() string
	ContentType() string
	OccurredAt() time.Time
}

func (r *BasicResource) applyResourceCreatedEvent(event ResourceCreatedEventInterface) {
	r.data = event.Data()
	r.contentType = event.ContentType()
	r.lastModified = event.OccurredAt()
	r.etag = "" // Reset ETag so it will be recalculated
}

func (r *BasicResource) applyResourceURIAssignedEvent(event *event.ResourceURIAssignedEvent) {
	r.uri = event.URI()
	r.lastModified = event.OccurredAt()
	r.etag = "" // Reset ETag so it will be recalculated
}

func (r *BasicResource) applyResourceUpdatedEvent(event *event.ResourceUpdatedEvent) {
	r.data = event.NewData()
	r.contentType = event.ContentType()
	r.lastModified = event.OccurredAt()
	r.etag = "" // Reset ETag so it will be recalculated
}

func (r *BasicResource) applyResourceDeletedEvent(event *event.ResourceDeletedEvent) {
	r.lastModified = event.OccurredAt()
	r.etag = "" // Reset ETag so it will be recalculated
}

// LoadFromHistory reconstructs the resource state from events (for event sourcing)
func (r *BasicResource) LoadFromHistory(events []domain.Event) {
	// Ensure RDF validator is initialized if not already set
	if r.rdfValidator == nil {
		r.rdfValidator = service.NewStandardRDFValidationService()
	}

	for _, evt := range events {
		switch e := evt.(type) {
		case *event.ResourceCreatedEvent:
			r.applyResourceCreatedEvent(e)
		case *event.ResourceURIAssignedEvent:
			r.applyResourceURIAssignedEvent(e)
		case *event.ResourceUpdatedEvent:
			r.applyResourceUpdatedEvent(e)
		case *event.ResourceDeletedEvent:
			r.applyResourceDeletedEvent(e)
		}
	}
	// Call base implementation to update version and sequence
	r.Entity.LoadFromHistory(events)
}
