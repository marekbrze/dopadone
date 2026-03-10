# Detailed Implementation Plan: Task-68.3 TUI Delete Integration

**Task**: TASK-68.3 - TUI Delete Integration  
**Status**: To Do → In Progress  
**Dependencies**: TASK-68.1 (Done), TASK-68.2 (Done)  
**Estimated Effort**: 4-6 hours  

---

## Overview

This document provides a detailed implementation plan for **Task-68.3: TUI Delete Integration**, which wires up delete functionality in the TUI with the confirmation modal created in Task-68.1 and the cascade soft delete service from Task-68.2.

## Prerequisites Verification

✅ **Task-68.1 Complete**: Confirmation modal component ready  
✅ **Task-68.2 Complete**: Cascade soft delete service ready  
✅ **Dependencies Met**: Can proceed with implementation

---

## Architecture Context

### Current TUI Architecture
- **Model** (`internal/tui/model.go`): Contains application state including focus, modals, and data
- **Messages** (`internal/tui/messages.go`): Defines tea.Msg types for async operations
- **Commands** (`internal/tui/commands.go`): Implements tea.Cmd functions for async operations
- **Handlers** (`internal/tui/handlers.go`): Processes messages and updates state
- **Constants** (`internal/tui/constants.go`): Defines keybindings and constants
- **Renderer** (`internal/tui/renderer.go`): Renders the UI including footer

### Integration Points
1. **Model struct** needs new fields for confirmation modal state
2. **Message types** need delete-specific messages
3. **Commands** need delete operations for each entity type
4. **Key handlers** need to handle 'd' key in all focus columns
5. **Message handlers** need to process confirm/cancel/success/error messages
6. **Footer** needs to display 'd: delete' shortcut

---

## Implementation Phases

### Phase 1: Update Model Struct (15 min)
**Priority**: CRITICAL - Foundation for all other work

#### File: `internal/tui/model.go`

**Changes Required**:
```go
type Model struct {
    // ... existing fields ...
    
    // Delete confirmation modal
    confirmModal       *confirmmodal.Modal
    isConfirmModalOpen bool
}
```

**Rationale**:
- Follows existing modal pattern (see `isAreaModalOpen`, `isHelpOpen`)
- Pointer allows nil check (no modal when closed)
- Boolean flag for quick state checks

**Implementation Steps**:
1. Add `confirmModal *confirmmodal.Modal` field
2. Add `isConfirmModalOpen bool` field
3. Ensure fields are initialized to nil/false in constructor

**Validation**:
- [ ] Fields added to Model struct
- [ ] Zero values handled correctly (nil/false)
- [ ] Code compiles without errors

---

### Phase 2: Add Delete Messages (20 min)
**Priority**: CRITICAL - Defines message contracts

#### File: `internal/tui/messages.go`

**New Message Types**:
```go
// Delete operation messages
type DeleteConfirmedMsg struct {
    EntityType confirmmodal.EntityType
    EntityID   string
    EntityName string
}

type DeleteSuccessMsg struct {
    EntityType confirmmodal.EntityType
    EntityName string
}

type DeleteErrorMsg struct {
    Err        error
    EntityType confirmmodal.EntityType
    EntityName string
}
```

**Rationale**:
- `DeleteConfirmedMsg`: Dispatched when user presses 'y' in confirmation modal
- `DeleteSuccessMsg`: Dispatched when delete operation succeeds
- `DeleteErrorMsg`: Dispatched when delete operation fails
- Includes entity name for user-friendly toasts

**Implementation Steps**:
1. Define `DeleteConfirmedMsg` struct with entity type, ID, and name
2. Define `DeleteSuccessMsg` struct with entity type and name
3. Define `DeleteErrorMsg` struct with error, entity type, and name
4. Ensure proper ordering (group delete messages together)

**Validation**:
- [ ] Message types defined
- [ ] Fields match expected usage
- [ ] Follow existing naming conventions
- [ ] Code compiles without errors

---

### Phase 3: Add Delete Commands (30 min)
**Priority**: CRITICAL - Implements async delete operations

