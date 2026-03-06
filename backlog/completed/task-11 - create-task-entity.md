---
id: TASK-11
title: create task entity
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 10:53'
updated_date: '2026-03-03 11:14'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Task is a single unit of work within a project. Tasks can be created, listed, updated, and deleted. Each task has title, description, dates, priority, context, estimated duration, and status.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create database migration for tasks table with fields: id, project_id (FK), title, description, start_date, deadline, priority (critical/high/medium/low), context (string), estimated_duration (minutes as int), status (todo/in_progress/waiting/done), created_at, updated_at, deleted_at
- [x] #2 Add value objects for TaskStatus (todo, in_progress, waiting, done), TaskPriority (critical, high, medium, low), TaskDuration (5, 15, 30, 60, 120, 240, 480 minutes) in internal/domain/value_objects.go
- [x] #3 Add domain entity Task in internal/domain/task.go with NewTask factory and validation
- [x] #4 Add SQL queries in queries/tasks.sql for CreateTask, GetTask, ListTasks, UpdateTask, DeleteTask (soft delete)
- [x] #5 Run sqlc generate to create internal/db/tasks.sql.go
- [x] #6 Add CLI commands in cmd/dopa/tasks.go: task create, task list, task get, task update, task delete
- [x] #7 Add is_next boolean field to tasks table (default false) and domain entity for marking priority tasks
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan for Task Entity

### Phase 1: Database Migration (AC#1)
1. Create new migration file: `migrations/YYYYMMDDHHMMSS_add_tasks_table.sql`
2. Define `tasks` table schema with all required fields:
   - id (TEXT PRIMARY KEY)
   - project_id (TEXT, FK to projects)
   - title (TEXT NOT NULL)
   - description (TEXT)
   - start_date (TIMESTAMP)
   - deadline (TIMESTAMP)
   - priority (TEXT with CHECK constraint: critical/high/medium/low)
   - context (TEXT)
   - estimated_duration (INTEGER, minutes)
   - status (TEXT with CHECK constraint: todo/in_progress/waiting/done)
   - is_next (BOOLEAN DEFAULT FALSE) - toggle for marking priority/focused tasks
   - created_at, updated_at, deleted_at (TIMESTAMP)
3. Add appropriate indexes (project_id, status, deadline, is_next)
4. Add down migration to drop tasks table

### Phase 2: Value Objects (AC#2)
1. Add to `internal/domain/value_objects.go`:
   - TaskStatus type with constants (todo, in_progress, waiting, done)
   - TaskPriority type with constants (critical, high, medium, low)
   - TaskDuration type with constants (5, 15, 30, 60, 120, 240, 480 minutes)
   - Validation methods (IsValid()) for each
   - Parse methods (ParseTaskStatus, ParseTaskPriority, ParseTaskDuration)
   - String() methods for each
2. Add error variables for invalid values

### Phase 3: Domain Entity (AC#3)
1. Create `internal/domain/task.go`:
   - Define Task struct with all fields from migration (including IsNext bool)
   - Create NewTask factory function with validation:
     * Title not empty
     * ProjectID not empty
     * Valid status, priority, duration
     * Valid date range (if both start_date and deadline provided)
     * Deadline requires start_date
   - Add helper methods:
     * IsDeleted()
     * IsCompleted()
     * MarkCompleted()
     * SetStatus(), SetPriority() with validation
     * SetNext() / ClearNext() - toggle is_next flag
2. Add comprehensive error variables

### Phase 4: SQL Queries (AC#4)
1. Create `queries/tasks.sql` with sqlc annotations:
   - CreateTask: INSERT with all fields
   - GetTask: SELECT by ID
   - ListTasks: SELECT all by project_id
   - ListNextTasks: SELECT all where is_next = true
   - UpdateTask: UPDATE all mutable fields
   - DeleteTask (soft delete): UPDATE deleted_at
   - ListTasksByStatus: Filter by status
   - ListTasksByPriority: Filter by priority
   - ToggleNext: UPDATE is_next field

### Phase 5: Code Generation (AC#5)
1. Run `sqlc generate` to create:
   - `internal/db/tasks.sql.go` (query implementations)
   - Update `internal/db/models.go` with Task struct
   - Update `internal/db/querier.go` with new interface methods

### Phase 6: CLI Commands (AC#6)
1. Create `cmd/dopa/tasks.go`:
   - Root command: `tasks`
   - Subcommands:
     * `create`: Create new task (--project-id required, --title, --status, --priority, --next, etc.)
     * `list`: List tasks (--project-id, --status, --priority, --next, --format)
     * `next`: List all tasks marked as "next" (shortcut for list --next)
     * `get <id>`: Get single task
     * `update <id>`: Update task fields (including --next/--no-next toggle)
     * `delete <id>`: Soft delete task
2. Implement command handlers following pattern from projects.go:
   - Parse and validate flags
   - Connect to database
   - Call generated queries
   - Format and display output
3. Add helper functions:
   - taskToMap() for output formatting
   - Parsing helpers for task-specific types

### Testing Strategy
- Run existing tests to ensure no regressions
- Test database migration (apply and rollback)
- Test each CLI command manually
- Verify value object validation
- Test domain entity validation
- Test "next" toggle functionality (set, list, clear)

### Dependencies & Risks
- Must complete before dependent tasks (task listing, filtering)
- Migration must be idempotent
- All datetime handling must be consistent with existing code
- Foreign key constraints must be properly handled
- is_next flag should be clearly visible in list output
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Task entity with full CRUD operations.

Changes:
- Created migration 20260303110742_add_tasks_table.sql with tasks table including id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next fields
- Added TaskStatus (todo/in_progress/waiting/done), TaskPriority (critical/high/medium/low), TaskDuration (5/15/30/60/120/240/480 min) value objects to internal/domain/value_objects.go
- Created internal/domain/task.go with Task entity, NewTask factory, and validation
- Added queries/tasks.sql with CreateTask, GetTaskByID, ListTasksByProject, ListNextTasks, ListTasksByStatus, ListTasksByPriority, UpdateTask, SoftDeleteTask, ToggleIsNext queries
- Generated internal/db/tasks.sql.go via sqlc
- Created cmd/dopa/tasks.go with task create/list/next/get/update/delete CLI commands
- Added task parser functions to internal/cli/validation.go
- Updated db_test.go to verify tasks table and indexes

Tests: go test ./... passes all tests
<!-- SECTION:FINAL_SUMMARY:END -->
