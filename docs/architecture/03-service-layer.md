# Service Layer

## Overview

The Service Layer orchestrates business logic, acting as the bridge between presentation (CLI/TUI) and data access (repository). It implements the **Provider Pattern** where interfaces are defined alongside implementations.

**Key Characteristics**:
- **Interface-first design**: All services expose interfaces for testability
- **Dependency injection**: Repository injected via constructor
- **Context-first methods**: All methods accept `context.Context` as first parameter
- **Business logic encapsulation**: Service handles business rules, domain handles validation

## Service Interface Pattern

### Provider Pattern

Interfaces are defined **alongside implementations** (not in consumer packages):

```go
// internal/service/interfaces.go
package service

// AreaServiceInterface defines the contract for area business operations.
// Areas are top-level organizational units that contain subareas, projects, and tasks.
type AreaServiceInterface interface {
    // List retrieves all non-deleted areas sorted by sort_order.
    List(ctx context.Context) ([]domain.Area, error)
    
    // GetByID retrieves a single area by its unique identifier.
    GetByID(ctx context.Context, id string) (*domain.Area, error)
    
    // Create creates a new area with the given name and color.
    Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error)
    
    // Update modifies an existing area's name and color.
    Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error)
    
    // SoftDelete marks an area as deleted without removing it from the database.
    SoftDelete(ctx context.Context, id string) error
    
    // HardDelete permanently removes an area and all its children.
    HardDelete(ctx context.Context, id string) error
}
```

**Why Provider Pattern?**
- Keeps interfaces close to implementations for easier maintenance
- Allows consumers to define their own interfaces if needed
- Simplifies the dependency graph
- Enables straightforward mocking for tests

### Interface Design Principles

1. **Small, focused interfaces**: Each service has a single responsibility
2. **Context-first**: All methods accept `context.Context` as first parameter
3. **Domain types**: Use domain types in interfaces, not DB types
4. **Clear contracts**: Document what each method does

---

## Service Interfaces

### TaskServiceInterface

```go
type TaskServiceInterface interface {
    Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)
    GetByID(ctx context.Context, id string) (*domain.Task, error)
    ListByProject(ctx context.Context, projectID string) ([]domain.Task, error)
    ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
    ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Task, error)
    ListNext(ctx context.Context) ([]domain.Task, error)
    ListAll(ctx context.Context) ([]domain.Task, error)
    ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error)
    Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error)
    SoftDelete(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
    ToggleIsNext(ctx context.Context, id string) (*domain.Task, error)
}
```

### ProjectServiceInterface

```go
type ProjectServiceInterface interface {
    Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error)
    GetByID(ctx context.Context, id string) (*domain.Project, error)
    ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error)
    ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)
    ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)
    ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error)
    ListAll(ctx context.Context) ([]domain.Project, error)
    Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)
    SoftDelete(ctx context.Context, id string) error
    SoftDeleteWithCascade(ctx context.Context, id string) error
    HardDelete(ctx context.Context, id string) error
}
```

---

## Dependency Injection

### Service Construction

Services are created with repository interfaces injected:

```go
// internal/service/task_service.go
type TaskService struct {
    repo db.Querier
}

func NewTaskService(repo db.Querier, tm *db.TransactionManager) *TaskService {
    return &TaskService{repo: repo}
}
```

### Service Container

All services are created once in a container:

```go
// cmd/dopa/main.go
type ServiceContainer struct {
    Projects  *service.ProjectService
    Tasks     *service.TaskService
    Subareas  *service.SubareaService
    Areas     *service.AreaService
}

func GetServices() (*ServiceContainer, error) {
    db, err := GetDB()
    if err != nil {
        return nil, err
    }
    
    querier := db.New(db)
    txManager := db.NewTransactionManager()
    
    return &ServiceContainer{
        Projects:  service.NewProjectService(querier, txManager),
        Tasks:     service.NewTaskService(querier, txManager),
        Subareas:  service.NewSubareaService(querier, txManager),
        Areas:     service.NewAreaService(querier, txManager),
    }, nil
}
```

---

## Business Logic Patterns

### Pattern 1: Validation via Domain Factory

Service uses domain factory for validation:

```go
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    // 1. Use domain factory for validation
    task, err := domain.NewTask(domain.NewTaskParams{
        ProjectID:         params.ProjectID,
        Title:             params.Title,
        Description:       params.Description,
        StartDate:         params.StartDate,
        Deadline:          params.Deadline,
        Priority:          params.Priority,
        Context:           params.Context,
        EstimatedDuration: params.EstimatedDuration,
        Status:            params.Status,
        IsNext:            params.IsNext,
    })
    if err != nil {
        return nil, err // Domain validation failed
    }
    
    // 2. Convert to DB params
    var isNext int64
    if task.IsNext {
        isNext = 1
    }
    
    dbParams := db.CreateTaskParams{
        ID:          task.ID,
        ProjectID:   task.ProjectID,
        Title:       task.Title,
        Description: sql.NullString{String: task.Description, Valid: task.Description != ""},
        Priority:    task.Priority.String(),
        Status:      task.Status.String(),
        IsNext:      isNext,
        CreatedAt:   task.CreatedAt,
        UpdatedAt:   task.UpdatedAt,
    }
    
    // 3. Call repository
    dbTask, err := s.repo.CreateTask(ctx, dbParams)
    if err != nil {
        return nil, err
    }
    
    // 4. Convert back to domain
    result := converter.DbTaskToDomain(dbTask)
    return &result, nil
}
```

### Pattern 2: Error Handling with Context

Service wraps errors with context:

```go
func (s *TaskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
    res, err := s.repo.GetTaskByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrTaskNotFound
        }
        return nil, fmt.Errorf("get task by id: %w", err)
    }
    
    result := converter.DbTaskToDomain(res)
    return &result, nil
}
```

### Pattern 3: Business Rules in Service

Some business logic lives in service, not domain:

```go
// Example: Auto-assign sort order
func (s *AreaService) Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
    // Business rule: Auto-assign sort order
    areas, err := s.repo.ListAreas(ctx)
    if err != nil {
        return nil, err
    }
    sortOrder := len(areas)
    
    area, err := domain.NewArea(name, color, sortOrder)
    if err != nil {
        return nil, err
    }
    
    // ... create in repository
}
```

### Pattern 4: Context Usage

Context is passed through all layers:

```go
func (s *TaskService) Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error) {
    // Context enables:
    // - Cancellation (ctx.Done())
    // - Timeouts (context.WithTimeout)
    // - Tracing (opentelemetry.FromContext)
    
    existingTask, err := s.GetByID(ctx, params.ID)
    if err != nil {
        return nil, err
    }
    
    // ... update logic
    
    return s.repo.UpdateTask(ctx, dbParams)
}
```

---

## Error Handling

### Service-Level Errors

Services define their own sentinel errors:

```go
// internal/service/task_service.go
var (
    ErrTaskNotFound = errors.New("task not found")
)

// internal/service/project_service.go
var (
    ErrProjectNotFound = errors.New("project not found")
    ErrProjectHasChildren = errors.New("project has child projects")
)
```

### Error Wrapping Pattern

```go
func (s *TaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error) {
    // Wrap repository errors with context
    dbTasks, err := s.repo.ListTasksByProject(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("list tasks by project: %w", err)
    }
    
    // Convert to domain
    tasks := make([]domain.Task, len(dbTasks))
    for i, dbTask := range dbTasks {
        tasks[i] = converter.DbTaskToDomain(dbTask)
    }
    
    return tasks, nil
}
```

### Error Checking in Presentation Layer

```go
// cmd/dopa/tasks.go
task, err := services.Tasks.GetByID(ctx, id)
if err != nil {
    if errors.Is(err, service.ErrTaskNotFound) {
        return fmt.Errorf("task %s not found", id)
    }
    return fmt.Errorf("failed to get task: %w", err)
}
```

---

## Service Method Examples

### Create Method (Full Flow)

