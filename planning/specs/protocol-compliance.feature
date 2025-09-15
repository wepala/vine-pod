Feature: Solid Protocol and LDP Compliance
  As a pod implementer
  I want to ensure full compliance with Solid and LDP protocols
  So that my pod interoperates with the broader Solid ecosystem

  Background:
    Given I have a Solid pod implementation
    And the pod supports both Solid Protocol and LDP specifications

  @solid-compliance @http-methods
  Scenario Outline: Support required HTTP methods
    Given a resource exists at "https://alice.example.com/test-resource.ttl"
    When I send a "<method>" request to the resource
    Then the server should handle the request appropriately
    And return a valid HTTP status code
    And include appropriate headers for the method

    Examples:
      | method  |
      | GET     |
      | HEAD    |
      | OPTIONS |
      | PUT     |
      | POST    |
      | PATCH   |
      | DELETE  |

  @solid-compliance @content-negotiation
  Scenario: Support required RDF formats
    Given a resource with RDF content exists
    When I request the resource with different Accept headers
    Then the server must support at least:
      | Format   | Media Type           | Required |
      | Turtle   | text/turtle         | Yes      |
      | JSON-LD  | application/ld+json | Yes      |
      | RDF/XML  | application/rdf+xml | Optional |
    And return the appropriate Content-Type header
    And preserve semantic equivalence across formats

  @ldp-compliance @container-types
  Scenario: Support LDP container types
    Given I want to test LDP container compliance
    When I create containers of different types
    Then the server must support:
      | Container Type    | LDP Type                           |
      | Basic Container   | ldp:BasicContainer                |
      | Direct Container  | ldp:DirectContainer               |
      | Indirect Container| ldp:IndirectContainer             |
    And each container type should behave according to LDP specifications

  @ldp-compliance @link-headers
  Scenario: Include appropriate Link headers
    Given various types of resources exist in the pod
    When I request these resources
    Then each response should include Link headers indicating:
      | Header Purpose      | Example Value                      |
      | Resource Type       | rel="type"                        |
      | LDP Version         | rel="http://www.w3.org/ns/ldp#"   |
      | Constraints         | rel="http://www.w3.org/ns/ldp#constrainedBy" |
    And the headers should accurately reflect the resource characteristics

  @solid-compliance @cors-support
  Scenario: Enable CORS for cross-origin requests
    Given a client application from a different origin
    When the client makes a cross-origin request to my pod
    Then the server should include appropriate CORS headers:
      | Header                       | Purpose                           |
      | Access-Control-Allow-Origin  | Allow cross-origin access         |
      | Access-Control-Allow-Methods | List supported HTTP methods       |
      | Access-Control-Allow-Headers | List allowed request headers      |
    And preflight requests should be handled correctly

  @solid-compliance @etag-support
  Scenario: Support ETags for caching and conflict prevention
    Given a resource exists with current content
    When I request the resource
    Then the response should include an ETag header
    And subsequent requests with If-None-Match should return 304 if unchanged
    And PUT requests with If-Match should validate the ETag
    And conflicting updates should be rejected with 412 Precondition Failed

  @ldp-compliance @containment-triples
  Scenario: Maintain accurate containment information
    Given a container exists with multiple resources
    When I request the container's RDF representation
    Then it should include containment triples using ldp:contains
    And the triples should accurately reflect the actual container contents
    And adding/removing resources should update containment triples automatically

  @solid-compliance @webid-authentication
  Scenario: Support WebID-based authentication
    Given a user with a valid WebID "https://alice.example.com/profile#me"
    When they authenticate using their WebID
    Then the server should verify the WebID document
    And extract the user's public key for verification
    And associate the authenticated identity with the session
    And apply access controls based on the WebID

  @protocol-compliance @uri-handling
  Scenario: Handle URIs according to web standards
    Given various URI formats and patterns
    When resources are created or accessed using these URIs
    Then the server should:
      | Requirement              | Behavior                          |
      | Preserve URI semantics   | Don't modify meaningful URI parts |
      | Handle percent-encoding  | Decode/encode appropriately       |
      | Support fragment IDs     | Handle #fragment identifiers      |
      | Maintain URI consistency | Same resource = same URI          |
    And URI canonicalization should be consistent

  @ldp-compliance @server-managed-properties
  Scenario: Handle server-managed properties correctly
    Given I attempt to modify server-managed properties
    When I send a PUT or PATCH request with server-managed triples
    Then the server should:
      | Property Type        | Action                           |
      | dcterms:modified     | Ignore client value, set server |
      | dcterms:created      | Ignore client value, preserve   |
      | ldp:contains         | Manage automatically             |
      | Custom properties    | Allow client modification        |
    And return appropriate error messages for violations

  @solid-compliance @notifications
  Scenario: Support Solid Notifications Protocol
    Given a client subscribes to changes on a resource
    When the resource is modified
    Then subscribed clients should receive notifications
    And notifications should include:
      | Information          | Details                          |
      | Change type          | Created, updated, deleted        |
      | Resource URI         | Which resource changed           |
      | Timestamp           | When change occurred             |
      | Actor               | Who made the change (if known)   |
    And notification delivery should be reliable

  @protocol-compliance @error-handling
  Scenario Outline: Return appropriate error responses
    Given specific error conditions occur
    When a client makes a request that triggers the error
    Then the server should return the correct HTTP status code
    And include helpful error information

    Examples:
      | Error Condition           | Expected Status | Description           |
      | Resource not found        | 404            | Not Found             |
      | Unauthorized access       | 401            | Unauthorized          |
      | Insufficient permissions  | 403            | Forbidden             |
      | Invalid request format    | 400            | Bad Request           |
      | Unsupported media type    | 415            | Unsupported Media Type|
      | Precondition failed       | 412            | Precondition Failed   |
      | Internal server error     | 500            | Internal Server Error |

  @ldp-compliance @constraints
  Scenario: Enforce LDP constraints
    Given a container with specific constraints
    When clients attempt to violate these constraints
    Then the server should reject invalid operations
    And return 409 Conflict with constraint violation details
    And the Link header should reference the constraint document
    And clients should be able to discover and understand constraints

  @solid-compliance @discovery
  Scenario: Enable protocol and capability discovery
    Given a client wants to discover pod capabilities
    When they request the pod's root or well-known endpoints
    Then they should be able to discover:
      | Capability            | Discovery Method                 |
      | Supported protocols   | Link headers and metadata        |
      | Authentication methods| OIDC configuration endpoints     |
      | Storage quotas       | Custom headers or metadata       |
      | Notification support | Subscription endpoints           |
    And the discovery information should be machine-readable

  @protocol-compliance @interoperability
  Scenario: Ensure interoperability with other Solid implementations
    Given another compliant Solid pod implementation
    When resources are shared or synchronized between pods
    Then both pods should interpret the data identically
    And access control decisions should be consistent
    And protocol negotiations should succeed
    And users should have a seamless cross-pod experience

  @solid-compliance @security-headers
  Scenario: Include appropriate security headers
    Given any HTTP response from the pod
    When security-sensitive operations are performed
    Then responses should include security headers such as:
      | Header                    | Purpose                          |
      | Strict-Transport-Security | Enforce HTTPS usage              |
      | Content-Security-Policy   | Prevent XSS attacks              |
      | X-Frame-Options          | Prevent clickjacking             |
      | X-Content-Type-Options   | Prevent MIME type confusion      |
    And headers should follow security best practices