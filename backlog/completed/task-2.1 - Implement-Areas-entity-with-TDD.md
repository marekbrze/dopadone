---
id: TASK-2.1
title: Implement Areas entity with TDD
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 08:25'
updated_date: '2026-03-03 08:57'
labels:
  - backend
  - database
  - tdd
dependencies:
  - TASK-2
parent_task_id: TASK-2
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement Areas entity CRUD operations using TDD approach. Areas are top-level organizational units (Work, Home, Personal Projects).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Write unit tests for Area value object validation (Color hex format)
- [x] #2 Write test queries in queries/areas.sql (test-driven): Create, GetByID, List, Update, SoftDelete
- [x] #3 Run sqlc generate and verify internal/db/models.go contains Area struct
- [x] #4 Run sqlc generate and verify internal/db/queries.sql.go contains Area queries
- [x] #5 Write integration test: Create area, fetch by ID, verify all fields
- [x] #6 Write integration test: List areas, verify ordering
- [x] #7 Write integration test: Update area, verify updated_at changes
- [x] #8 Write integration test: Soft delete area, verify deleted_at is set
- [x] #9 Verify all tests pass with go test ./...
- [x] #10 Document Area entity in internal/domain/area.go with domain types
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
TDD IMPLEMENTATION PLAN:

Phase 1: Domain Model (Test First)
1. Create internal/domain/area.go
2. Define Area struct (domain model)
3. Create internal/domain/value_objects.go
4. Define Color value object with validation
5. Write unit tests for Color validation (TestColorValidation)
6. Run tests: go test ./internal/domain -v (should fail initially)
7. Implement Color validation logic
8. Run tests again (should pass)

Phase 2: SQL Queries (Test-Driven)
1. Create queries/areas.sql
2. Write Create query with parameters
3. Write GetByID query
4. Write List query (active areas, ordered by name)
5. Write Update query
6. Write SoftDelete query
7. Create placeholder test file: queries/areas_test.go (will use sqlc generated code)

Phase 3: Code Generation
1. Run sqlc generate
2. Verify internal/db/models.go has Area struct
3. Verify internal/db/queries.sql.go has all area queries
4. Run go build ./internal/db to verify compilation

Phase 4: Integration Tests
1. Create internal/db/areas_test.go
2. Write TestCreateArea:
   - Setup test DB
   - Create area
   - Fetch by ID
   - Verify all fields match
3. Write TestListAreas:
   - Create multiple areas
   - Call List
   - Verify ordering and soft-delete filtering
4. Write TestUpdateArea:
   - Create area
   - Update name/color
   - Verify updated_at changed
5. Write TestSoftDeleteArea:
   - Create area
   - Soft delete
   - Verify deleted_at is set
   - Verify it doesn't appear in List
6. Run all tests: go test ./internal/db -v

Phase 5: Verification
1. All unit tests pass
2. All integration tests pass
3. Code compiles without errors
4. sqlc generated code is clean
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Areas entity with full TDD approach:

**Domain Layer:**
- Created `internal/domain/value_objects.go` with Color hex validation (#RRGGBB format)
- Created `internal/domain/area.go` with Area struct and NewArea constructor
- Unit tests for Color validation (13 test cases) in `value_objects_test.go`

**Database Layer:**
- Created SQL queries in `queries/areas.sql`: CreateArea, GetAreaByID, ListAreas, UpdateArea, SoftDelete
- Generated sqlc code in `internal/db/models.go` (Area struct) and `internal/db/areas.sql.go` (queries)
- Integration tests in `internal/db/areas_test.go`:
  - TestCreateArea: Create and fetch by ID, verify all fields
  - TestListAreas: Verify alphabetical ordering by name
  - TestUpdateArea: Update name/color, verify updated_at changes
  - TestSoftDeleteArea: Soft delete, verify deleted_at set and excluded from List

**All tests pass:** `go test ./...`
<!-- SECTION:FINAL_SUMMARY:END -->
