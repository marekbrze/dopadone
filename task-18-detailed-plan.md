# Task-18: Detailed Implementation Plan
## Navigation & State Persistence

### Context & Dependencies

**Completed Foundation (Ready to Build On):**
- ✅ Task-15: Core TUI Framework - Focus management, column navigation (h/l/Tab)
- ✅ Task-20: Tree Rendering Package - **Already provides:**
  - `GetNextVisibleNode()`, `GetPrevVisibleNode()` - Navigation helpers
  - `TreeNode` with `IsExpanded` state management
  - `Render()` with selection highlighting
  - Unlimited nesting support
- ✅ Task-21: Data Loading & Integration - **Already provides:**
  - Cascade loading: Area → Subarea → Projects → Tasks
  - Selection indices tracking (`selectedAreaIndex`, `selectedSubareaIndex`, etc.)
  - Empty state messages
  - Repository integration
  - Spinner feedback

**What Task-18 Actually Needs to Add:**
1. **State Persistence Layer** - Save/restore selections and tree expand state per area
2. **j/k Navigation** - Wire up existing tree helpers to key handlers
3. **Area Switching with State** - [/] keys with state save/restore
4. **Visual Polish** - Bold + inverted styling for selections
5. **Testing** - Comprehensive unit and integration tests

---

## Architecture Overview

### Current State (Task-21 Completion)
```
Model Structure:
├── focus: FocusColumn (which column is focused)
├── selectedAreaIndex, selectedSubareaIndex, selectedProjectIndex, selectedTaskIndex
├── areas[], subareas[], projects[], tasks[]
├── isLoading flags
└── spinner.Model

Update Loop:
- h/l/Tab: Column focus switching (already implemented)
- Area selection: Triggers cascade data load (already implemented)
```

### New Architecture (Task-18 Addition)
```
Model Structure (Additions):
├── areaStates: map[string]*AreaState  // Keyed by area ID
│   ├── selectedSubareaIndex: int
│   ├── selectedProjectIndex: int  
│   ├── selectedTaskIndex: int
│   └── expandedProjects: map[string]bool  // Project IDs -> expanded state
└── selectedProjectID: string  // For tree navigation (ID, not index)

Update Loop (Additions):
- j/k: In-column navigation (wrap-around, tree-aware for projects)
- Enter/Space: Toggle expand/collapse (projects column only)
- [/]: Area switching with state save/restore
```

---

## Implementation Plan

### Phase 1: State Management Layer (2 hours)

#### Track A: Per-Area State Storage (1 hour)
**File**: `internal/tui/model.go`

1. **Define AreaState struct:**
```go
type AreaState struct {
    SelectedSubareaIndex int
    SelectedProjectIndex int
    SelectedTaskIndex    int
    ExpandedProjects     map[string]bool  // Project ID -> is expanded
}
```

2. **Extend Model struct:**
```go
type Model struct {
    // ... existing fields ...
    
    // State persistence
    areaStates         map[string]*AreaState  // Keyed by area.ID
    selectedProjectID  string                  // Current selected project (for tree nav)
    
    // Tree reference (for navigation)
    projectTree        *tree.TreeNode
}
```

3. **Add state management methods:**
```go
func (m *Model) getAreaState(areaID string) *AreaState {
    if m.areaStates[areaID] == nil {
        m.areaStates[areaID] = &AreaState{
            SelectedSubareaIndex: 0,
            SelectedProjectIndex: 0,
            SelectedTaskIndex:    0,
            ExpandedProjects:     make(map[string]bool),
        }
    }
    return m.areaStates[areaID]
}

func (m *Model) saveCurrentAreaState() {
    areaID := m.areas[m.selectedAreaIndex].ID
    state := m.getAreaState(areaID)
    state.SelectedSubareaIndex = m.selectedSubareaIndex
    state.SelectedProjectIndex = m.selectedProjectIndex
    state.SelectedTaskIndex = m.selectedTaskIndex
    // Save tree expand state
    m.saveTreeExpandState(state)
}

func (m *Model) restoreAreaState(areaID string) {
    state := m.getAreaState(areaID)
    m.selectedSubareaIndex = state.SelectedSubareaIndex
    m.selectedProjectIndex = state.SelectedProjectIndex
    m.selectedTaskIndex = state.SelectedTaskIndex
    // Restore tree expand state
    m.restoreTreeExpandState(state)
}
```

