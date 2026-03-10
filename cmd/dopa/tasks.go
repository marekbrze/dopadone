package main

import (
	"context"
	"fmt"
	"time"

	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/cli/filter"
	"github.com/marekbrze/dopadone/internal/cli/output"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/service"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:     "tasks",
	Short:   "Manage tasks",
	Long:    "Manage tasks in the project database. Tasks belong to projects and can be prioritized.",
	Aliases: []string{"task"},
}

var tasksCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	Long:  "Create a new task under a project.",
	Example: `  # Create a basic task
  dopa tasks create --project-id "proj-123" --title "Write documentation"

  # Create a task with all options
  dopa tasks create --project-id "proj-123" --title "API Integration" \
    --description "Integrate with external API" \
    --status in_progress --priority high \
    --start-date 2024-01-15 --deadline 2024-01-31 \
    --context "backend" --duration 60 --next`,
	Run: runTasksCreate,
}

var tasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  "List all tasks, optionally filtered by project, status, priority, or next flag.",
	Example: `  # List all tasks
  dopa tasks list

  # List tasks by project
  dopa tasks list --project-id "proj-123"

  # List tasks marked as next
  dopa tasks list --next

  # List tasks by status
  dopa tasks list --status in_progress

  # Output as JSON
  dopa tasks list --json`,
	Run: runTasksList,
}

var tasksNextCmd = &cobra.Command{
	Use:   "next",
	Short: "List tasks marked as 'next'",
	Long:  "List all tasks marked with the --next flag (priority/focused tasks).",
	Example: `  # List all next tasks
  dopa tasks next`,
	Run: runTasksNext,
}

var tasksGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a task by ID",
	Long:  "Display details of a specific task by its ID.",
	Example: `  # Get a task by ID
  dopa tasks get "task-123"`,
	Args: cobra.ExactArgs(1),
	Run:  runTasksGet,
}

var tasksUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task",
	Long:  "Update a task's fields. All fields are optional.",
	Example: `  # Update task title
  dopa tasks update "task-123" --title "New title"

  # Update status and priority
  dopa tasks update "task-123" --status done --priority critical

  # Mark task as next
  dopa tasks update "task-123" --next

  # Remove next flag
  dopa tasks update "task-123" --no-next`,
	Args: cobra.ExactArgs(1),
	Run:  runTasksUpdate,
}

var tasksDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Long:  "Delete a task by ID. By default performs a soft delete.",
	Example: `  # Soft delete a task (can be recovered)
  dopa tasks delete "task-123"

  # Permanently delete a task
  dopa tasks delete "task-123" --permanent`,
	Args: cobra.ExactArgs(1),
	Run:  runTasksDelete,
}

var (
	taskCreateProjectID   string
	taskCreateTitle       string
	taskCreateDescription string
	taskCreateStatus      string
	taskCreatePriority    string
	taskCreateStartDate   string
	taskCreateDeadline    string
	taskCreateContext     string
	taskCreateDuration    int
	taskCreateNext        bool

	taskUpdateTitle       string
	taskUpdateDescription string
	taskUpdateStatus      string
	taskUpdatePriority    string
	taskUpdateStartDate   string
	taskUpdateDeadline    string
	taskUpdateContext     string
	taskUpdateDuration    int
	taskUpdateNext        bool
	taskUpdateNoNext      bool

	taskListProjectID string
	taskListStatus    string
	taskListPriority  string
	taskListNext      bool
	taskListJSON      bool
	taskListFormat    string
	taskListFilter    string

	taskPermanent bool
)

