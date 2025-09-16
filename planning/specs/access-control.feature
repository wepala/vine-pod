Feature: Access Control and Permissions
  As a pod owner
  I want to control who can access my resources
  So that I can maintain privacy and share data selectively

  Background:
    Given I have a valid WebID "https://alice.example.com/profile#me"
    And my pod is located at "https://alice.example.com/"
    And there are other users with WebIDs:
      | User | WebID                                 |
      | Bob  | https://bob.example.com/profile#me    |
      | Carol| https://carol.example.com/profile#me  |

  @R027 @public-access-configuration
  Scenario: Configure pod for public access
    Given my pod is currently private
    When I configure my pod to allow public read access
    Then anonymous users should be able to read public resources
    And anonymous users should not be able to write to any resources
    And the public access setting should be reflected in the pod's configuration

  @R027 @private-access-configuration
  Scenario: Configure pod for private access only
    Given my pod currently allows public access
    When I configure my pod to be completely private
    Then anonymous users should receive "401 Unauthorized" for any access attempt
    And only authenticated users with explicit permissions should have access
    And the access control policies should be updated accordingly

  @R032 @R020 @permission-granting @shared-resource-access
  Scenario: Grant read permission to specific user
    Given I have a private resource at "https://alice.example.com/personal/diary.ttl"
    When I grant read permission to Bob's WebID "https://bob.example.com/profile#me"
    Then Bob should be able to read the resource
    And Bob should receive "200 OK" when requesting the resource
    But Bob should receive "403 Forbidden" when attempting to modify it

  @R032 @R021 @R022 @permission-granting @collaborative-creation @collaborative-updates
  Scenario: Grant write permission to specific user
    Given I have a resource at "https://alice.example.com/shared/document.ttl"
    When I grant write permission to Bob's WebID "https://bob.example.com/profile#me"
    Then Bob should be able to read and modify the resource
    And Bob should be able to create new resources in the same container
    But Bob should not be able to delete the resource without control permission

  @R032 @permission-granting
  Scenario: Grant control permission to specific user
    Given I have a resource at "https://alice.example.com/collaborative/project.ttl"
    When I grant control permission to Carol's WebID "https://carol.example.com/profile#me"
    Then Carol should be able to read, write, and delete the resource
    And Carol should be able to modify access control permissions for the resource
    And Carol should have full administrative control over the resource

  @R033 @permission-revocation
  Scenario: Revoke user permissions
    Given Bob has read and write permissions for "https://alice.example.com/shared/notes.ttl"
    When I revoke Bob's write permission but keep read permission
    Then Bob should still be able to read the resource
    But Bob should receive "403 Forbidden" when attempting to modify it
    And the access control list should reflect the updated permissions

  @R033 @permission-revocation
  Scenario: Revoke all permissions from user
    Given Carol has full control permission for "https://alice.example.com/project/data.ttl"
    When I revoke all permissions from Carol
    Then Carol should receive "403 Forbidden" for any access attempt
    And Carol should no longer appear in the resource's access control list

  @access-inheritance @R034
  Scenario: Inherit permissions from parent container
    Given a container "https://alice.example.com/family/" has read permission for Bob
    And inheritance is enabled for the container
    When I create a new resource "https://alice.example.com/family/photo.jpg"
    Then Bob should automatically have read permission for the new resource
    And the resource should inherit the container's access control settings

  @access-inheritance @R034
  Scenario: Override inherited permissions
    Given a container "https://alice.example.com/work/" has read permission for Bob
    And a resource "https://alice.example.com/work/confidential.ttl" exists within it
    When I explicitly deny access to Bob for the confidential resource
    Then Bob should be able to access other resources in the container
    But Bob should receive "403 Forbidden" for the confidential resource
    And the explicit denial should override the inherited permission

  @authentication @solid-oidc
  Scenario: Authenticate using Solid-OIDC
    Given Bob wants to access my protected resource
    When Bob presents valid Solid-OIDC credentials
    Then the system should verify his WebID "https://bob.example.com/profile#me"
    And Bob should be granted access according to his permissions
    And the authentication should be cached for the session duration

  @authentication @invalid-credentials
  Scenario: Reject invalid authentication
    Given someone attempts to access my protected resource
    When they present invalid or expired credentials
    Then they should receive "401 Unauthorized" response
    And they should be redirected to the authentication provider
    And no access should be granted to any protected resources

  @group-permissions @advanced-access-control
  Scenario: Grant permissions to a group
    Given I have created a group "https://alice.example.com/groups/family"
    And the group contains members:
      | https://bob.example.com/profile#me   |
      | https://carol.example.com/profile#me |
    When I grant read permission to the family group for "https://alice.example.com/photos/"
    Then all group members should have read access to the photos container
    And new group members should automatically inherit the permissions
    And removing someone from the group should revoke their access

  @access-control-discovery @R019 @R026 @access-level-awareness
  Scenario: Discover access control information
    Given I have a resource with specific access controls
    When an authenticated user queries the access control information
    Then they should see their own permission level for the resource
    And they should see who else has access (if they have control permission)
    And they should understand what actions they are allowed to perform

  @conditional-access @advanced-access-control
  Scenario: Time-based access control
    Given I want to grant temporary access to a resource
    When I grant read permission to Bob with expiration "2024-12-31T23:59:59Z"
    Then Bob should have access until the expiration date
    And access should be automatically revoked after expiration
    And Bob should receive "403 Forbidden" after the access expires

  @cross-pod-access @federation
  Scenario: Access resources across different pods
    Given Bob has his own pod at "https://bob.example.com/"
    And I grant read permission to Bob for "https://alice.example.com/shared/data.ttl"
    When Bob accesses the resource from his pod's interface
    Then he should be able to read the resource seamlessly
    And the cross-pod authentication should work transparently
    And the access should be logged on both pods

  @access-logging @security-audit
  Scenario: Audit access attempts
    Given I have enabled access logging for sensitive resources
    When users attempt to access "https://alice.example.com/sensitive/document.ttl"
    Then successful and failed access attempts should be logged
    And the logs should include timestamps, WebIDs, and actions attempted
    And I should be able to review access patterns and security events

  @emergency-access @security-features
  Scenario: Emergency access revocation
    Given multiple users have various permissions across my pod
    When I trigger an emergency access revocation
    Then all external access should be immediately revoked
    And only my own WebID should retain access to all resources
    And affected users should be notified of the access changes
    And I should be able to selectively restore access afterwards