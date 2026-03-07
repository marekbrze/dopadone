# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Nested Task Grouping Feature (Task-51)

**Hierarchical task display with expandable/collapsible subproject groups**

Implemented comprehensive nested task grouping functionality that displays tasks from parent projects and their nested subprojects in a hierarchical, visually organized format.

**Overview**:
When a parent project is selected in the Projects column, the Tasks column now displays:
1. Direct tasks belonging to the selected project (at the top, ungrouped)
2. Tasks from nested subprojects, grouped under collapsible headers
3. Visual indicators showing expand/collapse state (▸ / ▾)
4. Task counts per subproject group

**Architecture**:

**Phase 1 - Service Layer (Task-52)**:
- Added `ListByProjectRecursive` method to `TaskService` using WITH RECURSIVE CTE
- SQL query in `queries/tasks.sql` traverses project hierarchy
- Returns tasks from selected project and all nested subprojects
- Filters out deleted projects and tasks
- Orders by: is_next DESC, priority DESC, deadline ASC, title ASC
- Comprehensive test coverage: 12 test cases including edge cases

**Phase 2 - Domain Model (Task-54)**:
- Created `TaskGroup` struct with ProjectID, ProjectName, Tasks, IsExpanded
- Created `GroupedTasks` struct with DirectTasks, Groups, TotalCount
- Factory method `NewGroupedTasks()` with graceful edge case handling
- Mutation methods: `AddTask()`, `RemoveTask()`, `ToggleGroup()`, `Clear()`
- Order preservation for tasks and groups
- Test coverage: 96.4%+ (25+ test cases)

**Phase 3 - Service Integration (Task-57)**:
- Added `GetGroupedTasks` method to `TaskService`
- Batch loading of project names (eliminates N+1 queries)
- Backward compatible - maintains flat `Tasks` field
- State persistence for expand/collapse across navigation
- Test coverage: 85%+

**Phase 4 - TUI Rendering (Task-58)**:
- Renders tasks with group headers showing subproject name
- Visual indicators for expanded/collapsed groups
- Proper indentation (2 spaces) for nested tasks
- Task count display per group
- Subtle styling for headers (dimmed, no reverse highlight)
- Text truncation to prevent wrapping
- Performance: O(n) rendering

**Phase 5 - TUI Interaction (Task-56)**:
- Keyboard shortcuts: Enter/Space to toggle groups
- Navigation skips headers when collapsed
- State persistence in `expandedTaskGroups` map
- Selection adjustment when collapsing groups
- Helper methods for navigation and rendering
- State saved in `AreaState.ExpandedTaskGroups`

**Phase 6 - Error Handling (Task-55)**:
- Graceful handling of empty projects
- Handling of orphaned subprojects
- User-friendly error messages
- Error state tracking and rendering

**Phase 7 - Performance Optimization (Task-53)**:
- Batch loading of project names (O(1) queries)
- Single SQL query for recursive loading
- O(n) time complexity for grouping
- Benchmarks: 100, 1000, 10000 tasks
- Target: <100ms for 1000 tasks achieved

**Files Modified**:
- `queries/tasks.sql` - Added `ListTasksByProjectRecursive` query
- `queries/projects.sql` - Added `ListProjectsByIDs` query for batch loading
- `internal/domain/task_group.go` - New file with GroupedTasks domain model
- `internal/domain/errors.go` - Centralized error types
- `internal/service/task_service.go` - Added `ListByProjectRecursive` and `GetGroupedTasks`
- `internal/service/project_service.go` - Added `ListByIDs` method
- `internal/tui/commands.go` - Updated `LoadTasksCmd` to use `GetGroupedTasks`
- `internal/tui/renderer.go` - Grouped task rendering with headers
- `internal/tui/navigator.go` - Navigation logic for grouped tasks
- `internal/tui/handlers.go` - Expand/collapse interaction handlers
- `internal/tui/state.go` - State persistence for expanded groups
- `internal/tui/constants.go` - Error message constants

**Test Coverage**:
- Domain Layer: 96.4% (task_group.go)
- Service Layer: 85%+ (task_service_test.go)
- TUI Layer: 80%+ (navigation, rendering tests)

**Documentation**:
- Created `docs/FEATURE_NESTED_TASK_GROUPING.md` with comprehensive implementation guide
- Updated `docs/architecture/02-domain-layer.md` with GroupedTasks patterns
- Updated `docs/architecture/03-service-layer.md` with recursive loading patterns
- Updated `docs/TUI.md` with grouped task rendering and interaction

**Benefits**:
- **Improved Organization**: Clear visual hierarchy of tasks across nested projects
- **Better UX**: Expand/collapse groups to focus on relevant tasks
- **Performance**: Optimized O(n) algorithms with no N+1 queries
- **State Persistence**: Expanded/collapsed state remembered across navigation
- **Scalability**: Handles 1000+ tasks efficiently (<100ms)

**User Impact**:
Users can now see all tasks from a project and its subprojects in one view, organized hierarchically with collapsible groups. This makes managing complex project structures much more intuitive.

#### Centralized Error Handling System (Task-55)

**Comprehensive error handling across all application layers**

Implemented a centralized error handling system with domain error types, service layer error wrapping, and TUI error state management.

**Changes**:

**Domain Layer** (`internal/domain/errors.go`):
- Added centralized sentinel errors: `ErrNotFound`, `ErrInvalidInput`, `ErrDatabaseError`, `ErrEmptyID`
- Created custom error types: `ValidationError`, `DatabaseError`, `NotFoundError`
- Implemented error factory functions: `NewValidationError()`, `NewDatabaseError()`, `NewNotFoundError()`
- Added helper functions for type-safe error checking: `IsNotFound()`, `IsDatabaseError()`, `IsValidationError()`
- Implemented error wrapping/unwrapping with `Unwrap()` methods for compatibility with `errors.Is()` and `errors.As()`

**Service Layer**:
- Updated all services to use domain error types instead of generic errors
- Replaced generic "not found" errors with `domain.NewNotFoundError()`
- Implemented graceful handling of empty results (e.g., missing parent projects return empty, not error)
- Added context-aware error wrapping with operation details
- Mapped `sql.ErrNoRows` to domain-specific not found errors

**TUI Layer**:
- Added error state tracking for each column: `areaLoadError`, `subareaLoadError`, `projectLoadError`, `taskLoadError`
- Implemented user-friendly error message formatting with `formatUserError()`
- Added error rendering with visual indicators (red error messages with ✗ icon)
- Defined user-friendly error message constants: `ErrMsgDatabase`, `ErrMsgTimeout`, `ErrMsgCancelled`, `ErrMsgNotFound`
- Enhanced error recovery with retry mechanisms and graceful degradation
- Added comprehensive error handling tests (`internal/domain/errors_test.go`, `internal/tui/task_navigation_test.go`)

**Documentation**:
- Updated `docs/architecture/02-domain-layer.md` with "Error Handling Patterns" section
- Updated `docs/architecture/03-service-layer.md` with "Error Wrapping Best Practices" section
- Updated `docs/TUI.md` with "Error State Management" section
- Documented error checking patterns, error wrapping principles, and best practices

**Benefits**:
- **Consistency**: All layers use the same error types and patterns
- **Type Safety**: Custom error types provide structured error information
- **Error Chaining**: `Unwrap()` enables `errors.Is()` and `errors.As()` compatibility
- **User-Friendly**: Technical errors mapped to clear user messages
- **Testability**: Easy to check for specific error types in tests
- **Debugging**: Error context preserved for troubleshooting while showing clean messages to users
- **Graceful Degradation**: Application continues functioning when some data fails to load

**Files Modified**:
- `internal/domain/errors.go` (NEW)
- `internal/domain/errors_test.go` (NEW)
- `internal/service/area_service.go`
- `internal/service/project_service.go`
- `internal/service/subarea_service.go`
- `internal/service/task_service.go`
- `internal/tui/app.go`
- `internal/tui/constants.go`
- `internal/tui/handlers.go`
- `internal/tui/model.go`
- `internal/tui/navigator.go`
- `internal/tui/renderer.go`
- `internal/tui/state.go`
- `internal/tui/task_navigation_test.go` (NEW)
- `docs/architecture/02-domain-layer.md`
- `docs/architecture/03-service-layer.md`
- `docs/TUI.md`

**Testing**:
- All error handling tests passing
- Domain error type tests verify error messages and wrapping behavior
- TUI error navigation tests verify grouped task handling with expand/collapse functionality
- Service layer tests updated to use new error types

**Backward Compatibility**:
- All existing functionality preserved
- Error handling is now more robust and informative
- No breaking changes to public APIs

### Changed

#### Tree Visual Design (Task-45)

**Modernized project tree rendering with arrow-based indicators**

The project tree component now uses a clean, minimalist design with arrow indicators instead of traditional box-drawing characters.

**Changes**:
- Replaced box-drawing characters (├─└│) with simple 2-space indentation
- Replaced expand/collapse indicators `[-]`/`[+]` with arrows `▾`/`▸`
- Removed vertical connector lines for cleaner visual appearance
- Improved readability with consistent indentation at all depth levels

**Files Modified**:
- `internal/tui/tree/constants.go`: Updated tree character constants
- `internal/tui/tree/renderer.go`: Simplified rendering logic
- `internal/tui/tree/renderer_test.go`: Updated test expectations
- `docs/TUI.md`: Added documentation for new tree styling

**Visual Comparison**:

Before (box-drawing):
```
├─ Project A
│  ├─ Subproject A1
│  └─ Subproject A2
└─ Project B
```

After (arrow indicators):
```
▾ Project A
  Subproject A1
  ▸ Subproject A2
Project B
```

**Benefits**:
- Reduced visual clutter with no vertical connector lines
- Clearer expand/collapse state with intuitive arrow indicators
- Modern, minimalist appearance
- Better readability on high-DPI displays
- Customizable through `TreeStyle` struct

**Backward Compatibility**:
- All existing tree navigation and functionality preserved
- Only visual rendering changed, no API changes
- Custom tree styles can still be applied via `TreeStyle` struct

**Testing**:
- All 45 tree renderer tests updated and passing
- Visual verification in TUI confirms modern appearance
- Expand/collapse functionality verified working
- Navigation preserved across all tree operations

**Documentation**:
- Added tree styling section to `docs/TUI.md`
- Updated tree rendering examples
- Documented customization options via `TreeStyle`
