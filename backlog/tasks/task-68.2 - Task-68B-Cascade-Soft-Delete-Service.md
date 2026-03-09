---
id: TASK-68.2
title: 'Task-68B: Cascade Soft Delete Service'
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-09 19:04'
updated_date: '2026-03-09 20:50'
labels: []
dependencies: []
references:
  - '# Parent: TASK-68 - Add option to delete subareas'
  - projects and tasks in tui
  - '# Sibling: TASK-68.1 - Confirmation Modal Component'
  - '# Dependent: TASK-68.3 - TUI Delete Integration'
  - '# Milestone: m-2 Deleting items'
parent_task_id: TASK-68
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement cascade soft delete for projects with subprojects. Recursively soft delete all child projects and their tasks in a transaction. Add SoftDeleteWithCascade method to ProjectService.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Add SoftDeleteWithCascade method to ProjectService
- [x] #2 Implement softDeleteRecursive helper for nested projects
- [x] #3 Soft delete all tasks within deleted projects
- [x] #4 Use transaction for atomic cascade delete
- [x] #5 Add SoftDeleteTasksByProject SQL query
- [x] #6 Write service layer tests with 100% coverage
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Task 68.2: Cascade Soft Delete Service

## Task Overview
**Parent Task**: TASK-68 - Add option to delete subareas, projects and tasks in TUI  
**Dependencies**: 
- Parallel with: TASK-68.1 (Confirmation Modal Component)
- Blocks: TASK-68.3 (TUI Delete Integration)

**Estimated Effort**: 8.5 hours (~1-2 days)

## Implementation Phases

### Phase 1: Database Layer Foundation (SEQUENTIAL - Start Here)
**Duration**: 30 minutes  
**Priority**: HIGH - Blocks all other work

#### Subtask 1.1: Add SQL Query
**File**: `queries/projects.sql`

```sql
-- name: SoftDeleteTasksByProject :exec
-- Soft deletes all tasks within a project
-- Used during cascade soft delete to remove tasks in bulk
UPDATE tasks
SET deleted_at = ?
WHERE project_id = ? AND deleted_at IS NULL;
```

**Rationale**:
- Follows existing query patterns (see `SoftDeleteProject`, `SoftDeleteTask`)
- Uses `:exec` since we don't need return values (efficient)
- Checks `deleted_at IS NULL` to avoid re-deleting (idempotent)
- Bulk delete for performance (one query per project, not per task)

**Validation**:
- [ ] Query syntax is correct
- [ ] Follows project naming conventions
- [ ] Includes descriptive comment
- [ ] Idempotent (can run multiple times safely)

#### Subtask 1.2: Generate Go Code
**Command**: `make sqlc-generate`

**Generated Files**:
- `internal/db/projects.sql.go` - Query implementation
- `internal/db/querier.go` - Interface update (adds method to Querier interface)

**Validation**:
- [ ] `make sqlc-generate` runs without errors
- [ ] Generated code compiles
- [ ] New method appears in `db.Querier` interface
- [ ] Run `go build ./...` to verify

**Phase 1 Complete When**: SQL query added, code generated, builds successfully

---

### Phase 2: Service Layer Implementation (SEQUENTIAL - Depends on Phase 1)
**Duration**: 2 hours  
**Priority**: HIGH

#### Subtask 2.1: Implement Public Method - SoftDeleteWithCascade
**File**: `internal/service/project_service.go`

**Method Signature**:
```go
// SoftDeleteWithCascade soft deletes a project and all its descendants (child projects and tasks).
// The operation is atomic - either all deletions succeed or none do.
// Returns ErrProjectNotFound if the project doesn't exist.
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, id string) error
```

**Implementation Pattern** (follow existing `HardDelete()` pattern):
```go
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, id string) error {
    // Step 1: Validate project exists
    project, err := s.repo.GetProjectByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrProjectNotFound
        }
        return fmt.Errorf("get project %s: %w", id, err)
    }
    
    // Step 2: Check if already deleted (idempotent)
    if project.DeletedAt.Valid {
        return nil // Already soft deleted, no-op
    }

    // Step 3: Execute in transaction if available
    if s.tm == nil {
        return s.softDeleteRecursive(ctx, s.repo, id)
    }

    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        return s.softDeleteRecursive(ctx, tx, id)
    })
}
```

**Key Design Decisions**:
- **Validation first**: Check project exists before expensive operations
- **Idempotency**: Safe to call multiple times (no error if already deleted)
- **Transaction support**: Works with or without transaction manager
- **Error wrapping**: Provides context for debugging

**Validation**:
- [ ] Method follows existing code patterns
- [ ] Error messages are descriptive
- [ ] Context propagation (ctx as first param)
- [ ] Transaction handling (both with and without tm)
- [ ] Proper godoc comment

#### Subtask 2.2: Implement Private Helper - softDeleteRecursive
**File**: `internal/service/project_service.go`

