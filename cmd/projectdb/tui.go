package main

import (
	"fmt"

	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI interface",
	Long:  "Launch an interactive terminal user interface for managing projects, areas, and subareas.",
	Run: func(cmd *cobra.Command, args []string) {
		database, err := GetDB()
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Failed to connect to database: %v\n", err)
			return
		}
		defer database.Close()

		repo := db.New(database)
		p := tui.New(repo)
		if _, err := p.Run(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error running TUI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
