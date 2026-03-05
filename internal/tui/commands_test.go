package tui

import (
	"testing"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/mocks"
)

func TestLoadAreasCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mockAreaSvc, _, _, _ := mocks.NewMockServices()
		expectedAreas := []domain.Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
		}
		mocks.SetupMockAreaSuccess(mockAreaSvc, expectedAreas)

		cmd := LoadAreasCmd(mockAreaSvc)
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
		_, mockSubareaSvc, _, _ := mocks.NewMockServices()
		expectedSubareas := []domain.Subarea{
			{ID: "1", Name: "Subarea 1", AreaID: "area-1"},
		}
		mocks.SetupMockSubareaSuccess(mockSubareaSvc, expectedSubareas)

		cmd := LoadSubareasCmd(mockSubareaSvc, "area-1")
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
		_, _, mockProjectSvc, _ := mocks.NewMockServices()
		expectedProjects := []domain.Project{
			{ID: "1", Name: "Project 1", SubareaID: ptrToString("subarea-1")},
		}
		mocks.SetupMockProjectSuccess(mockProjectSvc, expectedProjects)

		subareaID := "subarea-1"
		cmd := LoadProjectsCmd(mockProjectSvc, &subareaID)
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
		_, _, _, mockTaskSvc := mocks.NewMockServices()
		expectedTasks := []domain.Task{
			{ID: "1", Title: "Task 1", ProjectID: "project-1"},
		}
		mocks.SetupMockTaskSuccess(mockTaskSvc, expectedTasks)

		cmd := LoadTasksCmd(mockTaskSvc, "project-1")
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

func ptrToString(s string) *string {
	return &s
}
