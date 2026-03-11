package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/pressly/goose/v3"
)

//go:embed migrations
var migrationsFS embed.FS

func Run(db *sql.DB, command string) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	fsys, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to sub migrations: %w", err)
	}
	goose.SetBaseFS(fsys)

	var runErr error
	switch command {
	case "up":
		runErr = goose.Up(db, ".")
	case "down":
		runErr = goose.Down(db, ".")
	case "status":
		runErr = goose.Status(db, ".")
	case "reset":
		runErr = goose.Reset(db, ".")
	default:
		return fmt.Errorf("unknown migration command: %s", command)
	}

	goose.SetBaseFS(nil)

	if runErr != nil {
		return fmt.Errorf("migration %s failed: %w", command, runErr)
	}

	return nil
}

func RunOnDriver(d driver.DatabaseDriver, command string) error {
	db := d.GetDB()
	if db == nil {
		return fmt.Errorf("driver not connected")
	}

	if err := Run(db, command); err != nil {
		return err
	}

	if d.Type() == driver.DriverTursoReplica {
		if syncer, ok := d.(interface{ Sync() error }); ok {
			log.Printf("[Migration] Syncing schema to remote for embedded replica")
			if err := syncer.Sync(); err != nil {
				return fmt.Errorf("failed to sync schema to remote: %w", err)
			}
			log.Printf("[Migration] Schema sync completed")
		}
	}

	return nil
}
