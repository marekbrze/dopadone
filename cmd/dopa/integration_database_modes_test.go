package main

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

func TestGetDB_SQLiteMode(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	origDBMode := dbMode
	defer func() {
		dbPath = origDBPath
		tursoURL = origTursoURL
		tursoToken = origTursoToken
		dbMode = origDBMode
	}()

	dbPath = testDBPath
	tursoURL = ""
	tursoToken = ""
	dbMode = ""

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Errorf("db.Ping() error = %v", err)
	}
}

func TestGetDB_SQLiteExplicitMode(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	origDBMode := dbMode
	defer func() {
		dbPath = origDBPath
		tursoURL = origTursoURL
		tursoToken = origTursoToken
		dbMode = origDBMode
	}()

	dbPath = testDBPath
	tursoURL = ""
	tursoToken = ""
	dbMode = "local"

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := db.Ping(); err != nil {
		t.Errorf("db.Ping() error = %v", err)
	}
}

func TestGetServices_SQLiteMode(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	origDBMode := dbMode
	defer func() {
		dbPath = origDBPath
		tursoURL = origTursoURL
		tursoToken = origTursoToken
		dbMode = origDBMode
	}()

	dbPath = testDBPath
	tursoURL = ""
	tursoToken = ""
	dbMode = ""

	services, err := GetServices()
	if err != nil {
		t.Fatalf("GetServices() error = %v", err)
	}
	defer func() { _ = services.Close() }()

	if services.Projects == nil {
		t.Error("Projects service is nil")
	}
	if services.Tasks == nil {
		t.Error("Tasks service is nil")
	}
	if services.Areas == nil {
		t.Error("Areas service is nil")
	}
	if services.Subareas == nil {
		t.Error("Subareas service is nil")
	}
}

func TestGetDB_InvalidMode(t *testing.T) {
	origDBMode := dbMode
	defer func() { dbMode = origDBMode }()

	dbMode = "invalid-mode"

	_, err := GetDB()
	if err == nil {
		t.Error("GetDB() should return error for invalid mode")
	}
}

func TestGetDB_MissingCredentialsForRemote(t *testing.T) {
	origDBMode := dbMode
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	defer func() {
		dbMode = origDBMode
		tursoURL = origTursoURL
		tursoToken = origTursoToken
	}()

	dbMode = "remote"
	tursoURL = ""
	tursoToken = ""

	_, err := GetDB()
	if err == nil {
		t.Error("GetDB() should return error for remote mode without credentials")
	}
}

func TestGetDB_MissingCredentialsForReplica(t *testing.T) {
	origDBMode := dbMode
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	origDBPath := dbPath
	defer func() {
		dbMode = origDBMode
		tursoURL = origTursoURL
		tursoToken = origTursoToken
		dbPath = origDBPath
	}()

	dbMode = "replica"
	tursoURL = ""
	tursoToken = ""
	dbPath = "/tmp/test.db"

	_, err := GetDB()
	if err == nil {
		t.Error("GetDB() should return error for replica mode without credentials")
	}
}

func TestGetDB_AutoDetection_LocalSQLite(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origTursoURL := tursoURL
	origTursoToken := tursoToken
	origDBMode := dbMode
	defer func() {
		dbPath = origDBPath
		tursoURL = origTursoURL
		tursoToken = origTursoToken
		dbMode = origDBMode
	}()

	dbPath = testDBPath
	tursoURL = ""
	tursoToken = ""
	dbMode = ""

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	var mode string
	row := db.QueryRow("SELECT 1")
	if err := row.Scan(&mode); err != nil {
		t.Errorf("Failed to query database: %v", err)
	}
}

func TestGetDB_WithSyncInterval(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origSyncInterval := syncInterval
	defer func() {
		dbPath = origDBPath
		syncInterval = origSyncInterval
	}()

	dbPath = testDBPath
	syncInterval = "30s"

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()
}

func TestGetDB_InvalidSyncInterval(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	origSyncInterval := syncInterval
	defer func() {
		dbPath = origDBPath
		syncInterval = origSyncInterval
	}()

	dbPath = testDBPath
	syncInterval = "invalid"

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()
}

