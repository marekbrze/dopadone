---
id: TASK-16
title: 'TUI 14B: Data Layer & Tree Rendering'
status: Done
assignee: []
created_date: '2026-03-03 12:30'
updated_date: '2026-03-03 14:36'
labels:
  - tui
  - mvp
  - phase2
dependencies:
  - TASK-15
references:
  - internal/domain/project.go
  - internal/domain/subarea.go
  - internal/domain/task.go
  - internal/domain/area.go
  - internal/db/querier.go
  - internal/tui/app.go
documentation:
  - 'https://github.com/charmbracelet/bubbletea'
  - 'https://github.com/charmbracelet/bubbles/tree/main/spinner'
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Implement data loading from database and project tree rendering with unlimited nesting support via parent_id. Includes async data loading with loading states, interactive tree view with expand/collapse, and proper state management for selections across areas.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Subareas load and display for selected area
- [ ] #2 Projects display in hierarchical tree structure with visual indicators (├─, └─, │) for parent/child relationships
- [ ] #3 Nested projects support expand/collapse behavior with visual indicators (+/-)
- [ ] #4 Tasks load and display for selected project
- [ ] #5 Unlimited project nesting depth supported via recursive tree building using parent_id
- [ ] #6 Empty state messages shown with contextual hints (e.g., 'No subareas - press a to add')
- [ ] #7 Loading spinner displayed while fetching data from database using bubbles/spinner component
- [ ] #8 First area auto-selected on app initialization and its data loads automatically
- [ ] #9 Tree building logic isolated in internal/tui/tree package for reusability and testability
- [ ] #10 Data loading follows clean architecture: TUI depends on domain types, uses repository interfaces, no direct DB access
- [ ] #11 Unit tests for tree building with various nesting scenarios (0, 1, 2, 5+ levels)
- [ ] #12 Unit tests for data loading with mock database/repository
- [ ] #13 Unit tests for expand/collapse state management
- [ ] #14 Loading state management prevents duplicate data fetches
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Task-16 Implementation Plan (Refined)

## Assessment: Split Recommended

Current task has 14 ACs spanning two distinct concerns:
1. **Tree rendering** (pure logic, no I/O) 
2. **Data loading** (database integration, async operations)

**Recommendation**: Split into 2 sequential subtasks for better parallelization and cleaner PRs.

---

## Option A: Split Into Subtasks (Recommended)

### Subtask 16A: Tree Rendering Package
**Scope**: ACs 2, 3, 5, 9, 11
**Duration**: 3-4 hours
**Dependencies**: None (can start immediately, parallel with task-19)

**Deliverables**:
- internal/tui/tree/node.go - TreeNode struct
- internal/tui/tree/builder.go - BuildTree() function
- internal/tui/tree/renderer.go - Render() function
- internal/tui/tree/tree_test.go - Comprehensive tests

**Tests**:
- TestBuildEmptyTree
- TestBuildSingleLevel
- TestBuildMultiLevel (2, 3, 5+ depths)
- TestTreeIndicators (├─ └─ │)
- TestExpandCollapse

**Documentation**:
- Package doc comment for internal/tui/tree
- Godoc for exported functions

---

### Subtask 16B: Data Loading & Integration
**Scope**: ACs 1, 4, 6, 7, 8, 10, 12, 13, 14
**Duration**: 5-6 hours
**Dependencies**: Task-16A (tree package), Task-15 (completed)

**Deliverables**:
- internal/tui/loader.go - Data loading commands
- internal/tui/messages.go - Async message types
- Updated app.go - Model expansion, Update handlers
- Updated tui_test.go - Integration tests

**Tests**:
- TestLoadAreasCmd (mock repo)
- TestLoadSubareasCmd (mock repo)
- TestLoadProjectsCmd (mock repo)
- TestAutoSelectFirstArea
- TestSelectionCascades
- TestLoadingStatePreventsDuplicates

**Documentation**:
- Update internal/tui/README.md (or create) with architecture notes
- Comment message flow in Update()

---

## Option B: Single Task (Alternative)

If keeping as single task, use the existing 8-phase plan with these modifications:

### Sequential Phases (Must be done in order)
1. **Phase 1: Tree Package** (1.5-2h) - Foundation, no dependencies
2. **Phase 2: Data Loader** (1.5-2h) - Depends on Phase 1 completion
3. **Phase 3: Spinner** (30min) - Can start anytime after Phase 1
4. **Phase 4: Model State** (1.5h) - Depends on Phase 2
5. **Phase 5: Integration** (1.5h) - Depends on Phase 4
6. **Phase 6: Navigation Prep** (1h) - Depends on Phase 5
7. **Phase 7: Testing** (2h) - Can start incrementally after each phase
8. **Phase 8: Validation** (30min) - Final verification

