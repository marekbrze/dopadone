---
id: TASK-85
title: Space menu config actions don't work
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-12 06:05'
updated_date: '2026-03-12 06:26'
labels:
  - bug
  - tui
  - spacemenu
dependencies: []
references:
  - internal/tui/spacemenu/spacemenu.go
  - internal/tui/spacemenu/types.go
  - internal/tui/app.go
  - internal/tui/area_handlers.go
  - internal/tui/handlers.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When using the Space menu to access config (Space → c), the area management actions (n: New, e: Edit, d: Delete) don't trigger any action. The spacemenu.Update() function switches to StateConfig on 'c' key but doesn't handle 'n', 'e', 'd' keys to emit ActionMsg.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Pressing 'n' in config menu emits ActionCreateArea and opens area modal in create mode
- [x] #2 Pressing 'e' in config menu emits ActionEditArea and opens area modal in edit mode for current area
- [x] #3 Pressing 'd' in config menu emits ActionDeleteArea and opens area modal in delete mode for current area
- [x] #4 Pressing 'c' in main menu still switches to config submenu (existing behavior preserved)
- [x] #5 Unit tests cover all config menu key handlers
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze current implementation and confirm the bug
   - Review spacemenu.go Update() function - verify n/e/d keys not handled in StateConfig
   - Review areamodal modes (ModeCreate, ModeEdit, ModeDeleteConfirm)
   - Review app.go handleSpaceMenuAction() to understand current flow

2. Fix spacemenu key handling (internal/tui/spacemenu/spacemenu.go)
   - Add handlers for n, e, d keys when state is StateConfig
   - Each key should emit ActionMsg with corresponding action (ActionCreateArea, ActionEditArea, ActionDeleteArea)
   - Close the menu after emitting action (return CloseMsg along with ActionMsg)

3. Update app.go to handle area modal modes (internal/tui/app.go)
   - Modify handleOpenAreaModal() to accept optional mode parameter
   - OR create new handlers: handleOpenAreaModalCreate(), handleOpenAreaModalEdit(), handleOpenAreaModalDelete()
   - Update handleSpaceMenuAction() to call appropriate handler based on action type
   - For Edit/Delete: pass current area ID to pre-select in modal

4. Add unit tests for spacemenu (internal/tui/spacemenu/spacemenu_test.go)
   - Test n key in StateConfig emits ActionCreateArea
   - Test e key in StateConfig emits ActionEditArea
   - Test d key in StateConfig emits ActionDeleteArea
   - Test c key in StateMain still switches to StateConfig (regression test)
   - Test keys in StateMain dont trigger actions (only c works)

5. Manual verification
   - Run TUI and test Space -> c -> n (should open area modal in create mode)
   - Test Space -> c -> e (should open area modal in edit mode with current area)
   - Test Space -> c -> d (should open area modal in delete confirm mode)
   - Verify Space -> c still shows config menu

6. Update documentation (if needed)
   - Check docs/TUI.md for Space menu documentation
   - Update if behavior changes warrant documentation updates
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Technical Analysis:
- Root cause: spacemenu.Update() handles c key to switch to StateConfig but does not handle n/e/d keys to emit ActionMsg
- Actions are already defined in types.go: ActionCreateArea, ActionEditArea, ActionDeleteArea
- areamodal already supports modes: ModeCreate, ModeEdit, ModeDeleteConfirm
- handleSpaceMenuAction() in app.go already handles these actions but routes them all to handleOpenAreaModal()
- Need to: (1) emit ActionMsg from spacemenu for n/e/d, (2) open areamodal in correct mode based on action
- Task is small enough to implement in one go - no need to split

Fix: Edit mode now properly pre-fills current area name and color in the input field
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Fixed Space menu config actions (n/e/d keys) to properly emit ActionMsg and open area modal in the correct mode.

## Changes

1. **internal/tui/areamodal/area_modal.go**: Added `SetMode()` and `SetSelectedIndex()` methods to allow setting initial mode and selected index from outside the modal.

2. **internal/tui/spacemenu/spacemenu.go**: Added handlers for `n`, `e`, `d` keys when in `StateConfig` state to emit `ActionMsg` with corresponding actions:
   - `n` → `ActionCreateArea`
   - `e` → `ActionEditArea`  
   - `d` → `ActionDeleteArea`

3. **internal/tui/handlers.go**: Added `handleOpenAreaModalWithMode()` function that accepts an initial mode parameter and sets the modal mode accordingly. For delete mode, also triggers stats loading.

4. **internal/tui/app.go**: Updated `handleSpaceMenuAction()` to route each action to `handleOpenAreaModalWithMode()` with the correct mode:
   - `ActionCreateArea` → `ModeCreate`
   - `ActionEditArea` → `ModeEdit`
   - `ActionDeleteArea` → `ModeDeleteConfirm`

5. **internal/tui/spacemenu/spacemenu_test.go**: Added comprehensive tests for config menu key handlers covering all scenarios.

## Testing

- All unit tests pass (`go test ./...`)
- Build passes (`go build ./...`)
- Lint passes (`golangci-lint run`)

## Manual Testing Required

User should test in TUI:
- Space → c → n (opens area modal in create mode)
- Space → c → e (opens area modal in edit mode with current area)
- Space → c → d (opens area modal in delete confirm mode)

## Additional Fix

Edit mode now properly pre-fills the current area name and color in the input field before saving. Added `SetupForEdit()`, `SetupForCreate()`, and `SetupForDelete()` methods to areamodal for proper mode initialization.

- Added unit tests for areamodal (`internal/tui/areamodal/area_modal_test.go`)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Unit tests pass (go test ./internal/tui/spacemenu/...)
- [ ] #2 Manual testing in TUI confirms all key handlers work
- [x] #3 No regressions in existing Space menu behavior
<!-- DOD:END -->
