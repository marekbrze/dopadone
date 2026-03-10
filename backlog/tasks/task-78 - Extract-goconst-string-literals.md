---
id: TASK-78
title: Extract goconst string literals
status: Done
assignee:
  - '@claude'
created_date: '2026-03-10 19:13'
updated_date: '2026-03-10 20:47'
labels:
  - lint
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Extract repeated strings to constants: enter, esc, root, windows
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create constant for "enter" key string
- [x] #2 Create constant for "esc" key string
- [x] #3 Create constant for "root" node name
- [x] #4 Create constant for "windows" OS string
- [x] #5 golangci-lint reports 0 goconst issues
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
## Implementation Plan for TASK-78: Extract goconst string literals

### Task Assessment
**Complexity:** LOW - Simple refactoring task
**Estimated Effort:** 1-2 hours
**Splitting Required:** NO - Single cohesive task

This is a straightforward code quality improvement. All acceptance criteria relate to the same goal: fix goconst lint warnings.

### Analysis

**Current State:**
- Constants KeyEnter, KeyEsc, RootNodeName already exist but are NOT used consistently
- "windows" constant does not exist yet
- goconst config: min-len=3, min-occurrences=3, excluded in test files

**Files Requiring Updates:**

| String | Production Files | Test Files (for consistency) |
|--------|------------------|------------------------------|
| "enter" | areamodal/area_modal.go (3), modal/modal.go (1) | None |
| "esc" | areamodal, modal, confirmmodal, help, spacemenu | confirmmodal_test, help_test, spacemenu_test |
| "root" | tree/builder.go, tree/navigation.go, tree/renderer.go | Multiple tree test files |
| "windows" | version/version.go (3) | None |

### Implementation Phases

#### Phase 1: Create Missing Constant (5 min) - SEQUENTIAL FIRST
1. Add OSWindows constant to internal/tui/constants.go:
   - const OSWindows = "windows"
   - Optionally add OSLinux and OSDarwin for consistency
2. Tests: None needed - constants are simple string values

#### Phase 2: Update Production Code (30 min) - PARALLEL with Phase 3

**2a. Update KeyEnter usage**
- Files: internal/tui/areamodal/area_modal.go, internal/tui/modal/modal.go
- Replace: case "enter": → case tui.KeyEnter:

**2b. Update KeyEsc usage**
- Files: areamodal, modal, confirmmodal, help, spacemenu
- Replace: "esc" → tui.KeyEsc

**2c. Update RootNodeName usage**
- Files: internal/tui/tree/builder.go, internal/tui/tree/navigation.go, internal/tui/tree/renderer.go
- Replace: "root" → RootNodeName (handle both string comparisons and constructor calls)

**2d. Update OSWindows usage**
- Files: internal/version/version.go
- Replace: "windows" → tui.OSWindows

Tests: Run go test ./internal/tui/... and go test ./internal/version/... - no new tests needed

#### Phase 3: Update Test Files (15 min) - PARALLEL with Phase 2

**3a. Update KeyEsc in tests**
- Files: confirmmodal/modal_test.go, help/help_test.go, spacemenu/spacemenu_test.go
- Replace hardcoded "esc" with tui.KeyEsc

**3b. Update RootNodeName in tests**
- Files: Multiple tree test files (node_test.go, builder_test.go, renderer_test.go, navigation_test.go)
- Replace hardcoded "root" with tree.RootNodeName
- Note: Tests excluded from goconst but updating ensures consistency

Tests: Run go test ./... - all tests should pass

#### Phase 4: Verification (10 min) - SEQUENTIAL AFTER Phase 2 & 3

1. Run golangci-lint: golangci-lint run ./... (expected: 0 goconst issues)
2. Run full test suite: go test -race -cover ./... (expected: all pass)
3. Check for regressions: go build ./cmd/dopa (expected: successful build)

### Test Strategy

**No New Tests Required**
- This is a refactoring task with no behavior changes
- Existing tests verify functionality
- Linter validates the change

**Regression Prevention:**
- golangci-lint will fail if string literals reappear
- goconst config ensures strings with 3+ occurrences are flagged

### Documentation Updates

Check and update if these docs exist:
1. docs/CODE_QUALITY.md - Add section on string literal constants
2. docs/TUI.md - Document available key constants if relevant section exists

### Dependencies & Parallelization

Dependency graph:
Phase 1 (Create Constants)
    ↓
    ├─→ Phase 2 (Update Production) ──→ Phase 4 (Verify Lint)
    └─→ Phase 3 (Update Tests)     ──→ Phase 4 (Verify Lint)

**Can run in parallel:** Phase 2 and Phase 3 (after Phase 1 completes)
**Must run sequentially:** Phase 1 → Phase 2/3 → Phase 4

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Import cycle | Low | Medium | Constants in separate packages (tui, tree) |
| Missed occurrences | Low | Low | goconst will catch them |
| Test failures | Very Low | Low | No behavior changes, only string replacements |

### Acceptance Criteria Checklist

- [ ] AC #1: Create constant for "enter" key string → tui.KeyEnter already exists, just need to USE it
- [ ] AC #2: Create constant for "esc" key string → tui.KeyEsc already exists, just need to USE it
- [ ] AC #3: Create constant for "root" node name → tree.RootNodeName already exists, just need to USE it
- [ ] AC #4: Create constant for "windows" OS string → CREATE tui.OSWindows
- [ ] AC #5: golangci-lint reports 0 goconst issues → VERIFY with golangci-lint run
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Extracted repeated string literals to constants to eliminate goconst lint warnings.

## Changes

### Key Constants Added/Fixed
- Added `KeyQ`, `KeyCtrlC`, `KeyEnter` to `internal/tui/keyboard_handler.go` local constants

### Existing Constants Verified
- `KeyEnter`, `KeyEsc` in `internal/tui/internal/constants/keys.go`
- `RootNodeName` in `internal/tui/tree/constants.go`
- `OSWindows` in `internal/constants/os.go`

### Documentation
- Updated `docs/CODE_QUALITY.md` with available string constants reference

## Verification
- All tests pass: `go test ./...`
- No goconst issues: `golangci-lint run ./... | grep goconst` returns 0 issues
<!-- SECTION:FINAL_SUMMARY:END -->
