---
id: TASK-14
title: Implement TUI - 3-column browser with quick-add
status: Done
assignee: []
created_date: '2026-03-03 12:17'
updated_date: '2026-03-04 09:03'
labels:
  - tui
  - ux
  - mvp
dependencies: []
references:
  - internal/domain/area.go
  - internal/domain/project.go
  - internal/domain/task.go
  - internal/db/querier.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Build an interactive terminal UI (TUI) for browsing and managing the project hierarchy. The TUI provides a visual interface with areas as tabs at the top and a 3-column browser below (Subareas | Projects | Tasks). Users can navigate with keyboard and quickly add new items with the 'a' key. This is an MVP focused on viewing and quick-add; full editing will be a follow-up task.

Tech stack: bubbletea + bubbles for TUI framework. The TUI should follow clean architecture principles with the UI as a delivery mechanism that depends on domain/use cases.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Given the TUI is running, When viewing areas, Then areas are displayed as a horizontal tab bar at the top of the screen
- [ ] #2 Given an area is selected, When viewing the 3-column layout, Then the left column shows all subareas belonging to the selected area
- [ ] #3 Given a subarea is selected, When viewing the center column, Then all projects for that subarea are displayed in a tree structure (supporting unlimited nesting via parent_id)
- [ ] #4 Given a project is selected, When viewing the right column, Then all tasks for that project are displayed in a list
- [ ] #5 Given focus is on any column, When pressing 'h' or left arrow, Then focus moves to the column on the left (wrapping from subareas to tasks)
- [ ] #6 Given focus is on any column, When pressing 'l' or right arrow, Then focus moves to the column on the right (wrapping from tasks to subareas)
- [ ] #7 Given focus is on any column, When pressing Tab, Then focus cycles through columns in order (subareas → projects → tasks → subareas)
- [ ] #8 Given focus is on a column with items, When pressing up/down arrows or j/k, Then the selected item changes within that column
- [ ] #9 Given focus is on any column, When pressing 'a' key, Then a quick-add modal appears centered on screen with a single title input field
- [ ] #10 Given the quick-add modal is open, When user types a title and presses Enter, Then a new item is created in the focused column's context (subarea/project/task) and the modal closes
- [ ] #11 Given the quick-add modal is open, When user presses Escape, Then the modal closes without creating an item
- [ ] #12 Given the user switches areas via tab bar, When returning to a previously selected area, Then the last selected subarea, project, and task are restored for that area
- [ ] #13 Given the TUI is started, When there are existing areas in the database, Then the first area is automatically selected and its subareas are loaded
- [ ] #14 Given the TUI displays a project tree in the center column, When a project has child projects, Then they are displayed indented under their parent with visual tree indicators (├─, └─, etc.)
- [ ] #15 Given focus is on a column, When there are no items to display, Then an empty state message is shown (e.g., 'No subareas - press a to add')
- [ ] #16 Given the TUI is running, When pressing 'q' or Ctrl+C, Then the application exits cleanly
<!-- AC:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Split into subtasks: TASK-15 (14A Core), TASK-16 (14B Data/Tree), TASK-18 (14C Navigation), TASK-19 (14D Modal), TASK-17 (14E Polish). Original task kept for reference.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## TUI Implementation Complete

Successfully implemented a 3-column browser TUI with quick-add functionality, split into 6 subtasks (15-21, with 16 superseded by 20+21).

### Subtasks Completed:
- **Task-15** (14A Core): Framework, tabs, 3-column layout, column navigation, exit handling
- **Task-20** (16A Tree): Tree rendering package with unlimited nesting, lipgloss styling
- **Task-21** (16B Data): Data loading, spinner, empty states, cascade loading, auto-select
- **Task-18** (14C Nav): In-column j/k navigation, area switching ([/]), state persistence
- **Task-19** (14D Modal): Quick-add modal with context-aware creation, validation
- **Task-17** (14E Polish): Help modal (?), toast notifications, footer, docs
- **Task-16**: Superseded by splitting into Tasks 20+21

### Key Features:
- Area tabs with keyboard navigation ([/])
- 3-column browser: Subareas | Projects | Tasks
- Hierarchical project tree with expand/collapse
- State persistence per area (selections, tree expansion)
- Quick-add modal for creating items (a key)
- Help modal with all shortcuts (? key)
- Toast notifications for errors
- Loading spinner during data fetch

### Test Coverage:
- internal/tui: 82.4% (40+ tests)
- internal/tui/tree: 95.0% (51 tests)
- All 127 TUI tests passing

### Architecture:
- Clean architecture: TUI → Domain (no reverse dependencies)
- Repository injection pattern (DIP)
- Model-Update-View pattern (bubbletea)
- No magic numbers (all constants named)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All acceptance criteria are implemented and manually tested
- [ ] #2 Unit tests exist for key components (navigation, state management, tree rendering)
- [ ] #3 TUI follows clean architecture: UI components depend on domain/use cases, not vice versa
- [ ] #4 Code uses bubbletea best practices (Model-Update-View pattern)
- [ ] #5 No external dependencies beyond bubbletea, bubbles, and existing project dependencies
- [ ] #6 Keyboard shortcuts are documented in help screen (accessible via '?' key)
- [ ] #7 Error states are handled gracefully (database errors, empty states)
- [ ] #8 Code is linted and passes go vet
<!-- DOD:END -->
