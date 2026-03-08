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

var projectsCmd = &cobra.Command{
	Use:     "projects",
	Short:   "Manage projects",
	Long:    "Manage projects in the project database. Projects can be nested under subareas or other projects.",
	Aliases: []string{"project", "proj"},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  "Create a new project. Must specify either --subarea-id OR --parent-id (not both).",
	Example: `  # Create a project under a subarea (root project)
  dopa projects create --name "Website Redesign" --subarea-id "subarea-123"

  # Create a nested project under another project
  dopa projects create --name "Backend API" --parent-id "project-456" --priority high

  # Create a project with all optional fields
  dopa projects create --name "Q4 Campaign" --subarea-id "subarea-123" \
    --status active --priority urgent --progress 25 \
    --start-date 2024-10-01 --deadline 2024-12-31 \
    --color "#FF5733" --goal "Launch campaign by year end" \
    --description "Marketing campaign for Q4"`,
	Run: runProjectsCreate,
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "List all projects, optionally filtered by status, priority, subarea, or parent.",
	Example: `  # List all projects
  dopa projects list

  # List projects by status
  dopa projects list --status active

  # List high priority projects
  dopa projects list --priority high

  # List projects under a subarea
  dopa projects list --subarea-id "subarea-123"

  # List nested projects under a parent
  dopa projects list --parent-id "project-456"

  # Output as JSON
  dopa projects list --json

  # Output as YAML
  dopa projects list --format=yaml

  # Filter with query syntax
  dopa projects list --filter 'status=active AND priority>=high'`,
	Run: runProjectsList,
}

var projectsGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a project by ID",
	Long:  "Display details of a specific project by its ID.",
	Example: `  # Get a project by ID
  dopa projects get "project-123"`,
	Args: cobra.ExactArgs(1),
	Run:  runProjectsGet,
}

var projectsUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a project",
	Long:  "Update a project's fields. All fields are optional.",
	Example: `  # Update a project's name
  dopa projects update "project-123" --name "New Project Name"

  # Update status and priority
  dopa projects update "project-123" --status completed --priority high

  # Update progress and deadline
  dopa projects update "project-123" --progress 75 --deadline 2024-12-31

  # Update multiple fields
  dopa projects update "project-123" --name "Updated" --status on_hold --color "#00FF00"`,
	Args: cobra.ExactArgs(1),
	Run:  runProjectsUpdate,
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a project",
	Long:  "Delete a project by ID. By default performs a soft delete.",
	Example: `  # Soft delete a project (can be recovered)
  dopa projects delete "project-123"

  # Permanently delete a project
  dopa projects delete "project-123" --permanent`,
	Args: cobra.ExactArgs(1),
	Run:  runProjectsDelete,
}

var (
	projCreateName        string
	projCreateSubareaID   string
	projCreateParentID    string
	projCreateStatus      string
	projCreatePriority    string
	projCreateProgress    int
	projCreateDeadline    string
	projCreateStartDate   string
	projCreateColor       string
	projCreateGoal        string
	projCreateDescription string

	projUpdateName        string
	projUpdateSubareaID   string
	projUpdateParentID    string
	projUpdateStatus      string
	projUpdatePriority    string
	projUpdateProgress    int
	projUpdateDeadline    string
	projUpdateStartDate   string
	projUpdateColor       string
	projUpdateGoal        string
	projUpdateDescription string

	projListStatus    string
	projListPriority  string
	projListSubareaID string
	projListParentID  string
	projListJSON      bool
	projListFormat    string
	projListFilter    string

	projPermanent bool
)

