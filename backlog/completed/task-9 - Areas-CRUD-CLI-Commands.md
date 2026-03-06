---
id: TASK-9
title: Areas CRUD CLI Commands
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 09:39'
updated_date: '2026-03-03 10:13'
labels:
  - cli
  - crud
  - areas
dependencies:
  - TASK-8
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement full CRUD operations for Areas entity via CLI using existing domain validation and database layer.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 areas create: --name (required), --color flags with validation
- [x] #2 areas list: default table output, --json flag, --format=table|json
- [x] #3 areas get <id>: display single area details
- [x] #4 areas update <id>: --name, --color flags (at least one required)
- [x] #5 areas delete <id>: --permanent flag for hard delete vs soft delete
- [x] #6 Auto-generated help with examples via Cobra
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create cmd/dopa/areas.go with cobra command struct
2. Implement areas create subcommand with flags
3. Implement areas list subcommand with output formatting
4. Implement areas get subcommand
5. Implement areas update subcommand
6. Implement areas delete subcommand with --permanent flag
7. Register all subcommands under areas parent
8. Add examples to help text
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented full Areas CRUD CLI commands in cmd/dopa/areas.go following existing patterns from projects.go and subareas.go.

Commands implemented:
- `areas create`: Creates new area with --name (required) and optional --color flags
- `areas list`: Lists all areas with table output (default) or JSON via --json/--format=json
- `areas get <id>`: Displays single area details as JSON
- `areas update <id>`: Updates area with --name and/or --color (at least one required)
- `areas delete <id>`: Soft delete by default, --permanent flag for hard delete

All commands include auto-generated help with examples via Cobra. Registered areasCmd in main.go alongside existing subareas and projects commands.
<!-- SECTION:FINAL_SUMMARY:END -->
