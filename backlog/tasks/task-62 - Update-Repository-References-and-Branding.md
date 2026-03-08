---
id: TASK-62
title: Update Repository References and Branding
status: In Progress
assignee:
  - '@assistant'
created_date: '2026-03-07 21:49'
updated_date: '2026-03-07 22:04'
labels:
  - release
  - refactor
  - documentation
dependencies: []
references:
  - backlog/tasks/task-61 - release-first-version-of-the-app-on-github.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Replace all placeholder repository URLs and fix branding inconsistencies across the codebase to use the actual GitHub repository marekbrze/dopadone. This is part of the v1.0.0 release preparation (task-61).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Replace github.com/example/dopa with github.com/marekbrze/dopadone in: README.md (lines 17, 62, 68), Makefile (line 18 LDFLAGS), scripts/install.sh (line 7 REPO variable), go.mod module name, and all imports in .go files
- [x] #2 Fix version package to reference 'dopa' instead of 'projectdb' in internal/version/version.go (lines 35, 41, 94, 130, 186)
- [x] #3 Search for any remaining 'example' or 'projectdb' references: grep -r 'example/dopa' . && grep -r 'projectdb' .
- [x] #4 Run 'go mod tidy' after module name change
- [x] #5 Verify builds successfully with 'make build'
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Update go.mod module name from github.com/example/dopadone to github.com/marekbrze/dopadone
2. Update all import statements in .go files to use new module path
3. Update README.md repository URLs (lines 17, 62, 68)
4. Update Makefile LDFLAGS (line 18-20)
5. Update scripts/install.sh REPO variable (line 7)
6. Fix internal/version/version.go to reference 'dopa' instead of 'projectdb' (lines 35, 41, 92, 94, 130, 132, 143, 186, 286, 321)
7. Run go mod tidy
8. Search for any remaining 'example' or 'projectdb' references
9. Run make build to verify
10. Run make test to verify all tests pass
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Successfully updated all repository references and branding:

1. Updated go.mod module path from github.com/example/dopadone to github.com/marekbrze/dopadone
2. Updated all Go import statements across 52+ files using sed
3. Updated README.md repository URLs (3 occurrences)
4. Updated Makefile LDFLAGS to use correct module path
5. Updated scripts/install.sh REPO variable and URL comment
6. Fixed internal/version/version.go to use "dopa" instead of "projectdb" (8 references)
7. Updated documentation files (TUI.md, RELEASE.md, REBRANDING.md) to reflect final repository
8. Updated decision document to use cmd/dopa instead of cmd/projectdb
9. Ran go mod tidy successfully
10. All tests pass (make test)
11. Build succeeds (make build)
12. Binary correctly shows "dopa" when running ./bin/dopa version --all

Remaining references in backlog tasks and historical documentation are intentional for tracking purposes.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Updated repository references and branding from placeholder URLs to actual GitHub repository

## Changes

**Module and Imports:**
- Changed go.mod module path: `github.com/example/dopadone` → `github.com/marekbrze/dopadone`
- Updated all Go import statements across 52+ files
- Ran `go mod tidy` to update dependencies

**Documentation:**
- Updated README.md: 3 repository URL references
- Updated docs/TUI.md: 1 import example
- Updated docs/RELEASE.md: 5 GitHub URLs
- Updated docs/REBRANDING.md: Module path and example URLs
- Updated decision document: cmd/projectdb → cmd/dopa

**Build System:**
- Updated Makefile LDFLAGS to use correct module path for version injection
- Updated scripts/install.sh REPO variable and installation URL

**Branding Fixes:**
- Fixed internal/version/version.go: replaced all "projectdb" references with "dopa" (8 occurrences)
- Includes: BuildInfo output, GitHub API URL, asset names, binary names, temp directory names, error messages

## Verification

✅ All tests pass: `make test` (120+ tests)
✅ Build succeeds: `make build`
✅ Binary shows correct name: `./bin/dopa version --all` outputs "dopa" instead of "projectdb"
✅ No placeholder references remain in production code

## Notes

Remaining "example" references in backlog tasks and task descriptions are intentional for historical tracking. All production code and documentation now reference the actual repository: `github.com/marekbrze/dopadone`
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 All tests pass: make test
- [x] #2 Build succeeds: make build
- [x] #3 No 'example' or 'projectdb' references remain in production code
- [x] #4 Running ./bin/dopa version --all shows correct project name
<!-- DOD:END -->