func init() {
	tasksCmd.AddCommand(tasksCreateCmd)
	tasksCmd.AddCommand(tasksListCmd)
	tasksCmd.AddCommand(tasksNextCmd)
	tasksCmd.AddCommand(tasksGetCmd)
	tasksCmd.AddCommand(tasksUpdateCmd)
	tasksCmd.AddCommand(tasksDeleteCmd)

	tasksCreateCmd.Flags().StringVar(&taskCreateProjectID, "project-id", "", "parent project ID (required)")
	tasksCreateCmd.Flags().StringVar(&taskCreateTitle, "title", "", "task title (required)")
	tasksCreateCmd.Flags().StringVar(&taskCreateDescription, "description", "", "task description")
	tasksCreateCmd.Flags().StringVar(&taskCreateStatus, "status", "todo", "task status (todo|in_progress|waiting|done)")
	tasksCreateCmd.Flags().StringVar(&taskCreatePriority, "priority", "medium", "task priority (critical|high|medium|low)")
	tasksCreateCmd.Flags().StringVar(&taskCreateStartDate, "start-date", "", "start date (YYYY-MM-DD)")
	tasksCreateCmd.Flags().StringVar(&taskCreateDeadline, "deadline", "", "deadline date (YYYY-MM-DD)")
	tasksCreateCmd.Flags().StringVar(&taskCreateContext, "context", "", "task context (e.g., 'backend', 'frontend')")
	tasksCreateCmd.Flags().IntVar(&taskCreateDuration, "duration", 0, "estimated duration in minutes (5|15|30|60|120|240|480)")
	tasksCreateCmd.Flags().BoolVar(&taskCreateNext, "next", false, "mark as priority/next task")
	if err := tasksCreateCmd.MarkFlagRequired("project-id"); err != nil {
		panic(err)
	}
	if err := tasksCreateCmd.MarkFlagRequired("title"); err != nil {
		panic(err)
	}

	tasksListCmd.Flags().StringVar(&taskListProjectID, "project-id", "", "filter by project ID")
	tasksListCmd.Flags().StringVar(&taskListStatus, "status", "", "filter by status")
	tasksListCmd.Flags().StringVar(&taskListPriority, "priority", "", "filter by priority")
	tasksListCmd.Flags().BoolVar(&taskListNext, "next", false, "show only next tasks")
	tasksListCmd.Flags().BoolVar(&taskListJSON, "json", false, "output as JSON")
	tasksListCmd.Flags().StringVar(&taskListFormat, "format", "table", "output format (table|json|yaml)")
	tasksListCmd.Flags().StringVar(&taskListFilter, "filter", "", "filter expression")

	tasksUpdateCmd.Flags().StringVar(&taskUpdateTitle, "title", "", "new task title")
	tasksUpdateCmd.Flags().StringVar(&taskUpdateDescription, "description", "", "new description")
	tasksUpdateCmd.Flags().StringVar(&taskUpdateStatus, "status", "", "new status (todo|in_progress|waiting|done)")
	tasksUpdateCmd.Flags().StringVar(&taskUpdatePriority, "priority", "", "new priority (critical|high|medium|low)")
	tasksUpdateCmd.Flags().StringVar(&taskUpdateStartDate, "start-date", "", "new start date (YYYY-MM-DD)")
	tasksUpdateCmd.Flags().StringVar(&taskUpdateDeadline, "deadline", "", "new deadline date (YYYY-MM-DD)")
	tasksUpdateCmd.Flags().StringVar(&taskUpdateContext, "context", "", "new context")
	tasksUpdateCmd.Flags().IntVar(&taskUpdateDuration, "duration", 0, "new estimated duration in minutes")
	tasksUpdateCmd.Flags().BoolVar(&taskUpdateNext, "next", false, "mark as next task")
	tasksUpdateCmd.Flags().BoolVar(&taskUpdateNoNext, "no-next", false, "remove next flag")

	tasksDeleteCmd.Flags().BoolVar(&taskPermanent, "permanent", false, "permanently delete (cannot be recovered)")
}

