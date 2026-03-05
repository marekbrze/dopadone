---
id: doc-1
title: Data Layer Architecture
type: technical
created_date: '2026-03-03 09:21'
---

# Data Layer Architecture

## Overview

This document describes the data layer architecture for the ADHD-friendly project management system, implemented using Domain-Driven Design principles with SQLite as the database.

## Technology Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Database | SQLite | Local file-based storage, zero configuration |
| Migrations | Goose | Schema versioning and migrations |
| Query Generation | sqlc | Type-safe SQL-to-Go code generation |

## Domain Model

### Hierarchy

```
Area (Top-level)
└── Subarea (Second-level)
    └── Project (Recursively nestable)
        └── Sub-project
```

### Entity Relationships

```
┌─────────┐       ┌───────────┐       ┌──────────┐
│  Area   │──1:N──│  Subarea  │──1:N──│ Project  │
└─────────┘       └───────────┘       └──────────┘
                                            │
                                            │ 1:N (self-ref)
                                            ▼
                                      ┌──────────┐
                                      │ Project  │ (nested)
                                      └──────────┘
```

## Tables

### areas

Top-level organization containers.

| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| name | TEXT | NOT NULL |
| color | TEXT | NULL (hex code) |
| created_at | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP |
| updated_at | TIMESTAMP | NOT NULL, DEFAULT CURRENT_TIMESTAMP |
| deleted_at | TIMESTAMP | NULL (soft delete) |

### subareas

Second-level grouping under areas.

| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| name | TEXT | NOT NULL |
| area_id | TEXT | NOT NULL → areas.id |
| color | TEXT | NULL (inherits from area if null) |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |
| deleted_at | TIMESTAMP | NULL (soft delete) |

### projects

Goal-oriented task groups, recursively nestable.

| Column | Type | Constraints |
|--------|------|-------------|
| id | TEXT | PRIMARY KEY |
| name | TEXT | NOT NULL |
| description | TEXT | NULL (markdown support) |
| goal | TEXT | NULL |
| status | TEXT | NOT NULL, CHECK: active/completed/on_hold/archived |
| priority | TEXT | NOT NULL, CHECK: low/medium/high/urgent |
| progress | INTEGER | NOT NULL, DEFAULT 0 |
| deadline | TIMESTAMP | NULL |
| color | TEXT | NULL |
| parent_id | TEXT | NULL → projects.id (self-reference) |
| subarea_id | TEXT | NULL → subareas.id |
| position | INTEGER | NOT NULL, DEFAULT 0 |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |
| completed_at | TIMESTAMP | NULL |
| deleted_at | TIMESTAMP | NULL (soft delete) |

**Constraint:** `(parent_id IS NOT NULL) OR (subarea_id IS NOT NULL)` — project must have either parent or subarea.

## Indexes

| Table | Index | Columns |
|-------|-------|---------|
| projects | idx_projects_deadline | deadline |
| projects | idx_projects_status_priority | status, priority |
| projects | idx_projects_parent_id | parent_id |
| projects | idx_projects_subarea_id | subarea_id |
| subareas | idx_subareas_area_id | area_id |

## Advanced Queries

### Recursive CTE for Hierarchical Filtering

The `ListProjectsBySubareaRecursive` query uses a recursive Common Table Expression (CTE) to efficiently filter hierarchical project data:

**Purpose**: Retrieve all projects belonging to a subarea, including nested projects whose parent chain leads to the subarea.

**Performance**: Server-side filtering reduces memory from O(n) to O(k) where k = filtered results.

**SQL Structure**:
```sql
-- name: ListProjectsBySubareaRecursive :many
WITH RECURSIVE project_hierarchy AS (
    -- Base case: projects directly in subarea
    SELECT 
        id, name, description, goal, status, priority, progress,
        deadline, color, parent_id, subarea_id, position,
        created_at, updated_at, completed_at, deleted_at
    FROM projects
    WHERE subarea_id = sqlc.narg('subarea_id') AND deleted_at IS NULL
    
    UNION ALL
    
    -- Recursive case: children of projects already in hierarchy
    SELECT 
        p.id, p.name, p.description, p.goal, p.status, p.priority, p.progress,
        p.deadline, p.color, p.parent_id, p.subarea_id, p.position,
        p.created_at, p.updated_at, p.completed_at, p.deleted_at
    FROM projects p
    INNER JOIN project_hierarchy ph ON p.parent_id = ph.id
    WHERE p.deleted_at IS NULL
)
SELECT * FROM project_hierarchy
ORDER BY position ASC, name ASC;
```

**Key Features**:
- **Base Case**: Selects projects with direct `subarea_id` membership
- **Recursive Case**: Joins to find children of projects already in the hierarchy
- **Soft Delete Filtering**: Applies `deleted_at IS NULL` at each recursion level
- **Ordering**: Results sorted by position, then name
- **No Depth Limit**: SQLite handles cycle detection automatically

