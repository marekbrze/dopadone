# Dopadone - CLI Project Management for Developers

[![CI](https://github.com/marekbrze/dopadone/workflows/CI/badge.svg)](https://github.com/marekbrze/dopadone/actions/workflows/ci.yml)
[![Release](https://github.com/marekbrze/dopadone/workflows/Release/badge.svg)](https://github.com/marekbrze/dopadone/actions/workflows/release.yml)

**Organize your projects, tasks, and workflows from the command line.**

Dopadone is a lightweight, SQLite-based project management tool designed for developers who prefer staying in the terminal. It provides a clean hierarchical structure for organizing work without the overhead of bloated project management software.

**The problem:** You need to track projects, sub-projects, and tasks across multiple areas of your life (work, personal, side projects), but existing tools are either too heavy, require a web browser, or don't fit a developer's workflow.

**The solution:** A CLI-first project database that stores everything locally in SQLite, supports hierarchical organization, and integrates naturally with your existing terminal-based workflow.

## Quick Start

Get started in 30 seconds with a complete workflow using **Areas** (top-level containers):

```bash
# Install (one of these methods)
go install github.com/marekbrze/dopadone/cmd/dopa@latest
# OR: download binary from releases
# OR: make build

# Initialize database with migrations
dopa migrate up

# CREATE: Add a new area
dopa area create --name "Work" --color "#3B82F6"

# READ: List all areas
dopa area list

# UPDATE: Modify an area
dopa area update <area-id> --name "Professional Work"

# DELETE: Remove an area (soft delete by default)
dopa area delete <area-id>
```

> **Note:** All `get`, `update`, and `delete` commands require the entity's **UUID** (e.g., `b3e76f50-3640-4dfa-be85-c5401dd18555`), not its name. Use `list` commands to find UUIDs.

Once you understand areas, the same CRUD pattern applies to subareas, projects, and tasks.

## Installation

### Option 1: Quick Install (Recommended)

The fastest way to install Dopadone on Linux or macOS:

```bash
curl -sSL https://raw.githubusercontent.com/marekbrze/dopadone/main/scripts/install.sh | sh
```

This script will:
- Detect your operating system and architecture
- Download the latest release from GitHub
- Install the binary to `/usr/local/bin/dopa` (requires sudo)
- Verify the installation

#### Advanced Options

```bash
# Dry-run mode (test without installing)
./install.sh --dry-run

# Skip installation verification
./install.sh --no-verify

# Auto-confirm upgrade (no prompts)
./install.sh --yes

# Custom installation directory
INSTALL_DIR=$HOME/bin ./install.sh

# Show help and all options
./install.sh --help
```

### Option 2: Download Pre-built Binaries

Download the latest release for your platform:

| Platform | Architecture | Binary |
|----------|--------------|--------|
| Linux | amd64 | `dopa-linux-amd64` |
| macOS | amd64 | `dopa-darwin-amd64` |
| macOS | arm64 (M1/M2) | `dopa-darwin-arm64` |
| Windows | amd64 | `dopa-windows-amd64.exe` |

After downloading, make it executable and move to your PATH:

```bash
chmod +x dopa-*
sudo mv dopa-* /usr/local/bin/dopa
```

### Option 2: Install with Go

```bash
go install github.com/marekbrze/dopadone/cmd/dopa@latest
```

### Option 3: Build from Source

```bash
git clone https://github.com/marekbrze/dopadone.git
cd dopa
make build
# Binary will be at bin/dopa
```

## Usage

### Data Hierarchy

```
Areas (top-level categories, e.g., "Work", "Personal")
└── Subareas (subcategories, e.g., "Backend", "Frontend" under Work)
    └── Projects (e.g., "Website Redesign" linked to a subarea)
        └── Sub-projects (nested projects, e.g., "API Design" under Website Redesign)
            └── Tasks (individual work items, e.g., "Write tests")
```

### Areas

Top-level containers for organizing your work:

```bash
# Create
dopa area create --name "Personal" --color "#10B981"
dopa area create --name "Work" --color "#3B82F6"

# List
dopa area list
dopa area list --format json
dopa area list --filter 'name=Work'

# Get details
dopa area get <area-id>

# Update
dopa area update <area-id> --name "Professional"
dopa area update <area-id> --color "#8B5CF6"

# Delete (soft delete by default, recoverable)
dopa area delete <area-id>
dopa area delete <area-id> --permanent  # Hard delete
```

### Subareas

Subdivisions within areas:

```bash
# Create (requires parent area)
dopa subarea create --name "Backend" --area-id <area-id>
dopa subarea create --name "Frontend" --area-id <area-id> --color "#EC4899"

# List
dopa subarea list
dopa subarea list --filter 'area_id=<area-id>'

# Get, Update, Delete (same pattern as areas)
dopa subarea get <subarea-id>
dopa subarea update <subarea-id> --name "API Development"
dopa subarea delete <subarea-id>
```

### Projects

Projects can be root-level (linked to subarea) or nested (linked to parent project):

```bash
# Create root project (linked to subarea)
dopa project create --name "Website Redesign" --subarea-id <subarea-id>

# Create nested sub-project (linked to parent project)
dopa project create --name "API Integration" --parent-id <project-id>

# Create with all options
dopa project create --name "Q4 Campaign" --subarea-id <subarea-id> \
  --status active --priority high --progress 25 \
  --start-date 2024-10-01 --deadline 2024-12-31 \
  --goal "Launch by year end" --description "Marketing campaign"

# List with filters
dopa project list
dopa project list --status active
dopa project list --priority high
dopa project list --subarea-id <subarea-id>
dopa project list --parent-id <project-id>
dopa project list --filter 'status=active AND priority>=high'

# Update
dopa project update <project-id> --status completed --progress 100
dopa project update <project-id> --priority urgent --deadline 2024-12-31

# Delete
dopa project delete <project-id>
```

**Project Status Options:** `active`, `completed`, `on_hold`, `archived`

**Project Priority Options:** `low`, `medium`, `high`, `urgent`

### Tasks

Individual work items within projects:

```bash
# Create
dopa task create --project-id <project-id> --title "Write documentation"
dopa task create --project-id <project-id> --title "API Integration" \
  --description "Integrate with external API" \
  --status in_progress --priority high \
  --start-date 2024-01-15 --deadline 2024-01-31 \
  --context "backend" --duration 60 --next

# List
dopa task list
dopa task list --project-id <project-id>
dopa task list --status in_progress
dopa task list --next  # Show only priority tasks

# Show next/priority tasks
dopa task next

# Update
dopa task update <task-id> --status done
dopa task update <task-id> --next     # Mark as priority
dopa task update <task-id> --no-next  # Remove priority flag

# Delete
dopa task delete <task-id>
```

**Task Status Options:** `todo`, `in_progress`, `waiting`, `done`

**Task Priority Options:** `critical`, `high`, `medium`, `low`

**Duration Options:** `5`, `15`, `30`, `60`, `120`, `240`, `480` (minutes)

## Output Formats

Control output format with `--format` or `-o`:

```bash
# Table format (default)
dopa area list

# JSON format
dopa area list --format json
dopa area list -o json

# YAML format
dopa area list --format yaml
```

### Table Output Example

```
 ID                                    NAME     COLOR    CREATED    
 c37fd550-dee9-4966-9173-eff71dbebc70  Work     #3B82F6  2024-03-03 
```

### JSON Output Example

```json
{
  "id": "c37fd550-dee9-4966-9173-eff71dbebc70",
  "name": "Work",
  "color": "#3B82F6",
  "created_at": "2024-03-03T12:08:47Z",
  "updated_at": "2024-03-03T12:08:47Z"
}
```

## Global Flags

These flags work with all commands:

| Flag | Description | Default |
|------|-------------|---------|
| `--db` | Path to SQLite database file | `./dopa.db` |
| `-o, --output` | Output format (`table`, `json`) | `table` |
| `--format` | Extended output format (`table`, `json`, `yaml`) | `table` |

Examples:

```bash
# Use a custom database location
dopa --db /path/to/my.db area list

# Output as JSON for scripting
dopa -o json project list --status active

# Combine with other tools
dopa --format json task list | jq '.[] | select(.priority=="high")'
```

## Terminal User Interface (TUI)

Dopadone includes an interactive terminal UI for visual browsing and management:

```bash
# Launch the TUI
dopa tui
```

The TUI provides:
- **Area tabs** at the top for quick switching between work areas
- **3-column browser** (Subareas | Projects | Tasks) with keyboard navigation
- **Focus-aware borders** showing which column is active
- **Terminal resize support** - adapts to window size changes

### TUI Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `h`, `←` | Move focus left (wraps right-to-left) |
| `l`, `→` | Move focus right (wraps left-to-right) |
| `Tab` | Cycle focus through columns |
| `j`, `↓` | Navigate down in current column (wraps to top) |
| `k`, `↑` | Navigate up in current column (wraps to bottom) |
| `[` | Switch to previous area tab (wraps to last) |
| `]` | Switch to next area tab (wraps to first) |
| `Enter`, `Space` | Toggle expand/collapse for project tree nodes |
| `a` | Open quick-add modal for creating items |
| `?` | Show help modal with all keyboard shortcuts |
| `q`, `Ctrl+C` | Exit TUI |

**Quick-Add Modal:**
- Press `a` to open context-aware modal
- Creates subarea, project, or task based on focused column
- Shows parent context (e.g., "New Project in: Work Tasks")
- Type title and press `Enter` to create
- Press `Escape` to cancel

**Help Modal:**
- Press `?` to open help modal with all keyboard shortcuts
- Shortcuts are grouped by category (Navigation, Actions, General)
- Press `?`, `Escape`, or `q` to close

**Note:** The TUI is fully implemented (Task-14 complete). All core features including quick-add modal, help system, and state persistence are production-ready.

For comprehensive TUI documentation including architecture, components, and implementation details, see [docs/TUI.md](docs/TUI.md).

## Database Migrations

Dopadone manages its own schema via embedded migrations:

```bash
# Apply migrations (run once after install)
dopa migrate up

# Check migration status
dopa migrate status

# Rollback last migration
dopa migrate down

# Reset database (warning: destroys data)
dopa migrate reset
```

---

## Development

### Prerequisites

- Go 1.21+
- SQLite3
- [goose](https://github.com/pressly/goose) (for migrations)
- [sqlc](https://sqlc.dev/) (for query generation)

### Build Commands

```bash
make build          # Build binary to bin/dopa
make build-all      # Cross-compile for all platforms
make dist           # Build + create distribution archives
make clean          # Remove build artifacts
```

### Development Workflow

```bash
make run            # Build and run
make dev            # Run with go run (faster iteration)
make test           # Run all tests
make lint           # Run linter
make seed           # Seed database with test data
```

#### Quick Development Script

Use the `scripts/dev.sh` script for common development tasks:

```bash
# Seed database with unique contextual tasks
./scripts/dev.sh seed

# Start TUI with seeded data
./scripts/dev.sh tui

# Run tests
./scripts/dev.sh test
```

The seed script creates **unique, contextual tasks** for each project type:
- Software projects: "Define requirements", "Write code", "Write tests"
- Learning projects: "Read documentation", "Complete exercises", "Take practice tests"
- Home projects: "Research materials", "Get quotes", "Set budget"
- Mobile apps: "Design UI mockups", "Implement core functionality", "Test on devices"

See [backlog/docs/doc-4 - Development Workflow.md](backlog/docs/doc-4%20-%20Development%20Workflow.md) for detailed development guidelines.

### Database Development

```bash
make migrate-up     # Apply all migrations
make migrate-down   # Rollback last migration
make migrate-status # Check migration status
make migrate-reset  # Reset database (down + up)
make sqlc-generate  # Generate sqlc code after query changes
```

### Testing

```bash
# Run all tests
go test ./... -v

# Run specific test suites
go test ./internal/db/... -v -run "TestCompleteHierarchy|TestSoftDeleteCascade"
```

### Schema Reference

| Table | Description |
|-------|-------------|
| `areas` | Top-level organizational categories |
| `subareas` | Subcategories within areas |
| `projects` | Projects with nesting via `parent_id` or linking to `subarea_id` |
| `tasks` | Individual work items linked to projects |

### Key Constraints

1. **Project Hierarchy**: A project must have either `parent_id` OR `subarea_id` (not both, not neither)
2. **Soft Delete**: Deletes are soft by default; child entities remain when parent is soft-deleted
3. **Foreign Key Cascade**: Hard deletes cascade at the database level
4. **Transaction Support**: Multi-entity operations use transactions for atomicity and consistency

### Transaction Support

Dopadone uses database transactions for operations that modify multiple entities:

- **HardDelete operations**: Cascade deletes wrapped in transactions (tasks→projects→subareas→area)
- **Batch operations**: Sort order updates and bulk changes are atomic
- **Serializable isolation**: Strongest consistency guarantees
- **Automatic rollback**: On error or panic, all changes are rolled back

For detailed documentation on transaction usage patterns, integration examples, and best practices, see [docs/TRANSACTIONS.md](docs/TRANSACTIONS.md).

### Tech Stack

- **Database**: SQLite (via modernc.org/sqlite - pure Go)
- **Migrations**: goose v3
- **Query Generation**: sqlc (type-safe SQL)
- **UUIDs**: google/uuid

### Indexes

| Index | Purpose |
|-------|---------|
| `idx_projects_deadline` | Fast deadline-based queries |
| `idx_projects_status_priority` | Composite index for status+priority filtering |
| `idx_projects_parent_id` | Efficient child project lookups |
| `idx_projects_subarea_id` | Fast subarea project listings |
| `idx_subareas_area_id` | Efficient area-based subarea queries |
