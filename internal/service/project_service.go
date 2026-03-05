package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/example/projectdb/internal/converter"
	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

var (
	ErrProjectNotFound    = errors.New("project not found")
	ErrCircularReference  = errors.New("circular reference detected: project cannot be its own ancestor")
	ErrProjectHasChildren = errors.New("project has non-deleted children and cannot be hard deleted")
	ErrNotImplemented     = errors.New("method not yet implemented")
)

type ProjectStats struct {
	TaskCount    int64
	ProjectCount int64
}

type ProjectService struct {
	repo db.Querier
	tm   *db.TransactionManager
}

func NewProjectService(repo db.Querier, tm *db.TransactionManager) *ProjectService {
	return &ProjectService{repo: repo, tm: tm}
}

type CreateProjectParams struct {
	Name        string
	Description string
	Goal        string
	Status      domain.ProjectStatus
	Priority    domain.Priority
	Progress    domain.Progress
	StartDate   *time.Time
	Deadline    *time.Time
	Color       domain.Color
	ParentID    *string
	SubareaID   *string
	Position    int
}

func (s *ProjectService) Create(ctx context.Context, params CreateProjectParams) (*domain.Project, error) {
	if params.ParentID != nil {
		if err := s.ValidateParentHierarchy(ctx, *params.ParentID, ""); err != nil {
			return nil, err
		}
	}

	project, err := domain.NewProject(domain.NewProjectParams{
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
	})
	if err != nil {
		return nil, err
	}

	dbParams := db.CreateProjectParams{
		ID:          project.ID,
		Name:        project.Name,
		Description: sql.NullString{String: project.Description, Valid: project.Description != ""},
		Goal:        sql.NullString{String: project.Goal, Valid: project.Goal != ""},
		Status:      project.Status.String(),
		Priority:    project.Priority.String(),
		Progress:    int64(project.Progress.Int()),
		Deadline:    project.Deadline,
		Color:       sql.NullString{String: string(project.Color), Valid: project.Color != ""},
		ParentID:    sql.NullString{String: stringPtr(project.ParentID), Valid: project.ParentID != nil},
		SubareaID:   sql.NullString{String: stringPtr(project.SubareaID), Valid: project.SubareaID != nil},
		Position:    int64(project.Position),
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		CompletedAt: project.CompletedAt,
		DeletedAt:   project.DeletedAt,
	}

	dbResult, err := s.repo.CreateProject(ctx, dbParams)
	if err != nil {
		return nil, err
	}

	result := converter.DbProjectToDomain(dbResult)
	return &result, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	dbResult, err := s.repo.GetProjectByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	result := converter.DbProjectToDomain(dbResult)
	return &result, nil
}

func (s *ProjectService) ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error) {
	rows, err := s.repo.ListProjectsBySubarea(ctx, sql.NullString{String: subareaID, Valid: true})
	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectToDomain(row)
	}
	return projects, nil
}

func (s *ProjectService) ListByParent(ctx context.Context, parentID string) ([]domain.Project, error) {
	rows, err := s.repo.ListProjectsByParent(ctx, sql.NullString{String: parentID, Valid: true})
	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectToDomain(row)
	}
	return projects, nil
}

func (s *ProjectService) ListAll(ctx context.Context) ([]domain.Project, error) {
	rows, err := s.repo.ListAllProjects(ctx)
	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectToDomain(row)
	}
	return projects, nil
}

func (s *ProjectService) ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error) {
	rows, err := s.repo.GetProjectsByStatus(ctx, status.String())
	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectToDomain(row)
	}
	return projects, nil
}

func (s *ProjectService) ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error) {
	rows, err := s.repo.ListProjectsByPriority(ctx, priority.String())
	if err != nil {
		return nil, err
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectToDomain(row)
	}
	return projects, nil
}

