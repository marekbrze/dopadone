# SQLite to Turso Data Migration Guide

This guide walks you through migrating your existing Dopadone SQLite database to Turso cloud. It covers data export, import, verification, and rollback procedures.

## Overview

This guide focuses on **data migration** - moving your existing SQLite data to Turso. For schema migrations (database structure changes), see the [Turso Migrations Guide](TURSO_MIGRATIONS.md).

### When to Use This Guide

- You have an existing Dopadone SQLite database with data
- You want to move to Turso for cloud backup or multi-device sync
- You're switching from local SQLite to Turso Remote or Replica mode

### Migration Paths

| From | To | Recommended Method |
|------|----|--------------------|
| Local SQLite | Turso Remote | Export + Create from dump |
| Local SQLite | Turso Replica | Export + Create from dump |

## Prerequisites

Before starting, ensure you have:

- [ ] Turso CLI installed: `turso --version`
- [ ] Turso account: `turso auth status`
- [ ] Existing Dopadone SQLite database
- [ ] Sufficient disk space for backup (at least 2x database size)
- [ ] Network connectivity to Turso

### Turso Account Setup

If you haven't set up Turso yet:

```bash
# Install Turso CLI (macOS)
brew install tursodatabase/tap/turso

# Sign up for account
turso auth signup

# Verify login
turso auth status
```

See the [Turso Setup Guide](TURSO_SETUP.md) for detailed account setup instructions.

---

## Migration Checklist

Use this checklist to track your migration progress:

### Pre-Migration
- [ ] Backup existing SQLite database
- [ ] Verify SQLite database integrity
- [ ] Check database size (≤2GB for Turso)
- [ ] Ensure all migrations applied
- [ ] Document current record counts

### Turso Setup
- [ ] Create Turso database (or have credentials ready)
- [ ] Generate Turso auth token
- [ ] Verify Turso connectivity

### Migration
- [ ] Export data from SQLite
- [ ] Create Turso database from dump OR import to existing database
- [ ] Verify data integrity

