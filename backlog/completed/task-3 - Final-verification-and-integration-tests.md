---
id: TASK-3
title: Final verification and integration tests
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 08:26'
updated_date: '2026-03-03 09:19'
labels:
  - backend
  - database
  - integration
dependencies:
  - TASK-2.1.1.1
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Final verification of the complete database layer. Run all migrations, verify schema integrity, and test end-to-end hierarchical data operations.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Run goose up on fresh SQLite database and verify all tables exist
- [x] #2 Verify all foreign key constraints are properly created
- [x] #3 Verify all indexes are created (use SQLite PRAGMA index_list)
- [x] #4 Write integration test: Complete hierarchy: Area → Subarea → Project → Sub-project
- [x] #5 Write integration test: Cascade behavior when soft-deleting parent (verify orphaned records handled correctly)
- [x] #6 Write integration test: Query performance on indexed fields (deadline, status, priority)
- [x] #7 Run goose down and verify clean rollback
- [x] #8 Run goose up again and verify idempotency
- [x] #9 Verify sqlc-generated code compiles: go build ./...
- [x] #10 Run full test suite: go test ./... -v
- [x] #11 Document any assumptions or constraints in README or schema docs
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
TDD IMPLEMENTATION PLAN:

Phase 1: Migration Verification
1. Create fresh SQLite database
2. Run goose up
3. Query sqlite_master to verify all tables exist:
   - areas, subareas, projects
4. Run PRAGMA commands:
   - PRAGMA foreign_key_list for each table
   - PRAGMA index_list for each table
5. Verify all expected indexes exist
6. Write test: TestMigrationIdempotency
   - Run goose down
   - Run goose up again
   - Verify no errors

Phase 2: Schema Integrity Tests
1. Write TestForeignKeyConstraints:
   - Try insert subarea with invalid area_id
   - Expect FK violation
   - Try insert project with invalid subarea_id
   - Expect FK violation
2. Write TestCheckConstraints:
   - Try insert project with progress > 100
   - Expect CHECK violation
   - Try insert project with both parent_id and subarea_id NULL
   - Expect CHECK violation
3. Write TestEnumConstraints:
   - Try insert project with invalid status
   - Expect CHECK violation

Phase 3: End-to-End Hierarchy Test
1. Write TestCompleteHierarchy:
   - Create Area: "Home"
   - Create Subarea: "Travel"
   - Create Project: "Trip to Japan"
   - Create Sub-project: "Hotel booking"
   - Create Sub-project: "Flight tickets"
   - Verify all FK relationships
   - Verify queries return correct results
2. Write TestSoftDeleteCascade:
   - Create full hierarchy
   - Soft delete area
   - Verify subarea still exists (soft delete, not hard delete)
   - Document orphaned record behavior

Phase 4: Performance Tests
1. Write TestIndexPerformance:
   - Insert 1000 projects with various statuses/deadlines
   - Query by status
   - Query by deadline range
   - Verify queries use indexes (EXPLAIN QUERY PLAN)

Phase 5: Rollback Test
1. Write TestMigrationRollback:
   - Run goose down to zero
   - Verify database is empty
   - Run goose up
   - Verify schema recreated

Phase 6: Final Build & Test
1. Run go build ./... (all packages)
2. Run go test ./... -v -cover
3. Verify coverage report
4. Document test results
5. Update README with setup instructions
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Final verification and integration tests completed.

**New Integration Tests Added** (internal/db/integration_test.go):
- TestCompleteHierarchy: Tests Area → Subarea → Project → Sub-project hierarchy with FK verification
- TestSoftDeleteCascadeBehavior: Documents that soft-deletes do NOT cascade (by design)
- TestIndexPerformance: Verifies indexes are used for status, deadline, and status+priority queries
- TestMigrationIdempotency: Confirms goose up is idempotent

**Existing Tests Verified**:
- TestMigrationUp: Tables, FKs, indexes, and CHECK constraints
- TestMigrationDown: Clean rollback

**Documentation Updated**:
- README.md with schema overview, constraints, soft delete behavior, and setup instructions

**All Acceptance Criteria Met**:
1. goose up on fresh SQLite - verified
2. FK constraints created - verified
3. Indexes created (PRAGMA index_list) - verified
4. Complete hierarchy test - TestCompleteHierarchy
5. Soft-delete cascade test - TestSoftDeleteCascadeBehavior
6. Query performance test - TestIndexPerformance
7. Rollback test - TestMigrationDown
8. Idempotency test - TestMigrationIdempotency
9. Build passes - go build ./...
10. Full test suite passes - go test ./... -v
11. Documentation updated - README.md
<!-- SECTION:FINAL_SUMMARY:END -->
