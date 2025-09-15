Feature: Container Management
  As a pod owner
  I want to organize my resources using containers
  So that I can structure my data hierarchically

  Background:
    Given I have a valid WebID and access credentials
    And my pod is located at "https://alice.example.com/"

  @R008 @container-creation
  Scenario: Create a basic container
    Given I want to create a container for organizing my photos
    When I create a container at "https://alice.example.com/photos/"
    Then the container should be created successfully
    And I should receive a "201 Created" response
    And the container should be an LDP Basic Container
    And the Link header should include 'rel="type"' with "http://www.w3.org/ns/ldp#BasicContainer"

  @R009 @container-deletion
  Scenario: Delete empty container
    Given an empty container exists at "https://alice.example.com/temp/"
    When I delete the container
    Then the container should be removed successfully
    And I should receive a "204 No Content" response
    And subsequent GET requests should return "404 Not Found"

  @R009 @container-deletion-non-empty
  Scenario: Delete non-empty container
    Given a container exists at "https://alice.example.com/photos/"
    And it contains resources:
      | https://alice.example.com/photos/sunset.jpg |
      | https://alice.example.com/photos/beach.jpg  |
    When I attempt to delete the container
    Then I should receive a "409 Conflict" response
    And the container should remain with its contents intact

  @R010 @container-listing
  Scenario: List resources in a container
    Given a container exists at "https://alice.example.com/documents/"
    And it contains the following resources:
      | URI                                           | Type           | Size   |
      | https://alice.example.com/documents/cv.pdf   | Non-RDF Source | 245KB  |
      | https://alice.example.com/documents/notes.ttl | RDF Source     | 1.2KB  |
      | https://alice.example.com/documents/archive/  | Container      | -      |
    When I request the contents of the container
    Then I should see 3 resources
    And each resource should include its URI and type information
    And the response should include LDP containment triples

  @R011 @pseudo-folders-by-type
  Scenario Outline: View resources organized by type
    Given my pod contains resources of different types:
      | URI                                      | Type           | RDF Type        |
      | https://alice.example.com/profile.ttl   | RDF Source     | foaf:Person     |
      | https://alice.example.com/friend.ttl    | RDF Source     | foaf:Person     |
      | https://alice.example.com/photo.jpg     | Non-RDF Source | -               |
      | https://alice.example.com/note.ttl      | RDF Source     | schema:Note     |
    When I request resources filtered by "<type_filter>"
    Then I should see only resources of type "<expected_type>"
    And the resources should be presented as if in a virtual container

    Examples:
      | type_filter                | expected_type   |
      | http://xmlns.com/foaf/0.1/Person | foaf:Person     |
      | http://schema.org/Note     | schema:Note     |
      | ldp:NonRDFSource          | Non-RDF Source  |

  @R012 @container-membership
  Scenario: Add resource to container membership
    Given a container exists at "https://alice.example.com/bookmarks/"
    And a resource exists at "https://alice.example.com/articles/interesting-read.ttl"
    When I add the resource to the container's membership
    Then the container should include the resource in its containment triples
    And I should be able to see the resource when listing container contents

  @R012 @container-membership-removal
  Scenario: Remove resource from container membership
    Given a container exists at "https://alice.example.com/favorites/"
    And it contains a resource "https://alice.example.com/articles/old-favorite.ttl"
    When I remove the resource from the container's membership
    Then the container should no longer list the resource
    And the containment triple should be removed
    But the original resource should still exist at its location

  @container-types @ldp-compliance
  Scenario: Create Direct Container with membership rules
    Given I want to create a Direct Container for managing team members
    When I create a container at "https://alice.example.com/team/" with membership rules:
      | membershipResource    | https://alice.example.com/team/      |
      | hasMemberRelation     | http://schema.org/member             |
    Then the container should be created as an LDP Direct Container
    And the Link header should include 'rel="type"' with "http://www.w3.org/ns/ldp#DirectContainer"
    And the container should enforce the specified membership rules

  @container-types @ldp-compliance
  Scenario: Create Indirect Container with content relation
    Given I want to create an Indirect Container for managing project resources
    When I create a container at "https://alice.example.com/projects/" with rules:
      | membershipResource      | https://alice.example.com/projects/   |
      | hasMemberRelation       | http://schema.org/hasPart             |
      | insertedContentRelation | http://schema.org/contentLocation     |
    Then the container should be created as an LDP Indirect Container
    And the Link header should include 'rel="type"' with "http://www.w3.org/ns/ldp#IndirectContainer"
    And new resources should create membership triples according to the rules

  @R010 @container-navigation
  Scenario: Navigate hierarchical container structure
    Given I have a hierarchical container structure:
      """
      /
      ├── documents/
      │   ├── personal/
      │   │   ├── diary.ttl
      │   │   └── photos/
      │   │       └── vacation.jpg
      │   └── work/
      │       └── reports/
      │           └── quarterly.pdf
      └── settings/
          └── preferences.ttl
      """
    When I navigate to "https://alice.example.com/documents/personal/"
    Then I should see 2 items: "diary.ttl" and "photos/"
    And I should be able to navigate deeper into "photos/"
    And I should see breadcrumb information showing the current path

  @container-metadata @R006
  Scenario: View container metadata and statistics
    Given a container exists at "https://alice.example.com/music/"
    And it contains 150 resources with a total size of 2.5GB
    When I request metadata for the container
    Then I should see metadata including:
      | Field              | Value                    |
      | Container Type     | LDP Basic Container      |
      | Created Date       | 2024-01-15T10:30:00Z    |
      | Last Modified      | 2024-03-20T14:45:00Z    |
      | Resource Count     | 150                      |
      | Total Size         | 2.5GB                    |
      | Child Containers   | 5                        |

  @R013 @container-search
  Scenario: Search within container contents
    Given a container exists at "https://alice.example.com/documents/"
    And it contains resources with various metadata and content
    When I search within the container for resources containing "project proposal"
    Then I should see only resources matching the search criteria
    And the results should maintain the container context
    And I should be able to refine the search with additional filters

  @container-permissions @access-control
  Scenario: Container-level access control
    Given I have a private container at "https://alice.example.com/private/"
    And I want to grant read access to "https://bob.example.com/profile#me"
    When I set access control for the container
    Then Bob should be able to read the container contents
    But Bob should not be able to modify the container
    And the access permissions should apply to all resources within the container