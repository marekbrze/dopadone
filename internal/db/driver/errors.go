package driver

import (
	"errors"
	"fmt"
)

var (
	ErrDriverNotRegistered = errors.New("driver not registered")
	ErrConnectionFailed    = errors.New("connection failed")
	ErrInvalidConfig       = errors.New("invalid configuration")
	ErrDriverAlreadyClosed = errors.New("driver already closed")
)

type DriverError struct {
	Driver DriverType
	Op     string
	Err    error
}

func (e *DriverError) Error() string {
	if e.Driver != "" {
		return fmt.Sprintf("driver %s: %s: %v", e.Driver, e.Op, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *DriverError) Unwrap() error {
	return e.Err
}

func NewDriverError(driver DriverType, op string, err error) *DriverError {
	return &DriverError{
		Driver: driver,
		Op:     op,
		Err:    err,
	}
}
