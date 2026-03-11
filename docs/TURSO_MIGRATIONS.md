# Migration Guide for libSQL (Turso)

This document covers migration considerations when using libSQL/Turso with Dopadone.

## Overview

Dopadone supports three database modes:
- **Local SQLite**: Standard local database file
- **Turso Remote**: Direct connection to Turso cloud database
- **Turso Replica**: Embedded replica with local file and sync to cloud

All modes use **goose v3** for migrations with the `sqlite3` dialect.

## Migration Commands

### Standard Migrations

```bash
# Apply all pending migrations
dopa migrate up

# Rollback last migration
dopa migrate down

# Show migration status
dopa migrate status

# Reset database (rollback all, then apply all)
dopa migrate reset
```

### Schema Verification

```bash
# Verify database schema consistency
dopa migrate verify
```

This command checks:
- All expected tables exist
- Goose version table is present
- Schema is consistent with migrations

## libSQL-Specific Considerations

### Goose Dialect Compatibility

The `sqlite3` dialect works correctly with libSQL because libSQL maintains SQLite compatibility. No dialect changes are required.

### Migration Sync Strategy

For **Turso Replica** mode:

1. Migrations run locally against the embedded replica file
2. After successful migration, the schema is automatically synced to Turso cloud
3. This ensures local and remote schemas stay consistent

```bash
# Migrations automatically sync for replica mode
dopa --db-mode replica migrate up
```

For **Turso Remote** mode:
- Migrations run directly against the remote database
- No sync is needed (direct remote writes)

### Schema Drift Detection

Use `dopa migrate verify` to detect schema drift between local and remote:

```bash
# For replica mode, verify local schema matches expectations
dopa migrate verify
```

### Transaction Behavior

libSQL handles transactions similarly to SQLite:
- DDL statements (CREATE, ALTER, DROP) are transactional
- Migration rollbacks work as expected
- No special handling required

## Best Practices

### 1. Run Migrations Before Sync

For embedded replica mode, always run migrations locally first:

```bash
# Correct order
dopa migrate up    # Runs locally, then syncs
```

### 2. Verify After Migration

Always verify schema after migrations:

```bash
dopa migrate verify
```

### 3. Backup Before Major Changes

For production databases:

```bash
# Turso automatically maintains backups
# For local SQLite, copy the file before migration
cp dopadone.db dopadone.db.backup
dopa migrate up
```

### 4. Handle Sync Failures

If sync fails after migration:

1. The migration is applied locally
2. Check network connectivity
3. Retry with `dopa migrate status` to verify
4. Manually trigger sync if needed (replica mode)

## Troubleshooting

For comprehensive migration troubleshooting, see the [Turso Troubleshooting Guide](TURSO_TROUBLESHOOTING.md) which covers:

- Migration errors specific to libSQL
- Schema drift detection and resolution
- Rollback procedures
- Sync failures after migration

### Quick Migration Troubleshooting

**Migration Fails** (`goose: migration failed`):
1. Check SQL syntax is SQLite-compatible
2. Run `dopa migrate status` to see current state
3. See [Migration Errors](TURSO_TROUBLESHOOTING.md#migration-errors) for details

**Sync Fails After Migration**:
1. Check Turso credentials
2. Verify network connectivity
3. See [Replica Mode Sync Issues](TURSO_TROUBLESHOOTING.md#replica-mode-sync-issues)

**Schema Drift Detected**:
1. Check `goose_db_version` table versions
2. Run `dopa migrate verify`
3. See [Migration Errors](TURSO_TROUBLESHOOTING.md#migration-errors)

## Migration Files

Migrations are stored in `internal/migrate/migrations/`:

| File | Description |
|------|-------------|
| `20240301000000_initial_schema.sql` | Initial schema (areas, subareas, projects) |
| `20260303110742_add_tasks_table.sql` | Tasks table |
| `20260304120000_add_sort_order_to_areas.sql` | Sort order for areas |

### Creating New Migrations

1. Create a new file: `YYYYMMDDHHMMSS_description.sql`
2. Include `-- +goose Up` and `-- +goose Down` sections
3. Use standard SQLite syntax (libSQL compatible)
4. Test locally before deploying

## Environment Variables

| Variable | Description |
|----------|-------------|
| `TURSO_DATABASE_URL` | Turso database URL |
| `TURSO_AUTH_TOKEN` | Turso authentication token |
| `DOPA_DB_MODE` | Database mode: `local`, `remote`, `replica`, `auto` |
| `DOPA_DB_PATH` | Local database file path |

## Related Documentation

- [Turso Setup Guide](TURSO_SETUP.md) - Step-by-step Turso account and database setup
- [Database Modes](DATABASE_MODES.md) - Detailed mode explanations
- [Turso Documentation](https://docs.turso.tech) - Official Turso docs
- [Goose Documentation](https://github.com/pressly/goose) - Migration tool docs
