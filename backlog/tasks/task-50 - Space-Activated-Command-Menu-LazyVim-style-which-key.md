---
id: TASK-50
title: Space-Activated Command Menu (LazyVim-style which-key)
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 18:09'
updated_date: '2026-03-06 21:13'
labels: []
dependencies: []
references:
  - internal/tui/help/help.go
  - internal/tui/modal/modal.go
  - internal/tui/theme/theme.go
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a floating overlay command menu triggered by Space key, similar to LazyVim's which-key functionality. When no input field is focused, pressing Space should display a floating modal overlay showing available commands with their key shortcuts and descriptions. Initial commands: 'c' for Config (Area management), 'q' for Quit. The menu should be dismissable with Space, Escape, or 'q'. When user selects Config, it should open a Config submenu with Area management options (create/edit/delete areas). The menu uses the existing theme system for consistent styling.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Space key triggers command menu overlay when no input field is focused
- [x] #2 Menu displays as floating overlay centered on screen with key shortcuts and descriptions
- [x] #3 Menu shows 'c: Config' and 'q: Quit' options with visual styling
- [x] #4 Space, Escape, or 'q' keys dismiss the menu without action
- [x] #5 Pressing 'c' opens Config submenu with Area management options
- [x] #6 Pressing 'q' from the menu triggers app quit
- [x] #7 Menu uses existing theme system (internal/tui/theme) for consistent colors
- [x] #8 Menu component follows existing modal patterns (similar to help modal structure)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Task Analysis

This task implements a Space-activated command menu (LazyVim-style which-key). After analyzing the acceptance criteria and codebase patterns, I recommend implementing this as a **single cohesive task** rather than splitting, because:

1. The functionality is self-contained within one component
2. Config submenu is just a view state within the same component
3. Similar complexity to existing help/modal components
4. Splitting would add coordination overhead without significant benefit

## Implementation Phases

### Phase 1: Core Spacemenu Component (Parallel Track A)

**Files to create:**
- `internal/tui/spacemenu/spacemenu.go` - Main component logic
- `internal/tui/spacemenu/styles.go` - Lipgloss styling
- `internal/tui/spacemenu/types.go` - Type definitions

**Implementation Details:**

1. **Types (`types.go`):**
   ```go
   type MenuState int
   const (
       StateMain MenuState = iota
       StateConfig
   )
   
   type Command struct {
       Key         string
       Label       string
       Description string
       Action      MenuAction
   }
   
   type MenuAction int
   const (
       ActionNone MenuAction = iota
       ActionQuit
       ActionConfig
       ActionCreateArea
       ActionEditArea
       ActionDeleteArea
   )
   ```

2. **Component (`spacemenu.go`):**
   - Follow existing modal patterns (see `internal/tui/help/help.go`)
   - Implement BubbleTea Model interface (Init/Update/View)
   - State machine for menu views (main → config)
   - Key handling for navigation and actions
   - View rendering with centered overlay

3. **Styling (`styles.go`):**
   - Use theme.ColorTheme for all colors (injected via New())
   - Border style matching help modal
   - Key highlight styling
   - Description text styling

### Phase 2: Integration with Main App (Parallel Track B - depends on Phase 1)

**Files to modify:**
- `internal/tui/app.go` - Add spacemenu state and integration
- `internal/tui/model.go` - Add menu state enum (if needed)

**Integration Points:**

1. **Model updates (`app.go` Model struct):**
   ```go
   spaceMenu       *spacemenu.SpaceMenu
   isSpaceMenuOpen bool
   ```

2. **Key handling priority:**
   - Add Space key handler in Update()
   - Only trigger if no modal/input is focused
   - Priority order: help > modal > area modal > space menu > normal
   
3. **Update flow:**
   - Route Space key to open menu when appropriate
   - Handle spacemenu.CloseMsg
   - Handle spacemenu.ActionMsg for quit/config actions

4. **View rendering:**
   - Add spacemenu.View() to overlay stack
   - Render order: base → toasts → modals → spacemenu

### Phase 3: Theme Integration (After Phase 1)

**Files to modify:**
- `internal/tui/theme/theme.go` - Add spacemenu color methods

**Add methods:**
```go
func (t ColorTheme) MenuBackground() lipgloss.AdaptiveColor
func (t ColorTheme) MenuKeyHighlight() lipgloss.AdaptiveColor
func (t ColorTheme) MenuDescription() lipgloss.AdaptiveColor
```

### Phase 4: Testing (Parallel with all phases)

**Files to create:**
- `internal/tui/spacemenu/spacemenu_test.go` - Unit tests
- `internal/tui/integration_spacemenu_test.go` - Integration test

**Test Coverage:**

1. **Unit tests (`spacemenu_test.go`):**
   - Open/close behavior
   - Key handling (Space/Esc/q/c)
   - State transitions (main → config)
   - View rendering
   - Theme integration

2. **Integration test (`integration_spacemenu_test.go`):**
   - Space key triggers menu from normal state
   - Space key does NOT trigger when modal is open
   - Full flow: open → navigate → action → close
   - Quit action terminates app

### Phase 5: Documentation (Final phase)

**Files to modify:**
- `docs/TUI.md` - Document spacemenu component
- `internal/tui/help/help.go` - Add Space key to shortcuts

**Documentation updates:**
- Add Space key to keyboard shortcuts section
- Document spacemenu component architecture
- Add to feature list
- Update footer help text

