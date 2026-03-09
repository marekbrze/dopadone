---
id: TASK-68
title: 'Add option to delete subareas, projects and tasks in tui'
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-09 18:15'
updated_date: '2026-03-09 19:06'
labels: []
dependencies: []
references:
  - '# Task-68.1: Confirmation Modal Component'
  - '# Task-68.2: Cascade Soft Delete Service'
  - '# Task-68.3: TUI Delete Integration'
  - '# Task-68.4: Delete Documentation & Testing'
  - '# Milestone: m-2 Deleting items'
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add delete functionality for the focused element (subarea, project, or task) in the TUI. When user presses 'd' on a selected item, show a confirmation dialog. On confirmation, perform soft delete. For projects with subprojects, cascade delete children. Update footer to show the new 'd: delete' shortcut.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Pressing 'd' when subarea is focused opens confirmation dialog with subarea name
- [ ] #2 Pressing 'd' when project is focused opens confirmation dialog with project name
- [ ] #3 Pressing 'd' when task is focused opens confirmation dialog with task title
- [ ] #4 Confirmation dialog shows item name and 'y/n' options (y=confirm, n/esc=cancel)
- [ ] #5 Pressing 'y' in confirmation dialog triggers soft delete via service layer
- [ ] #6 Pressing 'n' or Escape in confirmation dialog cancels and returns to normal view
- [ ] #7 Successful deletion shows success toast and refreshes the column data
- [ ] #8 Failed deletion shows error toast and does not change UI state
- [ ] #9 Deleting a project with subprojects cascades delete to all children (subprojects and their tasks)
- [ ] #10 Pressing 'd' on empty column does nothing (no-op)
- [ ] #11 Footer is updated to include 'd: delete' in help text
- [ ] #12 Delete confirmation modal follows existing modal patterns (lipgloss styling, centered overlay)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
IMPLEMENTATION PLAN FOR TASK-68: Add Delete Functionality in TUI

## TASK COMPLEXITY ASSESSMENT

**Medium-High Complexity** - This task involves multiple concerns:
- New UI component (confirmation modal)
- Service layer enhancement (cascade soft delete)
- Keyboard event handling
- State management and refresh
- Comprehensive testing
- Documentation updates

**Decision: Split into 4 subtasks for better manageability**

---

## PHASE 1: FOUNDATION - Confirmation Modal Component (Task-68A)

**Goal:** Create reusable confirmation modal following existing patterns

### Files to Create/Modify:
1. `internal/tui/confirmmodal/confirm_modal.go` - New file
2. `internal/tui/confirmmodal/styles.go` - New file
3. `internal/tui/confirmmodal/confirm_modal_test.go` - New file

### Implementation Details:

**confirm_modal.go:**
```go
package confirmmodal

type EntityType string

const (
    EntityTypeSubarea EntityType = "subarea"
    EntityTypeProject EntityType = "project"
    EntityTypeTask    EntityType = "task"
)

type Modal struct {
    message     string
    itemName    string
    entityType  EntityType
    entityID    string
    width       int
    height      int
}

type ConfirmMsg struct {
    EntityType EntityType
    EntityID   string
}

type CancelMsg struct{}

func New(itemName string, entityType EntityType, entityID string) *Modal
func (m *Modal) Init() tea.Cmd
func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd)
func (m *Modal) View() string
```

**Key Features:**
- Centered overlay modal (follow modal/modal.go pattern)
- Lipgloss styling with theme integration
- Display item name and type in message
- Keyboard handling: 'y' (confirm), 'n'/Escape (cancel)
- Auto-sizing based on content

**Styles (styles.go):**
- Use theme colors (error for destructive action warning)
- Match existing modal patterns from modal/styles.go
- Responsive width/height

**Testing (confirm_modal_test.go):**
- Test modal creation with different entity types
- Test keyboard handling (y/n/escape)
- Test message generation
- Test view rendering

**Dependencies:** None (standalone component)

---

## PHASE 2: SERVICE LAYER - Cascade Soft Delete (Task-68B)

**Goal:** Implement cascade soft delete for projects with subprojects

### Files to Modify:
1. `internal/service/project_service.go`
2. `internal/service/project_service_test.go`

### Implementation Details:

**Add to project_service.go:**
```go
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, id string) error {
    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        return s.softDeleteRecursive(ctx, tx, id)
    })
}

func (s *ProjectService) softDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error {
    // Get all child projects
    children, err := q.ListProjectsByParent(ctx, sql.NullString{
        String: projectID,
        Valid:  true,
    })
    if err != nil {
        return err
    }
    
    // Recursively soft delete children
    for _, child := range children {
        if err := s.softDeleteRecursive(ctx, q, child.ID); err != nil {
            return err
        }
    }
    
    // Soft delete tasks in this project
    now := time.Now()
    if err := q.SoftDeleteTasksByProject(ctx, projectID, &now); err != nil {
        return err
    }
    
    // Soft delete this project
    params := db.SoftDeleteProjectParams{
        ID:        projectID,
        DeletedAt: &now,
    }
    _, err = q.SoftDeleteProject(ctx, params)
    return err
}
```

