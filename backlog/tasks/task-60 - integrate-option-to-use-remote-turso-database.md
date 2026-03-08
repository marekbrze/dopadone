---
id: TASK-60
title: integrate option to use remote turso database
status: In Progress
assignee:
  - '@assistant'
created_date: '2026-03-07 15:49'
updated_date: '2026-03-08 19:03'
labels:
  - feature
  - database
  - turso
  - libsql
dependencies: []
references:
  - backlog/tasks/task-60.6 - Integration-tests-for-database-modes.md
  - >-
    backlog/tasks/task-60.7 -
    Integration-and-refactoring-wire-up-database-abstraction.md
  - backlog/tasks/task-60.8 - Configuration-system-for-database-modes.md
  - >-
    backlog/tasks/task-60.9 -
    Documentation-Turso-setup-and-configuration-guide.md
documentation:
  - 'https://docs.turso.tech/features/embedded-replicas'
  - 'https://www.meetgor.com/turso-libsql-embedded-replicas-golang/'
  - 'https://github.com/tursodatabase/libsql-client-go'
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add support for remote Turso database with automatic mode detection. Support three connection modes: local SQLite (current), remote Turso, and embedded replica (local with cloud sync). Maintain full backward compatibility - local SQLite remains the default. Configuration via CLI flags, environment variables, or config file with precedence order: CLI > env > config.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Support three database modes: local SQLite, remote Turso, embedded replica with auto-sync
- [ ] #2 Implement configuration via CLI flags (--turso-url, --turso-token), environment variables (TURSO_DATABASE_URL, TURSO_AUTH_TOKEN), and config file with precedence CLI > env > config
- [ ] #3 Auto-detect database mode based on configuration presence
- [ ] #4 Implement auto-sync with configurable interval for embedded replica mode
- [ ] #5 Ensure existing goose migrations work with libSQL
- [ ] #6 Migration sync: run migrations locally, sync schema to Turso via embedded replica
- [ ] #7 Maintain backward compatibility: local SQLite is default, Turso is opt-in
- [ ] #8 Fail fast on Turso connection failures when remote/replica mode is configured
- [ ] #9 Add connection status indicator in TUI showing connected/syncing/offline/local-only status
- [ ] #10 Add database connection abstraction layer to support multiple drivers
- [ ] #11 Update documentation with Turso setup guide and configuration examples
- [ ] #12 Add integration tests for all three database modes
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for Turso Database Integration

This task is split into 11 subtasks that can be worked on in parallel or sequentially based on dependencies.

## Task Breakdown & Dependencies

### Phase 1: Foundation (Sequential - Must be done first)

**TASK-60.1: Database abstraction layer and driver interface** [HIGH PRIORITY]
- Creates the core DatabaseDriver interface
- Defines factory pattern for driver creation
- Must be completed FIRST - all other drivers depend on this
- Estimated: 4-6 hours
- Files: internal/db/driver/interface.go, internal/db/driver/factory.go

**TASK-60.8: Configuration system for database modes** [HIGH PRIORITY]
- Implements config precedence: CLI > env > config > defaults
- Adds CLI flags: --turso-url, --turso-token, --db-mode, --sync-interval
- Can be done IN PARALLEL with 60.1
- Estimated: 3-4 hours
- Files: internal/config/database.go, cmd/dopa/main.go (flags)

### Phase 2: Driver Implementations (Parallel - After Phase 1)

All three drivers can be implemented IN PARALLEL after Phase 1 is complete.

**TASK-60.2: SQLite driver implementation** [HIGH PRIORITY]
- Implements local SQLite driver (preserves existing functionality)
- Depends on: TASK-60.1
- Estimated: 2-3 hours
- Files: internal/db/driver/sqlite.go

**TASK-60.3: Turso remote driver implementation** [HIGH PRIORITY]
- Implements direct remote Turso connection
- Depends on: TASK-60.1, TASK-60.8
- Estimated: 4-5 hours
- Files: internal/db/driver/turso_remote.go

