package service

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/deiu/rdf2go"
	"github.com/knakk/rdf"
	"github.com/piprate/json-gold/ld"
)

// StandardRDFValidationService implements the RDFValidationService interface
// using professional RDF libraries for proper validation and parsing
type StandardRDFValidationService struct {
	jsonLDProcessor *ld.JsonLdProcessor
}

// NewStandardRDFValidationService creates a new instance of StandardRDFValidationService
func NewStandardRDFValidationService() RDFValidationService {
	return &StandardRDFValidationService{
		jsonLDProcessor: ld.NewJsonLdProcessor(),
	}
}

// ValidateJSONLD validates JSON-LD data and extracts the @id field
func (s *StandardRDFValidationService) ValidateJSONLD(data string) (resourceID string, err error) {
	// Parse as JSON first
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return "", NewValidationError(FormatJSONLD, "Invalid JSON syntax", err)
	}

	// Validate JSON-LD structure using json-gold
	expanded, err := s.jsonLDProcessor.Expand(jsonData, nil)
	if err != nil {
		return "", NewValidationError(FormatJSONLD, "JSON-LD expansion failed", err)
	}

	// Convert to RDF to ensure it's valid RDF
	_, err = s.jsonLDProcessor.ToRDF(expanded, nil)
	if err != nil {
		return "", NewValidationError(FormatJSONLD, "JSON-LD to RDF conversion failed", err)
	}

	// Extract @id from the original JSON-LD
	if jsonMap, ok := jsonData.(map[string]interface{}); ok {
		if id, exists := jsonMap["@id"]; exists {
			if idStr, ok := id.(string); ok {
				return idStr, nil
			}
		}
	}

	// The expanded form is an array of objects in JSON-LD 1.1
	// Try to extract @id from the expanded form if not found in original

	return "", NewValidationError(FormatJSONLD, "No @id field found in JSON-LD", nil)
}

// ValidateTurtle validates Turtle data and extracts the subject URI
func (s *StandardRDFValidationService) ValidateTurtle(data string) (resourceID string, err error) {
	graph := rdf2go.NewGraph("")
	err = graph.Parse(strings.NewReader(data), "text/turtle")
	if err != nil {
		return "", NewValidationError(FormatTurtle, "Turtle parsing failed", err)
	}

	// Find the first subject URI
	for triple := range graph.IterTriples() {
		if triple.Subject.RawValue() != "" {
			// Return the first subject URI we find
			return triple.Subject.RawValue(), nil
		}
	}

	return "", NewValidationError(FormatTurtle, "No subject URI found in Turtle data", nil)
}

// ValidateRDFXML validates RDF/XML data and extracts the rdf:about URI
func (s *StandardRDFValidationService) ValidateRDFXML(data string) (resourceID string, err error) {
	decoder := rdf.NewTripleDecoder(strings.NewReader(data), rdf.RDFXML)

	// Try to decode the first triple to validate the syntax
	triple, err := decoder.Decode()
	if err != nil && err != io.EOF {
		return "", NewValidationError(FormatRDFXML, "RDF/XML parsing failed", err)
	}

	// If we got a valid triple, extract the subject
	if triple.Subj.Type() == rdf.TermIRI {
		return triple.Subj.String(), nil
	}

	return "", NewValidationError(FormatRDFXML, "No rdf:about URI found in RDF/XML data", nil)
}

// ValidateN3 validates N3 data and extracts the subject URI
func (s *StandardRDFValidationService) ValidateN3(data string) (resourceID string, err error) {
	// N3 is an extension of Turtle, so we can use the same parser
	graph := rdf2go.NewGraph("")
	err = graph.Parse(strings.NewReader(data), "text/n3")
	if err != nil {
		return "", NewValidationError(FormatN3, "N3 parsing failed", err)
	}

	// Find the first subject URI
	for triple := range graph.IterTriples() {
		if triple.Subject.RawValue() != "" {
			return triple.Subject.RawValue(), nil
		}
	}

	return "", NewValidationError(FormatN3, "No subject URI found in N3 data", nil)
}

