package service

import (
	"context"

	"github.com/example/dopadone/internal/domain"
)

// Package service provides business logic for project management operations.
//
// This package defines service interfaces for dependency injection and testability.
// Interfaces are defined in the provider package (alongside implementations) following
// the "accept interfaces, return structs" principle.
//
// Design Decision: Provider Pattern
//
// We define interfaces where they're implemented (provider pattern) rather than where
// they're consumed (consumer pattern). This approach:
//   - Keeps interfaces close to implementations for easier maintenance
//   - Allows consumers to define their own interfaces if needed
//   - Simplifies the dependency graph
//   - Enables straightforward mocking for tests
//
// For consumers that need different interface shapes, they can define their own
// narrower interfaces following Go's interface composition patterns.
//
// Context-First Design
//
// All service methods accept context.Context as their first parameter following Go best
// practices. This enables:
//   - Request cancellation and timeout support
//   - Tracing and distributed context propagation
//   - Consistent API design across all services

// AreaServiceInterface defines the contract for area business operations.
// Areas are top-level organizational units that contain subareas, projects, and tasks.
type AreaServiceInterface interface {
	// List retrieves all non-deleted areas sorted by sort_order.
	List(ctx context.Context) ([]domain.Area, error)

	// GetByID retrieves a single area by its unique identifier.
	// Returns an error if the area is not found or has been soft-deleted.
	GetByID(ctx context.Context, id string) (*domain.Area, error)

	// Create creates a new area with the given name and color.
	// The sort order is automatically assigned based on existing areas.
	Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error)

	// Update modifies an existing area's name and color.
	// Returns an error if the area is not found.
	Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error)

	// UpdateSortOrder changes the sort position of a single area.
	UpdateSortOrder(ctx context.Context, id string, sortOrder int) error

	// ReorderAll updates the sort order of all areas based on their positions in the list.
	// The order in the slice determines the new sort order.
	ReorderAll(ctx context.Context, areaIDs []string) error

	// SoftDelete marks an area as deleted without removing it from the database.
	// Note: Children (subareas, projects, tasks) become orphaned but are not deleted.
	SoftDelete(ctx context.Context, id string) error

	// HardDelete permanently removes an area and all its children from the database.
	// This operation cannot be undone.
	HardDelete(ctx context.Context, id string) error

	// GetStats retrieves statistics about an area's children (subareas, projects, tasks).
	GetStats(ctx context.Context, id string) (*AreaStats, error)
}

// SubareaServiceInterface defines the contract for subarea business operations.
// Subareas are organizational units within areas that can contain projects and tasks.
type SubareaServiceInterface interface {
	// Create creates a new subarea within the specified area.
	Create(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error)

	// GetByID retrieves a single subarea by its unique identifier.
	// Returns an error if the subarea is not found or has been soft-deleted.
	GetByID(ctx context.Context, id string) (*domain.Subarea, error)

	// ListByArea retrieves all non-deleted subareas within the specified area.
	ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error)

	// Update modifies an existing subarea's properties.
	Update(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error)

	// SoftDelete marks a subarea as deleted without removing it from the database.
	SoftDelete(ctx context.Context, id string) error

	// HardDelete permanently removes a subarea from the database.
	HardDelete(ctx context.Context, id string) error

	// GetStats retrieves statistics about a subarea's children (projects).
	GetStats(ctx context.Context, id string) (*SubareaStats, error)

	// GetEffectiveColor returns the color to use for a subarea, falling back to the
	// parent area's color if the subarea has no color set.
	// The context parameter is included for future-proofing (e.g., caching, tracing)
	// but is not currently used in the implementation.
	GetEffectiveColor(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color

	// ListAll retrieves all non-deleted subareas across all areas.
	ListAll(ctx context.Context) ([]domain.Subarea, error)
}

