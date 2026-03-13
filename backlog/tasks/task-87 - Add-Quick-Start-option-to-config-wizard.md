---
id: TASK-87
title: Add Quick Start option to config wizard
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-12 18:52'
updated_date: '2026-03-13 07:20'
labels:
  - tui
  - onboarding
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the first-run wizard to offer a Quick Start option on the welcome screen. Users should see: (1) Quick Start - use local SQLite with defaults, (2) Custom Setup - full wizard for path/mode selection, (3) Exit. This is the foundation for simplified onboarding.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Welcome screen shows Quick Start and Custom Setup options
- [x] #2 Quick Start bypasses mode selection and config screens
- [x] #3 Quick Start uses default local SQLite path
- [x] #4 Custom Setup enters existing full wizard flow
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan for Quick Start Option

### Overview
Add Quick Start option to the config wizard welcome screen, allowing users to quickly start with local SQLite defaults or choose custom setup.

### Task Scope Analysis
This is a well-contained task focused on modifying the config wizard component. No task splitting needed as:
- All changes are within the configwizard package
- Estimated time: ~3-4 hours implementation + testing
- Changes are cohesive around the welcome screen flow

---

### Phase 1: Type and State Changes (Sequential)

**File: internal/tui/configwizard/types.go**

1. Add new welcome option type:
```go
type WelcomeOption int

const (
    WelcomeOptionQuickStart WelcomeOption = iota
    WelcomeOptionCustomSetup
    WelcomeOptionExit
)
```

2. Add selection state to Wizard struct:
- Add `selectedWelcomeOption int` field to track user selection

---

### Phase 2: Welcome Screen UI (Sequential, after Phase 1)

**File: internal/tui/configwizard/wizard.go**

1. Update New() constructor:
   - Initialize selectedWelcomeOption = 0 (Quick Start as default)

2. Update renderWelcome() function:
   - Replace single "Press Enter" with three selectable options:
     ```
     ▸ Quick Start - Use local SQLite with defaults (recommended)
       Custom Setup - Choose database mode and configure
       Exit
     ```
   - Show selection indicator (▸) for selected option
   - Update hint text to: "↑/↓: Navigate • Enter: Select"

3. Follow BubbleTea Golden Rules:
   - Account for borders in text width calculations
   - Truncate text explicitly to prevent wrapping
   - Use weight-based layout if needed

---

### Phase 3: Navigation and Flow Handling (Sequential, after Phase 2)

**File: internal/tui/configwizard/wizard.go**

1. Update handleUp() for welcome step:
   - Navigate between welcome options (wrap around)

2. Update handleDown() for welcome step:
   - Navigate between welcome options (wrap around)

3. Update handleEnter() for stepWelcome:
   - Quick Start (option 0):
     - Set mode = ModeLocal
     - Set localPath to cli.DefaultDBPath()
     - Skip to verification (stepSuccess via verifyAndSave())
   - Custom Setup (option 1):
     - Enter existing flow (stepModeSelection)
   - Exit (option 2):
     - Return tea.Quit

4. Ensure handleBack() still works:
   - ESC on welcome step should exit (existing behavior)

---

### Phase 4: Testing (Can be done in parallel with Phases 1-3 for TDD)

**File: internal/tui/configwizard/wizard_test.go**

1. Add tests for welcome option navigation:
   - Test navigation down through all options
   - Test navigation up through all options
   - Test wrap-around behavior

2. Add tests for Quick Start flow:
   - Test that Quick Start sets ModeLocal
   - Test that Quick Start uses default path
   - Test that Quick Start skips mode selection

3. Add tests for Custom Setup flow:
   - Test that Custom Setup enters stepModeSelection
   - Test that mode selection still works

4. Add tests for Exit option:
   - Test that Exit returns tea.Quit command

5. Update existing tests if needed:
   - Ensure TestWizardModeSelection_* tests still pass
   - Ensure TestWizardCancelOnWelcome still passes

---

### Phase 5: Documentation Updates (Sequential, after Phase 4)

