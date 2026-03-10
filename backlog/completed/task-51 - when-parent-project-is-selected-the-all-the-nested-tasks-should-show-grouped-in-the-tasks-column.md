---
id: TASK-51
title: >-
  when parent project is selected the all the nested tasks should show grouped
  in the tasks column
status: Done
assignee: []
created_date: '2026-03-06 20:53'
updated_date: '2026-03-07 20:59'
labels:
  - feature
  - tui
  - service-layer
  - enhancement
dependencies: []
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
When a parent project is selected in the Projects column, display all tasks from the project and its nested subprojects grouped by subproject in the Tasks column. Tasks should be organized hierarchically with collapsible groups, showing direct project tasks first, followed by nested subproject tasks with indented headers.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [ ] #1 Service layer: Add ListByProjectRecursive method to TaskService that retrieves tasks from a project and all its nested subprojects recursively
- [ ] #2 Service layer: Implement belongsToProject helper method for recursive project hierarchy traversal
- [ ] #3 Service layer: Add ListByProjectRecursive to TaskServiceInterface
- [ ] #4 Service layer: Write comprehensive unit tests with 80%+ coverage for recursive task loading
- [ ] #5 Data model: Create GroupedTasks struct to represent tasks organized by subproject with group metadata (project ID, name, expanded state)
- [ ] #6 TUI commands: Update LoadTasksCmd to use ListByProjectRecursive instead of ListByProject
- [ ] #7 TUI rendering: Render tasks grouped by subproject with indented headers showing subproject name
- [ ] #8 TUI rendering: Show direct project tasks at the top without header (ungrouped)
- [ ] #9 TUI state: Add expandedTaskGroups map to Model to track which groups are expanded/collapsed
- [ ] #10 TUI interaction: Support expand/collapse of task groups with Enter or Space key
- [ ] #11 TUI interaction: Update task navigation to work with grouped structure (skip headers when navigating)
- [ ] #12 Visual design: Use indentation (2 spaces) for tasks under each subproject group header
- [ ] #13 Visual design: Use subtle styling for group headers (dimmed color, no reverse highlight)
- [ ] #14 Error handling: Handle empty project (no tasks at any level) gracefully
- [ ] #15 Error handling: Handle orphaned subprojects (parent doesn't exist) gracefully
- [ ] #16 Performance: Ensure O(n) time complexity where n = total tasks across all subprojects
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
IMPLEMENTATION PLAN: Nested Task Grouping by Subproject

## Decision: SPLIT into Subtasks

This task is too large (16 ACs, 13-19 hours) and will be split into 7 subtasks (51A-51G) for better parallelization and tracking.

## Task Breakdown & Dependencies

### Wave 1 (PARALLEL - No Dependencies)
**Task 51A: Service Layer - Recursive Task Loading**
- AC: #1-4 (Service methods + unit tests)
- File: internal/service/task_service.go
- Methods: ListByProjectRecursive, belongsToProject
- Tests: 85%+ coverage with table-driven tests
- Effort: 3-4 hours

**Task 51B: Data Model - GroupedTasks Structure**
- AC: #5 (GroupedTasks struct)
- File: internal/domain/task_group.go (new)
- Structs: TaskGroup, GroupedTasks
- Tests: Unit tests for grouping logic
- Effort: 2-3 hours

### Wave 2 (Sequential - Depends on Wave 1)
**Task 51C: TUI Commands - Update LoadTasksCmd**
- AC: #6 (Use ListByProjectRecursive)
- Depends on: 51A, 51B
- File: internal/tui/commands.go
- Update: LoadTasksCmd to use recursive loading
- Effort: 1 hour

### Wave 3 (Sequential - Depends on Wave 2)
**Task 51D: TUI Rendering - Grouped Task Display**
- AC: #7-8, #12-13 (Rendering + visual design)
- Depends on: 51B, 51C
- Files: internal/tui/renderer.go, internal/tui/tui.go
- Features: Grouped rendering, indentation, styling
- Effort: 2-3 hours

