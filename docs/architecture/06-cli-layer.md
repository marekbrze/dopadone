# CLI Layer

## Overview

The CLI Layer provides the command-line interface using **Cobra**. It handles user input, calls service methods, and formats output. The CLI layer depends on the Service Layer, not the repository directly.

**Key Characteristics**:
- **Cobra framework**: Industry-standard CLI library for Go
- **Service injection**: Services injected via ServiceContainer
- **Flag parsing**: User-friendly flags with validation
- **Output formatting**: Multiple formats (table, JSON, YAML)

---

## Command Structure

### Root Command

```go
// cmd/dopa/main.go
var rootCmd = &cobra.Command{
    Use:   "dopa",
    Short: "CLI project management for developers",
    Long:  "Organize your projects, tasks, and workflows from the command line.",
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

### Resource Commands

Each resource (area, subarea, project, task) has its own command group:

```go
// cmd/dopa/tasks.go
var tasksCmd = &cobra.Command{
    Use:     "tasks",
    Short:   "Manage tasks",
    Long:    "Manage tasks in the project database.",
    Aliases: []string{"task"},
}

var tasksCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a new task",
    Run:   runTasksCreate,
}

var tasksListCmd = &cobra.Command{
    Use:   "list",
    Short: "List all tasks",
    Run:   runTasksList,
}
```

### Command Registration
```go
func init() {
    tasksCmd.AddCommand(tasksCreateCmd)
    tasksCmd.AddCommand(tasksListCmd)
    tasksCmd.AddCommand(tasksGetCmd)
    tasksCmd.AddCommand(tasksUpdateCmd)
    tasksCmd.AddCommand(tasksDeleteCmd)
    
    rootCmd.AddCommand(tasksCmd)
    rootCmd.AddCommand(projectsCmd)
    rootCmd.AddCommand(areasCmd)
    rootCmd.AddCommand(subareasCmd)
}
```

---

## Service Injection

### Service Container
Services are created once and accessed globally:
```go
// cmd/dopa/main.go
var services *ServiceContainer

type ServiceContainer struct {
    Projects  *service.ProjectService
    Tasks     *service.TaskService
    Subareas  *service.SubareaService
    Areas     *service.AreaService
}

func GetServices() (*ServiceContainer, error) {
    if services != nil {
        return services, nil
    }
    
    db, err := GetDB()
    if err != nil {
        return nil, err
    }
    
    querier := db.New(db)
    txManager := db.NewTransactionManager()
    
    services = &ServiceContainer{
        Projects:  service.NewProjectService(querier, txManager),
        Tasks:     service.NewTaskService(querier, txManager),
        Subareas:  service.NewSubareaService(querier, txManager),
        Areas:     service.NewAreaService(querier, txManager),
    }
    
    return services, nil
}
```

### Using Services in Commands
```go
func runTasksCreate(cmd *cobra.Command, args []string) error {
    services, err := GetServices()
    if err != nil {
        return err
    }
    
    title, _ := cmd.Flags().GetString("title")
    projectID, _ := cmd.Flags().GetString("project-id")
    
    task, err := services.Tasks.Create(context.Background(), service.CreateTaskParams{
        ProjectID: projectID,
        Title:     title,
        Status:    domain.TaskStatusTodo,
    })
    
    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }
    
    return output.Write(task)
}
```

---

## Flag Parsing

### Required Flags
```go
func init() {
    tasksCreateCmd.Flags().String("project-id", "", "Project ID (required)")
    tasksCreateCmd.Flags().String("title", "", "Task title (required)")
    tasksCreateCmd.MarkFlagRequired("project-id")
    tasksCreateCmd.MarkFlagRequired("title")
}
```

### Optional Flags
```go
func init() {
    tasksCreateCmd.Flags().String("description", "", "Task description")
    tasksCreateCmd.Flags().String("status", "todo", "Task status")
    tasksCreateCmd.Flags().String("priority", "medium", "Task priority")
    tasksCreateCmd.Flags().String("start-date", "", "Start date (YYYY-MM-DD)")
    tasksCreateCmd.Flags().String("deadline", "", "Deadline (YYYY-MM-DD)")
    tasksCreateCmd.Flags().Bool("next", false, "Mark as priority task")
}
```

---

## Output Formatting

### Format Flag
```go
// Global flags in rootCmd
rootCmd.PersistentFlags().StringP("output", "table", "Output format (table, json, yaml)")
rootCmd.PersistentFlags().StringP("format", "table", "Output format (table, json, yaml)")
```

### Output Writer
```go
// internal/cli/output/formatter.go
func Write(v interface{}) error {
    format, _ := getFormat()
    
    switch format {
    case "json":
        return writeJSON(v)
    case "yaml":
        return writeYAML(v)
    default:
        return writeTable(v)
    }
}

