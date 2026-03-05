---
id: TASK-30
title: 'Replace interface{} types with proper nullable types in DB models'
status: To Do
assignee: []
created_date: '2026-03-04 16:59'
labels:
  - architecture
  - refactoring
  - db
dependencies: []
priority: low
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
DB models use interface{} for nullable fields (DeletedAt, Deadline) instead of proper types like *time.Time. This reduces type safety and can cause runtime errors. Refactor to use proper nullable types.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Replace DeletedAt interface{} with *time.Time
- [ ] #2 Replace Deadline interface{} with *time.Time
- [ ] #3 Update sqlc.yaml to generate proper nullable types
- [ ] #4 Update all usages of these fields across codebase
- [ ] #5 All tests pass after type changes
<!-- AC:END -->
