package driver

import (
	"sync"
)

type DriverFactory func(config *Config) (DatabaseDriver, error)

var registry = struct {
	mu        sync.RWMutex
	factories map[DriverType]DriverFactory
}{
	factories: make(map[DriverType]DriverFactory),
}

func RegisterDriver(dt DriverType, factory DriverFactory) error {
	if factory == nil {
		return NewDriverError(dt, "register", ErrInvalidConfig)
	}

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if _, exists := registry.factories[dt]; exists {
		return NewDriverError(dt, "register", ErrInvalidConfig)
	}

	registry.factories[dt] = factory
	return nil
}

func GetFactory(dt DriverType) (DriverFactory, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	factory, ok := registry.factories[dt]
	return factory, ok
}

func RegisteredDrivers() []DriverType {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	drivers := make([]DriverType, 0, len(registry.factories))
	for dt := range registry.factories {
		drivers = append(drivers, dt)
	}
	return drivers
}

func UnregisterDriver(dt DriverType) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	delete(registry.factories, dt)
}

func ResetRegistry() {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.factories = make(map[DriverType]DriverFactory)
}
