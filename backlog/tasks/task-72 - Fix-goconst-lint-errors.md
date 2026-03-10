---
id: TASK-72
title: Fix goconst lint errors
status: To Do
assignee: []
created_date: '2026-03-10 12:35'
labels:
  - lint
  - code-quality
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
String literals appear multiple times and should be extracted as constants to improve maintainability.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Create constant for 'enter' string in internal/tui/areamodal/area_modal.go (3 occurrences)
- [ ] #2 Create constant for 'esc' string in internal/tui/areamodal/area_modal.go (4 occurrences)
- [ ] #3 Use existing RootNodeName constant for 'root' string in internal/tui/tree/navigation.go
- [ ] #4 Create constant for 'windows' string in internal/version/version.go (3 occurrences)
- [ ] #5 Run make lint to verify all goconst errors are resolved
<!-- AC:END -->
