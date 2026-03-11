package driver

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	ResetRegistry()
	_ = RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
		return NewTursoReplicaDriver(config)
	})
}

func TestNewTursoReplicaDriver(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoURL:     "libsql://example.turso.io",
				TursoToken:   "test-token",
				SyncInterval: 60 * time.Second,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name: "missing database path",
			config: &Config{
				Type:       DriverTursoReplica,
				TursoURL:   "libsql://example.turso.io",
				TursoToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing url",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoToken:   "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoURL:     "libsql://example.turso.io",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, err := NewTursoReplicaDriver(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTursoReplicaDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && driver == nil {
				t.Error("NewTursoReplicaDriver() returned nil driver")
			}
		})
	}
}

func TestTursoReplicaDriver_Type(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.Type() != DriverTursoReplica {
		t.Errorf("Type() = %v, want %v", driver.Type(), DriverTursoReplica)
	}
}

func TestTursoReplicaDriver_Status(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.Status() != StatusDisconnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusDisconnected)
	}
}

func TestTursoReplicaDriver_GetDB_BeforeConnect(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.GetDB() != nil {
		t.Error("GetDB() should return nil before connect")
	}
}

func TestTursoReplicaDriver_Close_WithoutConnect(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if err := driver.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestTursoReplicaDriver_Ping_WithoutConnect(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	err = driver.Ping(context.Background())
	if err == nil {
		t.Error("Ping() should return error when not connected")
	}

	var driverErr *DriverError
	if !errors.As(err, &driverErr) {
		t.Error("Ping() should return DriverError")
	}
}

func TestTursoReplicaDriver_Sync_BeforeConnect(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	err = driver.Sync()
	if err == nil {
		t.Error("Sync() should return error when not connected")
	}

	var driverErr *DriverError
	if !errors.As(err, &driverErr) {
		t.Error("Sync() should return DriverError")
	}
}

func TestTursoReplicaDriver_SyncInfo(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	info := driver.SyncInfo()
	if info.Status != SyncStatusOffline {
		t.Errorf("SyncInfo.Status = %v, want %v", info.Status, SyncStatusOffline)
	}
}

func TestTursoReplicaDriver_LastSyncTime(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if !driver.LastSyncTime().IsZero() {
		t.Error("LastSyncTime() should be zero before any sync")
	}
}

func TestTursoReplicaDriver_Connect_ContextCancellation(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   "/tmp/test.db",
		TursoURL:       "libsql://example.turso.io",
		TursoToken:     "test-token",
		ConnectTimeout: 5 * time.Second,
		MaxRetries:     10,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = driver.Connect(ctx)
	if err == nil {
		t.Error("Connect() should return error with cancelled context")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("Connect() error = %v, want context.Canceled", err)
	}
}

func TestTursoReplicaDriver_Init(t *testing.T) {
	_ = RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
		return NewTursoReplicaDriver(config)
	})

	factory, ok := GetFactory(DriverTursoReplica)
	if !ok {
		t.Error("TursoReplicaDriver should be registered")
	}
	if factory == nil {
		t.Error("Factory should not be nil")
	}

	drv, err := factory(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Errorf("Factory returned error: %v", err)
	}
	if drv == nil {
		t.Error("Factory returned nil driver")
	}
	if drv.Type() != DriverTursoReplica {
		t.Errorf("Driver type = %v, want %v", drv.Type(), DriverTursoReplica)
	}
}

func TestTursoReplicaDriver_Connect_ZeroRetries(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = driver.Connect(ctx)
	if err == nil {
		_ = driver.Close()
		t.Skip("Connection succeeded unexpectedly, skipping test")
	}
}

func TestTursoReplicaDriver_DefaultsApplied(t *testing.T) {
	config := &Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	}

	driver, err := NewTursoReplicaDriver(config)
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.config.Type != DriverTursoReplica {
		t.Errorf("Type = %v, want %v", driver.config.Type, DriverTursoReplica)
	}
}

func TestTursoReplicaDriver_StatusTransitions(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.Status() != StatusDisconnected {
		t.Errorf("Initial Status() = %v, want %v", driver.Status(), StatusDisconnected)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = driver.Connect(ctx)

	if driver.Status() != StatusConnected && driver.Status() != StatusError {
		t.Errorf("After Connect() Status() = %v, want %v or %v", driver.Status(), StatusConnected, StatusError)
	}
}

func TestTursoReplicaDriver_Close_AfterFailedConnect(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = driver.Connect(ctx)

	if err := driver.Close(); err != nil {
		t.Errorf("Close() after failed connect error = %v", err)
	}

	if driver.Status() != StatusDisconnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusDisconnected)
	}
}

