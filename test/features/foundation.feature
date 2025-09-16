Feature: Foundation Infrastructure
  As a developer
  I want to verify that the core infrastructure is working
  So that I can build domain features on top of it

  Background:
    Given the vine-pod service is configured
    And the database connection is available

  @foundation @health
  Scenario: Health endpoint responds correctly
    When I send a GET request to "/health"
    Then the response status should be 200
    And the response should contain "status" as "healthy"
    And the response should contain "service" as "vine-pod"

  @foundation @version
  Scenario: Version endpoint provides build information
    When I send a GET request to "/version"
    Then the response status should be 200
    And the response should contain a "version" field
    And the response should contain a "go_version" field
    And the response should contain a "platform" field

  @foundation @database
  Scenario: Database connection is established
    Given the application is starting up
    When the database module initializes
    Then the database connection should be successful
    And the database should be pingable

  @foundation @fx
  Scenario: Fx dependency injection is working
    Given the Fx container is configured
    When the application starts
    Then all dependencies should be resolved
    And no circular dependencies should exist
    And the application should start without errors

  @foundation @logging
  Scenario: Structured logging is operational
    Given the Zap logger is configured
    When the application logs messages
    Then logs should be in JSON format
    And logs should contain proper fields
    And log levels should be respected