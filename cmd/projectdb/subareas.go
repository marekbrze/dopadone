package main

import (
	"context"
	"fmt"
	"time"

	"github.com/example/projectdb/internal/cli"
	"github.com/example/projectdb/internal/cli/filter"
	"github.com/example/projectdb/internal/cli/output"
	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/service"
	"github.com/spf13/cobra"
)

var subareasCmd = &cobra.Command{
	Use:     "subareas",
	Short:   "Manage subareas",
	Long:    "Manage subareas in the project database. Subareas are subdivisions of areas.",
	Aliases: []string{"subarea", "sub"},
}

var subareasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new subarea",
	Long:  "Create a new subarea under a specified area.",
	Example: `  # Create a subarea with required fields
  projectdb subareas create --name "Backend" --area-id "area-123"

  # Create a subarea with a color
  projectdb subareas create --name "Frontend" --area-id "area-123" --color "#FF5733"`,
	Run: runSubareasCreate,
}

var subareasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subareas",
	Long:  "List all subareas, optionally filtered by area ID.",
	Example: `  # List all subareas
  projectdb subareas list

  # List subareas for a specific area
  projectdb subareas list --area-id "area-123"

  # Output as JSON
  projectdb subareas list --json
  projectdb subareas list --format=json

  # Output as YAML
  projectdb subareas list --format=yaml

  # Filter subareas
  projectdb subareas list --filter 'name=Backend'`,
	Run: runSubareasList,
}

var subareasGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a subarea by ID",
	Long:  "Display details of a specific subarea by its ID.",
	Example: `  # Get a subarea by ID
  projectdb subareas get "subarea-123"`,
	Args: cobra.ExactArgs(1),
	Run:  runSubareasGet,
}

var subareasUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a subarea",
	Long:  "Update a subarea's name and/or color.",
	Example: `  # Update a subarea's name
  projectdb subareas update "subarea-123" --name "New Name"

  # Update a subarea's color
  projectdb subareas update "subarea-123" --color "#00FF00"

  # Update both name and color
  projectdb subareas update "subarea-123" --name "Updated" --color "#0000FF"`,
	Args: cobra.ExactArgs(1),
	Run:  runSubareasUpdate,
}

var subareasDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a subarea",
	Long:  "Delete a subarea by ID. By default performs a soft delete.",
	Example: `  # Soft delete a subarea (can be recovered)
  projectdb subareas delete "subarea-123"

  # Permanently delete a subarea
  projectdb subareas delete "subarea-123" --permanent`,
	Args: cobra.ExactArgs(1),
	Run:  runSubareasDelete,
}

var (
	subareaName       string
	subareaAreaID     string
	subareaColor      string
	subareaListAreaID string
	subareaPermanent  bool
	subareaJSON       bool
	subareaFormat     string
	subareaFilter     string
)

func init() {
	subareasCmd.AddCommand(subareasCreateCmd)
	subareasCmd.AddCommand(subareasListCmd)
	subareasCmd.AddCommand(subareasGetCmd)
	subareasCmd.AddCommand(subareasUpdateCmd)
	subareasCmd.AddCommand(subareasDeleteCmd)

	subareasCreateCmd.Flags().StringVar(&subareaName, "name", "", "subarea name (required)")
	subareasCreateCmd.Flags().StringVar(&subareaAreaID, "area-id", "", "parent area ID (required)")
	subareasCreateCmd.Flags().StringVar(&subareaColor, "color", "", "color in hex format (e.g., #FF5733)")
	subareasCreateCmd.MarkFlagRequired("name")
	subareasCreateCmd.MarkFlagRequired("area-id")

	subareasListCmd.Flags().StringVar(&subareaListAreaID, "area-id", "", "filter by area ID")
	subareasListCmd.Flags().BoolVar(&subareaJSON, "json", false, "output as JSON")
	subareasListCmd.Flags().StringVar(&subareaFormat, "format", "table", "output format (table|json|yaml)")
	subareasListCmd.Flags().StringVar(&subareaFilter, "filter", "", "filter expression (e.g., 'name=Backend')")

	subareasUpdateCmd.Flags().StringVar(&subareaName, "name", "", "new subarea name")
	subareasUpdateCmd.Flags().StringVar(&subareaColor, "color", "", "new color in hex format (e.g., #FF5733)")

	subareasDeleteCmd.Flags().BoolVar(&subareaPermanent, "permanent", false, "permanently delete (cannot be recovered)")
}

