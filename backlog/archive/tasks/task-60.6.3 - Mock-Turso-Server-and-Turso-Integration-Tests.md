---
id: TASK-60.6.3
title: Mock Turso Server and Turso Integration Tests
status: To Do
assignee: []
created_date: '2026-03-11 13:20'
updated_date: '2026-03-11 14:23'
labels:
  - testing
  - integration
  - turso
  - mock
dependencies: []
references:
  - backlog/tasks/task-60.6 - Integration-tests-for-database-modes.md
parent_task_id: TASK-60.6
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create mock Turso server for CI testing and comprehensive integration tests for remote and replica modes. Part of TASK-60.6 integration test suite. Depends on 60.6.2 for fail-fast test patterns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Create in-process mock HTTP server simulating Turso API
- [ ] #2 Test remote mode connection with valid/invalid credentials
- [ ] #3 Test remote mode query execution and transactions
- [ ] #4 Test replica mode initial sync on connect
- [ ] #5 Test replica mode periodic auto-sync
- [ ] #6 Test replica mode manual sync trigger
- [ ] #7 Test replica mode offline operation
- [ ] #8 Test replica mode local writes visible immediately
- [ ] #9 All Turso tests run without external dependencies in CI
<!-- AC:END -->
