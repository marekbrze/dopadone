---
id: TASK-2.1.1
title: Implement Subareas entity with TDD
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 08:26'
updated_date: '2026-03-03 09:01'
labels:
  - backend
  - database
  - tdd
dependencies:
  - TASK-2.1
parent_task_id: TASK-2.1
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement Subareas entity CRUD operations using TDD. Subareas are second-level organization children of Areas (e.g., Home → Car, Work → Marketing).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Write test queries in queries/subareas.sql: Create, GetByID, ListByArea, Update, SoftDelete
- [x] #2 Run sqlc generate and verify internal/db/models.go contains Subarea struct
- [x] #3 Run sqlc generate and verify internal/db/queries.sql.go contains Subarea queries
- [x] #4 Write integration test: Create subarea with FK to area, verify relationship
- [x] #5 Write integration test: List subareas by area_id, verify correct filtering
- [x] #6 Write integration test: Update subarea, verify updated_at changes
- [x] #7 Write integration test: Soft delete subarea, verify deleted_at is set
- [x] #8 Write unit test: Color inheritance from parent Area if subarea.color is null
- [x] #9 Verify all tests pass with go test ./...
- [x] #10 Document Subarea entity in internal/domain/subarea.go
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
TDD IMPLEMENTATION PLAN:

Phase 1: Domain Model
1. Add Subarea struct to internal/domain/subarea.go
2. Add color inheritance logic (null color inherits from parent Area)
3. Write unit test: TestSubareaColorInheritance
4. Implement inheritance logic
5. Run tests: go test ./internal/domain -v

Phase 2: SQL Queries (Test-Driven)
1. Create queries/subareas.sql
2. Write Create query (requires area_id FK)
3. Write GetByID query (join with areas for color inheritance)
4. Write ListByArea query (filter by area_id, exclude soft-deleted)
5. Write Update query
6. Write SoftDelete query

Phase 3: Code Generation
1. Run sqlc generate
2. Verify internal/db/models.go has Subarea struct
3. Verify internal/db/queries.sql.go has all subarea queries
4. Run go build ./internal/db

Phase 4: Integration Tests
1. Create internal/db/subareas_test.go
2. Write TestCreateSubarea:
   - Create parent area first
   - Create subarea with FK
   - Fetch by ID
   - Verify relationship
3. Write TestListSubareasByArea:
   - Create area
   - Create multiple subareas for that area
   - Create subarea for different area
   - Call ListByArea
   - Verify only correct subareas returned
4. Write TestUpdateSubarea:
   - Create area + subarea
   - Update subarea
   - Verify updated_at changed
5. Write TestSoftDeleteSubarea:
   - Create area + subarea
   - Soft delete
   - Verify deleted_at set
6. Run tests: go test ./internal/db -v

Phase 5: Verification
1. All tests pass
2. FK constraints work correctly
3. Code compiles
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Subareas entity with TDD following existing patterns.

Changes:
- Created internal/domain/subarea.go with Subarea struct, NewSubarea constructor, and GetEffectiveColor method for color inheritance from parent Area
- Created internal/domain/subarea_test.go with unit tests for color inheritance and validation
- Created queries/subareas.sql with CRUD queries: CreateSubarea, GetSubareaByID, ListSubareasByArea, UpdateSubarea, SoftDeleteSubarea
- Generated sqlc code in internal/db/subareas.sql.go and updated models.go/querier.go
- Created internal/db/subareas_test.go with integration tests for all CRUD operations

Tests:
- go test ./internal/domain -v (color inheritance unit tests)
- go test ./internal/db -v (integration tests for FK constraints, filtering, updates, soft deletes)
- go test ./... (all tests pass)
- go build ./... (build succeeds)
<!-- SECTION:FINAL_SUMMARY:END -->
