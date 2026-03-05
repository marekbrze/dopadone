---
id: TASK-12
title: prepare a distribution flow
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 11:12'
updated_date: '2026-03-03 11:36'
labels: []
dependencies: []
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement a distribution and upgrade flow for the Go CLI/TUI application. Users install via GitHub Releases and manually trigger updates. CI/CD should handle semantic versioning and release automation.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Document release process: how to publish new versions (CI/CD workflow)
- [x] #2 Document user upgrade instructions: how users check for and install updates
- [x] #3 Implement Go build pipeline for cross-platform binaries (windows/amd64, darwin/amd64, darwin/arm64, linux/amd64)
- [x] #4 Create GitHub Actions workflow for automated releases on version tags
- [x] #5 Implement 'version' command showing current version and checking for latest release
- [x] #6 Add update command/instructions for manual upgrade process
- [x] #7 Generate release notes automatically from commits or changelog
- [x] #8 Create AGENTS.md guidelines for implementing release workflow (version bumping, tagging, CI triggers)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create internal/version package with ldflags-injectable version info
2. Update rootCmd.Version and version command to use dynamic version + GitHub API check
3. Create .github/workflows/release.yml for cross-platform builds on tag push
4. Create scripts/generate-changelog.sh for auto-generating release notes
5. Update Makefile with cross-compile targets (windows/amd64, darwin/amd64, darwin/arm64, linux/amd64)
6. Create docs/RELEASE.md documenting release process and user upgrade instructions
7. Update AGENTS.md with release workflow guidelines (version bumping, tagging, CI triggers)
8. Test the workflow by creating a test release
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented complete distribution and upgrade flow for projectdb CLI.

Changes:
- Created internal/version package with ldflags-injectable version info (Version, GitCommit, BuildDate)
- Updated version command to show detailed build info with --all flag
- Added new update command that checks GitHub for latest release and shows upgrade instructions
- Created .github/workflows/release.yml for automated cross-platform releases on version tags
- Added scripts/generate-changelog.sh for auto-generating release notes from git commits
- Added scripts/install.sh for quick installation via curl
- Updated Makefile with cross-compile targets (linux/amd64, darwin/amd64, darwin/arm64, windows/amd64)
- Created docs/RELEASE.md documenting release process and user upgrade instructions
- Updated AGENTS.md with release workflow guidelines (version bumping, tagging, CI triggers)

Supported platforms:
- Linux (amd64)
- macOS (amd64, arm64/M1/M2)
- Windows (amd64)

Commands added:
- projectdb version [--all] - Show version info
- projectdb update - Check for updates and show upgrade instructions

Tests run: go test ./... (all passing)
Build verification: make build-all (all 4 platforms)

Updated: Added embedded migrations and fully automated upgrade command.
- Migrations are now embedded in the binary using go:embed
- Added 'migrate' command with subcommands: up, down, status, reset
- 'upgrade' command now downloads, replaces binary, AND runs migrations automatically
- Users just run: projectdb upgrade
<!-- SECTION:FINAL_SUMMARY:END -->