func runTasksCreate(cmd *cobra.Command, args []string) {
	if err := cli.ValidateTaskProjectID(taskCreateProjectID); err != nil {
		cli.ExitWithError(err)
	}

	if err := cli.ValidateTaskTitle(taskCreateTitle); err != nil {
		cli.ExitWithError(err)
	}

	status, err := cli.ParseTaskStatus(taskCreateStatus)
	if err != nil {
		cli.ExitWithError(err)
	}

	priority, err := cli.ParseTaskPriority(taskCreatePriority)
	if err != nil {
		cli.ExitWithError(err)
	}

	startDate, deadline, err := cli.ParseDate(taskCreateStartDate, taskCreateDeadline)
	if err != nil {
		cli.ExitWithError(err)
	}

	var durationVal domain.TaskDuration
	if taskCreateDuration > 0 {
		duration, err := cli.ParseTaskDuration(taskCreateDuration)
		if err != nil {
			cli.ExitWithError(err)
		}
		durationVal = duration
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer cli.CloseWithLog(services, "services")

	ctx := context.Background()

	params := service.CreateTaskParams{
		ProjectID:         taskCreateProjectID,
		Title:             taskCreateTitle,
		Description:       taskCreateDescription,
		StartDate:         startDate,
		Deadline:          deadline,
		Priority:          priority,
		Context:           taskCreateContext,
		EstimatedDuration: durationVal,
		Status:            status,
		IsNext:            taskCreateNext,
	}

	task, err := services.Tasks.Create(ctx, params)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to create task"))
	}

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
		output.PrintSuccess(fmt.Sprintf("Task created with ID: %s%s", task.ID, nextFlag))
	}
}

func runTasksList(cmd *cobra.Command, args []string) {
	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer cli.CloseWithLog(services, "services")

	ctx := context.Background()

	var tasks []domain.Task

	if taskListNext {
		tasks, err = services.Tasks.ListNext(ctx)
	} else if taskListProjectID != "" {
		tasks, err = services.Tasks.ListByProject(ctx, taskListProjectID)
	} else if taskListStatus != "" {
		status, parseErr := cli.ParseTaskStatus(taskListStatus)
		if parseErr != nil {
			cli.ExitWithError(parseErr)
		}
		tasks, err = services.Tasks.ListByStatus(ctx, status)
	} else if taskListPriority != "" {
		priority, parseErr := cli.ParseTaskPriority(taskListPriority)
		if parseErr != nil {
			cli.ExitWithError(parseErr)
		}
		tasks, err = services.Tasks.ListByPriority(ctx, priority)
	} else {
		tasks, err = services.Tasks.ListAll(ctx)
	}

	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to list tasks"))
	}

	var taskMaps []map[string]interface{}
	for _, t := range tasks {
		taskMaps = append(taskMaps, domainTaskToMap(t))
	}

	if taskListFilter != "" {
		taskMaps, err = filter.EvaluateFilter(taskListFilter, taskMaps)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to apply filter"))
		}
	}

	useJSON := taskListJSON || taskListFormat == cli.FormatJSON
	useYAML := taskListFormat == cli.FormatYAML

	if useJSON {
		formatter := output.NewJSONFormatter()
		for _, m := range taskMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output JSON"))
		}
		return
	}

	if useYAML {
		formatter := output.NewYAMLFormatter()
		for _, m := range taskMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output YAML"))
		}
		return
	}

	if len(taskMaps) == 0 {
		output.PrintInfo("No tasks found")
		return
	}

	formatter, _ := GetFormatter()
	formatter.PrintHeader([]string{"ID", "TITLE", "STATUS", "PRIORITY", "PROJECT", "NEXT", "DEADLINE"})
	for _, m := range taskMaps {
		deadlineStr := ""
		if v, ok := m["deadline"]; ok && v != nil {
			deadlineStr = v.(string)
		}
		nextStr := ""
		if v, ok := m["is_next"]; ok && v.(bool) {
			nextStr = "*"
		}
		formatter.PrintRow([]string{
			m["id"].(string),
			truncate(m["title"].(string), 30),
			m["status"].(string),
			m["priority"].(string),
			truncate(m["project_id"].(string), 12),
			nextStr,
			deadlineStr,
		})
	}
	if err := formatter.Flush(); err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to flush output"))
	}
}

