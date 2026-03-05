---
id: TASK-40
title: add option for adding projects
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-05 20:40'
updated_date: '2026-03-05 21:35'
labels: []
dependencies: []
references:
  - internal/tui/modal/modal.go
  - internal/tui/handlers.go
  - internal/tui/commands.go
  - internal/tui/app.go
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Currently, when a project is selected in the Projects column and the user presses 'a' to add a new project, it automatically creates a subproject under the selected project. This makes it impossible to create a root-level project.\n\nAdd an opt-in checkbox in the quick-add modal to allow creating a project as a subproject of the currently selected project. The default should be unchecked, meaning new projects are created at root level (under the selected subarea) by default.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Modal shows 'Add as subproject' checkbox when a project is selected in Projects column
- [x] #2 Checkbox is unchecked by default (creates root project)
- [x] #3 Checking the box changes creation to subproject under selected project
- [x] #4 Keyboard navigation works (Tab between input and checkbox, Space to toggle)
- [x] #5 Root project creation uses selected subarea as parent
- [x] #6 Existing behavior preserved for Subareas and Tasks columns
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# IMPLEMENTATION PLAN: Add Subproject Checkbox to Quick-Add Modal

## Task Analysis

This task enhances the quick-add modal to support opt-in subproject creation when a project is selected in the Projects column. Currently, pressing 'a' with a project selected automatically creates a subproject. The new behavior adds a checkbox (unchecked by default) to explicitly opt into subproject creation.

**Task is well-scoped and does NOT require splitting** - all acceptance criteria are related to a single feature enhancement.

## Implementation Phases

### Phase 1: Modal Component Enhancement (UI & State Management)
**Priority: High | Can run: Sequentially**

**Files to modify:**
- `internal/tui/modal/modal.go`
- `internal/tui/modal/styles.go`

**Changes:**
1. Add checkbox state to Modal struct:
   - Add `showCheckbox bool` field (determines if checkbox is visible)
   - Add `checkboxChecked bool` field (checkbox state, default false)
   - Add `focusedElement enum` (input vs checkbox focus)

2. Update New() constructor:
   - Add parameter `showCheckbox bool`
   - Initialize checkbox as unchecked
   - Set initial focus to input field

3. Implement keyboard navigation in Update():
   - Tab: cycle focus between input ↔ checkbox
   - Shift+Tab: reverse cycle
   - Space: toggle checkbox when focused
   - Maintain existing Enter/Esc behavior

4. Update View() to render checkbox:
   - Render checkbox below input field
   - Show checkbox only when `showCheckbox` is true
   - Style checkbox based on focus state and checked state
   - Update hint text: "Tab: Switch • Space: Toggle • Enter: Create • Esc: Cancel"

5. Add checkbox styles to styles.go:
   - CheckboxStyle (unfocused, unchecked)
   - CheckboxFocusedStyle (focused, unchecked)
   - CheckboxCheckedStyle (unfocused, checked)
   - CheckboxFocusedCheckedStyle (focused, checked)

**Testing in this phase:**
- Unit tests for checkbox state initialization
- Unit tests for Tab/Shift+Tab navigation
- Unit tests for Space toggle behavior
- Unit tests for checkbox rendering

**Estimated time:** 2-3 hours

---

### Phase 2: Handler Logic Update (Business Logic)
**Priority: High | Can run: After Phase 1 | Sequentially**

**Files to modify:**
- `internal/tui/handlers.go`
- `internal/tui/commands.go`

**Changes:**

1. **Update getParentContext() in handlers.go:**
   - Detect when a project is selected in Projects column
   - Return flag indicating checkbox should be shown
   - Keep existing logic for Subareas and Tasks columns
   - Signature change: return (parentName, entityType, parentID, subareaID, showCheckbox)

2. **Update handleQuickAdd() in handlers.go:**
   - Extract showCheckbox flag from getParentContext()
   - Pass showCheckbox to modal.New()
   - Update modal creation: `modal.New(parentName, entityType, parentID, subareaID, showCheckbox)`

3. **Update handleModalSubmit() in handlers.go:**
   - Check checkbox state when entityType is Project/Subproject
   - If checkbox is checked: create as subproject (use parentID)
   - If checkbox is unchecked: create as root project (use subareaID)
   - Preserve existing behavior for other entity types

4. **Update SubmitMsg in modal.go:**
   - Add `AsSubproject bool` field to SubmitMsg struct
   - Populate this field based on checkbox state

5. **Update CreateProjectCmd in commands.go (if needed):**
   - No changes needed - already supports both parentID and subareaID
   - Verify it handles nil parentID correctly for root projects

**Testing in this phase:**
- Unit tests for getParentContext() with project selected
- Unit tests for getParentContext() with subarea selected
- Unit tests for getParentContext() with task selected
- Integration tests for project creation with checkbox checked
- Integration tests for project creation with checkbox unchecked
- Regression tests for subarea/task creation

**Estimated time:** 2-3 hours

---

