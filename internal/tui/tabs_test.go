package tui

import (
	"testing"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/mocks"
)

func TestUpdateTabsFromAreasEmpty(t *testing.T) {
	tabs := updateTabsFromAreas([]domain.Area{}, 0)
	if len(tabs) != 0 {
		t.Errorf("Expected empty tabs for empty areas, got %d tabs", len(tabs))
	}
}

func TestUpdateTabsFromAreasSingle(t *testing.T) {
	areas := []domain.Area{
		{ID: "area-1", Name: "Personal"},
	}
	tabs := updateTabsFromAreas(areas, 0)

	if len(tabs) != 1 {
		t.Fatalf("Expected 1 tab, got %d", len(tabs))
	}

	if tabs[0].Name != "Personal" {
		t.Errorf("Expected tab name 'Personal', got '%s'", tabs[0].Name)
	}

	if tabs[0].ID != "area-1" {
		t.Errorf("Expected tab ID 'area-1', got '%s'", tabs[0].ID)
	}

	if !tabs[0].IsActive {
		t.Error("Expected tab to be active")
	}
}

func TestUpdateTabsFromAreasMultiple(t *testing.T) {
	areas := []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
		{ID: "area-3", Name: "Side Projects"},
	}
	tabs := updateTabsFromAreas(areas, 1)

	if len(tabs) != 3 {
		t.Fatalf("Expected 3 tabs, got %d", len(tabs))
	}

	expected := []struct {
		name     string
		id       string
		isActive bool
	}{
		{"Personal", "area-1", false},
		{"Work", "area-2", true},
		{"Side Projects", "area-3", false},
	}

	for i, exp := range expected {
		if tabs[i].Name != exp.name {
			t.Errorf("Tab %d: expected name '%s', got '%s'", i, exp.name, tabs[i].Name)
		}
		if tabs[i].ID != exp.id {
			t.Errorf("Tab %d: expected ID '%s', got '%s'", i, exp.id, tabs[i].ID)
		}
		if tabs[i].IsActive != exp.isActive {
			t.Errorf("Tab %d: expected IsActive=%v, got %v", i, exp.isActive, tabs[i].IsActive)
		}
	}
}

func TestUpdateTabsFromAreasFirstSelected(t *testing.T) {
	areas := []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
	}
	tabs := updateTabsFromAreas(areas, 0)

	if len(tabs) != 2 {
		t.Fatalf("Expected 2 tabs, got %d", len(tabs))
	}

	if !tabs[0].IsActive {
		t.Error("Expected first tab to be active")
	}

	if tabs[1].IsActive {
		t.Error("Expected second tab to be inactive")
	}
}

func TestUpdateTabsFromAreasLastSelected(t *testing.T) {
	areas := []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
		{ID: "area-3", Name: "Side Projects"},
	}
	tabs := updateTabsFromAreas(areas, 2)

	if len(tabs) != 3 {
		t.Fatalf("Expected 3 tabs, got %d", len(tabs))
	}

	for i := 0; i < 2; i++ {
		if tabs[i].IsActive {
			t.Errorf("Expected tab %d to be inactive", i)
		}
	}

	if !tabs[2].IsActive {
		t.Error("Expected last tab to be active")
	}
}

func TestAreasLoadedUpdatesTabs(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)
	m.areas = []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
	}

	msg := AreasLoadedMsg{
		Areas: m.areas,
	}

	updatedModel, _ := m.Update(msg)
	model := updatedModel.(Model)

	if len(model.tabs) != 2 {
		t.Fatalf("Expected 2 tabs after AreasLoadedMsg, got %d", len(model.tabs))
	}

	if model.tabs[0].Name != "Personal" {
		t.Errorf("Expected first tab name 'Personal', got '%s'", model.tabs[0].Name)
	}

	if model.tabs[1].Name != "Work" {
		t.Errorf("Expected second tab name 'Work', got '%s'", model.tabs[1].Name)
	}

	if model.selectedTab != 0 {
		t.Errorf("Expected selectedTab to be 0, got %d", model.selectedTab)
	}
}

