package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wepala/vine-pod/internal/domain/service"
)

func TestStandardRDFValidationService_ValidateJSONLD(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("validates correct JSON-LD with @id", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@id": "https://alice.example.com/notes/note1",
			"@type": "schema:Note",
			"schema:text": "My first note"
		}`

		// Act
		resourceID, err := rdfService.ValidateJSONLD(jsonLD)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://alice.example.com/notes/note1", resourceID)
	})

	t.Run("returns error for invalid JSON syntax", func(t *testing.T) {
		// Arrange
		invalidJSON := `{"@id": "incomplete`

		// Act
		resourceID, err := rdfService.ValidateJSONLD(invalidJSON)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resourceID)
		assert.Contains(t, err.Error(), "Invalid JSON syntax")
	})

	t.Run("returns error for JSON-LD without @id", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@type": "schema:Note",
			"schema:text": "Note without @id"
		}`

		// Act
		resourceID, err := rdfService.ValidateJSONLD(jsonLD)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resourceID)
		assert.Contains(t, err.Error(), "No @id field found")
	})

	t.Run("validates complex JSON-LD with context", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@context": {
				"schema": "http://schema.org/",
				"foaf": "http://xmlns.com/foaf/0.1/"
			},
			"@id": "https://bob.example.com/profile",
			"@type": "foaf:Person",
			"foaf:name": "Bob Smith",
			"schema:jobTitle": "Software Engineer"
		}`

		// Act
		resourceID, err := rdfService.ValidateJSONLD(jsonLD)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://bob.example.com/profile", resourceID)
	})
}

func TestStandardRDFValidationService_ValidateTurtle(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("validates correct Turtle data", func(t *testing.T) {
		// Arrange
		turtle := `@prefix foaf: <http://xmlns.com/foaf/0.1/> .
@prefix schema: <http://schema.org/> .

<https://alice.example.com/profile/card#me>
    a foaf:Person ;
    foaf:name "Alice Smith" ;
    schema:jobTitle "Software Developer" .`

		// Act
		resourceID, err := rdfService.ValidateTurtle(turtle)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://alice.example.com/profile/card#me", resourceID)
	})

	t.Run("returns error for invalid Turtle syntax", func(t *testing.T) {
		// Arrange
		invalidTurtle := `@prefix foaf: <http://xmlns.com/foaf/0.1/> .
<https://alice.example.com/profile/card#me>
    a foaf:Person
    foaf:name "Alice Smith" # Missing semicolon`

		// Act
		resourceID, err := rdfService.ValidateTurtle(invalidTurtle)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resourceID)
		assert.Contains(t, err.Error(), "Turtle parsing failed")
	})

	t.Run("validates minimal Turtle data", func(t *testing.T) {
		// Arrange
		turtle := `<https://example.com/resource1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .`

		// Act
		resourceID, err := rdfService.ValidateTurtle(turtle)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/resource1", resourceID)
	})
}

func TestStandardRDFValidationService_ValidateRDFXML(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("validates correct RDF/XML data", func(t *testing.T) {
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
		resourceID, err := rdfService.ValidateRDFXML(rdfxml)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://alice.example.com/contacts/bob", resourceID)
	})

	t.Run("returns error for invalid RDF/XML syntax", func(t *testing.T) {
		// Arrange
		invalidRDFXML := `<?xml version="1.0" encoding="UTF-8"?>
<rdf:RDF
  xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about="https://example.com/resource">
    <unclosed:tag>
</rdf:RDF>`

		// Act
		resourceID, err := rdfService.ValidateRDFXML(invalidRDFXML)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resourceID)
		assert.Contains(t, err.Error(), "RDF/XML parsing failed")
	})

	t.Run("validates minimal RDF/XML data", func(t *testing.T) {
		// Arrange
		rdfxml := `<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about="https://example.com/resource1">
    <rdf:type rdf:resource="http://xmlns.com/foaf/0.1/Person"/>
  </rdf:Description>
</rdf:RDF>`

		// Act
		resourceID, err := rdfService.ValidateRDFXML(rdfxml)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/resource1", resourceID)
	})
}

func TestStandardRDFValidationService_ValidateNTriples(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("validates correct N-Triples data", func(t *testing.T) {
		// Arrange
		ntriples := `<https://example.com/resource1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .
<https://example.com/resource1> <http://xmlns.com/foaf/0.1/name> "Alice Smith" .`

		// Act
		resourceID, err := rdfService.ValidateNTriples(ntriples)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "https://example.com/resource1", resourceID)
	})

	t.Run("returns error for invalid N-Triples syntax", func(t *testing.T) {
		// Arrange
		invalidNTriples := `<https://example.com/resource1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> # Missing object`

		// Act
		resourceID, err := rdfService.ValidateNTriples(invalidNTriples)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resourceID)
		assert.Contains(t, err.Error(), "N-Triples parsing failed")
	})
}

func TestStandardRDFValidationService_ConvertFormat(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("converts JSON-LD to Turtle", func(t *testing.T) {
		// Arrange
		jsonLD := `{
			"@id": "https://example.com/resource1",
			"@type": "http://xmlns.com/foaf/0.1/Person"
		}`

		// Act
		result, err := rdfService.ConvertFormat(jsonLD, string(service.FormatJSONLD), string(service.FormatTurtle))

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, "https://example.com/resource1")
		// Note: The simplified implementation creates a basic triple
		assert.Contains(t, result, "Resource")
	})

	t.Run("converts Turtle to JSON-LD", func(t *testing.T) {
		// Arrange
		turtle := `<https://example.com/resource1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://xmlns.com/foaf/0.1/Person> .`

		// Act
		result, err := rdfService.ConvertFormat(turtle, string(service.FormatTurtle), string(service.FormatJSONLD))

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Contains(t, result, "https://example.com/resource1")
	})

	t.Run("returns error for unsupported source format", func(t *testing.T) {
		// Arrange
		data := `some data`

		// Act
		result, err := rdfService.ConvertFormat(data, "unsupported/format", string(service.FormatTurtle))

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "unsupported source format")
	})

	t.Run("returns error for unsupported target format", func(t *testing.T) {
		// Arrange
		jsonLD := `{"@id": "https://example.com/resource1"}`

		// Act
		result, err := rdfService.ConvertFormat(jsonLD, string(service.FormatJSONLD), "unsupported/format")

		// Assert
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "unsupported target format")
	})
}

func TestStandardRDFValidationService_SupportedFormats(t *testing.T) {
	rdfService := service.NewStandardRDFValidationService()

	t.Run("returns list of supported formats", func(t *testing.T) {
		// Act
		formats := rdfService.SupportedFormats()

		// Assert
		assert.NotEmpty(t, formats)
		assert.Contains(t, formats, string(service.FormatJSONLD))
		assert.Contains(t, formats, string(service.FormatTurtle))
		assert.Contains(t, formats, string(service.FormatRDFXML))
		assert.Contains(t, formats, string(service.FormatN3))
		assert.Contains(t, formats, string(service.FormatNTriples))
	})
}

func TestValidationError(t *testing.T) {
	t.Run("creates error with basic message", func(t *testing.T) {
		err := service.NewValidationError(service.FormatJSONLD, "test error", nil)

		assert.Equal(t, "application/ld+json validation error: test error", err.Error())
	})

	t.Run("creates error with position information", func(t *testing.T) {
		err := service.NewValidationErrorWithPosition(service.FormatTurtle, 5, 10, "syntax error", nil)

		assert.Contains(t, err.Error(), "line 5, column 10")
		assert.Contains(t, err.Error(), "syntax error")
	})
}
