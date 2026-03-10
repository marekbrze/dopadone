---
id: TASK-71
title: Fix lint errors
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-10 12:03'
updated_date: '2026-03-10 12:12'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix all 50 lint errors reported by golangci-lint to improve code quality and CI pipeline compliance.

Categories to fix:
- errcheck (34): Unchecked return values on defer Close(), MarkFlagRequired, Flush, fmt.Fprintf
- goconst (8): Repeated string literals (json, yaml, ctrl+c, enter, esc, Task, root, windows)
- gocyclo (2): High complexity in runTasksUpdate (34) and Model.Update (94)
- ineffassign (3): Assignments to err/cmd that are never used
- staticcheck (2): Empty branch, ineffective field assignment
- unused (2): Dead code (projectsToRows, containsSubstring functions)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 All 34 errcheck warnings resolved with proper error handling (log or handle explicitly, not silent ignore)
- [ ] #2 All 8 goconst warnings resolved by extracting string literals into constants
- [ ] #3 runTasksUpdate function refactored to reduce complexity from 34 to <=30
- [ ] #4 Model.Update function refactored by extracting message handlers to reduce complexity from 94 to <=30
- [ ] #5 All 3 ineffassign warnings resolved by using or removing assignments
- [ ] #6 Both staticcheck warnings resolved (SA9003 empty branch, SA4005 ineffective assignment)
- [ ] #7 Both unused function warnings resolved (remove or use projectsToRows, containsSubstring)
- [ ] #8 All tests pass after changes (go test ./...)
- [ ] #9 golangci-lint run reports 0 issues
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Task-71: Fix Lint Errors

## Task Splitting Strategy

This task will be split into **6 subtasks** for safer, incremental implementation:

### Phase 1: Quick Wins (Parallel)
**Can be done simultaneously:**
- **71-1**: Extract string constants (goconst) - 1-2h
- **71-2**: Fix ineffassign warnings - 30m
- **71-3**: Fix staticcheck warnings - 30m
- **71-4**: Remove unused code - 20m

### Phase 2: Error Handling (Sequential)
- **71-5**: Fix errcheck warnings - 2-3h
  - Depends on: 71-1 (uses constants)
  - Fix 33 unchecked return values
  - Categories: defer Close(), MarkFlagRequired, Flush, fmt.Fprintf

### Phase 3: Complexity Refactoring (Sequential, Critical)
- **71-6**: Refactor high complexity functions - 4-6h
  - Depends on: 71-1, 71-5
  - runTasksUpdate: complexity 34 → ≤30
  - Model.Update: complexity 94 → ≤30 (MAJOR REFACTORING)
  - Extract handlers using Elm architecture pattern

## Execution Order

**Sequential Path:**
1. Start Phase 1 tasks (71-2, 71-3 first, then 71-1, 71-4)
2. Complete 71-5 (error handling)
3. Complete 71-6 (complexity - most critical)

**Parallel Opportunities:**
- Phase 1 tasks can run concurrently
- After 71-1 complete, 71-2/71-3/71-4 can finish in parallel

## Critical Refactoring: Model.Update

The Model.Update function has complexity 94 (target ≤30):

**Strategy:** Extract message handlers
- Create update_key.go (keyboard handlers)
- Create update_mouse.go (mouse handlers)
- Create update_resize.go (resize handlers)
- Create update_data.go (data loading handlers)
- Create update_modal.go (modal handlers)

**Testing:** Extensive manual TUI testing required
- All navigation
- All modals (create/edit/delete)
- All keyboard shortcuts
- All mouse interactions

## Testing Strategy

**After each subtask:**
- Run: `go test ./...`
- Verify: `golangci-lint run`
- Check specific lint category resolved

**After 71-6 (critical):**
- Full TUI manual testing
- Integration tests
- Race detector: `go test -race ./...`

## Documentation Updates

- docs/architecture/06-cli-layer.md: Add error handling patterns
- docs/TUI.md: Document new handler architecture
- internal/constants/README.md: Explain constants organization

## Success Criteria

- [ ] All 50 lint errors resolved
- [ ] golangci-lint run reports 0 issues
- [ ] All tests pass
- [ ] TUI fully functional
- [ ] No regressions
- [ ] Documentation updated

## Total Estimated Time

8-12 hours over 2 weeks
<!-- SECTION:PLAN:END -->