**Method Signature**:
```go
// softDeleteRecursive recursively soft deletes a project, its children, and their tasks.
// Uses depth-first traversal: delete children first, then parent.
func (s *ProjectService) softDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error
```

**Implementation Pattern** (follow existing `hardDeleteRecursive()` pattern):
```go
func (s *ProjectService) softDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error {
    // Step 1: Get all direct children
    children, err := q.ListProjectsByParent(ctx, sql.NullString{
        String: projectID,
        Valid:  true,
    })
    if err != nil {
        return fmt.Errorf("list child projects of %s: %w", projectID, err)
    }

    // Step 2: Recursively delete each child (depth-first)
    for _, child := range children {
        if err := s.softDeleteRecursive(ctx, q, child.ID); err != nil {
            return fmt.Errorf("cascade delete child %s: %w", child.ID, err)
        }
    }

    // Step 3: Soft delete all tasks in current project
    now := time.Now()
    if err := q.SoftDeleteTasksByProject(ctx, db.SoftDeleteTasksByProjectParams{
        DeletedAt: sql.NullTime{
            Time:  now,
            Valid: true,
        },
        ProjectID: projectID,
    }); err != nil {
        return fmt.Errorf("soft delete tasks in project %s: %w", projectID, err)
    }

    // Step 4: Soft delete current project
    if err := q.SoftDeleteProject(ctx, db.SoftDeleteProjectParams{
        DeletedAt: sql.NullTime{
            Time:  now,
            Valid: true,
        },
        ID: projectID,
    }); err != nil {
        return fmt.Errorf("soft delete project %s: %w", projectID, err)
    }

    return nil
}
```

**Key Design Decisions**:
- **Depth-first traversal**: Delete children before parent (maintains referential integrity)
- **Error wrapping**: Each operation wrapped with context
- **Same timestamp**: All deletions in cascade use same `now` (consistent audit trail)
- **Accepts interface**: Uses `db.Querier` not concrete type (testability)

**Validation**:
- [ ] Follows existing recursive pattern
- [ ] Proper base case (no children = just delete tasks + project)
- [ ] Error context at each level
- [ ] Depth-first traversal (children before parent)
- [ ] Consistent timestamp across all deletions

**Phase 2 Complete When**: Both methods implemented, code compiles, follows patterns

---

### Phase 3: Unit Testing (PARALLEL after Phase 2)
**Duration**: 3 hours  
**Priority**: HIGH  
**Can Run In Parallel With**: Nothing (depends on Phase 2)

#### Subtask 3.1: Update Mock Implementation
**File**: `internal/service/project_service_test.go`

**Add to mock struct**:
```go
type mockProjectQuerier struct {
    // ... existing fields ...
    softDeleteTasksByProjectFunc func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error
}

func (m *mockProjectQuerier) SoftDeleteTasksByProject(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
    if m.softDeleteTasksByProjectFunc != nil {
        return m.softDeleteTasksByProjectFunc(ctx, arg)
    }
    return nil
}
```

#### Subtask 3.2: Write Unit Tests - Success Scenarios
**File**: `internal/service/project_service_test.go`

**Test Cases**:

1. **Single project with no children**:
```go
t.Run("single project no children", func(t *testing.T) {
    // Setup: project exists, no children, 3 tasks
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify: 
    //   - SoftDeleteTasksByProject called once
    //   - SoftDeleteProject called once
    //   - No ListProjectsByParent calls
    //   - No error returned
})
```

2. **Project with direct children (1 level)**:
```go
t.Run("project with direct children", func(t *testing.T) {
    // Setup: parent + 2 children, each has tasks
    // Execute: SoftDeleteWithCascade(parentID)
    // Verify:
    //   - ListProjectsByParent called 3 times (parent + 2 children)
    //   - SoftDeleteTasksByProject called 3 times
    //   - SoftDeleteProject called 3 times (depth-first: child1, child2, parent)
    //   - Correct call order verified
    //   - No error returned
})
```

3. **Deeply nested project (3+ levels)**:
```go
t.Run("deeply nested hierarchy", func(t *testing.T) {
    // Setup: Level 0 → Level 1 → Level 2 → Level 3
    // Execute: SoftDeleteWithCascade(level0ID)
    // Verify:
    //   - All levels processed
    //   - Depth-first order verified
    //   - All tasks deleted
    //   - No error returned
})
```

4. **Multiple children at same level**:
```go
t.Run("multiple siblings", func(t *testing.T) {
    // Setup: Parent with 5 direct children
    // Execute: SoftDeleteWithCascade(parentID)
    // Verify:
    //   - All 5 children processed
    //   - All tasks deleted
    //   - No error returned
})
```

5. **Already soft-deleted project (idempotency)**:
```go
t.Run("already deleted project", func(t *testing.T) {
    // Setup: Project with deleted_at set
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - No recursive calls
    //   - No delete calls
    //   - No error returned (idempotent)
})
```

