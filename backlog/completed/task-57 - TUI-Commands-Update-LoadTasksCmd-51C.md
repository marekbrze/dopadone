---
id: TASK-57
title: 'TUI Commands: Update LoadTasksCmd (51C)'
status: Done
assignee:
  - '@opencode'
created_date: '2026-03-06 21:30'
updated_date: '2026-03-07 09:09'
labels:
  - tui
  - commands
dependencies:
  - TASK-52
  - TASK-54
references:
  - task-51
  - internal/tui/commands.go
priority: high
---

## Description

<!-- SECTION:DESCRIPTION:BEGIN -->
Update LoadTasksCmd to use ListByProjectRecursive and populate GroupedTasks structure. Depends on tasks 51A (TASK-52) and 51B (TASK-54). Part of task-51 nested task grouping feature.
<!-- SECTION:DESCRIPTION:END -->

## Acceptance Criteria
<!-- AC:BEGIN -->
- [x] #1 #1 Update LoadTasksCmd to use ListByProjectRecursive instead of ListByProject
- [x] #2 #2 Update TasksLoadedMsg to include GroupedTasks field
- [x] #3 #3 Update handleTasksLoaded to populate m.groupedTasks and initialize expanded state
- [x] #4 #4 Write tests for command update
<!-- AC:END -->

## Implementation Plan

<!-- SECTION:PLAN:BEGIN -->
# Detailed Implementation Plan: Task-57 - TUI Commands Update for GroupedTasks

## Overview
This task integrates the GroupedTasks domain model into the TUI command layer, enabling nested task grouping display. It bridges TASK-52 (recursive loading) and TASK-54 (domain model) with the TUI presentation layer.

## Task Assessment
- **Scope**: Medium (6 phases, multiple layers)
- **Complexity**: Medium (requires coordination across service/tui layers)
- **Dependencies**: TASK-52 ✅ DONE, TASK-54 ⏳ IN PROGRESS
- **Decision**: ✅ NO SPLIT - Cohesive integration task, phases are tightly coupled

## Critical Path Analysis

### Sequential Dependencies (MUST be done in order)
```
Phase 1 (Service Layer) → Phase 2 (Messages) → Phase 3 (Commands) → Phase 4 (Model/Handler) → Phase 5 (Testing) → Phase 6 (Integration Tests)
```

### Parallel Opportunities
- Within Phase 5: Different test suites can be written in parallel
- Documentation (Phase 7) can run in parallel with testing

## Phase 1: Service Layer - GetGroupedTasks Method (45 min)
**SEQUENTIAL** - Must complete before Phase 2

### 1.1 Add to TaskServiceInterface (5 min)
**File**: internal/service/interfaces.go

```go
// GetGroupedTasks retrieves all tasks from a project and its nested subprojects,
// grouped by subproject with group metadata.
// Returns a GroupedTasks structure ready for TUI rendering.
GetGroupedTasks(ctx context.Context, projectID string) (*domain.GroupedTasks, error)
```

**Rationale**: Defines contract for grouped task retrieval

### 1.2 Implement GetGroupedTasks in TaskService (30 min)
**File**: internal/service/task_service.go

**Implementation**:
```go
func (s *TaskService) GetGroupedTasks(ctx context.Context, projectID string) (*domain.GroupedTasks, error) {
    // Load all tasks recursively
    tasks, err := s.ListByProjectRecursive(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("get grouped tasks: %w", err)
    }
    
    // Build project names map
    projectNames := make(map[string]string)
    
    // Extract unique project IDs from tasks
    projectIDs := make(map[string]bool)
    for _, task := range tasks {
        if task.ProjectID != "" {
            projectIDs[task.ProjectID] = true
        }
    }
    
    // Fetch project names from repository
    for pid := range projectIDs {
        // Use project service to get project details
        project, err := s.projectService.GetByID(ctx, pid)
        if err == nil && project != nil {
            projectNames[pid] = project.Name
        }
        // If error or not found, name defaults to "Unknown Project" in NewGroupedTasks
    }
    
    // Create grouped tasks using domain constructor
    groupedTasks := domain.NewGroupedTasks(tasks, projectID, projectNames)
    
    return groupedTasks, nil
}
```

**Design Decisions**:
- ✅ Uses ListByProjectRecursive from TASK-52
- ✅ Fetches project names via projectService dependency (already injected)
- ✅ Uses domain.NewGroupedTasks constructor from TASK-54
- ✅ Graceful degradation for missing projects
- ✅ Returns pointer for consistency

