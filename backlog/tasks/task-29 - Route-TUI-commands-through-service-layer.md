---
id: TASK-29
title: Route TUI commands through service layer
status: To Do
assignee: []
created_date: '2026-03-04 16:59'
updated_date: '2026-03-04 17:00'
labels:
  - architecture
  - refactoring
  - tui
dependencies:
  - TASK-25
  - TASK-27
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
TUI commands in internal/tui/commands.go directly use db.Querier, bypassing the service layer. Refactor to use services for consistent architecture and to enable proper testing of TUI business logic.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TUI commands.go uses ProjectService for project operations
- [ ] #2 TUI commands.go uses TaskService for task operations
- [ ] #3 TUI commands.go uses SubareaService for subarea operations
- [ ] #4 TUI commands.go uses AreaService for area operations
- [ ] #5 All TUI tests pass after refactoring
- [ ] #6 No direct db.Querier usage in TUI layer
<!-- AC:END -->
