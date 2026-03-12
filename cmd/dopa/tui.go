package main

import (
	"fmt"

	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/service"
	"github.com/marekbrze/dopadone/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI interface",
	Long:  "Launch an interactive terminal user interface for managing projects, areas, and subareas.",
	Run: func(cmd *cobra.Command, args []string) {
		drv, err := GetDriver()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to get database driver: %v\n", err)
			return
		}
		defer func() {
			if err := cli.CloseDriver(drv); err != nil {
				_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Failed to close driver: %v\n", err)
			}
		}()

		dbConn := drv.GetDB()
		queries := db.New(dbConn)
		tm := db.NewTransactionManager(dbConn)

		projectService := service.NewProjectService(queries, tm)
		taskService := service.NewTaskService(queries, tm, projectService)

		services := &ServiceContainer{
			Projects: projectService,
			Tasks:    taskService,
			Subareas: service.NewSubareaService(queries, tm),
			Areas:    service.NewAreaService(queries, tm),
			db:       dbConn,
		}
		defer cli.CloseWithLog(services, "services")

		p := tui.New(services.Areas, services.Subareas, services.Projects, services.Tasks, drv)
		if _, err := p.Run(); err != nil {
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Error running TUI: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
