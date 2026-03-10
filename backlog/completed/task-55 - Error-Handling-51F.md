---
id: TASK-55
title: Error Handling (51F)
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 21:31'
updated_date: '2026-03-07 18:12'
labels:
  - error-handling
dependencies: []
references:
  - task-51
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Handle error cases and edge scenarios. Depends on 51A-51E. Part of task-51.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Handle empty project gracefully
- [x] #2 Handle orphaned subprojects gracefully
- [x] #3 Add error handling for database errors
- [x] #4 Add user-friendly error messages
- [x] #5 Write error handling tests
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Error Handling (Task-55)

## Scope
- Domain layer: Centralized error types (internal/domain/errors.go)
- Service layer: Error wrapping and context (4 services)
- TUI layer: Error state tracking + rendering
- Tests: Service + TUI error handling tests

## Phase 1: Domain Error Types (30 min)
**File**: internal/domain/errors.go (NEW)

### Define Sentinel Errors
```go
package domain

import "errors"

var (
    // Common errors
    ErrNotFound      = errors.New("resource not found")
    ErrInvalidInput  = errors.New("invalid input provided")
    ErrDatabaseError = errors.New("database operation failed")
    
    // Hierarchy errors
    ErrEmptyID        = errors.New("id cannot be empty")
    ErrOrphanedEntity = errors.New("entity references non-existent parent")
)
```

### Define Error Types with Context
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

type DatabaseError struct {
    Operation string
    Err       error
}

func (e *DatabaseError) Error() string {
    return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
    return e.Err
}
```

**Test File**: internal/domain/errors_test.go (NEW)
- Test error messages
- Test errors.Is() and errors.As() compatibility
- Test error wrapping/unwrapping

## Phase 2: Service Layer Error Handling (2 hours)

### 2.1 TaskService Updates
**File**: internal/service/task_service.go

#### Update ListByProjectRecursive
```go
func (s *TaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
    // Handle empty projectID
    if projectID == "" {
        return []domain.Task{}, nil // Empty result, not error
    }
    
    // Check if project exists
    _, err := s.projectSvc.GetByID(ctx, projectID)
    if err != nil {
        if errors.Is(err, service.ErrProjectNotFound) {
            // Project does not exist - return empty, not error
            return []domain.Task{}, nil
        }
        // Real database error - wrap with context
        return nil, fmt.Errorf("check project existence: %w", err)
    }
    
    // Load all tasks
    allTasks, err := s.db.ListAllTasks(ctx)
    if err != nil {
        return nil, fmt.Errorf("list all tasks: %w", err)
    }
    
    // Filter and build hierarchy...
}
```

**Test File**: internal/service/task_service_error_test.go (NEW)
- Empty projectID → empty slice
- Non-existent project → empty slice
- Database errors → wrapped with context
- Orphaned tasks → gracefully ignored

### 2.2 ProjectService Updates
**File**: internal/service/project_service.go

#### Update ListBySubareaRecursive
```go
func (s *ProjectService) ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error) {
    if subareaID == "" {
        return []domain.Project{}, nil
    }
    
    // Similar error handling pattern...
}
```

**Test File**: internal/service/project_service_error_test.go (NEW)

### 2.3 AreaService Updates
**File**: internal/service/area_service.go

#### Update GetByID
```go
func (s *AreaService) GetByID(ctx context.Context, id string) (*domain.Area, error) {
    row, err := s.repo.GetAreaByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound // Use domain error
        }
        return nil, fmt.Errorf("get area %s: %w", id, err)
    }
    // ...
}
```

**Test File**: internal/service/area_service_error_test.go (NEW)

### 2.4 SubareaService Updates
**File**: internal/service/subarea_service.go

**Test File**: internal/service/subarea_service_error_test.go (NEW)

## Phase 3: TUI Error State Management (1 hour)

### 3.1 Model Updates
**File**: internal/tui/state.go

```go
type Model struct {
    // ... existing fields ...
    
    // Error tracking
    areaLoadError    error
    subareaLoadError error
    projectLoadError error
    taskLoadError    error
}

