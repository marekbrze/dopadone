---
id: TASK-33
title: Create comprehensive architecture documentation
status: Done
assignee: []
created_date: '2026-03-04 17:00'
updated_date: '2026-03-05 20:29'
labels:
  - documentation
  - architecture
dependencies:
  - TASK-25
  - TASK-26
  - TASK-27
  - TASK-28
  - TASK-29
references:
  - internal/domain/task.go
  - internal/domain/project.go
  - internal/domain/value_objects.go
  - internal/service/interfaces.go
  - internal/service/task_service.go
  - internal/converter/converter.go
  - internal/db/querier.go
  - cmd/dopa/tasks.go
  - cmd/dopa/main.go
  - internal/tui/commands.go
  - docs/TUI.md
documentation:
  - 'https://www.domainlanguage.com/ddd/'
  - 'https://github.com/charmbracelet/bubbletea'
  - 'https://go.dev/blog/error-handling-and-go'
  - 'https://go.dev/blog/table-driven-tests'
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive architecture documentation for AI assistants and future developers. Document the layered architecture (Domain → Service → Repository), design decisions, patterns, and testing strategies. Include real code examples from the codebase to illustrate patterns. The documentation should enable AI assistants to understand the project structure and follow established patterns when implementing new features.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create docs/architecture.md with high-level layer diagram (Domain, Service, Repository, Converter, CLI/TUI layers) showing data flow and dependencies
- [x] #2 Document Domain-Driven Design patterns: entities (Task, Project, Area, Subarea), value objects (TaskStatus, Priority, Color), domain validation, factory methods (NewTask, NewProject)
- [x] #3 Document Service Layer Architecture: service interfaces pattern, dependency injection approach, context-first design, business logic encapsulation, error handling patterns
- [x] #4 Document Bubble Tea TUI Architecture: Elm pattern (Model-Update-View), command system, message types, state management, key handling, modal system
- [x] #5 Document Testing Strategy: table-driven tests, interface mocking, unit test patterns, benchmark patterns, golden files, test helpers
- [x] #6 Include real code examples from codebase for each pattern (domain validation, service creation, TUI commands, table-driven tests)
- [x] #7 Create high-level ASCII diagrams showing layer interactions and data flow between components
- [x] #8 Update README.md with architecture overview section linking to docs/architecture.md
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Task-33: Architecture Documentation

## Executive Summary

**Deliverable**: Modular documentation suite with START HERE guide + 8 chapter documents
**Approach**: Create navigation hub + focused chapter documents for better accessibility
**Estimated Time**: 15-20 hours (6-10 work sessions)
**Decision**: SPLIT into multiple documents for easier navigation and maintenance

## Why Split Into Multiple Documents

1. **Better Navigation**: Quick access to specific topics without scrolling
2. **Easier Maintenance**: Update one chapter without touching others
3. **Faster Loading**: Load only what you need (especially for AI assistants)
4. **Clearer Organization**: Each document has single focus
5. **Better for AI**: Smaller, focused documents easier to process
6. **Progressive Learning**: START HERE → pick your chapter

---

## Document Structure

```
docs/
├── START_HERE.md                    # Entry point, quick start, navigation
├── architecture/
│   ├── README.md                    # Architecture overview (links to chapters)
│   ├── 01-overview.md              # High-level architecture & layers
│   ├── 02-domain-layer.md          # Domain-Driven Design patterns
│   ├── 03-service-layer.md         # Service interfaces & business logic
│   ├── 04-converter-layer.md       # Type transformations
│   ├── 05-repository-layer.md      # Data access & SQL queries
│   ├── 06-cli-layer.md             # Cobra commands & CLI patterns
│   ├── 07-testing-strategy.md      # Testing patterns & examples
│   └── 08-deployment.md            # Build, release, deployment (optional)
├── TUI.md                          # Already exists - TUI documentation
└── images/                          # Diagrams and images (if needed)
```

---

## Task Breakdown: Create 9 Documents

### Subtask 33-A: START_HERE.md (Entry Point) ⏱️ 1-2 hours
**Priority**: HIGHEST (first document users/AI see)
**Dependencies**: None
**Blocks**: All other documents (provides navigation structure)

**Objective**: Create welcoming entry point with quick start and navigation

