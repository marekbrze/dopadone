---
id: TASK-72
title: Fix errcheck lint errors
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-10 12:35'
updated_date: '2026-03-10 14:55'
labels:
  - lint
  - code-quality
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Error return values are not checked in multiple files. Need to handle errors from Close(), Flush(), MarkFlagRequired(), Fprintf(), RemoveAll() calls.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Fix unchecked services.Close() errors in cmd/dopa/areas.go
- [x] #2 Fix unchecked formatter.Flush() errors in cmd/dopa/areas.go, projects.go, subareas.go
- [x] #3 Fix unchecked MarkFlagRequired() errors in cmd/dopa/projects.go, subareas.go, tasks.go
- [x] #4 Fix unchecked fmt.Fprintf() errors in cmd/dopa/tui.go
- [x] #5 Fix unchecked db.Close() and os.RemoveAll() errors in internal/db/*.go tests
- [x] #6 Fix unchecked rows.Close() errors in internal/db/db_test.go and integration_test.go
- [x] #7 Fix unchecked database.Close() errors in internal/tui/*_test.go
- [x] #8 Run make lint to verify all errcheck errors are resolved
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create internal/cli/helpers.go with CloseWithLog() helper for defer cleanup (foundation)
2. Fix CLI command files (parallel work possible):
   - cmd/dopa/areas.go: Use CloseWithLog() for services.Close(), check formatter.Flush()
   - cmd/dopa/projects.go: Panic on MarkFlagRequired() error, check formatter.Flush()
   - cmd/dopa/subareas.go: Panic on MarkFlagRequired() errors, check formatter.Flush()
   - cmd/dopa/tasks.go: Panic on MarkFlagRequired() errors
3. Fix cmd/dopa/tui.go: Use log/slog for error output or ignore Fprintf errors with _=
4. Fix cmd/dopa/main.go: Use CloseWithLog() for db.Close() in migrate commands
5. Fix test files (parallel work possible):
   - internal/db/areas_test.go: Use t.Cleanup() instead of manual cleanup
   - internal/db/db_test.go: Check rows.Close() errors or use t.Cleanup()
   - internal/db/integration_test.go: Check rows.Close() errors
   - internal/tui/complete_test.go: Check database.Close() in t.Cleanup()
   - internal/tui/db_test.go: Check database.Close() in t.Cleanup()
6. Run make lint to verify all errcheck errors resolved
7. Run tests to ensure nothing broken

Dependencies:
- Step 1 must complete before steps 2, 4 (requires helper function)
- Steps 2, 3, 4, 5 can be done in parallel after step 1
- Step 6 must be last (verification)

Test Strategy:
- Run go test ./... after each file modification
- Focus on affected packages to verify behavior unchanged

Documentation:
- No documentation updates needed (internal code quality fix)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- services.Close() errors: Use helper function that logs warning on error
- formatter.Flush() errors: Use ExitWithError consistent with existing pattern
- MarkFlagRequired() errors: Panic in init() since these are programming errors
- Test cleanup: Refactor to use t.Cleanup() with proper error handling
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed all errcheck lint errors across CLI commands and test files.

Changes:
- Created internal/cli/helpers.go with CloseWithLog() helper for defer cleanup
- Fixed cmd/dopa/main.go: db.Close() calls now use cli.CloseWithLog()
- Fixed cmd/dopa/tasks.go: services.Close() uses cli.CloseWithLog(), MarkFlagRequired() panics on error, formatter.Flush() checked
- Fixed cmd/dopa/tui.go: Fprintf errors explicitly ignored with _=
- Fixed internal/cli/db.go: db.Close() on Ping error explicitly ignored
- Fixed internal/db/areas_test.go: os.RemoveAll() and db.Close() in error paths use _= for cleanup
- Fixed internal/db/db_test.go: rows.Close() calls use defer func() { _ = rows.Close() }() pattern
- Fixed internal/db/integration_test.go: rows.Close() calls use defer func() { _ = rows.Close() }() pattern
- Fixed internal/tui/complete_test.go: database.Close() uses defer func() { _ = database.Close() }()
- Fixed internal/tui/db_test.go: database.Close() uses defer func() { _ = database.Close() }()

Test Strategy:
- Verified with make lint - all errcheck errors resolved
- Error handling follows existing patterns (panic for programming errors, log warnings for cleanup)
<!-- SECTION:FINAL_SUMMARY:END -->
