---
id: doc-5
title: Service Layer Architecture
type: technical
created_date: '2026-03-05'
updated_date: '2026-03-05'
---

# Service Layer Architecture

## Overview

This document describes the service layer architecture introduced to separate business logic from CLI and TUI command handlers. The service layer encapsulates domain operations and provides a clean API for data access.

## Architecture Layers

```
┌─────────────────────────────────────────┐
│          Presentation Layer             │
│  (CLI commands, TUI, Future REST API)   │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│           Service Layer                 │
│  (Business logic, validation,           │
│   domain operations)                    │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│         Data Access Layer               │
│  (sqlc-generated Querier interface)     │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│            SQLite Database              │
└─────────────────────────────────────────┘
```

## Service Layer Components

### Domain Types

Located in `internal/domain/`, these types represent business entities with proper type safety:

```go
type TaskStatus string
const (
    TaskStatusTodo       TaskStatus = "todo"
    TaskStatusInProgress TaskStatus = "in_progress"
    TaskStatusWaiting    TaskStatus = "waiting"
    TaskStatusDone       TaskStatus = "done"
)

type Priority string
const (
    PriorityCritical Priority = "critical"
    PriorityHigh     Priority = "high"
    PriorityMedium   Priority = "medium"
    PriorityLow      Priority = "low"
)

type ProjectStatus string
const (
    ProjectStatusActive    ProjectStatus = "active"
    ProjectStatusCompleted ProjectStatus = "completed"
    ProjectStatusOnHold    ProjectStatus = "on_hold"
    ProjectStatusArchived  ProjectStatus = "archived"
)

type Progress int
type Color string
```

### Services

Each entity has a dedicated service in `internal/service/`:

#### ProjectService

```go
type ProjectService struct {
    db db.Querier
}

func (s *ProjectService) Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error)
func (s *ProjectService) GetByID(ctx context.Context, id string) (*domain.Project, error)
func (s *ProjectService) ListAll(ctx context.Context) ([]domain.Project, error)
func (s *ProjectService) ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)
func (s *ProjectService) ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error)
func (s *ProjectService) ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error)
func (s *ProjectService) ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error)
func (s *ProjectService) ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)
func (s *ProjectService) Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)
func (s *ProjectService) SoftDelete(ctx context.Context, id string) error
func (s *ProjectService) HardDelete(ctx context.Context, id string) error
```

**ListBySubareaRecursive**: Returns all projects belonging to a subarea, including nested projects whose parent chain leads to the subarea. This method recursively traverses the project hierarchy to find all descendant projects.

- **Performance**: O(n) time complexity where n = total projects
- **Algorithm**: 
  1. Loads all non-deleted projects via `ListAll(ctx)`
  2. Builds a project map for O(1) parent lookups
  3. Filters projects that belong to the subarea (direct membership or via parent chain)
