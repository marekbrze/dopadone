# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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