#### File: `internal/tui/commands.go`

**New Command Functions**:

**1. DeleteSubareaCmd**:
```go
func DeleteSubareaCmd(subareaSvc service.SubareaServiceInterface, id string, name string) tea.Cmd {
    return func() tea.Msg {
        err := subareaSvc.SoftDelete(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{
                Err:        err,
                EntityType: confirmmodal.EntityTypeSubarea,
                EntityName: name,
            }
        }
        return DeleteSuccessMsg{
            EntityType: confirmmodal.EntityTypeSubarea,
            EntityName: name,
        }
    }
}
```

**2. DeleteProjectCmd**:
```go
func DeleteProjectCmd(projectSvc service.ProjectServiceInterface, id string, name string) tea.Cmd {
    return func() tea.Msg {
        err := projectSvc.SoftDeleteWithCascade(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{
                Err:        err,
                EntityType: confirmmodal.EntityTypeProject,
                EntityName: name,
            }
        }
        return DeleteSuccessMsg{
            EntityType: confirmmodal.EntityTypeProject,
            EntityName: name,
        }
    }
}
```

**3. DeleteTaskCmd**:
```go
func DeleteTaskCmd(taskSvc service.TaskServiceInterface, id string, name string) tea.Cmd {
    return func() tea.Msg {
        err := taskSvc.SoftDelete(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{
                Err:        err,
                EntityType: confirmmodal.EntityTypeTask,
                EntityName: name,
            }
        }
        return DeleteSuccessMsg{
            EntityType: confirmmodal.EntityTypeTask,
            EntityName: name,
        }
    }
}
```

**Rationale**:
- Each command calls the appropriate service method
- Uses context.Background() for proper context propagation
- Returns typed success/error messages for proper handling
- Includes entity name for toast messages
- Follows existing command patterns (see `LoadSubareasCmd`, `CreateSubareaCmd`)

**Implementation Steps**:
1. Import `confirmmodal` package at2. Implement `DeleteSubareaCmd` with error handling
3. Implement `DeleteProjectCmd` with cascade delete call
4. Implement `DeleteTaskCmd` with error handling
5. Add proper godoc comments

**Validation**:
- [ ] All three command functions implemented
- [ ] Proper error handling and all functions
- [ ] Context propagation correct
- [ ] Code compiles without errors
- [ ] Follows existing command patterns

---

### Phase 4: Add Delete Key Binding (15 min)
**Priority**: HIGH - User interaction entry point

#### File: `internal/tui/constants.go`

**Changes Required**:
```go
type keymap struct {
    // ... existing bindings ...
    delete key.Binding
}
```

**Update defaultKeyMap**:
```go
var defaultKeyMap = keymap{
    // ... existing keys ...
    delete: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "delete"),
    ),
}
```

**Rationale**:
- Consistent with existing keybinding pattern
- Single key 'd' for simplicity
- Help text follows convention: "key", "description"

**Implementation Steps**:
1. Add `delete` field to keymap struct
2. Add binding to defaultKeyMap initialization
3. Ensure binding is accessible via m.keys.delete

**Validation**:
- [ ] Key binding added
- [ ] Help text is correct
- [ ] Accessible via Model.keys.delete
- [ ] Code compiles without errors

---

### Phase 5: Implement Delete Key Handler (45 min)
**Priority**: CRITICAL - Core functionality

#### File: `internal/tui/handlers.go` or separate file

