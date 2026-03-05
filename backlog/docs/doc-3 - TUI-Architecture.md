---
id: doc-3
title: TUI Architecture
type: technical
created_date: '2026-03-03'
---

# TUI Architecture

## Overview

This document describes the Terminal User Interface (TUI) architecture for ProjectDB, built using the bubbletea framework with Model-Update-View (MVU) pattern.

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Framework | bubbletea v1.3.10 | TUI framework with MVU pattern |
| Components | bubbles v1.0.0 | Pre-built UI components |
| Styling | lipgloss v1.1.0 | Terminal styling and layouts |
| Architecture | Clean Architecture | Delivery layer depends on domain |

## Architecture Pattern

The TUI follows the **Model-Update-View (MVU)** pattern:

```
┌──────────────┐
│    Model     │  Application state
└──────┬───────┘
       │
       ▼
┌──────────────┐      ┌──────────────┐
│    Update    │◄────►│   Commands   │
└──────┬───────┘      └──────────────┘
       │
       ▼
┌──────────────┐
│     View     │  Render to string
└──────────────┘
```

### Model-Update-View Flow

1. **Model**: Contains all application state (focus, selections, data)
2. **Update**: Handles messages (keyboard, resize) and updates model
3. **View**: Renders model to string (no side effects)
4. **Commands**: Async operations (future: data loading, etc.)

## Package Structure

```
internal/tui/
├── app.go              # Main Model struct and tea.Model implementation
├── model.go            # FocusColumn enum and state definitions (AreaState)
├── tui.go              # Exported New() function for program creation
├── messages.go         # Message types for async data loading
├── commands.go         # Loader and CRUD commands using service layer interfaces
                          # Loader commands (Task-38):
                          # - LoadAreasCmd: AreaServiceInterface.List()
                          # - LoadSubareasCmd: SubareaServiceInterface.ListByArea()
                          # - LoadProjectsCmd: ProjectServiceInterface.ListBySubareaRecursive()
                          # - LoadTasksCmd: TaskServiceInterface.ListByProject()
                          # CRUD commands (Task-39):
                          # - CreateSubareaCmd: SubareaServiceInterface.Create()
                          # - CreateProjectCmd: ProjectServiceInterface.Create()
                          # - CreateTaskCmd: TaskServiceInterface.Create()
                          # - CreateAreaCmd: AreaServiceInterface.Create()
                          # - UpdateAreaCmd: AreaServiceInterface.Update()
                          # - DeleteAreaCmd: AreaServiceInterface.SoftDelete/HardDelete()
                          # - ReorderAreasCmd: AreaServiceInterface.ReorderAll()
                          # - LoadAreaStatsCmd: AreaServiceInterface.GetStats()
├── constants.go        # Named constants for TUI (no magic numbers)
├── app_test.go         # Unit tests for app functionality
├── messages_test.go    # Unit tests for message types
├── commands_test.go    # Unit tests for loader commands with mocks
├── navigation_test.go  # Unit tests for navigation (Task-18)
├── state_test.go       # Unit tests for state persistence (Task-18)
├── integration_test.go # Integration tests for navigation flow (Task-18)
├── tree/               # Tree rendering package (Task-20)
│   ├── node.go         # TreeNode model with expand/collapse
│   ├── builder.go      # Build hierarchical tree from flat list
│   ├── renderer.go     # Lipgloss-styled tree rendering
│   ├── navigation.go   # Tree navigation helpers
│   ├── constants.go    # Tree characters and styling constants
│   └── *_test.go       # Comprehensive unit tests (51 tests, 95% coverage)
└── views/
    ├── tabs.go         # Area tabs component
    ├── columns.go      # 3-column layout component
    └── styles.go       # Shared lipgloss styles
```

## Core Components

### 1. Model Struct

