---
id: TASK-82.2
title: Welcome Modal Integration and Testing
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-11 20:50'
updated_date: '2026-03-12 05:54'
labels:
  - tui
  - integration
  - onboarding
  - testing
dependencies:
  - TASK-82.1
parent_task_id: TASK-82
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Integrate the welcome modal into the main TUI app flow when no areas exist, implement auto-selection of created area, load data, and write comprehensive tests. This is part of the larger 'Prompt user for initial area' feature (TASK-82).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 When TUI loads and areas list is empty, show welcome modal instead of main UI
- [x] #2 After area creation, auto-select the new area and trigger data loading (subareas, projects, tasks)
- [x] #3 Existing area creation via Space menu continues to work unchanged
- [x] #4 Unit tests cover: welcome modal rendering, validation, and area creation flow
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
PHASE 1: State & Types (Sequential) - ~30 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Add welcome modal fields to Model struct (app.go):
   - welcomeModal *welcome.Modal
   - isWelcomeOpen bool

2. Update InitialModel() to initialize:
   - welcomeModal: nil
   - isWelcomeOpen: false

3. Import welcome package in app.go

TEST CHECKPOINT: go build ./internal/tui/...

PHASE 2: Message Handling in data_loader_handlers.go (Sequential after Phase 1) - ~45 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Modify handleAreasLoaded():
   - If msg.Err != nil: handle error (no change)
   - If len(msg.Areas) == 0:
     * Initialize welcome modal: m.welcomeModal = welcome.New()
     * Set isWelcomeOpen = true
     * Return early (skip normal loading flow)
   - If len(msg.Areas) > 0: existing behavior (no change)

2. Pass window size to welcome modal for proper centering

TEST CHECKPOINT: go test -run TestEmptyAreasShowsWelcome ./internal/tui/...

PHASE 3: Welcome Event Handlers (Sequential after Phase 2) - ~1 hour
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Create internal/tui/welcome_handlers.go:
   - handleWelcomeSubmit(msg welcome.SubmitMsg):
     * Close welcome modal
     * Create area via CreateAreaCmd()
     * Store pending area selection flag for auto-select
   - handleWelcomeExit(msg welcome.ExitMsg):
     * Return tea.Quit to exit application
   - handleWelcomeMessages(msg interface{}):
     * Router function returning (Model, tea.Cmd, bool)
     * Handle welcome.SubmitMsg
     * Handle welcome.ExitMsg

2. Add welcome message routing in app.go Update():
   - Early check: if model, cmd, handled := m.handleWelcomeMessages(msg); handled { return model, cmd }

3. Add auto-selection after AreaCreatedMsg:
   - Track if this is first area creation from welcome modal
   - Auto-select the new area (m.selectedAreaIndex = len(areas) - 1)
   - Trigger LoadSubareasCmd immediately

TEST CHECKPOINT: go test -run TestWelcomeFlow ./internal/tui/...

PHASE 4: Keyboard Event Routing (Sequential after Phase 3) - ~30 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Update keyboard_handler.go handleKeyPress():
   - Add check for isWelcomeOpen BEFORE other modal checks:
     if m.isWelcomeOpen && m.welcomeModal != nil {
         return m.handleWelcomeKeyPress(msg)
     }

2. Create handleWelcomeKeyPress():
   - Route key events to welcome modal Update()
   - Handle ctrl+c/ctrl+q to also exit app

3. Handle window size events:
   - Pass WindowSizeMsg to welcome modal when open

TEST CHECKPOINT: go test -run TestWelcomeKeyboard ./internal/tui/...

PHASE 5: View Rendering (Sequential after Phase 4) - ~30 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Update app.go View() method:
   - Add check BEFORE base view rendering:
     if m.isWelcomeOpen && m.welcomeModal != nil {
         return m.welcomeModal.View()
     }
   - Welcome modal uses its own centering (already in welcome.View())

TEST CHECKPOINT: Manual test - run app with empty DB

PHASE 6: Space Menu Compatibility (Sequential after Phase 5) - ~30 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Verify Space menu area creation path unchanged:
   - handleOpenAreaModal() creates areamodal.Modal
   - areamodal.SubmitMsg handled by existing handleAreaModalSubmit()
   - CreateAreaCmd called, AreaCreatedMsg handled

2. Add integration test:
   - Test that Space menu can create area when areas already exist
   - Test that welcome modal appears when areas empty
   - Test both paths produce valid AreaCreatedMsg

TEST CHECKPOINT: go test -run TestSpaceMenuCompat ./internal/tui/...

PHASE 7: Comprehensive Tests (Parallel with Phases 3-6) - ~2 hours
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Create internal/tui/welcome_integration_test.go:

1. TestEmptyAreasShowsWelcomeModal:
   - Model with mock services returning empty areas list
   - Send AreasLoadedMsg with empty areas
   - Assert isWelcomeOpen == true
   - Assert welcomeModal != nil

2. TestWelcomeModalSubmitCreatesArea:
   - Welcome modal open
   - Send welcome.SubmitMsg with name and color
   - Assert CreateAreaCmd returned
   - Assert modal closed

