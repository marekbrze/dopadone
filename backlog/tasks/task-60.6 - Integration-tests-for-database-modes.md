---
id: TASK-60.6
title: Integration tests for database modes
status: To Do
assignee: []
created_date: '2026-03-08 19:01'
updated_date: '2026-03-10 07:30'
labels:
  - testing
  - integration
  - database
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive integration tests for all three database modes (local SQLite, remote Turso, embedded replica). Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test suite for local SQLite mode - all CRUD operations
- [ ] #2 Test suite for remote Turso mode - connection, queries, error handling
- [ ] #3 Test suite for embedded replica mode - sync, local writes, remote reads
- [ ] #4 Test configuration precedence (CLI > env > config)
- [ ] #5 Test fail-fast behavior on connection failures
- [ ] #6 Test backward compatibility - defaults to local SQLite
- [ ] #7 Use testcontainers or mock Turso server for CI
<!-- AC:END -->
