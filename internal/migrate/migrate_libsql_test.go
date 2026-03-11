package migrate

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/tursodatabase/go-libsql"
)

func TestMigrationsWithLibSQL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping libSQL integration test in short mode")
	}

	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"embedded_replica_compatible", testMigrationsWithEmbeddedReplica},
		{"remote_connection_compatible", testMigrationsWithRemoteConnection},
		{"idempotent_migrations", testIdempotentMigrations},
		{"migration_status", testMigrationStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.fn)
	}
}

func testMigrationsWithEmbeddedReplica(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		t.Fatalf("Failed to open libSQL database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping libSQL database: %v", err)
	}

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Failed to run migrations up: %v", err)
	}

	verifyExpectedTables(t, db)
}

func testMigrationsWithRemoteConnection(t *testing.T) {
	tursoURL := os.Getenv("TURSO_TEST_URL")
	tursoToken := os.Getenv("TURSO_TEST_TOKEN")

	if tursoURL == "" || tursoToken == "" {
		t.Skip("Skipping remote test: TURSO_TEST_URL and TURSO_TEST_TOKEN required")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "remote.db")

	connStr := buildLibSQLConnectionString(dbPath, tursoURL, tursoToken)

	db, err := sql.Open("libsql", connStr)
	if err != nil {
		t.Fatalf("Failed to open libSQL remote connection: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping libSQL remote: %v", err)
	}

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Failed to run migrations on remote: %v", err)
	}

	verifyExpectedTables(t, db)
}

func testIdempotentMigrations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		t.Fatalf("Failed to open libSQL database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := Run(db, "up"); err != nil {
		t.Fatalf("First migration up failed: %v", err)
	}

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Second migration up (idempotent) failed: %v", err)
	}

	verifyExpectedTables(t, db)
}

func testMigrationStatus(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		t.Fatalf("Failed to open libSQL database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := Run(db, "status"); err != nil {
		t.Logf("Migration status on empty db: %v", err)
	}

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Migration up failed: %v", err)
	}

	if err := Run(db, "status"); err != nil {
		t.Fatalf("Migration status after up failed: %v", err)
	}
}

func TestMigrationDownAndReset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping libSQL integration test in short mode")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("libsql", "file:"+dbPath)
	if err != nil {
		t.Fatalf("Failed to open libSQL database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Migration up failed: %v", err)
	}

	if err := Run(db, "down"); err != nil {
		t.Logf("Migration down: %v (may fail if last migration)", err)
	}

	if err := Run(db, "up"); err != nil {
		t.Fatalf("Migration up after down failed: %v", err)
	}

	if err := Run(db, "reset"); err != nil {
		t.Fatalf("Migration reset failed: %v", err)
	}

	verifyExpectedTables(t, db)
}

func verifyExpectedTables(t *testing.T, db *sql.DB) {
	t.Helper()

	expectedTables := []string{"areas", "subareas", "projects", "tasks", "goose_db_version"}

	rows, err := db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
	`)
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	defer func() { _ = rows.Close() }()

	existingTables := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("Failed to scan table name: %v", err)
		}
		existingTables[name] = true
	}

	for _, expected := range expectedTables {
		if !existingTables[expected] {
			t.Errorf("Expected table %q not found", expected)
		}
	}
}

func buildLibSQLConnectionString(dbPath, url, token string) string {
	if token != "" {
		return dbPath + "?url=" + url + "&authToken=" + token
	}
	return dbPath
}
