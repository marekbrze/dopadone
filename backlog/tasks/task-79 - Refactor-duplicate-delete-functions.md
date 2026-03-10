---
id: TASK-79
title: Refactor duplicate delete functions
status: To Do
assignee: []
created_date: '2026-03-10 19:13'
updated_date: '2026-03-10 19:13'
labels:
  - refactor
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
runProjectsDelete and runSubareasDelete have 42 lines of duplicate code. Extract common delete logic into shared helper.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Create shared delete helper function
- [ ] #2 Refactor runProjectsDelete to use helper
- [ ] #3 Refactor runSubareasDelete to use helper
- [ ] #4 golangci-lint dupl warning resolved
<!-- AC:END -->
