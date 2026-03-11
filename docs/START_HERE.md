# Dopadone - Start Here

Welcome to Dopadone! This guide helps you quickly understand the project architecture and find the documentation you need.

## Quick Start

**For AI Assistants**: Start with [Architecture Overview](architecture/01-overview.md)

**For New Developers**: Follow the learning path below

**For Quick Reference**: Jump to specific chapter

## What is Dopadone?

Dopadone is a lightweight, SQLite-based CLI project management tool designed for developers who prefer staying in the terminal. It provides:

- **Hierarchical Organization**: Areas → Subareas → Projects → Tasks
- **Terminal-First Workflow**: Full CLI and interactive TUI
- **Rich Domain Model**: Domain-Driven Design with validation and business rules
- **Local Storage**: SQLite database for offline-first work

The architecture follows **layered design principles** with clear separation of concerns, making it easy to understand, test, and extend.

## Architecture at a Glance

```
┌─────────────────────────────────────────────────────┐
│                   Presentation Layer                 │
│  ┌──────────────────┐      ┌────────────────────┐  │
│  │   CLI (Cobra)    │      │  TUI (Bubble Tea)  │  │
│  └──────────────────┘      └────────────────────┘  │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│                   Service Layer                      │
│        Business Logic & Validation                   │
│     (AreaService, ProjectService, TaskService)      │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│                   Converter Layer                    │
│          Type Transformations (DB ↔ Domain)          │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│                  Repository Layer                    │
│           Data Access (db.Querier)                   │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│                   Domain Layer                       │
│   Entities, Value Objects, Factory Methods           │
│         (Task, Project, Area, Subarea)              │
└─────────────────────────────────────────────────────┘
```

**Data Flow**: Request → Presentation → Service → Converter → Repository → Domain

## Learning Paths

### For New Developers
1. [Architecture Overview](architecture/01-overview.md) - Understand the big picture
2. [Domain Layer](architecture/02-domain-layer.md) - Learn domain models and validation
3. [Service Layer](architecture/03-service-layer.md) - Understand business logic patterns
4. [Testing Strategy](architecture/07-testing-strategy.md) - Write great tests
5. [TUI Documentation](TUI.md) - Explore terminal UI architecture

### For AI Assistants
1. [Architecture Overview](architecture/01-overview.md) - System structure and layer interactions
2. [Domain Layer](architecture/02-domain-layer.md) - Domain patterns and validation rules
3. [Service Layer](architecture/03-service-layer.md) - Service interfaces and business logic
4. [Testing Strategy](architecture/07-testing-strategy.md) - Testing patterns and examples

### For Specific Tasks

| Task | Start Here |
|------|------------|
| **Adding a new entity** | [Domain Layer](architecture/02-domain-layer.md) |
| **Creating a service** | [Service Layer](architecture/03-service-layer.md) |
| **Writing tests** | [Testing Strategy](architecture/07-testing-strategy.md) |
| **Building TUI features** | [TUI Documentation](TUI.md) |
| **Adding CLI commands** | [CLI Layer](architecture/06-cli-layer.md) |
| **Understanding data flow** | [Converter Layer](architecture/04-converter-layer.md) → [Repository Layer](architecture/05-repository-layer.md) |

## Documentation Index

### Architecture Documentation

| Document | Description | Time |
|----------|-------------|------|
| [01 - Architecture Overview](architecture/01-overview.md) | High-level structure, layers, data flow | 10-15 min |
| [02 - Domain Layer](architecture/02-domain-layer.md) | DDD patterns, entities, validation | 15-20 min |
| [03 - Service Layer](architecture/03-service-layer.md) | Business logic, dependency injection | 20-25 min |
| [04 - Converter Layer](architecture/04-converter-layer.md) | Type transformations between layers | 10-15 min |
| [05 - Repository Layer](architecture/05-repository-layer.md) | Data access, SQL queries, transactions | 15-20 min |
| [06 - CLI Layer](architecture/06-cli-layer.md) | Cobra commands, service injection | 15-20 min |
| [07 - Testing Strategy](architecture/07-testing-strategy.md) | Test patterns, mocking, coverage | 20-25 min |

### Other Documentation

| Document | Description |
|----------|-------------|
| [TUI Documentation](TUI.md) | Terminal UI architecture, components, Elm pattern |
| [Transaction Handling](TRANSACTIONS.md) | Database transactions, atomicity, best practices |
| [Database Modes](DATABASE_MODES.md) | Local SQLite, Turso remote, and embedded replica modes |
| [Turso Setup Guide](TURSO_SETUP.md) | Turso account signup, CLI installation, and credential setup |
| [Turso Migrations](TURSO_MIGRATIONS.md) | Schema migration guide for libSQL/Turso integration |
| [Turso Data Migration](TURSO_DATA_MIGRATION.md) | Step-by-step guide for migrating data from SQLite to Turso |
| [Turso Performance](TURSO_PERFORMANCE.md) | Performance optimization, benchmarks, and best practices for each mode |
| [Turso Troubleshooting](TURSO_TROUBLESHOOTING.md) | Comprehensive error solutions and diagnostic toolkit |
| [Release Process](RELEASE.md) | Versioning, tagging, deployment workflow |
| [CI/CD Pipeline](CI-CD.md) | GitHub Actions workflows, build process, automation |
| [Code Quality and Linting](CODE_QUALITY.md) | Linting standards, error handling patterns, best practices |
| [Rebranding Documentation](REBRANDING.md) | Project rename from ProjectDB to Dopadone |
| [Repository Migration](REPOSITORY_MIGRATION.md) | Migration from placeholder to production repository |