func writeJSON(v interface{}) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(v)
}

func writeYAML(v interface{}) error {
    data, err := yaml.Marshal(v)
    if err != nil {
        return err
    }
    fmt.Println(string(data))
    return nil
}

func writeTable(v interface{}) error {
    // Use tablewriter to format as table
    // ...
}
```

### Usage Examples
```bash
# Table format (default)
dopa tasks list

# JSON format
dopa tasks list --format json
dopa tasks list -o json

# YAML format
dopa tasks list --format yaml
```

---

## Error Handling

### CLI Error Wrapping
```go
// internal/cli/errors.go
func WrapError(err error, message string) error {
    return fmt.Errorf("%s: %w", message, err)
}
```

### Cleanup Error Handling

Cleanup operations (Close, Flush, RemoveAll) should not fail the main operation. Use these patterns:

**1. CloseWithLog helper for defer cleanup**:
```go
// internal/cli/helpers.go
func CloseWithLog(closer Closer, name string) {
    if closer == nil {
        return
    }
    if err := closer.Close(); err != nil {
        slog.Warn("failed to close resource", "name", name, "error", err)
    }
}

// Usage in commands
func runTasksList(cmd *cobra.Command, args []string) error {
    services, err := GetServices()
    if err != nil {
        return err
    }
    defer cli.CloseWithLog(services, "services")
    
    // ... command logic
}
```

**2. Explicit ignore for cleanup errors**:
```go
// When cleanup error is acceptable to ignore
defer func() { _ = rows.Close() }()
defer func() { _ = db.Close() }()
```

**3. Panic for programming errors**:
```go
// MarkFlagRequired errors are programming errors
func init() {
    tasksCreateCmd.Flags().String("title", "", "Task title")
    if err := tasksCreateCmd.MarkFlagRequired("title"); err != nil {
        panic(fmt.Sprintf("failed to mark flag required: %v", err))
    }
}
```

### Error Display
```go
func runTasksGet(cmd *cobra.Command, args []string) error {
    id := args[0]
    
    services, err := GetServices()
    if err != nil {
        return err
    }
    
    task, err := services.Tasks.GetByID(context.Background(), id)
    if err != nil {
        if errors.Is(err, service.ErrTaskNotFound) {
            return fmt.Errorf("task %s not found", id)
        }
        return fmt.Errorf("failed to get task: %w", err)
    }
    
    return output.Write(task)
}
```

---

## Command Examples

### Create Task
```go
func runTasksCreate(cmd *cobra.Command, args []string) error {
    services, err := GetServices()
    if err != nil {
        return err
    }
    
    title, _ := cmd.Flags().GetString("title")
    projectID, _ := cmd.Flags().GetString("project-id")
    description, _ := cmd.Flags().GetString("description")
    statusStr, _ := cmd.Flags().GetString("status")
    priorityStr, _ := cmd.Flags().GetString("priority")
    
    status, err := domain.ParseTaskStatus(statusStr)
    if err != nil {
        return fmt.Errorf("invalid status: %w", err)
    }
    
    priority, err := domain.ParsePriority(priorityStr)
    if err != nil {
        return fmt.Errorf("invalid priority: %w", err)
    }
    
    task, err := services.Tasks.Create(context.Background(), service.CreateTaskParams{
        ProjectID:   projectID,
        Title:       title,
        Description: description,
        Status:      status,
        Priority:    priority,
    })
    
    if err != nil {
        if errors.Is(err, domain.ErrTaskTitleEmpty) {
            return fmt.Errorf("task title cannot be empty")
        }
        return fmt.Errorf("failed to create task: %w", err)
    }
    
    return output.Write(task)
}
```

### List Tasks
```go
func runTasksList(cmd *cobra.Command, args []string) error {
    services, err := GetServices()
    if err != nil {
        return err
    }
    
    var tasks []domain.Task
    
    projectID, _ := cmd.Flags().GetString("project-id")
    if projectID != "" {
        tasks, err = services.Tasks.ListByProject(context.Background(), projectID)
    } else {
        tasks, err = services.Tasks.ListAll(context.Background())
    }
    
    if err != nil {
        return fmt.Errorf("failed to list tasks: %w", err)
    }
    
    return output.Write(tasks)
}
```

---

## Best Practices

### 1. Use Cobra Idiom
```go
// ✅ Good: Use Cobra commands
var tasksCmd = &cobra.Command{
    Use:   "tasks",
    Short: "Manage tasks",
}

// ❌ Bad: Manual flag parsing
```
// ✅ Good: Use Cobra flags
title := cmd.Flag("title").(string)

// ❌ Bad: Manual flag parsing
os.Args
```