**Tests** (`internal/tui/state_test.go`):
- `TestAreaStateInitialization`: New area gets default state
- `TestAreaStateSaveRestore`: State persists correctly
- `TestAreaStateIsolation`: Different areas have independent states
- `TestTreeExpandStatePersistence`: Expand/collapse survives area switches

---

#### Track B: Selection Index Helpers (1 hour, can parallel with Track A)
**File**: `internal/tui/app.go`

1. **Navigation helpers with wrap-around:**
```go
func (m *Model) navigateUp(column FocusColumn) {
    switch column {
    case FocusSubareas:
        if len(m.subareas) == 0 {
            return  // No-op on empty
        }
        if m.selectedSubareaIndex == 0 {
            m.selectedSubareaIndex = len(m.subareas) - 1  // Wrap to last
        } else {
            m.selectedSubareaIndex--
        }
    case FocusProjects:
        m.navigateTreeUp()  // Special tree navigation
    case FocusTasks:
        if len(m.tasks) == 0 {
            return
        }
        if m.selectedTaskIndex == 0 {
            m.selectedTaskIndex = len(m.tasks) - 1
        } else {
            m.selectedTaskIndex--
        }
    }
}

func (m *Model) navigateDown(column FocusColumn) {
    // Similar logic with wrap-around
}

func (m *Model) navigateTreeUp() {
    if m.projectTree == nil {
        return
    }
    prevNode := tree.GetPrevVisibleNode(m.projectTree, m.selectedProjectID)
    if prevNode != nil {
        m.selectedProjectID = prevNode.ID
        m.syncTreeSelectionToIndex()
    }
}

func (m *Model)navigateTreeDown() {
    // Similar using GetNextVisibleNode
}

func (m *Model) syncTreeSelectionToIndex() {
    // Convert selectedProjectID to index for rendering
    visibleNodes := tree.GetAllVisibleNodes(m.projectTree)
    for i, node := range visibleNodes {
        if node.ID == m.selectedProjectID {
            m.selectedProjectIndex = i
            break
        }
    }
}
```

2. **Empty column helper:**
```go
func (m *Model) isEmpty(column FocusColumn) bool {
    switch column {
    case FocusSubareas:
        return len(m.subareas) == 0
    case FocusProjects:
        return m.projectTree == nil
    case FocusTasks:
        return len(m.tasks) == 0
    }
    return true
}
```

**Tests** (`internal/tui/navigation_test.go`):
- `TestNavigateUpWrap`: Last → First wrapping works
- `TestNavigateDownWrap`: First → Last wrapping works
- `TestNavigateEmptyColumn`: No-op on empty columns
- `TestNavigateTreeWrap`: Tree navigation wraps correctly
- `TestTreeSelectionSynchronization`: ID ↔ Index conversion works

---

### Phase 2: Key Handler Integration (2 hours)

#### Track A: Navigation Keys (1 hour)
**File**: `internal/tui/app.go` (Update method)

1. **Add key handlers:**
```go
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        // ... existing h/l/Tab handlers ...
        
        case "j", "down":
            if !m.isEmpty(m.focus) {
                m.navigateDown(m.focus)
            }
        case "k", "up":
            if !m.isEmpty(m.focus) {
                m.navigateUp(m.focus)
            }
        case "enter", " ":  // Space
            if m.focus == FocusProjects {
                m.toggleTreeExpand()
            }
        case "[":
            m.switchToPreviousArea()
        case "]":
            m.switchToNextArea()
        }
    }
    // ... rest of Update ...
}

func (m *Model) toggleTreeExpand() {
    if m.selectedProjectID == "" {
        return
    }
    node := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
    if node != nil && !node.IsLeaf() {
        node.ToggleExpanded()
        // Update state
        areaID := m.areas[m.selectedAreaIndex].ID
        state := m.getAreaState(areaID)
        state.ExpandedProjects[m.selectedProjectID] = node.IsExpanded
    }
}
```

**Tests** (`internal/tui/navigation_test.go`):
- `TestToggleTreeExpand`: Expand/collapse updates tree state
- `TestExpandStatePersists`: Expand state survives navigation
- `TestSpaceKeyInProjects`: Space toggles expand in projects column
- `TestSpaceKeyInOtherColumns`: Space does nothing in subareas/tasks

