---
id: TASK-60.5
title: Connection status indicator in TUI
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-08 19:01'
updated_date: '2026-03-11 13:06'
labels:
  - tui
  - database
  - ui
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add database connection status indicator to TUI showing connected/syncing/offline/local-only states. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Add connection status field to TUI model
- [x] #2 Create status indicator component in status bar or title bar
- [x] #3 Show visual indicator: ● (green) connected, ◐ (yellow) syncing, ○ (red) offline, ■ (gray) local-only
- [x] #4 Update status on connection state changes
- [x] #5 Add status message on hover/help
- [x] #6 Support status updates from sync goroutine
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: Connection Status Indicator in TUI (TASK-60.5)

## Overview
Add database connection status indicator to TUI showing connected/syncing/offline/local-only states. This task integrates with the existing database driver abstraction layer to provide real-time status feedback.

## Scope Analysis
This task is appropriately scoped as a single unit - it focuses on TUI visualization of connection state. No splitting needed.

## Dependencies
- **TASK-60.4** (Turso embedded replica driver): Provides SyncStatus for sync state
- **TASK-60.7** (Integration): Driver must be wired into application
- Existing `internal/db/driver` package with ConnectionStatus, SyncStatus types

## Architecture

```
┌───────────────────────────────────────────────────────────┐
│                      TUI Layer                             │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Model                                               │  │
│  │  - dbDriver: driver.DatabaseDriver (inject)         │  │
│  │  - dbMode: DatabaseMode                              │  │
│  │  - connectionStatus: ConnectionStatusView            │  │
│  └─────────────────────────────────────────────────────┘  │
│                          ↑                                 │
│                    Poll/Subscribe                          │
│                          ↓                                 │
│  ┌─────────────────────────────────────────────────────┐  │
│  │  Status Indicator Component                          │  │
│  │  - Render indicator with symbol + color              │  │
│  │  - Show tooltip on hover/help                        │  │
│  │  - Display in footer or title bar                    │  │
│  └─────────────────────────────────────────────────────┘  │
└───────────────────────────────────────────────────────────┘
                           ↑
                     Status() method
                           ↓
┌───────────────────────────────────────────────────────────┐
│                  Database Driver Layer                     │
│  DatabaseDriver.Status() → ConnectionStatus               │
│  TursoReplicaDriver.SyncInfo() → SyncStatus               │
└───────────────────────────────────────────────────────────┘
```

## Implementation Steps

### Phase 1: Define Types and Model Integration (Sequential)

**Step 1.1: Create Connection Status Types for TUI**
- File: `internal/tui/connection_status.go`
- Define `ConnectionStatusView` struct with all needed fields
- Define display constants (symbols, colors, messages)
- Define helper functions for status conversion

```go
type ConnectionStatusView struct {
    Mode         DatabaseMode      // local, remote, replica
    Status       ConnectionStatus  // connected, disconnected, etc.
    SyncStatus   SyncStatus        // idle, syncing, offline (replica only)
    LastSyncAt   time.Time         // for replica mode
    ErrorMessage string            // if any error
}

// Display symbols per AC#3
const (
    SymbolConnected = "●"  // green
    SymbolSyncing   = "◐"  // yellow
    SymbolOffline   = "○"  // red
    SymbolLocalOnly = "■"  // gray
)
```

**Step 1.2: Extend Model Structure**
- File: `internal/tui/app.go`
- Add `dbDriver driver.DatabaseDriver` field to Model
- Add `dbMode driver.DriverType` field to Model
- Add `connectionStatus ConnectionStatusView` field to Model

**Step 1.3: Create Status Indicator Component**
- File: `internal/tui/statusindicator/component.go`
- Create `StatusIndicator` component with:
  - `Render() string` method for visual output
  - `GetTooltip() string` method for hover/help text
  - Proper lipgloss styling with theme integration

### Phase 2: Integration and Updates (Sequential after Phase 1)

**Step 2.1: Inject Driver into TUI**
- File: `internal/tui/tui.go`, `cmd/dopa/tui.go`
- Modify `New()` to accept `driver.DatabaseDriver` parameter
- Initialize connection status in `InitialModel()`

