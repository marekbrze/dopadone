---
id: TASK-60.6.2
title: Config Precedence and Fail-Fast Tests
status: To Do
assignee: []
created_date: '2026-03-11 13:19'
labels:
  - testing
  - integration
  - config
dependencies: []
references:
  - backlog/tasks/task-60.6 - Integration-tests-for-database-modes.md
parent_task_id: TASK-60.6
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create tests for configuration precedence (CLI > env > config) and fail-fast behavior on connection failures. Part of TASK-60.6 integration test suite. Can run in parallel with 60.6.1.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test CLI flags override environment variables
- [ ] #2 Test CLI flags override config file values
- [ ] #3 Test environment variables override config file values
- [ ] #4 Test default values when nothing specified
- [ ] #5 Test remote mode fails fast with invalid URL
- [ ] #6 Test remote mode fails fast with invalid token
- [ ] #7 Test replica mode timeout on unreachable primary
- [ ] #8 Test connection retry logic
<!-- AC:END -->
