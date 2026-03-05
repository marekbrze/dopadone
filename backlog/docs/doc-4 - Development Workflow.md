# Development Workflow Guide

## Overview

This guide covers the complete development workflow for ProjectDB, including database seeding, testing, and common development tasks.

## Quick Start

### Development Script

The `dev.sh` script provides shortcuts for common development tasks:

```bash
# Seed database with unique contextual tasks
./dev.sh seed

# Start TUI with seeded data
./dev.sh tui

# Run tests
./dev.sh test
```

## Database Seeding

### Why Contextual Tasks?

The seed script creates **unique, contextual tasks** based on project type, making the TUI more realistic and easier to test. Instead of generic "Task 1", "Task 2", "Task 3" for every project, each project gets tasks that make sense for its domain.

### Task Categories

The seed script uses pattern matching to create appropriate tasks:

| Project Pattern | Example Tasks |
|----------------|---------------|
| **Software Development** (`*API*`, `*Backend*`, `*Frontend*`, `*UI*`, `*CLI*`) | "Define requirements", "Write code", "Write tests" |
| **Learning** (`*AWS*`, `*Rust*`, `*Learn*`, `*Book*`, `*Exam*`) | "Read documentation", "Complete exercises", "Take practice tests" |
| **Home Improvement** (`*Garden*`, `*Renovation*`) | "Research materials", "Get quotes from contractors", "Set budget" |
| **Organization** (`*Basement*`, `*Organization*`) | "Declutter items", "Buy storage containers", "Organize by category" |
| **Business/Analytics** (`*E-commerce*`, `*Analytics*`, `*Dashboard*`, `*Pipeline*`) | "Gather requirements", "Design system architecture", "Implement MVP" |
| **Mobile Apps** (`*Habit*`, `*Tracker*`, `*Note*`, `*App*`, `*Mobile*`) | "Design UI mockups", "Implement core functionality", "Test on devices" |
| **Open Source** (`*ProjectDB*`, `*Open*`, `*Source*`) | "Write documentation", "Set up CI/CD", "Create release" |
| **Fitness** (`*Training*`, `*Running*`, `*Base*`, `*Speed*`, `*Long*`) | "Plan workout schedule", "Buy running gear", "Track progress" |
| **Meal Prep** (`*Meal*`, `*Prep*`) | "Plan weekly menu", "Create shopping list", "Prep ingredients" |
| **Default/Generic** | "Research options", "Create action plan", "Execute plan" |

### Seeded Data Structure

The seed script creates:

- **3 Areas**: Personal, Work, Side Projects
- **7 Subareas**: Health & Fitness, Learning, Home, Client A, Client B, Open Source, Mobile App
- **25 Projects** (with hierarchy):
  - Root projects linked to subareas
  - Child projects linked to parent projects
- **75 Tasks**: 3 unique tasks per project, contextually appropriate

### Running the Seed Script

```bash
# Method 1: Using dev.sh (recommended)
./dev.sh seed

# Method 2: Using Makefile
make seed

# Method 3: Direct script call
./scripts/seed-test-data.sh projectdb.db
```

### Customizing Seed Data

To add new task patterns, edit `scripts/seed-test-data.sh`:

```bash
case "$PROJ_NAME" in
    *YourPattern*)
        $CMD tasks create --title "Your custom task 1" --project-id "$PROJ_ID" > /dev/null 2>&1
        $CMD tasks create --title "Your custom task 2" --project-id "$PROJ_ID" > /dev/null 2>&1
        $CMD tasks create --title "Your custom task 3" --project-id "$PROJ_ID" > /dev/null 2>&1
        ;;
esac
```

## Testing Workflow

### Running Tests

```bash
# All tests
make test

# Or using dev.sh
./dev.sh test

# Specific package
go test ./internal/db/... -v

# With coverage
go test ./... -cover
```

### Test Database

Tests use isolated databases:
- `test-*.db` - Temporary test databases (auto-cleaned)
- `projectdb.db` - Main development database (used by seed script)

## TUI Development

### Starting the TUI

```bash
# Method 1: Using dev.sh
./dev.sh tui

# Method 2: Using go run
go run ./cmd/projectdb tui

# Method 3: Using make
make run
```

### TUI Development Tips

1. **Seed before testing**: Always run `./dev.sh seed` before testing the TUI to ensure realistic data
2. **Navigate the hierarchy**: Use arrow keys to explore parent/child projects
3. **Check task display**: Verify that different projects show different contextual tasks
4. **Test area switching**: Press `[` and `]` to switch between areas and verify data loads correctly

## Common Development Tasks

### Adding a New Entity

1. **Create migration**:
   ```bash
   # Create new migration file
   goose -dir migrations sqlite3 projectdb.db create add_entity_name sql
   ```

2. **Write schema**:
   ```sql
   -- +goose Up
   CREATE TABLE entity_name (
       id TEXT PRIMARY KEY,
       -- ... fields
   );
   
   -- +goose Down
   DROP TABLE entity_name;
   ```

3. **Generate sqlc code**:
   ```bash
   make sqlc-generate
   ```

