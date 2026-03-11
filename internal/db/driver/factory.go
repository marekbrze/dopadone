package driver

func NewDriver(opts ...Option) (DatabaseDriver, error) {
	config := DefaultConfig()
	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	factory, ok := GetFactory(config.Type)
	if !ok {
		return nil, NewDriverError(config.Type, "create", ErrDriverNotRegistered)
	}

	return factory(config)
}