// ProjectServiceInterface defines the contract for project business operations.
// Projects are work items that can be organized hierarchically within subareas.
type ProjectServiceInterface interface {
	// Create creates a new project with the specified parameters.
	// Validates parent hierarchy if parentID is provided.
	Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error)

	// GetByID retrieves a single project by its unique identifier.
	GetByID(ctx context.Context, id string) (*domain.Project, error)

	// ListBySubarea retrieves all non-deleted projects directly within the specified subarea.
	ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error)

	// ListByParent retrieves all non-deleted projects that are children of the specified project.
	ListByParent(ctx context.Context, parentID string) ([]domain.Project, error)

	// ListAll retrieves all non-deleted projects across all subareas.
	ListAll(ctx context.Context) ([]domain.Project, error)

	// ListByStatus retrieves all non-deleted projects with the specified status.
	ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)

	// ListByPriority retrieves all non-deleted projects with the specified priority.
	ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error)

	// ListBySubareaRecursive retrieves all projects within a subarea, including nested projects.
	// This method will be implemented in Task-29B.
	ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error)

	// Update modifies an existing project's properties.
	// Validates parent hierarchy if parentID is being changed.
	Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error)

	// SoftDelete marks a project as deleted without removing it from the database.
	SoftDelete(ctx context.Context, id string) error

	// HardDelete permanently removes a project from the database.
	// Returns an error if the project has non-deleted children.
	HardDelete(ctx context.Context, id string) error

	// GetStats retrieves statistics about a project's children (tasks and sub-projects).
	GetStats(ctx context.Context, id string) (*ProjectStats, error)

	// ValidateParentHierarchy ensures that setting parentID as the parent of projectID
	// would not create a circular reference in the project hierarchy.
	ValidateParentHierarchy(ctx context.Context, parentID string, projectID string) error
}

// TaskServiceInterface defines the contract for task business operations.
// Tasks are actionable items within projects.
type TaskServiceInterface interface {
	// Create creates a new task with the specified parameters.
	Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error)

	// GetByID retrieves a single task by its unique identifier.
	GetByID(ctx context.Context, id string) (*domain.Task, error)

	// ListByProject retrieves all non-deleted tasks within the specified project.
	ListByProject(ctx context.Context, projectID string) ([]domain.Task, error)

	// ListByProjectRecursive retrieves all non-deleted tasks from a project and all its nested subprojects.
	ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error)

	// ListByStatus retrieves all non-deleted tasks with the specified status.
	ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)

	// ListByPriority retrieves all non-deleted tasks with the specified priority.
	ListByPriority(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error)

	// ListNext retrieves all non-deleted tasks marked with is_next=true.
	ListNext(ctx context.Context) ([]domain.Task, error)

	// ListAll retrieves all non-deleted tasks across all projects.
	ListAll(ctx context.Context) ([]domain.Task, error)

	// Update modifies an existing task's properties.
	Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error)

	// SoftDelete marks a task as deleted without removing it from the database.
	SoftDelete(ctx context.Context, id string) error

	// HardDelete permanently removes a task from the database.
	HardDelete(ctx context.Context, id string) error

	// SetStatus updates a task's status to the specified value.
	SetStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)

	// MarkCompleted marks a task as completed by setting its status to TaskStatusDone.
	MarkCompleted(ctx context.Context, id string) (*domain.Task, error)

	// SetPriority updates a task's priority to the specified value.
	SetPriority(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error)

	// ToggleIsNext flips the is_next flag on a task.
	ToggleIsNext(ctx context.Context, id string) (*domain.Task, error)
}

// Compile-time interface satisfaction checks ensure that our service implementations
// correctly implement their respective interfaces. This catches interface/implementation
// mismatches at compile time rather than runtime.
var (
	_ AreaServiceInterface    = (*AreaService)(nil)
	_ SubareaServiceInterface = (*SubareaService)(nil)
	_ ProjectServiceInterface = (*ProjectService)(nil)
	_ TaskServiceInterface    = (*TaskService)(nil)
)
