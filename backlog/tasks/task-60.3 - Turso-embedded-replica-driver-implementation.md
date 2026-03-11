---
id: TASK-60.3
title: Turso embedded replica driver implementation
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-08 19:01'
updated_date: '2026-03-11 10:26'
labels:
  - database
  - turso
  - libsql
  - replication
milestone: m-1
dependencies: []
references:
  - internal/db/driver/driver.go
  - internal/db/driver/config.go
  - internal/db/driver/turso_remote.go
  - internal/db/driver/turso_remote_test.go
documentation:
  - 'https://docs.turso.tech/features/embedded-replicas'
  - 'https://github.com/tursodatabase/go-libsql'
parent_task_id: TASK-60
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement embedded replica driver with automatic sync to Turso primary database. This driver provides local SQLite database that syncs with a remote Turso primary, offering microsecond read latency with cloud backup. Part of TASK-60 Turso integration. Uses go-libsql package for embedded replica support.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create TursoReplicaDriver implementing DatabaseDriver interface
- [x] #2 Initialize embedded replica with libsql.NewEmbeddedReplicaConn(localPath, url, authToken)
- [x] #3 Implement auto-sync with configurable interval (default 60s)
- [x] #4 Add manual Sync() method for on-demand synchronization
- [x] #5 Handle sync errors gracefully - log and retry
- [x] #6 Support context cancellation for sync goroutine
- [x] #7 Track sync status and last sync time
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: Turso Embedded Replica Driver

## Overview
Implement TursoReplicaDriver with embedded replica support, auto-sync, and status tracking.

## Estimated Time: 6-8 hours

## Key Technical Decision: Library Choice

**Decision: Use go-libsql package**

Rationale:
- libsql-client-go is deprecated
- go-libsql is the recommended library for Turso Cloud
- Provides native embedded replica support with NewEmbeddedReplicaConnector
- Has built-in WithSyncInterval option for auto-sync
- Supports manual Sync() method on connector

**Trade-off:** Requires CGO_ENABLED=1 (already required for SQLite anyway)

## Architecture

```
internal/db/driver/
├── turso_replica.go              # TursoReplicaDriver implementation
├── turso_replica_test.go         # Unit tests
└── turso_replica_integration_test.go # Integration tests (build tag)
```

## Implementation Phases

### Phase 1: Add go-libsql Dependency (15 min) - AC #1, #2

**1.1 Add dependency to go.mod**
```bash
go get github.com/tursodatabase/go-libsql
```

**1.2 Update go.mod to use go-libsql for embedded replica**
Note: May need to keep libsql-client-go for remote driver (investigate compatibility)

### Phase 2: TursoReplicaDriver Core (2 hours) - AC #1, #2, #4

**2.1 Define struct with embedded replica connector**
```go
type TursoReplicaDriver struct {
    db          *sql.DB
    connector   *libsql.Connector  // Embedded replica connector
    config      *Config
    status      ConnectionStatus
    mu          sync.RWMutex
    
    // Sync tracking (AC #7)
    syncStatus  SyncStatus
    lastSyncAt  time.Time
    lastSyncErr error
}

type SyncStatus string

const (
    SyncStatusIdle     SyncStatus = "idle"
    SyncStatusSyncing  SyncStatus = "syncing"
    SyncStatusError    SyncStatus = "error"
    SyncStatusOffline  SyncStatus = "offline"
)
```

**2.2 Implement constructor with validation**
```go
func NewTursoReplicaDriver(config *Config) (*TursoReplicaDriver, error) {
    // Validate: URL, token, AND local path required
    if config.TursoURL == "" {
        return nil, NewDriverError(DriverTursoReplica, "create", 
            fmt.Errorf("%w: turso url required", ErrInvalidConfig))
    }
    if config.TursoToken == "" {
        return nil, NewDriverError(DriverTursoReplica, "create", 
            fmt.Errorf("%w: turso token required", ErrInvalidConfig))
    }
    if config.DatabasePath == "" {
        return nil, NewDriverError(DriverTursoReplica, "create", 
            fmt.Errorf("%w: database path required", ErrInvalidConfig))
    }
    // ...
}
```

