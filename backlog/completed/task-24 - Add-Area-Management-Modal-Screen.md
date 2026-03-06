---
id: TASK-24
title: Add Area Management Modal Screen
status: Done
assignee:
  - '@agent'
created_date: '2026-03-04 12:16'
updated_date: '2026-03-04 14:58'
labels:
  - tui
  - area-management
  - enhancement
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a modal screen for managing areas including: reordering (u/d keys), creating new areas, editing name/color, and deleting areas with cascade delete of all subareas/projects/tasks. Deletion should offer user choice between soft and permanent delete. Allow deletion of the last area resulting in empty state.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 List all areas with visual indicators for selection and order
- [x] #2 Add new area functionality with name and color input
- [x] #3 Edit existing area (name and color fields)
- [x] #4 Delete area with confirmation dialog showing cascade warning
- [x] #5 Delete modal offers choice: soft delete or permanent delete
- [x] #6 Permanent delete cascades to all subareas, projects, and tasks
- [x] #7 Soft delete marks only the area record (children become orphaned)
- [x] #8 Allow deletion of last area resulting in empty state UI
- [x] #9 Toast notifications for success/error feedback
- [x] #10 Database schema updated with sort_order column on areas table
- [x] #11 Unit tests for area service methods (CRUD + reorder)
- [x] #12 Integration tests for cascade delete behavior
- [x] #13 Modal screen accessible via Ctrl+A keybinding
- [x] #14 u/d keys reorder areas; changes batch-saved only on modal close (OK), discarded on Esc
- [x] #15 Color picker with 12 predefined colors (dropdown)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add sort_order column to areas table (migration)\n2. Update queries/areas.sql with new queries (sort_order, hard delete, stats counts)\n3. Run sqlc generate to update db package\n4. Update domain/area.go with SortOrder field\n5. Create internal/service/area_service.go with CRUD + reorder + stats methods\n6. Create internal/tui/areamodal/area_modal.go (list, create, edit, delete-confirm modes)\n7. Add messages to internal/tui/messages.go\n8. Add commands to internal/tui/commands.go\n9. Integrate modal into internal/tui/app.go with Ctrl+A keybinding\n10. Write unit tests for area_service\n11. Test the full flow manually
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
SPECIFICATION FINALIZED:

Keybinding: Ctrl+A to open area management modal

Predefined Colors (12 options):
- Blue (#3B82F6), Green (#10B981), Yellow (#F59E0B), Red (#EF4444)
- Purple (#8B5CF6), Pink (#EC4899), Orange (#F97316), Teal (#14B8A6)
- Indigo (#6366F1), Cyan (#06B6D4), Gray (#6B7280), Brown (#92400E)

Reorder: u/d keys, batch-saved ONLY when modal closes (OK/Enter)
- Cancel (Esc) discards any pending reorder changes

Delete behavior:
- Soft delete: marks area deleted_at only (children orphaned)
- Permanent delete: CASCADE removes all subareas/projects/tasks
- Show stats: "3 subareas, 12 projects, 45 tasks will be affected"

Modal states: List → Create/Edit/DeleteConfirm → back to List

Empty state: When no areas, show "Press a to create your first area"

Step 1-3 complete: Added migration, updated SQL queries, ran sqlc generate, updated domain/area.go with SortOrder, updated converters.go

Implementation progress:
- Area modal UI complete: list, create, edit, delete-confirm, reorder modes
- Ctrl+A keybinding opens modal
- u/d keys for reorder (batch-saved on Enter, discarded on Esc)
- Stats loading on delete confirm
- Toast notifications for all operations
- 12 predefined colors with Tab navigation
- Empty state UI for no areas

Remaining:
- Unit tests for area_service (AC#11)
- Integration tests for cascade delete (AC#12)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented area management modal screen with full CRUD operations, reorder functionality,  and u/d keys batch reorder), delete confirm dialog with soft and permanent delete options ( and toast notifications for success/error feedback.

Features:
- Ctrl+A opens modal
- List mode: shows areas with color preview and visual indicators
- Create mode: form with name input and Tab for color navigation
- Edit mode: populates name and color fields,- Delete mode: shows confirmation with stats and soft vs permanent delete options
- Reorder mode: u/d keys to reorder (changes saved on Enter, discarded on Esc)
- Empty state UI when no areas
- Database: sort_order column added to areas table
- Commands for CRUD + reorder + stats
- Messages for commands for app.go
- Area service with unit tests
- Integration with app.go and- Modal UI in areamodal/area_modal.go

Stats are loaded via LoadAreaStatsCmd when entering delete confirm mode.
<!-- SECTION:FINAL_SUMMARY:END -->