### Wave 4 (Sequential - Depends on Wave 3)
**Task 51E: TUI Interaction - Expand/Collapse & Navigation**
- AC: #9-11 (State + interaction)
- Depends on: 51D
- Files: internal/tui/navigator.go, internal/tui/state.go
- Features: Group toggle, navigation skip headers
- Effort: 2-3 hours

### Wave 5 (Sequential - Depends on Wave 4)
**Task 51F: Error Handling & Edge Cases**
- AC: #14-15 (Error handling)
- Depends on: 51A-51E
- Files: All layers
- Features: Graceful handling of edge cases
- Effort: 1 hour

### Wave 6 (Sequential - Depends on Wave 5)
**Task 51G: Performance Optimization**
- AC: #16 (Performance)
- Depends on: 51A-51F
- Files: Service + TUI
- Features: Ensure O(n) complexity, benchmark tests
- Effort: 1-2 hours

## Parallelization Strategy

**Maximum Parallelization: 2 developers**
- Developer 1: Service layer (51A)
- Developer 2: Domain model (51B)
- Then sequential for TUI layers (51C-51G)

**Single Developer:**
- Complete 51A → 51B in parallel or sequentially
- Then complete 51C → 51D → 51E → 51F → 51G sequentially

## Testing Strategy (golang-testing patterns)

### Service Layer Tests (Task 51A)
**File**: internal/service/task_service_test.go

**Test Cases** (table-driven):
1. Empty project (no tasks at any level)
2. Direct tasks only (no subprojects)
3. Tasks in nested subprojects (1 level deep)
4. Tasks in deeply nested subprojects (3+ levels)
5. Mixed: direct + nested tasks
6. Soft-deleted tasks excluded
7. Empty projectID input
8. Database error handling
9. Performance with 1000+ tasks
10. Circular project hierarchy (prevented by existing validation)

**Coverage Target**: 85%+ for ListByProjectRecursive and GroupTasksByProject

### Domain Model Tests (Task 51B)
**File**: internal/domain/task_group_test.go

**Test Cases**:
1. Empty GroupedTasks initialization
2. Group tasks by project correctly
3. Direct tasks vs nested tasks separation
4. TotalCount calculation
5. Group ordering

**Coverage Target**: 80%+

### TUI Component Tests (Tasks 51C-51E)
**Files**: internal/tui/renderer_test.go, navigator_test.go

**Test Cases**:
1. Render empty grouped tasks
2. Render direct tasks only
3. Render groups with expanded state
4. Render groups with collapsed state
5. Mixed direct + groups
6. Task completion styling
7. Selected task highlighting
8. Text truncation in narrow columns
9. Group toggle (expand/collapse)
10. Navigation skip group headers

**Coverage Target**: 75%+

### Integration Tests (Task 51F)
**File**: internal/tui/integration_test.go

**User Flows**:
1. Select parent project → see grouped tasks
2. Collapse a group → tasks hidden
3. Expand a group → tasks shown
4. Navigate through grouped tasks
5. Toggle task completion in grouped view
6. Switch between areas → state persists
7. Error: orphaned subproject
8. Error: empty project tree

### Performance Tests (Task 51G)
**File**: internal/service/task_service_bench_test.go

**Benchmarks**:
1. BenchmarkListByProjectRecursive_100tasks
2. BenchmarkListByProjectRecursive_1000tasks
3. BenchmarkListByProjectRecursive_10000tasks
4. BenchmarkGroupTasksByProject_100tasks
5. BenchmarkGroupTasksByProject_1000tasks

**Performance Targets**:
- ListByProjectRecursive: O(n) where n = total tasks
- GroupTasksByProject: O(n) where n = total tasks
- No N+1 queries
- <100ms for 1000 tasks

## Documentation Updates

### Files to Update:

1. **docs/TUI.md** (Task 51E)
   - Add section: "Task Grouping by Subproject"
   - Keyboard shortcuts: Enter/Space to toggle groups
   - Visual design rationale
   - Performance characteristics