## Key Concepts

### Layered Architecture

Each layer has a **single responsibility** and depends only on layers below it:

- **Domain**: Core business entities with validation (no dependencies)
- **Repository**: Data access abstraction
- **Converter**: Type transformations
- **Service**: Business logic orchestration
- **Presentation**: User interface (CLI/TUI)

→ See [Architecture Overview](architecture/01-overview.md) for details

### Domain-Driven Design

Rich domain model with:

- **Factory methods** (`NewTask`, `NewProject`) for validation
- **Value objects** (`TaskStatus`, `Priority`, `Color`)
- **Domain validation** enforcing business invariants

→ See [Domain Layer](architecture/02-domain-layer.md) for details

### Dependency Injection

Services are **injected**, not created internally:

```go
type ServiceContainer struct {
    Projects  *service.ProjectService
    Tasks     *service.TaskService
    Subareas  *service.SubareaService
    Areas     *service.AreaService
}
```

Benefits: Testability, flexibility, loose coupling

→ See [Service Layer](architecture/03-service-layer.md) for details

### Testing Philosophy

Comprehensive testing with:

- **Table-driven tests** for clarity and coverage
- **Interface mocking** for isolation
- **Test helpers** for common patterns
- **Golden files** for snapshot testing

→ See [Testing Strategy](architecture/07-testing-strategy.md) for details

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.21+ | Primary implementation language |
| **Database** | SQLite (modernc.org/sqlite) | Local storage, zero dependencies |
| **Migrations** | goose v3 | Schema versioning |
| **Query Generation** | sqlc | Type-safe SQL queries |
| **CLI Framework** | Cobra | Command-line interface |
| **TUI Framework** | Bubble Tea | Interactive terminal UI |
| **UUIDs** | google/uuid | Unique identifiers |
| **Testing** | Go testing package | Unit tests, benchmarks |

## Project Structure

```
dopa/
├── cmd/dopa/          # CLI commands (Presentation Layer)
│   ├── main.go             # Entry point, service container
│   ├── tasks.go            # Task commands
│   └── projects.go         # Project commands
│
├── internal/
│   ├── domain/             # Domain Layer
│   │   ├── task.go         # Task entity with factory method
│   │   ├── project.go      # Project entity
│   │   ├── area.go         # Area entity
│   │   └── value_objects.go # Value objects (Status, Priority, Color)
│   │
│   ├── service/            # Service Layer
│   │   ├── interfaces.go   # Service interfaces
│   │   ├── task_service.go # Task business logic
│   │   └── project_service.go
│   │
│   ├── converter/          # Converter Layer
│   │   └── converter.go    # DB ↔ Domain conversions
│   │
│   ├── db/                 # Repository Layer
│   │   ├── querier.go      # Repository interface
│   │   ├── tasks.sql.go    # Task queries (generated)
│   │   └── models.go       # DB models (generated)
│   │
│   └── tui/                # TUI (Presentation Layer)
│       ├── app.go          # Main application
│       └── commands.go     # TUI commands
│
└── docs/                   # Documentation
    ├── START_HERE.md       # This file
    ├── architecture/       # Architecture documentation
    └── TUI.md              # TUI documentation
```

## Getting Help

- **Issues**: Report bugs or request features via GitHub Issues
- **Documentation**: You're reading it! Start with [Architecture Overview](architecture/01-overview.md)
- **Code Examples**: Each architecture chapter includes real code examples from the codebase

## Quick Reference

### Common Patterns

**Creating a domain entity**:
```go
task, err := domain.NewTask(domain.NewTaskParams{
    ProjectID: projectID,
    Title:     "Write documentation",
    Status:    domain.TaskStatusTodo,
    Priority:  domain.PriorityHigh,
})
```

**Using a service**:
```go
services := GetServices()
task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
    ProjectID: projectID,
    Title:     "Write documentation",
    Status:    domain.TaskStatusTodo,
})
```

**Writing a table-driven test**:
```go
tests := []struct {
    name    string
    params  domain.NewTaskParams
    wantErr error
}{
    {
        name: "valid task",
        params: domain.NewTaskParams{
            Title:     "Test",
            ProjectID: "proj-123",
            Status:    domain.TaskStatusTodo,
            Priority:  domain.PriorityMedium,
        },
        wantErr: nil,
    },
    // ...
}
```

---

**Next Step**: Start with [Architecture Overview](architecture/01-overview.md) to understand the system structure.