```go
type Model struct {
    // Service layer interfaces (Task-38)
    areaSvc     service.AreaServiceInterface
    subareaSvc  service.SubareaServiceInterface
    projectSvc  service.ProjectServiceInterface
    taskSvc     service.TaskServiceInterface
    
    focus       FocusColumn      // Current focused column
    width       int              // Terminal width
    height      int              // Terminal height
    ready       bool             // Lazy initialization flag
    tabs        []views.Tab      // Area tabs (placeholder)
    selectedTab int              // Selected area index
    
    // Loaded data
    areas       []domain.Area
    subareas    []domain.Subarea
    projects    []domain.Project
    tasks       []domain.Task
    
    // Selection tracking (simplified to int indices)
    selectedAreaIndex    int
    selectedSubareaIndex int
    selectedProjectIndex int
    selectedTaskIndex    int
    
    // Loading states
    isLoadingAreas    bool
    isLoadingSubareas bool
    isLoadingProjects bool
    isLoadingTasks    bool
    
    // UI components
    spinner     spinner.Model    // Loading indicator
    
    // State persistence (Task-18)
    areaStates  map[string]*AreaState  // Per-area state (keyed by area ID)
}

// AreaState tracks navigation state per area (Task-18)
type AreaState struct {
    SelectedSubareaIndex int            // Selected subarea index
    SelectedProjectIndex int            // Selected project index  
    SelectedTaskIndex    int            // Selected task index
    ExpandedProjects     map[string]bool // Project tree expansion state
}
```

### 2. FocusColumn Enum

```go
type FocusColumn int

const (
    FocusSubareas FocusColumn = iota
    FocusProjects
    FocusTasks
)
```

Navigation methods:
- `Prev()` - Move left with wrapping (Subareas → Tasks → Projects → Subareas)
- `Next()` - Move right with wrapping (Subareas → Projects → Tasks → Subareas)
- `String()` - Human-readable column name

### 3. View Components

#### Tabs (views/tabs.go)

Renders area tabs with visual highlighting:

```go
type Tab struct {
    Name     string
    ID       string
    IsActive bool
}

func TabsView(tabs []Tab, selectedIndex int) string
```

#### Columns (views/columns.go)

Renders 3-column layout with focus-aware borders:

```go
type Column struct {
    Title     string
    Content   string
    IsFocused bool
    Width     int
    Height    int
}

func ColumnView(col Column) string
func Layout(columns []Column, width, height int) string
func LayoutWithTabs(tabs string, columns []Column, width, height int) string
```

#### Styles (views/styles.go)

Shared lipgloss styles:
- `ActiveTabStyle` - Highlighted tab styling
- `InactiveTabStyle` - Dimmed tab styling
- `FocusedColumnStyle` - Thick border for active column
- `UnfocusedColumnStyle` - Normal border for inactive columns
- `ColumnHeaderStyle` - Bold headers
- `EmptyContentStyle` - Italic placeholder text

### 4. Tree Package (internal/tui/tree/)

The tree package provides hierarchical project display with unlimited nesting support.

#### TreeNode Model

```go
type TreeNode struct {
    ID         string
    Name       string
    Depth      int
    IsExpanded bool
    Children   []*TreeNode
    Parent     *TreeNode
    Data       interface{}  // Stores domain.Project
}

func (n *TreeNode) IsLeaf() bool
func (n *TreeNode) HasChildren() bool
func (n *TreeNode) ToggleExpanded()
```

#### Tree Builder

```go
func BuildFromProjects(projects []domain.Project) *TreeNode
```

Transforms flat project list to hierarchical tree:
- Separates roots (SubareaID != nil) from children
- Sorts by Position field
- Recursively attaches children to parents
- Handles orphans (logs warning and skips)
- Supports unlimited nesting depth

#### Tree Renderer

```go
func Render(root *TreeNode, selectedID string) string
```

Renders tree with lipgloss styling:
- Tree characters: ├─ └─ │ (constants, no magic strings)
- Indentation: 2 spaces per depth level
- Expand/collapse: [+]/[-] indicators
- Selected node: lipgloss highlighting
- Skips children of collapsed nodes

#### Navigation Helpers

```go
func GetNextVisibleNode(node *TreeNode) *TreeNode
func GetPrevVisibleNode(root, node *TreeNode) *TreeNode
func GetAllVisibleNodes(root *TreeNode) []*TreeNode
func FindNodeByID(root *TreeNode, id string) *TreeNode
func ExpandAll(root *TreeNode)
func CollapseAll(root *TreeNode)
```

Navigation respects collapsed state and skips hidden children.

#### Constants

```go
const (
    TreeIndent      = "  "
    TreeBranch      = "├─ "
    TreeLast        = "└─ "
    TreeVertical    = "│  "
    ExpandedIcon    = "[-]"
    CollapsedIcon   = "[+]"
)
```

All tree characters use named constants (no magic strings).

#### Test Coverage

- 51 unit tests, 95.0% coverage
- Tests cover: empty tree, flat list, nested (1-5+ levels), orphans, position ordering
- All exported functions have godoc comments
- Follows F.I.R.S.T. testing principles

