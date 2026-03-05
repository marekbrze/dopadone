package tui

import (
	"context"
	"database/sql"
	"testing"

	"github.com/example/projectdb/internal/db"
)

type MockQuerier struct {
	areas    []db.ListAreasRow
	subareas []db.Subarea
	projects []db.Project
	tasks    []db.Task
	err      error
}

func (m *MockQuerier) CreateArea(ctx context.Context, arg db.CreateAreaParams) (db.CreateAreaRow, error) {
	return db.CreateAreaRow{}, m.err
}

func (m *MockQuerier) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	return db.Project{}, m.err
}

func (m *MockQuerier) CreateSubarea(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, m.err
}

func (m *MockQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return db.Task{}, m.err
}

func (m *MockQuerier) GetAreaByID(ctx context.Context, id string) (db.GetAreaByIDRow, error) {
	return db.GetAreaByIDRow{}, m.err
}

func (m *MockQuerier) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return db.Project{}, m.err
}

func (m *MockQuerier) GetProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	return m.projects, m.err
}

func (m *MockQuerier) GetSubareaByID(ctx context.Context, id string) (db.Subarea, error) {
	return db.Subarea{}, m.err
}

func (m *MockQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
	return db.Task{}, m.err
}

func (m *MockQuerier) ListAreas(ctx context.Context) ([]db.ListAreasRow, error) {
	return m.areas, m.err
}

func (m *MockQuerier) ListAllProjects(ctx context.Context) ([]db.Project, error) {
	return m.projects, m.err
}

func (m *MockQuerier) ListNextTasks(ctx context.Context) ([]db.Task, error) {
	return m.tasks, m.err
}

func (m *MockQuerier) ListProjectsByParent(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
	return m.projects, m.err
}

func (m *MockQuerier) ListProjectsBySubarea(ctx context.Context, subareaID sql.NullString) ([]db.Project, error) {
	return m.projects, m.err
}

func (m *MockQuerier) ListSubareasByArea(ctx context.Context, areaID string) ([]db.Subarea, error) {
	return m.subareas, m.err
}

func (m *MockQuerier) ListTasksByPriority(ctx context.Context, priority string) ([]db.Task, error) {
	return m.tasks, m.err
}

func (m *MockQuerier) ListTasksByProject(ctx context.Context, projectID string) ([]db.Task, error) {
	return m.tasks, m.err
}

func (m *MockQuerier) ListTasksByStatus(ctx context.Context, status string) ([]db.Task, error) {
	return m.tasks, m.err
}

func (m *MockQuerier) SoftDeleteArea(ctx context.Context, arg db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
	return db.SoftDeleteAreaRow{}, m.err
}

func (m *MockQuerier) SoftDeleteProject(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
	return db.Project{}, m.err
}

func (m *MockQuerier) SoftDeleteSubarea(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, m.err
}

func (m *MockQuerier) SoftDeleteTask(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error) {
	return db.Task{}, m.err
}

func (m *MockQuerier) ToggleIsNext(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
	return db.Task{}, m.err
}

func (m *MockQuerier) UpdateArea(ctx context.Context, arg db.UpdateAreaParams) (db.UpdateAreaRow, error) {
	return db.UpdateAreaRow{}, m.err
}

func (m *MockQuerier) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	return db.Project{}, m.err
}

func (m *MockQuerier) UpdateSubarea(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, m.err
}

func (m *MockQuerier) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	return db.Task{}, m.err
}

func (m *MockQuerier) CountProjectsByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, m.err
}

func (m *MockQuerier) CountSubareasByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, m.err
}

func (m *MockQuerier) CountTasksByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, m.err
}

func (m *MockQuerier) DeleteProjectsBySubarea(ctx context.Context, areaID string) error {
	return m.err
}

func (m *MockQuerier) DeleteSubareasByArea(ctx context.Context, areaID string) error {
	return m.err
}

func (m *MockQuerier) DeleteTasksByProject(ctx context.Context, areaID string) error {
	return m.err
}

func (m *MockQuerier) HardDeleteArea(ctx context.Context, id string) error {
	return m.err
}

func (m *MockQuerier) UpdateAreaSortOrder(ctx context.Context, arg db.UpdateAreaSortOrderParams) error {
	return m.err
}

func TestLoadAreasCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mock := &MockQuerier{
			areas: []db.ListAreasRow{
				{ID: "1", Name: "Area 1"},
				{ID: "2", Name: "Area 2"},
			},
		}

		cmd := LoadAreasCmd(mock)
		msg := cmd()

		loaded, ok := msg.(AreasLoadedMsg)
		if !ok {
			t.Fatal("Expected AreasLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Areas) != 2 {
			t.Errorf("Expected 2 areas, got %d", len(loaded.Areas))
		}
	})
}

func TestLoadSubareasCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mock := &MockQuerier{
			subareas: []db.Subarea{
				{ID: "1", Name: "Subarea 1", AreaID: "area-1"},
			},
		}

		cmd := LoadSubareasCmd(mock, "area-1")
		msg := cmd()

		loaded, ok := msg.(SubareasLoadedMsg)
		if !ok {
			t.Fatal("Expected SubareasLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Subareas) != 1 {
			t.Errorf("Expected 1 subarea, got %d", len(loaded.Subareas))
		}
	})
}

func TestLoadProjectsCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mock := &MockQuerier{
			projects: []db.Project{
				{ID: "1", Name: "Project 1", SubareaID: sql.NullString{String: "subarea-1", Valid: true}},
			},
		}

		subareaID := "subarea-1"
		cmd := LoadProjectsCmd(mock, &subareaID)
		msg := cmd()

		loaded, ok := msg.(ProjectsLoadedMsg)
		if !ok {
			t.Fatal("Expected ProjectsLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Projects) != 1 {
			t.Errorf("Expected 1 project, got %d", len(loaded.Projects))
		}
	})
}

func TestLoadTasksCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mock := &MockQuerier{
			tasks: []db.Task{
				{ID: "1", Title: "Task 1", ProjectID: "project-1"},
			},
		}

		cmd := LoadTasksCmd(mock, "project-1")
		msg := cmd()

		loaded, ok := msg.(TasksLoadedMsg)
		if !ok {
			t.Fatal("Expected TasksLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Tasks) != 1 {
			t.Errorf("Expected 1 task, got %d", len(loaded.Tasks))
		}
	})
}
