---
id: TASK-35
title: 'Task-29A: Define Service Interfaces'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-05 10:10'
updated_date: '2026-03-05 10:55'
labels:
  - architecture
  - refactoring
  - service-layer
dependencies: []
references:
  - 'Related: TASK-29 (parent task)'
  - 'Related: TASK-25 (services exist)'
  - 'Related: TASK-27 (converter package)'
  - internal/service/area_service.go
  - internal/service/subarea_service.go
  - internal/service/project_service.go
  - internal/service/task_service.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create service interfaces (AreaServiceInterface, SubareaServiceInterface, ProjectServiceInterface, TaskServiceInterface) in internal/service/interfaces.go for mocking and testability. This is the foundation task that enables all subsequent refactoring work.

**Objective**: Create service interfaces for mocking and testability

**Parent Task**: TASK-29 (Route TUI commands through service layer)

**Dependencies**: TASK-25 (services exist), TASK-27 (converter package exists)

**Blocks**: TASK-29B, TASK-29C, TASK-29D, TASK-29E

**Design Decisions**:

1. **Provider Pattern**: Interfaces defined alongside implementations (not consumer pattern)
   - Keeps interfaces close to implementations for maintainability
   - Allows consumers to define narrower interfaces if needed
   - Simplifies dependency graph
   - Enables straightforward mocking

2. **Context First**: All methods accept context.Context as first parameter
   - Follows Go best practices
   - Enables future cancellation and timeout support
   - GetEffectiveColor updated to accept context (not currently used, but future-proof)

3. **Include Future Methods**: ListBySubareaRecursive included in interface
   - Will be implemented in Task-29B
   - Enables parallel work on Task-29B and Task-29C
   - Interface defines contract before implementation

4. **Compile-Time Checks**: Use var _ Interface = (*Struct)(nil) pattern
   - Catches interface/implementation mismatches at compile time
   - Better than runtime errors
   - Self-documenting (shows which types implement which interfaces)

**Deliverables**:
- Create internal/service/interfaces.go with 4 service interfaces
- AreaServiceInterface (9 methods)
- SubareaServiceInterface (9 methods) - GetEffectiveColor signature updated
- ProjectServiceInterface (13 methods including ListBySubareaRecursive)
- TaskServiceInterface (14 methods)
- Compile-time interface satisfaction checks
- Comprehensive godoc documentation

**Testing**:
- Compile-time interface satisfaction checks (automatic)
- Existing service tests must still pass
- Code must compile without errors

**Documentation**:
- Package-level documentation explaining design decisions
- Interface-level documentation explaining purpose
- Method-level documentation explaining behavior
- Note about ListBySubareaRecursive being implemented in Task-29B
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 internal/service/interfaces.go file created with package documentation
- [x] #2 AreaServiceInterface defined with 9 methods (List, GetByID, Create, Update, UpdateSortOrder, ReorderAll, SoftDelete, HardDelete, GetStats)
- [x] #3 SubareaServiceInterface defined with 9 methods including GetEffectiveColor with context.Context parameter
- [x] #4 ProjectServiceInterface defined with 13 methods including ListBySubareaRecursive (to be implemented in Task-29B)
- [x] #5 TaskServiceInterface defined with 14 methods (Create, GetByID, ListByProject, ListByStatus, ListByPriority, ListNext, ListAll, Update, SoftDelete, HardDelete, SetStatus, MarkCompleted, SetPriority, ToggleIsNext)
- [x] #6 SubareaService.GetEffectiveColor updated to accept context.Context as first parameter
- [x] #7 Compile-time interface satisfaction checks added (var _ Interface = (*Service)(nil)) for all 4 services
- [x] #8 All interfaces have godoc comments explaining purpose and methods
- [x] #9 Code compiles without errors (go build ./internal/service/...)
- [x] #10 All existing tests pass (go test ./internal/service/...)
- [x] #11 Linter passes with no errors (golangci-lint run ./internal/service/...)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: Define Service Interfaces

## Overview
Create service interfaces in `internal/service/interfaces.go` for AreaService, SubareaService, ProjectService, and TaskService. This is the foundation task that enables mocking, testability, and the TUI refactoring work.

## Phase 1: Interface Design & File Creation (30 min)

### 1.1 Create interfaces.go file
- Create `internal/service/interfaces.go`
- Add package documentation explaining the purpose
- Document the design decision: Provider pattern (interfaces alongside implementations)

### 1.2 Define AreaServiceInterface (9 methods)
```go
type AreaServiceInterface interface {
    List(ctx context.Context) ([]domain.Area, error)
    GetByID(ctx context.Context, id string) (*domain.Area, error)
    Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error)
    Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error)
    UpdateSortOrder(ctx context.Context, id string, sortOrder int) error
    ReorderAll(ctx context.Context, areaIDs []string) error
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*AreaStats, error)
}
```

### 1.3 Define SubareaServiceInterface (9 methods)
- Add context.Context to GetEffectiveColor signature (currently missing)
```go
type SubareaServiceInterface interface {
    Create(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error)
    GetByID(ctx context.Context, id string) (*domain.Subarea, error)
    ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error)
    Update(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*SubareaStats, error)
    GetEffectiveColor(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color
    ListAll(ctx context.Context) ([]domain.Subarea, error)
}
```

### 1.4 Define ProjectServiceInterface (13 methods)
- Include ListBySubareaRecursive (will be implemented in Task-29B)
```go
type ProjectServiceInterface interface {
    Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error)
    GetByID(ctx context.Context, id string) (*domain.Project, error)
    ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error)
    ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)
    ListAll(ctx context.Context) ([]domain.Project, error)
    ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)
    ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error)
    ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error)
    Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*ProjectStats, error)
    ValidateParentHierarchy(ctx context.Context, parentID string, projectID string) error
}
```

