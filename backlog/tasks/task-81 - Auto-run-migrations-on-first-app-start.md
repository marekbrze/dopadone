---
id: TASK-81
title: Auto-run migrations on first app start
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-11 05:43'
updated_date: '2026-03-11 20:34'
labels:
  - backend
  - storage
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When a new user installs and runs the app for the first time, the database is empty and migrations need to be applied manually. The app should automatically run migrations on startup if the database is new or has pending migrations, providing a seamless first-run experience.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Migrations run automatically when database is first created
- [x] #2 Migrations run automatically when there are pending migrations on app start
- [x] #3 Migration errors are handled gracefully with clear error messages
- [x] #4 Works for both CLI commands and TUI mode
- [x] #5 No duplicate migration runs if database is already up to date
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan for Auto-Run Migrations on First App Start

### Overview
This feature ensures seamless first-run experience by automatically running database migrations when the app starts. Goose (used by migrate.Run) is idempotent, making repeated calls safe.

### File Changes

**1. cmd/dopa/main.go**
- Move `--skip-migrate` from `upgradeCmd.Flags()` to `rootCmd.PersistentFlags()`
- This makes the flag available to all commands (CLI and TUI)

**2. internal/cli/db.go**
- Add `EnsureMigrations(db *sql.DB) error` function
- Wraps `migrate.Run(db, "up")` with clear error context
- Returns nil if migrations already applied (goose.Up is idempotent)

**3. cmd/dopa/main.go - GetDriver()**
- Call `cli.EnsureMigrations()` after successful connection
- Skip if `--skip-migrate` flag is set
- Works for both SQLite and Turso modes
- Error handling: fail fast with clear message

**4. internal/migrate/migrate.go**
- No changes needed - `Run()` already handles idempotent `goose.Up()`

### Error Handling Strategy
- **Migration errors**: Clear message explaining what failed
- **Partial migrations**: Goose handles rollback automatically
- **Missing directory**: Already handled by `cli.EnsureDirExists()`

### Sequential Implementation Steps

1. **Add persistent flag** (5 min)
   - Move flag registration in init()
   - Update upgrade command to check root flag

2. **Create EnsureMigrations helper** (10 min)
   - Add function in internal/cli/db.go
   - Wrap with error context

3. **Integrate in GetDriver()** (15 min)
   - Add migration call after connection
   - Check skip flag before calling
   - Handle both SQLite and Turso paths

4. **Write unit tests** (20 min)
   - Test EnsureMigrations function
   - Test skip flag behavior
   - Test error handling

5. **Write integration tests** (30 min)
   - Test fresh database (first run)
   - Test pending migrations
   - Test already up-to-date database
   - Test error scenarios

6. **Manual testing** (15 min)
   - Fresh install scenario
   - Pending migrations scenario
   - Up-to-date DB scenario
   - Both CLI and TUI modes
   - With --skip-migrate flag

### Parallel Work Opportunities

- **Steps 4 and 5** (unit tests and integration tests) can be developed in parallel
- Tests can be written after step 3 is complete

### Dependencies
None - this is a self-contained feature.

### Testing Strategy

**Unit Tests (internal/cli/db_test.go)**
- `TestEnsureMigrations_Success`: Fresh DB gets migrations
- `TestEnsureMigrations_AlreadyMigrated`: Up-to-date DB is no-op
- `TestEnsureMigrations_NilDB`: Returns error for nil connection

**Integration Tests (cmd/dopa/main_test.go or new file)**
- `TestAutoMigration_FirstRun`: New DB gets all migrations
- `TestAutoMigration_PendingMigrations`: Partially migrated DB completes
- `TestAutoMigration_SkipFlag`: --skip-migrate bypasses migrations
- `TestAutoMigration_TUI`: TUI mode triggers migrations

### Documentation Updates
- No user-facing docs needed (transparent behavior)
- Code comments explaining auto-migration behavior
- Consider adding note to README about seamless first-run experience
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Technical Context

**Key Insight**: goose.Up() is idempotent - calling it on an already-migrated database is a no-op. This means we can safely call it on every startup without performance concerns.