```go
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    // 1. Domain validation (factory method)
    task, err := domain.NewTask(domain.NewTaskParams{
        ProjectID:         params.ProjectID,
        Title:             params.Title,
        Status:            params.Status,
        Priority:          params.Priority,
        Description:       params.Description,
        StartDate:         params.StartDate,
        Deadline:          params.Deadline,
        Context:           params.Context,
        EstimatedDuration: params.EstimatedDuration,
        IsNext:            params.IsNext,
    })
    if err != nil {
        return nil, err
    }
    
    // 2. Prepare DB params (handle NULL values)
    var isNext int64
    if task.IsNext {
        isNext = 1
    }
    
    var estimatedDuration sql.NullInt64
    if task.EstimatedDuration != 0 {
        estimatedDuration = sql.NullInt64{
            Int64: int64(task.EstimatedDuration.Int()),
            Valid: true,
        }
    }
    
    dbParams := db.CreateTaskParams{
        ID:                task.ID,
        ProjectID:         task.ProjectID,
        Title:             task.Title,
        Description:       sql.NullString{String: task.Description, Valid: task.Description != ""},
        StartDate:         task.StartDate,
        Deadline:          task.Deadline,
        Priority:          task.Priority.String(),
        Context:           sql.NullString{String: task.Context, Valid: task.Context != ""},
        EstimatedDuration: estimatedDuration,
        Status:            task.Status.String(),
        IsNext:            isNext,
        CreatedAt:         task.CreatedAt,
        UpdatedAt:         task.UpdatedAt,
        DeletedAt:         task.DeletedAt,
    }
    
    // 3. Call repository
    dbResult, err := s.repo.CreateTask(ctx, dbParams)
    if err != nil {
        return nil, err
    }
    
    // 4. Convert back to domain
    result := converter.DbTaskToDomain(dbResult)
    return &result, nil
}
```

### List Method

```go
func (s *TaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error) {
    dbTasks, err := s.repo.ListTasksByProject(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("list tasks by project: %w", err)
    }
    
    tasks := make([]domain.Task, len(dbTasks))
    for i, dbTask := range dbTasks {
        tasks[i] = converter.DbTaskToDomain(dbTask)
    }
    
    return tasks, nil
}
```

### Update Method

```go
func (s *TaskService) Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error) {
    existing, err := s.GetByID(ctx, params.ID)
    if err != nil {
        return nil, err
    }
    
    // Apply updates (only non-nil fields)
    if params.Title != nil {
        existing.Title = *params.Title
    }
    if params.Status != nil {
        existing.Status = *params.Status
    }
    if params.Priority != nil {
        existing.Priority = *params.Priority
    }
    if params.Description != nil {
        existing.Description = *params.Description
    }
    
    existing.UpdatedAt = time.Now()
    
    // Convert to DB params
    dbParams := db.UpdateTaskParams{
        ID:          existing.ID,
        Title:       existing.Title,
        Description: sql.NullString{String: existing.Description, Valid: existing.Description != ""},
        Status:      existing.Status.String(),
        Priority:    existing.Priority.String(),
        UpdatedAt:   existing.UpdatedAt,
    }
    
    dbResult, err := s.repo.UpdateTask(ctx, dbParams)
    if err != nil {
        return nil, err
    }
    
    result := converter.DbTaskToDomain(dbResult)
    return &result, nil
}
```

---

## Advanced Patterns

### Recursive Data Loading

For hierarchical data structures, services can provide recursive loading methods:

```go
// ListByProjectRecursive retrieves tasks from a project and all nested subprojects
func (s *TaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
    // Use SQL CTE (Common Table Expression) for efficient recursive loading
    dbTasks, err := s.repo.ListTasksByProjectRecursive(ctx, projectID)
    if err != nil {
        return nil, fmt.Errorf("list tasks recursively: %w", err)
    }
    
    // Convert to domain
    tasks := make([]domain.Task, len(dbTasks))
    for i, dbTask := range dbTasks {
        tasks[i] = converter.DbTaskToDomain(dbTask)
    }
    
    return tasks, nil
}
```

**SQL Pattern (WITH RECURSIVE CTE):**

```sql
-- queries/tasks.sql
-- name: ListTasksByProjectRecursive :many
WITH RECURSIVE project_tree AS (
    -- Base case: start with the root project
    SELECT id, parent_id, subarea_id
    FROM projects
    WHERE id = ? AND deleted_at IS NULL
    
    UNION ALL
    
    -- Recursive case: find all child projects
    SELECT p.id, p.parent_id, p.subarea_id
    FROM projects p
    INNER JOIN project_tree pt ON p.parent_id = pt.id
    WHERE p.deleted_at IS NULL
)
SELECT t.*
FROM tasks t
INNER JOIN project_tree pt ON t.project_id = pt.id
WHERE t.deleted_at IS NULL
ORDER BY t.is_next DESC, t.priority DESC, t.deadline ASC, t.title ASC;
```

**Benefits:**
- **Single query**: Fetches all tasks in one database call
- **Performance**: Database engine optimizes the recursive CTE
- **Consistency**: All data fetched at the same point in time
- **Simplicity**: No N+1 query problem

**When to Use:**
- Hierarchical data (projects → subprojects, areas → subareas)
- Tree traversal (categories, tags, organizational structures)
- Graph queries (when depth is bounded)