**New Handler Function**:
```go
func (m *Model) handleDeleteKey() (tea.Model, tea.Cmd) {
    // If confirmation modal is already open, let it modal handle it key
    if m.isConfirmModalOpen && m.confirmModal != nil {
        newModal, cmd := m.confirmModal.Update(msg)
        m.confirmModal = newModal
        return m, cmd
    }

    // Handle based on focus column
    switch m.focus {
    case FocusSubareas:
        return m.handleDeleteSubarea()
    case FocusProjects:
        return m.handleDeleteProject()
    case FocusTasks:
        return m.handleDeleteTask()
    }
    return m, nil
}

func (m *Model) handleDeleteSubarea() (tea.Model, tea.Cmd) {
    // Check for empty column
    if len(m.subareas) == 0 || m.selectedSubareaIndex >= len(m.subareas) {
        return m, nil // No-op
    }

    subarea := m.subareas[m.selectedSubareaIndex]
    m.confirmModal = confirmmodal.New(
        subarea.Name,
        confirmmodal.EntityTypeSubarea,
        subarea.ID,
    )
    m.isConfirmModalOpen = true
    return m, nil
}

func (m *Model) handleDeleteProject() (tea.Model, tea.Cmd) {
    // Check for empty column or no selected project
    if m.projectTree == nil || m.selectedProjectID == "" {
        return m, nil // No-op
    }

    node := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
    if node == nil {
        return m, nil
    }

    m.confirmModal = confirmmodal.New(
        node.Name,
        confirmmodal.EntityTypeProject,
        node.ID,
    )
    m.isConfirmModalOpen = true
    return m, nil
}

func (m *Model) handleDeleteTask() (tea.Model, tea.Cmd) {
    // Check for empty column or no selected task
    if m.groupedTasks == nil || len(m.groupedTasks.AllLines) == 0 {
        return m, nil // No-op
    }

    if m.selectedTaskIndex < 0 || m.selectedTaskIndex >= len(m.groupedTasks.AllLines) {
        return m, nil
    }

    // Skip group headers
    if m.isLineGroupHeader(m.selectedTaskIndex) {
        return m, nil
    }

    task := m.getTaskAtLine(m.selectedTaskIndex)
    if task == nil {
        return m, nil
    }

    m.confirmModal = confirmmodal.New(
        task.Title,
        confirmmodal.EntityTypeTask,
        task.ID,
    )
    m.isConfirmModalOpen = true
    return m, nil
}
```

**Key Design Decisions**:
- **Modal first**: If modal open, delegate to modal's Update
- **Empty check**: Return early if no items or column (no-op)
- **Group header skip**: Can't delete group headers (only tasks)
- **Tree navigation**: Use `tree.FindNodeByID` for project tree
- **Task navigation**: Use `getTaskAtLine` for grouped tasks

**Implementation Steps**:
1. Import `confirmmodal` and `tree` packages
2. Implement `handleDeleteKey` as router function
3. Implement `handleDeleteSubarea` with validation
4. Implement `handleDeleteProject` with tree navigation
5. Implement `handleDeleteTask` with grouped tasks
6. Add helper methods if needed (e.g., `getTaskAtLine`)

**Validation**:
- [ ] All four handler functions implemented
- [ ] Empty column handling correct
- [ ] Group header skip correct
- [ ] Modal delegation working
- [ ] Code compiles without errors

---

### Phase 6: Add Message Handlers (30 min)
**Priority**: CRITICAL - Processes modal and delete messages

#### File: `internal/tui/handlers.go` or separate file

**New Handler Functions**:

**1. Handle Confirmation Modal Messages**:
```go
func (m *Model) handleConfirmModalMessages(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case confirmmodal.ConfirmMsg:
        return m, func() tea.Msg {
            return DeleteConfirmedMsg{
                EntityType: msg.EntityType,
                EntityID:   msg.EntityID,
                EntityName: msg.EntityName,
            }
        }
    
    case confirmmodal.CancelMsg:
        m.isConfirmModalOpen = false
        m.confirmModal = nil
        return m, nil
    }
    return m, nil
}
```

**2. Handle Delete Confirmed**:
```go
func (m *Model) handleDeleteConfirmed(msg DeleteConfirmedMsg) (tea.Model, tea.Cmd) {
    // Close modal
    m.isConfirmModalOpen = false
    m.confirmModal = nil
    
    // Dispatch appropriate delete command
    switch msg.EntityType {
    case confirmmodal.EntityTypeSubarea:
        return m, DeleteSubareaCmd(m.subareaSvc, msg.EntityID, msg.EntityName)
    case confirmmodal.EntityTypeProject:
        return m, DeleteProjectCmd(m.projectSvc, msg.EntityID, msg.EntityName)
    case confirmmodal.EntityTypeTask:
        return m, DeleteTaskCmd(m.taskSvc, msg.EntityID, msg.EntityName)
    default:
        return m, nil
    }
}
```