---

#### Track B: Area Switching with State (1 hour)
**File**: `internal/tui/app.go`

1. **Area switching functions:**
```go
func (m *Model) switchToPreviousArea() {
    if len(m.areas) == 0 {
        return
    }
    // Save current state
    m.saveCurrentAreaState()
    
    // Switch area with wrap-around
    if m.selectedAreaIndex == 0 {
        m.selectedAreaIndex = len(m.areas) - 1
    } else {
        m.selectedAreaIndex--
    }
    
    // Restore state for new area
    areaID := m.areas[m.selectedAreaIndex].ID
    m.restoreAreaState(areaID)
    
    // Trigger data reload for new area
    m.loadAreaData(areaID)
}

func (m *Model) switchToNextArea() {
    // Similar logic with wrap to first
}

func (m *Model) loadAreaData(areaID string) tea.Cmd {
    // Reuse existing data loading from Task-21
    // This triggers the cascade: subareas → projects → tasks
    return LoadSubareasCmd(m.repo, areaID)
}
```

**Tests** (`internal/tui/navigation_test.go`):
- `TestAreaSwitchWraps`: [ on first area → last area, ] on last → first
- `TestAreaSwitchSavesState`: Selection indices are saved
- `TestAreaSwitchRestoresState`: State is restored when returning
- `TestAreaSwitchReloadsData`: Data loads for new area
- `TestTreeExpandStateAcrossAreas`: Tree state persists per area

---

### Phase 3: Visual Feedback & Styling (1.5 hours)

#### Track A: Selected Item Styling (45 min)
**File**: `internal/tui/app.go` (render methods)

1. **Update renderSubareas():**
```go
func (m *Model) renderSubareas() string {
    if m.isLoadingSubareas {
        return m.spinner.View() + " Loading subareas..."
    }
    
    if len(m.subareas) == 0 {
        return "No subareas\nPress 'a' to add one"
    }
    
    var lines []string
    for i, subarea := range m.subareas {
        line := subarea.Name
        if i == m.selectedSubareaIndex && m.focus == FocusSubareas {
            // Selected + focused: Bold + Inverted
            line = lipgloss.NewStyle().
                Bold(true).
                Reverse(true).
                Render(line)
        }
        lines = append(lines, line)
    }
    return strings.Join(lines, "\n")
}
```

2. **Update renderProjects():**
```go
func (m *Model) renderProjects() string {
    if m.isLoadingProjects {
        return m.spinner.View() + " Loading projects..."
    }
    
    if m.projectTree == nil {
        return "No projects\nPress 'a' to add one"
    }
    
    // Tree.Render() already supports selected node highlighting
    // Pass selectedProjectID for styling
    return tree.Render(m.projectTree, m.selectedProjectID)
}
```

3. **Update renderTasks():** Similar to renderSubareas()

**Tests** (`internal/tui/render_test.go`):
- `TestSelectedSubareaStyling`: Selected subarea has bold+inverted
- `TestSelectedTaskStyling`: Selected task has bold+inverted
- `TestUnselectedItemsNoStyling`: Unselected items have no styling
- `TestTreeViewSelection`: Tree renders with selection highlighting

---

#### Track B: Area Tab Highlighting (30 min)
**File**: `internal/tui/views/tabs.go`

1. **Update TabsView():**
```go
func TabsView(tabs []Tab, selectedIndex int) string {
    var renderedTabs []string
    for i, tab := range tabs {
        style := InactiveTabStyle
        if i == selectedIndex {
            // Active tab: Bold + Inverted
            style = ActiveTabStyle.
                Bold(true).
                Reverse(true)
        }
        renderedTabs = append(renderedTabs, style.Render(tab.Name))
    }
    return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
```

**Tests** (`internal/tui/views/tabs_test.go`):
- `TestActiveTabStyling`: Active tab has bold+inverted
- `TestInactiveTabStyling`: Inactive tabs have normal styling

---

#### Track C: Scroll Behavior (15 min)
**File**: `internal/tui/app.go`

For MVP scope:
- Rely on terminal's natural scroll behavior
- No custom viewport tracking needed for <100 items
- Add comment noting future enhancement opportunity

**Documentation:**
```go
// Scroll Behavior (MVP):
// - Relies on terminal's natural scroll for selected items
// - Future enhancement: track viewport offset and ensure selected item is visible
// - Current approach works well for trees with <100 nodes
```

