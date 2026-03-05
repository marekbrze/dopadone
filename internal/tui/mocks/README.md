# TUI Service Mocks

This directory contains mock implementations of service interfaces used for testing the TUI layer.

## Overview

The mocks follow the **func-field pattern** for maximum flexibility in testing. Each mock struct contains function fields that can be overridden per-test to configure specific behaviors.

## Mock Implementations

### Service Mocks

- **MockAreaService** - Implements `service.AreaServiceInterface` (9 methods)
- **MockSubareaService** - Implements `service.SubareaServiceInterface` (9 methods)
- **MockProjectService** - Implements `service.ProjectServiceInterface` (13 methods)
- **MockTaskService** - Implements `service.TaskServiceInterface` (14 methods)

### Design Pattern

Each mock uses the **func-field pattern**:

```go
type MockAreaService struct {
    ListFunc func(ctx context.Context) ([]domain.Area, error)
    // ... other function fields
}

func (m *MockAreaService) List(ctx context.Context) ([]domain.Area, error) {
    if m.ListFunc != nil {
        return m.ListFunc(ctx)
    }
    return []domain.Area{}, nil  // Safe default
}
```

**Benefits**:
- Maximum flexibility: Tests can override specific methods as needed
- Safe defaults: Unset function fields return zero values (no panics)
- Table-driven testing: Easy to configure different scenarios per test case
- Type-safe: Compile-time verification of mock implementations

## Helper Functions

### NewMockServices()

Creates all 4 mock services with default implementations:

```go
mocks := mocks.NewMockServices()
model := InitialModel(
    mocks.AreaSvc,
    mocks.SubareaSvc,
    mocks.ProjectSvc,
    mocks.TaskSvc,
)
```

### Setup Helpers

Pre-configured helpers for common scenarios:

```go
// Success scenarios
mocks.SetupMockAreaSuccess(areas)
mocks.SetupMockSubareaSuccess(subareas)
mocks.SetupMockProjectSuccess(projects)
mocks.SetupMockTaskSuccess(tasks)

// Error scenarios
mocks.SetupMockAreaError(errors.New("database error"))
mocks.SetupMockProjectError(errors.New("query failed"))

// Create operation helpers
mocks.SetupMockAreaCreate()
mocks.SetupMockProjectCreate()
```

## Usage Examples

### Basic Test with Mocks

```go
func TestLoadAreas(t *testing.T) {
    mocks := mocks.NewMockServices()
    
    expectedAreas := []domain.Area{
        {ID: "area-1", Name: "Work"},
        {ID: "area-2", Name: "Personal"},
    }
    
    mocks.SetupMockAreaSuccess(expectedAreas)
    
    model := InitialModel(
        mocks.AreaSvc,
        mocks.SubareaSvc,
        mocks.ProjectSvc,
        mocks.TaskSvc,
    )
    
    // Test logic here
}
```

### Custom Mock Behavior

```go
func TestAreaError(t *testing.T) {
    mocks := mocks.NewMockServices()
    
    // Custom error scenario
    mocks.AreaSvc.ListFunc = func(ctx context.Context) ([]domain.Area, error) {
        return nil, errors.New("database connection failed")
    }
    
    model := InitialModel(
        mocks.AreaSvc,
        mocks.SubareaSvc,
        mocks.ProjectSvc,
        mocks.TaskSvc,
    )
    
    // Test error handling
}
```

### Table-Driven Tests

```go
func TestAreaScenarios(t *testing.T) {
    tests := []struct {
        name     string
        setup    func(*mocks.MockServices)
        expected []domain.Area
        hasError bool
    }{
        {
            name: "success",
            setup: func(m *mocks.MockServices) {
                m.SetupMockAreaSuccess([]domain.Area{{ID: "1", Name: "Test"}})
            },
            expected: []domain.Area{{ID: "1", Name: "Test"}},
            hasError: false,
        },
        {
            name: "database error",
            setup: func(m *mocks.MockServices) {
                m.SetupMockAreaError(errors.New("db error"))
            },
            expected: nil,
            hasError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mocks := mocks.NewMockServices()
            tt.setup(&mocks)
            // Test logic
        })
    }
}
```

## Testing Strategy

The TUI tests use these mocks to:

1. **Isolate TUI layer**: No database dependencies in unit tests
2. **Test edge cases**: Network errors, empty results, malformed data
3. **Verify interactions**: Ensure correct service methods are called
4. **Enable parallel testing**: Each test has its own mock instances
5. **Speed up test suite**: No I/O operations, pure in-memory testing

## Migration from db.Querier

Previous tests used `MockQuerier` (database layer). Migration steps:

1. Replace `MockQuerier` with `mocks.NewMockServices()`
2. Pass 4 service interfaces to `InitialModel()`
3. Configure mocks using helper functions or custom `*Func` fields
4. Remove direct database dependencies

## Related Documentation

- [Service Layer Documentation](../../../internal/service/README.md)
- [TUI Architecture](../../../docs/TUI.md)
- [Testing Guidelines](../../../docs/TESTING.md)
