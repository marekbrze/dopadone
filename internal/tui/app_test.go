package tui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/modal"
	"github.com/example/projectdb/internal/tui/tree"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel(nil)

	if model.repo != nil {
		t.Error("Expected repo to be nil")
	}
	if model.focus != FocusSubareas {
		t.Errorf("Expected initial focus to be FocusSubareas, got %v", model.focus)
	}
	if len(model.areas) != 0 {
		t.Error("Expected areas to be empty")
	}
	if len(model.subareas) != 0 {
		t.Error("Expected subareas to be empty")
	}
	if len(model.projects) != 0 {
		t.Error("Expected projects to be empty")
	}
	if len(model.tasks) != 0 {
		t.Error("Expected tasks to be empty")
	}
}

func TestModelInitWithNilRepo(t *testing.T) {
	model := InitialModel(nil)
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil when repo is nil")
	}
}

func TestModelUpdateAreasLoaded(t *testing.T) {
	model := InitialModel(nil)
	areas := []domain.Area{
		{ID: "1", Name: "Area 1"},
		{ID: "2", Name: "Area 2"},
	}

	msg := AreasLoadedMsg{Areas: areas}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.areas) != 2 {
		t.Errorf("Expected 2 areas, got %d", len(m.areas))
	}
	if m.isLoadingAreas {
		t.Error("Expected isLoadingAreas to be false after loading")
	}
}

func TestModelUpdateAreasLoadedWithError(t *testing.T) {
	model := InitialModel(nil)

	msg := AreasLoadedMsg{Err: errors.New("database error")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.areas) != 0 {
		t.Errorf("Expected 0 areas on error, got %d", len(m.areas))
	}
	if m.isLoadingAreas {
		t.Error("Expected isLoadingAreas to be false after error")
	}
}

func TestModelUpdateSubareasLoaded(t *testing.T) {
	model := InitialModel(nil)
	subareas := []domain.Subarea{
		{ID: "1", Name: "Subarea 1", AreaID: "area-1"},
	}

	msg := SubareasLoadedMsg{Subareas: subareas}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.subareas) != 1 {
		t.Errorf("Expected 1 subarea, got %d", len(m.subareas))
	}
	if m.isLoadingSubareas {
		t.Error("Expected isLoadingSubareas to be false after loading")
	}
}

func TestModelUpdateProjectsLoaded(t *testing.T) {
	model := InitialModel(nil)
	projects := []domain.Project{
		{ID: "1", Name: "Project 1"},
	}

	msg := ProjectsLoadedMsg{Projects: projects}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.projects) != 1 {
		t.Errorf("Expected 1 project, got %d", len(m.projects))
	}
	if m.isLoadingProjects {
		t.Error("Expected isLoadingProjects to be false after loading")
	}
}

func TestModelUpdateTasksLoaded(t *testing.T) {
	model := InitialModel(nil)
	tasks := []domain.Task{
		{ID: "1", Title: "Task 1", ProjectID: "project-1"},
	}

	msg := TasksLoadedMsg{Tasks: tasks}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(m.tasks))
	}
	if m.isLoadingTasks {
		t.Error("Expected isLoadingTasks to be false after loading")
	}
}

func TestModelUpdateKeyPress(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		initial  FocusColumn
		expected FocusColumn
	}{
		{"left from subareas", "left", FocusSubareas, FocusTasks},
		{"right from subareas", "right", FocusSubareas, FocusProjects},
		{"tab from subareas", "tab", FocusSubareas, FocusProjects},
		{"left from projects", "left", FocusProjects, FocusSubareas},
		{"right from projects", "right", FocusProjects, FocusTasks},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel(nil)
			model.focus = tt.initial

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			updatedModel, _ := model.Update(msg)
			m := updatedModel.(Model)

			if m.focus != tt.expected {
				t.Errorf("Expected focus %v, got %v", tt.expected, m.focus)
			}
		})
	}
}

func TestRenderSubareasEmpty(t *testing.T) {
	model := InitialModel(nil)
	result := model.renderSubareas()

	if result != EmptyStateNoSubareas {
		t.Errorf("Expected empty state message, got %s", result)
	}
}

