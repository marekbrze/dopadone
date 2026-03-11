---
id: TASK-60.9
title: 'Documentation: Turso setup and configuration guide'
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-08 19:02'
updated_date: '2026-03-11 15:51'
labels:
  - documentation
  - turso
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create comprehensive documentation for Turso database integration including setup instructions, configuration examples, and troubleshooting. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Document Turso account setup and database creation
- [ ] #2 Document CLI configuration examples for all three modes
- [ ] #3 Document environment variable configuration
- [ ] #4 Document config file YAML structure
- [ ] #5 Add migration guide from local SQLite to Turso
- [ ] #6 Add troubleshooting section for common issues
- [x] #7 Include performance considerations and best practices
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for TASK-60.9: Documentation: Turso setup and configuration guide

## Analysis Summary

The task has been split into 5 subtasks due to:
1. **AC#4 (YAML config)** requires implementation first - documented in 60.9.1
2. Documentation is extensive and benefits from focused subtasks
3. Clear dependencies enable parallel work where possible

## Subtask Breakdown

### TASK-60.9.1: YAML Config File Support (Prerequisite)
**Priority**: HIGH | **Type**: Coding + Documentation
**AC Coverage**: AC#4 (config file YAML structure)

**Implementation Steps**:
1. Research Go YAML config libraries (viper, yaml.v3)
2. Define YAML schema for database configuration
3. Implement config file discovery (./dopadone.yaml, ~/.config/dopadone/config.yaml)
4. Implement config file parser with validation
5. Integrate with existing precedence: CLI > config > env > defaults
6. Add --config flag for custom config path
7. Write unit tests for parsing and validation
8. Write integration tests for precedence chain
9. Document YAML schema in DATABASE_MODES.md

**Tests Required**:
- Unit: YAML parsing, validation, precedence
- Integration: Full config chain with all modes

**Dependencies**: None (can start immediately)

---

### TASK-60.9.2: Turso Getting Started Guide
**Priority**: MEDIUM | **Type**: Documentation
**AC Coverage**: AC#1 (Turso account setup and database creation)

**Implementation Steps**:
1. Create docs/TURSO_SETUP.md
2. Document Turso account signup
3. Document Turso CLI installation
4. Document database creation (CLI + web UI)
5. Document token generation
6. Document finding database URL
7. Add quick start examples
8. Cross-reference with DATABASE_MODES.md

**Tests Required**:
- Manual: Verify all documented steps work
- Review: Check links to official docs are valid

**Dependencies**: TASK-60.9.1 (for YAML config examples)

---

### TASK-60.9.5: SQLite to Turso Migration Guide
**Priority**: MEDIUM | **Type**: Documentation
**AC Coverage**: AC#5 (migration guide from local SQLite to Turso)

**Implementation Steps**:
1. Enhance docs/TURSO_MIGRATIONS.md
2. Create migration checklist
3. Document data export (sqlite3 .dump, .backup)
4. Document Turso database setup
5. Document data import methods
6. Document verification steps
7. Document rollback procedure
8. Add common pitfalls section

**Tests Required**:
- Manual: Walk through complete migration
- Validation: Test rollback procedure

**Dependencies**: TASK-60.9.2 (needs setup guide)

---

### TASK-60.9.4: Performance Best Practices
**Priority**: MEDIUM | **Type**: Documentation
**AC Coverage**: AC#7 (performance considerations and best practices)

**Implementation Steps**:
1. Create docs/TURSO_PERFORMANCE.md or add section to DATABASE_MODES.md
2. Document performance characteristics per mode
3. Document sync interval tuning
4. Document connection pooling
5. Create benchmark results table
6. Document mode selection guide
7. Add optimization tips

**Tests Required**:
- Benchmark: Run benchmarks for each mode
- Validation: Verify recommendations are accurate

**Dependencies**: TASK-60.9.1 (needs YAML config for examples)

---

### TASK-60.9.3: Comprehensive Troubleshooting Guide
**Priority**: MEDIUM | **Type**: Documentation
**AC Coverage**: AC#6 (troubleshooting section for common issues)

**Implementation Steps**:
1. Create docs/TURSO_TROUBLESHOOTING.md
2. Document connection errors
3. Document authentication failures
4. Document network issues
5. Document sync failures
6. Document migration errors
7. Create error message reference
8. Add diagnostic commands

**Tests Required**:
- Manual: Verify diagnostic commands work
- Review: Ensure all common errors are covered

**Dependencies**: TASK-60.9.5 (needs migration guide for migration errors)

---

## Execution Order (Parallel Tracks)

