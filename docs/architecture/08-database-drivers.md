# Database Driver Architecture

This document describes the database driver abstraction layer that enables multiple database backends (SQLite, Turso remote, Turso embedded replica) while maintaining backward compatibility with existing code.

## Overview

The driver abstraction layer follows Go best practices with:
- **Interface-based design**: `DatabaseDriver` interface for all database implementations
- **Factory pattern**: Driver creation based on configuration
- **Dependency injection**: Drivers are injected, not created internally
- **Thread-safe registry**: Concurrent-safe driver registration and lookup

## Architecture

```
internal/db/driver/
├── driver.go                      # DatabaseDriver interface definition
├── registry.go                    # Driver registry for registration
├── factory.go                     # Factory pattern for driver creation
├── config.go                      # Configuration types and defaults
├── errors.go                      # Driver-specific error types
├── turso_remote.go                # Turso remote driver implementation
├── turso_remote_test.go           # Unit tests for Turso remote
├── turso_remote_integration_test.go # Integration tests (requires real DB)
├── turso_replica.go               # Turso embedded replica driver implementation
├── turso_replica_test.go          # Unit tests for Turso replica
├── turso_replica_integration_test.go # Integration tests (requires real DB)
└── driver_test.go                 # Core unit tests
```

### Data Flow

```
Configuration → Factory → Registry → Driver → *sql.DB
```

## DatabaseDriver Interface

All database drivers must implement the `DatabaseDriver` interface:

```go
type DatabaseDriver interface {
    // Connect establishes connection to the database
    Connect(ctx context.Context) error
    
    // Close closes the database connection
    Close() error
    
    // GetDB returns the underlying *sql.DB for compatibility
    GetDB() *sql.DB
    
    // Ping verifies the connection is alive
    Ping(ctx context.Context) error
    
    // Type returns the driver type identifier
    Type() DriverType
    
    // Status returns current connection status
    Status() ConnectionStatus
}
```

### Key Design Decisions

1. **sql.DB Compatibility**: The `GetDB()` method returns `*sql.DB` to maintain compatibility with sqlc-generated code
2. **Context Support**: All connection operations accept `context.Context` for cancellation and timeouts
3. **Status Tracking**: Each driver tracks its connection status for monitoring

## Driver Types

| Type | Constant | Description |
|------|----------|-------------|
| SQLite | `DriverSQLite` | Local SQLite database |
| Turso Remote | `DriverTursoRemote` | Remote Turso database (cloud) |
| Turso Replica | `DriverTursoReplica` | Embedded replica with sync |

## Turso Remote Driver

The `TursoRemoteDriver` connects directly to a remote Turso database using go-libsql.

### Features

- **Direct cloud connection**: Uses local temp file for connection handling
- **Authentication**: Uses auth tokens for secure connections
- **Fail-fast behavior**: Immediate connection validation
- **Retry logic**: Configurable retry attempts with intervals
- **Timeout support**: Connection timeout configuration

## Turso Embedded Replica Driver

The `TursoReplicaDriver` provides an embedded replica with automatic sync to a Turso primary database. This driver offers microsecond read latency with cloud backup.

### Features

- **Local SQLite replica**: Fast local reads with cloud backup
- **Auto-sync**: Configurable sync interval (default: 60s)
- **Manual sync**: On-demand synchronization via `Sync()` method
- **Sync status tracking**: Monitor sync state, last sync time, and errors
- **Graceful error handling**: Logs errors and continues retry attempts
- **Context cancellation**: Clean shutdown of sync goroutines
- **Offline operation**: Can work offline, syncs when connection available

### Sync Status Types

| Status | Description |
|--------|-------------|
| `SyncStatusIdle` | No sync in progress |
| `SyncStatusSyncing` | Sync operation in progress |
| `SyncStatusError` | Last sync failed |
| `SyncStatusOffline` | Driver not connected |

### Usage

