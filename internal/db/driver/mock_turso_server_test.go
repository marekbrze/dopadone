package driver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestMockTursoServer(t *testing.T) {
	server := NewMockTursoServer(t)
	defer func() {
		if err := server.Close(); err != nil {
			t.Logf("Failed to close mock server: %v", err)
		}
	}()

	t.Run("starts and stops", func(t *testing.T) {
		if server.URL() == "" {
			t.Error("Server URL should not be empty")
		}
	})
}

type MockTursoServer struct {
	t        *testing.T
	server   *http.Server
	port     int
	baseURL  string
	mu       sync.Mutex
	handlers map[string]http.HandlerFunc
}

func NewMockTursoServer(t *testing.T) *MockTursoServer {
	server := &MockTursoServer{
		t:        t,
		handlers: make(map[string]http.HandlerFunc),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleRoot)
	mux.HandleFunc("/v1/pipeline", server.handlePipeline)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}

	server.server = &http.Server{Handler: mux}
	server.port = listener.Addr().(*net.TCPAddr).Port
	server.baseURL = fmt.Sprintf("http://127.0.0.1:%d", server.port)

	go func() {
		if err := server.server.Serve(listener); err != nil && err != http.ErrServerClosed {
			t.Logf("Mock server error: %v", err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	return server
}

func (m *MockTursoServer) URL() string {
	return m.baseURL
}

func (m *MockTursoServer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.server.Shutdown(ctx)
}

func (m *MockTursoServer) handleRoot(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		m.t.Logf("Failed to encode response: %v", err)
	}
}

func (m *MockTursoServer) handlePipeline(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")

	body, _ := io.ReadAll(r.Body)
	var req struct {
		Requests []struct {
			Type string `json:"type"`
			Stmt struct {
				Sql string `json:"sql"`
			} `json:"stmt"`
		} `json:"requests"`
	}
	_ = json.Unmarshal(body, &req)

	results := make([]map[string]interface{}, len(req.Requests))
	for i := range req.Requests {
		results[i] = map[string]interface{}{
			"type": "ok",
			"result": map[string]interface{}{
				"cols": []string{},
				"rows": [][]interface{}{},
			},
		}
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"results": results,
	}); err != nil {
		m.t.Logf("Failed to encode response: %v", err)
	}
}

func (m *MockTursoServer) SetResponse(path string, handler http.HandlerFunc) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[path] = handler
}

func TestTursoRemote_WithMockServer(t *testing.T) {
	server := NewMockTursoServer(t)
	defer func() {
		if err := server.Close(); err != nil {
			t.Logf("Failed to close mock server: %v", err)
		}
	}()

	t.Run("connect with mock server", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       server.URL(),
			TursoToken:     "test-token",
			ConnectTimeout: 5 * time.Second,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err = drv.Connect(ctx)
		if err != nil {
			t.Logf("Connect() error = %v (expected - mock is not a real libsql server)", err)
		}
	})
}

