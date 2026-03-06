---
id: doc-2
title: CLI CRUD Operations Guide
type: user-guide
created_date: '2026-03-03'
---

# CLI CRUD Operations Guide

## Overview

The `dopa` CLI provides a complete command-line interface for managing projects, areas, and subareas in an ADHD-friendly project management system. Built with Cobra, it offers intuitive CRUD operations with colored table output, JSON/YAML support, and advanced filtering capabilities.

## Installation

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Building from Source

```bash
# Clone the repository
git clone <repository-url>
cd <repository-name>

# Install dependencies
go mod download

# Build the binary
go build -o dopa cmd/dopa/main.go

# Or run directly
go run cmd/dopa/main.go [command]
```

### Binary Location

After building, the `dopa` binary will be in your project root or you can install it to your PATH:

```bash
go install cmd/dopa/main.go
```

## Quick Start

```bash
# Initialize database (run migrations)
./dopa --help

# Create an area
./dopa areas create --name "Work" --color "#FF5733"

# Create a subarea under that area
./dopa subareas create --name "Backend" --area-id "area-123" --color "#3498DB"

# Create a project under the subarea
./dopa projects create --name "API Redesign" --subarea-id "subarea-456" --priority high

# List all projects
./dopa projects list

# List projects with filters
./dopa projects list --status active --priority high --format json
```

## Global Flags

These flags are available for all commands:

| Flag | Default | Description |
|------|---------|-------------|
| `--db` | `./dopa.db` | Path to SQLite database file |
| `-o, --output` | `table` | Output format: `table`, `json`, or `yaml` |
| `-h, --help` | - | Help for any command |

## Output Formats

The CLI supports three output formats:

### Table Format (Default)

```bash
./dopa areas list
```

Output:
```
ID          NAME         COLOR      CREATED AT           UPDATED AT
area-abc123 Engineering  #4A90E2    2024-01-15 10:30:00  2024-01-15 10:30:00
area-def456 Marketing    #E74C3C    2024-01-15 11:00:00  2024-01-15 11:00:00
```

### JSON Format

```bash
./dopa areas list --json
# or
./dopa areas list --format json
```

Output:
```json
[
  {
    "id": "area-abc123",
    "name": "Engineering",
    "color": "#4A90E2",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### YAML Format

```bash
./dopa areas list --format yaml
```

Output:
```yaml
- id: area-abc123
  name: Engineering
  color: "#4A90E2"
  created_at: 2024-01-15T10:30:00Z
  updated_at: 2024-01-15T10:30:00Z
```

---

# Areas CRUD Operations

Areas are top-level containers for organizing your work. Each area can contain multiple subareas.

## Create Area

Create a new top-level area.

### Syntax

```bash
./dopa areas create --name <name> [--color <hex-color>]
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--name` | Yes | - | Area name |
| `--color` | No | - | Color in hex format (e.g., `#FF5733`) |

### Examples

```bash
# Create area with required name
./dopa areas create --name "Engineering"

# Create area with color
./dopa areas create --name "Marketing" --color "#FF5733"

# Create personal area
./dopa areas create --name "Personal Projects" --color "#3498DB"
```

### Output

```json
{
  "id": "area-abc123",
  "name": "Engineering",
  "color": "#4A90E2",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## List Areas

List all areas in the database.

### Syntax

```bash
./dopa areas list [--json] [--format table|json|yaml]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` | false | Output as JSON (shorthand for `--format json`) |
| `--format` | table | Output format: `table`, `json`, or `yaml` |

### Examples

```bash
# List all areas (table format)
./dopa areas list

# List as JSON
./dopa areas list --json