func TestTursoReplicaDriver_AutoSync_StartStop(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
		SyncInterval: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	driver.startAutoSync()

	if driver.syncCancel == nil {
		t.Error("syncCancel should not be nil after startAutoSync")
	}

	err = driver.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if driver.syncCancel != nil {
		t.Error("syncCancel should be nil after Close()")
	}
}

func TestTursoReplicaDriver_AutoSync_ContextCancellation(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
		SyncInterval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	driver.startAutoSync()

	done := make(chan bool, 1)
	go func() {
		_ = driver.Close()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Error("Close() should complete within 2 seconds")
	}
}

func TestTursoReplicaDriver_ConcurrentSyncAccess(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	done := make(chan bool, 1)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = driver.SyncInfo()
				_ = driver.LastSyncTime()
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestTursoReplicaDriver_ConcurrentStatusAccess(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	done := make(chan bool, 1)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = driver.Status()
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestTursoReplicaDriver_ConcurrentGetDBAccess(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	done := make(chan bool, 1)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = driver.GetDB()
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestTursoReplicaDriver_Close_StopsAutoSync(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
		SyncInterval: 10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	driver.startAutoSync()

	start := time.Now()
	err = driver.Close()
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	if elapsed > 500*time.Millisecond {
		t.Errorf("Close() took too long: %v, expected quick exit", elapsed)
	}
}

func TestTursoReplicaDriver_NoAutoSync_WhenZeroInterval(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
		SyncInterval: 0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	if driver.syncCancel != nil {
		t.Error("syncCancel should be nil when SyncInterval is 0")
	}
}

func TestTursoReplicaDriver_ConfigOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Option
		config *Config
	}{
		{
			name: "with sync interval",
			opts: []Option{
				WithSyncInterval(30 * time.Second),
			},
			config: &Config{SyncInterval: 30 * time.Second},
		},
		{
			name: "with connect timeout",
			opts: []Option{
				WithConnectTimeout(30 * time.Second),
			},
			config: &Config{ConnectTimeout: 30 * time.Second},
		},
		{
			name: "with max retries",
			opts: []Option{
				WithMaxRetries(5),
			},
			config: &Config{MaxRetries: 5},
		},
		{
			name: "with retry interval",
			opts: []Option{
				WithRetryInterval(2 * time.Second),
			},
			config: &Config{RetryInterval: 2 * time.Second},
		},
		{
			name: "all timeout options",
			opts: []Option{
				WithConnectTimeout(15 * time.Second),
				WithMaxRetries(3),
				WithRetryInterval(500 * time.Millisecond),
				WithSyncInterval(120 * time.Second),
			},
			config: &Config{
				ConnectTimeout: 15 * time.Second,
				MaxRetries:     3,
				RetryInterval:  500 * time.Millisecond,
				SyncInterval:   120 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := DefaultConfig()
			for _, opt := range tt.opts {
				opt(config)
			}

			if tt.config.ConnectTimeout != 0 && config.ConnectTimeout != tt.config.ConnectTimeout {
				t.Errorf("ConnectTimeout = %v, want %v", config.ConnectTimeout, tt.config.ConnectTimeout)
			}
			if tt.config.MaxRetries != 0 && config.MaxRetries != tt.config.MaxRetries {
				t.Errorf("MaxRetries = %v, want %v", config.MaxRetries, tt.config.MaxRetries)
			}
			if tt.config.RetryInterval != 0 && config.RetryInterval != tt.config.RetryInterval {
				t.Errorf("RetryInterval = %v, want %v", config.RetryInterval, tt.config.RetryInterval)
			}
			if tt.config.SyncInterval != 0 && config.SyncInterval != tt.config.SyncInterval {
				t.Errorf("SyncInterval = %v, want %v", config.SyncInterval, tt.config.SyncInterval)
			}
		})
	}
}

func TestTursoReplicaDriver_Connect_RetryAttempts(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 10 * time.Millisecond,
		MaxRetries:     2,
		RetryInterval:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	err = driver.Connect(ctx)
	elapsed := time.Since(start)

	if err == nil {
		_ = driver.Close()
		t.Skip("Connection succeeded unexpectedly, skipping test")
	}

	minExpected := 3*10*time.Millisecond + 2*5*time.Millisecond
	if elapsed < minExpected-5*time.Millisecond {
		t.Errorf("Connect took %v, expected at least %v for 3 attempts with retry interval", elapsed, minExpected)
	}
}

func TestTursoReplicaDriver_Connect_CancelDuringRetry(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 100 * time.Millisecond,
		MaxRetries:     10,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err = driver.Connect(ctx)
	if err == nil {
		_ = driver.Close()
		t.Skip("Connection succeeded unexpectedly, skipping test")
	}

	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Errorf("Connect() error = %v, want context.DeadlineExceeded or context.Canceled", err)
	}
}

