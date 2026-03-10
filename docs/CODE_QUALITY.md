# Code Quality and Linting Standards

## Overview

This document describes the code quality standards and linting requirements for the Dopadone project. All code must pass the configured linters before being merged.

## Quick Reference

```bash
# Run all linters
make lint

# Run specific linter
golangci-lint run --disable-all --enable=errcheck

# Auto-fix issues where possible
golangci-lint run --fix
```

## Enabled Linters

The project uses golangci-lint with the following linters enabled:

### Critical Linters

| Linter | Purpose | Severity |
|--------|---------|----------|
| `errcheck` | Checks for unchecked error returns | High |
| `govet` | Reports suspicious constructs | High |
| `staticcheck` | Advanced static analysis | High |
| `typecheck` | Type checking | High |

### Code Style Linters

| Linter | Purpose | Severity |
|--------|---------|----------|
| `gofmt` | Code formatting | Medium |
| `goimports` | Import statement organization | Medium |
| `gosimple` | Code simplification suggestions | Medium |
| `ineffassign` | Detects ineffective assignments | Medium |

### Code Quality Linters

| Linter | Purpose | Severity |
|--------|---------|----------|
| `goconst` | Detects repeated strings that could be constants | Medium |
| `gocyclo` | Cyclomatic complexity (threshold: 20) | Medium |
| `dupl` | Code clone detection (threshold: 150) | Low |

## Linter-Specific Guidelines

### errcheck - Error Handling

**Rule**: Never ignore error return values.

#### Production Code

In production code, all errors must be handled explicitly:

```go
// ✅ GOOD: Handle error explicitly
result, err := someFunction()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// ✅ GOOD: Log and continue if appropriate
if err := services.Close(); err != nil {
    log.Printf("warning: failed to close services: %v", err)
}

// ✅ GOOD: Handle in CLI layer with user message
if err := cmd.MarkFlagRequired("project-id"); err != nil {
    return fmt.Errorf("failed to mark flag as required: %w", err)
}
```

#### Test Code

In test code, cleanup operations may use the blank identifier pattern to explicitly indicate intent:

```go
// ✅ GOOD: Explicitly ignore cleanup errors in defer
func TestSomething(t *testing.T) {
    db, err := createTestDB()
    if err != nil {
        t.Fatal(err)
    }
    defer func() { _ = db.Close() }()
    
    rows, err := db.Query("SELECT ...")
    if err != nil {
        t.Fatal(err)
    }
    defer func() { _ = rows.Close() }()
    
    // ... test logic
}
```

**Why ignore cleanup errors in tests?**
- Cleanup failures should not mask the actual test result
- In defer, errors cannot be returned
- Logging in tests adds noise for non-critical cleanup
- The `_ =` pattern makes intent explicit to readers and linters

**Alternative: Using t.Cleanup()**

```go
func TestDatabase(t *testing.T) {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    
    t.Cleanup(func() {
        _ = db.Close()
    })
    
    // ... test logic
}
```

### gocyclo - Cyclomatic Complexity

**Rule**: Keep cyclomatic complexity under 20.

**Patterns to reduce complexity:**
1. Extract helper methods
2. Use early returns
3. Use table-driven tests
4. Use switch statements instead of if-else chains

See [CLI Layer Documentation](architecture/06-cli-layer.md) for detailed examples.

### goconst - String Constants

**Rule**: Extract repeated strings (3+ occurrences) as constants.

```go
// ❌ BAD: Repeated strings
func processStatus(status string) {
    if status == "pending" {
        // ...
    }
}

func validateStatus(status string) error {
    if status != "pending" && status != "completed" {
        return errors.New("invalid status")
    }
}

// ✅ GOOD: Use constants
const (
    StatusPending   = "pending"
    StatusCompleted = "completed"
)

func processStatus(status string) {
    if status == StatusPending {
        // ...
    }
}
```

### dupl - Code Duplication

