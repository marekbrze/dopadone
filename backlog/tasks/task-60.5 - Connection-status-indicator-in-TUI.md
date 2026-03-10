---
id: TASK-60.5
title: Connection status indicator in TUI
status: To Do
assignee: []
created_date: '2026-03-08 19:01'
updated_date: '2026-03-10 07:30'
labels:
  - tui
  - database
  - ui
milestone: m-1
dependencies: []
parent_task_id: TASK-60
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add database connection status indicator to TUI showing connected/syncing/offline/local-only states. Part of task-60 Turso integration.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Add connection status field to TUI model
- [ ] #2 Create status indicator component in status bar or title bar
- [ ] #3 Show visual indicator: ● (green) connected, ◐ (yellow) syncing, ○ (red) offline, ■ (gray) local-only
- [ ] #4 Update status on connection state changes
- [ ] #5 Add status message on hover/help
- [ ] #6 Support status updates from sync goroutine
<!-- AC:END -->
