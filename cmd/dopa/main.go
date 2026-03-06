package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/example/dopadone/internal/cli"
	"github.com/example/dopadone/internal/cli/output"
	"github.com/example/dopadone/internal/db"
	"github.com/example/dopadone/internal/migrate"
	"github.com/example/dopadone/internal/service"
	"github.com/example/dopadone/internal/version"
	"github.com/spf13/cobra"
)

var (
	dbPath       string
	outputFormat string
	showAll      bool
	skipMigrate  bool
)

var rootCmd = &cobra.Command{
	Use:     "dopa",
	Short:   "Dopadone - ADHD-friendly task and project management",
	Long:    "A CLI tool for managing projects, areas, and subareas in a SQLite database.",
	Version: version.Version,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		if showAll {
			fmt.Println(version.BuildInfo())
		} else {
			fmt.Printf("dopa %s\n", version.Version)
		}
	},
}

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version and run migrations",
	Long:  "Downloads the latest release, replaces the binary, and runs database migrations.",
	Run: func(cmd *cobra.Command, args []string) {
		opts := version.UpgradeOptions{
			DBPath:      dbPath,
			SkipMigrate: skipMigrate,
		}
		if err := version.PerformUpgrade(opts); err != nil {
			fmt.Fprintf(os.Stderr, "Upgrade failed: %v\n", err)
			os.Exit(1)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Long:  "Run embedded database migrations. Commands: up, down, status, reset",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if err := migrate.Run(db, "up"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations applied successfully.")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback the last migration",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if err := migrate.Run(db, "down"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migration rolled back successfully.")
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if err := migrate.Run(db, "status"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	},
}

var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset database (rollback all, then apply all)",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := GetDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer db.Close()

		if err := migrate.Run(db, "reset"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database reset successfully.")
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "./dopadone.db", "path to database file")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "table", "output format (table|json)")
	versionCmd.Flags().BoolVar(&showAll, "all", false, "show detailed build information")
	upgradeCmd.Flags().BoolVar(&skipMigrate, "skip-migrate", false, "skip running migrations after upgrade")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateStatusCmd)
	migrateCmd.AddCommand(migrateResetCmd)
	rootCmd.AddCommand(areasCmd)
	rootCmd.AddCommand(subareasCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(tasksCmd)
}

func GetDB() (*sql.DB, error) {
	db, err := cli.Connect(dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetFormatter() (output.Formatter, error) {
	return output.NewFormatter(outputFormat)
}

type ServiceContainer struct {
	Projects *service.ProjectService
	Tasks    *service.TaskService
	Subareas *service.SubareaService
	Areas    *service.AreaService
	db       *sql.DB
}

func (s *ServiceContainer) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func GetServices() (*ServiceContainer, error) {
	dbConn, err := GetDB()
	if err != nil {
		return nil, err
	}
	queries := db.New(dbConn)
	tm := db.NewTransactionManager(dbConn)

	return &ServiceContainer{
		Projects: service.NewProjectService(queries, tm),
		Tasks:    service.NewTaskService(queries, tm),
		Subareas: service.NewSubareaService(queries, tm),
		Areas:    service.NewAreaService(queries, tm),
		db:       dbConn,
	}, nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(cli.ExitError)
	}
}
