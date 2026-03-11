# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Delete Functionality in TUI (Task-68)

Comprehensive delete functionality for the Terminal User Interface with confirmation flow:

- **Delete Key Binding**: Press `d` on selected item to open confirmation modal
- **Confirmation Modal**: Reusable `confirmmodal` component with y/n/escape handling
  - Displays entity type and item name
  - Theme-aware warning styling (red border for destructive actions)
  - Long names truncated with ellipsis (max 40 chars)
- **Cascade Soft Delete for Projects**: `SoftDeleteWithCascade` service method
  - Recursively soft deletes all child projects
  - Soft deletes all tasks within deleted projects
  - Transaction-wrapped for atomicity
- **Entity Support**: Subareas, Projects (with cascade), Tasks
- **Toast Notifications**: Success/error feedback after deletion
- **Column Refresh**: Automatic data reload after successful delete
- **No-op on Empty Columns**: Safe behavior when no item selected
- **Footer Update**: Shows `d: delete` shortcut

**Implementation Files**:
- `internal/tui/confirmmodal/` - Confirmation modal component
- `internal/tui/delete_handlers.go` - Delete key and message handlers
- `internal/service/project_service.go` - Cascade soft delete method
- `internal/tui/renderer_footer.go` - Updated footer with delete shortcut

### Changed

#### Database Storage Location (Task-80)

Default database path changed from current working directory to user config directory:

- **New Default Paths**:
  - Linux: `~/.config/dopadone/dopadone.db`
  - macOS: `~/Library/Application Support/dopadone/dopadone.db`
  - Windows: `%APPDATA%/dopadone/dopadone.db`
- **Automatic Directory Creation**: Directory created automatically if it doesn't exist
- **Fallback Behavior**: Falls back to `./dopadone.db` if user config directory is unavailable
- **Backward Compatibility**: Users can still override with `--db` flag, `DOPA_DB_PATH` env, or config file
- **Migration Support**: `MigrateFromOldPath()` utility function available for future migration

**Implementation Files**:
- `internal/cli/dbpath.go` - Core path logic (DefaultDBPath, EnsureDirExists, MigrateFromOldPath)
- `internal/cli/db.go` - Updated Connect() to create directory automatically
- `cmd/dopa/config.go` - Updated resolveDBPath() to use DefaultDBPath()
- `cmd/dopa/main.go` - Changed --db flag default to empty string

### Deprecated

### Removed

### Fixed

### Security

## [1.0.0] - 2026-03-09

**Initial Release** - Dopadone is a lightweight, SQLite-based CLI project management tool designed for developers who prefer staying in the terminal.

### Added

#### Core Application

- **Hierarchical Data Model**: Areas → Subareas → Projects → Tasks with full CRUD operations
- **SQLite Storage**: Local-first database with embedded migrations
- **CLI Interface**: Complete command-line interface using Cobra framework
  - CRUD commands for areas, subareas, projects, and tasks
  - Filtering and listing with `--filter` flag
  - Multiple output formats: table, JSON, YAML
- **Terminal User Interface (TUI)**: Interactive Bubble Tea-based TUI
  - 3-column browser layout (Subareas | Projects | Tasks)
  - Area tabs for quick navigation
  - Quick-add modal for creating items
  - Help modal with keyboard shortcuts
  - Focus-aware borders and visual feedback

#### Service Layer Architecture

Comprehensive service layer with dependency injection for testability and maintainability:

- **Service Interfaces**: `AreaServiceInterface`, `SubareaServiceInterface`, `ProjectServiceInterface`, `TaskServiceInterface`
- **Business Logic Layer**: All validation and business rules in services
- **Dependency Injection**: Services injected into CLI and TUI layers
- **Mock Support**: Service mocks for unit testing

#### TUI Features

- **Space-Activated Command Menu (Task-50)**: LazyVim-style which-key command palette
  - Press Space to reveal available commands
  - Context-aware command suggestions
  - Visual keyboard shortcut hints
- **Task Completion Toggle (Task-49)**: Press `x` to mark tasks as done/todo
- **Adaptive Theme System**: Color scheme adapts to terminal capabilities
- **Responsive Layout**:
  - Proportional column widths (Task-42): Columns scale based on terminal width
  - Stacked layout for narrow terminals (Task-43): Vertical layout when width < 120 chars
- **Project Tree Navigation**: 
  - Arrow-based indicators (`▸`/`▾`) for expand/collapse
  - 2-space indentation for hierarchy
  - Modern minimalist design
- **Quick-Add Modal**: Context-aware creation with `a` key

#### Nested Task Grouping (Task-51)

Comprehensive nested task grouping functionality for hierarchical task display:

- **Recursive Task Loading**: `ListByProjectRecursive` using WITH RECURSIVE CTE
- **GroupedTasks Domain Model**: `TaskGroup` and `GroupedTasks` structs with factory methods
- **Expand/Collapse UI**: Visual indicators and keyboard navigation
- **Performance Optimized**: Batch loading of project names (no N+1 queries)
- **State Persistence**: Expanded/collapsed state remembered across navigation

#### Error Handling System (Task-55)

Centralized error handling across all application layers:

- **Domain Error Types**: `ErrNotFound`, `ErrInvalidInput`, `ErrDatabaseError`
- **Custom Error Types**: `ValidationError`, `DatabaseError`, `NotFoundError`
- **Error Wrapping**: `Unwrap()` for `errors.Is()` and `errors.As()` compatibility
- **User-Friendly Messages**: Technical errors mapped to clear user messages
- **Graceful Degradation**: Application continues when some data fails to load

#### CLI Features

- **JSON Output Support (Task-23)**: All create commands support `--format json`
- **Filtering**: Advanced filter syntax with `--filter` flag
- **Task Priority**: Mark tasks as "next" for priority focus
- **Soft Delete**: Recoverable deletes by default
- **Hard Delete**: Permanent deletion with `--permanent` flag

#### Database Features

- **Transaction Support (Task-31)**: ACID compliance for multi-entity operations
  - Hard delete cascade operations wrapped in transactions
  - Serializable isolation level
  - Automatic rollback on error
- **Server-Side Filtering (Task-32)**: Queries execute filters at database level
- **Nullable Types (Task-30)**: Proper nullable timestamp handling
- **Database Migrations**: Embedded migrations via goose

#### GitHub Actions Release Workflow (Task-65)

Automated CI/CD pipeline for building and publishing releases:

- **Multi-Platform Builds**: Linux (amd64), macOS (amd64, arm64), Windows (amd64)
- **Version Injection**: Compile-time version info via ldflags
- **Distribution Archives**: `.tar.gz` (Unix) and `.zip` (Windows)
- **SHA256 Checksums**: Verification checksums for all binaries
- **Pre-release Support**: Tags with hyphens marked as pre-release

#### Development Tools

- **Makefile**: Build, test, lint, and distribution targets
- **Cross-Compilation**: Build for all platforms from any OS
- **Test Data Seeder**: Contextual seed data for development
- **Development Script**: `scripts/dev.sh` for common tasks

#### Documentation

- **Architecture Documentation**: Comprehensive 7-part series covering all layers
- **TUI Documentation**: Complete TUI architecture and components guide
- **Transaction Documentation**: Database transaction patterns and best practices
- **Release Process**: Versioning, tagging, and deployment workflow
- **CI/CD Documentation**: GitHub Actions workflow details

### Changed

#### Project Rebranding (Task-47)

Comprehensive rebranding from **ProjectDB** to **Dopadone**:

- **Product Name**: Dopadone (CLI: `dopa`)
- **Module Path**: `github.com/marekbrze/dopadone`
- **Binary**: `dopa` (was: `projectdb`)
- **Database**: `dopadone.db` (was: `projectdb.db`)
- **All imports, scripts, and documentation updated**

#### Repository Reference Migration (Task-62)

Updated all repository references from placeholder URLs to production:

- **Go Module**: `github.com/example/dopadone` → `github.com/marekbrze/dopadone`
- **All import statements** across 52+ Go files
- **Build system**: Makefile, installation scripts
- **Documentation**: All markdown files updated
- **Version info**: Correct binary names and GitHub URLs

#### Tree Visual Design (Task-45)

Modernized project tree rendering:

- Replaced box-drawing characters with simple 2-space indentation
- Arrow indicators (`▸`/`▾`) instead of `[-]`/`[+]`
- Removed vertical connector lines for cleaner appearance

#### Text Truncation (Task-44)

Improved ellipsis handling in narrow columns:

- Ellipsis replaces only overflow text, not entire value
- Preserves as much content as possible in constrained space

### Deprecated

Nothing deprecated in this release.

### Removed

Nothing removed in this release.

### Fixed

- **CLI JSON Output (Task-23)**: Create commands now properly support JSON output format
- **Test Failures (Task-48)**: Fixed `getParentContext` assignment mismatch in `app_test.go`
- **Area Tabs Styling (Task-46)**: Fixed color scheme compatibility with area tabs

### Security

- **Reproducible Builds**: Fixed Go version, `-trimpath` flag, pure Go build
- **Checksum Verification**: SHA256 checksums for all release binaries
- **No External Dependencies**: Pure Go SQLite driver (modernc.org/sqlite)

## Version History

[1.0.0]: https://github.com/marekbrze/dopadone/releases/tag/v1.0.0
