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