**Content Structure**:
```markdown
# Dopadone - Start Here

Welcome to Dopadone! This guide helps you quickly understand the project architecture and find the documentation you need.

## Quick Start

**For AI Assistants**: Start with [Architecture Overview](architecture/01-overview.md)

**For New Developers**: Follow the learning path below

**For Quick Reference**: Jump to specific chapter

## What is Dopadone?

[2-3 paragraphs describing the project]

## Architecture at a Glance

[High-level ASCII diagram showing all layers]

```
Presentation Layer (CLI/TUI)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Domain Layer (Entities & Value Objects)
```

## Learning Paths

### For New Developers
1. [Architecture Overview](architecture/01-overview.md) - Understand the big picture
2. [Domain Layer](architecture/02-domain-layer.md) - Learn domain models
3. [Service Layer](architecture/03-service-layer.md) - Understand business logic
4. [Testing Strategy](architecture/07-testing-strategy.md) - Write great tests
5. [TUI Documentation](TUI.md) - Explore terminal UI

### For AI Assistants
1. [Architecture Overview](architecture/01-overview.md) - System structure
2. [Domain Layer](architecture/02-domain-layer.md) - Domain patterns
3. [Service Layer](architecture/03-service-layer.md) - Service patterns
4. [Testing Strategy](architecture/07-testing-strategy.md) - Testing patterns

### For Specific Topics
- **Adding a new entity**: [Domain Layer](architecture/02-domain-layer.md)
- **Creating a service**: [Service Layer](architecture/03-service-layer.md)
- **Writing tests**: [Testing Strategy](architecture/07-testing-strategy.md)
- **Building TUI features**: [TUI Documentation](TUI.md)
- **CLI commands**: [CLI Layer](architecture/06-cli-layer.md)

## Documentation Index

### Architecture Documentation
- [01 - Architecture Overview](architecture/01-overview.md) - High-level structure
- [02 - Domain Layer](architecture/02-domain-layer.md) - DDD patterns & entities
- [03 - Service Layer](architecture/03-service-layer.md) - Business logic
- [04 - Converter Layer](architecture/04-converter-layer.md) - Type transformations
- [05 - Repository Layer](architecture/05-repository-layer.md) - Data access
- [06 - CLI Layer](architecture/06-cli-layer.md) - Command-line interface
- [07 - Testing Strategy](architecture/07-testing-strategy.md) - Test patterns

### Other Documentation
- [TUI Documentation](TUI.md) - Terminal UI architecture & components
- [API Reference](api/README.md) - API documentation (if exists)
- [Contributing Guide](../CONTRIBUTING.md) - How to contribute

## Key Concepts

### Layered Architecture
Brief explanation with link to [01-overview.md](architecture/01-overview.md)

### Domain-Driven Design
Brief explanation with link to [02-domain-layer.md](architecture/02-domain-layer.md)

### Dependency Injection
Brief explanation with link to [03-service-layer.md](architecture/03-service-layer.md)

### Testing Philosophy
Brief explanation with link to [07-testing-strategy.md](architecture/07-testing-strategy.md)

## Technology Stack

- **Language**: Go 1.21+
- **Database**: SQLite with sqlc
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea
- **Testing**: Go testing package, table-driven tests

## Getting Help

- **Issues**: [GitHub Issues](link)
- **Discussions**: [GitHub Discussions](link)
- **Documentation**: You're reading it! Start with [Architecture Overview](architecture/01-overview.md)

---

**Next Step**: Start with [Architecture Overview](architecture/01-overview.md) to understand the system structure.
```

**Deliverables**:
- [ ] Create docs/START_HERE.md
- [ ] Add quick start section
- [ ] Add learning paths for different audiences
- [ ] Add documentation index with descriptions
- [ ] Add high-level architecture diagram
- [ ] Add technology stack overview

**Testing**:
- [ ] All internal links work
- [ ] Learning paths are logical
- [ ] New developer can find what they need in < 2 minutes

---

### Subtask 33-B: architecture/README.md (Navigation Hub) ⏱️ 1 hour
**Priority**: HIGH (provides architecture section navigation)
**Dependencies**: None (can be created in parallel with 33-A)
**Blocks**: 33-C through 33-H

**Objective**: Create architecture section table of contents

**Content Structure**:
```markdown
# Architecture Documentation

This directory contains comprehensive documentation of Dopadone's architecture.

## Overview

Dopadone follows a **layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│  ┌────────────────┐  ┌────────────────┐│
│  │  CLI (Cobra)   │  │ TUI (BubbleTea)││
│  └────────────────┘  └────────────────┘│
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│          Service Layer                  │
│  Business Logic & Validation            │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│        Repository Layer                 │
│  Data Access (db.Querier)               │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│          Domain Layer                   │
│  Entities, Value Objects, Validation    │
└─────────────────────────────────────────┘
```

## Chapters

### [01 - Architecture Overview](01-overview.md)
**What you'll learn**: High-level architecture, layer responsibilities, data flow

**Topics covered**:
- Layered architecture principles
- Dependency flow and inversion
- Module organization (internal/ structure)
- Service container pattern
- Dependency injection approach

**Time to read**: 10-15 minutes

---

### [02 - Domain Layer](02-domain-layer.md)
**What you'll learn**: Domain-Driven Design patterns, entity design, validation

**Topics covered**:
- Entity design (Task, Project, Area, Subarea)
- Value objects (Status, Priority, Color, Duration)
- Factory methods (NewTask, NewProject)
- Domain validation patterns
- Business invariants

**Time to read**: 15-20 minutes

---

### [03 - Service Layer](03-service-layer.md)
**What you'll learn**: Service interfaces, business logic, error handling

**Topics covered**:
- Service interface pattern
- Dependency injection in services
- Context-first design
- Business logic encapsulation
- Error handling and wrapping
- Service method examples

**Time to read**: 20-25 minutes

---

### [04 - Converter Layer](04-converter-layer.md)
**What you'll learn**: Type transformations between layers

**Topics covered**:
- DB types to Domain types
- Null handling patterns
- Bidirectional conversions
- Converter testing strategies

**Time to read**: 10-15 minutes

---

### [05 - Repository Layer](05-repository-layer.md)
**What you'll learn**: Data access patterns, SQL queries, transactions

**Topics covered**:
- db.Querier interface
- sqlc code generation
- Query patterns
- Transaction handling
- Soft delete implementation

**Time to read**: 15-20 minutes

---

### [06 - CLI Layer](06-cli-layer.md)
**What you'll learn**: Cobra commands, CLI patterns, output formatting

**Topics covered**:
- Command structure (projects, tasks, areas)
- Service injection in CLI
- Flag parsing and validation
- Output formatting (table/JSON/YAML)
- Error handling in CLI

**Time to read**: 15-20 minutes

---

### [07 - Testing Strategy](07-testing-strategy.md)
**What you'll learn**: Testing patterns, mocking strategies, best practices