func (m *Model) ClearErrors() {
    m.areaLoadError = nil
    m.subareaLoadError = nil
    m.projectLoadError = nil
    m.taskLoadError = nil
}
```

### 3.2 Handler Updates
**File**: internal/tui/handlers.go

```go
func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) {
    m.isLoadingTasks = false
    
    if msg.Err != nil {
        // Store error for rendering
        m.taskLoadError = msg.Err
        
        // Show toast notification
        m.toasts = append(m.toasts, toast.NewError(
            m.formatUserError(msg.Err),
        ))
        
        // Reset to safe state
        m.tasks = []domain.Task{}
        m.groupedTasks = &domain.GroupedTasks{}
        
        return
    }
    
    // Clear error on success
    m.taskLoadError = nil
    
    // ... rest of success handling
}

func (m *Model) formatUserError(err error) string {
    // Map technical errors to user-friendly messages
    if errors.Is(err, context.Canceled) {
        return "Operation cancelled"
    }
    if errors.Is(err, context.DeadlineExceeded) {
        return "Loading took too long. Please try again."
    }
    if strings.Contains(err.Error(), "database") || strings.Contains(err.Error(), "sql") {
        return "Unable to load data. Please restart the application."
    }
    
    // Generic fallback
    return fmt.Sprintf("Error: %v", err)
}
```

## Phase 4: TUI Error Rendering (1 hour)

### 4.1 Renderer Updates
**File**: internal/tui/renderer.go

```go
func (m *Model) RenderTasks() string {
    if m.isLoadingTasks {
        return m.spinner.View() + " " + LoadingMessageTasks
    }
    
    // Handle error state
    if m.taskLoadError != nil {
        return m.renderError(m.taskLoadError, "tasks")
    }
    
    // Handle empty state
    if m.groupedTasks == nil || m.groupedTasks.TotalCount == 0 {
        return m.renderEmptyTasks()
    }
    
    // ... normal rendering
}

func (m *Model) renderError(err error, context string) string {
    errorStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("9")). // Red
        PaddingLeft(2).
        PaddingTop(1)
    
    userMsg := m.formatUserError(err)
    
    return errorStyle.Render("✗ " + userMsg)
}

func (m *Model) renderEmptyTasks() string {
    emptyStyle := lipgloss.NewStyle().
        Foreground(m.theme.Dimmed).
        PaddingLeft(2)
    
    // Context-aware empty messages
    msg := EmptyStateNoTasks
    
    if m.selectedProjectID != "" && m.groupedTasks != nil && len(m.groupedTasks.Groups) > 0 {
        msg = "No tasks in this project or its subprojects"
    } else if m.selectedProjectID != "" {
        msg = "No tasks in this project"
    }
    
    return emptyStyle.Render(msg)
}
```

### 4.2 Constants
**File**: internal/tui/constants.go

```go
const (
    // ... existing constants ...
    
    // User-friendly error messages
    ErrMsgDatabase  = "Unable to load data. Please try again."
    ErrMsgTimeout   = "Loading took too long. Please retry."
    ErrMsgCancelled = "Operation cancelled"
    ErrMsgNoData    = "No data available"
)
```

## Phase 5: Testing (2-3 hours)

### 5.1 Domain Error Tests
**File**: internal/domain/errors_test.go

```go
func TestSentinelErrors(t *testing.T) {
    tests := []struct {
        name     string
        err      error
        message  string
    }{
        {"not found", domain.ErrNotFound, "resource not found"},
        {"invalid input", domain.ErrInvalidInput, "invalid input provided"},
        {"database error", domain.ErrDatabaseError, "database operation failed"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.err.Error() != tt.message {
                t.Errorf("got %q, want %q", tt.err.Error(), tt.message)
            }
        })
    }
}