## Event Handling

### Keyboard Messages

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "h", "left":
            m.focus = m.focus.Prev()
        case "l", "right":
            m.focus = m.focus.Next()
        case "tab":
            m.focus = m.focus.Next()
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.ready = true
    }
    return m, nil
}
```

### Message Types

| Message | Source | Purpose |
|---------|--------|---------|
| `tea.KeyMsg` | Keyboard input | Navigation, actions |
| `tea.WindowSizeMsg` | Terminal resize | Update dimensions |
| `tea.Quit` | Exit command | Clean shutdown |
| `spinner.TickMsg` | Spinner tick | Update loading animation |
| `LoadAreasMsg` | Data command | Trigger area loading |
| `AreasLoadedMsg` | Data result | Areas loaded from DB |
| `LoadSubareasMsg` | Data command | Trigger subarea loading |
| `SubareasLoadedMsg` | Data result | Subareas loaded from DB |
| `LoadProjectsMsg` | Data command | Trigger project loading |
| `ProjectsLoadedMsg` | Data result | Projects loaded from DB |
| `LoadTasksMsg` | Data command | Trigger task loading |
| `TasksLoadedMsg` | Data result | Tasks loaded from DB |

## Data Loading Architecture

### Async Operations Flow

```
User selects Area
       │
       ▼
LoadSubareasCmd() ──► Database Query
       │                      │
       │                      ▼
       │              SubareasLoadedMsg
       │                      │
       ▼                      ▼
Update() ◄────────────────────┘
       │
       ▼
Auto-select first Subarea
       │
       ▼
LoadProjectsCmd() ──► Database Query
       │                      │
       │                      ▼
       │              ProjectsLoadedMsg
       │                      │
       ▼                      ▼
Update() ◄────────────────────┘
       │
       ▼
Build tree from projects
       │
       ▼
Auto-select first Project
       │
       ▼
LoadTasksCmd() ──► Database Query
       │                      │
       │                      ▼
       │              TasksLoadedMsg
       │                      │
       ▼                      ▼
Update() ◄────────────────────┘
```

### Service Injection Pattern

```go
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
        spinner:    spinner.New(spinner.WithSpinner(spinner.Dot)),
        // ... other fields
    }
}
```

Clean architecture boundaries:
- TUI depends on service interfaces, not concrete implementations or database layer
- Domain entities have no framework dependencies
- Dependency Inversion Principle followed
- Services encapsulate business logic and return domain types directly

### Loader Commands

All loader functions use service layer interfaces instead of direct database access:

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

// LoadProjectsCmd uses hierarchical retrieval for nested projects
func LoadProjectsCmd(projectSvc service.ProjectServiceInterface, subareaID *string) tea.Cmd {
    return func() tea.Msg {
        var projects []domain.Project
        var err error
        
        if subareaID != nil {
            projects, err = projectSvc.ListBySubareaRecursive(context.Background(), *subareaID)
        } else {
            projects, err = projectSvc.ListAll(context.Background())
        }
        
        if err != nil {
            return ProjectsLoadedMsg{Err: err}
        }
        return ProjectsLoadedMsg{Projects: projects}
    }
}
```

Similar pattern for: LoadSubareasCmd, LoadTasksCmd

### CRUD Commands

All CRUD commands use service layer interfaces for data modifications:

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

func DeleteAreaCmd(areaSvc service.AreaServiceInterface, id string, hard bool) tea.Cmd {
    return func() tea.Msg {
        var err error
        if hard {
            err = areaSvc.HardDelete(context.Background(), id)
        } else {
            err = areaSvc.SoftDelete(context.Background(), id)
        }
        if err != nil {
            return AreaDeletedMsg{Err: err}
        }
        return AreaDeletedMsg{ID: id}
    }
}
```

Similar pattern for: CreateSubareaCmd, CreateProjectCmd, CreateTaskCmd, ReorderAreasCmd, LoadAreaStatsCmd

**Benefits**:
- Services handle business logic and validation
- Consistent error handling through service layer
- Easy mocking for tests with service interfaces
- Type-safe domain operations
- No converter layer needed (services return domain types)

**Benefits**:
- Services return domain types directly (no converter layer needed in commands)
- Consistent error handling through service layer
- Easy mocking for tests with service interfaces
- Business logic centralized in services

### Type Converters

Type converters exist in `internal/converter/` and are used by the service layer to convert DB types to domain types:

```go
func DbAreaToDomain(dbArea db.Area) domain.Area
func DbSubareaToDomain(dbSubarea db.Subarea) domain.Subarea
func DbProjectToDomain(dbProject db.Project) domain.Project
func DbTaskToDomain(dbTask db.Task) domain.Task
```

**TUI commands don't need converters** because services return domain types directly. This simplifies the TUI layer and keeps it focused on presentation logic.

### Loading State Management

```go
// Prevent duplicate loads
if m.isLoadingSubareas {
    return m, nil  // Skip if already loading
}