3. TestWelcomeExitQuitsApp:
   - Welcome modal open
   - Send welcome.ExitMsg
   - Assert tea.Quit returned

4. TestAutoSelectionAfterFirstAreaCreated:
   - Welcome flow active
   - AreaCreatedMsg received
   - Assert new area auto-selected
   - Assert LoadSubareasCmd triggered

5. TestSpaceMenuAreaCreationStillWorks:
   - Existing areas present
   - Open space menu, select create area
   - Submit area modal
   - Assert area created via areamodal path

6. TestWelcomeModalNotShownWhenAreasExist:
   - Model with areas already loaded
   - AreasLoadedMsg with 2+ areas
   - Assert isWelcomeOpen == false

7. TestKeyboardRoutingToWelcomeModal:
   - Welcome modal open
   - Send character key
   - Assert routed to welcome modal

TEST COMMANDS:
- go test -v ./internal/tui/... -run Welcome
- go test -race ./internal/tui/...
- golangci-lint run ./internal/tui/...

PHASE 8: Documentation Updates (Sequential after all phases) - ~30 min
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Update docs/TUI.md:
   - Add "First-Time User Flow" section:
     * When no areas exist, welcome modal appears
     * User must create first area to proceed
     * ESC exits the application
   - Update keyboard shortcuts table:
     * Add note that ESC behavior differs in welcome modal
   - Update "Data Loading" section:
     * Document welcome modal intercept on empty areas
     * Document auto-selection flow after first area

2. Update docs/START_HERE.md (if needed):
   - Add note about first-time user experience

DEPENDENCIES & EXECUTION ORDER:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DEPENDS ON: TASK-82.1 (welcome.Modal component - DONE)

SEQUENTIAL EXECUTION:
Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5 → Phase 6 → Phase 8

PARALLEL EXECUTION:
Phase 7 (Tests) can run in parallel with Phases 3-6

ESTIMATED TOTAL TIME: ~5-6 hours

FILES TO CREATE:
- internal/tui/welcome_handlers.go
- internal/tui/welcome_integration_test.go

FILES TO MODIFY:
- internal/tui/app.go
- internal/tui/data_loader_handlers.go
- internal/tui/keyboard_handler.go
- docs/TUI.md
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
IMPLEMENTATION PATTERNS TO FOLLOW:

1. State management (app.go):
   - Follow existing modal patterns: areaModal/isAreaModalOpen
   - welcomeModal field is pointer (*welcome.Modal) for nil check
   - isWelcomeOpen bool flag for quick state check

2. Message routing pattern:
   - Use handleAreaMessages() pattern for handleWelcomeMessages()
   - Return (Model, tea.Cmd, bool) for router functions
   - Check early in Update() switch before other handlers

3. Event handling:
   - welcome.SubmitMsg → CreateAreaCmd → AreaCreatedMsg
   - welcome.ExitMsg → tea.Quit (app exit)
   - WindowSizeMsg → pass to welcome modal for centering

4. Auto-selection after welcome:
   - Add isFromWelcomeFlow bool flag to track origin
   - In handleAreaCreated(): if isFromWelcomeFlow, auto-select + load
   - Clear flag after handling

5. Test patterns:
   - Use mock services from mocks/services.go
   - Follow integration_test.go patterns
   - Test both welcome modal path and Space menu path

6. Compatibility checks:
   - Space menu area creation must work unchanged
   - areamodal.SubmitMsg still handled by handleAreaModalSubmit()
   - CreateAreaCmd is shared between both paths

KEY FILES TO REFERENCE:
- internal/tui/app.go (Model struct, View method)
- internal/tui/keyboard_handler.go (key routing)
- internal/tui/data_loader_handlers.go (areas loading)
- internal/tui/area_handlers.go (area creation patterns)
- internal/tui/welcome/welcome.go (modal component)
- internal/tui/mocks/services.go (test mocks)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Integrated welcome modal into TUI app flow for first-time users.

Changes:
- Added welcome modal state fields to Model struct (welcomeModal, isWelcomeOpen, isFromWelcomeFlow)
- Modified handleAreasLoaded to show welcome modal when areas list is empty
- Created welcome_handlers.go with handlers for Submit/Exit/Auto-selection messages
- Added keyboard routing to welcome modal in handleKeyPress()
- Updated View() to render welcome modal when open
- Updated handleWindowSize() to pass resize events to welcome modal
- Added comprehensive integration tests in welcome_integration_test.go
- Updated docs/TUI.md with First-Time User Flow section

Auto-selection flow:
- When user creates first area via welcome modal, area is auto-selected
- Subareas are loaded immediately for the new area
- User transitions directly to main TUI after first area creation

Space menu compatibility:
- Existing area creation via Space menu continues to work unchanged
- Both paths use CreateAreaCmd but only welcome flow triggers auto-selection

All tests pass: go test -race ./internal/tui/...
Lint passes: golangci-lint run ./internal/tui/...
Build passes: go build ./...
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All tests pass: go test -race ./internal/tui/...
- [x] #2 Lint passes: golangci-lint run ./internal/tui/...
- [x] #3 Manual test: run app with empty database
- [x] #4 Docs updated: docs/TUI.md includes first-time user flow
<!-- DOD:END -->
