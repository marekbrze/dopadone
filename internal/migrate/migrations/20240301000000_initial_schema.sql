-- +goose Up
-- Initial schema for project management database

-- areas table
CREATE TABLE areas (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    color TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- subareas table
CREATE TABLE subareas (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    area_id TEXT NOT NULL,
    color TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE
);

-- projects table
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    goal TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    priority TEXT NOT NULL DEFAULT 'medium',
    progress INTEGER NOT NULL DEFAULT 0,
    deadline TIMESTAMP NULL,
    color TEXT,
    parent_id TEXT,
    subarea_id TEXT,
    position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (parent_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (subarea_id) REFERENCES subareas(id) ON DELETE CASCADE,
    -- Project must have either parent_id or subarea_id
    CHECK ((parent_id IS NOT NULL) OR (subarea_id IS NOT NULL)),
    -- Status ENUM constraint
    CHECK (status IN ('active', 'completed', 'on_hold', 'archived')),
    -- Priority ENUM constraint
    CHECK (priority IN ('low', 'medium', 'high', 'urgent'))
);

-- Indexes for performance
CREATE INDEX idx_projects_deadline ON projects(deadline);
CREATE INDEX idx_projects_status_priority ON projects(status, priority);
CREATE INDEX idx_projects_parent_id ON projects(parent_id);
CREATE INDEX idx_projects_subarea_id ON projects(subarea_id);
CREATE INDEX idx_subareas_area_id ON subareas(area_id);

-- +goose Down
-- Down migration for initial schema
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS subareas;
DROP TABLE IF EXISTS areas;