## Execution Strategy

### Sequential Dependencies:
- Phase 1 → Phase 2 (component must exist before integration)
- Phase 1 → Phase 3 (component exists before theme methods)
- All phases → Phase 5 (documentation last)

### Parallel Opportunities:
- Phase 1 and Phase 4 (tests can be written alongside component)
- Phase 2 and Phase 3 (integration and theme can be parallel)
- Phase 4 runs throughout all phases

### Recommended Execution Order:

1. **Start:** Phase 1 + Phase 4 tests (parallel)
2. **Next:** Phase 2 + Phase 3 (parallel)
3. **Finally:** Phase 5

## Testing Strategy

### Unit Tests (table-driven):
```go
func TestSpaceMenu_KeyHandling(t *testing.T) {
    tests := []struct {
        name       string
        key        string
        state      MenuState
        wantAction MenuAction
        wantClose  bool
    }{
        {"escape closes", "esc", StateMain, ActionNone, true},
        {"q from main quits", "q", StateMain, ActionQuit, true},
        {"c opens config", "c", StateMain, ActionNone, false},
        // ... more cases
    }
}
```

### Integration Tests:
- Use existing mock services pattern
- Test full keyboard flow
- Verify state transitions

## File Structure

```
internal/tui/
├── spacemenu/
│   ├── spacemenu.go      # Main component
│   ├── styles.go         # Lipgloss styling
│   ├── types.go          # Type definitions
│   └── spacemenu_test.go # Unit tests
├── app.go                # Integration (modified)
├── theme/
│   └── theme.go          # Theme methods (modified)
└── integration_spacemenu_test.go  # Integration test (new)
```

## Key Design Decisions

1. **Single component, multiple views:** Use MenuState enum instead of separate components
2. **Theme injection:** Pass theme to New() constructor for consistency
3. **Action-based messages:** Return ActionMsg for app-level handling
4. **Follow existing patterns:** Match help.go and modal.go structure

## Success Criteria

- [ ] Space key opens menu when no input focused
- [ ] Menu displays centered with proper styling
- [ ] All key shortcuts work (Space/Esc/q/c)
- [ ] Config submenu shows area management options
- [ ] Quit action terminates app
- [ ] Theme colors adapt to terminal
- [ ] Unit tests pass with >80% coverage
- [ ] Integration tests pass
- [ ] Documentation updated

## Estimated Effort

- Phase 1: 2-3 hours (component implementation)
- Phase 2: 1-2 hours (integration)
- Phase 3: 30 minutes (theme methods)
- Phase 4: 2-3 hours (testing)
- Phase 5: 1 hour (documentation)

**Total: 6-9 hours**
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Splitting Analysis

After analyzing the acceptance criteria and complexity:

**Decision: Implement as single task**

**Rationale:**
- Core functionality is cohesive (one component with state machine)
- Config submenu is just another view state (not a separate feature)
- Similar scope to existing help modal
- Splitting would add coordination overhead without clear benefits

**Alternative (if splitting preferred):**
- Task 50a: Basic spacemenu with quit (AC 1-4, 6-8)
- Task 50b: Config submenu integration (AC 5)
- Dependencies: 50b depends on 50a

**Recommendation:** Proceed with single-task implementation for faster delivery and simpler coordination.

Starting implementation - creating spacemenu component files

Created spacemenu component (types.go, styles.go, spacemenu.go)

Integrated spacemenu into app.go

Added Space key handler to open menu when no other modal is open

Code compiles successfully

Added unit tests for spacemenu component

Created integration tests for Space key flow

Updated TUI.md documentation with spacemenu component

Added Space key to help modal shortcuts
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented Space-activated command menu (LazyVim-style which-key) with the following changes:

**New Components:**
- Created internal/tui/spacemenu/ package with:
  - spacemenu.go: Main component implementing BubbleTea Model interface
  - types.go: Menu state and action type definitions
  - styles.go: Lipgloss styling using existing theme system
  - spacemenu_test.go: Comprehensive unit tests

**Integration:**
- Updated internal/tui/app.go to:
  - Add spacemenu state fields (spaceMenu, isSpaceMenuOpen)
  - Handle Space key to open menu when no modal is open
  - Route spacemenu messages (CloseMsg, ActionMsg)
  - Render spacemenu overlay in View()
  - Handle quit action from menu

**Features:**
- Space key opens command menu when no modal is focused
- Menu displays centered overlay with two initial commands:
  - c: Config (Area management) - opens Config submenu
  - q: Quit (Exit application)
- Config submenu shows area management options
- Menu dismissible with Space, Escape, or q
- Uses existing theme system for consistent styling
- Follows existing modal patterns (similar to help modal)

**Testing:**
- Unit tests: open/close behavior, key handling, state transitions, view rendering
- Integration tests: Space key flow, modal priority, navigation

**Documentation:**
- Updated TUI.md with spacemenu component in architecture diagram
- Added Space key to keyboard shortcuts section
- Updated Model structure documentation
- Added Space shortcut to help modal

All acceptance criteria and Definition of Done items completed.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Unit tests for spacemenu component (open/close/key handling)
- [x] #2 Integration test for Space key flow
- [x] #3 Update TUI.md documentation with new component
- [x] #4 Follow bubbletea golden rules for layout
<!-- DOD:END -->