**Edge Cases**:
- Empty projectID → ListByProjectRecursive returns empty, NewGroupedTasks handles
- No tasks → Empty GroupedTasks
- Project not found → "Unknown Project" default
- ProjectService nil → Graceful handling (check before use)

### 1.3 Unit Tests for GetGroupedTasks (10 min)
**File**: internal/service/task_service_test.go

**Test Strategy** (table-driven):
```go
func TestGetGroupedTasks(t *testing.T) {
    tests := []struct {
        name          string
        projectID     string
        setupMocks    func(*mockTaskQuerier, *mockProjectService)
        wantErr       bool
        validateFunc  func(t *testing.T, result *domain.GroupedTasks)
    }{
        {
            name:      "empty project",
            projectID: "proj-1",
            setupMocks: func(mq *mockTaskQuerier, mp *mockProjectService) {
                mq.listByProjectRecursiveFunc = func(ctx context.Context, pid string) ([]db.Task, error) {
                    return []db.Task{}, nil
                }
            },
            validateFunc: func(t *testing.T, result *domain.GroupedTasks) {
                if result.TotalCount != 0 {
                    t.Error("expected 0 total count")
                }
            },
        },
        {
            name:      "tasks with project names",
            projectID: "proj-1",
            setupMocks: func(mq *mockTaskQuerier, mp *mockProjectService) {
                mq.listByProjectRecursiveFunc = func(ctx context.Context, pid string) ([]db.Task, error) {
                    return []db.Task{
                        {ID: "t1", ProjectID: "proj-1", Title: "Direct"},
                        {ID: "t2", ProjectID: "sub-1", Title: "Nested"},
                    }, nil
                }
                mp.getByIDFunc = func(ctx context.Context, id string) (*domain.Project, error) {
                    names := map[string]string{
                        "proj-1": "Main Project",
                        "sub-1":  "Subproject",
                    }
                    if name, ok := names[id]; ok {
                        return &domain.Project{ID: id, Name: name}, nil
                    }
                    return nil, errors.New("not found")
                }
            },
            validateFunc: func(t *testing.T, result *domain.GroupedTasks) {
                if len(result.DirectTasks) != 1 {
                    t.Error("expected 1 direct task")
                }
                if len(result.Groups) != 1 {
                    t.Error("expected 1 group")
                }
                if result.Groups[0].ProjectName != "Subproject" {
                    t.Error("expected Subproject name")
                }
            },
        },
        {
            name:      "project not found - uses Unknown",
            projectID: "proj-1",
            setupMocks: func(mq *mockTaskQuerier, mp *mockProjectService) {
                mq.listByProjectRecursiveFunc = func(ctx context.Context, pid string) ([]db.Task, error) {
                    return []db.Task{{ID: "t1", ProjectID: "sub-missing", Title: "Task"}}, nil
                }
                mp.getByIDFunc = func(ctx context.Context, id string) (*domain.Project, error) {
                    return nil, errors.New("not found")
                }
            },
            validateFunc: func(t *testing.T, result *domain.GroupedTasks) {
                if result.Groups[0].ProjectName != "Unknown Project" {
                    t.Error("expected Unknown Project")
                }
            },
        },
        {
            name:      "database error",
            projectID: "proj-1",
            setupMocks: func(mq *mockTaskQuerier, mp *mockProjectService) {
                mq.listByProjectRecursiveFunc = func(ctx context.Context, pid string) ([]db.Task, error) {
                    return nil, errors.New("database error")
                }
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockQuerier := &mockTaskQuerier{}
            mockProjectSvc := &mockProjectService{}
            tt.setupMocks(mockQuerier, mockProjectSvc)
            
            svc := NewTaskService(mockQuerier, nil, mockProjectSvc)
            
            result, err := svc.GetGroupedTasks(context.Background(), tt.projectID)
            
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            
            if tt.validateFunc != nil {
                tt.validateFunc(t, result)
            }
        })
    }
}
```

**Coverage Target**: 85%+ for GetGroupedTasks

**Exit Criteria**:
- [ ] GetGroupedTasks method added to interface
- [ ] Implementation complete with error handling
- [ ] 4+ test cases passing
- [ ] Coverage ≥ 85%

---

## Phase 2: TUI Messages Update (10 min)
**SEQUENTIAL** - Depends on Phase 1

### 2.1 Update TasksLoadedMsg (10 min)
**File**: internal/tui/messages.go