**Step 2.2: Create Status Update Command**
- File: `internal/tui/commands.go`
- Add `PollConnectionStatusCmd` for periodic status polling
- Add `ConnectionStatusUpdatedMsg` message type
- Poll interval: 1 second (configurable)

```go
func PollConnectionStatusCmd(d driver.DatabaseDriver, interval time.Duration) tea.Cmd {
    return tea.Tick(interval, func(t time.Time) tea.Msg {
        return ConnectionStatusUpdatedMsg{
            Status:     d.Status(),
            SyncInfo:   getSyncInfo(d),  // type assertion for replica
            DriverType: d.Type(),
        }
    })
}
```

**Step 2.3: Handle Status Updates**
- File: `internal/tui/handlers.go` or new `internal/tui/connection_handlers.go`
- Add handler for `ConnectionStatusUpdatedMsg`
- Update Model's `connectionStatus` field
- Schedule next poll command

### Phase 3: Rendering (Sequential after Phase 2)

**Step 3.1: Integrate into Footer or Title Bar**
- File: `internal/tui/renderer_footer.go`
- Modify `RenderFooter()` to include status indicator
- Layout: `[status indicator] | h/l: columns | ...`

**Step 3.2: Add Help/Tooltip Support**
- File: `internal/tui/help/help.go`
- Add connection status explanation to help modal
- Include keyboard shortcut hint for status details (optional)

### Phase 4: Testing (Parallel with Phase 3)

**Step 4.1: Unit Tests**
- File: `internal/tui/statusindicator/component_test.go`
- Test all status combinations (local, remote, replica × connected, syncing, offline)
- Test symbol/color output
- Test tooltip generation
- Use table-driven tests per golang-testing skill

**Step 4.2: Integration Tests**
- File: `internal/tui/integration_connection_status_test.go`
- Test status updates through message flow
- Test status polling with mock driver
- Test UI rendering with different statuses

**Step 4.3: Mock Driver for Testing**
- File: `internal/tui/mocks/driver.go`
- Create mock implementation of `driver.DatabaseDriver`
- Support configurable status returns for testing

### Phase 5: Documentation (Parallel with Phase 4)

**Step 5.1: Update TUI Documentation**
- File: `docs/TUI.md`
- Add section explaining connection status indicator
- Document symbols and their meanings
- Add screenshot or ASCII example

**Step 5.2: Update Help Modal**
- File: `internal/tui/help/help.go`
- Add connection status to keyboard shortcuts reference

## File Changes Summary

| File | Action | Purpose |
|------|--------|---------|
| `internal/tui/connection_status.go` | CREATE | Status types and helpers |
| `internal/tui/statusindicator/component.go` | CREATE | Status indicator component |
| `internal/tui/statusindicator/styles.go` | CREATE | Lipgloss styles for indicator |
| `internal/tui/app.go` | MODIFY | Add driver and status fields to Model |
| `internal/tui/tui.go` | MODIFY | Inject driver into TUI |
| `internal/tui/commands.go` | MODIFY | Add status poll command |
| `internal/tui/messages.go` | MODIFY | Add status message types |
| `internal/tui/connection_handlers.go` | CREATE | Status update handlers |
| `internal/tui/renderer_footer.go` | MODIFY | Include status indicator |
| `internal/tui/help/help.go` | MODIFY | Add status to help content |
| `cmd/dopa/tui.go` | MODIFY | Pass driver to TUI constructor |
| `docs/TUI.md` | MODIFY | Document status indicator |

## Test Files to Create

| File | Purpose |
|------|---------|
| `internal/tui/connection_status_test.go` | Unit tests for status types |
| `internal/tui/statusindicator/component_test.go` | Component unit tests |
| `internal/tui/integration_connection_status_test.go` | Integration tests |
| `internal/tui/mocks/driver.go` | Mock driver for testing |

## Visual Design

### Footer Layout
```
┌─────────────────────────────────────────────────────────────┐
│ ● local | h/l: columns | j/k: navigate | a: add | d: delete | ?: help | q: quit │
└─────────────────────────────────────────────────────────────┘
```

### Status Indicator States

