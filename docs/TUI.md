# Terminal User Interface (TUI) Documentation

This document provides comprehensive documentation for the Dopadone Terminal User Interface, including architecture, components, and implementation details.

## Overview

The TUI provides an interactive, keyboard-driven interface for managing projects, subareas, and tasks. Built with the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, it follows The Elm Architecture for state management.

## Architecture

### Core Components

```
internal/tui/
├── app.go              # Main application logic and Update/View functions
├── model.go            # Application state and data structures
├── tui.go              # Entry point and initialization
├── constants.go        # Key bindings and constants
├── messages.go         # Bubble Tea message types
├── commands.go         # Bubble Tea commands
├── converters.go       # Data conversion utilities
├── mocks/              # Service mock implementations for testing
│   ├── services.go     # Mock implementations of service interfaces
│   ├── helpers.go      # Mock setup helper functions
│   └── README.md       # Mock usage documentation
├── help/               # Help modal component
│   ├── help.go         # Help modal logic
│   └── styles.go       # Lipgloss styling
├── spacemenu/          # Command menu component (LazyVim-style which-key)
│   ├── spacemenu.go    # Command menu logic
│   ├── types.go        # Menu state and action types
│   └── styles.go       # Lipgloss styling
├── toast/              # Toast notification component
│   ├── toast.go        # Toast logic with auto-dismiss
│   └── styles.go       # Error/success styling
├── confirmmodal/       # Confirmation modal component
│   ├── modal.go        # Confirmation modal logic (y/n/esc)
│   └── styles.go       # Destructive action warning styling
├── modal/              # Quick-add modal component
│   ├── modal.go        # Modal dialog logic
│   ├── styles.go       # Modal styling
│   └── validation.go   # Input validation
├── tree/               # Project tree navigation
│   ├── node.go         # Tree node representation
│   ├── builder.go      # Tree construction from data
│   ├── navigation.go   # Navigation logic
│   └── renderer.go     # Tree rendering
└── views/              # UI view components
    ├── tabs.go         # Area tabs component
    ├── columns.go      # 3-column browser layout
    └── styles.go       # View styling
```

### Service Layer Integration

The TUI follows **dependency injection** principles and depends on service layer interfaces rather than direct database access:

```go
type Model struct {
    // Service layer dependencies (injected)
    areaSvc     service.AreaServiceInterface
    subareaSvc  service.SubareaServiceInterface
    projectSvc  service.ProjectServiceInterface
    taskSvc     service.TaskServiceInterface
    
    // ... other fields
}
```

**Benefits**:
- **Testability**: Services can be individually mocked for unit tests
- **Loose coupling**: TUI layer isolated from database implementation details
- **Flexibility**: Service implementations can be swapped without TUI changes
- **Clear boundaries**: Well-defined interfaces between layers

**Initialization**:
```go
// In cmd/dopa/tui.go
func runTUI() error {
    // Get service container
    container := db.GetServices()
    defer container.Close()
    
    // Create TUI with injected services
    model := tui.New(
        container.AreaService(),
        container.SubareaService(),
        container.ProjectService(),
        container.TaskService(),
    )
    
    // Start program
    p := tea.NewProgram(model)
    return p.Start()
}
```

**See Also**: [Service Layer Documentation](../internal/service/README.md)

### Service Layer Architecture

The TUI interacts with the data layer exclusively through **service interfaces**, not directly with the database:

```
┌─────────────────────────────────────────┐
│            TUI Layer (app.go)           │
│  - User interface logic                 │
│  - State management                     │
│  - Event handling                       │
└────────────────┬────────────────────────┘
                 │ depends on
                 ▼
┌─────────────────────────────────────────┐
│        Service Layer Interfaces         │
│  - AreaServiceInterface                 │
│  - SubareaServiceInterface              │
│  - ProjectServiceInterface              │
│  - TaskServiceInterface                 │
└────────────────┬────────────────────────┘
                 │ implemented by
                 ▼
┌─────────────────────────────────────────┐
│      Service Implementations            │
│  - Business logic                       │
│  - Validation                           │
│  - Data transformation                  │
└────────────────┬────────────────────────┘
                 │ uses
                 ▼
┌─────────────────────────────────────────┐
│      Data Layer (db.Querier)            │
│  - Database queries                     │
│  - SQL operations                       │
└─────────────────────────────────────────┘
```

**Commands Use Services**:

All TUI commands (`commands.go`) use service interfaces instead of `db.Querier`:

**Loader Commands (Task-38):**
```go
func LoadAreasCmd(areaSvc service.AreaServiceInterface) tea.Cmd {
    return func() tea.Msg {
        areas, err := areaSvc.List(context.Background())
        if err != nil {
            return AreasLoadedMsg{Err: err}
        }
        return AreasLoadedMsg{Areas: areas}
    }
}
```

**CRUD Commands (Task-39):**
```go
func CreateAreaCmd(areaSvc service.AreaServiceInterface, name string, color domain.Color) tea.Cmd {
    return func() tea.Msg {
        area, err := areaSvc.Create(context.Background(), name, color)
        if err != nil {
            return AreaCreatedMsg{Err: err}
        }
        return AreaCreatedMsg{Area: area}
    }
}

func UpdateAreaCmd(areaSvc service.AreaServiceInterface, id string, name string, color domain.Color) tea.Cmd {
    return func() tea.Msg {
        area, err := areaSvc.Update(context.Background(), id, name, color)
        if err != nil {
            return AreaUpdatedMsg{Err: err}
        }
        return AreaUpdatedMsg{Area: area}
    }
}
```

Similar patterns for: CreateSubareaCmd, CreateProjectCmd, CreateTaskCmd, DeleteAreaCmd, ReorderAreasCmd, LoadAreaStatsCmd