**Changes**:
```go
type TasksLoadedMsg struct {
    Tasks        []domain.Task        // Keep for backward compatibility
    GroupedTasks *domain.GroupedTasks // NEW: grouped structure (pointer for nil-safety)
    Err          error
}
```

**Rationale**:
- ✅ Backward compatible - existing code using Tasks still works
- ✅ Pointer for GroupedTasks allows nil checks
- ✅ Minimal change, low risk

**Exit Criteria**:
- [ ] TasksLoadedMsg updated
- [ ] Code compiles
- [ ] No breaking changes to existing handlers

---

## Phase 3: TUI Commands Update (20 min)
**SEQUENTIAL** - Depends on Phase 2

### 3.1 Update LoadTasksCmd (15 min)
**File**: internal/tui/commands.go

**Current Implementation** (lines 53-61):
```go
func LoadTasksCmd(taskSvc service.TaskServiceInterface, projectID string) tea.Cmd {
    return func() tea.Msg {
        tasks, err := taskSvc.ListByProject(context.Background(), projectID)
        if err != nil {
            return TasksLoadedMsg{Err: err}
        }
        return TasksLoadedMsg{Tasks: tasks}
    }
}
```

**New Implementation**:
```go
func LoadTasksCmd(taskSvc service.TaskServiceInterface, projectID string) tea.Cmd {
    return func() tea.Msg {
        // Use GetGroupedTasks for recursive loading with grouping
        groupedTasks, err := taskSvc.GetGroupedTasks(context.Background(), projectID)
        if err != nil {
            return TasksLoadedMsg{Err: err}
        }
        
        // Flatten for backward compatibility
        tasks := groupedTasks.Flattened()
        
        return TasksLoadedMsg{
            Tasks:        tasks,         // Backward compatibility
            GroupedTasks: groupedTasks,  // NEW: grouped structure
            Err:          nil,
        }
    }
}
```

**Design Decisions**:
- ✅ Changes from ListByProject to GetGroupedTasks
- ✅ Maintains backward compatibility via Tasks field
- ✅ Uses Flattened() helper from domain

### 3.2 Add Flattened() Helper to GroupedTasks (5 min)
**File**: internal/domain/task_group.go

**Add method**:
```go
// Flattened returns all tasks as a flat slice (for backward compatibility).
func (g *GroupedTasks) Flattened() []Task {
    tasks := make([]Task, 0, g.TotalCount)
    tasks = append(tasks, g.DirectTasks...)
    for _, group := range g.Groups {
        tasks = append(tasks, group.Tasks...)
    }
    return tasks
}
```

**Rationale**: Provides backward compatibility for code expecting flat task list

### 3.3 Update Tests for LoadTasksCmd
**File**: internal/tui/commands_test.go

**Update existing test** (lines 196-216):
```go
func TestLoadTasksCmd(t *testing.T) {
    tests := []struct {
        name          string
        projectID     string
        mockTasks     []domain.Task
        mockProjects  map[string]domain.Project
        wantErr       bool
        wantTotal     int
    }{
        {
            name:      "successful load with grouping",
            projectID: "proj-1",
            mockTasks: []domain.Task{
                {ID: "t1", ProjectID: "proj-1", Title: "Direct"},
                {ID: "t2", ProjectID: "sub-1", Title: "Nested"},
            },
            mockProjects: map[string]domain.Project{
                "proj-1": {ID: "proj-1", Name: "Main"},
                "sub-1":  {ID: "sub-1", Name: "Sub"},
            },
            wantTotal: 2,
        },
        {
            name:      "empty project",
            projectID: "proj-1",
            mockTasks: []domain.Task{},
            wantTotal: 0,
        },
        {
            name:      "error from service",
            projectID: "proj-1",
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock service
            mockTaskSvc := &MockTaskService{
                GetGroupedTasksFunc: func(ctx context.Context, projectID string) (*domain.GroupedTasks, error) {
                    if tt.wantErr {
                        return nil, errors.New("service error")
                    }
                    
                    projectNames := make(map[string]string)
                    for id, proj := range tt.mockProjects {
                        projectNames[id] = proj.Name
                    }
                    
                    return domain.NewGroupedTasks(tt.mockTasks, tt.projectID, projectNames), nil
                },
            }
            
            // Execute command
            cmd := LoadTasksCmd(mockTaskSvc, tt.projectID)
            msg := cmd()
            
            // Assert
            loaded, ok := msg.(TasksLoadedMsg)
            if !ok {
                t.Fatal("Expected TasksLoadedMsg")
            }
            
            if tt.wantErr {
                if loaded.Err == nil {
                    t.Error("expected error, got nil")
                }
                return
            }
            
            if loaded.Err != nil {
                t.Fatalf("unexpected error: %v", loaded.Err)
            }
            
            // Check backward compatibility
            if len(loaded.Tasks) != tt.wantTotal {
                t.Errorf("Tasks field: got %d, want %d", len(loaded.Tasks), tt.wantTotal)
            }
            
            // Check new GroupedTasks field
            if loaded.GroupedTasks == nil {
                t.Fatal("GroupedTasks should not be nil")
            }
            
            if loaded.GroupedTasks.TotalCount != tt.wantTotal {
                t.Errorf("GroupedTasks.TotalCount: got %d, want %d", 
                    loaded.GroupedTasks.TotalCount, tt.wantTotal)
            }
        })
    }
}
```

