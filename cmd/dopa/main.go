package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/cli/output"
	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/migrate"
	"github.com/marekbrze/dopadone/internal/service"
	"github.com/marekbrze/dopadone/internal/version"
	"github.com/spf13/cobra"
)

var (
	dbPath       string
	outputFormat string
	showAll      bool
	skipMigrate  bool
	tursoURL     string
	tursoToken   string
	dbMode       string
	syncInterval string
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
		defer cli.CloseWithLog(db, "database")

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
		defer cli.CloseWithLog(db, "database")

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
		defer cli.CloseWithLog(db, "database")

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
		defer cli.CloseWithLog(db, "database")

		if err := migrate.Run(db, "reset"); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database reset successfully.")
	},
}

var migrateVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify database schema consistency",
	Long:  "Check that all expected tables exist and schema is consistent with migrations",
	Run: func(cmd *cobra.Command, args []string) {
		drv, err := GetDriver()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get database driver: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = drv.Close() }()

		verification, err := migrate.VerifyConsistency(drv)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Schema verification failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(verification.String())

		if !verification.Consistent {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbPath, "db", "./dopadone.db", "path to database file")
	rootCmd.PersistentFlags().StringVar(&tursoURL, "turso-url", "", "Turso database URL (env: TURSO_DATABASE_URL)")
	rootCmd.PersistentFlags().StringVar(&tursoToken, "turso-token", "", "Turso auth token (env: TURSO_AUTH_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&dbMode, "db-mode", "", "Database mode: local|remote|replica|auto (env: DOPA_DB_MODE, default: auto)")
	rootCmd.PersistentFlags().StringVar(&syncInterval, "sync-interval", "60s", "Sync interval for embedded replica mode")
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
	migrateCmd.AddCommand(migrateVerifyCmd)
	rootCmd.AddCommand(areasCmd)
	rootCmd.AddCommand(subareasCmd)
	rootCmd.AddCommand(projectsCmd)
	rootCmd.AddCommand(tasksCmd)
}

func GetDB() (*sql.DB, error) {
	drv, err := GetDriver()
	if err != nil {
		return nil, err
	}
	return drv.GetDB(), nil
}

func GetDriver() (driver.DatabaseDriver, error) {
	syncDur, err := time.ParseDuration(syncInterval)
	if err != nil {
		syncDur = 60 * time.Second
	}

	cfg := LoadConfig(dbPath, tursoURL, tursoToken, dbMode, syncDur)
	driverCfg := cfg.ToDriverConfig()

	result, err := driver.DetectOrExplicitMode(driverCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to detect database mode: %w", err)
	}

	log.Printf("[Database] Mode: %s (%s)", result.Type, result.Reason)

	if err := driver.ValidateConfigForMode(driverCfg, result.Type); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	switch result.Type {
	case driver.DriverSQLite:
		db, err := cli.Connect(cfg.DatabasePath)
		if err != nil {
			return nil, err
		}
		return &sqlDriverWrapper{db: db}, nil
	case driver.DriverTursoRemote, driver.DriverTursoReplica:
		drv, err := cli.ConnectWithDriver(
			driver.WithDriverType(result.Type),
			driver.WithDatabasePath(cfg.DatabasePath),
			driver.WithTurso(cfg.TursoURL, cfg.TursoToken),
			driver.WithSyncInterval(cfg.SyncInterval),
		)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := drv.Connect(ctx); err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		log.Printf("[Database] Connected successfully in %s mode", result.Type)
		return drv, nil
	default:
		return nil, fmt.Errorf("unsupported database mode: %s", result.Type)
	}
}

type sqlDriverWrapper struct {
	db *sql.DB
}

func (w *sqlDriverWrapper) Connect(ctx context.Context) error { return nil }
func (w *sqlDriverWrapper) Close() error                      { return w.db.Close() }
func (w *sqlDriverWrapper) GetDB() *sql.DB                    { return w.db }
func (w *sqlDriverWrapper) Ping(ctx context.Context) error    { return w.db.PingContext(ctx) }
func (w *sqlDriverWrapper) Type() driver.DriverType           { return driver.DriverSQLite }
func (w *sqlDriverWrapper) Status() driver.ConnectionStatus   { return driver.StatusConnected }

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

	projectService := service.NewProjectService(queries, tm)
	taskService := service.NewTaskService(queries, tm, projectService)

	return &ServiceContainer{
		Projects: projectService,
		Tasks:    taskService,
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
