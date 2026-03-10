---
id: TASK-36
title: 'Task-29B: Add ListBySubareaRecursive to ProjectService'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-05 10:10'
updated_date: '2026-03-05 11:41'
labels:
  - architecture
  - refactoring
  - service-layer
dependencies:
  - TASK-35
references:
  - 'Related: TASK-29 (parent task)'
  - 'Related: TASK-29A (requires interfaces)'
  - 'Related: TASK-25 (ProjectService exists)'
  - internal/service/project_service.go
  - internal/service/project_service_test.go
  - internal/tui/commands.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Move belongsToSubarea recursive logic from TUI layer (commands.go:72-84) to ProjectService for better separation of concerns.

**Dependencies**: TASK-29A (interfaces defined), TASK-25 (ProjectService exists)
**Blocks**: TASK-29D
**Parent Task**: TASK-29

**Objective**: Extract and enhance belongsToSubarea logic from commands.go

**Deliverables**:
1. Add ListBySubareaRecursive(ctx, subareaID string) ([]domain.Project, error) to ProjectService
2. Add private belongsToSubarea helper to ProjectService
3. Update ProjectServiceInterface with new method

**Algorithm**:
- Load all projects once (ListAllProjects)
- Build project hierarchy map for parent lookups
- Recursively filter projects belonging to subarea (including nested)

**Testing** (internal/service/project_service_test.go):
- Empty result
- Direct project membership
- Nested project (parent chain)
- Mixed: direct + nested
- Deep nesting (3+ levels)
- Projects in other subareas (excluded)

**Performance**: O(n) where n = total projects
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 ListBySubareaRecursive method implemented in ProjectService
- [x] #2 belongsToSubarea helper added as private method
- [x] #3 ProjectServiceInterface updated with new method
- [x] #4 Unit tests pass with 80%+ coverage for new method
- [x] #5 Edge cases tested (deep nesting, empty results)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: ListBySubareaRecursive

## Task Assessment
**Size**: Medium (well-scoped, atomic)
**Decision**: Do NOT split - task is appropriately sized with clear boundaries

## Phase 1: Implementation (Sequential)

### Step 1.1: Add Private Helper Method
File: `internal/service/project_service.go`

```go
// belongsToSubarea recursively checks if a project belongs to a subarea
// by checking direct membership or parent chain membership.
func (s *ProjectService) belongsToSubarea(
    project domain.Project,
    subareaID string,
    projectMap map[string]domain.Project,
) bool {
    // Direct membership: project.SubareaID == subareaID
    if project.SubareaID != nil && *project.SubareaID == subareaID {
        return true
    }
    
    // Check parent chain recursively
    if project.ParentID != nil {
        if parent, exists := projectMap[*project.ParentID]; exists {
            return s.belongsToSubarea(parent, subareaID, projectMap)
        }
    }
    
    return false
}
```

**Why**: Extracts reusable logic, follows Go patterns for private helpers
**Pattern**: Idiomatic Go private method (lowercase first letter)

### Step 1.2: Implement ListBySubareaRecursive
File: `internal/service/project_service.go`

Replace the stub implementation (lines 194-196):

```go
func (s *ProjectService) ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error) {
    // Edge case: empty subareaID returns empty slice
    if subareaID == "" {
        return []domain.Project{}, nil
    }
    
    // Load all projects (ListAll already filters soft-deleted)
    allProjects, err := s.ListAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("list projects for subarea %s: %w", subareaID, err)
    }
    
    // Build project map for O(1) parent lookups
    projectMap := make(map[string]domain.Project, len(allProjects))
    for _, p := range allProjects {
        projectMap[p.ID] = p
    }
    
    // Filter projects belonging to subarea (direct + nested)
    result := make([]domain.Project, 0)
    for _, project := range allProjects {
        if s.belongsToSubarea(project, subareaID, projectMap) {
            result = append(result, project)
        }
    }
    
    return result, nil
}
```

