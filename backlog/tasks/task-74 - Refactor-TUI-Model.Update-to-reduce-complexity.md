---
id: TASK-74
title: Refactor TUI Model.Update to reduce complexity
status: Done
assignee:
  - '@ai-agent'
created_date: '2026-03-10 12:37'
updated_date: '2026-03-10 18:27'
labels:
  - lint
  - code-quality
  - refactor
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The (Model).Update function in internal/tui/app.go has cyclomatic complexity of 94 (limit is 30). Need to refactor by extracting handlers for different message types.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Analyze current Update method structure and message types
- [x] #2 Extract handler functions for distinct message types (e.g., keyMsgHandler, mouseMsgHandler, etc.)
- [x] #3 Reduce complexity to below 30
- [x] #4 Ensure all tests pass after refactoring
- [x] #5 Run make lint to verify gocyclo error is resolved
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze current Update method - map all 25+ message types to handler categories
2. Create data_loader_handlers.go - Extract: AreasLoadedMsg, SubareasLoadedMsg, ProjectsLoadedMsg, TasksLoadedMsg handlers
3. Create keyboard_handler.go - Extract entire tea.KeyMsg switch block with nested modal/key handling
4. Create spacemenu_handler.go - Extract spacemenu.CloseMsg and spacemenu.ActionMsg handlers
5. Create window_handler.go - Extract tea.WindowSizeMsg handler
6. Create spinner_handler.go - Extract spinner.TickMsg handler
7. Create toast_handler.go - Extract ToastTickMsg handler
8. Create task_handlers.go - Extract TaskStatusToggledMsg handler
9. Refactor Update method to dispatcher pattern - Route to extracted handlers
10. Run gocyclo to verify complexity < 30
11. Run all tests to ensure no regressions
12. Update TUI.md documentation if needed
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created 6 new handler files to extract logic from Update method
- Refactored Update method to simple dispatcher pattern
- Reduced complexity from 94 to 5
- All tests passing
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Successfully refactored the TUI Model.Update method to reduce cyclomatic complexity from 94 to below 30.

Changes:
- Created 6 new handler files to extract message handling logic
- data_loader_handlers.go - Handles AreasLoadedMsg, SubareasLoadedMsg, ProjectsLoadedMsg, TasksLoadedMsg
- keyboard_handler.go - Routes keyboard input through modal/key handlers
- spinner_handler.go - Handles spinner.TickMsg
- toast_handler.go - Handles ToastTickMsg
- task_handlers.go - Handles TaskStatusToggledMsg
- window_handler.go - Handles tea.WindowSizeMsg

- Refactored Update method to use dispatcher pattern routing messages to appropriate handlers
- All existing tests pass without regression
- Reduced cyclomatic complexity from 94 to 5 (verified with gocyclo)
- No lint errors
- Code is ready for review (no commit/push as requested)

Verification:
- Cyclomatic complexity: 94 → 5
- All tests: PASS
- Lint check: PASS
- Compilation: Successful
<!-- SECTION:FINAL_SUMMARY:END -->
