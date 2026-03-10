package cli

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Connect(dbPath string) (*sql.DB, error) {
	if dbPath == "" {
		return nil, NewValidationError("db", "database path cannot be empty")
	}

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, WrapError(err, "failed to resolve database path")
	}

	dir := filepath.Dir(absPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, NewValidationError("db", fmt.Sprintf("database directory does not exist: %s", dir))
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
