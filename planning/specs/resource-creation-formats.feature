Feature: Resource Creation with RDF Serialization Formats
  As a pod owner or application developer
  I want to create resources using different RDF serialization formats
  So that I can work with my preferred RDF syntax and ensure interoperability

  Background:
    Given I have a valid WebID and access credentials
    And my pod is located at "https://alice.example.com/"
    And the container "https://alice.example.com/data/" exists

  @turtle-creation @R002
  Scenario: Create resource with Turtle format
    Given I have Turtle data:
      """
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix schema: <http://schema.org/> .
      @base <https://alice.example.com/> .

      <profile/card#me>
        a foaf:Person ;
        foaf:name "Alice Smith" ;
        foaf:mbox <mailto:alice@example.com> ;
        schema:jobTitle "Software Developer" .
      """
    When I create a resource with this Turtle data at "https://alice.example.com/profile/card"
    Then the resource should be created successfully
    And I should receive a "201 Created" response
    And the Content-Type should be "text/turtle"
    And the resource should contain the RDF triples from the Turtle data
    And the resource ID should be extracted from the subject URI "https://alice.example.com/profile/card#me"

  @turtle-creation @prefixes
  Scenario: Create resource with Turtle format using multiple prefixes
    Given I have Turtle data with multiple prefixes:
      """
      @prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
      @prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix dcterms: <http://purl.org/dc/terms/> .
      @prefix ldp: <http://www.w3.org/ns/ldp#> .

      <https://alice.example.com/data/project1>
        a ldp:Resource ;
        dcterms:title "My Project" ;
        dcterms:creator <https://alice.example.com/profile/card#me> ;
        dcterms:created "2023-01-15T10:30:00Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> ;
        rdfs:comment "A project about semantic web technologies" .
      """
    When I create a resource with this Turtle data
    Then the resource should be created at "https://alice.example.com/data/project1"
    And all prefixes should be properly resolved
    And the resource should preserve the semantic relationships

  @ntriples-creation @R002
  Scenario: Create resource with N-Triples format
    Given I have N-Triples data:
      """
      <https://alice.example.com/data/document1> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Document> .
      <https://alice.example.com/data/document1> <http://schema.org/name> "Important Document" .
      <https://alice.example.com/data/document1> <http://schema.org/author> <https://alice.example.com/profile/card#me> .
      <https://alice.example.com/data/document1> <http://purl.org/dc/terms/created> "2023-12-01T14:30:00Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> .
      """
    When I create a resource with this N-Triples data
    Then the resource should be created at "https://alice.example.com/data/document1"
    And I should receive a "201 Created" response
    And the Content-Type should be "application/n-triples"
    And each triple should be properly parsed and stored
    And the resource should be queryable using the RDF data

  @rdfxml-creation @R002
  Scenario: Create resource with RDF/XML format
    Given I have RDF/XML data:
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <rdf:RDF
        xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
        xmlns:foaf="http://xmlns.com/foaf/0.1/"
        xmlns:schema="http://schema.org/">

        <foaf:Person rdf:about="https://alice.example.com/contacts/bob">
          <foaf:name>Bob Johnson</foaf:name>
          <foaf:mbox rdf:resource="mailto:bob@example.com"/>
          <schema:knows rdf:resource="https://alice.example.com/profile/card#me"/>
          <foaf:homepage rdf:resource="https://bob.example.com/"/>
        </foaf:Person>
      </rdf:RDF>
      """
    When I create a resource with this RDF/XML data
    Then the resource should be created at "https://alice.example.com/contacts/bob"
    And I should receive a "201 Created" response
    And the Content-Type should be "application/rdf+xml"
    And the rdf:about attribute should determine the resource URI
    And all XML namespaces should be properly handled

  @rdfa-creation @R002
  Scenario: Create resource with RDFa embedded in HTML
    Given I have HTML with RDFa markup:
      """
      <!DOCTYPE html>
      <html>
      <head>
        <title>Alice's Blog Post</title>
      </head>
      <body vocab="http://schema.org/" typeof="BlogPosting">
        <article resource="https://alice.example.com/blog/post1">
          <h1 property="headline">Introduction to Semantic Web</h1>
          <p>Published by
            <span property="author" typeof="Person">
              <span property="name">Alice Smith</span>
            </span>
            on <time property="datePublished" datetime="2023-12-01">December 1, 2023</time>
          </p>
          <div property="articleBody">
            <p>The semantic web is a vision of the web where data has well-defined meaning...</p>
          </div>
          <div property="keywords">semantic web, RDF, linked data</div>
        </article>
      </body>
      </html>
      """
    When I create a resource with this RDFa data
    Then the resource should be created at "https://alice.example.com/blog/post1"
    And the RDF data should be extracted from the RDFa markup
    And the resource should contain structured data about the blog post
    And the HTML content should be preserved alongside the extracted RDF

  @format-conversion @content-negotiation
  Scenario: Create resource in one format and retrieve in another
    Given I have Turtle data:
      """
      @prefix vcard: <http://www.w3.org/2006/vcard/ns#> .

      <https://alice.example.com/contacts/charlie>
        a vcard:Individual ;
        vcard:fn "Charlie Brown" ;
        vcard:email <mailto:charlie@example.com> ;
        vcard:tel <tel:+1-555-123-4567> .
      """
    When I create a resource with this Turtle data
    And I request the resource with Accept header "application/ld+json"
    Then I should receive the resource in JSON-LD format
    And the semantic content should be equivalent
    And the Content-Type should be "application/ld+json"

  @server-assigned-uri @turtle
  Scenario: Create resource with Turtle data using server-assigned URI
    Given I have Turtle data without a specific subject URI:
      """
      @prefix schema: <http://schema.org/> .

      []
        a schema:Event ;
        schema:name "Weekly Team Meeting" ;
        schema:startDate "2023-12-15T14:00:00Z"^^<http://www.w3.org/2001/XMLSchema#dateTime> ;
        schema:location "Conference Room A" .
      """
    When I POST this Turtle data to container "https://alice.example.com/events/"
    Then a new resource should be created with a server-assigned URI
    And the URI should start with "https://alice.example.com/events/"
    And the blank node should be replaced with the assigned URI
    And I should receive a "201 Created" response

  @error-handling @turtle-errors
  Scenario: Handle invalid Turtle syntax
    Given I have invalid Turtle data:
      """
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      # Missing closing quote and period
      <https://alice.example.com/invalid>
        foaf:name "Incomplete Name
      """
    When I attempt to create a resource with this invalid Turtle data
    Then the creation should fail
    And I should receive a "400 Bad Request" response
    And the error message should indicate "Invalid Turtle syntax"
    And the error should specify the line number and nature of the syntax error

  @error-handling @ntriples-errors
  Scenario: Handle invalid N-Triples syntax
    Given I have invalid N-Triples data:
      """
      <https://alice.example.com/invalid> <http://schema.org/name> "Valid triple" .
      invalid-uri <http://schema.org/name> "Invalid URI format" .
      <https://alice.example.com/another> <http://schema.org/name> "Missing period"
      """
    When I attempt to create a resource with this invalid N-Triples data
    Then the creation should fail
    And I should receive a "400 Bad Request" response
    And the error message should indicate "Invalid N-Triples syntax"
    And the error should specify which triple is malformed

  @error-handling @rdfxml-errors
  Scenario: Handle invalid RDF/XML syntax
    Given I have invalid RDF/XML data:
      """
      <?xml version="1.0" encoding="UTF-8"?>
      <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
        <rdf:Description rdf:about="https://alice.example.com/invalid">
          <invalidNamespace:property>Value without proper namespace</invalidNamespace:property>
        </rdf:Description>
      </rdf:RDF>
      """
    When I attempt to create a resource with this invalid RDF/XML data
    Then the creation should fail
    And I should receive a "400 Bad Request" response
    And the error message should indicate "Invalid RDF/XML: undefined namespace"

  @error-handling @rdfa-errors
  Scenario: Handle invalid RDFa markup
    Given I have HTML with invalid RDFa:
      """
      <!DOCTYPE html>
      <html>
      <body vocab="http://schema.org/">
        <div typeof="Person" resource="https://alice.example.com/person1">
          <span property="invalidProperty">Invalid property not in schema.org</span>
          <span property="name">Valid Name</span>
        </div>
      </body>
      </html>
      """
    When I attempt to create a resource with this invalid RDFa
    Then the creation should succeed with warnings
    And I should receive a "201 Created" response
    And the response should include warnings about unknown properties
    And only valid RDF triples should be stored

  @data-integrity @round-trip
  Scenario: Verify data integrity across format conversions
    Given I have complex Turtle data with various RDF constructs:
      """
      @prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix schema: <http://schema.org/> .

      <https://alice.example.com/complex/data>
        a foaf:Person, schema:Person ;
        foaf:name "Alice Smith"@en, "Alice Schmidt"@de ;
        schema:birthDate "1990-05-15"^^<http://www.w3.org/2001/XMLSchema#date> ;
        foaf:knows ( <https://bob.example.com/> <https://charlie.example.com/> ) ;
        schema:address [
          a schema:PostalAddress ;
          schema:streetAddress "123 Main St" ;
          schema:addressLocality "Anytown" ;
          schema:postalCode "12345"
        ] .
      """
    When I create a resource with this Turtle data
    And I retrieve the resource in N-Triples format
    And I retrieve the resource in RDF/XML format
    And I retrieve the resource in JSON-LD format
    Then all retrieved formats should contain semantically equivalent data
    And language tags should be preserved
    And datatype annotations should be maintained
    And blank nodes should be consistently handled
    And RDF lists should be properly represented in each format

  @large-data @performance
  Scenario: Handle large RDF datasets efficiently
    Given I have a large Turtle file with 10,000 triples
    When I create a resource with this large dataset
    Then the creation should complete within 30 seconds
    And I should receive a "201 Created" response
    And the resource should be queryable
    And memory usage should remain reasonable during processing

  @concurrent-creation @race-conditions
  Scenario: Handle concurrent resource creation with different formats
    Given multiple clients are creating resources simultaneously
    When client A creates a resource with Turtle data
    And client B creates a different resource with JSON-LD data at the same time
    And client C creates a third resource with RDF/XML data concurrently
    Then all three resources should be created successfully
    And each resource should maintain its format-specific characteristics
    And no data corruption should occur
    And all clients should receive appropriate success responses