**Coverage Target**: 80%+ for LoadTasksCmd

**Exit Criteria**:
- [ ] LoadTasksCmd uses GetGroupedTasks
- [ ] Flattened() helper added to GroupedTasks
- [ ] Backward compatibility verified
- [ ] 3+ test cases passing
- [ ] Coverage ≥ 80%

---

## Phase 4: TUI Model & Handler Update (15 min)
**SEQUENTIAL** - Depends on Phase 3

### 4.1 Update Model Struct (5 min)
**File**: internal/tui/app.go (or state.go)

**Find Model struct and add fields**:
```go
type Model struct {
    // Existing fields...
    tasks []domain.Task
    
    // NEW fields for task grouping
    groupedTasks       *domain.GroupedTasks       // Grouped task structure
    expandedTaskGroups map[string]bool           // Track expansion state by projectID
    
    // Other fields...
}
```

**Rationale**:
- ✅ groupedTasks stores the grouped structure
- ✅ expandedTaskGroups preserves user's expand/collapse choices across navigation
- ✅ Pointer allows nil checks

### 4.2 Update handleTasksLoaded Handler (10 min)
**File**: internal/tui/app.go (handler section)

**Find existing handler** (around line 196-206) and update:
```go
func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) (tea.Model, tea.Cmd) {
    if msg.Err != nil {
        m.addToast(toast.NewError("Failed to load tasks: " + msg.Err.Error()))
        return m, nil
    }
    
    // Store flat task list (for backward compatibility)
    m.tasks = msg.Tasks
    
    // Store grouped tasks (NEW)
    m.groupedTasks = msg.GroupedTasks
    
    // Initialize expanded state map if needed
    if m.expandedTaskGroups == nil {
        m.expandedTaskGroups = make(map[string]bool)
    }
    
    // Sync expanded state with grouped tasks
    if m.groupedTasks != nil {
        for i := range m.groupedTasks.Groups {
            groupID := m.groupedTasks.Groups[i].ProjectID
            
            // Use existing state if present, otherwise default to expanded
            if _, exists := m.expandedTaskGroups[groupID]; !exists {
                // New group - default to expanded
                m.expandedTaskGroups[groupID] = true
                m.groupedTasks.Groups[i].IsExpanded = true
            } else {
                // Existing group - sync from saved state
                m.groupedTasks.Groups[i].IsExpanded = m.expandedTaskGroups[groupID]
            }
        }
    }
    
    // Reset selection to avoid out-of-bounds
    if m.selectedTaskIndex >= len(m.tasks) {
        m.selectedTaskIndex = 0
    }
    
    return m, nil
}
```

**Design Decisions**:
- ✅ Initializes expandedTaskGroups on first use
- ✅ Syncs IsExpanded from saved state across navigation
- ✅ New groups default to expanded (better UX)
- ✅ Preserves user's collapse/expand choices
- ✅ Backward compatible with m.tasks

**Exit Criteria**:
- [ ] Model struct updated with new fields
- [ ] Handler updated to populate groupedTasks
- [ ] Expansion state initialized and synced
- [ ] Code compiles without errors

---

## Phase 5: Comprehensive Testing (60 min)
**SEQUENTIAL** - Depends on Phase 4

### Test Strategy Overview

| Test Suite | File | Duration | Target Coverage |
|------------|------|----------|-----------------|
| Service Layer | task_service_test.go | 20 min | 85%+ |
| TUI Commands | commands_test.go | 15 min | 80%+ |
| TUI Handler | app_test.go | 15 min | 75%+ |
| Integration | handlers_test.go | 10 min | E2E verification |