### Post-Migration
- [ ] Update Dopadone configuration
- [ ] Test application connectivity
- [ ] Verify all operations work
- [ ] Archive old database (don't delete yet)

### Rollback Plan (if needed)
- [ ] Keep backup accessible
- [ ] Document rollback steps
- [ ] Test rollback procedure

---

## Step 1: Prepare and Backup

### 1.1 Verify Database Integrity

Before migrating, verify your SQLite database is healthy:

```bash
# Check database integrity
sqlite3 dopadone.db "PRAGMA integrity_check;"
# Expected output: ok

# Check for corruption
sqlite3 dopadone.db "PRAGMA quick_check;"
# Expected output: ok
```

If you see any errors, run:

```bash
# Attempt to recover
sqlite3 dopadone.db ".recover" > recovered.sql
sqlite3 dopadone-fixed.db < recovered.sql
# Use dopadone-fixed.db for migration
```

### 1.2 Check Database Size

Turso has a 2GB database size limit. Check your database size:

```bash
# Check file size
ls -lh dopadone.db

# If near 2GB, consider:
# 1. Cleaning up old data
# 2. Archiving historical records
# 3. Using Turso in replica mode (syncs incrementally)
```

### 1.3 Verify Migrations Applied

Ensure all schema migrations are applied:

```bash
# Check migration status
dopa migrate status

# Expected: All migrations show as "applied"
```

If pending migrations exist:

```bash
# Apply pending migrations
dopa migrate up

# Verify
dopa migrate verify
```

### 1.4 Document Current State

Record current record counts for verification later:

```bash
# Count records in each table
echo "=== Record Counts Before Migration ===" > migration-counts.txt
sqlite3 dopadone.db "SELECT 'areas: ' || COUNT(*) FROM areas;" >> migration-counts.txt
sqlite3 dopadone.db "SELECT 'subareas: ' || COUNT(*) FROM subareas;" >> migration-counts.txt
sqlite3 dopadone.db "SELECT 'projects: ' || COUNT(*) FROM projects;" >> migration-counts.txt
sqlite3 dopadone.db "SELECT 'tasks: ' || COUNT(*) FROM tasks;" >> migration-counts.txt
sqlite3 dopadone.db "SELECT 'goose_db_version: ' || COUNT(*) FROM goose_db_version;" >> migration-counts.txt

# View counts
cat migration-counts.txt
```

### 1.5 Create Backup

**Always backup before migration:**

```bash
# Method 1: Simple file copy (recommended for most cases)
cp dopadone.db dopadone-backup-$(date +%Y%m%d-%H%M%S).db

# Method 2: SQLite backup (safer for databases in use)
sqlite3 dopadone.db ".backup dopadone-backup-$(date +%Y%m%d-%H%M%S).db"

# Verify backup
ls -lh dopadone-backup-*.db
```

---

## Step 2: Export Data from SQLite

Choose the export method that best fits your situation:

### Method 1: SQL Dump (Recommended)

The `.dump` command creates a complete SQL script of your database:

```bash
# Create SQL dump
sqlite3 dopadone.db .dump > dopadone-export.sql

# Verify dump was created
ls -lh dopadone-export.sql

# Quick check dump content
head -50 dopadone-export.sql
tail -20 dopadone-export.sql
```

**Advantages:**
- Human-readable SQL
- Works with Turso's `--from-dump` option
- Can inspect and modify before import

**Best for:** Most migration scenarios

### Method 2: Binary Database File

Use the database file directly:

```bash
# No export needed - use dopadone.db directly
# Turso can create from file:
# turso db create dopadone --from-file dopadone.db
```

**Advantages:**
- Faster (no export step)
- Preserves exact database state

**Best for:** Quick migration, databases <500MB

### Method 3: Online Backup API

For production databases that must stay online:

```bash
# Create a consistent backup while database is in use
sqlite3 dopadone.db ".backup 'dopadone-online-backup.db'"

# Then dump the backup
sqlite3 dopadone-online-backup.db .dump > dopadone-export.sql
```

**Best for:** Zero-downtime migration

### Export Verification

After exporting, verify the dump:

```bash
# Check dump size is reasonable
ls -lh dopadone-export.sql

# Should be similar to or larger than database
ls -lh dopadone.db

# Verify dump contains expected tables
grep "CREATE TABLE" dopadone-export.sql

# Verify dump contains goose version table (for migration tracking)
grep "goose_db_version" dopadone-export.sql
```

---

## Step 3: Create Turso Database

### Option A: Create New Database from Dump

Create a fresh Turso database with your data:

```bash
# Create database from SQL dump
turso db create dopadone --from-dump dopadone-export.sql

# Or create from database file
turso db create dopadone --from-file dopadone.db
```

**When to use:** New Turso deployment, no existing database

### Option B: Create Empty Database First

Create the database, then import:

```bash
# Create empty database
turso db create dopadone

# Get database URL
turso db show dopadone --url
# Example: libsql://dopadone-myorg.turso.io

# Create auth token
turso db tokens create dopadone --expiration never
```

**When to use:** Need to verify database creation before import

### Get Turso Credentials

After creating the database, get your credentials:

```bash
# Get database URL
turso db show dopadone --url
# Output: libsql://dopadone-organization.turso.io

# Create auth token (use non-expiring for migration)
turso db tokens create dopadone --expiration never
# Output: long JWT token string

# Save these for configuration
export TURSO_DATABASE_URL="libsql://dopadone-organization.turso.io"
export TURSO_AUTH_TOKEN="your-token-here"
```

---

## Step 4: Import Data to Turso

If you created the database with `--from-dump` or `--from-file`, data is already imported. Skip to verification.

### Import to Existing Database

If you have an existing Turso database:

```bash
# Method 1: Using Turso shell (for small imports)
turso db shell dopadone < dopadone-export.sql

# Method 2: Using Turso shell interactively
turso db shell dopadone
# At prompt:
.read dopadone-export.sql
.quit
```

### Import Considerations

| Scenario | Solution |
|----------|----------|
| Large import (>100MB) | Use `--from-dump` when creating database |
| Import to existing database | Use `turso db shell` |
| Import fails partway | Check error, fix dump, re-create database |
| Schema conflict | Ensure migrations match before export |

---

## Step 5: Verify Data Integrity

Comprehensive verification ensures your migration was successful.

### 5.1 Schema Verification

Verify the schema was migrated correctly:

```bash
# Check migration status in Turso
dopa --turso-url "$TURSO_DATABASE_URL" \
     --turso-auth-token "$TURSO_AUTH_TOKEN" \
     --db-mode remote \
     migrate status

# Should show all migrations as applied

# Verify schema
dopa --db-mode remote migrate verify
```

### 5.2 Record Count Verification

Compare record counts between source and destination:

```bash
# Source database counts (SQLite)
echo "=== Source (SQLite) ===" 
sqlite3 dopadone.db "SELECT 'areas: ' || COUNT(*) FROM areas;"
sqlite3 dopadone.db "SELECT 'subareas: ' || COUNT(*) FROM subareas;"
sqlite3 dopadone.db "SELECT 'projects: ' || COUNT(*) FROM projects;"
sqlite3 dopadone.db "SELECT 'tasks: ' || COUNT(*) FROM tasks;"

# Destination database counts (Turso)
echo "=== Destination (Turso) ==="
dopa --db-mode remote areas list | wc -l
dopa --db-mode remote projects list | wc -l
dopa --db-mode remote tasks list | wc -l
```

**Important:** Counts should match exactly. If they don't:
1. Check for import errors in Turso shell
2. Verify dump file completeness
3. Re-run import if necessary

### 5.3 Sample Data Verification

Verify actual data content:

```bash
# Check a few records in source
sqlite3 dopadone.db "SELECT id, title FROM tasks LIMIT 3;"

# Check same records in destination
turso db shell dopadone "SELECT id, title FROM tasks LIMIT 3;"
```

### 5.4 Functional Verification

Test that Dopadone works with the Turso database:

```bash
# Test read operations
dopa --db-mode remote areas list
dopa --db-mode remote projects list
dopa --db-mode remote tasks list

# Test write operation (create and delete a test task)
TEST_PROJECT_ID=$(dopa --db-mode remote projects list --json | jq -r '.[0].id')
TEST_TASK_ID=$(dopa --db-mode remote tasks create "Migration test" --project-id "$TEST_PROJECT_ID" --json | jq -r '.id')
dopa --db-mode remote tasks delete "$TEST_TASK_ID"

echo "Functional verification complete"
```

### 5.5 Checksum Verification (Optional)

For critical data, perform checksum comparison:

```bash
# Export from Turso
turso db shell dopadone .dump > turso-export.sql

# Compare dumps (expect differences in order, not content)
# Sort both files and compare
sort dopadone-export.sql > source-sorted.sql
sort turso-export.sql > dest-sorted.sql
diff source-sorted.sql dest-sorted.sql

# Or compare specific tables
sqlite3 dopadone.db "SELECT md5(group_concat(id || title)) FROM tasks ORDER BY id;"
turso db shell dopadone "SELECT md5(group_concat(id || title)) FROM tasks ORDER BY id;"
```

---

## Step 6: Update Configuration

### 6.1 Configure Dopadone for Turso

Update your Dopadone configuration to use Turso:

**Option 1: Environment Variables**

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export TURSO_DATABASE_URL="libsql://dopadone-yourorg.turso.io"
export TURSO_AUTH_TOKEN="your-token-here"
export DOPA_DB_MODE=remote

# Or for replica mode (recommended)
export DOPA_DB_PATH="./dopadone-replica.db"
export DOPA_DB_MODE=replica
```

**Option 2: YAML Configuration File**

Create or update `dopadone.yaml`:

```yaml
# For Remote mode
database:
  mode: remote
  turso:
    url: libsql://dopadone-yourorg.turso.io
    token: ${TURSO_AUTH_TOKEN}

# OR for Replica mode (recommended for offline support)
database:
  mode: replica
  path: ./dopadone-replica.db
  sync_interval: 60s
  turso:
    url: libsql://dopadone-yourorg.turso.io
    token: ${TURSO_AUTH_TOKEN}
```

### 6.2 Verify Configuration

Test the new configuration:

```bash
# Test with explicit mode
dopa --db-mode remote areas list

# Test with config file (if using)
dopa areas list

# For replica mode, verify sync
dopa --db-mode replica areas list
# Check logs for "[TursoReplica] Synced" messages
```

### 6.3 Update Multiple Machines

If using Dopadone on multiple machines:

1. **Copy configuration file** to each machine
2. **Set environment variables** or create config file
3. **For replica mode**: Each machine gets its own local replica file
4. **Test connectivity** on each machine

---

## Step 7: Rollback Procedure

If migration fails or you need to revert, follow these procedures.

### Scenario 1: Before Import (Easy)

If you haven't imported to Turso yet:

```bash
# Nothing to do - original database is intact
# Just proceed with fixing the export/import issue
```

### Scenario 2: Import Failed (Medium)

If the import to Turso failed:

```bash
# 1. Check error messages
turso db shell dopadone < dopadone-export.sql 2>&1 | tee import-errors.log

# 2. Review errors
cat import-errors.log

# 3. Fix the dump file if needed
# Common issues:
#   - Schema mismatch: Apply migrations before export
#   - Duplicate data: Clear destination first
#   - Invalid SQL: Edit dump file manually

# 4. Delete Turso database and recreate
turso db destroy dopadone
turso db create dopadone --from-dump dopadone-export-fixed.sql
```

### Scenario 3: After Config Change (Medium)

If you've updated configuration and want to use SQLite again:

```bash
# Option A: Change environment variable
unset TURSO_DATABASE_URL
unset TURSO_AUTH_TOKEN
export DOPA_DB_MODE=local

# Option B: Use CLI flag
dopa --db-mode local --db ./dopadone-backup-TIMESTAMP.db areas list

# Option C: Update config file
# Edit dopadone.yaml and set mode: local
```

### Scenario 4: Data Corruption (Hard)

If Turso data is corrupted:

```bash
# 1. Destroy corrupted database
turso db destroy dopadone

# 2. Recreate from backup
turso db create dopadone --from-dump dopadone-export.sql

# OR restore from Turso backups (if enabled)
# Check Turso dashboard for backup options
```

### Scenario 5: Complete Rollback to SQLite

To completely revert to local SQLite:

```bash
# 1. Restore from backup
cp dopadone-backup-TIMESTAMP.db dopadone.db

# 2. Update configuration
export DOPA_DB_MODE=local
unset TURSO_DATABASE_URL
unset TURSO_AUTH_TOKEN

# 3. Verify
dopa areas list

# 4. (Optional) Delete Turso database
turso db destroy dopadone
```

---

## Common Pitfalls and Solutions

### Pitfall 1: Database Too Large

**Symptom:** Import fails with size error

```
Error: database file too large (max 2GB)
```

**Solutions:**

```bash
# 1. Check size
ls -lh dopadone.db

# 2. If over limit, clean up old data
# Example: Archive old completed tasks
sqlite3 dopadone.db "DELETE FROM tasks WHERE status = 'done' AND updated_at < date('now', '-1 year');"

# 3. Vacuum to reclaim space
sqlite3 dopadone.db "VACUUM;"

# 4. Re-check size and retry
ls -lh dopadone.db

# 5. If still too large, consider replica mode
# (syncs incrementally rather than full import)
```

### Pitfall 2: Schema Mismatch

**Symptom:** Import fails with SQL errors

```
Error: no such table: areas
Error: table areas has no column named new_column
```

**Solution:**

```bash
# 1. Ensure migrations applied to SQLite before export
dopa migrate status
dopa migrate up

# 2. Verify schema version
sqlite3 dopadone.db "SELECT * FROM goose_db_version ORDER BY id;"

# 3. Re-export after migrations
sqlite3 dopadone.db .dump > dopadone-export.sql

# 4. If Turso already has partial schema, destroy and recreate
turso db destroy dopadone
turso db create dopadone --from-dump dopadone-export.sql
```

### Pitfall 3: Special Characters in Data

**Symptom:** Import partially fails, some records missing

**Solutions:**

```bash
# 1. Use binary-safe methods
# Avoid: Manual copy-paste of dump content
# Use: --from-dump or --from-file

# 2. Check for encoding issues
file dopadone.db
file dopadone-export.sql

# 3. If encoding is wrong, re-export with proper encoding
sqlite3 dopadone.db ".dump" | iconv -f UTF-8 -t UTF-8 > dopadone-export-fixed.sql
```

### Pitfall 4: Foreign Key Violations

**Symptom:** Import fails on foreign key constraints

```
Error: FOREIGN KEY constraint failed
```

**Solutions:**

```bash
# 1. The dump should handle this correctly, but if issues occur:
# Verify dump has tables in correct order (areas before projects before tasks)
head -100 dopadone-export.sql | grep "CREATE TABLE"

# 2. If order is wrong, manually reorder or use --from-file instead
turso db create dopadone --from-file dopadone.db
```

### Pitfall 5: Token Expired During Import

**Symptom:** Import stops mid-way with authentication error

**Solutions:**

```bash
# 1. Use non-expiring token for migration
turso db tokens create dopadone --expiration never

# 2. If using --from-dump, token not needed for creation
# Only needed for subsequent operations

# 3. If import fails, re-run (Turso will error on duplicate tables)
# Better to destroy and recreate:
turso db destroy dopadone
turso db create dopadone --from-dump dopadone-export.sql
```

### Pitfall 6: Network Interruption

**Symptom:** Incomplete import, timeout errors

**Solutions:**

```bash
# 1. For large databases, use --from-file (local file upload)
# Turso handles retries internally

# 2. Check network stability
ping -c 10 turso.io

# 3. If network is unreliable:
# - Use smaller dump file
# - Or use replica mode which handles network gracefully

# 4. Retry import after destroying failed database
turso db destroy dopadone
turso db create dopadone --from-dump dopadone-export.sql
```

### Pitfall 7: goose_db_version Missing

**Symptom:** Migration status shows wrong version after import

```
Migration status shows: 0 migrations applied (expected: 3)
```

**Solutions:**

```bash
# 1. Verify dump includes goose_db_version
grep "goose_db_version" dopadone-export.sql

# 2. If missing, re-export
sqlite3 dopadone.db .dump > dopadone-export.sql

# 3. Verify table content is included
grep -A5 "INSERT INTO goose_db_version" dopadone-export.sql

# 4. Re-import with correct dump
turso db destroy dopadone
turso db create dopadone --from-dump dopadone-export.sql
```

### Pitfall 8: Replica Mode Initial Sync Timeout

**Symptom:** First sync in replica mode times out

```
[TursoReplica] Sync failed: context deadline exceeded
```

**Solutions:**

```bash
# 1. Wait and retry - initial sync can take time for large databases
dopa --db-mode replica areas list

# 2. If persistent, check logs for specifics
# Look for frame sync progress

# 3. As fallback, use remote mode first, then switch to replica
# After migration works in remote mode:
export DOPA_DB_MODE=replica
dopa areas list  # Will create fresh replica from remote
```

---

## Migration Timeline

A typical migration follows this timeline:

| Phase | Duration | Activities |
|-------|----------|------------|
| Preparation | 5-10 min | Backup, verify integrity, document state |
| Export | 1-5 min | Create SQL dump or prepare file |
| Turso Setup | 2-5 min | Create database, get credentials |
| Import | 1-10 min | Import data (depends on size) |
| Verification | 5-10 min | Verify counts, test functionality |
| Configuration | 2-5 min | Update config, test connectivity |
| **Total** | **15-45 min** | Complete migration |

---

## Post-Migration Checklist

After completing the migration:

- [ ] All record counts match between SQLite and Turso
- [ ] Read operations work (list areas, projects, tasks)
- [ ] Write operations work (create, update, delete)
- [ ] Migration status shows all migrations applied
- [ ] Configuration is saved and tested
- [ ] Backup is archived (don't delete for 1-2 weeks)
- [ ] Documentation updated for team (if applicable)

---

## Related Documentation

- [Turso Setup Guide](TURSO_SETUP.md) - Account and credential setup
- [Turso Migrations](TURSO_MIGRATIONS.md) - Schema migration procedures
- [Turso Troubleshooting](TURSO_TROUBLESHOOTING.md) - Error solutions
- [Database Modes](DATABASE_MODES.md) - Mode configuration and comparison
- [Turso Performance](TURSO_PERFORMANCE.md) - Performance optimization

---

## Getting Help

If you encounter issues:

1. **Check this guide** - Most common issues are documented in Common Pitfalls
2. **Run verification** - Follow Step 5 to diagnose the issue
3. **Check Turso status** - [status.turso.tech](https://status.turso.tech)
4. **Turso Discord** - [tur.so/discord](https://tur.so/discord)
5. **GitHub Issues** - Report bugs in the Dopadone repository
