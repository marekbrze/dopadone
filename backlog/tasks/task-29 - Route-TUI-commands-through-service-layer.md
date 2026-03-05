---
id: TASK-29
title: Route TUI commands through service layer
status: To Do
assignee: []
created_date: '2026-03-04 16:59'
updated_date: '2026-03-05 10:12'
labels:
  - architecture
  - refactoring
  - tui
dependencies:
  - TASK-25
  - TASK-27
references:
  - internal/tui/commands.go
  - internal/tui/app.go
  - internal/tui/tui.go
  - internal/tui/commands_test.go
  - internal/service/area_service.go
  - internal/service/subarea_service.go
  - internal/service/project_service.go
  - internal/service/task_service.go
  - 'Related: TASK-25 (services exist)'
  - 'Related: TASK-27 (converter package)'
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
TUI commands in internal/tui/commands.go directly use db.Querier, bypassing the service layer. Refactor to use services for consistent architecture and to enable proper testing of TUI business logic.\n\n**Design Decisions:**\n- Individual service fields in Model (not a container)\n- Move LoadProjectsCmd filtering logic (belongsToSubarea) to ProjectService\n- Create service interfaces for better test mocking\n- Accept services in New() function (dependency injection)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 AreaServiceInterface is defined with all required methods
- [ ] #2 SubareaServiceInterface is defined with all required methods
- [ ] #3 ProjectServiceInterface is defined with all required methods
- [ ] #4 TaskServiceInterface is defined with all required methods
- [ ] #5 Model struct in app.go has individual service interface fields (areaSvc, subareaSvc, projectSvc, taskSvc)
- [ ] #6 All 12 TUI command functions use service interfaces instead of db.Querier
- [ ] #7 ProjectService has new ListBySubareaRecursive method for filtering logic
- [ ] #8 New() in tui.go accepts services as parameters
- [ ] #9 All TUI tests pass with new service-level mocks
- [ ] #10 No direct db.Querier usage in TUI layer (grep verified)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Task-29: Route TUI Commands Through Service Layer
## Implementation Plan & Task Splitting Strategy

### Overview
This is a **LARGE refactoring task** that requires careful planning and should be split into 5 sequential subtasks to:
- Minimize risk through incremental changes
- Enable thorough testing at each phase
- Allow for parallel work on independent parts
- Maintain working state throughout refactoring

### Current State Analysis
- **12 command functions** in `internal/tui/commands.go` use `db.Querier` directly
- **Model struct** has single `repo db.Querier` field (app.go:20)
- **4 services exist**: AreaService (9 methods), SubareaService (9 methods), ProjectService (12 methods), TaskService (14 methods)
- **Services follow consistent pattern**: Dependency injection, domain types, error handling
- **Converter package exists**: Centralized DB-to-domain conversions (task-27)

### Task Splitting Strategy

#### Task-29A: Define Service Interfaces (Foundation) ⏱️ 2-3 hours
**Priority**: HIGH (enables all subsequent work)
**Risk**: LOW
**Dependencies**: TASK-25 (services exist), TASK-27 (converter package exists)
**Blocks**: TASK-29B, TASK-29C

**Objective**: Create service interfaces for mocking and testability

**Deliverables**:
1. Create `internal/service/interfaces.go` with:
   - `AreaServiceInterface` (9 methods)
   - `SubareaServiceInterface` (9 methods)
   - `ProjectServiceInterface` (12 methods + ListBySubareaRecursive)
   - `TaskServiceInterface` (14 methods)

**Implementation Steps**:
```go
// internal/service/interfaces.go
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

// Similar for SubareaServiceInterface, ProjectServiceInterface, TaskServiceInterface
```

**Tests**:
- Verify interfaces are satisfied by concrete implementations
- Compile-time check: `var _ AreaServiceInterface = (*AreaService)(nil)`

**Documentation**:
- Add doc comments to each interface
- Document design decision: "Define interfaces where they're used" pattern

---

#### Task-29B: Add ListBySubareaRecursive Method (Prerequisite) ⏱️ 3-4 hours
**Priority**: HIGH (required for Task-29D)
**Risk**: MEDIUM (new business logic)
**Dependencies**: TASK-29A (interfaces defined), TASK-25 (ProjectService exists)
**Blocks**: TASK-29D

**Objective**: Move belongsToSubarea recursive logic from TUI to ProjectService