**3. Handle Delete Success**:
```go
func (m *Model) handleDeleteSuccess(msg DeleteSuccessMsg) (tea.Model, tea.Cmd) {
    // Show success toast
    entityName := string(msg.EntityType)
    m.addToast(toast.NewSuccess(fmt.Sprintf("%s '%s' deleted successfully", entityName, msg.EntityName)))
    
    // Refresh appropriate column
    switch msg.EntityType {
    case confirmmodal.EntityTypeSubarea:
        if len(m.areas) > 0 {
            areaID := m.areas[m.selectedAreaIndex].ID
            return m, LoadSubareasCmd(m.subareaSvc, areaID)
        }
    case confirmmodal.EntityTypeProject:
        if len(m.subareas) > 0 {
            subareaID := m.subareas[m.selectedSubareaIndex].ID
            return m, LoadProjectsCmd(m.projectSvc, &subareaID)
        }
    case confirmmodal.EntityTypeTask:
        if m.selectedProjectID != "" {
            return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
        }
    }
    return m, nil
}
```

**4. Handle Delete Error**:
```go
func (m *Model) handleDeleteError(msg DeleteErrorMsg) (tea.Model, tea.Cmd) {
    // Show error toast with entity name
    entityName := string(msg.EntityType)
    m.addToast(toast.NewError(fmt.Sprintf("Failed to delete %s '%s': %v", entityName, msg.EntityName, msg.Err)))
    
    // No state change needed, modal already closed
    return m, nil
}
```

**Key Design Decisions**:
- **Immediate modal close**: Close modal before async operation
- **Entity name in toasts**: Makes messages user-friendly
- **Column refresh**: Reload appropriate column after successful delete
- **No refresh on error**: Keep current state on error

**Implementation Steps**:
1. Implement `handleConfirmModalMessages` for modal delegation
2. Implement `handleDeleteConfirmed` for command dispatch
3. Implement `handleDeleteSuccess` with toast and refresh
4. Implement `handleDeleteError` with toast
5. Update Update() method to call these handlers

**Validation**:
- [ ] All four message handlers implemented
- [ ] Modal delegation correct
- [ ] Toast messages include entity name
- [ ] Column refresh logic correct
- [ ] Code compiles without errors

---

### Phase 7: Update View Rendering (15 min)
**Priority**: MEDIUM - UI integration

#### File: `internal/tui/renderer.go` or separate file

**Changes Required**:

**1. Render Confirmation Modal Overlay**:
```go
func (m *Model) View() string {
    // ... existing view rendering ...
    
    // Render confirmation modal on top if open
    if m.isConfirmModalOpen && m.confirmModal != nil {
        modalView := m.confirmModal.View()
        // Overlay modal centered on screen
        return lipgloss.Place(
            m.width, m.height,
            lipgloss.WithWhitespaceChars(" "),
            lipgloss.WithForegroundColor(lipgloss.Color("#000000")),
        ).Render(modalView, lipgloss.Center, lipgloss.Center)
    }
    
    return view
}
```

**2. Update Footer**:
```go
func (m *Model) RenderFooter() string {
    shortcuts := []string{
        "h/l: columns",
        "j/k: nav",
        "a: add",
        "x: toggle",
        "d: delete", // NEW
        "?: help",
        "q: quit",
    }
    // ... rest of footer rendering ...
}
```

**Key Design Decisions**:
- **Modal overlay**: Rendered on top of main view, centered
- **Footer order**: 'd' comes before 'x' (toggle) as per decision
- **Consistent styling**: Use existing lipgloss patterns

**Implementation Steps**:
1. Update View() to render confirmation modal overlay
2. Update RenderFooter() to include 'd: delete' shortcut
3. Ensure modal styling matches existing modal patterns
4. Test rendering in different window sizes