**Use Case**: Optimized filtering for ProjectService.ListBySubareaRecursive (Task-32)

**Benchmark Results** (10% filter ratio):
- 100 projects → 1.4μs, 8.6KB, 22 allocs
- 500 projects → 7μs, 44.8KB, 102 allocs
- 1000 projects → 13.4μs, 89.7KB, 202 allocs

## Go Models (sqlc-generated)

### Area

```go
type Area struct {
    ID        string         `json:"id"`
    Name      string         `json:"name"`
    Color     sql.NullString `json:"color"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt *time.Time     `json:"deleted_at"`
}
```

### Subarea

```go
type Subarea struct {
    ID        string         `json:"id"`
    Name      string         `json:"name"`
    AreaID    string         `json:"area_id"`
    Color     sql.NullString `json:"color"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt *time.Time     `json:"deleted_at"`
}
```

### Project

```go
type Project struct {
    ID          string         `json:"id"`
    Name        string         `json:"name"`
    Description sql.NullString `json:"description"`
    Goal        sql.NullString `json:"goal"`
    Status      string         `json:"status"`
    Priority    string         `json:"priority"`
    Progress    int64          `json:"progress"`
    Deadline    *time.Time     `json:"deadline"`
    Color       sql.NullString `json:"color"`
    ParentID    sql.NullString `json:"parent_id"`
    SubareaID   sql.NullString `json:"subarea_id"`
    Position    int64          `json:"position"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    CompletedAt *time.Time     `json:"completed_at"`
    DeletedAt   *time.Time     `json:"deleted_at"`
}
```

## Querier Interface

```go
type Querier interface {
    // Areas
    CreateArea(ctx, CreateAreaParams) (Area, error)
    GetAreaByID(ctx, id string) (Area, error)
    ListAreas(ctx) ([]Area, error)
    UpdateArea(ctx, UpdateAreaParams) (Area, error)
    SoftDeleteArea(ctx, SoftDeleteAreaParams) (Area, error)

    // Subareas
    CreateSubarea(ctx, CreateSubareaParams) (Subarea, error)
    GetSubareaByID(ctx, id string) (Subarea, error)
    ListSubareasByArea(ctx, areaID string) ([]Subarea, error)
    UpdateSubarea(ctx, UpdateSubareaParams) (Subarea, error)
    SoftDeleteSubarea(ctx, SoftDeleteSubareaParams) (Subarea, error)

    // Projects
    CreateProject(ctx, CreateProjectParams) (Project, error)
    GetProjectByID(ctx, id string) (Project, error)
    GetProjectsByStatus(ctx, status string) ([]Project, error)
    ListProjectsByParent(ctx, parentID sql.NullString) ([]Project, error)
    ListProjectsBySubarea(ctx, subareaID sql.NullString) ([]Project, error)
    ListProjectsBySubareaRecursive(ctx, subareaID sql.NullString) ([]ListProjectsBySubareaRecursiveRow, error)
    UpdateProject(ctx, UpdateProjectParams) (Project, error)
    SoftDeleteProject(ctx, SoftDeleteProjectParams) (Project, error)
}
```

## File Structure

```
migrations/
  20240301000000_initial_schema.sql    # Up/down migrations

queries/
  areas.sql      # Area CRUD queries
  subareas.sql   # Subarea CRUD queries
  projects.sql   # Project CRUD queries

internal/db/
  models.go          # sqlc-generated structs
  querier.go         # sqlc-generated interface
  areas.sql.go       # sqlc-generated area queries
  subareas.sql.go    # sqlc-generated subarea queries
  projects.sql.go    # sqlc-generated project queries
  db.go              # Database connection

sqlc.yaml            # sqlc configuration
```

## Usage Example

```go
import (
    "database/sql"
    "github.com/yourproject/internal/db"
)

// Open database
database, err := sql.Open("sqlite3", "./data.db")
if err != nil {
    panic(err)
}

// Create queries instance
queries := db.New(database)

// Create area
area, err := queries.CreateArea(ctx, db.CreateAreaParams{
    ID:    uuid.New().String(),
    Name:  "Work",
    Color: sql.NullString{String: "#FF5733", Valid: true},
})

// List projects by subarea
projects, err := queries.ListProjectsBySubarea(ctx, sql.NullString{
    String: subareaID,
    Valid:  true,
})
```

## Design Decisions

1. **Soft Delete Pattern**: Records are never hard-deleted; `deleted_at` is set to mark deletion
2. **UUID Strings**: Using string-based UUIDs for portability
3. **Recursive Projects**: `parent_id` allows unlimited nesting of projects
4. **Flexible Status**: Status stored as TEXT with CHECK constraint for flexibility
5. **Color Inheritance**: Subareas/projects can inherit color from parent if not specified
