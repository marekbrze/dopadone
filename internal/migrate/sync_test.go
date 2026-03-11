package migrate

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/marekbrze/dopadone/internal/db/driver"
	_ "modernc.org/sqlite"
)

type mockDriver struct {
	db       *sql.DB
	synced   bool
	syncErr  error
	dataType driver.DriverType
}

func (m *mockDriver) Connect(ctx context.Context) error { return nil }
func (m *mockDriver) Close() error                      { return nil }
func (m *mockDriver) GetDB() *sql.DB                    { return m.db }
func (m *mockDriver) Ping(ctx context.Context) error    { return nil }
func (m *mockDriver) Type() driver.DriverType           { return m.dataType }
func (m *mockDriver) Status() driver.ConnectionStatus {
	return driver.StatusConnected
}

func (m *mockDriver) Sync() error {
	m.synced = true
	return m.syncErr
}

func TestMigrationSyncer_RunAndSync(t *testing.T) {
	tests := []struct {
		name       string
		driverType driver.DriverType
		wantSync   bool
	}{
		{"sqlite_no_sync", driver.DriverSQLite, false},
		{"turso_remote_no_sync", driver.DriverTursoRemote, false},
		{"turso_replica_syncs", driver.DriverTursoReplica, true},
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

			mock := &mockDriver{
				db:       db,
				dataType: tt.driverType,
			}

			syncer := NewMigrationSyncer(mock)

			if err := syncer.RunAndSync(db, "up"); err != nil {
				t.Fatalf("RunAndSync failed: %v", err)
			}

			if mock.synced != tt.wantSync {
				t.Errorf("synced = %v, want %v", mock.synced, tt.wantSync)
			}
		})
	}
}

func TestMigrationSyncer_RunAndSyncWithError(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	mock := &mockDriver{
		db:       db,
		dataType: driver.DriverTursoReplica,
		syncErr:  os.ErrNotExist,
	}

	syncer := NewMigrationSyncer(mock)

	if err := syncer.RunAndSync(db, "up"); err == nil {
		t.Error("Expected error when sync fails, got nil")
	}
}