func init() {
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsGetCmd)
	projectsCmd.AddCommand(projectsUpdateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)

	projectsCreateCmd.Flags().StringVar(&projCreateName, "name", "", "project name (required)")
	projectsCreateCmd.Flags().StringVar(&projCreateSubareaID, "subarea-id", "", "parent subarea ID (required if parent-id not set)")
	projectsCreateCmd.Flags().StringVar(&projCreateParentID, "parent-id", "", "parent project ID (required if subarea-id not set)")
	projectsCreateCmd.Flags().StringVar(&projCreateStatus, "status", "active", "project status (active|completed|on_hold|archived)")
	projectsCreateCmd.Flags().StringVar(&projCreatePriority, "priority", "medium", "project priority (low|medium|high|urgent)")
	projectsCreateCmd.Flags().IntVar(&projCreateProgress, "progress", 0, "completion percentage (0-100)")
	projectsCreateCmd.Flags().StringVar(&projCreateDeadline, "deadline", "", "deadline date (YYYY-MM-DD)")
	projectsCreateCmd.Flags().StringVar(&projCreateStartDate, "start-date", "", "start date (YYYY-MM-DD)")
	projectsCreateCmd.Flags().StringVar(&projCreateColor, "color", "", "color in hex format (e.g., #FF5733)")
	projectsCreateCmd.Flags().StringVar(&projCreateGoal, "goal", "", "project goal/outcome")
	projectsCreateCmd.Flags().StringVar(&projCreateDescription, "description", "", "project description")
	projectsCreateCmd.MarkFlagRequired("name")

	projectsListCmd.Flags().StringVar(&projListStatus, "status", "", "filter by status")
	projectsListCmd.Flags().StringVar(&projListPriority, "priority", "", "filter by priority")
	projectsListCmd.Flags().StringVar(&projListSubareaID, "subarea-id", "", "filter by subarea ID")
	projectsListCmd.Flags().StringVar(&projListParentID, "parent-id", "", "filter by parent project ID")
	projectsListCmd.Flags().BoolVar(&projListJSON, "json", false, "output as JSON")
	projectsListCmd.Flags().StringVar(&projListFormat, "format", "table", "output format (table|json|yaml)")
	projectsListCmd.Flags().StringVar(&projListFilter, "filter", "", "filter expression (e.g., 'status=active AND priority>=high')")

	projectsUpdateCmd.Flags().StringVar(&projUpdateName, "name", "", "new project name")
	projectsUpdateCmd.Flags().StringVar(&projUpdateStatus, "status", "", "new status (active|completed|on_hold|archived)")
	projectsUpdateCmd.Flags().StringVar(&projUpdatePriority, "priority", "", "new priority (low|medium|high|urgent)")
	projectsUpdateCmd.Flags().IntVar(&projUpdateProgress, "progress", 0, "new completion percentage (0-100)")
	projectsUpdateCmd.Flags().StringVar(&projUpdateDeadline, "deadline", "", "new deadline date (YYYY-MM-DD)")
	projectsUpdateCmd.Flags().StringVar(&projUpdateStartDate, "start-date", "", "new start date (YYYY-MM-DD)")
	projectsUpdateCmd.Flags().StringVar(&projUpdateColor, "color", "", "new color in hex format")
	projectsUpdateCmd.Flags().StringVar(&projUpdateGoal, "goal", "", "new goal")
	projectsUpdateCmd.Flags().StringVar(&projUpdateDescription, "description", "", "new description")
	projectsUpdateCmd.Flags().StringVar(&projUpdateSubareaID, "subarea-id", "", "new subarea ID")
	projectsUpdateCmd.Flags().StringVar(&projUpdateParentID, "parent-id", "", "new parent project ID")

	projectsDeleteCmd.Flags().BoolVar(&projPermanent, "permanent", false, "permanently delete (cannot be recovered)")
}