**Dependencies:**
- Service may need to inject other services for complex operations
- SQL query uses WITH RECURSIVE (supported by SQLite, PostgreSQL)

### Cascade Soft Delete Operations

For hierarchical entities with nested children, services provide cascade soft delete functionality that recursively marks all descendants as deleted in a single atomic operation.

**Implementation Pattern**:

```go
// SoftDeleteWithCascade soft deletes a project and all its descendants (child projects and tasks).
// The operation is atomic - either all deletions succeed or none do.
func (s *ProjectService) SoftDeleteWithCascade(ctx context.Context, id string) error {
    // Step 1: Validate project exists
    project, err := s.repo.GetProjectByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrProjectNotFound
        }
        return fmt.Errorf("get project %s: %w", id, err)
    }
    
    // Step 2: Check if already deleted (idempotent)
    if project.DeletedAt.Valid {
        return nil // Already soft deleted, no-op
    }

    // Step 3: Execute in transaction for atomicity
    if s.tm == nil {
        return s.softDeleteRecursive(ctx, s.repo, id)
    }

    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        return s.softDeleteRecursive(ctx, tx, id)
    })
}
```

**Recursive Helper**:

```go
// softDeleteRecursive recursively soft deletes a project, its children, and their tasks.
// Uses depth-first traversal: delete children first, then parent.
func (s *ProjectService) softDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error {
    // Step 1: Get all direct children
    children, err := q.ListProjectsByParent(ctx, sql.NullString{
        String: projectID,
        Valid:  true,
    })
    if err != nil {
        return fmt.Errorf("list child projects of %s: %w", projectID, err)
    }

    // Step 2: Recursively delete each child (depth-first)
    for _, child := range children {
        if err := s.softDeleteRecursive(ctx, q, child.ID); err != nil {
            return fmt.Errorf("cascade delete child %s: %w", child.ID, err)
        }
    }

    // Step 3: Soft delete all tasks in current project
    now := time.Now()
    if err := q.SoftDeleteTasksByProject(ctx, db.SoftDeleteTasksByProjectParams{
        DeletedAt: sql.NullTime{Time: now, Valid: true},
        ProjectID: projectID,
    }); err != nil {
        return fmt.Errorf("soft delete tasks in project %s: %w", projectID, err)
    }

    // Step 4: Soft delete current project
    if err := q.SoftDeleteProject(ctx, db.SoftDeleteProjectParams{
        DeletedAt: sql.NullTime{Time: now, Valid: true},
        ID:        projectID,
    }); err != nil {
        return fmt.Errorf("soft delete project %s: %w", projectID, err)
    }

    return nil
}
```

**Key Design Decisions**:
- **Depth-first traversal**: Children deleted before parent to maintain referential integrity
- **Transaction support**: All operations in single transaction for atomicity
- **Idempotency**: Safe to call multiple times (no error if already deleted)
- **Error wrapping**: Each operation wrapped with context for debugging
- **Same timestamp**: All deletions in cascade use same `now` value for consistent audit trail

**SQL Query**:

```sql
-- name: SoftDeleteTasksByProject :exec
-- Soft deletes all tasks within a project
UPDATE tasks
SET deleted_at = ?
WHERE project_id = ? AND deleted_at IS NULL;
```

**When to Use**:
- Deleting a parent project that has subprojects
- Ensuring data consistency across hierarchy
- When all descendant data should be marked as deleted together

**See Also**: `internal/service/project_service.go` for reference implementation

**Dependency Injection for Complex Services:**

When services need to coordinate, inject dependencies via constructor:

```go
type TaskService struct {
    repo            db.Querier
    projectService  ProjectServiceInterface  // Injected dependency
}

func NewTaskService(repo db.Querier, tm *db.TransactionManager, projectService ProjectServiceInterface) *TaskService {
    return &TaskService{
        repo:           repo,
        projectService: projectService,
    }
}
```

---

## Error Wrapping Best Practices

The service layer plays a crucial role in translating low-level errors into domain-appropriate errors while preserving error context.

### Error Translation Pattern

Services translate repository/database errors into domain errors:

```go
func (s *TaskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
    res, err := s.repo.GetTaskByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Map database error to domain error
            return nil, domain.NewNotFoundError("task", id)
        }
        // Wrap with operation context
        return nil, domain.NewDatabaseError("GetTaskByID", err)
    }
    
    result := converter.DbTaskToDomain(res)
    return &result, nil
}
```

