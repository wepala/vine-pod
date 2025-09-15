Feature: Resource Management
  As a pod owner
  I want to manage resources in my data pod
  So that I can store and organize my personal data

  Background:
    Given I have a valid WebID and access credentials
    And my pod is located at "https://alice.example.com/"

  @R001 @resource-listing
  Scenario: List available resources in pod
    Given my pod contains the following resources:
      | URI                                      | Type           | Content-Type         |
      | https://alice.example.com/profile.ttl   | RDF Source     | text/turtle          |
      | https://alice.example.com/photo.jpg     | Non-RDF Source | image/jpeg           |
      | https://alice.example.com/data/         | Container      | application/ld+json  |
    When I request a list of resources in my pod
    Then I should see 3 resources
    And each resource should include its URI, type, and metadata

  @R002 @resource-creation
  Scenario: Create new resource with JSON-LD data
    Given I have valid JSON-LD data:
      """
      {
        "@id": "https://alice.example.com/notes/note1",
        "@type": "schema:Note",
        "schema:text": "My first note"
      }
      """
    When I create a new resource with this data
    Then the resource should be created at "https://alice.example.com/notes/note1"
    And the resource should have content type "application/ld+json"
    And I should receive a "201 Created" response
    And the Location header should contain "https://alice.example.com/notes/note1"

  @R002 @intermediate-containers
  Scenario: Create resource with deep path creates intermediate containers
    Given the path "https://alice.example.com/documents/work/projects/" does not exist
    And I have valid JSON-LD data:
      """
      {
        "@id": "https://alice.example.com/documents/work/projects/proposal.ttl",
        "@type": "schema:Document",
        "schema:name": "Project Proposal"
      }
      """
    When I create a new resource with this data
    Then the resource should be created at "https://alice.example.com/documents/work/projects/proposal.ttl"
    And intermediate containers should be automatically created:
      | https://alice.example.com/documents/      |
      | https://alice.example.com/documents/work/ |
      | https://alice.example.com/documents/work/projects/ |
    And each intermediate container should be a valid LDP Basic Container

  @R003 @resource-update
  Scenario: Update existing resource with new JSON-LD data
    Given a resource exists at "https://alice.example.com/profile.ttl"
    And I have updated JSON-LD data:
      """
      {
        "@id": "https://alice.example.com/profile.ttl",
        "@type": "foaf:Person",
        "foaf:name": "Alice Smith",
        "foaf:mbox": "mailto:alice@example.com"
      }
      """
    When I update the resource with this data
    Then the resource should be updated successfully
    And I should receive a "200 OK" response
    And the resource should contain the new data

  @R003 @conditional-update
  Scenario: Update resource with ETag validation
    Given a resource exists at "https://alice.example.com/profile.ttl" with ETag "abc123"
    And I have updated data with ETag "abc123"
    When I update the resource with If-Match header "abc123"
    Then the resource should be updated successfully
    And I should receive a "200 OK" response
    And the response should include a new ETag

  @R003 @stale-update-prevention
  Scenario: Prevent update with stale ETag
    Given a resource exists at "https://alice.example.com/profile.ttl" with ETag "abc123"
    And I have updated data with stale ETag "old456"
    When I attempt to update the resource with If-Match header "old456"
    Then the update should be rejected
    And I should receive a "412 Precondition Failed" response

  @R004 @resource-deletion
  Scenario: Delete existing resource
    Given a resource exists at "https://alice.example.com/old-note.ttl"
    When I delete the resource
    Then the resource should be removed
    And I should receive a "204 No Content" response
    And subsequent GET requests should return "404 Not Found"

  @R005 @external-resource-linking
  Scenario: Link to external resources
    Given I want to reference an external resource "https://bob.example.com/profile.ttl"
    And I have JSON-LD data that references this external resource:
      """
      {
        "@id": "https://alice.example.com/contacts/bob",
        "@type": "foaf:Person",
        "foaf:knows": "https://bob.example.com/profile.ttl"
      }
      """
    When I create this resource in my pod
    Then the resource should be created successfully
    And the link to the external resource should be preserved
    And I should be able to follow the link to the external resource

  @R006 @resource-metadata
  Scenario: View resource metadata
    Given a resource exists at "https://alice.example.com/document.ttl"
    When I request metadata for this resource
    Then I should see metadata including:
      | Field          | Description              |
      | Creation Date  | When resource was made   |
      | Modified Date  | Last modification time   |
      | Content Type   | MIME type of resource    |
      | Size          | Size in bytes            |
      | ETag          | Entity tag for caching   |

  @R007 @content-negotiation
  Scenario Outline: Download resource in different formats
    Given a resource exists at "https://alice.example.com/profile.ttl" with RDF data
    When I request the resource with Accept header "<accept_header>"
    Then I should receive the resource in "<format>" format
    And the Content-Type header should be "<content_type>"
    And I should receive a "200 OK" response

    Examples:
      | accept_header          | format   | content_type           |
      | text/turtle           | Turtle   | text/turtle            |
      | application/ld+json   | JSON-LD  | application/ld+json    |
      | application/rdf+xml   | RDF/XML  | application/rdf+xml    |

  @R007 @content-negotiation-preference
  Scenario: Handle multiple Accept headers with quality values
    Given a resource exists at "https://alice.example.com/data.ttl"
    When I request the resource with Accept header "application/ld+json;q=0.9, text/turtle;q=1.0"
    Then I should receive the resource in Turtle format
    And the Content-Type header should be "text/turtle"

  @R002 @server-assigned-uri
  Scenario: Create resource with server-assigned URI
    Given I want to create a resource in container "https://alice.example.com/notes/"
    And I have JSON-LD data without a specific @id:
      """
      {
        "@type": "schema:Note",
        "schema:text": "Auto-assigned URI note"
      }
      """
    When I POST this data to the container
    Then a new resource should be created with a server-assigned URI
    And the URI should start with "https://alice.example.com/notes/"
    And I should receive a "201 Created" response
    And the Location header should contain the assigned URI