---

### Phase 4: Comprehensive Testing (2 hours)

#### Test Suite Structure
```
internal/tui/
├── state_test.go          - State persistence tests
├── navigation_test.go     - Navigation and key handler tests
├── render_test.go         - Visual feedback tests
└── integration_test.go    - Full flow integration tests
```

#### Unit Tests (1.5 hours)

**state_test.go:**
- `TestAreaStateInitialization`
- `TestAreaStateSaveRestore`
- `TestAreaStateIsolation`
- `TestTreeExpandStatePersistence`
- `TestMultipleAreasIndependentStates`

**navigation_test.go:**
- `TestNavigateUpWrap`
- `TestNavigateDownWrap`
- `TestNavigateEmptyColumn`
- `TestNavigateTreeUp`
- `TestNavigateTreeDown`
- `TestNavigateTreeWrap`
- `TestTreeSelectionSynchronization`
- `TestToggleTreeExpand`
- `TestSpaceKeyInProjects`
- `TestAreaSwitchWraps`
- `TestAreaSwitchSavesState`
- `TestAreaSwitchRestoresState`

**render_test.go:**
- `TestSelectedSubareaStyling`
- `TestSelectedTaskStyling`
- `TestUnselectedItemsNoStyling`
- `TestTreeViewSelection`
- `TestActiveTabStyling`

**Edge Cases:**
- `TestSingleItemWrap`: j/k on single item stays on same item
- `TestAllCollapsedTree`: Navigation skips collapsed children
- `TestEmptyArea`: Area with no data doesn't crash
- `TestRapidAreaSwitching`: Multiple quick area switches don't corrupt state

---

#### Integration Tests (30 min)

**integration_test.go:**
```go
func TestFullNavigationFlow(t *testing.T) {
    // Setup model with test data
    model := setupTestModel()
    
    // Navigate through all columns
    model.navigateDown(FocusSubareas)
    model.navigateDown(FocusSubareas)
    assert.Equal(t, 2, model.selectedSubareaIndex)
    
    // Switch to projects
    model.focus = FocusProjects
    model.navigateTreeDown()
    assert.NotEmpty(t, model.selectedProjectID)
    
    // Expand/collapse project
    model.toggleTreeExpand()
    node := tree.FindNodeByID(model.projectTree, model.selectedProjectID)
    assert.True(t, node.IsExpanded)
    
    // Switch area and verify state saved
    model.switchToNextArea()
    // ... verify state changed
}

func TestAreaSwitchingStatePersistence(t *testing.T) {
    model := setupTestModel()
    
    // Navigate and expand in area 1
    model.selectedAreaIndex = 0
    model.selectedSubareaIndex = 2
    model.selectedProjectID = "project-1"
    model.toggleTreeExpand()
    
    // Switch to area 2
    model.switchToNextArea()
    assert.Equal(t, 1, model.selectedAreaIndex)
    
    // Navigate in area 2
    model.selectedSubareaIndex = 1
    model.selectedProjectID = "project-5"
    
    // Switch back to area 1
    model.switchToPreviousArea()
    assert.Equal(t, 0, model.selectedAreaIndex)
    
    // Verify area 1 state restored
    assert.Equal(t, 2, model.selectedSubareaIndex)
    assert.Equal(t, "project-1", model.selectedProjectID)
    node := tree.FindNodeByID(model.projectTree, "project-1")
    assert.True(t, node.IsExpanded)
}
```

---

#### Test Coverage Target
- **Target**: >85% coverage for all new code
- **Focus areas**: State management, navigation logic, key handlers
- **Run command**: `go test ./internal/tui -v -cover`

---

### Phase 5: Documentation & Validation (1.5 hours)

#### Code Documentation (45 min)

1. **Package documentation** (`internal/tui/doc.go`):
```go
/*
Package tui implements the terminal user interface for projectdb.

Architecture:
- Model-Update-View pattern (bubbletea)
- State persistence per area
- Clean architecture: TUI depends on domain types and repository interfaces

Navigation:
- Column navigation: h/l/Tab (already implemented in Task-15)
- In-column navigation: j/k (this task)
- Area switching: [/] (this task)
- Tree expand/collapse: Enter/Space (this task)

State Management:
- Per-area state storage (selections, tree expand state)
- Automatic save/restore on area switches
- Independent state for each area
*/
package tui
```