func runProjectsCreate(cmd *cobra.Command, args []string) {
	if err := cli.ValidateProjectName(projCreateName); err != nil {
		cli.ExitWithError(err)
	}

	if projCreateSubareaID == "" && projCreateParentID == "" {
		cli.ExitWithError(cli.NewValidationError("", "either --subarea-id or --parent-id must be provided"))
	}
	if projCreateSubareaID != "" && projCreateParentID != "" {
		cli.ExitWithError(cli.NewValidationError("", "cannot specify both --subarea-id and --parent-id, choose one"))
	}

	status, err := cli.ParseProjectStatus(projCreateStatus)
	if err != nil {
		cli.ExitWithError(err)
	}

	priority, err := cli.ParsePriority(projCreatePriority)
	if err != nil {
		cli.ExitWithError(err)
	}

	progress, err := cli.ParseProgress(projCreateProgress)
	if err != nil {
		cli.ExitWithError(err)
	}

	color, err := cli.ParseColor(projCreateColor)
	if err != nil {
		cli.ExitWithError(err)
	}

	startDate, deadline, err := cli.ParseDate(projCreateStartDate, projCreateDeadline)
	if err != nil {
		cli.ExitWithError(err)
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	var parentIDPtr *string
	if projCreateParentID != "" {
		parentIDPtr = &projCreateParentID
	}

	var subareaIDPtr *string
	if projCreateSubareaID != "" {
		subareaIDPtr = &projCreateSubareaID
	}

	params := service.CreateProjectParams{
		Name:        projCreateName,
		Description: projCreateDescription,
		Goal:        projCreateGoal,
		Status:      domain.ProjectStatus(status),
		Priority:    domain.Priority(priority),
		Progress:    domain.Progress(progress.Int()),
		StartDate:   startDate,
		Deadline:    deadline,
		Color:       domain.Color(color),
		ParentID:    parentIDPtr,
		SubareaID:   subareaIDPtr,
		Position:    0,
	}

	project, err := services.Projects.Create(ctx, params)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to create project"))
	}

	formatter, err := GetFormatter()
	if err != nil {
		cli.ExitWithError(err)
	}

	if jsonFormatter, ok := formatter.(*output.JSONFormatter); ok {
		if err := jsonFormatter.PrintObject(domainProjectToMap(*project)); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output project"))
		}
	} else {
		output.PrintSuccess(fmt.Sprintf("Project created with ID: %s", project.ID))
	}
}

func runProjectsList(cmd *cobra.Command, args []string) {
	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	var projects []domain.Project

	if projListStatus != "" {
		status, err := cli.ParseProjectStatus(projListStatus)
		if err != nil {
			cli.ExitWithError(err)
		}
		projects, err = services.Projects.ListByStatus(ctx, status)
	} else if projListPriority != "" {
		priority, err := cli.ParsePriority(projListPriority)
		if err != nil {
			cli.ExitWithError(err)
		}
		projects, err = services.Projects.ListByPriority(ctx, priority)
	} else if projListSubareaID != "" {
		projects, err = services.Projects.ListBySubarea(ctx, projListSubareaID)
	} else if projListParentID != "" {
		projects, err = services.Projects.ListByParent(ctx, projListParentID)
	} else {
		projects, err = services.Projects.ListAll(ctx)
	}

	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to list projects"))
	}

	var projectMaps []map[string]interface{}
	for _, p := range projects {
		projectMaps = append(projectMaps, domainProjectToMap(p))
	}

	if projListFilter != "" {
		projectMaps, err = filter.EvaluateFilter(projListFilter, projectMaps)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to apply filter"))
		}
	}

	useJSON := projListJSON || projListFormat == "json"
	useYAML := projListFormat == "yaml"

	if useJSON {
		formatter := output.NewJSONFormatter()
		for _, m := range projectMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output JSON"))
		}
		return
	}

	if useYAML {
		formatter := output.NewYAMLFormatter()
		for _, m := range projectMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output YAML"))
		}
		return
	}

	if len(projectMaps) == 0 {
		output.PrintInfo("No projects found")
		return
	}

	formatter, _ := GetFormatter()
	formatter.PrintHeader([]string{"ID", "NAME", "STATUS", "PRIORITY", "PROGRESS", "DEADLINE"})
	for _, m := range projectMaps {
		deadlineStr := ""
		if v, ok := m["deadline"]; ok && v != nil {
			deadlineStr = v.(string)
		}
		formatter.PrintRow([]string{
			m["id"].(string),
			m["name"].(string),
			m["status"].(string),
			m["priority"].(string),
			fmt.Sprintf("%d%%", m["progress"]),
			deadlineStr,
		})
	}
	formatter.Flush()
}

func runProjectsGet(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	project, err := services.Projects.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrProjectNotFound {
			cli.ExitWithError(fmt.Errorf("project not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get project"))
	}

	formatter := output.NewJSONFormatter()
	if err := formatter.PrintObject(domainProjectToMap(*project)); err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to output project"))
	}
}