### Parallel Opportunities (Single Task Flow)
- **Phase 1** can be done by one dev while another works on **task-19** (Quick-Add Modal)
- **Phase 3** (spinner) can be done independently during Phase 2-5
- **Phase 7** tests can be written alongside each phase (TDD)

---

## Test Plan (Detailed)

### Tree Package Tests (tree_test.go)
```go
// Group 1: Tree Building
TestBuildTreeFromEmptyList
TestBuildTreeFromFlatList (single level)
TestBuildTreeFromNestedList (multiple levels)
TestBuildTreeWithMultipleRoots
TestBuildTreeDeepNesting (5+ levels)
TestBuildTreePreservesOrder (Position field)

// Group 2: Tree Rendering
TestRenderTreeBasic
TestRenderTreeWithIndicators
TestRenderTreeWithExpandCollapse
TestRenderTreeEmpty

// Group 3: Tree Navigation (prepare for task-18)
TestGetNextVisibleNode
TestGetPrevVisibleNode
TestSkipCollapsedNodes
```

### Data Loading Tests (loader_test.go)
```go
// Group 1: Commands
TestLoadAreasCmdReturnsCorrectMsg
TestLoadSubareasCmdWithAreaID
TestLoadProjectsCmdWithSubareaID
TestLoadTasksCmdWithProjectID
TestLoadCmdReturnsErrorOnFailure

// Group 2: State Management
TestAutoSelectFirstArea
TestAutoSelectFirstSubarea
TestCascadingDataLoads
TestPreventDuplicateLoads
```

### Integration Tests (app_test.go or tui_test.go)
```go
// Group 1: Initialization
TestInitialModelLoadsFirstArea
TestFirstAreaDataLoadsOnStartup

// Group 2: Selection Cascades
TestSubareaSelectionLoadsProjects
TestProjectSelectionLoadsTasks
TestAreaSwitchReloadsData

// Group 3: Loading States
TestSpinnerShowsDuringLoad
TestSpinnerHidesAfterLoad
TestEmptyStateMessages
```

---

## Documentation Updates

### Code Documentation
1. **internal/tui/tree/doc.go**
   - Package overview
   - Usage examples
   - Architecture decision (why separate package)

2. **internal/tui/loader.go**
   - Comment each LoadCmd function
   - Document message flow

3. **internal/tui/app.go**
   - Update Model struct comments
   - Document Update() message handling flow

4. **internal/tui/messages.go** (new file)
   - Document each message type
   - Explain async pattern

### Project Documentation (if needed)
5. **internal/tui/README.md** (create if not exists)
   - Architecture overview
   - Data flow diagram (text-based)
   - How to add new columns/features

---

## Parallel Work Strategy

### With Split (Recommended)
```
Time →
Dev 1: [16A: Tree Package    ][16B: Data Loading]
Dev 2:         [Task-19: Quick-Add Modal        ]
Dev 3:                   [Task-18: Navigation   ]
```

### Without Split
```
Time →
Dev 1: [Task-16 (all phases, 8-10 hours)        ]
Dev 2: [Wait for 16A completion...][Task-19    ]
Dev 3: [Wait for 16B...          ][Task-18    ]
```

---

## Dependency Graph

```
Task-15 (Done) ─┬─→ Task-16A (Tree) ─┬─→ Task-16B (Data) ─→ Task-18 (Nav)
                │                     │
                └─→ Task-19 (Modal) ─┘                      ↓
                                                        Task-17 (Polish)
```

**Parallel paths**:
- Task-16A can start immediately (no dependencies)
- Task-19 can start immediately (only depends on Task-15)
- Task-18 must wait for Task-16B
- Task-17 waits for all others

---

## Recommended Approach

### Decision: Split Task-16 into 16A and 16B

**Rationale**:
1. Smaller, focused PRs (easier review)
2. Enables parallelization (16A + task-19 simultaneously)
3. Tree package is pure logic (fully testable in isolation)
4. Clear acceptance criteria split
5. Reduces risk of large PR conflicts

**Next Steps**:
1. Create Task-16A: Tree Rendering Package
2. Create Task-16B: Data Loading & Integration
3. Mark Task-16B as dependent on Task-16A
4. Update Task-18 dependency: 15, 16A, 16B
5. Archive or update original Task-16

**Acceptance Criteria Distribution**:
- Task-16A: ACs 2, 3, 5, 9, 11 (tree rendering)
- Task-16B: ACs 1, 4, 6, 7, 8, 10, 12, 13, 14 (data loading)
<!-- SECTION:PLAN:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
Superseded by TASK-20 (Tree Rendering Package) and TASK-21 (Data Loading & Integration). All 14 ACs split between the two new tasks.
<!-- SECTION:FINAL_SUMMARY:END -->
