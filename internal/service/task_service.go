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
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService struct {
	repo db.Querier
}

func NewTaskService(repo db.Querier, tm *db.TransactionManager) *TaskService {
	return &TaskService{repo: repo}
}

type CreateTaskParams struct {
	ProjectID         string
	Title             string
	Description       string
	StartDate         *time.Time
	Deadline          *time.Time
	Priority          domain.TaskPriority
	Context           string
	EstimatedDuration domain.TaskDuration
	Status            domain.TaskStatus
	IsNext            bool
}

func (s *TaskService) Create(ctx context.Context, params CreateTaskParams) (*domain.Task, error) {
	task, err := domain.NewTask(domain.NewTaskParams{
		ProjectID:         params.ProjectID,
		Title:             params.Title,
		Description:       params.Description,
		StartDate:         params.StartDate,
		Deadline:          params.Deadline,
		Priority:          params.Priority,
		Context:           params.Context,
		EstimatedDuration: params.EstimatedDuration,
		Status:            params.Status,
		IsNext:            params.IsNext,
	})
	if err != nil {
		return nil, err
	}

	var isNext int64
	if task.IsNext {
		isNext = 1
	}

	var estimatedDuration sql.NullInt64
	if task.EstimatedDuration != 0 {
		estimatedDuration = sql.NullInt64{Int64: int64(task.EstimatedDuration.Int()), Valid: true}
	}

	dbParams := db.CreateTaskParams{
		ID:                task.ID,
		ProjectID:         task.ProjectID,
		Title:             task.Title,
		Description:       sql.NullString{String: task.Description, Valid: task.Description != ""},
		StartDate:         task.StartDate,
		Deadline:          task.Deadline,
		Priority:          task.Priority.String(),
		Context:           sql.NullString{String: task.Context, Valid: task.Context != ""},
		EstimatedDuration: estimatedDuration,
		Status:            task.Status.String(),
		IsNext:            isNext,
		CreatedAt:         task.CreatedAt,
		UpdatedAt:         task.UpdatedAt,
		DeletedAt:         task.DeletedAt,
	}

	dbResult, err := s.repo.CreateTask(ctx, dbParams)
	if err != nil {
		return nil, err
	}

	result := converter.DbTaskToDomain(dbResult)
	return &result, nil
}

func (s *TaskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	res, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	result := converter.DbTaskToDomain(res)
	return &result, nil
}

func (s *TaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error) {
	rows, err := s.repo.ListTasksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = converter.DbTaskToDomain(row)
	}
	return tasks, nil
}

func (s *TaskService) ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error) {
	rows, err := s.repo.ListTasksByStatus(ctx, status.String())
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = converter.DbTaskToDomain(row)
	}
	return tasks, nil
}

func (s *TaskService) ListByPriority(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error) {
	rows, err := s.repo.ListTasksByPriority(ctx, priority.String())
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = converter.DbTaskToDomain(row)
	}
	return tasks, nil
}

func (s *TaskService) ListNext(ctx context.Context) ([]domain.Task, error) {
	rows, err := s.repo.ListNextTasks(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = converter.DbTaskToDomain(row)
	}
	return tasks, nil
}

func (s *TaskService) ListAll(ctx context.Context) ([]domain.Task, error) {
	rows, err := s.repo.ListAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, len(rows))
	for i, row := range rows {
		tasks[i] = converter.DbTaskToDomain(row)
	}
	return tasks, nil
}

type UpdateTaskParams struct {
	ID                string
	Title             string
	Description       string
	StartDate         *time.Time
	Deadline          *time.Time
	Priority          domain.TaskPriority
	Context           string
	EstimatedDuration domain.TaskDuration
	Status            domain.TaskStatus
	IsNext            bool
}

func (s *TaskService) Update(ctx context.Context, params UpdateTaskParams) (*domain.Task, error) {
	var isNext int64
	if params.IsNext {
		isNext = 1
	}

	var estimatedDuration sql.NullInt64
	if params.EstimatedDuration != 0 {
		estimatedDuration = sql.NullInt64{Int64: int64(params.EstimatedDuration.Int()), Valid: true}
	}

	dbParams := db.UpdateTaskParams{
		ID:                params.ID,
		Title:             params.Title,
		Description:       sql.NullString{String: params.Description, Valid: params.Description != ""},
		StartDate:         params.StartDate,
		Deadline:          params.Deadline,
		Priority:          params.Priority.String(),
		Context:           sql.NullString{String: params.Context, Valid: params.Context != ""},
		EstimatedDuration: estimatedDuration,
		Status:            params.Status.String(),
		IsNext:            isNext,
		UpdatedAt:         time.Now(),
	}

	dbResult, err := s.repo.UpdateTask(ctx, dbParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	result := converter.DbTaskToDomain(dbResult)
	return &result, nil
}

func (s *TaskService) SoftDelete(ctx context.Context, id string) error {
	now := time.Now()
	params := db.SoftDeleteTaskParams{
		ID:        id,
		DeletedAt: &now,
	}
	_, err := s.repo.SoftDeleteTask(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

func (s *TaskService) HardDelete(ctx context.Context, id string) error {
	return s.repo.HardDeleteTask(ctx, id)
}

func (s *TaskService) SetStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error) {
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := task.SetStatus(status); err != nil {
		return nil, err
	}

	return s.Update(ctx, UpdateTaskParams{
		ID:                task.ID,
		Title:             task.Title,
		Description:       task.Description,
		StartDate:         task.StartDate,
		Deadline:          task.Deadline,
		Priority:          task.Priority,
		Context:           task.Context,
		EstimatedDuration: task.EstimatedDuration,
		Status:            task.Status,
		IsNext:            task.IsNext,
	})
}

func (s *TaskService) MarkCompleted(ctx context.Context, id string) (*domain.Task, error) {
	return s.SetStatus(ctx, id, domain.TaskStatusDone)
}

func (s *TaskService) SetPriority(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error) {
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := task.SetPriority(priority); err != nil {
		return nil, err
	}

	return s.Update(ctx, UpdateTaskParams{
		ID:                task.ID,
		Title:             task.Title,
		Description:       task.Description,
		StartDate:         task.StartDate,
		Deadline:          task.Deadline,
		Priority:          task.Priority,
		Context:           task.Context,
		EstimatedDuration: task.EstimatedDuration,
		Status:            task.Status,
		IsNext:            task.IsNext,
	})
}

func (s *TaskService) ToggleIsNext(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var newIsNext int64
	if task.IsNext {
		newIsNext = 0
	} else {
		newIsNext = 1
	}

	params := db.ToggleIsNextParams{
		ID:        id,
		IsNext:    newIsNext,
		UpdatedAt: time.Now(),
	}

	dbResult, err := s.repo.ToggleIsNext(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	result := converter.DbTaskToDomain(dbResult)
	return &result, nil
}
