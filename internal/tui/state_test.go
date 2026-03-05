package tui

import (
	"testing"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/tree"
)

func TestNewAreaState(t *testing.T) {
	state := NewAreaState()
	if state.SelectedSubareaIndex != 0 {
		t.Errorf("expected 0, got %d", state.SelectedSubareaIndex)
	}
	if state.SelectedProjectIndex != 0 {
		t.Errorf("expected 0, got %d", state.SelectedProjectIndex)
	}
	if state.SelectedTaskIndex != 0 {
		t.Errorf("expected 0, got %d", state.SelectedTaskIndex)
	}
	if state.ExpandedProjects == nil {
		t.Error("expected ExpandedProjects to be initialized")
	}
}

func TestGetAreaStateCreatesNew(t *testing.T) {
	m := Model{
		areaStates: make(map[string]*AreaState),
	}

	state := m.GetAreaState("area-1")
	if state == nil {
		t.Fatal("expected state to be created")
	}
	if m.areaStates["area-1"] == nil {
		t.Error("expected state to be stored in map")
	}
}

func TestGetAreaStateReturnsExisting(t *testing.T) {
	existingState := &AreaState{
		SelectedSubareaIndex: 5,
	}
	m := Model{
		areaStates: map[string]*AreaState{
			"area-1": existingState,
		},
	}

	state := m.GetAreaState("area-1")
	if state != existingState {
		t.Error("expected existing state to be returned")
	}
}

func TestSaveCurrentAreaState(t *testing.T) {
	m := Model{
		areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
		},
		areaStates:           make(map[string]*AreaState),
		selectedAreaIndex:    0,
		selectedSubareaIndex: 2,
		selectedProjectIndex: 3,
		selectedTaskIndex:    1,
	}

	m.SaveCurrentAreaState()

	state := m.areaStates["area-1"]
	if state == nil {
		t.Fatal("expected state to be saved")
	}
	if state.SelectedSubareaIndex != 2 {
		t.Errorf("expected 2, got %d", state.SelectedSubareaIndex)
	}
	if state.SelectedProjectIndex != 3 {
		t.Errorf("expected 3, got %d", state.SelectedProjectIndex)
	}
	if state.SelectedTaskIndex != 1 {
		t.Errorf("expected 1, got %d", state.SelectedTaskIndex)
	}
}

func TestRestoreAreaState(t *testing.T) {
	m := Model{
		areaStates: map[string]*AreaState{
			"area-1": {
				SelectedSubareaIndex: 4,
				SelectedProjectIndex: 2,
				SelectedTaskIndex:    6,
			},
		},
	}

	m.RestoreAreaState("area-1")

	if m.selectedSubareaIndex != 4 {
		t.Errorf("expected 4, got %d", m.selectedSubareaIndex)
	}
	if m.selectedProjectIndex != 2 {
		t.Errorf("expected 2, got %d", m.selectedProjectIndex)
	}
	if m.selectedTaskIndex != 6 {
		t.Errorf("expected 6, got %d", m.selectedTaskIndex)
	}
}

func TestSaveRestoreAreaStateRoundTrip(t *testing.T) {
	m := Model{
		areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
			{ID: "area-2", Name: "Area 2"},
		},
		areaStates:           make(map[string]*AreaState),
		selectedAreaIndex:    0,
		selectedSubareaIndex: 3,
		selectedProjectIndex: 7,
		selectedTaskIndex:    2,
	}

	m.SaveCurrentAreaState()

	m.selectedSubareaIndex = 0
	m.selectedProjectIndex = 0
	m.selectedTaskIndex = 0

	m.RestoreAreaState("area-1")

	if m.selectedSubareaIndex != 3 {
		t.Errorf("expected 3, got %d", m.selectedSubareaIndex)
	}
	if m.selectedProjectIndex != 7 {
		t.Errorf("expected 7, got %d", m.selectedProjectIndex)
	}
	if m.selectedTaskIndex != 2 {
		t.Errorf("expected 2, got %d", m.selectedTaskIndex)
	}
}

func TestAreaStateIsolation(t *testing.T) {
	m := Model{
		areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
			{ID: "area-2", Name: "Area 2"},
		},
		areaStates:           make(map[string]*AreaState),
		selectedAreaIndex:    0,
		selectedSubareaIndex: 1,
		selectedProjectIndex: 2,
		selectedTaskIndex:    3,
	}

	m.SaveCurrentAreaState()

	m.selectedAreaIndex = 1
	m.selectedSubareaIndex = 5
	m.selectedProjectIndex = 6
	m.selectedTaskIndex = 7
	m.SaveCurrentAreaState()

	state1 := m.areaStates["area-1"]
	state2 := m.areaStates["area-2"]

	if state1.SelectedSubareaIndex != 1 {
		t.Errorf("area-1 subarea: expected 1, got %d", state1.SelectedSubareaIndex)
	}
	if state2.SelectedSubareaIndex != 5 {
		t.Errorf("area-2 subarea: expected 5, got %d", state2.SelectedSubareaIndex)
	}
}

func TestTreeExpandStatePersistence(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	parent := tree.NewTreeNode("parent", "Parent", nil)
	child := tree.NewTreeNode("child", "Child", nil)
	parent.AddChild(child)
	root.AddChild(parent)

	parent.IsExpanded = true

	m := Model{
		areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
		},
		areaStates:        make(map[string]*AreaState),
		selectedAreaIndex: 0,
		projectTree:       root,
	}

	m.SaveCurrentAreaState()

	state := m.areaStates["area-1"]
	if !state.ExpandedProjects["parent"] {
		t.Error("expected parent to be marked as expanded in state")
	}
}

func TestTreeExpandStateRestore(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	parent := tree.NewTreeNode("parent", "Parent", nil)
	child := tree.NewTreeNode("child", "Child", nil)
	parent.AddChild(child)
	root.AddChild(parent)

	parent.IsExpanded = false

	m := Model{
		areaStates: map[string]*AreaState{
			"area-1": {
				ExpandedProjects: map[string]bool{"parent": true},
			},
		},
		projectTree: root,
	}

	m.restoreTreeExpandState(m.areaStates["area-1"])

	if !parent.IsExpanded {
		t.Error("expected parent to be restored as expanded")
	}
}

func TestSaveCurrentAreaStateEmptyAreas(t *testing.T) {
	m := Model{
		areas:             []domain.Area{},
		areaStates:        make(map[string]*AreaState),
		selectedAreaIndex: 0,
	}

	m.SaveCurrentAreaState()

	if len(m.areaStates) != 0 {
		t.Error("expected no state to be saved for empty areas")
	}
}

func TestSaveCurrentAreaStateIndexOutOfRange(t *testing.T) {
	m := Model{
		areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
		},
		areaStates:        make(map[string]*AreaState),
		selectedAreaIndex: 5,
	}

	m.SaveCurrentAreaState()

	if len(m.areaStates) != 0 {
		t.Error("expected no state to be saved for out of range index")
	}
}