**Database Query Needed:**
Add to `internal/db/projects.sql`:
```sql
-- name: SoftDeleteTasksByProject :exec
UPDATE tasks
SET deleted_at = ?
WHERE project_id = ? AND deleted_at IS NULL;
```

**Testing (project_service_test.go):**
- Test cascade soft delete with nested projects
- Test transaction rollback on error
- Test soft delete of tasks within project
- Test error handling (not found, db errors)

**Dependencies:** None (can be done in parallel with Task-68A)

---

## PHASE 3: TUI INTEGRATION - Delete Handler & State Management (Task-68C)

**Goal:** Wire up delete functionality in TUI with confirmation modal

### Files to Modify:
1. `internal/tui/model.go` - Add confirmation modal state
2. `internal/tui/app.go` - Add keyboard handler for 'd' key
3. `internal/tui/commands.go` - Add delete commands
4. `internal/tui/messages.go` - Add delete messages
5. `internal/tui/renderer.go` - Update footer
6. `internal/tui/constants.go` - Add delete key binding

### Implementation Details:

**model.go additions:**
```go
type Model struct {
    // ... existing fields ...
    
    // Delete confirmation modal
    confirmModal       *confirmmodal.Modal
    isConfirmModalOpen bool
    deleteEntityType   confirmmodal.EntityType
    deleteEntityID     string
}
```

**messages.go additions:**
```go
type DeleteConfirmedMsg struct {
    EntityType confirmmodal.EntityType
    EntityID   string
}

type DeleteSuccessMsg struct {
    EntityType confirmmodal.EntityType
}

type DeleteErrorMsg struct {
    Err        error
    EntityType confirmmodal.EntityType
}
```

**commands.go additions:**
```go
func DeleteSubareaCmd(subareaSvc service.SubareaServiceInterface, id string) tea.Cmd {
    return func() tea.Msg {
        err := subareaSvc.SoftDelete(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{Err: err, EntityType: confirmmodal.EntityTypeSubarea}
        }
        return DeleteSuccessMsg{EntityType: confirmmodal.EntityTypeSubarea}
    }
}

func DeleteProjectCmd(projectSvc service.ProjectServiceInterface, id string) tea.Cmd {
    return func() tea.Msg {
        err := projectSvc.SoftDeleteWithCascade(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{Err: err, EntityType: confirmmodal.EntityTypeProject}
        }
        return DeleteSuccessMsg{EntityType: confirmmodal.EntityTypeProject}
    }
}

func DeleteTaskCmd(taskSvc service.TaskServiceInterface, id string) tea.Cmd {
    return func() tea.Msg {
        err := taskSvc.SoftDelete(context.Background(), id)
        if err != nil {
            return DeleteErrorMsg{Err: err, EntityType: confirmmodal.EntityTypeTask}
        }
        return DeleteSuccessMsg{EntityType: confirmmodal.EntityTypeTask}
    }
}
```

**app.go keyboard handler (in Update function):**
```go
case key.Matches(msg, m.keys.delete):
    // Check if confirmation modal is open
    if m.isConfirmModalOpen {
        // Let modal handle the key
        newModal, cmd := m.confirmModal.Update(msg)
        m.confirmModal = newModal
        return m, cmd
    }
    
    // Get currently selected item based on focus
    switch m.focus {
    case FocusSubareas:
        if len(m.subareas) == 0 || m.selectedSubareaIndex >= len(m.subareas) {
            return m, nil // No-op on empty column
        }
        subarea := m.subareas[m.selectedSubareaIndex]
        m.confirmModal = confirmmodal.New(
            subarea.Name,
            confirmmodal.EntityTypeSubarea,
            subarea.ID,
        )
        m.isConfirmModalOpen = true
        
    case FocusProjects:
        if m.projectTree == nil {
            return m, nil // No-op on empty column
        }
        project := m.getSelectedProject()
        if project == nil {
            return m, nil
        }
        m.confirmModal = confirmmodal.New(
            project.Name,
            confirmmodal.EntityTypeProject,
            project.ID,
        )
        m.isConfirmModalOpen = true
        
    case FocusTasks:
        // Get current task from grouped tasks
        task := m.getSelectedTask()
        if task == nil {
            return m, nil // No-op on empty column
        }
        m.confirmModal = confirmmodal.New(
            task.Title,
            confirmmodal.EntityTypeTask,
            task.ID,
        )
        m.isConfirmModalOpen = true
    }
    return m, nil
```

