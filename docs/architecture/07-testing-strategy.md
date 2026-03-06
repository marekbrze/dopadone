# Testing Strategy

## Overview

Dopadone follows a **comprehensive testing approach** combining table-driven tests, interface mocking, test helpers, and golden files. The testing is integrated directly with the domain and service layers.

**Key Characteristics**:
- **Table-driven tests**: Clear, declarative test cases
- **Interface mocking**: Isolate units for testing
- **Test helpers**: Reduce boilerplate
- **Golden files**: Snapshot testing for complex output
- **Benchmarks**: Performance testing

---

## Table-Driven Tests

### Pattern Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input  InputType
        want   ExpectedType
        err    error
    }{
        // Setup
        for _, tt := range tests {
            t.Run(tt.name, func(t *testing.T) {
                // Execute test
                got, err := FunctionName(tt.input)
                
                // Assert
                if !reflect.DeepEqual(got, tt.want) {
                    t.Errorf("got = %v, want %v", got, tt.want)
                }
                if !errors.Is(err, tt.err) {
                    t.Errorf("got error %v, want %v", err, tt.err)
                }
            })
        }
    })
}
```

### Example: Domain Entity Test
```go
// internal/domain/task_test.go (if it existed)
func TestNewTask(t *testing.T) {
    tests := []struct {
        name    string
        params  domain.NewTaskParams
        want   *domain.Task
        err    error
    }{
        {
            name: "valid task with all fields",
            params: domain.NewTaskParams{
                ProjectID: "proj-123",
                Title:     "Test Task",
                Status:    domain.TaskStatusTodo,
                Priority: domain.PriorityHigh,
            },
            want: &domain.Task{
                ID:        "expected-uuid",
                ProjectID: "proj-123",
                Title:     "Test Task",
                Status:    domain.TaskStatusTodo,
                Priority: domain.PriorityHigh,
            },
            err: nil,
        },
        {
            name: "empty title",
            params: domain.NewTaskParams{
                ProjectID: "proj-123",
                Title:     "",
                Status:    domain.TaskStatusTodo,
                Priority: domain.PriorityMedium,
            },
            want:  nil,
            err:   domain.ErrTaskTitleEmpty,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            task, err := domain.NewTask(tt.params)
            
            if !errors.Is(err, tt.err) {
                t.Errorf("expected error %v, got %v", err, tt.err)
            }
            
            if !reflect.DeepEqual(task, tt.want) {
                t.Errorf("expected %v, got %v", task, tt.want)
            }
        })
    }
}
```

**Benefits**:
- **Clarity**: Test cases are self-documenting
- **coverage**: Easy to add new cases
- **maintainability**: Table format separates test data from test logic
- **subtests**: Run specific cases independently with `t.Run()`

---

## Interface Mocking

### Mock Pattern

```go
// internal/service/task_service_test.go
type mockTaskQuerier struct {
    createTaskFunc func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
    getTaskByIDFunc func(ctx context.Context, id string) (db.Task, error)
}

```

**Implementation**:
```go
func TestTaskService_Create(t *testing.T) {
    // Create mock
    mock := &mockTaskQuerier{
        createTaskFunc: func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
            return db.Task{
                ID:        "test-task-id",
                ProjectID: arg.ProjectID,
                Title:     arg.Title,
                CreatedAt: time.Now(),
                UpdatedAt: time.Now(),
            }, nil
        },
    }
    
    service := NewTaskService(mock, nil)
    
    // Call service
    task, err := service.Create(context.Background(), service.CreateTaskParams{
        ProjectID: "proj-123",
        Title:     "Test Task",
        Status:    domain.TaskStatusTodo,
    })
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if task.ID != "test-task-id" {
        t.Errorf("expected test-task-id, got %s", task.ID)
    }
}
```

### Mock Implementation Details
```go
type mockTaskQuerier struct {
    createTaskFunc          func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error)
    getTaskByIDFunc         func(ctx context.Context, id string) (db.Task, error)
    listTasksByProjectFunc  func(ctx context.Context, projectID string) ([]db.Task, error)
    updateTaskFunc          func(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error)
    hardDeleteTaskFunc      func(ctx context.Context, id string) error
}

func (m *mockTaskQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
    if m.createTaskFunc != nil {
        return m.createTaskFunc(ctx, arg)
    }
    return db.Task{}, nil
}

func (m *mockTaskQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
    if m.getTaskByIDFunc != nil {
        return m.getTaskByIDFunc(ctx, id)
    }
    return db.Task{}, nil
}