**Topics covered**:
- Table-driven tests pattern
- Interface mocking
- Test helpers
- Golden files
- Benchmark patterns
- Test coverage goals

**Time to read**: 20-25 minutes

---

## Reading Order

### For New Developers
```
01-overview.md → 02-domain-layer.md → 03-service-layer.md → 07-testing-strategy.md
```

### For AI Assistants
```
01-overview.md → 02-domain-layer.md → 03-service-layer.md → 07-testing-strategy.md
```

### For Specific Tasks
- **Adding a new entity**: 02-domain-layer.md → 03-service-layer.md
- **Creating a CLI command**: 06-cli-layer.md → 03-service-layer.md
- **Writing tests**: 07-testing-strategy.md
- **Understanding data flow**: 01-overview.md → 04-converter-layer.md → 05-repository-layer.md

## Design Principles

### 1. Layered Architecture
Each layer has a single responsibility and depends only on layers below it.

### 2. Dependency Injection
Services are injected, not created internally, for testability and flexibility.

### 3. Domain-Driven Design
Rich domain model with validation in factory methods.

### 4. Interface Segregation
Small, focused interfaces defined alongside implementations.

### 5. Comprehensive Testing
Table-driven tests with interface mocking for reliability.

## Key Files Reference

Quick reference to key files in each layer:

- **Domain**: `internal/domain/*.go`
- **Services**: `internal/service/*.go`
- **Converters**: `internal/converter/*.go`
- **Repository**: `internal/db/*.go`
- **CLI**: `cmd/dopa/*.go`
- **TUI**: `internal/tui/*.go`

---

**Ready to dive in?** Start with [01 - Architecture Overview](01-overview.md)
```

**Deliverables**:
- [ ] Create docs/architecture/README.md
- [ ] Add chapter descriptions with time estimates
- [ ] Add reading order recommendations
- [ ] Add quick reference to key files

---

### Subtask 33-C: 01-overview.md (High-Level Architecture) ⏱️ 2-3 hours
**Priority**: HIGH (foundation for all other chapters)
**Dependencies**: 33-B (README exists)
**Blocks**: All subsequent chapters (sets the context)

**Objective**: Document high-level architecture, layers, and data flow

**Content Structure**:
```markdown
# Architecture Overview

## Layered Architecture

Dopadone follows a **layered architecture** pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────┐
│                   Presentation Layer                 │
│  ┌──────────────────┐      ┌────────────────────┐  │
│  │   CLI (Cobra)    │      │  TUI (Bubble Tea)  │  │
│  │  cmd/dopa/  │      │   internal/tui/    │  │
│  └────────┬─────────┘      └──────────┬─────────┘  │
└───────────┼───────────────────────────┼────────────┘
            │                           │
            └───────────┬───────────────┘
                        ▼
┌─────────────────────────────────────────────────────┐
│                   Service Layer                      │
│  ┌─────────────────────────────────────────────┐   │
│  │  AreaService | SubareaService               │   │
│  │  ProjectService | TaskService               │   │
│  │  internal/service/                          │   │
│  └────────────────────┬────────────────────────┘   │
└───────────────────────┼────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────┐
│                  Converter Layer                     │
│  ┌─────────────────────────────────────────────┐   │
│  │  DB Types ↔ Domain Types                    │   │
│  │  internal/converter/                        │   │
│  └────────────────────┬────────────────────────┘   │
└───────────────────────┼────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────┐
│                  Repository Layer                    │
│  ┌─────────────────────────────────────────────┐   │
│  │  db.Querier Interface                       │   │
│  │  SQL Queries (sqlc generated)               │   │
│  │  internal/db/                               │   │
│  └────────────────────┬────────────────────────┘   │
└───────────────────────┼────────────────────────────┘
                        ▼
┌─────────────────────────────────────────────────────┐
│                  Domain Layer                        │
│  ┌─────────────────────────────────────────────┐   │
│  │  Entities: Task, Project, Area, Subarea     │   │
│  │  Value Objects: Status, Priority, Color     │   │
│  │  Factory Methods: NewTask, NewProject       │   │
│  │  internal/domain/                           │   │
│  └─────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Domain Layer
**Purpose**: Core business entities and value objects

**Responsibilities**:
- Define entities (Task, Project, Area, Subarea)
- Define value objects (Status, Priority, Color, Duration)
- Enforce business invariants through factory methods
- Provide domain validation

**Key Files**:
- `internal/domain/task.go` - Task entity
- `internal/domain/project.go` - Project entity
- `internal/domain/value_objects.go` - Value objects

**Dependencies**: None (bottom layer)

**See Also**: [Domain Layer](02-domain-layer.md)

---

### 2. Repository Layer
**Purpose**: Data access abstraction

**Responsibilities**:
- Define db.Querier interface
- Execute SQL queries (generated by sqlc)
- Handle database transactions
- Implement soft delete pattern

**Key Files**:
- `internal/db/querier.go` - Repository interface
- `internal/db/tasks.sql.go` - Task queries
- `internal/db/projects.sql.go` - Project queries
- `internal/db/transaction.go` - Transaction handling

**Dependencies**: Domain layer (for type conversions)

**See Also**: [Repository Layer](05-repository-layer.md)

---

### 3. Converter Layer
**Purpose**: Transform types between layers

**Responsibilities**:
- Convert database types to domain types
- Handle NULL values
- Provide bidirectional conversions
- Centralize transformation logic

**Key Files**:
- `internal/converter/converter.go` - Conversion functions

**Dependencies**: Domain layer, DB types