**Benefits of Service Layer**:
- **Separation of concerns**: TUI focuses on UI, services handle business logic
- **Testability**: Mock services enable isolated unit testing
- **Flexibility**: Service implementations can change without TUI modifications
- **Maintainability**: Clear boundaries make code easier to understand and modify

### The Elm Architecture

The TUI follows The Elm Architecture pattern:

1. **Model**: Application state stored in `Model` struct (model.go)
2. **Update**: State transitions in `Update()` function (app.go)
3. **View**: Rendering logic in `View()` function (app.go)
4. **Commands**: Asynchronous operations via `tea.Cmd` (commands.go)

## Features

### 1. Area Tabs

Top-level navigation showing all areas as clickable tabs.

**Implementation**: `views/tabs.go`

- Horizontal tab bar at top of screen
- Active tab highlighted with border
- Keyboard navigation with `[` and `]` keys
- Wraps around at boundaries

### 2. Three-Column Browser

**Columns**: Subareas | Projects | Tasks

**Implementation**: `views/columns.go`, `tree/` package

- Focus-aware borders (thick border = active column)
- Independent navigation per column
- Project tree with expand/collapse
- Synchronized selection state

#### Proportional Column Layout

The three-column browser uses weight-based proportional widths to optimize screen real estate:

**Column Widths**:
- **Subareas**: 25% of available width (minimum 20 characters)
- **Projects**: 25% of available width (minimum 20 characters)
- **Tasks**: 50% of available width (minimum 40 characters)

**Layout Algorithm**:

The layout calculation follows these steps:

1. **Calculate available width**: `totalWidth - gaps` (6 characters for 3 gaps between columns)
2. **Apply weight distribution**: Subareas=1, Projects=1, Tasks=2 (25/25/50 ratio)
3. **Enforce minimum widths**: Each column respects its minimum character constraint
4. **Handle narrow terminals**: Below 80 columns, columns may overlap (stacked layout planned for future)

**Example Width Calculations**:

| Terminal Width | Subareas | Projects | Tasks |
|---------------|----------|----------|-------|
| 80 cols | 20 chars | 20 chars | 40 chars |
| 120 cols | 28 chars | 28 chars | 58 chars |
| 160 cols | 38 chars | 38 chars | 78 chars |

**Responsive Behavior**:

- Column widths recalculate instantly on terminal resize
- No animation or transition delays
- Minimum width constraints prevent unusable narrow columns
- Text truncation ensures clean borders at all sizes

### Responsive Layout Modes

The three-column browser automatically switches between two layout modes based on terminal width:

#### Stacked Layout (width < 120 cols)

For narrow terminals, the layout stacks Subareas and Projects vertically:

```
┌──────────────┬──────┐
│  Subareas    │      │
│  (Top)       │      │
├──────────────┤ Tasks│
│  Projects    │      │
│  (Bottom)    │      │
└──────────────┴──────┘
```

**Column Widths**:
- **Left side (Subareas+Projects combined)**: 25% of width
- **Right side (Tasks)**: 75% of width

**Height Distribution**:
- Subareas and Projects share equal height
- Tasks column uses full available height

#### Side-by-Side Layout (width >= 120 cols)

For wide terminals, all three columns are side-by-side:

```
┌─────┬────────┬──────────┐
│Sub- │Projects│  Tasks   │
│areas│        │          │
└─────┴────────┴──────────┘
```

**Column Widths**: Proportional 25/25/50 split (see Proportional Column Layout section)

#### Layout Switching

- **Threshold**: 120 columns
- **Behavior**: Instant, no animation
- **Trigger**: Automatic on terminal resize

**Testing**: Resize terminal from 119→120 cols to see instant layout switch.

#### Text Truncation