```go
driver, err := driver.NewTursoReplicaDriver(&driver.Config{
    Type:         driver.DriverTursoReplica,
    DatabasePath: "/path/to/local/replica.db",
    TursoURL:     "libsql://your-database.turso.io",
    TursoToken:   "your-auth-token",
    SyncInterval: 60 * time.Second,
})
if err != nil {
    return err
}

ctx := context.Background()
if err := driver.Connect(ctx); err != nil {
    return err
}
defer driver.Close()

// Check sync status
info := driver.SyncInfo()
fmt.Printf("Status: %s, Last sync: %v\n", info.Status, info.LastSyncAt)

// Manual sync
if err := driver.Sync(); err != nil {
    log.Printf("Sync failed: %v", err)
}
```

### SyncInfo Structure

```go
type SyncInfo struct {
    Status     SyncStatus  // Current sync status
    LastSyncAt time.Time   // Time of last successful sync
    LastError  error       // Last sync error, if any
}
```

### Usage

```go
driver, err := driver.NewDriver(
    driver.WithDriverType(driver.DriverTursoRemote),
    driver.WithTurso("libsql://your-database.turso.io", "your-auth-token"),
    driver.WithConnectTimeout(10*time.Second),
    driver.WithMaxRetries(3),
    driver.WithRetryInterval(1*time.Second),
)
if err != nil {
    return err
}

ctx := context.Background()
if err := driver.Connect(ctx); err != nil {
    return err
}
defer driver.Close()

// Use the database
db := driver.GetDB()
rows, err := db.QueryContext(ctx, "SELECT * FROM tasks")
```

### Configuration Options

| Option | Default | Description |
|--------|---------|-------------|
| `WithConnectTimeout(d)` | 10s | Maximum time for connection attempt |
| `WithMaxRetries(n)` | 3 | Number of retry attempts |
| `WithRetryInterval(d)` | 1s | Time between retries |

### Auto-Registration

The driver auto-registers on package import via `init()`:

```go
func init() {
    err := RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
        return NewTursoRemoteDriver(config)
    })
    if err != nil {
        panic(fmt.Sprintf("failed to register turso-remote driver: %v", err))
    }
}
```

## Configuration

Use functional options to configure drivers:

```go
// SQLite configuration
driver, err := driver.NewDriver(
    driver.WithDriverType(driver.DriverSQLite),
    driver.WithDatabasePath("/path/to/database.db"),
)

// Turso remote configuration
driver, err := driver.NewDriver(
    driver.WithDriverType(driver.DriverTursoRemote),
    driver.WithTurso("https://your-database.turso.io", "your-auth-token"),
)

// Turso replica configuration
driver, err := driver.NewDriver(
    driver.WithDriverType(driver.DriverTursoReplica),
    driver.WithDatabasePath("/path/to/local/replica.db"),
    driver.WithTurso("https://your-database.turso.io", "your-auth-token"),
    driver.WithSyncInterval(30 * time.Second),
)
```

### Configuration Options

| Option | Description |
|--------|-------------|
| `WithDriverType(dt)` | Set the driver type |
| `WithDatabasePath(path)` | Set local database path |
| `WithTurso(url, token)` | Set Turso URL and auth token |
| `WithSyncInterval(d)` | Set sync interval for replica |
| `WithConnectTimeout(d)` | Set connection timeout (default: 10s) |
| `WithMaxRetries(n)` | Set max retry attempts (default: 3) |
| `WithRetryInterval(d)` | Set time between retries (default: 1s) |

## Driver Registry

Drivers must be registered before use:

```go
// Register a new driver
err := driver.RegisterDriver(driver.DriverSQLite, sqliteDriverFactory)
if err != nil {
    // Handle registration error
}

// Get registered driver types
drivers := driver.RegisteredDrivers()

// Check if driver is registered
factory, ok := driver.GetFactory(driver.DriverSQLite)
```

### Thread Safety

The registry is thread-safe and can be accessed concurrently from multiple goroutines.

## Error Handling

The driver package provides custom error types:

```go
// Sentinel errors
var (
    ErrDriverNotRegistered = errors.New("driver not registered")
    ErrConnectionFailed    = errors.New("connection failed")
    ErrInvalidConfig       = errors.New("invalid configuration")
    ErrDriverAlreadyClosed = errors.New("driver already closed")
)

// DriverError for detailed error context
type DriverError struct {
    Driver DriverType
    Op     string
    Err    error
}
```

