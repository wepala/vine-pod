package event

import (
	"time"

	"github.com/wepala/vine-os/core/pericarp/pkg/domain"
)

// ResourceCreatedEvent is emitted when a resource is created from RDF data
type ResourceCreatedEvent struct {
	resourceID  string
	data        string
	contentType string
	extractedID string
	occurredAt  time.Time
	version     int
}

// NewResourceCreatedEvent creates a new ResourceCreatedEvent
func NewResourceCreatedEvent(resourceID, data, contentType, extractedID string) *ResourceCreatedEvent {
	return &ResourceCreatedEvent{
		resourceID:  resourceID,
		data:        data,
		contentType: contentType,
		extractedID: extractedID,
		occurredAt:  time.Now(),
		version:     1,
	}
}

// NewResourceCreatedFromJSONLDEvent creates a new ResourceCreatedEvent for JSON-LD data
func NewResourceCreatedFromJSONLDEvent(resourceID, data, extractedID string) *ResourceCreatedEvent {
	return NewResourceCreatedEvent(resourceID, data, "application/ld+json", extractedID)
}

// NewResourceCreatedFromTurtleEvent creates a new ResourceCreatedEvent for Turtle data
func NewResourceCreatedFromTurtleEvent(resourceID, data, extractedID string) *ResourceCreatedEvent {
	return NewResourceCreatedEvent(resourceID, data, "text/turtle", extractedID)
}

// NewResourceCreatedFromRDFXMLEvent creates a new ResourceCreatedEvent for RDF/XML data
func NewResourceCreatedFromRDFXMLEvent(resourceID, data, extractedID string) *ResourceCreatedEvent {
	return NewResourceCreatedEvent(resourceID, data, "application/rdf+xml", extractedID)
}

// EventType returns the event type identifier
func (e *ResourceCreatedEvent) EventType() string {
	return "resource.created"
}

// AggregateID returns the ID of the aggregate that generated this event
func (e *ResourceCreatedEvent) AggregateID() string {
	return e.resourceID
}

// Version returns the version of the aggregate when this event occurred
func (e *ResourceCreatedEvent) Version() int {
	return e.version
}

// OccurredAt returns the timestamp when this event occurred
func (e *ResourceCreatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// SetVersion sets the event version (called by Entity when adding event)
func (e *ResourceCreatedEvent) SetVersion(version int) {
	e.version = version
}

// Data returns the resource data
func (e *ResourceCreatedEvent) Data() string {
	return e.data
}

// ContentType returns the content type
func (e *ResourceCreatedEvent) ContentType() string {
	return e.contentType
}

// ExtractedID returns the ID extracted from the resource data
func (e *ResourceCreatedEvent) ExtractedID() string {
	return e.extractedID
}

// ResourceURIAssignedEvent is emitted when a URI is assigned to a resource
type ResourceURIAssignedEvent struct {
	resourceID string
	uri        string
	occurredAt time.Time
	version    int
}

// NewResourceURIAssignedEvent creates a new ResourceURIAssignedEvent
func NewResourceURIAssignedEvent(resourceID, uri string) *ResourceURIAssignedEvent {
	return &ResourceURIAssignedEvent{
		resourceID: resourceID,
		uri:        uri,
		occurredAt: time.Now(),
		version:    1,
	}
}

// EventType returns the event type identifier
func (e *ResourceURIAssignedEvent) EventType() string {
	return "resource.uri_assigned"
}

// AggregateID returns the ID of the aggregate that generated this event
func (e *ResourceURIAssignedEvent) AggregateID() string {
	return e.resourceID
}

// Version returns the version of the aggregate when this event occurred
func (e *ResourceURIAssignedEvent) Version() int {
	return e.version
}

// OccurredAt returns the timestamp when this event occurred
func (e *ResourceURIAssignedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// SetVersion sets the event version (called by Entity when adding event)
func (e *ResourceURIAssignedEvent) SetVersion(version int) {
	e.version = version
}

// URI returns the assigned URI
func (e *ResourceURIAssignedEvent) URI() string {
	return e.uri
}

// ResourceUpdatedEvent is emitted when a resource is updated
type ResourceUpdatedEvent struct {
	resourceID   string
	previousData string
	newData      string
	contentType  string
	occurredAt   time.Time
	version      int
}

// NewResourceUpdatedEvent creates a new ResourceUpdatedEvent
func NewResourceUpdatedEvent(resourceID, previousData, newData, contentType string) *ResourceUpdatedEvent {
	return &ResourceUpdatedEvent{
		resourceID:   resourceID,
		previousData: previousData,
		newData:      newData,
		contentType:  contentType,
		occurredAt:   time.Now(),
		version:      1,
	}
}

// EventType returns the event type identifier
func (e *ResourceUpdatedEvent) EventType() string {
	return "resource.updated"
}

// AggregateID returns the ID of the aggregate that generated this event
func (e *ResourceUpdatedEvent) AggregateID() string {
	return e.resourceID
}

// Version returns the version of the aggregate when this event occurred
func (e *ResourceUpdatedEvent) Version() int {
	return e.version
}

// OccurredAt returns the timestamp when this event occurred
func (e *ResourceUpdatedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// SetVersion sets the event version (called by Entity when adding event)
func (e *ResourceUpdatedEvent) SetVersion(version int) {
	e.version = version
}

// PreviousData returns the previous data
func (e *ResourceUpdatedEvent) PreviousData() string {
	return e.previousData
}

// NewData returns the new data
func (e *ResourceUpdatedEvent) NewData() string {
	return e.newData
}

// ContentType returns the content type
func (e *ResourceUpdatedEvent) ContentType() string {
	return e.contentType
}

// ResourceDeletedEvent is emitted when a resource is deleted
type ResourceDeletedEvent struct {
	resourceID string
	uri        string
	occurredAt time.Time
	version    int
}

// NewResourceDeletedEvent creates a new ResourceDeletedEvent
func NewResourceDeletedEvent(resourceID, uri string) *ResourceDeletedEvent {
	return &ResourceDeletedEvent{
		resourceID: resourceID,
		uri:        uri,
		occurredAt: time.Now(),
		version:    1,
	}
}

// EventType returns the event type identifier
func (e *ResourceDeletedEvent) EventType() string {
	return "resource.deleted"
}

// AggregateID returns the ID of the aggregate that generated this event
func (e *ResourceDeletedEvent) AggregateID() string {
	return e.resourceID
}

// Version returns the version of the aggregate when this event occurred
func (e *ResourceDeletedEvent) Version() int {
	return e.version
}

// OccurredAt returns the timestamp when this event occurred
func (e *ResourceDeletedEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// SetVersion sets the event version (called by Entity when adding event)
func (e *ResourceDeletedEvent) SetVersion(version int) {
	e.version = version
}

// URI returns the URI of the deleted resource
func (e *ResourceDeletedEvent) URI() string {
	return e.uri
}

// Ensure all events implement the domain.Event interface
var _ domain.Event = (*ResourceCreatedEvent)(nil)
var _ domain.Event = (*ResourceURIAssignedEvent)(nil)
var _ domain.Event = (*ResourceUpdatedEvent)(nil)
var _ domain.Event = (*ResourceDeletedEvent)(nil)
