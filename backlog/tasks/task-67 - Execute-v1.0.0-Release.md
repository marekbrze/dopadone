---
id: TASK-67
title: Execute v1.0.0 Release
status: In Progress
assignee:
  - '@{myself}'
created_date: '2026-03-07 21:50'
updated_date: '2026-03-09 21:05'
labels:
  - release
  - deployment
  - verification
dependencies:
  - TASK-62
  - TASK-63
  - TASK-64
  - TASK-65
references:
  - backlog/tasks/task-61 - release-first-version-of-the-app-on-github.md
  - backlog/tasks/task-62 - Update-Repository-References-and-Branding.md
  - backlog/tasks/task-63 - Prepare-CHANGELOG-for-v1.0.0-Release.md
  - backlog/tasks/task-64 - Update-and-Test-Installation-Script.md
  - backlog/tasks/task-65 - Create-GitHub-Actions-Release-Workflow.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Execute the actual v1.0.0 release by creating and pushing the version tag, monitoring the automated release process, and verifying the release is successful. This is the final step of task-61 and MUST be executed after all preparation tasks (62-65) are complete.

This task MUST be executed sequentially after tasks 62, 63, 64, and 65 are complete.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Verify all previous tasks (62-65) are complete and merged to main branch
- [ ] #2 Run final verification: make test passes, make lint passes, make build-all succeeds
- [ ] #3 Verify version injection: VERSION=v1.0.0 make build-versioned && ./bin/dopa version --all shows v1.0.0
- [ ] #4 Create annotated git tag: git tag -a v1.0.0 -m 'Release v1.0.0: First stable release'
- [ ] #5 Push tag to trigger release: git push origin v1.0.0
- [ ] #6 Monitor GitHub Actions workflow completion without errors
- [ ] #7 Verify release appears at https://github.com/marekbrze/dopadone/releases/tag/v1.0.0
- [ ] #8 Verify all platform binaries (Linux, macOS Intel/ARM, Windows) are attached to release
- [ ] #9 Test installation script: curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | sh
- [ ] #10 Test downloaded binary shows correct version: dopa version --all
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Fix test build issues (mock queriers)\n2. Verify build and lint pass\n3. Update CHANGELOG with actual release date\n4. Commit all pending changes\n5. Create v1.0.0 git tag\n6. Push tag to trigger release\n7. Monitor GitHub Actions workflow\n8. Verify release on GitHub\n9. Test installation script
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Release is publicly accessible on GitHub
- [ ] #2 All binary archives are downloadable
- [ ] #3 Installation script successfully installs the binary
- [ ] #4 Installed binary shows correct version: v1.0.0
- [ ] #5 GitHub release notes are properly formatted
- [ ] #6 No errors or warnings in GitHub Actions logs
<!-- DOD:END -->
