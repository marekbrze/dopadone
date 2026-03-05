---
id: TASK-7
title: Projects CRUD CLI Commands
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 09:38'
updated_date: '2026-03-03 10:12'
labels:
  - cli
  - crud
  - projects
dependencies:
  - TASK-8
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement full CRUD operations for Projects entity via CLI - the most complex entity with many fields and recursive nesting.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 projects create: --name (required), --subarea-id OR --parent-id (one required), --status, --priority, --progress, --deadline, --start-date, --color, --goal, --description
- [x] #2 projects list: --status, --priority, --subarea-id, --parent-id filters, --json flag
- [x] #3 projects get <id>: display single project with all fields
- [x] #4 projects update <id>: all editable fields as optional flags
- [x] #5 projects delete <id>: --permanent flag
- [x] #6 Validation: subarea_id XOR parent_id constraint, valid status/priority/progress values
- [x] #7 Auto-generated help with examples via Cobra
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create cmd/projectdb/projects.go file with cobra command structure (following subareas.go pattern)
2. Implement 'projects create' command with all required/optional flags:
   - Required: --name
   - One required (XOR): --subarea-id OR --parent-id
   - Optional: --status, --priority, --progress, --deadline, --start-date, --color, --goal, --description
   - Add XOR validation for subarea-id vs parent-id
   - Use existing validation helpers from internal/cli/validation.go
3. Implement 'projects list' command with filter flags:
   - --status, --priority, --subarea-id, --parent-id filters
   - --json flag for JSON output
   - Support listing all projects when no filters provided
4. Implement 'projects get <id>' command to display single project with all fields
5. Implement 'projects update <id>' command with all editable fields as optional flags
6. Implement 'projects delete <id>' command with --permanent flag
7. Register all subcommands under projects parent command in main.go
8. Add comprehensive examples to help text for each command
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented full CRUD operations for Projects entity via CLI following the same pattern as subareas.go. Created cmd/projectdb/projects.go with all required subcommands. Fixed critical bug where create and update commands were sharing the same flag variables, causing validation failures. Separated variables into projCreate* and projUpdate* prefixes to prevent cross-contamination. All validation working correctly including XOR constraint for subarea-id vs parent-id, status/priority/progress values, and date range validation.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented complete CRUD CLI commands for Projects entity with all required functionality:

**Changes:**
- Created cmd/projectdb/projects.go with 5 subcommands: create, list, get, update, delete
- Implemented comprehensive flag support for all project fields
- Added XOR validation for subarea-id vs parent-id constraint
- Integrated all existing validation helpers (status, priority, progress, color, date range)
- Added filter support in list command (status, priority, subarea-id, parent-id)
- Implemented soft delete with --permanent flag for hard delete
- Added comprehensive help examples for all commands

**Key Implementation Details:**
- Separated create and update flag variables to prevent cross-contamination (projCreate* vs projUpdate*)
- Reused validation functions from internal/cli/validation.go
- Followed existing codebase patterns from subareas.go
- All fields validated before database operations
- JSON output supported via --json flag

**Testing:**
- Manual testing of all CRUD operations
- Verified validation for all field types
- Tested nested project creation (parent-id vs subarea-id)
- Confirmed soft delete and permanent delete functionality
- All existing tests passing
<!-- SECTION:FINAL_SUMMARY:END -->