### Track 1 (Sequential - Critical Path)
```
60.9.1 (YAML config) → 60.9.2 (Setup guide) → 60.9.5 (Migration) → 60.9.3 (Troubleshooting)
```

### Track 2 (Parallel - After 60.9.1)
```
60.9.1 (YAML config) → 60.9.4 (Performance)
```

### Final Integration
```
All tasks complete → Update DATABASE_MODES.md → Validate → Mark 60.9 Done
```

## Dependency Graph
```
    ┌─────────────────┐
    │   60.9.1        │  HIGH - YAML Config (Prerequisite)
    │   (YAML Config) │
    └────────┬────────┘
             │
     ┌───────┴───────┐
     ▼               ▼
┌─────────────┐ ┌─────────────┐
│  60.9.2     │ │  60.9.4     │  MEDIUM - Can run in parallel
│  (Setup)    │ │  (Perf)     │
└──────┬──────┘ └─────────────┘
       │
       ▼
┌─────────────┐
│  60.9.5     │  MEDIUM - Migration Guide
│  (Migration)│
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  60.9.3     │  MEDIUM - Troubleshooting
│  (Trouble)  │
└─────────────┘
```

## AC Coverage Mapping

| AC | Subtask | Status |
|----|----------|--------|
| #1 Turso account setup | 60.9.2 | Pending |
| #2 CLI configuration | Existing (DATABASE_MODES.md) | Done |
| #3 Environment variables | Existing (DATABASE_MODES.md) | Done |
| #4 Config file YAML | 60.9.1 | Pending |
| #5 Migration guide | 60.9.5 | Pending |
| #6 Troubleshooting | 60.9.3 | Pending |
| #7 Performance | 60.9.4 | Pending |

## Documentation Updates Required

### Files to Create
- docs/TURSO_SETUP.md (60.9.2)
- docs/TURSO_TROUBLESHOOTING.md (60.9.3)
- docs/TURSO_PERFORMANCE.md or section in DATABASE_MODES.md (60.9.4)

### Files to Update
- docs/DATABASE_MODES.md - Add YAML config section, cross-references
- docs/TURSO_MIGRATIONS.md - Enhance with migration guide (60.9.5)
- docs/START_HERE.md - Update documentation index

## Testing Strategy

### Unit Tests (60.9.1)
- YAML parsing
- Config validation
- Precedence chain

### Integration Tests (60.9.1)
- Full config chain with all modes
- CLI flag override tests

### Manual Testing (All)
- Walk through all documented procedures
- Verify examples work as documented
- Test error scenarios for troubleshooting

## Estimated Time

| Subtask | Type | Hours |
|---------|------|-------|
| 60.9.1 | Coding + Tests | 4-6 |
| 60.9.2 | Documentation | 2-3 |
| 60.9.3 | Documentation | 2-3 |
| 60.9.4 | Documentation | 2-3 |
| 60.9.5 | Documentation | 2-3 |
| Integration | Documentation | 1-2 |
| **Total** | | **13-20 hours** |

## Success Criteria

- ✅ All 7 acceptance criteria met
- ✅ YAML config file implemented and documented
- ✅ All documentation files created/updated
- ✅ All examples validated and working
- ✅ Cross-references between documents
- ✅ START_HERE.md updated with new docs
- ✅ No orphaned documentation
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Analysis (2026-03-11)

### Gap Identified
AC#4 (Document config file YAML structure) cannot be completed because **YAML config file support is not implemented**. The current implementation only supports:
- CLI flags (--db, --turso-url, --turso-token, --db-mode, --sync-interval)
- Environment variables (DOPA_DB_PATH, TURSO_DATABASE_URL, TURSO_AUTH_TOKEN, DOPA_DB_MODE)

### Decision: Split Task
Task split into 5 subtasks:
- TASK-60.9.1: YAML config file implementation (HIGH - prerequisite)
- TASK-60.9.2: Turso Getting Started Guide
- TASK-60.9.3: Comprehensive Troubleshooting Guide
- TASK-60.9.4: Performance Best Practices Guide
- TASK-60.9.5: SQLite to Turso Migration Guide

### Existing Documentation
- DATABASE_MODES.md covers AC#2, AC#3 (CLI + env vars)
- TURSO_MIGRATIONS.md partially covers AC#5
- 08-database-drivers.md covers technical architecture

### Parallel Execution Possible
- 60.9.2 and 60.9.4 can run in parallel after 60.9.1
- 60.9.5 depends on 60.9.2
- 60.9.3 depends on 60.9.5
<!-- SECTION:NOTES:END -->