#### Subtask 3.3: Write Unit Tests - Error Scenarios
**File**: `internal/service/project_service_test.go`

**Test Cases**:

1. **Non-existent project**:
```go
t.Run("non-existent project", func(t *testing.T) {
    // Setup: GetProjectByID returns sql.ErrNoRows
    // Execute: SoftDeleteWithCascade("nonexistent")
    // Verify:
    //   - Returns ErrProjectNotFound
    //   - No delete operations attempted
})
```

2. **Database error on GetProjectByID**:
```go
t.Run("database error on get", func(t *testing.T) {
    // Setup: GetProjectByID returns generic error
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - Returns wrapped error
    //   - Error contains "get project" context
    //   - No delete operations attempted
})
```

3. **Database error on ListProjectsByParent**:
```go
t.Run("database error on list children", func(t *testing.T) {
    // Setup: ListProjectsByParent returns error
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - Returns wrapped error
    //   - Error contains "list child projects" context
    //   - No delete operations attempted
})
```

4. **Database error on SoftDeleteTasksByProject**:
```go
t.Run("database error on delete tasks", func(t *testing.T) {
    // Setup: SoftDeleteTasksByProject returns error
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - Returns wrapped error
    //   - Error contains "soft delete tasks" context
    //   - SoftDeleteProject NOT called
})
```

5. **Database error on SoftDeleteProject**:
```go
t.Run("database error on delete project", func(t *testing.T) {
    // Setup: SoftDeleteProject returns error
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - Returns wrapped error
    //   - Error contains "soft delete project" context
})
```

6. **Error in nested child**:
```go
t.Run("error propagates from nested child", func(t *testing.T) {
    // Setup: 3-level hierarchy, error on level 2
    // Execute: SoftDeleteWithCascade(level0ID)
    // Verify:
    //   - Error propagates up
    //   - Error contains full context chain
    //   - Level 3 NOT processed (early exit)
})
```

#### Subtask 3.4: Write Unit Tests - Transaction Scenarios
**File**: `internal/service/project_service_test.go`

**Test Cases**:

1. **With transaction manager**:
```go
t.Run("with transaction manager", func(t *testing.T) {
    // Setup: Service with mock transaction manager
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - WithTransaction called
    //   - Operations use transaction querier
    //   - Transaction committed on success
})
```

2. **Without transaction manager**:
```go
t.Run("without transaction manager", func(t *testing.T) {
    // Setup: Service with nil transaction manager
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - Operations use direct repo
    //   - No transaction overhead
    //   - Succeeds (for SQLite, not distributed)
})
```

3. **Transaction rollback on error** (integration test needed):
```go
// Note: This requires integration test with real DB
// See Phase 4 for integration test plan
```

#### Subtask 3.5: Write Unit Tests - Cascade Verification
**File**: `internal/service/project_service_test.go`

**Test Cases**:

1. **Correct deletion order (depth-first)**:
```go
t.Run("verifies depth-first order", func(t *testing.T) {
    // Setup: Complex hierarchy with call recorder
    // Execute: SoftDeleteWithCascade(parentID)
    // Verify:
    //   - Children deleted before parent
    //   - Tasks deleted before their project
    //   - Order is consistent
})
```

2. **All levels processed**:
```go
t.Run("all hierarchy levels processed", func(t *testing.T) {
    // Setup: 4-level hierarchy with markers at each level
    // Execute: SoftDeleteWithCascade(rootID)
    // Verify:
    //   - Level 0, 1, 2, 3 all processed
    //   - No levels skipped
    //   - No duplicates
})
```

3. **Consistent timestamps**:
```go
t.Run("consistent deleted_at timestamps", func(t *testing.T) {
    // Setup: Hierarchy with timestamp capture
    // Execute: SoftDeleteWithCascade(projectID)
    // Verify:
    //   - All deletions use same timestamp
    //   - Timestamp is recent (within 1 second)
    //   - Consistent across tasks and projects
})
```

**Phase 3 Complete When**: All unit tests pass, 100% coverage for new code

---

### Phase 4: Integration Testing (SEQUENTIAL after Phase 3)
**Duration**: 2 hours  
**Priority**: MEDIUM  
**Depends On**: Phase 3 (unit tests passing)

#### Subtask 4.1: Create Integration Test File
**File**: `internal/service/project_service_cascade_integration_test.go`

**File Header**:
```go
//go:build integration
// +build integration

package service_test

import (
    "context"
    "testing"
    "time"
    
    "github.com/yourusername/dopadone/internal/db"
    "github.com/yourusername/dopadone/internal/service"
    "github.com/yourusername/dopadone/internal/domain"
)

// Integration tests for cascade soft delete functionality
// Run with: go test -tags=integration ./internal/service/...
```

#### Subtask 4.2: Write Integration Test Helper
**File**: `internal/service/project_service_cascade_integration_test.go`