func runProjectsUpdate(cmd *cobra.Command, args []string) {
	id := args[0]

	flags := []string{projUpdateName, projUpdateStatus, projUpdatePriority, projUpdateColor, projUpdateGoal, projUpdateDescription, projUpdateStartDate, projUpdateDeadline, projUpdateSubareaID, projUpdateParentID}
	hasChanges := false
	for _, f := range flags {
		if f != "" {
			hasChanges = true
			break
		}
	}
	if !hasChanges && cmd.Flags().Changed("progress") {
		hasChanges = true
	}

	if !hasChanges {
		cli.ExitWithError(cli.NewValidationError("", "at least one field must be provided to update"))
	}

	startDate, deadline, err := cli.ParseDate(projUpdateStartDate, projUpdateDeadline)
	if err != nil {
		cli.ExitWithError(err)
	}

	var parentID *string
	var subareaID *string
	if projUpdateSubareaID != "" && projUpdateParentID != "" {
		cli.ExitWithError(cli.NewValidationError("", "cannot specify both --subarea-id and --parent-id, choose one"))
	}
	if projUpdateSubareaID != "" {
		subareaID = &projUpdateSubareaID
		parentID = nil
	} else if projUpdateParentID != "" {
		parentID = &projUpdateParentID
		subareaID = nil
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	existing, err := services.Projects.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrProjectNotFound {
			cli.ExitWithError(fmt.Errorf("project not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get project"))
	}

	name := projUpdateName
	if name == "" {
		name = existing.Name
	}

	status := existing.Status
	if projUpdateStatus != "" {
		parsedStatus, err := cli.ParseProjectStatus(projUpdateStatus)
		if err != nil {
			cli.ExitWithError(err)
		}
		status = parsedStatus
	}

	priority := existing.Priority
	if projUpdatePriority != "" {
		parsedPriority, err := cli.ParsePriority(projUpdatePriority)
		if err != nil {
			cli.ExitWithError(err)
		}
		priority = parsedPriority
	}

	progress := existing.Progress
	if cmd.Flags().Changed("progress") {
		parsedProgress, err := cli.ParseProgress(projUpdateProgress)
		if err != nil {
			cli.ExitWithError(err)
		}
		progress = parsedProgress
	}

	color := existing.Color
	if projUpdateColor != "" {
		parsedColor, err := cli.ParseColor(projUpdateColor)
		if err != nil {
			cli.ExitWithError(err)
		}
		color = parsedColor
	}

	params := service.UpdateProjectParams{
		ID:          id,
		Name:        name,
		Description: projUpdateDescription,
		Goal:        projUpdateGoal,
		Status:      status,
		Priority:    priority,
		Progress:    progress,
		StartDate:   startDate,
		Deadline:    deadline,
		Color:       color,
		ParentID:    parentID,
		SubareaID:   subareaID,
	}

	project, err := services.Projects.Update(ctx, params)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to update project"))
	}

	output.PrintSuccess(fmt.Sprintf("Project updated: %s", project.ID))
}

func runProjectsDelete(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	_, err = services.Projects.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrProjectNotFound {
			cli.ExitWithError(fmt.Errorf("project not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get project"))
	}

	if projPermanent {
		err := services.Projects.HardDelete(ctx, id)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to permanently delete project"))
		}
		output.PrintSuccess(fmt.Sprintf("Project permanently deleted: %s", id))
		return
	}

	err = services.Projects.SoftDelete(ctx, id)
	if err != nil {
		if err == service.ErrProjectNotFound {
			cli.ExitWithError(fmt.Errorf("project not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to delete project"))
	}

	output.PrintSuccess(fmt.Sprintf("Project deleted: %s", id))
}

func domainProjectToMap(p domain.Project) map[string]interface{} {
	result := map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"status":     p.Status.String(),
		"priority":   p.Priority.String(),
		"progress":   p.Progress.Int(),
		"position":   p.Position,
		"created_at": p.CreatedAt.Format(time.RFC3339),
		"updated_at": p.UpdatedAt.Format(time.RFC3339),
	}

	if p.Description != "" {
		result["description"] = p.Description
	}
	if p.Goal != "" {
		result["goal"] = p.Goal
	}
	if p.Deadline != nil {
		result["deadline"] = p.Deadline.Format("2006-01-02")
	}
	if p.Color != "" {
		result["color"] = string(p.Color)
	}
	if p.ParentID != nil {
		result["parent_id"] = *p.ParentID
	}
	if p.SubareaID != nil {
		result["subarea_id"] = *p.SubareaID
	}
	if p.CompletedAt != nil {
		result["completed_at"] = p.CompletedAt.Format(time.RFC3339)
	}

	return result
}
