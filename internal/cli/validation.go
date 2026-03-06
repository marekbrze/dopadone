package cli

import (
	"time"

	"github.com/example/dopadone/internal/domain"
)

func ParseProjectStatus(s string) (domain.ProjectStatus, error) {
	status, err := domain.ParseProjectStatus(s)
	if err != nil {
		return "", NewValidationError("status", err.Error())
	}
	return status, nil
}

func ParsePriority(s string) (domain.Priority, error) {
	priority, err := domain.ParsePriority(s)
	if err != nil {
		return "", NewValidationError("priority", err.Error())
	}
	return priority, nil
}

func ParseProgress(n int) (domain.Progress, error) {
	progress, err := domain.ParseProgress(n)
	if err != nil {
		return 0, NewValidationError("progress", err.Error())
	}
	return progress, nil
}

func ParseColor(s string) (domain.Color, error) {
	color, err := domain.ParseColor(s)
	if err != nil {
		return "", NewValidationError("color", err.Error())
	}
	return color, nil
}

func ParseDate(startDateStr, deadlineStr string) (*time.Time, *time.Time, error) {
	var startDate *time.Time
	var deadline *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, nil, NewValidationError("start_date", "invalid date format, use YYYY-MM-DD")
		}
		startDate = &parsed
	}

	if deadlineStr != "" {
		parsed, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			return nil, nil, NewValidationError("deadline", "invalid date format, use YYYY-MM-DD")
		}
		deadline = &parsed
	}

	_, err := domain.NewDateRange(startDate, deadline)
	if err != nil {
		return nil, nil, NewValidationError("date_range", err.Error())
	}

	return startDate, deadline, nil
}

func ValidateProjectName(name string) error {
	if name == "" {
		return NewValidationError("name", "project name cannot be empty")
	}
	return nil
}

func ParseTaskStatus(s string) (domain.TaskStatus, error) {
	status, err := domain.ParseTaskStatus(s)
	if err != nil {
		return "", NewValidationError("status", err.Error())
	}
	return status, nil
}

func ParseTaskPriority(s string) (domain.TaskPriority, error) {
	priority, err := domain.ParseTaskPriority(s)
	if err != nil {
		return "", NewValidationError("priority", err.Error())
	}
	return priority, nil
}

func ParseTaskDuration(n int) (domain.TaskDuration, error) {
	duration, err := domain.ParseTaskDuration(n)
	if err != nil {
		return 0, NewValidationError("estimated_duration", err.Error())
	}
	return duration, nil
}

func ValidateTaskTitle(title string) error {
	if title == "" {
		return NewValidationError("title", "task title cannot be empty")
	}
	return nil
}

func ValidateTaskProjectID(projectID string) error {
	if projectID == "" {
		return NewValidationError("project_id", "task project_id cannot be empty")
	}
	return nil
}
