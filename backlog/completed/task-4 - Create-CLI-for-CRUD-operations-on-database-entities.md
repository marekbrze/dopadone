---
id: TASK-4
title: Create CLI for CRUD operations on database entities
status: Done
assignee: []
created_date: '2026-03-03 09:28'
updated_date: '2026-03-03 10:41'
labels:
  - cli
  - crud
  - cobra
dependencies: []
references:
  - internal/domain/area.go
  - internal/domain/subarea.go
  - internal/domain/project.go
  - internal/domain/value_objects.go
  - internal/db/querier.go
  - internal/db/db.go
documentation:
  - backlog/docs/doc-1 - Data-Layer-Architecture.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
SUPERSEDED - This task has been split into more atomic tasks: TASK-8 (CLI Foundation), TASK-9 (Areas CRUD), TASK-5 (Subareas CRUD), TASK-7 (Projects CRUD), TASK-6 (Advanced Features). Use those tasks instead.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 CLI structure: cmd/dopa/ with cobra root command, subcommands for areas/subareas/projects
- [ ] #2 areas command: create (name, color flags), list (--json, --format=table), get <id>, update <id> (--name, --color), delete <id> (--permanent)
- [ ] #3 subareas command: create (name, area-id, color flags), list (--area-id, --json), get <id>, update <id> (--name, --color), delete <id> (--permanent)
- [ ] #4 projects command: create (name, subarea-id OR parent-id required, status, priority, progress, deadline, start-date, color, goal, description), list (--status, --priority, --subarea-id, --parent-id, --json), get <id>, update <id> (all editable fields), delete <id> (--permanent)
- [ ] #5 Advanced query: list commands support --filter flag with query syntax (e.g., --filter 'status=active AND priority>=high')
- [ ] #6 Output formatting: Default table view with colored headers, --json for machine-readable output, --format=json|table|yaml
- [ ] #7 Validation: Use existing domain types (ProjectStatus, Priority, Progress, Color) with proper error messages
- [ ] #8 Error handling: User-friendly error messages, exit codes (0=success, 1=error, 2=validation error), no stack traces in output
- [ ] #9 Database path: --db flag to specify database file (default: ./dopa.db), error if not found
- [ ] #10 Help documentation: Cobra auto-generated help with examples for each command
- [ ] #11 Unit tests: Test command flag parsing, validation, and output formatting
<!-- AC:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Task superseded by atomic tasks TASK-5 through TASK-9. All CRUD operations for Areas, Subareas, and Projects have been successfully implemented and documented. Documentation available in doc-2 - CLI CRUD Operations Guide.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All acceptance criteria verified and checked
- [ ] #2 go build ./... succeeds
- [ ] #3 go test ./... passes
- [ ] #4 go vet ./... passes
- [ ] #5 CLI binary builds and runs
- [ ] #6 Help text is clear and includes examples
<!-- DOD:END -->
