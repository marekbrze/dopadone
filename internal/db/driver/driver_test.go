package driver

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"testing"
	"time"
)

type mockDriver struct {
	db         *sql.DB
	status     ConnectionStatus
	driverType DriverType
	connectErr error
	pingErr    error
	closeErr   error
}

func newMockDriver(config *Config) (DatabaseDriver, error) {
	return &mockDriver{
		status:     StatusDisconnected,
		driverType: config.Type,
	}, nil
}

func (m *mockDriver) Connect(ctx context.Context) error {
	if m.connectErr != nil {
		m.status = StatusError
		return m.connectErr
	}
	m.status = StatusConnected
	return nil
}

func (m *mockDriver) Close() error {
	if m.closeErr != nil {
		return m.closeErr
	}
	m.status = StatusDisconnected
	return nil
}

func (m *mockDriver) GetDB() *sql.DB {
	return m.db
}

func (m *mockDriver) Ping(ctx context.Context) error {
	return m.pingErr
}

func (m *mockDriver) Type() DriverType {
	return m.driverType
}

func (m *mockDriver) Status() ConnectionStatus {
	return m.status
}

func TestDatabaseDriverInterface(t *testing.T) {
	tests := []struct {
		name       string
		driver     DatabaseDriver
		wantType   DriverType
		wantStatus ConnectionStatus
	}{
		{
			name:       "mock driver initial state",
			driver:     &mockDriver{status: StatusDisconnected, driverType: DriverSQLite},
			wantType:   DriverSQLite,
			wantStatus: StatusDisconnected,
		},
		{
			name:       "mock driver connected",
			driver:     &mockDriver{status: StatusConnected, driverType: DriverTursoRemote},
			wantType:   DriverTursoRemote,
			wantStatus: StatusConnected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.driver.Type(); got != tt.wantType {
				t.Errorf("Type() = %v, want %v", got, tt.wantType)
			}
			if got := tt.driver.Status(); got != tt.wantStatus {
				t.Errorf("Status() = %v, want %v", got, tt.wantStatus)
			}
		})
	}
}

func TestDriverRegistration(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	t.Run("register new driver", func(t *testing.T) {
		err := RegisterDriver("test-driver", newMockDriver)
		if err != nil {
			t.Errorf("RegisterDriver() error = %v", err)
		}

		factory, ok := GetFactory("test-driver")
		if !ok {
			t.Error("GetFactory() driver not found")
		}
		if factory == nil {
			t.Error("GetFactory() returned nil factory")
		}
	})

	t.Run("register duplicate driver returns error", func(t *testing.T) {
		err := RegisterDriver("test-driver-dup", newMockDriver)
		if err != nil {
			t.Fatalf("first registration failed: %v", err)
		}

		err = RegisterDriver("test-driver-dup", newMockDriver)
		if err == nil {
			t.Error("duplicate registration should return error")
		}
	})

	t.Run("register nil factory returns error", func(t *testing.T) {
		err := RegisterDriver("nil-factory", nil)
		if err == nil {
			t.Error("nil factory registration should return error")
		}
	})

	t.Run("get unregistered factory", func(t *testing.T) {
		_, ok := GetFactory("unregistered")
		if ok {
			t.Error("GetFactory() should return false for unregistered driver")
		}
	})

	t.Run("registered drivers list", func(t *testing.T) {
		ResetRegistry()
		_ = RegisterDriver("driver-a", newMockDriver)
		_ = RegisterDriver("driver-b", newMockDriver)

		drivers := RegisteredDrivers()
		if len(drivers) != 2 {
			t.Errorf("RegisteredDrivers() returned %d drivers, want 2", len(drivers))
		}
	})

	t.Run("unregister driver", func(t *testing.T) {
		ResetRegistry()
		_ = RegisterDriver("to-remove", newMockDriver)

		UnregisterDriver("to-remove")

		_, ok := GetFactory("to-remove")
		if ok {
			t.Error("driver should be unregistered")
		}
	})
}

func TestConcurrentRegistry(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	var wg sync.WaitGroup
	errCh := make(chan error, 100)

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			driverType := DriverType("concurrent-driver")
			err := RegisterDriver(driverType, newMockDriver)
			if err == nil {
				UnregisterDriver(driverType)
			}
		}(i)

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = RegisteredDrivers()
		}()

		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, _ = GetFactory(DriverType("some-driver"))
		}(i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			t.Errorf("concurrent access error: %v", err)
		}
	}
}