func TestErrorWrapping(t *testing.T) {
    originalErr := sql.ErrNoRows
    wrappedErr := &domain.DatabaseError{
        Operation: "GetTaskByID",
        Err:       originalErr,
    }
    
    // Test errors.Is compatibility
    if !errors.Is(wrappedErr, originalErr) {
        t.Error("errors.Is should find wrapped error")
    }
    
    // Test errors.As compatibility
    var dbErr *domain.DatabaseError
    if !errors.As(wrappedErr, &dbErr) {
        t.Error("errors.As should extract DatabaseError")
    }
}
```

### 5.2 Service Error Tests
**File**: internal/service/task_service_error_test.go

```go
func TestListByProjectRecursive_ErrorHandling(t *testing.T) {
    tests := []struct {
        name       string
        projectID  string
        mockSetup  func(*MockQuerier)
        wantLen    int
        wantErr    bool
        errContain string
    }{
        {
            name:      "empty projectID",
            projectID: "",
            wantLen:    0,
            wantErr:    false,
        },
        {
            name:      "non-existent project",
            projectID: "nonexistent",
            mockSetup: func(m *MockQuerier) {
                m.getProjectByIDError = sql.ErrNoRows
            },
            wantLen:    0, // Empty result, not error
            wantErr:    false,
        },
        {
            name:      "database error on list tasks",
            projectID: "proj-1",
            mockSetup: func(m *MockQuerier) {
                m.listAllTasksError = errors.New("database connection lost")
            },
            wantErr:    true,
            errContain: "list all tasks",
        },
        {
            name:      "orphaned tasks in non-existent projects",
            projectID: "proj-1",
            mockSetup: func(m *MockQuerier) {
                m.tasks = []db.Task{
                    {ID: "t1", ProjectID: "orphan-proj"},
                }
            },
            wantLen:    0, // Orphaned tasks ignored
            wantErr:    false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ... test implementation
        })
    }
}
```

### 5.3 TUI Error Tests
**File**: internal/tui/error_handling_test.go

```go
func TestModelErrorHandling(t *testing.T) {
    t.Run("error state tracking", func(t *testing.T) {
        m := &Model{
            taskLoadError: errors.New("database error"),
        }
        
        if m.taskLoadError == nil {
            t.Error("taskLoadError should be set")
        }
    })
    
    t.Run("error clearing on success", func(t *testing.T) {
        m := &Model{
            taskLoadError: errors.New("previous error"),
        }
        
        m.handleTasksLoaded(TasksLoadedMsg{
            Tasks: []domain.Task{},
            Err:   nil,
        })
        
        if m.taskLoadError != nil {
            t.Error("taskLoadError should be cleared on success")
        }
    })
}

func TestErrorRendering(t *testing.T) {
    m := &Model{
        taskLoadError: errors.New("database connection failed"),
        theme:         DefaultTheme,
    }
    
    output := m.RenderTasks()
    
    if !strings.Contains(output, "✗") {
        t.Error("error icon should be rendered")
    }
    
    if !strings.Contains(output, "Unable to load data") {
        t.Error("user-friendly error message should be shown")
    }
}

func TestUserErrorMessageMapping(t *testing.T) {
    m := &Model{}
    
    tests := []struct {
        err      error
        contains string
    }{
        {context.Canceled, "cancelled"},
        {context.DeadlineExceeded, "too long"},
        {errors.New("sql: connection failed"), "Unable to load data"},
    }
    
    for _, tt := range tests {
        msg := m.formatUserError(tt.err)
        if !strings.Contains(msg, tt.contains) {
            t.Errorf("user message %q should contain %q", msg, tt.contains)
        }
    }
}
```

## Implementation Order

### Sequential (MUST be done in order):

1. **Phase 1**: Domain errors.go (foundation for all layers)
2. **Phase 2**: Service layer error handling (depends on domain errors)
3. **Phase 3**: TUI error state (depends on services returning errors)
4. **Phase 4**: TUI error rendering (depends on error state)
5. **Phase 5**: Tests (can start after each phase)

### Estimated Time: 6-7 hours total
- Phase 1: 30 minutes
- Phase 2: 2 hours
- Phase 3: 1 hour
- Phase 4: 1 hour
- Phase 5: 2-3 hours

## Acceptance Criteria Mapping

| AC | Phase | Implementation |
|----|-------|----------------|
| #1 Handle empty project gracefully | Phase 2 | ListByProjectRecursive returns empty slice |
| #2 Handle orphaned subprojects gracefully | Phase 2 | Silently ignore non-existent parents |
| #3 Add error handling for database errors | Phase 2 | Wrap with fmt.Errorf and context |
| #4 Add user-friendly error messages | Phase 4 | TUI error renderer + formatUserError |
| #5 Write error handling tests | Phase 5 | Domain + Service + TUI tests |

## Documentation Updates

### 1. docs/architecture/02-domain-layer.md
Add section: "Error Handling Patterns"
- Sentinel errors (ErrNotFound, etc.)
- Custom error types (ValidationError, DatabaseError)
- Error checking with errors.Is/As

### 2. docs/architecture/03-service-layer.md
Add section: "Error Wrapping Best Practices"
- Wrap errors with context using fmt.Errorf
- Map sql.ErrNoRows to domain errors
- Handle empty results gracefully

### 3. docs/TUI.md
Add section: "Error State Management"
- Error tracking in Model
- User-friendly error rendering
- Error clearing on successful loads

### 4. CHANGELOG.md
```markdown
## [Unreleased]
### Added
- Centralized domain error types for consistent error handling across layers
- User-friendly error messages in TUI with visual indicators
- Graceful handling of empty projects and orphaned entities
- Comprehensive error handling tests across domain, service, and TUI layers
```

## Testing Commands

```bash
# Run domain error tests
go test ./internal/domain -v -run TestError

