---
id: TASK-65
title: Create GitHub Actions Release Workflow
status: Done
assignee:
  - '@ai-assistant'
created_date: '2026-03-07 21:50'
updated_date: '2026-03-08 16:27'
labels:
  - release
  - ci-cd
  - github-actions
dependencies: []
references:
  - backlog/tasks/task-61 - release-first-version-of-the-app-on-github.md
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a GitHub Actions workflow that automatically builds cross-platform binaries, creates GitHub releases, and uploads assets when a version tag is pushed. This is the core CI/CD infrastructure for task-61 (v1.0.0 release).

This task can be executed in parallel with tasks 62-64.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Create .github/workflows/release.yml with workflow configuration
- [x] #2 Configure workflow to trigger on tags matching 'v*' pattern
- [x] #3 Implement build jobs for all platforms: Linux amd64, macOS amd64, macOS arm64, Windows amd64
- [x] #4 Inject version information via ldflags during build (Version, GitCommit, BuildDate)
- [x] #5 Create distribution archives (.tar.gz for Unix, .zip for Windows)
- [x] #6 Create GitHub release with auto-generated release notes
- [x] #7 Upload all binary archives as release assets
- [x] #8 Calculate and upload SHA256 checksums for all binaries
- [x] #9 Mark pre-releases for tags with hyphen (e.g., v1.0.0-beta.1)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create .github/workflows directory structure
2. Create release.yml workflow file with:
   - Trigger on v* tags and workflow_dispatch for testing
   - Setup Go environment (1.21+)
   - Build matrix for all platforms (linux-amd64, darwin-amd64, darwin-arm64, windows-amd64)
   - Inject version info via ldflags (Version, GitCommit, BuildDate)
   - Create distribution archives (.tar.gz for Unix, .zip for Windows)
3. Add release creation step:
   - Generate release notes automatically
   - Mark pre-releases for tags with hyphen (e.g., v1.0.0-beta.1)
   - Upload all binary archives as assets
   - Calculate and upload SHA256 checksums
4. Test workflow configuration:
   - Validate YAML syntax
   - Ensure naming matches expected format: dopa-{os}-{arch}.{ext}
   - Verify all acceptance criteria are met
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Created comprehensive GitHub Actions release workflow with the following features:

- Automated builds for all target platforms (Linux amd64, macOS amd64/arm64, Windows amd64)
- Version injection via ldflags matching existing Makefile configuration
- Platform-specific archive creation (.tar.gz for Unix, .zip for Windows)
- SHA256 checksum generation for all releases
- Automatic release notes generation
- Pre-release detection based on hyphen in tag name
- Manual testing capability via workflow_dispatch trigger

The workflow follows GitHub Actions best practices and integrates seamlessly with the existing build infrastructure.
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary

Created `.github/workflows/release.yml` providing automated cross-platform binary releases for the Dopadone project.

## Changes

- **Workflow Triggers**: Push on `v*` tags and manual `workflow_dispatch` for testing
- **Build Matrix**: Configured parallel builds for Linux amd64, macOS (amd64/arm64), Windows amd64
- **Version Injection**: Implemented ldflags injection matching Makefile pattern (Version, GitCommit, BuildDate)
- **Archive Creation**: Platform-specific archives (.tar.gz for Unix, .zip for Windows)
- **Checksums**: SHA256 generation for all releases with platform-specific commands
- **Release Management**: Auto-generated release notes, pre-release detection via hyphen in tag
- **Asset Upload**: All binaries, archives, and checksums uploaded as release assets

## Technical Details

- Go version: 1.21+ (as required)
- Archive naming: `dopa-{os}-{arch}.{ext}` (matches specification)
- Module path: `github.com/marekbrze/dopadone`
- Build flags: `-trimpath` with `-ldflags "-s -w"` for optimized binaries

## Testing

Workflow can be tested via:
1. Manual trigger: GitHub Actions → Release workflow → Run workflow
2. Tag push: `git tag v0.0.1-test && git push origin v0.0.1-test`

## Dependencies

Ready for task-61 (v1.0.0 release) - no blockers.
<!-- SECTION:FINAL_SUMMARY:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [x] #1 Workflow file is valid YAML and passes validation
- [x] #2 Workflow can be tested with workflow_dispatch trigger
- [x] #3 All build steps use correct Go version (1.21+)
- [x] #4 Archives match expected naming: dopa-{os}-{arch}.{ext}
- [x] #5 Release notes are properly formatted
<!-- DOD:END -->