func TestFactoryPattern(t *testing.T) {
	ResetRegistry()
	defer ResetRegistry()

	_ = RegisterDriver(DriverSQLite, newMockDriver)

	t.Run("create driver with valid config", func(t *testing.T) {
		driver, err := NewDriver(WithDriverType(DriverSQLite), WithDatabasePath("/tmp/test.db"))
		if err != nil {
			t.Fatalf("NewDriver() error = %v", err)
		}
		if driver == nil {
			t.Fatal("NewDriver() returned nil driver")
		}
		if driver.Type() != DriverSQLite {
			t.Errorf("Type() = %v, want %v", driver.Type(), DriverSQLite)
		}
	})

	t.Run("create driver with unregistered type returns error", func(t *testing.T) {
		_, err := NewDriver(WithDriverType("unregistered"))
		if err == nil {
			t.Error("NewDriver() should return error for unregistered driver")
		}
	})

	t.Run("create driver with invalid config returns error", func(t *testing.T) {
		_, err := NewDriver(WithDriverType(DriverSQLite))
		if err == nil {
			t.Error("NewDriver() should return error for missing database path")
		}
	})
}

func TestConfigOptions(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		config := DefaultConfig()
		if config.Type != DriverSQLite {
			t.Errorf("default Type = %v, want %v", config.Type, DriverSQLite)
		}
		if config.SyncInterval != 60*time.Second {
			t.Errorf("default SyncInterval = %v, want %v", config.SyncInterval, 60*time.Second)
		}
	})

	t.Run("with database path", func(t *testing.T) {
		config := DefaultConfig()
		WithDatabasePath("/custom/path.db")(config)
		if config.DatabasePath != "/custom/path.db" {
			t.Errorf("DatabasePath = %v, want /custom/path.db", config.DatabasePath)
		}
	})

	t.Run("with turso config", func(t *testing.T) {
		config := DefaultConfig()
		WithTurso("https://example.turso.io", "token123")(config)
		if config.TursoURL != "https://example.turso.io" {
			t.Errorf("TursoURL = %v, want https://example.turso.io", config.TursoURL)
		}
		if config.TursoToken != "token123" {
			t.Errorf("TursoToken = %v, want token123", config.TursoToken)
		}
	})

	t.Run("with sync interval", func(t *testing.T) {
		config := DefaultConfig()
		WithSyncInterval(30 * time.Second)(config)
		if config.SyncInterval != 30*time.Second {
			t.Errorf("SyncInterval = %v, want 30s", config.SyncInterval)
		}
	})

	t.Run("with driver type", func(t *testing.T) {
		config := DefaultConfig()
		WithDriverType(DriverTursoRemote)(config)
		if config.Type != DriverTursoRemote {
			t.Errorf("Type = %v, want %v", config.Type, DriverTursoRemote)
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		config := DefaultConfig()
		WithDriverType(DriverSQLite)(config)
		WithDatabasePath("/path/to/db.sqlite")(config)
		WithSyncInterval(120 * time.Second)(config)

		if config.Type != DriverSQLite {
			t.Errorf("Type = %v, want %v", config.Type, DriverSQLite)
		}
		if config.DatabasePath != "/path/to/db.sqlite" {
			t.Errorf("DatabasePath = %v, want /path/to/db.sqlite", config.DatabasePath)
		}
		if config.SyncInterval != 120*time.Second {
			t.Errorf("SyncInterval = %v, want 120s", config.SyncInterval)
		}
	})
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid sqlite config",
			config: &Config{
				Type:         DriverSQLite,
				DatabasePath: "/tmp/test.db",
			},
			wantErr: false,
		},
		{
			name: "invalid sqlite config - missing path",
			config: &Config{
				Type: DriverSQLite,
			},
			wantErr: true,
		},
		{
			name: "valid turso remote config",
			config: &Config{
				Type:       DriverTursoRemote,
				TursoURL:   "https://example.turso.io",
				TursoToken: "token123",
			},
			wantErr: false,
		},
		{
			name: "invalid turso remote config - missing url",
			config: &Config{
				Type:       DriverTursoRemote,
				TursoToken: "token123",
			},
			wantErr: true,
		},
		{
			name: "invalid turso remote config - missing token",
			config: &Config{
				Type:     DriverTursoRemote,
				TursoURL: "https://example.turso.io",
			},
			wantErr: true,
		},
		{
			name: "valid turso replica config",
			config: &Config{
				Type:         DriverTursoReplica,
				DatabasePath: "/tmp/replica.db",
				TursoURL:     "https://example.turso.io",
				TursoToken:   "token123",
			},
			wantErr: false,
		},
		{
			name: "invalid turso replica config - missing path",
			config: &Config{
				Type:       DriverTursoReplica,
				TursoURL:   "https://example.turso.io",
				TursoToken: "token123",
			},
			wantErr: true,
		},
		{
			name: "invalid driver type",
			config: &Config{
				Type: "unknown",
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

func TestErrors(t *testing.T) {
	t.Run("driver error formatting", func(t *testing.T) {
		err := NewDriverError(DriverSQLite, "connect", ErrConnectionFailed)
		expected := "driver sqlite: connect: connection failed"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("driver error without driver type", func(t *testing.T) {
		err := &DriverError{Op: "test", Err: ErrInvalidConfig}
		expected := "test: invalid configuration"
		if err.Error() != expected {
			t.Errorf("Error() = %v, want %v", err.Error(), expected)
		}
	})

	t.Run("error unwrapping", func(t *testing.T) {
		err := NewDriverError(DriverSQLite, "connect", ErrConnectionFailed)
		if !errors.Is(err, ErrConnectionFailed) {
			t.Error("errors.Is should match ErrConnectionFailed")
		}

		var driverErr *DriverError
		if !errors.As(err, &driverErr) {
			t.Error("errors.As should extract DriverError")
		}
	})

	t.Run("sentinel errors", func(t *testing.T) {
		tests := []struct {
			name string
			err  error
			str  string
		}{
			{"ErrDriverNotRegistered", ErrDriverNotRegistered, "driver not registered"},
			{"ErrConnectionFailed", ErrConnectionFailed, "connection failed"},
			{"ErrInvalidConfig", ErrInvalidConfig, "invalid configuration"},
			{"ErrDriverAlreadyClosed", ErrDriverAlreadyClosed, "driver already closed"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.err.Error() != tt.str {
					t.Errorf("error string = %v, want %v", tt.err.Error(), tt.str)
				}
			})
		}
	})
}

func TestMockDriverConnect(t *testing.T) {
	t.Run("successful connect", func(t *testing.T) {
		driver := &mockDriver{status: StatusDisconnected, driverType: DriverSQLite}
		err := driver.Connect(context.Background())
		if err != nil {
			t.Errorf("Connect() error = %v", err)
		}
		if driver.Status() != StatusConnected {
			t.Errorf("Status() = %v, want %v", driver.Status(), StatusConnected)
		}
	})

	t.Run("connect with error", func(t *testing.T) {
		driver := &mockDriver{
			status:     StatusDisconnected,
			driverType: DriverSQLite,
			connectErr: errors.New("connection refused"),
		}
		err := driver.Connect(context.Background())
		if err == nil {
			t.Error("Connect() should return error")
		}
		if driver.Status() != StatusError {
			t.Errorf("Status() = %v, want %v", driver.Status(), StatusError)
		}
	})
}

func TestMockDriverClose(t *testing.T) {
	t.Run("successful close", func(t *testing.T) {
		driver := &mockDriver{status: StatusConnected, driverType: DriverSQLite}
		err := driver.Close()
		if err != nil {
			t.Errorf("Close() error = %v", err)
		}
		if driver.Status() != StatusDisconnected {
			t.Errorf("Status() = %v, want %v", driver.Status(), StatusDisconnected)
		}
	})

	t.Run("close with error", func(t *testing.T) {
		driver := &mockDriver{
			status:     StatusConnected,
			driverType: DriverSQLite,
			closeErr:   errors.New("close failed"),
		}
		err := driver.Close()
		if err == nil {
			t.Error("Close() should return error")
		}
	})
}

func TestMockDriverPing(t *testing.T) {
	t.Run("successful ping", func(t *testing.T) {
		driver := &mockDriver{driverType: DriverSQLite}
		err := driver.Ping(context.Background())
		if err != nil {
			t.Errorf("Ping() error = %v", err)
		}
	})

	t.Run("ping with error", func(t *testing.T) {
		driver := &mockDriver{
			driverType: DriverSQLite,
			pingErr:    errors.New("ping failed"),
		}
		err := driver.Ping(context.Background())
		if err == nil {
			t.Error("Ping() should return error")
		}
	})
}

func TestConnectionStatus(t *testing.T) {
	tests := []struct {
		name   string
		status ConnectionStatus
		str    string
	}{
		{"disconnected", StatusDisconnected, "disconnected"},
		{"connected", StatusConnected, "connected"},
		{"connecting", StatusConnecting, "connecting"},
		{"error", StatusError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.str {
				t.Errorf("status string = %v, want %v", tt.status, tt.str)
			}
		})
	}
}

func TestDriverType(t *testing.T) {
	tests := []struct {
		name       string
		driverType DriverType
		str        string
	}{
		{"sqlite", DriverSQLite, "sqlite"},
		{"turso-remote", DriverTursoRemote, "turso-remote"},
		{"turso-replica", DriverTursoReplica, "turso-replica"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.driverType) != tt.str {
				t.Errorf("driver type string = %v, want %v", tt.driverType, tt.str)
			}
		})
	}
}