2. **Inline documentation** (in relevant files):
- Document AreaState struct and its purpose
- Comment navigation helper functions
- Document state save/restore flow
- Note clean architecture boundaries

3. **Update existing comments** in:
- `app.go`: Add navigation section to Model struct comment
- `model.go`: Document new state fields

---

#### User Documentation (30 min)

1. **Update README.md:**
```markdown
## TUI Keyboard Shortcuts

### Navigation
- `h/l` or `←/→`: Switch between columns (Subareas | Projects | Tasks)
- `j/k` or `↓/↑`: Navigate within current column
- `Tab`: Cycle through columns
- `[`/`]`: Switch to previous/next area

### Actions
- `Enter` or `Space`: Expand/collapse project (Projects column)
- `a`: Add new item (context-aware)
- `q` or `Ctrl+C`: Exit

### State Persistence
- Selections are automatically saved per area
- Tree expand/collapse state persists when switching areas
```

2. **Create `docs/tui-navigation.md`** (optional, if detailed docs needed):
- Full navigation guide with examples
- State persistence behavior explained
- Troubleshooting section

---

#### Final Validation Checklist (15 min)

**Acceptance Criteria Verification:**
- [ ] #1 j/k navigate within columns with wrap-around
- [ ] #2 Project tree navigation respects expand/collapse
- [ ] #3 Enter/Space toggles expand/collapse
- [ ] #4 [ and ] navigate between areas with wrapping
- [ ] #5 Selected area tab has bold + inverted styling
- [ ] #6 Last selected index restored per area
- [ ] #7 Tree expand/collapse state persisted per area
- [ ] #8 Scroll behavior works (minimal implementation)
- [ ] #9 j/k on empty column is no-op
- [ ] #10 Selected item has bold + inverted styling
- [ ] #11 Unit tests for navigation boundary cases
- [ ] #12 Unit tests for state persistence
- [ ] #13 Integration test for full navigation flow
- [ ] #14 All navigation functions < 20 lines

**Code Quality:**
- [ ] `go test ./internal/tui -v -cover` passes with >85%
- [ ] `go vet ./internal/tui` passes
- [ ] `golangci-lint run` passes (if configured)
- [ ] All functions < 20 lines (verify with `gocyclo` or manual)
- [ ] No magic numbers (all constants named)
- [ ] Clean code principles applied

**Manual Testing:**
```bash
# Run seed data
bash scripts/seed-test-data.sh

# Launch TUI
go run ./cmd/projectdb tui

# Test checklist:
1. [ ] j/k navigation in each column (wrap-around works)
2. [ ] Tree navigation respects collapsed nodes
3. [ ] Enter/Space toggles expand/collapse
4. [ ] [/] switches areas with wrapping
5. [ ] Selection state persists across area switches
6. [ ] Tree expand state persists across area switches
7. [ ] Selected items show bold + inverted styling
8. [ ] Active area tab shows bold + inverted styling
9. [ ] j/k on empty columns does nothing
10. [ ] All visual feedback correct
```

---

## Parallelization Strategy

### Sequential Dependencies (Must Complete in Order)
```
Phase 1 → Phase 2 → Phase 3 → Phase 4 → Phase 5
```

### Parallel Opportunities Within Phases

**Phase 1 (2 hours):**
- Track A (State Storage) and Track B (Selection Helpers) can be done simultaneously
- With 2 developers: ~1.5 hours wall-clock

**Phase 3 (1.5 hours):**
- Track A (Item Styling), Track B (Tab Styling), Track C (Scroll) can be done simultaneously
- With 3 developers: ~45 min wall-clock

**Phase 4 (2 hours):**
- Unit tests can be written in parallel by multiple developers
- Different test files (state_test, navigation_test, render_test)

### Team Execution Timeline

**Single Developer:**
- 9 hours total (sequential flow)

**2 Developers:**
- Phase 1: 1.5h (parallel tracks)
- Phase 2: 2h (sequential)
- Phase 3: 1h (parallel tracks)
- Phase 4: 1.5h (parallel test writing)
- Phase 5: 1h (collaborative validation)
- **Total: ~7 hours wall-clock**

