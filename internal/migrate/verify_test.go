package migrate

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/marekbrze/dopadone/internal/db/driver"
	_ "modernc.org/sqlite"
)

func TestVerifySchema(t *testing.T) {
	tests := []struct {
		name           string
		setupDB        func(db *sql.DB) error
		wantConsistent bool
		wantErr        bool
	}{
		{
			name: "empty_database",
			setupDB: func(db *sql.DB) error {
				return nil
			},
			wantConsistent: true,
			wantErr:        false,
		},
		{
			name: "with_migrations",
			setupDB: func(db *sql.DB) error {
				return Run(db, "up")
			},
			wantConsistent: true,
			wantErr:        false,
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

			mock := &mockVerifyDriver{db: db}
			verification, err := VerifySchema(mock)

			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if verification.Consistent != tt.wantConsistent {
				t.Errorf("VerifySchema().Consistent = %v, want %v", verification.Consistent, tt.wantConsistent)
			}
		})
	}
}

func TestVerifyConsistency(t *testing.T) {
	tests := []struct {
		name           string
		setupDB        func(db *sql.DB) error
		wantConsistent bool
	}{
		{
			name: "missing_tables",
			setupDB: func(db *sql.DB) error {
				return nil
			},
			wantConsistent: false,
		},
		{
			name: "all_tables_present",
			setupDB: func(db *sql.DB) error {
				return Run(db, "up")
			},
			wantConsistent: true,
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

			mock := &mockVerifyDriver{db: db}
			verification, err := VerifyConsistency(mock)

			if err != nil {
				t.Errorf("VerifyConsistency() error = %v", err)
				return
			}

			if verification.Consistent != tt.wantConsistent {
				t.Errorf("VerifyConsistency().Consistent = %v, want %v", verification.Consistent, tt.wantConsistent)
			}
		})
	}
}

func TestSchemaVerification_String(t *testing.T) {
	tests := []struct {
		name           string
		verification   *SchemaVerification
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "consistent_schema",
			verification: &SchemaVerification{
				LocalVersion: 20260304120000,
				Consistent:   true,
				Tables: []TableInfo{
					{Name: "areas", Columns: 6, Indexes: 0},
					{Name: "projects", Columns: 16, Indexes: 4},
				},
			},
			wantContains: []string{"Local Version:", "Consistent", "areas", "projects"},
		},
		{
			name: "schema_with_errors",
			verification: &SchemaVerification{
				LocalVersion: 0,
				Consistent:   false,
				Errors:       []string{"Missing table: areas"},
			},
			wantContains: []string{"Issues detected", "Missing table"},
		},
		{
			name: "schema_with_warnings",
			verification: &SchemaVerification{
				LocalVersion: 20260304120000,
				Consistent:   true,
				Warnings:     []string{"Could not determine remote version"},
			},
			wantContains: []string{"Warnings", "Could not determine"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.verification.String()

			for _, want := range tt.wantContains {
				if !contains(result, want) {
					t.Errorf("String() missing expected substring %q", want)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if contains(result, notWant) {
					t.Errorf("String() should not contain %q", notWant)
				}
			}
		})
	}
}

func TestGetGooseVersion(t *testing.T) {
	tests := []struct {
		name    string
		setupDB func(db *sql.DB) error
		want    int64
		wantErr bool
	}{
		{
			name: "no_goose_table",
			setupDB: func(db *sql.DB) error {
				return nil
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "with_migrations",
			setupDB: func(db *sql.DB) error {
				return Run(db, "up")
			},
			want:    20260304120000,
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

			got, err := getGooseVersion(db)

			if (err != nil) != tt.wantErr {
				t.Errorf("getGooseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("getGooseVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMissingTables(t *testing.T) {
	tests := []struct {
		name     string
		setupDB  func(db *sql.DB) error
		expected []string
		want     []string
	}{
		{
			name: "all_tables_present",
			setupDB: func(db *sql.DB) error {
				return Run(db, "up")
			},
			expected: expectedTables,
			want:     nil,
		},
		{
			name: "missing_some_tables",
			setupDB: func(db *sql.DB) error {
				_, err := db.Exec("CREATE TABLE areas (id TEXT PRIMARY KEY)")
				return err
			},
			expected: []string{"areas", "subareas", "projects"},
			want:     []string{"subareas", "projects"},
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

			got := findMissingTables(db, tt.expected)

			if len(got) != len(tt.want) {
				t.Errorf("findMissingTables() = %v, want %v", got, tt.want)
				return
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("findMissingTables()[%d] = %v, want %v", i, v, tt.want[i])
				}
			}
		})
	}
}

type mockVerifyDriver struct {
	db *sql.DB
}

func (m *mockVerifyDriver) Connect(ctx context.Context) error { return nil }
func (m *mockVerifyDriver) Close() error                      { return nil }
func (m *mockVerifyDriver) GetDB() *sql.DB                    { return m.db }
func (m *mockVerifyDriver) Ping(ctx context.Context) error    { return nil }
func (m *mockVerifyDriver) Type() driver.DriverType           { return driver.DriverSQLite }
func (m *mockVerifyDriver) Status() driver.ConnectionStatus   { return driver.StatusConnected }

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && s[:len(substr)] == substr ||
			(len(s) > len(substr) && contains(s[1:], substr)))
}
