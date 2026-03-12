---
id: TASK-83
title: initial config when user installs
status: In Progress
assignee:
  - '@opencode'
created_date: '2026-03-11 17:08'
updated_date: '2026-03-12 06:31'
labels:
  - tui
  - onboarding
  - config
  - wizard
dependencies:
  - TASK-80
  - TASK-81
  - TASK-82
references:
  - cmd/dopa/main.go
  - internal/config/file.go
  - internal/db/driver/config.go
  - internal/tui/app.go
documentation:
  - docs/DATABASE_MODES.md
  - docs/TUI.md
  - .agents/skills/bubbletea/SKILL.md
  - .agents/skills/golang-patterns/SKILL.md
  - 'https://github.com/charmbracelet/bubbletea'
  - 'https://github.com/charmbracelet/lipgloss'
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When a user installs dopadone and runs it for the first time, they should be guided through an interactive configuration wizard that sets up their database connection. The wizard should present three options: (1) local SQLite database for single-device use, (2) Turso remote database for cloud-only access, or (3) Turso embedded replica for offline-capable cloud sync. Based on the user's choice, the wizard should collect necessary information, test the connection, run initial migrations, and persist the configuration. The goal is to provide a seamless 'set it and forget it' experience where users don't need to manually configure database paths or run migration commands. This task integrates with TASK-80 (default database path), TASK-81 (auto-migrations), and TASK-82 (initial area prompt) to create a complete first-run onboarding experience.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 TUI config wizard launches automatically on first run (no config file exists)
- [ ] #2 Wizard displays Dopadone branding with clear welcome message
- [ ] #3 User can choose between three database modes: Local, Turso Remote, Turso Replica
- [ ] #4 Each mode shows appropriate description and requirements
- [ ] #5 Local mode: displays default path (from TASK-80), allows custom path, creates directory automatically
- [ ] #6 Turso Remote mode: prompts for URL and auth token, validates connection before proceeding
- [ ] #7 Turso Replica mode: prompts for local path, URL, and token, validates connection before proceeding
- [ ] #8 Wizard tests database connection and runs migrations (TASK-81 integration) before saving config
- [ ] #9 On connection failure: shows clear error message and allows retry in same wizard session
- [ ] #10 On user cancellation: creates minimal local config as fallback so app is usable
- [ ] #11 Successful configuration is saved to standard location (dopadone.yaml in user config dir)
- [ ] #12 Config file has correct permissions (0600) for security
- [ ] #13 --skip-init flag allows advanced users to skip wizard and use defaults
- [ ] #14 Subsequent app runs detect existing config and skip wizard
- [ ] #15 Works correctly on all platforms (Linux, macOS, Windows)
- [ ] #16 Unit tests cover: first-run detection, wizard component, validation, retry logic
- [ ] #17 Integration tests cover: complete first-run flow for all three database modes
- [ ] #18 After config wizard completes, TUI starts normally; welcome modal appears automatically on first run (TASK-82 auto-trigger when no areas exist)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
PHASE 0: Prerequisites (before starting TASK-83)
1. Complete TASK-80 (database in user config directory) - establishes default database path
2. Understand TASK-81 (auto-run migrations) - will be integrated into config wizard
3. TASK-82 (initial area prompt) - COMPLETE, auto-triggers when TUI starts with empty database

PHASE 1: Configuration Wizard TUI Component
1. Create internal/tui/configwizard/ package with Bubble Tea component
2. Define wizard states: welcome -> choose_database_mode -> configure_connection -> verify_and_save
3. Implement step-by-step wizard with back/next navigation
4. Add progress indicator and help text
5. Support all three database modes: local, remote, replica
6. Write unit tests for wizard component

PHASE 2: First-Run Detection and Trigger
1. Create internal/config/first_run.go with IsFirstRun() function
2. Check for config file existence in standard locations
3. Modify cmd/dopa/main.go to check first run before any command
4. If first run, launch TUI config wizard before executing command
5. Support --skip-init flag for advanced users

PHASE 3: Configuration Wizard Flow
Step 1 - Welcome Screen:
- Dopadone branding and welcome message
- Explain the setup process
- Show three database mode options with descriptions

Step 2 - Mode Selection:
- Local SQLite (recommended for single device)
- Turso Remote (cloud-only, always online)
- Turso Replica (hybrid: local + cloud sync)

Step 3 - Mode-Specific Configuration:
Local mode:
- Show default path (from TASK-80: ~/.local/share/dopadone/dopadone.db)
- Allow custom path selection
- Auto-create directory