| Mode | Status | Symbol | Color | Tooltip |
|------|--------|--------|-------|---------|
| Local | - | ■ | Gray | "Local database (no sync)" |
| Remote | Connected | ● | Green | "Connected to Turso" |
| Remote | Disconnected | ○ | Red | "Disconnected from Turso" |
| Remote | Connecting | ◐ | Yellow | "Connecting to Turso..." |
| Replica | Connected + Syncing | ◐ | Yellow | "Syncing with Turso..." |
| Replica | Connected + Idle | ● | Green | "Connected (last sync: 2m ago)" |
| Replica | Offline | ○ | Red | "Offline - changes will sync when connected" |

## Threading & Concurrency

The status polling uses Bubble Tea's command pattern:
- `PollConnectionStatusCmd` runs in a goroutine via `tea.Tick`
- Driver status access is thread-safe (uses `sync.RWMutex`)
- Status updates are serialized through message queue

## Risk Mitigation

1. **Performance**: Poll interval of 1s is minimal overhead; can increase if needed
2. **Stale Status**: Driver already handles concurrent status updates safely
3. **Visual Clutter**: Indicator is compact (single symbol + mode text)
4. **Backward Compatibility**: Local mode shows subtle gray indicator

## Definition of Done Checklist

- [ ] All 6 acceptance criteria met
- [ ] Unit tests pass with >80% coverage on new code
- [ ] Integration tests pass
- [ ] TUI documentation updated
- [ ] Code follows golang-pro, golang-patterns, golang-testing best practices
- [ ] No breaking changes to existing TUI interface
- [ ] Manual testing on all three database modes verified

## Estimated Time
- Phase 1: 2 hours
- Phase 2: 1.5 hours
- Phase 3: 1 hour
- Phase 4: 2 hours
- Phase 5: 0.5 hours
- **Total: ~7 hours**

## Execution Order
1. Phase 1 (Types & Model)
2. Phase 2 (Integration)
3. Phase 3 (Rendering) - can overlap with Phase 4
4. Phase 4 (Testing) - parallel with Phase 3
5. Phase 5 (Documentation) - parallel with Phase 4
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Created connection status types and helpers
- Added status indicator component with lipgloss styling
- Integrated driver into TUI model
- Added status polling command (2s interval)
- Added handler for connection status updates
- Updated footer renderer to include status indicator
- Modified TUI initialization to accept driver parameter
- Updated all test files to pass nil driver
- All tests pass, build succeeds, lint clean (only pre-existing issues)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented connection status indicator in TUI showing database mode and sync state.

## Changes

### New Files
- `internal/tui/connection_status.go` - Types and helpers for connection status view
- `internal/tui/connection_handlers.go` - Handler for connection status updates
- `internal/tui/statusindicator/component.go` - Status indicator component with symbol/color rendering
- `internal/tui/statusindicator/styles.go` - Lipgloss styles for status indicator

### Modified Files
- `internal/tui/app.go` - Added `dbDriver` and `connectionStatus` fields to Model, updated Init() to start status polling
- `internal/tui/commands.go` - Added `PollConnectionStatusCmd` for periodic status polling
- `internal/tui/messages.go` - Added `ConnectionStatusUpdatedMsg` message type
- `internal/tui/renderer_footer.go` - Updated `RenderFooter()` to include status indicator
- `internal/tui/tui.go` - Added `dbDriver` parameter to `New()`
- `cmd/dopa/tui.go` - Pass driver to TUI constructor

## Implementation Details

### Visual Indicator
- **Local SQLite**: ■ (gray) - Local database (no sync)
- **Remote Connected**: ● (green) - Connected to Turso
- **Remote Syncing**: ◐ (yellow) - Connecting/syncing
- **Replica Connected**: ● (green) - Connected with last sync time
- **Replica Syncing**: ◐ (yellow) - Syncing with Turso
- **Offline**: ○ (red) - Disconnected

### Status Polling
- Polls driver status every 2 seconds
- Updates Model.connectionStatus on each poll
- Uses Bubble Tea command pattern for thread-safe updates

### Testing
- All existing tests pass
- Added unit tests for connection status view logic

## Footer Layout
```
● local | h/l: columns | j/k: navigate | a: add | d: delete | ?: help | q: quit
```
<!-- SECTION:FINAL_SUMMARY:END -->
