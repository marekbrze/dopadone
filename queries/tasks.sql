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