**TASK-60.4: Turso embedded replica driver** [HIGH PRIORITY]
- Implements embedded replica with auto-sync
- Depends on: TASK-60.1, TASK-60.8
- Estimated: 6-8 hours (most complex)
- Files: internal/db/driver/turso_replica.go

### Phase 3: Integration Features (Parallel - After Phase 2)

**TASK-60.8: Database mode auto-detection** [MEDIUM PRIORITY]
- Auto-detects mode based on config presence
- Depends on: TASK-60.8 (config), TASK-60.2, TASK-60.3, TASK-60.4
- Estimated: 2-3 hours
- Files: internal/db/driver/detector.go

**TASK-60.5: Migration compatibility with libSQL** [MEDIUM PRIORITY]
- Ensures goose migrations work with libSQL
- Depends on: TASK-60.3, TASK-60.4
- Can be done IN PARALLEL with 60.8
- Estimated: 3-4 hours
- Files: internal/migrate/libsql.go, tests

**TASK-60.6: Connection status indicator in TUI** [MEDIUM PRIORITY]
- Adds status indicator to TUI
- Depends on: TASK-60.4 (replica driver for sync status)
- Can be done IN PARALLEL with 60.5, 60.8
- Estimated: 3-4 hours
- Files: internal/tui/status.go, internal/tui/view.go

### Phase 4: Integration & Wiring (Sequential - After Phase 3)

**TASK-60.7: Integration and refactoring** [HIGH PRIORITY]
- Wires up all drivers into main application
- Refactors GetDB() and GetServices()
- Depends on: ALL previous tasks
- Must be done AFTER all drivers are implemented
- Estimated: 4-5 hours
- Files: cmd/dopa/main.go, internal/cli/db.go

### Phase 5: Testing & Documentation (Parallel - After Phase 4)

