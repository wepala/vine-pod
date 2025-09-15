# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with the vine-pod service codebase.

## Project Overview & Architectural Vision

**vine-pod** is a Linked Data Platform (LDP) service implementing the Solid Protocol for decentralized data sovereignty. This is NOT a traditional CRUD application - it's a semantic web service with specific architectural constraints.

### Core Architectural Principles

1. **Semantic Web First**: Every resource has semantic meaning, URIs are identities, RDF is the native data model
2. **Event Sourcing by Design**: State is derived from strongly typed events, never stored directly in entities
3. **Protocol Compliance**: Solid Protocol and LDP specifications are non-negotiable constraints
4. **Domain-Driven Architecture**: Business logic lives in domain entities, infrastructure is pluggable
5. **Interface Segregation**: Concrete types hidden behind well-defined interfaces

## Technology Stack

- **Language**: Go 1.24+ (latest version)
- **Architecture**: Domain-Driven Design (DDD) with Clean Architecture
- **Dependencies**:
  - **Kratos**: HTTP framework and middleware (write code in Kratos idiomatic way)
  - **GORM**: ORM for database operations
  - **Fx**: Dependency injection framework
  - **Zap**: Structured logging
  - **Pericarp**: Core DDD library for domain entities and event bus
  - **Moq**: Mock generation for testing
- **Testing**: TDD (Test-Driven Development) and BDD (Behavior-Driven Development)
- **Database**: PostgreSQL (production), SQLite (development/testing)

## Development Commands

```bash
# Core development workflow
make build              # Build the application
make run               # Build and run the application
make test              # Run all tests
make test-cover        # Run tests with coverage report
make lint              # Run golangci-lint
make fmt               # Format code

# Development setup
make dev-setup         # Setup development environment (install tools)
make tidy              # Tidy Go modules

# Docker operations
make docker-build      # Build Docker image
make docker-run        # Build and run in Docker

# Multi-platform builds
make build-all         # Build for multiple platforms (Linux, macOS, Windows)
```

## Domain-Driven Design Structure

The codebase follows strict DDD principles with a clean architecture folder structure:

```
internal/
├── application/           # Application Services & Handlers
│   ├── command/          # Command handlers (CQRS)
│   ├── query/            # Query handlers (CQRS)
│   ├── service/          # Application services
│   └── event/            # Event handlers
├── domain/               # Domain Layer (Business Logic)
│   ├── entity/           # Domain entities/aggregates
│   ├── service/          # Domain services
│   ├── repository/       # Repository interfaces
│   ├── event/            # Domain events
│   └── value/            # Value objects
└── infrastructure/       # Infrastructure Layer
    ├── persistence/      # GORM models and repositories
    │   ├── model/        # GORM models (strict DDD - no business logic)
    │   └── repository/   # Repository implementations
    ├── http/             # HTTP handlers and middleware
    ├── config/           # Configuration
    └── logging/          # Logging setup
```

## Architectural Decision Framework

When implementing features, follow this decision hierarchy:

### 1. Protocol Compliance Check
- **Question**: Does this violate Solid Protocol or LDP specifications?
- **Action**: If yes, reject the approach. Find protocol-compliant alternative.
- **Example**: Never store business logic in GORM models - they're just persistence projections

### 2. Domain-First Design
- **Question**: Can this be modeled as domain behavior with events?
- **Action**: Start with strongly typed events, then build entity methods
- **Pattern**: `new(Resource).FromJSONLD(data).WithURI(uri)` → `ResourceCreatedFromJSONLDEvent` + `ResourceURIAssignedEvent`

### 3. Interface Boundary Decisions
- **Question**: Should this be exposed as an interface?
- **Action**: If it's used by multiple layers or needs testing isolation, create interface
- **Rule**: Domain entities, application services, and repositories = interfaces

### 4. Event vs State Trade-offs
- **Question**: Should this be stored as state or derived from events?
- **Action**: Default to events unless performance requires materialized views
- **Trade-off**: Events = audit trail + flexibility, State = query performance

### 5. Pericarp Integration Points
- **Question**: Is this using pericarp idiomatically?
- **Checklist**:
  - ✅ Entities embed `domain.Entity`
  - ✅ From/With methods return interface types
  - ✅ Events are strongly typed structs
  - ✅ Errors added via `AddError()`, not returned

## Critical Design Constraints

### Non-Negotiable Rules
1. **Never store state in domain entities** - Only events and behavior
2. **IDs come from RDF content** - Extract @id, subject URIs, rdf:about
3. **All operations are atomic** - One entity, one event, one transaction
4. **Test-first development** - No code without tests
5. **Interface-driven design** - Concrete types are implementation details

### Solid Protocol Constraints
- HTTP methods have specific semantics (GET = safe, PUT = idempotent)
- Container membership is managed via LDP containment rules
- Access control via WAC (Web Access Control) or ACP
- Content negotiation required for RDF formats
- URI space must follow LDP conventions

