# Database Modes

This document describes the different database connection modes available in Dopadone and how to configure them.

## Overview

Dopadone supports three database modes to accommodate different workflows:

| Mode | Description | Best For |
|------|-------------|----------|
| **SQLite** | Local SQLite database | Single-device usage, offline work |
| **Turso Remote** | Direct connection to Turso cloud | Always-online, cloud-first workflow |
| **Turso Replica** | Local replica with cloud sync | Offline-capable with cloud backup |

## SQLite Mode (Default)

SQLite mode uses a local database file. This is the default and requires no additional configuration.

### Usage

```bash
# Default: uses ./dopadone.db
dopa tasks list

# Specify custom database path
dopa --db /path/to/my-database.db tasks list
```

### Environment Variables

```bash
export DOPA_DB_PATH=/path/to/database.db
dopa tasks list  # Uses DOPA_DB_PATH
```

### When to Use

- Single device workflow
- No internet dependency required
- Maximum performance (local file access)
- Simple setup with no external services

## Turso Remote Mode

Remote mode connects directly to a Turso cloud database. All operations require an internet connection.

### Prerequisites

1. A Turso account and database: [turso.tech](https://turso.tech)
2. Database URL (e.g., `libsql://your-db.turso.io`)
3. Authentication token

### Usage

```bash
# Using CLI flags
dopa --turso-url "libsql://your-db.turso.io" \
     --turso-auth-token "your-auth-token" \
     --db-mode remote \
     tasks list

# Using environment variables
export TURSO_DATABASE_URL="libsql://your-db.turso.io"
export TURSO_AUTH_TOKEN="your-auth-token"
dopa --db-mode remote tasks list
```

### When to Use

- Always-online environments
- Multiple devices sharing the same data
- No need for offline access
- Cloud-first workflow with automatic backups

### Limitations

- Requires internet connection for all operations
- Network latency affects performance
- Connection failures block operations

## Turso Replica Mode

Replica mode maintains a local SQLite database that automatically syncs with a Turso primary database. This combines local performance with cloud backup.

### Prerequisites

Same as Turso Remote mode.

### Usage

```bash
# Using CLI flags
dopa --db ./local-replica.db \
     --turso-url "libsql://your-db.turso.io" \
     --turso-auth-token "your-auth-token" \
     --db-mode replica \
     --sync-interval 60s \
     tasks list

# Using environment variables
export DOPA_DB_PATH=./local-replica.db
export TURSO_DATABASE_URL="libsql://your-db.turso.io"
export TURSO_AUTH_TOKEN="your-auth-token"
export DOPA_DB_MODE=replica
dopa tasks list
```

### Sync Behavior

| Aspect | Behavior |
|--------|----------|
| **Read operations** | Always from local file (microsecond latency) |
| **Write operations** | Written locally, then synced to remote |
| **Background sync** | Automatic at configured interval (default: 60s) |
| **Offline mode** | Works offline, syncs when connection available |

### Sync Interval Configuration

```bash
# Sync every 30 seconds
dopa --sync-interval 30s tasks list

# Sync every 5 minutes
dopa --sync-interval 5m tasks list
```

### When to Use

- Need offline capability with cloud backup
- Want local performance with remote sync
- Traveling or intermittent connectivity
- Working with large datasets (local reads are fast)

## Configuration Reference

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--db` | Local database path | `./dopadone.db` |
| `--turso-url` | Turso database URL | - |
| `--turso-auth-token` | Turso authentication token | - |
| `--db-mode` | Database mode: `local`, `remote`, `replica`, `auto` | `auto` |
| `--sync-interval` | Sync interval for replica mode | `60s` |
| `--config` | Path to YAML config file | Auto-discovered |

### Environment Variables

| Variable | Equivalent Flag |
|----------|-----------------|
| `DOPA_DB_PATH` | `--db` |
| `TURSO_DATABASE_URL` | `--turso-url` |
| `TURSO_AUTH_TOKEN` | `--turso-auth-token` |
| `DOPA_DB_MODE` | `--db-mode` |

### YAML Configuration File

Dopadone supports configuration via a YAML file for persistent settings. Create a `dopadone.yaml` file:

```yaml
database:
  path: ./dopadone.db          # Local database path
  mode: auto                    # local|remote|replica|auto
  sync_interval: 60s            # Sync interval for replica mode
  turso:
    url: libsql://xxx.turso.io  # Turso database URL
    token: xxx                  # Turso auth token (or use env)
```

#### Config File Discovery Order

Dopadone automatically searches for configuration files in this order:

1. **Explicit path** - `--config /path/to/config.yaml`
2. **Current directory** - `./dopadone.yaml`
3. **XDG config home** - `$XDG_CONFIG_HOME/dopadone/config.yaml`
4. **Default XDG** - `~/.config/dopadone/config.yaml`
5. **Home directory** - `~/.dopadone.yaml`

The first file found is used. If no config file exists, CLI flags and environment variables are used.

#### Example Configurations

**Local SQLite:**
```yaml
database:
  path: ~/projects/my-project/dopadone.db
  mode: local
```

**Turso Remote:**
```yaml
database:
  mode: remote
  turso:
    url: libsql://my-db.turso.io
    token: ${TURSO_AUTH_TOKEN}  # Reference env variable
```

**Turso Replica:**
```yaml
database:
  path: ./local-replica.db
  mode: replica
  sync_interval: 30s
  turso:
    url: libsql://my-db.turso.io
    token: ${TURSO_AUTH_TOKEN}
```

### Configuration Precedence

When multiple configuration sources are present, Dopadone uses this precedence order:

```
CLI flags > Environment variables > Config file > Defaults
```

This means:
- CLI flags always override everything else
- Environment variables override config file values
- Config file values override defaults
- If nothing is specified, sensible defaults are used

#### Precedence Examples

```bash
# Config file has mode: local, but CLI overrides
dopa --db-mode remote tasks list  # Uses "remote"

# Config file has turso.url, but env overrides
export TURSO_DATABASE_URL="libsql://other.turso.io"
dopa tasks list  # Uses env variable, not config file

# Partial config: file provides path, env provides token
# Config file:
#   database:
#     path: ./mydb.db
#     turso:
#       url: libsql://my.turso.io
# Env: TURSO_AUTH_TOKEN=secret
dopa tasks list  # Uses path and URL from file, token from env
```

## Auto-Detection

When `--db-mode` is not specified (or set to `auto`), Dopadone automatically detects the mode:

| Configuration | Detected Mode |
|--------------|---------------|
| Only `--db` specified | SQLite (local) |
| `--turso-url` + `--turso-auth-token` (no `--db`) | Turso Remote |
| `--db` + `--turso-url` + `--turso-auth-token` | Turso Replica |

### Examples

```bash
# Auto-detects SQLite (local mode)
dopa --db ./mydb.db tasks list

# Auto-detects Turso Remote
dopa --turso-url "libsql://db.turso.io" \
     --turso-auth-token "token" \
     tasks list

# Auto-detects Turso Replica
dopa --db ./replica.db \
     --turso-url "libsql://db.turso.io" \
     --turso-auth-token "token" \
     tasks list
```

## Mode Comparison

| Feature | SQLite | Turso Remote | Turso Replica |
|---------|--------|--------------|---------------|
| **Offline access** | Yes | No | Yes |
| **Cloud backup** | No | Yes | Yes |
| **Multi-device sync** | Manual | Automatic | Automatic |
| **Read latency** | Microseconds | Network | Microseconds |
| **Write latency** | Microseconds | Network | Microseconds + sync |
| **Setup complexity** | None | Requires Turso account | Requires Turso account |
| **Internet required** | No | Yes | Only for sync |

## Troubleshooting

### Connection Errors

```bash
# Check if Turso credentials are correct
dopa --turso-url "libsql://your-db.turso.io" \
     --turso-auth-token "your-token" \
     --db-mode remote \
     areas list
```

### Sync Issues (Replica Mode)

If sync is not working:
1. Check internet connectivity
2. Verify Turso credentials are valid
3. Check sync logs in the console output

### Performance Issues

- **SQLite mode is slow**: Check disk I/O, ensure database is on fast storage
- **Remote mode is slow**: Network latency is expected; consider replica mode
- **Replica mode sync lag**: Reduce `--sync-interval` for faster syncs

## Related Documentation

- [Turso Setup Guide](TURSO_SETUP.md) - Step-by-step Turso account and database setup
- [Database Driver Architecture](architecture/08-database-drivers.md) - Technical implementation details
- [Architecture Overview](architecture/01-overview.md) - System architecture
- [TUI Documentation](TUI.md) - Terminal UI features including connection status indicator

## TUI Connection Status

When using the TUI interface (`dopa tui`), a connection status indicator is displayed in the footer:

| Indicator | Meaning |
|-----------|---------|
| `■ local` | SQLite local mode (gray) |
| `● remote` | Connected to Turso remote (green) |
| `◐ remote` | Connecting to Turso (yellow) |
| `○ remote` | Disconnected from Turso (red) |
| `● replica` | Replica synced with Turso (green) |
| `◐ replica` | Replica syncing with Turso (yellow) |
| `○ replica` | Replica offline (red) |

The status updates automatically every 2 seconds to reflect the current connection state.