### Graceful Handling of Empty Results

Services handle edge cases gracefully without errors:

```go
func (s *TaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
    // Validate input
    if projectID == "" {
        return []domain.Task{}, nil  // Empty result, not error
    }
    
    // Check if project exists
    _, err := s.projectSvc.GetByID(ctx, projectID)
    if err != nil {
        if domain.IsNotFound(err) {
            // Project doesn't exist - return empty, not error
            return []domain.Task{}, nil
        }
        // Database error - wrap with context
        return nil, fmt.Errorf("check project existence: %w", err)
    }
    
    // Load all tasks
    allTasks, err := s.db.ListAllTasks(ctx)
    if err != nil {
        return nil, fmt.Errorf("list tasks: %w", err)
    }
    
    // Filter and build hierarchy...
}
```

### Error Wrapping Principles

**1. Use domain error types**:
```go
// ✅ Good: Use domain errors
return nil, domain.NewNotFoundError("task", id)

// ❌ Bad: Use generic errors
return nil, fmt.Errorf("task not found")
```

**2. Add operation context**:
```go
// ✅ Good: Wrap with context
return nil, fmt.Errorf("get task by id: %w", err)

// ❌ Bad: Return raw error
return nil, err
```

**3. Handle edge cases gracefully**:
```go
// ✅ Good: Empty result for non-existent parent
if domain.IsNotFound(err) {
    return []domain.Task{}, nil
}

// ❌ Bad: Return error for missing parent
if err != nil {
    return nil, err
}
```

**4. Use sentinel errors for type checking**:
```go
// ✅ Good: Check error type
if domain.IsNotFound(err) {
    // Handle not found
}

// ❌ Bad: String comparison
if err.Error() == "not found" {
    // Fragile and error-prone
}
```

### Error Context Preservation

Always preserve the original error for debugging:

```go
func (s *TaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error) {
    dbTasks, err := s.repo.ListTasksByProject(ctx, projectID)
    if err != nil {
        // ✅ Preserves original error with %w
        return nil, fmt.Errorf("list tasks by project %s: %w", projectID, err)
    }
    
    // ... rest of implementation
}
```

### Handling Orphaned Data

Services handle orphaned entities gracefully:

```go
func (s *TaskService) buildProjectHierarchy(ctx context.Context) map[string][]string {
    projects, err := s.projectSvc.ListAll(ctx)
    if err != nil {
        return map[string][]string{}  // Empty map, not error
    }
    
    hierarchy := make(map[string][]string)
    
    for _, project := range projects {
        if project.ParentID != nil {
            parentID := *project.ParentID
            
            // Include even if parent doesn't exist (orphaned)
            // These won't appear in recursive loading
            hierarchy[parentID] = append(hierarchy[parentID], project.ID)
        }
    }
    
    return hierarchy
}
```

### Service-to-TUI Error Flow

```
Repository Error (sql.ErrNoRows)
        ↓
Service wraps (domain.NewNotFoundError)
        ↓
TUI checks type (domain.IsNotFound)
        ↓
User message ("Resource not found")
```

**Example Flow**:

```go
// Repository (internal/db/tasks.sql.go)
func (q *Queries) GetTaskByID(ctx context.Context, id string) (Task, error) {
    // Returns sql.ErrNoRows if not found
}

// Service (internal/service/task_service.go)
func (s *TaskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
    task, err := s.repo.GetTaskByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.NewNotFoundError("task", id)
        }
        return nil, domain.NewDatabaseError("GetTaskByID", err)
    }
    return converter.DbTaskToDomain(task), nil
}

// TUI (internal/tui/handlers.go)
func (m *Model) handleTaskLoaded(msg TaskLoadedMsg) {
    if msg.Err != nil {
        if domain.IsNotFound(msg.Err) {
            m.showError("Task not found")
        } else if domain.IsDatabaseError(msg.Err) {
            m.showError("Unable to load task. Please try again.")
        } else {
            m.showError("An error occurred")
        }
        return
    }
    // ... success handling
}
```

### Error Handling Checklist

When implementing service methods, ensure:

- ✅ **Use domain errors**: Return domain.NewNotFoundError(), not generic errors
- ✅ **Wrap with context**: Use `fmt.Errorf("operation: %w", err)` to preserve stack
- ✅ **Handle empty gracefully**: Return empty slices/maps, not errors, for missing parents
- ✅ **Check error types**: Use domain.IsNotFound(), not string comparison
- ✅ **Preserve original error**: Always use `%w` format verb for error wrapping
- ✅ **Add timeouts**: Use context.WithTimeout for long operations
- ✅ **Log internally**: Log technical errors for debugging, return clean errors to callers