**Validation**:
- [ ] Modal overlay renders correctly
- [ ] Modal is centered on screen
- [ ] Footer includes 'd: delete' shortcut
- [ ] Footer order is correct (d before x)
- [ ] Code compiles without errors

---

### Phase 8: Integration in Update Method (20 min)
**Priority**: CRITICAL - Wires everything together

#### File: `internal/tui/app.go` or `internal/tui/update.go`

**Changes Required**:

**1. Add Delete Key Handler to Update()**:
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd
    
    // Handle confirmation modal messages FIRST (before other handlers)
    if m.isConfirmModalOpen {
        if modalMsg, ok := m.handleConfirmModalMessages(msg); modalMsg.Ok {
            return modalMsg, modalCmd
        }
    }
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch {
        case key.Matches(msg, m.keys.delete):
            return m.handleDeleteKey()
        // ... existing key handlers ...
        }
    
    case confirmmodal.ConfirmMsg:
        return m.handleDeleteConfirmed(msg)
    
    case confirmmodal.CancelMsg:
        return m.handleConfirmModalMessages(msg)
    
    case DeleteConfirmedMsg:
        return m.handleDeleteConfirmed(msg)
    
    case DeleteSuccessMsg:
        return m.handleDeleteSuccess(msg)
    
    case DeleteErrorMsg:
        return m.handleDeleteError(msg)
    
    // ... existing message handlers ...
    }
    
    return m, tea.Batch(cmds...)
}
```

**Key Design Decisions**:
- **Modal priority**: If modal open, handle modal messages first
- **Message ordering**: Confirmation modal messages before delete messages
- **Existing handlers**: Don't break existing functionality

**Implementation Steps**:
1. Add confirmation modal message handling at top of Update()
2. Add 'd' key to key switch
3. Add message cases for delete messages
4. Ensure existing handlers still work
5. Test with both keyboard and mouse events

**Validation**:
- [ ] Delete key handler integrated
- [ ] Modal message handling correct
- [ ] Message routing correct
- [ ] Existing functionality still works
- [ ] Code compiles without errors

---

### Phase 9: Testing Strategy (60 min)
**Priority**: HIGH - Quality assurance

#### File: `internal/tui/app_delete_test.go` (NEW FILE)

**Test Categories**:

**1. Key Handler Tests**:
- Test 'd' key opens confirmation modal for subarea
- Test 'd' key opens confirmation modal for project
- Test 'd' key opens confirmation modal for task
- Test 'd' key does nothing on empty columns
- Test 'd' key skips group headers in tasks column

- Test modal shows correct item name

**2. Confirmation Modal Tests**:
- Test 'y' key triggers delete
- Test 'n' key cancels delete
- Test Escape key cancels delete
- Test modal closes on modal state cleared

**3. Delete Command Tests**:
- Test DeleteSubareaCmd success
- Test DeleteSubareaCmd error
- Test DeleteProjectCmd success (- Test DeleteProjectCmd error
- Test DeleteTaskCmd success
- Test DeleteTaskCmd error

**4. Message Handler Tests**:
- Test DeleteConfirmedMsg dispatches correct command
- Test DeleteSuccessMsg shows toast and refreshes
- Test DeleteErrorMsg shows toast
- Test handleConfirmModalMessages with ConfirmMsg
- Test handleConfirmModalMessages with CancelMsg

**5. Integration Tests**:
- Test full delete flow for subarea
- Test full delete flow for project
- Test full delete flow for task
- Test error handling
- Test column refresh after delete
- Test modal state management

**Test Pattern**:
```go
func TestDeleteKey(t *testing.T) {
    tests := []struct {
        name          string
        focus         FocusColumn
        setup         func(*Model)
        initialItems int
        pressKey       string
        expectModal    bool
        expectName     string
    }{
        {
            name: "opens modal for subarea",
            focus: FocusSubareas,
            setup: func(m *Model) {
                m.subareas = []domain.Subarea{
                    {ID: "sub-1", Name: "Test Subarea"},
                    {ID: "sub-2", Name: "Another Subarea"},
                }
                m.selectedSubareaIndex = 0
            initialItems: 22,
            pressKey: "d",
            expectModal: true,
            expectName: "Test Subarea",
        },
        // ... more tests
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup model
            m := tt.setup(t)
            
            // Press key
            m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyMsg, Runes: tt.pressKey, chars: []rune{tt.pressKey})
            if !tt.expectModal {
                t.Fatalf("expected modal to open, got none")
            }
            
            // Check modal properties
            if !m.isConfirmModalOpen || m.confirmModal == nil {
                t.Fatal("expected modal to be open")
            }
            
            if m.confirmModal.(type) == nil {
                t.Fatal("expected modal to be set")
            }
            
            // Verify modal has correct item name
            // (Access private fields via reflection or add getter methods)
            if m.confirmModal.ItemName() != tt.expectName {
                t.Errorf("expected item name %q, got %q", tt.expectName, m.confirmModal.ItemName())
            }
        })
    }
}
```

**Implementation Steps**:
1. Create test file `internal/tui/app_delete_test.go`
2. Write key handler tests for all focus columns
3. Write confirmation modal tests
4. Write delete command tests with mocks
5. Write message handler tests
6. Write integration tests
7. Add helper functions for common setup/assertions

**Test Coverage Goals**:
- **Key handlers**: 100% coverage
- **Confirmation modal**: 100% coverage
- **Delete commands**: 100% coverage
- **Message handlers**: 100% coverage
- **Integration tests**: All major flows paths covered

**Validation**:
- [ ] Test file created
- [ ] All test categories implemented
- [ ] Helper functions for common setup
- [ ] Mock services for testing
- [ ] 100% coverage for new code
- [ ] Tests pass with `go test -v`
- [ ] Tests pass with `go test -race`

---

### Phase 10: Documentation Updates (15 min)
**Priority**: LOW - User documentation

#### File: `docs/TUI.md`

**Changes Required**:

**1. Update Keyboard Shortcuts Section**:
```markdown
### Actions