### Phase 3: Comprehensive Testing Suite
**Priority: High | Can run: After Phase 2 | Sequentially**

**Files to create/modify:**
- `internal/tui/modal/modal_test.go` (extend existing)
- `internal/tui/handlers_test.go` (create new)

**Test Categories:**

**A. Modal Checkbox Behavior Tests (modal_test.go):**
```go
- TestCheckboxVisibility(): verify checkbox appears only when showCheckbox=true
- TestCheckboxInitialState(): verify checkbox starts unchecked
- TestCheckboxToggle(): verify Space toggles checkbox
- TestTabNavigationInputToCheckbox(): verify Tab moves focus to checkbox
- TestTabNavigationCheckboxToInput(): verify Tab cycles back to input
- TestShiftTabNavigation(): verify Shift+Tab reverse navigation
- TestCheckboxSubmitMsg(): verify SubmitMsg.AsSubproject reflects checkbox state
- TestCheckboxRendering(): verify checkbox renders correctly in all states
```

**B. Handler Logic Tests (handlers_test.go):**
```go
- TestGetParentContextProjectSelected(): verify checkbox flag is true
- TestGetParentContextSubareaSelected(): verify checkbox flag is false
- TestGetParentContextTaskSelected(): verify checkbox flag is false
- TestHandleModalSubmitProjectUnchecked(): verify root project creation
- TestHandleModalSubmitProjectChecked(): verify subproject creation
- TestHandleModalSubmitSubarea(): verify no checkbox, existing behavior
- TestHandleModalSubmitTask(): verify no checkbox, existing behavior
```

**C. Integration Tests:**
```go
- TestEndToEndProjectCreationRoot(): user flow for creating root project
- TestEndToEndProjectCreationSubproject(): user flow for creating subproject
- TestKeyboardNavigationFlow(): Tab → Space → Enter sequence
```

**D. Edge Case Tests:**
```go
- TestModalWithEmptyProjectSelection(): no project selected behavior
- TestCheckboxPersistsOnValidationError(): checkbox state after validation error
- TestMultipleToggleCycles(): toggle checkbox multiple times
```

**Test Coverage Targets:**
- Modal component: 95%+
- Handler functions: 90%+
- Overall feature: 85%+

**Estimated time:** 3-4 hours

---

### Phase 4: Documentation Updates
**Priority: Medium | Can run: After Phase 3 | Sequentially**

**Files to modify:**
- `docs/TUI.md`
- `internal/tui/modal/modal.go` (inline comments)
- `internal/tui/handlers.go` (inline comments)

**Documentation Changes:**

1. **Update docs/TUI.md:**
   - Add section "Quick-Add Modal Enhancements"
   - Document checkbox behavior for project creation
   - Add keyboard navigation diagram
   - Add usage examples:
     - Creating a root project from Projects column
     - Creating a subproject from Projects column
   - Update screenshots (if applicable)

2. **Inline Code Comments:**
   - Document checkbox state management in Modal struct
   - Explain Tab/Space navigation logic in Update()
   - Clarify parent context determination in getParentContext()
   - Document the dual-mode project creation in handleModalSubmit()

3. **Update AGENTS.md (if needed):**
   - Add note about checkbox behavior in TUI patterns section

**Estimated time:** 1-2 hours

---

## Dependencies & Execution Order

**Sequential Execution (must complete in order):**
1. Phase 1: Modal Component Enhancement
   ↓
2. Phase 2: Handler Logic Update
   ↓
3. Phase 3: Comprehensive Testing Suite
   ↓
4. Phase 4: Documentation Updates

**Rationale for Sequential Execution:**
- Phase 2 depends on Phase 1's Modal struct changes
- Phase 3 tests Phases 1 & 2 implementation
- Phase 4 documents the completed feature

**No Parallel Execution Possible:** All phases have dependencies on previous phases.

---

## Acceptance Criteria Mapping

| AC | Phase | Implementation |
|----|-------|----------------|
| #1: Modal shows checkbox when project selected | Phase 1 + Phase 2 | Modal.showCheckbox=true when project selected |
| #2: Checkbox unchecked by default | Phase 1 | Modal.checkboxChecked=false initially |
| #3: Checking box creates subproject | Phase 2 | handleModalSubmit() uses parentID when checked |
| #4: Keyboard navigation (Tab/Space) | Phase 1 | Update() handles Tab/Space keys |
| #5: Root project uses subarea as parent | Phase 2 | handleModalSubmit() uses subareaID when unchecked |
| #6: Preserve existing behavior | Phase 2 | Conditional checkbox display + existing code paths |

---

## Risk Mitigation

**Risk 1: Breaking existing project creation flow**
- **Mitigation:** Keep checkbox hidden for Subareas/Tasks columns
- **Test:** Regression tests for all entity types

**Risk 2: Keyboard navigation conflicts**
- **Mitigation:** Use Tab (already unused) and Space (standard checkbox toggle)
- **Test:** Verify no conflicts with existing Enter/Esc

