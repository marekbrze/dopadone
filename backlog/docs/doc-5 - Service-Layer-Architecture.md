---
id: doc-5
title: Service Layer Architecture
type: technical
created_date: '2026-03-05'
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
func (s *ProjectService) ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)
func (s *ProjectService) Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)
func (s *ProjectService) SoftDelete(ctx context.Context, id string) error
func (s *ProjectService) HardDelete(ctx context.Context, id string) error
```

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

1. **Separation of Concerns**: Business logic isolated from presentation layer
2. **Type Safety**: Domain types prevent invalid states
3. **Testability**: Services can be mocked/stubbed for testing
4. **Reusability**: Same service layer can be used by CLI, TUI, and future REST API
5. **Validation**: Centralized business validation rules
6. **Error Handling**: Domain-specific errors with proper context
7. **Maintainability**: Clear boundaries between layers

## File Structure

```
internal/
  domain/
    project.go    # Project domain types and constants
    task.go       # Task domain types and constants
    subarea.go    # Subarea domain types and constants
    area.go       # Area domain types and constants
  
  service/
    project.go    # ProjectService
    task.go       # TaskService
    subarea.go    # SubareaService
    area.go       # AreaService
    errors.go     # Domain errors
  
  converter/
    project.go    # DB → Domain converters
    task.go       # DB → Domain converters
    subarea.go    # DB → Domain converters
    area.go       # DB → Domain converters

cmd/projectdb/
  main.go         # ServiceContainer and GetServices()
  projects.go     # CLI commands using ProjectService
  tasks.go        # CLI commands using TaskService
  subareas.go     # CLI commands using SubareaService
  areas.go        # CLI commands using AreaService
```

## Testing Strategy

### Unit Tests

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
