package service

import "fmt"

//go:generate moq -out rdf_validation_service_mock.go . RDFValidationService

// RDFValidationService defines the interface for validating RDF data formats
// and extracting resource identifiers from different RDF serializations
type RDFValidationService interface {
	// ValidateJSONLD validates JSON-LD data and extracts the @id field
	ValidateJSONLD(data string) (resourceID string, err error)

	// ValidateTurtle validates Turtle data and extracts the subject URI
	ValidateTurtle(data string) (resourceID string, err error)

	// ValidateRDFXML validates RDF/XML data and extracts the rdf:about URI
	ValidateRDFXML(data string) (resourceID string, err error)

	// ValidateN3 validates N3 data and extracts the subject URI
	ValidateN3(data string) (resourceID string, err error)

	// ValidateNTriples validates N-Triples data and extracts the subject URI
	ValidateNTriples(data string) (resourceID string, err error)

	// ConvertFormat converts RDF data from one format to another
	ConvertFormat(data string, fromFormat, toFormat string) (string, error)

	// SupportedFormats returns a list of supported RDF formats
	SupportedFormats() []string
}

// RDFFormat represents supported RDF serialization formats
type RDFFormat string

const (
	FormatJSONLD   RDFFormat = "application/ld+json"
	FormatTurtle   RDFFormat = "text/turtle"
	FormatRDFXML   RDFFormat = "application/rdf+xml"
	FormatN3       RDFFormat = "text/n3"
	FormatNTriples RDFFormat = "application/n-triples"
	FormatRDFJSON  RDFFormat = "application/rdf+json"
)

// ValidationError represents an RDF validation error with format-specific details
type ValidationError struct {
	Format  RDFFormat
	Line    int
	Column  int
	Message string
	Cause   error
}

func (e *ValidationError) Error() string {
	if e.Line > 0 && e.Column > 0 {
		return string(e.Format) + " validation error at line " +
			fmt.Sprintf("%d", e.Line) + ", column " +
			fmt.Sprintf("%d", e.Column) + ": " + e.Message
	}
	return string(e.Format) + " validation error: " + e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a new ValidationError
func NewValidationError(format RDFFormat, message string, cause error) *ValidationError {
	return &ValidationError{
		Format:  format,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationErrorWithPosition creates a new ValidationError with line/column information
func NewValidationErrorWithPosition(format RDFFormat, line, column int, message string, cause error) *ValidationError {
	return &ValidationError{
		Format:  format,
		Line:    line,
		Column:  column,
		Message: message,
		Cause:   cause,
	}
}