To prevent text wrapping in bordered panels (BubbleTea Golden Rule #2), all text is automatically truncated with intelligent partial-content preservation:

**Implementation**: `views/columns.go` - `truncateString()` function

**Behavior**:
- **Partial content preservation**: Shows as many characters as possible before appending ellipsis (…)
- **Example**: "Very Long Project Name" with maxLen=15 → "Very Long Pro…" (not just "…")
- **Maximum text width**: `columnWidth - 4` (accounting for 2 border chars + 2 padding chars)
- **Applies to**: Column titles, item names, and all content lines

**Advanced Features**:
- **ANSI escape code preservation**: Colored text remains visible in truncated portion (e.g., `"\x1b[31mRed Text\x1b[0m"` preserves color codes)
- **Unicode-aware truncation**: Multi-byte characters (emojis 🎉, CJK characters 日本語) handled correctly without breaking character boundaries
- **Edge case handling**: Very narrow columns (maxLen ≤ 1) show first character + ellipsis for minimal context

**Implementation Details**:
- Uses rune-based iteration for proper Unicode handling
- Tracks ANSI escape state to preserve color/formatting codes
- Calculates visible character count separately from total byte length

This ensures:
- No horizontal scrolling needed
- Clean, aligned borders
- Maximum readability at all terminal sizes
- Users can differentiate between long values even in narrow columns

**Test Coverage**: Comprehensive tests in `columns_test.go` covering basic truncation, ANSI codes, Unicode characters, and edge cases (44 tests total).

### 3. Project Tree Navigation

Hierarchical display of projects and sub-projects with expand/collapse functionality.

**Implementation**: `tree/` package

- **node.go**: Tree node data structure
- **builder.go**: Builds tree from flat project list
- **navigation.go**: Handles up/down navigation with wrapping
- **renderer.go**: Renders tree with indentation and expand/collapse indicators
- **constants.go**: Tree styling constants and configuration

**Visual Design**:

The tree uses a modern, minimalist design with arrow indicators:

- **Arrow Indicators**: 
  - `▾` (down triangle) for expanded nodes
  - `▸` (right triangle) for collapsed nodes
  - No indicator for leaf nodes (projects without subprojects)
- **Indentation**: Simple 2-space indentation per depth level (no vertical connector lines)
- **Clean Appearance**: Minimalist design without box-drawing characters for reduced visual clutter

**Example Tree Rendering**:
```
▾ Project A
  Subproject A1
  ▸ Subproject A2
Project B
▾ Project C
  Subproject C1
```

**Navigation**:
- `j`/`↓`: Move down (wraps to top)
- `k`/`↑`: Move up (wraps to bottom)
- `Enter`/`Space`: Toggle expand/collapse

### 4. Help Modal

Context-sensitive help showing all keyboard shortcuts grouped by category.

**Implementation**: `help/help.go`, `help/styles.go`

**Categories**:
- **Navigation**: h/l/j/k, arrows, Tab, [/]
- **Actions**: a, Enter, Space
- **General**: q, ?, Ctrl+C

**Trigger**: Press `?` key

**Behavior**:
- Modal overlay on main view
- Grouped shortcuts with descriptions
- Close with `?`, `Escape`, or `q`

### 5. Toast Notifications

Non-intrusive notifications for errors and success messages.

**Implementation**: `toast/toast.go`, `toast/styles.go`

**Features**:
- Auto-dismiss after configurable duration
- Stack multiple toasts
- Error (red) and success (green) variants
- Smooth fade-in/fade-out animations

**Usage in Code**:
```go
// Show error toast
m.toasts = append(m.toasts, toast.NewError("Database error: connection failed"))

// Show success toast
m.toasts = append(m.toasts, toast.NewSuccess("Task created successfully"))
```

### 6. Quick-Add Modal

Context-aware modal for creating new items.

**Implementation**: `modal/modal.go`, `modal/styles.go`, `modal/validation.go`

**Behavior**:
- Press `a` to open
- Context determined by focused column:
  - Subarea column → Create subarea in current area
  - Projects column → Create project in current subarea
  - Tasks column → Create task in current project
- Shows parent context (e.g., "New Project in: Work Tasks")
- Input validation for non-empty title
- Press `Enter` to create, `Escape` to cancel

**Subproject Checkbox** (Projects Column Only):

When a project is selected in the Projects column, the modal displays a checkbox labeled "Add as subproject" below the input field:

- **Unchecked (default)**: Creates a root-level project under the currently selected subarea
- **Checked**: Creates a subproject under the currently selected project
- **Keyboard navigation**: 
  - `Tab`/`Shift+Tab`: Navigate between input field and checkbox
  - `Space`: Toggle checkbox when focused
  - `Enter`: Submit with current checkbox state

This feature provides explicit control over project hierarchy, allowing users to create both root projects and subprojects when a project is selected.

**Note**: The checkbox only appears when a project is selected in the Projects column. Behavior remains unchanged for Subareas and Tasks columns.

### 7. Task Completion Toggle

Quick keyboard shortcut to mark tasks as complete with visual feedback.

**Implementation**: `internal/tui/app.go`, `internal/tui/renderer.go`

**Behavior**:
- Press `x` to toggle task completion status (Tasks column only)
- Smart toggle logic:
  - `todo` → `done`
  - `in_progress` → `done`
  - `waiting` → `done`
  - `done` → `todo`
- Immediate optimistic UI update
- Database persistence in background
- Automatic rollback with error toast if database operation fails

**Visual Feedback**:

Completed tasks display with three visual indicators:

1. **Checkmark Icon**: `✓` prefix before task title
2. **Strikethrough Text**: Horizontal line through text
3. **Muted Color**: Dimmed text using theme's muted color

**Example**:
```
  Incomplete task
✓ Completed task (with strikethrough and muted color)
  In progress task
✓ Another completed task
```

**Error Handling**:
- Optimistic UI update happens immediately
- If database operation fails:
  - UI state reverts to original status
  - Error toast notification appears
  - User can retry the operation

**Theme Integration**:
- Completed task colors automatically adapt to terminal theme
- Light terminal: Muted gray text (`#9CA3AF`)
- Dark terminal: Slightly lighter gray text (`#6B7280`)
- See Theme System section for details

**Accessibility**:
- Multiple visual cues (icon + strikethrough + color) ensure visibility
- Works with both light and dark terminal themes
- No reliance on color alone for status indication

### 8. Grouped Task Display

Hierarchical task organization showing tasks grouped by subproject with expandable/collapsible groups.

**Implementation**: `internal/tui/renderer.go`

**Overview**:
When a parent project with subprojects is selected, tasks are displayed in a grouped format rather than a flat list. Each subproject becomes a collapsible group header showing its tasks indented underneath.

**Visual Design**:
```
  Direct task 1
  Direct task 2

▸ Backend API (3 tasks)
  Database work (2 tasks)
▾ Frontend (1 task)
    ✓ Implement login form
```

**Key Features**:

1. **Direct Tasks**: Tasks belonging directly to the selected project appear at the top without grouping
2. **Group Headers**: Each subproject displays as a header with:
   - Expand/collapse icon: `▸` (collapsed) or `▾` (expanded)
   - Subproject name
   - Task count with proper pluralization
   - Dimmed styling using theme's secondary color
3. **Indentation**: Tasks under groups are indented 2 spaces deeper than direct tasks
4. **Text Truncation**: Long task titles are truncated with `…` to prevent wrapping
5. **Selection State**: Works with existing task navigation and completion toggle

**Implementation Details**:

The rendering logic uses the `GroupedTasks` domain model:

```go
type GroupedTasks struct {
    DirectTasks []Task        // Tasks from the selected project
    Groups      []TaskGroup   // Subproject groups
    TotalCount  int           // Total tasks across all groups
}

type TaskGroup struct {
    ProjectID   string
    ProjectName string
    Tasks       []Task
    IsExpanded  bool
}
```

**Rendering Algorithm**:
1. Render direct tasks (no header, minimal indentation)
2. Add blank line separator if both direct and grouped tasks exist
3. For each subproject group:
   - Render group header with icon and task count
   - If expanded: render indented tasks with proper styling

**Text Truncation**:
- Calculates available width based on column size
- Accounts for indentation, prefix (✓ or spaces), and borders
- Truncates with `…` character to prevent horizontal scrolling
- Preserves maximum content visibility

**Theme Integration**:
- Group headers use `theme.Secondary` color (dimmed appearance)
- No reverse highlighting on headers (subtle visual distinction)
- Task styling follows existing theme rules (completed, selected, etc.)

**Performance**:
- O(n) rendering where n = total visible tasks
- No nested loops or complex calculations
- String builder pattern for efficiency

### 9. Footer with Quick Reference

Persistent footer showing common keyboard shortcuts.

**Display**: `h/l: columns | j/k: nav | a: add | x: toggle | ?: help | q: quit`

**Implementation**: Footer rendered in main `View()` function

### 10. Confirmation Modal

Reusable confirmation dialog for destructive actions (delete operations).

**Implementation**: `confirmmodal/modal.go`, `confirmmodal/styles.go`

**Features**:
- Centered overlay modal using lipgloss.Place
- Warning styling with theme.Error color (red border/text)
- Displays entity type and item name
- Long names truncated with ellipsis (max 40 chars)
- Keyboard shortcuts: `y` (confirm), `n`/`Escape` (cancel)

**Entity Types Supported**:
- `EntityTypeSubarea`: Delete subarea confirmation
- `EntityTypeProject`: Delete project confirmation
- `EntityTypeTask`: Delete task confirmation

**Message Types**:
- `ConfirmMsg`: Returned on `y` key press, contains EntityType, EntityID, EntityName
- `CancelMsg`: Returned on `n` or `Escape` key press

**Usage**:
```go
// Create confirmation modal
modal := confirmmodal.New("Project Name", confirmmodal.EntityTypeProject, "proj-123")

// Handle messages in Update()
case confirmmodal.ConfirmMsg:
    // Execute delete operation
    return m, DeleteProjectCmd(m.projectSvc, msg.EntityID)
    
case confirmmodal.CancelMsg:
    // Close modal, no action taken
    m.confirmModal = nil
    m.isConfirmModalOpen = false
```

 **Theme Integration**:
- Border uses `theme.Default.Error` for destructive action warning
- Title uses bold + error color
- Hint text uses muted color for keyboard hints
- Automatic adaptation to light/dark terminal themes

**Integration Points**:
- Model needs: `confirmModal *confirmmodal.Modal`, `isConfirmModalOpen bool`
- Open trigger: `d` key in tree navigation (Task-68.3)
- Close on: ConfirmMsg, CancelMsg
- Use with: DeleteSubareaCmd, DeleteProjectCmd, DeleteTaskCmd

**Supported Entity Types**:
- Subareas: Deletes subarea and all its within it
- Projects: Deletes project with CASCADE (soft-deletes child projects and tasks recursively)
- Tasks: Deletes individual tasks (no cascade)

**Cascade Delete Behavior**:
When deleting a project that the service:
   1. Marks project as soft-deleted
   2. Finds all child projects and soft-deletes them
   3. Finds all tasks belonging to those child projects and soft-deletes them
   4. Updates parent references: remaining children now point to root projects
   5. Saves deletion to database

**Cascade Soft Delete Service**: Implemented in `internal/service/project_service.go`

```go
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, projectID string) error {
    tx, err := s.db.Begin()
    if err != nil {
        tx.Rollback()
        return nil, domain.NewDatabaseError("SoftDeleteWithCascade", err)
    }
    
    // Mark parent as soft-deleted
    err := s.db.MarkParentSoftDeleted(ctx, projectID)
    if err != nil {
        return nil, domain.NewDatabaseError("MarkParentSoftDeleted", err)
    }
    
    // Find and delete child projects recursively
    childProjects, err := s.db.GetProjectsByParentID(ctx, projectID)
    if err != nil {
        return nil, domain.NewDatabaseError("GetProjectsByParentID", err)
    }
    
    for _, childProject := range childProjects {
        if err := s.softDeleteWithCascade(ctx, childProject.ID); err != nil {
            return err
        }
    }
    
    // Find and delete tasks recursively
    childTasks, err := s.db.GetTasksByProjectID(ctx, projectID)
    if err != nil {
        return nil, domain.NewDatabaseError("GetTasksByProjectID", err)
    }
    
    for _, childTask := range childTasks {
        if err := s.taskSvc.SoftDelete(ctx, childTask.ID); err != nil {
            return err
        }
    }
    
    // Mark parent as soft-deleted
    err := s.db.MarkParentSoftDeleted(ctx, projectID)
    if err != nil {
        return nil, domain.NewDatabaseError("MarkParentSoftDeleted", err)
    }
    
    return nil
}
```

This ensures:
- **Referential integrity**: No orphaned tasks or projects remain in the database
- **Consistent UX**: Users see confirmation before deletion
- **Safety**: Accidental deletion prevented by requiring explicit confirmation

## Keyboard Shortcuts

### Navigation

| Key | Action | Description |
|-----|--------|-------------|
| `h`, `←` | Focus Left | Move to previous column (wraps right-to-left) |
| `l`, `→` | Focus Right | Move to next column (wraps left-to-right) |
| `Tab` | Cycle Focus | Move through columns in order |
| `j`, `↓` | Navigate Down | Move down in current column (wraps to top) |
| `k`, `↑` | Navigate Up | Move up in current column (wraps to bottom) |
| `[` | Previous Area | Switch to previous area tab (wraps to last) |
| `]` | Next Area | Switch to next area tab (wraps to first) |

### Actions

| Key | Action | Description |
|-----|--------|-------------|
| `Enter`, `Space` | Toggle Expand/Collapse | Expand or collapse project tree nodes |
| `a` | Quick Add | Open modal to create new item |
| `x` | Toggle Task Completion | Mark task as done/undone (Tasks column only) |
| `d` | Delete | Delete selected item (shows confirmation modal) |
| `Tab`, `Shift+Tab` | Navigate Modal | In quick-add modal: cycle between input and checkbox (when visible) |
| `Space` | Toggle Checkbox | In quick-add modal: toggle checkbox when focused |

### General

| Key | Action | Description |
|-----|--------|-------------|
| `Space` | Command Menu | Open command menu when no modal is open (LazyVim-style which-key) |
| `?` | Help | Show help modal with all shortcuts |
| `q`, `Ctrl+C` | Quit | Exit the TUI |
| `Escape` | Cancel/Close | Close modal or cancel operation |

## State Management

### Model Structure

```go
type Model struct {
    // Service layer dependencies (injected)
    areaSvc     service.AreaServiceInterface
    subareaSvc  service.SubareaServiceInterface
    projectSvc  service.ProjectServiceInterface
    taskSvc     service.TaskServiceInterface
    
    // Core data
    areas    []domain.Area
    subareas []domain.Subarea
    projects []domain.Project
    tasks    []domain.Task
    
    // UI state
    focus           FocusColumn
    selectedTab     int
    selectedAreaIndex    int
    selectedSubareaIndex int
    selectedProjectIndex int
    selectedTaskIndex    int
    
    // Tree state
    projectTree       *tree.TreeNode
    areaStates        map[string]*AreaState
    selectedProjectID string
    
    // Loading states
    isLoadingAreas    bool
    isLoadingSubareas bool
    isLoadingProjects bool
    isLoadingTasks    bool
    
    // Modal state
    modal         *modal.Modal
    isModalOpen   bool
    areaModal     *areamodal.Modal
    isAreaModalOpen bool
    spaceMenu     *spacemenu.SpaceMenu
    isSpaceMenuOpen bool
    helpModal     *help.HelpModal
    isHelpOpen    bool
    
    // Toasts
    toasts []toast.Toast
    
    // Terminal size
    width  int
    height int
    ready  bool
}
```

### Dependency Injection Pattern

The TUI uses **constructor injection** to receive service dependencies:

```go
// InitialModel creates a new Model with injected services
func InitialModel(
    areaSvc service.AreaServiceInterface,
    subareaSvc service.SubareaServiceInterface,
    projectSvc service.ProjectServiceInterface,
    taskSvc service.TaskServiceInterface,
) Model {
    return Model{
        areaSvc:    areaSvc,
        subareaSvc: subareaSvc,
        projectSvc: projectSvc,
        taskSvc:    taskSvc,
        // ... initialize other fields
    }
}
```

This pattern enables:
- Easy testing with mock services
- Clear documentation of dependencies
- Decoupling from database layer

### Focus States

```go
type FocusColumn int

const (
    FocusSubareas FocusColumn = iota
    FocusProjects
    FocusTasks
)
```

### Key Handling Flow

1. User presses key
2. `Update()` receives `tea.KeyMsg`
3. Switch on key type:
   - Global keys (`q`, `?`) handled first
   - Modal-specific keys if modal open
   - Column-specific navigation keys
   - Action keys (`a`, `Enter`)
4. Return new model and optional commands

## Testing

### Test Coverage

Comprehensive test suite covering:

1. **Unit Tests**:
   - `*_test.go` files alongside implementation
   - Test individual components in isolation
   - Use mock services for predictable testing
   - Located in `internal/tui/` directory

2. **Mock Services** (`internal/tui/mocks/`):
   - Mock implementations of all 4 service interfaces
   - Func-field pattern for maximum flexibility
   - Helper functions for common setup scenarios
   - See [mocks/README.md](../internal/tui/mocks/README.md) for details

3. **Integration Tests** (`integration_test.go`):
   - End-to-end user flows
   - Keyboard input sequences
   - State transitions
   - Modal interactions

4. **Navigation Tests** (`navigation_test.go`):
   - Column focus switching
   - Tab navigation
   - Tree navigation with wrapping

### Testing with Mocks

All tests use mock services instead of database connections:

```go
import "github.com/marekbrze/dopadone/internal/tui/mocks"

func TestLoadAreas(t *testing.T) {
    // Create mock services
    m := mocks.NewMockServices()
    
    // Configure mock behavior
    expectedAreas := []domain.Area{
        {ID: "area-1", Name: "Work"},
        {ID: "area-2", Name: "Personal"},
    }
    m.SetupMockAreaSuccess(expectedAreas)
    
    // Create model with mocked services
    model := InitialModel(
        m.AreaSvc,
        m.SubareaSvc,
        m.ProjectSvc,
        m.TaskSvc,
    )
    
    // Test logic here...
}
```

### Running Tests

```bash
# Run all TUI tests
go test ./internal/tui/... -v

# Run specific test
go test ./internal/tui/... -v -run TestHelpModal

# Run with coverage
go test ./internal/tui/... -cover

# Run with race detector
go test -race ./internal/tui/...
```

### Test Patterns

**Success Scenario**:
```go
m.SetupMockAreaSuccess(areas)
m.SetupMockProjectSuccess(projects)
```

**Error Scenario**:
```go
m.SetupMockAreaError(errors.New("database connection failed"))
```

**Custom Behavior**:
```go
m.AreaSvc.ListFunc = func(ctx context.Context) ([]domain.Area, error) {
    // Custom logic for specific test case
    return []domain.Area{{ID: "1", Name: "Test"}}, nil
}
```

### Manual Verification Checklist

For each release, manually verify:

- [ ] All keyboard shortcuts work as documented
- [ ] Tab navigation wraps correctly
- [ ] Column focus indicators are visible
- [ ] Tree expand/collapse works
- [ ] Help modal shows all shortcuts
- [ ] Toast notifications appear and auto-dismiss
- [ ] All keyboard shortcuts work as documented
- [ ] Tab navigation wraps correctly
- [ ] Column focus indicators are visible
- [ ] Tree expand/collapse works
- [ ] Help modal shows all shortcuts
- [ ] Toast notifications appear and auto-dismiss

- [ ] Quick-add modal creates items in correct context

## Delete Functionality

- [ ] Press `d` on selected item opens confirmation modal
- [ ] Press `y` confirms deletion
- [ ] Press `n` or `Escape` cancels deletion
- [ ] Verify success toast appears
- [ ] Press `d` on empty column does does nothing (no-op)
- [ ] Verify cascade delete works for projects with subprojects
    - Deleting a project also deletes all child projects and their tasks
    - Verify error toast appears on failure

- [ ] Verify footer shows `d: delete` shortcut

The TUI includes a comprehensive theming system that automatically adapts colors based on the terminal's background color, ensuring readability across different terminal themes.

### Overview

The theme system uses [lipgloss.AdaptiveColor](https://github.com/charmbracelet/lipgloss) to define colors that automatically switch between light and dark variants based on the terminal's color scheme. This eliminates the need for hardcoded ANSI color codes and ensures UI elements remain readable in any terminal theme.

**Key Features**:
- **Semantic Color Roles**: Colors defined by purpose (primary, secondary, success, error, warning, muted)
- **Automatic Adaptation**: Colors switch automatically based on terminal background
- **Configurable**: Theme selection via `backlog/config.yml`
- **Zero Hardcoded Colors**: All 48 previously hardcoded colors replaced with theme references

### Theme Package

**Location**: `internal/tui/theme/`

**Components**:
- `theme.go`: Core theme definitions and semantic color roles
- `loader.go`: Configuration loader and theme mode selection
- `theme_test.go`: Comprehensive theme tests

### Color Theme Structure

```go
type ColorTheme struct {
    Primary    lipgloss.AdaptiveColor  // Primary UI color (tabs, active borders)
    Secondary  lipgloss.AdaptiveColor  // Secondary text and accents
    Success    lipgloss.AdaptiveColor  // Success states and confirmations
    Error      lipgloss.AdaptiveColor  // Error states and warnings
    Warning    lipgloss.AdaptiveColor  // Warning messages
    Muted      lipgloss.AdaptiveColor  // Muted/disabled text
    Dimmed     lipgloss.AdaptiveColor  // Dimmed borders and separators
    Background lipgloss.AdaptiveColor  // Background color
    Foreground lipgloss.AdaptiveColor  // Default text color
}
```

### Semantic Color Methods

The `ColorTheme` struct provides helper methods for common UI elements:

```go
theme := theme.Default

// Tab styling
theme.TabActiveBackground()      // Active tab background
theme.TabActiveForeground()      // Active tab text
theme.TabInactiveBackground()    // Inactive tab background
theme.TabInactiveForeground()    // Inactive tab text

// Column styling
theme.ColumnFocusedBorder()      // Focused column border
theme.ColumnUnfocusedBorder()    // Unfocused column border
theme.ColumnHeader()             // Column header text

// General UI elements
theme.EmptyText()                // Empty state text
theme.FooterForeground()         // Footer text color
theme.FooterBackground()         // Footer background color
```

### Theme Modes

Three theme modes are supported:

1. **auto** (default): Automatically adapts colors based on terminal background
2. **light**: Forces light theme colors
3. **dark**: Forces dark theme colors

### Configuration

Theme mode is configured in `backlog/config.yml`:

```yaml
project_name: "adhd-coach-v2"
default_status: "To Do"
# ... other config ...
theme: "auto"  # Options: "auto", "light", "dark"
```

### Theme Loading

Themes are loaded during TUI initialization:

```go
// In internal/tui/tui.go
func New(
    areaSvc service.AreaServiceInterface,
    subareaSvc service.SubareaServiceInterface,
    projectSvc service.ProjectServiceInterface,
    taskSvc service.TaskServiceInterface,
) tea.Model {
    // Load theme from config
    theme, err := theme.LoadTheme("backlog/config.yml")
    if err != nil {
        theme = theme.Default
    }
    
    return Model{
        // ... services ...
        theme: theme,
        // ... other fields ...
    }
}
```

### Default Theme Colors

The default theme uses carefully selected colors for optimal readability:

**Light Terminal Background**:
- Primary: `#0066CC` (Blue)
- Secondary: `#6B7280` (Gray)
- Success: `#059669` (Green)
- Error: `#DC2626` (Red)
- Warning: `#D97706` (Amber)
- Muted: `#9CA3AF` (Light Gray)

**Dark Terminal Background**:
- Primary: `#4D9FFF` (Light Blue)
- Secondary: `#9CA3AF` (Light Gray)
- Success: `#10B981` (Light Green)
- Error: `#EF4444` (Light Red)
- Warning: `#F59E0B` (Light Amber)
- Muted: `#6B7280` (Gray)

### Implementation Pattern

All UI components access theme colors through the `Model.theme` field:

```go
// Before (hardcoded):
activeTabStyle := lipgloss.NewStyle().
    Background(lipgloss.Color("#0066CC")).
    Foreground(lipgloss.Color("#FFFFFF"))

// After (theme-aware):
activeTabStyle := lipgloss.NewStyle().
    Background(m.theme.TabActiveBackground()).
    Foreground(m.theme.TabActiveForeground())
```

### Component Integration

All TUI components use the theme system:

**Views (`views/styles.go`)**:
- Tab active/inactive styles
- Column borders and headers
- Empty state text

**Modals (`modal/styles.go`)**:
- Modal borders and backgrounds
- Input field styling
- Button styles

**Toasts (`toast/styles.go`)**:
- Error toast styling
- Success toast styling

**Help (`help/styles.go`)**:
- Help modal borders
- Category headers
- Shortcut text

**Tree (`tree/renderer.go`)**:
- Tree node styling
- Expand/collapse indicators

**Footer (`renderer_footer.go`)**:
- Footer text and background

### Testing

Theme system includes comprehensive tests:

```bash
# Run theme tests
go test ./internal/tui/theme/... -v

# Check coverage
go test ./internal/tui/theme/... -cover
```

**Test Coverage**: 60%+ coverage including:
- Theme loading from config
- Theme mode selection (auto/light/dark)
- Color theme validation
- Default theme values

### Benefits

1. **Readability**: UI elements automatically adapt to terminal theme
2. **Maintainability**: Centralized color management, no scattered hardcoded values
3. **Flexibility**: Easy to add new themes or modify existing ones
4. **User Control**: Users can override automatic detection with manual theme selection
5. **Consistency**: All components use semantic color roles for uniform appearance

### Future Enhancements

Potential theme system improvements:
- Custom color palettes via config
- Multiple theme presets (beyond default)
- Theme export/import functionality
- Per-area theme customization

## Styling

### Lipgloss Usage

All styling uses [Lipgloss](https://github.com/charmbracelet/lipgloss) for consistent, declarative styling.

**Color Palette**:
- Active border: Bold white
- Inactive border: Dim gray
- Selected item: Cyan background
- Error: Red (#FF6B6B)
- Success: Green (#4ECDC4)
- Help modal: Purple accent (#B794F6)

**Style Files**:
- `help/styles.go`: Help modal styles
- `toast/styles.go`: Toast notification styles
- `modal/styles.go`: Quick-add modal styles
- `views/styles.go`: Main view styles
- `tree/constants.go`: Tree rendering characters and indicators
- `theme/theme.go`: Theme definitions and semantic color roles

### Tree Styling

The project tree uses a customizable styling system defined in `tree/constants.go`:

**Character Constants**:
- `TreeIndent`: 2-space indentation per depth level
- `ExpandedIcon`: `▾` (down triangle) for expanded nodes with children
- `CollapsedIcon`: `▸` (right triangle) for collapsed nodes with children
- No indicator for leaf nodes

**Customization**:

The `TreeStyle` struct allows custom tree rendering characters:

```go
type TreeStyle struct {
    Branch   string // Non-last child prefix
    Last     string // Last child prefix
    Vertical string // Vertical continuation
    Indent   string // Depth indentation
}
```

To customize the tree appearance, create a custom `TreeStyle` and pass it to the renderer:

```go
customStyle := tree.TreeStyle{
    Branch:   "├─ ",
    Last:     "└─ ",
    Vertical: "│  ",
    Indent:   "  ",
}
renderer := tree.NewRenderer()
renderer.SetStyle(customStyle)
```

**Current Design Rationale**:
- Simple indentation (no vertical lines) reduces visual clutter
- Arrow indicators (▸/▾) provide clear expand/collapse state
- 2-space indentation ensures proper alignment at all depths
- Unicode characters provide modern, clean appearance

### Responsive Design

The TUI adapts to terminal size:
- Minimum supported: 80x24
- Columns resize proportionally
- Text truncation with ellipsis for long names
- Proper wrapping for narrow terminals

## Error Handling

The TUI implements comprehensive error handling across three dimensions: **error state tracking**, **error rendering**, and **user-friendly messaging**.

### Error State Management

The Model tracks errors for each data loading operation:

```go
// internal/tui/model.go

type Model struct {
    // ... existing fields ...
    
    // Error tracking for each column
    areaLoadError    error
    subareaLoadError error
    projectLoadError error
    taskLoadError    error
}
```

**Benefits**:
- Errors persist in model state for rendering
- Allows retry mechanisms
- Clear separation between loading state and error state
- Type-safe error checking with domain helpers

**Error Clearing**:
```go
func (m *Model) ClearErrors() {
    m.areaLoadError = nil
    m.subareaLoadError = nil
    m.projectLoadError = nil
    m.taskLoadError = nil
}
```

### Error Handling in Message Handlers

Handlers store errors and show user-friendly messages:

```go
// internal/tui/handlers.go

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
```

### User-Friendly Error Messages

The TUI maps technical errors to user-friendly messages:

```go
func (m *Model) formatUserError(err error) string {
    // Check specific error types first
    if domain.IsNotFound(err) {
        return "Resource not found"
    }
    
    // Check context errors
    if errors.Is(err, context.Canceled) {
        return "Operation cancelled"
    }
    if errors.Is(err, context.DeadlineExceeded) {
        return "Loading took too long. Please try again."
    }
    
    // Check database errors
    if domain.IsDatabaseError(err) || 
       strings.Contains(err.Error(), "database") || 
       strings.Contains(err.Error(), "sql") {
        return "Unable to load data. Please restart the application."
    }
    
    // Generic fallback
    return fmt.Sprintf("Error: %v", err)
}
```

**Message Constants** (internal/tui/constants.go):
```go
const (
    ErrMsgDatabase  = "Unable to load data. Please restart the application."
    ErrMsgTimeout   = "Loading took too long. Please try again."
    ErrMsgCancelled = "Operation cancelled"
    ErrMsgNotFound  = "Resource not found"
)
```

### Error Rendering in Views

Renderers check error state before rendering content:

```go
// internal/tui/renderer.go

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
```

**Error Display**:
```
┌─────────────────────────────┐
│ Projects    Subareas  Tasks │
│                        [Col]│
│ ▾ Project A                │
│   Subproject A1            │
│                            │
│                            │
│ ✗ Unable to load data.     │
│   Please restart the       │
│   application.             │
└─────────────────────────────┘
```

### Error Type Checking

Use domain helper functions for type-safe error checking:

```go
// ✅ Good: Use domain helpers
if domain.IsNotFound(err) {
    // Handle not found
} else if domain.IsDatabaseError(err) {
    // Handle database error
}

// ❌ Bad: String comparison
if strings.Contains(err.Error(), "not found") {
    // Fragile and error-prone
}
```

### Empty State vs Error State

**Empty State**: Valid condition (no data)
```go
func (m *Model) renderEmptyTasks() string {
    emptyStyle := lipgloss.NewStyle().
        Foreground(m.theme.Dimmed).
        PaddingLeft(2)
    
    msg := "No tasks in this project"
    
    if len(m.groupedTasks.Groups) > 0 {
        msg = "No tasks in this project or its subprojects"
    }
    
    return emptyStyle.Render(msg)
}
```

**Error State**: Something went wrong
```go
func (m *Model) renderError(err error, context string) string {
    // User-friendly error with red styling
    return errorStyle.Render("✗ " + userMsg)
}
```

**Visual Difference**:
```
Empty State (dimmed):
  No tasks in this project

Error State (red):
  ✗ Unable to load data. Please restart the application.
```

### Error Recovery Strategies

**1. Retry on transient errors**:
```go
// User can press 'r' to retry loading
case key.Matches(msg, m.keys.retry):
    if m.taskLoadError != nil {
        return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
    }
```

**2. Clear errors on navigation**:
```go
func (m *Model) handleProjectSelected(msg ProjectSelectedMsg) {
    m.ClearErrors()  // Clear previous errors
    m.selectedProjectID = msg.ProjectID
    return m, LoadTasksCmd(m.taskSvc, msg.ProjectID)
}
```

**3. Graceful degradation**:
```go
// If tasks fail to load, still show project tree
if m.taskLoadError != nil {
    // Show error in tasks column
    // Keep projects column functional
}
```

### Database Errors

All database errors are caught and displayed as toast notifications:

```go
result, err := db.CreateProject(...)
if err != nil {
    return m, func() tea.Msg {
        return toast.NewError("Failed to create project: " + err.Error())
    }
}
```

### User Input Validation

The quick-add modal validates input:
- Non-empty title required
- Max length enforcement
- Whitespace trimming

### Error Handling Best Practices

1. **Never expose technical details**: Map to user-friendly messages
2. **Use error state tracking**: Store errors in model for rendering
3. **Type-safe checking**: Use domain.IsNotFound(), not string comparison
4. **Clear errors appropriately**: Clear on successful loads or navigation
5. **Provide recovery options**: Allow retry for transient errors
6. **Graceful degradation**: Keep other columns functional when one fails
7. **Visual distinction**: Use different styling for empty vs error states

### Testing Error Handling

Test error scenarios with mock services:

```go
func TestTaskLoadingError(t *testing.T) {
    mockSvc := &mockTaskService{
        getError: domain.NewDatabaseError("ListTasks", errors.New("connection failed")),
    }
    
    model := New(mockSvc, ...)
    model.selectedProjectID = "proj-1"
    
    // Trigger load
    msg := LoadTasksMsg{ProjectID: "proj-1"}
    result, _ := model.Update(msg)
    
    // Verify error state
    if result.taskLoadError == nil {
        t.Error("expected taskLoadError to be set")
    }
    
    // Verify user message
    rendered := result.RenderTasks()
    if !strings.Contains(rendered, "Unable to load data") {
        t.Error("expected user-friendly error message")
    }
}
```

## Performance Considerations

### Efficient Rendering

- Only re-render visible portions of tree
- Lazy loading of task lists
- Minimize state updates

### Memory Management

- Reuse buffers where possible
- Limit toast history
- Clean up expanded nodes cache

## Future Enhancements

Potential improvements for future versions:

1. **Search/Filter**: Add `/` key for searching projects/tasks
2. **Bulk Operations**: Multi-select with visual indication
3. **Drag & Drop**: Reorder items with keyboard
4. **Custom Themes**: User-configurable color schemes
5. **Split View**: Side-by-side task detail view
6. **Undo/Redo**: Command history with `u` and `Ctrl+R`
7. **Export**: Save current view to file
8. **Keyboard Macros**: Record and replay key sequences

## Troubleshooting

### Common Issues

1. **TUI doesn't render correctly**:
   - Check terminal supports true color
   - Verify minimum size (80x24)
   - Try setting `TERM=xterm-256color`

2. **Keys not responding**:
   - Check for keyboard layout issues
   - Verify terminal doesn't intercept keys
   - Try different terminal emulator

3. **Performance lag**:
   - Reduce number of expanded tree nodes
   - Check database query performance
   - Enable debug logging

### Debug Mode

Enable debug logging:
```bash
export TUI_DEBUG=1
dopa tui
```

## Related Documentation

- [README.md](../README.md) - User-facing TUI documentation
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