// ... implement other methods
```

**Key Points**:
- Mock implements the same interface as production code
- Function fields allow customizable behavior
- Default implementations return zero values
- Can be used across multiple tests

---

## Test Helpers

### Helper Pattern
```go
// internal/test/helpers.go
func assertNoError(t *testing.T, got error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got = %v, want %v", got, want)
    }
}

func assertError(t *testing.T, got error, want error) {
    t.Helper()
    if !errors.Is(got, want) {
        t.Errorf("got error %v, want %v", got, want)
    }
}
```

### Usage
```go
func TestTaskService_Create(t *testing.T) {
    task, err := service.Create(ctx, params)
    
    assertNoError(t, err)
    assertEqual(t, task.Title, "Test Task")
}
```

**Benefits**:
- **Reduces boilerplate**: Common assertions in one place
- **t.Helper()**: Marks function as helper for better error messages
- **Consistency**: Same patterns across all tests

---

## Golden Files

### Pattern
Golden files store expected output for comparison:

```
testdata/
├── task_create.golden
├── task_list.golden
└── task_update.golden
```

### Test Example
```go
func TestTaskListGolden(t *testing.T) {
    // Setup
    tasks := []domain.Task{...}
    
    // Generate output
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    encoder.Encode(tasks)
    
    // Read golden file
    golden := filepath.Join("testdata", "task_list.golden")
    want, err := os.ReadFile(golden)
    assertNoError(t, err)
    
    // Compare
    if !bytes.Equal(buf.Bytes(), want) {
        t.Errorf("output does not match golden file")
    }
}
```

### Updating Golden Files
```bash
# Update golden files when output changes intentionally
go test ./internal/service -update
```

---

## Benchmarks

### Benchmark Pattern
```go
func BenchmarkTaskServiceCreate(b *testing.B) {
    // Setup
    service := NewTaskService(mockRepo)
    
    // Reset timer (exclude setup)
    b.ResetTimer()
    
    // Run benchmark
    for i := 0; i < b.N; i++ {
        service.Create(context.Background(), service.CreateTaskParams{
            ProjectID: "proj-123",
            Title:     fmt.Sprintf("Task %d", i),
            Status:    domain.TaskStatusTodo,
            Priority: domain.PriorityMedium,
        })
    }
}
```

### Memory Benchmarks
```bash
# Run with memory profiling
go test -bench=. -benchmem ./internal/service
```

---

## Test Coverage

### Coverage Goals
- **80%+ general coverage** for most packages
- **100% critical paths** (domain validation, business logic)

 
```bash
# Run coverage
go test ./... -coverprofile=coverage.out

# View HTML report
go tool cover -html=coverage.out -o coverage.html
```

### Coverage by Package
| Package | Target Coverage | Notes |
|---------|----------------|------|
| `internal/domain` | 100% | Critical validation logic |
| `internal/service` | 85% | Business logic |
| `internal/converter` | 90% | Type conversions |
| `internal/db` | 70% | Repository (lower priority) |
| `cmd/dopa` | 75% | CLI commands |
| `internal/tui` | 80% | TUI components |

---

## Test Organization

### Unit Tests
- Located in same package as code: `*_test.go`
- Test individual functions/methods
- Fast execution

- No external dependencies (use mocks)

```go
// internal/domain/task_test.go
// internal/service/task_service_test.go
```

### Integration Tests
- Test multiple components together
- Use real database (in-memory SQLite)
- Located in `*_integration_test.go`
```go
// internal/db/integration_test.go
// internal/tui/integration_test.go
```

### Benchmarks
- Performance tests
- Located in `*_benchmark_test.go`
```go
// internal/service/project_service_benchmark_test.go
```

---

## Best Practices

### 1. Write Table-Driven Tests
✅ **Do**: Use table-driven tests for clarity
❌ **Don't**: Write individual test functions for each case
```go
// ✅ Good
func TestNewTask(t *testing.T) {
    tests := []struct {
        name    string
        params  domain.NewTaskParams
        err    error
    }{
        {"empty title", domain.NewTaskParams{Title: ""}, domain.ErrTaskTitleEmpty},
        {"valid task", domain.NewTaskParams{Title: "Test"}, nil},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := domain.NewTask(tt.params)
            assertError(t, err, tt.err)
        })
    }
}
```

### 2. Use Interface Mocking
✅ **Do**: Create mock implementations of interfaces
❌ **Don't**: Use concrete implementations in tests
```go
// ✅ Good: Mock allows isolation
mockRepo := &mockTaskQuerier{
    CreateFunc: func(ctx context.Context, arg CreateTaskParams) (Task, error) {
        return Task{ID: "test"}, nil
    },
}
```

### 3. Keep Tests Independent
✅ **Do**: Each test should be self-contained
❌ **Don't**: Tests should not depend on shared state

```go
// ✅ Good: Clean setup for each test
func TestTaskService_Create(t *testing.T) {
    mock := newMockTaskQuerier()
    service := NewTaskService(mock)
    
    // Test is isolated
    task, err := service.Create(ctx, params)
}
```

### 4. Use Subtests
✅ **Do**: Group related test cases with `t.Run()`
```go
// ✅ Good: Subtests for organization
t.Run("valid input", func(t *testing.T) {
    // Test valid input
})

