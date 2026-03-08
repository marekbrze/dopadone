---
id: TASK-60.8
title: Database mode auto-detection
status: To Do
assignee: []
created_date: '2026-03-08 19:02'
labels:
  - database
  - configuration
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement automatic database mode detection based on configuration presence. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Auto-detect mode: if turso-url present and db-mode not set, use remote
- [ ] #2 Auto-detect mode: if turso-url + local path set, use embedded replica
- [ ] #3 Auto-detect mode: if only db-path set, use local SQLite (default)
- [ ] #4 Add validation for required configuration per mode
- [ ] #5 Log detected mode at startup for visibility
<!-- AC:END -->
