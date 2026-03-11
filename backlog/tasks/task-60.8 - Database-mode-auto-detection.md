---
id: TASK-60.8
title: Database mode auto-detection
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-08 19:02'
updated_date: '2026-03-11 10:44'
labels:
  - database
  - configuration
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement automatic database mode detection based on configuration presence. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Auto-detect mode: if turso-url present and db-mode not set, use remote
- [x] #2 Auto-detect mode: if turso-url + local path set, use embedded replica
- [x] #3 Auto-detect mode: if only db-path set, use local SQLite (default)
- [x] #4 Add validation for required configuration per mode
- [x] #5 Log detected mode at startup for visibility
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan: Database Mode Auto-Detection

## Overview
Implement automatic database mode detection based on configuration presence, with comprehensive validation and logging.

## Estimated Time: 3-4 hours

## Scope Assessment
**Appropriately scoped** - Focused on detection logic, configuration loading, and application wiring. All 5 ACs are tightly related to the same feature.

## Architecture

```
internal/db/driver/
├── detector.go           # NEW: Auto-detection logic
├── detector_test.go      # NEW: Unit tests for detection
└── config.go             # MODIFIED: Add config loading helpers

cmd/dopa/
├── main.go               # MODIFIED: Add CLI flags, wire detection
└── config.go             # NEW: Config loading from CLI/env/file
```

## Implementation Phases

### Phase 1: Configuration Loading (1 hour) - AC #4

**1.1 Create cmd/dopa/config.go for configuration management**
```go
type Config struct {
    DatabasePath string
    TursoURL     string
    TursoToken   string
    DBMode       string // "local", "remote", "replica", or "" (auto-detect)
    SyncInterval time.Duration
}

func LoadConfig() *Config {
    cfg := &Config{
        DatabasePath: getDatabasePath(),
        TursoURL:     getTursoURL(),
        TursoToken:   getTursoToken(),
        DBMode:       getDBMode(),
        SyncInterval: getSyncInterval(),
    }
    return cfg
}

// Precedence: CLI > env > default
func getDatabasePath() string { ... }
func getTursoURL() string { ... }
func getTursoToken() string { ... }
```

**1.2 Add CLI flags to main.go**
```go
var (
    // Existing flags
    dbPath string
    
    // New flags
    tursoURL     string
    tursoToken   string
    dbMode       string
    syncInterval string
)

func init() {
    rootCmd.PersistentFlags().StringVar(&tursoURL, "turso-url", "", "Turso database URL (env: TURSO_DATABASE_URL)")
    rootCmd.PersistentFlags().StringVar(&tursoToken, "turso-token", "", "Turso auth token (env: TURSO_AUTH_TOKEN)")
    rootCmd.PersistentFlags().StringVar(&dbMode, "db-mode", "", "Database mode: local|remote|replica (auto-detect if not set)")
    rootCmd.PersistentFlags().StringVar(&syncInterval, "sync-interval", "60s", "Sync interval for embedded replica mode")
}
```

**1.3 Support environment variables**
```go
func getTursoURL() string {
    if tursoURL != "" {
        return tursoURL
    }
    return os.Getenv("TURSO_DATABASE_URL")
}

func getTursoToken() string {
    if tursoToken != "" {
        return tursoToken
    }
    return os.Getenv("TURSO_AUTH_TOKEN")
}
```

### Phase 2: Auto-Detection Logic (1.5 hours) - AC #1, #2, #3