**See Also**: [Converter Layer](04-converter-layer.md)

---

### 4. Service Layer
**Purpose**: Business logic orchestration

**Responsibilities**:
- Implement business rules
- Orchestrate repository calls
- Apply business validation
- Handle errors with context
- Provide service interfaces for testability

**Key Files**:
- `internal/service/interfaces.go` - Service interfaces
- `internal/service/task_service.go` - Task business logic
- `internal/service/project_service.go` - Project business logic

**Dependencies**: Repository layer (via db.Querier interface)

**See Also**: [Service Layer](03-service-layer.md)

---

### 5. Presentation Layer
**Purpose**: User interface (CLI and TUI)

**CLI Responsibilities**:
- Parse command-line arguments
- Call service methods
- Format and display output
- Handle user-friendly errors

**TUI Responsibilities**:
- Interactive terminal interface
- Keyboard/mouse event handling
- Real-time updates
- Modal dialogs and forms

**Key Files**:
- `cmd/dopa/tasks.go` - Task CLI commands
- `cmd/dopa/projects.go` - Project CLI commands
- `internal/tui/app.go` - TUI main application
- `internal/tui/commands.go` - TUI command functions

**Dependencies**: Service layer (via service interfaces)

**See Also**: [CLI Layer](06-cli-layer.md), [TUI Documentation](../TUI.md)

---

## Data Flow

### Request Flow (Create Task Example)

```
User Input (CLI/TUI)
    ↓
1. Presentation Layer
   - Parse input flags/fields
   - Validate UI-specific constraints
   - Create service params
    ↓
2. Service Layer
   - Apply business validation
   - Execute business logic
   - Call repository
    ↓
3. Converter Layer
   - Convert domain params → DB params
    ↓
4. Repository Layer
   - Execute SQL INSERT
   - Return DB type
    ↓
5. Converter Layer
   - Convert DB type → domain type
    ↓
6. Service Layer
   - Return domain entity
    ↓
7. Presentation Layer
   - Format domain entity
   - Display to user
```

### Code Example: Create Task Flow

```go
// 1. CLI Layer (cmd/dopa/tasks.go)
func runTasksCreate(cmd *cobra.Command, args []string) error {
    title, _ := cmd.Flags().GetString("title")
    projectID, _ := cmd.Flags().GetString("project-id")
    
    // Call service
    services := GetServices()
    task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
        ProjectID: projectID,
        Title:     title,
        Status:    domain.TaskStatusTodo,
    })
    
    if err != nil {
        return cli.WrapError(err, "failed to create task")
    }
    
    return output.Write(task)
}

// 2. Service Layer (internal/service/task_service.go)
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    // Business validation using domain factory
    task, err := domain.NewTask(domain.NewTaskParams{
        ProjectID: params.ProjectID,
        Title:     params.Title,
        Status:    params.Status,
    })
    if err != nil {
        return nil, fmt.Errorf("create task: %w", err)
    }
    
    // Convert domain → DB
    dbParams := converter.DomainTaskToDb(task)
    
    // Call repository
    dbTask, err := s.repo.CreateTask(ctx, dbParams)
    if err != nil {
        return nil, fmt.Errorf("create task: %w", err)
    }
    
    // Convert DB → domain
    return converter.DbTaskToDomain(dbTask), nil
}

// 3. Repository Layer (internal/db/tasks.sql.go)
func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
    // Generated by sqlc
    row := q.db.QueryRowContext(ctx, createTask, arg.ID, arg.Title, ...)
    // ...
}
```

---

## Dependency Injection

### Service Container Pattern

All services are created once and injected where needed:

```go
// cmd/dopa/main.go
type ServiceContainer struct {
    Projects  *service.ProjectService
    Tasks     *service.TaskService
    Subareas  *service.SubareaService
    Areas     *service.AreaService
}

func GetServices() (*ServiceContainer, error) {
    db, err := GetDB()
    if err != nil {
        return nil, err
    }
    
    querier := db.New(db)
    
    return &ServiceContainer{
        Projects:  service.NewProjectService(querier),
        Tasks:     service.NewTaskService(querier),
        Subareas:  service.NewSubareaService(querier),
        Areas:     service.NewAreaService(querier),
    }, nil
}
```

### Benefits

1. **Testability**: Inject mocks for testing
2. **Flexibility**: Swap implementations easily
3. **Loose Coupling**: Depend on interfaces, not implementations
4. **Clear Dependencies**: Explicit in constructors

### Injection Examples

**CLI**:
```go
services := GetServices()
task, err := services.Tasks.Create(ctx, params)
```

**TUI**:
```go
func New(
    areaSvc service.AreaServiceInterface,
    subareaSvc service.SubareaServiceInterface,
    projectSvc service.ProjectServiceInterface,
    taskSvc service.TaskServiceInterface,
) tea.Model {
    return Model{
        areaSvc:    areaSvc,
        subareaSvc: subareaSvc,
        projectSvc: projectSvc,
        taskSvc:    taskSvc,
    }
}
```

**Testing**:
```go
mockTaskSvc := &MockTaskService{
    CreateFunc: func(ctx context.Context, params service.CreateTaskParams) (*domain.Task, error) {
        return &domain.Task{ID: "test"}, nil
    },
}

services := &ServiceContainer{
    Tasks: mockTaskSvc,
}
```

---

## Module Organization

### Directory Structure

