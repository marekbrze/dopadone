-- name: CreateTask :one
INSERT INTO tasks (id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at;

-- name: GetTaskByID :one
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE id = ? AND deleted_at IS NULL;

-- name: ListTasksByProject :many
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE project_id = ? AND deleted_at IS NULL
ORDER BY is_next DESC, priority DESC, deadline ASC, title ASC;

-- name: ListNextTasks :many
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE is_next = 1 AND deleted_at IS NULL
ORDER BY priority DESC, deadline ASC, title ASC;

-- name: ListTasksByStatus :many
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE status = ? AND deleted_at IS NULL
ORDER BY is_next DESC, priority DESC, deadline ASC, title ASC;

-- name: ListTasksByPriority :many
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE priority = ? AND deleted_at IS NULL
ORDER BY is_next DESC, deadline ASC, title ASC;

-- name: UpdateTask :one
UPDATE tasks
SET title = ?, description = ?, start_date = ?, deadline = ?, priority = ?, context = ?, estimated_duration = ?, status = ?, is_next = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at;

-- name: SoftDeleteTask :one
UPDATE tasks
SET deleted_at = ?
WHERE id = ?
RETURNING id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at;

-- name: ToggleIsNext :one
UPDATE tasks
SET is_next = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at;

-- name: HardDeleteTask :exec
DELETE FROM tasks WHERE id = ?;

-- name: ListAllTasks :many
SELECT id, project_id, title, description, start_date, deadline, priority, context, estimated_duration, status, is_next, created_at, updated_at, deleted_at
FROM tasks
WHERE deleted_at IS NULL
ORDER BY is_next DESC, priority DESC, deadline ASC, title ASC;

-- name: ListTasksByProjectRecursive :many
WITH RECURSIVE project_tree AS (
    SELECT projects.id FROM projects 
    WHERE projects.id = sqlc.narg('project_id') AND projects.deleted_at IS NULL
    
    UNION ALL
    
    SELECT p.id FROM projects p
    INNER JOIN project_tree pt ON p.parent_id = pt.id
    WHERE p.deleted_at IS NULL
)
SELECT t.id, t.project_id, t.title, t.description, t.start_date, t.deadline, t.priority, t.context, t.estimated_duration, t.status, t.is_next, t.created_at, t.updated_at, t.deleted_at
FROM tasks t
INNER JOIN project_tree pt ON t.project_id = pt.id
WHERE t.deleted_at IS NULL
ORDER BY t.is_next DESC, t.priority DESC, t.deadline ASC, t.title ASC;