## Test-Driven Development (TDD)

**CRITICAL**: Always generate/update tests FIRST, then implement functionality.

### TDD Workflow:
1. **Red**: Write failing test first
2. **Green**: Write minimal code to pass test
3. **Refactor**: Improve code while keeping tests green

### Testing Structure:
```
internal/
├── domain/
│   ├── entity/
│   │   ├── resource.go
│   │   └── resource_test.go    # Unit tests for domain logic
├── application/
│   ├── service/
│   │   ├── resource_service.go
│   │   └── resource_service_test.go
└── infrastructure/
    └── persistence/
        └── repository/
            ├── resource_repository.go
            └── resource_repository_test.go
```

### Testing Conventions:
- **Test Package Naming**: All test packages MUST end with `_test` (e.g., `package domain_test`, `package service_test`)
- **Mock Generation**: Use `moq` to generate mocks for interfaces:
  ```bash
  //go:generate moq -out repository_mock.go . Repository
  ```
- **Test Files**: Test files must end with `_test.go`

### BDD Testing:
- Use Gherkin scenarios for acceptance tests
- Place feature files in `test/features/`
- Implement step definitions in `test/steps/`

## Implementation Strategy & Patterns

### Feature Development Strategy

**Phase 1: Domain Modeling (Always Start Here)**
1. **Identify the Domain Event**: What business event is occurring?
2. **Design Strongly Typed Event**: Create struct with all necessary data
3. **Model Entity Interface**: What behaviors does this entity need?
4. **Implement From/With Methods**: How is this entity constructed?
5. **Write Tests First**: TDD for all entity behavior

**Phase 2: Application Orchestration**
1. **Command/Query Handlers**: Orchestrate domain operations
2. **Application Services**: Coordinate multiple entities
3. **Event Handlers**: React to domain events for projections

**Phase 3: Infrastructure Implementation**
1. **Repository Interfaces**: Define what persistence needs
2. **GORM Models**: Pure data structures for persistence
3. **HTTP Handlers**: Kratos-idiomatic REST endpoints

### Code Organization Patterns

```bash
# Domain-first development order
internal/domain/entity/resource_test.go         # 1. Write failing test
internal/domain/entity/resource.go             # 2. Make test pass
internal/domain/repository/resource.go         # 3. Define interface
internal/application/command/create_resource_test.go # 4. Application test
internal/application/command/create_resource.go     # 5. Application logic
internal/infrastructure/repository/resource_test.go # 6. Infrastructure test
internal/infrastructure/repository/resource.go     # 7. Infrastructure impl
```

### Common Anti-Patterns to Avoid
1. **State in Entities**: Domain entities storing data instead of emitting events
2. **Generic Events**: Using `map[string]interface{}` instead of typed structs
3. **Logic in Models**: GORM models with business logic
4. **Infrastructure Leakage**: Domain depending on infrastructure types
5. **Missing Interfaces**: Concrete types in application layer

## Logging with Zap

Use structured logging throughout:

```go
logger.Info("resource created",
    zap.String("resource_id", resource.ID),
    zap.String("user_id", userID),
    zap.String("type", resource.Type),
)
```

## Error Handling

Follow Go best practices:
- Return errors, don't panic
- Use custom error types for domain errors
- Log errors at appropriate levels
- Use error wrapping for context

## Configuration

Environment variables (see `.env.example`):
- `SERVER_HOST`, `SERVER_PORT`: Server configuration
- `LOG_LEVEL`: Zap logging level
- `SOLID_DATA_PATH`: Data storage path
- Database connection strings for GORM

## Dependency Management

- Use `go mod tidy` to manage dependencies
- Pin specific versions in `go.mod`
- Use Fx modules to organize dependency injection

## Key Commands for Development

```bash
# TDD workflow
go test -v ./internal/domain/...           # Test domain layer first
go test -v ./internal/application/...      # Test application layer
go test -v ./internal/infrastructure/...   # Test infrastructure layer
go test -v ./...                          # Run all tests

# BDD workflow
go test -v ./test/features/...            # Run BDD scenarios

# Code quality
make lint                                 # Run linter
make fmt                                  # Format code
go vet ./...                             # Static analysis

# Mock generation
go generate ./...                         # Generate mocks using moq

# Coverage analysis
make test-cover                           # Generate coverage report
open coverage.html                        # View coverage in browser
```

## Pericarp Integration

Use the pericarp core library for domain modeling:

### Domain Entities (Interface-Based Atomic Approach):
```go
// Resource interface defines the contract for all LDP resources
type Resource interface {
    // Core pericarp methods
    ID() string
    AddEvent(event domain.Event)
    AddError(err error)
    HasErrors() bool
    GetErrors() []error

    // Resource construction methods
    FromJSONLD(data string) Resource
    FromTurtle(data string) Resource
    FromRDFXML(data string) Resource
    WithURI(uri string) Resource
    WithContainer(containerURI string) Resource

    // Resource operations
    Update(data string, contentType string) Resource
    Delete() Resource

    // Resource metadata
    GetResourceType() ResourceType
    IsContainer() bool
}

// BasicResource - concrete implementation of Resource interface
type BasicResource struct {
    domain.Entity  // Embedded pericarp entity (not pointer)
}

// From/With constructor pattern for atomic instantiation
func (r *BasicResource) FromJSONLD(data string) Resource {
    if data == "" {
        r.AddError(errors.New("empty JSON-LD data"))
        return r
    }

    // Parse JSON-LD to extract the @id field (resource URI)
    var jsonLD map[string]interface{}
    if err := json.Unmarshal([]byte(data), &jsonLD); err != nil {
        r.AddError(fmt.Errorf("invalid JSON-LD: %w", err))
        return r
    }

    // Extract @id from JSON-LD content
    resourceID, ok := jsonLD["@id"].(string)
    if !ok || resourceID == "" {
        r.AddError(errors.New("missing or invalid @id in JSON-LD"))
        return r
    }

    // Initialize entity with ID extracted from the resource content
    if r.ID() == "" {
        r.Entity = domain.NewEntity(resourceID)
    }

    // Use strongly typed events for better type safety and domain semantics
    event := ResourceCreatedFromJSONLDEvent{
        ResourceID:  resourceID,
        Data:        data,
        ContentType: "application/ld+json",
        CreatedAt:   time.Now(),
    }

    r.AddEvent(event)  // Automatically manages version and sequence
    return r
}

func (r *BasicResource) WithURI(uri string) Resource {
    if uri == "" {
        r.AddError(errors.New("invalid URI"))
        return r
    }

    event := ResourceURIAssignedEvent{
        ResourceID: r.ID(),
        URI:        uri,
        AssignedAt: time.Now(),
    }
    r.AddEvent(event)
    return r
}

// Container interface extends Resource with container-specific methods
type Container interface {
    Resource // Embeds all Resource interface methods

    // Container-specific operations
    AddMember(memberURI string) Container
    RemoveMember(memberURI string) Container
    WithMembershipRules(membershipResource, hasMemberRelation, insertedContentRelation string) Container
    GetContainerType() ContainerType

    // Container construction methods
    FromType(containerType string) Container
}

// BasicContainer - concrete implementation of Container interface
type BasicContainer struct {
    domain.Entity  // Embedded pericarp entity
}

// Factory functions for creating different resource types
func NewBasicResource() Resource {
    return &BasicResource{}
}

func NewBasicContainer() Container {
    return &BasicContainer{}
}

// Strongly typed events provide better type safety and domain semantics
type ResourceCreatedFromJSONLDEvent struct {
    ResourceID  string
    Data        string
    ContentType string
    CreatedAt   time.Time
}

type ResourceURIAssignedEvent struct {
    ResourceID string
    URI        string
    AssignedAt time.Time
}

type ContainerMemberAddedEvent struct {
    ContainerID string
    MemberURI   string
    AddedAt     time.Time
}

// Usage patterns:
// Basic resource: NewBasicResource().FromJSONLD("{}").WithURI("/data/resource1")
// Container: NewBasicContainer().FromType("ldp:BasicContainer").WithURI("/data/")
// Mixed: NewBasicContainer().FromJSONLD("{}").AddMember("/data/resource1")
}
```

### Event Bus Setup:
```go
// Use pericarp modules in Fx
var PericarpModule = fx.Module("pericarp",
    pericarp.Module, // Include pericarp's Fx module
    fx.Provide(
        // Your domain services that use pericarp event bus
    ),
)
```

## Kratos Integration

Follow Kratos idioms for HTTP layer:

### HTTP Handlers:
```go
// Kratos service definition
type ResourceService struct {
    pb.UnimplementedResourceServiceServer
    usecase ResourceUsecase
}

// Kratos HTTP handler
func (s *ResourceService) CreateResource(ctx context.Context, req *pb.CreateResourceRequest) (*pb.CreateResourceReply, error) {
    // Use application services, not direct domain access
    result, err := s.usecase.CreateResource(ctx, req)
    if err != nil {
        return nil, err
    }
    return result, nil
}
```

### Middleware:
```go
// Kratos middleware pattern
func LoggingMiddleware(logger log.Logger) middleware.Middleware {
    return logging.Server(logger)
}
```

## Integration Testing

- Use testcontainers for database integration tests
- Mock external services in tests
- Use Fx testing utilities for dependency injection in tests

## Performance Considerations

- Use GORM's batch operations for bulk inserts
- Implement proper database indexing
- Use Zap's structured logging for performance
- Monitor with appropriate metrics

## Security

- Validate all inputs at application layer
- Use proper authentication/authorization with Kratos
- Sanitize data before persistence
- Use HTTPS in production
- Follow OWASP security guidelines