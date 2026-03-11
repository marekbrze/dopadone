package driver

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ResetRegistry()
	_ = RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
		return NewTursoRemoteDriver(config)
	})
	m.Run()
}

func TestNewTursoRemoteDriver(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Type:       DriverTursoRemote,
				TursoURL:   "libsql://example.turso.io",
				TursoToken: "test-token",
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
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
				TursoURL: "libsql://example.turso.io",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, err := NewTursoRemoteDriver(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTursoRemoteDriver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && driver == nil {
				t.Error("NewTursoRemoteDriver() returned nil driver")
			}
		})
	}
}

func TestTursoRemoteDriver_Type(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	if driver.Type() != DriverTursoRemote {
		t.Errorf("Type() = %v, want %v", driver.Type(), DriverTursoRemote)
	}
}

func TestTursoRemoteDriver_Status(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	if driver.Status() != StatusDisconnected {
		t.Errorf("Status() = %v, want %v", driver.Status(), StatusDisconnected)
	}
}

func TestTursoRemoteDriver_GetDB_BeforeConnect(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	if driver.GetDB() != nil {
		t.Error("GetDB() should return nil before connect")
	}
}

func TestTursoRemoteDriver_Close_WithoutConnect(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	if err := driver.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestTursoRemoteDriver_Ping_WithoutConnect(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_Connect_ContextCancellation(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://example.turso.io",
		TursoToken:     "test-token",
		ConnectTimeout: 5 * time.Second,
		MaxRetries:     10,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_Init(t *testing.T) {
	_ = RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
		return NewTursoRemoteDriver(config)
	})

	factory, ok := GetFactory(DriverTursoRemote)
	if !ok {
		t.Error("TursoRemoteDriver should be registered")
	}
	if factory == nil {
		t.Error("Factory should not be nil")
	}

	drv, err := factory(&Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	})
	if err != nil {
		t.Errorf("Factory returned error: %v", err)
	}
	if drv == nil {
		t.Error("Factory returned nil driver")
	}
	if drv.Type() != DriverTursoRemote {
		t.Errorf("Driver type = %v, want %v", drv.Type(), DriverTursoRemote)
	}
}

func TestTursoRemoteDriver_Connect_ZeroRetries(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = driver.Connect(ctx)
	if err == nil {
		_ = driver.Close()
		t.Skip("Connection succeeded unexpectedly, skipping test")
	}
}

func TestTursoRemoteDriver_DefaultsApplied(t *testing.T) {
	config := &Config{
		Type:       DriverTursoRemote,
		TursoURL:   "libsql://example.turso.io",
		TursoToken: "test-token",
	}

	driver, err := NewTursoRemoteDriver(config)
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	if driver.config.Type != DriverTursoRemote {
		t.Errorf("Type = %v, want %v", driver.config.Type, DriverTursoRemote)
	}
}

func TestTursoRemoteDriver_StatusTransitions(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://example.turso.io",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_Close_AfterFailedConnect(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 1 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_Connect_RetryAttempts(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 10 * time.Millisecond,
		MaxRetries:     2,
		RetryInterval:  5 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_Connect_CancelDuringRetry(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://nonexistent.invalid.example",
		TursoToken:     "test-token",
		ConnectTimeout: 100 * time.Millisecond,
		MaxRetries:     10,
		RetryInterval:  1 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
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

func TestTursoRemoteDriver_ConcurrentStatusAccess(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://example.turso.io",
		TursoToken:     "test-token",
		ConnectTimeout: 10 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	done := make(chan bool)
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

func TestTursoRemoteDriver_ConcurrentGetDBAccess(t *testing.T) {
	driver, err := NewTursoRemoteDriver(&Config{
		Type:           DriverTursoRemote,
		TursoURL:       "libsql://example.turso.io",
		TursoToken:     "test-token",
		ConnectTimeout: 10 * time.Millisecond,
		MaxRetries:     0,
		RetryInterval:  0,
	})
	if err != nil {
		t.Fatalf("NewTursoRemoteDriver() error = %v", err)
	}

	done := make(chan bool)
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

func TestTursoRemoteDriver_ConfigOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []Option
		config *Config
	}{
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
			},
			config: &Config{
				ConnectTimeout: 15 * time.Second,
				MaxRetries:     3,
				RetryInterval:  500 * time.Millisecond,
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
		})
	}
}
