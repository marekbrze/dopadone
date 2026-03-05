package migrate

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

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

	switch command {
	case "up":
		if err := goose.Up(db, "."); err != nil {
			return fmt.Errorf("migration up failed: %w", err)
		}
	case "down":
		if err := goose.Down(db, "."); err != nil {
			return fmt.Errorf("migration down failed: %w", err)
		}
	case "status":
		if err := goose.Status(db, "."); err != nil {
			return fmt.Errorf("migration status failed: %w", err)
		}
	case "reset":
		if err := goose.Reset(db, "."); err != nil {
			return fmt.Errorf("migration reset failed: %w", err)
		}
	default:
		return fmt.Errorf("unknown migration command: %s", command)
	}

	return nil
}