```go
func setupIntegrationTest(t *testing.T) (*sql.DB, *service.ProjectService, func()) {
    t.Helper()
    
    // Create in-memory SQLite database
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }
    
    // Run migrations
    if err := runMigrations(db); err != nil {
        t.Fatalf("failed to run migrations: %v", err)
    }
    
    // Create service with transaction manager
    repo := db.New(db)
    tm := db.NewTransactionManager(db)
    svc := service.NewProjectService(repo, tm)
    
    // Cleanup function
    cleanup := func() {
        db.Close()
    }
    
    return db, svc, cleanup
}

func createTestHierarchy(t *testing.T, db *sql.DB) (rootID string, childIDs []string, taskIDs []string) {
    t.Helper()
    // Create: Area → Subarea → Root Project → 2 Children → Tasks
    // Return IDs for verification
}
```

#### Subtask 4.3: Write Integration Tests
**File**: `internal/service/project_service_cascade_integration_test.go`

**Test Cases**:

1. **Full cascade with real database**:
```go
func TestCascadeSoftDeleteIntegration(t *testing.T) {
    db, svc, cleanup := setupIntegrationTest(t)
    defer cleanup()
    
    rootID, childIDs, taskIDs := createTestHierarchy(t, db)
    
    // Execute cascade delete
    err := svc.SoftDeleteWithCascade(context.Background(), rootID)
    if err != nil {
        t.Fatalf("cascade delete failed: %v", err)
    }
    
    // Verify all projects soft deleted
    for _, id := range append([]string{rootID}, childIDs...) {
        project, err := db.GetProjectByID(context.Background(), id)
        if err != nil {
            t.Errorf("failed to get project %s: %v", id, err)
        }
        if !project.DeletedAt.Valid {
            t.Errorf("project %s not soft deleted", id)
        }
    }
    
    // Verify all tasks soft deleted
    for _, taskID := range taskIDs {
        task, err := db.GetTaskByID(context.Background(), taskID)
        if err != nil {
            t.Errorf("failed to get task %s: %v", taskID, err)
        }
        if !task.DeletedAt.Valid {
            t.Errorf("task %s not soft deleted", taskID)
        }
    }
}
```

2. **Transaction rollback verification**:
```go
func TestCascadeSoftDeleteTransactionRollback(t *testing.T) {
    // Setup database with constraints that will fail
    // Execute cascade delete that will fail mid-way
    // Verify: No partial deletes (all or nothing)
}
```

3. **Data integrity after delete**:
```go
func TestDataIntegrityAfterCascadeDelete(t *testing.T) {
    // Create hierarchy
    // Delete parent
    // Verify:
    //   - Foreign keys still valid
    //   - Can query deleted items with deleted filter
    //   - Non-deleted siblings unaffected
    //   - Timestamps consistent
}
```

4. **Performance with deep nesting**:
```go
func TestCascadeDeletePerformance(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping performance test in short mode")
    }
    
    // Create 10-level deep hierarchy
    // Measure delete time
    // Assert: < 1 second for 100 tasks across 10 levels
}
```

**Phase 4 Complete When**: All integration tests pass, run with `go test -tags=integration ./internal/service/...`

---

### Phase 5: Documentation Updates (PARALLEL with Phases 3-4)
**Duration**: 30 minutes  
**Priority**: MEDIUM  
**Can Run In Parallel With**: Phase 3 (Unit Testing), Phase 4 (Integration Testing)

#### Subtask 5.1: Update Service Documentation
**File**: `internal/service/project_service.go`

**Update godoc comment**:
```go
// ProjectService provides business logic for project operations.
//
// Cascade Delete Behavior:
// - SoftDeleteWithCascade: Recursively soft deletes a project and all descendants
//   including child projects and their tasks. Operation is atomic via transaction.
// - SoftDelete: Deletes only the specified project, no cascade.
// - HardDelete: Permanently removes project and all descendants from database.
//
// All cascade operations use depth-first traversal to maintain referential integrity.
type ProjectService struct {
    // ...
}
```

#### Subtask 5.2: Update Architecture Documentation (if needed)
**File**: `docs/architecture/03-service-layer.md` (if exists)

**Add section on cascade operations**:
```markdown
### Cascade Operations

The service layer supports cascade delete operations for hierarchical entities:

**Soft Delete Cascade**:
- Marks entities as deleted without removing data
- Preserves referential integrity
- Reversible operation
- Uses transactions for atomicity

**Implementation Pattern**:
1. Validate entity exists
2. Start transaction (if available)
3. Recursively process children (depth-first)
4. Delete child tasks
5. Delete child projects
6. Delete current entity
7. Commit transaction

See `ProjectService.SoftDeleteWithCascade` for reference implementation.
```

#### Subtask 5.3: Update README/API Documentation
**File**: `README.md` or `docs/API.md`

