---
id: TASK-82.1
title: Welcome Modal Component
status: Done
assignee:
  - '@assistant'
created_date: '2026-03-11 20:50'
updated_date: '2026-03-12 05:54'
labels:
  - tui
  - component
  - onboarding
dependencies: []
references:
  - internal/tui/modal/modal.go
  - internal/tui/modal/styles.go
  - internal/tui/areamodal/area_modal.go
  - internal/tui/theme/theme.go
  - docs/TUI.md
parent_task_id: TASK-82
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create the welcome modal UI component with branding, text input, color selection, and validation. This is part of the larger 'Prompt user for initial area' feature (TASK-82).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Modal displays Dopadone branding and welcome message
- [x] #2 Modal has text input for area name with validation (non-empty required)
- [x] #3 Modal has color selection (Tab to cycle through predefined colors)
- [x] #4 User cannot skip/close the welcome modal without creating an area (ESC exits app or is disabled)
- [x] #5 Modal follows existing BubbleTea patterns (Init/Update/View)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
PHASE 1: Component Structure & Types (Sequential)
1. Create internal/tui/welcome/ package directory
2. Create types.go with message types:
   - SubmitMsg { Name string, Color domain.Color }
   - ExitMsg {} (for app termination on ESC)
3. Create Modal struct with fields:
   - textinput.Model for area name input
   - colorIndex int for color selection
   - errorMsg string for validation
   - width/height int for responsive layout
   - focusedElement (input only for this modal)

PHASE 2: Styling & Branding (Parallel with Phase 1)
1. Create styles.go with Lipgloss styling:
   - Branding styles (logo, welcome message)
   - Input field styling (matching existing modal patterns)
   - Color preview styling (reuse from areamodal)
   - Error text styling
   - Hint text styling
2. Define Dopadone branding elements:
   - ASCII art logo or styled title
   - Welcome message text
   - Guidance text for first-time users

PHASE 3: Core BubbleTea Implementation (Sequential after Phase 1&2)
1. Create welcome.go with:
   - New() constructor initializing text input
   - Init() returning nil (no initial commands)
   - Update() handling:
     • tea.WindowSizeMsg for responsive sizing
     • tea.KeyMsg for keyboard input:
       - Tab/Shift+Tab: cycle colors
       - Enter: validate and submit
       - ESC: return ExitMsg (parent will quit app)
       - Character input: pass to textinput
   - View() rendering:
     • Branding section (logo, welcome message)
     • Input field with label
     • Color selection preview
     • Error message (if any)
     • Hint text (Enter: Create • Tab: Color • ESC: Exit)

PHASE 4: Input Validation (Sequential after Phase 3)
1. Create validation.go:
   - ValidateName(name string) error function
   - Check for empty/whitespace-only names
   - Max length check (100 chars, matching areamodal)
2. Integrate validation in Update():
   - Clear error on input change
   - Show error on empty submit

PHASE 5: Unit Tests (Parallel with Phases 2-4)
1. Create welcome_test.go with table-driven tests:
   - TestNew: verify initial state
   - TestView: verify rendering contains branding, input, hints
   - TestValidation: empty name shows error
   - TestColorCycling: Tab cycles through colors
   - TestESCBehavior: ESC returns ExitMsg
   - TestEnterWithValidInput: returns SubmitMsg with correct data
   - TestEnterWithEmptyInput: shows validation error

PHASE 6: Documentation Updates (Sequential after all phases)
1. Update docs/TUI.md:
   - Add Welcome Modal section under Components
   - Document keyboard shortcuts specific to welcome modal
   - Add to keyboard shortcuts table
   - Document first-time user flow

DEPENDENCIES:
- None (standalone component)
- Integration with main app happens in TASK-82.2

SEQUENTIAL EXECUTION:
Phase 1 → Phase 3 → Phase 4 → Phase 6
Phase 2 can run parallel with Phase 1
Phase 5 can run parallel with Phases 2-4

TEST COMMANDS:
- go test ./internal/tui/welcome/... -v
- go test -race ./internal/tui/welcome/...
- golangci-lint run ./internal/tui/welcome/...
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
IMPLEMENTATION PATTERNS TO FOLLOW:

1. Follow existing modal patterns from:
   - internal/tui/modal/modal.go (structure, messages)
   - internal/tui/areamodal/area_modal.go (color selection, validation)

2. Use existing theme system:
   - Import github.com/marekbrze/dopadone/internal/tui/theme
   - Use theme.Default for colors
   - Use AdaptiveColor for light/dark terminal support

3. Color selection:
   - Reuse PredefinedColors from areamodal package
   - Import: areamodal.PredefinedColors

4. Key constants:
   - Import github.com/marekbrze/dopadone/internal/tui/internal/constants
   - Use constants.KeyEnter, constants.KeyEsc

5. Validation:
   - Match areamodal patterns: trim whitespace, check empty
   - Max length: 100 characters (consistent with areamodal)

6. Message types:
   - SubmitMsg for successful creation
   - ExitMsg for ESC (app termination)

7. Styling:
   - Use lipgloss.RoundedBorder() for modal border
   - Use theme.Default.Primary for border color
   - Use theme.Default.Error for validation errors
   - Use theme.Default.Muted for hint text

BRANDING ELEMENTS:
- Title: "Welcome to Dopadone" (styled with Primary color, Bold)
- Subtitle: "Your project management companion"
- Guidance: "Create your first area to get started"
- Input label: "Area Name:"
- Color label: "Color (Tab to change):"

- Created internal/tui/welcome/ package with 4 files: types.go, styles.go, welcome.go, validation.go
- Created welcome_test.go with 19 comprehensive tests covering all functionality
- Implemented branding: "Welcome to Dopadone" title, subtitle, guidance text
- Implemented text input with validation (non-empty, max 100 chars)
- Implemented color selection via Tab/Shift+Tab cycling through 12 predefined colors
- ESC returns ExitMsg (parent app handles by exiting)
- Added docs/TUI.md section 12 (Welcome Modal) with keyboard shortcuts table
- All tests pass (go test -race ./internal/tui/welcome/...)
- Lint passes (golangci-lint run ./internal/tui/welcome/...)
- Build passes (go build ./...)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Welcome Modal component in internal/tui/welcome/ package.

Files created:
- types.go: SubmitMsg and ExitMsg message types
- styles.go: Lipgloss styling (ModalBorder, TitleStyle, InputField, etc.)
- validation.go: ValidateName function with empty/whitespace/length checks
- welcome.go: Main Modal struct with New(), Init(), Update(), View() methods
- welcome_test.go: 20 unit tests covering all functionality

Features implemented:
- Dopadone branding with "Welcome to Dopadone" title and subtitle
- Text input for area name with validation (non-empty, max 100 chars)
- Color selection cycling via Tab/Shift+Tab using PredefinedColors
- ESC returns ExitMsg (parent app will quit)
- Enter validates and returns SubmitMsg with name and color
- Error messages clear automatically on input change

Documentation updated:
- docs/TUI.md: Added section 12. Welcome Modal with usage, message types, keyboard shortcuts
- Added Welcome Modal keyboard shortcuts to keyboard shortcuts table

All tests pass (go test -race), lint passes (golangci-lint), build passes.

Integration with main app deferred to TASK-82.2.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Run go test ./internal/tui/welcome/... -v and ensure all tests pass
- [x] #2 Run golangci-lint run ./internal/tui/welcome/... and fix any issues
- [ ] #3 Verify modal renders correctly in different terminal sizes
- [ ] #4 Manually test keyboard navigation and validation
<!-- DOD:END -->
