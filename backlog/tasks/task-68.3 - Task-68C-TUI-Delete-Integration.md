---
id: TASK-68.3
title: 'Task-68C: TUI Delete Integration'
status: Done
assignee:
  - '@assistant'
created_date: '2026-03-09 19:04'
updated_date: '2026-03-10 08:00'
labels: []
dependencies:
  - TASK-68.1
  - TASK-68.2
references:
  - '# Parent: TASK-68 - Add option to delete subareas'
  - projects and tasks in tui
  - '# Dependency: TASK-68.1 - Confirmation Modal Component'
  - '# Dependency: TASK-68.2 - Cascade Soft Delete Service'
  - '# Dependent: TASK-68.4 - Delete Documentation & Testing'
  - '# Milestone: m-2 Deleting items'
parent_task_id: TASK-68
priority: high
---

## Description
<!-- SECTION:DESCRIPTION:BEGIN -->
Wire up delete functionality in TUI with confirmation modal. Add 'd' key binding to all columns. Handle confirmation, execute delete, show toast notifications, refresh columns.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Add confirmModal state to Model struct
- [x] #2 Add delete key binding 'd' to all focus columns
- [x] #3 Open confirmation modal with correct item name when 'd' pressed
- [x] #4 Execute appropriate delete command on 'y' confirmation
- [x] #5 Show success toast and refresh column after delete
- [x] #6 Show error toast on delete failure
- [x] #7 Handle 'n' and Escape to cancel delete
- [x] #8 No-op when pressing 'd' on empty columns
- [x] #9 Update footer to show 'd: delete' shortcut
- [x] #10 Write integration tests for delete flow
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Wire up delete functionality in TUI with confirmation modal. Add 'd' key binding to all focus columns. Open confirmation modal with correct item name when 'd' pressed. Handle confirmation, execute appropriate delete command on 'y' confirmation, Show success toast and refresh column after delete. Handle 'n' and Escape to cancel delete. No-op when pressing 'd' on empty columns. Update footer to show 'd: delete' shortcut. Write integration tests for delete flow
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
\n## Part of\n\nThis is subtask 3 of 4 for TASK-68. Integration layer - requires Task-68.1 and Task-68.2 to be complete first.

Scoping complete

Confirmed approach:
- Keep all 10 ACs separate
- Follow existing modal pattern (confirmModal + isConfirmModalOpen)
- Use existing services (SoftDelete, SoftDeleteWithCascade)
- Show item name in toast
- d before x: toggle in footer

Test strategy:
- Separate test file (app_delete_test.go)
- 100% coverage for new code

Estimated effort: 4-6 hours
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:SUMMARY:BEGIN -->
Implemented TUI delete integration with confirmation modal for subareas, projects, and tasks.

**Changes:**
- Added `confirmModal` and `isConfirmModalOpen` fields to Model struct
- Added DeleteSuccessMsg/DeleteErrorMsg message types
- Created DeleteSubareaCmd, DeleteProjectCmd (cascade), DeleteTaskCmd commands
- Added 'd' key handling in app.go for all focus columns
- Implemented delete_handlers.go with handlers for confirm/cancel/success/error
- Updated footer to include 'd: delete' shortcut
- Added SoftDeleteWithCascade to MockProjectService interface and implementation
- Created mock helpers for subarea/project/task delete operations
- Added comprehensive unit tests for all delete commands

**Key Implementation Details:**
- Projects use SoftDeleteWithCascade for recursive deletion of child projects/tasks
- Confirmation modal displays item name for user clarity
- Success/error toast notifications include entity name and type
- Column refresh logic after successful delete (subarea → load subareas, project → load projects, task → load tasks)
- Empty column check prevents modal from opening on no selection
- 'n' and Escape keys cancel delete operation
- Footer shows 'd: delete' before 'x: toggle'

**Tests:**
- Unit tests for DeleteSubareaCmd, DeleteProjectCmd, DeleteTaskCmd
- Test coverage for success and error scenarios
- All tests passing
<!-- SECTION:SUMMARY:END -->
