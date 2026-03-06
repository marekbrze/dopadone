---
id: TASK-27
title: Eliminate code duplication in DB-to-domain converters
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-04 16:59'
updated_date: '2026-03-04 18:09'
labels:
  - architecture
  - refactoring
  - code-quality
dependencies: []
references:
  - internal/tui/converters.go
  - internal/tui/converters_test.go
  - internal/service/area_service.go
  - internal/service/task_service.go
  - internal/service/subarea_service.go
  - internal/service/project_service.go
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Similar dbXxxToDomain conversion functions are repeated across multiple files (converters.go, area_service.go, CLI commands). Create a centralized converter package to eliminate duplication and ensure consistency.

**Scope:**
- Direction: DB-to-domain only (not FromDomain)
- Focus: Pure data transformation, no validation/sanitization
- Pattern: Simple concrete functions (no interfaces for mocking)
- Location: internal/converter package
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create internal/converter package with all DB-to-domain conversions
- [x] #2 All existing tests pass after consolidation
- [x] #3 Remove duplicate converter functions from TUI layer (internal/tui/converters.go)
- [x] #4 Remove duplicate converter functions from service layer (area_service.go, task_service.go, subarea_service.go, project_service.go)
- [x] #5 Fix bug in project_service.go (unused startDate variable)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
PHASE 1: Create internal/converter package
1. Create converter.go with nullStringToString helper and all converter functions
2. Use domain.Parse*() functions for safer parsing
3. Return values (not pointers) for immutability

PHASE 2: Port tests
1. Port existing tests from internal/tui/converters_test.go
2. Add edge case tests

PHASE 3: Update service layer
1. area_service.go - remove 4 duplicate converters
2. task_service.go - remove dbTaskToDomain
3. subarea_service.go - remove dbSubareaToDomain
4. project_service.go - remove dbProjectToDomain + fix startDate bug

PHASE 4: Update TUI layer
1. Delete internal/tui/converters.go and converters_test.go
2. Update commands.go to use new converter package

PHASE 5: Verify
1. Run tests: go test ./...
2. Run linter: golangci-lint run
3. Mark all ACs complete
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
DESIGN DECISIONS (confirmed with user):

1. Return Types: All converters will return VALUES (not pointers)
   - Rationale: Memory efficient, no nil checks needed, immutable domain types
   - Pattern: func dbXxxToDomain(db db.Xxx) domain.Xxx

2. Type Parsing: Use domain.Parse*() functions (not direct casting)
   - Rationale: Safer with validation, handles invalid DB data gracefully
   - Examples: domain.ParseTaskStatus(), domain.ParseColor(), domain.ParsePriority()

3. NULL Handling: Keep nullStringToString helper function
   - Rationale: Cleaner code, reusable across all converters
   - Location: Will be unexported helper in internal/converter package

4. Bug Fix: Will fix unused startDate in project_service.go as part of this task
   - Lines 314-319 in project_service.go
   - Result should be used in returned struct

SCOPE CLARIFICATION:
- CLI layer has NO converters (verified via grep search)
- Original AC #4 updated to reflect actual scope
- Total duplicate functions: 13 across 5 files
- Estimated line reduction: ~400 lines of duplicate code

FILES TO BE CREATED:
- internal/converter/converter.go (~200 lines)
- internal/converter/converter_test.go (~150 lines, ported from TUI)

FILES TO BE MODIFIED:
- internal/service/area_service.go (remove ~94 lines)
- internal/service/task_service.go (remove ~50 lines)
- internal/service/subarea_service.go (remove ~19 lines)
- internal/service/project_service.go (remove ~63 lines, fix bug)
- internal/tui/commands.go (update imports)

FILES TO BE DELETED:
- internal/tui/converters.go (191 lines)
- internal/tui/converters_test.go (moved to converter package)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
All acceptance criteria complete. Ready for final review.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Code compiles without errors
- [x] #2 No regressions in existing functionality
- [x] #3 golangci-lint passes
<!-- DOD:END -->