// Set loading flag before command
m.isLoadingSubareas = true
return m, LoadSubareasCmd(m.repo)
```

### Empty State Handling

Contextual empty state messages with keyboard hints:

```go
func (m Model) renderSubareas() string {
    if m.isLoadingSubareas {
        return spinnerView("Loading subareas...")
    }
    if len(m.subareas) == 0 {
        return "No subareas\n\nPress 'a' to add one"
    }
    // ... render subareas
}
```

### Cascade Loading

Data loads cascade from top to bottom:
1. App starts → LoadAreas
2. Area selected → LoadSubareas
3. Subarea selected → LoadProjects
4. Project selected → LoadTasks

Each step auto-selects the first item, triggering the next load.

## Navigation & State Persistence (Task-18)

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `j` / `↓` | Navigate down in current column (wrap to top at bottom) |
| `k` / `↑` | Navigate up in current column (wrap to bottom at top) |
| `Enter` / `Space` | Toggle expand/collapse for project tree nodes |
| `[` | Switch to previous area tab (wrap to last) |
| `]` | Switch to next area tab (wrap to first) |

### Navigation Helpers

All navigation functions are under 20 lines following SRP:

```go
func (m *Model) navigateUp()
func (m *Model) navigateDown()
func (m *Model) navigateTreeUp()
func (m *Model) navigateTreeDown()
func (m *Model) isEmpty() bool
```

### State Persistence

State is saved and restored per area when switching:

```go
type AreaState struct {
    SelectedSubareaIndex int
    SelectedProjectIndex int
    SelectedTaskIndex    int
    ExpandedProjects     map[string]bool
}

func (m *Model) getAreaState(areaID string) *AreaState
func (m *Model) saveCurrentAreaState()
func (m *Model) restoreAreaState(areaID string)
```

When user switches areas via `[`/`]`:
1. Current selections and tree expansion saved to `areaStates[oldAreaID]`
2. New area data loaded (subareas, projects, tasks)
3. Selections and tree expansion restored from `areaStates[newAreaID]`

### Visual Feedback

| Element | Styling |
|---------|---------|
| Selected items | Bold + Inverted colors |
| Active tab | Bold + Inverted background |
| Inactive tabs | Normal styling |
| Empty column navigation | No-op (no selection change) |

### Tree Navigation

Project tree navigation respects collapsed state:

```go
// Uses tree package helpers
func GetNextVisibleNode(node *TreeNode) *TreeNode
func GetPrevVisibleNode(root, node *TreeNode) *TreeNode
```

- Skips hidden children of collapsed nodes
- Wraps from last visible node to first
- Works with unlimited nesting depth

### Test Coverage

- Unit tests: navigation_test.go (boundary cases, wrap-around, empty columns)
- Unit tests: state_test.go (state persistence, area switching)
- Integration tests: integration_test.go (full navigation flow)
- Coverage: 82.4% for tui package

### Test Coverage

- MockQuerier for isolated testing
- Unit tests for messages, commands, converters, and app
- Test coverage: 85.5% (exceeds 85% requirement)
- All loader functions under 20 lines
- No magic numbers (all constants named)

## View Rendering

```go
func (m Model) View() string {
    if !m.ready {
        return "\n  Initializing..."
    }
    
    // Compose tabs
    tabs := views.TabsView(m.tabs, m.selectedTab)
    
    // Build columns with focus state
    columns := []views.Column{
        {Title: "Subareas", IsFocused: m.focus == FocusSubareas},
        {Title: "Projects", IsFocused: m.focus == FocusProjects},
        {Title: "Tasks", IsFocused: m.focus == FocusTasks},
    }
    
    // Render layout
    return views.LayoutWithTabs(tabs, columns, m.width, m.height)
}
```

**Important**: View() has no side effects - it only renders state.

## Clean Architecture Integration

```
cmd/projectdb/
    └── tui.go         CLI command (delivery layer)
         │
         ▼