```
dopa/
├── cmd/
│   └── dopa/          # CLI commands
│       ├── main.go         # Entry point, service container
│       ├── tasks.go        # Task commands
│       ├── projects.go     # Project commands
│       ├── areas.go        # Area commands
│       └── subareas.go     # Subarea commands
│
├── internal/
│   ├── domain/             # Domain layer
│   │   ├── task.go         # Task entity
│   │   ├── project.go      # Project entity
│   │   ├── area.go         # Area entity
│   │   ├── subarea.go      # Subarea entity
│   │   └── value_objects.go # Value objects
│   │
│   ├── db/                 # Repository layer
│   │   ├── db.go           # Database connection
│   │   ├── querier.go      # Repository interface
│   │   ├── models.go       # DB models (generated)
│   │   ├── tasks.sql.go    # Task queries (generated)
│   │   ├── projects.sql.go # Project queries (generated)
│   │   └── queries/        # SQL files (sqlc source)
│   │
│   ├── converter/          # Converter layer
│   │   └── converter.go    # Type conversions
│   │
│   ├── service/            # Service layer
│   │   ├── interfaces.go   # Service interfaces
│   │   ├── task_service.go # Task business logic
│   │   ├── project_service.go
│   │   ├── area_service.go
│   │   └── subarea_service.go
│   │
│   ├── tui/                # TUI layer
│   │   ├── app.go          # Main application
│   │   ├── commands.go     # TUI commands
│   │   └── ...             # Other TUI components
│   │
│   └── cli/                # CLI utilities
│       ├── output/         # Output formatting
│       ├── filter/         # Filter parsing
│       └── errors.go       # CLI error handling
│
├── docs/                   # Documentation
│   ├── START_HERE.md
│   ├── architecture/
│   └── TUI.md
│
└── testdata/               # Test fixtures
```

### Package Dependencies

```
cmd/dopa → internal/service → internal/db → internal/domain
                    ↓
             internal/converter → internal/domain
                    ↓
              internal/tui → internal/service
```

**Rules**:
- `internal/domain` has NO dependencies (bottom layer)
- `internal/db` depends only on `internal/domain`
- `internal/service` depends on `internal/db` and `internal/domain`
- `cmd/dopa` and `internal/tui` depend on `internal/service`

---

## Design Principles

### 1. Separation of Concerns
Each layer has a single, well-defined responsibility.

### 2. Dependency Inversion
High-level modules (services) don't depend on low-level modules (database). Both depend on abstractions (interfaces).

### 3. Interface Segregation
Small, focused interfaces. Services define what they need, not what implementations provide.

### 4. Single Responsibility
Each service, entity, and module has one reason to change.

### 5. Open/Closed
Open for extension (new services, new entities), closed for modification (existing interfaces stable).

---

## Common Patterns

### Repository Pattern
Abstraction over data access:
```go
type Querier interface {
    CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
    GetTaskByID(ctx context.Context, id string) (Task, error)
    // ...
}
```

### Factory Pattern
Create validated entities:
```go
func NewTask(params NewTaskParams) (*Task, error) {
    if params.Title == "" {
        return nil, ErrTaskTitleEmpty
    }
    return &Task{ID: uuid.New().String(), ...}, nil
}
```

### Strategy Pattern (Services)
Different implementations for same interface:
```go
type TaskServiceInterface interface {
    Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)
}

// Production implementation
type TaskService struct { ... }

// Mock implementation for testing
type MockTaskService struct { ... }
```

---

## Next Steps

Now that you understand the high-level architecture:

1. **Learn about domain models**: [Domain Layer](02-domain-layer.md)
2. **Understand business logic**: [Service Layer](03-service-layer.md)
3. **See how data flows**: [Converter Layer](04-converter-layer.md)
4. **Explore testing**: [Testing Strategy](07-testing-strategy.md)

---

**Navigation**: [← Back to Architecture](README.md) | [Next: Domain Layer →](02-domain-layer.md)
```

**Deliverables**:
- [ ] Create docs/architecture/01-overview.md
- [ ] Add layered architecture diagram
- [ ] Document layer responsibilities
- [ ] Document data flow with examples
- [ ] Document dependency injection pattern
- [ ] Document module organization
- [ ] Add design principles section
- [ ] Add common patterns section

**Testing**:
- [ ] Diagrams render correctly
- [ ] Code examples compile
- [ ] Cross-references work
- [ ] Code examples match actual implementation

---

### Subtask 33-D: 02-domain-layer.md (Domain-Driven Design) ⏱️ 3-4 hours
**Priority**: HIGH (core architectural pattern)
**Dependencies**: 33-C (overview exists)
**Blocks**: 33-E (service layer references domain)

**Objective**: Document Domain-Driven Design patterns, entities, value objects, validation

**Content Outline**:
```markdown
# Domain Layer

## Overview
[Explain DDD approach in this project]

## Entities
### Task Entity
- Structure and fields
- Factory method (NewTask)
- Validation rules

### Project Entity
- Structure with hierarchy
- Parent-child relationships

### Area & Subarea Entities
- Organizational hierarchy

## Value Objects
### TaskStatus
- Enum pattern
- Parse validation

### TaskPriority
- Priority levels
- Comparison logic

### Color
- Hex validation
- Parse method

### TaskDuration
- Duration types
- Time parsing

## Factory Methods
### NewTask Pattern
- Validation in factory
- Error handling
- Default values

### NewProject Pattern
- Hierarchy validation
- Business rules

## Domain Validation
### Validation Errors
- Sentinel errors
- Error messages

### Business Invariants
- Date range validation
- Required fields
- Business rules

## Code Examples
[Real examples from internal/domain/]

## Testing Domain Logic
[Patterns for testing domain validation]

## Best Practices
[DDD principles followed]