**Deliverables**:
1. Add `ListBySubareaRecursive(ctx, subareaID string) ([]domain.Project, error)` to:
   - `internal/service/project_service.go`
   - `internal/service/interfaces.go` (ProjectServiceInterface)

2. Extract and enhance belongsToSubarea logic (commands.go:72-84):
   - Move to ProjectService as private helper
   - Optimize with single database query (ListAllProjects)
   - Build project hierarchy map for recursive filtering

**Tests** (internal/service/project_service_test.go):
- Empty result
- Direct project membership
- Nested project (parent chain)
- Mixed: direct + nested
- Deep nesting (3+ levels)
- Projects in other subareas (excluded)

**Documentation**:
- Document recursive algorithm in godoc
- Add inline comments explaining hierarchy traversal
- Document performance characteristics (O(n) where n = total projects)

---

#### Task-29C: Update Model Structure (Infrastructure) ⏱️ 2-3 hours
**Priority**: HIGH (foundation for command refactoring)
**Risk**: MEDIUM (structural change)
**Dependencies**: TASK-29A (interfaces defined)
**Blocks**: TASK-29D
**Parallel with**: TASK-29B (can work simultaneously)

**Objective**: Replace single `repo` field with 4 service interface fields

**Deliverables**:
1. Update `internal/tui/app.go` Model struct:
```go
type Model struct {
    // Remove: repo db.Querier
    
    // Add service interfaces
    areaSvc    service.AreaServiceInterface
    subareaSvc service.SubareaServiceInterface
    projectSvc service.ProjectServiceInterface
    taskSvc    service.TaskServiceInterface
    
    // ... rest of fields unchanged
}
```

2. Update `internal/tui/tui.go`:
   - Update InitialModel() signature to accept 4 services
   - Update New() signature to accept 4 services

3. Update caller code (cmd/root.go or main.go):
   - Create service instances
   - Pass to tui.New()

**Tests**:
- Update existing TUI tests to create mock services
- Verify Model initialization works correctly

**Documentation**:
- Update architecture docs to reflect service layer usage
- Document dependency injection pattern

---

#### Task-29D: Refactor Load Commands (Phase 1) ⏱️ 4-5 hours
**Priority**: HIGH (core functionality)
**Risk**: MEDIUM (changing core commands)
**Dependencies**: TASK-29B (ListBySubareaRecursive), TASK-29C (Model structure)
**Blocks**: TASK-29E

**Objective**: Refactor 4 load commands to use service layer

**Commands to Refactor**:
1. LoadAreasCmd → AreaServiceInterface.List()
2. LoadSubareasCmd → SubareaServiceInterface.ListByArea()
3. LoadProjectsCmd → ProjectServiceInterface.ListBySubareaRecursive()
4. LoadTasksCmd → TaskServiceInterface.ListByProject()

**Tests** (internal/tui/commands_test.go):
- Create mock implementations of service interfaces
- Table-driven tests for success and error cases
- Verify service method calls with correct parameters

**Documentation**:
- Document service layer benefits (testability, separation of concerns)
- Update inline comments in command functions

---

#### Task-29E: Refactor CRUD Commands (Phase 2) ⏱️ 5-6 hours
**Priority**: HIGH (complete refactoring)
**Risk**: MEDIUM (changing remaining commands)
**Dependencies**: TASK-29D (load commands done)
**Blocks**: None (final task)

**Objective**: Refactor remaining 8 commands to use service layer

**Commands to Refactor**:
1. CreateSubareaCmd → SubareaServiceInterface.Create()
2. CreateProjectCmd → ProjectServiceInterface.Create()
3. CreateTaskCmd → TaskServiceInterface.Create()
4. CreateAreaCmd → AreaServiceInterface.Create()
5. UpdateAreaCmd → AreaServiceInterface.Update()
6. DeleteAreaCmd → AreaServiceInterface.SoftDelete/HardDelete()
7. ReorderAreasCmd → AreaServiceInterface.ReorderAll()
8. LoadAreaStatsCmd → AreaServiceInterface.GetStats()

**Tests** (internal/tui/commands_test.go):
- Table-driven tests for all 8 commands
- Test success and error paths
- Verify service method calls with correct parameters

**Verification**:
```bash
# Ensure no direct db.Querier usage remains
grep -r "db.Querier" internal/tui/commands.go  # Should return nothing
grep -r "repo\." internal/tui/commands.go      # Should return nothing

# Run all tests
go test ./internal/tui/... -v

# Run linter
golangci-lint run ./internal/tui/...
```

