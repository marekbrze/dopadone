---
id: TASK-34
title: Add missing ListAll and ListByPriority service methods
status: Done
assignee:
  - '@assistant'
created_date: '2026-03-04 20:10'
updated_date: '2026-03-04 20:21'
labels:
  - service-layer
  - backend
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add ListAll() methods to TaskService and SubareaService, and ListByPriority() to ProjectService. These methods are needed for CLI refactoring to fully eliminate direct db.Querier usage.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 TaskService.ListAll() returns all non-deleted tasks
- [x] #2 SubareaService.ListAll() returns all non-deleted subareas
- [x] #3 ProjectService.ListByPriority() returns projects filtered by priority
- [x] #4 All methods follow existing service patterns (domain types, error handling)
- [x] #5 Unit tests for all new methods with 80%+ coverage
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add SQL queries (ListAllTasks, ListAllSubareas, ListProjectsByPriority)
2. Run sqlc generate to create Go code
3. Add TaskService.ListAll() method
4. Add SubareaService.ListAll() method
5. Add ProjectService.ListByPriority() method
6. Update mock queriers in tests to support new repo methods
7. Add unit tests for all three new service methods
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Added ListAll and ListByPriority service methods to support CLI refactoring.

Changes:
- Added ListAllTasks SQL query (queries/tasks.sql)
- Added ListAllSubareas SQL query (queries/subareas.sql)
- Added ListProjectsByPriority SQL query (queries/projects.sql)
- Added TaskService.ListAll() method
- Added SubareaService.ListAll() method
- Added ProjectService.ListByPriority() method
- Updated mock queriers in all test files to support new repository methods
- Added unit tests for all three new service methods with 80%+ coverage

All service tests pass successfully.
<!-- SECTION:FINAL_SUMMARY:END -->