t.Run("invalid input", func(t *testing.T) {
    // Test invalid input
})
```

### 5. Test Error Cases
✅ **Do**: Test both success and error paths
```go
// ✅ Good: Test all cases
func TestTaskService_Create(t *testing.T) {
    t.Run("success", func(t *testing.T) {
        // Test successful creation
    })
    
    t.Run("validation error", func(t *testing.T) {
        // Test domain validation failure
    })
    
    t.Run("repository error", func(t *testing.T) {
        // Test database error
    })
}
```

---

## Common Test Patterns

### Testing Domain Validation
```go
func TestNewTask_Validation(t *testing.T) {
    tests := []struct {
        name        string
        params      domain.NewTaskParams
        expectedErr error
    }{
        {
            name: "empty title",
            params: domain.NewTaskParams{
                ProjectID: "proj-123",
                Title:     "",
                Status:    domain.TaskStatusTodo,
            },
            expectedErr: domain.ErrTaskTitleEmpty,
        },
        {
            name: "invalid date range",
            params: domain.NewTaskParams{
                ProjectID: "proj-123",
                Title:     "Test",
                Status:    domain.TaskStatusTodo,
                StartDate: ptr(time.Now().Add(48 * time.Hour)),
                Deadline:  ptr(time.Now().Add(24 * time.Hour)), // Before start
            },
            expectedErr: domain.ErrTaskInvalidDateRange,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := domain.NewTask(tt.params)
            assertError(t, err, tt.expectedErr)
        })
    }
}
```

### Testing Service Methods
```go
func TestTaskService_Create(t *testing.T) {
    mock := &mockTaskQuerier{
        CreateTaskFunc: func(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
            return db.Task{
                ID:        "task-123",
                ProjectID: arg.ProjectID,
                Title:     arg.Title,
                CreatedAt: time.Now(),
            }, nil
        },
    }
    
    service := NewTaskService(mock)
    
    task, err := service.Create(context.Background(), CreateTaskParams{
        ProjectID: "proj-123",
        Title:     "Test Task",
        Status:    domain.TaskStatusTodo,
    })
    
    assertNoError(t, err)
    if task.ID != "task-123" {
        t.Errorf("expected task-123, got %s", task.ID)
    }
}
```

### Testing CLI Commands
```go
func TestTasksCreateCommand(t *testing.T) {
    // Setup mock service
    mockService := &MockTaskService{
        CreateFunc: func(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
            return &domain.Task{
                ID:        "task-123",
                ProjectID: params.ProjectID,
                Title:     params.Title,
            }, nil
        },
    }
    
    // Create command with mock
    cmd := newcobra.Command{}
    cmd.SetContext(&cobra.Context{Values: map[string]interface{}{
        "title":      "Test Task",
        "project-id": "proj-123",
    }})
    
    // Execute command
    err := runTasksCreate(cmd, []string{})
    
    // Assert
    assertNoError(t, err)
}
```

---

## Key Files

| File | Purpose |
|------|---------|
| `internal/domain/*_test.go` | Domain validation tests |
| `internal/service/*_test.go` | Service tests with mocks |
| `internal/converter/*_test.go` | Converter unit tests |
| `internal/db/*_test.go` | Repository integration tests |
| `cmd/dopa/*_test.go` | CLI command tests |
| `internal/tui/*_test.go` | TUI component tests |
| `testdata/*.golden` | Golden files for snapshot testing |
| `internal/test/helpers.go` | Test helper functions |

---

**Navigation**: [← CLI Layer](06-cli-layer.md) | [Back to Architecture →](README.md)
