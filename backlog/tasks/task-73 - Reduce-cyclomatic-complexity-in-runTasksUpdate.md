---
id: TASK-73
title: Reduce cyclomatic complexity in runTasksUpdate
status: To Do
assignee: []
created_date: '2026-03-10 12:36'
updated_date: '2026-03-10 15:32'
labels:
  - lint
  - code-quality
  - refactor
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The runTasksUpdate function in cmd/dopa/tasks.go has cyclomatic complexity of 34 (limit is 30). Need to refactor by extracting helper functions.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Analyze current runTasksUpdate function structure
- [x] #2 Extract helper functions for complex conditional logic
- [x] #3 Reduce complexity to below 30
- [ ] #4 Ensure all tests pass after refactoring
- [ ] #5 Run make lint to verify gocyclo error is resolved
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create ValidateUpdateFlags() in internal/cli to check if --next/--no-next flags are mutually exclusive
2. Create prepareTaskUpdateParams() helper in cmd/dopa/tasks.go to:
  - Accept existing task, flags, and service
  - Return UpdateTaskParams with resolved values

3. Extract validation to internal/cli:
  - ValidateUpdateFlags() - validate flags provided and flag conflicts
4. Create tests for new helpers

5. Refactor runTasksUpdate to use new helpers

6. Run make lint to verify gocyclo error is resolved
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
1. Created ValidateUpdateFlags() to internal/cli
2. Extracted private helpers has hasUpdate function (reduce complexity). Started unit testing approach
<!-- SECTION:NOTES:END -->
