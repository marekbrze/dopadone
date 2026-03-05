# Terminal User Interface (TUI) Documentation

This document provides comprehensive documentation for the ProjectDB Terminal User Interface, including architecture, components, and implementation details.

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
├── help/               # Help modal component
│   ├── help.go         # Help modal logic
│   └── styles.go       # Lipgloss styling
├── toast/              # Toast notification component
│   ├── toast.go        # Toast logic with auto-dismiss
│   └── styles.go       # Error/success styling
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

### 3. Project Tree Navigation

Hierarchical display of projects and sub-projects with expand/collapse functionality.

**Implementation**: `tree/` package

- **node.go**: Tree node data structure
- **builder.go**: Builds tree from flat project list
- **navigation.go**: Handles up/down navigation with wrapping
- **renderer.go**: Renders tree with indentation and expand/collapse indicators

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

### 7. Footer with Quick Reference

Persistent footer showing common keyboard shortcuts.

**Display**: `h/l: columns | j/k: nav | a: add | ?: help | q: quit`

**Implementation**: Footer rendered in main `View()` function

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

### General

| Key | Action | Description |
|-----|--------|-------------|
| `?` | Help | Show help modal with all shortcuts |
| `q`, `Ctrl+C` | Quit | Exit the TUI |
| `Escape` | Cancel/Close | Close modal or cancel operation |

## State Management

### Model Structure

```go
type Model struct {
    // Core data
    areas    []Area
    subareas []Subarea
    projects []Project
    tasks    []Task
    
    // UI state
    focusedColumn  FocusColumn
    selectedArea   int
    selectedSubarea int
    selectedProject int
    selectedTask   int
    
    // Tree state
    tree          *tree.Node
    expandedNodes map[string]bool
    
    // Modal state
    showHelp  bool
    showAddModal bool
    addModalContext AddModalContext
    
    // Toasts
    toasts []toast.Toast
    
    // Terminal size
    width  int
    height int
}
```

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
   - Mock data for predictable testing

2. **Integration Tests** (`integration_test.go`):
   - End-to-end user flows
   - Keyboard input sequences
   - State transitions
   - Modal interactions

3. **Navigation Tests** (`navigation_test.go`):
   - Column focus switching
   - Tab navigation
   - Tree navigation with wrapping

### Running Tests

```bash
# Run all TUI tests
go test ./internal/tui/... -v

# Run specific test
go test ./internal/tui/... -v -run TestHelpModal

# Run with coverage
go test ./internal/tui/... -cover
```

### Manual Verification Checklist

For each release, manually verify:

- [ ] All keyboard shortcuts work as documented
- [ ] Tab navigation wraps correctly
- [ ] Column focus indicators are visible
- [ ] Tree expand/collapse works
- [ ] Help modal shows all shortcuts
- [ ] Toast notifications appear and auto-dismiss
- [ ] Quick-add modal creates items in correct context
- [ ] Terminal resize handling
- [ ] No visual artifacts or rendering issues

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

### Responsive Design

The TUI adapts to terminal size:
- Minimum supported: 80x24
- Columns resize proportionally
- Text truncation with ellipsis for long names
- Proper wrapping for narrow terminals

## Error Handling

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
projectdb tui
```

## Related Documentation

- [README.md](../README.md) - User-facing TUI documentation
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Lipgloss Documentation](https://github.com/charmbracelet/lipgloss)
