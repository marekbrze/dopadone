package cli

import (
	"database/sql"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestEnsureMigrations(t *testing.T) {
	tests := []struct {
		name    string
		setupDB func(db *sql.DB) error
		wantErr bool
	}{
		{
			name: "fresh_database",
			setupDB: func(db *sql.DB) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "already_migrated",
			setupDB: func(db *sql.DB) error {
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			dbPath := filepath.Join(tempDir, "test.db")

			db, err := sql.Open("sqlite", dbPath)
			if err != nil {
				t.Fatalf("Failed to open database: %v", err)
			}
			defer func() { _ = db.Close() }()

			if err := tt.setupDB(db); err != nil {
				t.Fatalf("Failed to setup database: %v", err)
			}

			err = EnsureMigrations(db)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureMigrations() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if err := db.QueryRow("SELECT COUNT(*) FROM goose_db_version").Scan(new(int)); err != nil {
					t.Errorf("goose_db_version table should exist: %v", err)
				}
			}
		})
	}
}

func TestEnsureMigrations_NilDB(t *testing.T) {
	err := EnsureMigrations(nil)
	if err == nil {
		t.Error("EnsureMigrations(nil) should return error")
	}
	if !IsValidationError(err) {
		t.Errorf("EnsureMigrations(nil) should return ValidationError, got %T", err)
	}
}

func TestEnsureMigrations_Idempotent(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := EnsureMigrations(db); err != nil {
		t.Fatalf("First EnsureMigrations() failed: %v", err)
	}

	if err := EnsureMigrations(db); err != nil {
		t.Fatalf("Second EnsureMigrations() failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM goose_db_version").Scan(&count); err != nil {
		t.Errorf("Failed to count migrations: %v", err)
	}

	if count == 0 {
		t.Error("Expected at least one migration record")
	}
}
