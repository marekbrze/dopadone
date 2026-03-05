---
id: TASK-19
title: 'TUI 14D: Quick-Add Modal'
status: Done
assignee:
  - '@assistant'
created_date: '2026-03-03 12:31'
updated_date: '2026-03-04 06:00'
labels:
  - tui
  - mvp
  - phase2
dependencies:
  - TASK-15
  - TASK-21
references:
  - internal/tui/app.go
  - internal/tui/commands.go
  - internal/tui/messages.go
  - internal/db/querier.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement quick-add modal for creating subareas, projects, and tasks with keyboard.

**Context-Aware Creation:**
- Uses focused column + selected parent to determine what to create
- Subareas → created under selected Area
- Projects → created under selected Subarea
- Tasks → created under selected Project

**UX Flow:**
1. Press "a" key → modal opens centered
2. Modal shows parent context (e.g., "New Project in: Work Tasks")
3. Single title input field with focus
4. Enter creates item, Escape cancels
5. On success: close modal, refresh column, focus new item
6. On error: inline error message in modal

**Validation:**
- Title required (non-empty after trim)
- Disallow newlines and control characters only

**Visual Design:**
- Centered modal overlay
- Responsive width: 40-60% of terminal
- Clear visual hierarchy with border
- Parent context shown in header

**Assumptions:**
- A valid parent is always selected (app ensures selection on data load)
- Small datasets (no pagination needed)
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Press "a" key opens quick-add modal centered on screen
- [x] #2 Modal displays parent context (e.g., "New Project in: [Parent Name]")
- [x] #3 Modal has single title input field with cursor focus
- [x] #4 Enter key creates item in focused column context with validation
- [x] #5 Escape key closes modal without creating (no changes)
- [x] #6 Column refreshes and focuses newly created item after successful creation
- [x] #7 Inline error message displayed in modal for creation failures (e.g., validation, DB errors)
- [x] #8 Input validation: title required (non-empty after trim), no newlines/control chars
- [x] #9 Modal width is responsive (40-60% of terminal width)
- [x] #10 Unit tests for modal open/close behavior
- [x] #11 Unit tests for item creation with valid/invalid input
- [x] #12 Unit tests for error handling and display
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Task Assessment

**Complexity:** Moderate (12 ACs, 6 phases)
**Size:** Appropriate - cohesive feature (quick-add modal) that should be implemented together
**Recommendation:** Do NOT split - the modal, validation, and creation logic are tightly coupled

**Dependencies Status:**
- ✅ Task-15 (Core TUI Framework) - Complete
- ✅ Task-21 (Data Loading) - Complete

Both dependencies are complete, so task-19 can start immediately.

---

## Phase 1: Modal Component Foundation (2-3 hours)

**Can start immediately - no dependencies**

### Track 1A: Modal Structure & Styling
Create `internal/tui/modal/modal.go`:
- Define Modal struct with fields: title, input (textinput.Model), errorMsg, width, height
- Implement `New(parentContext, entityType string) *Modal`
- Implement `Update(msg tea.Msg) (*Modal, tea.Cmd)` handling:
  - textinput events (typing)
  - Enter key (submit)
  - Escape key (cancel)
  - Character limit (max 255 chars)
- Implement `View() string` rendering:
  - Centered modal overlay
  - Responsive width (40-60% terminal width)
  - Border with title
  - Input field with cursor
  - Error message area (red text)
  - Hint text at bottom

Create `internal/tui/modal/styles.go`:
- ModalBorder style (rounded, accent color)
- InputField style (focused)
- ErrorText style (red foreground)
- HintText style (dimmed)
- TitleStyle (bold)

**Testing:**
- Unit test: modal creation with parent context
- Unit test: View() output format
- Visual test: render modal in isolation

### Track 1B: Input Validation (Parallel with 1A)
Create `internal/tui/modal/validation.go`:
- Implement `ValidateTitle(title string) error`:
  - Check non-empty after trim
  - Check no newlines (\n, \r)
  - Check no control characters (using unicode.IsControl)
  - Check max length (255 chars)
- Return specific errors:
  - ErrTitleEmpty
  - ErrTitleTooLong
  - ErrTitleInvalidChars

**Testing:**
- Unit test: valid titles pass
- Unit test: empty titles fail
- Unit test: titles with newlines fail
- Unit test: titles with control chars fail
- Unit test: whitespace-only titles fail

