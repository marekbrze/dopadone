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
│        Converter Layer                  │
│  Type Transformations (DB ↔ Domain)     │
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

| Layer | Directory | Key Files |
|-------|-----------|-----------|
| **Domain** | `internal/domain/` | `task.go`, `project.go`, `area.go`, `value_objects.go` |
| **Services** | `internal/service/` | `interfaces.go`, `task_service.go`, `project_service.go` |
| **Converters** | `internal/converter/` | `converter.go` |
| **Repository** | `internal/db/` | `querier.go`, `tasks.sql.go`, `models.go` |
| **CLI** | `cmd/dopa/` | `main.go`, `tasks.go`, `projects.go` |
| **TUI** | `internal/tui/` | `app.go`, `commands.go` |

## Related Documentation

- **[Start Here](../START_HERE.md)** - Entry point with learning paths
- **[TUI Documentation](../TUI.md)** - Terminal UI architecture
- **[Transactions](../TRANSACTIONS.md)** - Database transaction patterns

---

**Ready to dive in?** Start with [01 - Architecture Overview](01-overview.md)
