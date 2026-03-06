# Dopadone Rebranding Documentation

## Overview

This document details the comprehensive rebranding of the project management CLI tool from **ProjectDB** to **Dopadone**, with the CLI command changing to `dopa`.

**Date**: March 2026  
**Task**: TASK-47  
**Status**: Completed

## Brand Identity

### New Names

- **Product Name**: Dopadone
- **CLI Command**: `dopa`
- **Module Path**: `github.com/example/dopadone`
- **Database**: `dopadone.db`
- **Binary**: `dopa`

### Rationale

- **"Dopa"**: Short, memorable, ADHD-friendly command name (easy to type, easy to remember)
- **"Dopadone"**: Full product name that combines "dopa" (referring to dopamine) with "done" (task completion)
- **Identity**: Emphasizes the dopamine hit of completing tasks and managing projects effectively

## Changes Summary

### 1. Code Structure

#### Module Path
- **Before**: `github.com/example/projectdb`
- **After**: `github.com/example/dopadone`

#### Command Directory
- **Before**: `cmd/projectdb/`
- **After**: `cmd/dopa/`

#### Binary Name
- **Before**: `projectdb`
- **After**: `dopa`

#### Database
- **Before**: `projectdb.db`
- **After**: `dopadone.db`

### 2. Import Updates

All Go files updated with new import paths:

```go
// Before
import "github.com/example/projectdb/internal/domain"
import "github.com/example/projectdb/internal/service"
import "github.com/example/projectdb/internal/tui"

// After
import "github.com/example/dopadone/internal/domain"
import "github.com/example/dopadone/internal/service"
import "github.com/example/dopadone/internal/tui"
```

**Files Updated**:
- `cmd/dopa/*.go` (6 files)
- `internal/cli/*.go`
- `internal/converter/*.go`
- `internal/service/*.go` (8 files)
- `internal/tui/**/*.go` (30+ files)
- All test files

### 3. Build System

#### Makefile Changes

```makefile
# Before
BINARY_NAME=projectdb
DB_PATH=projectdb.db

# After
BINARY_NAME=dopa
DB_PATH=dopadone.db
```

**Updated Targets**:
- `build`: Binary name changed
- `build-linux`, `build-darwin`, `build-windows`: Output names updated
- `dist`: Archive names changed to `dopa-<platform>-<arch>`
- LDFLAGS updated with new module path
- All references to `projectdb` replaced with `dopa`

### 4. Documentation

#### Updated Files

**Root Documentation**:
- `README.md`: All references, examples, installation instructions
- `docs/START_HERE.md`: Product name, project structure, code examples

**Architecture Documentation**:
- `docs/architecture/01-overview.md`: System overview
- `docs/architecture/02-domain-layer.md`: Domain examples
- `docs/architecture/03-service-layer.md`: Service examples
- `docs/architecture/04-converter-layer.md`: Converter examples
- `docs/architecture/05-repository-layer.md`: Repository examples
- `docs/architecture/06-cli-layer.md`: CLI examples
- `docs/architecture/07-testing-strategy.md`: Test examples

**Other Documentation**:
- `docs/TUI.md`: TUI documentation
- `docs/RELEASE.md`: Release process
- `docs/TRANSACTIONS.md`: Transaction handling

**Backlog Documentation**:
- All task files in `backlog/tasks/` and `backlog/completed/` updated

### 5. Scripts

#### Installation Script
`scripts/install.sh`:
- Binary download URLs updated
- Binary names changed from `projectdb-*` to `dopa-*`
- Installation paths updated

#### Test Data Script
`scripts/seed-test-data.sh`:
- Database path references updated
- Command invocations changed from `projectdb` to `dopa`

#### Changelog Script
`scripts/generate-changelog.sh`:
- Repository references updated
- Binary name references updated

### 6. Configuration Files

**Go Module**:
- `go.mod`: Module path changed
- `go.sum`: Dependency checksums updated

**Development Tools**:
- `dev.sh`: Development script updated
- `sqlc.yaml`: SQLC configuration updated

## Installation Changes

### Before

```bash
# Install
go install github.com/example/projectdb/cmd/projectdb@latest

# Run
projectdb area list
projectdb task create --title "Example"
```

### After

```bash
# Install
go install github.com/example/dopadone/cmd/dopa@latest

# Run
dopa area list
dopa task create --title "Example"
```

### Binary Downloads

**Old URLs**:
```
https://github.com/example/projectdb/releases/download/v1.0.0/projectdb-linux-amd64
https://github.com/example/projectdb/releases/download/v1.0.0/projectdb-darwin-amd64
```