**Add to API section**:
```markdown
#### Cascade Soft Delete

```go
err := services.Projects.SoftDeleteWithCascade(ctx, projectID)
```

Soft deletes a project and all its descendants:
- All child projects (recursive)
- All tasks in those projects
- Atomic operation (all or nothing)
- Returns `ErrProjectNotFound` if project doesn't exist
- Idempotent (safe to call multiple times)

**Use Case**: When deleting a parent project, ensure all related data is also marked as deleted.
```

**Phase 5 Complete When**: Documentation updated and reviewed

---

### Phase 6: Code Quality & Final Verification (SEQUENTIAL at end)
**Duration**: 1 hour  
**Priority**: HIGH  
**Depends On**: All previous phases

#### Subtask 6.1: Run Linters
**Commands**:
```bash
# Format code
gofmt -w ./internal/service/
goimports -w ./internal/service/

# Run linters
golangci-lint run ./internal/service/...
go vet ./internal/service/...
staticcheck ./internal/service/...
```

**Validation**:
- [ ] No linting errors
- [ ] No compiler warnings
- [ ] Code formatted correctly

#### Subtask 6.2: Run Tests with Coverage
**Commands**:
```bash
# Run all unit tests
go test -v -race ./internal/service/...

# Run with coverage
go test -cover -coverprofile=coverage.out ./internal/service/...

# Check coverage percentage
go tool cover -func=coverage.out | grep total

# Verify 100% coverage for new methods
go tool cover -html=coverage.out
```

**Validation**:
- [ ] All tests pass
- [ ] No race conditions detected
- [ ] Coverage >= 100% for `SoftDeleteWithCascade` method
- [ ] Coverage >= 100% for `softDeleteRecursive` method
- [ ] Coverage >= 100% for SQL query execution

#### Subtask 6.3: Run Integration Tests
**Commands**:
```bash
# Run integration tests
go test -tags=integration -v ./internal/service/...

# Run all tests (unit + integration)
go test -v -race -tags=integration ./internal/service/...
```

**Validation**:
- [ ] Integration tests pass
- [ ] Real database operations work correctly
- [ ] Transaction behavior verified

#### Subtask 6.4: Manual Testing
**Steps**:
1. Build application: `make build`
2. Create test hierarchy via CLI
3. Delete parent project with cascade
4. Verify all descendants marked as deleted
5. Verify non-deleted items unaffected
6. Test error scenarios (non-existent ID, etc.)

**Validation**:
- [ ] Manual testing successful
- [ ] No unexpected behavior
- [ ] Error messages clear and helpful

#### Subtask 6.5: Code Review Checklist
**Self-Review**:
- [ ] Code follows existing patterns
- [ ] Error messages are descriptive
- [ ] Context propagation correct (ctx parameter)
- [ ] Transaction handling (both with and without tm)
- [ ] Idempotency verified (can call multiple times)
- [ ] No code duplication
- [ ] Proper godoc comments on all exported functions
- [ ] No commented-out code
- [ ] No TODO comments without tickets
- [ ] Variable names clear and consistent
- [ ] No magic numbers (use constants)

**Phase 6 Complete When**: All quality checks pass, ready for PR

---

## Task Dependencies and Execution Order

### Sequential Dependencies (Must Complete in Order)
```
Phase 1 (Database)
    ↓
Phase 2 (Service Layer)
    ↓
Phase 3 (Unit Tests)
    ↓
Phase 4 (Integration Tests)
    ↓
Phase 6 (Final Verification)
```

### Parallel Work Opportunities
```
Phase 2 (Service Layer)
    ↓
    ├──→ Phase 3 (Unit Tests) [SEQUENTIAL - needs implementation]
    └──→ Phase 5 (Documentation) [PARALLEL - can start after Phase 2]

Phase 3 (Unit Tests)
    ↓
    ├──→ Phase 4 (Integration Tests) [SEQUENTIAL - needs unit tests]
    └──→ Phase 5 (Documentation) [PARALLEL - can continue]
```

### Optimal Execution Timeline
```
Hour 0-0.5:   Phase 1 (Database Layer)
Hour 0.5-2.5: Phase 2 (Service Layer)
Hour 2.5-5.5: Phase 3 (Unit Tests) + Phase 5 (Documentation) [PARALLEL]
Hour 5.5-7.5: Phase 4 (Integration Tests)
Hour 7.5-8.5: Phase 6 (Final Verification)
```

---

## Risk Mitigation Strategies

### Risk 1: Performance with Deep Nesting
**Mitigation**:
- Integration test with 10+ levels
- Monitor query execution time
- Consider batch operations if needed (future enhancement)

### Risk 2: Transaction Deadlocks
**Mitigation**:
- SQLite handles locking at database level (unlikely)
- No foreign key constraints (soft delete safe)
- Test concurrent deletes in integration tests

### Risk 3: Inconsistent State on Error
**Mitigation**:
- All operations in transaction
- Transaction rollback on any error
- Integration tests verify rollback behavior

### Risk 4: Breaking Changes
**Mitigation**:
- Existing `SoftDelete()` method unchanged
- New `SoftDeleteWithCascade()` is opt-in
- No changes to public API