---

## Testing Services

Services are tested using **mock implementations** of the repository interface.

### Mock Pattern

```go
// internal/service/task_service_test.go
type mockTaskQuerier struct {
    createTaskFunc          func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
    getTaskByIDFunc         func(ctx context.Context, id string) (db.Task, error)
    listTasksByProjectFunc  func(ctx context.Context, projectID string) ([]db.Task, error)
    updateTaskFunc          func(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error)
    // ... other methods
}

func (m *mockTaskQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
    if m.createTaskFunc != nil {
        return m.createTaskFunc(ctx, arg)
    }
    return db.Task{}, nil
}

func (m *mockTaskQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
    if m.getTaskByIDFunc != nil {
        return m.getTaskByIDFunc(ctx, id)
    }
    return db.Task{}, nil
}

// ... implement all interface methods
```

### Test Example

```go
func TestTaskService_Create(t *testing.T) {
    tests := []struct {
        name    string
        params  service.CreateTaskParams
        mockFn  func(*mockTaskQuerier)
        wantErr error
    }{
        {
            name: "valid task",
            params: service.CreateTaskParams{
                ProjectID: "proj-123",
                Title:     "Test Task",
                Status:    domain.TaskStatusTodo,
                Priority:  domain.PriorityMedium,
            },
            mockFn: func(m *mockTaskQuerier) {
                m.createTaskFunc = func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
                    return db.Task{
                        ID:        arg.ID,
                        Title:     arg.Title,
                        ProjectID: arg.ProjectID,
                    }, nil
                }
            },
            wantErr: nil,
        },
        {
            name: "empty title fails validation",
            params: service.CreateTaskParams{
                ProjectID: "proj-123",
                Title:     "",
                Status:    domain.TaskStatusTodo,
                Priority:  domain.PriorityMedium,
            },
            mockFn:  func(m *mockTaskQuerier) {},
            wantErr: domain.ErrTaskTitleEmpty,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &mockTaskQuerier{}
            tt.mockFn(mock)
            
            svc := service.NewTaskService(mock, nil)
            
            result, err := svc.Create(context.Background(), tt.params)
            
            if tt.wantErr != nil {
                if !errors.Is(err, tt.wantErr) {
                    t.Errorf("expected error %v, got %v", tt.wantErr, err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
                if result == nil {
                    t.Error("expected result, got nil")
                }
            }
        })
    }
}
```

---

## Best Practices

### 1. Accept Interfaces, Return Structs

```go
// ✅ Good: Service accepts interface
func NewTaskService(repo db.Querier) *TaskService {
    return &TaskService{repo: repo}
}

// ✅ Good: Service returns concrete type
func (s *TaskService) Create(...) (*domain.Task, error) {
    // ...
}
```

### 2. Context-First Design

```go
// ✅ Good: Context is first parameter
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    // ...
}

// ❌ Bad: Context not first
func (s *TaskService) Create(params CreateTaskParams, ctx context.Context) (*domain.Task, error) {
    // ...
}
```

### 3. Domain Validation in Factory

```go
// ✅ Good: Validation in domain factory
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    task, err := domain.NewTask(domain.NewTaskParams{...})
    if err != nil {
        return nil, err // Domain validation failed
    }
    // ...
}

// ❌ Bad: Validation scattered in service
func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
    if params.Title == "" {
        return nil, errors.New("title empty")
    }
    if !params.Status.IsValid() {
        return nil, errors.New("invalid status")
    }
    // ... more validation
}
```

### 4. Wrap Errors with Context

```go
// ✅ Good: Error with context
return nil, fmt.Errorf("create task: %w", err)

// ❌ Bad: Lost context
return nil, err
```

---

## Key Files

| File | Purpose |
|------|---------|
| `internal/service/interfaces.go` | Service interfaces |
| `internal/service/task_service.go` | TaskService implementation |
| `internal/service/project_service.go` | ProjectService implementation |
| `internal/service/area_service.go` | AreaService implementation |
| `internal/service/subarea_service.go` | SubareaService implementation |
| `cmd/dopa/main.go` | Service container and DI |

---

**Navigation**: [← Domain Layer](02-domain-layer.md) | [Next: Converter Layer →](04-converter-layer.md)
