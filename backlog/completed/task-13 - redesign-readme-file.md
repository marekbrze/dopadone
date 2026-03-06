---
id: TASK-13
title: redesign readme file
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 11:56'
updated_date: '2026-03-03 12:11'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
readme file should mainly focus on using the cli app from the user perspective. development docs should be at the bottom
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Add user-focused Quick Start section showing full CRUD workflow for one entity (areas, subareas, or projects)
- [x] #2 Add Installation section after Quick Start covering all methods equally (binary downloads, go install, build from source)
- [x] #3 Add comprehensive Usage section with inline command examples for all entities (areas, subareas, projects, tasks)
- [x] #4 Relocate build/compile instructions to bottom Development section
- [x] #5 Relocate testing/CI information to bottom Development section
- [x] #6 Replace schema-heavy introduction with user-focused problem/solution description (what problem does Dopadone solve?)
- [x] #7 Add output format examples showing table vs JSON output
- [x] #8 Document global flags (--db, --output) in a dedicated section
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze current README structure and identify user-focused vs developer-focused content
2. Draft new problem/solution introduction (AC #6)
3. Create Quick Start section with complete CRUD workflow for areas entity (AC #1)
4. Create Installation section covering binary downloads, go install, and build from source (AC #2)
5. Build comprehensive Usage section with inline examples for all entities (AC #3):
   - Areas: create, list, get, update, delete
   - Subareas: create, list, get, update, delete  
   - Projects: create (root and nested), list, get, update, delete
   - Tasks: create, list, get, update, delete, next
6. Add Output Format section with table vs JSON examples (AC #7)
7. Add Global Flags section documenting --db and --output (AC #8)
8. Relocate development content to bottom Development section (AC #4, #5):
   - Build/compile instructions
   - Testing/CI information
   - Schema details
   - Technical constraints
9. Review and refine all examples for accuracy
10. Verify markdown formatting and structure
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created new user-focused README with problem/solution intro
- Added Quick Start section with complete area CRUD workflow
- Added Installation section covering all three methods
- Added comprehensive Usage section for all entities
- Added Output Formats section with table/JSON examples
- Added Global Flags section documenting --db and --output
- Relocated all development content to bottom Development section
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Redesigned README.md to be user-focused with the following structure:

**User-Facing Sections (top):**
- Problem/solution introduction explaining what Dopadone solves
- Quick Start with complete CRUD workflow for areas entity
- Installation covering binary downloads, go install, and build from source equally
- Comprehensive Usage section with inline examples for areas, subareas, projects, and tasks
- Output Formats section showing table vs JSON output examples
- Global Flags section documenting --db, --output, and --format
- Database Migrations section for end-user migration commands

**Developer Sections (bottom):**
- Prerequisites, build commands, development workflow
- Database development, testing instructions
- Schema reference, constraints, and tech stack details
- Index documentation

All command examples verified against actual CLI output. Development docs properly isolated at the bottom for users who want to contribute.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Review README for clarity and completeness
- [x] #2 Verify all command examples are accurate and tested
- [x] #3 Ensure proper markdown formatting and links
- [x] #4 Check that development docs are properly isolated at bottom
<!-- DOD:END -->