---

## Success Criteria

### Functional Requirements
- [ ] Cascade soft delete works for nested projects
- [ ] Tasks deleted in all child projects
- [ ] Transaction ensures atomicity
- [ ] Validation before deletion (ErrProjectNotFound)
- [ ] Proper error handling with context

### Quality Requirements
- [ ] 100% test coverage for new code
- [ ] No race conditions detected
- [ ] No compiler warnings
- [ ] Follows project coding standards
- [ ] Documentation complete

### Performance Requirements
- [ ] Delete 100 tasks in < 100ms
- [ ] Support 10 levels of nesting
- [ ] No memory leaks

### Integration Requirements
- [ ] Works with existing transaction manager
- [ ] Compatible with current database schema
- [ ] No breaking changes to existing code

---

## Deliverables

### Code Files
1. `queries/projects.sql` - New SQL query
2. `internal/db/projects.sql.go` - Generated query code (auto)
3. `internal/db/querier.go` - Updated interface (auto)
4. `internal/service/project_service.go` - 2 new methods
5. `internal/service/project_service_test.go` - Unit tests
6. `internal/service/project_service_cascade_integration_test.go` - Integration tests

### Documentation Files
1. `internal/service/project_service.go` - Godoc updates
2. `docs/architecture/03-service-layer.md` - Pattern documentation (if exists)
3. `README.md` or `docs/API.md` - API documentation

### Test Artifacts
1. Coverage report (100% for new code)
2. Race detection report (clean)
3. Integration test results

---

## Post-Implementation Tasks

### Follow-up Tasks (Create in Backlog)
1. **Performance Monitoring**: Add metrics for cascade delete operations
2. **Audit Logging**: Log cascade delete operations for compliance
3. **Soft Delete UI**: Add UI indicators for soft-deleted items
4. **Restore Functionality**: Implement undo/restore for soft-deleted items
5. **Hard Delete Job**: Create background job to permanently delete old soft-deleted items

### Documentation Updates
1. Update user guide with cascade delete behavior
2. Add troubleshooting guide for common errors
3. Document transaction behavior and limitations

---

## Commands Summary

```bash
# Generate SQL code
make sqlc-generate

# Build project
go build ./...

# Run unit tests
go test -v -race ./internal/service/...

# Run with coverage
go test -cover -coverprofile=coverage.out ./internal/service/...
go tool cover -html=coverage.out

# Run integration tests
go test -tags=integration -v ./internal/service/...

# Run all quality checks
make lint
make test

# Manual build and test
make build
./dopadone project create --name "Test"
./dopadone project delete --cascade <project-id>
```

---

## Acceptance Criteria Mapping

| AC | Implementation Phase | Test Coverage |
|----|---------------------|---------------|
| #1: Add SoftDeleteWithCascade method | Phase 2.1 | Unit + Integration |
| #2: Implement softDeleteRecursive helper | Phase 2.2 | Unit + Integration |
| #3: Soft delete all tasks within deleted projects | Phase 2.2 | Unit + Integration |
| #4: Use transaction for atomic cascade delete | Phase 2.1 | Integration |
| #5: Add SoftDeleteTasksByProject SQL query | Phase 1.1 | Integration |
| #6: Write service layer tests with 100% coverage | Phase 3 + 4 | Coverage report |

All acceptance criteria will be completed and verified through the phased implementation.
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
# Task 68.2 Scope: Cascade Soft Delete Service

## Overview
Implement cascade soft delete functionality for projects with nested subprojects. When a parent project is soft deleted, all child projects and their associated tasks should also be soft deleted within a single transaction.

## Technical Analysis

### Current State
- **SoftDelete method**: Exists but only deletes a single project, no cascade
- **HardDelete method**: Already implements cascade delete pattern (lines 286-320)
- **Transaction support**: Available via `s.tm.WithTransaction()`
- **Recursive pattern**: Already exists in `hardDeleteRecursive()` method

### Gap Analysis
1. **Missing SQL query**: Need `SoftDeleteTasksByProject` to soft delete tasks in bulk
2. **Missing service method**: Need `SoftDeleteWithCascade` public method
3. **Missing recursive helper**: Need `softDeleteRecursive` private method
4. **Missing tests**: No tests exist for cascade soft delete

## Implementation Approach

### Phase 1: Database Layer
**File**: `queries/projects.sql`

Add new query:
```sql
-- name: SoftDeleteTasksByProject :exec
UPDATE tasks
SET deleted_at = ?
WHERE project_id = ? AND deleted_at IS NULL;
```

**Why this approach**:
- Follows existing query patterns (see `SoftDeleteProject`, `SoftDeleteTask`)
- Uses `:exec` since we don't need return values (efficient)
- Checks `deleted_at IS NULL` to avoid re-deleting
- ProjectID parameter for bulk delete (performance)

**After adding**: Run `make sqlc-generate` to generate Go code

### Phase 2: Service Layer
**File**: `internal/service/project_service.go`

