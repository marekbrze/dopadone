---
id: TASK-60.4
title: Migration compatibility with libSQL
status: To Do
assignee: []
created_date: '2026-03-08 19:01'
updated_date: '2026-03-10 07:30'
labels:
  - database
  - migrations
  - turso
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Ensure goose migrations work with libSQL drivers and implement migration sync strategy for embedded replicas. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Test existing goose migrations against libSQL remote connection
- [ ] #2 Test existing goose migrations against libSQL embedded replica
- [ ] #3 Implement migration sync: run locally, sync to Turso via embedded replica
- [ ] #4 Add migration verification command to check schema consistency
- [ ] #5 Document any libSQL-specific migration considerations
<!-- AC:END -->