func TestRenderProjectsEmpty(t *testing.T) {
	model := InitialModel(nil)
	result := model.renderProjects()

	if result != EmptyStateNoProjects {
		t.Errorf("Expected empty state message, got %s", result)
	}
}

func TestRenderTasksEmpty(t *testing.T) {
	model := InitialModel(nil)
	result := model.renderTasks()

	if result != EmptyStateNoTasks {
		t.Errorf("Expected empty state message, got %s", result)
	}
}

func TestRenderSubareasWithSelection(t *testing.T) {
	model := InitialModel(nil)
	model.subareas = []domain.Subarea{
		{ID: "1", Name: "Subarea 1"},
		{ID: "2", Name: "Subarea 2"},
	}
	model.selectedSubareaIndex = 0

	result := model.renderSubareas()

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestJoinLines(t *testing.T) {
	lines := []string{"line1", "line2", "line3"}
	result := joinLines(lines)

	expected := "line1\nline2\nline3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestJoinLinesEmpty(t *testing.T) {
	lines := []string{}
	result := joinLines(lines)

	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestModelViewInit(t *testing.T) {
	model := InitialModel(nil)
	model.ready = true
	result := model.View()

	if result == "" {
		t.Error("Expected non-empty view")
	}
	if result == "\n  Initializing..." {
		t.Error("View should not show initializing message when ready")
	}
}

func TestModelViewNotReady(t *testing.T) {
	model := InitialModel(nil)
	model.ready = false
	result := model.View()

	if result != "\n  Initializing..." {
		t.Errorf("Expected initializing message, got '%s'", result)
	}
}

func TestRenderSubareasLoading(t *testing.T) {
	model := InitialModel(nil)
	model.isLoadingSubareas = true
	result := model.renderSubareas()

	if result == "" {
		t.Error("Expected non-empty loading message")
	}
}

func TestRenderProjectsLoading(t *testing.T) {
	model := InitialModel(nil)
	model.isLoadingProjects = true
	result := model.renderProjects()

	if result == "" {
		t.Error("Expected non-empty loading message")
	}
}

func TestRenderTasksLoading(t *testing.T) {
	model := InitialModel(nil)
	model.isLoadingTasks = true
	result := model.renderTasks()

	if result == "" {
		t.Error("Expected non-empty loading message")
	}
}

func TestRenderProjectsWithData(t *testing.T) {
	model := InitialModel(nil)
	model.projects = []domain.Project{
		{ID: "1", Name: "Project 1"},
		{ID: "2", Name: "Project 2"},
	}
	model.selectedProjectIndex = 1

	result := model.renderProjects()

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestRenderTasksWithData(t *testing.T) {
	model := InitialModel(nil)
	model.tasks = []domain.Task{
		{ID: "1", Title: "Task 1", ProjectID: "p1"},
		{ID: "2", Title: "Task 2", ProjectID: "p1"},
	}
	model.selectedTaskIndex = 0

	result := model.renderTasks()

	if result == "" {
		t.Error("Expected non-empty result")
	}
}

func TestCascadeLoadingAreasToSubareas(t *testing.T) {
	model := InitialModel(nil)
	model.areas = []domain.Area{
		{ID: "area-1", Name: "Area 1"},
	}

	msg := AreasLoadedMsg{Areas: model.areas}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.selectedAreaIndex != 0 {
		t.Errorf("Expected selectedAreaIndex to be 0, got %d", m.selectedAreaIndex)
	}
}

func TestCascadeLoadingSubareasToProjects(t *testing.T) {
	model := InitialModel(nil)
	model.subareas = []domain.Subarea{
		{ID: "subarea-1", Name: "Subarea 1", AreaID: "area-1"},
	}

	msg := SubareasLoadedMsg{Subareas: model.subareas}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.selectedSubareaIndex != 0 {
		t.Errorf("Expected selectedSubareaIndex to be 0, got %d", m.selectedSubareaIndex)
	}
}

func TestCascadeLoadingProjectsToTasks(t *testing.T) {
	model := InitialModel(nil)
	model.projects = []domain.Project{
		{ID: "project-1", Name: "Project 1"},
	}

	msg := ProjectsLoadedMsg{Projects: model.projects}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.selectedProjectIndex != 0 {
		t.Errorf("Expected selectedProjectIndex to be 0, got %d", m.selectedProjectIndex)
	}
}

func TestModelUpdateWindowSize(t *testing.T) {
	model := InitialModel(nil)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 100 {
		t.Errorf("Expected width 100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("Expected height 50, got %d", m.height)
	}
	if !m.ready {
		t.Error("Expected ready to be true")
	}
}

func TestModelUpdateQuitCommand(t *testing.T) {
	model := InitialModel(nil)

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}
	updatedModel, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Expected quit command")
	}
	_ = updatedModel
}

func TestModelUpdateCtrlC(t *testing.T) {
	model := InitialModel(nil)

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(msg)

	if cmd == nil {
		t.Error("Expected quit command for Ctrl+C")
	}
	_ = updatedModel
}

func TestModelUpdateSubareasLoadedWithError(t *testing.T) {
	model := InitialModel(nil)

	msg := SubareasLoadedMsg{Err: errors.New("database error")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.subareas) != 0 {
		t.Errorf("Expected 0 subareas on error, got %d", len(m.subareas))
	}
	if m.isLoadingSubareas {
		t.Error("Expected isLoadingSubareas to be false after error")
	}
}

func TestModelUpdateProjectsLoadedWithError(t *testing.T) {
	model := InitialModel(nil)

	msg := ProjectsLoadedMsg{Err: errors.New("database error")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.projects) != 0 {
		t.Errorf("Expected 0 projects on error, got %d", len(m.projects))
	}
	if m.isLoadingProjects {
		t.Error("Expected isLoadingProjects to be false after error")
	}
}

func TestModelUpdateTasksLoadedWithError(t *testing.T) {
	model := InitialModel(nil)

	msg := TasksLoadedMsg{Err: errors.New("database error")}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if len(m.tasks) != 0 {
		t.Errorf("Expected 0 tasks on error, got %d", len(m.tasks))
	}
	if m.isLoadingTasks {
		t.Error("Expected isLoadingTasks to be false after error")
	}
}

func TestGetParentContextForSubprojects(t *testing.T) {
	model := InitialModel(nil)
	model.focus = FocusProjects
	model.subareas = []domain.Subarea{
		{ID: "subarea-1", Name: "Test Subarea"},
	}
	model.selectedSubareaIndex = 0

	t.Run("returns_subproject_when_project_selected", func(t *testing.T) {
		model.selectedProjectID = "project-1"
		model.projectTree = &tree.TreeNode{
			ID:   "project-1",
			Name: "Test Project",
		}

		parentName, entityType, parentID, subareaID := model.getParentContext()

		if parentName != "Test Project" {
			t.Errorf("Expected parent name 'Test Project', got '%s'", parentName)
		}
		if entityType != modal.EntityTypeSubproject {
			t.Errorf("Expected entity type Subproject, got %s", entityType)
		}
		if parentID != "project-1" {
			t.Errorf("Expected parent ID 'project-1', got '%s'", parentID)
		}
		if subareaID != nil {
			t.Error("Expected subareaID to be nil for subprojects")
		}
	})

	t.Run("returns_project_when_no_project_selected", func(t *testing.T) {
		model.selectedProjectID = ""
		model.selectedSubareaIndex = 0

		parentName, entityType, parentID, subareaID := model.getParentContext()

		if parentName != "Test Subarea" {
			t.Errorf("Expected parent name 'Test Subarea', got '%s'", parentName)
		}
		if entityType != modal.EntityTypeProject {
			t.Errorf("Expected entity type Project, got %s", entityType)
		}
		if parentID != "" {
			t.Errorf("Expected parent ID to be empty, got '%s'", parentID)
		}
		if subareaID == nil || *subareaID != "subarea-1" {
			t.Error("Expected subareaID to be 'subarea-1'")
		}
	})
}