#### Method 1: SoftDeleteWithCascade (public)
```go
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, id string) error
```

**Responsibilities**:
- Validate project exists (return `ErrProjectNotFound` if not)
- Wrap operation in transaction if `s.tm` available
- Call `softDeleteRecursive` for cascade logic

**Pattern to follow**: Same as `HardDelete()` (lines 286-300)

#### Method 2: softDeleteRecursive (private)
```go
func (s *ProjectService) softDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error
```

**Responsibilities**:
- Get all child projects via `ListProjectsByParent`
- Recursively call itself for each child
- Soft delete tasks in current project via `SoftDeleteTasksByProject`
- Soft delete current project via `SoftDeleteProject`
- Proper error wrapping with context

**Pattern to follow**: Similar to `hardDeleteRecursive()` (lines 302-320)
**Key difference**: Use soft delete instead of hard delete

**Error handling**:
- Wrap errors with `fmt.Errorf("operation: %w", err)` for context
- Use descriptive messages: "list child projects", "cascade delete child", etc.

### Phase 3: Testing Strategy

#### Unit Tests
**File**: `internal/service/project_service_test.go`

Update mock to implement new interface method:
```go
func (m *mockProjectQuerier) SoftDeleteTasksByProject(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
    return nil
}
```

Add `softDeleteTasksByProjectFunc` field to mock for testing.

**Test cases needed**:

1. **Success scenarios**:
   - Single project with no children
   - Project with direct children (1 level)
   - Deeply nested project (3+ levels)
   - Multiple children at same level
   
2. **Validation scenarios**:
   - Non-existent project ID → `ErrProjectNotFound`
   - Already soft-deleted project → should still work (idempotent)
   
3. **Error scenarios**:
   - Database error on GetProjectByID → wrapped error
   - Database error on ListProjectsByParent → wrapped error
   - Database error on SoftDeleteTasksByProject → wrapped error
   - Database error on SoftDeleteProject → wrapped error
   
4. **Transaction scenarios**:
   - With transaction manager → uses transaction
   - Without transaction manager → uses direct repo
   - Error in middle → transaction rolls back (integration test)

5. **Cascade verification**:
   - Verify tasks deleted in correct order (children first)
   - Verify all levels processed
   - Verify correct deleted_at timestamps

**Test pattern**: Use table-driven tests (see existing tests in file)

#### Integration Tests
**New file**: `internal/service/project_service_cascade_integration_test.go`

**Why separate file**: Integration tests need real database, tag with `//go:build integration`

**Test setup**:
- Create in-memory SQLite database
- Run migrations
- Create test data hierarchy (area → subarea → projects → tasks)

**Test cases**:
1. Full cascade flow with real DB
2. Transaction rollback on error
3. Verify data integrity after delete
4. Performance test with deep nesting (10+ levels)

### Phase 4: Code Quality

**Must pass**:
- `go build ./...` - No compilation errors
- `go test ./internal/service/... -v` - All tests pass
- `go test -race ./internal/service/...` - No race conditions
- `go test -cover ./internal/service/...` - 100% coverage for new code

**Code review checklist**:
- [ ] Error messages are descriptive
- [ ] Context propagation (ctx parameter)
- [ ] Transaction handling (both with and without tm)
- [ ] Idempotency (can call multiple times safely)
- [ ] Follows existing code patterns
- [ ] No code duplication
- [ ] Proper godoc comments

## Files to Modify

### New Files
1. `internal/service/project_service_cascade_integration_test.go` - Integration tests

### Modified Files
1. `queries/projects.sql` - Add SoftDeleteTasksByProject query
2. `internal/db/projects.sql.go` - Auto-generated by sqlc
3. `internal/db/querier.go` - Auto-generated by sqlc (interface update)
4. `internal/service/project_service.go` - Add 2 new methods
5. `internal/service/project_service_test.go` - Add unit tests, update mock

## Dependencies

### Internal Dependencies
- `db.Querier` interface (data access)
- `db.TransactionManager` (transaction support)
- `domain` package (error types)

### External Dependencies
- Standard library only (context, database/sql, fmt, time)

### Parallel Work
- **Can run in parallel with**: TASK-68.1 (Confirmation Modal Component)
- **Blocks**: TASK-68.3 (TUI Delete Integration) - needs this service method

## Risks and Mitigations

### Risk 1: Performance with Deep Nesting
**Impact**: Slow cascade delete with 10+ levels
**Mitigation**: 
- Accept performance cost (soft delete is fast enough)
- Add depth limit warning in logs if needed
- Monitor in production

### Risk 2: Transaction Deadlock
**Impact**: Database deadlock with concurrent deletes
**Mitigation**:
- SQLite handles locking at database level
- Unlikely with soft delete (no foreign key constraints)
- Test with concurrent operations

### Risk 3: Inconsistent State on Error
**Impact**: Partial delete if transaction fails
**Mitigation**:
- All operations wrapped in transaction
- Transaction rollback on any error
- Integration tests verify rollback

