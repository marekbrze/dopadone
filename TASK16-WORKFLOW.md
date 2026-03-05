# Task-16 Implementation Workflow & Task Splitting Summary

## Decision: Split Task-16 into Two Subtasks

**Original Task-16** had 14 acceptance criteria covering two distinct concerns:
1. Tree rendering (pure logic)
2. Data loading (database integration)

**Split into**:
- **Task-20 (16A)**: Tree Rendering Package (5 ACs)
- **Task-21 (16B)**: Data Loading & Integration (9 ACs)

---

## Dependency Graph

```
Task-15 (Done) ──┬──→ Task-20 (Tree) ──┬──→ Task-21 (Data) ──┬──→ Task-18 (Navigation)
                 │                      │                      │
                 └──→ Task-19 (Modal) ─┘                      └──→ Task-17 (Polish)
```

### Parallel Work Opportunities

**After Task-15**:
- ✅ Task-20 (Tree) - Can start immediately (no dependencies)
- ✅ Task-19 (Modal) - Can start immediately (only depends on Task-15)

**After Task-20**:
- ✅ Task-21 (Data) - Can start (depends on Task-20)
- ✅ Task-18 (Navigation) - Can start parallel with Task-21

**After Task-21**:
- ✅ Task-18 (Navigation) - Can continue/complete
- ✅ Task-17 (Polish) - Can start (depends on 15, 16, 18, 19)

---

## Task Summary

### Task-20: Tree Rendering Package
**Priority**: High  
**Duration**: 5-6 hours  
**Dependencies**: None  
**ACs**: 5 (tree rendering, expand/collapse, unlimited depth, testing)  

**Files**:
- internal/tui/tree/node.go
- internal/tui/tree/builder.go
- internal/tui/tree/renderer.go
- internal/tui/tree/tree_test.go

**Key Features**:
- TreeNode struct with unlimited nesting
- BuildFromProjects() for converting flat list to tree
- Render() with tree characters (├─ └─ │)
- Navigation helpers (prep for task-18)
- Comprehensive unit tests (>90% coverage)

---

### Task-21: Data Loading & Integration
**Priority**: High  
**Duration**: 8-9 hours  
**Dependencies**: Task-20 (tree package)  
**ACs**: 9 (data loading, spinner, empty states, cascade selection)  

**Files**:
- internal/tui/messages.go (new)
- internal/tui/loader.go (new)
- internal/tui/loader_test.go (new)
- internal/tui/app.go (updates)
- cmd/projectdb/tui.go (updates)

**Key Features**:
- Repository injection pattern
- Async data loading commands
- Loading spinner integration
- Empty state messages
- Cascade loading (Area → Subarea → Projects → Tasks)
- Auto-select first item in each column
- Loading state prevents duplicate fetches

---

## Sequential vs Parallel Work

### Within Each Task (Sequential)
Both Task-20 and Task-21 have internal phases that **must be done sequentially**:
- Foundation → Implementation → Integration → Testing

### Between Tasks (Parallel)
- **Task-20** and **Task-19** can be done simultaneously
- **Task-21** must wait for **Task-20**
- **Task-18** can start after **Task-20** (parallel with Task-21)

---

## Test Strategy

### Task-20 Tests (Tree Package)
- Tree Building: 7+ test cases
- Tree Rendering: 6+ test cases
- Navigation: 4+ test cases
- **Coverage Goal**: >90%

### Task-21 Tests (Data Loading)
- Data Loaders: 5+ test cases (with mock repos)
- Message Handling: 5+ test cases
- Integration: 3+ test cases
- **Coverage Goal**: >85%

---

## Documentation Updates

### Code Documentation
- Package doc for internal/tui/tree
- Godoc for all exported functions
- Comment message flow in Update()
- Architecture notes in loader.go

### Project Documentation
- internal/tui/README.md (create if needed)
- Architecture overview
- Data flow diagram

---

## Risk Mitigation

### Task-20 Risks
1. **Performance with deep trees**: Test with 10+ levels, lazy rendering
2. **Circular references**: Detect cycles during build
3. **Memory usage**: Benchmark with 1000+ nodes

### Task-21 Risks
1. **Loading state complexity**: Clear state machine, single flag
2. **Error handling**: Basic here, full polish in Task-17
3. **Empty database**: Guide users with empty state messages

---

## Effort Estimation

### Single Developer (Sequential)
- Task-20: 5-6 hours
- Task-21: 8-9 hours
- **Total**: 13-15 hours