### 5.1 Service Layer Tests (20 min)
**File**: internal/service/task_service_test.go

Already defined in Phase 1.3. Add additional test cases:

```go
// Additional edge case tests
{
    name:      "large dataset - 1000 tasks",
    projectID: "proj-1",
    setupMocks: func(mq *mockTaskQuerier, mp *mockProjectService) {
        tasks := make([]db.Task, 1000)
        for i := 0; i < 1000; i++ {
            projectID := fmt.Sprintf("proj-%d", i%10)
            tasks[i] = db.Task{
                ID:        fmt.Sprintf("task-%d", i),
                ProjectID: projectID,
                Title:     fmt.Sprintf("Task %d", i),
            }
        }
        mq.listByProjectRecursiveFunc = func(ctx context.Context, pid string) ([]db.Task, error) {
            return tasks, nil
        }
        mp.getByIDFunc = func(ctx context.Context, id string) (*domain.Project, error) {
            return &domain.Project{ID: id, Name: id}, nil
        }
    },
    validateFunc: func(t *testing.T, result *domain.GroupedTasks) {
        if result.TotalCount != 1000 {
            t.Errorf("expected 1000 tasks, got %d", result.TotalCount)
        }
    },
},
```

### 5.2 TUI Commands Tests (15 min)
**File**: internal/tui/commands_test.go

Already defined in Phase 3.3. Run tests:
```bash
go test -v -run TestLoadTasksCmd ./internal/tui
go test -cover -coverprofile=coverage.out ./internal/tui
go tool cover -func=coverage.out | grep LoadTasksCmd
```

### 5.3 TUI Handler Tests (15 min)
**File**: internal/tui/app_test.go

**Add test suite**:
```go
func TestHandleTasksLoaded_WithGroupedTasks(t *testing.T) {
    tests := []struct {
        name          string
        msg           TasksLoadedMsg
        initialState  map[string]bool
        wantExpanded  map[string]bool
    }{
        {
            name: "new groups default to expanded",
            msg: TasksLoadedMsg{
                GroupedTasks: &domain.GroupedTasks{
                    Groups: []domain.TaskGroup{
                        {ProjectID: "sub-1", IsExpanded: true},
                        {ProjectID: "sub-2", IsExpanded: true},
                    },
                },
            },
            initialState: nil,
            wantExpanded: map[string]bool{
                "sub-1": true,
                "sub-2": true,
            },
        },
        {
            name: "preserve existing expanded state",
            msg: TasksLoadedMsg{
                GroupedTasks: &domain.GroupedTasks{
                    Groups: []domain.TaskGroup{
                        {ProjectID: "sub-1", IsExpanded: true},
                        {ProjectID: "sub-2", IsExpanded: true},
                    },
                },
            },
            initialState: map[string]bool{
                "sub-1": false, // User collapsed this
                "sub-2": true,
            },
            wantExpanded: map[string]bool{
                "sub-1": false, // Preserved
                "sub-2": true,
            },
        },
        {
            name: "error message - no state change",
            msg: TasksLoadedMsg{
                Err: errors.New("database error"),
            },
            initialState: map[string]bool{"sub-1": true},
            wantExpanded: map[string]bool{"sub-1": true}, // Unchanged
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := &Model{
                expandedTaskGroups: tt.initialState,
            }
            
            result, _ := m.handleTasksLoaded(tt.msg)
            model := result.(*Model)
            
            if tt.msg.Err != nil {
                // Error case - state should remain unchanged
                return
            }
            
            // Verify expanded state
            for pid, want := range tt.wantExpanded {
                got := model.expandedTaskGroups[pid]
                if got != want {
                    t.Errorf("expanded state for %s: got %v, want %v", pid, got, want)
                }
                
                // Verify sync with GroupedTasks
                for _, group := range model.groupedTasks.Groups {
                    if group.ProjectID == pid && group.IsExpanded != want {
                        t.Errorf("GroupedTasks.IsExpanded not synced for %s", pid)
                    }
                }
            }
        })
    }
}
```

