---
id: TASK-60.4
title: Migration compatibility with libSQL
status: Done
assignee:
  - '@ai-agent'
created_date: '2026-03-08 19:01'
updated_date: '2026-03-11 12:28'
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
- [x] #1 Test existing goose migrations against libSQL remote connection
- [x] #2 Test existing goose migrations against libSQL embedded replica
- [x] #3 Implement migration sync: run locally, sync to Turso via embedded replica
- [x] #4 Add migration verification command to check schema consistency
- [x] #5 Document any libSQL-specific migration considerations
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create libSQL migration test suite for AC#1 and AC#2
2. Implement migration sync strategy for AC#3
3. Add migration verification command for AC#4
4. Document libSQL-specific considerations for AC#5
5. Run tests, lint, build to verify quality
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created libSQL migration test suite in internal/migrate/migrate_libsql_test.go
- Implemented migration sync functionality in internal/migrate/sync.go
- Added schema verification in internal/migrate/verify.go
- Added dopa migrate verify CLI command
- Created documentation in docs/TURSO_MIGRATIONS.md
- All tests pass, lint clean, build successful
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented migration compatibility with libSQL for Turso integration.

Changes:
- Created libSQL migration test suite (migrate_libsql_test.go) for testing migrations against embedded replica and remote connections
- Added migration sync functionality (sync.go) that automatically syncs schema to Turso after migrations for embedded replica mode
- Implemented schema verification (verify.go) to check database consistency and detect schema drift
- Added `dopa migrate verify` CLI command to verify schema consistency
- Extended cmd/dopa/main.go with GetDriver() function and sqlDriverWrapper for unified driver interface
- Created comprehensive documentation (docs/TURSO_MIGRATIONS.md) covering libSQL-specific considerations, migration strategies per driver mode, and troubleshooting

Tests:
- go test ./internal/migrate/... -v -short (all pass)
- golangci-lint run ./internal/migrate/... (0 issues)
- go build ./... (successful)

The implementation maintains backward compatibility with existing SQLite local mode while adding full support for Turso remote and embedded replica modes.
<!-- SECTION:FINAL_SUMMARY:END -->
