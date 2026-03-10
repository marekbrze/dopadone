---
id: TASK-75
title: Handle unchecked error returns (32 instances)
status: To Do
assignee: []
created_date: '2026-03-10 12:37'
labels:
  - lint
  - code-quality
  - errcheck
dependencies: []
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Error return values are not checked in multiple files across cmd/dopa and internal packages. Need to handle errors from Close(), Flush(), MarkFlagRequired(), Fprintf(), RemoveAll() calls.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Fix unchecked services.Close() errors in cmd/dopa/areas.go (4 instances)
- [ ] #2 Fix unchecked formatter.Flush() errors in cmd/dopa/ (3 instances)
- [ ] #3 Fix unchecked MarkFlagRequired() errors in cmd/dopa/ (5 instances)
- [ ] #4 Fix unchecked fmt.Fprintf() errors in cmd/dopa/tui.go (2 instances)
- [ ] #5 Fix unchecked db.Close() and os.RemoveAll() errors in internal/db/ tests (10 instances)
- [ ] #6 Fix unchecked rows.Close() errors in internal/db/ tests (6 instances)
- [ ] #7 Fix unchecked database.Close() errors in internal/tui/ tests (3 instances)
- [ ] #8 Run make lint to verify all 32 errcheck errors are resolved
<!-- AC:END -->