**Documentation**:
- Update architecture diagrams to show service layer
- Document migration from direct db.Querier to service layer
- Add inline comments explaining service layer benefits

---

### Task Dependencies Graph

```
TASK-25 (services exist) ──┐
TASK-27 (converters) ──────┤
                            ↓
                      Task-29A: Define Interfaces
                            ↓
                 ┌──────────┴──────────┐
                 ↓                     ↓
         Task-29B: Recursive     Task-29C: Model
                 ↓                     ↓
                 └──────────┬──────────┘
                            ↓
                      Task-29D: Load Commands
                            ↓
                      Task-29E: CRUD Commands
```

### Sequential vs Parallel Execution

**Sequential (Must be done in order)**:
1. TASK-25 → Task-29A (services must exist)
2. Task-29A → Task-29B (interfaces needed)
3. Task-29A → Task-29C (interfaces needed)
4. Task-29B + Task-29C → Task-29D (both prerequisites)
5. Task-29D → Task-29E (load commands first)

**Parallel Work Opportunities**:
- **Task-29B and Task-29C**: Can be developed in parallel after Task-29A
  
- **Within Task-29D**: Load commands can be refactored independently
  - LoadAreasCmd, LoadSubareasCmd, LoadProjectsCmd, LoadTasksCmd
  
- **Within Task-29E**: Commands can be batched
  - Batch 1: Create* commands (4 functions)
  - Batch 2: Update/Delete commands (4 functions)

### Testing Strategy

**Unit Tests** (for each subtask):
1. **Task-29A**: Interface satisfaction checks (compile-time)
2. **Task-29B**: ListBySubareaRecursive unit tests with mock db.Querier
3. **Task-29C**: Model initialization tests with mock services
4. **Task-29D**: Table-driven tests for each load command
5. **Task-29E**: Table-driven tests for each CRUD command

**Integration Tests**:
- Test command flow: TUI → Service → Database
- Verify error propagation
- Test concurrent command execution

**Regression Testing**:
- Run full TUI test suite after each task
- Manual testing of TUI interactions
- Verify no behavioral changes

### Documentation Updates

**After Each Task**:
1. Update inline code comments
2. Update godoc for new methods/interfaces
3. Update architecture decision records (if needed)

**After Task-29E (Final)**:
1. Update README.md architecture section
2. Create migration guide for future refactoring
3. Document service layer patterns and best practices
4. Update CLAUDE.md or AGENTS.md with new architecture

### Risk Mitigation

**High-Risk Areas**:
1. **Task-29C**: Model structure change affects all commands
   - Mitigation: Comprehensive test coverage before refactoring
   
2. **Task-29D**: LoadProjectsCmd with recursive logic
   - Mitigation: Thorough unit tests for edge cases (deep nesting, circular refs)
   
3. **Task-29E**: Breaking existing functionality
   - Mitigation: Incremental refactoring with tests at each step

**Rollback Strategy**:
- Each task is a separate git commit
- Can revert to previous working state if issues found
- Feature flags not needed (direct refactoring)

### Estimated Timeline

- **Task-29A**: 2-3 hours (1 session)
- **Task-29B**: 3-4 hours (1-2 sessions)
- **Task-29C**: 2-3 hours (1 session)
- **Task-29D**: 4-5 hours (2 sessions)
- **Task-29E**: 5-6 hours (2-3 sessions)

**Total**: 16-21 hours (7-9 work sessions)

### Recommendation
**Split into 5 subtasks** for:
- Lower risk through incremental changes
- Better test coverage at each phase
- Easier code review
- Ability to parallelize work
- Clearer progress tracking
- Safer rollback if needed
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Subtasks Created

This task has been split into 5 sequential subtasks for safer incremental refactoring:

1. **TASK-35** (Task-29A): Define Service Interfaces (2-3h)
2. **TASK-36** (Task-29B): Add ListBySubareaRecursive (3-4h)
3. **TASK-37** (Task-29C): Update Model Structure (2-3h) - Can work in parallel with 29B
4. **TASK-38** (Task-29D): Refactor Load Commands (4-5h)
5. **TASK-39** (Task-29E): Refactor CRUD Commands (5-6h)

**Critical Path**: 29A → [29B + 29C] → 29D → 29E
**Total Estimate**: 16-21 hours

See individual subtasks for detailed implementation steps and acceptance criteria.
<!-- SECTION:NOTES:END -->
