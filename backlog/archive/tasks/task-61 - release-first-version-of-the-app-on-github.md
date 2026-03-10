---
id: TASK-61
title: release first version of the app on github
status: To Do
assignee: []
created_date: '2026-03-07 21:40'
updated_date: '2026-03-07 21:50'
labels:
  - release
  - ci-cd
  - github-actions
  - documentation
dependencies: []
references:
  - backlog/tasks/task-62 - Update-Repository-References-and-Branding.md
  - backlog/tasks/task-63 - Prepare-CHANGELOG-for-v1.0.0-Release.md
  - backlog/tasks/task-64 - Update-and-Test-Installation-Script.md
  - backlog/tasks/task-65 - Create-GitHub-Actions-Release-Workflow.md
  - backlog/tasks/task-66 - Create-GitHub-Actions-CI-Workflow.md
  - backlog/tasks/task-67 - Execute-v1.0.0-Release.md
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Release the first stable version (v1.0.0) of Dopadone on GitHub with automated CI/CD pipeline for cross-platform binary distribution. This includes final testing, documentation updates, GitHub Actions setup, and creation of user-friendly installation methods.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Manual verification: Test binary builds locally using 'make build-all' and verify version command shows correct version info with 'make build-versioned VERSION=v1.0.0 && ./bin/dopa version --all'
- [ ] #2 Final testing and verification: Run 'make test' and 'make lint' with all tests passing and no linting errors
- [ ] #3 Update documentation: Replace all placeholder URLs (github.com/example/dopa) with actual repository URL (github.com/marekbrze/dopa) in README.md and any other documentation files
- [ ] #4 Prepare CHANGELOG: Ensure CHANGELOG.md has v1.0.0 section with all features, changes, and improvements documented following Keep a Changelog format
- [ ] #5 Setup GitHub Actions: Create .github/workflows/release.yml with automated build pipeline for cross-platform binaries (Linux amd64, macOS amd64/arm64, Windows amd64) triggered by version tags (v*)
- [ ] #6 Configure release workflow: Workflow must build all platform binaries, create GitHub release with auto-generated release notes, and upload binary archives
- [ ] #7 Create installation script: Create scripts/install.sh for easy one-line installation (curl | sh) that downloads correct binary for user's platform
- [ ] #8 Create version tag: Create annotated git tag v1.0.0 with appropriate release message
- [ ] #9 Push and verify: Push tag to GitHub, verify GitHub Actions workflow completes successfully, and confirm release appears at github.com/marekbrze/dopa/releases
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Implementation Plan for v1.0.0 Release

This task has been broken down into 6 subtasks organized in 2 parallel tracks:

## Track 1: Pre-Release Preparation (Sequential)

1. **Task-62**: Update Repository References and Branding (30-45 min)
   - Replace example/dopa with marekbrze/dopadone
   - Fix projectdb → dopa branding
   - Can run in parallel with Track 2

2. **Task-63**: Prepare CHANGELOG for v1.0.0 Release (20-30 min)
   - Move [Unreleased] → [1.0.0]
   - Add comparison links
   - Can run in parallel with Track 2

3. **Task-64**: Update and Test Installation Script (30-40 min)
   - Update repository URLs
   - Test platform detection
   - Depends on Task-62

## Track 2: CI/CD Infrastructure (Parallel with Track 1)

4. **Task-65**: Create GitHub Actions Release Workflow (60-90 min)
   - Build cross-platform binaries
   - Create releases automatically
   - Upload assets
   - Can run in parallel with Track 1

5. **Task-66**: Create GitHub Actions CI Workflow (30-40 min)
   - Run tests on PR/push
   - Optional but recommended
   - Can run in parallel with Track 1

## Track 3: Final Release (Sequential, after Tracks 1 & 2)

6. **Task-67**: Execute v1.0.0 Release (20-30 min)
   - Create and push v1.0.0 tag
   - Monitor release process
   - Verify installation
   - Depends on Tasks 62, 63, 64, 65

## Execution Strategy

**Phase 1: Preparation (Parallel)** - 1.5-2 hours
- Execute Tasks 62, 63, 65 simultaneously
- Then execute Task 64 (depends on 62)
- Task 66 can be done anytime

**Phase 2: Integration & Testing** - 30 minutes
- Merge all changes to main
- Run complete test suite
- Review all changes

**Phase 3: Release** - 30 minutes
- Execute Task 67
- Monitor and verify

**Total Estimated Time**: 2.5-3 hours

## Key Points

- Tasks 62, 63, 65, 66 can run in PARALLEL
- Task 64 must run AFTER Task 62
- Task 67 must run AFTER Tasks 62, 63, 64, 65
- Task 66 (CI workflow) is optional but recommended
<!-- SECTION:PLAN:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 Installation tested on at least one platform (Linux or macOS)
- [ ] #2 All acceptance criteria completed and verified
- [ ] #3 GitHub release v1.0.0 is publicly available and accessible at github.com/marekbrze/dopa/releases/tag/v1.0.0
- [ ] #4 Documentation reviewed and updated for v1.0.0 release
- [ ] #5 Installation tested on at least one platform (Linux or macOS)
- [ ] #6 Release includes binaries for all supported platforms (Linux, macOS Intel/ARM, Windows)
<!-- DOD:END -->
