package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidColorFormat   = errors.New("invalid color format: must be a valid hex color (e.g., #FF0000)")
	ErrInvalidProjectStatus = errors.New("invalid project status: must be one of active, completed, on_hold, archived")
	ErrInvalidPriority      = errors.New("invalid priority: must be one of low, medium, high, urgent")
	ErrInvalidProgress      = errors.New("invalid progress: must be between 0 and 100")
	ErrInvalidDateRange     = errors.New("invalid date range: start_date must be before deadline")
	ErrDeadlineWithoutStart = errors.New("deadline cannot be set without start_date")
)

type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusOnHold    ProjectStatus = "on_hold"
	ProjectStatusArchived  ProjectStatus = "archived"
)

func (s ProjectStatus) IsValid() bool {
	switch s {
	case ProjectStatusActive, ProjectStatusCompleted, ProjectStatusOnHold, ProjectStatusArchived:
		return true
	default:
		return false
	}
}

func (s ProjectStatus) String() string {
	return string(s)
}

func ParseProjectStatus(s string) (ProjectStatus, error) {
	status := ProjectStatus(s)
	if !status.IsValid() {
		return "", ErrInvalidProjectStatus
	}
	return status, nil
}

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

func (p Priority) String() string {
	return string(p)
}

func ParsePriority(s string) (Priority, error) {
	priority := Priority(s)
	if !priority.IsValid() {
		return "", ErrInvalidPriority
	}
	return priority, nil
}

type Progress int

func (p Progress) IsValid() bool {
	return p >= 0 && p <= 100
}

func (p Progress) Int() int {
	return int(p)
}

func ParseProgress(n int) (Progress, error) {
	p := Progress(n)
	if !p.IsValid() {
		return 0, ErrInvalidProgress
	}
	return p, nil
}

type DateRange struct {
	StartDate *time.Time
	Deadline  *time.Time
}

func (d DateRange) IsValid() bool {
	if d.StartDate == nil && d.Deadline == nil {
		return true
	}
	if d.Deadline != nil && d.StartDate == nil {
		return false
	}
	if d.StartDate != nil && d.Deadline != nil {
		return d.StartDate.Before(*d.Deadline)
	}
	return true
}

func NewDateRange(startDate, deadline *time.Time) (DateRange, error) {
	dr := DateRange{StartDate: startDate, Deadline: deadline}
	if !dr.IsValid() {
		return DateRange{}, ErrInvalidDateRange
	}
	return dr, nil
}

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type Color string

func (c Color) IsValid() bool {
	if c == "" {
		return true
	}
	return hexColorRegex.MatchString(string(c))
}

func (c Color) String() string {
	return string(c)
}

func ParseColor(s string) (Color, error) {
	if s == "" {
		return Color(""), nil
	}
	if !hexColorRegex.MatchString(s) {
		return "", ErrInvalidColorFormat
	}
	return Color(s), nil
}

var (
	ErrInvalidTaskStatus   = errors.New("invalid task status: must be one of todo, in_progress, waiting, done")
	ErrInvalidTaskPriority = errors.New("invalid task priority: must be one of critical, high, medium, low")
	ErrInvalidTaskDuration = errors.New("invalid task duration: must be one of 5, 15, 30, 60, 120, 240, 480 minutes")
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusWaiting    TaskStatus = "waiting"
	TaskStatusDone       TaskStatus = "done"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusTodo, TaskStatusInProgress, TaskStatusWaiting, TaskStatusDone:
		return true
	default:
		return false
	}
}

func (s TaskStatus) String() string {
	return string(s)
}

func ParseTaskStatus(s string) (TaskStatus, error) {
	status := TaskStatus(s)
	if !status.IsValid() {
		return "", ErrInvalidTaskStatus
	}
	return status, nil
}

type TaskPriority string

const (
	TaskPriorityCritical TaskPriority = "critical"
	TaskPriorityHigh     TaskPriority = "high"
	TaskPriorityMedium   TaskPriority = "medium"
	TaskPriorityLow      TaskPriority = "low"
)

func (p TaskPriority) IsValid() bool {
	switch p {
	case TaskPriorityCritical, TaskPriorityHigh, TaskPriorityMedium, TaskPriorityLow:
		return true
	default:
		return false
	}
}

func (p TaskPriority) String() string {
	return string(p)
}

func ParseTaskPriority(s string) (TaskPriority, error) {
	priority := TaskPriority(s)
	if !priority.IsValid() {
		return "", ErrInvalidTaskPriority
	}
	return priority, nil
}

type TaskDuration int

const (
	TaskDuration5   TaskDuration = 5
	TaskDuration15  TaskDuration = 15
	TaskDuration30  TaskDuration = 30
	TaskDuration60  TaskDuration = 60
	TaskDuration120 TaskDuration = 120
	TaskDuration240 TaskDuration = 240
	TaskDuration480 TaskDuration = 480
)

func (d TaskDuration) IsValid() bool {
	switch d {
	case TaskDuration5, TaskDuration15, TaskDuration30, TaskDuration60, TaskDuration120, TaskDuration240, TaskDuration480:
		return true
	default:
		return false
	}
}

func (d TaskDuration) Int() int {
	return int(d)
}

func ParseTaskDuration(n int) (TaskDuration, error) {
	duration := TaskDuration(n)
	if !duration.IsValid() {
		return 0, ErrInvalidTaskDuration
	}
	return duration, nil
}