### Error Checking

```go
import "errors"

driver, err := driver.NewDriver(opts...)
if err != nil {
    if errors.Is(err, driver.ErrDriverNotRegistered) {
        // Driver not registered
    }
    
    var driverErr *driver.DriverError
    if errors.As(err, &driverErr) {
        // Access detailed error info
        log.Printf("Driver: %s, Operation: %s", driverErr.Driver, driverErr.Op)
    }
}
```

## Implementing a New Driver

To implement a new database driver:

### 1. Create the Driver Struct

```go
type MyDriver struct {
    db     *sql.DB
    status driver.ConnectionStatus
    config *driver.Config
}
```

### 2. Implement the Interface

```go
func (d *MyDriver) Connect(ctx context.Context) error {
    db, err := sql.Open("my-driver", d.config.DatabasePath)
    if err != nil {
        d.status = driver.StatusError
        return driver.NewDriverError(d.Type(), "connect", err)
    }
    d.db = db
    d.status = driver.StatusConnected
    return nil
}

func (d *MyDriver) Close() error {
    if d.db == nil {
        return driver.ErrDriverAlreadyClosed
    }
    err := d.db.Close()
    d.db = nil
    d.status = driver.StatusDisconnected
    return err
}

func (d *MyDriver) GetDB() *sql.DB {
    return d.db
}

func (d *MyDriver) Ping(ctx context.Context) error {
    if d.db == nil {
        return driver.ErrConnectionFailed
    }
    return d.db.PingContext(ctx)
}

func (d *MyDriver) Type() driver.DriverType {
    return "my-driver"
}

func (d *MyDriver) Status() driver.ConnectionStatus {
    return d.status
}
```

### 3. Create Factory Function

```go
func NewMyDriver(config *driver.Config) (driver.DatabaseDriver, error) {
    return &MyDriver{
        status: driver.StatusDisconnected,
        config: config,
    }, nil
}
```

### 4. Register the Driver

```go
func init() {
    err := driver.RegisterDriver("my-driver", NewMyDriver)
    if err != nil {
        panic(err)
    }
}
```

## Backward Compatibility

The abstraction layer maintains full backward compatibility:

1. **Existing `Connect()` function** remains unchanged in `internal/cli/db.go`
2. **sqlc-generated code** works without modifications (uses `*sql.DB`)
3. **TransactionManager** continues to work (accepts `*sql.DB`)
4. **All existing tests** pass without modifications

### Migration Path

Existing code can continue using the old API:

```go
// Old API (still works)
db, err := cli.Connect("/path/to/database.db")

// New API (for multi-driver support)
drv, err := cli.ConnectWithDriver(
    driver.WithDriverType(driver.DriverSQLite),
    driver.WithDatabasePath("/path/to/database.db"),
)
if err != nil {
    return err
}
db := drv.GetDB() // Get underlying *sql.DB
```

## Testing

The driver package includes comprehensive tests:

- **Interface contract tests**: Verify implementations satisfy the interface
- **Registry tests**: Test registration, lookup, and thread-safety
- **Factory tests**: Test driver creation via factory
- **Configuration tests**: Test functional options
- **Error tests**: Test error types and wrapping
- **Turso remote tests**: Connection, retry logic, timeout handling
- **Integration tests**: Real database tests (build tag: `integration`)

Run tests with:

```bash
# Unit tests
go test ./internal/db/driver/... -v
go test ./internal/db/driver/... -cover

# Integration tests (requires TURSO_TEST_URL and TURSO_TEST_TOKEN)
TURSO_TEST_URL=libsql://your-db.turso.io \
TURSO_TEST_TOKEN=your-token \
go test ./internal/db/driver/... -tags=integration -v
```

## Related Tasks

- **TASK-60.1**: Database abstraction layer (completed)
- **TASK-60.2**: Turso remote driver implementation (completed)
- **TASK-60.3**: Turso embedded replica driver implementation (completed)
- **TASK-60.7**: Integration and wiring
- **TASK-60.8**: Configuration system

## References

- [Architecture Overview](01-overview.md)
- [Service Layer](03-service-layer.md)
- [Repository Layer](05-repository-layer.md)
