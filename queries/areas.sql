-- name: CreateArea :one
INSERT INTO areas (id, name, color, sort_order, created_at, updated_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING id, name, color, sort_order, created_at, updated_at, deleted_at;

-- name: GetAreaByID :one
SELECT id, name, color, sort_order, created_at, updated_at, deleted_at
FROM areas
WHERE id = ? AND deleted_at IS NULL;

-- name: ListAreas :many
SELECT id, name, color, sort_order, created_at, updated_at, deleted_at
FROM areas
WHERE deleted_at IS NULL
ORDER BY sort_order ASC;

-- name: UpdateArea :one
UPDATE areas
SET name = ?, color = ?, updated_at = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING id, name, color, sort_order, created_at, updated_at, deleted_at;

-- name: UpdateAreaSortOrder :exec
UPDATE areas
SET sort_order = ?, updated_at = ?
WHERE id = ?;

-- name: SoftDeleteArea :one
UPDATE areas
SET deleted_at = ?
WHERE id = ?
RETURNING id, name, color, sort_order, created_at, updated_at, deleted_at;

-- name: HardDeleteArea :exec
DELETE FROM areas
WHERE id = ?;

-- name: CountSubareasByArea :one
SELECT COUNT(*) FROM subareas WHERE area_id = ?;

-- name: CountProjectsByArea :one
SELECT COUNT(*) FROM projects WHERE subarea_id IN (SELECT id FROM subareas WHERE area_id = ?);

-- name: CountTasksByArea :one
SELECT COUNT(*) FROM tasks WHERE project_id IN (
    SELECT id FROM projects WHERE subarea_id IN (SELECT id FROM subareas WHERE area_id = ?)
);

-- name: DeleteSubareasByArea :exec
DELETE FROM subareas WHERE area_id = ?;

-- name: DeleteProjectsBySubarea :exec
DELETE FROM projects WHERE subarea_id IN (SELECT id FROM subareas WHERE area_id = ?);

-- name: DeleteTasksByProject :exec
DELETE FROM tasks WHERE project_id IN (
    SELECT id FROM projects WHERE subarea_id IN (SELECT id FROM subareas WHERE area_id = ?)
);