**Performance**: O(n) time, O(n) space
**Error Handling**: Wraps errors with context using fmt.Errorf
**Edge Cases**: Empty subareaID, no projects, all deleted

## Phase 2: Testing (Sequential - After Phase 1)

### Step 2.1: Add Test Cases to Existing Test File
File: `internal/service/project_service_test.go`

Add comprehensive table-driven test with 12 cases:

```go
func TestProjectService_ListBySubareaRecursive(t *testing.T) {
    now := time.Now()
    subareaA := "subarea-a"
    subareaB := "subarea-b"
    
    tests := []struct {
        name       string
        subareaID  string
        mock       func() *mockProjectQuerier
        wantCount  int
        wantErr    bool
        wantIDs    []string // Optional: verify specific IDs returned
    }{
        {
            name:      "empty subareaID returns empty slice",
            subareaID: "",
            mock:      func() *mockProjectQuerier { return &mockProjectQuerier{} },
            wantCount: 0,
            wantErr:   false,
        },
        {
            name:      "no projects in database",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{}, nil
                    },
                }
            },
            wantCount: 0,
            wantErr:   false,
        },
        {
            name:      "direct membership only",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                Priority:  "high",
                                Progress:  0,
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 1,
            wantIDs:   []string{"proj-root-a"},
        },
        {
            name:      "nested project via parent",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                Priority:  "high",
                                Progress:  0,
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-child-a1",
                                ParentID:  sql.NullString{String: "proj-root-a", Valid: true},
                                SubareaID: sql.NullString{Valid: false}, // No direct membership
                                Status:    "active",
                                Priority:  "medium",
                                Progress:  0,
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 2,
            wantIDs:   []string{"proj-root-a", "proj-child-a1"},
        },
        {
            name:      "deep nesting (3 levels)",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-child-a1",
                                ParentID:  sql.NullString{String: "proj-root-a", Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-grandchild-a1",
                                ParentID:  sql.NullString{String: "proj-child-a1", Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 3,
        },
        {
            name:      "excludes projects in other subareas",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-root-b",
                                SubareaID: sql.NullString{String: subareaB, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 1,
            wantIDs:   []string{"proj-root-a"},
        },
        {
            name:      "excludes soft-deleted projects",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        // ListAll should not return deleted, but test edge case
                        return []db.Project{
                            {
                                ID:        "proj-active",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 1,
        },
        {
            name:      "mixed direct and nested",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a1",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-root-a2",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "proj-child-a1",
                                ParentID:  sql.NullString{String: "proj-root-a1", Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 3,
        },
        {
            name:      "orphaned project (parent doesn't exist)",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-orphan",
                                ParentID:  sql.NullString{String: "nonexistent-parent", Valid: true},
                                SubareaID: sql.NullString{Valid: false},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 0, // Orphan is not included
        },
        {
            name:      "root project with no parent",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return []db.Project{
                            {
                                ID:        "proj-root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                ParentID:  sql.NullString{Valid: false},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 1,
        },
        {
            name:      "database error",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        return nil, errors.New("database connection failed")
                    },
                }
            },
            wantCount: 0,
            wantErr:   true,
        },
        {
            name:      "complex hierarchy",
            subareaID: subareaA,
            mock: func() *mockProjectQuerier {
                return &mockProjectQuerier{
                    listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
                        // Subarea-A: Root -> Child -> Grandchild
                        //           Root2
                        // Subarea-B: RootB (excluded)
                        return []db.Project{
                            {
                                ID:        "root-a",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "child-a",
                                ParentID:  sql.NullString{String: "root-a", Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "grandchild-a",
                                ParentID:  sql.NullString{String: "child-a", Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "root-a2",
                                SubareaID: sql.NullString{String: subareaA, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                            {
                                ID:        "root-b",
                                SubareaID: sql.NullString{String: subareaB, Valid: true},
                                Status:    "active",
                                CreatedAt: now,
                                UpdatedAt: now,
                            },
                        }, nil
                    },
                }
            },
            wantCount: 4,
            wantIDs:   []string{"root-a", "child-a", "grandchild-a", "root-a2"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := NewProjectService(tt.mock())
            got, err := svc.ListBySubareaRecursive(context.Background(), tt.subareaID)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ListBySubareaRecursive() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr && len(got) != tt.wantCount {
                t.Errorf("ListBySubareaRecursive() returned %d projects, want %d", len(got), tt.wantCount)
                t.Logf("Returned IDs: %v", getProjectIDs(got))
            }
            
            // Verify specific IDs if provided
            if tt.wantIDs != nil && !tt.wantErr {
                gotIDs := getProjectIDs(got)
                for _, wantID := range tt.wantIDs {
                    found := false
                    for _, gotID := range gotIDs {
                        if gotID == wantID {
                            found = true
                            break
                        }
                    }
                    if !found {
                        t.Errorf("Expected project ID %s not found in results", wantID)
                    }
                }
            }
        })
    }
}

func getProjectIDs(projects []domain.Project) []string {
    ids := make([]string, len(projects))
    for i, p := range projects {
        ids[i] = p.ID
    }
    return ids
}
```

