package tui

import (
	"github.com/example/dopadone/internal/domain"
)

type LoadAreasMsg struct{}

type AreasLoadedMsg struct {
	Areas []domain.Area
	Err   error
}

type LoadSubareasMsg struct {
	AreaID string
}

type SubareasLoadedMsg struct {
	Subareas []domain.Subarea
	Err      error
}

type LoadProjectsMsg struct {
	SubareaID string
}

type ProjectsLoadedMsg struct {
	Projects []domain.Project
	Err      error
}

type LoadTasksMsg struct {
	ProjectID string
}

type TasksLoadedMsg struct {
	Tasks        []domain.Task
	GroupedTasks *domain.GroupedTasks
	Err          error
}

type SubareaCreatedMsg struct {
	Subarea domain.Subarea
	Err     error
}

type ProjectCreatedMsg struct {
	Project domain.Project
	Err     error
}

type TaskCreatedMsg struct {
	Task domain.Task
	Err  error
}

type ToastTickMsg struct{}

type LoadAreaStatsMsg struct {
	AreaID string
}

type AreaStatsLoadedMsg struct {
	Stats struct {
		Subareas int64
		Projects int64
		Tasks    int64
	}
	Err error
}

type AreaCreatedMsg struct {
	Area domain.Area
	Err  error
}

type AreaUpdatedMsg struct {
	Area domain.Area
	Err  error
}

type AreaDeletedMsg struct {
	AreaID string
	Hard   bool
	Err    error
}

type AreasReorderedMsg struct {
	Err error
}

type TaskStatusToggledMsg struct {
	Task           *domain.Task
	OriginalStatus domain.TaskStatus
	TaskIndex      int
	Err            error
}
