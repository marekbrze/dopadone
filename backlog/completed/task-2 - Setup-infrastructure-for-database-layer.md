---
id: TASK-2
title: Setup infrastructure for database layer
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 08:25'
updated_date: '2026-03-03 08:51'
labels:
  - infrastructure
  - setup
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Set up directory structure, Goose migrations, and sqlc configuration. This is the foundation task that all other database tasks depend on.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create directory structure: migrations/, queries/, internal/db/, internal/domain/
- [x] #2 Create sqlc.yaml configuration for SQLite with proper settings
- [x] #3 Create initial Goose migration file: 20240301000000_initial_schema.up.sql with all table definitions (areas, subareas, projects)
- [x] #4 Add CHECK constraint: project must have either parent_id or subarea_id
- [x] #5 Add ENUM types (or TEXT with CHECK) for project_status and priority
- [x] #6 Add soft delete support: deleted_at TIMESTAMP NULL on all tables
- [x] #7 Add indexes on: projects(deadline), projects(status,priority), projects(parent_id), projects(subarea_id), subareas(area_id)
- [x] #8 Add metadata fields: created_at, updated_at on all tables; completed_at on projects
- [x] #9 Create down migration: 20240301000000_initial_schema.down.sql
- [x] #10 Write integration test that runs goose up and verifies schema exists
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
TDD IMPLEMENTATION PLAN:

Phase 1: Directory Structure
1. Create migrations/ directory
2. Create queries/ directory
3. Create internal/db/ directory
4. Create internal/domain/ directory
5. Initialize go.mod if not exists

Phase 2: sqlc Configuration
1. Create sqlc.yaml with SQLite configuration
2. Configure generate section with queries/*.sql
3. Configure sql package path (internal/db)
4. Verify sqlc is installed (go install github.com/kyleconroy/sqlc/cmd/sqlc@latest)

Phase 3: Migration Files (Schema)
1. Create 20240301000000_initial_schema.up.sql:
   - areas table (id, name, color, timestamps, deleted_at)
   - subareas table (id, name, area_id FK, color, timestamps, deleted_at)
   - projects table (id, name, description, goal, status ENUM, priority ENUM, progress, dates, color, parent_id FK, subarea_id FK, position, timestamps, deleted_at)
   - CHECK constraints for ENUMs
   - CHECK constraint: (parent_id IS NOT NULL) OR (subarea_id IS NOT NULL)
   - All indexes
2. Create 20240301000000_initial_schema.down.sql:
   - DROP TABLE projects
   - DROP TABLE subareas
   - DROP TABLE areas

Phase 4: Integration Test (TDD)
1. Create internal/db/db_test.go
2. Write TestMigrationUp:
   - Create temp SQLite file
   - Run goose up
   - Query sqlite_master to verify tables exist
   - Verify foreign keys with PRAGMA foreign_key_list
   - Verify indexes with PRAGMA index_list
3. Write TestMigrationDown:
   - Run goose down
   - Verify tables are dropped
4. Run tests: go test ./internal/db -v

Phase 5: Verification
1. Install goose: go install github.com/pressly/goose/v3/cmd/goose@latest
2. Run goose up on test database
3. Verify with sqlite3 CLI: .schema, .tables
4. Verify all tests pass
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created directory structure: migrations/, queries/, internal/db/, internal/domain/
- Created sqlc.yaml configuration for SQLite
- Created combined migration file with goose annotations (-- +goose Up/Down) instead of separate files
- All tables have soft delete support (deleted_at), timestamps (created_at, updated_at), and projects has completed_at
- CHECK constraints for status and priority enums
- CHECK constraint ensuring project has parent_id OR subarea_id
- All required indexes created
- Integration tests verify schema creation and migration rollback
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Set up database infrastructure with Goose migrations, sqlc configuration, and comprehensive integration tests.

### Changes
- Created directory structure: migrations/, queries/, internal/db/, internal/domain/
- Added sqlc.yaml with SQLite configuration for type-safe SQL generation
- Created migration file 20240301000000_initial_schema.sql with:
  - areas, subareas, projects tables with proper foreign keys
  - Soft delete support (deleted_at on all tables)
  - Timestamps (created_at, updated_at) on all tables
  - completed_at on projects
  - CHECK constraints for status (active/completed/on_hold/archived) and priority (low/medium/high/urgent)
  - CHECK constraint ensuring project has parent_id OR subarea_id
  - Indexes on projects(deadline), projects(status,priority), projects(parent_id), projects(subarea_id), subareas(area_id)
- Created integration tests (internal/db/db_test.go) that verify:
  - Goose up/down migrations work correctly
  - All tables and indexes are created
  - Foreign key relationships are established
  - CHECK constraints enforce valid values

### Tests
```bash
go test ./internal/db -v
```
All tests pass.
<!-- SECTION:FINAL_SUMMARY:END -->
