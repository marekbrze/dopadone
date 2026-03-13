//go:build !windows && !(darwin && amd64)

package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	libsql "github.com/tursodatabase/go-libsql"
)

type SyncStatus string

const (
	SyncStatusIdle    SyncStatus = "idle"
	SyncStatusSyncing SyncStatus = "syncing"
	SyncStatusError   SyncStatus = "error"
	SyncStatusOffline SyncStatus = "offline"
)

type SyncInfo struct {
	Status     SyncStatus
	LastSyncAt time.Time
	LastError  error
}

type TursoReplicaDriver struct {
	db        *sql.DB
	connector *libsql.Connector
	config    *Config
	status    ConnectionStatus
	mu        sync.RWMutex

	syncStatus  SyncStatus
	lastSyncAt  time.Time
	lastSyncErr error

	syncCtx    context.Context
	syncCancel context.CancelFunc
	syncWg     sync.WaitGroup
}

func NewTursoReplicaDriver(config *Config) (*TursoReplicaDriver, error) {
	if config == nil {
		return nil, NewDriverError(DriverTursoReplica, "create", ErrInvalidConfig)
	}

	if config.TursoURL == "" {
		return nil, NewDriverError(DriverTursoReplica, "create", fmt.Errorf("%w: turso url required", ErrInvalidConfig))
	}

	if config.TursoToken == "" {
		return nil, NewDriverError(DriverTursoReplica, "create", fmt.Errorf("%w: turso token required", ErrInvalidConfig))
	}

	if config.DatabasePath == "" {
		return nil, NewDriverError(DriverTursoReplica, "create", fmt.Errorf("%w: database path required", ErrInvalidConfig))
	}

	return &TursoReplicaDriver{
		config:     config,
		status:     StatusDisconnected,
		syncStatus: SyncStatusOffline,
	}, nil
}

func (d *TursoReplicaDriver) Connect(ctx context.Context) error {
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

func (d *TursoReplicaDriver) tryConnect(ctx context.Context) error {
	connectCtx, cancel := context.WithTimeout(ctx, d.config.ConnectTimeout)
	defer cancel()

	connector, err := libsql.NewEmbeddedReplicaConnector(
		d.config.DatabasePath,
		d.config.TursoURL,
		libsql.WithAuthToken(d.config.TursoToken),
	)
	if err != nil {
		return NewDriverError(DriverTursoReplica, "connect", fmt.Errorf("%w: %v", ErrConnectionFailed, err))
	}

	db := sql.OpenDB(connector)
	if err := db.PingContext(connectCtx); err != nil {
		_ = connector.Close()
		return NewDriverError(DriverTursoReplica, "connect", fmt.Errorf("%w: %v", ErrConnectionFailed, err))
	}

	d.mu.Lock()
	d.db = db
	d.connector = connector
	d.syncStatus = SyncStatusIdle
	d.mu.Unlock()

	d.setStatus(StatusConnected)

	if err := d.Sync(); err != nil {
		log.Printf("[TursoReplica] Initial sync failed: %v", err)
	}

	if d.config.SyncInterval > 0 {
		d.startAutoSync()
	}

	return nil
}

func (d *TursoReplicaDriver) Sync() error {
	d.mu.Lock()
	if d.connector == nil {
		d.mu.Unlock()
		return NewDriverError(DriverTursoReplica, "sync", ErrDriverAlreadyClosed)
	}
	d.syncStatus = SyncStatusSyncing
	d.mu.Unlock()

	replicated, err := d.connector.Sync()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err != nil {
		d.syncStatus = SyncStatusError
		d.lastSyncErr = err
		log.Printf("[TursoReplica] Sync failed: %v", err)
		return NewDriverError(DriverTursoReplica, "sync", err)
	}

	d.syncStatus = SyncStatusIdle
	d.lastSyncAt = time.Now()
	d.lastSyncErr = nil

	log.Printf("[TursoReplica] Synced %d frames (frame_no: %d)", replicated.FramesSynced, replicated.FrameNo)

	return nil
}

func (d *TursoReplicaDriver) startAutoSync() {
	d.syncCtx, d.syncCancel = context.WithCancel(context.Background())

	d.syncWg.Add(1)
	go func() {
		defer d.syncWg.Done()

		ticker := time.NewTicker(d.config.SyncInterval)
		defer ticker.Stop()

		for {
			select {
			case <-d.syncCtx.Done():
				return
			case <-ticker.C:
				if err := d.Sync(); err != nil {
					log.Printf("[TursoReplica] Auto-sync failed: %v", err)
				}
			}
		}
	}()
}

func (d *TursoReplicaDriver) Close() error {
	if d.syncCancel != nil {
		d.syncCancel()
		d.syncWg.Wait()
		d.syncCancel = nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.status = StatusDisconnected
	d.syncStatus = SyncStatusOffline

	if d.db == nil {
		return nil
	}

	var errs []error

	if d.connector != nil {
		if err := d.connector.Close(); err != nil {
			errs = append(errs, err)
		}
		d.connector = nil
	}

	if err := d.db.Close(); err != nil {
		errs = append(errs, err)
	}
	d.db = nil

	if len(errs) > 0 {
		return NewDriverError(DriverTursoReplica, "close", errorsJoin(errs...))
	}
	return nil
}

func (d *TursoReplicaDriver) GetDB() *sql.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db
}

func (d *TursoReplicaDriver) Ping(ctx context.Context) error {
	d.mu.RLock()
	db := d.db
	d.mu.RUnlock()

	if db == nil {
		return NewDriverError(DriverTursoReplica, "ping", ErrDriverAlreadyClosed)
	}

	if err := db.PingContext(ctx); err != nil {
		return NewDriverError(DriverTursoReplica, "ping", err)
	}
	return nil
}

func (d *TursoReplicaDriver) Type() DriverType {
	return DriverTursoReplica
}

func (d *TursoReplicaDriver) Status() ConnectionStatus {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.syncStatus == SyncStatusSyncing {
		return StatusConnecting
	}

	return d.status
}

func (d *TursoReplicaDriver) setStatus(status ConnectionStatus) {
	d.mu.Lock()
	d.status = status
	d.mu.Unlock()
}

func (d *TursoReplicaDriver) SyncInfo() SyncInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return SyncInfo{
		Status:     d.syncStatus,
		LastSyncAt: d.lastSyncAt,
		LastError:  d.lastSyncErr,
	}
}

func (d *TursoReplicaDriver) LastSyncTime() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.lastSyncAt
}

func errorsJoin(errs ...error) error {
	var msg string
	for _, err := range errs {
		if msg != "" {
			msg += "; "
		}
		msg += err.Error()
	}
	return fmt.Errorf("%s", msg)
}

func init() {
	err := RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
		return NewTursoReplicaDriver(config)
	})
	if err != nil {
		panic(fmt.Sprintf("failed to register turso-replica driver: %v", err))
	}
}
