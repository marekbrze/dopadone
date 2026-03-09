---
id: TASK-68.3
title: 'Task-68C: TUI Delete Integration'
status: To Do
assignee: []
created_date: '2026-03-09 19:04'
updated_date: '2026-03-09 19:06'
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
- [ ] #1 Add confirmModal state to Model struct
- [ ] #2 Add delete key binding 'd' to all focus columns
- [ ] #3 Open confirmation modal with correct item name when 'd' pressed
- [ ] #4 Execute appropriate delete command on 'y' confirmation
- [ ] #5 Show success toast and refresh column after delete
- [ ] #6 Show error toast on delete failure
- [ ] #7 Handle 'n' and Escape to cancel delete
- [ ] #8 No-op when pressing 'd' on empty columns
- [ ] #9 Update footer to show 'd: delete' shortcut
- [ ] #10 Write integration tests for delete flow
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
\n## Part of\n\nThis is subtask 3 of 4 for TASK-68. Integration layer - requires Task-68.1 and Task-68.2 to be complete first.
<!-- SECTION:NOTES:END -->