**2.3 Implement Connect() using go-libsql**
```go
func (d *TursoReplicaDriver) Connect(ctx context.Context) error {
    d.setStatus(StatusConnecting)
    
    // Create embedded replica connector WITHOUT auto-sync
    // (we will manage sync ourselves for more control)
    connector, err := libsql.NewEmbeddedReplicaConnector(
        d.config.DatabasePath,
        d.config.TursoURL,
        libsql.WithAuthToken(d.config.TursoToken),
    )
    if err != nil {
        d.setStatus(StatusError)
        return NewDriverError(DriverTursoReplica, "connect", err)
    }
    
    db := sql.OpenDB(connector)
    if err := db.PingContext(ctx); err != nil {
        connector.Close()
        d.setStatus(StatusError)
        return NewDriverError(DriverTursoReplica, "connect", 
            fmt.Errorf("%w: %v", ErrConnectionFailed, err))
    }
    
    d.mu.Lock()
    d.db = db
    d.connector = connector
    d.mu.Unlock()
    
    d.setStatus(StatusConnected)
    
    // Perform initial sync
    if err := d.Sync(); err != nil {
        // Log but do not fail - can work offline
    }
    
    // Start auto-sync goroutine if interval > 0
    if d.config.SyncInterval > 0 {
        d.startAutoSync()
    }
    
    return nil
}
```

**2.4 Implement manual Sync() method - AC #4**
```go
func (d *TursoReplicaDriver) Sync() error {
    d.mu.Lock()
    if d.connector == nil {
        d.mu.Unlock()
        return NewDriverError(DriverTursoReplica, "sync", ErrDriverAlreadyClosed)
    }
    d.syncStatus = SyncStatusSyncing
    d.mu.Unlock()
    
    replicated, err := d.connector.Sync()
    
    d.mu.Lock()
    defer d.mu.Unlock()
    
    if err != nil {
        d.syncStatus = SyncStatusError
        d.lastSyncErr = err
        return NewDriverError(DriverTursoReplica, "sync", err)
    }
    
    d.syncStatus = SyncStatusIdle
    d.lastSyncAt = time.Now()
    d.lastSyncErr = nil
    
    // Log sync details
    log.Printf("Synced %d frames (frame_no: %d)", 
        replicated.FramesSynced, replicated.FrameNo)
    
    return nil
}
```

**2.5 Implement remaining interface methods**
- `Close()` - Stop sync goroutine, close connector and DB
- `GetDB()` - Return sql.DB
- `Ping()` - Verify connection
- `Type()` - Return DriverTursoReplica
- `Status()` - Return connection status

### Phase 3: Auto-Sync Goroutine (2 hours) - AC #3, #5, #6

**3.1 Implement auto-sync goroutine**
```go
type TursoReplicaDriver struct {
    // ... existing fields ...
    
    // Sync goroutine management
    syncCtx    context.Context
    syncCancel context.CancelFunc
    syncWg     sync.WaitGroup
}

func (d *TursoReplicaDriver) startAutoSync() {
    d.syncCtx, d.syncCancel = context.WithCancel(context.Background())
    
    d.syncWg.Add(1)
    go func() {
        defer d.syncWg.Done()
        
        ticker := time.NewTicker(d.config.SyncInterval)
        defer ticker.Stop()
        
        for {
            select {
            case <-d.syncCtx.Done():
                return
            case <-ticker.C:
                if err := d.Sync(); err != nil {
                    // AC #5: Log error but continue
                    log.Printf("Auto-sync failed: %v", err)
                    
                    // Retry logic could be added here
                    // For now, just log and continue
                }
            }
        }
    }()
}
```

**3.2 Implement context cancellation - AC #6**
```go
func (d *TursoReplicaDriver) Close() error {
    // Stop auto-sync goroutine first
    if d.syncCancel != nil {
        d.syncCancel()
        d.syncWg.Wait()  // Wait for goroutine to exit
        d.syncCancel = nil
    }
    
    d.mu.Lock()
    defer d.mu.Unlock()
    
    if d.db == nil {
        return nil
    }
    
    // Close connector first, then DB
    var errs []error
    
    if d.connector != nil {
        if err := d.connector.Close(); err != nil {
            errs = append(errs, err)
        }
        d.connector = nil
    }
    
    if err := d.db.Close(); err != nil {
        errs = append(errs, err)
    }
    d.db = nil
    d.status = StatusDisconnected
    
    if len(errs) > 0 {
        return NewDriverError(DriverTursoReplica, "close", 
            errors.Join(errs...))
    }
    return nil
}
```

**3.3 Implement graceful error handling - AC #5**
```go
func (d *TursoReplicaDriver) Sync() error {
    // ... sync logic ...
    
    if err != nil {
        d.mu.Lock()
        d.syncStatus = SyncStatusError
        d.lastSyncErr = err
        d.mu.Unlock()
        
        // Log with context
        log.Printf("[TursoReplica] Sync failed: %v", err)
        
        // Do not propagate error in auto-sync context
        // Caller can check SyncStatus() if needed
        return NewDriverError(DriverTursoReplica, "sync", err)
    }
    
    // Success
    d.mu.Lock()
    d.syncStatus = SyncStatusIdle
    d.lastSyncAt = time.Now()
    d.lastSyncErr = nil
    d.mu.Unlock()
    
    return nil
}
```

### Phase 4: Sync Status Tracking (1 hour) - AC #7