---

## Phase 2: App Integration (1-2 hours)

**Sequential dependency: Phase 1 must complete**

### Track 2A: Model Extension
Update `internal/tui/app.go`:
- Add field: `modal *modal.Modal` (pointer, nil when closed)
- Add field: `isModalOpen bool` (state flag)

Update `InitialModel()`:
- Initialize modal = nil, isModalOpen = false

### Track 2B: Event Routing
Update `Update()` in `app.go`:
- Add "a" key handler:
  - Check focused column
  - Determine parent context and entity type
  - Create modal with NewCreateModal()
  - Set isModalOpen = true
  - Return cmd: textinput.Focus()
- When isModalOpen == true:
  - Route all key events to modal.Update()
  - Handle modal.CloseMsg: close modal, isModalOpen = false
  - Handle modal.SubmitMsg: trigger creation command

### Track 2C: View Overlay
Update `View()` in `app.go`:
- After rendering base layout, if isModalOpen:
  - Calculate modal position (centered)
  - Render modal.View() as overlay using lipgloss.Place()
  - Use lipgloss overlay technique

**Testing:**
- Unit test: "a" key opens modal
- Unit test: modal state isolation
- Unit test: event routing when modal open
- Visual test: modal renders centered

---

## Phase 3: Create Commands (1-2 hours)

**Can start in parallel with Phase 2**

### Track 3A: Create Commands
Update `internal/tui/commands.go`:
- Implement `CreateSubareaCmd(repo, name, areaID) tea.Cmd`:
  - Generate UUID
  - Create domain.Subarea using domain.NewSubarea()
  - Convert to db.CreateSubareaParams
  - Call repo.CreateSubarea()
  - Return SubareaCreatedMsg
- Implement `CreateProjectCmd(repo, name, parentID, subareaID) tea.Cmd`:
  - Generate UUID
  - Create domain.Project using domain.NewProject()
  - Set defaults: Status=active, Priority=medium, Progress=0
  - Convert to db.CreateProjectParams
  - Call repo.CreateProject()
  - Return ProjectCreatedMsg
- Implement `CreateTaskCmd(repo, title, projectID) tea.Cmd`:
  - Generate UUID
  - Create domain.Task using domain.NewTask()
  - Set defaults: Status=todo, Priority=medium, IsNext=false
  - Convert to db.CreateTaskParams
  - Call repo.CreateTask()
  - Return TaskCreatedMsg

### Track 3B: Create Messages
Update `internal/tui/messages.go`:
- Add SubareaCreatedMsg { Subarea domain.Subarea; Err error }
- Add ProjectCreatedMsg { Project domain.Project; Err error }
- Add TaskCreatedMsg { Task domain.Task; Err error }

**Testing:**
- Unit test: each create command with mock repository
- Unit test: error handling for each command
- Unit test: domain validation errors propagate

---

## Phase 4: Post-Creation Flow (1 hour)

**Sequential dependency: Phase 2 and Phase 3 must complete**

### Track 4A: Success Handling
Update `Update()` in `app.go`:
- Handle SubareaCreatedMsg:
  - On success: close modal, reload subareas, find and select new subarea
  - On error: set modal error, keep open
- Handle ProjectCreatedMsg:
  - On success: close modal, reload projects, rebuild tree, select new project
  - On error: set modal error, keep open
- Handle TaskCreatedMsg:
  - On success: close modal, reload tasks, select new task
  - On error: set modal error, keep open

### Track 4B: Helper Functions
Add to `internal/tui/app.go`:
- Implement `getParentContext() (name, entityType string)`:
  - Based on focused column and selected item
  - Subareas: return selected area name, "Subarea"
  - Projects: return selected subarea name, "Project"
  - Tasks: return selected project name, "Task"
- Implement `reloadColumnData(column FocusColumn) tea.Cmd`:
  - Return appropriate Load*Cmd based on column
- Implement `findAndSelectNewItem(name, column)`:
  - Find item by name in data slice
  - Update selection index
  - Set focus to correct column

**Testing:**
- Unit test: success flow closes modal
- Unit test: error flow keeps modal open with error
- Unit test: new item selection works
- Integration test: end-to-end creation flow

---

## Phase 5: Validation Integration (30 min)

**Sequential dependency: Phase 1B and Phase 2 must complete**