# List as YAML
./dopa areas list --format yaml
```

## Get Area

Display details of a single area by ID.

### Syntax

```bash
./dopa areas get <area-id>
```

### Arguments

- `area-id`: The ID of the area to retrieve

### Examples

```bash
./dopa areas get area-abc123
```

### Output

```json
{
  "id": "area-abc123",
  "name": "Engineering",
  "color": "#4A90E2",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Update Area

Update an existing area's name and/or color.

### Syntax

```bash
./dopa areas update <area-id> [--name <name>] [--color <hex-color>]
```

### Arguments

- `area-id`: The ID of the area to update

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | No* | New area name |
| `--color` | No* | New color in hex format |

*At least one flag is required

### Examples

```bash
# Update name
./dopa areas update area-abc123 --name "Engineering Team"

# Update color
./dopa areas update area-abc123 --color "#9B59B6"

# Update both
./dopa areas update area-abc123 --name "Engineering" --color "#3498DB"
```

## Delete Area

Delete an area. By default, performs a soft delete (marks as deleted but retains in database). Use `--permanent` for hard delete.

### Syntax

```bash
./dopa areas delete <area-id> [--permanent]
```

### Arguments

- `area-id`: The ID of the area to delete

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--permanent` | false | Permanently delete from database (hard delete) |

### Examples

```bash
# Soft delete (default)
./dopa areas delete area-abc123

# Permanent delete
./dopa areas delete area-abc123 --permanent
```

**Warning**: Permanent delete cannot be undone!

---

# Subareas CRUD Operations

Subareas are second-level containers that belong to areas. They help further organize projects within an area.

## Create Subarea

Create a new subarea under an area.

### Syntax

```bash
./dopa subareas create --name <name> --area-id <area-id> [--color <hex-color>]
```

### Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--name` | Yes | - | Subarea name |
| `--area-id` | Yes | - | Parent area ID |
| `--color` | No | - | Color in hex format (inherits from area if not set) |

### Examples

```bash
# Create subarea with required fields
./dopa subareas create --name "Backend" --area-id "area-abc123"

# Create subarea with custom color
./dopa subareas create --name "Frontend" --area-id "area-abc123" --color "#3498DB"

# Create multiple subareas under same area
./dopa subareas create --name "DevOps" --area-id "area-abc123" --color "#2ECC71"
```

### Output

```json
{
  "id": "subarea-xyz789",
  "name": "Backend",
  "area_id": "area-abc123",
  "color": "#3498DB",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

## List Subareas

List all subareas or filter by area.

### Syntax

```bash
./dopa subareas list [--area-id <area-id>] [--json] [--format table|json|yaml]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--area-id` | - | Filter by parent area ID |
| `--json` | false | Output as JSON |
| `--format` | table | Output format: `table`, `json`, or `yaml` |

### Examples

```bash
# List all subareas
./dopa subareas list

# List subareas under specific area
./dopa subareas list --area-id "area-abc123"

# List as JSON
./dopa subareas list --json

# List filtered subareas as YAML
./dopa subareas list --area-id "area-abc123" --format yaml
```

## Get Subarea

Display details of a single subarea by ID.

### Syntax

```bash
./dopa subareas get <subarea-id>
```

### Arguments

- `subarea-id`: The ID of the subarea to retrieve

### Examples

```bash
./dopa subareas get subarea-xyz789
```

## Update Subarea

Update an existing subarea's name and/or color.

### Syntax

```bash
./dopa subareas update <subarea-id> [--name <name>] [--color <hex-color>]
```

### Arguments

- `subarea-id`: The ID of the subarea to update

### Flags

| Flag | Required | Description |
|------|----------|-------------|
| `--name` | No* | New subarea name |
| `--color` | No* | New color in hex format |

*At least one flag is required

### Examples

```bash
# Update name
./dopa subareas update subarea-xyz789 --name "Backend Services"

# Update color
./dopa subareas update subarea-xyz789 --color "#E74C3C"

# Update both
./dopa subareas update subarea-xyz789 --name "Backend" --color "#3498DB"
```

## Delete Subarea

Delete a subarea (soft delete by default).

### Syntax

```bash
./dopa subareas delete <subarea-id> [--permanent]
```

### Arguments

- `subarea-id`: The ID of the subarea to delete

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--permanent` | false | Permanently delete from database |

### Examples

```bash
# Soft delete
./dopa subareas delete subarea-xyz789

# Permanent delete
./dopa subareas delete subarea-xyz789 --permanent
```

---

# Projects CRUD Operations

Projects are goal-oriented task containers that can be nested recursively. They can belong to either a subarea (root project) or another project (nested project).

## Create Project

Create a new project. Must specify either `--subarea-id` OR `--parent-id` (not both).

### Syntax

```bash
./dopa projects create --name <name> (--subarea-id <id> | --parent-id <id>) [options]
```

### Required Flags

| Flag | Description | Notes |
|------|-------------|-------|
| `--name` | Project name | Required |
| `--subarea-id` | Parent subarea ID | Required if `--parent-id` not set |
| `--parent-id` | Parent project ID | Required if `--subarea-id` not set |

### Optional Flags

| Flag | Default | Description | Valid Values |
|------|---------|-------------|--------------|
| `--status` | active | Project status | `active`, `completed`, `on_hold`, `archived` |
| `--priority` | medium | Project priority | `low`, `medium`, `high`, `urgent` |
| `--progress` | 0 | Completion percentage | 0-100 |
| `--deadline` | - | Deadline date | Format: `YYYY-MM-DD` |
| `--start-date` | - | Start date | Format: `YYYY-MM-DD` |
| `--color` | - | Color hex code | e.g., `#FF5733` |
| `--goal` | - | Project goal/outcome | Text |
| `--description` | - | Project description | Text (markdown supported) |

### Examples

```bash
# Create root project under subarea (minimum)
./dopa projects create --name "Website Redesign" --subarea-id "subarea-123"

# Create nested project under another project
./dopa projects create --name "Backend API" --parent-id "project-456" --priority high

# Create project with all fields
./dopa projects create --name "Q4 Campaign" --subarea-id "subarea-123" \
  --status active \
  --priority urgent \
  --progress 25 \
  --start-date 2024-10-01 \
  --deadline 2024-12-31 \
  --color "#FF5733" \
  --goal "Launch campaign by year end" \
  --description "Marketing campaign for Q4 targeting new demographics"

# Create on-hold project
./dopa projects create --name "Research Project" --subarea-id "subarea-456" \
  --status on_hold \
  --priority low

# Create nested project with description
./dopa projects create --name "Database Migration" --parent-id "project-789" \
  --priority high \
  --description "## Objective\nMigrate from PostgreSQL to SQLite\n\n## Tasks\n- [ ] Export data\n- [ ] Transform schema"
```

### Output

```json
{
  "id": "project-def456",
  "name": "Website Redesign",
  "subarea_id": "subarea-123",
  "status": "active",
  "priority": "medium",
  "progress": 0,
  "created_at": "2024-01-15T10:40:00Z",
  "updated_at": "2024-01-15T10:40:00Z"
}
```

## List Projects

List all projects with optional filtering.

### Syntax

```bash
./dopa projects list [options]
```

### Filter Flags

| Flag | Description |
|------|-------------|
| `--status` | Filter by status (`active`, `completed`, `on_hold`, `archived`) |
| `--priority` | Filter by priority (`low`, `medium`, `high`, `urgent`) |
| `--subarea-id` | Filter by parent subarea ID |
| `--parent-id` | Filter by parent project ID |
| `--filter` | Advanced filter query (see Advanced Filtering section) |

### Output Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--json` | false | Output as JSON |
| `--format` | table | Output format: `table`, `json`, or `yaml` |

### Examples

```bash
# List all projects
./dopa projects list

# List active projects only
./dopa projects list --status active

# List high priority projects
./dopa projects list --priority high

# List projects under a subarea
./dopa projects list --subarea-id "subarea-123"

# List nested projects under a parent
./dopa projects list --parent-id "project-456"

# Combine filters
./dopa projects list --status active --priority high

# Output as JSON
./dopa projects list --json

# Output as YAML
./dopa projects list --format yaml

# Advanced filtering
./dopa projects list --filter 'status=active AND priority>=high'
./dopa projects list --filter 'progress>=50 OR status=completed'
```

## Get Project

Display details of a single project by ID.

### Syntax

```bash
./dopa projects get <project-id>
```

### Arguments

- `project-id`: The ID of the project to retrieve

### Examples

```bash
./dopa projects get project-def456
```

### Output

```json
{
  "id": "project-def456",
  "name": "Website Redesign",
  "description": "Complete redesign of company website",
  "goal": "Launch new website by Q2",
  "status": "active",
  "priority": "high",
  "progress": 45,
  "deadline": "2024-06-30T00:00:00Z",
  "color": "#3498DB",
  "subarea_id": "subarea-123",
  "parent_id": null,
  "created_at": "2024-01-15T10:40:00Z",
  "updated_at": "2024-02-20T14:30:00Z"
}
```

## Update Project

Update any editable field of a project.

### Syntax

```bash
./dopa projects update <project-id> [options]
```

### Arguments

- `project-id`: The ID of the project to update

### Editable Flags

All optional flags from create are available for update:

| Flag | Description |
|------|-------------|
| `--name` | Update project name |
| `--status` | Update status |
| `--priority` | Update priority |
| `--progress` | Update completion percentage |
| `--deadline` | Update deadline |
| `--start-date` | Update start date |
| `--color` | Update color |
| `--goal` | Update goal |
| `--description` | Update description |

**Note**: You cannot change `--subarea-id` or `--parent-id` after creation.

### Examples

```bash
# Update name
./dopa projects update project-123 --name "Website Redesign v2"

# Update status to completed
./dopa projects update project-123 --status completed

# Update progress
./dopa projects update project-123 --progress 75

# Update priority and progress
./dopa projects update project-123 --priority urgent --progress 80

# Update multiple fields
./dopa projects update project-123 \
  --name "New Website" \
  --status active \
  --priority high \
  --progress 50 \
  --deadline 2024-12-31

# Mark project as completed
./dopa projects update project-123 --status completed --progress 100
```

## Delete Project

Delete a project (soft delete by default).

### Syntax

```bash
./dopa projects delete <project-id> [--permanent]
```

### Arguments

- `project-id`: The ID of the project to delete

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--permanent` | false | Permanently delete from database |

### Examples

```bash
# Soft delete (default)
./dopa projects delete project-123

# Permanent delete
./dopa projects delete project-123 --permanent
```

**Warning**: Permanent delete cannot be undone and will also delete nested projects!

---

# Advanced Features

## Advanced Filtering

The `--filter` flag supports complex query expressions for list commands.

### Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equals | `status=active` |
| `!=` | Not equals | `status!=archived` |
| `<` | Less than | `progress<50` |
| `<=` | Less than or equal | `progress<=75` |
| `>` | Greater than | `progress>25` |
| `>=` | Greater than or equal | `priority>=high` |
| `AND` | Logical AND | `status=active AND priority=high` |
| `OR` | Logical OR | `status=completed OR progress=100` |

### Grouping

Use parentheses for complex conditions:

```bash
./dopa projects list --filter '(status=active OR status=on_hold) AND priority>=medium'
```

### Filter Examples

```bash
# Active projects with high priority
./dopa projects list --filter 'status=active AND priority>=high'

# Projects with progress over 50% or already completed
./dopa projects list --filter 'progress>50 OR status=completed'

# Not archived and not low priority
./dopa projects list --filter 'status!=archived AND priority!=low'

# Complex condition with grouping
./dopa projects list --filter '(status=active AND progress<25) OR priority=urgent'

# Progress range
./dopa projects list --filter 'progress>=25 AND progress<=75'
```

### Valid Field Names

- Areas: `name`, `color`, `created_at`, `updated_at`
- Subareas: `name`, `area_id`, `color`, `created_at`, `updated_at`
- Projects: `name`, `status`, `priority`, `progress`, `deadline`, `created_at`, `updated_at`

## Database Management

### Specifying Database Path

```bash
# Use default database (./dopa.db)
./dopa areas list

# Use custom database path
./dopa --db /path/to/custom.db areas list

# Use environment-specific database
./dopa --db ./databases/production.db projects list
```

### Database Location Best Practices

```bash
# Development
./dopa --db ./data/dev.db areas create --name "Dev Area"

# Testing
./dopa --db ./data/test.db projects list

# Production
./dopa --db /var/lib/dopa/prod.db projects list --format json
```

## Working with Nested Projects

Projects can be nested recursively to create hierarchies:

```bash
# Create root project
./dopa projects create --name "Website Redesign" --subarea-id "subarea-123"

# Create nested project (Phase 1)
./dopa projects create --name "Phase 1: Planning" --parent-id "project-root" --priority high

# Create sub-project under Phase 1
./dopa projects create --name "Requirements Gathering" --parent-id "project-phase1" --priority high

# Create another nested project (Phase 2)
./dopa projects create --name "Phase 2: Development" --parent-id "project-root" --priority medium

# List all projects under parent
./dopa projects list --parent-id "project-root"
```

### Project Hierarchy Example

```
Website Redesign (root)
├── Phase 1: Planning
│   ├── Requirements Gathering
│   └── Design Review
├── Phase 2: Development
│   ├── Backend API
│   └── Frontend UI
└── Phase 3: Deployment
    ├── Testing
    └── Production Release
```

---

# Use Cases and Workflows

## Setting Up a New Project Structure

```bash
# 1. Create top-level areas
./dopa areas create --name "Work" --color "#3498DB"
./dopa areas create --name "Personal" --color "#2ECC71"

# 2. Create subareas for Work
./dopa subareas create --name "Client Projects" --area-id "work-area-id"
./dopa subareas create --name "Internal Tools" --area-id "work-area-id"

# 3. Create projects
./dopa projects create --name "Client ABC Website" --subarea-id "client-subarea-id" --priority high
./dopa projects create --name "Internal Dashboard" --subarea-id "internal-subarea-id" --priority medium
```

## Tracking Project Progress

```bash
# Create project
./dopa projects create --name "Q1 Initiative" --subarea-id "sub-123" \
  --priority high --goal "Complete Q1 objectives"

# Update progress as work progresses
./dopa projects update project-id --progress 25
./dopa projects update project-id --progress 50
./dopa projects update project-id --progress 75

# Mark as completed
./dopa projects update project-id --status completed --progress 100
```

## Reviewing and Reporting

```bash
# Review all active projects
./dopa projects list --status active --format json > active-projects.json

# Find overdue projects (high priority not completed)
./dopa projects list --filter 'priority>=high AND status!=completed'

# Generate report of completed work
./dopa projects list --status completed --format yaml > completed-report.yaml

# Review project hierarchy
./dopa projects list --parent-id "parent-project-id" --format json
```

## Managing Priorities

```bash
# Find all urgent projects
./dopa projects list --priority urgent

# Escalate project priority
./dopa projects update project-id --priority urgent

# Review high priority items
./dopa projects list --filter 'priority=high OR priority=urgent' --status active
```

## Bulk Operations Workflow

```bash
# Create multiple areas
for area in "Engineering" "Marketing" "Sales"; do
  ./dopa areas create --name "$area"
done

# List all and save to file
./dopa areas list --json > areas-backup.json

# Filter and export
./dopa projects list --status active --format yaml > active-projects.yaml
```

---

# Best Practices

## 1. Use Meaningful Names

```bash
# Good
./dopa projects create --name "Q4 Marketing Campaign" --subarea-id "sub-123"

# Avoid
./dopa projects create --name "Project 1" --subarea-id "sub-123"
```

## 2. Leverage Colors for Visual Organization

```bash
# Use consistent color coding
./dopa areas create --name "Engineering" --color "#3498DB"  # Blue
./dopa areas create --name "Marketing" --color "#E74C3C"   # Red
./dopa areas create --name "Sales" --color "#2ECC71"       # Green
```

## 3. Set Clear Goals and Descriptions

```bash
./dopa projects create --name "API Migration" --subarea-id "sub-123" \
  --goal "Migrate all services to new API by end of Q2" \
  --description "## Objectives\n- Update authentication\n- Migrate endpoints\n- Update documentation"
```

## 4. Use Priority and Status Consistently

```bash
# Priority indicates importance
--priority urgent  # Critical, needs immediate attention
--priority high    # Important, should be done soon
--priority medium  # Normal priority
--priority low     # Nice to have, can wait

# Status indicates progress
--status active      # Currently being worked on
--status on_hold     # Paused temporarily
--status completed   # Finished
--status archived    # No longer relevant
```

## 5. Regular Reviews

```bash
# Weekly review of active projects
./dopa projects list --status active

# Find stale projects (on hold for too long)
./dopa projects list --status on_hold

# Review high priority items
./dopa projects list --filter 'priority>=high AND status=active'
```

## 6. Use Soft Delete by Default

```bash
# Safe: soft delete (can be recovered)
./dopa projects delete project-id

# Dangerous: permanent delete (cannot be undone)
./dopa projects delete project-id --permanent
```

## 7. Export Regularly

```bash
# Backup database regularly
cp ./dopa.db ./backups/dopa-$(date +%Y%m%d).db

# Export critical data
./dopa projects list --format json > backups/projects-$(date +%Y%m%d).json
./dopa areas list --format json > backups/areas-$(date +%Y%m%d).json
```

---

# Troubleshooting

## Common Errors

### Database Not Found

```bash
Error: database file not found: ./dopa.db
```

**Solution**: Ensure the database file exists or specify the correct path:
```bash
./dopa --db /correct/path/to/database.db areas list
```

### Invalid Color Format

```bash
Error: invalid color format: must be hex (e.g., #FF5733)
```

**Solution**: Use proper hex color format with `#` prefix:
```bash
./dopa areas create --name "Test" --color "#FF5733"
```

### Missing Required Fields

```bash
Error: name is required
```

**Solution**: Provide all required flags:
```bash
./dopa areas create --name "My Area"
```

### Invalid Status/Priority Values

```bash
Error: invalid status: must be one of [active completed on_hold archived]
```

**Solution**: Use valid enum values:
```bash
./dopa projects create --name "Test" --subarea-id "sub-123" --status active --priority high
```

### Subarea-Id and Parent-Id Conflict

```bash
Error: cannot specify both subarea-id and parent-id
```

**Solution**: Specify only one parent reference:
```bash
# Root project
./dopa projects create --name "Test" --subarea-id "sub-123"

# Nested project
./dopa projects create --name "Test" --parent-id "proj-456"
```

## Getting Help

```bash
# General help
./dopa --help

# Command-specific help
./dopa areas --help
./dopa projects create --help
./dopa projects list --help
```

---

# Exit Codes

The CLI uses standard exit codes:

| Code | Meaning | Example |
|------|---------|---------|
| 0 | Success | Operation completed successfully |
| 1 | Error | Database error, unexpected failure |
| 2 | Validation Error | Invalid input, missing required field |

Use these codes in scripts:

```bash
#!/bin/bash
./dopa areas create --name "Test"
if [ $? -eq 0 ]; then
  echo "Success!"
elif [ $? -eq 2 ]; then
  echo "Validation error"
else
  echo "System error"
fi
```

---

# Reference

## Command Aliases

| Full Command | Aliases |
|--------------|---------|
| `areas` | `area` |
| `subareas` | `subarea`, `sub` |
| `projects` | `project`, `proj` |

## Quick Reference Card

### Areas

```bash
areas create   --name <name> [--color <hex>]
areas list     [--json] [--format table|json|yaml]
areas get      <id>
areas update   <id> [--name <name>] [--color <hex>]
areas delete   <id> [--permanent]
```

### Subareas

```bash
subareas create  --name <name> --area-id <id> [--color <hex>]
subareas list    [--area-id <id>] [--json] [--format table|json|yaml]
subareas get     <id>
subareas update  <id> [--name <name>] [--color <hex>]
subareas delete  <id> [--permanent]
```

### Projects

```bash
projects create  --name <name> (--subarea-id <id> | --parent-id <id>) [options...]
projects list    [--status <s>] [--priority <p>] [--subarea-id <id>] [--parent-id <id>] [--filter <query>] [--json] [--format table|json|yaml]
projects get     <id>
projects update  <id> [--name <name>] [--status <s>] [--priority <p>] [--progress <n>] [options...]
projects delete  <id> [--permanent]
```

---

# See Also

- [Data Layer Architecture](doc-1 - Data-Layer-Architecture.md) - Technical details of database schema and implementation
- [Backlog Tasks](../tasks/) - Task breakdown and implementation details

---

# Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2026-03-03 | Initial release with full CRUD operations for Areas, Subareas, and Projects |
