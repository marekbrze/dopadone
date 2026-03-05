package main

import (
	"context"
	"fmt"
	"time"

	"github.com/example/projectdb/internal/cli"
	"github.com/example/projectdb/internal/cli/filter"
	"github.com/example/projectdb/internal/cli/output"
	"github.com/example/projectdb/internal/domain"
	"github.com/spf13/cobra"
)

var areasCmd = &cobra.Command{
	Use:     "areas",
	Short:   "Manage areas",
	Long:    "Manage areas in the project database. Areas are top-level containers for subareas and projects.",
	Aliases: []string{"area"},
}

var areasCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new area",
	Long:  "Create a new area with a name and optional color.",
	Example: `  # Create an area with required name
  projectdb areas create --name "Engineering"

  # Create an area with a color
  projectdb areas create --name "Marketing" --color "#FF5733"`,
	Run: runAreasCreate,
}

var areasListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all areas",
	Long:  "List all areas in the database.",
	Example: `  # List all areas
  projectdb areas list

  # Output as JSON
  projectdb areas list --json
  projectdb areas list --format=json

  # Output as YAML
  projectdb areas list --format=yaml

  # Filter areas
  projectdb areas list --filter 'name=Engineering'`,
	Run: runAreasList,
}

var areasGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get an area by ID",
	Long:  "Display details of a specific area by its ID.",
	Example: `  # Get an area by ID
  projectdb areas get "area-123"`,
	Args: cobra.ExactArgs(1),
	Run:  runAreasGet,
}

var areasUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an area",
	Long:  "Update an area's name and/or color.",
	Example: `  # Update an area's name
  projectdb areas update "area-123" --name "New Name"

  # Update an area's color
  projectdb areas update "area-123" --color "#00FF00"

  # Update both name and color
  projectdb areas update "area-123" --name "Updated" --color "#0000FF"`,
	Args: cobra.ExactArgs(1),
	Run:  runAreasUpdate,
}

var areasDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an area",
	Long:  "Delete an area by ID. By default performs a soft delete.",
	Example: `  # Soft delete an area (can be recovered)
  projectdb areas delete "area-123"

  # Permanently delete an area
  projectdb areas delete "area-123" --permanent`,
	Args: cobra.ExactArgs(1),
	Run:  runAreasDelete,
}

var (
	areaName      string
	areaColor     string
	areaJSON      bool
	areaFormat    string
	areaFilter    string
	areaPermanent bool
)

func init() {
	areasCmd.AddCommand(areasCreateCmd)
	areasCmd.AddCommand(areasListCmd)
	areasCmd.AddCommand(areasGetCmd)
	areasCmd.AddCommand(areasUpdateCmd)
	areasCmd.AddCommand(areasDeleteCmd)

	areasCreateCmd.Flags().StringVar(&areaName, "name", "", "area name (required)")
	areasCreateCmd.Flags().StringVar(&areaColor, "color", "", "color in hex format (e.g., #FF5733)")
	areasCreateCmd.MarkFlagRequired("name")

	areasListCmd.Flags().BoolVar(&areaJSON, "json", false, "output as JSON")
	areasListCmd.Flags().StringVar(&areaFormat, "format", "table", "output format (table|json|yaml)")
	areasListCmd.Flags().StringVar(&areaFilter, "filter", "", "filter expression (e.g., 'name=Engineering')")

	areasUpdateCmd.Flags().StringVar(&areaName, "name", "", "new area name")
	areasUpdateCmd.Flags().StringVar(&areaColor, "color", "", "new color in hex format (e.g., #FF5733)")

	areasDeleteCmd.Flags().BoolVar(&areaPermanent, "permanent", false, "permanently delete (cannot be recovered)")
}

func runAreasCreate(cmd *cobra.Command, args []string) {
	if areaName == "" {
		cli.ExitWithError(cli.NewValidationError("name", "area name cannot be empty"))
	}

	color, err := cli.ParseColor(areaColor)
	if err != nil {
		cli.ExitWithError(err)
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	area, err := services.Areas.Create(ctx, areaName, color)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to create area"))
	}

	output.PrintSuccess(fmt.Sprintf("Area created with ID: %s", area.ID))
}