### Track 5A: Modal Validation
Update `internal/tui/modal/modal.go`:
- On Enter key:
  - Run ValidateTitle(input.Value())
  - If error: set modal.errorMsg, keep modal open
  - If valid: return SubmitMsg with title

### Track 5B: Error Display
Update modal View():
- If errorMsg != "":
  - Render error above input field
  - Use ErrorText style (red, bold)

**Testing:**
- Unit test: Enter with empty input shows error
- Unit test: Enter with invalid chars shows error
- Unit test: Enter with valid input triggers submit

---

## Phase 6: Comprehensive Testing (2-3 hours)

**Throughout implementation and final validation**

### Track 6A: Modal Component Tests
Create `internal/tui/modal/modal_test.go`:
- TestNewModal: creation with various contexts
- TestModalView: rendering output
- TestModalInput: typing, backspace
- TestModalValidation: valid/invalid inputs
- TestModalSubmit: Enter key behavior
- TestModalCancel: Escape key behavior
- TestModalError: error display

### Track 6B: Create Command Tests
Create `internal/tui/create_test.go`:
- TestCreateSubareaCmd: success with mock repo
- TestCreateSubareaCmdError: DB error handling
- TestCreateProjectCmd: success with subarea parent
- TestCreateProjectCmd: success with project parent (nested)
- TestCreateProjectCmdError: DB error handling
- TestCreateTaskCmd: success with mock repo
- TestCreateTaskCmdError: DB error handling

### Track 6C: Integration Tests
Update `internal/tui/integration_test.go`:
- TestQuickAddSubareaFlow: "a" → type → Enter → verify created
- TestQuickAddProjectFlow: "a" → type → Enter → verify created
- TestQuickAddTaskFlow: "a" → type → Enter → verify created
- TestQuickAddCancel: "a" → Escape → verify no changes
- TestQuickAddValidationError: "a" → Enter (empty) → verify error shown
- TestQuickAddDBError: mock repo error → verify error shown

### Track 6D: Edge Cases
- Test no parent selected (empty column)
- Test rapid "a" key presses (debouncing)
- Test very long titles (truncation)
- Test Unicode characters in titles
- Test modal with small terminal (responsive)

**Coverage Target:** >85% for modal package, >90% for create commands

---

## Phase 7: Documentation (1 hour, parallel with Phase 6)

### Track 7A: Code Documentation
- Add package doc for `internal/tui/modal`
- Document Modal struct fields
- Document validation rules
- Document create command parameters
- Add inline comments for complex logic

### Track 7B: Architecture Documentation
Update `internal/tui/README.md` (if exists):
- Add modal component overview
- Document quick-add flow
- Document validation rules
- Add diagram of message flow

### Track 7C: User Documentation
Create/update user guide:
- Document "a" key shortcut
- Explain context-aware creation
- List validation rules
- Show example workflows

---

## Phase 8: Final Validation (30 min)

**Sequential - all phases must complete**

1. Run all tests: `go test ./internal/tui/... -v -cover`
2. Run linting: `golangci-lint run` or `go vet ./...`
3. Build verification: `go build ./cmd/projectdb`
4. Manual testing:
   - AC #1-9: All functionality tests
   - AC #10-12: Test coverage verification
5. Code review checklist:
   - No side effects in View()
   - Clean architecture (TUI → domain)
   - All errors handled
   - No magic numbers
   - Constants named appropriately

---

## Parallelization Strategy

**Sequential Dependencies:**
1. Phase 1 → Phase 2 (modal needed for integration)
2. Phase 3 → Phase 4 (commands needed for flow)
3. Phase 2 + Phase 3 → Phase 4 (both needed)
4. Phase 4 → Phase 5 (flow needed for validation)
5. Phase 1-5 → Phase 8 (all must complete)

**Parallel Opportunities:**
- Phase 1A and 1B can run simultaneously (2 developers)
- Phase 2 and Phase 3 can run simultaneously after Phase 1
- Phase 6 and Phase 7 can run simultaneously
- Track 6A, 6B, 6C can run in parallel (3 developers)

**Single Developer Timeline:**
- Phase 1: 2-3 hours
- Phase 2: 1-2 hours
- Phase 3: 1-2 hours (can overlap with Phase 2)
- Phase 4: 1 hour
- Phase 5: 30 min
- Phase 6: 2-3 hours
- Phase 7: 1 hour (can overlap with Phase 6)
- Phase 8: 30 min
- **Total: 8-11 hours**