2. **docs/architecture/03-service-layer.md** (Task 51A)
   - Document ListByProjectRecursive method
   - Explain recursive loading pattern
   - Performance considerations

3. **internal/service/README.md** (Task 51A)
   - Update TaskServiceInterface documentation
   - Add usage examples

4. **internal/domain/README.md** (Task 51B)
   - Document TaskGroup and GroupedTasks
   - Explain grouping structure

5. **CHANGELOG.md** (Final)
   - Entry: "Feature: Nested task grouping by subproject"

## Implementation Guidelines (from golang-patterns)

### Code Quality:
- Use gofmt and golangci-lint on all code
- Add context.Context to all blocking operations
- Handle all errors explicitly (no naked returns)
- Write table-driven tests with subtests
- Document all exported functions and types
- Propagate errors with fmt.Errorf("%w", err)

### Service Layer (Task 51A):
- Use existing patterns from ProjectService.ListBySubareaRecursive
- Return early on errors
- Keep methods focused (single responsibility)
- Use meaningful error messages with context

### Domain Model (Task 51B):
- Make zero value useful
- Use factory methods if needed
- Keep structs immutable where possible
- Add validation methods if needed

### TUI Layer (Tasks 51C-51E):
- Follow Bubble Tea patterns (Elm architecture)
- Use weight-based sizing for flexible layouts
- Always truncate text explicitly (never auto-wrap)
- Account for borders in height calculations
- Match mouse detection to layout orientation

## Success Criteria

Each subtask must:
- [ ] All assigned acceptance criteria checked
- [ ] Unit tests passing with required coverage
- [ ] Integration tests passing (if applicable)
- [ ] Linting passing (golangci-lint)
- [ ] Race detector passing (go test -race)
- [ ] Documentation updated
- [ ] Code reviewed

Final criteria:
- [ ] All 7 subtasks completed
- [ ] All 16 ACs satisfied
- [ ] Full test suite passing
- [ ] Manual testing complete
- [ ] Performance benchmarks passing
- [ ] Documentation complete

## Risk Mitigation

1. **Integration Issues**: Clear interfaces between layers, contract tests
2. **Performance**: Profile early with 1000+ tasks, optimize if needed
3. **State Management**: Thorough testing of expand/collapse state
4. **Navigation**: Test edge cases (empty groups, all collapsed)
5. **Visual Design**: User feedback on indentation/styling

## Estimated Timeline

**Total Effort**: 13-19 hours
**With Parallelization**: 10-15 hours (2 developers)
**Sequential**: 13-19 hours (1 developer)

**Suggested Timeline (Single Developer)**:
- Day 1: Tasks 51A + 51B (5-7 hours)
- Day 2: Tasks 51C + 51D (3-6 hours)
- Day 3: Tasks 51E + 51F + 51G (4-6 hours)

**Suggested Timeline (2 Developers)**:
- Day 1 Morning: Dev 1 = 51A (3-4h), Dev 2 = 51B (2-3h)
- Day 1 Afternoon: Dev 1 = 51C (1h), Dev 2 = support 51A tests
- Day 2 Morning: Dev 1 = 51D (2-3h), Dev 2 = 51D support
- Day 2 Afternoon: Dev 1 = 51E (2-3h), Dev 2 = 51E support
- Day 3: Dev 1 = 51F (1h) + 51G (1-2h), Dev 2 = documentation
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
## Task Split Completed

This task has been split into 7 subtasks for better parallelization and tracking:

### Subtasks Created:

**Wave 1 (PARALLEL - Can start immediately):**
- TASK-52: Service Layer - Recursive Task Loading (51A)
  - AC: #1-4
  - Effort: 3-4 hours
  
- TASK-54: Data Model - GroupedTasks Structure (51B)
  - AC: #5
  - Effort: 2-3 hours

**Wave 2 (Sequential - After Wave 1):**
- TASK-57: TUI Commands - Update LoadTasksCmd (51C)
  - AC: #6
  - Depends on: TASK-52, TASK-54
  - Effort: 1 hour

