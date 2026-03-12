//go:build integration

package driver

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTursoReplicaDriver_Integration_Connect(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: dbPath,
		TursoURL:     url,
		TursoToken:   token,
		SyncInterval: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer driver.Close()

	if driver.Status() != StatusConnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusConnected)
	}

	if driver.GetDB() == nil {
		t.Error("GetDB() should not return nil after connect")
	}
}

func TestTursoReplicaDriver_Integration_Sync(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: dbPath,
		TursoURL:     url,
		TursoToken:   token,
		SyncInterval: 0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer driver.Close()

	info := driver.SyncInfo()
	if info.Status == SyncStatusError {
		t.Errorf("Initial sync status = %v, should not be error", info.Status)
	}

	if err := driver.Sync(); err != nil {
		t.Errorf("Sync() error = %v", err)
	}

	info = driver.SyncInfo()
	if info.Status != SyncStatusIdle {
		t.Errorf("SyncInfo.Status = %v, want %v", info.Status, SyncStatusIdle)
	}

	if info.LastSyncAt.IsZero() {
		t.Error("LastSyncAt should not be zero after sync")
	}
}

func TestTursoReplicaDriver_Integration_AutoSync(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: dbPath,
		TursoURL:     url,
		TursoToken:   token,
		SyncInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer driver.Close()

	initialSyncTime := driver.LastSyncTime()

	time.Sleep(2 * time.Second)

	newSyncTime := driver.LastSyncTime()

	if newSyncTime.Before(initialSyncTime) {
		t.Errorf("Auto-sync should have occurred. Initial: %v, New: %v", initialSyncTime, newSyncTime)
	}
}

func TestTursoReplicaDriver_Integration_ReadWrite(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: dbPath,
		TursoURL:     url,
		TursoToken:   token,
		SyncInterval: 0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer driver.Close()

	db := driver.GetDB()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS test_integration (id INTEGER PRIMARY KEY, name TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.ExecContext(ctx, "INSERT INTO test_integration (name) VALUES (?)", "test_value")
	if err != nil {
		t.Fatalf("Failed to insert: %v", err)
	}

	var name string
	err = db.QueryRowContext(ctx, "SELECT name FROM test_integration WHERE id = last_insert_rowid()").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if name != "test_value" {
		t.Errorf("name = %v, want %v", name, "test_value")
	}

	_, err = db.ExecContext(ctx, "DROP TABLE test_integration")
	if err != nil {
		t.Logf("Warning: failed to drop table: %v", err)
	}
}
