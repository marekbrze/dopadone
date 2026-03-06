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

---

**Navigation**: [← Repository Layer](05-repository-layer.md) | [Back to Architecture →](README.md)