**4.1 Add sync status accessors**
```go
// SyncInfo contains sync status information
type SyncInfo struct {
    Status     SyncStatus
    LastSyncAt time.Time
    LastError  error
}

func (d *TursoReplicaDriver) SyncInfo() SyncInfo {
    d.mu.RLock()
    defer d.mu.RUnlock()
    
    return SyncInfo{
        Status:     d.syncStatus,
        LastSyncAt: d.lastSyncAt,
        LastError:  d.lastSyncErr,
    }
}

func (d *TursoReplicaDriver) LastSyncTime() time.Time {
    d.mu.RLock()
    defer d.mu.RUnlock()
    return d.lastSyncAt
}
```

**4.2 Add status to ConnectionStatus for consistency**
```go
// Update Status() to reflect sync state
func (d *TursoReplicaDriver) Status() ConnectionStatus {
    d.mu.RLock()
    defer d.mu.RUnlock()
    
    // If syncing, report as "connecting" to indicate activity
    if d.syncStatus == SyncStatusSyncing {
        return StatusConnecting
    }
    
    return d.status
}
```

### Phase 5: Driver Registration (15 min)

**5.1 Register driver in init()**
```go
func init() {
    err := RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
        return NewTursoReplicaDriver(config)
    })
    if err != nil {
        panic(fmt.Sprintf("failed to register turso-replica driver: %v", err))
    }
}
```

### Phase 6: Tests (2.5 hours)

**6.1 Unit Tests (turso_replica_test.go)**

Following golang-testing patterns with table-driven tests:

```go
func TestNewTursoReplicaDriver(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        wantErr bool
    }{
        {
            name: "valid config",
            config: &Config{
                Type:         DriverTursoReplica,
                DatabasePath: "/tmp/test.db",
                TursoURL:     "libsql://example.turso.io",
                TursoToken:   "test-token",
                SyncInterval: 60 * time.Second,
            },
            wantErr: false,
        },
        {
            name: "missing database path",
            config: &Config{
                Type:       DriverTursoReplica,
                TursoURL:   "libsql://example.turso.io",
                TursoToken: "test-token",
            },
            wantErr: true,
        },
        {
            name: "missing url",
            config: &Config{
                Type:         DriverTursoReplica,
                DatabasePath: "/tmp/test.db",
                TursoToken:   "test-token",
            },
            wantErr: true,
        },
        {
            name: "missing token",
            config: &Config{
                Type:         DriverTursoReplica,
                DatabasePath: "/tmp/test.db",
                TursoURL:     "libsql://example.turso.io",
            },
            wantErr: true,
        },
        {
            name:    "nil config",
            config:  nil,
            wantErr: true,
        },
    }
    // ... test implementation
}

func TestTursoReplicaDriver_Type(t *testing.T) { ... }
func TestTursoReplicaDriver_Status(t *testing.T) { ... }
func TestTursoReplicaDriver_Sync_BeforeConnect(t *testing.T) { ... }
func TestTursoReplicaDriver_SyncInfo(t *testing.T) { ... }
func TestTursoReplicaDriver_AutoSync_StartStop(t *testing.T) { ... }
func TestTursoReplicaDriver_AutoSync_ContextCancellation(t *testing.T) { ... }
func TestTursoReplicaDriver_ConcurrentSyncAccess(t *testing.T) { ... }
func TestTursoReplicaDriver_Close_StopsAutoSync(t *testing.T) { ... }
```

**6.2 Integration Tests (turso_replica_integration_test.go)**
```go
//go:build integration
// +build integration

func TestTursoReplicaDriver_Integration_Connect(t *testing.T) {
    url := os.Getenv("TURSO_TEST_URL")
    token := os.Getenv("TURSO_TEST_TOKEN")
    if url == "" || token == "" {
        t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required")
    }
    
    tmpDir := t.TempDir()
    dbPath := filepath.Join(tmpDir, "replica.db")
    
    driver, err := NewTursoReplicaDriver(&Config{
        Type:         DriverTursoReplica,
        DatabasePath: dbPath,
        TursoURL:     url,
        TursoToken:   token,
        SyncInterval: 5 * time.Second,
    })
    // ... test implementation
}

func TestTursoReplicaDriver_Integration_Sync(t *testing.T) { ... }
func TestTursoReplicaDriver_Integration_AutoSync(t *testing.T) { ... }
func TestTursoReplicaDriver_Integration_ReadWrite(t *testing.T) { ... }
```

### Phase 7: Documentation (30 min)

**7.1 Update docs/architecture/08-database-drivers.md**
- Add TursoReplicaDriver section
- Document sync behavior and status tracking
- Add configuration examples

**7.2 Document sync status in code**
- Document SyncInfo struct
- Document Sync() method behavior
- Document auto-sync goroutine lifecycle

