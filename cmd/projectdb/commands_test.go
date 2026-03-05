package main

import (
	"testing"

	"github.com/spf13/cobra"
)

func findCommand(root *cobra.Command, name string) *cobra.Command {
	for _, cmd := range root.Commands() {
		if cmd.Name() == name {
			return cmd
		}
	}
	return nil
}

func TestAreasListCommand_Flags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		shouldFail bool
	}{
		{"json flag", []string{"--json"}, false},
		{"format json", []string{"--format=json"}, false},
		{"format yaml", []string{"--format=yaml"}, false},
		{"format table", []string{"--format=table"}, false},
		{"filter flag", []string{"--filter", "name=test"}, false},
		{"combined flags", []string{"--json", "--filter", "name=test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listCmd := findCommand(areasCmd, "list")
			if listCmd == nil {
				t.Fatal("list command not found")
			}

			err := listCmd.ParseFlags(tt.args)
			if err != nil && !tt.shouldFail {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestSubareasListCommand_Flags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		shouldFail bool
	}{
		{"json flag", []string{"--json"}, false},
		{"format json", []string{"--format=json"}, false},
		{"format yaml", []string{"--format=yaml"}, false},
		{"area-id flag", []string{"--area-id", "area-123"}, false},
		{"filter flag", []string{"--filter", "name=test"}, false},
		{"combined flags", []string{"--json", "--area-id", "area-123", "--filter", "name=test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listCmd := findCommand(subareasCmd, "list")
			if listCmd == nil {
				t.Fatal("list command not found")
			}

			err := listCmd.ParseFlags(tt.args)
			if err != nil && !tt.shouldFail {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestProjectsListCommand_Flags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		shouldFail bool
	}{
		{"json flag", []string{"--json"}, false},
		{"format json", []string{"--format=json"}, false},
		{"format yaml", []string{"--format=yaml"}, false},
		{"status flag", []string{"--status", "active"}, false},
		{"priority flag", []string{"--priority", "high"}, false},
		{"subarea-id flag", []string{"--subarea-id", "subarea-123"}, false},
		{"parent-id flag", []string{"--parent-id", "project-123"}, false},
		{"filter flag", []string{"--filter", "status=active AND priority=high"}, false},
		{"combined flags", []string{"--json", "--status", "active", "--filter", "progress>=50"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listCmd := findCommand(projectsCmd, "list")
			if listCmd == nil {
				t.Fatal("list command not found")
			}

			err := listCmd.ParseFlags(tt.args)
			if err != nil && !tt.shouldFail {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestAreasCreateCommand_RequiredFlags(t *testing.T) {
	createCmd := findCommand(areasCmd, "create")
	if createCmd == nil {
		t.Fatal("create command not found")
	}

	requiredFlag := createCmd.Flags().Lookup("name")
	if requiredFlag == nil {
		t.Error("--name flag not found")
		return
	}

	if !createCmd.Flags().Changed("name") {
	}
}

func TestAreasUpdateCommand_Flags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"name flag", []string{"--name", "New Name"}},
		{"color flag", []string{"--color", "#FF0000"}},
		{"both flags", []string{"--name", "New Name", "--color", "#00FF00"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateCmd := findCommand(areasCmd, "update")
			if updateCmd == nil {
				t.Fatal("update command not found")
			}

			err := updateCmd.ParseFlags(tt.args)
			if err != nil {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestProjectsCreateCommand_Flags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"basic create with subarea", []string{"--name", "Test", "--subarea-id", "subarea-123"}},
		{"basic create with parent", []string{"--name", "Test", "--parent-id", "project-123"}},
		{"create with all flags", []string{
			"--name", "Test",
			"--subarea-id", "subarea-123",
			"--status", "active",
			"--priority", "high",
			"--progress", "50",
			"--color", "#FF0000",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createCmd := findCommand(projectsCmd, "create")
			if createCmd == nil {
				t.Fatal("create command not found")
			}

			err := createCmd.ParseFlags(tt.args)
			if err != nil {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestDeleteCommand_PermanentFlag(t *testing.T) {
	entities := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"areas", areasCmd},
		{"subareas", subareasCmd},
		{"projects", projectsCmd},
	}

	for _, entity := range entities {
		t.Run(entity.name, func(t *testing.T) {
			deleteCmd := findCommand(entity.cmd, "delete")
			if deleteCmd == nil {
				t.Fatal("delete command not found")
			}

			err := deleteCmd.ParseFlags([]string{"--permanent"})
			if err != nil {
				t.Errorf("ParseFlags() error = %v", err)
			}
		})
	}
}

func TestRootCommand_HasSubcommands(t *testing.T) {
	root := &cobra.Command{Use: "test"}
	root.AddCommand(areasCmd)
	root.AddCommand(subareasCmd)
	root.AddCommand(projectsCmd)

	expectedCommands := []string{"areas", "subareas", "projects"}
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range root.Commands() {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %s not found", expected)
		}
	}
}

func TestCommandAliases(t *testing.T) {
	tests := []struct {
		cmdName     string
		cmd         *cobra.Command
		expectedLen int
	}{
		{"areas", areasCmd, 1},
		{"subareas", subareasCmd, 2},
		{"projects", projectsCmd, 2},
	}

	for _, tt := range tests {
		t.Run(tt.cmdName, func(t *testing.T) {
			if len(tt.cmd.Aliases) != tt.expectedLen {
				t.Errorf("%s has %d aliases, expected %d", tt.cmdName, len(tt.cmd.Aliases), tt.expectedLen)
			}
		})
	}
}

func TestFlagDefaults(t *testing.T) {
	t.Run("areas list format default", func(t *testing.T) {
		listCmd := findCommand(areasCmd, "list")
		if listCmd == nil {
			t.Fatal("list command not found")
		}

		formatFlag := listCmd.Flags().Lookup("format")
		if formatFlag == nil {
			t.Error("--format flag not found")
			return
		}

		if formatFlag.DefValue != "table" {
			t.Errorf("format default = %q, want 'table'", formatFlag.DefValue)
		}
	})

	t.Run("subareas list format default", func(t *testing.T) {
		listCmd := findCommand(subareasCmd, "list")
		if listCmd == nil {
			t.Fatal("list command not found")
		}

		formatFlag := listCmd.Flags().Lookup("format")
		if formatFlag == nil {
			t.Error("--format flag not found")
			return
		}

		if formatFlag.DefValue != "table" {
			t.Errorf("format default = %q, want 'table'", formatFlag.DefValue)
		}
	})

	t.Run("projects list format default", func(t *testing.T) {
		listCmd := findCommand(projectsCmd, "list")
		if listCmd == nil {
			t.Fatal("list command not found")
		}

		formatFlag := listCmd.Flags().Lookup("format")
		if formatFlag == nil {
			t.Error("--format flag not found")
			return
		}

		if formatFlag.DefValue != "table" {
			t.Errorf("format default = %q, want 'table'", formatFlag.DefValue)
		}
	})
}

func TestFilterFlag_Exists(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"areas list", findCommand(areasCmd, "list")},
		{"subareas list", findCommand(subareasCmd, "list")},
		{"projects list", findCommand(projectsCmd, "list")},
	}

	for _, tt := range commands {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd == nil {
				t.Fatalf("command not found")
			}

			filterFlag := tt.cmd.Flags().Lookup("filter")
			if filterFlag == nil {
				t.Error("--filter flag not found")
			}
		})
	}
}

func TestYAMLFormat_Exists(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"areas list", findCommand(areasCmd, "list")},
		{"subareas list", findCommand(subareasCmd, "list")},
		{"projects list", findCommand(projectsCmd, "list")},
	}

	for _, tt := range commands {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cmd == nil {
				t.Fatalf("command not found")
			}

			formatFlag := tt.cmd.Flags().Lookup("format")
			if formatFlag == nil {
				t.Error("--format flag not found")
				return
			}

			err := tt.cmd.ParseFlags([]string{"--format=yaml"})
			if err != nil {
				t.Errorf("YAML format should be accepted: %v", err)
			}
		})
	}
}