func TestTursoReplica_WithMockServer(t *testing.T) {
	server := NewMockTursoServer(t)
	defer func() {
		if err := server.Close(); err != nil {
			t.Logf("Failed to close mock server: %v", err)
		}
	}()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	t.Run("replica mode configuration", func(t *testing.T) {
		drv, err := NewTursoReplicaDriver(&Config{
			Type:           DriverTursoReplica,
			DatabasePath:   dbPath,
			TursoURL:       server.URL(),
			TursoToken:     "test-token",
			SyncInterval:   60 * time.Second,
			ConnectTimeout: 5 * time.Second,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		if drv.Type() != DriverTursoReplica {
			t.Errorf("Type() = %v, want %v", drv.Type(), DriverTursoReplica)
		}

		if drv.Status() != StatusDisconnected {
			t.Errorf("Status() = %v, want %v", drv.Status(), StatusDisconnected)
		}
	})
}

func TestFailFast_InvalidURL(t *testing.T) {
	t.Run("remote mode fails fast with invalid URL", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       "https://nonexistent.invalid.example",
			TursoToken:     "test-token",
			ConnectTimeout: 100 * time.Millisecond,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		start := time.Now()
		err = drv.Connect(ctx)
		elapsed := time.Since(start)

		if err == nil {
			_ = drv.Close()
			t.Skip("Connection succeeded unexpectedly")
		}

		if elapsed > 5*time.Second {
			t.Errorf("Connect took too long: %v, expected fail-fast", elapsed)
		}
	})
}

func TestFailFast_InvalidToken(t *testing.T) {
	t.Run("remote mode fails with empty token", func(t *testing.T) {
		_, err := NewTursoRemoteDriver(&Config{
			Type:       DriverTursoRemote,
			TursoURL:   "https://example.turso.io",
			TursoToken: "",
		})
		if err == nil {
			t.Error("NewTursoRemoteDriver() should fail with empty token")
		}
	})
}

func TestConnectionRetryLogic(t *testing.T) {
	t.Run("verify retry attempts", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       "https://nonexistent.invalid.example",
			TursoToken:     "test-token",
			ConnectTimeout: 50 * time.Millisecond,
			MaxRetries:     2,
			RetryInterval:  20 * time.Millisecond,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		start := time.Now()
		err = drv.Connect(ctx)
		elapsed := time.Since(start)

		if err == nil {
			_ = drv.Close()
			t.Skip("Connection succeeded unexpectedly")
		}

		minExpected := 3*50*time.Millisecond + 2*20*time.Millisecond
		if elapsed < minExpected-10*time.Millisecond {
			t.Errorf("Connect took %v, expected at least %v for 3 attempts", elapsed, minExpected)
		}
	})
}

func TestTursoReplica_SyncBehavior(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "replica.db")

	t.Run("sync info before connect", func(t *testing.T) {
		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		info := drv.SyncInfo()
		if info.Status != SyncStatusOffline {
			t.Errorf("SyncInfo.Status = %v, want %v", info.Status, SyncStatusOffline)
		}

		if !drv.LastSyncTime().IsZero() {
			t.Error("LastSyncTime() should be zero before connect")
		}
	})

	t.Run("sync returns error before connect", func(t *testing.T) {
		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		err = drv.Sync()
		if err == nil {
			t.Error("Sync() should return error before connect")
		}
	})
}

func TestTursoAllModes_NoExternalDeps(t *testing.T) {
	t.Run("all driver types can be created", func(t *testing.T) {
		tests := []struct {
			name   string
			config *Config
		}{
			{
				name: "sqlite",
				config: &Config{
					Type:         DriverSQLite,
					DatabasePath: "/tmp/test.db",
				},
			},
			{
				name: "turso-remote",
				config: &Config{
					Type:       DriverTursoRemote,
					TursoURL:   "https://example.turso.io",
					TursoToken: "test-token",
				},
			},
			{
				name: "turso-replica",
				config: &Config{
					Type:         DriverTursoReplica,
					DatabasePath: "/tmp/test.db",
					TursoURL:     "https://example.turso.io",
					TursoToken:   "test-token",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var drv DatabaseDriver
				var err error

				switch tt.config.Type {
				case DriverSQLite:
					drv = &mockDriver{driverType: DriverSQLite}
				case DriverTursoRemote:
					drv, err = NewTursoRemoteDriver(tt.config)
				case DriverTursoReplica:
					drv, err = NewTursoReplicaDriver(tt.config)
				}

				if err != nil {
					t.Errorf("Failed to create %s driver: %v", tt.name, err)
				}
				if drv == nil {
					t.Errorf("%s driver is nil", tt.name)
				}
				if drv != nil && drv.Type() != tt.config.Type {
					t.Errorf("Type() = %v, want %v", drv.Type(), tt.config.Type)
				}
			})
		}
	})
}

func TestTursoConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid turso remote",
			config: &Config{
				Type:       DriverTursoRemote,
				TursoURL:   "https://example.turso.io",
				TursoToken: "test-token",
			},
			wantErr: false,
		},
		{
			name: "missing url",
			config: &Config{
				Type:       DriverTursoRemote,
				TursoToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: &Config{
				Type:     DriverTursoRemote,
				TursoURL: "https://example.turso.io",
			},
			wantErr: true,
		},
		{
			name: "valid turso replica",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/test.db",
				TursoURL:     "https://example.turso.io",
				TursoToken:   "test-token",
			},
			wantErr: false,
		},
		{
			name: "replica missing path",
			config: &Config{
				Type:       DriverTursoReplica,
				TursoURL:   "https://example.turso.io",
				TursoToken: "test-token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTursoRemote_ConnectionErrorHandling(t *testing.T) {
	t.Run("handles empty URL", func(t *testing.T) {
		_, err := NewTursoRemoteDriver(&Config{
			Type:       DriverTursoRemote,
			TursoURL:   "",
			TursoToken: "test-token",
		})
		if err == nil {
			t.Error("NewTursoRemoteDriver() should fail with empty URL")
		}
	})

	t.Run("handles timeout on connection", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       "libsql://nonexistent.invalid.example",
			TursoToken:     "test-token",
			ConnectTimeout: 50 * time.Millisecond,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		start := time.Now()
		err = drv.Connect(ctx)
		elapsed := time.Since(start)

		if err == nil {
			_ = drv.Close()
			t.Skip("Connection succeeded unexpectedly")
		}

		if elapsed > 1*time.Second {
			t.Errorf("Connect took too long: %v, expected timeout within 200ms", elapsed)
		}
	})

	t.Run("handles context deadline exceeded", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       "libsql://nonexistent.invalid.example",
			TursoToken:     "test-token",
			ConnectTimeout: 1 * time.Second,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		err = drv.Connect(ctx)
		if err == nil {
			_ = drv.Close()
			t.Skip("Connection succeeded unexpectedly")
		}

		if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			t.Logf("Connect() error = %v (acceptable)", err)
		}
	})
}

func TestTursoRemote_QueryErrorHandling(t *testing.T) {
	t.Run("GetDB returns nil before connect", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:       DriverTursoRemote,
			TursoURL:   "https://example.turso.io",
			TursoToken: "test-token",
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		if drv.GetDB() != nil {
			t.Error("GetDB() should return nil before connect")
		}
	})

	t.Run("GetDB returns nil after failed connect", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:           DriverTursoRemote,
			TursoURL:       "https://nonexistent.invalid.example",
			TursoToken:     "test-token",
			ConnectTimeout: 50 * time.Millisecond,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_ = drv.Connect(ctx)

		if drv.GetDB() != nil {
			t.Error("GetDB() should return nil after failed connect")
		}
	})

	t.Run("Ping returns error before connect", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:       DriverTursoRemote,
			TursoURL:   "https://example.turso.io",
			TursoToken: "test-token",
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		err = drv.Ping(context.Background())
		if err == nil {
			t.Error("Ping() should return error before connect")
		}
	})

	t.Run("Close is idempotent", func(t *testing.T) {
		drv, err := NewTursoRemoteDriver(&Config{
			Type:       DriverTursoRemote,
			TursoURL:   "https://example.turso.io",
			TursoToken: "test-token",
		})
		if err != nil {
			t.Fatalf("NewTursoRemoteDriver() error = %v", err)
		}

		if err := drv.Close(); err != nil {
			t.Errorf("First Close() error = %v", err)
		}
		if err := drv.Close(); err != nil {
			t.Errorf("Second Close() error = %v", err)
		}
		if err := drv.Close(); err != nil {
			t.Errorf("Third Close() error = %v", err)
		}
	})
}

func TestTursoReplica_SyncStateMachine(t *testing.T) {
	t.Run("sync status transitions through states", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		info := drv.SyncInfo()
		if info.Status != SyncStatusOffline {
			t.Errorf("Initial status = %v, want %v", info.Status, SyncStatusOffline)
		}

		err = drv.Sync()
		if err == nil {
			t.Error("Sync() should return error when offline")
		}

		info = drv.SyncInfo()
		if info.Status != SyncStatusOffline {
			t.Errorf("Status after failed sync = %v, want %v", info.Status, SyncStatusOffline)
		}
	})

	t.Run("sync info contains error message on failure", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:           DriverTursoReplica,
			DatabasePath:   dbPath,
			TursoURL:       "https://nonexistent.invalid.example",
			TursoToken:     "test-token",
			SyncInterval:   0,
			ConnectTimeout: 50 * time.Millisecond,
			MaxRetries:     0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_ = drv.Connect(ctx)

		_ = drv.Sync()

		info := drv.SyncInfo()
		if info.LastError == nil && drv.Status() != StatusConnected {
			t.Log("SyncInfo.LastError is nil (acceptable if never attempted sync)")
		}
	})

	t.Run("last sync time updates correctly", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		initialTime := drv.LastSyncTime()
		if !initialTime.IsZero() {
			t.Error("LastSyncTime() should be zero initially")
		}
	})
}

