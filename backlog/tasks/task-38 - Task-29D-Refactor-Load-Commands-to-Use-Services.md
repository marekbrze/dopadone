---
id: TASK-38
title: 'Task-29D: Refactor Load Commands to Use Services'
status: To Do
assignee: []
created_date: '2026-03-05 10:11'
labels:
  - architecture
  - refactoring
  - tui
dependencies:
  - TASK-36
  - TASK-37
references:
  - 'Related: TASK-29 (parent task)'
  - 'Related: TASK-29B (requires recursive method)'
  - 'Related: TASK-29C (requires Model structure)'
  - internal/tui/commands.go
  - internal/tui/commands_test.go
  - internal/tui/app.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Refactor 4 load commands in internal/tui/commands.go to use service layer instead of direct db.Querier access.

**Dependencies**: TASK-29B (ListBySubareaRecursive), TASK-29C (Model structure)
**Blocks**: TASK-29E

**Objective**: Replace direct db.Querier usage with service layer calls

**Commands to Refactor**:
1. LoadAreasCmd → AreaServiceInterface.List()
2. LoadSubareasCmd → SubareaServiceInterface.ListByArea()
3. LoadProjectsCmd → ProjectServiceInterface.ListBySubareaRecursive()
4. LoadTasksCmd → TaskServiceInterface.ListByProject()

**Changes**:
- Update command function signatures to accept service interfaces
- Replace repo.ListXxx() calls with service.ListXxx() calls
- Remove converter logic (services return domain types)
- Update command invocations in app.go/update.go

**Testing** (internal/tui/commands_test.go):
- Create mock implementations of service interfaces
- Table-driven tests for success and error cases
- Verify service method calls with correct parameters
- Test all 4 load commands

**Verification**:
- Run full TUI test suite
- Manual TUI testing shows correct data loading
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 LoadAreasCmd uses AreaServiceInterface
- [ ] #2 LoadSubareasCmd uses SubareaServiceInterface
- [ ] #3 LoadProjectsCmd uses ProjectServiceInterface.ListBySubareaRecursive()
- [ ] #4 LoadTasksCmd uses TaskServiceInterface
- [ ] #5 All 4 load command tests pass with mocks
<!-- AC:END -->