### 2. Validate Flags
```go
// ✅ Good: Validate input
status, err := domain.ParseTaskStatus(statusStr)
if err != nil {
    return fmt.Errorf("invalid status: %w", err)
}

// ❌ Bad: Pass invalid input to service
status := domain.TaskStatus(statusStr) // Could be invalid!
```

### 3. Wrap Errors with Context
```go
// ✅ Good: Error with context
return fmt.Errorf("failed to create task: %w", err)

// ❌ Bad: Lost context
return err
```

### 4. Use Global Service Container
```go
// ✅ Good: Single service instance
services, err := GetServices()
if err != nil {
    return err
}

// ❌ Bad: Create service in each command
db := GetDB()
service := NewTaskService(db)
```

---

## Reducing Cyclomatic Complexity

When command functions become too complex (gocyclo limit: 30), extract helper functions to improve maintainability.

### Pattern: Extract Helper Functions

Break down complex command logic into focused helper functions:

```go
// ✅ Good: Complex command with helper functions
func runTasksUpdate(cmd *cobra.Command, args []string) {
    id := args[0]

    if err := validateUpdateFlags(cmd); err != nil {
        cli.ExitWithError(err)
    }

    services, err := GetServices()
    if err != nil {
        cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
    }
    defer cli.CloseWithLog(services, "services")

    ctx := context.Background()

    existing, err := services.Tasks.GetByID(ctx, id)
    if err != nil {
        if err == service.ErrTaskNotFound {
            cli.ExitWithError(fmt.Errorf("task not found: %s", id))
        }
        cli.ExitWithError(cli.WrapError(err, "failed to get task"))
    }

    params, err := prepareTaskUpdateParams(cmd, existing)
    if err != nil {
        cli.ExitWithError(err)
    }

    task, err := services.Tasks.Update(ctx, params)
    if err != nil {
        cli.ExitWithError(cli.WrapError(err, "failed to update task"))
    }

    printTaskUpdateSuccess(task)
}

// Helper: Validate update flags
func validateUpdateFlags(cmd *cobra.Command) error {
    return cli.ValidateUpdateFlags(cmd, cli.UpdateFlagValues{
        Title:       taskUpdateTitle,
        Description: taskUpdateDescription,
        Status:      taskUpdateStatus,
        Priority:    taskUpdatePriority,
        StartDate:   taskUpdateStartDate,
        Deadline:    taskUpdateDeadline,
        Context:     taskUpdateContext,
        Duration:    taskUpdateDuration,
        Next:        taskUpdateNext,
        NoNext:      taskUpdateNoNext,
    })
}

// Helper: Prepare update parameters
func prepareTaskUpdateParams(cmd *cobra.Command, existing *domain.Task) (service.UpdateTaskParams, error) {
    params := service.UpdateTaskParams{
        ID:                existing.ID,
        Title:             existing.Title,
        Description:       existing.Description,
        StartDate:         existing.StartDate,
        Deadline:          existing.Deadline,
        Priority:          existing.Priority,
        Context:           existing.Context,
        EstimatedDuration: existing.EstimatedDuration,
        Status:            existing.Status,
        IsNext:            existing.IsNext,
    }

    if taskUpdateTitle != "" {
        params.Title = taskUpdateTitle
    }
    if taskUpdateStatus != "" {
        status, err := cli.ParseTaskStatus(taskUpdateStatus)
        if err != nil {
            return service.UpdateTaskParams{}, err
        }
        params.Status = status
    }
    // ... more field updates

    return params, nil
}

// Helper: Print success message
func printTaskUpdateSuccess(task *domain.Task) {
    formatter, err := GetFormatter()
    if err != nil {
        cli.ExitWithError(err)
    }

    if jsonFormatter, ok := formatter.(*output.JSONFormatter); ok {
        if err := jsonFormatter.PrintObject(domainTaskToMap(*task)); err != nil {
            cli.ExitWithError(cli.WrapError(err, "failed to output task"))
        }
    } else {
        nextFlag := ""
        if task.IsNext {
            nextFlag = " [NEXT]"
        }
        output.PrintSuccess(fmt.Sprintf("Task updated: %s%s", task.ID, nextFlag))
    }
}
```

### Benefits of This Pattern

1. **Reduced Complexity**: Each function has a single responsibility
2. **Testability**: Helper functions can be unit tested independently
3. **Readability**: Main command logic is easier to follow
4. **Maintainability**: Changes to validation/output don't affect core logic

### When to Extract

- Function exceeds gocyclo limit of 30
- Multiple nested conditionals
- Repeated patterns across commands
- Complex flag validation logic

---

## Reusable Helper Patterns

### Delete Helper Pattern

When multiple commands share identical logic patterns, extract them into reusable helpers with interfaces. This reduces code duplication and improves maintainability.

#### Interface-Based Abstraction

