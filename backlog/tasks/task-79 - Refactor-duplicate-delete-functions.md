---
id: TASK-79
title: Refactor duplicate delete functions
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-10 19:13'
updated_date: '2026-03-10 21:03'
labels:
  - refactor
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
runProjectsDelete and runSubareasDelete have 42 lines of duplicate code. Extract common delete logic into shared helper.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create shared delete helper function
- [x] #2 Refactor runProjectsDelete to use helper
- [x] #3 Refactor runSubareasDelete to use helper
- [x] #4 golangci-lint dupl warning resolved
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze the duplicate pattern in runProjectsDelete and runSubareasDelete functions
   - Both functions share identical structure: get services, check entity exists, soft/hard delete, print success
   - Differences only in: entity name (Project/Subarea), service used, and error type

2. Design a generic delete helper
   - Create Deleteable interface in internal/cli package with methods:
     - GetByID(ctx, id) (to verify existence)
     - SoftDelete(ctx, id) error
     - HardDelete(ctx, id) error
   - Create RunDelete helper function that accepts:
     - Deleteable interface
     - ID
     - permanent flag
     - entity name (for error messages)
     - notFoundErr (for checking specific not found error)

3. Implement the helper in internal/cli/delete.go
   - Add proper error handling
   - Support context propagation
   - Return appropriate errors

4. Refactor existing delete commands
   - Create wrapper types that implement Deleteable interface for each entity
   - Update runProjectsDelete to use helper
   - Update runSubareasDelete to use helper
   - Also refactor runAreasDelete and runTasksDelete (same pattern exists)

5. Write comprehensive tests
   - Unit tests for RunDelete helper
   - Test cases: soft delete, hard delete, not found errors, service errors
   - Test all entity types (area, subarea, project, task)
   - Verify error messages are correct

6. Verification
   - Run golangci-lint to verify dupl warning is resolved
   - Run existing tests to ensure no regressions
   - Manual testing of delete commands

7. Documentation updates
   - Add GoDoc comments to RunDelete explaining usage
   - Document Deleteable interface
   - No new user-facing documentation needed (internal refactor)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
2026-03-10:21:02: Implementation completed successfully. All acceptance criteria met and unit tests passing, golangci-lint dupl warning resolved.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Extracted duplicate delete functions in CLI commands by extracting common delete logic into a shared helper in internal/cli/delete.go. Created Deleteable interface and RunDelete helper function. All four delete commands (projects, subareas, areas, tasks) now use the helper.
<!-- SECTION:FINAL_SUMMARY:END -->
