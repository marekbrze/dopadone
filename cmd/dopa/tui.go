package main

import (
	"fmt"

	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI interface",
	Long:  "Launch an interactive terminal user interface for managing projects, areas, and subareas.",
	Run: func(cmd *cobra.Command, args []string) {
		services, err := GetServices()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to initialize services: %v\n", err)
			return
		}
		defer cli.CloseWithLog(services, "services")

		p := tui.New(services.Areas, services.Subareas, services.Projects, services.Tasks)
		if _, err := p.Run(); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error running TUI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