func runSubareasCreate(cmd *cobra.Command, args []string) {
	if subareaName == "" {
		cli.ExitWithError(cli.NewValidationError("name", "subarea name cannot be empty"))
	}
	if subareaAreaID == "" {
		cli.ExitWithError(cli.NewValidationError("area-id", "area ID cannot be empty"))
	}

	color, err := cli.ParseColor(subareaColor)
	if err != nil {
		cli.ExitWithError(err)
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	subarea, err := services.Subareas.Create(ctx, subareaName, subareaAreaID, color)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to create subarea"))
	}

	formatter, err := GetFormatter()
	if err != nil {
		cli.ExitWithError(err)
	}

	if jsonFormatter, ok := formatter.(*output.JSONFormatter); ok {
		if err := jsonFormatter.PrintObject(domainSubareaToMap(*subarea)); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output subarea"))
		}
	} else {
		output.PrintSuccess(fmt.Sprintf("Subarea created with ID: %s", subarea.ID))
	}
}

func runSubareasList(cmd *cobra.Command, args []string) {
	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	var subareas []domain.Subarea

	if subareaListAreaID != "" {
		subareas, err = services.Subareas.ListByArea(ctx, subareaListAreaID)
	} else {
		subareas, err = services.Subareas.ListAll(ctx)
	}

	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to list subareas"))
	}

	var subareaMaps []map[string]interface{}
	for _, s := range subareas {
		subareaMaps = append(subareaMaps, domainSubareaToMap(s))
	}

	if subareaFilter != "" {
		subareaMaps, err = filter.EvaluateFilter(subareaFilter, subareaMaps)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to apply filter"))
		}
	}

	useJSON := subareaJSON || subareaFormat == "json"
	useYAML := subareaFormat == "yaml"

	if useJSON {
		formatter := output.NewJSONFormatter()
		for _, m := range subareaMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output JSON"))
		}
		return
	}

	if useYAML {
		formatter := output.NewYAMLFormatter()
		for _, m := range subareaMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output YAML"))
		}
		return
	}

	if len(subareaMaps) == 0 {
		output.PrintInfo("No subareas found")
		return
	}

	formatter, _ := GetFormatter()
	formatter.PrintHeader([]string{"ID", "NAME", "AREA ID", "COLOR", "CREATED"})
	for _, m := range subareaMaps {
		formatter.PrintRow([]string{
			m["id"].(string),
			m["name"].(string),
			m["area_id"].(string),
			colorStrOrNil(m, "color"),
			m["created_at"].(string)[:10],
		})
	}
	formatter.Flush()
}

func runSubareasGet(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	subarea, err := services.Subareas.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrSubareaNotFound {
			cli.ExitWithError(fmt.Errorf("subarea not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get subarea"))
	}

	formatter := output.NewJSONFormatter()
	if err := formatter.PrintObject(domainSubareaToMap(*subarea)); err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to output subarea"))
	}
}

func runSubareasUpdate(cmd *cobra.Command, args []string) {
	id := args[0]

	if subareaName == "" && subareaColor == "" {
		cli.ExitWithError(cli.NewValidationError("", "at least one of --name or --color must be provided"))
	}

	color, err := cli.ParseColor(subareaColor)
	if err != nil {
		cli.ExitWithError(err)
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	existing, err := services.Subareas.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrSubareaNotFound {
			cli.ExitWithError(fmt.Errorf("subarea not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get subarea"))
	}

	newName := existing.Name
	if subareaName != "" {
		newName = subareaName
	}

	newColor := existing.Color
	if subareaColor != "" {
		newColor = color
	}

	subarea, err := services.Subareas.Update(ctx, id, newName, existing.AreaID, newColor)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to update subarea"))
	}

	output.PrintSuccess(fmt.Sprintf("Subarea updated: %s", subarea.ID))
}

func runSubareasDelete(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	_, err = services.Subareas.GetByID(ctx, id)
	if err != nil {
		if err == service.ErrSubareaNotFound {
			cli.ExitWithError(fmt.Errorf("subarea not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to get subarea"))
	}

	if subareaPermanent {
		err := services.Subareas.HardDelete(ctx, id)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to permanently delete subarea"))
		}
		output.PrintSuccess(fmt.Sprintf("Subarea permanently deleted: %s", id))
		return
	}

	err = services.Subareas.SoftDelete(ctx, id)
	if err != nil {
		if err == service.ErrSubareaNotFound {
			cli.ExitWithError(fmt.Errorf("subarea not found: %s", id))
		}
		cli.ExitWithError(cli.WrapError(err, "failed to delete subarea"))
	}

	output.PrintSuccess(fmt.Sprintf("Subarea deleted: %s", id))
}

func domainSubareaToMap(s domain.Subarea) map[string]interface{} {
	result := map[string]interface{}{
		"id":         s.ID,
		"name":       s.Name,
		"area_id":    s.AreaID,
		"created_at": s.CreatedAt.Format(time.RFC3339),
		"updated_at": s.UpdatedAt.Format(time.RFC3339),
	}

	if s.Color != "" {
		result["color"] = string(s.Color)
	}

	return result
}

func colorStrOrNil(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
