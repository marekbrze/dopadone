---
id: TASK-82
title: Prompt user for initial area
status: Done
assignee: []
created_date: '2026-03-11 10:18'
updated_date: '2026-03-12 05:54'
labels:
  - tui
  - onboarding
  - feature
dependencies:
  - TASK-82.1
  - TASK-82.2
references:
  - internal/tui/areamodal/area_modal.go
  - internal/tui/modal/modal.go
  - internal/tui/data_loader_handlers.go
  - docs/TUI.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When the user opens the TUI app for the first time and no areas exist, show a dedicated welcome modal with branding that prompts the user to create their first area. The user must create at least one area to continue (no skip option). After creating the initial area, automatically select it and load its data (subareas, projects, tasks).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 When TUI loads and areas list is empty, show a dedicated welcome modal instead of the main UI
- [ ] #2 Welcome modal displays Dopadone branding, welcome message, and guidance text
- [ ] #3 Welcome modal has text input for area name with validation (non-empty required)
- [ ] #4 Welcome modal has color selection (Tab to cycle through predefined colors)
- [ ] #5 User cannot skip/close the welcome modal without creating an area (ESC exits app or is disabled)
- [ ] #6 After area creation, auto-select the new area and trigger data loading (subareas, projects, tasks)
- [ ] #7 Existing area creation via Space menu continues to work unchanged
- [ ] #8 Unit tests cover: welcome modal rendering, validation, and area creation flow
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create internal/tui/welcome/ package with WelcomeModal component
2. Implement Modal struct with text input, color selection, branding styles
3. Add welcome modal state to Model struct (welcomeModal *welcome.Modal, isWelcomeOpen bool)
4. Add WelcomeAreaCreatedMsg message type
5. Modify handleAreasLoaded to check empty areas and show welcome modal
6. Create welcome_handlers.go for handling welcome modal messages
7. Update View() in renderer.go to render welcome modal
8. Write unit tests in welcome/welcome_test.go
9. Integration test for first-time user flow
10. Update TUI.md documentation with welcome modal section

Subtasks:
- TASK-82.1: Welcome Modal Component (prerequisite)
- TASK-82.2: Integration and Testing (depends on 82.1)

Execution order: TASK-82.1 → TASK-82.2
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Split into two subtasks for better tracking:\n- TASK-82.1: Welcome Modal Component (standalone UI component)\n- TASK-82.2: Integration and Testing (depends on 82.1)\n\nBoth subtasks must be completed to fulfill the original requirements.
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Run go test ./... and ensure all tests pass
- [ ] #2 Run golangci-lint run and fix any issues
- [ ] #3 Manually test the welcome flow by starting with empty database
- [ ] #4 Update TUI.md documentation if the architecture changes significantly
<!-- DOD:END -->