---

**Navigation**: [← Overview](01-overview.md) | [Next: Service Layer →](03-service-layer.md)
```

**Reference Files**:
- internal/domain/task.go:20-89
- internal/domain/project.go
- internal/domain/value_objects.go
- internal/domain/task_test.go
- internal/domain/project_test.go

---

### Subtask 33-E: 03-service-layer.md (Service Interfaces) ⏱️ 3-4 hours
**Priority**: HIGH (business logic layer)
**Dependencies**: 33-D (domain layer documented)
**Blocks**: 33-F (converter references services)

**Objective**: Document service interface pattern, dependency injection, business logic

**Content Outline**:
```markdown
# Service Layer

## Overview
[Explain service layer role]

## Service Interface Pattern
### Provider Pattern
- Define interfaces alongside implementations
- Compile-time satisfaction checks

### Interface Design
- Small, focused interfaces
- Context-first methods

## Service Interfaces
### AreaServiceInterface
- Method list with docs

### SubareaServiceInterface
- Method list with docs

### ProjectServiceInterface
- Method list with docs
- ListBySubareaRecursive method

### TaskServiceInterface
- Method list with docs

## Dependency Injection
### Service Construction
- Constructor pattern
- Interface injection

### Service Container
- Container pattern
- GetServices() function

## Business Logic Patterns
### Validation
- Business vs domain validation
- Service-level checks

### Error Handling
- Error wrapping
- Domain-specific errors
- Error checking patterns

### Context Usage
- Request cancellation
- Timeout support
- Tracing

## Code Examples
[Real examples from internal/service/]

## Testing Services
[Mocking repository layer]

## Best Practices
[Service layer principles]

---

**Navigation**: [← Domain Layer](02-domain-layer.md) | [Next: Converter Layer →](04-converter-layer.md)
```

**Reference Files**:
- internal/service/interfaces.go
- internal/service/task_service.go
- internal/service/project_service.go
- internal/service/task_service_test.go
- cmd/dopa/main.go (ServiceContainer)

---

### Subtask 33-F: 04-converter-layer.md (Type Conversions) ⏱️ 1-2 hours
**Priority**: MEDIUM (supporting layer)
**Dependencies**: 33-D (domain types), 33-E (service types)
**Blocks**: 33-G (repository uses converters)

**Objective**: Document type transformation patterns between DB and domain

**Content Outline**:
```markdown
# Converter Layer

## Overview
[Explain converter role]

## DB → Domain Conversions
### DbTaskToDomain
- NULL handling
- Value object parsing

### DbProjectToDomain
- Hierarchy fields
- Status conversion

## Domain → DB Conversions
### DomainTaskToDb
- Required fields
- NULL value handling

## Null Handling Patterns
- sql.NullString
- sql.NullTime
- Helper functions

## Code Examples
[Real examples from internal/converter/]

## Testing Converters
[Unit testing conversions]

## Best Practices
[Converter patterns]

---

**Navigation**: [← Service Layer](03-service-layer.md) | [Next: Repository Layer →](05-repository-layer.md)
```

**Reference Files**:
- internal/converter/converter.go
- internal/converter/converter_test.go

---

### Subtask 33-G: 05-repository-layer.md (Data Access) ⏱️ 2-3 hours
**Priority**: MEDIUM (data access layer)
**Dependencies**: 33-F (converter documented)
**Blocks**: 33-H (CLI uses repository indirectly)

**Objective**: Document repository pattern, sqlc, transactions

**Content Outline**:
```markdown
# Repository Layer

## Overview
[Explain repository pattern]

## db.Querier Interface
### Interface Definition
- Method list

### Implementation
- sqlc generated code

## SQL Queries with sqlc
### Query Files
- queries/tasks.sql
- queries/projects.sql

### Code Generation
- sqlc.yaml config
- Generated files

## Transaction Handling
### Transaction Pattern
- Begin/Commit/Rollback
- Context propagation

## Soft Delete Pattern
- DeletedAt field
- Filter queries

## Code Examples
[Real examples from internal/db/]

## Testing Repository
[Integration testing with SQLite]

## Best Practices
[Repository patterns]

---

**Navigation**: [← Converter Layer](04-converter-layer.md) | [Next: CLI Layer →](06-cli-layer.md)
```

**Reference Files**:
- internal/db/querier.go
- internal/db/tasks.sql.go
- internal/db/transaction.go
- sqlc.yaml

---

### Subtask 33-H: 06-cli-layer.md (Cobra Commands) ⏱️ 2-3 hours
**Priority**: MEDIUM (CLI documentation)
**Dependencies**: 33-E (service layer documented)
**Blocks**: None

**Objective**: Document CLI patterns, service injection, output formatting

**Content Outline**:
```markdown
# CLI Layer

## Overview
[Explain CLI architecture]

## Cobra Command Structure
### Root Command
- main.go entry point

### Resource Commands
- tasks.go
- projects.go
- areas.go
- subareas.go

## Service Injection
### Service Container
- GetServices() pattern

### Service Usage
- Calling service methods
- Error handling

## Flag Parsing
### Required Flags
- Validation

### Optional Flags
- Default values

## Output Formatting
### Table Format
- tablewriter usage

### JSON/YAML Format
- json.Marshal
- yaml.Marshal

## Error Handling
### CLI Errors
- cli.WrapError
- User-friendly messages

## Code Examples
[Real examples from cmd/dopa/]

## Testing CLI
[Command testing patterns]

## Best Practices
[CLI patterns]

---