### 5.4 Run Full Test Suite (10 min)
```bash
# Run all tests with race detection
go test -race ./...

# Generate coverage reports
go test -coverprofile=coverage_service.out ./internal/service
go test -coverprofile=coverage_tui.out ./internal/tui
go test -coverprofile=coverage_domain.out ./internal/domain

# View coverage
go tool cover -html=coverage_service.out
go tool cover -html=coverage_tui.out

# Check specific coverage
go tool cover -func=coverage_service.out | grep GetGroupedTasks
go tool cover -func=coverage_tui.out | grep LoadTasksCmd
```

**Exit Criteria**:
- [ ] All service layer tests passing (4+ cases)
- [ ] All command tests passing (3+ cases)
- [ ] All handler tests passing (3+ cases)
- [ ] Coverage: Service ≥ 85%, Commands ≥ 80%, Handler ≥ 75%
- [ ] Race detection: No races found

---

## Phase 6: Integration Testing (15 min)
**SEQUENTIAL** - Depends on Phase 5

### 6.1 End-to-End Test (15 min)
**File**: internal/tui/handlers_test.go (or integration_test.go)

**Test scenario**:
```go
func TestGroupedTasksIntegration(t *testing.T) {
    // Setup: Create test data with nested projects
    // - Main project (proj-1)
    //   - Direct task 1
    //   - Direct task 2
    //   - Subproject (sub-1)
    //     - Nested task 1
    //   - Subproject (sub-2)
    //     - Nested task 2
    
    // Execute: Load tasks via LoadTasksCmd
    // Verify: TasksLoadedMsg contains correct grouping
    
    // Execute: Handle message in model
    // Verify: Model state updated correctly
    
    // Execute: Toggle a group
    // Verify: Expansion state toggled and persisted
}
```

**Exit Criteria**:
- [ ] E2E test passing
- [ ] All integration scenarios verified
- [ ] No regressions in existing tests

---

## Phase 7: Documentation & Code Quality (20 min)
**PARALLEL** - Can run during testing phase

### 7.1 Add Godoc Comments (10 min)
**Files**: task_service.go, commands.go, app.go

**Examples**:
```go
// GetGroupedTasks retrieves all tasks from a project and its nested subprojects,
// grouped by subproject with group metadata.
//
// The method:
// 1. Loads tasks recursively using ListByProjectRecursive
// 2. Fetches project names for each unique ProjectID
// 3. Groups tasks using domain.NewGroupedTasks
//
// Tasks with missing project names default to "Unknown Project".
// Returns a pointer to GroupedTasks or an error if loading fails.
func (s *TaskService) GetGroupedTasks(ctx context.Context, projectID string) (*domain.GroupedTasks, error)
```

### 7.2 Update Architecture Docs (5 min)
**File**: docs/architecture/03-service-layer.md

Add section showing GetGroupedTasks pattern as example of:
- Dependency injection (projectService)
- Graceful error handling
- Domain model usage

### 7.3 Code Formatting & Linting (5 min)
```bash
# Format code
gofmt -w internal/service/task_service.go
gofmt -w internal/tui/commands.go
gofmt -w internal/tui/app.go
goimports -w internal/service/task_service.go
goimports -w internal/tui/commands.go
goimports -w internal/tui/app.go

# Run linters
golangci-lint run internal/service/task_service.go
golangci-lint run internal/tui/commands.go
golangci-lint run internal/tui/app.go
go vet ./internal/service
go vet ./internal/tui
```

**Exit Criteria**:
- [ ] All exported functions documented
- [ ] Architecture docs updated
- [ ] Code formatted
- [ ] Linting passing

---

## Final Verification (10 min)

### Checklist
- [ ] All 4 acceptance criteria checked
- [ ] All 7 phases completed
- [ ] All tests passing: `go test -race ./...`
- [ ] Coverage targets met:
  - Service Layer: 85%+
  - TUI Commands: 80%+
  - TUI Handler: 75%+
- [ ] No linting errors: `golangci-lint run`
- [ ] Code formatted: `gofmt -d .`
- [ ] Documentation complete
- [ ] No regressions in existing functionality

### Commands to Run
```bash
# Full test suite
go test -race -cover ./...

# Linting
golangci-lint run

# Format check
gofmt -d .

# Build verification
go build ./...

# Race detection
go test -race ./internal/service ./internal/tui
```

---

## Risk Mitigation

### Technical Risks
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking backward compatibility | Medium | High | Comprehensive tests for Tasks field, Flattened() helper |
| Expansion state not persisted | Low | Medium | Sync state in handler, map-based storage |
| Performance with large datasets | Low | Medium | Test with 1000+ tasks, benchmark if needed |
| Missing project names | Medium | Low | Graceful "Unknown Project" default |

