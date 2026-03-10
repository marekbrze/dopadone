---
id: TASK-68.4
title: 'Task-68D: Delete Documentation & Testing'
status: To Do
assignee: []
created_date: '2026-03-09 19:05'
updated_date: '2026-03-10 07:31'
labels: []
milestone: m-2
dependencies:
  - TASK-68.3
references:
  - '# Parent: TASK-68 - Add option to delete subareas'
  - projects and tasks in tui
  - '# Dependency: TASK-68.3 - TUI Delete Integration'
  - '# Milestone: m-2 Deleting items'
parent_task_id: TASK-68
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update documentation and comprehensive testing for delete functionality. Update TUI.md with keyboard shortcuts, add mock delete methods, ensure 80%+ test coverage.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Update docs/TUI.md with 'd: delete' keyboard shortcut
- [ ] #2 Document delete confirmation modal behavior
- [ ] #3 Document cascade delete for projects in TUI.md
- [ ] #4 Add mock delete methods to mocks/services.go
- [ ] #5 Add delete scenarios to integration tests
- [ ] #6 Verify 80%+ test coverage for all delete-related code
- [ ] #7 Manual testing checklist: all entity types, empty columns, errors
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
\n## Part of\n\nThis is subtask 4 of 4 for TASK-68. Documentation and testing - requires Task-68.3 to be complete. Documentation work can start early in parallel with Task-68.3.
<!-- SECTION:NOTES:END -->