### Risk 4: Breaking Change
**Impact**: Existing code using SoftDelete won't cascade
**Mitigation**:
- Keep existing `SoftDelete()` method unchanged
- New method `SoftDeleteWithCascade()` is opt-in
- No breaking changes to public API

## Success Criteria

### Functional Requirements
- [ ] Cascade soft delete works for nested projects
- [ ] Tasks deleted in all child projects
- [ ] Transaction ensures atomicity
- [ ] Validation before deletion
- [ ] Proper error handling

### Quality Requirements
- [ ] 100% test coverage for new code
- [ ] No race conditions
- [ ] No compiler warnings
- [ ] Follows project coding standards
- [ ] Documentation comments on public methods

### Performance Requirements
- [ ] Delete 100 tasks in < 100ms
- [ ] Support 10 levels of nesting
- [ ] No memory leaks

## Estimated Effort

**Breakdown**:
- SQL query + sqlc generation: 30 minutes
- Service method implementation: 2 hours
- Unit tests: 3 hours
- Integration tests: 2 hours
- Code review + fixes: 1 hour
- **Total**: 8.5 hours (~1-2 days)

**Assumptions**:
- Developer familiar with codebase
- No blockers or major issues
- Test infrastructure already exists

## Open Questions

1. ✅ **RESOLVED**: Should cascade delete also delete tasks in child projects?
   - **Answer**: Yes, delete all tasks in parent + children

2. ✅ **RESOLVED**: Should we validate project exists before deletion?
   - **Answer**: Yes, return ErrProjectNotFound if not found

3. ✅ **RESOLVED**: Unit tests only or also integration tests?
   - **Answer**: Both unit and integration tests

4. **NEW**: Should we add audit logging for cascade deletes?
   - **Recommendation**: Out of scope for this task, create follow-up task if needed

5. **NEW**: Should we add a depth limit for recursion?
   - **Recommendation**: No limit needed for SQLite, monitor in production

## Next Steps

1. Review and approve this scope
2. Create feature branch: `feature/task-68.2-cascade-soft-delete`
3. Implement Phase 1 (SQL query)
4. Implement Phase 2 (Service methods)
5. Implement Phase 3 (Tests)
6. Run quality checks
7. Create pull request
8. Update task status to Done

## Appendix: Code References

### Existing Pattern: HardDelete
```go
func (s *ProjectService) HardDelete(ctx context.Context, id string) error {
    if s.tm == nil {
        if err := s.hardDeleteRecursive(ctx, s.repo, id); err != nil {
            return err
        }
        return s.repo.HardDeleteProject(ctx, id)
    }

    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        if err := s.hardDeleteRecursive(ctx, tx, id); err != nil {
            return err
        }
        return tx.HardDeleteProject(ctx, id)
    })
}
```

### Existing Pattern: hardDeleteRecursive
```go
func (s *ProjectService) hardDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error {
    children, err := q.ListProjectsByParent(ctx, sql.NullString{String: projectID, Valid: true})
    if err != nil {
        return err
    }

    for _, child := range children {
        if err := s.hardDeleteRecursive(ctx, q, child.ID); err != nil {
            return err
        }
        if err := q.DeleteTasksByProjectID(ctx, child.ID); err != nil {
            return err
        }
        if err := q.HardDeleteProject(ctx, child.ID); err != nil {
            return err
        }
    }

    return q.DeleteTasksByProjectID(ctx, projectID)
}
```

### Golang Patterns to Follow

From golang-patterns skill:
- **Error wrapping**: Use `fmt.Errorf("context: %w", err)` for error chain
- **Context propagation**: Always pass ctx as first parameter
- **Interface acceptance**: Accept `db.Querier` interface, not concrete type
- **Table-driven tests**: Use for comprehensive test coverage
- **Helper functions**: Mark with `t.Helper()` for better error messages

From golang-testing skill:
- **Subtests**: Use `t.Run()` for organized test cases
- **Test helpers**: Create helper functions for common setup
- **Mocking**: Use interface-based mocks for isolation
- **Coverage**: Aim for 100% on new code

From golang-pro skill:
- **Transaction management**: Support both with and without transaction manager
- **Recursive algorithms**: Ensure proper base case and error handling
- **Validation first**: Check preconditions before expensive operations

## Scope Review Complete

**Clarifying Questions Resolved:**
- ✅ Cascade delete includes tasks in all child projects
- ✅ Validate project exists before deletion
- ✅ Both unit and integration tests required
- ✅ No audit logging (out of scope, can add later)
- ✅ No recursion depth limit (SQLite can handle reasonable depth)

**Scope Approved:** Ready for implementation

**Key Decisions:**
1. Follow existing HardDelete pattern for consistency
2. Use transaction for atomicity
3. 100% test coverage required
4. Estimated effort: 8.5 hours (~1-2 days)
<!-- SECTION:NOTES:END -->
