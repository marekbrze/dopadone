package tui

import (
	"context"
	"database/sql"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/converter"
	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

func LoadAreasCmd(repo db.Querier) tea.Cmd {
	return func() tea.Msg {
		areas, err := repo.ListAreas(context.Background())
		if err != nil {
			return AreasLoadedMsg{Err: err}
		}
		domainAreas := make([]domain.Area, len(areas))
		for i, a := range areas {
			domainAreas[i] = converter.DbListAreasRowToDomain(a)
		}
		return AreasLoadedMsg{Areas: domainAreas}
	}
}

func LoadSubareasCmd(repo db.Querier, areaID string) tea.Cmd {
	return func() tea.Msg {
		subareas, err := repo.ListSubareasByArea(context.Background(), areaID)
		if err != nil {
			return SubareasLoadedMsg{Err: err}
		}
		domainSubareas := make([]domain.Subarea, len(subareas))
		for i, s := range subareas {
			domainSubareas[i] = converter.DbSubareaToDomain(s)
		}
		return SubareasLoadedMsg{Subareas: domainSubareas}
	}
}

func LoadProjectsCmd(repo db.Querier, subareaID *string) tea.Cmd {
	return func() tea.Msg {
		allProjects, err := repo.ListAllProjects(context.Background())
		if err != nil {
			return ProjectsLoadedMsg{Err: err}
		}

		var filteredProjects []domain.Project
		if subareaID != nil {
			projectMap := make(map[string]db.Project)
			for _, p := range allProjects {
				projectMap[p.ID] = p
			}

			for _, p := range allProjects {
				if belongsToSubarea(p, *subareaID, projectMap) {
					filteredProjects = append(filteredProjects, converter.DbProjectToDomain(p))
				}
			}
		} else {
			filteredProjects = make([]domain.Project, len(allProjects))
			for i, p := range allProjects {
				filteredProjects[i] = converter.DbProjectToDomain(p)
			}
		}

		return ProjectsLoadedMsg{Projects: filteredProjects}
	}
}

func belongsToSubarea(project db.Project, subareaID string, projectMap map[string]db.Project) bool {
	if project.SubareaID.Valid && project.SubareaID.String == subareaID {
		return true
	}

	if project.ParentID.Valid {
		if parent, exists := projectMap[project.ParentID.String]; exists {
			return belongsToSubarea(parent, subareaID, projectMap)
		}
	}

	return false
}

func LoadTasksCmd(repo db.Querier, projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := repo.ListTasksByProject(context.Background(), projectID)
		if err != nil {
			return TasksLoadedMsg{Err: err}
		}
		domainTasks := make([]domain.Task, len(tasks))
		for i, t := range tasks {
			domainTasks[i] = converter.DbTaskToDomain(t)
		}
		return TasksLoadedMsg{Tasks: domainTasks}
	}
}