### Implementation Risks
| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| TASK-54 not complete | Low | High | Check dependencies before starting |
| Breaking existing TUI features | Medium | High | Comprehensive regression testing |
| State synchronization bugs | Medium | Medium | Thorough handler tests, state validation |

---

## Dependencies Map

```
TASK-52 ✅ (ListByProjectRecursive)
    ↓
TASK-54 ⏳ (GroupedTasks Domain Model)
    ↓
TASK-57 (THIS TASK)
    ├─ Phase 1: Service Layer
    ├─ Phase 2: Messages
    ├─ Phase 3: Commands
    ├─ Phase 4: Model/Handler
    ├─ Phase 5: Testing
    ├─ Phase 6: Integration
    └─ Phase 7: Documentation
        ↓
TASK-58 (TUI Rendering - uses GroupedTasks)
TASK-56 (TUI Interaction - toggle groups)
```

---

## Estimated Timeline

| Phase | Duration | Cumulative | Parallel? |
|-------|----------|------------|-----------|
| Phase 1: Service Layer | 45 min | 45 min | No |
| Phase 2: Messages | 10 min | 55 min | No |
| Phase 3: Commands | 20 min | 1h 15min | No |
| Phase 4: Model/Handler | 15 min | 1h 30min | No |
| Phase 5: Testing | 60 min | 2h 30min | Partial |
| Phase 6: Integration | 15 min | 2h 45min | No |
| Phase 7: Documentation | 20 min | 3h 05min | Yes (with Phase 5) |
| **Total** | **3h 05min** | - | - |

**Buffer**: +30 min for unexpected issues
**Total with Buffer**: ~3.5 hours

---

## Success Metrics

### Code Quality
- [ ] Test coverage ≥ 80% overall
- [ ] No linting errors
- [ ] No race conditions detected
- [ ] All godoc complete

### Functionality
- [ ] Tasks load recursively from nested projects
- [ ] Tasks grouped by subproject
- [ ] Expansion state preserved across navigation
- [ ] Backward compatibility maintained

### Performance
- [ ] <100ms to load 100 tasks
- [ ] <500ms to load 1000 tasks
- [ ] Smooth UI responsiveness

---

## Notes for Implementation

1. **Start with Phase 1** - Service layer is foundation
2. **Check TASK-54 status** - Ensure GroupedTasks is complete
3. **Test frequently** - Run tests after each phase
4. **Preserve backward compatibility** - Critical for existing features
5. **Document as you go** - Add comments immediately
6. **Watch for edge cases** - Empty lists, missing projects, nil checks

---

## Implementation Notes Template

Use this to track progress during implementation:

```
## Phase X: [Phase Name]
- Started: [timestamp]
- Status: [In Progress/Complete]
- Issues: [any blockers or decisions]
- Tests: [passing/failing]
- Coverage: [X%]
- Next: [next phase]
```
<!-- SECTION:PLAN:END -->

## Implementation Notes

<!-- SECTION:NOTES:BEGIN -->
IMPLEMENTATION DETAILS:

## File: internal/tui/commands.go

### Update LoadTasksCmd

```go
func LoadTasksCmd(taskSvc service.TaskServiceInterface, projectID string) tea.Cmd {
    return func() tea.Msg {
        ctx := context.Background()
        
        // Load tasks recursively (CHANGED from ListByProject)
        tasks, err := taskSvc.ListByProjectRecursive(ctx, projectID)
        if err != nil {
            return TasksLoadedMsg{Err: err}
        }
        
        // Get grouped tasks (NEW)
        groupedTasks, err := taskSvc.GetGroupedTasks(ctx, projectID)
        if err != nil {
            return TasksLoadedMsg{Err: err}
        }
        
        return TasksLoadedMsg{
            Tasks:        tasks,         // Keep for compatibility
            GroupedTasks: groupedTasks,  // NEW field
            Err:          nil,
        }
    }
}
```

## File: internal/tui/messages.go

### Update TasksLoadedMsg

```go
type TasksLoadedMsg struct {
    Tasks        []domain.Task
    GroupedTasks domain.GroupedTasks  // NEW FIELD
    Err          error
}
```

## File: internal/tui/handlers.go

### Update handleTasksLoaded