**Team Timeline (3 developers):**
- Phase 1: 1.5 hours (dev 1: 1A, dev 2: 1B)
- Phase 2-3: 2 hours (dev 1: 2, dev 2: 3, dev 3: review/prep)
- Phase 4-5: 1 hour (dev 1)
- Phase 6: 1.5 hours (dev 1: 6A, dev 2: 6B, dev 3: 6C)
- Phase 7: 30 min (dev 2, parallel)
- Phase 8: 30 min (team)
- **Total: 5-6 hours wall-clock**

---

## Risk Assessment

**High Risk:**
1. **Modal overlay rendering issues**
   - Mitigation: Use lipgloss.Place() with careful positioning
   - Test with various terminal sizes

**Medium Risk:**
2. **Validation edge cases**
   - Mitigation: Comprehensive validation tests (Phase 1B)
   - Test Unicode, emojis, special chars

3. **Tree rebuild after project creation**
   - Mitigation: Reuse existing tree building logic from task-21
   - Test nested projects

**Low Risk:**
4. **Focus management after creation**
   - Mitigation: Clear helper function with tests
   - Fallback to first item if not found

---

## Files to Create/Modify

**New Files:**
- `internal/tui/modal/modal.go` (~150 lines)
- `internal/tui/modal/styles.go` (~50 lines)
- `internal/tui/modal/validation.go` (~40 lines)
- `internal/tui/modal/modal_test.go` (~200 lines)
- `internal/tui/create_test.go` (~150 lines)

**Modified Files:**
- `internal/tui/app.go` (+80 lines)
- `internal/tui/commands.go` (+90 lines)
- `internal/tui/messages.go` (+15 lines)
- `internal/tui/constants.go` (+10 lines)
- `internal/tui/integration_test.go` (+100 lines)

**Total Impact:** ~880 lines of code + tests

---

## Success Criteria

- All 12 acceptance criteria verified
- Test coverage >85% for new code
- Clean architecture maintained
- No regressions in existing tests
- Manual testing passes all edge cases
- Code passes linting and formatting checks
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Phase 2-3: Integrated modal into app and created commands

- Added modal field to Model struct in app.go

- Implemented event routing: 'a' key opens modal, events routed to modal when open

- Created CreateSubareaCmd, CreateProjectCmd, CreateTaskCmd in commands.go

- Added SubareaCreatedMsg, ProjectCreatedMsg, TaskCreatedMsg to messages.go

Phase 1: Created modal package with Modal component

- Implemented modal.go with New(), Update(), View() methods

- Added validation.go with ValidateTitle() for input validation

- Created styles.go with lipgloss styling for modal UI

- Modal supports 3 entity types: Subarea, Project, Task
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented quick-add modal for creating subareas, projects, and tasks with keyboard shortcut 'a'.

**Features Implemented:**
- Context-aware creation: Modal determines entity type based on focused column (Subareas→Subarea, Projects→Project, Tasks→Task)
- Parent context display: Modal shows parent name in title (e.g., "New Project in: Work Tasks")
- Single input field: Focused text input with 255 char limit
- Validation: Non-empty after trim, no newlines/control characters
- Responsive width: 40-60% of terminal width
- Inline error display: Errors shown in modal without closing
- Success flow: Closes modal, reloads column, focuses new item

**Technical Implementation:**
- Created modal package with Modal, validation, and styles
- Integrated modal into app.go with event routing and view overlay
- Added create commands (CreateSubareaCmd, CreateProjectCmd, CreateTaskCmd)
- Implemented post-creation handlers with error handling
- Comprehensive test coverage: modal (17 tests), validation (30 tests), create commands (14 tests)

**Files Modified:**
- internal/tui/app.go: Added modal integration, event routing, creation handlers
- internal/tui/commands.go: Added 3 create commands
- internal/tui/messages.go: Added 3 creation messages

**Files Created:**
- internal/tui/modal/modal.go: Modal component with Update/View methods
- internal/tui/modal/styles.go: Lipgloss styling
- internal/tui/modal/validation.go: Input validation
- internal/tui/modal/modal_test.go: Modal tests
- internal/tui/modal/validation_test.go: Validation tests
- internal/tui/create_test.go: Create command tests

**Testing:**
- All 127 TUI tests passing
- Build verification successful
- No regressions in existing tests
<!-- SECTION:FINAL_SUMMARY:END -->
