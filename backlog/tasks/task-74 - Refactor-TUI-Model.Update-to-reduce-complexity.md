---
id: TASK-74
title: Refactor TUI Model.Update to reduce complexity
status: To Do
assignee: []
created_date: '2026-03-10 12:37'
labels:
  - lint
  - code-quality
  - refactor
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
The (Model).Update function in internal/tui/app.go has cyclomatic complexity of 94 (limit is 30). Need to refactor by extracting handlers for different message types.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Analyze current Update method structure and message types
- [ ] #2 Extract handler functions for distinct message types (e.g., keyMsgHandler, mouseMsgHandler, etc.)
- [ ] #3 Reduce complexity to below 30
- [ ] #4 Ensure all tests pass after refactoring
- [ ] #5 Run make lint to verify gocyclo error is resolved
<!-- AC:END -->
