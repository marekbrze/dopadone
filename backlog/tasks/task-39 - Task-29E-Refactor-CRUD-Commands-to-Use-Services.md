---
id: TASK-39
title: 'Task-29E: Refactor CRUD Commands to Use Services'
status: To Do
assignee: []
created_date: '2026-03-05 10:11'
labels:
  - architecture
  - refactoring
  - tui
dependencies:
  - TASK-38
references:
  - 'Related: TASK-29 (parent task)'
  - 'Related: TASK-29D (requires load commands done)'
  - internal/tui/commands.go
  - internal/tui/commands_test.go
  - internal/tui/app.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Refactor remaining 8 CRUD commands in internal/tui/commands.go to use service layer instead of direct db.Querier access.

**Dependencies**: TASK-29D (load commands done)
**Blocks**: None (final task in series)

**Objective**: Complete refactoring of all TUI commands to service layer

**Commands to Refactor**:
1. CreateSubareaCmd → SubareaServiceInterface.Create()
2. CreateProjectCmd → ProjectServiceInterface.Create()
3. CreateTaskCmd → TaskServiceInterface.Create()
4. CreateAreaCmd → AreaServiceInterface.Create()
5. UpdateAreaCmd → AreaServiceInterface.Update()
6. DeleteAreaCmd → AreaServiceInterface.SoftDelete/HardDelete()
7. ReorderAreasCmd → AreaServiceInterface.ReorderAll()
8. LoadAreaStatsCmd → AreaServiceInterface.GetStats()

**Changes**:
- Update command function signatures to accept service interfaces
- Replace repo.Xxx() calls with service.Xxx() calls
- Remove converter logic (services return domain types)
- Update command invocations in app.go/update.go

**Testing** (internal/tui/commands_test.go):
- Table-driven tests for all 8 commands
- Test success and error paths
- Verify service method calls with correct parameters

**Verification**:
```bash
# Ensure no direct db.Querier usage remains
grep -r 'db.Querier' internal/tui/commands.go  # Should return nothing
grep -r 'repo\.' internal/tui/commands.go      # Should return nothing

# Run all tests
go test ./internal/tui/... -v

# Run linter
golangci-lint run ./internal/tui/...
```

**Documentation**:
- Update architecture diagrams to show service layer
- Document migration from direct db.Querier to service layer
- Add inline comments explaining service layer benefits
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CreateSubareaCmd uses SubareaServiceInterface
- [ ] #2 CreateProjectCmd uses ProjectServiceInterface
- [ ] #3 CreateTaskCmd uses TaskServiceInterface
- [ ] #4 CreateAreaCmd uses AreaServiceInterface
- [ ] #5 UpdateAreaCmd uses AreaServiceInterface
- [ ] #6 DeleteAreaCmd uses AreaServiceInterface
- [ ] #7 ReorderAreasCmd uses AreaServiceInterface
- [ ] #8 LoadAreaStatsCmd uses AreaServiceInterface
- [ ] #9 All 8 CRUD command tests pass with mocks
- [ ] #10 No direct db.Querier usage in commands.go (grep verified)
- [ ] #11 Full TUI test suite passes
- [ ] #12 golangci-lint passes
<!-- AC:END -->