**Rule**: Keep code duplication under 150 tokens.

**Strategies:**
1. Extract common logic into helper functions
2. Use composition over copy-paste
3. Create reusable test helpers

## CI Integration

### Pre-commit Checks

The CI pipeline runs these checks on every push:

1. **Test Job**: `go test ./... -race -coverprofile=coverage.out`
2. **Lint Job**: 
   - `go vet ./...`
   - `golangci-lint run`
3. **Build Job**: `go build ./...`

### Coverage Gate

- Minimum coverage: 20%
- Measured on entire codebase
- Blocks PR merge if below threshold

## Configuration

The linter configuration is in `.golangci.yml`:

```yaml
linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - ineffassign
    - typecheck
    - gosimple
    - goconst
    - gocyclo
    - dupl

linters-settings:
  gocyclo:
    min-complexity: 20
  dupl:
    threshold: 150
  goconst:
    min-len: 3
    min-occurrences: 3
```

## Common Issues and Solutions

### Issue: errcheck warnings in defer

**Problem**: `defer file.Close()` triggers errcheck warning.

**Solution**: Use blank identifier to explicitly ignore:
```go
defer func() { _ = file.Close() }()
```

### Issue: High cyclomatic complexity

**Problem**: Function has too many branches.

**Solution**: Extract helper methods or use table-driven approach:
```go
// Before: complex if-else chain
func process(cmd string) error {
    if cmd == "create" {
        // 20 lines
    } else if cmd == "update" {
        // 20 lines
    } else if cmd == "delete" {
        // 20 lines
    }
    return nil
}

// After: dispatch table
func process(cmd string) error {
    handlers := map[string]func() error{
        "create": handleCreate,
        "update": handleUpdate,
        "delete": handleDelete,
    }
    
    handler, ok := handlers[cmd]
    if !ok {
        return fmt.Errorf("unknown command: %s", cmd)
    }
    return handler()
}
```

### Issue: Duplicate code in tests

**Problem**: Similar test setup code repeated across files.

**Solution**: Create test helpers:
```go
// internal/test/helpers.go
func SetupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to create test db: %v", err)
    }
    t.Cleanup(func() { _ = db.Close() })
    return db
}

// Usage in tests
func TestSomething(t *testing.T) {
    db := test.SetupTestDB(t)
    // ... use db
}
```

## Best Practices

### 1. Run Linters Locally

Always run `make lint` before pushing:
```bash
# Quick check
make lint

# With auto-fix
golangci-lint run --fix

# Specific package
golangci-lint run ./internal/service/...
```

### 2. Fix Issues Immediately

Don't let lint issues accumulate. Fix them as they appear.

### 3. Understand the Rule

Before disabling a linter rule, understand why it exists:
```bash
# Check linter documentation
golangci-lint linters
```

### 4. Use Meaningful Variable Names

```go
// ❌ BAD: Cryptic names
if err := f(); err != nil {
    _ = db.C()  // What is C()?
}

// ✅ GOOD: Clear names
if err := fetchData(); err != nil {
    _ = database.Close()  // Clear intent
}
```

### 5. Document Exceptions

If you must disable a linter, document why:
```go
//nolint:errcheck // Test cleanup - error not critical
defer func() { _ = rows.Close() }()
```

## Related Documentation

- [Testing Strategy](architecture/07-testing-strategy.md) - Test patterns and helpers
- [CI/CD Pipeline](CI-CD.md) - CI workflow details
- [Release Process](RELEASE.md) - Release requirements

## Getting Help

- Run `golangci-lint linters` to see all available linters
- Run `golangci-lint run --help` for CLI options
- Check [golangci-lint documentation](https://golangci-lint.run/)
- Review existing code for patterns

---

**Last Updated**: 2026-03-10  
**Related**: TASK-75 (errcheck fixes), TASK-73 (gocyclo), TASK-74 (TUI complexity)