| Key | Action | Description |
|-----|--------|-------------|
| `h/l` | Columns | Switch between Subareas, Projects, and Tasks columns |
| `j/k` | Navigate | Move up/down within the current column |
| `a` | Add | Add new item to current column |
| `x` | Toggle | Toggle task completion or project expansion |
| `d` | Delete | Delete focused item (shows confirmation) | ← NEW

| `?` | Help | Show keyboard shortcuts help |
| `q` | Quit | Exit the application |
```

**2. Update Quick Reference Section**:
Add note about delete confirmation modal and cascade delete behavior.

**Implementation Steps**:
1. Add 'd: delete' to keyboard shortcuts table
2. Add description of delete behavior
3. Document confirmation modal
4. Document cascade delete for projects
5. Update any examples if needed

**Validation**:
- [ ] Keyboard shortcuts table updated
- [ ] Delete behavior documented
- [ ] Confirmation modal mentioned
- [ ] Documentation is clear and accurate

---

## Task Split Summary

This task is **NOT split** - it's implemented as a single cohesive unit with 10 phases that build on each other.

**Rationale for NOT splitting**:
- Task is well-scoped and focused (TUI integration only)
- Dependencies are already complete (68.1 and 68.2)
- All phases are sequential but build on each other
- Natural testing boundaries (unit tests vs integration tests)
- Estimated effort is 4-6 hours (manageable in one session)

---

## Execution Timeline

### Sequential Phases (Must Complete in Order)
```
Phase 1 (Model Struct)           →  15 min
Phase 2 (Messages)             →  20 min
Phase 3 (Commands)              →  30 min
Phase 4 (Key Binding)            →  15 min
Phase 5 (Key Handler)            →  45 min
Phase 6 (Message Handlers)        →  30 min
Phase 7 (View Rendering)           →  15 min
Phase 8 (Update Integration)      →  20 min
Phase 9 (Testing)                 →  60 min
Phase 10 (Documentation)           →  15 min
```

**Total Estimated Time**: 4-6 hours (265 minutes)

**Buffer**: 30 minutes (for code review, fixes, and testing)

**Total with Buffer**: 5-7 hours

---

## Parallel Work Opportunities

While phases are mostly sequential, some work can be done in parallel:

- **Testing + Documentation**: Phase 9 and 10 can run in parallel (different focus areas)
- **Unit tests**: Can write tests for different handlers in parallel during Phase 9
- **Code quality checks**: Can run linters while writing tests

---

## Dependencies and Integration Points

### Internal Dependencies
- **confirmmodal package**: Modal component from Task-68.1
- **tree package**: Tree navigation for project selection
- **service package**: Soft delete methods from Task-68.2
- **toast package**: Toast notifications
- **lipgloss**: Styling for modal overlay

### Integration Points
1. **Model struct**: Add confirmation modal state
2. **Messages**: Define delete message types
3. **Commands**: Implement delete operations
4. **Handlers**: Process delete key and messages
5. **View**: Render modal overlay,6. **Footer**: Show delete shortcut
7. **Update**: Route messages to handlers

---

## Risks and Mitigations

### Risk 1: Modal State Conflicts
**Impact**: Multiple modals open at same time  
**Mitigation**:
- Use single `isConfirmModalOpen` boolean flag
- Check modal state before opening new modal
- Ensure only one modal can be active at a time
- Follow existing modal pattern (area modal, help modal)

**Validation**: Test opening confirmation modal when other modals are open

### Risk 2: Column Refresh Timing
**Impact**: UI might flash or show stale data briefly  
**Mitigation**:
- Delete operations are async (tea.Cmd pattern)
- Show loading state during refresh (implicit via column reload)
- Refresh only affected column, not entire UI
- Use existing load commands (LoadSubareasCmd, etc.)

**Validation**: Test column refresh after delete

### Risk 3: Toast Message Clarity
**Impact**: Users might not see feedback on delete operation  
**Mitigation**:
- Include entity name in toast messages
- Use clear success/error distinction
- Toast auto-dismisses after timeout (existing behavior)
- Show entity type for context (subarea, project, task)

**Validation**: Test toast messages appear and content is correct

### Risk 4: Empty Column Edge Case
**Impact**: User might try to delete when no items exist  
**Mitigation**:
- Check for empty columns before opening modal
- Return early (no-op) if column is empty
- Check for valid selection index
- Handle nil tree nodes gracefully

**Validation**: Test with empty columns in all three focus states

### Risk 5: Group Header Selection
**Impact**: User might try to delete group header instead of task  
**Mitigation**:
- Check if selected line is group header
- Return early (no-op) if group header
- Only allow deletion of actual tasks

**Validation**: Test with group headers in tasks column

---

## Acceptance Criteria Mapping

| AC # | Implementation Phase | Test Coverage |
|------|----------------------|---------------|
| #1: Add confirmModal state to Model struct | Phase 1 | Unit test |
| #2: Add delete key binding 'd' to all focus columns | Phase 4, Unit test |
| #3: Open confirmation modal with correct item name when 'd' pressed | Phase 5 | Unit test |
| #4: Execute appropriate delete command on 'y' confirmation | Phase 6 | Unit test |
| #5: Show success toast and refresh column after delete | Phase 6 | Integration test |
| #6: Show error toast on delete failure | Phase 6 | Unit test |
| #7: Handle 'n' and Escape to cancel delete | Phase 6 | Unit test |
| #8: No-op when pressing 'd' on empty columns | Phase 5 | Unit test |
| #9: Update footer to show 'd: delete' shortcut | Phase 7 | Visual inspection |
| #10: Write integration tests for delete flow | Phase 9 | Integration test |

---

## Files to Modify

### New Files
1. `internal/tui/app_delete_test.go` - Integration tests for delete functionality (~200 lines)

### Modified Files
1. `internal/tui/model.go` - Add confirmation modal state (2 lines)
2. `internal/tui/messages.go` - Add delete message types (15 lines)
3. `internal/tui/commands.go` - Add delete commands (60 lines)
4. `internal/tui/constants.go` - Add delete key binding (5 lines)
5. `internal/tui/handlers.go` or new file - Add delete handlers (150 lines)
6. `internal/tui/renderer.go` - Update footer and modal overlay (15 lines)
7. `internal/tui/app.go` or `update.go` - Integrate delete handlers (30 lines)
8. `docs/TUI.md` - Update keyboard shortcuts (10 lines)

**Total Lines Changed**: ~500 lines (excluding tests)

---

## Success Criteria

### Functional Requirements
- [ ] Pressing 'd' opens confirmation modal in all columns
- [ ] Confirmation modal shows correct item name
- [ ] Pressing 'y' triggers delete operation
- [ ] Pressing 'n' or Escape cancels delete
- [ ] Success toast shows entity name
- [ ] Error toast shows entity name and error
- [ ] Column refreshes after successful delete
- [ ] No-op on empty columns
- [ ] Footer shows 'd: delete' shortcut

### Quality Requirements
- [ ] 100% test coverage for new code
- [ ] No race conditions
- [ ] No compiler warnings
- [ ] Follows project coding standards
- [ ] Consistent with existing modal patterns
- [ ] Proper error handling

### Integration Requirements
- [ ] Works with existing modal system
- [ ] Doesn't break existing functionality
- [ ] Toast notifications display correctly
- [ ] Column refresh works properly
- [ ] Modal overlay renders correctly

---

## Post-Implementation Verification

### Manual Testing Checklist
- [ ] Start TUI: `make run`
- [ ] Select subarea, press 'd', verify modal shows subarea name
- [ ] Press 'y', verify toast and column refresh
- [ ] Select project, press 'd', verify modal shows project name
- [ ] Press 'n', verify modal closes
- [ ] Select task, press 'd', verify modal shows task title
- [ ] Press Escape, verify modal closes
- [ ] Press 'd' on empty column, verify no-op
- [ ] Verify footer shows 'd: delete'
- [ ] Verify cascade delete works for project with children
- [ ] Verify error toast shows on database error

### Automated Testing
```bash
# Run all tests
go test ./internal/tui/... -v

