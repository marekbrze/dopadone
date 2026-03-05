package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/example/projectdb/internal/converter"
	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

var (
	ErrSubareaNotFound = errors.New("subarea not found")
)

type SubareaStats struct {
	ProjectCount int64
}

type SubareaService struct {
	repo db.Querier
}

func NewSubareaService(repo db.Querier) *SubareaService {
	return &SubareaService{repo: repo}
}

func (s *SubareaService) Create(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
	subarea, err := domain.NewSubarea(name, areaID, color)
	if err != nil {
		return nil, err
	}

	params := db.CreateSubareaParams{
		ID:        subarea.ID,
		Name:      subarea.Name,
		AreaID:    subarea.AreaID,
		Color:     sql.NullString{String: string(subarea.Color), Valid: subarea.Color != ""},
		CreatedAt: subarea.CreatedAt,
		UpdatedAt: subarea.UpdatedAt,
		DeletedAt: nil,
	}

	dbResult, err := s.repo.CreateSubarea(ctx, params)
	if err != nil {
		return nil, err
	}

	result := converter.DbSubareaToDomain(dbResult)
	return &result, nil
}

func (s *SubareaService) GetByID(ctx context.Context, id string) (*domain.Subarea, error) {
	res, err := s.repo.GetSubareaByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubareaNotFound
		}
		return nil, err
	}
	result := converter.DbSubareaToDomain(res)
	return &result, nil
}

func (s *SubareaService) ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error) {
	rows, err := s.repo.ListSubareasByArea(ctx, areaID)
	if err != nil {
		return nil, err
	}

	subareas := make([]domain.Subarea, len(rows))
	for i, row := range rows {
		res := converter.DbSubareaToDomain(row)
		subareas[i] = res
	}
	return subareas, nil
}

func (s *SubareaService) Update(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
	params := db.UpdateSubareaParams{
		ID:        id,
		Name:      name,
		AreaID:    areaID,
		Color:     sql.NullString{String: string(color), Valid: color != ""},
		UpdatedAt: time.Now(),
	}

	dbResult, err := s.repo.UpdateSubarea(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSubareaNotFound
		}
		return nil, err
	}

	result := converter.DbSubareaToDomain(dbResult)
	return &result, nil
}

func (s *SubareaService) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	params := db.SoftDeleteSubareaParams{
		ID:        id,
		DeletedAt: now,
	}
	_, err := s.repo.SoftDeleteSubarea(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrSubareaNotFound
		}
		return err
	}
	return nil
}

func (s *SubareaService) HardDelete(ctx context.Context, id string) error {
	return s.repo.HardDeleteSubarea(ctx, id)
}

func (s *SubareaService) GetStats(ctx context.Context, id string) (*SubareaStats, error) {
	projectCount, err := s.repo.CountProjectsBySubarea(ctx, sql.NullString{String: id, Valid: true})
	if err != nil {
		return nil, err
	}

	return &SubareaStats{
		ProjectCount: projectCount,
	}, nil
}

func (s *SubareaService) GetEffectiveColor(subarea *domain.Subarea, parentArea *domain.Area) domain.Color {
	return subarea.GetEffectiveColor(parentArea)
}

func (s *SubareaService) ListAll(ctx context.Context) ([]domain.Subarea, error) {
	rows, err := s.repo.ListAllSubareas(ctx)
	if err != nil {
		return nil, err
	}

	subareas := make([]domain.Subarea, len(rows))
	for i, row := range rows {
		res := converter.DbSubareaToDomain(row)
		subareas[i] = res
	}
	return subareas, nil
}
