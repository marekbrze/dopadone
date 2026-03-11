package driver

import (
	"context"
	"database/sql"
)

type ConnectionStatus string

const (
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusConnected    ConnectionStatus = "connected"
	StatusConnecting   ConnectionStatus = "connecting"
	StatusError        ConnectionStatus = "error"
)

type DatabaseDriver interface {
	Connect(ctx context.Context) error
	Close() error
	GetDB() *sql.DB
	Ping(ctx context.Context) error
	Type() DriverType
	Status() ConnectionStatus
}