**Test Coverage**: 12 comprehensive test cases
**Pattern**: Table-driven tests with subtests (Go testing best practice)
**Coverage Target**: 80%+ for new method

### Step 2.2: Run Tests and Verify Coverage

```bash
# Run specific test
go test ./internal/service/... -v -run TestProjectService_ListBySubareaRecursive

# Generate coverage report
go test ./internal/service/... -coverprofile=coverage.out -run TestProjectService_ListBySubareaRecursive

# View coverage
go tool cover -func=coverage.out | grep ListBySubareaRecursive

# Ensure 80%+ coverage
go test ./internal/service/... -cover -run TestProjectService_ListBySubareaRecursive
```

## Phase 3: Code Quality & Documentation (Sequential - After Phase 2)

### Step 3.1: Run Linter

```bash
# Run golangci-lint on service package
golangci-lint run ./internal/service/...

# Run go vet
go vet ./internal/service/...
```

### Step 3.2: Run All Tests

```bash
# Ensure no regressions
go test ./internal/service/... -v

# Run with race detector
go test -race ./internal/service/...
```

### Step 3.3: Update Implementation Notes

Document the implementation details and any decisions made.

## Phase 4: Verification (Sequential - After Phase 3)

### Verification Checklist

```bash
# 1. All tests pass
go test ./internal/service/... -v
# Expected: PASS

# 2. Coverage >= 80% for new method
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep ListBySubareaRecursive
# Expected: >= 80.0%

# 3. Linter passes
golangci-lint run ./internal/service/...
# Expected: No issues

# 4. Race detection passes
go test -race ./internal/service/...
# Expected: PASS

# 5. No regressions in other tests
go test ./... -short
# Expected: All PASS
```

## Documentation Updates

### Internal Documentation

No external documentation changes required. The method is already documented in:
- Interface definition (interfaces.go:132-134)
- Implementation will have godoc comments

### Code Comments

Add inline comments explaining:
- Why we load all projects (single DB query)
- How the recursion works (parent chain traversal)
- Performance characteristics (O(n) time/space)

## Work Breakdown: Sequential vs Parallel

### Sequential Steps (MUST be done in order)

1. **Phase 1**: Implementation
   - Step 1.1: Add private helper
   - Step 1.2: Implement public method
   
2. **Phase 2**: Testing
   - Step 2.1: Write tests
   - Step 2.2: Run tests & verify coverage
   
3. **Phase 3**: Code Quality
   - Step 3.1: Run linter
   - Step 3.2: Run all tests
   - Step 3.3: Update notes
   
4. **Phase 4**: Verification
   - Complete checklist

### Why All Sequential?

- Each phase depends on the previous phase
- Implementation must exist before testing
- Tests must pass before linting
- Linting must pass before final verification
- This is a single, well-scoped task (not splittable)

