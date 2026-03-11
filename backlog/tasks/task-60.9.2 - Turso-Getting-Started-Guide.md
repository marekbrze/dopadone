---
id: TASK-60.9.2
title: Turso Getting Started Guide
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-11 14:18'
updated_date: '2026-03-11 15:20'
labels:
  - documentation
  - turso
dependencies:
  - TASK-60.9.1
parent_task_id: TASK-60.9
priority: medium
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Create a comprehensive guide for setting up Turso account, creating databases, and obtaining credentials. This addresses AC#1 of TASK-60.9. Part of task-60.9 documentation effort.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Document Turso account signup process at turso.tech
- [x] #2 Document Turso CLI installation (turso cli)
- [x] #3 Document database creation via CLI and web UI
- [x] #4 Document how to generate authentication tokens
- [x] #5 Document how to find database URL
- [x] #6 Include screenshots or ASCII diagrams where helpful
- [x] #7 Link to official Turso documentation for advanced topics
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Research Turso documentation
   - Review official Turso docs at docs.turso.tech for account setup, CLI, and token generation
   - Verify all URLs and commands are current
   - Note any Turso platform changes

2. Create docs/TURSO_SETUP.md structure
   - Prerequisites section
   - Quick Start overview (5-minute path)
   - Detailed sections for each AC

3. Document AC#1: Account Signup (Section 1)
   - Navigate to turso.tech
   - Sign up options (GitHub, Google, email)
   - Account verification steps
   - Free tier limits and features
   - Add ASCII diagram of signup flow

4. Document AC#2: CLI Installation (Section 2)
   - macOS: Homebrew
   - Linux: curl script
   - Windows: PowerShell
   - Verify installation: turso --version
   - Add ASCII diagram of installation paths

5. Document AC#3: Database Creation (Section 3)
   - CLI method: turso db create
   - Web UI method (Turso dashboard)
   - Database naming best practices
   - Region selection
   - Add ASCII diagram of creation flow

6. Document AC#4: Token Generation (Section 4)
   - turso auth token create
   - Token scopes (read-only, full access)
   - Token expiration options
   - Security best practices
   - Add ASCII diagram of token flow

7. Document AC#5: Finding Database URL (Section 5)
   - turso db show <name> --url
   - Web UI: Database details page
   - URL format explanation
   - Add ASCII diagram of URL structure

8. AC#6: ASCII Diagrams
   - Create flow diagrams for signup, install, create
   - Add visual aids where helpful
   - Use consistent diagram style

9. AC#7: Link to Official Docs
   - Add "Learn More" section with links
   - Link to specific Turso docs pages
   - Add disclaimer about Turso changes

10. Create Quick Start Examples
    - Minimal setup (3 commands)
    - Full YAML config example
    - Common use case patterns

11. Add Dopadone Integration Section
    - How to configure Dopadone with Turso
    - dopadone.yaml example
    - Environment variable example
    - CLI flag example

12. Add Troubleshooting Section
    - Common signup issues
    - CLI installation problems
    - Token authentication errors

13. Cross-Reference Updates
    - Update DATABASE_MODES.md with link
    - Update START_HERE.md documentation index
    - Add link from TURSO_MIGRATIONS.md

14. Manual Verification
    - Walk through entire guide
    - Verify all commands work
    - Verify all links are valid
    - Test examples on clean setup

## File to Create
- docs/TURSO_SETUP.md (new)

## Files to Update
- docs/DATABASE_MODES.md (add cross-reference)
- docs/START_HERE.md (add to documentation index)

## Dependencies
- TASK-60.9.1: DONE ✅ (YAML config examples can be referenced)

## Sequential Steps (Must do in order)
1-2, 3-7, 8-9, 10-14

## Parallel Opportunities
None - Documentation is linear by nature

## Estimated Time: 3-4 hours
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Created comprehensive Turso Getting Started Guide (docs/TURSO_SETUP.md) covering:

**New file created:**
- docs/TURSO_SETUP.md - Complete guide with 7 acceptance criteria addressed

**Content includes:**
- Account signup process (CLI and Web)
- CLI installation for macOS/Linux/Windows
- Database creation via CLI and Web UI
- Authentication token generation with scopes and expiration
- Database URL retrieval methods
- ASCII diagrams for visual flows (signup, installation, token generation, URL structure)
- Links to official Turso documentation
- Dopadone-specific configuration examples (YAML, env vars, CLI flags)
- Troubleshooting section

**Cross-references updated:**
- docs/DATABASE_MODES.md - Added link to TURSO_SETUP.md
- docs/START_HERE.md - Added to documentation index
- docs/TURSO_MIGRATIONS.md - Added link to TURSO_SETUP.md

**Quality checks:**
- Build: PASS
- Tests: Pre-existing failure in internal/migrate (unrelated to docs)
- Lint: Pre-existing issues (unrelated to docs)
<!-- SECTION:FINAL_SUMMARY:END -->
