package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/example/dopadone/internal/converter"
	"github.com/example/dopadone/internal/db"
	"github.com/example/dopadone/internal/domain"
)

// AreaStats contains statistics about an area's children
type AreaStats struct {
	SubareaCount int64
	ProjectCount int64
	TaskCount    int64
}

// AreaService provides business logic for area operations
type AreaService struct {
	repo db.Querier
	tm   *db.TransactionManager
}

// NewAreaService creates a new AreaService
func NewAreaService(repo db.Querier, tm *db.TransactionManager) *AreaService {
	return &AreaService{repo: repo, tm: tm}
}

// List retrieves all non-deleted areas sorted by sort_order
func (s *AreaService) List(ctx context.Context) ([]domain.Area, error) {
	rows, err := s.repo.ListAreas(ctx)
	if err != nil {
		return nil, err
	}

	areas := make([]domain.Area, len(rows))
	for i, row := range rows {
		areas[i] = converter.DbListAreasRowToDomain(row)
	}
	return areas, nil
}

// GetByID retrieves a single area by ID
func (s *AreaService) GetByID(ctx context.Context, id string) (*domain.Area, error) {
	row, err := s.repo.GetAreaByID(ctx, id)
	if err != nil {
		return nil, err
	}
	area := converter.DbGetAreaByIDRowToDomain(row)
	return &area, nil
}

// Create creates a new area
func (s *AreaService) Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
	// Get the next sort order
	areas, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	nextSortOrder := len(areas)

	area, err := domain.NewArea(name, color, nextSortOrder)
	if err != nil {
		return nil, err
	}

	params := db.CreateAreaParams{
		ID:        area.ID,
		Name:      area.Name,
		Color:     sql.NullString{String: string(area.Color), Valid: area.Color != ""},
		SortOrder: int64(area.SortOrder),
		CreatedAt: area.CreatedAt,
		UpdatedAt: area.UpdatedAt,
		DeletedAt: nil,
	}

	row, err := s.repo.CreateArea(ctx, params)
	if err != nil {
		return nil, err
	}

	result := converter.DbCreateAreaRowToDomain(row)
	return &result, nil
}

// Update updates an existing area's name and color
func (s *AreaService) Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error) {
	params := db.UpdateAreaParams{
		ID:        id,
		Name:      name,
		Color:     sql.NullString{String: string(color), Valid: color != ""},
		UpdatedAt: time.Now(),
	}

	row, err := s.repo.UpdateArea(ctx, params)
	if err != nil {
		return nil, err
	}

	result := converter.DbUpdateAreaRowToDomain(row)
	return &result, nil
}

// UpdateSortOrder updates the sort order of a single area
func (s *AreaService) UpdateSortOrder(ctx context.Context, id string, sortOrder int) error {
	params := db.UpdateAreaSortOrderParams{
		ID:        id,
		SortOrder: int64(sortOrder),
		UpdatedAt: time.Now(),
	}
	return s.repo.UpdateAreaSortOrder(ctx, params)
}

// ReorderAll updates the sort order of all areas based on their positions in the list
func (s *AreaService) ReorderAll(ctx context.Context, areaIDs []string) error {
	if s.tm == nil {
		for i, id := range areaIDs {
			if err := s.UpdateSortOrder(ctx, id, i); err != nil {
				return err
			}
		}
		return nil
	}

	return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
		for i, id := range areaIDs {
			params := db.UpdateAreaSortOrderParams{
				ID:        id,
				SortOrder: int64(i),
				UpdatedAt: time.Now(),
			}
			if err := tx.UpdateAreaSortOrder(ctx, params); err != nil {
				return err
			}
		}
		return nil
	})
}

// SoftDelete marks an area as deleted (children become orphaned)
func (s *AreaService) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	params := db.SoftDeleteAreaParams{
		ID:        id,
		DeletedAt: &now,
	}
	_, err := s.repo.SoftDeleteArea(ctx, params)
	return err
}

// HardDelete permanently deletes an area and all its children
func (s *AreaService) HardDelete(ctx context.Context, id string) error {
	if s.tm == nil {
		if err := s.repo.DeleteTasksByProject(ctx, id); err != nil {
			return err
		}
		if err := s.repo.DeleteProjectsBySubarea(ctx, id); err != nil {
			return err
		}
		if err := s.repo.DeleteSubareasByArea(ctx, id); err != nil {
			return err
		}
		return s.repo.HardDeleteArea(ctx, id)
	}

	return s.tm.WithTransaction(ctx, func(ctx context.Context, tx db.Querier) error {
		if err := tx.DeleteTasksByProject(ctx, id); err != nil {
			return err
		}
		if err := tx.DeleteProjectsBySubarea(ctx, id); err != nil {
			return err
		}
		if err := tx.DeleteSubareasByArea(ctx, id); err != nil {
			return err
		}
		return tx.HardDeleteArea(ctx, id)
	})
}

// GetStats retrieves statistics about an area's children
func (s *AreaService) GetStats(ctx context.Context, id string) (*AreaStats, error) {
	subareaCount, err := s.repo.CountSubareasByArea(ctx, id)
	if err != nil {
		return nil, err
	}

	projectCount, err := s.repo.CountProjectsByArea(ctx, id)
	if err != nil {
		return nil, err
	}

	taskCount, err := s.repo.CountTasksByArea(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AreaStats{
		SubareaCount: subareaCount,
		ProjectCount: projectCount,
		TaskCount:    taskCount,
	}, nil
}
