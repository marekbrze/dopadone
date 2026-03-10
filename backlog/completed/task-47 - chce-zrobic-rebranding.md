---
id: TASK-47
title: chce zrobic rebranding
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 16:37'
updated_date: '2026-03-06 17:49'
labels: []
dependencies: []
references:
  - README.md
  - docs/START_HERE.md
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
chce zrobic rebranding aplikacji na dopadone i glowna komenda to powinno byc dopa
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 User can execute rebranding via CLI command
- [x] #2 All command help text and examples changed from 'dopa' to 'dopa'
- [x] #3 All imports in Go files changed from 'github.com/example/dopa' to 'github.com/example/dopadone'
- [x] #4 Database default path changed from './dopa.db' to './dopadone.db'
- [x] #5 Documentation updated (README.md, docs/, backlog/)
- [x] #6 Makefile updated (binary name and ldflags)
- [x] #7 Command directory moved from cmd/dopa/ to cmd/dopa/
- [x] #8 All references in documentation updated from 'Dopadone'/'dopa' to 'Dopadone'/'dopa'
- [x] #9 Scripts updated (if they contain dopa references)
- [x] #10 Rebranding checklist created with all items verified completion
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Phase 1: Foundation (Sequential - Must Complete First)
## Task 47.1: Core Module Rebranding
1. Update go.mod module path from github.com/example/dopa to github.com/example/dopadone
2. Run go mod tidy to update dependencies
3. Use gomvpkg or manual find/replace to update all import paths in Go files
4. Move cmd/dopa/ directory to cmd/dopa/
5. Update all Go files to import from new module path
6. Update internal package imports
7. Run go build ./... to verify compilation
8. Run go test ./... to verify all tests pass

## Task 47.2: Build System Updates (Depends on 47.1)
1. Update Makefile BINARY_NAME from dopa to dopa
2. Update Makefile DB_PATH from dopa.db to dopadone.db
3. Update LDFLAGS module path references
4. Update all build targets (build-linux, build-darwin, etc.)
5. Update dist target archive names
6. Run make clean && make build to verify
7. Run make test to verify tests still work

# Phase 2: Documentation (Can Run in Parallel with Phase 1)
## Task 47.3: Documentation Rebranding
1. Update README.md: Change Dopadone → Dopadone, dopa → dopa
2. Update docs/START_HERE.md title and references
3. Update docs/architecture/*.md files
4. Update docs/TUI.md, docs/RELEASE.md, etc.
5. Update all command examples in documentation
6. Update project structure diagram in docs
7. Update installation instructions
8. Verify all internal doc links still work

# Phase 3: Scripts and Configuration
## Task 47.4: Scripts and Utilities (Depends on 47.1)
1. Update scripts/install.sh references to dopa
2. Update scripts/seed-test-data.sh database path references
3. Update scripts/generate-changelog.sh if needed
4. Check for any other shell scripts
5. Test scripts with new binary name

# Phase 4: Verification and Finalization
## Task 47.5: Comprehensive Verification
1. Create REBRANDING_CHECKLIST.md with all changes verified
2. Run make clean && make build-all to verify cross-platform builds
3. Run full test suite: make test
4. Run linter: make lint
5. Test CLI commands manually: ./bin/dopa --help, ./bin/dopa task list, etc.
6. Verify TUI works: ./bin/dopa tui
7. Check database operations with new path
8. Verify all acceptance criteria are met

# Dependencies:
- 47.1 must complete before 47.2, 47.4, 47.5
- 47.3 can run in parallel with 47.1 and 47.2
- 47.5 must run after all other tasks

# Parallel Execution Strategy:
- Phase 1 (47.1) → Sequential, foundation work
- Phase 2 (47.3) → Can start immediately, documentation only
- Phase 1.5 (47.2) → After 47.1 completes
- Phase 3 (47.4) → After 47.1 completes, can run in parallel with 47.2
- Phase 4 (47.5) → Final verification after all tasks
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Task 47 complete - AC1-10 complete

- Phase 1 (AC1-10 complete): go.mod module path updated
- Moved cmd/dopa/ directory to cmd/dopa/
- Updated all Go imports
- Updated Makefile
- Build and tests verified successfully
- Phase 2: Updated docs/
    - README.md
    - docs/START_HERE.md
    - docs/architecture/*.md
- Phase 3: Updated scripts/
    - scripts/install.sh
    - scripts/seed-test-data.sh
    - scripts/generate-changelog.sh
- Phase 4: Manual verification complete
    - Created REBRANDING_CHECKLIST.md
    - Build and tests pass
    - CLI works correctly (./bin/dopa --help, task list, TUI)
    - Database operations work with new path (./dopadone.db)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Comprehensive rebranding from Dopadone to Dopadone, transforming the CLI tool from 'dopa' to 'dopa'. This involves updating all source code imports, command-line interface, documentation, build configuration, and scripts. The rebranding creates a cohesive ADHD-friendly identity with the short, memorable command name 'dopa' and the descriptive name 'dopadone' in documentation and longer-form contexts.
<!-- SECTION:FINAL_SUMMARY:END -->
