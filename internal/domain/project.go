package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProjectNameEmpty        = errors.New("project name cannot be empty")
	ErrProjectInvalidStatus    = errors.New("project status is invalid")
	ErrProjectInvalidPriority  = errors.New("project priority is invalid")
	ErrProjectInvalidProgress  = errors.New("project progress must be between 0 and 100")
	ErrProjectNoParent         = errors.New("project must have either parent_id or subarea_id")
	ErrProjectInvalidDateRange = errors.New("project deadline must be after start date")
)

type Project struct {
	ID          string        // Unique identifier (UUID)
	Name        string        // Display name of the project
	Description string        // Optional description of the project
	Goal        string        // Optional goal/outcome of the project
	Status      ProjectStatus // Current status (active, completed, on_hold, archived)
	Priority    Priority      // Priority level (low, medium, high, urgent)
	Progress    Progress      // Completion percentage (0-100)
	StartDate   *time.Time    // Optional start date for the project
	Deadline    *time.Time    // Optional deadline for the project
	Color       Color         // Optional color for UI display (hex format: #RRGGBB)
	ParentID    *string       // Foreign key to parent Project (for nested projects)
	SubareaID   *string       // Foreign key to parent Subarea (for root projects)
	Position    int           // Position for ordering within parent
	CreatedAt   time.Time     // Timestamp when the project was created
	UpdatedAt   time.Time     // Timestamp when the project was last updated
	CompletedAt *time.Time    // Timestamp when the project was completed (nil if not completed)
	DeletedAt   *time.Time    // Timestamp when the project was soft-deleted (nil if not deleted)
}

type NewProjectParams struct {
	Name        string
	Description string
	Goal        string
	Status      ProjectStatus
	Priority    Priority
	Progress    Progress
	StartDate   *time.Time
	Deadline    *time.Time
	Color       Color
	ParentID    *string
	SubareaID   *string
	Position    int
}

func NewProject(params NewProjectParams) (*Project, error) {
	if params.Name == "" {
		return nil, ErrProjectNameEmpty
	}

	if !params.Status.IsValid() {
		return nil, ErrProjectInvalidStatus
	}

	if !params.Priority.IsValid() {
		return nil, ErrProjectInvalidPriority
	}

	if !params.Progress.IsValid() {
		return nil, ErrProjectInvalidProgress
	}

	if params.ParentID == nil && params.SubareaID == nil {
		return nil, ErrProjectNoParent
	}

	if !params.Color.IsValid() {
		return nil, ErrInvalidColorFormat
	}

	if params.StartDate != nil && params.Deadline != nil {
		if !params.StartDate.Before(*params.Deadline) {
			return nil, ErrProjectInvalidDateRange
		}
	}
	if params.Deadline != nil && params.StartDate == nil {
		return nil, ErrProjectInvalidDateRange
	}

	now := time.Now()
	return &Project{
		ID:          uuid.New().String(),
		Name:        params.Name,
		Description: params.Description,
		Goal:        params.Goal,
		Status:      params.Status,
		Priority:    params.Priority,
		Progress:    params.Progress,
		StartDate:   params.StartDate,
		Deadline:    params.Deadline,
		Color:       params.Color,
		ParentID:    params.ParentID,
		SubareaID:   params.SubareaID,
		Position:    params.Position,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	}, nil
}

func (p *Project) IsDeleted() bool {
	return p.DeletedAt != nil
}

func (p *Project) IsCompleted() bool {
	return p.Status == ProjectStatusCompleted
}

func (p *Project) IsNested() bool {
	return p.ParentID != nil
}

func (p *Project) MarkCompleted(completedAt time.Time) {
	p.Status = ProjectStatusCompleted
	p.Progress = Progress(100)
	p.CompletedAt = &completedAt
	p.UpdatedAt = time.Now()
}

func (p *Project) SetProgress(progress Progress) error {
	if !progress.IsValid() {
		return ErrProjectInvalidProgress
	}
	p.Progress = progress
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) SetPriority(priority Priority) error {
	if !priority.IsValid() {
		return ErrProjectInvalidPriority
	}
	p.Priority = priority
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Project) SetStatus(status ProjectStatus) error {
	if !status.IsValid() {
		return ErrProjectInvalidStatus
	}
	p.Status = status
	p.UpdatedAt = time.Now()
	return nil
}
