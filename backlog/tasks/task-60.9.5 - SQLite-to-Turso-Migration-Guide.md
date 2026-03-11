---
id: TASK-60.9.5
title: SQLite to Turso Migration Guide
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-11 14:19'
updated_date: '2026-03-11 16:16'
labels:
  - documentation
  - turso
  - migration
dependencies:
  - TASK-60.9.2
parent_task_id: TASK-60.9
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Enhance the migration documentation with step-by-step guide for migrating from local SQLite to Turso. This addresses AC#5 of TASK-60.9. Part of task-60.9 documentation effort.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create step-by-step migration checklist
- [x] #2 Document data export from SQLite
- [x] #3 Document Turso database creation and setup
- [x] #4 Document data import to Turso
- [x] #5 Document verification steps to ensure data integrity
- [x] #6 Document rollback procedure if migration fails
- [x] #7 Include common pitfalls and how to avoid them
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for TASK-60.9.5: SQLite to Turso Migration Guide

## Task Analysis

### Scope Clarification
This task creates a **DATA migration guide** (moving existing SQLite data to Turso), NOT a schema migration guide (goose migrations are already documented in TURSO_MIGRATIONS.md and TURSO_TROUBLESHOOTING.md).

### Dependencies Status
- TASK-60.9.2 (Turso Getting Started Guide): COMPLETE ✓
- Turso setup documentation exists - can reference for AC#3

### Existing Documentation Audit
| Document | Content | Gap |
|----------|---------|-----|
| TURSO_MIGRATIONS.md | Schema migrations (goose) | No data migration |
| TURSO_TROUBLESHOOTING.md | Migration error troubleshooting | No migration procedure |
| TURSO_SETUP.md | Turso account/DB setup | No import from existing data |
| DATABASE_MODES.md | Mode configuration | No migration path |

### Task Size Assessment
- 7 ACs, all documentation-focused
- Linear process (export → setup → import → verify → rollback)
- Estimated: 3-4 hours
- **Decision: No splitting needed** - cohesive single documentation effort

---

## PHASE 1: Research and Preparation (Sequential)

### 1.1. Research Turso Data Import Methods
- Review Turso CLI: `turso db create --from-file`, `--from-dump`
- Review Turso shell: `turso db shell` for SQL imports
- Research sqlite3 export methods: `.dump`, `.backup`, `.read`
- Document file size limits (Turso: 2GB max)
- Research libSQL compatibility with SQLite exports

### 1.2. Analyze Existing Codebase
- Review internal/db/ for database file paths
- Review internal/migrate/ for schema information
- Check for any migration-related CLI commands
- Identify data types and tables to migrate

### 1.3. Create Test Migration Environment
- Create test SQLite database with sample data
- Test export methods locally
- Prepare test cases for validation

---

## PHASE 2: Document Migration Procedures (Sequential)