4. **Implement domain logic** in `internal/domain/`
5. **Add CLI commands** in `cmd/projectdb/`
6. **Write tests** in corresponding test files

### Modifying Existing Entities

1. **Create migration** for schema changes
2. **Update sqlc queries** in `internal/db/*.sql`
3. **Regenerate code**: `make sqlc-generate`
4. **Update domain models** if needed
5. **Update CLI commands** if needed
6. **Update tests**

### Debugging Database Issues

```bash
# Check database state
sqlite3 projectdb.db ".schema"
sqlite3 projectdb.db "SELECT COUNT(*) FROM tasks;"

# Check task distribution
sqlite3 projectdb.db "SELECT p.name, COUNT(t.id) FROM projects p LEFT JOIN tasks t ON p.id = t.project_id GROUP BY p.id;"

# Verify task-project links
sqlite3 projectdb.db "SELECT t.title, p.name FROM tasks t JOIN projects p ON t.project_id = p.id LIMIT 10;"
```

## Development Environment Setup

### Prerequisites

- **Go 1.21+**: `go version`
- **SQLite3**: `sqlite3 --version`
- **goose**: `go install github.com/pressly/goose/v3/cmd/goose@latest`
- **sqlc**: Download from https://sqlc.dev/

### Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd projectdb

# Install dependencies
make install-deps

# Build binary
make build

# Initialize database
make migrate-up

# Seed database
./dev.sh seed

# Run tests
make test
```

## Code Organization

```
projectdb/
├── cmd/projectdb/          # CLI commands
│   ├── main.go            # Entry point
│   ├── areas.go           # Area CRUD commands
│   ├── subareas.go        # Subarea CRUD commands
│   ├── projects.go        # Project CRUD commands
│   ├── tasks.go           # Task CRUD commands
│   └── tui.go             # TUI command
├── internal/
│   ├── db/                # Generated database code
│   │   ├── querier.go     # Query interface
│   │   ├── models.go      # Database models
│   │   └── *.sql.go       # Generated query implementations
│   ├── domain/            # Business logic
│   │   ├── area.go        # Area domain logic
│   │   ├── subarea.go     # Subarea domain logic
│   │   ├── project.go     # Project domain logic
│   │   └── task.go        # Task domain logic
│   ├── cli/               # CLI utilities
│   │   ├── output/        # Formatters (table, JSON, YAML)
│   │   └── filter/        # Query parser & evaluator
│   └── tui/               # Terminal UI
│       ├── app.go         # Main TUI model
│       ├── commands.go    # Database commands
│       ├── tree/          # Tree rendering
│       └── views/         # UI components
├── migrations/            # Database migrations
├── scripts/               # Development scripts
│   └── seed-test-data.sh  # Database seeding
├── Makefile              # Build automation
├── dev.sh                # Development helper
└── README.md             # User documentation
```

## Best Practices

### Database Operations

1. **Always use transactions** for multi-step operations
2. **Use soft delete** by default (set `deleted_at` timestamp)
3. **Test migrations** before applying to production
4. **Back up database** before destructive operations

### Testing

1. **Use isolated databases** for tests (auto-cleanup)
2. **Test domain logic** separately from database
3. **Use table-driven tests** for comprehensive coverage
4. **Test error paths** not just happy path

### TUI Development

1. **Seed database** before testing TUI
2. **Test navigation** with keyboard shortcuts
3. **Verify data loading** on area/project switches
4. **Check state persistence** when switching contexts

### Code Quality

1. **Run linter** before committing: `make lint`
2. **Format code**: `gofmt -w .`
3. **Update documentation** when adding features
4. **Write clear commit messages**

## Troubleshooting

### Common Issues

#### "Tasks not showing in TUI"

**Cause**: Database not seeded or seed script created tasks with invalid project IDs

**Solution**:
```bash
# Reseed database
rm projectdb.db
./dev.sh seed

# Verify tasks exist
sqlite3 projectdb.db "SELECT COUNT(*) FROM tasks;"

# Check task-project links
sqlite3 projectdb.db "SELECT t.title, p.name FROM tasks t JOIN projects p ON t.project_id = p.id LIMIT 5;"
```

#### "Foreign key constraint failed"

**Cause**: Trying to create entity with invalid parent reference

**Solution**: Verify parent ID exists before creating child entity

#### "Migration failed"

**Cause**: Database state doesn't match expected migration state

**Solution**:
```bash
# Reset database (WARNING: destroys data)
make migrate-reset

# Or manually check state
goose -dir migrations sqlite3 projectdb.db status
```

## Additional Resources

- [Data Layer Architecture](doc-1%20-%20Data-Layer-Architecture.md)
- [CLI CRUD Operations Guide](doc-2%20-%20CLI-CRUD-Operations-Guide.md)
- [TUI Architecture](doc-3%20-%20TUI-Architecture.md)
- [Main README](../../README.md)

## Contributing

When contributing changes:

1. **Update documentation** if adding/modifying features
2. **Add tests** for new functionality
3. **Run full test suite**: `make test`
4. **Seed and test TUI**: `./dev.sh seed && ./dev.sh tui`
5. **Update README** if user-facing changes