Define an interface that captures the common operations:

```go
// internal/cli/delete.go

// Deleteable defines the interface for entities that support soft and hard deletion.
type Deleteable interface {
    // GetByID retrieves an entity to verify its existence before deletion.
    GetByID(ctx context.Context, id string) (any, error)

    // SoftDelete marks an entity as deleted without removing it from the database.
    SoftDelete(ctx context.Context, id string) error

    // HardDelete permanently removes an entity from the database.
    HardDelete(ctx context.Context, id string) error
}

// DeleteParams contains the parameters needed for the RunDelete helper.
type DeleteParams struct {
    ID          string
    Permanent   bool
    EntityName  string
    NotFoundErr error
}

// RunDelete executes a delete operation (soft or hard) on an entity using the provided service.
func RunDelete(ctx context.Context, svc Deleteable, params DeleteParams) error {
    // Verify entity exists
    _, err := svc.GetByID(ctx, params.ID)
    if err != nil {
        if params.NotFoundErr != nil && err == params.NotFoundErr {
            return fmt.Errorf("%s not found: %s", params.EntityName, params.ID)
        }
        return WrapError(err, fmt.Sprintf("failed to get %s", params.EntityName))
    }

    // Perform deletion
    if params.Permanent {
        if err := svc.HardDelete(ctx, params.ID); err != nil {
            return WrapError(err, fmt.Sprintf("failed to permanently delete %s", params.EntityName))
        }
        output.PrintSuccess(fmt.Sprintf("%s permanently deleted: %s", params.EntityName, params.ID))
        return nil
    }

    if err := svc.SoftDelete(ctx, params.ID); err != nil {
        if params.NotFoundErr != nil && err == params.NotFoundErr {
            return fmt.Errorf("%s not found: %s", params.EntityName, params.ID)
        }
        return WrapError(err, fmt.Sprintf("failed to delete %s", params.EntityName))
    }

    output.PrintSuccess(fmt.Sprintf("%s deleted: %s", params.EntityName, params.ID))
    return nil
}
```

#### Using the Helper in Commands

Implement the interface for each entity type and use the helper:

```go
// cmd/dopa/projects.go

// projectDeleter wraps ProjectService to implement Deleteable interface
type projectDeleter struct {
    *service.ProjectService
}

func (p *projectDeleter) GetByID(ctx context.Context, id string) (any, error) {
    return p.ProjectService.GetByID(ctx, id)
}

func runProjectsDelete(cmd *cobra.Command, args []string) {
    id := args[0]
    permanent, _ := cmd.Flags().GetBool("permanent")

    services, err := GetServices()
    if err != nil {
        cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
    }
    defer cli.CloseWithLog(services, "services")

    ctx := context.Background()
    params := cli.DeleteParams{
        ID:          id,
        Permanent:   permanent,
        EntityName:  "project",
        NotFoundErr: service.ErrProjectNotFound,
    }

    if err := cli.RunDelete(ctx, &projectDeleter{services.Projects}, params); err != nil {
        cli.ExitWithError(err)
    }
}
```

#### Benefits of This Pattern

1. **Eliminates Duplicate Code**: Removes 40+ lines of duplicate logic from each delete command
2. **Consistent Error Handling**: All delete operations handle errors the same way
3. **Testability**: Helper logic is tested once, not duplicated across multiple command tests
4. **Maintainability**: Bug fixes and improvements apply to all entity types
5. **Linting Compliance**: Resolves `dupl` warnings about duplicate code

#### When to Use

- Multiple commands share identical logic patterns
- Code duplication detected by linters (e.g., `dupl`, `goconst`)
- Common operations across entity types (CRUD operations, validation, formatting)
- Before/After hooks that repeat across commands

#### Related Patterns

- **Validator Helpers**: Common input validation logic
- **Formatter Helpers**: Shared output formatting patterns
- **Middleware/Interceptors**: Cross-cutting concerns like logging, metrics

---

## Key Files

| File | Purpose |
|------|---------|
| `cmd/dopa/main.go` | Entry point, service container, root command |
| `cmd/dopa/tasks.go` | Task commands |
| `cmd/dopa/projects.go` | Project commands |
| `cmd/dopa/areas.go` | Area commands |
| `cmd/dopa/subareas.go` | Subarea commands |
| `internal/cli/output/formatter.go` | Output formatting (table/JSON/YAML) |
| `internal/cli/filter/parser.go` | Filter parsing |
| `internal/cli/errors.go` | CLI error handling |
| `internal/cli/helpers.go` | Cleanup helpers (CloseWithLog) |
| `internal/cli/delete.go` | Delete helper (Deleteable interface, RunDelete) |

---

**Navigation**: [← Repository Layer](05-repository-layer.md) | [Back to Architecture →](README.md)