# Run with race detection
go test -race ./internal/tui/...

# Run with coverage
go test -cover ./internal/tui/...
go tool cover -html=coverage.out

# Run linting
golangci-lint run ./internal/tui/...
```

---

## Commands Summary

```bash
# Build project
go build ./...

# Run tests
go test ./internal/tui/... -v

# Run with race detection
go test -race ./internal/tui/...

# Run with coverage
go test -cover -coverprofile=coverage.out ./internal/tui/...
go tool cover -html=coverage.out

# Run linting
golangci-lint run ./internal/tui/...

# Format code
gofmt -w ./internal/tui/
goimports -w ./internal/tui/

# Manual testing
make run
```

---

## Next Steps

1. ✅ **Verify Prerequisites**: Task-68.1 and 68.2 complete
2. **Start Implementation**: Begin with Phase 1 (Model Struct)
3. **Follow Phases**: Complete phases 2-8 sequentially
4. **Write Tests**: Phase 9 in parallel with phases 6-8
5. **Update Docs**: Phase 10 after implementation complete
6. **Run Quality Checks**: Linting, formatting, race detection
7. **Manual Testing**: Verify all functionality manually
8. **Create Pull Request**: Reference Task-68.3
9. **Update Task Status**: Mark as Done after review

---

## Notes

- **No Task Splitting**: This task is well-scoped and focused on TUI integration
- **Sequential Phases**: Phases must complete in order as they build on each other
- **Testing Priority**: Comprehensive testing is critical for quality
- **Documentation**: Update keyboard shortcuts in TUI.md
- **Dependencies Met**: Task-68.1 and 68.2 provide foundation
- **Estimated Effort**: 4-6 hours (manageable in one session)

**Key Decision**: Keep this as a single task because the phases are sequential and build naturally on each other, making it easier to track progress and maintain consistency.
