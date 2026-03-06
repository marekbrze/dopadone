---
id: TASK-6
title: Advanced CLI Features - Filtering and Output Polish
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 09:37'
updated_date: '2026-03-03 10:29'
labels:
  - cli
  - testing
  - advanced
dependencies:
  - TASK-8
  - TASK-9
  - TASK-5
  - TASK-7
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Add advanced query filtering, YAML output format, and comprehensive unit tests for all CLI commands.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 All list commands support --filter flag with query syntax (e.g., --filter 'status=active AND priority>=high')
- [x] #2 Add YAML output format: --format=yaml option
- [x] #3 Unit tests for command flag parsing across all entities
- [x] #4 Unit tests for validation logic
- [x] #5 Unit tests for output formatting (table, JSON, YAML)
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
1. Create internal/cli/filter package with parser.go (AST + lexer + parser)\n2. Implement evaluator.go for filter evaluation\n3. Add YAML formatter to output package (yaml.go)\n4. Add --filter and --format=yaml flags to all list commands\n5. Create comprehensive tests for filter parser\n6. Create tests for YAML and all output formats\n7. Create tests for command flag parsing\n8. Run all tests to verify
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Implemented advanced CLI features with query filtering and YAML output support.

Changes:
- Created internal/cli/filter package with lexer, parser, and AST for query expressions
- Implemented filter evaluator supporting: =, !=, <, <=, >, >=, AND, OR, parentheses
- Added YAML output formatter using gopkg.in/yaml.v3
- Updated all list commands (areas, subareas, projects) with --filter and --format=yaml flags
- Created comprehensive tests:
  - filter/parser_test.go: 50+ tests for lexer, parser, evaluator
  - output/formatter_test.go: Tests for table, JSON, YAML formatters
  - cmd/dopa/commands_test.go: Tests for flag parsing across all commands

Files created/modified:
- internal/cli/filter/parser.go (new)
- internal/cli/filter/evaluator.go (new)
- internal/cli/filter/parser_test.go (new)
- internal/cli/output/yaml.go (new)
- internal/cli/output/formatter.go (updated)
- internal/cli/output/formatter_test.go (updated)
- cmd/dopa/areas.go (updated)
- cmd/dopa/subareas.go (updated)
- cmd/dopa/projects.go (updated)
- cmd/dopa/commands_test.go (new)

All tests passing: go test ./...
<!-- SECTION:FINAL_SUMMARY:END -->