**New URLs**:
```
https://github.com/example/dopadone/releases/download/v1.0.0/dopa-linux-amd64
https://github.com/example/dopadone/releases/download/v1.0.0/dopa-darwin-amd64
```

## Database Migration

The database filename changed but the schema remains the same:

- **Old**: `./projectdb.db` (or path specified by `--db` flag)
- **New**: `./dopadone.db` (or path specified by `--db` flag)

**Migration Path**:
```bash
# If you have existing data, simply rename the database file
mv projectdb.db dopadone.db

# Or specify the old database path explicitly
dopa --db ./projectdb.db area list
```

## API Compatibility

### Breaking Changes

1. **Binary Name**: Must use `dopa` instead of `projectdb`
2. **Module Path**: Go imports changed (affects extensions/plugins)
3. **Database Default**: Default database path changed

### Non-Breaking Changes

1. **Database Schema**: Identical, no migration needed
2. **CLI Flags**: All flags remain the same
3. **Command Structure**: All commands work identically
4. **Configuration**: Same configuration options

### Command Examples

All commands work the same, just with the new binary name:

```bash
# Areas
dopa area create --name "Work" --color "#3B82F6"
dopa area list
dopa area update <id> --name "Personal"

# Projects
dopa project create --name "Website" --area-id <id>
dopa project list --area-id <id>

# Tasks
dopa task create --project-id <id> --title "Write docs"
dopa task list --project-id <id>
dopa task complete <id>

# TUI
dopa tui
```

## Developer Impact

### For Contributors

1. **Clone**: Repository URL may change (check with maintainers)
2. **Build**: Use `make build` to create `bin/dopa`
3. **Test**: All tests updated, run `make test`
4. **Imports**: All import statements updated automatically

### For Users

1. **Install**: New binary name `dopa` instead of `projectdb`
2. **Scripts**: Update any scripts that call `projectdb` to use `dopa`
3. **Aliases**: Consider adding: `alias projectdb='dopa'` for transition
4. **Database**: Rename existing database file if desired

## Verification

### Build Verification

```bash
# Clean and rebuild
make clean
make build

# Verify binary
./bin/dopa version
./bin/dopa --help

# Run tests
make test
make lint
```

### Functional Verification

```bash
# Test basic operations
dopa migrate up
dopa area create --name "Test" --color "#FF0000"
dopa area list
dopa tui
```

### Cross-Platform Builds

```bash
make build-all
# Creates:
# - bin/dopa-linux-amd64
# - bin/dopa-darwin-amd64
# - bin/dopa-darwin-arm64
# - bin/dopa-windows-amd64.exe
```

## Rollback Plan

If issues arise, rollback is straightforward:

1. Revert the commit
2. Rebuild with old binary name
3. Restore database: `mv dopadone.db projectdb.db`

## Timeline

- **Planning**: Task created and acceptance criteria defined
- **Implementation**: All phases completed (Core, Docs, Scripts, Verification)
- **Testing**: Manual verification of all functionality
- **Completion**: All acceptance criteria met
- **Documentation**: This document created

## Checklist

- [x] Module path updated in `go.mod`
- [x] All Go imports updated
- [x] Command directory moved (`cmd/projectdb/` → `cmd/dopa/`)
- [x] All Go files compile successfully
- [x] All tests pass
- [x] Makefile updated (binary name, db path, ldflags)
- [x] Build targets verified
- [x] README.md updated
- [x] docs/START_HERE.md updated
- [x] docs/architecture/*.md updated (7 files)
- [x] docs/TUI.md updated
- [x] docs/RELEASE.md updated
- [x] docs/TRANSACTIONS.md updated
- [x] scripts/install.sh updated
- [x] scripts/seed-test-data.sh updated
- [x] scripts/generate-changelog.sh updated
- [x] Manual CLI testing successful
- [x] Manual TUI testing successful
- [x] Database operations verified
- [x] Cross-platform builds tested
- [x] Documentation created (this file)

## Future Considerations

1. **GitHub Repository**: May want to rename repository from `projectdb` to `dopadone`
2. **GitHub Actions**: Update workflow references if repository renamed
3. **Docker Images**: Update image names if published
4. **Package Registries**: Update package names in registries
5. **User Migration Guide**: Create blog post or announcement for users

## Support

For questions or issues related to the rebranding:

1. Check this documentation
2. Review updated README.md
3. Check GitHub Issues for known issues
4. Contact maintainers if needed

---

**Note**: This rebranding represents a significant milestone in the project's evolution, establishing a clear, memorable identity that emphasizes the product's focus on productivity and task completion for developers.
