---
id: TASK-32
title: Optimize data filtering with server-side queries
status: To Do
assignee: []
created_date: '2026-03-04 17:00'
labels:
  - performance
  - optimization
  - db
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Current implementation loads all data into memory and filters in Go (e.g., LoadProjectsCmd loads ALL projects). Optimize by moving filtering logic to SQL queries using sqlc for better performance with large datasets.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Identify all in-memory filtering operations
- [ ] #2 Create filtered SQL queries in queries/*.sql files
- [ ] #3 Update sqlc to generate filtered query functions
- [ ] #4 Update TUI commands to use filtered queries
- [ ] #5 Add tests for filtered queries
- [ ] #6 Performance benchmark shows improvement for large datasets
<!-- AC:END -->