**Message handlers:**
```go
func (m *Model) handleConfirmModalClose(msg confirmmodal.CancelMsg) {
    m.isConfirmModalOpen = false
    m.confirmModal = nil
}

func (m *Model) handleConfirmModalConfirm(msg confirmmodal.ConfirmMsg) {
    m.isConfirmModalOpen = false
    m.confirmModal = nil
    
    // Dispatch appropriate delete command
    switch msg.EntityType {
    case confirmmodal.EntityTypeSubarea:
        return m, DeleteSubareaCmd(m.subareaSvc, msg.EntityID)
    case confirmmodal.EntityTypeProject:
        return m, DeleteProjectCmd(m.projectSvc, msg.EntityID)
    case confirmmodal.EntityTypeTask:
        return m, DeleteTaskCmd(m.taskSvc, msg.EntityID)
    }
}

func (m *Model) handleDeleteSuccess(msg DeleteSuccessMsg) {
    // Show success toast
    entityName := string(msg.EntityType)
    m.toasts = append(m.toasts, toast.NewSuccess(
        fmt.Sprintf("%s deleted successfully", entityName),
    ))
    
    // Refresh the appropriate column
    switch msg.EntityType {
    case confirmmodal.EntityTypeSubarea:
        return m, LoadSubareasCmd(m.subareaSvc, m.areas[m.selectedTab].ID)
    case confirmmodal.EntityTypeProject:
        return m, LoadProjectsCmd(m.projectSvc, m.subareas[m.selectedSubareaIndex].ID)
    case confirmmodal.EntityTypeTask:
        return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
    }
}

func (m *Model) handleDeleteError(msg DeleteErrorMsg) {
    // Show error toast
    m.toasts = append(m.toasts, toast.NewError(
        fmt.Sprintf("Failed to delete %s: %v", msg.EntityType, msg.Err),
    ))
}
```

**renderer.go footer update:**
```go
func (m *Model) renderFooter() string {
    shortcuts := []string{
        "h/l: columns",
        "j/k: nav",
        "a: add",
        "x: toggle",
        "d: delete",  // NEW
        "?: help",
        "q: quit",
    }
    // ... rest of footer rendering
}
```

**constants.go:**
```go
var defaultKeyMap = keymap{
    // ... existing keys ...
    delete: key.NewBinding(
        key.WithKeys("d"),
        key.WithHelp("d", "delete"),
    ),
}
```

**Testing:**
- `internal/tui/app_delete_test.go` - New test file
  - Test 'd' key on each column
  - Test confirmation modal opens with correct item name
  - Test 'y' confirmation triggers delete
  - Test 'n' and Escape cancel
  - Test success toast appears
  - Test error toast on failure
  - Test column refresh after delete
  - Test no-op on empty columns

**Dependencies:** Requires Task-68A and Task-68B

---

## PHASE 4: DOCUMENTATION & TESTING (Task-68D)

**Goal:** Update documentation and add comprehensive tests

### Files to Modify:
1. `docs/TUI.md` - Update keyboard shortcuts section
2. `internal/tui/mocks/services.go` - Add mock delete methods
3. Test files for all new components

### Documentation Updates:

**docs/TUI.md:**
- Add 'd: delete' to Actions keyboard shortcuts table
- Document delete confirmation modal behavior
- Document cascade delete for projects
- Add to Quick Reference section
- Update footer example

**Example addition to keyboard shortcuts:**
```markdown
### Actions

| Key | Action | Description |
|-----|--------|-------------|
| `d` | Delete | Delete focused item (shows confirmation) |
```

**Mock services update (mocks/services.go):**
```go
type MockSubareaService struct {
    // ... existing fields ...
    SoftDeleteFunc func(ctx context.Context, id string) error
}

func (m *MockSubareaService) SoftDelete(ctx context.Context, id string) error {
    if m.SoftDeleteFunc != nil {
        return m.SoftDeleteFunc(ctx, id)
    }
    return nil
}

// Similar for MockProjectService and MockTaskService
```

### Test Coverage Goals:

**Unit Tests:**
- confirmmodal package: 90%+ coverage
- Delete commands: 100% coverage
- Message handlers: 100% coverage

**Integration Tests:**
- End-to-end delete flow for each entity type
- Error scenarios (database failures)
- Edge cases (empty columns, not found errors)

**Test Files:**
1. `internal/tui/confirmmodal/confirm_modal_test.go`
2. `internal/service/project_service_cascade_test.go`
3. `internal/tui/app_delete_test.go`

