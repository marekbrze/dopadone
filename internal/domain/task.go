package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTaskTitleEmpty       = errors.New("task title cannot be empty")
	ErrTaskProjectIDEmpty   = errors.New("task project_id cannot be empty")
	ErrTaskInvalidStatus    = errors.New("task status is invalid")
	ErrTaskInvalidPriority  = errors.New("task priority is invalid")
	ErrTaskInvalidDuration  = errors.New("task estimated_duration is invalid")
	ErrTaskInvalidDateRange = errors.New("task deadline must be after start date")
	ErrTaskDeadlineNoStart  = errors.New("task deadline cannot be set without start_date")
)

type Task struct {
	ID                string
	ProjectID         string
	Title             string
	Description       string
	StartDate         *time.Time
	Deadline          *time.Time
	Priority          TaskPriority
	Context           string
	EstimatedDuration TaskDuration
	Status            TaskStatus
	IsNext            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type NewTaskParams struct {
	ProjectID         string
	Title             string
	Description       string
	StartDate         *time.Time
	Deadline          *time.Time
	Priority          TaskPriority
	Context           string
	EstimatedDuration TaskDuration
	Status            TaskStatus
	IsNext            bool
}

func NewTask(params NewTaskParams) (*Task, error) {
	if params.Title == "" {
		return nil, ErrTaskTitleEmpty
	}

	if params.ProjectID == "" {
		return nil, ErrTaskProjectIDEmpty
	}

	if !params.Status.IsValid() {
		return nil, ErrTaskInvalidStatus
	}

	if !params.Priority.IsValid() {
		return nil, ErrTaskInvalidPriority
	}

	if params.EstimatedDuration != 0 && !params.EstimatedDuration.IsValid() {
		return nil, ErrTaskInvalidDuration
	}

	if params.Deadline != nil && params.StartDate == nil {
		return nil, ErrTaskDeadlineNoStart
	}

	if params.StartDate != nil && params.Deadline != nil {
		if !params.StartDate.Before(*params.Deadline) {
			return nil, ErrTaskInvalidDateRange
		}
	}

	now := time.Now()
	return &Task{
		ID:                uuid.New().String(),
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
		CreatedAt:         now,
		UpdatedAt:         now,
		DeletedAt:         nil,
	}, nil
}

func (t *Task) IsDeleted() bool {
	return t.DeletedAt != nil
}

func (t *Task) IsCompleted() bool {
	return t.Status == TaskStatusDone
}

func (t *Task) MarkCompleted(completedAt time.Time) {
	t.Status = TaskStatusDone
	t.UpdatedAt = time.Now()
}

func (t *Task) SetStatus(status TaskStatus) error {
	if !status.IsValid() {
		return ErrTaskInvalidStatus
	}
	t.Status = status
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) SetPriority(priority TaskPriority) error {
	if !priority.IsValid() {
		return ErrTaskInvalidPriority
	}
	t.Priority = priority
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Task) SetNext() {
	t.IsNext = true
	t.UpdatedAt = time.Now()
}

func (t *Task) ClearNext() {
	t.IsNext = false
	t.UpdatedAt = time.Now()
}