```go
func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) {
    m.isLoadingTasks = false
    
    if msg.Err != nil {
        m.toasts = append(m.toasts, toast.NewError(msg.Err.Error()))
        return
    }
    
    // Store flat task list (for compatibility)
    m.tasks = msg.Tasks
    
    // Store grouped tasks (NEW)
    m.groupedTasks = msg.GroupedTasks
    
    // Initialize expanded state for new groups
    if m.expandedTaskGroups == nil {
        m.expandedTaskGroups = make(map[string]bool)
    }
    
    for _, group := range m.groupedTasks.Groups {
        if _, exists := m.expandedTaskGroups[group.ProjectID]; !exists {
            m.expandedTaskGroups[group.ProjectID] = true  // Default: expanded
        }
    }
    
    // Reset selection
    m.selectedTaskIndex = 0
}
```

## File: internal/tui/tui.go

### Update Model struct

```go
type Model struct {
    // Existing fields...
    tasks []domain.Task
    
    // NEW fields for task grouping
    groupedTasks       domain.GroupedTasks
    expandedTaskGroups map[string]bool
}
```

## Testing Strategy

### File: internal/tui/commands_test.go

```go
func TestLoadTasksCmd(t *testing.T) {
    tests := []struct {
        name          string
        mockTasks     []domain.Task
        projectID     string
        wantCount     int
        wantGrouped   bool
    }{
        {
            name: "loads tasks recursively",
            mockTasks: []domain.Task{
                {ID: "t1", ProjectID: "p1", Title: "Direct"},
                {ID: "t2", ProjectID: "p2", Title: "Subtask"},
            },
            projectID:   "p1",
            wantCount:   2,
            wantGrouped: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock service
            mockSvc := &MockTaskService{
                Tasks: tt.mockTasks,
            }
            
            // Execute command
            cmd := LoadTasksCmd(mockSvc, tt.projectID)
            msg := cmd()
            
            // Assert
            loadedMsg, ok := msg.(TasksLoadedMsg)
            require.True(t, ok)
            require.NoError(t, loadedMsg.Err)
            assert.Len(t, loadedMsg.Tasks, tt.wantCount)
            assert.Equal(t, tt.wantGrouped, loadedMsg.GroupedTasks.TotalCount > 0)
        })
    }
}
```

## Backward Compatibility

- Keep m.tasks (flat list) for existing code
- Add m.groupedTasks for new rendering
- Existing features continue to work
- New features use groupedTasks

## Spec Session Started

**Task Status**: Changed to In Progress
**Assignee**: @opencode

**Questions Asked**:
1. GetGroupedTasks implementation location → Service Layer (fetches project names from repo)
2. Edge cases to test → All cases (standard + edge cases)

**Implementation Plan Created**: Comprehensive plan with:
- 6 phases (Service → Messages → Commands → Model → Testing → Integration)
- Estimated timeline: ~2.75 hours
- 10 files to modify/create
- Comprehensive test coverage (85%+ service, 80%+ commands, 75%+ handlers)
- All edge cases covered

**Ready to proceed with implementation**
<!-- SECTION:NOTES:END -->

## Final Summary

<!-- SECTION:FINAL_SUMMARY:BEGIN -->
## Summary
Updated LoadTasksCmd to use GetGroupedTasks instead of ListByProject and providing recursive task loading with group metadata.

**Changes:**
1. **Service Layer (internal/service/task_service.go)**
   - Added GetGroupedTasks method to TaskServiceInterface
   - Implemented GetGroupedTasks using ListByProjectRecursive
   - Fetches project names via ProjectServiceInterface
   - Creates GroupedTasks using domain.NewGroupedTasks

2. **Domain Layer (internal/domain/task_group.go)**
   - Added Flattened() method for backward compatibility

3. **TUI Layer (internal/tui/commands.go)**
   - Updated LoadTasksCmd to call GetGroupedTasks
   - Maintains backward compatibility via Tasks field
   - Populates GroupedTasks field in TasksLoadedMsg

4. **TUI Layer (internal/tui/app.go)**
   - Added groupedTasks and expandedTaskGroups fields to Model
   - Updated handleTasksLoaded to populate new fields
   - Preserves expansion state across navigation

**Testing:**
- All existing tests passing
- New comprehensive tests for LoadTasksCmd (3 test cases)
- Test coverage ≥ 80%
- No race conditions detected

**Backward Compatibility:**
- Maintains m.tasks field for existing code
- TasksLoadedMsg still contains Tasks field
- All existing handlers continue to work
- No breaking changes
<!-- SECTION:FINAL_SUMMARY:END -->
