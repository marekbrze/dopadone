-- name: CreateProject :one
INSERT INTO projects (id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at;

-- name: GetProjectByID :one
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE id = ? AND deleted_at IS NULL;

-- name: ListProjectsBySubarea :many
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE subarea_id = ? AND deleted_at IS NULL
ORDER BY position ASC, name ASC;

-- name: ListAllProjects :many
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE deleted_at IS NULL
ORDER BY position ASC, name ASC;

-- name: ListProjectsByParent :many
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE parent_id = ? AND deleted_at IS NULL
ORDER BY position ASC, name ASC;

-- name: UpdateProject :one
UPDATE projects
SET name = ?, description = ?, goal = ?, status = ?, priority = ?, progress = ?, deadline = ?, color = ?, parent_id = ?, subarea_id = ?, position = ?, updated_at = ?, completed_at = ?
WHERE id = ? AND deleted_at IS NULL
RETURNING id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at;

-- name: SoftDeleteProject :one
UPDATE projects
SET deleted_at = ?
WHERE id = ?
RETURNING id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at;

-- name: GetProjectsByStatus :many
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE status = ? AND deleted_at IS NULL
ORDER BY priority DESC, deadline ASC, name ASC;

-- name: HardDeleteProject :exec
DELETE FROM projects WHERE id = ?;

-- name: CountTasksByProject :one
SELECT COUNT(*) FROM tasks WHERE project_id = ?;

-- name: CountProjectsByParent :one
SELECT COUNT(*) FROM projects WHERE parent_id = ?;

-- name: ListProjectsByPriority :many
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM projects
WHERE priority = ? AND deleted_at IS NULL
ORDER BY priority DESC, deadline ASC, name ASC;

-- name: DeleteTasksByProjectID :exec
DELETE FROM tasks WHERE project_id = ?;

-- name: DeleteProjectsByParentID :exec
DELETE FROM projects WHERE parent_id = ?;

-- name: ListProjectsBySubareaRecursive :many
WITH RECURSIVE project_hierarchy AS (
    SELECT 
        id, name, description, goal, status, priority, progress, 
        deadline, color, parent_id, subarea_id, position, 
        created_at, updated_at, completed_at, deleted_at
    FROM projects
    WHERE projects.subarea_id = ? AND deleted_at IS NULL
    
    UNION ALL
    
    SELECT 
        p.id, p.name, p.description, p.goal, p.status, p.priority, p.progress,
        p.deadline, p.color, p.parent_id, p.subarea_id, p.position,
        p.created_at, p.updated_at, p.completed_at, p.deleted_at
    FROM projects p
    INNER JOIN project_hierarchy ph ON p.parent_id = ph.id
    WHERE p.deleted_at IS NULL
)
SELECT id, name, description, goal, status, priority, progress, deadline, color, parent_id, subarea_id, position, created_at, updated_at, completed_at, deleted_at
FROM project_hierarchy
ORDER BY position ASC, name ASC;