**Navigation**: [← Repository Layer](05-repository-layer.md) | [Next: Testing Strategy →](07-testing-strategy.md)
```

**Reference Files**:
- cmd/dopa/main.go
- cmd/dopa/tasks.go
- cmd/dopa/projects.go
- internal/cli/output/
- cmd/dopa/commands_test.go

---

### Subtask 33-I: 07-testing-strategy.md (Testing Patterns) ⏱️ 3-4 hours
**Priority**: HIGH (critical for quality)
**Dependencies**: 33-C through 33-H (all layers documented)
**Blocks**: None (final chapter)

**Objective**: Document testing patterns, mocking, test helpers

**Content Outline**:
```markdown
# Testing Strategy

## Overview
[Explain testing philosophy]

## Table-Driven Tests
### Pattern
- Test structure
- Subtests with t.Run()

### Examples
[Real test examples]

## Interface Mocking
### Mock Pattern
- Function fields
- Default implementations

### Service Mocks
- MockTaskService
- MockProjectService

### Repository Mocks
- MockQuerier

## Test Helpers
### Helper Functions
- t.Helper() usage
- Cleanup with t.Cleanup()

### Assertion Helpers
- assertNoError
- assertEqual

## Golden Files
### Pattern
- testdata/ directory
- Update flag

### Usage
[Real golden file examples]

## Benchmarks
### Benchmark Pattern
- b.N loop
- b.ResetTimer()

### Memory Benchmarks
- -benchmem flag

## Test Coverage
### Coverage Goals
- 80%+ general
- 100% critical paths

### Running Coverage
- go test -cover
- Coverage reports

## Code Examples
[Real test files]

## Best Practices
[Testing principles]

## Test Organization
### Unit Tests
- *_test.go files

### Integration Tests
- integration_test.go

### Benchmarks
- *_benchmark_test.go

---

**Navigation**: [← CLI Layer](06-cli-layer.md) | [Back to Architecture →](README.md)
```

**Reference Files**:
- internal/service/task_service_test.go
- internal/tui/commands_test.go
- internal/domain/task_test.go
- internal/tui/mocks/
- testdata/

---

### Subtask 33-J: README.md Update ⏱️ 1 hour
**Priority**: HIGH (entry point visibility)
**Dependencies**: All subtasks 33-A through 33-I
**Blocks**: None

**Objective**: Add architecture section to main README

**Content**:
```markdown
## Architecture

Dopadone follows a **layered architecture** with clear separation of concerns:

```
Presentation Layer (CLI/TUI)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Domain Layer (Entities & Value Objects)
```

**Key Design Principles**:
- **Domain-Driven Design**: Rich domain model with factory methods and validation
- **Dependency Injection**: Services injected into CLI/TUI for testability
- **Interface Segregation**: Small, focused interfaces
- **Comprehensive Testing**: Table-driven tests with interface mocking

**Documentation**:
- **[Get Started](docs/START_HERE.md)** - Quick start guide
- **[Architecture Overview](docs/architecture/01-overview.md)** - System structure
- **[TUI Documentation](docs/TUI.md)** - Terminal UI architecture

**Quick Links**:
- [Domain Layer](docs/architecture/02-domain-layer.md) - DDD patterns
- [Service Layer](docs/architecture/03-service-layer.md) - Business logic
- [Testing Strategy](docs/architecture/07-testing-strategy.md) - Test patterns
```

---

## Task Dependencies & Execution Order

### Sequential Dependencies (Must Follow Order)

```
33-A (START_HERE) ──┐
33-B (arch/README) ─┤
                     ↓
                 33-C (overview)
                     ↓
                 33-D (domain)
                     ↓
                 33-E (service)
                     ↓
                 33-F (converter)
                     ↓
                 33-G (repository)
                     ↓
                 33-H (CLI)
                     ↓
                 33-I (testing)
                     ↓
                 33-J (README update)