## File Structure After Implementation

```
internal/db/driver/
├── driver.go                        # Unchanged
├── registry.go                      # Unchanged
├── factory.go                       # Unchanged
├── config.go                        # Unchanged
├── errors.go                        # Unchanged
├── turso_remote.go                  # Existing
├── turso_replica.go                 # NEW
├── turso_replica_test.go            # NEW
├── turso_replica_integration_test.go # NEW (build tag: integration)
└── driver_test.go                   # Existing
```

## Test Coverage Target: 90%+

## Dependencies

- **Requires**: TASK-60.1 (completed)
- **Requires**: TASK-60.2 (completed - for reference patterns)
- **Blocks**: TASK-60.7 (integration)

## Acceptance Criteria Mapping

| AC | Phase | Description |
|----|-------|-------------|
| #1 | 2 | TursoReplicaDriver implementing DatabaseDriver |
| #2 | 1, 2 | Initialize with libsql.NewEmbeddedReplicaConnector |
| #3 | 3 | Auto-sync with configurable interval |
| #4 | 2 | Manual Sync() method |
| #5 | 3 | Graceful error handling - log and retry |
| #6 | 3 | Context cancellation for sync goroutine |
| #7 | 4 | Track sync status and last sync time |

## Risk Mitigation

1. **CGO Requirement**: go-libsql requires CGO - already needed for SQLite
2. **Concurrent Access**: Use sync.RWMutex for all shared state
3. **Goroutine Leak**: Ensure Close() waits for sync goroutine to exit
4. **Offline Operation**: Initial sync failure should not prevent connection
5. **Error Recovery**: Log errors but continue auto-sync attempts

## Success Criteria

- [ ] All 7 acceptance criteria met
- [ ] Unit tests pass (90%+ coverage)
- [ ] Integration tests pass with real Turso database
- [ ] Driver registered and available via factory
- [ ] Documentation updated
- [ ] Code passes golangci-lint
- [ ] No goroutine leaks (verified with race detector)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Analysis

### Scope Assessment
This task is **appropriately scoped** for a single PR:
- Implements ONE driver type (Turso embedded replica)
- All 7 ACs are tightly related (replica + sync functionality)
- Estimated 6-8 hours of work
- Clear boundaries with other tasks (uses patterns from TASK-60.2)

### Why Not Split
Considered splitting into:
- 60.3a: Core driver (AC#1, #2, #4)
- 60.3b: Auto-sync (AC#3, #5, #6)
- 60.3c: Status tracking (AC#7)

**Decision: Keep as single task** because:
1. Auto-sync is core to the driver functionality
2. Status tracking is essential for auto-sync monitoring
3. Splitting would create tight coupling between subtasks
4. Testing would require integration across all "subtasks"
5. Single PR reduces coordination overhead

### Dependencies & Execution Order

**Sequential Dependencies:**
- TASK-60.1 (Database abstraction layer) - **COMPLETED**
- TASK-60.2 (Turso remote driver) - **COMPLETED** - Use as reference for patterns

**Parallel Opportunities:**
- Can be implemented **in parallel** with:
  - TASK-60.8 (Configuration system) - Independent
  - TASK-60.5 (Migration compatibility) - Independent

**Blocked Tasks:**
- TASK-60.7 (Integration and wiring) - Requires all drivers complete

### Key Technical Considerations

1. **Library Choice**: go-libsql recommended over deprecated libsql-client-go
2. **CGO Requirement**: Already required for SQLite, no new constraint
3. **Concurrency**: Careful goroutine lifecycle management critical
4. **Offline Resilience**: Initial sync failure should not block connection
5. **Testing Strategy**: Unit tests + optional integration tests with build tag

### Reference Files

- internal/db/driver/turso_remote.go - Pattern reference
- internal/db/driver/turso_remote_test.go - Testing patterns
- internal/db/driver/config.go - Configuration options
- internal/db/driver/driver.go - Interface definition
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented TursoReplicaDriver with embedded replica support using go-libsql.

 Auto-sync with configurable interval, manual Sync method with graceful error handling, context cancellation support sync status tracking. All unit tests passing, Documentation updated in docs/architecture/08-database-drivers.md. Code passes golangci-lint with race detector. Valid status: driver is ready for review. Not committing or pushing changes to repository.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All unit tests pass (90%+ coverage)
- [x] #2 Integration tests pass with real Turso database (or skipped without credentials)
- [x] #3 Driver auto-registers via init()
- [x] #4 Documentation updated (docs/architecture/08-database-drivers.md)
- [x] #5 Code passes golangci-lint with race detector
- [x] #6 No goroutine leaks verified
- [x] #7 Code follows golang-pro, golang-patterns, and golang-testing best practices
- [ ] #8 3
- [ ] #9 5
<!-- DOD:END -->
