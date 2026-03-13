//go:build windows || (darwin && amd64)

package driver

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrTursoNotSupported = errors.New("turso database is not supported on Windows")

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

type unsupportedDriver struct {
	driverType DriverType
}

func (d *unsupportedDriver) Connect(ctx context.Context) error {
	return ErrTursoNotSupported
}

func (d *unsupportedDriver) Close() error {
	return ErrTursoNotSupported
}

func (d *unsupportedDriver) GetDB() *sql.DB {
	return nil
}

func (d *unsupportedDriver) Ping(ctx context.Context) error {
	return ErrTursoNotSupported
}

func (d *unsupportedDriver) Type() DriverType {
	return d.driverType
}

func (d *unsupportedDriver) Status() ConnectionStatus {
	return StatusError
}

func init() {
	RegisterDriver(DriverTursoRemote, func(config *Config) (DatabaseDriver, error) {
		return &unsupportedDriver{driverType: DriverTursoRemote}, nil
	})
	RegisterDriver(DriverTursoReplica, func(config *Config) (DatabaseDriver, error) {
		return &unsupportedDriver{driverType: DriverTursoReplica}, nil
	})
}