func TestAreasLoadedUpdatesTabsCorrectSelection(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)
	m.selectedAreaIndex = 1
	m.areas = []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
		{ID: "area-3", Name: "Side Projects"},
	}

	msg := AreasLoadedMsg{
		Areas: m.areas,
	}

	updatedModel, _ := m.Update(msg)
	model := updatedModel.(Model)

	if len(model.tabs) != 3 {
		t.Fatalf("Expected 3 tabs, got %d", len(model.tabs))
	}

	for i, tab := range model.tabs {
		expectedActive := (i == 1)
		if tab.IsActive != expectedActive {
			t.Errorf("Tab %d: expected IsActive=%v, got %v", i, expectedActive, tab.IsActive)
		}
	}

	if model.selectedTab != 1 {
		t.Errorf("Expected selectedTab to be 1, got %d", model.selectedTab)
	}
}

func TestInitialModelHasEmptyTabs(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)

	if len(m.tabs) != 0 {
		t.Errorf("Expected InitialModel to have empty tabs, got %d tabs", len(m.tabs))
	}
}

func TestSwitchToNextAreaUpdatesTabs(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)
	m.areas = []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
		{ID: "area-3", Name: "Side Projects"},
	}
	m.selectedAreaIndex = 0
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)

	cmd := m.SwitchToNextArea()

	if cmd == nil {
		t.Error("Expected switchToNextArea to return a command")
	}

	if m.selectedAreaIndex != 1 {
		t.Errorf("Expected selectedAreaIndex to be 1, got %d", m.selectedAreaIndex)
	}

	if m.selectedTab != 1 {
		t.Errorf("Expected selectedTab to be 1, got %d", m.selectedTab)
	}

	if len(m.tabs) != 3 {
		t.Fatalf("Expected 3 tabs, got %d", len(m.tabs))
	}

	for i, tab := range m.tabs {
		expectedActive := (i == 1)
		if tab.IsActive != expectedActive {
			t.Errorf("Tab %d: expected IsActive=%v, got %v", i, expectedActive, tab.IsActive)
		}
	}
}

func TestSwitchToPreviousAreaUpdatesTabs(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)
	m.areas = []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
		{ID: "area-3", Name: "Side Projects"},
	}
	m.selectedAreaIndex = 2
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)

	cmd := m.SwitchToPreviousArea()

	if cmd == nil {
		t.Error("Expected switchToPreviousArea to return a command")
	}

	if m.selectedAreaIndex != 1 {
		t.Errorf("Expected selectedAreaIndex to be 1, got %d", m.selectedAreaIndex)
	}

	if m.selectedTab != 1 {
		t.Errorf("Expected selectedTab to be 1, got %d", m.selectedTab)
	}

	if len(m.tabs) != 3 {
		t.Fatalf("Expected 3 tabs, got %d", len(m.tabs))
	}

	for i, tab := range m.tabs {
		expectedActive := (i == 1)
		if tab.IsActive != expectedActive {
			t.Errorf("Tab %d: expected IsActive=%v, got %v", i, expectedActive, tab.IsActive)
		}
	}
}

func TestSwitchAreaWrapsTabs(t *testing.T) {
	mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc := mocks.NewMockServices()
	m := InitialModel(mockAreaSvc, mockSubareaSvc, mockProjectSvc, mockTaskSvc)
	m.areas = []domain.Area{
		{ID: "area-1", Name: "Personal"},
		{ID: "area-2", Name: "Work"},
	}
	m.selectedAreaIndex = 1
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)

	_ = m.SwitchToNextArea()

	if m.selectedAreaIndex != 0 {
		t.Errorf("Expected selectedAreaIndex to wrap to 0, got %d", m.selectedAreaIndex)
	}

	if !m.tabs[0].IsActive {
		t.Error("Expected first tab to be active after wrap")
	}

	if m.tabs[1].IsActive {
		t.Error("Expected second tab to be inactive after wrap")
	}

	m.selectedAreaIndex = 0
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)

	_ = m.SwitchToPreviousArea()

	if m.selectedAreaIndex != 1 {
		t.Errorf("Expected selectedAreaIndex to wrap to 1, got %d", m.selectedAreaIndex)
	}

	if m.tabs[0].IsActive {
		t.Error("Expected first tab to be inactive after wrap")
	}

	if !m.tabs[1].IsActive {
		t.Error("Expected second tab to be active after wrap")
	}
}