**3 Developers:**
- Phase 1: 1.5h (2 parallel + 1 helper)
- Phase 2: 2h (sequential)
- Phase 3: 45min (3 parallel tracks)
- Phase 4: 1h (parallel test writing)
- Phase 5: 45min (collaborative validation)
- **Total: ~6 hours wall-clock**

---

## Risk Mitigation

### Technical Risks

**Risk 1: Tree navigation complexity**
- **Likelihood**: Low (Task-20 already provides helpers)
- **Mitigation**: Reuse `GetNextVisibleNode`/`GetPrevVisibleNode` directly
- **Contingency**: Add debug logging if tree sync issues occur

**Risk 2: State synchronization bugs**
- **Likelihood**: Medium (state management is tricky)
- **Mitigation**: Comprehensive unit tests, clear state model
- **Contingency**: Add state validation assertions in Update()

**Risk 3: Performance with large trees**
- **Likelihood**: Low (targeting <100 nodes)
- **Mitigation**: Use existing efficient tree helpers
- **Contingency**: Defer optimization to Task-17 (Polish)

**Risk 4: Integration issues with existing code**
- **Likelihood**: Low (clean interfaces from Task-15/20/21)
- **Mitigation**: Incremental testing, frequent builds
- **Contingency**: Revert to last known good state, debug incrementally

### Schedule Risks

**Risk 1: Underestimating testing effort**
- **Likelihood**: Medium (14 ACs require thorough testing)
- **Mitigation**: Start testing early (TDD approach), parallel test writing
- **Contingency**: Prioritize critical path tests, defer edge cases to Task-17

**Risk 2: Scope creep (adding features)**
- **Likelihood**: Medium (easy to add "nice-to-haves")
- **Mitigation**: Strict adherence to ACs, defer enhancements to Task-17
- **Contingency**: Document enhancements as future tasks, don't implement now

---

## Success Criteria

### Definition of Done

**Code Complete:**
- [ ] All 14 acceptance criteria implemented and verified
- [ ] All functions < 20 lines (SRP followed)
- [ ] Test coverage >85%
- [ ] Code passes `go vet` and lint checks
- [ ] No regressions in existing TUI features

**Quality Gates:**
- [ ] Unit tests pass: `go test ./internal/tui -v -cover`
- [ ] Integration tests pass
- [ ] Manual testing checklist complete
- [ ] Code review approved (if applicable)
- [ ] Documentation updated (code comments + README)

**User Experience:**
- [ ] Navigation feels smooth and intuitive
- [ ] State persistence works transparently
- [ ] Visual feedback is clear and consistent
- [ ] No confusing behaviors or edge case failures

---

## Next Steps After Completion

### Immediate Follow-ups
- Task-19: Quick-Add Modal (can start in parallel with Task-18)
- Task-17: Help, Errors & Polish (depends on Task-18 completion)

### Future Enhancements (Out of Scope for Task-18)
- Custom viewport tracking for large trees
- Keyboard shortcut customization
- Search/filter within columns
- Multi-select functionality

---

## Appendix: File Changes Summary

### New Files
- `internal/tui/state_test.go` - State persistence tests
- `internal/tui/navigation_test.go` - Navigation and key handler tests
- `internal/tui/render_test.go` - Visual feedback tests
- `internal/tui/integration_test.go` - Full flow tests

### Modified Files
- `internal/tui/model.go` - Add AreaState struct and state management methods
- `internal/tui/app.go` - Add navigation helpers, key handlers, render updates
- `internal/tui/views/tabs.go` - Update tab styling for active tab
- `README.md` - Add keyboard shortcuts documentation

### Test Coverage Impact
- Current: ~85% (Task-21 completion)
- Target: >85% (maintain or improve)
- New code: >90% coverage required

---

## Conclusion

Task-18 is well-scoped as a single task because:
1. **Foundation is ready**: Tree helpers and data loading already complete
2. **Clear focus**: Add state persistence and wire up navigation
3. **Reasonable size**: 9 hours, 14 clear ACs
4. **Manageable complexity**: Integration with existing clean architecture

The implementation follows a clear progression:
1. **Build state layer** (foundation)
2. **Wire navigation** (integration)
3. **Polish visuals** (user experience)
4. **Test thoroughly** (quality assurance)
5. **Validate and document** (completion)

With the detailed plan above, Task-18 can be implemented efficiently with high quality and comprehensive testing.
