---
id: TASK-20
title: 'TUI 16A: Tree Rendering Package'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-03 13:48'
updated_date: '2026-03-03 15:19'
labels:
  - tui
  - mvp
  - phase2
dependencies: []
references:
  - internal/domain/project.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Build tree package for hierarchical project display with unlimited nesting. Pure logic package with no database dependencies. Includes lipgloss styling, expand/collapse state management, and navigation helpers for small trees (<100 nodes).
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 Projects display in hierarchical tree structure with visual indicators (├─, └─, │) and lipgloss styling
- [x] #2 Nested projects support expand/collapse behavior with +/- indicators, state managed in TreeNode.IsExpanded
- [x] #3 Unlimited nesting depth via recursive tree building using domain.Project.ParentID
- [x] #4 Tree package includes: node model, tree builder, renderer with lipgloss, and navigation helpers
- [x] #5 Unit tests for tree building with scenarios: empty, flat, 1-level nested, 5+ levels deep, orphans, position ordering
- [x] #6 Unit tests for tree rendering: indicators, collapsed nodes, selected node highlighting
- [x] #7 Unit tests for navigation: GetNextVisibleNode, GetPrevVisibleNode, skip collapsed children, boundary cases
- [x] #8 Performance validated for trees up to 100 nodes (no optimization needed for small trees)
- [x] #9 Package has >90% test coverage with fast, independent, repeatable tests following F.I.R.S.T. principles
- [x] #10 All exported types and functions have godoc comments explaining purpose and usage
- [x] #11 No magic numbers/strings: use constants for tree characters (├─, └─, │, +, -) and indentation
- [x] #12 Source code dependencies point inward: tree package depends only on domain types, lipgloss, and stdlib
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
Phase 1: Package Structure & Node Model (1 hour)
- Create internal/tui/tree/ with clean architecture boundaries
- node.go: TreeNode struct with ID, Name, Depth, IsExpanded, Children, Data (interface{})
- Methods: IsLeaf(), HasChildren(), ToggleExpanded()
- Clean code: meaningful names, small methods, no magic numbers
- Tests: TestTreeNodeBasics, TestTreeNodeToggle

Phase 2: Tree Builder (1.5 hours)
- builder.go: BuildFromProjects([]domain.Project) *TreeNode
- Algorithm: separate roots (SubareaID != nil) from children, sort by Position, recursive attach
- Edge cases: empty list → nil, orphans → log warning & skip, circular refs → detect & prevent
- Clean architecture: pure logic, no DB dependencies, depends only on domain types
- TDD approach: write tests first for all scenarios (empty, flat, nested, deep, orphans, position)
- Tests: TestBuildEmpty, TestBuildFlat, TestBuildNested, TestBuildDeep, TestBuildOrphans

Phase 3: Tree Renderer with Lipgloss (1.5 hours)
- renderer.go: Render(root *TreeNode, selectedID string) string
- Tree characters: ├─ └─ │ with constants (no magic strings)
- Indentation: 2 spaces per depth level
- Expand/collapse: [+]/[-] indicators based on IsExpanded
- Selected node: lipgloss styling for highlighting
- Skip children of collapsed nodes
- Clean code: small helper functions, self-documenting names
- Tests: TestRenderEmpty, TestRenderSingleNode, TestRenderMultiLevel, TestRenderCollapsed, TestRenderSelected

Phase 4: Navigation Helpers (45 min)
- navigation.go: GetNextVisibleNode, GetPrevVisibleNode, GetAllVisibleNodes
- Logic: respect collapsed state, skip hidden children, handle boundaries
- Clean architecture: separate navigation concern but within tree package
- Tests: TestNavigationNext, TestNavigationPrev, TestNavigationSkipsCollapsed, TestNavigationBoundaries

Phase 5: Integration & Polish (30 min)
- Run full test suite: go test ./internal/tui/tree -v -cover
- Target: >90% coverage
- Code review checklist:
  * No external dependencies (only domain + lipgloss + stdlib)
  * All exports have godoc
  * Constants for tree characters, indentation
  * Error handling for edge cases
  * Clean code principles applied
- Verify performance with 100-node tree benchmark

Dependencies:
- internal/domain (Project type) - already exists
- github.com/charmbracelet/lipgloss - add if not present

Enables:
- Task-21 (Data Loading) - will use tree.BuildFromProjects
- Task-18 (Navigation & State) - will use tree navigation helpers

Estimated effort: 5-6 hours for single developer
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
Phase 1-5 complete:
- Created internal/tui/tree/ package with clean architecture
- node.go: TreeNode struct with ID, Name, Depth, IsExpanded, Children, Parent, Data
- builder.go: BuildFromProjects transforms flat project list to hierarchical tree
- renderer.go: Render with lipgloss styling for selected nodes
- navigation.go: GetNextVisibleNode, GetPrevVisibleNode, GetAllVisibleNodes, FindNodeByID
- constants.go: TreeIndent, TreeBranch, TreeLast, TreeVertical, ExpandedIcon, CollapsedIcon

Tests: 51 tests, 95.0% coverage
All AC criteria met: tree indicators, expand/collapse, unlimited nesting, godoc comments
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary
Implemented tree package for hierarchical project display with unlimited nesting support.

## Changes
- Created `internal/tui/tree/` package with clean architecture (no DB dependencies)
- **node.go**: TreeNode model with expand/collapse state, depth tracking, parent/child refs
- **builder.go**: BuildFromProjects() transforms flat domain.Project list to tree structure
- **renderer.go**: lipgloss-styled rendering with tree indicators (├─ └─ │) and selection highlighting
- **navigation.go**: GetNextVisibleNode, GetPrevVisibleNode, FindNodeByID, ExpandAll, CollapseAll
- **constants.go**: Named constants for all tree characters (no magic strings)

## Test Coverage
- 51 unit tests covering all scenarios
- 95.0% statement coverage (exceeds 90% target)
- Fast, independent tests following F.I.R.S.T. principles

## Acceptance Criteria Met
- Hierarchical display with visual indicators
- Expand/collapse with [+]/[-] icons
- Unlimited nesting via recursive building
- Navigation helpers skip collapsed nodes
- Full godoc documentation on all exports
- Inward dependency: domain → tree → lipgloss/stdlib only
<!-- SECTION:FINAL_SUMMARY:END -->