Turso Remote:
- Ask for Turso URL (libsql://your-db.turso.io)
- Ask for auth token
- Test connection before proceeding

Turso Replica:
- Show default local replica path
- Ask for Turso URL and token
- Test connection before proceeding

Step 4 - Verification and Setup:
- Test database connection
- Run migrations (integrate TASK-81)
- Show success/failure message
- Save configuration file

PHASE 4: Integration and Error Handling
1. Integrate auto-migration from TASK-81 into wizard
2. If connection fails, show error and allow retry
3. If user cancels, create minimal local config (fallback)
4. Save config to standard location (from TASK-80)
5. After config saves, TUI starts normally - TASK-82 welcome modal auto-triggers if no areas exist

PHASE 5: Configuration Persistence
1. Create dopadone.yaml in user config directory
2. Write selected mode and all configuration
3. Set file permissions (0600 for security)
4. Create default database directory if needed

PHASE 6: Testing
1. Unit tests for first-run detection
2. Unit tests for wizard component (each step)
3. Integration test: first run -> config -> migrations -> TUI starts -> welcome modal appears (auto)
4. Test all three database modes
5. Test cancellation and fallback behavior
6. Test retry on connection failure
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Design Decisions (from scoping session)

### Trigger Timing
- **When**: First run after install (when no config file exists)
- **Detection**: Check for dopadone.yaml in standard locations
- **Bypass**: --skip-init flag for advanced users

### Interface
- **Type**: TUI wizard using Bubble Tea framework
- **Style**: Interactive with arrow keys, forms, and visual feedback
- **Fallback**: None - TUI is required for first-run experience

### Validation Strategy
- **Connection testing**: Required before saving config
- **Migration execution**: Automatically run migrations (TASK-81 integration)
- **Error handling**: Clear error messages with retry option

### Failure Handling
- **Connection failure**: Show error, allow retry with different settings
- **User cancellation**: Create minimal local config as fallback
- **Goal**: App is always usable after first run, even with fallback config

## Dependencies on Other Tasks

### TASK-80: Store database in user config directory (MUST COMPLETE FIRST)
- Provides default database path: ~/.local/share/dopadone/dopadone.db
- Provides config file location: ~/.config/dopadone/config.yaml
- This task establishes the standard paths that TASK-83 will use

### TASK-81: Auto-run migrations on first app start (INTEGRATE)
- TASK-83 wizard will call migration logic as part of setup
- No need for separate auto-migration on first run - wizard handles it
- Use EnsureMigrations() function from TASK-81

### TASK-82: Prompt user for initial area (HANDS OFF TO)
- After TASK-83 wizard completes successfully, trigger TASK-82
- Seamless transition from database config to initial area creation
- User sees continuous onboarding flow

## Implementation Order

1. **BLOCKER**: Wait for TASK-80 to be completed
2. **PARALLEL**: Can start TUI component design while TASK-80 is in progress
3. **SEQUENTIAL**: Complete TASK-83 before starting TASK-82
4. **TESTING**: Integration test should cover complete onboarding flow

## Key Files to Create/Modify

### New Files
- internal/config/first_run.go - First-run detection logic
- internal/tui/configwizard/wizard.go - Main wizard component
- internal/tui/configwizard/welcome.go - Welcome screen
- internal/tui/configwizard/mode_selection.go - Database mode selection
- internal/tui/configwizard/local_config.go - Local mode configuration
- internal/tui/configwizard/turso_config.go - Turso mode configuration
- internal/tui/configwizard/verification.go - Connection testing and migrations
- internal/tui/configwizard/styles.go - Lipgloss styles for wizard

### Modified Files
- cmd/dopa/main.go - Add first-run check and wizard trigger
- internal/config/file.go - Add config saving with permissions
- docs/DATABASE_MODES.md - Document first-run wizard
- docs/START_HERE.md - Document onboarding flow

## Technical Considerations

### Bubble Tea Patterns (from skill)
- Follow 4 Golden Rules for layout (golden-rules.md)
- Use weight-based panel sizing
- Always truncate text explicitly
- Account for borders in height calculations

### Go Patterns (from skills)
- Accept interfaces, return structs
- Proper error wrapping with context
- Context for cancellation
- Table-driven tests for comprehensive coverage

### Security
- Config file permissions: 0600 (owner read/write only)
- Tokens stored in config file, not in environment
- Validate all user input before saving

## Updated Dependency on TASK-82 (2026-03-12)

TASK-82 is now COMPLETE. The welcome modal auto-triggers in handleAreasLoaded() when len(areas) == 0.

**No explicit hand-off code needed from TASK-83.** After config wizard saves the config file, the TUI starts normally. If the database has no areas, TASK-82's welcome modal appears automatically.

Integration flow:
1. TASK-83 config wizard completes → saves config
2. TUI starts (normal app launch)
3. AreasLoadedMsg arrives with empty list
4. TASK-82 handleAreasLoaded() detects empty → shows welcome modal
5. User creates first area via welcome modal

This means TASK-83 only needs to save the config file correctly - no special integration with TASK-82 required.
<!-- SECTION:NOTES:END -->

## Definition of Done
<!-- DOD:BEGIN -->
- [ ] #1 All unit tests pass (go test ./...)
- [ ] #2 No linting errors (make lint or golangci-lint run)
- [ ] #3 Manual testing: complete first-run flow for all three database modes
- [ ] #4 Manual testing: cancellation and fallback behavior
- [ ] #5 Manual testing: connection failure and retry behavior
- [ ] #6 Manual testing: subsequent runs skip wizard correctly
- [ ] #7 Update DATABASE_MODES.md with first-run wizard documentation
- [ ] #8 Update START_HERE.md with onboarding flow description
- [ ] #9 Code review: check Bubble Tea patterns follow golden rules
- [ ] #10 Code review: verify error handling is comprehensive
- [ ] #11 Cross-platform testing: verify paths work on Linux, macOS, Windows
- [ ] #12 Security review: verify config file permissions are correct
- [ ] #13 UX review: wizard flow is intuitive and helpful
- [ ] #14 Integration test with TASK-82: verify hand-off to initial area prompt works
<!-- DOD:END -->
