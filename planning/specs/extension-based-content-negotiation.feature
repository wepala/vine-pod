Feature: Extension-Based Content Negotiation for RDF Resources
  As a pod owner, content consumer, or application developer
  I want to access the same RDF resource in different serialization formats using file extensions
  So that I can work with my preferred RDF format without complex Accept headers

  Background:
    Given I have a valid WebID and access credentials
    And my pod is located at "https://alice.example.com/"
    And a resource exists with RDF data at the base URI "https://alice.example.com/data/profile"

  @extension-negotiation @R016
  Scenario Outline: Access same resource with different extensions
    Given the base resource "https://alice.example.com/data/profile" contains RDF data:
      """
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix schema: <http://schema.org/> .

      <https://alice.example.com/data/profile#me>
        a foaf:Person ;
        foaf:name "Alice Smith" ;
        foaf:mbox <mailto:alice@example.com> ;
        schema:jobTitle "Software Developer" .
      """
    When I request the resource with extension "<extension>" at "https://alice.example.com/data/profile<extension>"
    Then I should receive a "200 OK" response
    And the Content-Type should be "<content_type>"
    And the response should contain the same RDF data in "<format>" format
    And the semantic content should be equivalent to the original

    Examples:
      | extension | format   | content_type           |
      | .ttl      | Turtle   | text/turtle            |
      | .json     | JSON-LD  | application/ld+json    |
      | .jsonld   | JSON-LD  | application/ld+json    |
      | .nt       | N-Triples| application/n-triples  |
      | .rdf      | RDF/XML  | application/rdf+xml    |
      | .xml      | RDF/XML  | application/rdf+xml    |

  @extension-negotiation @default-format
  Scenario: Access resource without extension returns default format
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile" without any extension
    And I do not specify an Accept header
    Then I should receive a "200 OK" response
    And the Content-Type should be "text/turtle" (default format)
    And the response should contain the RDF data in Turtle format

  @extension-negotiation @default-with-accept
  Scenario: Access resource without extension but with Accept header
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile" without any extension
    And I include Accept header "application/ld+json"
    Then I should receive a "200 OK" response
    And the Content-Type should be "application/ld+json"
    And the response should contain the RDF data in JSON-LD format

  @extension-priority @R016
  Scenario: Extension takes precedence over Accept header
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile.ttl"
    And I include Accept header "application/ld+json"
    Then I should receive a "200 OK" response
    And the Content-Type should be "text/turtle" (from extension, not Accept header)
    And the response should contain the RDF data in Turtle format
    And the Accept header should be ignored in favor of the extension

  @extension-creation @R002
  Scenario Outline: Create resource using extension-specific URI
    Given I have RDF data in "<input_format>" format:
      """
      <rdf_data_content>
      """
    When I create a resource at "https://alice.example.com/data/document<extension>"
    Then the resource should be created successfully
    And I should receive a "201 Created" response
    And the resource should be accessible at "https://alice.example.com/data/document"
    And the resource should also be accessible at "https://alice.example.com/data/document<extension>"
    And both URIs should return the same semantic content

    Examples:
      | extension | input_format |
      | .ttl      | Turtle       |
      | .json     | JSON-LD      |
      | .nt       | N-Triples    |
      | .rdf      | RDF/XML      |

  @extension-negotiation @multiple-extensions
  Scenario: Access resource with multiple supported extensions
    Given the base resource "https://alice.example.com/data/contact" contains contact information
    When I request "https://alice.example.com/data/contact.ttl"
    Then I should receive Turtle format
    When I request "https://alice.example.com/data/contact.json"
    Then I should receive JSON-LD format
    When I request "https://alice.example.com/data/contact.nt"
    Then I should receive N-Triples format
    When I request "https://alice.example.com/data/contact.rdf"
    Then I should receive RDF/XML format
    And all responses should contain semantically equivalent data

  @extension-negotiation @complex-data
  Scenario: Handle complex RDF structures with extensions
    Given the base resource "https://alice.example.com/data/project" contains complex RDF:
      """
      @prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix schema: <http://schema.org/> .

      <https://alice.example.com/data/project#main>
        a schema:Project ;
        schema:name "Semantic Web Project"@en, "Projekt Web Semantyczny"@pl ;
        schema:startDate "2023-01-15"^^<http://www.w3.org/2001/XMLSchema#date> ;
        schema:member [
          a foaf:Person ;
          foaf:name "Alice Smith" ;
          schema:role "Lead Developer"
        ] ;
        schema:member [
          a foaf:Person ;
          foaf:name "Bob Johnson" ;
          schema:role "Designer"
        ] .
      """
    When I request "https://alice.example.com/data/project.json"
    Then the JSON-LD response should preserve language tags
    And the JSON-LD response should preserve datatype annotations
    And the JSON-LD response should properly represent blank nodes
    When I request "https://alice.example.com/data/project.nt"
    Then the N-Triples response should contain individual triples
    And all complex structures should be properly serialized

  @extension-error-handling @unsupported-extension
  Scenario: Handle unsupported file extensions
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile.xyz" with an unsupported extension
    Then I should receive a "406 Not Acceptable" response
    And the error message should indicate "Unsupported file extension '.xyz'"
    And the response should include supported extensions in the error details

  @extension-error-handling @non-rdf-extension
  Scenario: Handle non-RDF file extensions for RDF resources
    Given the base resource "https://alice.example.com/data/profile" contains RDF data
    When I request "https://alice.example.com/data/profile.pdf"
    Then I should receive a "406 Not Acceptable" response
    And the error message should indicate "Cannot serve RDF resource in 'pdf' format"
    And the response should suggest supported RDF extensions

  @extension-error-handling @missing-resource
  Scenario: Handle missing resource with extension
    Given no resource exists at "https://alice.example.com/data/nonexistent"
    When I request "https://alice.example.com/data/nonexistent.ttl"
    Then I should receive a "404 Not Found" response
    And the error should indicate the base resource does not exist
    And the extension should not affect the 404 response

  @extension-negotiation @case-insensitive
  Scenario: Handle case-insensitive extensions
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile.TTL"
    Then I should receive a "200 OK" response
    And the Content-Type should be "text/turtle"
    When I request "https://alice.example.com/data/profile.JSON"
    Then I should receive a "200 OK" response
    And the Content-Type should be "application/ld+json"

  @extension-negotiation @redirect-behavior
  Scenario: Handle redirect from extensionless to default extension
    Given the server is configured to redirect to default extensions
    And the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile" with Accept header "*/*"
    Then I should receive a "301 Moved Permanently" or "302 Found" response
    And the Location header should be "https://alice.example.com/data/profile.ttl"
    When I follow the redirect
    Then I should receive a "200 OK" response
    And the Content-Type should be "text/turtle"

  @extension-negotiation @content-length
  Scenario: Verify proper Content-Length headers for different formats
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile.ttl"
    Then the Content-Length header should match the Turtle serialization size
    When I request "https://alice.example.com/data/profile.json"
    Then the Content-Length header should match the JSON-LD serialization size
    And different formats should have different Content-Length values

  @extension-negotiation @caching
  Scenario: Verify proper caching headers for extension-based requests
    Given the base resource "https://alice.example.com/data/profile" exists
    When I request "https://alice.example.com/data/profile.ttl"
    Then the response should include appropriate caching headers
    And the ETag should be format-specific
    When I request "https://alice.example.com/data/profile.json"
    Then the ETag should be different from the Turtle format ETag
    And both responses should support conditional requests

  @extension-negotiation @head-requests
  Scenario: Support HEAD requests with extensions
    Given the base resource "https://alice.example.com/data/profile" exists
    When I send a HEAD request to "https://alice.example.com/data/profile.ttl"
    Then I should receive a "200 OK" response
    And the Content-Type should be "text/turtle"
    And the Content-Length should match the GET response
    And the response body should be empty

  @extension-negotiation @options-requests
  Scenario: Support OPTIONS requests showing available formats
    Given the base resource "https://alice.example.com/data/profile" exists
    When I send an OPTIONS request to "https://alice.example.com/data/profile"
    Then I should receive a "200 OK" response
    And the response should include available content types in headers
    And the response should indicate supported extensions

  @extension-negotiation @concurrent-access
  Scenario: Handle concurrent access with different extensions
    Given the base resource "https://alice.example.com/data/profile" exists
    When multiple clients simultaneously request:
      | Client | URI                                           | Expected Format |
      | A      | https://alice.example.com/data/profile.ttl    | Turtle          |
      | B      | https://alice.example.com/data/profile.json   | JSON-LD         |
      | C      | https://alice.example.com/data/profile.nt     | N-Triples       |
    Then all clients should receive "200 OK" responses
    And each client should receive the correct format
    And there should be no format mixing or corruption

  @extension-negotiation @performance
  Scenario: Efficient format conversion for extension requests
    Given the base resource "https://alice.example.com/data/large-dataset" contains 1000 RDF triples
    When I request "https://alice.example.com/data/large-dataset.ttl"
    Then the response should be received within 5 seconds
    When I request "https://alice.example.com/data/large-dataset.json"
    Then the response should be received within 5 seconds
    And format conversion should not significantly impact response time

  @extension-negotiation @advanced-formats
  Scenario: Support for advanced RDF formats
    Given the base resource "https://alice.example.com/data/graph" contains named graph data
    When I request "https://alice.example.com/data/graph.nq"
    Then I should receive the data in N-Quads format
    And the Content-Type should be "application/n-quads"
    When I request "https://alice.example.com/data/graph.trig"
    Then I should receive the data in TriG format
    And the Content-Type should be "application/trig"
    And named graph information should be preserved