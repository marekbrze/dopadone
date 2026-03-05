---
id: TASK-5
title: Subareas CRUD CLI Commands
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 09:37'
updated_date: '2026-03-03 09:53'
labels:
  - cli
  - crud
  - subareas
dependencies:
  - TASK-8
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement full CRUD operations for Subareas entity via CLI with area relationship handling.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 subareas create: --name (required), --area-id (required), --color flags
- [x] #2 subareas list: --area-id filter, --json flag, --format=table|json
- [x] #3 subareas get <id>: display single subarea details
- [x] #4 subareas update <id>: --name, --color flags
- [x] #5 subareas delete <id>: --permanent flag
- [x] #6 Auto-generated help with examples via Cobra
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create cmd/projectdb/subareas.go with Cobra command structure\n2. Implement subareas create command with --name, --area-id, --color flags\n3. Implement subareas list command with --area-id filter (lists all when not specified)\n4. Implement subareas get command for single subarea display\n5. Implement subareas update command with --name, --color flags\n6. Implement subareas delete command with --permanent flag (soft delete by default)\n7. Add examples to all commands for help text\n8. Register subareas parent command in main.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented complete CRUD CLI commands for Subareas entity.

Changes:
- Created cmd/projectdb/subareas.go with 5 subcommands (create, list, get, update, delete)
- Implemented subareas create with --name, --area-id (required) and --color flags
- Implemented subareas list with --area-id filter, --json flag, and --format=table|json
- Implemented subareas get for single subarea display (JSON output)
- Implemented subareas update with --name and --color flags
- Implemented subareas delete with --permanent flag (soft delete by default, permanent delete via raw SQL)
- Added listAllSubareas helper for listing all subareas when no area filter is provided
- Added examples to all commands for auto-generated Cobra help
- Registered subareas command in main.go

All commands include:
- Proper validation using existing CLI validation utilities
- Consistent error handling with exit codes
- Table and JSON output format support
- Colored output using lipgloss
<!-- SECTION:FINAL_SUMMARY:END -->