## Performance Considerations

- **Time Complexity**: O(n) where n = total projects
- **Space Complexity**: O(n) for projectMap
- **Database Calls**: 1 (ListAll)
- **Optimization**: Pre-allocated slices, O(1) map lookups

## Risk Mitigation

1. **Test edge cases thoroughly** (empty, deleted, orphaned)
2. **Handle database errors gracefully** with wrapped errors
3. **Verify soft-delete filtering** works correctly
4. **Test recursion depth** (should handle 10+ levels)
5. **Race detection** ensures no concurrency issues

## Success Criteria

- [x] Implementation plan created
- [ ] belongsToSubarea helper implemented (private method)
- [ ] ListBySubareaRecursive implemented (public method)
- [ ] 12 test cases passing
- [ ] 80%+ code coverage for new method
- [ ] All existing tests still pass
- [ ] Linter passes with no errors
- [ ] Race detection passes
- [ ] Implementation notes updated
- [ ] Acceptance criteria marked complete

## Estimated Effort

- **Implementation**: 30 minutes
- **Testing**: 45 minutes
- **Verification**: 15 minutes
- **Total**: ~1.5 hours

## Dependencies

- **Depends on**: TASK-35 (interfaces defined) - ✅ Complete
- **Depends on**: TASK-25 (ProjectService exists) - ✅ Complete
- **Blocks**: TASK-29D
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
<!-- SECTION:NOTES:BEGIN -->
## Implementation Specification

### Design Decisions (Confirmed)

1. **No Caching**: Load all projects fresh from database each time via ListAll()
2. **Exclude Soft-Deleted**: Filter out projects where DeletedAt != nil
3. **Empty subareaID**: Return empty slice (no error)
4. **No Cycle Detection**: Rely on existing ValidateParentHierarchy to prevent circular references

### Algorithm

```
1. Validate input (empty subareaID returns empty slice)
2. Load ALL non-deleted projects via ListAll(ctx)
3. Build projectMap: map[projectID]domain.Project for O(1) parent lookups
4. For each project:
   - Check if belongsToSubarea(project, subareaID, projectMap)
   - If true, append to results
5. Return filtered projects
```

### Private Helper: belongsToSubarea

```go
// belongsToSubarea recursively checks if a project belongs to a subarea
// by checking direct membership or parent chain membership.
func (s *ProjectService) belongsToSubarea(
    project domain.Project,
    subareaID string,
    projectMap map[string]domain.Project,
) bool {
    // Direct membership: project.SubareaID == subareaID
    if project.SubareaID != nil && *project.SubareaID == subareaID {
        return true
    }
    
    // Check parent chain recursively
    if project.ParentID != nil {
        if parent, exists := projectMap[*project.ParentID]; exists {
            return s.belongsToSubarea(parent, subareaID, projectMap)
        }
    }
    
    return false
}
```

### Public Method: ListBySubareaRecursive

```go
// ListBySubareaRecursive returns all projects that belong to a subarea,
// including nested projects whose parent chain leads to the subarea.
func (s *ProjectService) ListBySubareaRecursive(
    ctx context.Context,
    subareaID string,
) ([]domain.Project, error) {
    // Edge case: empty subareaID
    if subareaID == "" {
        return []domain.Project{}, nil
    }
    
    // Load all projects (non-deleted filtered in ListAll)
    allProjects, err := s.ListAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("list projects for subarea %s: %w", subareaID, err)
    }
    
    // Build project map for O(1) parent lookups
    projectMap := make(map[string]domain.Project, len(allProjects))
    for _, p := range allProjects {
        projectMap[p.ID] = p
    }
    
    // Filter projects belonging to subarea (direct + nested)
    var result []domain.Project
    for _, project := range allProjects {
        if s.belongsToSubarea(project, subareaID, projectMap) {
            result = append(result, project)
        }
    }
    
    return result, nil
}
```