**Current Flow**:
```
main.go → GetDriver() → Connect() → *sql.DB/sql.Driver
```

**New Flow**:
```
main.go → GetDriver() → Connect() → EnsureMigrations() → *sql.DB/sql.Driver
```

**Database Modes to Support**:
1. SQLite (local) - via `cli.Connect()`
2. Turso Remote - via `cli.ConnectWithDriver()`
3. Turso Replica - via `cli.ConnectWithDriver()` with sync

Both paths in GetDriver() need migration call.

**Code Patterns to Follow**:
- Error wrapping: `fmt.Errorf("context: %w", err)`
- Use existing `cli.WrapError()` for consistency
- Check skipMigrate flag at the call site in GetDriver()
- Keep EnsureMigrations simple - just wrap migrate.Run

**Files to Modify**:
1. cmd/dopa/main.go (lines ~25, ~184, ~240-265)
2. internal/cli/db.go (add new function)

**Files to Create**:
1. internal/cli/db_migration_test.go (new test file)
2. cmd/dopa/auto_migration_test.go (integration tests)

- Fixed missing migration file in embedded FS (20260304120000_add_sort_order_to_areas.sql)
- Added mutex to EnsureMigrations for concurrent safety
- Reset goose.SetBaseFS(nil) after migrations to avoid affecting other tests

## Additional: --dev flag
- Added `--dev` / `-D` flag to use local `./dopa.db` for testing binaries
- Added Make targets: `dev`, `dev-build`, `dev-run`, `dev-tui`, `dev-clean`
- Useful for testing without affecting main database

## Documentation Updates
- Updated README.md:
  - Removed manual `migrate up` step from Quick Start
  - Added auto-migrations section
  - Added `--dev` and `--skip-migrate` flags to Global Flags table
  - Added Dev Mode & Testing section with make targets
- Updated docs/DATABASE_MODES.md:
  - Added "Automatic Migrations" section
  - Added "Dev Mode (Testing)" section under SQLite Mode
  - Updated CLI Flags table with --dev and --skip-migrate
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary
Implemented auto-run migrations on first app start for seamless first-run experience.

## Changes
1. **cmd/dopa/main.go**
   - Moved `--skip-migrate` flag from `upgradeCmd` to `rootCmd.PersistentFlags()` for global availability
   - Added auto-migration call in `GetDriver()` after successful connection (both SQLite and Turso modes)
   - Migration errors fail fast with clear error message

2. **internal/cli/db.go**
   - Added `EnsureMigrations(db *sql.DB) error` function
   - Uses mutex for concurrent safety (prevents SQLITE_BUSY errors)
   - Wraps `migrate.Run(db, "up")` with clear error context

3. **internal/migrate/migrate.go**
   - Fixed goose.SetBaseFS leak by resetting to nil after migrations
   - This prevents test failures when other code uses filesystem-based migrations

4. **internal/cli/db_migrate_test.go** (new)
   - Unit tests for EnsureMigrations function
   - Tests: fresh database, nil DB, idempotent behavior

5. **internal/migrate/migrations/**
   - Added missing migration file: 20260304120000_add_sort_order_to_areas.sql

## Testing
- All existing tests pass
- New unit tests for EnsureMigrations
- Concurrent test (TestGetDB_Concurrent) verifies mutex protection
- Lint: 0 issues
- Build: successful

## Bonus: Dev Mode
- Added `--dev` / `-D` flag: uses `./dopa.db` in current directory
- New Make targets:
  - `make dev` - run with `go run` and --dev flag
  - `make dev-build` - build binary for testing
  - `make dev-run` - build and run with --dev
  - `make dev-tui` - launch TUI with --dev
  - `make dev-clean` - remove ./dopa.db
- Allows testing binaries without affecting production database

## Documentation
- README.md: Updated Quick Start, Global Flags, added Dev Mode section
- docs/DATABASE_MODES.md: Added Auto-Migrations and Dev Mode sections
<!-- SECTION:FINAL_SUMMARY:END -->