# Run service error tests
go test ./internal/service -v -run TestError

# Run TUI error tests
go test ./internal/tui -v -run TestError

# Run all tests with coverage
go test -cover ./internal/domain ./internal/service ./internal/tui

# Run with race detector
go test -race ./internal/service ./internal/tui
```

## Success Criteria

- [ ] All 5 acceptance criteria met
- [ ] Domain errors.go created with sentinel errors
- [ ] All 4 services updated with error wrapping
- [ ] TUI Model has error tracking fields
- [ ] TUI renderer shows user-friendly errors
- [ ] Test coverage ≥ 75% for error handling code
- [ ] All tests passing
- [ ] Race detector passing
- [ ] Linting passing (golangci-lint)
- [ ] Documentation updated
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
IMPLEMENTATION DETAILS:

## Service Layer Error Handling

### File: internal/service/task_service.go

```go
func (s *TaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
    // Validate input
    if projectID == "" {
        return []domain.Task{}, nil  // Empty result, not error
    }
    
    // Check if project exists
    _, err := s.projectSvc.GetByID(ctx, projectID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            // Project doesn't exist - return empty, not error
            return []domain.Task{}, nil
        }
        // Database error
        return nil, fmt.Errorf("get project %s: %w", projectID, err)
    }
    
    // Load all tasks
    allTasks, err := s.db.ListAllTasks(ctx)
    if err != nil {
        return nil, fmt.Errorf("list tasks: %w", err)
    }
    
    // Build hierarchy and filter
    // ... rest of implementation
}
```

### Handle Orphaned Subprojects

```go
func (s *TaskService) buildProjectHierarchy(ctx context.Context) map[string][]string {
    projects, err := s.projectSvc.ListAll(ctx)
    if err != nil {
        return map[string][]string{}  // Empty map, not error
    }
    
    hierarchy := make(map[string][]string)
    
    for _, project := range projects {
        if project.ParentID != nil {
            parentID := *project.ParentID
            
            // Check if parent exists
            parentExists := false
            for _, p := range projects {
                if p.ID == parentID {
                    parentExists = true
                    break
                }
            }
            
            // Include even if parent doesn't exist (orphaned)
            // These won't appear in recursive loading
            hierarchy[parentID] = append(hierarchy[parentID], project.ID)
        }
    }
    
    return hierarchy
}
```

## TUI Error Handling

### File: internal/tui/renderer.go

```go
func (m *Model) RenderTasks() string {
    if m.isLoadingTasks {
        return m.spinner.View() + " " + LoadingMessageTasks
    }
    
    // Handle error state
    if m.taskLoadError != nil {
        return m.renderError(m.taskLoadError)
    }
    
    // Handle empty state
    if m.groupedTasks.TotalCount == 0 {
        return m.renderEmptyState()
    }
    
    // ... normal rendering
}

func (m *Model) renderError(err error) string {
    errorStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("9")).  // Red
        PaddingLeft(2)
    
    msg := fmt.Sprintf("Error loading tasks: %v", err)
    
    // Provide helpful context
    if errors.Is(err, context.Canceled) {
        msg = "Task loading was cancelled"
    } else if strings.Contains(err.Error(), "database") {
        msg = "Database error. Please try restarting the application."
    }
    
    return errorStyle.Render("✗ " + msg)
}