func runAreasList(cmd *cobra.Command, args []string) {
	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	areas, err := services.Areas.List(ctx)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to list areas"))
	}

	areaMaps := make([]map[string]interface{}, len(areas))
	for i, a := range areas {
		areaMaps[i] = domainAreaToMap(a)
	}

	if areaFilter != "" {
		areaMaps, err = filter.EvaluateFilter(areaFilter, areaMaps)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to apply filter"))
		}
	}

	useJSON := areaJSON || areaFormat == "json"
	useYAML := areaFormat == "yaml"

	if useJSON {
		formatter := output.NewJSONFormatter()
		for _, m := range areaMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output JSON"))
		}
		return
	}

	if useYAML {
		formatter := output.NewYAMLFormatter()
		for _, m := range areaMaps {
			formatter.AddObject(m)
		}
		if err := formatter.Flush(); err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to output YAML"))
		}
		return
	}

	if len(areaMaps) == 0 {
		output.PrintInfo("No areas found")
		return
	}

	formatter, _ := GetFormatter()
	formatter.PrintHeader([]string{"ID", "NAME", "COLOR", "CREATED"})
	for _, m := range areaMaps {
		colorStr := ""
		if v, ok := m["color"]; ok && v != nil {
			if s, ok := v.(string); ok {
				colorStr = s
			}
		}
		formatter.PrintRow([]string{
			m["id"].(string),
			m["name"].(string),
			colorStr,
			m["created_at"].(string)[:10],
		})
	}
	formatter.Flush()
}

func runAreasGet(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	area, err := services.Areas.GetByID(ctx, id)
	if err != nil {
		cli.ExitWithError(fmt.Errorf("area not found: %s", id))
	}

	formatter := output.NewJSONFormatter()
	if err := formatter.PrintObject(domainAreaToMap(*area)); err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to output area"))
	}
}

func runAreasUpdate(cmd *cobra.Command, args []string) {
	id := args[0]

	if areaName == "" && areaColor == "" {
		cli.ExitWithError(cli.NewValidationError("", "at least one of --name or --color must be provided"))
	}

	color, err := cli.ParseColor(areaColor)
	if err != nil {
		cli.ExitWithError(err)
	}

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	existing, err := services.Areas.GetByID(ctx, id)
	if err != nil {
		cli.ExitWithError(fmt.Errorf("area not found: %s", id))
	}

	newName := existing.Name
	if areaName != "" {
		newName = areaName
	}

	newColor := existing.Color
	if color != "" {
		newColor = color
	}

	area, err := services.Areas.Update(ctx, id, newName, newColor)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to update area"))
	}

	output.PrintSuccess(fmt.Sprintf("Area updated: %s", area.ID))
}

func runAreasDelete(cmd *cobra.Command, args []string) {
	id := args[0]

	services, err := GetServices()
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to connect to database"))
	}
	defer services.Close()

	ctx := context.Background()

	_, err = services.Areas.GetByID(ctx, id)
	if err != nil {
		cli.ExitWithError(fmt.Errorf("area not found: %s", id))
	}

	if areaPermanent {
		err := services.Areas.HardDelete(ctx, id)
		if err != nil {
			cli.ExitWithError(cli.WrapError(err, "failed to permanently delete area"))
		}
		output.PrintSuccess(fmt.Sprintf("Area permanently deleted: %s", id))
		return
	}

	err = services.Areas.SoftDelete(ctx, id)
	if err != nil {
		cli.ExitWithError(cli.WrapError(err, "failed to delete area"))
	}

	output.PrintSuccess(fmt.Sprintf("Area deleted: %s", id))
}

func domainAreaToMap(a domain.Area) map[string]interface{} {
	result := map[string]interface{}{
		"id":         a.ID,
		"name":       a.Name,
		"sort_order": a.SortOrder,
		"created_at": a.CreatedAt.Format(time.RFC3339),
		"updated_at": a.UpdatedAt.Format(time.RFC3339),
	}

	if a.Color != "" {
		result["color"] = string(a.Color)
	}

	return result
}