**2.1 Create internal/db/driver/detector.go**
```go
type DetectionResult struct {
    Type    DriverType
    Reason  string
    Config  *Config
}

// DetectMode determines driver type based on configuration
// AC #1, #2, #3: Auto-detect based on presence of turso-url and db-path
func DetectMode(cfg *Config) DetectionResult {
    // AC #3: Only db-path set → use local SQLite (default)
    if cfg.TursoURL == "" && cfg.TursoToken == "" && cfg.DatabasePath != "" {
        return DetectionResult{
            Type:   DriverSQLite,
            Reason: "local SQLite (no Turso configuration found)",
        }
    }
    
    // AC #2: turso-url + local path set → use embedded replica
    if cfg.TursoURL != "" && cfg.DatabasePath != "" {
        return DetectionResult{
            Type:   DriverTursoReplica,
            Reason: "embedded replica (Turso URL + local path configured)",
        }
    }
    
    // AC #1: turso-url present, db-mode not set → use remote
    if cfg.TursoURL != "" && cfg.DatabasePath == "" {
        return DetectionResult{
            Type:   DriverTursoRemote,
            Reason: "remote Turso (Turso URL configured without local path)",
        }
    }
    
    // Default fallback
    return DetectionResult{
        Type:   DriverSQLite,
        Reason: "local SQLite (default fallback)",
    }
}

// ValidateConfig validates configuration based on detected mode
// AC #4: Validation for required configuration per mode
func ValidateConfig(cfg *Config, detectedType DriverType) error {
    switch detectedType {
    case DriverSQLite:
        if cfg.DatabasePath == "" {
            return NewDriverError(detectedType, "validate", 
                fmt.Errorf("%w: database path required", ErrInvalidConfig))
        }
    
    case DriverTursoRemote:
        if cfg.TursoURL == "" || cfg.TursoToken == "" {
            return NewDriverError(detectedType, "validate", 
                fmt.Errorf("%w: turso URL and token required", ErrInvalidConfig))
        }
    
    case DriverTursoReplica:
        if cfg.TursoURL == "" || cfg.TursoToken == "" || cfg.DatabasePath == "" {
            return NewDriverError(detectedType, "validate", 
                fmt.Errorf("%w: turso URL, token, and database path required", ErrInvalidConfig))
        }
    
    default:
        return NewDriverError(detectedType, "validate", ErrDriverNotRegistered)
    }
    
    return nil
}
```

**2.2 Support explicit mode override**
```go
// ParseExplicitMode parses user-specified mode
func ParseExplicitMode(mode string) (DriverType, error) {
    switch mode {
    case "", "auto":
        return "", nil // Signal auto-detect
    case "local":
        return DriverSQLite, nil
    case "remote":
        return DriverTursoRemote, nil
    case "replica":
        return DriverTursoReplica, nil
    default:
        return "", fmt.Errorf("invalid database mode: %s (valid: local, remote, replica, auto)", mode)
    }
}

// DetectOrExplicitMode returns driver type, preferring explicit mode if set
func DetectOrExplicitMode(cfg *Config) (DetectionResult, error) {
    explicitMode, err := ParseExplicitMode(cfg.DBMode)
    if err != nil {
        return DetectionResult{}, err
    }
    
    if explicitMode != "" {
        return DetectionResult{
            Type:   explicitMode,
            Reason: fmt.Sprintf("explicit mode: %s", cfg.DBMode),
        }, nil
    }
    
    return DetectMode(cfg), nil
}
```

### Phase 3: Logging and Integration (1 hour) - AC #5

**3.1 Add logging in main.go**
```go
func GetDB() (*sql.DB, error) {
    cfg := LoadConfig()
    
    // Detect mode
    result, err := driver.DetectOrExplicitMode(&driver.Config{
        Type:         driver.DriverType(cfg.DBMode),
        DatabasePath: cfg.DatabasePath,
        TursoURL:     cfg.TursoURL,
        TursoToken:   cfg.TursoToken,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to detect database mode: %w", err)
    }
    
    // AC #5: Log detected mode at startup
    log.Printf("[Database] Mode: %s (%s)", result.Type, result.Reason)
    
    // Validate configuration
    if err := driver.ValidateConfig(&driver.Config{
        Type:         result.Type,
        DatabasePath: cfg.DatabasePath,
        TursoURL:     cfg.TursoURL,
        TursoToken:   cfg.TursoToken,
    }, result.Type); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %w", err)
    }
    
    // Create driver
    drv, err := cli.ConnectWithDriver(
        driver.WithDriverType(result.Type),
        driver.WithDatabasePath(cfg.DatabasePath),
        driver.WithTurso(cfg.TursoURL, cfg.TursoToken),
        driver.WithSyncInterval(cfg.SyncInterval),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create database driver: %w", err)
    }
    
    // Connect
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := drv.Connect(ctx); err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    log.Printf("[Database] Connected successfully in %s mode", result.Type)
    
    return drv.GetDB(), nil
}
```

### Phase 4: Testing (1.5 hours)

**4.1 Unit Tests (internal/db/driver/detector_test.go)**
Following golang-testing patterns with table-driven tests:

```go
func TestDetectMode(t *testing.T) {
    tests := []struct {
        name     string
        config   *Config
        expected DriverType
        reason   string
    }{
        // AC #3: Only db-path → local SQLite
        {
            name: "local_sqlite_default",
            config: &Config{
                DatabasePath: "/tmp/test.db",
            },
            expected: DriverSQLite,
            reason:   "local SQLite",
        },
        
        // AC #1: turso-url only → remote
        {
            name: "remote_turso",
            config: &Config{
                TursoURL: "libsql://test.turso.io",
                TursoToken: "test-token",
            },
            expected: DriverTursoRemote,
            reason:   "remote Turso",
        },
        
        // AC #2: turso-url + local path → embedded replica
        {
            name: "embedded_replica",
            config: &Config{
                DatabasePath: "/tmp/test.db",
                TursoURL:     "libsql://test.turso.io",
                TursoToken:   "test-token",
            },
            expected: DriverTursoReplica,
            reason:   "embedded replica",
        },
        
        // Edge cases
        {
            name: "empty_config",
            config: &Config{},
            expected: DriverSQLite,
            reason:   "default fallback",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := DetectMode(tt.config)
            if result.Type != tt.expected {
                t.Errorf("DetectMode() = %v, want %v", result.Type, tt.expected)
            }
            if !strings.Contains(result.Reason, tt.reason) {
                t.Errorf("Reason = %v, want to contain %v", result.Reason, tt.reason)
            }
        })
    }
}

func TestValidateConfig(t *testing.T) {
    tests := []struct {
        name    string
        config  *Config
        drvType DriverType
        wantErr bool
    }{
        {
            name: "valid_sqlite",
            config: &Config{DatabasePath: "/tmp/test.db"},
            drvType: DriverSQLite,
            wantErr: false,
        },
        {
            name: "invalid_sqlite_missing_path",
            config: &Config{DatabasePath: ""},
            drvType: DriverSQLite,
            wantErr: true,
        },
        {
            name: "valid_remote",
            config: &Config{
                TursoURL: "libsql://test.turso.io",
                TursoToken: "token",
            },
            drvType: DriverTursoRemote,
            wantErr: false,
        },
        {
            name: "invalid_remote_missing_token",
            config: &Config{TursoURL: "libsql://test.turso.io"},
            drvType: DriverTursoRemote,
            wantErr: true,
        },
        {
            name: "valid_replica",
            config: &Config{
                DatabasePath: "/tmp/test.db",
                TursoURL:     "libsql://test.turso.io",
                TursoToken:   "token",
            },
            drvType: DriverTursoReplica,
            wantErr: false,
        },
        {
            name: "invalid_replica_missing_path",
            config: &Config{
                TursoURL: "libsql://test.turso.io",
                TursoToken: "token",
            },
            drvType: DriverTursoReplica,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateConfig(tt.config, tt.drvType)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}

func TestParseExplicitMode(t *testing.T) {
    tests := []struct {
        input    string
        expected DriverType
        wantErr  bool
    }{
        {"", "", false},
        {"auto", "", false},
        {"local", DriverSQLite, false},
        {"remote", DriverTursoRemote, false},
        {"replica", DriverTursoReplica, false},
        {"invalid", "", true},
    }
    // ...
}

func TestDetectOrExplicitMode(t *testing.T) {
    // Test that explicit mode overrides auto-detection
    // ...
}
```

**4.2 Integration Tests**
```go
func TestDetectMode_Integration_AllModes(t *testing.T) {
    // Test that detection works with real config loading
    // ...
}
```

### Phase 5: Documentation (30 min)

**5.1 Update docs/architecture/08-database-drivers.md**
- Document auto-detection logic
- Document configuration precedence
- Add examples for each mode
- Document CLI flags and environment variables

**5.2 Update docs/START_HERE.md**
- Add reference to database mode auto-detection

**5.3 Create docs/DATABASE_MODES.md**
- Comprehensive guide for all database modes
- Configuration examples
- Troubleshooting section

## File Structure After Implementation

```
internal/db/driver/
├── detector.go           # NEW
├── detector_test.go      # NEW
├── config.go             # UNCHANGED
└── ...

cmd/dopa/
├── main.go               # MODIFIED
├── config.go             # NEW
└── ...
```

## Test Coverage Target: 90%+

## Dependencies

- **Requires**: TASK-60.1 (COMPLETED), TASK-60.2 (DONE), TASK-60.3 (DONE)
- **Blocks**: TASK-60.7 (integration)

## Acceptance Criteria Mapping

| AC | Phase | Description |
|----|-------|-------------|
| #1 | 2 | Auto-detect: turso-url present → remote |
| #2 | 2 | Auto-detect: turso-url + local path → embedded replica |
| #3 | 2 | Auto-detect: only db-path → local SQLite |
| #4 | 2 | Validation for required config per mode |
| #5 | 3 | Log detected mode at startup |

## Parallel Opportunities