func TestTursoReplicaDriver_SyncStatusTransitions(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	info := driver.SyncInfo()
	if info.Status != SyncStatusOffline {
		t.Errorf("Initial SyncInfo.Status = %v, want %v", info.Status, SyncStatusOffline)
	}

	driver.mu.Lock()
	driver.syncStatus = SyncStatusIdle
	driver.mu.Unlock()

	info = driver.SyncInfo()
	if info.Status != SyncStatusIdle {
		t.Errorf("SyncInfo.Status = %v, want %v", info.Status, SyncStatusIdle)
	}

	driver.mu.Lock()
	driver.syncStatus = SyncStatusSyncing
	driver.mu.Unlock()

	if driver.Status() != StatusConnecting {
		t.Errorf("Status() during sync = %v, want %v", driver.Status(), StatusConnecting)
	}

	driver.mu.Lock()
	driver.syncStatus = SyncStatusError
	driver.lastSyncErr = errors.New("test error")
	driver.mu.Unlock()

	info = driver.SyncInfo()
	if info.Status != SyncStatusError {
		t.Errorf("SyncInfo.Status = %v, want %v", info.Status, SyncStatusError)
	}
	if info.LastError == nil {
		t.Error("SyncInfo.LastError should not be nil when status is error")
	}
}

func TestTursoReplicaDriver_DoubleClose(t *testing.T) {
	driver, err := NewTursoReplicaDriver(&Config{
		Type:         DriverTursoReplica,
		DatabasePath: "/tmp/test.db",
		TursoURL:     "libsql://example.turso.io",
		TursoToken:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	err = driver.Close()
	if err != nil {
		t.Errorf("First Close() error = %v", err)
	}

	err = driver.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}

func TestTursoReplicaDriver_WithTursoOption(t *testing.T) {
	config := DefaultConfig()
	WithTurso("libsql://test.turso.io", "my-token")(config)

	if config.TursoURL != "libsql://test.turso.io" {
		t.Errorf("TursoURL = %v, want %v", config.TursoURL, "libsql://test.turso.io")
	}
	if config.TursoToken != "my-token" {
		t.Errorf("TursoToken = %v, want %v", config.TursoToken, "my-token")
	}
}

func TestTursoReplicaDriver_WithDatabasePathOption(t *testing.T) {
	config := DefaultConfig()
	WithDatabasePath("/custom/path.db")(config)

	if config.DatabasePath != "/custom/path.db" {
		t.Errorf("DatabasePath = %v, want %v", config.DatabasePath, "/custom/path.db")
	}
}

func TestTursoReplicaDriver_WithDriverTypeOption(t *testing.T) {
	config := DefaultConfig()
	WithDriverType(DriverTursoReplica)(config)

	if config.Type != DriverTursoReplica {
		t.Errorf("Type = %v, want %v", config.Type, DriverTursoReplica)
	}
}

func TestConfig_Validate_TursoReplica(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid turso-replica config",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoURL:     "libsql://example.turso.io",
				TursoToken:   "test-token",
			},
			wantErr: false,
		},
		{
			name: "missing database path",
			config: &Config{
				Type:       DriverTursoReplica,
				TursoURL:   "libsql://example.turso.io",
				TursoToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing url",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoToken:   "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoURL:     "libsql://example.turso.io",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTursoReplicaDriver_NegativeMaxRetries(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	driver, err := NewTursoReplicaDriver(&Config{
		Type:           DriverTursoReplica,
		DatabasePath:   dbPath,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     -1,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoReplicaDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = driver.Connect(ctx)

	if driver.syncCancel != nil {
		_ = driver.Close()
	}
}

func TestErrorsJoin(t *testing.T) {
	tests := []struct {
		name    string
		errs    []error
		wantMsg string
		wantNil bool
	}{
		{
			name:    "no errors",
			errs:    []error{},
			wantMsg: "",
		},
		{
			name:    "single error",
			errs:    []error{errors.New("error1")},
			wantMsg: "error1",
		},
		{
			name:    "multiple errors",
			errs:    []error{errors.New("error1"), errors.New("error2")},
			wantMsg: "error1; error2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorsJoin(tt.errs...)
			if tt.wantMsg == "" && err.Error() != "" {
				t.Errorf("errorsJoin() = %v, want empty", err.Error())
			}
			if tt.wantMsg != "" && err.Error() != tt.wantMsg {
				t.Errorf("errorsJoin() = %v, want %v", err.Error(), tt.wantMsg)
			}
		})
	}
}