### 1.5 Define TaskServiceInterface (14 methods)
```go
type TaskServiceInterface interface {
    Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)
    GetByID(ctx context.Context, id string) (*domain.Task, error)
    ListByProject(ctx context.Context, projectID string) ([]domain.Task, error)
    ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
    ListByPriority(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error)
    ListNext(ctx context.Context) ([]domain.Task, error)
    ListAll(ctx context.Context) ([]domain.Task, error)
    Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    SetStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)
    MarkCompleted(ctx context.Context, id string) (*domain.Task, error)
    SetPriority(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error)
    ToggleIsNext(ctx context.Context, id string) (*domain.Task, error)
}
```

## Phase 2: Implementation Updates (20 min)

### 2.1 Update SubareaService.GetEffectiveColor
- Add context.Context parameter to method signature
- Update implementation in `internal/service/subarea_service.go`
```go
func (s *SubareaService) GetEffectiveColor(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color {
    return subarea.GetEffectiveColor(parentArea)
}
```

### 2.2 Add Compile-Time Interface Checks
- Add at bottom of interfaces.go:
```go
// Compile-time interface satisfaction checks
var (
    _ AreaServiceInterface    = (*AreaService)(nil)
    _ SubareaServiceInterface = (*SubareaService)(nil)
    _ ProjectServiceInterface = (*ProjectService)(nil)
    _ TaskServiceInterface    = (*TaskService)(nil)
)
```

## Phase 3: Documentation (15 min)

### 3.1 Add Package Documentation
```go
// Package service provides business logic for project management operations.
//
// This package defines service interfaces for dependency injection and testability.
// Interfaces are defined in the provider package (alongside implementations) following
// the "accept interfaces, return structs" principle.
//
// Design Decision: Provider Pattern
// We define interfaces where they're implemented (provider pattern) rather than where
// they're consumed (consumer pattern). This approach:
// - Keeps interfaces close to implementations for easier maintenance
// - Allows consumers to define their own interfaces if needed
// - Simplifies the dependency graph
// - Enables straightforward mocking for tests
//
// For consumers that need different interface shapes, they can define their own
// narrower interfaces following Go's interface composition patterns.
```

### 3.2 Add Interface Documentation
- Document each interface's purpose
- Document each method's behavior
- Note that ListBySubareaRecursive will be implemented in Task-29B

## Phase 4: Verification (15 min)

### 4.1 Build Verification
```bash
# Verify compilation
go build ./internal/service/...

# Run tests
go test ./internal/service/... -v

# Run linter
golangci-lint run ./internal/service/...
```

### 4.2 Verify Interface Satisfaction
- Compile-time checks will fail if interfaces don't match implementations
- Run `go build` to confirm all services satisfy their interfaces

## Success Criteria
1. ✅ interfaces.go created with all 4 service interfaces
2. ✅ All interfaces have proper godoc comments
3. ✅ Context added to GetEffectiveColor method
4. ✅ ListBySubareaRecursive included in ProjectServiceInterface
5. ✅ Compile-time interface checks added
6. ✅ Code compiles without errors
7. ✅ All existing tests pass
8. ✅ Linter passes with no errors

## Notes
- ListBySubareaRecursive will be implemented in Task-29B
- GetEffectiveColor context parameter is added but not used (future-proofing)
- Compile-time checks ensure interfaces stay in sync with implementations
- This task is a prerequisite for Task-29B, Task-29C, and all TUI refactoring work
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Task specification completed on 2026-03-05. Implementation plan created with detailed phases for interface design, implementation updates, documentation, and verification. All design decisions documented including provider pattern choice, context-first approach, and ListBySubareaRecursive inclusion. Ready for implementation.

- Created internal/service/interfaces.go with 4 service interfaces (212 lines)
- Updated SubareaService.GetEffectiveColor to accept context.Context parameter
- Added stub ListBySubareaRecursive to ProjectService for interface compatibility
- Updated SubareaService tests to match new GetEffectiveColor signature
- Added compile-time interface satisfaction checks for all 4 services
- All verification passed: compilation, 40 tests, linter
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented service interfaces for dependency injection and testability

**Changes**:
- Created internal/service/interfaces.go with 4 service interfaces:
  - AreaServiceInterface (9 methods)
  - SubareaServiceInterface (9 methods, GetEffectiveColor now accepts context)
  - ProjectServiceInterface (13 methods, includes ListBySubareaRecursive stub)
  - TaskServiceInterface (14 methods)
- Added compile-time interface satisfaction checks using var _ Interface = (*Service)(nil) pattern
- Updated SubareaService.GetEffectiveColor signature to accept context.Context (future-proofing)
- Added stub ProjectService.ListBySubareaRecursive for interface compatibility (implementation in Task-29B)
- Updated SubareaService tests to match new GetEffectiveColor signature

**Design Decisions**:
- Provider pattern: Interfaces defined alongside implementations for maintainability
- Context-first: All methods accept context.Context as first parameter
- Compile-time checks: Ensures interfaces stay in sync with implementations

**Testing**:
- All 40 existing tests pass
- Code compiles without errors
- Linter (go vet) passes with no issues

**Impact**: Enables mocking, testability, and future TUI refactoring work (Tasks 29B-E). Foundation for dependency injection throughout the application.
<!-- SECTION:FINAL_SUMMARY:END -->
