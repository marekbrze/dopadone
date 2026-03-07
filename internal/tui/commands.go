package tui

import (
	"context"

	"github.com/charmbracelet/bubbletea"
	"github.com/example/dopadone/internal/domain"
	"github.com/example/dopadone/internal/service"
)

func LoadAreasCmd(areaSvc service.AreaServiceInterface) tea.Cmd {
	return func() tea.Msg {
		areas, err := areaSvc.List(context.Background())
		if err != nil {
			return AreasLoadedMsg{Err: err}
		}
		return AreasLoadedMsg{Areas: areas}
	}
}

func LoadSubareasCmd(subareaSvc service.SubareaServiceInterface, areaID string) tea.Cmd {
	return func() tea.Msg {
		subareas, err := subareaSvc.ListByArea(context.Background(), areaID)
		if err != nil {
			return SubareasLoadedMsg{Err: err}
		}
		return SubareasLoadedMsg{Subareas: subareas}
	}
}

// LoadProjectsCmd loads projects for a subarea using hierarchical retrieval.
// When subareaID is provided, uses ListBySubareaRecursive to include nested projects.
// When subareaID is nil, loads all projects using ListAll.
func LoadProjectsCmd(projectSvc service.ProjectServiceInterface, subareaID *string) tea.Cmd {
	return func() tea.Msg {
		var projects []domain.Project
		var err error

		if subareaID != nil {
			projects, err = projectSvc.ListBySubareaRecursive(context.Background(), *subareaID)
		} else {
			projects, err = projectSvc.ListAll(context.Background())
		}

		if err != nil {
			return ProjectsLoadedMsg{Err: err}
		}

		return ProjectsLoadedMsg{Projects: projects}
	}
}

func LoadTasksCmd(taskSvc service.TaskServiceInterface, projectID string) tea.Cmd {
	return func() tea.Msg {
		groupedTasks, err := taskSvc.GetGroupedTasks(context.Background(), projectID)
		if err != nil {
			return TasksLoadedMsg{Err: err}
		}

		tasks := groupedTasks.Flattened()

		return TasksLoadedMsg{
			Tasks:        tasks,
			GroupedTasks: groupedTasks,
			Err:          nil,
		}
	}
}

func CreateSubareaCmd(subareaSvc service.SubareaServiceInterface, name string, areaID string) tea.Cmd {
	return func() tea.Msg {
		subarea, err := subareaSvc.Create(context.Background(), name, areaID, "")
		if err != nil {
			return SubareaCreatedMsg{Err: err}
		}

		return SubareaCreatedMsg{Subarea: *subarea}
	}
}

func CreateProjectCmd(projectSvc service.ProjectServiceInterface, name string, parentID *string, subareaID *string) tea.Cmd {
	return func() tea.Msg {
		params := service.CreateProjectParams{
			Name:        name,
			Description: "",
			Goal:        "",
			Status:      domain.ProjectStatusActive,
			Priority:    domain.PriorityMedium,
			Progress:    0,
			StartDate:   nil,
			Deadline:    nil,
			Color:       "",
			ParentID:    parentID,
			SubareaID:   subareaID,
			Position:    0,
		}

		project, err := projectSvc.Create(context.Background(), params)
		if err != nil {
			return ProjectCreatedMsg{Err: err}
		}

		return ProjectCreatedMsg{Project: *project}
	}
}

func CreateTaskCmd(taskSvc service.TaskServiceInterface, title string, projectID string) tea.Cmd {
	return func() tea.Msg {
		params := service.CreateTaskParams{
			ProjectID:         projectID,
			Title:             title,
			Description:       "",
			StartDate:         nil,
			Deadline:          nil,
			Priority:          domain.TaskPriorityMedium,
			Context:           "",
			EstimatedDuration: 0,
			Status:            domain.TaskStatusTodo,
			IsNext:            false,
		}

		task, err := taskSvc.Create(context.Background(), params)
		if err != nil {
			return TaskCreatedMsg{Err: err}
		}

		return TaskCreatedMsg{Task: *task}
	}
}

func CreateAreaCmd(areaSvc service.AreaServiceInterface, name string, color domain.Color) tea.Cmd {
	return func() tea.Msg {
		area, err := areaSvc.Create(context.Background(), name, color)
		if err != nil {
			return AreaCreatedMsg{Err: err}
		}

		return AreaCreatedMsg{Area: *area}
	}
}

func UpdateAreaCmd(areaSvc service.AreaServiceInterface, id string, name string, color domain.Color) tea.Cmd {
	return func() tea.Msg {
		area, err := areaSvc.Update(context.Background(), id, name, color)
		if err != nil {
			return AreaUpdatedMsg{Err: err}
		}

		return AreaUpdatedMsg{Area: *area}
	}
}

func DeleteAreaCmd(areaSvc service.AreaServiceInterface, id string, hard bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if hard {
			err = areaSvc.HardDelete(context.Background(), id)
		} else {
			err = areaSvc.SoftDelete(context.Background(), id)
		}

		return AreaDeletedMsg{AreaID: id, Hard: hard, Err: err}
	}
}

func ReorderAreasCmd(areaSvc service.AreaServiceInterface, areaIDs []string) tea.Cmd {
	return func() tea.Msg {
		err := areaSvc.ReorderAll(context.Background(), areaIDs)
		if err != nil {
			return AreasReorderedMsg{Err: err}
		}
		return AreasReorderedMsg{}
	}
}

func LoadAreaStatsCmd(areaSvc service.AreaServiceInterface, areaID string) tea.Cmd {
	return func() tea.Msg {
		stats, err := areaSvc.GetStats(context.Background(), areaID)
		if err != nil {
			return AreaStatsLoadedMsg{Err: err}
		}

		return AreaStatsLoadedMsg{
			Stats: struct {
				Subareas int64
				Projects int64
				Tasks    int64
			}{
				Subareas: stats.SubareaCount,
				Projects: stats.ProjectCount,
				Tasks:    stats.TaskCount,
			},
		}
	}
}

func ToggleTaskStatusCmd(
	taskSvc service.TaskServiceInterface,
	taskID string,
	newStatus domain.TaskStatus,
	originalStatus domain.TaskStatus,
	taskIndex int,
) tea.Cmd {
	return func() tea.Msg {
		task, err := taskSvc.SetStatus(context.Background(), taskID, newStatus)

		return TaskStatusToggledMsg{
			Task:           task,
			OriginalStatus: originalStatus,
			TaskIndex:      taskIndex,
			Err:            err,
		}
	}
}
