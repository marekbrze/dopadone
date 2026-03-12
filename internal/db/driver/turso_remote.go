//go:build !windows

package driver

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	libsql "github.com/tursodatabase/go-libsql"
)

type TursoRemoteDriver struct {
	db        *sql.DB
	connector *libsql.Connector
	config    *Config
	status    ConnectionStatus
	mu        sync.RWMutex
	tempDir   string
}

func NewTursoRemoteDriver(config *Config) (*TursoRemoteDriver, error) {
	if config == nil {
		return nil, NewDriverError(DriverTursoRemote, "create", ErrInvalidConfig)
	}

	if config.TursoURL == "" {
		return nil, NewDriverError(DriverTursoRemote, "create", fmt.Errorf("%w: turso url required", ErrInvalidConfig))
	}

	if config.TursoToken == "" {
		return nil, NewDriverError(DriverTursoRemote, "create", fmt.Errorf("%w: turso token required", ErrInvalidConfig))
	}

	return &TursoRemoteDriver{
		config: config,
		status: StatusDisconnected,
	}, nil
}

func (d *TursoRemoteDriver) Connect(ctx context.Context) error {
	d.setStatus(StatusConnecting)

	maxRetries := d.config.MaxRetries
	if maxRetries < 0 {
		maxRetries = 0
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				d.setStatus(StatusError)
				return ctx.Err()
			case <-time.After(d.config.RetryInterval):
			}
		}

		err := d.tryConnect(ctx)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	d.setStatus(StatusError)
	return lastErr
}

func (d *TursoRemoteDriver) tryConnect(ctx context.Context) error {
	connectCtx, cancel := context.WithTimeout(ctx, d.config.ConnectTimeout)
	defer cancel()

	tempDir, err := os.MkdirTemp("", "turso-remote-*")
	if err != nil {
		return NewDriverError(DriverTursoRemote, "connect", fmt.Errorf("%w: failed to create temp dir: %v", ErrConnectionFailed, err))
	}

	dbPath := filepath.Join(tempDir, "remote.db")

	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath,
		d.config.TursoURL,
		libsql.WithAuthToken(d.config.TursoToken),
		libsql.WithSyncInterval(0),
	)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return NewDriverError(DriverTursoRemote, "connect", fmt.Errorf("%w: %v", ErrConnectionFailed, err))
	}

	db := sql.OpenDB(connector)
	if err := db.PingContext(connectCtx); err != nil {
		_ = connector.Close()
		_ = os.RemoveAll(tempDir)
		return NewDriverError(DriverTursoRemote, "connect", fmt.Errorf("%w: %v", ErrConnectionFailed, err))
	}

	d.mu.Lock()
	d.db = db
	d.connector = connector
	d.tempDir = tempDir
	d.mu.Unlock()

	d.setStatus(StatusConnected)
	return nil
}

func (d *TursoRemoteDriver) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.status = StatusDisconnected

	if d.db == nil {
		return nil
	}

	var errs []error

	if err := d.db.Close(); err != nil {
		errs = append(errs, err)
	}
	d.db = nil

	if d.connector != nil {
		if err := d.connector.Close(); err != nil {
			errs = append(errs, err)
		}
		d.connector = nil
	}

	if d.tempDir != "" {
		if err := os.RemoveAll(d.tempDir); err != nil {
			errs = append(errs, err)
		}
		d.tempDir = ""
	}

	if len(errs) > 0 {
		return NewDriverError(DriverTursoRemote, "close", errorsJoin(errs...))
	}
	return nil
}

func (d *TursoRemoteDriver) GetDB() *sql.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db
}

func (d *TursoRemoteDriver) Ping(ctx context.Context) error {
	d.mu.RLock()
	db := d.db
	d.mu.RUnlock()

	if db == nil {
		return NewDriverError(DriverTursoRemote, "ping", ErrDriverAlreadyClosed)
	}

	if err := db.PingContext(ctx); err != nil {
		return NewDriverError(DriverTursoRemote, "ping", err)
	}
	return nil
}

func (d *TursoRemoteDriver) Type() DriverType {
	return DriverTursoRemote
}

func (d *TursoRemoteDriver) Status() ConnectionStatus {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.status
}

func (d *TursoRemoteDriver) setStatus(status ConnectionStatus) {
	d.mu.Lock()
	d.status = status
	d.mu.Unlock()
}

func init() {
	err := RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
		return NewTursoRemoteDriver(config)
	})
	if err != nil {
		panic(fmt.Sprintf("failed to register turso-remote driver: %v", err))
	}
}