**File: docs/TUI.md**

1. Update Config Wizard section (if exists):
   - Document the three welcome options
   - Document Quick Start behavior
   - Document navigation keys

2. Update keyboard shortcuts table:
   - Add welcome screen navigation

---

### Testing Strategy

1. **Unit Tests**:
   - All new functionality covered by table-driven tests
   - Test both success and edge cases
   - Test navigation wrap-around behavior

2. **Integration Tests**:
   - Manual verification of complete flows:
     - Quick Start → Success
     - Custom Setup → Local → Success
     - Custom Setup → Remote → Success
     - Custom Setup → Replica → Success
     - Exit from welcome screen

3. **Visual Verification**:
   - Test rendering in different terminal sizes
   - Verify text truncation works correctly
   - Check selection indicator visibility

---

### Acceptance Criteria Mapping

| AC | Implementation |
|----|----------------|
| #1 Welcome screen shows Quick Start and Custom Setup options | Phase 2: renderWelcome() update |
| #2 Quick Start bypasses mode selection and config screens | Phase 3: handleEnter() Quick Start path |
| #3 Quick Start uses default local SQLite path | Phase 3: Uses cli.DefaultDBPath() |
| #4 Custom Setup enters existing full wizard flow | Phase 3: handleEnter() Custom Setup path |

---

### Dependencies

- **No task dependencies**: This task is self-contained
- **Code dependencies**: 
  - cli.DefaultDBPath() (already exists)
  - Existing verifyAndSave() logic (reuse)
  - Existing styles (reuse)

---

### Parallel vs Sequential Work

**Sequential**:
- Phase 1 → Phase 2 → Phase 3 → Phase 5
- Phase 4 can be done TDD-style alongside Phases 1-3

**Can be parallelized**:
- Phase 4 tests can be written before implementation (TDD)
- Phase 5 documentation can start after Phase 2 is complete

---

### Estimated Time

| Phase | Time |
|-------|------|
| Phase 1: Types | 15 min |
| Phase 2: UI | 45 min |
| Phase 3: Flow | 30 min |
| Phase 4: Tests | 60 min |
| Phase 5: Docs | 15 min |
| **Total** | **2.75 hours** |

---

### Risk Assessment

| Risk | Mitigation |
|------|------------|
| Breaking existing wizard flows | Comprehensive test coverage before changes |
| UI layout issues on small terminals | Test with narrow terminal, use text truncation |
| State management complexity | Keep state changes minimal, reuse existing patterns |
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Added WelcomeOption type (QuickStart, CustomSetup, Exit) to types.go
- Added selectedWelcomeOption field to Wizard struct
- Updated renderWelcome() to show 3 selectable options with navigation
- Updated handleUp/handleDown for welcome screen navigation with wrap-around
- Updated handleEnter to handle Quick Start (skip to verification), Custom Setup (enter wizard), Exit (quit)
- All existing tests pass
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Added Quick Start option to the config wizard welcome screen, allowing users to quickly start with local SQLite defaults.

### Changes
- **types.go**: Added `WelcomeOption` enum with `QuickStart`, `CustomSetup`, `Exit` options
- **wizard.go**: 
  - Added `selectedWelcomeOption` field for tracking user selection
  - Updated `renderWelcome()` to display 3 selectable options with descriptions
  - Updated `handleUp()`/`handleDown()` for welcome screen navigation with wrap-around
  - Updated `handleEnter()` to handle all 3 options:
    - Quick Start: Sets mode to local with default path, skips to verification
    - Custom Setup: Enters existing full wizard flow (mode selection)
    - Exit: Quits the wizard

### Flow Changes
- Welcome screen now has navigable options instead of a single "Press Enter" prompt
- Quick Start bypasses mode selection and config screens entirely
- Custom Setup preserves existing behavior (mode → config → verify)

### Testing
- All 24 existing tests pass
- Build and lint pass
<!-- SECTION:FINAL_SUMMARY:END -->
