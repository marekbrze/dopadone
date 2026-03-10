---
id: TASK-60.7
title: 'Integration and refactoring: wire up database abstraction'
status: To Do
assignee: []
created_date: '2026-03-08 19:01'
updated_date: '2026-03-10 07:30'
labels:
  - integration
  - refactoring
  - database
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Integrate database abstraction layer into main application, refactor connection initialization, and update service container. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Refactor cmd/dopa/main.go to use driver factory
- [ ] #2 Update GetDB() to use configured driver
- [ ] #3 Update GetServices() to handle different connection types
- [ ] #4 Ensure all CLI commands work with new abstraction
- [ ] #5 Ensure TUI works with new abstraction
- [ ] #6 Maintain backward compatibility - no breaking changes
<!-- AC:END -->