This task can be done **in parallel** with:
- TASK-60.5 (Migration compatibility) - Independent
- TASK-60.6 (Integration tests) - Independent
- TASK-60.9 (Documentation) - Can start in parallel

## Sequential Dependencies

This task **must be done before**:
- TASK-60.7 (Integration and refactoring) - Needs auto-detection working

## Success Criteria

- [ ] All 5 acceptance criteria met
- [ ] Unit tests pass (90%+ coverage)
- [ ] Integration tests pass with real configuration
- [ ] Documentation updated
- [ ] Code passes golangci-lint
- [ ] Code follows golang-pro, golang-patterns, golang-testing best practices
- [ ] Backward compatibility maintained (SQLite is still default)

## Risk Mitigation

1. **Backward Compatibility**: Default to local SQLite if no Turso config present
2. **Configuration Validation**: Clear error messages for missing required fields
3. **Logging**: Always log detected mode for debugging
4. **Environment Variables**: Support both CLI flags and env vars with proper precedence
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Analysis

### Scope Assessment
This task is **appropriately scoped** for a single PR:
- Focused on auto-detection logic and configuration loading
- All 5 ACs are tightly related to the same feature
- Estimated 3-4 hours of work
- Clear boundaries with other tasks

### Why Not Split
Considered splitting into:
- 60.8a: Detection logic
- 60.8b: Configuration loading
- 60.8c: Logging and integration

**Decision: Keep as single task** because:
1. Detection logic and config loading are tightly coupled
2. Logging is integral to the detection process
3. Splitting would create artificial boundaries
4. Testing would require integration across all parts
5. Single PR reduces coordination overhead

### Dependencies & Execution Order

**Sequential Dependencies:**
- TASK-60.1 (Database abstraction layer) - **COMPLETED**
- TASK-60.2 (Turso remote driver) - **DONE**
- TASK-60.3 (Turso embedded replica driver) - **DONE**

**Parallel Opportunities:**
- Can be implemented **in parallel** with:
  - TASK-60.5 (Migration compatibility) - Independent
  - TASK-60.6 (Integration tests) - Independent
  - TASK-60.9 (Documentation) - Independent

**Blocked Tasks:**
- TASK-60.7 (Integration and refactoring) - Requires auto-detection working

### Key Technical Decisions

1. **Detection Logic**: Pure function with clear rules based on config presence
2. **Configuration Precedence**: CLI flags > env vars > defaults
3. **Explicit Override**: Support --db-mode flag to override auto-detection
4. **Validation**: Validate config AFTER detection, not before
5. **Logging**: Always log detected mode for observability
6. **Backward Compatibility**: Default to SQLite if no Turso config

### Reference Files

- internal/db/driver/config.go - Configuration types
- internal/db/driver/driver.go - Interface definition
- cmd/dopa/main.go - Application entry point

Created detector.go with DetectMode, DetectOrExplicitMode, ParseExplicitMode, and ValidateConfigForMode functions. Added detector_test.go with comprehensive table-driven tests covering all modes and edge cases. Added cmd/dopa/config.go for configuration loading with CLI/env precedence. Modified cmd/dopa/main.go with new flags: --turso-url, --turso-token, --db-mode, --sync-interval. All tests pass and linting passes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented automatic database mode detection for Turso integration.

## Changes

- Added `internal/db/driver/detector.go` with:
  - `DetectMode()`: Auto-detects driver type based on config presence
  - `DetectOrExplicitMode()`: Supports explicit mode override via `--db-mode`
  - `ParseExplicitMode()`: Parses user-specified mode strings
  - `ValidateConfigForMode()`: Validates required config per detected mode

- Added `cmd/dopa/config.go` for configuration loading with CLI/env precedence

- Modified `cmd/dopa/main.go`:
  - Added CLI flags: `--turso-url`, `--turso-token`, `--db-mode`, `--sync-interval`
  - Updated `GetDB()` to use detection logic with logging

## Detection Rules

| Config | Detected Mode |
|--------|---------------|
| `--db` only | local SQLite (default) |
| `--turso-url` + `--turso-token` | remote Turso |
| All three flags | embedded replica |

## Environment Variables

- `TURSO_DATABASE_URL`: Turso database URL
- `TURSO_AUTH_TOKEN`: Turso auth token
- `DOPA_DB_MODE`: Database mode override
- `DOPA_DB_PATH`: Database path

## Tests

- Added `detector_test.go` with table-driven tests covering all modes and edge cases
- All existing tests pass
- Linting passes with 0 issues
<!-- SECTION:FINAL_SUMMARY:END -->