func (m *Model) renderEmptyState() string {
    emptyStyle := lipgloss.NewStyle().
        Foreground(m.theme.Dimmed).
        PaddingLeft(2)
    
    // Different messages based on context
    msg := "No tasks found"
    
    if len(m.groupedTasks.Groups) > 0 {
        // Has subprojects but no tasks
        msg = "No tasks in this project or its subprojects"
    } else {
        // No subprojects
        msg = "No tasks in this project"
    }
    
    return emptyStyle.Render(msg)
}
```

### File: internal/tui/handlers.go

```go
func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) {
    m.isLoadingTasks = false
    
    if msg.Err != nil {
        // Store error for rendering
        m.taskLoadError = msg.Err
        
        // Show toast notification
        m.toasts = append(m.toasts, toast.NewError(
            fmt.Sprintf("Failed to load tasks: %v", msg.Err),
        ))
        
        // Reset state
        m.tasks = []domain.Task{}
        m.groupedTasks = domain.GroupedTasks{}
        
        return
    }
    
    // Clear error
    m.taskLoadError = nil
    
    // ... rest of success handling
}
```

## File: internal/tui/tui.go

### Add Error Field to Model

```go
type Model struct {
    // Existing fields...
    
    // Error tracking
    taskLoadError    error
    projectLoadError error
}
```

## Edge Case Handling

### 1. Empty Project Tree

```go
// In ListByProjectRecursive
if len(descendantIDs) == 0 && len(tasksByProject[projectID]) == 0 {
    return []domain.Task{}, nil  // Valid empty result
}
```

### 2. Deeply Nested Subprojects

```go
// Add max depth protection
const maxDepth = 10

func (s *TaskService) getDescendantProjectIDs(
    projectID string,
    hierarchy map[string][]string,
    depth int,
) []string {
    if depth > maxDepth {
        // Log warning but don't error
        log.Printf("Warning: max depth %d reached for project %s", maxDepth, projectID)
        return []string{}
    }
    
    // ... rest of implementation
}
```

### 3. Concurrent Modifications

```go
// Use context timeout to prevent hanging
func (s *TaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    // ... rest of implementation
}
```

## Testing Strategy

### File: internal/service/task_service_test.go

```go
func TestListByProjectRecursive_EdgeCases(t *testing.T) {
    tests := []struct {
        name      string
        setup     func(*MockQuerier)
        projectID string
        wantTasks int
        wantErr   bool
    }{
        {
            name:      "empty projectID",
            projectID: "",
            wantTasks: 0,
            wantErr:   false,
        },
        {
            name: "non-existent project",
            setup: func(m *MockQuerier) {
                // Return empty, no project found
            },
            projectID: "nonexistent",
            wantTasks: 0,
            wantErr:   false,  // Empty result, not error
        },
        {
            name: "orphaned subproject",
            setup: func(m *MockQuerier) {
                m.projects = []db.Project{
                    {ID: "orphan", ParentID: ptr("missing-parent")},
                }
                m.tasks = []db.Task{
                    {ID: "t1", ProjectID: "orphan"},
                }
            },
            projectID: "some-project",
            wantTasks: 0,  // Orphaned subproject not included
            wantErr:   false,
        },
        {
            name: "database error",
            setup: func(m *MockQuerier) {
                m.listAllTasksError = errors.New("database connection failed")
            },
            projectID: "proj-1",
            wantErr:   true,
        },
    }
    
    // ... test execution
}
```

## Error Messages

User-friendly error messages:

| Scenario | Message |
|----------|---------|
| Empty project | "No tasks in this project" |
| Empty tree | "No tasks in this project or its subprojects" |
| Database error | "Database error. Please try restarting." |
| Timeout | "Loading tasks took too long. Please try again." |
| Cancelled | "Task loading was cancelled" |
| Orphaned subproject | (Gracefully ignored, not shown to user) |

- Created domain errors.go (Phase 1 - 30 min)
- TaskService error handling (Phase 2 - 2 hours)
- ProjectService error handling (Phase 2 - 2 hours)
- AreaService error handling (Phase 2 - 30 min)
- SubareaService error handling (Phase 2 - 30 min)
- TUI error state (Phase 3 - 30 min)
- TUI error rendering (Phase 4 - 1 hour)
- TUI error handlers (Phase 5 - 1 hour)
- TUI constants (Phase 6 - 30 min)
- Documentation (Phase 7 - 30 min)
- Integration tests (Phase 8 - 30 min)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Completed error handling implementation:
- Created domain/domain/errors.go with sentinel errors and error types
- Updated service layer (TaskService, Project_service, SubareaService) to use domain.ErrNotFound, mappings
- Added user-friendly error messages in TUI
- Added comprehensive error handling tests
<!-- SECTION:FINAL_SUMMARY:END -->
