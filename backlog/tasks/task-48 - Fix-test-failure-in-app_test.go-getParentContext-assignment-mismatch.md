---
id: TASK-48
title: Fix test failure in app_test.go - getParentContext assignment mismatch
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 16:50'
updated_date: '2026-03-06 17:14'
labels: []
dependencies: []
references:
  - internal/tui/app_test.go
  - internal/tui/model.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Fix compilation error in internal/tui/app_test.go where getParentContext() returns 5 values (parentName, entityType, parentID, subareaID, showCheckbox) but tests only capture 4. This is blocking all TUI tests from running. The function signature was updated to return a 5th bool value for showCheckbox, but the tests at lines 492 and 512 were not updated. Need to update both test cases to capture all 5 values and verify the showCheckbox behavior is correct.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Tests in internal/tui/app_test.go compile without errors
- [x] #2 All tests in internal/tui package pass with go test ./internal/tui/...
- [x] #3 Both getParentContext calls (lines 492 and 512) capture all 5 return values
- [x] #4 Test assertions verify the showCheckbox bool value is correct
- [x] #5 No other test files have similar assignment mismatch issues
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Analyze the signature change in getParentContext() function
2. Understand the semantic meaning of showCheckbox bool return value
3. Update test case at line 492 to capture all 5 return values
4. Add meaningful assertion for showCheckbox value in first test case
5. Update test case at line 512 to capture all 5 return values
6. Add meaningful assertion for showCheckbox value in second test case
7. Verify test structure follows golang-testing patterns (table-driven if applicable)
8. Run go vet ./internal/tui/... for static analysis
9. Run go test ./internal/tui/... to verify compilation and execution
10. Run go test -race ./internal/tui/... to check for race conditions
11. Search codebase for similar assignment mismatches in other test files
12. Update any relevant documentation about getParentContext behavior
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Error message: internal/tui/app_test.go:492:50: assignment mismatch: 4 variables but model.getParentContext returns 5 values. Same error at line 512. This is blocking all TUI tests from running. The function signature was likely updated to return an additional value, but tests were not updated accordingly.

---
Root cause analysis:
The getParentContext() function signature was changed from returning 4 values to 5 values in commit aec6bdf (feat(tui): add subproject checkbox to quick-add modal). The function now returns a showCheckbox boolean to indicate whether the user should be shown a checkbox for creating a subproject.

The function behavior also changed - it no longer returns EntityTypeSubproject, instead it returns EntityTypeProject with showCheckbox=true to indicate that a checkbox should be shown for the subproject option.

Test updates:
1. Line 492 test (renamed from "returns_subproject_when_project_selected" to "returns_project_with_checkbox_when_project_selected"):
   - Now captures showCheckbox (5th value)
   - Expects EntityTypeProject (not EntityTypeSubproject)
   - Expects subareaID to be set (not nil)
   - Expects showCheckbox=true

2. Line 512 test (returns_project_when_no_project_selected):
   - Now captures showCheckbox (5th value)
   - Expects showCheckbox=false
   - Other assertions unchanged
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed test failures in internal/tui/app_test.go caused by getParentContext() signature change from 4 to 5 return values.\n\nChanges made:\n- Updated both test cases (lines 492 and 512) to capture all 5 return values (parentName, entityType, parentID, subareaID, showCheckbox)\n- Fixed test expectations to match new function behavior that returns EntityTypeProject instead of EntityTypeSubproject\n- First test now expects showCheckbox=true when project is selected with subarea context\n- Second test expects showCheckbox=false for regular projects without project selection\n- Added assertions for showCheckbox boolean value in both tests\n\nVerification:\n- All tests in internal/tui package compile and pass\n- No assignment mismatches found in other test files\n- go vet ./internal/tui/... passes\n- go fmt ./internal/tui/... passes\n- Race detector shows no issues (go test -race)
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Run go vet ./internal/tui/...
- [x] #2 Run go fmt ./internal/tui/...
- [x] #3 Code review: verify test assertions are meaningful
- [x] #4 No race conditions detected in tests (go test -race)
<!-- DOD:END -->