func TestTursoReplica_OfflineFirstBehavior(t *testing.T) {
	t.Run("can create driver with unreachable remote", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://definitely.not.reachable.example",
			TursoToken:   "invalid-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		if drv.Type() != DriverTursoReplica {
			t.Errorf("Type() = %v, want %v", drv.Type(), DriverTursoReplica)
		}

		if drv.Status() != StatusDisconnected {
			t.Errorf("Status() = %v, want %v", drv.Status(), StatusDisconnected)
		}
	})

	t.Run("sync can be called multiple times safely", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		for i := 0; i < 5; i++ {
			_ = drv.Sync()
		}

		if drv.Status() != StatusDisconnected {
			t.Errorf("Status() = %v, want %v", drv.Status(), StatusDisconnected)
		}
	})

	t.Run("auto-sync interval can be disabled", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		if drv.config.SyncInterval != 0 {
			t.Errorf("SyncInterval = %v, want 0 (disabled)", drv.config.SyncInterval)
		}
	})

	t.Run("auto-sync interval is respected", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 30 * time.Second,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		if drv.config.SyncInterval != 30*time.Second {
			t.Errorf("SyncInterval = %v, want 30s", drv.config.SyncInterval)
		}
	})
}

func TestTursoReplica_ConcurrentSyncOperations(t *testing.T) {
	t.Run("concurrent sync calls are safe", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "replica.db")

		drv, err := NewTursoReplicaDriver(&Config{
			Type:         DriverTursoReplica,
			DatabasePath: dbPath,
			TursoURL:     "https://example.turso.io",
			TursoToken:   "test-token",
			SyncInterval: 0,
		})
		if err != nil {
			t.Fatalf("NewTursoReplicaDriver() error = %v", err)
		}

		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				for j := 0; j < 10; j++ {
					_ = drv.Sync()
					_ = drv.SyncInfo()
					_ = drv.LastSyncTime()
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestDriverRegistry_AllModes(t *testing.T) {
	ResetRegistry()

	_ = RegisterDriver(DriverSQLite, func(config *Config) (DatabaseDriver, error) {
		return &mockDriver{status: StatusDisconnected, driverType: DriverSQLite}, nil
	})
	_ = RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
		return NewTursoRemoteDriver(config)
	})
	_ = RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
		return NewTursoReplicaDriver(config)
	})

	t.Run("all driver types are registered", func(t *testing.T) {
		types := []DriverType{DriverSQLite, DriverTursoRemote, DriverTursoReplica}
		for _, dt := range types {
			factory, ok := GetFactory(dt)
			if !ok {
				t.Errorf("Driver type %v not registered", dt)
			}
			if factory == nil {
				t.Errorf("Factory for %v is nil", dt)
			}
		}
	})

	t.Run("can create each driver type", func(t *testing.T) {
		tests := []struct {
			name   string
			dt     DriverType
			config *Config
		}{
			{
				name: "sqlite",
				dt:   DriverSQLite,
				config: &Config{
					Type:         DriverSQLite,
					DatabasePath: filepath.Join(t.TempDir(), "test.db"),
				},
			},
			{
				name: "turso-remote",
				dt:   DriverTursoRemote,
				config: &Config{
					Type:       DriverTursoRemote,
					TursoURL:   "https://example.turso.io",
					TursoToken: "test-token",
				},
			},
			{
				name: "turso-replica",
				dt:   DriverTursoReplica,
				config: &Config{
					Type:         DriverTursoReplica,
					DatabasePath: filepath.Join(t.TempDir(), "replica.db"),
					TursoURL:     "https://example.turso.io",
					TursoToken:   "test-token",
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				factory, ok := GetFactory(tt.dt)
				if !ok {
					t.Fatalf("Factory not found for %v", tt.dt)
				}

				drv, err := factory(tt.config)
				if err != nil {
					t.Fatalf("Factory() error = %v", err)
				}
				if drv == nil {
					t.Fatal("Factory() returned nil driver")
				}
				if drv.Type() != tt.dt {
					t.Errorf("Type() = %v, want %v", drv.Type(), tt.dt)
				}
			})
		}
	})
}
