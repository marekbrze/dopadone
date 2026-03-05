---
id: TASK-17
title: 'TUI 14E: Help, Errors & Polish'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 12:31'
updated_date: '2026-03-04 08:42'
labels:
  - tui
  - mvp
  - phase4
dependencies:
  - TASK-15
  - TASK-18
  - TASK-19
  - TASK-20
  - TASK-21
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Final polish including help screen, error handling, documentation updates, and full integration testing.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 ? key opens help modal with all keyboard shortcuts
- [x] #2 Help groups shortcuts by category (navigation, actions)
- [x] #3 Database errors displayed gracefully as toast/notification
- [x] #4 Quick reference shortcuts shown in footer
- [x] #5 README updated with tui command usage
- [x] #6 Keyboard shortcuts documented
- [x] #7 All code passes go vet and lint
- [x] #8 Integration tests for key user flows
- [x] #9 All TUI ACs from tasks 15, 18-21 verified manually
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Phase 1: Help Modal Component (ACs 1, 2)\n- Create internal/tui/help/help.go - Help modal component with all shortcuts\n- Create internal/tui/help/styles.go - Lipgloss styling\n- Add ? key handler in app.go to open help modal\n- Group shortcuts by category: Navigation (h/l/j/k/arrows/Tab/[/]), Actions (a/Enter/Space), General (q/?)\n\nPhase 2: Toast Notification System (AC 3)\n- Create internal/tui/toast/toast.go - Toast component with auto-dismiss\n- Create internal/tui/toast/styles.go - Error/success styling\n- Add toasts field to Model struct\n- Update error handlers to show toast notifications\n\nPhase 3: Footer with Quick Reference (AC 4)\n- Add footer to View() in app.go showing key shortcuts\n- Display: h/l: columns | j/k: nav | a: add | ?: help | q: quit\n\nPhase 4: Documentation Updates (ACs 5, 6)\n- Update README.md TUI section with ? key for help\n- Ensure keyboard shortcuts table is complete\n\nPhase 5: Testing & Validation (ACs 7, 8, 9)\n- Run go vet and lint checks\n- Add integration tests for help modal and toast\n- Manual verification of all 62 ACs from tasks 15, 17-21
<!-- SECTION:PLAN:END -->