**TASK-60.7: Integration tests** [MEDIUM PRIORITY]
- Tests all three database modes
- Depends on: TASK-60.7 (integration)
- Can be done IN PARALLEL with documentation
- Estimated: 6-8 hours
- Files: internal/db/driver/*_test.go, integration tests

**TASK-60.9: Documentation** [MEDIUM PRIORITY]
- Turso setup guide and configuration examples
- Can be done IN PARALLEL with testing
- Estimated: 3-4 hours
- Files: docs/TURSO.md, docs/architecture/08-database-drivers.md

## Total Estimated Time
- Sequential work (critical path): 20-28 hours
- Parallel work: 15-20 hours
- **Total: 35-48 hours** (can be reduced to 25-35 hours with parallel execution)

## Dependency Graph
```
60.1 (interface) ────┬─────────────┬────────────┐
                     ↓             ↓            ↓
60.8 (config) ──┬─→ 60.2 (SQLite) 60.3 (remote) 60.4 (replica)
                │    │             │            │
                │    └─────┬───────┴────────────┤
                │          ↓                    ↓
                └────→ 60.8 (detect)       60.5 (migration)
                           │                    │
                           └────────┬───────────┤
                                    ↓           ↓
                              60.6 (TUI)   60.7 (integration)
                                    │           │
                                    └─────┬─────┘
                                          ↓
                                  60.7 (tests) + 60.9 (docs)
```

## Implementation Order

**Option 1: Sequential (safest)**
1. TASK-60.1 → 2. TASK-60.8 → 3. TASK-60.2 → 4. TASK-60.3 → 5. TASK-60.4 → 6. TASK-60.8 → 7. TASK-60.5 → 8. TASK-60.6 → 9. TASK-60.7 → 10. TASK-60.7 → 11. TASK-60.9

**Option 2: Optimized (faster, recommended)**
Day 1: TASK-60.1 + TASK-60.8 (parallel)
Day 2: TASK-60.2 + TASK-60.3 (parallel)
Day 3: TASK-60.4
Day 4: TASK-60.8 + TASK-60.5 + TASK-60.6 (parallel)
Day 5: TASK-60.7
Day 6: TASK-60.7 + TASK-60.9 (parallel)

## Key Technical Decisions

1. **Driver Interface Design**: Follow Go standard library sql.DB compatibility
2. **libSQL SDK**: Use github.com/tursodatabase/libsql-client-go
3. **Sync Strategy**: Background goroutine with configurable interval for replica mode
4. **Error Handling**: Fail-fast on remote/replica mode connection failures
5. **Backward Compatibility**: Local SQLite remains default, zero config changes needed
6. **Configuration**: Viper or custom config with YAML support
7. **Testing**: Use testcontainers or mock server for CI

## Success Criteria

- ✅ All 12 acceptance criteria from main task met
- ✅ All existing tests pass
- ✅ New integration tests pass for all modes
- ✅ Documentation complete
- ✅ Zero breaking changes to CLI/TUI
- ✅ Performance impact < 5% for local SQLite mode
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Implementation Strategy

### Task Decomposition Rationale

This task has been decomposed into 11 subtasks based on:
1. **Single Responsibility Principle**: Each task focuses on one component
2. **Parallel Execution**: Independent tasks can be worked on simultaneously
3. **Incremental Integration**: Changes can be tested in isolation
4. **Clear Dependencies**: Dependency graph shows execution order

### Architecture Overview

```
┌─────────────────────────────────────────────┐
│         Application Layer (main.go)         │
│  GetDB() → Driver Factory → DatabaseDriver  │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│         Driver Abstraction Layer            │
│         (internal/db/driver/)               │
│  - DatabaseDriver interface                 │
│  - DriverFactory                            │
│  - Configuration                            │
└─────────────────────────────────────────────┘
                    ↓
┌──────────────┬──────────────┬───────────────┐
│ SQLiteDriver │ TursoRemote  │ TursoReplica  │
│ (modernc)    │ (libsql)     │ (libsql+sync) │
└──────────────┴──────────────┴───────────────┘
```

### Key Design Patterns

1. **Strategy Pattern**: DatabaseDriver interface with multiple implementations
2. **Factory Pattern**: Driver creation based on configuration
3. **Dependency Injection**: Drivers injected into services
4. **Observer Pattern**: TUI observes connection status changes

### Testing Strategy

- **Unit Tests**: Each driver tested independently with mocks
- **Integration Tests**: End-to-end tests for all three modes
- **Backward Compatibility Tests**: Ensure existing functionality preserved
- **Performance Tests**: Benchmark connection overhead

### Risk Mitigation

1. **Backward Compatibility**: Default to local SQLite, no breaking changes
2. **Fail-Fast**: Remote/replica modes fail immediately on connection errors
3. **Incremental Rollout**: Each subtask can be merged independently
4. **Feature Flag**: Consider adding feature flag for Turso integration

### Performance Considerations

- **Local SQLite**: No performance impact (same code path)
- **Remote Turso**: Network latency added (acceptable for cloud use)
- **Embedded Replica**: Best of both worlds - local read speed + cloud sync

### Configuration Example

```yaml
database:
  mode: replica  # local | remote | replica (auto-detected if not set)
  local_path: ./dopadone.db
  turso:
    url: libsql://your-db.turso.io
    auth_token: your-token-here
  replica:
    sync_interval: 60s  # sync every 60 seconds
```

### CLI Usage Examples

```bash
# Local SQLite (default)
dopa --db ./mydb.db

# Remote Turso
dopa --turso-url libsql://db.turso.io --turso-token xxx

# Embedded Replica
dopa --db ./local.db --turso-url libsql://db.turso.io --turso-token xxx

# Environment variables
export TURSO_DATABASE_URL=libsql://db.turso.io
export TURSO_AUTH_TOKEN=xxx
dopa
```

### Migration Considerations

- Goose migrations run locally first
- Embedded replica syncs schema to Turso
- Remote mode requires manual migration on Turso side
- Consider adding `dopa migrate turso` command
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All tests pass including new integration tests
- [ ] #2 Documentation updated with Turso setup instructions
- [ ] #3 Code follows golang-pro, golang-patterns, and golang-testing best practices
- [ ] #4 No breaking changes to existing CLI/TUI interface
<!-- DOD:END -->
