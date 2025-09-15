Feature: Content Consumer Access
  As an external user or application
  I want to access and consume data from pods
  So that I can use shared information in my applications

  Background:
    Given there is a pod at "https://alice.example.com/" with public resources
    And I am an external content consumer

  @R014 @public-resource-discovery
  Scenario: Discover publicly available resources
    Given Alice's pod contains both public and private resources:
      | Resource                                    | Access Level |
      | https://alice.example.com/profile.ttl      | Public       |
      | https://alice.example.com/blog/            | Public       |
      | https://alice.example.com/private/diary.ttl| Private      |
      | https://alice.example.com/photos/vacation/ | Public       |
    When I request a list of publicly available resources
    Then I should see only the public resources
    And I should not see any private resources
    And each resource should include its basic metadata

  @R015 @public-resource-access
  Scenario: Access public resource content
    Given a public resource exists at "https://alice.example.com/profile.ttl"
    When I request the resource content
    Then I should receive the full resource data
    And I should get a "200 OK" response
    And the content should be properly formatted
    And I should not need authentication

  @R016 @format-specific-access
  Scenario Outline: Request resource in specific format
    Given a public RDF resource exists at "https://alice.example.com/data.ttl"
    When I request the resource with Accept header "<accept_header>"
    Then I should receive the resource in the requested format
    And the Content-Type should be "<content_type>"
    And the semantic content should be equivalent across formats

    Examples:
      | accept_header          | content_type           |
      | text/turtle           | text/turtle            |
      | application/ld+json   | application/ld+json    |
      | application/rdf+xml   | application/rdf+xml    |

  @R017 @resource-metadata-access
  Scenario: View resource metadata as content consumer
    Given public resources exist with various metadata
    When I request metadata for a public resource
    Then I should see safe metadata including:
      | Metadata Field    | Included | Reason                        |
      | Creation Date     | Yes      | Public information            |
      | Last Modified     | Yes      | Public information            |
      | Content Type      | Yes      | Necessary for processing      |
      | Size             | Yes      | Public information            |
      | Owner WebID      | No       | Private information           |
      | Access Controls   | No       | Security information          |
    And sensitive metadata should be excluded

  @R018 @linked-data-navigation
  Scenario: Follow linked data connections
    Given a public resource contains links to other resources:
      """
      @prefix foaf: <http://xmlns.com/foaf/0.1/> .
      @prefix schema: <http://schema.org/> .

      <https://alice.example.com/profile.ttl>
        foaf:knows <https://bob.example.com/profile.ttl> ;
        schema:worksFor <https://company.example.com/about> .
      """
    When I follow the links to connected resources
    Then I should be able to access linked resources if they are public
    And I should receive appropriate error responses for private linked resources
    And I should be able to traverse the linked data graph

  @R019 @search-and-filter
  Scenario: Search and filter public resources
    Given Alice's pod contains multiple public resources with different types and content
    When I search for resources containing "machine learning"
    Then I should see only public resources matching the search criteria
    And the results should be ranked by relevance
    And I should be able to apply additional filters such as:
      | Filter Type       | Example                          |
      | Content Type      | Only RDF sources                 |
      | Date Range        | Resources from last month        |
      | Resource Type     | Only schema:Article              |
    And private resources should never appear in search results

  @content-consumer-applications
  Scenario: Integrate with third-party application
    Given I am developing a research application
    And Alice has made her academic papers publicly available
    When my application accesses her published papers
    Then I should be able to:
      | Action                    | Description                      |
      | Fetch paper metadata      | Get titles, abstracts, authors  |
      | Download full content     | Access PDF or HTML versions     |
      | Follow citation links     | Navigate to referenced works    |
      | Aggregate information     | Combine data from multiple pods  |
    And the integration should work seamlessly without authentication

  @R015 @cross-origin-access
  Scenario: Access resources from web applications
    Given I have a web application running on "https://myapp.example.com"
    And Alice's pod allows cross-origin access to public resources
    When my web application requests Alice's public data
    Then the pod should include appropriate CORS headers
    And my application should be able to access the data client-side
    And browser security policies should be satisfied

  @content-consumer-caching
  Scenario: Efficient content caching
    Given I frequently access the same public resources
    When I request resources I've accessed before
    Then the pod should support HTTP caching mechanisms:
      | Mechanism      | Behavior                         |
      | ETag headers   | Enable conditional requests      |
      | Cache-Control  | Specify caching policies         |
      | Last-Modified  | Support time-based validation    |
    And I should receive "304 Not Modified" for unchanged content
    And my application should be able to cache content appropriately

  @R017 @resource-discovery-feeds
  Scenario: Subscribe to public content updates
    Given Alice regularly publishes new blog posts
    When I subscribe to updates from her blog container
    Then I should be notified when new posts are published
    And I should receive information about updated content
    And I should be able to process the updates in my application
    And the notification mechanism should be efficient and reliable

  @content-consumer-aggregation
  Scenario: Aggregate content from multiple pods
    Given I want to create a directory of public profiles
    When I access public profile data from multiple pods:
      | Pod                        | Profile Data Available    |
      | https://alice.example.com/ | Name, bio, interests      |
      | https://bob.example.com/   | Name, skills, projects    |
      | https://carol.example.com/ | Name, publications, CV    |
    Then I should be able to combine the information
    And create a unified directory view
    And respect the individual privacy settings of each pod

  @rate-limiting @respectful-access
  Scenario: Handle rate limiting appropriately
    Given Alice's pod implements rate limiting for public access
    When I make multiple requests to her public resources
    Then the pod should provide rate limit information in headers:
      | Header                | Purpose                          |
      | X-RateLimit-Limit     | Maximum requests per period      |
      | X-RateLimit-Remaining | Requests remaining in period     |
      | X-RateLimit-Reset     | When the rate limit resets       |
    And I should respect the rate limits
    And handle "429 Too Many Requests" responses gracefully

  @content-quality @data-validation
  Scenario: Validate accessed content quality
    Given I access various public resources with different quality levels
    When I process the content in my application
    Then I should validate:
      | Validation Aspect     | Check                            |
      | RDF syntax validity   | Proper turtle/JSON-LD format    |
      | Schema compliance     | Adherence to expected vocabularies|
      | Link integrity        | Referenced resources exist       |
      | Content freshness     | Last-modified dates are reasonable|
    And handle invalid or poor-quality content gracefully
    And provide meaningful error messages to users