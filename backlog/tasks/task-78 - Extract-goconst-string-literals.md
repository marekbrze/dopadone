---
id: TASK-78
title: Extract goconst string literals
status: To Do
assignee: []
created_date: '2026-03-10 19:13'
updated_date: '2026-03-10 19:13'
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
- [ ] #1 Create constant for "enter" key string
- [ ] #2 Create constant for "esc" key string
- [ ] #3 Create constant for "root" node name
- [ ] #4 Create constant for "windows" OS string
- [ ] #5 golangci-lint reports 0 goconst issues
<!-- AC:END -->
