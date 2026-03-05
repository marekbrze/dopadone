---
id: TASK-23
title: Fix CLI create commands to support JSON output
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-04 10:12'
updated_date: '2026-03-04 10:21'
labels:
  - backend
  - cli
  - bug
dependencies: []
references:
  - cmd/projectdb/areas.go
  - cmd/projectdb/subareas.go
  - cmd/projectdb/projects.go
  - cmd/projectdb/tasks.go
  - scripts/seed-test-data.sh
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The create commands for areas, subareas, projects, and tasks currently only output success messages and don't support JSON output format. This breaks the seed script which relies on --output json to extract created IDs. The create commands should respect the global --output flag and output created objects in JSON format when requested.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 areas create command supports --output json flag
- [x] #2 subareas create command supports --output json flag
- [x] #3 projects create command supports --output json flag
- [x] #4 tasks create command supports --output json flag
- [x] #5 JSON output includes the full created object with all fields
- [x] #6 Seed script (scripts/seed-test-data.sh) works correctly and creates all test data
- [x] #7 TUI displays seeded data correctly (areas, subareas, projects, tasks)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Modify areas create command to check --output flag and use JSON formatter when json is specified\n2. Modify subareas create command similarly\n3. Modify projects create command similarly\n4. Modify tasks create command similarly\n5. Test each command individually with --output json\n6. Run seed script to verify all data is created\n7. Run TUI to verify data displays correctly
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Implemented JSON output support for all create commands and fixed seed script ID extraction pattern.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed CLI create commands to support JSON output and fixed seed script to correctly extract created object IDs.

## Changes

1. CLI Create Commands - JSON Output Support
   - Modified areas create to respect --output json flag
   - Modified subareas create to respect --output json flag  
   - Modified projects create to respect --output json flag
   - Modified tasks create to respect --output json flag
   - All commands now output full JSON object with all fields when --output json is specified

2. Seed Script Fix (scripts/seed-test-data.sh)
   - Updated ID extraction pattern to handle JSON with spaces after colons
   - Applied to all area, subarea, and project create commands

## Testing

- Verified all create commands output correct JSON format
- Ran seed script successfully: creates 3 areas, 7 subareas, 25 projects
- Created integration tests proving TUI displays all seeded data correctly
- All existing tests pass

## Result

Seed script now works correctly. TUI displays all seeded data (areas, subareas, projects, tasks). CLI create commands support JSON output for programmatic use.
<!-- SECTION:FINAL_SUMMARY:END -->
