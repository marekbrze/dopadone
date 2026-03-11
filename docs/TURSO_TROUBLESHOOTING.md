# Comprehensive Troubleshooting Guide for Turso

This guide covers common issues, error messages, and solutions when using Turso database modes (Remote and Replica) with Dopadone.

## Quick Reference

| Problem | Quick Fix |
|---------|-----------|
| Connection timeout | Check network, verify URL, increase timeout |
| Invalid token | Create new token with `turso db tokens create` |
| Expired token | Create new token, consider non-expiring for production |
| Sync failures | Check connectivity, verify credentials, force sync |
| Migration errors | Verify SQLite compatibility, check goose dialect |
| Database locked | Check connections, increase timeout, reduce concurrency |
| Network unavailable | Use replica mode, wait for reconnection |

## Table of Contents

1. [Connection Issues](#connection-issues)
2. [Authentication Errors](#authentication-errors)
3. [Network & Offline Handling](#network--offline-handling)
4. [Replica Mode Sync Issues](#replica-mode-sync-issues)
5. [Migration Errors](#migration-errors)
6. [Database Lock Errors](#database-lock-errors)
7. [Error Message Reference](#error-message-reference)
8. [Diagnostic Toolkit](#diagnostic-toolkit)

---

## Connection Issues

### Connection Timeout

**Symptoms:**
- `driver turso-remote: connect: connection failed: context deadline exceeded`
- Long delays before errors appear
- Intermittent connection failures

**Causes:**
1. Network latency or unstable connection
2. Turso server experiencing high load
3. Firewall blocking outbound connections
4. DNS resolution issues
5. Default timeout too short for your network

**Solutions:**

1. **Increase connection timeout:**

```bash
# Via environment variable (if supported by your config)
export DOPA_CONNECT_TIMEOUT=30s

# Via YAML config file
# dopadone.yaml
database:
  connect_timeout: 30s  # Increase from default 10s
  turso:
    url: libsql://your-db.turso.io
    token: ${TURSO_AUTH_TOKEN}
```

2. **Check network connectivity:**

```bash
# Test basic connectivity to Turso
curl -v https://your-db.turso.io

# Test DNS resolution
nslookup your-db.turso.io
dig your-db.turso.io

# Test with specific port
nc -zv your-db.turso.io 443
```

3. **Verify firewall rules:**

```bash
# Ensure outbound HTTPS (443) is allowed
# On macOS/Linux, check if connection works:
curl -I https://turso.io

# Corporate proxy issues - set proxy if needed
export HTTPS_PROXY=http://proxy.example.com:8080
```

4. **Enable retry logic:**

```yaml
# dopadone.yaml
database:
  max_retries: 3        # Retry 3 times (default)
  retry_interval: 1s    # Wait 1s between retries (default)
  connect_timeout: 15s  # Per-connection timeout
```

### DNS Resolution Failures

**Symptoms:**
- `no such host` or `DNS lookup failed`
- Intermittent failures

**Solutions:**

```bash
# Check DNS is resolving
ping your-db.turso.io

# Use IP address (not recommended - can change)
# Instead, check DNS servers
cat /etc/resolv.conf

# Try alternative DNS
# macOS: System Preferences > Network > DNS
# Linux: Edit /etc/resolv.conf
# Add: nameserver 8.8.8.8
```

### TLS/SSL Errors

**Symptoms:**
- `certificate verify failed`
- `x509: certificate signed by unknown authority`
- `TLS handshake failed`

**Solutions:**

1. **Update CA certificates:**

```bash
# macOS (via Homebrew)
brew install ca-certificates
brew link --force ca-certificates

# Ubuntu/Debian
sudo apt-get update && sudo apt-get install -y ca-certificates
sudo update-ca-certificates

# CentOS/RHEL
sudo yum update ca-certificates
```

2. **Check system time:**

```bash
# TLS certificates require accurate time
date

# If incorrect, sync time
# macOS
sudo sntp -sS time.apple.com

# Linux
sudo ntpdate -s time.nist.gov
```

### Connection Refused

**Symptoms:**
- `connection refused`
- `network is unreachable`

**Solutions:**

1. **Verify database URL is correct:**

```bash
# Get correct URL from Turso
turso db show your-db-name --url

# URL format: libsql://dbname-organization.turso.io
```

2. **Check Turso service status:**

```bash
# Check Turso status page
curl -s https://status.turso.tech | grep -i "operational"

# Or visit: https://status.turso.tech
```

3. **Verify database exists:**

```bash
turso db list | grep your-db-name
```

---

## Authentication Errors

### Invalid Token

**Symptoms:**
- `authentication failed`
- `invalid token`
- `unauthorized`
- `driver turso-remote: connect: connection failed: ...`

**Causes:**
1. Token was copied incorrectly
2. Token was created for wrong database
3. Token was invalidated
4. Token format is corrupted

**Solutions:**

1. **Verify token format:**

```bash
# Token should be a long JWT-like string
# Example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

echo $TURSO_AUTH_TOKEN | wc -c
# Should be 100+ characters
```

2. **Create a new token:**

```bash
# Create fresh token
turso db tokens create your-db-name

# Copy the entire token, including any prefixes
```

3. **Verify token is for correct database:**

```bash
# List your databases
turso db list

# Create token for specific database
turso db tokens create correct-db-name
```

4. **Check if token was invalidated:**

```bash
# If tokens were invalidated, all previous tokens stop working
turso db tokens invalidate your-db-name  # This invalidates ALL tokens

# Create new token after invalidation
turso db tokens create your-db-name
```

5. **Store token correctly:**

```bash
# In shell profile (~/.bashrc, ~/.zshrc)
export TURSO_AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Verify it's set correctly
echo "$TURSO_AUTH_TOKEN" | head -c 20
```

### Expired Token

**Symptoms:**
- `token expired`
- `authentication failed` (after working previously)
- Sudden failures after days/weeks of use

**Causes:**
1. Token was created with expiration
2. System clock drift

**Solutions:**

1. **Create non-expiring token (for production):**

```bash
# Create token that never expires
turso db tokens create your-db-name --expiration never

# WARNING: Use with caution. Rotate periodically for security.
```

2. **Create long-lived token:**

```bash
# 30-day token
turso db tokens create your-db-name --expiration 30d

# 90-day token
turso db tokens create your-db-name --expiration 90d
```

3. **Set up token rotation:**

```bash
# Create new token
NEW_TOKEN=$(turso db tokens create your-db-name --expiration 30d)

# Update environment
export TURSO_AUTH_TOKEN="$NEW_TOKEN"

# Invalidate old tokens periodically
turso db tokens invalidate your-db-name
turso db tokens create your-db-name --expiration 30d
```

### Token Permission Issues

**Symptoms:**
- `permission denied` for writes
- Can read but not write

**Causes:**
1. Using read-only token for write operations

**Solutions:**

```bash
# Check if token is read-only
# Read-only tokens are created with:
turso db tokens create your-db-name --read-only

# Create full-access token for Dopadone
turso db tokens create your-db-name  # No --read-only flag
```

### Token Rotation Best Practices

1. **For Development:**
   - Use 7-day or 30-day tokens
   - Store in environment variables
   - Rotate when expired

2. **For Production:**
   - Use non-expiring or long-lived tokens (90+ days)
   - Store in secure secret manager
   - Rotate quarterly
   - Have a rotation procedure documented

3. **For CI/CD:**
   - Use repository secrets (GitHub, GitLab, etc.)
   - Use short-lived tokens when possible
   - Never log tokens

---

## Network & Offline Handling

### Network Unavailable

**Symptoms:**
- `network is unreachable`
- `no route to host`
- Complete failure to connect

**Detection:**

```bash
# Check if network is available
ping -c 3 8.8.8.8

# Check DNS
nslookup turso.io

# Check HTTPS
curl -I https://turso.io
```

**Solutions:**

1. **Use Replica Mode for offline support:**

```yaml
# dopadone.yaml
database:
  mode: replica
  path: ./dopadone-replica.db
  sync_interval: 60s
  turso:
    url: libsql://your-db.turso.io
    token: ${TURSO_AUTH_TOKEN}
```

With replica mode:
- Local file is always available
- Syncs when network is restored
- Graceful offline handling

2. **Graceful degradation pattern:**

```bash
# Try remote first, fallback to local
if dopa --db-mode remote areas list 2>/dev/null; then
    echo "Using remote connection"
else
    echo "Remote unavailable, using local"
    dopa --db-mode local areas list
fi
```

### Intermittent Connectivity

**Symptoms:**
- Works sometimes, fails other times
- Flaky connections

**Solutions:**

1. **Increase retries:**

```yaml
database:
  max_retries: 5        # More retries
  retry_interval: 2s    # Longer delay between retries
```

2. **Use replica mode with longer sync interval:**

```yaml
database:
  mode: replica
  sync_interval: 120s   # Sync every 2 minutes
```

### Offline Detection

**Check connection status programmatically:**

```bash
# Test connection quickly
dopa --db-mode remote areas list >/dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Online"
else
    echo "Offline or connection issue"
fi
```

**TUI connection indicator:**

When using `dopa tui`, the footer shows:
- `● remote` - Connected (green)
- `○ remote` - Disconnected (red)
- `◐ remote` - Connecting (yellow)

### Reconnection Strategies

1. **Automatic retry (built-in):**
   - Dopadone automatically retries connections
   - Default: 3 retries with 1s interval

2. **Manual reconnection:**

```bash
# If connection is stuck, restart the command
dopa areas list

# For replica mode, force sync
dopa --db-mode replica migrate status
```

3. **Connection health check:**

```bash
# Simple health check script
#!/bin/bash
MAX_ATTEMPTS=5
for i in $(seq 1 $MAX_ATTEMPTS); do
    if dopa --db-mode remote areas list >/dev/null 2>&1; then
        echo "Connection OK"
        exit 0
    fi
    echo "Attempt $i/$MAX_ATTEMPTS failed, retrying..."
    sleep 2
done
echo "Connection failed after $MAX_ATTEMPTS attempts"
exit 1
```

---

## Replica Mode Sync Issues

### Sync Timeout Failures

**Symptoms:**
- `[TursoReplica] Sync failed: context deadline exceeded`
- Sync takes too long
- Initial sync never completes

**Causes:**
1. Large database size
2. Slow network connection
3. High latency to Turso servers
4. Conflict during sync

**Solutions:**

1. **Check sync status:**

```bash
# View recent sync logs
dopa --db-mode replica areas list 2>&1 | grep -i sync

# Check sync info (if available in logs)
# Look for: [TursoReplica] Synced X frames (frame_no: Y)
```

2. **Manual sync trigger:**

```bash
# Restart to trigger fresh sync
dopa --db-mode replica areas list
```

3. **Increase sync interval for large databases:**

```yaml
database:
  mode: replica
  sync_interval: 300s   # 5 minutes for large DBs
```

4. **Check database size:**

```bash
# Check local replica file size
ls -lh ./dopadone-replica.db

# Large files (>100MB) may need longer sync times
```

### Sync Conflict Resolution

**Symptoms:**
- Data not syncing correctly
- Missing data after sync
- Duplicate or stale data

**Causes:**
1. Concurrent writes from multiple clients
2. Network interruption during sync
3. Schema mismatch

**Solutions:**

1. **Full re-sync:**

```bash
# Backup local data
cp ./dopadone-replica.db ./dopadone-replica.db.backup

# Delete local replica and re-download
rm ./dopadone-replica.db

# Restart to trigger fresh sync
dopa --db-mode replica areas list
```

2. **Verify schema consistency:**

```bash
# Check migration status
dopa migrate status

# Verify schema
dopa migrate verify
```

3. **Single writer pattern:**

When using replica mode:
- Designate one instance as primary writer
- Other instances should primarily read
- Sync conflicts are less likely with single writer

### Partial Sync States

**Symptoms:**
- Some data synced, some missing
- Inconsistent state between local and remote

**Solutions:**

```bash
# Check last sync time in logs
# Look for: [TursoReplica] Synced X frames (frame_no: Y)

# Force complete re-sync
rm ./dopadone-replica.db
dopa --db-mode replica areas list
```

### Sync Status Values

| Status | Meaning | Action |
|--------|---------|--------|
| `idle` | Sync complete, waiting for next interval | Normal operation |
| `syncing` | Currently syncing | Wait for completion |
| `error` | Last sync failed | Check logs, retry |
| `offline` | No network or not connected | Check connectivity |

---

## Migration Errors

### libSQL-Specific Migration Failures

**Symptoms:**
- `goose: migration failed`
- SQL syntax errors during migration
- Schema mismatch errors

**Common Issues:**

1. **SQLite feature not supported:**

Some SQLite features may behave differently in libSQL:

```sql
-- Problem: Some PRAGMA settings differ
PRAGMA journal_mode = WAL;  -- May behave differently

-- Solution: Test migrations on Turso before production
```

2. **Dialect issues:**

Dopadone uses `sqlite3` dialect with goose, which is compatible with libSQL.

```bash
# Verify migration works locally first
dopa migrate up

# Then apply to remote
dopa --db-mode remote migrate up
```

### Schema Incompatibility

**Symptoms:**
- `table already exists` but data missing
- `no such table` after migration
- Column type mismatches

**Solutions:**

1. **Check current schema version:**

```bash
dopa migrate status

# Output shows current version
# goose_db_version table tracks applied migrations
```

2. **Verify migration files:**

```bash
# List migration files
ls -la internal/migrate/migrations/

# Check migration content
cat internal/migrate/migrations/20240301000000_initial_schema.sql
```

3. **Rollback and re-apply:**

```bash
# WARNING: This may cause data loss
dopa migrate down    # Rollback last migration
dopa migrate up      # Re-apply
```

### Goose Dialect Issues

**Symptoms:**
- `unknown dialect`
- Migration parsing errors

**Solution:**

Dopadone uses `sqlite3` dialect which is correct for libSQL/Turso. No changes needed.

```bash
# If you see dialect errors, verify goose version
go version -m $(which goose) | grep goose
```

### Migration Rollback Procedures

1. **Safe rollback:**

```bash
# Check current version
dopa migrate status

# Rollback ONE migration
dopa migrate down

# Verify state
dopa migrate verify
```

2. **Full reset (destructive):**

```bash
# WARNING: Destroys all data
dopa migrate reset
```

3. **For replica mode:**

```bash
# Reset local replica
rm ./dopadone-replica.db

# Re-sync from remote
dopa --db-mode replica migrate status
```

---

## Database Lock Errors

### SQLITE_BUSY Errors

**Symptoms:**
- `database is locked`
- `SQLITE_BUSY: database is locked`
- `driver turso-replica: connect: database is locked`

**Causes:**
1. Another process has the database open
2. Long-running transaction
3. Too many concurrent connections
4. Network latency in remote mode

**Solutions:**

1. **Check for other processes:**

```bash
# Find processes using the database file
lsof ./dopadone-replica.db

# Or
fuser ./dopadone-replica.db
```

2. **Wait and retry:**

```bash
# The error may be transient
sleep 2
dopa areas list
```

3. **For replica mode - ensure single instance:**

```bash
# Only one Dopadone instance should use replica file
# Kill other instances
pkill -f dopa

# Try again
dopa areas list
```

### Concurrent Access Issues

**Symptoms:**
- Lock errors under load
- Timeouts with multiple clients

**Solutions:**

1. **Use connection pooling (built-in):**

Dopadone uses Go's `database/sql` package which handles connection pooling automatically.

2. **Reduce concurrent operations:**

```bash
# Instead of parallel operations, run sequentially
dopa tasks list
dopa projects list
```

3. **For replica mode - avoid concurrent writes:**

```bash
# Designate single writer
# Other replicas should sync passively
```

### Lock Timeout Solutions

1. **Increase timeout (local SQLite):**

```bash
# Set busy timeout in SQLite
# Dopadone handles this internally
```

2. **For remote mode - network timeout:**

```yaml
database:
  connect_timeout: 30s  # Increase connection timeout
```

### Connection Pool Tuning

Dopadone uses Go's default connection pool settings. For high-concurrency scenarios:

```yaml
# Connection pool is managed automatically
# No manual configuration needed for most use cases
```

If you need custom settings, modify the source code in `internal/db/driver/`.

---

## Error Message Reference

### Quick Reference Table

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `driver not registered` | Unknown database mode | Use: local, remote, or replica |
| `connection failed: context deadline exceeded` | Timeout | Increase timeout, check network |
| `connection failed: invalid token` | Bad credentials | Create new token |
| `turso url required` | Missing URL | Set TURSO_DATABASE_URL |
| `turso token required` | Missing token | Set TURSO_AUTH_TOKEN |
| `database path required` | Missing path | Set DOPA_DB_PATH |
| `driver already closed` | Use after close | Reconnect |
| `driver X: validate: invalid configuration` | Config error | Check config file/env vars |
| `authentication failed` | Invalid/expired token | Create new token |
| `network is unreachable` | No network | Check connectivity |
| `no such host` | DNS failure | Check DNS settings |
| `certificate verify failed` | TLS/CA issue | Update CA certificates |
| `database is locked` | Concurrent access | Close other connections |
| `SQLITE_BUSY` | Lock contention | Wait, reduce concurrency |
| `migration failed` | SQL error | Check migration syntax |
| `schema drift detected` | Schema mismatch | Run missing migrations |

### Dopadone-Specific Errors

#### `driver not registered`

**Full Error:**
```
driver : validate: driver not registered
```

**Cause:** Invalid database mode specified

**Solution:**
```bash
# Valid modes
dopa --db-mode local ...
dopa --db-mode remote ...
dopa --db-mode replica ...
dopa --db-mode auto ...  # Auto-detect
```

#### `invalid configuration`

**Full Error:**
```
driver turso-remote: create: invalid configuration: turso url required
driver turso-remote: create: invalid configuration: turso token required
driver turso-replica: create: invalid configuration: database path required
```

**Cause:** Missing required configuration

**Solution:**
```bash
# Set required environment variables
export TURSO_DATABASE_URL="libsql://your-db.turso.io"
export TURSO_AUTH_TOKEN="your-token"
export DOPA_DB_PATH="./replica.db"  # For replica mode
```

#### `connection failed`

**Full Error:**
```
driver turso-remote: connect: connection failed: <underlying error>
```

**Cause:** Network or authentication issue

**Solution:** See [Connection Issues](#connection-issues)

#### `driver already closed`

**Full Error:**
```
driver turso-remote: ping: driver already closed
driver turso-replica: sync: driver already closed
```

**Cause:** Operation attempted after Close() was called

**Solution:** Restart the application

### libSQL/Turso-Specific Errors

#### `AUTH_INVALID_TOKEN`

**Cause:** Token is malformed or invalid

**Solution:** Create new token with `turso db tokens create`

#### `AUTH_EXPIRED_TOKEN`

**Cause:** Token has expired

**Solution:**
```bash
# Create new token
turso db tokens create your-db --expiration 30d

# Or non-expiring
turso db tokens create your-db --expiration never
```

#### `DATABASE_NOT_FOUND`

**Cause:** Database doesn't exist or wrong URL

**Solution:**
```bash
# List databases
turso db list

# Get correct URL
turso db show your-db --url
```

#### `RATE_LIMIT_EXCEEDED`

**Cause:** Too many requests

**Solution:** Wait and retry. Consider replica mode to reduce remote calls.

### SQLite Errors (All Modes)

#### `SQLITE_BUSY` (5)

**Cause:** Database locked by another connection

**Solution:** Close other connections, wait, use replica mode

#### `SQLITE_LOCKED` (6)

**Cause:** Table locked

**Solution:** Wait and retry

#### `SQLITE_NOMEM` (7)

**Cause:** Out of memory

**Solution:** Free memory, reduce data size

#### `SQLITE_READONLY` (8)

**Cause:** Database is read-only

**Solution:** Check file permissions, use correct token (not read-only)

#### `SQLITE_IOERR` (10)

**Cause:** Disk I/O error

**Solution:** Check disk space, file permissions

#### `SQLITE_CORRUPT` (11)

**Cause:** Database file corrupted

**Solution:** Restore from backup or re-sync replica

#### `SQLITE_FULL` (13)

**Cause:** Disk full

**Solution:** Free disk space

#### `SQLITE_CANTOPEN` (14)

**Cause:** Cannot open database file

**Solution:** Check path, permissions

#### `SQLITE_PROTOCOL` (15)

**Cause:** Locking protocol error

**Solution:** Close all connections, restart

---

## Diagnostic Toolkit

### Connection Diagnostics

#### Test Basic Connectivity

```bash
#!/bin/bash
# test-connection.sh

echo "=== Network Test ==="
ping -c 3 8.8.8.8

echo "=== DNS Test ==="
nslookup turso.io

echo "=== HTTPS Test ==="
curl -I https://turso.io

echo "=== Turso API Test ==="
curl -s https://api.turso.tech/v1/health | head -5
```

#### Test Turso Connection

```bash
#!/bin/bash
# test-turso.sh

DB_URL="${TURSO_DATABASE_URL:-libsql://your-db.turso.io}"
TOKEN="${TURSO_AUTH_TOKEN}"

echo "Testing connection to: $DB_URL"

# Test with turso CLI
if command -v turso &> /dev/null; then
    turso db list
else
    echo "Turso CLI not installed"
fi

# Test with dopadone
dopa --turso-url "$DB_URL" \
     --turso-auth-token "$TOKEN" \
     --db-mode remote \
     areas list
```

### Token Validation

```bash
#!/bin/bash
# validate-token.sh

TOKEN="${TURSO_AUTH_TOKEN}"

if [ -z "$TOKEN" ]; then
    echo "ERROR: TURSO_AUTH_TOKEN not set"
    exit 1
fi

echo "Token length: ${#TOKEN} characters"

if [ ${#TOKEN} -lt 50 ]; then
    echo "WARNING: Token seems too short"
fi

# Test token with a simple query
dopa --db-mode remote areas list >/dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "Token is valid"
else
    echo "Token is INVALID or expired"
fi
```

### Sync Status Check

```bash
#!/bin/bash
# check-sync.sh

DB_PATH="${DOPA_DB_PATH:-./dopadone-replica.db}"

echo "=== Local Replica Status ==="
if [ -f "$DB_PATH" ]; then
    ls -lh "$DB_PATH"
    stat "$DB_PATH" | grep -i modify
else
    echo "Replica file not found: $DB_PATH"
fi

echo ""
echo "=== Testing Sync ==="
dopa --db-mode replica areas list 2>&1 | grep -i sync

echo ""
echo "=== Last Modified ==="
dopa --db-mode replica migrate status
```

### Log Analysis

```bash
#!/bin/bash
# analyze-logs.sh

# Run command and capture logs
dopa --db-mode replica areas list 2>&1 | tee /tmp/dopa-debug.log

echo ""
echo "=== Error Analysis ==="
grep -i error /tmp/dopa-debug.log
grep -i failed /tmp/dopa-debug.log
grep -i warning /tmp/dopa-debug.log

echo ""
echo "=== Sync Analysis ==="
grep -i sync /tmp/dopa-debug.log
```

### Performance Diagnostics

```bash
#!/bin/bash
# perf-test.sh

echo "=== Local SQLite Performance ==="
time dopa --db-mode local areas list > /dev/null

echo ""
echo "=== Remote Mode Performance ==="
time dopa --db-mode remote areas list > /dev/null

echo ""
echo "=== Replica Mode Performance ==="
time dopa --db-mode replica areas list > /dev/null
```

### Database Health Check

```bash
#!/bin/bash
# db-health.sh

echo "=== Database Health Check ==="

# Check migration status
echo "Migration Status:"
dopa migrate status

echo ""
echo "Schema Verification:"
dopa migrate verify

echo ""
echo "Database Integrity (SQLite):"
sqlite3 ./dopadone.db "PRAGMA integrity_check;"

echo ""
echo "Database Stats:"
sqlite3 ./dopadone.db "SELECT name FROM sqlite_master WHERE type='table';"
```

### Full Diagnostic Report

```bash
#!/bin/bash
# full-diagnostic.sh

echo "========================================="
echo "  DOPADONE DIAGNOSTIC REPORT"
echo "  $(date)"
echo "========================================="

echo ""
echo "=== System Info ==="
uname -a
date

echo ""
echo "=== Environment Variables ==="
env | grep -E "(DOPA|TURSO)" | sed 's/TOKEN=.*/TOKEN=***HIDDEN***/g'

echo ""
echo "=== Config File ==="
if [ -f "./dopadone.yaml" ]; then
    cat ./dopadone.yaml | sed 's/token:.*/token: ***HIDDEN***/g'
else
    echo "No config file found"
fi

echo ""
echo "=== Database Files ==="
ls -la *.db 2>/dev/null || echo "No database files in current directory"

echo ""
echo "=== Network Connectivity ==="
ping -c 1 8.8.8.8 && echo "Internet: OK" || echo "Internet: FAILED"
curl -s -o /dev/null -w "%{http_code}" https://turso.io && echo " Turso API: OK" || echo " Turso API: FAILED"

echo ""
echo "=== Turso CLI Status ==="
turso auth status 2>/dev/null || echo "Not logged in or CLI not installed"

echo ""
echo "=== Dopadone Connection Test ==="
dopa areas list >/dev/null 2>&1 && echo "Connection: OK" || echo "Connection: FAILED"

echo ""
echo "=== Migration Status ==="
dopa migrate status 2>/dev/null || echo "Migration check failed"

echo ""
echo "========================================="
echo "  END OF DIAGNOSTIC REPORT"
echo "========================================="
```

### Quick Diagnostic Commands

| Purpose | Command |
|---------|---------|
| Test connection | `dopa areas list` |
| Check migrations | `dopa migrate status` |
| Verify schema | `dopa migrate verify` |
| Check Turso auth | `turso auth status` |
| List databases | `turso db list` |
| Show database URL | `turso db show <name> --url` |
| Create token | `turso db tokens create <name>` |
| Test network | `curl -I https://turso.io` |
| Check DNS | `nslookup turso.io` |
| View replica status | `ls -lh ./dopadone-replica.db` |

---

## Related Documentation

- [Database Modes](DATABASE_MODES.md) - Mode explanations and configuration
- [Turso Setup Guide](TURSO_SETUP.md) - Account and credential setup
- [Turso Migrations](TURSO_MIGRATIONS.md) - Schema migration procedures
- [Turso Data Migration](TURSO_DATA_MIGRATION.md) - Step-by-step guide for migrating data from SQLite to Turso
- [Architecture Overview](architecture/01-overview.md) - System architecture
- [Turso Official Docs](https://docs.turso.tech) - Turso documentation

## Getting Help

1. **Check this guide** - Most common issues are documented here
2. **Run diagnostics** - Use the [Diagnostic Toolkit](#diagnostic-toolkit)
3. **Check Turso status** - [status.turso.tech](https://status.turso.tech)
4. **Turso Discord** - [tur.so/discord](https://tur.so/discord)
5. **GitHub Issues** - Report bugs in the Dopadone repository
