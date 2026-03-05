# Transaction Support

This document describes the transaction support system in ProjectDB, which ensures data consistency for multi-entity operations.

## Overview

ProjectDB uses database transactions to ensure **atomic** operations when modifying multiple entities. This guarantees that complex operations (like cascade deletes or batch updates) either complete entirely or roll back completely, leaving the database in a consistent state.

## Architecture

### TransactionManager

The `TransactionManager` in `internal/db/transaction.go` provides a callback-based API for managing transactions:

```go
type TransactionManager struct {
    db *sql.DB
}

func NewTransactionManager(db *sql.DB) *TransactionManager

func (tm *TransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx Querier) error) error
```

### Key Features

- **Callback Pattern**: Automatic commit/rollback eliminates manual error handling
- **Serializable Isolation**: Strongest consistency guarantees
- **Context-Based**: Transaction state passed via context for thread-safety
- **Panic Recovery**: Automatic rollback on panic

## Usage Patterns

### Basic Transaction

```go
func (s *MyService) UpdateMultiple(ctx context.Context) error {
    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        if err := tx.UpdateEntityA(ctx, ...); err != nil {
            return err // Will trigger rollback
        }
        if err := tx.UpdateEntityB(ctx, ...); err != nil {
            return err // Will trigger rollback
        }
        return nil // Will commit
    })
}
```

### Transaction-Aware Querier

Use `GetQuerierFromContext` to ensure code works both inside and outside transactions:

```go
func (s *MyService) DoWork(ctx context.Context) error {
    querier := db.GetQuerierFromContext(ctx, s.repo)
    return querier.SomeOperation(ctx, ...)
}
```

## Service Integration

### Constructor Pattern

Services accept an optional `*db.TransactionManager`:

```go
type AreaService struct {
    repo db.Querier
    tm   *db.TransactionManager
}

func NewAreaService(repo db.Querier, tm *db.TransactionManager) *AreaService {
    return &AreaService{repo: repo, tm: tm}
}
```

### When to Use Transactions

| Operation | Transaction Required | Reason |
|-----------|---------------------|--------|
| Single entity CRUD | No | Single SQL statement |
| HardDelete with cascade | Yes | Multiple dependent deletes |
| Batch updates | Yes | Multiple related updates |
| ReorderAll operations | Yes | Multiple sort order updates |

## Implementation Details

### Cascade Delete Order

For `HardDelete` operations, entities must be deleted in the correct order to respect foreign key constraints:

```
Area HardDelete:    tasks → projects → subareas → area
Subarea HardDelete: projects → tasks → subarea
Project HardDelete: tasks → nested_projects → project
```

### Isolation Level

We use `sql.LevelSerializable` which provides:
- **Read Committed**: No dirty reads
- **Repeatable Read**: Consistent reads within transaction
- **Serializable**: Full isolation (no phantom reads)

For SQLite, serializable isolation is efficient due to database-level locking.

## Error Handling

### Automatic Rollback

Transactions automatically roll back when:
1. The callback function returns an error
2. The callback function panics
3. The commit fails

### Error Wrapping

All transaction errors are wrapped with context:

```go
// Begin transaction failure
"begin transaction: <error>"

// Rollback failure
"rollback failed after error: <rollback_error> (original: <original_error>)"

// Commit failure
"commit transaction: <error>"
```

## Testing

### Unit Tests

Mock the `TransactionManager` to verify rollback behavior:

```go
type MockTransactionManager struct {
    WithTransactionFunc func(ctx context.Context, fn func(ctx context.Context, tx db.Querier) error) error
}

func (m *MockTransactionManager) WithTransaction(ctx context.Context, fn func(ctx context.Context, tx db.Querier) error) error {
    return m.WithTransactionFunc(ctx, fn)
}
```

### Integration Tests

Use a real SQLite database to test actual rollback behavior:

```go
func TestHardDeleteRollback(t *testing.T) {
    db := setupTestDB(t)
    tm := db.NewTransactionManager(db)
    
    // Create test data
    area := createTestArea(t, db)
    subarea := createTestSubarea(t, db, area.ID)
    project := createTestProject(t, db, subarea.ID)
    
    // Simulate error mid-delete
    err := tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
        tx.DeleteTasksByProject(ctx, project.ID)
        return errors.New("simulated error") // Triggers rollback
    })
    
    assert.Error(t, err)
    // Verify all entities still exist
    assertAreaExists(t, db, area.ID)
    assertSubareaExists(t, db, subarea.ID)
    assertProjectExists(t, db, project.ID)
}
```

## Examples

### AreaService.HardDelete

```go
func (s *AreaService) HardDelete(ctx context.Context, id uuid.UUID) error {
    if s.tm == nil {
        return s.hardDeleteWithoutTransaction(ctx, id)
    }
    
    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx Querier) error {
        // Delete in correct order: tasks → projects → subareas → area
        subareas, _ := tx.ListSubareasByArea(ctx, id)
        for _, subarea := range subareas {
            projects, _ := tx.ListProjectsBySubarea(ctx, subarea.ID)
            for _, project := range projects {
                if err := tx.DeleteTasksByProject(ctx, project.ID); err != nil {
                    return err
                }
            }
            if err := tx.DeleteProjectsBySubarea(ctx, subarea.ID); err != nil {
                return err
            }
        }
        if err := tx.DeleteSubareasByArea(ctx, id); err != nil {
            return err
        }
        return tx.HardDeleteArea(ctx, id)
    })
}
```

### AreaService.ReorderAll

```go
func (s *AreaService) ReorderAll(ctx context.Context, ids []uuid.UUID) error {
    if s.tm == nil {
        return s.reorderAllWithoutTransaction(ctx, ids)
    }
    
    return s.tm.WithTransaction(ctx, func(ctx context.Context, tx Querier) error {
        for i, id := range ids {
            if err := tx.UpdateAreaSortOrder(ctx, db.UpdateAreaSortOrderParams{
                ID:        id,
                SortOrder: int64(i),
            }); err != nil {
                return err
            }
        }
        return nil
    })
}
```

## Best Practices

1. **Keep Transactions Small**: Only wrap operations that must be atomic
2. **Don't Nest Transactions**: The callback pattern doesn't support savepoints
3. **Handle Errors Explicitly**: Return errors from the callback to trigger rollback
4. **Use Context**: Pass the context through the entire call chain
5. **Test Rollback Scenarios**: Verify that partial failures leave data consistent

## Migration Guide

To add transaction support to an existing service:

1. Add `tm *db.TransactionManager` field to service struct
2. Update constructor to accept optional `*db.TransactionManager`
3. Wrap multi-entity operations with `tm.WithTransaction`
4. Use `GetQuerierFromContext` for transaction-aware queries
5. Update service instantiation in `main.go`

Example:

```go
// Before
func NewAreaService(repo db.Querier) *AreaService {
    return &AreaService{repo: repo}
}

// After
func NewAreaService(repo db.Querier, tm *db.TransactionManager) *AreaService {
    return &AreaService{repo: repo, tm: tm}
}
```

## Performance Considerations

- **SQLite Locking**: SQLite uses database-level locks, so transactions serialize automatically
- **Lock Timeout**: Handled by the SQLite driver
- **Minimal Overhead**: Transaction overhead is negligible for SQLite workloads
- **Batch Operations**: Single transaction is more efficient than multiple individual transactions

## Troubleshooting

### "database is locked" errors

This indicates concurrent transaction conflicts. Solutions:
1. Reduce transaction scope
2. Increase SQLite busy timeout
3. Avoid long-running transactions

### "FOREIGN KEY constraint failed"

This indicates incorrect delete order. Ensure cascade deletes follow the dependency order:
- Tasks before projects
- Projects before subareas
- Subareas before areas

### Transaction not rolling back

Ensure you're:
1. Returning errors from the callback (not just logging)
2. Using the context passed to the callback
3. Not catching and suppressing panics