func CreateSubareaCmd(repo db.Querier, name string, areaID string) tea.Cmd {
	return func() tea.Msg {
		subarea, err := domain.NewSubarea(name, areaID, "")
		if err != nil {
			return SubareaCreatedMsg{Err: err}
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

		created, err := repo.CreateSubarea(context.Background(), params)
		if err != nil {
			return SubareaCreatedMsg{Err: err}
		}

		return SubareaCreatedMsg{Subarea: converter.DbSubareaToDomain(created)}
	}
}

func CreateProjectCmd(repo db.Querier, name string, parentID *string, subareaID *string) tea.Cmd {
	return func() tea.Msg {
		params := domain.NewProjectParams{
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

		project, err := domain.NewProject(params)
		if err != nil {
			return ProjectCreatedMsg{Err: err}
		}

		dbParams := db.CreateProjectParams{
			ID:          project.ID,
			Name:        project.Name,
			Description: sql.NullString{String: project.Description, Valid: project.Description != ""},
			Goal:        sql.NullString{String: project.Goal, Valid: project.Goal != ""},
			Status:      string(project.Status),
			Priority:    string(project.Priority),
			Progress:    int64(project.Progress),
			Deadline:    nil,
			Color:       sql.NullString{String: string(project.Color), Valid: project.Color != ""},
			ParentID:    sql.NullString{String: "", Valid: false},
			SubareaID:   sql.NullString{String: "", Valid: false},
			Position:    int64(project.Position),
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
			CompletedAt: nil,
			DeletedAt:   nil,
		}

		if project.ParentID != nil {
			dbParams.ParentID = sql.NullString{String: *project.ParentID, Valid: true}
		}
		if project.SubareaID != nil {
			dbParams.SubareaID = sql.NullString{String: *project.SubareaID, Valid: true}
		}

		created, err := repo.CreateProject(context.Background(), dbParams)
		if err != nil {
			return ProjectCreatedMsg{Err: err}
		}

		return ProjectCreatedMsg{Project: converter.DbProjectToDomain(created)}
	}
}

func CreateTaskCmd(repo db.Querier, title string, projectID string) tea.Cmd {
	return func() tea.Msg {
		params := domain.NewTaskParams{
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

		task, err := domain.NewTask(params)
		if err != nil {
			return TaskCreatedMsg{Err: err}
		}

		dbParams := db.CreateTaskParams{
			ID:                task.ID,
			ProjectID:         task.ProjectID,
			Title:             task.Title,
			Description:       sql.NullString{String: task.Description, Valid: task.Description != ""},
			StartDate:         nil,
			Deadline:          nil,
			Priority:          string(task.Priority),
			Context:           sql.NullString{String: task.Context, Valid: task.Context != ""},
			EstimatedDuration: sql.NullInt64{Int64: int64(task.EstimatedDuration), Valid: task.EstimatedDuration != 0},
			Status:            string(task.Status),
			IsNext:            0,
			CreatedAt:         task.CreatedAt,
			UpdatedAt:         task.UpdatedAt,
			DeletedAt:         nil,
		}

		if task.IsNext {
			dbParams.IsNext = 1
		}

		created, err := repo.CreateTask(context.Background(), dbParams)
		if err != nil {
			return TaskCreatedMsg{Err: err}
		}

		return TaskCreatedMsg{Task: converter.DbTaskToDomain(created)}
	}
}

func CreateAreaCmd(repo db.Querier, name string, color domain.Color) tea.Cmd {
	return func() tea.Msg {
		areas, err := repo.ListAreas(context.Background())
		if err != nil {
			return AreaCreatedMsg{Err: err}
		}
		nextSortOrder := len(areas)

		area, err := domain.NewArea(name, color, nextSortOrder)
		if err != nil {
			return AreaCreatedMsg{Err: err}
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

		row, err := repo.CreateArea(context.Background(), params)
		if err != nil {
			return AreaCreatedMsg{Err: err}
		}

		return AreaCreatedMsg{Area: domain.Area{
			ID:        row.ID,
			Name:      row.Name,
			Color:     domain.Color(row.Color.String),
			SortOrder: int(row.SortOrder),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}}
	}
}

func UpdateAreaCmd(repo db.Querier, id string, name string, color domain.Color) tea.Cmd {
	return func() tea.Msg {
		params := db.UpdateAreaParams{
			ID:        id,
			Name:      name,
			Color:     sql.NullString{String: string(color), Valid: color != ""},
			UpdatedAt: time.Now(),
		}

		row, err := repo.UpdateArea(context.Background(), params)
		if err != nil {
			return AreaUpdatedMsg{Err: err}
		}

		return AreaUpdatedMsg{Area: domain.Area{
			ID:        row.ID,
			Name:      row.Name,
			Color:     domain.Color(row.Color.String),
			SortOrder: int(row.SortOrder),
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
		}}
	}
}

func DeleteAreaCmd(repo db.Querier, id string, hard bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		if hard {
			err = repo.DeleteTasksByProject(context.Background(), id)
			if err != nil {
				return AreaDeletedMsg{AreaID: id, Hard: hard, Err: err}
			}
			err = repo.DeleteProjectsBySubarea(context.Background(), id)
			if err != nil {
				return AreaDeletedMsg{AreaID: id, Hard: hard, Err: err}
			}
			err = repo.DeleteSubareasByArea(context.Background(), id)
			if err != nil {
				return AreaDeletedMsg{AreaID: id, Hard: hard, Err: err}
			}
			err = repo.HardDeleteArea(context.Background(), id)
		} else {
			_, err = repo.SoftDeleteArea(context.Background(), db.SoftDeleteAreaParams{
				ID:        id,
				DeletedAt: time.Now(),
			})
		}

		return AreaDeletedMsg{AreaID: id, Hard: hard, Err: err}
	}
}

func ReorderAreasCmd(repo db.Querier, areaIDs []string) tea.Cmd {
	return func() tea.Msg {
		for i, id := range areaIDs {
			err := repo.UpdateAreaSortOrder(context.Background(), db.UpdateAreaSortOrderParams{
				ID:        id,
				SortOrder: int64(i),
				UpdatedAt: time.Now(),
			})
			if err != nil {
				return AreasReorderedMsg{Err: err}
			}
		}
		return AreasReorderedMsg{}
	}
}

func LoadAreaStatsCmd(repo db.Querier, areaID string) tea.Cmd {
	return func() tea.Msg {
		subareas, err := repo.CountSubareasByArea(context.Background(), areaID)
		if err != nil {
			return AreaStatsLoadedMsg{Err: err}
		}

		projects, err := repo.CountProjectsByArea(context.Background(), areaID)
		if err != nil {
			return AreaStatsLoadedMsg{Err: err}
		}

		tasks, err := repo.CountTasksByArea(context.Background(), areaID)
		if err != nil {
			return AreaStatsLoadedMsg{Err: err}
		}

		return AreaStatsLoadedMsg{
			Stats: struct {
				Subareas int64
				Projects int64
				Tasks    int64
			}{
				Subareas: subareas,
				Projects: projects,
				Tasks:    tasks,
			},
		}
	}
}
