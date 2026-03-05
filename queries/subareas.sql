-- name: CreateSubarea :one
INSERT INTO subareas (id, name, area_id, color, created_at, updated_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id, name, area_id, color, created_at, updated_at, deleted_at;

-- name: GetSubareaByID :one
SELECT id, name, area_id, color, created_at, updated_at, deleted_at
FROM subareas
WHERE id = ? AND deleted_at IS NULL;

-- name: ListSubareasByArea :many
SELECT id, name, area_id, color, created_at, updated_at, deleted_at
FROM subareas
WHERE area_id = ? AND deleted_at IS NULL
ORDER BY name ASC;

-- name: UpdateSubarea :one
UPDATE subareas
SET name = ?, area_id = ?, color = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING id, name, area_id, color, created_at, updated_at, deleted_at;

-- name: SoftDeleteSubarea :one
UPDATE subareas
SET deleted_at = ?
WHERE id = ?
RETURNING id, name, area_id, color, created_at, updated_at, deleted_at;

-- name: HardDeleteSubarea :exec
DELETE FROM subareas WHERE id = ?;

-- name: CountProjectsBySubarea :one
SELECT COUNT(*) FROM projects WHERE subarea_id = ?;

-- name: ListAllSubareas :many
SELECT id, name, area_id, color, created_at, updated_at, deleted_at
FROM subareas
WHERE deleted_at IS NULL
ORDER BY name ASC;

-- name: DeleteProjectsBySubareaID :exec
DELETE FROM projects WHERE subarea_id = ?;

-- name: DeleteTasksBySubareaID :exec
DELETE FROM tasks WHERE project_id IN (SELECT id FROM projects WHERE subarea_id = ?);
