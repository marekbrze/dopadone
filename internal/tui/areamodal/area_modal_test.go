package areamodal

import (
	"testing"

	"github.com/marekbrze/dopadone/internal/domain"
)

func TestModal_SetSelectedIndex(t *testing.T) {
	t.Run("sets valid index", func(t *testing.T) {
		areas := []Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
			{ID: "3", Name: "Area 3"},
		}
		m := New(areas)

		m.SetSelectedIndex(1)
		if m.selectedIndex != 1 {
			t.Errorf("expected selectedIndex 1, got %d", m.selectedIndex)
		}

		m.SetSelectedIndex(2)
		if m.selectedIndex != 2 {
			t.Errorf("expected selectedIndex 2, got %d", m.selectedIndex)
		}
	})

	t.Run("ignores negative index", func(t *testing.T) {
		areas := []Area{
			{ID: "1", Name: "Area 1"},
		}
		m := New(areas)
		originalIndex := m.selectedIndex

		m.SetSelectedIndex(-1)
		if m.selectedIndex != originalIndex {
			t.Errorf("expected selectedIndex to remain %d, got %d", originalIndex, m.selectedIndex)
		}
	})

	t.Run("ignores index out of range", func(t *testing.T) {
		areas := []Area{
			{ID: "1", Name: "Area 1"},
		}
		m := New(areas)
		originalIndex := m.selectedIndex

		m.SetSelectedIndex(5)
		if m.selectedIndex != originalIndex {
			t.Errorf("expected selectedIndex to remain %d, got %d", originalIndex, m.selectedIndex)
		}
	})
}

func TestModal_SetupForCreate(t *testing.T) {
	areas := []Area{
		{ID: "1", Name: "Area 1", Color: "#3B82F6"},
	}
	m := New(areas)

	m.input.SetValue("previous value")
	m.colorIndex = 5
	m.mode = ModeEdit

	m.SetupForCreate()

	if m.mode != ModeCreate {
		t.Errorf("expected mode ModeCreate, got %v", m.mode)
	}
	if m.input.Value() != "" {
		t.Errorf("expected empty input value, got %q", m.input.Value())
	}
	if m.colorIndex != 0 {
		t.Errorf("expected colorIndex 0, got %d", m.colorIndex)
	}
}

func TestModal_SetupForEdit(t *testing.T) {
	t.Run("sets up edit mode with area data", func(t *testing.T) {
		areas := []Area{
			{ID: "1", Name: "Area 1", Color: "#10B981"},
			{ID: "2", Name: "Area 2", Color: "#EF4444"},
		}
		m := New(areas)
		m.SetSelectedIndex(1)

		m.SetupForEdit()

		if m.mode != ModeEdit {
			t.Errorf("expected mode ModeEdit, got %v", m.mode)
		}
		if m.editAreaID != "2" {
			t.Errorf("expected editAreaID '2', got %q", m.editAreaID)
		}
		if m.input.Value() != "Area 2" {
			t.Errorf("expected input value 'Area 2', got %q", m.input.Value())
		}

		expectedColorIndex := -1
		for i, c := range PredefinedColors {
			if c == "#EF4444" {
				expectedColorIndex = i
				break
			}
		}
		if m.colorIndex != expectedColorIndex {
			t.Errorf("expected colorIndex %d, got %d", expectedColorIndex, m.colorIndex)
		}
	})

	t.Run("does nothing with empty areas", func(t *testing.T) {
		m := New([]Area{})

		m.SetupForEdit()

		if m.mode != ModeList {
			t.Errorf("expected mode to remain ModeList, got %v", m.mode)
		}
	})

	t.Run("does nothing with invalid selectedIndex", func(t *testing.T) {
		areas := []Area{
			{ID: "1", Name: "Area 1"},
		}
		m := New(areas)
		m.selectedIndex = 5

		m.SetupForEdit()

		if m.mode != ModeList {
			t.Errorf("expected mode to remain ModeList, got %v", m.mode)
		}
	})
}

func TestModal_SetupForDelete(t *testing.T) {
	areas := []Area{
		{ID: "1", Name: "Area 1"},
	}
	m := New(areas)
	m.mode = ModeList
	m.deleteChoice = DeleteChoiceSoft
	m.statsLoaded = true

	m.SetupForDelete()

	if m.mode != ModeDeleteConfirm {
		t.Errorf("expected mode ModeDeleteConfirm, got %v", m.mode)
	}
	if m.deleteChoice != DeleteChoiceNone {
		t.Errorf("expected deleteChoice DeleteChoiceNone, got %v", m.deleteChoice)
	}
	if m.statsLoaded {
		t.Error("expected statsLoaded to be false")
	}
}

func TestModal_SetStats(t *testing.T) {
	m := New([]Area{})

	stats := Stats{Subareas: 5, Projects: 10, Tasks: 25}
	m.SetStats(stats)

	if m.stats.Subareas != 5 {
		t.Errorf("expected Subareas 5, got %d", m.stats.Subareas)
	}
	if m.stats.Projects != 10 {
		t.Errorf("expected Projects 10, got %d", m.stats.Projects)
	}
	if m.stats.Tasks != 25 {
		t.Errorf("expected Tasks 25, got %d", m.stats.Tasks)
	}
	if !m.statsLoaded {
		t.Error("expected statsLoaded to be true")
	}
}

func TestModal_UpdateAreas(t *testing.T) {
	t.Run("updates areas list", func(t *testing.T) {
		m := New([]Area{})

		newAreas := []Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
		}
		m.UpdateAreas(newAreas)

		if len(m.areas) != 2 {
			t.Errorf("expected 2 areas, got %d", len(m.areas))
		}
	})

	t.Run("adjusts selectedIndex if out of bounds", func(t *testing.T) {
		m := New([]Area{})
		m.selectedIndex = 5

		newAreas := []Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
		}
		m.UpdateAreas(newAreas)

		if m.selectedIndex != 1 {
			t.Errorf("expected selectedIndex 1, got %d", m.selectedIndex)
		}
	})

	t.Run("keeps valid selectedIndex", func(t *testing.T) {
		m := New([]Area{})
		m.selectedIndex = 0

		newAreas := []Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
		}
		m.UpdateAreas(newAreas)

		if m.selectedIndex != 0 {
			t.Errorf("expected selectedIndex 0, got %d", m.selectedIndex)
		}
	})
}

func TestModal_New(t *testing.T) {
	areas := []Area{
		{ID: "1", Name: "Area 1", Color: domain.Color("#3B82F6")},
	}

	m := New(areas)

	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.mode != ModeList {
		t.Errorf("expected initial mode ModeList, got %v", m.mode)
	}
	if len(m.areas) != 1 {
		t.Errorf("expected 1 area, got %d", len(m.areas))
	}
	if m.selectedIndex != 0 {
		t.Errorf("expected selectedIndex 0, got %d", m.selectedIndex)
	}
}
