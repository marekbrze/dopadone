# Domain Layer

## Overview

The Domain Layer is the **core** of Dopadone, containing business entities, value objects, and validation logic. It follows **Domain-Driven Design (DDD)** principles with rich domain models that enforce business invariants through factory methods.

**Key Characteristics**:
- **Zero dependencies**: No imports from other project packages
- **Rich validation**: Business rules enforced in factory methods
- **Immutable creation**: Entities created through validated constructors
- **Value objects**: Encapsulated domain concepts with validation

## Entities

### Task Entity

The Task entity represents an individual work item within a project.

```go
// internal/domain/task.go
type Task struct {
    ID                string
    ProjectID         string
    Title             string
    Description       string
    StartDate         *time.Time
    Deadline          *time.Time
    Priority          TaskPriority
    Context           string
    EstimatedDuration TaskDuration
    Status            TaskStatus
    IsNext            bool
    CreatedAt         time.Time
    UpdatedAt         time.Time
    DeletedAt         *time.Time
}
```

**Factory Method**:
```go
func NewTask(params NewTaskParams) (*Task, error) {
    // Validation 1: Required fields
    if params.Title == "" {
        return nil, ErrTaskTitleEmpty
    }
    
    if params.ProjectID == "" {
        return nil, ErrTaskProjectIDEmpty
    }
    
    // Validation 2: Value object validation
    if !params.Status.IsValid() {
        return nil, ErrTaskInvalidStatus
    }
    
    if !params.Priority.IsValid() {
        return nil, ErrTaskInvalidPriority
    }
    
    // Validation 3: Business rules
    if params.Deadline != nil && params.StartDate == nil {
        return nil, ErrTaskDeadlineNoStart
    }
    
    if params.StartDate != nil && params.Deadline != nil {
        if !params.StartDate.Before(*params.Deadline) {
            return nil, ErrTaskInvalidDateRange
        }
    }
    
    // Create validated entity
    now := time.Now()
    return &Task{
        ID:          uuid.New().String(),
        Title:       params.Title,
        ProjectID:   params.ProjectID,
        Status:      params.Status,
        Priority:    params.Priority,
        CreatedAt:   now,
        UpdatedAt:   now,
        // ... other fields
    }, nil
}
```

**Validation Rules**:
- ✅ Title must not be empty
- ✅ ProjectID must not be empty
- ✅ Status must be a valid enum value
- ✅ Priority must be a valid enum value
- ✅ Deadline requires a start date
- ✅ Start date must be before deadline

---

### Project Entity

The Project entity represents a project that can be nested within a hierarchy.

```go
// internal/domain/project.go
type Project struct {
    ID          string
    Name        string
    Description string
    Goal        string
    Status      ProjectStatus
    Priority    Priority
    Progress    Progress
    StartDate   *time.Time
    Deadline    *time.Time
    Color       Color
    ParentID    *string  // For nested projects
    SubareaID   *string  // For root projects
    Position    int
    CreatedAt   time.Time
    UpdatedAt   time.Time
    CompletedAt *time.Time
    DeletedAt   *time.Time
}
```

**Factory Method**:
```go
func NewProject(params NewProjectParams) (*Project, error) {
    // Validation 1: Required fields
    if params.Name == "" {
        return nil, ErrProjectNameEmpty
    }
    
    // Validation 2: Value objects
    if !params.Status.IsValid() {
        return nil, ErrProjectInvalidStatus
    }
    
    if !params.Priority.IsValid() {
        return nil, ErrProjectInvalidPriority
    }
    
    if !params.Progress.IsValid() {
        return nil, ErrProjectInvalidProgress
    }
    
    // Validation 3: Hierarchy constraint
    if params.ParentID == nil && params.SubareaID == nil {
        return nil, ErrProjectNoParent
    }
    
    // Validation 4: Date logic
    if params.StartDate != nil && params.Deadline != nil {
        if !params.StartDate.Before(*params.Deadline) {
            return nil, ErrProjectInvalidDateRange
        }
    }
    
    now := time.Now()
    return &Project{
        ID:        uuid.New().String(),
        Name:      params.Name,
        Status:    params.Status,
        Priority:  params.Priority,
        Progress:  params.Progress,
        CreatedAt: now,
        UpdatedAt: now,
        // ... other fields
    }, nil
}
```