func runTasksNext(cmd *cobra.Command, args []string) {
	taskListNext = true
	runTasksList(cmd, args)
}

func runTasksGet(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer cli.CloseWithLog(services, "services")

	ctx := context.Background()

	task, err := services.Tasks.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrTaskNotFound {
			cli.ExitWithError(fmt.Errorf("task not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get task"))
	}

	formatter := output.NewJSONFormatter()
	if err := formatter.PrintObject(domainTaskToMap(*task)); err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to output task"))
	}
}

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

func runTasksDelete(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer cli.CloseWithLog(services, "services")

	ctx := context.Background()

	params := cli.DeleteParams{
		ID:          id,
		Permanent:   taskPermanent,
		EntityName:  "task",
		NotFoundErr: service.ErrTaskNotFound,
	}
	if err := cli.RunDelete(ctx, &taskDeleter{svc: services.Tasks}, params); err != nil {
		cli.ExitWithError(err)
	}
}

type taskDeleter struct {
	svc service.TaskServiceInterface
}

func (d *taskDeleter) GetByID(ctx context.Context, id string) (any, error) {
	return d.svc.GetByID(ctx, id)
}

func (d *taskDeleter) SoftDelete(ctx context.Context, id string) error {
	return d.svc.SoftDelete(ctx, id)
}

func (d *taskDeleter) HardDelete(ctx context.Context, id string) error {
	return d.svc.HardDelete(ctx, id)
}

func domainTaskToMap(t domain.Task) map[string]interface{} {
	result := map[string]interface{}{
		"id":         t.ID,
		"project_id": t.ProjectID,
		"title":      t.Title,
		"status":     t.Status.String(),
		"priority":   t.Priority.String(),
		"is_next":    t.IsNext,
		"created_at": t.CreatedAt.Format(time.RFC3339),
		"updated_at": t.UpdatedAt.Format(time.RFC3339),
	}

	if t.Description != "" {
		result["description"] = t.Description
	}
	if t.StartDate != nil {
		result["start_date"] = t.StartDate.Format("2006-01-02")
	}
	if t.Deadline != nil {
		result["deadline"] = t.Deadline.Format("2006-01-02")
	}
	if t.Context != "" {
		result["context"] = t.Context
	}
	if t.EstimatedDuration != 0 {
		result["estimated_duration"] = t.EstimatedDuration.Int()
	}

	return result
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

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
	if taskUpdateDescription != "" {
		params.Description = taskUpdateDescription
	}
	if taskUpdateStatus != "" {
		status, err := cli.ParseTaskStatus(taskUpdateStatus)
		if err != nil {
			return service.UpdateTaskParams{}, err
		}
		params.Status = status
	}
	if taskUpdatePriority != "" {
		priority, err := cli.ParseTaskPriority(taskUpdatePriority)
		if err != nil {
			return service.UpdateTaskParams{}, err
		}
		params.Priority = priority
	}
	if taskUpdateStartDate != "" || taskUpdateDeadline != "" {
		startDate, deadline, err := cli.ParseDate(taskUpdateStartDate, taskUpdateDeadline)
		if err != nil {
			return service.UpdateTaskParams{}, err
		}
		if startDate != nil {
			params.StartDate = startDate
		}
		if deadline != nil {
			params.Deadline = deadline
		}
	}
	if taskUpdateContext != "" {
		params.Context = taskUpdateContext
	}
	if cmd.Flags().Changed("duration") {
		if taskUpdateDuration > 0 {
			duration, err := cli.ParseTaskDuration(taskUpdateDuration)
			if err != nil {
				return service.UpdateTaskParams{}, err
			}
			params.EstimatedDuration = duration
		} else {
			params.EstimatedDuration = 0
		}
	}
	if taskUpdateNext {
		params.IsNext = true
	}
	if taskUpdateNoNext {
		params.IsNext = false
	}

	return params, nil
}

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
