//go:build integration

package driver

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestTursoRemoteDriver_Integration_Connect(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required for integration tests")
	}

	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       url,
		TursoToken:     token,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}
	defer driver.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if driver.Status() != StatusConnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusConnected)
	}
}

func TestTursoRemoteDriver_Integration_Query(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required for integration tests")
	}

	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       url,
		TursoToken:     token,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}
	defer driver.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	db := driver.GetDB()
	if db == nil {
		t.Fatal("GetDB() returned nil")
	}

	var result int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		t.Fatalf("QueryRowContext() error = %v", err)
	}

	if result != 1 {
		t.Errorf("Query result = %d, want 1", result)
	}
}

func TestTursoRemoteDriver_Integration_Transaction(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required for integration tests")
	}

	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       url,
		TursoToken:     token,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}
	defer driver.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	db := driver.GetDB()
	if db == nil {
		t.Fatal("GetDB() returned nil")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("BeginTx() error = %v", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "SELECT 1"); err != nil {
		t.Fatalf("ExecContext() error = %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Commit() error = %v", err)
	}
}

func TestTursoRemoteDriver_Integration_Ping(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required for integration tests")
	}

	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       url,
		TursoToken:     token,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}
	defer driver.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if err := driver.Ping(ctx); err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}

func TestTursoRemoteDriver_Integration_Close(t *testing.T) {
	url := os.Getenv("TURSO_TEST_URL")
	token := os.Getenv("TURSO_TEST_TOKEN")
	if url == "" || token == "" {
		t.Skip("TURSO_TEST_URL and TURSO_TEST_TOKEN required for integration tests")
	}

	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       url,
		TursoToken:     token,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := driver.Connect(ctx); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if err := driver.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if driver.Status() != StatusDisconnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusDisconnected)
	}

	if driver.GetDB() != nil {
		t.Error("GetDB() should return nil after close")
	}
}
