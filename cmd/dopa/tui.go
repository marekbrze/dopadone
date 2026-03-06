package main

import (
	"fmt"

	"github.com/example/dopadone/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI interface",
	Long:  "Launch an interactive terminal user interface for managing projects, areas, and subareas.",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := GetServices()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to initialize services: %v\n", err)
			return
		}
		defer services.Close()

		p := tui.New(services.Areas, services.Subareas, services.Projects, services.Tasks)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error running TUI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
