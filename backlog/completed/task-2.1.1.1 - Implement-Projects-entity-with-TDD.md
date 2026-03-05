---
id: TASK-2.1.1.1
title: Implement Projects entity with TDD
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 08:26'
updated_date: '2026-03-03 09:13'
labels:
  - backend
  - database
  - tdd
dependencies:
  - TASK-2.1.1
parent_task_id: TASK-2.1.1
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement Projects entity CRUD operations using TDD. Projects are goal-oriented task groups within Subareas, with recursive nesting support (parent_id self-reference).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Write test queries in queries/projects.sql: Create, GetByID, ListBySubarea, ListByParent, Update, SoftDelete, GetByStatus
- [x] #2 Run sqlc generate and verify internal/db/models.go contains Project struct
- [x] #3 Run sqlc generate and verify internal/db/queries.sql.go contains Project queries
- [x] #4 Write integration test: Create root project (with subarea_id), verify all fields
- [x] #5 Write integration test: Create nested project (with parent_id), verify hierarchy
- [x] #6 Write integration test: List projects by subarea_id, verify filtering
- [x] #7 Write integration test: List nested projects by parent_id, verify correct children
- [x] #8 Write integration test: Update project status to 'completed', verify completed_at is set
- [x] #9 Write integration test: Get projects by status, verify filtering
- [x] #10 Write integration test: Soft delete project, verify deleted_at is set
- [x] #11 Write unit test: CHECK constraint violation when both parent_id and subarea_id are NULL
- [x] #12 Write unit test: progress field validation (0-100 range)
- [x] #13 Write unit test: DateRange validation (start_date < deadline)
- [x] #14 Verify all tests pass with go test ./...
- [x] #15 Document Project entity in internal/domain/project.go with all value objects
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Phase 1: Add value objects (ProjectStatus, Priority, DateRange) with validation\nPhase 2: Write SQL queries in queries/projects.sql\nPhase 3: Run sqlc generate\nPhase 4: Write integration tests in internal/db/projects_test.go\nPhase 5: Document Project entity in internal/domain/project.go\nPhase 6: Verify all tests pass
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Projects entity with full TDD approach:

**Value Objects (internal/domain/value_objects.go)**:
- ProjectStatus enum: active, completed, on_hold, archived
- Priority enum: low, medium, high, urgent  
- Progress type with 0-100 validation
- DateRange value object with start_date < deadline validation

**SQL Queries (queries/projects.sql)**:
- CreateProject, GetProjectByID, UpdateProject, SoftDeleteProject
- ListProjectsBySubarea, ListProjectsByParent, GetProjectsByStatus

**Generated Code (sqlc)**:
- internal/db/models.go: Project struct with all fields
- internal/db/projects.sql.go: Query implementations
- internal/db/querier.go: Interface updated

**Integration Tests (internal/db/projects_test.go)**:
- TestCreateRootProject, TestCreateNestedProject
- TestListProjectsBySubarea, TestListProjectsByParent
- TestUpdateProjectStatus, TestGetProjectsByStatus
- TestSoftDeleteProject, TestConstraintViolationBothParentAndSubareaNull

**Domain Entity (internal/domain/project.go)**:
- Project struct with business logic methods
- NewProject constructor with full validation
- IsDeleted, IsCompleted, IsNested helper methods
- MarkCompleted, SetProgress, SetPriority, SetStatus mutators

**All tests pass**: go test ./...
<!-- SECTION:FINAL_SUMMARY:END -->