**Key Business Rules**:
- A project must have **either** `ParentID` **or** `SubareaID` (never both, never neither)
- Progress must be between 0 and 100
- Color must be valid hex format (#RRGGBB)

---

### Area & Subarea Entities

**Area**: Top-level organizational container
```go
type Area struct {
    ID        string
    Name      string
    Color     Color
    SortOrder int
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

**Subarea**: Subdivision within an area
```go
type Subarea struct {
    ID        string
    Name      string
    AreaID    string
    Color     Color
    SortOrder int
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

---

## Domain Structures

Domain structures organize entities for specific use cases, such as grouping tasks for hierarchical display.

### GroupedTasks Structure

GroupedTasks organizes tasks by subproject with group metadata for TUI rendering:

**Purpose**: Display tasks hierarchically when a parent project is selected, separating direct tasks from nested subproject tasks.

```go
// internal/domain/task_group.go

// TaskGroup represents a collection of tasks belonging to a subproject
type TaskGroup struct {
    ProjectID   string  // Unique identifier of the subproject
    ProjectName string  // Display name of the subproject
    Tasks       []Task  // Tasks belonging to this subproject
    IsExpanded  bool    // Expansion state for TUI (default: true)
}

// GroupedTasks organizes tasks by subproject with group metadata
type GroupedTasks struct {
    DirectTasks []Task      // Tasks belonging to the selected project
    Groups      []TaskGroup // Tasks from nested subprojects
    TotalCount  int         // Total tasks across all groups
}
```

**Constructor Pattern**:

```go
func NewGroupedTasks(tasks []Task, parentProjectID string, projectNames map[string]string) *GroupedTasks {
    // Handles nil/empty gracefully
    if tasks == nil {
        tasks = []Task{}
    }
    if projectNames == nil {
        projectNames = make(map[string]string)
    }
    
    grouped := &GroupedTasks{
        DirectTasks: []Task{},
        Groups:      []TaskGroup{},
    }
    
    // Group tasks by project while preserving order
    // ... grouping logic
    
    return grouped
}
```

**Mutation Methods**:

```go
// AddTask adds a task to the appropriate group or DirectTasks
func (g *GroupedTasks) AddTask(task Task)

// RemoveTask removes a task by ID, returns false if not found
func (g *GroupedTasks) RemoveTask(taskID string) bool

// ToggleGroup toggles expansion state, returns false if not found
func (g *GroupedTasks) ToggleGroup(projectID string) bool

// Clear resets all fields to empty/zero
func (g *GroupedTasks) Clear()
```

**Design Decisions**:
- ✅ **Mutable structure**: Supports CRUD operations after construction
- ✅ **Graceful handling**: No errors returned, handles nil/empty inputs safely
- ✅ **Order preservation**: Tasks maintain append order, groups maintain discovery order
- ✅ **Default expanded**: All groups default to `IsExpanded=true` for better UX
- ✅ **Separation of concerns**: Direct tasks vs. nested tasks clearly separated

**Use Case Example**:

```go
// In TUI when a project is selected
func (m *Model) loadTasks(projectID string) {
    // Get all tasks recursively
    tasks, _ := m.taskService.ListByProjectRecursive(ctx, projectID)
    
    // Get project names for display
    projectNames := m.getProjectNames()
    
    // Group tasks for display
    m.groupedTasks = domain.NewGroupedTasks(tasks, projectID, projectNames)
}
```

**Testing Pattern**:

```go
func TestNewGroupedTasks(t *testing.T) {
    tests := []struct {
        name            string
        tasks           []domain.Task
        parentProjectID string
        projectNames    map[string]string
        wantDirectCount int
        wantGroupCount  int
        wantTotalCount  int
    }{
        {
            name: "direct tasks only",
            tasks: []domain.Task{
                {ID: "t1", ProjectID: "proj-1", Title: "Task 1"},
                {ID: "t2", ProjectID: "proj-1", Title: "Task 2"},
            },
            parentProjectID: "proj-1",
            wantDirectCount: 2,
            wantGroupCount:  0,
            wantTotalCount:  2,
        },
        // ... more test cases
    }
    // ... test implementation
}
```

---

## Value Objects

Value objects encapsulate domain concepts with built-in validation.

### TaskStatus

```go
type TaskStatus string

const (
    TaskStatusTodo       TaskStatus = "todo"
    TaskStatusInProgress TaskStatus = "in_progress"
    TaskStatusWaiting    TaskStatus = "waiting"
    TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
    switch s {
    case TaskStatusTodo, TaskStatusInProgress, TaskStatusWaiting, TaskStatusDone:
        return true
    default:
        return false
    }
}

func ParseTaskStatus(s string) (TaskStatus, error) {
    status := TaskStatus(s)
    if !status.IsValid() {
        return "", ErrInvalidTaskStatus
    }
    return status, nil
}
```

### TaskPriority (and Project Priority)

```go
type Priority string

const (
    PriorityLow    Priority = "low"
    PriorityMedium Priority = "medium"
    PriorityHigh   Priority = "high"
    PriorityUrgent Priority = "urgent"
)

func (p Priority) IsValid() bool {
    switch p {
    case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
        return true
    default:
        return false
    }
}

func ParsePriority(s string) (Priority, error) {
    priority := Priority(s)
    if !priority.IsValid() {
        return "", ErrInvalidPriority
    }
    return priority, nil
}
```

### Color

```go
type Color string

func (c Color) IsValid() bool {
    if c == "" {
        return true // Optional field
    }
    matched, _ := regexp.MatchString(`^#[0-9A-Fa-f]{6}$`, string(c))
    return matched
}

func ParseColor(s string) (Color, error) {
    color := Color(s)
    if !color.IsValid() {
        return "", ErrInvalidColorFormat
    }
    return color, nil
}
```

**Valid formats**: `#FF0000`, `#00ff00`, `#0000FF`
**Invalid formats**: `FF0000`, `red`, `#FFF`

### Progress

```go
type Progress int

func (p Progress) IsValid() bool {
    return p >= 0 && p <= 100
}

func ParseProgress(n int) (Progress, error) {
    p := Progress(n)
    if !p.IsValid() {
        return 0, ErrInvalidProgress
    }
    return p, nil
}
```

### TaskDuration

```go
type TaskDuration int

const (
    Duration5    TaskDuration = 5
    Duration15   TaskDuration = 15
    Duration30   TaskDuration = 30
    Duration60   TaskDuration = 60
    Duration120  TaskDuration = 120
    Duration240  TaskDuration = 240
    Duration480  TaskDuration = 480
)

func (d TaskDuration) IsValid() bool {
    switch d {
    case Duration5, Duration15, Duration30, Duration60, Duration120, Duration240, Duration480:
        return true
    default:
        return false
    }
}
```

---

## Domain Validation

### Validation Errors

All domain validation errors are **sentinel errors** defined at package level:

```go
var (
    // Task errors
    ErrTaskTitleEmpty       = errors.New("task title cannot be empty")
    ErrTaskProjectIDEmpty   = errors.New("task project_id cannot be empty")
    ErrTaskInvalidStatus    = errors.New("task status is invalid")
    ErrTaskInvalidPriority  = errors.New("task priority is invalid")
    ErrTaskInvalidDuration  = errors.New("task estimated_duration is invalid")
    ErrTaskInvalidDateRange = errors.New("task deadline must be after start date")
    ErrTaskDeadlineNoStart  = errors.New("task deadline cannot be set without start_date")
    
    // Project errors
    ErrProjectNameEmpty        = errors.New("project name cannot be empty")
    ErrProjectInvalidStatus    = errors.New("project status is invalid")
    ErrProjectInvalidPriority  = errors.New("project priority is invalid")
    ErrProjectInvalidProgress  = errors.New("project progress must be between 0 and 100")
    ErrProjectNoParent         = errors.New("project must have either parent_id or subarea_id")
    ErrProjectInvalidDateRange = errors.New("project deadline must be after start date")
    
    // Value object errors
    ErrInvalidColorFormat   = errors.New("invalid color format: must be a valid hex color (e.g., #FF0000)")
    ErrInvalidProjectStatus = errors.New("invalid project status")
    ErrInvalidPriority      = errors.New("invalid priority")
    ErrInvalidProgress      = errors.New("invalid progress: must be between 0 and 100")
)
```

### Business Invariants

**Example 1: Date Range Validation**
```go
// A deadline cannot exist without a start date
if params.Deadline != nil && params.StartDate == nil {
    return nil, ErrTaskDeadlineNoStart
}

// Start date must be before deadline
if params.StartDate != nil && params.Deadline != nil {
    if !params.StartDate.Before(*params.Deadline) {
        return nil, ErrTaskInvalidDateRange
    }
}
```

**Example 2: Hierarchy Constraint**
```go
// A project must belong to either a parent project OR a subarea
if params.ParentID == nil && params.SubareaID == nil {
    return nil, ErrProjectNoParent
}
```

---

## Testing Domain Logic

Domain logic is tested directly without mocks since it has no dependencies.

**Example Test**:
```go
func TestNewTask_ValidParams(t *testing.T) {
    now := time.Now()
    startDate := now.Add(24 * time.Hour)
    deadline := now.Add(48 * time.Hour)
    
    task, err := domain.NewTask(domain.NewTaskParams{
        ProjectID:   "proj-123",
        Title:       "Write tests",
        Description: "Test the task entity",
        Status:      domain.TaskStatusTodo,
        Priority:    domain.PriorityHigh,
        StartDate:   &startDate,
        Deadline:    &deadline,
    })
    
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    
    if task.Title != "Write tests" {
        t.Errorf("expected title 'Write tests', got %s", task.Title)
    }
    
    if task.Status != domain.TaskStatusTodo {
        t.Errorf("expected status %s, got %s", domain.TaskStatusTodo, task.Status)
    }
}

func TestNewTask_EmptyTitle(t *testing.T) {
    _, err := domain.NewTask(domain.NewTaskParams{
        ProjectID: "proj-123",
        Title:     "",
        Status:    domain.TaskStatusTodo,
        Priority:  domain.PriorityMedium,
    })
    
    if err != domain.ErrTaskTitleEmpty {
        t.Errorf("expected ErrTaskTitleEmpty, got %v", err)
    }
}

func TestNewTask_InvalidDateRange(t *testing.T) {
    now := time.Now()
    startDate := now.Add(48 * time.Hour)
    deadline := now.Add(24 * time.Hour) // Before start date
    
    _, err := domain.NewTask(domain.NewTaskParams{
        ProjectID: "proj-123",
        Title:     "Test task",
        Status:    domain.TaskStatusTodo,
        Priority:  domain.PriorityMedium,
        StartDate: &startDate,
        Deadline:  &deadline,
    })
    
    if err != domain.ErrTaskInvalidDateRange {
        t.Errorf("expected ErrTaskInvalidDateRange, got %v", err)
    }
}
```

---

## Composite Structures

### GroupedTasks Structure

The `GroupedTasks` structure organizes tasks by subproject with group metadata, supporting nested task grouping in the TUI.

**Purpose**:
- Groups tasks from a project and its nested subprojects
- Separates direct tasks (selected project) from nested tasks (subprojects)
- Tracks expansion state for TUI rendering
- Preserves task order within groups

**Structure**:
```go
// internal/domain/task_group.go

type TaskGroup struct {
    ProjectID   string  // Unique identifier of the subproject
    ProjectName string  // Display name of the subproject
    Tasks       []Task  // Tasks belonging to this subproject
    IsExpanded  bool    // Expansion state for TUI (default: true)
}

type GroupedTasks struct {
    DirectTasks []Task      // Tasks belonging to the selected project
    Groups      []TaskGroup // Tasks from nested subprojects
    TotalCount  int         // Total tasks across all groups
    
    parentProjectID string // Private: used by AddTask to determine grouping
}
```

### Constructor Pattern

**Factory Method**:
```go
func NewGroupedTasks(tasks []Task, parentProjectID string, projectNames map[string]string) *GroupedTasks
```

**Usage**:
```go
tasks := []domain.Task{
    {ID: "t1", ProjectID: "proj-1", Title: "Direct Task"},
    {ID: "t2", ProjectID: "sub-1", Title: "Nested Task"},
}

projectNames := map[string]string{
    "proj-1": "Main Project",
    "sub-1":  "Subproject",
}

grouped := domain.NewGroupedTasks(tasks, "proj-1", projectNames)

// Result:
// grouped.DirectTasks = [{ID: "t1", ...}]
// grouped.Groups = [{ProjectID: "sub-1", ProjectName: "Subproject", Tasks: [{ID: "t2", ...}], IsExpanded: true}]
// grouped.TotalCount = 2
```

### Key Features

**1. Graceful Edge Case Handling**:
- `nil` or empty tasks → returns empty GroupedTasks
- `nil` projectNames → creates empty map
- Missing project name → defaults to "Unknown Project"
- No errors returned, always succeeds

**2. Order Preservation**:
```go
// Tasks maintain append order within groups
tasks := []Task{
    {ID: "t1", ProjectID: "sub-1", Title: "First"},
    {ID: "t2", ProjectID: "sub-1", Title: "Second"},
    {ID: "t3", ProjectID: "sub-1", Title: "Third"},
}

// Groups maintain discovery order
tasks := []Task{
    {ID: "t1", ProjectID: "sub-A", ...},
    {ID: "t2", ProjectID: "sub-B", ...},
    {ID: "t3", ProjectID: "sub-C", ...},
}
```

**3. Default Expansion State**:
- All groups default to `IsExpanded = true`
- Provides better UX by showing all tasks initially
- Users can collapse groups they don't need

### Mutation Methods

**AddTask** - Adds task to appropriate group:
```go
func (g *GroupedTasks) AddTask(task Task)
```

Behavior:
- Task.ProjectID matches parent → DirectTasks
- Group exists for ProjectID → existing group
- Otherwise → creates new group with IsExpanded=true

**RemoveTask** - Removes task by ID:
```go
func (g *GroupedTasks) RemoveTask(taskID string) bool
```

Returns:
- `true` if found and removed
- `false` if not found

**ToggleGroup** - Toggles group expansion:
```go
func (g *GroupedTasks) ToggleGroup(projectID string) bool
```

Returns:
- `true` if group found and toggled
- `false` if group not found

**Clear** - Resets to empty state:
```go
func (g *GroupedTasks) Clear()
```

### Performance Characteristics

**Time Complexity**:
| Operation | Complexity | Notes |
|-----------|------------|-------|
| NewGroupedTasks | O(n) | Single pass through tasks |
| AddTask | O(g) | g = number of groups |
| RemoveTask | O(n) | Linear search through all tasks |
| ToggleGroup | O(g) | g = number of groups |
| Clear | O(1) | Simple reset |

**Performance Targets**:
- NewGroupedTasks: <10ms for 1000 tasks
- All mutation operations: <10ms for typical datasets

### Integration with Service Layer

The `GroupedTasks` structure is created by the service layer:

```go
// internal/service/task_service.go

func (s *TaskService) GetGroupedTasks(ctx context.Context, projectID string) (*domain.GroupedTasks, error) {
    // 1. Load tasks recursively using ListByProjectRecursive
    tasks, err := s.ListByProjectRecursive(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("get grouped tasks: %w", err)
    }
    
    // 2. Build project names map by fetching from project service
    projectNames := make(map[string]string)
    for _, task := range tasks {
        if task.ProjectID != "" {
            project, err := s.projectService.GetByID(ctx, task.ProjectID)
            if err == nil && project != nil {
                projectNames[task.ProjectID] = project.Name
            }
        }
    }
    
    // 3. Create grouped structure using domain constructor
    return domain.NewGroupedTasks(tasks, projectID, projectNames), nil
}
```

### Usage in TUI Layer

**Loading Tasks**:
```go
// internal/tui/commands.go

func LoadTasksCmd(taskSvc service.TaskServiceInterface, projectID string) tea.Cmd {
    return func() tea.Msg {
        groupedTasks, err := taskSvc.GetGroupedTasks(context.Background(), projectID)
        if err != nil {
            return TasksLoadedMsg{Err: err}
        }
        
        return TasksLoadedMsg{
            Tasks:        groupedTasks.Flattened(), // Backward compatibility
            GroupedTasks: groupedTasks,             // New grouped structure
        }
    }
}
```

**State Management**:
```go
// internal/tui/app.go

type Model struct {
    tasks              []domain.Task
    groupedTasks       *domain.GroupedTasks  // NEW: grouped structure
    expandedTaskGroups map[string]bool       // NEW: expansion state persistence
}

func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) {
    m.tasks = msg.Tasks
    m.groupedTasks = msg.GroupedTasks
    
    // Initialize expansion state tracking
    if m.expandedTaskGroups == nil {
        m.expandedTaskGroups = make(map[string]bool)
    }
    
    // Sync expansion state from saved preferences
    for i := range m.groupedTasks.Groups {
        groupID := m.groupedTasks.Groups[i].ProjectID
        if _, exists := m.expandedTaskGroups[groupID]; !exists {
            m.expandedTaskGroups[groupID] = true  // Default: expanded
        }
        m.groupedTasks.Groups[i].IsExpanded = m.expandedTaskGroups[groupID]
    }
}
```

### Testing Pattern

Use table-driven tests with validateFunc for complex assertions:

```go
func TestNewGroupedTasks(t *testing.T) {
    tests := []struct {
        name            string
        tasks           []Task
        parentProjectID string
        projectNames    map[string]string
        wantDirectCount int
        wantGroupCount  int
        wantTotalCount  int
        validateFunc    func(t *testing.T, g *GroupedTasks)
    }{
        {
            name: "mixed direct and nested tasks",
            tasks: []Task{
                {ID: "t1", ProjectID: "proj-1", Title: "Direct"},
                {ID: "t2", ProjectID: "sub-1", Title: "Nested"},
            },
            parentProjectID: "proj-1",
            projectNames:    map[string]string{"proj-1": "Main", "sub-1": "Sub"},
            wantDirectCount: 1,
            wantGroupCount:  1,
            wantTotalCount:  2,
            validateFunc: func(t *testing.T, g *GroupedTasks) {
                if g.Groups[0].ProjectName != "Sub" {
                    t.Error("expected Sub project name")
                }
            },
        },
        // More test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := NewGroupedTasks(tt.tasks, tt.parentProjectID, tt.projectNames)
            
            if len(result.DirectTasks) != tt.wantDirectCount {
                t.Errorf("direct tasks: got %d, want %d", len(result.DirectTasks), tt.wantDirectCount)
            }
            
            if tt.validateFunc != nil {
                tt.validateFunc(t, result)
            }
        })
    }
}
```

### Design Decisions

**1. Mutable vs Immutable**:
- ✅ Mutable structure for CRUD operations after construction
- ✅ Allows dynamic task addition/removal without reconstruction
- ✅ Expansion state can be toggled by UI interactions

**2. No Error Returns**:
- ✅ Constructor never returns errors
- ✅ Graceful handling of edge cases (nil, empty, missing data)
- ✅ "Unknown Project" fallback for missing names
- ✅ Simplifies calling code (no error checking needed)

**3. Private State**:
- ✅ `parentProjectID` is private to encapsulate grouping logic
- ✅ External code doesn't need to know about parent/child relationship
- ✅ AddTask method uses private field to determine grouping

**4. Order Preservation**:
- ✅ Tasks maintain discovery order (as they appear in input)
- ✅ Groups maintain discovery order (first occurrence of each ProjectID)
- ✅ Predictable ordering for UI consistency

### When to Use GroupedTasks

**Use GroupedTasks when**:
- Displaying tasks from a project with nested subprojects
- Implementing hierarchical task views in TUI
- Supporting expand/collapse functionality for task groups
- Preserving user's expansion state across navigation

**Don't use GroupedTasks when**:
- Working with flat task lists (use `[]Task` directly)
- Tasks belong to a single project only
- You need database query results (use repository methods)

---
 
## Best Practices

### 1. Always Use Factory Methods

❌ **Don't**:
```go
task := &domain.Task{
    ID:     uuid.New().String(),
    Title:  "",
    Status: "invalid",
}
```

✅ **Do**:
```go
task, err := domain.NewTask(domain.NewTaskParams{
    Title:    "My Task",
    Status:   domain.TaskStatusTodo,
    Priority: domain.PriorityMedium,
})
```

### 2. Validate at Creation, Not Later

Validation happens **once** in the factory method, not scattered throughout the codebase.

### 3. Use Value Objects for Domain Concepts

Instead of raw strings/ints, use value objects that enforce validity:

❌ **Don't**:
```go
status := "todo" // Could be any string
```

✅ **Do**:
```go
status := domain.TaskStatusTodo // Type-safe, validated
```

### 4. Check Errors Explicitly

Domain errors are sentinel errors that should be checked explicitly:

```go
task, err := domain.NewTask(params)
if err != nil {
    if errors.Is(err, domain.ErrTaskTitleEmpty) {
        // Handle empty title
    } else if errors.Is(err, domain.ErrTaskInvalidDateRange) {
        // Handle invalid date range
    }
    return err
}
```

---

## Key Files

| File | Purpose |
|------|---------|
| `internal/domain/task.go` | Task entity with NewTask factory |
| `internal/domain/project.go` | Project entity with NewProject factory |
| `internal/domain/area.go` | Area entity |
| `internal/domain/subarea.go` | Subarea entity |
| `internal/domain/value_objects.go` | Status, Priority, Color, Duration, Progress |

---

**Navigation**: [← Overview](01-overview.md) | [Next: Service Layer →](03-service-layer.md)