**Running Tests:**
```bash
# Unit tests
go test ./internal/tui/confirmmodal/... -v
go test ./internal/service/... -v -run TestCascade

# Integration tests
go test ./internal/tui/... -v -run TestDelete

# Coverage
go test ./internal/tui/... -cover
go test ./internal/service/... -cover
```

**Dependencies:** Can run in parallel with Task-68C

---

## TASK SPLIT SUMMARY

### Task-68A: Confirmation Modal Component (FOUNDATION)
- **Priority:** HIGH
- **Dependencies:** None
- **Files:** 3 new files
- **Tests:** Unit tests for modal component
- **Estimated effort:** 4-6 hours

### Task-68B: Cascade Soft Delete Service (SERVICE LAYER)
- **Priority:** HIGH
- **Dependencies:** None
- **Files:** 2 modified files, 1 SQL query
- **Tests:** Service layer tests
- **Estimated effort:** 3-4 hours
- **Can run in parallel with:** Task-68A

### Task-68C: TUI Integration (INTEGRATION)
- **Priority:** HIGH
- **Dependencies:** Task-68A, Task-68B
- **Files:** 6 modified files
- **Tests:** Integration tests
- **Estimated effort:** 6-8 hours
- **Sequential after:** Task-68A and Task-68B

### Task-68D: Documentation & Testing (DOCUMENTATION)
- **Priority:** MEDIUM
- **Dependencies:** Task-68C
- **Files:** 2 modified files
- **Tests:** All test coverage
- **Estimated effort:** 3-4 hours
- **Can start partial work in parallel with:** Task-68C

---

## EXECUTION ORDER

### Parallel Track 1: UI Component
1. Task-68A: Confirmation Modal Component

### Parallel Track 2: Service Layer
1. Task-68B: Cascade Soft Delete

### Sequential Track:
3. Task-68C: TUI Integration (depends on 1 & 2)
4. Task-68D: Documentation & Testing (can start docs early, finish after 3)

**Optimal Execution:**
- Start Task-68A and Task-68B in parallel
- Once both complete, start Task-68C
- Task-68D documentation can start early, testing completes after Task-68C

---

## ACCEPTANCE CRITERIA MAPPING

- **AC #1-3:** Implemented in Task-68C (keyboard handler)
- **AC #4:** Implemented in Task-68A (confirmation dialog)
- **AC #5-6:** Implemented in Task-68C (y/n/escape handling)
- **AC #7-8:** Implemented in Task-68C (success/error toasts)
- **AC #9:** Implemented in Task-68B (cascade delete)
- **AC #10:** Implemented in Task-68C (no-op on empty column)
- **AC #11:** Implemented in Task-68C (footer update)
- **AC #12:** Implemented in Task-68A (modal patterns)

---

## RISKS & MITIGATIONS

**Risk 1: Cascade delete performance with deeply nested projects**
- Mitigation: Use transactions, add depth limit if needed
- Test with realistic data sets

**Risk 2: User accidentally deleting parent project**
- Mitigation: Clear confirmation message showing item name
- Future: Add undo functionality (separate task)

**Risk 3: State refresh after delete might cause UI flicker**
- Mitigation: Use optimistic updates where possible
- Show loading state during refresh

**Risk 4: Modal might conflict with existing modals**
- Mitigation: Follow existing modal state management pattern
- Ensure only one modal can be open at a time

---

## ROLLBACK PLAN

If issues arise:
1. Delete feature can be disabled by removing 'd' key binding
2. Revert to using only CLI delete commands
3. Modal component is isolated - easy to remove

---

## POST-IMPLEMENTATION CHECKLIST

- [ ] All acceptance criteria verified
- [ ] Unit tests pass with 80%+ coverage
- [ ] Integration tests pass
- [ ] Manual testing on all entity types
- [ ] Documentation updated
- [ ] No console errors or warnings
- [ ] Performance acceptable (no lag on delete)
- [ ] Accessibility: keyboard-only navigation works
- [ ] Toast notifications display correctly
- [ ] Footer shows delete shortcut
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Split

This task has been split into 4 subtasks for better manageability:

- **TASK-68.1**: Confirmation Modal Component (foundation)
- **TASK-68.2**: Cascade Soft Delete Service (service layer)
- **TASK-68.3**: TUI Delete Integration (integration)
- **TASK-68.4**: Delete Documentation & Testing (documentation)

See individual subtasks for detailed implementation details.

## Execution Strategy

**Parallel:** Start Task-68.1 and Task-68.2 simultaneously (no dependencies)

**Sequential:** Task-68.3 after both 68.1 and 68.2 complete, then Task-68.4

**Total Estimated Time:** 16-22 hours (reduced with parallel execution)
<!-- SECTION:NOTES:END -->
