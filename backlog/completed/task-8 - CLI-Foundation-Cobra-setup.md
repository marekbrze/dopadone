---
id: TASK-8
title: CLI Foundation - Cobra setup
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 09:39'
updated_date: '2026-03-03 09:47'
labels:
  - cli
  - cobra
  - foundation
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Set up Cobra CLI structure with root command, database connection handling, common output formatters (table/JSON), and error handling utilities.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Cobra root command at cmd/projectdb/main.go with --db flag (default: ./projectdb.db)
- [x] #2 Database connection helper validates path and provides *sql.DB
- [x] #3 Output formatter package supports table (colored headers) and JSON
- [x] #4 Error handling with exit codes: 0=success, 1=error, 2=validation
- [x] #5 Domain validation adapter using existing types
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Add Cobra and lipgloss dependencies to go.mod (go get github.com/spf13/cobra@latest github.com/charmbracelet/lipgloss)

2. Create cmd/projectdb/main.go with Cobra root command:
   - Persistent --db flag (default: ./projectdb.db)
   - Persistent --output/-o flag (table/json)
   - PreRun hook to initialize DB connection

3. Create internal/cli/db.go:
   - Connect(path string) (*sql.DB, error) - validates path exists, opens connection
   - RunMigrations(db, migrationsPath) - run goose migrations
   - Close(db) helper

4. Create internal/cli/output/formatter.go:
   - Formatter interface with Table(), JSON() methods
   - NewFormatter(format string) factory
   - TableFormatter using lipgloss for colored headers
   - JSONFormatter with pretty printing

5. Create internal/cli/errors.go:
   - Exit codes: ExitSuccess=0, ExitError=1, ExitValidation=2
   - ValidationError type for input validation failures
   - WrapError(err, context) for error context
   - ErrorHandler middleware for Cobra commands

6. Create internal/cli/validation.go:
   - Adapter functions wrapping domain types
   - ParseProjectStatus, ParsePriority, ParseProgress, ParseColor, ParseDateRange
   - Map domain errors to CLI validation errors

7. Update cmd/projectdb/main.go to wire everything together

8. Add tests:
   - internal/cli/output/formatter_test.go
   - internal/cli/errors_test.go
   - internal/cli/validation_test.go
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Successfully implemented Cobra CLI foundation with all required components.

Changes:
- Added Cobra (v1.10.2) and lipgloss (v1.1.0) dependencies
- Created cmd/projectdb/main.go with root command, --db flag (default: ./projectdb.db), and --output flag (table|json)
- Implemented internal/cli/db.go with Connect(), Close(), and RunMigrations() helpers for database connection management
- Implemented internal/cli/output/formatter.go with Formatter interface supporting table (colored headers using lipgloss) and JSON output formats
- Implemented internal/cli/errors.go with exit codes (0=success, 1=error, 2=validation), ValidationError type, and error handling utilities
- Implemented internal/cli/validation.go with adapter functions wrapping domain types (ParseProjectStatus, ParsePriority, ParseProgress, ParseColor, ParseDate, ValidateProjectName)

Tests:
- All tests passing: internal/cli/errors_test.go, internal/cli/validation_test.go, internal/cli/output/formatter_test.go
- Binary builds successfully and CLI commands work correctly
<!-- SECTION:FINAL_SUMMARY:END -->
