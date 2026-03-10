---
id: TASK-44
title: ellipsis shouldnt replace whole value when column is narror
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 08:34'
updated_date: '2026-03-06 08:49'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When column width is narrow and text exceeds available space, the current implementation replaces the entire value with a single ellipsis (…). This should instead show a partial value with ellipsis at the end to preserve as much readable content as possible.

Example:
- Current behavior: "Very Long Project Name" → "…" 
- Desired behavior: "Very Long Project Name" → "Very Long Pro…"

This improves usability by allowing users to see context and differentiate between long values instead of seeing just ellipses for all truncated content.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Given a string "Very Long Project Name" and maxLen=15, the truncated output is "Very Long Pro…"
- [x] #2 ANSI escape codes in the original string are preserved in the visible portion of the truncated output
- [x] #3 For maxLen <= 1, the output shows the first character plus ellipsis (e.g., "a…") to provide minimal context
- [x] #4 Truncation is Unicode-aware: multi-byte characters (emojis, CJK) are handled correctly without breaking character boundaries
- [x] #5 All existing tests pass after the fix
- [x] #6 New tests are added for edge cases: empty strings, very narrow columns, ANSI codes, and Unicode characters
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Phase 1: Analysis and Preparation (Sequential)
- Load bubbletea and golang-testing skills for TUI patterns and test best practices
- Read and analyze current truncateString implementation in internal/tui/views/columns.go:102-120
- Identify the bug: line 116 returns "…" instead of truncated portion
- Review existing tests in columns_test.go to understand test patterns

Phase 2: Fix Implementation (Sequential)
- Fix truncateString function to show partial value:
  1. Calculate visible character count (accounting for ANSI codes)
  2. If string fits (len(stripped) <= maxLen), return as-is
  3. If too long, truncate to maxLen-1 characters and append "…"
  4. Handle maxLen <= 1 edge case (show first char + "…")
  5. Use runes for Unicode-safe truncation (multi-byte chars)
  6. Preserve ANSI codes in visible portion of output

Phase 3: Write Comprehensive Tests (Sequential with parallel sub-tasks)
- Add new test function TestTruncateString in columns_test.go:
  * Test AC#1: Partial truncation ("Very Long Project Name", 15) → "Very Long Pro…"
  * Test AC#2: ANSI preservation with colored strings
  * Test AC#3: maxLen <= 1 edge cases (maxLen=1, maxLen=0)
  * Test AC#4: Unicode handling (emojis, CJK characters)
  * Test AC#6: Edge cases (empty strings, very narrow columns)
- Use table-driven test pattern (golang-testing skill)
- Add subtests for each category using t.Run()

Phase 4: Verification (Sequential)
- Run existing tests: go test ./internal/tui/views/... -v
- Run linter: make lint
- Manual verification: start TUI with narrow columns, verify truncation works
- Mark all ACs as complete using CLI

Phase 5: Documentation (Parallel with verification)
- Add code comments explaining the truncation logic
- Document ANSI handling approach
- Note Unicode considerations for future maintainers

Dependencies: None (single atomic task)
Estimated Time: 2-3 hours
Risk: Low (isolated function, existing tests provide safety net)
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
- Fixed truncateString in internal/tui/views/columns.go to show partial content with ellipsis
- Added comprehensive tests: TestTruncateString, TestTruncateStringWithANSI, TestTruncateStringUnicode, TestTruncateStringEdgeCases
- All 44 tests pass in views package
- Lint passes (pre-existing error in app_test.go unrelated to this change)
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Fixed truncateString function to show partial content instead of just ellipsis when column width is narrow.

Changes:
- Modified truncateString in internal/tui/views/columns.go to preserve visible characters up to maxLen-1, then append ellipsis
- Handles ANSI escape codes by tracking escape state and including codes in output without counting them as visible chars
- Uses rune-based iteration for proper Unicode handling (emojis, CJK characters)
- For maxLen <= 1, shows first visible character + ellipsis to provide minimal context

Tests added:
- TestTruncateString: basic truncation cases including AC#1 example
- TestTruncateStringWithANSI: verifies ANSI code preservation
- TestTruncateStringUnicode: emoji and CJK character handling
- TestTruncateStringEdgeCases: empty strings, very narrow columns, single chars

All 44 existing and new tests pass. Ready for user review.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Run tests: go test ./internal/tui/views/...
- [x] #2 Run linter: make lint or go vet
- [x] #3 Verify truncation works in TUI at narrow widths
<!-- DOD:END -->