func (s *ProjectService) ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error) {
	if subareaID == "" {
		return []domain.Project{}, nil
	}

	rows, err := s.repo.ListProjectsBySubareaRecursive(ctx, sql.NullString{
		String: subareaID,
		Valid:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("list projects for subarea %s: %w", subareaID, err)
	}

	projects := make([]domain.Project, len(rows))
	for i, row := range rows {
		projects[i] = converter.DbProjectRowToDomain(row)
	}
	return projects, nil
}

type UpdateProjectParams struct {
	ID          string
	Name        string
	Description string
	Goal        string
	Status      domain.ProjectStatus
	Priority    domain.Priority
	Progress    domain.Progress
	StartDate   *time.Time
	Deadline    *time.Time
	Color       domain.Color
	ParentID    *string
	SubareaID   *string
	Position    int
}

func (s *ProjectService) Update(ctx context.Context, params UpdateProjectParams) (*domain.Project, error) {
	if params.ParentID != nil {
		if err := s.ValidateParentHierarchy(ctx, *params.ParentID, params.ID); err != nil {
			return nil, err
		}
	}

	dbParams := db.UpdateProjectParams{
		ID:          params.ID,
		Name:        params.Name,
		Description: sql.NullString{String: params.Description, Valid: params.Description != ""},
		Goal:        sql.NullString{String: params.Goal, Valid: params.Goal != ""},
		Status:      params.Status.String(),
		Priority:    params.Priority.String(),
		Progress:    int64(params.Progress.Int()),
		Deadline:    params.Deadline,
		Color:       sql.NullString{String: string(params.Color), Valid: params.Color != ""},
		ParentID:    sql.NullString{String: stringPtr(params.ParentID), Valid: params.ParentID != nil},
		SubareaID:   sql.NullString{String: stringPtr(params.SubareaID), Valid: params.SubareaID != nil},
		Position:    int64(params.Position),
		UpdatedAt:   time.Now(),
		CompletedAt: nil,
	}

	dbResult, err := s.repo.UpdateProject(ctx, dbParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	result := converter.DbProjectToDomain(dbResult)
	return &result, nil
}

func (s *ProjectService) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	params := db.SoftDeleteProjectParams{
		ID:        id,
		DeletedAt: &now,
	}
	_, err := s.repo.SoftDeleteProject(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}
	return nil
}

func (s *ProjectService) HardDelete(ctx context.Context, id string) error {
	if s.tm == nil {
		if err := s.hardDeleteRecursive(ctx, s.repo, id); err != nil {
			return err
		}
		return s.repo.HardDeleteProject(ctx, id)
	}

	return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
		if err := s.hardDeleteRecursive(ctx, tx, id); err != nil {
			return err
		}
		return tx.HardDeleteProject(ctx, id)
	})
}

func (s *ProjectService) hardDeleteRecursive(ctx context.Context, q db.Querier, projectID string) error {
	children, err := q.ListProjectsByParent(ctx, sql.NullString{String: projectID, Valid: true})
	if err != nil {
		return err
	}

	for _, child := range children {
		if err := s.hardDeleteRecursive(ctx, q, child.ID); err != nil {
			return err
		}
		if err := q.DeleteTasksByProjectID(ctx, child.ID); err != nil {
			return err
		}
		if err := q.HardDeleteProject(ctx, child.ID); err != nil {
			return err
		}
	}

	return q.DeleteTasksByProjectID(ctx, projectID)
}

func (s *ProjectService) GetStats(ctx context.Context, id string) (*ProjectStats, error) {
	taskCount, err := s.repo.CountTasksByProject(ctx, id)
	if err != nil {
		return nil, err
	}

	projectCount, err := s.repo.CountProjectsByParent(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &ProjectStats{
		TaskCount:    taskCount,
		ProjectCount: projectCount,
	}, nil
}

func (s *ProjectService) ValidateParentHierarchy(ctx context.Context, parentID string, projectID string) error {
	if projectID != "" && projectID == parentID {
		return ErrCircularReference
	}

	currentID := parentID
	for {
		project, err := s.repo.GetProjectByID(ctx, currentID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		if projectID != "" && project.ID == projectID {
			return ErrCircularReference
		}

		if project.ParentID.Valid {
			currentID = project.ParentID.String
		} else {
			break
		}
	}

	return nil
}

func stringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