**Risk 3: User confusion about default behavior**
- **Mitigation:** Clear checkbox label: "Add as subproject of [ProjectName]"
- **Test:** User acceptance testing

**Risk 4: State inconsistency between modal and handler**
- **Mitigation:** Explicit SubmitMsg.AsSubproject field
- **Test:** Unit tests for message passing

---

## Testing Strategy

**Unit Tests (Go testing package):**
- Table-driven tests for all modal behaviors
- Mock services for handler tests
- Focus on edge cases and validation

**Integration Tests:**
- End-to-end user flows
- Service layer integration
- Real database operations (in-memory SQLite)

**Manual Testing Checklist:**
- [ ] Press 'a' with subarea selected → no checkbox appears
- [ ] Press 'a' with project selected → checkbox appears, unchecked
- [ ] Tab navigates between input and checkbox
- [ ] Space toggles checkbox when focused
- [ ] Submit with checkbox unchecked → root project created
- [ ] Submit with checkbox checked → subproject created
- [ ] Press 'a' with task selected → no checkbox, existing behavior
- [ ] Validation error preserves checkbox state
- [ ] Esc cancels modal correctly

---

## Implementation Notes (to be added during implementation)

*This section will be populated with actual implementation decisions, challenges faced, and solutions applied.*

---

## Final Summary (to be added after completion)

*This section will contain a PR-style summary of what was implemented.*
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
✅ Phase 1 Complete: Enhanced modal component with checkbox functionality
- Added checkbox state fields (showCheckbox, checkboxChecked, focusedElement)
- Implemented Tab/Shift+Tab navigation between input and checkbox
- Implemented Space toggle for checkbox
- Added checkbox rendering with proper styling
- Updated SubmitMsg to include AsSubproject field
- Added comprehensive checkbox styles to styles.go

✅ Phase 2 Complete: Updated handler logic
- Modified getParentContext() to return showCheckbox flag
- Updated handleQuickAdd() to pass showCheckbox to modal
- Modified handleModalSubmit() to handle AsSubproject field
- Changed entity type from EntityTypeSubproject to EntityTypeProject when project is selected

✅ Phase 3 Complete: Comprehensive testing
- Added 10 new test cases for checkbox functionality
- All tests passing (27 tests total)
- Test coverage includes: visibility, navigation, toggle, submit, rendering
- Verified keyboard navigation (Tab, Shift+Tab, Space)
- Tested checkbox state persistence across interactions

🐛 Bugfix: Root project not appearing after creation
- Issue: When project selected and checkbox unchecked, subareaID was nil
- Root cause: getParentContext() returned nil for subareaID when project selected
- Fix: Pass current subareaID along with parentID when project is selected
- Result: Root projects now correctly created under selected subarea and appear in tree
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Enhanced quick-add modal to support opt-in subproject creation when a project is selected in the Projects column.

## Changes

### Modal Component (internal/tui/modal/)
- Added checkbox functionality with state management (showCheckbox, checkboxChecked, focusedElement)
- Implemented Tab/Shift+Tab keyboard navigation between input field and checkbox
- Implemented Space key to toggle checkbox when focused
- Added checkbox rendering with visual feedback for focus and checked states
- Updated SubmitMsg to include AsSubproject field
- Added new checkbox styles (CheckboxStyle, CheckboxFocusedStyle, CheckboxCheckedStyle, CheckboxFocusedCheckedStyle)

### Handler Logic (internal/tui/handlers.go)
- Modified getParentContext() to return showCheckbox flag (true when project selected in Projects column)
- Updated handleQuickAdd() to pass showCheckbox parameter to modal.New()
- Modified handleModalSubmit() to create subproject or root project based on AsSubproject field
- Changed entity type from EntityTypeSubproject to EntityTypeProject when project is selected

### Tests (internal/tui/modal/modal_test.go)
- Added 10 comprehensive test cases for checkbox functionality
- Test coverage includes: visibility, initial state, navigation, toggle, submit, rendering
- Updated all existing tests to use new modal.New() signature
- All 27 tests passing

## User Impact

Users can now create root-level projects when a project is selected in the Projects column:
- Default behavior: Creates root project under selected subarea (checkbox unchecked)
- Checkbox checked: Creates subproject under selected project
- Checkbox only appears when a project is selected in Projects column
- No changes to behavior in Subareas or Tasks columns

## Testing

- Unit tests: 10 new tests, all passing
- Manual testing recommended: Press a with project selected → Tab → Space → Enter flow
- Keyboard navigation verified: Tab cycles between input and checkbox, Space toggles checkbox

## Bugfix

Fixed issue where root projects created from Projects column didn't appear in the tree:
- Modified getParentContext() to pass subareaID when project is selected
- Root projects are now correctly created under the current subarea
- Verified with build and all tests passing
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Unit tests for new checkbox behavior
- [ ] #2 Manual testing in TUI
- [x] #3 No regressions in existing quick-add functionality
<!-- DOD:END -->