// ValidateNTriples validates N-Triples data and extracts the subject URI
func (s *StandardRDFValidationService) ValidateNTriples(data string) (resourceID string, err error) {
	decoder := rdf.NewTripleDecoder(strings.NewReader(data), rdf.NTriples)

	// Try to decode the first triple to validate the syntax
	triple, err := decoder.Decode()
	if err != nil && err != io.EOF {
		return "", NewValidationError(FormatNTriples, "N-Triples parsing failed", err)
	}

	// If we got a valid triple, extract the subject
	if triple.Subj.Type() == rdf.TermIRI {
		return triple.Subj.String(), nil
	}

	return "", NewValidationError(FormatNTriples, "No subject URI found in N-Triples data", nil)
}

// ConvertFormat converts RDF data from one format to another
func (s *StandardRDFValidationService) ConvertFormat(data string, fromFormat, toFormat string) (string, error) {
	// First, validate and parse the source data to extract triples
	var triples []rdfTriple
	var err error

	switch fromFormat {
	case string(FormatJSONLD):
		triples, err = s.parseJSONLDToTriples(data)
	case string(FormatTurtle):
		triples, err = s.parseTurtleToTriples(data)
	case string(FormatRDFXML):
		triples, err = s.parseRDFXMLToTriples(data)
	case string(FormatN3):
		triples, err = s.parseN3ToTriples(data)
	case string(FormatNTriples):
		triples, err = s.parseNTriplesToTriples(data)
	default:
		return "", fmt.Errorf("unsupported source format: %s", fromFormat)
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse source data: %w", err)
	}

	// Convert triples to target format
	switch toFormat {
	case string(FormatJSONLD):
		return s.serializeTriplesToJSONLD(triples)
	case string(FormatTurtle):
		return s.serializeTriplesToTurtle(triples)
	case string(FormatRDFXML):
		return s.serializeTriplesToRDFXML(triples)
	case string(FormatN3):
		return s.serializeTriplesToN3(triples)
	case string(FormatNTriples):
		return s.serializeTriplesToNTriples(triples)
	default:
		return "", fmt.Errorf("unsupported target format: %s", toFormat)
	}
}

// SupportedFormats returns a list of supported RDF formats
func (s *StandardRDFValidationService) SupportedFormats() []string {
	return []string{
		string(FormatJSONLD),
		string(FormatTurtle),
		string(FormatRDFXML),
		string(FormatN3),
		string(FormatNTriples),
	}
}

// rdfTriple represents an RDF triple for internal processing
type rdfTriple struct {
	Subject   string
	Predicate string
	Object    string
	IsLiteral bool
}

// Helper methods for parsing different formats to triples

func (s *StandardRDFValidationService) parseJSONLDToTriples(data string) ([]rdfTriple, error) {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		return nil, err
	}

	expanded, err := s.jsonLDProcessor.Expand(jsonData, nil)
	if err != nil {
		return nil, err
	}

	_, err = s.jsonLDProcessor.ToRDF(expanded, nil)
	if err != nil {
		return nil, err
	}

	// For now, create a simple triple from the JSON-LD @id
	var triples []rdfTriple
	if jsonMap, ok := jsonData.(map[string]interface{}); ok {
		if id, exists := jsonMap["@id"]; exists {
			if idStr, ok := id.(string); ok {
				// Create a simple triple - this is a simplified implementation
				triple := rdfTriple{
					Subject:   idStr,
					Predicate: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type",
					Object:    "http://www.w3.org/2000/01/rdf-schema#Resource",
					IsLiteral: false,
				}
				triples = append(triples, triple)
			}
		}
	}

	return triples, nil
}

func (s *StandardRDFValidationService) parseTurtleToTriples(data string) ([]rdfTriple, error) {
	graph := rdf2go.NewGraph("")
	err := graph.Parse(strings.NewReader(data), "text/turtle")
	if err != nil {
		return nil, err
	}

	var triples []rdfTriple
	for triple := range graph.IterTriples() {
		t := rdfTriple{
			Subject:   triple.Subject.RawValue(),
			Predicate: triple.Predicate.RawValue(),
			Object:    triple.Object.RawValue(),
			IsLiteral: strings.HasPrefix(triple.Object.RawValue(), "\""),
		}
		triples = append(triples, t)
	}

	return triples, nil
}

