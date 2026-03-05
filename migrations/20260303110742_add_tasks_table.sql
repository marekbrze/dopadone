-- +goose Up
-- Add tasks table for task management

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    start_date TIMESTAMP NULL,
    deadline TIMESTAMP NULL,
    priority TEXT NOT NULL DEFAULT 'medium',
    context TEXT,
    estimated_duration INTEGER,
    status TEXT NOT NULL DEFAULT 'todo',
    is_next INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    -- Status ENUM constraint
    CHECK (status IN ('todo', 'in_progress', 'waiting', 'done')),
    -- Priority ENUM constraint
    CHECK (priority IN ('critical', 'high', 'medium', 'low'))
);

-- Indexes for performance
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_deadline ON tasks(deadline);
CREATE INDEX idx_tasks_is_next ON tasks(is_next);
CREATE INDEX idx_tasks_priority ON tasks(priority);

-- +goose Down
-- Down migration for tasks table
DROP TABLE IF EXISTS tasks;