func TestServiceContainer_Close(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	defer func() { dbPath = origDBPath }()

	dbPath = testDBPath

	services, err := GetServices()
	if err != nil {
		t.Fatalf("GetServices() error = %v", err)
	}

	if err := services.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestGetFormatter_Table(t *testing.T) {
	origOutputFormat := outputFormat
	defer func() { outputFormat = origOutputFormat }()

	outputFormat = "table"

	formatter, err := GetFormatter()
	if err != nil {
		t.Fatalf("GetFormatter() error = %v", err)
	}
	if formatter == nil {
		t.Error("GetFormatter() returned nil")
	}
}

func TestGetFormatter_JSON(t *testing.T) {
	origOutputFormat := outputFormat
	defer func() { outputFormat = origOutputFormat }()

	outputFormat = "json"

	formatter, err := GetFormatter()
	if err != nil {
		t.Fatalf("GetFormatter() error = %v", err)
	}
	if formatter == nil {
		t.Error("GetFormatter() returned nil")
	}
}

func TestGetFormatter_Invalid(t *testing.T) {
	origOutputFormat := outputFormat
	defer func() { outputFormat = origOutputFormat }()

	outputFormat = "invalid"

	_, err := GetFormatter()
	if err == nil {
		t.Error("GetFormatter() should return error for invalid format")
	}
}

func TestDatabaseModes_SQLiteOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	defer func() { dbPath = origDBPath }()

	dbPath = testDBPath

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS test (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec("INSERT INTO test (name) VALUES (?)", "test-value")
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	var name string
	err = db.QueryRow("SELECT name FROM test WHERE id = 1").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if name != "test-value" {
		t.Errorf("name = %v, want test-value", name)
	}
}

func TestDriverDetection_Integration(t *testing.T) {
	tests := []struct {
		name         string
		dbPath       string
		tursoURL     string
		tursoToken   string
		dbMode       string
		expectedType driver.DriverType
	}{
		{
			name:         "auto_detect_local",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "",
			expectedType: driver.DriverSQLite,
		},
		{
			name:         "explicit_local",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "local",
			expectedType: driver.DriverSQLite,
		},
		{
			name:         "explicit_sqlite_alias",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "sqlite",
			expectedType: driver.DriverSQLite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadConfig(LoadConfigParams{
				DBPath:       tt.dbPath,
				TursoURL:     tt.tursoURL,
				TursoToken:   tt.tursoToken,
				DBMode:       tt.dbMode,
				SyncInterval: 60 * time.Second,
			})
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}
			driverCfg := cfg.ToDriverConfig()

			result, err := driver.DetectOrExplicitMode(driverCfg)
			if err != nil {
				t.Fatalf("DetectOrExplicitMode() error = %v", err)
			}

			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, want %v", result.Type, tt.expectedType)
			}
		})
	}
}

func TestGetDB_Concurrent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	defer func() { dbPath = origDBPath }()

	dbPath = testDBPath

	const numGoroutines = 10
	errCh := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			db, err := GetDB()
			if err != nil {
				errCh <- err
				return
			}
			defer func() { _ = db.Close() }()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := db.PingContext(ctx); err != nil {
				errCh <- err
				return
			}
			errCh <- nil
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		if err := <-errCh; err != nil {
			t.Errorf("Concurrent GetDB() error: %v", err)
		}
	}
}

func TestConnectionPoolSettings(t *testing.T) {
	tmpDir := t.TempDir()
	testDBPath := filepath.Join(tmpDir, "test.db")

	origDBPath := dbPath
	defer func() { dbPath = origDBPath }()

	dbPath = testDBPath

	db, err := GetDB()
	if err != nil {
		t.Fatalf("GetDB() error = %v", err)
	}
	defer func() { _ = db.Close() }()

	stats := db.Stats()
	if stats.MaxOpenConnections < 0 {
		t.Error("MaxOpenConnections should not be negative")
	}
}

func runDatabaseModeTest(t *testing.T, name string, testFunc func(t *testing.T, db *sql.DB)) {
	t.Helper()

	t.Run(name, func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping integration test in short mode")
		}

		tmpDir := t.TempDir()
		testDBPath := filepath.Join(tmpDir, "test.db")

		origDBPath := dbPath
		defer func() { dbPath = origDBPath }()

		dbPath = testDBPath

		db, err := GetDB()
		if err != nil {
			t.Fatalf("GetDB() error = %v", err)
		}
		defer func() { _ = db.Close() }()

		testFunc(t, db)
	})
}

func TestDatabaseMode_AllModes(t *testing.T) {
	runDatabaseModeTest(t, "basic_query", func(t *testing.T, db *sql.DB) {
		var result int
		err := db.QueryRow("SELECT 1 + 1").Scan(&result)
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if result != 2 {
			t.Errorf("Query result = %d, want 2", result)
		}
	})

	runDatabaseModeTest(t, "transaction", func(t *testing.T, db *sql.DB) {
		tx, err := db.Begin()
		if err != nil {
			t.Fatalf("Begin() error = %v", err)
		}
		defer func() { _ = tx.Rollback() }()

		_, err = tx.Exec("CREATE TABLE IF NOT EXISTS tx_test (id INTEGER PRIMARY KEY)")
		if err != nil {
			t.Fatalf("Create table failed: %v", err)
		}

		if err := tx.Commit(); err != nil {
			t.Fatalf("Commit() error = %v", err)
		}
	})
}

func TestEnvironmentVariables_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, "test.db")

	if err := os.Setenv("DOPA_DB_PATH", envPath); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("DOPA_DB_PATH") }()

	origDBPath := dbPath
	defer func() { dbPath = origDBPath }()

	dbPath = "./dopadone.db"

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       dbPath,
		SyncInterval: 60 * time.Second,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != envPath {
		t.Errorf("DatabasePath = %v, want %v (from env)", cfg.DatabasePath, envPath)
	}
}