### 2.1. AC#1: Migration Checklist (Section 1)
Create a pre-flight checklist:
- [ ] Backup existing database
- [ ] Verify SQLite integrity
- [ ] Check database size (≤2GB for Turso)
- [ ] Create Turso account/database
- [ ] Generate auth token
- [ ] Export data from SQLite
- [ ] Import data to Turso
- [ ] Verify data integrity
- [ ] Update Dopadone configuration
- [ ] Test application connectivity
- [ ] Archive old database (don't delete yet)

### 2.2. AC#2: Data Export Documentation (Section 2)
Document multiple export methods:

**Method 1: SQLite Dump (Recommended)**
```bash
sqlite3 dopadone.db .dump > dopadone-export.sql
```

**Method 2: SQLite Backup**
```bash
sqlite3 dopadone.db ".backup dopadone-backup.db"
```

**Method 3: Online Backup API (for production)**
```bash
# Using sqlite3 shell with .backup
# Safe for databases in use
```

Include:
- Export command syntax
- File size considerations
- Verification of export
- Handling large databases

### 2.3. AC#3: Turso Database Setup Reference (Section 3)
**Cross-reference, don't duplicate:**
- Link to TURSO_SETUP.md for account creation
- Link to TURSO_SETUP.md for database creation
- Document Turso-specific import commands:
  ```bash
  # Create from existing dump
  turso db create dopadone --from-dump dopadone-export.sql
  
  # Create from existing file
  turso db create dopadone --from-file dopadone.db
  ```
- Document getting credentials for import

### 2.4. AC#4: Data Import Documentation (Section 4)
Document import methods:

**Method 1: Create Database from Dump**
```bash
turso db create dopadone --from-dump dopadone-export.sql
```

**Method 2: Create Database from File**
```bash
turso db create dopadone --from-file dopadone.db
```

**Method 3: Shell Import (for existing databases)**
```bash
# For appending to existing database
turso db shell dopadone < dopadone-export.sql
```

Include:
- Import command syntax
- Size limits and error handling
- Schema compatibility notes
- Progress monitoring

### 2.5. AC#5: Verification Steps (Section 5)
Document comprehensive verification:

**Schema Verification**
```bash
# Check migration status
dopa migrate status

# Verify schema
dopa migrate verify
```

**Data Integrity Verification**
```bash
# Count records in source
sqlite3 dopadone.db "SELECT COUNT(*) FROM areas;"
sqlite3 dopadone.db "SELECT COUNT(*) FROM projects;"
sqlite3 dopadone.db "SELECT COUNT(*) FROM tasks;"

# Count records in Turso (after import)
dopa areas list | wc -l
dopa projects list | wc -l
dopa tasks list | wc -l

# Compare counts - should match
```

**Functional Verification**
```bash
# Test CRUD operations
dopa areas list
dopa projects list
dopa tasks list

# Test write operations
dopa tasks create "Test task" --project-id <id>
```

**Checksum Verification (Optional)**
```bash
# Export both and compare
sqlite3 source.db .dump > source-dump.sql
turso db shell turso-db .dump > turso-dump.sql
diff source-dump.sql turso-dump.sql
```

### 2.6. AC#6: Rollback Procedure (Section 6)
Document rollback scenarios:

**Scenario 1: Before Import (Easy)**
- Original database is still intact
- Delete Turso database and retry

**Scenario 2: After Import (Medium)**
- Original database backed up
- Restore from backup
- Delete Turso database

**Scenario 3: After Config Change (Medium)**
- Revert Dopadone configuration
- Use original SQLite database

**Scenario 4: Data Corruption (Hard)**
- Delete Turso database
- Re-import from backup
- Or restore from Turso backups

Include rollback commands for each scenario.

### 2.7. AC#7: Common Pitfalls (Section 7)
Document pitfalls and solutions:

1. **Pitfall: Database too large**
   - Symptom: Import fails with size error
   - Solution: Split data, or use replica mode

2. **Pitfall: Schema mismatch**
   - Symptom: Import fails with SQL errors
   - Solution: Ensure migrations applied before export

3. **Pitfall: Special characters in data**
   - Symptom: Import partially fails
   - Solution: Use binary-safe export method

4. **Pitfall: Foreign key violations**
   - Symptom: Import fails on constraints
   - Solution: Import in correct order, or disable FK temporarily

5. **Pitfall: Token expired during import**
   - Symptom: Import stops mid-way
   - Solution: Use non-expiring token for migration

6. **Pitfall: Network interruption**
   - Symptom: Incomplete import
   - Solution: Use retry, or split into chunks

7. **Pitfall: goose_db_version table**
   - Symptom: Migration status shows wrong version
   - Solution: Ensure goose_db_version is exported

---

## PHASE 3: Document Structure and Integration (Sequential)

### 3.1. Create docs/TURSO_DATA_MIGRATION.md
Structure:
```markdown
# SQLite to Turso Data Migration Guide

## Overview
## Prerequisites
## Migration Checklist (AC#1)
## Step 1: Prepare and Backup
## Step 2: Export Data from SQLite (AC#2)
## Step 3: Create Turso Database (AC#3)
## Step 4: Import Data to Turso (AC#4)
## Step 5: Verify Data Integrity (AC#5)
## Step 6: Update Configuration
## Step 7: Rollback Procedure (AC#6)
## Common Pitfalls (AC#7)
## Related Documentation
```

### 3.2. Update docs/TURSO_MIGRATIONS.md
Add section at top:
```markdown
## Data Migration

To migrate your existing SQLite data to Turso, see:
- [SQLite to Turso Data Migration Guide](TURSO_DATA_MIGRATION.md)

This guide covers:
- Exporting data from local SQLite
- Importing to Turso cloud
- Verification and rollback procedures
```

### 3.3. Update docs/START_HERE.md
Add to documentation index:
```markdown
| [Turso Data Migration](TURSO_DATA_MIGRATION.md) | Step-by-step guide for migrating data from SQLite to Turso |
```

### 3.4. Update docs/DATABASE_MODES.md
Add cross-reference in related docs section:
```markdown
- [Turso Data Migration](TURSO_DATA_MIGRATION.md) - Migrating data from SQLite to Turso
```

### 3.5. Update docs/TURSO_TROUBLESHOOTING.md
Add reference to data migration:
```markdown
For data migration issues, see [SQLite to Turso Data Migration Guide](TURSO_DATA_MIGRATION.md#common-pitfalls).
```

---

## PHASE 4: Testing and Validation (Sequential)

### 4.1. Manual Migration Test
Walk through the complete guide:
1. Create test SQLite database with data
2. Follow export procedure
3. Import to Turso
4. Verify all data present
5. Test rollback procedure
6. Document any issues found

### 4.2. Edge Case Testing
- Empty database migration
- Large database (near 2GB limit)
- Special characters in data
- Unicode/emoji in data
- NULL values
- Foreign key relationships

### 4.3. Documentation Review
- Verify all commands work
- Check all links are valid
- Ensure consistent formatting
- Review for clarity

---

## TESTING STRATEGY

### Manual Tests (All ACs)
1. Complete migration walkthrough
2. Rollback procedure test
3. Pitfall reproduction and solution verification

### Validation Tests
- Command syntax verification
- Link validation (all cross-references)
- Markdown linting

### No Automated Tests Required
This is pure documentation with no code changes.

---

## FILES TO CREATE
- docs/TURSO_DATA_MIGRATION.md (main documentation, ~400 lines)

## FILES TO UPDATE
- docs/TURSO_MIGRATIONS.md (add data migration reference)
- docs/START_HERE.md (add to documentation index)
- docs/DATABASE_MODES.md (add cross-reference)
- docs/TURSO_TROUBLESHOOTING.md (add reference)

---

## PARALLEL OPPORTUNITIES
**None** - Documentation phases have sequential dependencies

---

## SEQUENTIAL EXECUTION ORDER
1. Phase 1: Research (1 hour)
2. Phase 2: Document procedures (2 hours)
3. Phase 3: Integration (30 minutes)
4. Phase 4: Testing (30 minutes)

**Total Estimated Time: 3-4 hours**

---

## SUCCESS CRITERIA
- ✅ All 7 ACs documented
- ✅ Complete migration walkthrough works
- ✅ Rollback procedure tested
- ✅ All pitfalls documented with solutions
- ✅ Cross-references updated
- ✅ Manual validation passes
- ✅ START_HERE.md updated
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Started implementation. Reading existing Turso documentation to understand current state and cross-references.

Created docs/TURSO_DATA_MIGRATION.md with all 7 ACs: migration checklist, data export, Turso setup, data import, verification steps, rollback procedure, and common pitfalls.

Updated cross-references in TURSO_MIGRATIONS.md, START_HERE.md, DATABASE_MODES.md, and TURSO_TROUBLESHOOTING.md.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created comprehensive SQLite to Turso data migration guide (docs/TURSO_DATA_MIGRATION.md) covering:

**New Documentation:**
- docs/TURSO_DATA_MIGRATION.md (~600 lines) with 7 complete sections:
  1. Migration checklist with pre-flight verification
  2. Data export methods (SQL dump, binary file, online backup)
  3. Turso database creation from dump/file
  4. Data import procedures for existing databases
  5. Comprehensive verification steps (schema, counts, functional, checksums)
  6. Rollback procedures for 5 different scenarios
  7. Common pitfalls with 8 detailed solutions

**Updated Cross-References:**
- docs/TURSO_MIGRATIONS.md - Added data migration reference section
- docs/START_HERE.md - Added to documentation index
- docs/DATABASE_MODES.md - Added related documentation link
- docs/TURSO_TROUBLESHOOTING.md - Added data migration reference

Build passes. Pre-existing lint issues and test failure in migrate_libsql_test.go are unrelated to this documentation change.
<!-- SECTION:FINAL_SUMMARY:END -->
