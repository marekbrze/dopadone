---
id: TASK-64
title: Update and Test Installation Script
status: To Do
assignee: []
created_date: '2026-03-07 21:49'
labels:
  - release
  - scripts
  - testing
dependencies:
  - TASK-62
references:
  - backlog/tasks/task-61 - release-first-version-of-the-app-on-github.md
  - backlog/tasks/task-62 - Update-Repository-References-and-Branding.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update the installation script with correct repository information and test it locally to ensure it works for the v1.0.0 release. This task depends on task-62 (repository URL updates) and is part of the v1.0.0 release preparation (task-61).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Update REPO variable in scripts/install.sh to 'marekbrze/dopadone'
- [ ] #2 Update BINARY_NAME to match actual binary name if needed
- [ ] #3 Test script dry-run for platform detection (Linux, macOS Intel, macOS ARM)
- [ ] #4 Add error handling for missing dependencies (curl, tar, unzip)
- [ ] #5 Add verification step after installation
- [ ] #6 Test archive extraction logic matches GitHub Actions output format
<!-- AC:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Script successfully detects platform on Linux, macOS Intel, and macOS ARM
- [ ] #2 Script can parse GitHub API release response
- [ ] #3 Installation verification works: dopa version
- [ ] #4 Script tested locally with platform detection and extraction logic
<!-- DOD:END -->