func (s *StandardRDFValidationService) parseRDFXMLToTriples(data string) ([]rdfTriple, error) {
	decoder := rdf.NewTripleDecoder(strings.NewReader(data), rdf.RDFXML)
	var triples []rdfTriple

	for {
		triple, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		t := rdfTriple{
			Subject:   triple.Subj.String(),
			Predicate: triple.Pred.String(),
			Object:    triple.Obj.String(),
			IsLiteral: triple.Obj.Type() == rdf.TermLiteral,
		}
		triples = append(triples, t)
	}

	return triples, nil
}

func (s *StandardRDFValidationService) parseN3ToTriples(data string) ([]rdfTriple, error) {
	// N3 can be parsed as Turtle
	return s.parseTurtleToTriples(data)
}

func (s *StandardRDFValidationService) parseNTriplesToTriples(data string) ([]rdfTriple, error) {
	decoder := rdf.NewTripleDecoder(strings.NewReader(data), rdf.NTriples)
	var triples []rdfTriple

	for {
		triple, err := decoder.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		t := rdfTriple{
			Subject:   triple.Subj.String(),
			Predicate: triple.Pred.String(),
			Object:    triple.Obj.String(),
			IsLiteral: triple.Obj.Type() == rdf.TermLiteral,
		}
		triples = append(triples, t)
	}

	return triples, nil
}

// Helper methods for serializing triples to different formats

func (s *StandardRDFValidationService) serializeTriplesToJSONLD(triples []rdfTriple) (string, error) {
	// Create a simple JSON-LD structure
	jsonLD := make(map[string]interface{})

	if len(triples) > 0 {
		jsonLD["@id"] = triples[0].Subject

		for _, triple := range triples {
			predicate := triple.Predicate
			if triple.IsLiteral {
				jsonLD[predicate] = triple.Object
			} else {
				jsonLD[predicate] = map[string]string{"@id": triple.Object}
			}
		}
	}

	result, err := json.MarshalIndent(jsonLD, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (s *StandardRDFValidationService) serializeTriplesToTurtle(triples []rdfTriple) (string, error) {
	var turtle strings.Builder

	for _, triple := range triples {
		turtle.WriteString(fmt.Sprintf("<%s> <%s> ", triple.Subject, triple.Predicate))

		if triple.IsLiteral {
			turtle.WriteString(fmt.Sprintf(`"%s"`, triple.Object))
		} else {
			turtle.WriteString(fmt.Sprintf("<%s>", triple.Object))
		}

		turtle.WriteString(" .\n")
	}

	return turtle.String(), nil
}

func (s *StandardRDFValidationService) serializeTriplesToRDFXML(triples []rdfTriple) (string, error) {
	var rdfxml strings.Builder

	rdfxml.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	rdfxml.WriteString(`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">` + "\n")

	if len(triples) > 0 {
		subject := triples[0].Subject
		rdfxml.WriteString(fmt.Sprintf(`  <rdf:Description rdf:about="%s">`, subject) + "\n")

		for _, triple := range triples {
			if triple.Subject == subject {
				predicate := triple.Predicate
				rdfxml.WriteString(fmt.Sprintf(`    <%s>`, predicate))

				if triple.IsLiteral {
					rdfxml.WriteString(triple.Object)
				} else {
					rdfxml.WriteString(fmt.Sprintf(` rdf:resource="%s"`, triple.Object))
				}

				rdfxml.WriteString(fmt.Sprintf(`</%s>`, predicate) + "\n")
			}
		}

		rdfxml.WriteString("  </rdf:Description>\n")
	}

	rdfxml.WriteString("</rdf:RDF>\n")

	return rdfxml.String(), nil
}

func (s *StandardRDFValidationService) serializeTriplesToN3(triples []rdfTriple) (string, error) {
	// N3 can be serialized as Turtle
	return s.serializeTriplesToTurtle(triples)
}

func (s *StandardRDFValidationService) serializeTriplesToNTriples(triples []rdfTriple) (string, error) {
	var ntriples strings.Builder

	for _, triple := range triples {
		ntriples.WriteString(fmt.Sprintf("<%s> <%s> ", triple.Subject, triple.Predicate))

		if triple.IsLiteral {
			ntriples.WriteString(fmt.Sprintf(`"%s"`, triple.Object))
		} else {
			ntriples.WriteString(fmt.Sprintf("<%s>", triple.Object))
		}

		ntriples.WriteString(" .\n")
	}

	return ntriples.String(), nil
}
