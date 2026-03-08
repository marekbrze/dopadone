package tui

import (
	"testing"

	"github.com/marekbrze/dopadone/internal/domain"
)

func TestLoadAreasMsg(t *testing.T) {
	msg := LoadAreasMsg{}
	if msg != (LoadAreasMsg{}) {
		t.Error("LoadAreasMsg should be an empty struct")
	}
}

func TestAreasLoadedMsg(t *testing.T) {
	areas := []domain.Area{
		{ID: "1", Name: "Area 1"},
		{ID: "2", Name: "Area 2"},
	}
	msg := AreasLoadedMsg{Areas: areas}

	if len(msg.Areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(msg.Areas))
	}
	if msg.Areas[0].Name != "Area 1" {
		t.Errorf("Expected first area name 'Area 1', got %s", msg.Areas[0].Name)
	}
}

func TestAreasLoadedMsgWithError(t *testing.T) {
	msg := AreasLoadedMsg{Err: domain.ErrAreaNameEmpty}

	if msg.Err == nil {
		t.Error("Expected error to be set")
	}
	if msg.Err.Error() != "area name cannot be empty" {
		t.Errorf("Unexpected error message: %s", msg.Err.Error())
	}
}

func TestLoadSubareasMsg(t *testing.T) {
	msg := LoadSubareasMsg{AreaID: "area-1"}

	if msg.AreaID != "area-1" {
		t.Errorf("Expected AreaID 'area-1', got %s", msg.AreaID)
	}
}

func TestSubareasLoadedMsg(t *testing.T) {
	subareas := []domain.Subarea{
		{ID: "1", Name: "Subarea 1", AreaID: "area-1"},
	}
	msg := SubareasLoadedMsg{Subareas: subareas}

	if len(msg.Subareas) != 1 {
		t.Errorf("Expected 1 subarea, got %d", len(msg.Subareas))
	}
}

func TestLoadProjectsMsg(t *testing.T) {
	subareaID := "subarea-1"
	msg := LoadProjectsMsg{SubareaID: subareaID}

	if msg.SubareaID != "subarea-1" {
		t.Errorf("Expected SubareaID 'subarea-1', got %s", msg.SubareaID)
	}
}

func TestProjectsLoadedMsg(t *testing.T) {
	projects := []domain.Project{
		{ID: "1", Name: "Project 1"},
	}
	msg := ProjectsLoadedMsg{Projects: projects}

	if len(msg.Projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(msg.Projects))
	}
}

func TestLoadTasksMsg(t *testing.T) {
	msg := LoadTasksMsg{ProjectID: "project-1"}

	if msg.ProjectID != "project-1" {
		t.Errorf("Expected ProjectID 'project-1', got %s", msg.ProjectID)
	}
}

func TestTasksLoadedMsg(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Title: "Task 1", ProjectID: "project-1"},
	}
	msg := TasksLoadedMsg{Tasks: tasks}

	if len(msg.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(msg.Tasks))
	}
}
