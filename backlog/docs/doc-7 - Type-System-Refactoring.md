---
id: doc-7
title: Type System Refactoring
type: technical
created_date: '2026-03-05'
---

# Type System Refactoring

## Overview

This document describes the type system refactoring completed in Task-30, which replaced `interface{}` types with proper `*time.Time` nullable types across all database models. This refactoring improves type safety, eliminates runtime type assertions, and makes the codebase more idiomatic.

## Problem

### Before Refactoring

The sqlc-generated database models used `interface{}` for nullable timestamp fields:

```go
type Area struct {
    DeletedAt interface{} `json:"deleted_at"`
}

type Project struct {
    Deadline    interface{} `json:"deadline"`
    CompletedAt interface{} `json:"completed_at"`
    DeletedAt   interface{} `json:"deleted_at"`
}

type Task struct {
    StartDate interface{} `json:"start_date"`
    Deadline  interface{} `json:"deadline"`
    DeletedAt interface{} `json:"deleted_at"`
}
```

This caused several issues:

1. **Type Safety**: No compile-time checking of timestamp types
2. **Runtime Panics**: Type assertions could fail at runtime
3. **Boilerplate Code**: Required type assertions in converter functions
4. **Code Duplication**: Same conversion pattern repeated across codebase

### Example of Problematic Code

```go
// Before: Type assertion with potential runtime panic
var deletedAt *time.Time
if dbArea.DeletedAt != nil {
    if t, ok := dbArea.DeletedAt.(time.Time); ok {
        deletedAt = &t
    }
}
```

## Solution

### Configuration Changes

Updated `sqlc.yaml` to configure proper type overrides for nullable timestamp columns:

```yaml
overrides:
  - column: "areas.deleted_at"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "subareas.deleted_at"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "projects.deadline"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "projects.completed_at"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "projects.deleted_at"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "tasks.start_date"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "tasks.deadline"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
  
  - column: "tasks.deleted_at"
    nullable: true
    go_type:
      import: "time"
      type: "Time"
      pointer: true
```

### After Refactoring

Database models now use `*time.Time` directly:

```go
type Area struct {
    DeletedAt *time.Time `json:"deleted_at"`
}

type Project struct {
    Deadline    *time.Time `json:"deadline"`
    CompletedAt *time.Time `json:"completed_at"`
    DeletedAt   *time.Time `json:"deleted_at"`
}

type Task struct {
    StartDate *time.Time `json:"start_date"`
    Deadline  *time.Time `json:"deadline"`
    DeletedAt *time.Time `json:"deleted_at"`
}
```

### Simplified Code

Converter functions no longer need type assertions:

```go
// After: Direct assignment
deletedAt := dbArea.DeletedAt
```

## Impact

### Type Safety

- **Before**: Runtime type assertions could panic
- **After**: Compile-time type checking catches errors early
- **Result**: Eliminated entire class of potential runtime errors

### Code Reduction

- **Converter Package**: Removed ~50 lines of type assertion code
- **Service Layer**: Removed ~50 lines of interface{} conversion code
- **Total**: ~100 lines of boilerplate removed

### Before & After Comparison

#### Converter Function (Before)

```go
func DbProjectToDomain(p db.Project) domain.Project {
    var deletedAt *time.Time
    if p.DeletedAt != nil {
        if t, ok := p.DeletedAt.(time.Time); ok {
            deletedAt = &t
        }
    }
    
    var deadline *time.Time
    if p.Deadline != nil {
        if t, ok := p.Deadline.(time.Time); ok {
            deadline = &t
        }
    }
    
    var completedAt *time.Time
    if p.CompletedAt != nil {
        if t, ok := p.CompletedAt.(time.Time); ok {
            completedAt = &t
        }
    }
    
    return domain.Project{
        DeletedAt:   deletedAt,
        Deadline:    deadline,
        CompletedAt: completedAt,
        // ... other fields
    }
}
```

#### Converter Function (After)

```go
func DbProjectToDomain(p db.Project) domain.Project {
    return domain.Project{
        DeletedAt:   p.DeletedAt,
        Deadline:    p.Deadline,
        CompletedAt: p.CompletedAt,
        // ... other fields
    }
}
```

## Files Changed

### Configuration

- `sqlc.yaml`: Added type overrides for nullable timestamp columns

### Generated Code

- `internal/db/models.go`: Auto-generated with `*time.Time` types

### Converter Package

- `internal/converter/converter.go`: Removed type assertions
  - `DbAreaToDomain`
  - `DbListAreasRowToDomain`
  - `DbGetAreaByIDRowToDomain`
  - `DbCreateAreaRowToDomain`
  - `DbUpdateAreaRowToDomain`
  - `DbSubareaToDomain`
  - `DbProjectToDomain`
  - `DbTaskToDomain`

### Service Layer

- `internal/service/project_service.go`: Removed interface{} conversion code
- `internal/service/task_service.go`: Removed interface{} conversion code
- `internal/service/area_service.go`: Updated SoftDelete method
- `internal/service/subarea_service.go`: Updated SoftDelete method

### Tests

- `internal/db/areas_test.go`: Updated to use `*time.Time` pointers
- `internal/db/projects_test.go`: Updated to use `*time.Time` pointers
- `internal/db/subareas_test.go`: Updated to use `*time.Time` pointers
- `internal/db/integration_test.go`: Updated to use `*time.Time` pointers

## Testing

All tests pass after refactoring:

- ✅ Converter tests (6 tests)
- ✅ Service tests (47 tests)
- ✅ Integration tests (DB layer)
- ✅ CLI tests (42 tests)
- ✅ Domain tests
- ✅ Race condition detection enabled
- ✅ `go vet ./...` passes with no errors

## Benefits

1. **Improved Type Safety**: Compile-time checking prevents type mismatches
2. **Reduced Boilerplate**: Eliminated ~100 lines of conversion code
3. **Better Maintainability**: Clearer intent, fewer places for bugs
4. **More Idiomatic**: Using `*time.Time` is standard Go practice
5. **No Runtime Panics**: Eliminated type assertion failures
6. **Better Performance**: Direct field access instead of type assertions

## Backward Compatibility

This is a pure refactoring with no API changes:

- Domain types unchanged (already used `*time.Time`)
- Service interfaces unchanged
- CLI commands unchanged
- TUI behavior unchanged
- Database schema unchanged

Only internal type representations improved.

## Related Documentation

- [Data Layer Architecture](doc-1 - Data-Layer-Architecture.md) - Updated with new type system
- [Service Layer Architecture](doc-5 - Service-Layer-Architecture.md) - Service layer usage
- Task-30: Complete implementation details and acceptance criteria

## Lessons Learned

1. **sqlc Configuration**: Column-specific overrides are more reliable than database-level type overrides
2. **Type Assertions**: Avoid `interface{}` when possible; use proper types from the start
3. **Testing**: Comprehensive test coverage enabled confident refactoring
4. **Gradual Migration**: Could have done this incrementally if needed

## Future Considerations

1. **Other Nullable Types**: Apply same pattern to other nullable fields if needed
2. **Custom Types**: Consider custom types for timestamps with validation
3. **Documentation**: Keep documentation synchronized with code changes