**Wave 3 (Sequential - After Wave 2):**
- TASK-58: TUI Rendering - Grouped Display (51D)
  - AC: #7-8, #12-13
  - Depends on: TASK-54, TASK-57
  - Effort: 2-3 hours

**Wave 4 (Sequential - After Wave 3):**
- TASK-56: TUI Interaction - Expand/Collapse (51E)
  - AC: #9-11
  - Depends on: TASK-58
  - Effort: 2-3 hours

**Wave 5 (Sequential - After Wave 4):**
- TASK-55: Error Handling (51F)
  - AC: #14-15
  - Depends on: TASK-52, TASK-54, TASK-57, TASK-58, TASK-56
  - Effort: 1 hour

**Wave 6 (Sequential - After Wave 5):**
- TASK-53: Performance Optimization (51G)
  - AC: #16
  - Depends on: All previous tasks
  - Effort: 1-2 hours

### AC Mapping:

| AC | Task | Status |
|----|------|--------|
| #1-4 | TASK-52 | Service layer methods |
| #5 | TASK-54 | Data model |
| #6 | TASK-57 | TUI commands |
| #7-8, #12-13 | TASK-58 | TUI rendering |
| #9-11 | TASK-56 | TUI interaction |
| #14-15 | TASK-55 | Error handling |
| #16 | TASK-53 | Performance |

### Parallelization Benefits:

✅ **2 developers can work in parallel** (TASK-52 + TASK-54)
✅ **Smaller, focused tasks** (2-4 hours each)
✅ **Clear dependencies** (Wave 1 → 2 → 3 → 4 → 5 → 6)
✅ **Earlier testing** (each task can be tested independently)
✅ **Easier code reviews** (smaller PRs)

### Timeline:

**Sequential (1 developer):**
- Day 1: TASK-52 + TASK-54 (5-7 hours)
- Day 2: TASK-57 + TASK-58 (3-6 hours)
- Day 3: TASK-56 + TASK-55 + TASK-53 (4-6 hours)
- **Total: 12-19 hours**

**Parallel (2 developers):**
- Dev 1: TASK-52 (3-4h) → TASK-57 (1h) → TASK-58 (2-3h)
- Dev 2: TASK-54 (2-3h) → support → TASK-56 (2-3h)
- Both: TASK-55 (1h) + TASK-53 (1-2h)
- **Total: 10-15 hours**

### Testing Coverage:

Each subtask includes detailed testing strategy:
- **Unit tests** with table-driven patterns
- **Integration tests** for end-to-end flows
- **Benchmark tests** for performance
- **Coverage targets**: 80-85%+

### Skills Used:

✅ **golang-pro** - Concurrency patterns, context usage
✅ **golang-patterns** - Error handling, interfaces, table-driven tests
✅ **bubbletea** - TUI architecture, rendering, interaction
✅ **golang-testing** - TDD, benchmarks, fuzzing, coverage

### Next Steps:

1. Start with TASK-52 and TASK-54 in parallel
2. Follow wave sequence (1 → 2 → 3 → 4 → 5 → 6)
3. Run tests after each task
4. Update documentation after completion
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Task Organization Complete

Task-51 has been successfully split into 7 subtasks (TASK-52 through TASK-58) with clear dependencies and detailed implementation plans.

**Key Achievements:**
- ✅ Identified parallel vs sequential work
- ✅ Created detailed implementation plans for each subtask
- ✅ Specified testing strategies (unit, integration, benchmarks)
- ✅ Documented error handling and edge cases
- ✅ Planned performance optimizations
- ✅ Included documentation updates

**Subtasks:**
- TASK-52: Service Layer (51A) - Foundation
- TASK-53: Performance (51G) - Optimization
- TASK-54: Data Model (51B) - Core structures
- TASK-55: Error Handling (51F) - Robustness
- TASK-56: TUI Interaction (51E) - User actions
- TASK-57: TUI Commands (51C) - Data loading
- TASK-58: TUI Rendering (51D) - Visual display

**Ready for implementation following the wave-based approach.**
<!-- SECTION:FINAL_SUMMARY:END -->