- **Edge Cases**: 
  - Empty subareaID returns empty slice
  - Orphaned projects (parent doesn't exist) are excluded
  - Soft-deleted projects are automatically excluded
- **Use Case**: When displaying all projects in a subarea tree view, this ensures nested projects are included even if they don't have a direct subareaID assignment

#### TaskService

```go
type TaskService struct {
    db db.Querier
}

func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)
func (s *TaskService) GetByID(ctx context.Context, id string) (*domain.Task, error)
func (s *TaskService) ListAll(ctx context.Context) ([]domain.Task, error)
func (s *TaskService) ListNext(ctx context.Context) ([]domain.Task, error)
func (s *TaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error)
func (s *TaskService) ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
func (s *TaskService) ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Task, error)
func (s *TaskService) Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error)
func (s *TaskService) SoftDelete(ctx context.Context, id string) error
func (s *TaskService) HardDelete(ctx context.Context, id string) error
```

#### SubareaService

```go
type SubareaService struct {
    db db.Querier
}

func (s *SubareaService) Create(ctx context.Context, name, areaID string, color domain.Color) (*domain.Subarea, error)
func (s *SubareaService) GetByID(ctx context.Context, id string) (*domain.Subarea, error)
func (s *SubareaService) ListAll(ctx context.Context) ([]domain.Subarea, error)
func (s *SubareaService) ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error)
func (s *SubareaService) Update(ctx context.Context, id, name, areaID string, color domain.Color) (*domain.Subarea, error)
func (s *SubareaService) SoftDelete(ctx context.Context, id string) error
func (s *SubareaService) HardDelete(ctx context.Context, id string) error
```

#### AreaService

```go
type AreaService struct {
    db db.Querier
}

func (s *AreaService) Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error)
func (s *AreaService) GetByID(ctx context.Context, id string) (*domain.Area, error)
func (s *AreaService) List(ctx context.Context) ([]domain.Area, error)
func (s *AreaService) Update(ctx context.Context, id, name string, color domain.Color) (*domain.Area, error)
func (s *AreaService) SoftDelete(ctx context.Context, id string) error
func (s *AreaService) HardDelete(ctx context.Context, id string) error
```

### Service Interfaces

All services implement corresponding interfaces defined in `internal/service/interfaces.go`. These interfaces enable:

- **Dependency Injection**: Services can be injected via interfaces rather than concrete types
- **Testability**: Easy mocking in unit tests without requiring database connections
- **Flexibility**: Consumers can define narrower interfaces if needed
- **Compile-time Safety**: Interface satisfaction is verified at compile time

#### Interface Definitions

```go
type AreaServiceInterface interface {
    List(ctx context.Context) ([]domain.Area, error)
    GetByID(ctx context.Context, id string) (*domain.Area, error)
    Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error)
    Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error)
    UpdateSortOrder(ctx context.Context, id string, sortOrder int) error
    ReorderAll(ctx context.Context, areaIDs []string) error
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*AreaStats, error)
}

type SubareaServiceInterface interface {
    Create(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error)
    GetByID(ctx context.Context, id string) (*domain.Subarea, error)
    ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error)
    Update(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*SubareaStats, error)
    GetEffectiveColor(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color
    ListAll(ctx context.Context) ([]domain.Subarea, error)
}

type ProjectServiceInterface interface {
    Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error)
    GetByID(ctx context.Context, id string) (*domain.Project, error)
    ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error)
    ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)
    ListAll(ctx context.Context) ([]domain.Project, error)
    ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)
    ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error)
    ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error)
    Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    GetStats(ctx context.Context, id string) (*ProjectStats, error)
    ValidateParentHierarchy(ctx context.Context, parentID string, projectID string) error
}

type TaskServiceInterface interface {
    Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)
    GetByID(ctx context.Context, id string) (*domain.Task, error)
    ListByProject(ctx context.Context, projectID string) ([]domain.Task, error)
    ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
    ListByPriority(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error)
    ListNext(ctx context.Context) ([]domain.Task, error)
    ListAll(ctx context.Context) ([]domain.Task, error)
    Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    SetStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)
    MarkCompleted(ctx context.Context, id string) (*domain.Task, error)
    SetPriority(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error)
    ToggleIsNext(ctx context.Context, id string) (*domain.Task, error)
}
```

#### Design Principles

1. **Provider Pattern**: Interfaces are defined alongside implementations (not in consumer packages)
   - Keeps interfaces close to implementations for maintainability
   - Simplifies dependency graph
   - Allows consumers to define narrower interfaces if needed

2. **Context-First**: All methods accept `context.Context` as first parameter
   - Follows Go best practices
   - Enables future cancellation and timeout support
   - Consistent API across all services

3. **Compile-Time Checks**: Interface satisfaction verified at compile time
   ```go
   var (
       _ AreaServiceInterface    = (*AreaService)(nil)
       _ SubareaServiceInterface = (*SubareaService)(nil)
       _ ProjectServiceInterface = (*ProjectService)(nil)
       _ TaskServiceInterface    = (*TaskService)(nil)
   )
   ```

### Service Container

Services are managed through a ServiceContainer in `cmd/projectdb/main.go`:

```go
type ServiceContainer struct {
    Projects *service.ProjectService
    Tasks    *service.TaskService
    Subareas *service.SubareaService
    Areas    *service.AreaService
    db       *sql.DB
}

func GetServices() (*ServiceContainer, error) {
    dbConn, err := GetDB()
    if err != nil {
        return nil, err
    }
    queries := db.New(dbConn)
    
    return &ServiceContainer{
        Projects: service.NewProjectService(queries),
        Tasks:    service.NewTaskService(queries),
        Subareas: service.NewSubareaService(queries),
        Areas:    service.NewAreaService(queries),
        db:       dbConn,
    }, nil
}

func (s *ServiceContainer) Close() error {
    if s.db != nil {
        return s.db.Close()
    }
    return nil
}
```

## Type Conversions

### DB → Domain Conversion

Located in `internal/converter/`, these functions convert sqlc types to domain types:

```go
func DbProjectToDomain(p db.Project) domain.Project
func DbTaskToDomain(t db.Task) domain.Task
func DbSubareaToDomain(s db.Subarea) domain.Subarea
func DbAreaToDomain(a db.Area) domain.Area
```

### CLI → Service Conversion

CLI commands parse user input and convert to service parameters:

```go
// CLI flags → domain types
status, err := cli.ParseTaskStatus(taskCreateStatus)
priority, err := cli.ParsePriority(taskCreatePriority)

// Domain types → service params
params := service.CreateTaskParams{
    ProjectID: taskCreateProjectID,
    Title:     taskCreateTitle,
    Status:    status,
    Priority:  priority,
    // ... other fields
}

// Call service
task, err := services.Tasks.Create(ctx, params)
```

## Error Handling

### Domain Errors

Services return domain-specific errors:

```go
var (
    ErrProjectNotFound = errors.New("project not found")
    ErrTaskNotFound    = errors.New("task not found")
    ErrSubareaNotFound = errors.New("subarea not found")
    ErrAreaNotFound    = errors.New("area not found")
    // ... other domain errors
)
```

### CLI Error Wrapping

CLI layer wraps service errors with context:

```go
task, err := services.Tasks.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, service.ErrTaskNotFound) {
        cli.ExitWithError(cli.NewNotFoundError("Task", id))
    }
    cli.ExitWithError(cli.WrapError(err, "Failed to get task"))
}
```

## Validation

### Service Layer Validation

Services handle business validation:

```go
func (s *ProjectService) Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error) {
    // Validate required fields
    if params.Name == "" {
        return nil, fmt.Errorf("%w: name is required", ErrValidation)
    }
    
    // Validate enums
    if !isValidProjectStatus(params.Status) {
        return nil, fmt.Errorf("%w: invalid status: %s", ErrValidation, params.Status)
    }
    
    // ... business logic
}
```

### CLI Layer Validation

CLI handles input validation:

```go
// Validate required flags
if err := cli.ValidateTaskTitle(taskCreateTitle); err != nil {
    cli.ExitWithError(err)
}

// Parse and validate enums
status, err := cli.ParseTaskStatus(taskCreateStatus)
if err != nil {
    cli.ExitWithError(err)
}
```

## Migration from Direct DB Access

### Before (Direct DB Access)

```go
func runProjectsCreate(cmd *cobra.Command, args []string) {
    db, err := GetDB()
    if err != nil {
        cli.ExitWithError(err)
    }
    defer db.Close()
    
    queries := db.New(db)
    project, err := queries.CreateProject(ctx, db.CreateProjectParams{
        // sqlc-generated params with sql.NullString types
    })
}
```

### After (Service Layer)

```go
func runProjectsCreate(cmd *cobra.Command, args []string) {
    services, err := GetServices()
    if err != nil {
        cli.ExitWithError(err)
    }
    defer services.Close()
    
    params := service.CreateProjectParams{
        // Domain types with proper type safety
        Name:     projectName,
        Status:   domain.ProjectStatusActive,
        Priority: domain.PriorityMedium,
    }
    
    project, err := services.Projects.Create(ctx, params)
}
```

## Benefits

1. **Separation of Concerns**: Business logic isolated from presentation layer (CLI and TUI)
2. **Type Safety**: Domain types prevent invalid states
3. **Testability**: Service interfaces enable easy mocking/stubbing for tests without database dependencies (used by both CLI and TUI tests)
4. **Reusability**: Same service layer can be used by CLI, TUI, and future REST API
5. **Validation**: Centralized business validation rules
6. **Error Handling**: Domain-specific errors with proper context
7. **Maintainability**: Clear boundaries between layers
8. **Dependency Injection**: Interfaces enable proper DI patterns and flexible component wiring
9. **Compile-Time Safety**: Interface satisfaction verified at compile time catches errors early
10. **Clean Architecture**: TUI depends on service interfaces, not database layer (Task-38, Task-39)

## File Structure

```
internal/
  domain/
    project.go    # Project domain types and constants
    task.go       # Task domain types and constants
    subarea.go    # Subarea domain types and constants
    area.go       # Area domain types and constants
  
  service/
    interfaces.go    # Service interfaces (AreaServiceInterface, etc.)
    project.go       # ProjectService implementation
    task.go          # TaskService implementation
    subarea.go       # SubareaService implementation
    area.go          # AreaService implementation
    errors.go        # Domain errors
  
  converter/
    project.go    # DB → Domain converters
    task.go       # DB → Domain converters
    subarea.go    # DB → Domain converters
    area.go       # DB → Domain converters
  
  tui/
    app.go         # TUI Model with service interfaces (Task-38)
    commands.go    # Loader & CRUD commands using services (Task-38, Task-39)
    mocks/         # Mock services for TUI testing
      helpers.go   # Mock setup helpers

cmd/projectdb/
  main.go         # ServiceContainer and GetServices()
  projects.go     # CLI commands using ProjectService
  tasks.go        # CLI commands using TaskService
  subareas.go     # CLI commands using SubareaService
  areas.go        # CLI commands using AreaService
```

## TUI Integration

The TUI uses service layer interfaces for both data loading and CRUD operations:

### TUI Commands Using Services (Task-38, Task-39)

**Loader Commands (Task-38):**
- `LoadAreasCmd(areaSvc AreaServiceInterface)` → `List(ctx)`
- `LoadSubareasCmd(subareaSvc SubareaServiceInterface, areaID)` → `ListByArea(ctx, areaID)`
- `LoadProjectsCmd(projectSvc ProjectServiceInterface, subareaID)` → `ListBySubareaRecursive(ctx, subareaID)`
- `LoadTasksCmd(taskSvc TaskServiceInterface, projectID)` → `ListByProject(ctx, projectID)`

**CRUD Commands (Task-39):**
- `CreateSubareaCmd(subareaSvc SubareaServiceInterface, ...)` → `Create(ctx, ...)`
- `CreateProjectCmd(projectSvc ProjectServiceInterface, ...)` → `Create(ctx, ...)`
- `CreateTaskCmd(taskSvc TaskServiceInterface, ...)` → `Create(ctx, ...)`
- `CreateAreaCmd(areaSvc AreaServiceInterface, ...)` → `Create(ctx, ...)`
- `UpdateAreaCmd(areaSvc AreaServiceInterface, ...)` → `Update(ctx, ...)`
- `DeleteAreaCmd(areaSvc AreaServiceInterface, id, hard)` → `SoftDelete(ctx, id)` or `HardDelete(ctx, id)`
- `ReorderAreasCmd(areaSvc AreaServiceInterface, areaIDs)` → `ReorderAll(ctx, areaIDs)`
- `LoadAreaStatsCmd(areaSvc AreaServiceInterface, areaID)` → `GetStats(ctx, areaID)`

All TUI commands receive service interfaces as parameters and call service methods without direct database access. Services return domain types directly, eliminating the need for converter logic in the TUI layer.

### TUI Model Structure

```go
type Model struct {
    areaSvc     service.AreaServiceInterface
    subareaSvc  service.SubareaServiceInterface
    projectSvc  service.ProjectServiceInterface
    taskSvc     service.TaskServiceInterface
    // ... other fields
}

func InitialModel(
    areaSvc service.AreaServiceInterface,
    subareaSvc service.SubareaServiceInterface,
    projectSvc service.ProjectServiceInterface,
    taskSvc service.TaskServiceInterface,
) Model {
    return Model{
        areaSvc:    areaSvc,
        subareaSvc: subareaSvc,
        projectSvc: projectSvc,
        taskSvc:    taskSvc,
        // ... other fields
    }
}
```

This pattern enables:
- Clean architecture boundaries (TUI → Services → Repository → Database)
- Easy mocking for unit tests
- No direct database dependency in TUI layer
- Services return domain types directly (no converter layer in TUI)

## Testing Strategy

### Unit Tests with Service Interfaces

Service interfaces enable easy mocking for unit tests. You can create mock implementations without requiring database connections:

```go
type MockAreaService struct {
    mock.Mock
}

func (m *MockAreaService) Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
    args := m.Called(ctx, name, color)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Area), args.Error(1)
}

// ... implement other interface methods

func TestTUIWithMockService(t *testing.T) {
    mockService := new(MockAreaService)
    mockService.On("List", mock.Anything).Return([]domain.Area{
        {ID: "1", Name: "Test Area", Color: "#FF5733"},
    }, nil)
    
    // Use mockService in TUI components without database
    model := NewModel(mockService)
    // ... test model behavior
}
```

### Service Layer Tests with Mock Querier

Services can be tested in isolation with mock Querier:

```go
type MockQuerier struct {
    // mock implementation
}

func TestProjectService_Create(t *testing.T) {
    mockDB := &MockQuerier{}
    service := NewProjectService(mockDB)
    
    params := CreateProjectParams{
        Name:     "Test Project",
        Status:   ProjectStatusActive,
        Priority: PriorityMedium,
    }
    
    project, err := service.Create(context.Background(), params)
    // ... assertions
}
```

### Integration Tests

CLI commands tested end-to-end:

```go
func TestProjectsCreateCommand(t *testing.T) {
    // Setup test database
    // Execute command
    // Verify output
}
```

## Future Enhancements

1. **Transaction Support**: Add transaction boundaries in service layer
2. **Caching**: Implement caching layer for frequently accessed data
3. **Event System**: Emit events for state changes (e.g., ProjectCompleted)
4. **Audit Logging**: Track all mutations through service layer
5. **Rate Limiting**: Add rate limiting for API usage
6. **Pagination**: Support pagination for large datasets

## Related Documentation

- [Data Layer Architecture](doc-1 - Data-Layer-Architecture.md) - Database and sqlc details
- [CLI CRUD Operations Guide](doc-2 - CLI-CRUD-Operations-Guide.md) - CLI command usage
- [TUI Architecture](doc-3 - TUI-Architecture.md) - TUI implementation details