internal/tui/
    ├── app.go         TUI application (delivery layer)
    ├── model.go       UI state
    └── views/         UI components
         │
         ▼
internal/domain/
    ├── area.go        Domain entities
    ├── project.go
    └── task.go
```

**Dependency Rule**: TUI (delivery) → Domain (core)
- ✅ TUI imports domain types (future: data loading)
- ❌ Domain never imports TUI

## Testing Strategy

### Unit Tests (tui_test.go)

Focus state machine tests:

```go
// Column transition tests
TestFocusColumnPrev()    // Wrapping left navigation
TestFocusColumnNext()    // Wrapping right navigation
TestFocusTransitionLeft() // Model method integration
TestFocusTransitionRight() // Model method integration

// State reachability tests
TestFocusStatesReachable() // All 3 columns reachable
TestFocusCycleTab()        // Tab cycles correctly
```

All tests: `go test ./internal/tui -v`

### Test Coverage Summary

| Package | Coverage | Test Count | Status |
|---------|----------|------------|--------|
| internal/tui | 82.4% | 40+ tests | ✅ Pass |
| internal/tui/tree | 95.0% | 51 tests | ✅ Pass |

Run all tests with coverage:
```bash
go test ./internal/tui/... -v -cover
```

### Integration Tests

Manual testing checklist:
- [ ] Tabs render with visual highlighting
- [ ] 3 columns display with borders
- [ ] h/l/arrows navigate with wrapping
- [ ] Tab cycles through columns
- [ ] j/k navigate within columns with wrap-around
- [ ] Enter/Space toggle tree expand/collapse
- [ ] `[` and `]` switch areas with wrapping
- [ ] Selection restored when switching areas
- [ ] Tree expansion persisted per area
- [ ] q/Ctrl+C exits cleanly
- [ ] Terminal resize adapts layout
- [ ] `projectdb tui` command works

## CLI Integration

```go
// cmd/projectdb/tui.go
var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Launch the TUI interface",
    Run: func(cmd *cobra.Command, args []string) {
        // Connect to database
        db, err := database.New(database.Config{
            Path: cfg.DatabasePath,
        })
        if err != nil {
            fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
            return
        }
        defer db.Close()
        
        // Create repository and TUI
        repo := db.New(db)
        p := tui.New(repo)
        
        if _, err := p.Run(); err != nil {
            fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
        }
    },
}
```

Launch: `projectdb tui`

## Implementation Status

### ✅ Completed: TUI 3-Column Browser with Quick-Add (Task-14)

Parent task split into 6 subtasks, all completed:
- Task-15 (14A Core): Framework, tabs, layout, column nav, exit
- Task-20 (16A Tree): Tree rendering package (superseded Task-16)
- Task-21 (16B Data): Data loading, spinner, empty states
- Task-18 (14C Nav): In-column nav, area switching, state persistence
- Task-19 (14D Modal): Quick-add modal with validation
- Task-17 (14E Polish): Help modal, toasts, footer, docs
- Task-38 (29D): Refactored load commands to use service layer interfaces

**Total**: 67 acceptance criteria across all tasks, 133 TUI tests passing
  - Comprehensive test coverage with mocked services

### ✅ Completed: Core TUI Framework (Task-15)
- 3-column layout with focus-aware borders
- Area tabs with visual highlighting
- Column navigation (h/l/arrows/tab)
- Clean architecture integration

### ✅ Completed: Tree Rendering Package (Task-20)
- Hierarchical project display with unlimited nesting
- Tree characters (├─ └─ │) with lipgloss styling
- Expand/collapse behavior with [+]/[-] indicators
- Navigation helpers for collapsed node handling
- 95% test coverage with 51 unit tests

### ✅ Completed: Data Loading & Integration (Task-21)
- Async data loading from database
- Cascade loading: Area → Subarea → Projects → Tasks
- Loading spinner with bubbles/spinner component
- Empty state messages with keyboard hints
- Service injection pattern (updated in Task-38)
- Type converters for DB → Domain (in service layer)
- 85.5% test coverage

### ✅ Completed: Navigation & State Persistence (Task-18)
- j/k/arrow keys for in-column navigation with wrap-around
- `[` and `]` keys for area switching with wrapping
- Enter/Space for tree expand/collapse toggle
- State persistence: selections and tree expansion saved per area
- Selected item styling: bold + inverted colors
- Active tab styling: bold + inverted background
- Navigation on empty columns is no-op
- Minimal scroll behavior for off-screen selections
- 82.4% test coverage (tui), 95.0% (tree)
- All functions under 20 lines following SRP

### ✅ Completed: Quick-Add Modal (Task-19)
- `a` key opens context-aware modal centered on screen
- Modal displays parent context (e.g., "New Project in: [Parent Name]")
- Single title input field with cursor focus and 255 char limit
- Enter key creates item in focused column context with validation
- Escape key closes modal without creating (no changes)
- Column refreshes and focuses newly created item after successful creation
- Inline error message displayed in modal for creation failures
- Input validation: title required (non-empty after trim), no newlines/control chars
- Modal width is responsive (40-60% of terminal width)
- Context-aware creation: Subareas→Area, Projects→Subarea, Tasks→Project
- Comprehensive test coverage: modal (17 tests), validation (30 tests), create commands (14 tests)
- All 127 TUI tests passing, no regressions

### ✅ Completed: Help, Errors & Polish (Task-17)
- `?` key opens help modal with all keyboard shortcuts
- Shortcuts grouped by category (Navigation, Actions, General)
- Toast notifications for database errors with auto-dismiss
- Footer with quick reference shortcuts
- README updated with TUI command usage
- Integration tests for key user flows
- All 62 ACs from tasks 15, 17-21 verified manually

### ✅ Completed: Refactor Load Commands to Use Services (Task-38)
- All 4 load commands refactored to use service layer interfaces
- LoadAreasCmd uses AreaServiceInterface.List()
- LoadSubareasCmd uses SubareaServiceInterface.ListByArea()
- LoadProjectsCmd uses ProjectServiceInterface.ListBySubareaRecursive() for hierarchical loading
- LoadTasksCmd uses TaskServiceInterface.ListByProject()
- Model structure updated to use service interfaces instead of db.Querier
- Comprehensive test coverage with mocked services (table-driven tests)
- Removed converter layer from TUI (services return domain types directly)
- Benefits: Better separation of concerns, easier mocking, centralized business logic

### ✅ Completed: Refactor CRUD Commands to Use Services (Task-39)
- All 8 CRUD commands refactored to use service layer interfaces
- CreateSubareaCmd uses SubareaServiceInterface.Create()
- CreateProjectCmd uses ProjectServiceInterface.Create()
- CreateTaskCmd uses TaskServiceInterface.Create()
- CreateAreaCmd uses AreaServiceInterface.Create()
- UpdateAreaCmd uses AreaServiceInterface.Update()
- DeleteAreaCmd uses AreaServiceInterface.SoftDelete/HardDelete()
- ReorderAreasCmd uses AreaServiceInterface.ReorderAll()
- LoadAreaStatsCmd uses AreaServiceInterface.GetStats()
- Removed all direct db.Querier usage from commands.go
- Comprehensive test coverage with mocked services (table-driven tests)
- All commands follow same pattern: service method calls, domain type returns, no converters needed
- Benefits: Consistent architecture, easier testing, centralized business logic, type safety

## Design Decisions

1. **bubbletea over other frameworks**: Simple MVU pattern, great community, pure Go
2. **lipgloss for styling**: Declarative styles, flexible layouts, excellent docs
3. **3-column layout**: Matches ADHD-friendly visual hierarchy
4. **Focus-aware borders**: Clear visual feedback for active column
5. **Wrapping navigation**: No dead-ends, continuous cycling
6. **Lazy initialization**: Wait for resize message before rendering
7. **No side effects in View**: Strict MVU pattern for testability

## Performance Considerations

- **Render on change only**: bubbletea only calls View() when model changes
- **Minimal allocations**: Reuse styles, avoid string concat in loops
- **Lazy init**: Wait for terminal size before first render
- **Efficient updates**: Return new model, don't mutate in-place

## Accessibility

- **High contrast colors**: Active elements clearly highlighted
- **Multiple keybindings**: h/l, arrows, Tab for same action
- **Visual feedback**: Borders change thickness on focus
- **Keyboard-first**: All actions accessible via keyboard

## Related Documentation

- [Data Layer Architecture](doc-1 - Data-Layer-Architecture.md)
- [CLI CRUD Operations Guide](doc-2 - CLI-CRUD-Operations-Guide.md)
- [bubbletea documentation](https://github.com/charmbracelet/bubbletea)
- [lipgloss documentation](https://github.com/charmbracelet/lipgloss)
