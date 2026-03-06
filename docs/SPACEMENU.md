# Space Menu Feature Documentation

## Overview

The Space Menu is a command palette feature inspired by LazyVim's which-key functionality. It provides a floating overlay command menu that appears when the Space key is pressed, offering quick access to various commands and configuration options.

## Features

### Keyboard-Driven Command Access

- **Space key trigger**: Press Space when no modal or input field is focused to open the command menu
- **Floating overlay**: Centered modal overlay displaying available commands
- **Keyboard shortcuts**: Each command has a single-key shortcut
- **Nested menus**: Support for hierarchical command organization

### Available Commands

#### Main Menu

| Key | Command | Description |
|-----|---------|-------------|
| `c` | Config | Open configuration submenu for Area management |
| `q` | Quit | Exit the application |

#### Config Submenu

| Key | Command | Description |
|-----|---------|-------------|
| `c` | Create Area | Create a new area |
| `e` | Edit Area | Edit selected area |
| `d` | Delete Area | Delete selected area |

### Dismissal Options

The menu can be dismissed in multiple ways:
- Press `Space` again
- Press `Escape`
- Press `q` (from main menu only)

## Architecture

### Component Structure

```
internal/tui/spacemenu/
├── spacemenu.go      # Main component logic (BubbleTea Model)
├── types.go          # Menu state and action type definitions
├── styles.go         # Lipgloss styling using theme system
└── spacemenu_test.go # Unit tests
```

### State Management

The menu uses a state machine pattern:

```go
type MenuState int

const (
    StateMain   MenuState = iota  // Main command menu
    StateConfig                    // Config submenu
)
```

### Message Types

```go
type CloseMsg struct{}                    // Close menu
type ActionMsg struct { Action MenuAction } // Execute action
```

### Integration Points

1. **Model Integration** (`internal/tui/app.go`):
   ```go
   type Model struct {
       spaceMenu       *spacemenu.SpaceMenu
       isSpaceMenuOpen bool
       // ... other fields
   }
   ```

2. **Key Handling Priority**:
   - Space key opens menu when no modal is active
   - Priority order: help > modal > area modal > space menu > normal

3. **View Rendering**:
   - Menu rendered as overlay on top of base UI
   - Uses z-index layering with other modals

### Theme Integration

The menu uses the existing theme system for consistent styling:

```go
theme := theme.Default

// Menu styling
menuStyle := lipgloss.NewStyle().
    Background(theme.Background()).
    Border(theme.Primary())
```

Colors automatically adapt to terminal background (light/dark mode).

## Implementation Details

### BubbleTea Pattern

The component follows the standard BubbleTea Model interface:

```go
func (m *SpaceMenu) Init() tea.Cmd
func (m *SpaceMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *SpaceMenu) View() string
```

### Key Handling

Keys are processed in the `Update` method:

```go
func (m *SpaceMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "escape", " ":
            return m, func() tea.Msg { return CloseMsg{} }
        case "c":
            if m.state == StateMain {
                m.state = StateConfig
            }
        case "q":
            if m.state == StateMain {
                return m, func() tea.Msg { return ActionMsg{Action: ActionQuit} }
            }
        }
    }
    return m, nil
}
```

### Action Handling

Actions are communicated back to the main app via messages:

```go
func (m *Model) handleSpaceMenuAction(msg spacemenu.ActionMsg) (tea.Model, tea.Cmd) {
    switch msg.Action {
    case spacemenu.ActionQuit:
        return m, tea.Quit
    case spacemenu.ActionConfig:
        // Handle config action
    }
    return m, nil
}
```

## Testing

### Unit Tests

Located in `internal/tui/spacemenu/spacemenu_test.go`:

- **Open/Close behavior**: Verify menu opens and closes correctly
- **Key handling**: Test all key shortcuts (Space/Esc/q/c)
- **State transitions**: Verify main → config submenu navigation
- **View rendering**: Ensure menu displays correctly
- **Theme integration**: Verify colors adapt to theme

### Integration Tests

Located in `internal/tui/integration_spacemenu_test.go`:

- **Space key flow**: Test Space key triggers menu from normal state
- **Modal priority**: Verify Space key doesn't trigger when other modals are open
- **Full navigation**: Test open → navigate → action → close flow
- **Quit action**: Verify quit action terminates app

### Running Tests

```bash
# Run spacemenu tests
go test ./internal/tui/spacemenu/... -v

# Run integration tests
go test ./internal/tui/... -v -run TestSpaceMenu

# Run with coverage
go test ./internal/tui/spacemenu/... -cover
```

## User Experience

### Visual Design

- **Centered overlay**: Menu appears in the center of the screen
- **Key highlighting**: Shortcut keys are visually emphasized
- **Descriptions**: Each command shows a brief description
- **Consistent styling**: Uses existing theme system for visual consistency

### Accessibility

- **Keyboard-only navigation**: No mouse required
- **Multiple dismissal options**: Escape, Space, or q to close
- **Clear visual hierarchy**: Key shortcuts are prominent
- **Context-aware**: Only appears when appropriate (no modal conflicts)

### User Flow Example

1. User presses `Space` in normal view → Menu opens
2. User presses `c` → Config submenu opens
3. User presses `c` again → Create Area modal opens
4. User creates area → Returns to normal view
5. Alternatively: User presses `Escape` at any point → Menu closes without action

## Future Enhancements

Potential improvements for future versions:

1. **Additional Commands**:
   - Search/filter functionality (`/`)
   - Theme selection (`t`)
   - Database operations (`d`)
   - Export/import (`e`)

2. **Enhanced Navigation**:
   - Arrow key navigation
   - Vim-style `j`/`k` movement
   - Search within commands

3. **Command History**:
   - Remember frequently used commands
   - Recent command list

4. **Customization**:
   - User-defined command shortcuts
   - Custom command groups
   - Command aliases

5. **Visual Improvements**:
   - Animated transitions
   - Command categories with icons
   - Fuzzy search highlighting

## Related Documentation

- [TUI.md](TUI.md) - Main TUI documentation
- [Architecture Overview](architecture/01-overview.md) - System architecture
- [BubbleTea Documentation](https://github.com/charmbracelet/bubbletea) - Framework docs

## Troubleshooting

### Menu doesn't open

- **Cause**: Another modal or input field is focused
- **Solution**: Close any open modals first

### Keys not responding

- **Cause**: Terminal may intercept certain keys
- **Solution**: Try different terminal emulator or check keyboard mapping

### Menu appears off-center

- **Cause**: Terminal size changed after menu opened
- **Solution**: Close and reopen menu to recalculate position

## Implementation Timeline

**Task**: TASK-50 - Space-Activated Command Menu (LazyVim-style which-key)

**Status**: Completed

**Implementation Date**: 2026-03-06

**Key Changes**:
- Created spacemenu component package
- Integrated with main app.go
- Added unit and integration tests
- Updated help modal and TUI documentation
- Followed BubbleTea golden rules for layout
- Used existing theme system for styling

**Acceptance Criteria**: All 8 criteria met ✓

**Definition of Done**: All 4 items completed ✓