```

### Parallel Work Opportunities

**Can be developed in parallel**:
- 33-A and 33-B (no dependencies between them)

**Can overlap**:
- 33-D (domain) can start while 33-C (overview) is being reviewed
- 33-F (converter) can start while 33-E (service) is being written
- 33-H (CLI) can start while 33-G (repository) is being written

**Independent work streams**:
- Stream 1: 33-A → 33-C → 33-D → 33-E → 33-I
- Stream 2: 33-B → 33-F → 33-G → 33-H → 33-J

---

## File Checklist

### Documents to Create (9 files)
- [ ] docs/START_HERE.md
- [ ] docs/architecture/README.md
- [ ] docs/architecture/01-overview.md
- [ ] docs/architecture/02-domain-layer.md
- [ ] docs/architecture/03-service-layer.md
- [ ] docs/architecture/04-converter-layer.md
- [ ] docs/architecture/05-repository-layer.md
- [ ] docs/architecture/06-cli-layer.md
- [ ] docs/architecture/07-testing-strategy.md

### Documents to Update (1 file)
- [ ] README.md (add architecture section)

---

## Acceptance Criteria Mapping

| AC # | Subtask(s) | Deliverable |
|------|------------|-------------|
| 1 | 33-C | High-level layer diagram in 01-overview.md |
| 2 | 33-D | Domain-Driven Design patterns in 02-domain-layer.md |
| 3 | 33-E | Service Layer Architecture in 03-service-layer.md |
| 4 | Already exists | Bubble Tea TUI in docs/TUI.md (reference in 33-A) |
| 5 | 33-I | Testing Strategy in 07-testing-strategy.md |
| 6 | All subtasks | Real code examples from codebase |
| 7 | 33-C | ASCII diagrams showing layer interactions |
| 8 | 33-J | README.md architecture overview section |

---

## Definition of Done Mapping

| DoD # | Subtask(s) | Validation |
|-------|------------|------------|
| 1 | 33-C through 33-I | Diagrams render correctly in all chapters |
| 2 | All subtasks | Code examples compile and are accurate |
| 3 | All subtasks | Documentation follows consistent style |
| 4 | 33-J | README.md section is concise (15-20 lines) |
| 5 | 33-A, 33-C | Structured for AI comprehension (clear sections, rationale) |

---

## Skills & References

### Go Development Skills
- **golang-pro**: Context propagation, production patterns
- **golang-patterns**: Idiomatic patterns, error handling, DI
- **golang-testing**: Table-driven tests, mocking, helpers
- **bubbletea**: Elm architecture, command system

### External References
- DDD: https://www.domainlanguage.com/ddd/
- Bubble Tea: https://github.com/charmbracelet/bubbletea
- Go Error Handling: https://go.dev/blog/error-handling-and-go
- Table-Driven Tests: https://go.dev/blog/table-driven-tests

---

## Estimated Timeline

### Single Developer (Sequential)
```
Day 1: 33-A, 33-B, 33-C (4-6 hours)
Day 2: 33-D, 33-E (6-8 hours)
Day 3: 33-F, 33-G, 33-H (5-7 hours)
Day 4: 33-I, 33-J (4-5 hours)
Total: 19-26 hours
```

### Two Developers (Parallel)
```
Developer A: 33-A → 33-C → 33-D → 33-E → 33-I → 33-J
Developer B: 33-B → 33-F → 33-G → 33-H
Total: 15-20 hours (with overlap)
```

---

## Success Criteria

✅ New developer finds documentation in < 2 minutes (via START_HERE.md)
✅ AI assistant can implement features following documented patterns
✅ All diagrams render correctly in GitHub, VS Code, Typora
✅ Code examples compile and match implementation
✅ Consistent style across all documents
✅ All 8 acceptance criteria met
✅ All 5 definition of done items complete

---

**Recommendation**: Execute with two developers in parallel, starting with 33-A and 33-B simultaneously.

**Total Estimated Time**: 15-20 hours (6-10 work sessions)
**Deliverable**: 9 new documents + 1 updated document
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Key Files to Reference:

Domain Layer:
- internal/domain/task.go - Task entity with factory method and validation
- internal/domain/project.go - Project entity with hierarchy
- internal/domain/value_objects.go - Value objects (Status, Priority, Color)

Service Layer:
- internal/service/interfaces.go - All service interfaces with documentation
- internal/service/task_service.go - TaskService implementation example
- internal/service/project_service.go - ProjectService with complex logic

Converter Layer:
- internal/converter/converter.go - DB to Domain conversions

Repository Layer:
- internal/db/querier.go - Repository interface
- internal/db/tasks.sql.go - SQL queries
- internal/db/transaction.go - Transaction handling

CLI Layer:
- cmd/dopa/tasks.go - CLI commands using TaskService
- cmd/dopa/main.go - Service injection pattern

TUI Layer (already documented):
- docs/TUI.md - Comprehensive TUI documentation
- internal/tui/commands.go - Command functions using services
- internal/tui/app.go - Model with service interfaces

Testing:
- internal/domain/task_test.go - Domain validation tests
- internal/service/task_service_test.go - Table-driven service tests
- internal/tui/commands_test.go - TUI command tests with mocks

Architecture Patterns:
- Dependency Injection via service interfaces
- Factory pattern for domain entities (NewTask, NewProject)
- Repository pattern for data access (db.Querier)
- Converter pattern for type transformations
- Elm Architecture for TUI (Model-Update-View)
- Table-driven tests for comprehensive coverage
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Created comprehensive architecture documentation suite with 9 documents:

### New Documents Created
1. **docs/START_HERE.md** - Entry point with quick start guide and learning paths
2. **docs/architecture/README.md** - Architecture navigation hub
3. **docs/architecture/01-overview.md** - High-level architecture, layers, data flow
4. **docs/architecture/02-domain-layer.md** - Domain-Driven Design patterns and validation
5. **docs/architecture/03-service-layer.md** - Service interfaces and business logic
6. **docs/architecture/04-converter-layer.md** - Type transformations between DB and Domain
7. **docs/architecture/05-repository-layer.md** - Data access patterns and SQL queries
8. **docs/architecture/06-cli-layer.md** - Cobra commands and CLI patterns
9. **docs/architecture/07-testing-strategy.md** - Testing patterns and best practices

### Updated Documents
1. **README.md** - Added architecture overview section with links

### Key Features
- Modular structure for easier navigation and maintenance
- Real code examples from the codebase in each chapter
- ASCII diagrams showing layer interactions and data flow
- Learning paths for different audiences (new developers, AI assistants)
- Clear separation of concerns and dependency flow
- Comprehensive coverage of all architectural layers

### Documentation Structure
- Entry point (START_HERE.md) guides users to appropriate content
- Architecture hub (architecture/README.md) provides navigation
- Individual chapters focus on specific layers and patterns
- Cross-references between related topics
- Consistent style following existing docs/TUI.md structure

All acceptance criteria completed.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All architecture diagrams render correctly in markdown viewers
- [x] #2 Code examples compile and are accurate
- [x] #3 Documentation follows existing docs/TUI.md structure and style
- [x] #4 README.md architecture section is concise and links to full docs
- [x] #5 Documentation is structured for AI comprehension (clear sections, examples, decision rationale)
<!-- DOD:END -->
