package driver

import (
	"time"
)

type DriverType string

const (
	DriverSQLite       DriverType = "sqlite"
	DriverTursoRemote  DriverType = "turso-remote"
	DriverTursoReplica DriverType = "turso-replica"
)

type Config struct {
	Type           DriverType
	DatabasePath   string
	TursoURL       string
	TursoToken     string
	SyncInterval   time.Duration
	ConnectTimeout time.Duration
	MaxRetries     int
	RetryInterval  time.Duration
}

type Option func(*Config)

func WithDatabasePath(path string) Option {
	return func(c *Config) {
		c.DatabasePath = path
	}
}

func WithTurso(url, token string) Option {
	return func(c *Config) {
		c.TursoURL = url
		c.TursoToken = token
	}
}

func WithSyncInterval(d time.Duration) Option {
	return func(c *Config) {
		c.SyncInterval = d
	}
}

func WithDriverType(dt DriverType) Option {
	return func(c *Config) {
		c.Type = dt
	}
}

func WithConnectTimeout(d time.Duration) Option {
	return func(c *Config) {
		c.ConnectTimeout = d
	}
}

func WithMaxRetries(n int) Option {
	return func(c *Config) {
		c.MaxRetries = n
	}
}

func WithRetryInterval(d time.Duration) Option {
	return func(c *Config) {
		c.RetryInterval = d
	}
}

func DefaultConfig() *Config {
	return &Config{
		Type:           DriverSQLite,
		SyncInterval:   60 * time.Second,
		ConnectTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryInterval:  1 * time.Second,
	}
}

func (c *Config) Validate() error {
	switch c.Type {
	case DriverSQLite:
		if c.DatabasePath == "" {
			return NewDriverError(c.Type, "validate", ErrInvalidConfig)
		}
	case DriverTursoRemote:
		if c.TursoURL == "" || c.TursoToken == "" {
			return NewDriverError(c.Type, "validate", ErrInvalidConfig)
		}
	case DriverTursoReplica:
		if c.TursoURL == "" || c.TursoToken == "" || c.DatabasePath == "" {
			return NewDriverError(c.Type, "validate", ErrInvalidConfig)
		}
	default:
		return NewDriverError(c.Type, "validate", ErrDriverNotRegistered)
	}
	return nil
}
