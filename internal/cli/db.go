package cli

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/migrate"
	_ "modernc.org/sqlite"
)

var migrationMutex sync.Mutex

func Connect(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		return nil, NewValidationError("db", "database path cannot be empty")
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, WrapError(err, "failed to resolve database path")
	}

	if err := EnsureDirExists(absPath); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", absPath)
	if err != nil {
		return nil, WrapError(err, "failed to open database")
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, WrapError(err, "failed to connect to database")
	}

	return db, nil
}

func Close(db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func ConnectWithDriver(opts ...driver.Option) (driver.DatabaseDriver, error) {
	drv, err := driver.NewDriver(opts...)
	if err != nil {
		return nil, WrapError(err, "failed to create database driver")
	}
	return drv, nil
}

func CloseDriver(drv driver.DatabaseDriver) error {
	if drv == nil {
		return nil
	}
	return drv.Close()
}

func RunMigrations(db *sql.DB, migrationsPath string) error {
	if db == nil {
		return NewValidationError("db", "database connection is nil")
	}

	if migrationsPath == "" {
		return NewValidationError("migrations", "migrations path cannot be empty")
	}

	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		return NewValidationError("migrations", fmt.Sprintf("migrations directory does not exist: %s", migrationsPath))
	}

	return nil
}

func EnsureMigrations(db *sql.DB) error {
	if db == nil {
		return NewValidationError("db", "database connection is nil")
	}

	migrationMutex.Lock()
	defer migrationMutex.Unlock()

	if err := migrate.Run(db, "up"); err != nil {
		return WrapError(err, "failed to run auto-migrations")
	}

	return nil
}