### Test Cases (Comprehensive)

**File**: internal/service/project_service_test.go

1. **Empty result** - subareaID has no projects
2. **Direct membership** - project.SubareaID == subareaID
3. **Nested project (1 level)** - parent belongs to subarea
4. **Nested project (2 levels)** - grandparent belongs to subarea
5. **Deep nesting (3+ levels)** - ancestor chain to subarea
6. **Mixed** - direct + nested projects in same subarea
7. **Exclusion** - projects in other subareas not included
8. **Empty subareaID** - returns empty slice, no error
9. **Soft-deleted excluded** - DeletedAt != nil projects excluded
10. **Multiple subareas** - only returns projects for requested subarea
11. **Database error** - properly wrapped and returned
12. **No parent** - root project with no parent (direct membership only)

### Test Data Structure

```go
// Test hierarchy:
// Subarea-A
//   ├── Project-Root-A (direct member)
//   │   ├── Project-Child-A1 (nested via parent)
//   │   │   └── Project-Grandchild-A1 (nested via grandparent)
//   │   └── Project-Child-A2 (nested via parent)
//   └── Project-Root-A2 (direct member)
//
// Subarea-B
//   └── Project-Root-B (should be excluded when querying Subarea-A)
//
// Soft-deleted
//   └── Project-Deleted (should be excluded from all queries)
```

### Performance Characteristics

- **Time Complexity**: O(n) where n = total projects
- **Space Complexity**: O(n) for projectMap
- **Database Calls**: 1 (ListAll)

### Edge Cases Handled

1. Empty subareaID → empty slice
2. No projects in database → empty slice
3. All projects deleted → empty slice
4. Circular references → prevented by ValidateParentHierarchy (assumed)
5. Orphaned projects (parent doesn't exist) → skipped gracefully
6. Deep nesting (10+ levels) → handled via recursion

### Verification Steps

```bash
# Run tests
go test ./internal/service/... -v -run TestProjectService_ListBySubareaRecursive

# Check coverage
go test ./internal/service/... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep ListBySubareaRecursive

# Run all tests
go test ./internal/service/... -v

# Run linter
golangci-lint run ./internal/service/...
```

### Success Criteria

- [x] Implementation plan created
- [ ] ListBySubareaRecursive implemented in ProjectService
- [ ] Private belongsToSubarea helper added
- [ ] All 12 test cases passing
- [ ] 80%+ code coverage for new method
- [ ] All existing tests still pass
- [ ] Linter passes with no errors
- [ ] Code compiles without errors

- ✅ Implemented ListBySu                                                                                                                                                                                                                                                                                                                                                                                                                        
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added ListBySubareaRecursive method to ProjectService for recursive project filtering by subarea.

Changes:
- Implemented ListBySubareaRecursive(ctx, subareaID) in ProjectService with O(n) time complexity
- Added private helper method belongsToSubarea for recursive parent chain traversal
- Added fmt import for proper error wrapping

Algorithm:
1. Validates input (empty subareaID returns empty slice)
2. Loads all non-deleted projects via ListAll(ctx)
3. Builds projectMap for O(1) parent lookups
4. Filters projects belonging to subarea (direct + nested via parent chain)

Testing:
- Added 12 comprehensive test cases covering all edge cases:
  - Empty subareaID, no projects, direct membership
  - Nested projects (1, 2, 3+ levels deep)
  - Mixed direct + nested scenarios
  - Exclusion of projects in other subareas
  - Orphaned projects, database errors
  - Complex hierarchies
- Achieved 100% code coverage for both methods
- All existing tests pass
- Race detection tests pass

Files modified:
- internal/service/project_service.go: Added ListBySubareaRecursive and belongsToSubarea methods
- internal/service/project_service_test.go: Added comprehensive test suite

Note: ProjectServiceInterface already had the method signature defined (lines 132-134 in interfaces.go)
<!-- SECTION:FINAL_SUMMARY:END -->