### Team (Parallel)
- Dev 1: Task-20 (5-6 hours)
- Dev 2: Task-19 (parallel, 4-5 hours)
- Dev 3: Task-21 (after Task-20, 8-9 hours)
- **Wall-clock**: 13-15 hours (Task-21 is bottleneck)

### Optimized Parallel
```
Time → 0----2----4----6----8----10---12---14 hours
Dev 1: [Task-20 Tree        ]
Dev 2:    [Task-19 Modal        ]
Dev 3:              [Task-21 Data         ]
Dev 1:                        [Task-18 Nav     ]
Dev 2:                        [Task-17 Polish  ]
```

**Wall-clock**: ~14-15 hours (Task-21 still bottleneck)

---

## Acceptance Criteria Distribution

### Original Task-16 (14 ACs)
1. ✅ Subareas load and display → Task-21
2. ✅ Projects display in tree → Task-20
3. ✅ Expand/collapse behavior → Task-20
4. ✅ Tasks load and display → Task-21
5. ✅ Unlimited nesting depth → Task-20
6. ✅ Empty state messages → Task-21
7. ✅ Loading spinner → Task-21
8. ✅ First area auto-selected → Task-21
9. ✅ Tree logic isolated → Task-20
10. ✅ Clean architecture → Task-21
11. ✅ Tree building tests → Task-20
12. ✅ Data loading tests → Task-21
13. ✅ Expand/collapse tests → Task-21
14. ✅ Prevent duplicate loads → Task-21

### Split Tasks
- **Task-20**: ACs 2, 3, 5, 9, 11 (5 ACs)
- **Task-21**: ACs 1, 4, 6, 7, 8, 10, 12, 13, 14 (9 ACs)

---

## Next Steps

1. ✅ Task-16 plan updated with split recommendation
2. ✅ Task-20 created with detailed plan
3. ✅ Task-21 created with detailed plan
4. ✅ Dependencies set: Task-21 depends on Task-20
5. ✅ Task-18 dependencies updated: depends on Task-20, Task-21
6. ✅ Task-20 completed (Tree Rendering Package)
7. ✅ Task-21 completed (Data Loading & Integration)
8. ✅ Documentation updated to reflect new architecture

---

## Commands to Verify Structure

```bash
# View task details
backlog task 20 --plain
backlog task 21 --plain

# Check dependencies
backlog task list --plain | grep -E "(TASK-20|TASK-21)"

# View implementation plans
backlog task 20 --plain  # Contains detailed plan
backlog task 21 --plain  # Contains detailed plan
```

---

## Task Completion Summary

### Task-20: Tree Rendering Package ✅ COMPLETE
**Completed**: 2026-03-03  
**Test Coverage**: 95.0%  
**Key Achievements**:
- Implemented internal/tui/tree/ package with clean architecture
- TreeNode model with expand/collapse, unlimited nesting support
- Tree builder converting flat domain.Project to hierarchical structure
- Lipgloss-styled renderer with tree indicators (├─ └─ │)
- Navigation helpers (GetNextVisibleNode, GetPrevVisibleNode, FindNodeByID)
- 51 unit tests covering all scenarios
- Full godoc documentation on all exports
- No external dependencies beyond domain + lipgloss + stdlib

### Task-21: Data Loading & Integration ✅ COMPLETE  
**Completed**: 2026-03-03  
**Test Coverage**: 85.5%  
**Key Achievements**:
- Created messages.go for async operations (LoadAreasMsg, LoadSubareasMsg, etc.)
- Created converters.go for DB→Domain type conversion
- Created commands.go with loader commands (all under 20 lines)
- Extended Model with repository, data slices, selections, loading flags, spinner
- Integrated bubbles/spinner for loading feedback
- Implemented cascade loading: Area → Subarea → Projects → Tasks
- Auto-selection of first item after data loads
- Loading states prevent duplicate fetches
- Empty state messages with keyboard hints
- Repository injection pattern following clean architecture
- 85.5% test coverage with MockQuerier

### Documentation Updates
- ✅ TUI Architecture document (doc-3) updated with tree package and data loading sections
- ✅ TASK16-WORKFLOW.md updated with completion status
- ✅ All AC criteria met for both tasks
- ✅ Clean architecture principles verified and maintained

---

## Original Task-16 Status

**Recommendation**: Archive Task-16 or mark as "Split into Task-20 and Task-21"

```bash
# Option 1: Archive
backlog task archive 16

# Option 2: Update description
backlog task edit 16 -d "Split into Task-20 (Tree Rendering) and Task-21 (Data Loading)"
backlog task edit 16 -s Cancelled
```

