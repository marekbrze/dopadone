package tui

import (
	"testing"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/tree"
)

func TestNavigateSubareasUpWrap(t *testing.T) {
	m := Model{
		subareas: []domain.Subarea{
			{ID: "s1", Name: "Subarea 1"},
			{ID: "s2", Name: "Subarea 2"},
			{ID: "s3", Name: "Subarea 3"},
		},
		selectedSubareaIndex: 0,
	}

	m.navigateSubareasUp()
	if m.selectedSubareaIndex != 2 {
		t.Errorf("expected index 2 (wrap to last), got %d", m.selectedSubareaIndex)
	}
}

func TestNavigateSubareasDownWrap(t *testing.T) {
	m := Model{
		subareas: []domain.Subarea{
			{ID: "s1", Name: "Subarea 1"},
			{ID: "s2", Name: "Subarea 2"},
			{ID: "s3", Name: "Subarea 3"},
		},
		selectedSubareaIndex: 2,
	}

	m.navigateSubareasDown()
	if m.selectedSubareaIndex != 0 {
		t.Errorf("expected index 0 (wrap to first), got %d", m.selectedSubareaIndex)
	}
}

func TestNavigateSubareasEmpty(t *testing.T) {
	m := Model{
		subareas:             []domain.Subarea{},
		selectedSubareaIndex: 0,
	}

	m.navigateSubareasUp()
	if m.selectedSubareaIndex != 0 {
		t.Errorf("expected no change on empty list, got %d", m.selectedSubareaIndex)
	}

	m.navigateSubareasDown()
	if m.selectedSubareaIndex != 0 {
		t.Errorf("expected no change on empty list, got %d", m.selectedSubareaIndex)
	}
}

func TestNavigateTasksUpWrap(t *testing.T) {
	m := Model{
		tasks: []domain.Task{
			{ID: "t1", Title: "Task 1"},
			{ID: "t2", Title: "Task 2"},
			{ID: "t3", Title: "Task 3"},
		},
		selectedTaskIndex: 0,
	}

	m.navigateTasksUp()
	if m.selectedTaskIndex != 2 {
		t.Errorf("expected index 2 (wrap to last), got %d", m.selectedTaskIndex)
	}
}

func TestNavigateTasksDownWrap(t *testing.T) {
	m := Model{
		tasks: []domain.Task{
			{ID: "t1", Title: "Task 1"},
			{ID: "t2", Title: "Task 2"},
			{ID: "t3", Title: "Task 3"},
		},
		selectedTaskIndex: 2,
	}

	m.navigateTasksDown()
	if m.selectedTaskIndex != 0 {
		t.Errorf("expected index 0 (wrap to first), got %d", m.selectedTaskIndex)
	}
}

func TestNavigateTasksEmpty(t *testing.T) {
	m := Model{
		tasks:             []domain.Task{},
		selectedTaskIndex: 0,
	}

	m.navigateTasksUp()
	if m.selectedTaskIndex != 0 {
		t.Errorf("expected no change on empty list, got %d", m.selectedTaskIndex)
	}

	m.navigateTasksDown()
	if m.selectedTaskIndex != 0 {
		t.Errorf("expected no change on empty list, got %d", m.selectedTaskIndex)
	}
}

func TestNavigateTreeUp(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	child1 := tree.NewTreeNode("p1", "Project 1", nil)
	child2 := tree.NewTreeNode("p2", "Project 2", nil)
	root.AddChild(child1)
	root.AddChild(child2)

	m := Model{
		projectTree:       root,
		selectedProjectID: "p2",
	}

	m.navigateTreeUp()
	if m.selectedProjectID != "p1" {
		t.Errorf("expected p1, got %s", m.selectedProjectID)
	}
}

func TestNavigateTreeDown(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	child1 := tree.NewTreeNode("p1", "Project 1", nil)
	child2 := tree.NewTreeNode("p2", "Project 2", nil)
	root.AddChild(child1)
	root.AddChild(child2)

	m := Model{
		projectTree:       root,
		selectedProjectID: "p1",
	}

	m.navigateTreeDown()
	if m.selectedProjectID != "p2" {
		t.Errorf("expected p2, got %s", m.selectedProjectID)
	}
}

func TestNavigateTreeUpWrap(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	child1 := tree.NewTreeNode("p1", "Project 1", nil)
	child2 := tree.NewTreeNode("p2", "Project 2", nil)
	root.AddChild(child1)
	root.AddChild(child2)

	m := Model{
		projectTree:       root,
		selectedProjectID: "p1",
	}

	m.navigateTreeUp()
	if m.selectedProjectID != "p2" {
		t.Errorf("expected p2 (wrap to last), got %s", m.selectedProjectID)
	}
}

func TestNavigateTreeDownWrap(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	child1 := tree.NewTreeNode("p1", "Project 1", nil)
	child2 := tree.NewTreeNode("p2", "Project 2", nil)
	root.AddChild(child1)
	root.AddChild(child2)

	m := Model{
		projectTree:       root,
		selectedProjectID: "p2",
	}

	m.navigateTreeDown()
	if m.selectedProjectID != "p1" {
		t.Errorf("expected p1 (wrap to first), got %s", m.selectedProjectID)
	}
}

func TestNavigateTreeEmpty(t *testing.T) {
	m := Model{
		projectTree:       nil,
		selectedProjectID: "",
	}

	m.navigateTreeUp()
	if m.selectedProjectID != "" {
		t.Errorf("expected no change on nil tree, got %s", m.selectedProjectID)
	}

	m.navigateTreeDown()
	if m.selectedProjectID != "" {
		t.Errorf("expected no change on nil tree, got %s", m.selectedProjectID)
	}
}

func TestNavigateTreeWithCollapsedNodes(t *testing.T) {
	root := tree.NewTreeNode("", "root", nil)
	parent := tree.NewTreeNode("parent", "Parent", nil)
	child := tree.NewTreeNode("child", "Child", nil)
	parent.AddChild(child)
	root.AddChild(parent)

	parent.IsExpanded = false

	m := Model{
		projectTree:       root,
		selectedProjectID: "parent",
	}

	m.navigateTreeDown()
	if m.selectedProjectID == "child" {
		t.Error("child should be skipped when parent is collapsed")
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		model    Model
		column   FocusColumn
		expected bool
	}{
		{
			name:     "empty subareas",
			model:    Model{subareas: []domain.Subarea{}},
			column:   FocusSubareas,
			expected: true,
		},
		{
			name:     "non-empty subareas",
			model:    Model{subareas: []domain.Subarea{{ID: "1"}}},
			column:   FocusSubareas,
			expected: false,
		},
		{
			name:     "empty tasks",
			model:    Model{tasks: []domain.Task{}},
			column:   FocusTasks,
			expected: true,
		},
		{
			name:     "non-empty tasks",
			model:    Model{tasks: []domain.Task{{ID: "1"}}},
			column:   FocusTasks,
			expected: false,
		},
		{
			name:     "nil project tree",
			model:    Model{projectTree: nil},
			column:   FocusProjects,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.model.IsEmpty(tt.column); got != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNavigateUpDispatch(t *testing.T) {
	m := Model{
		subareas: []domain.Subarea{
			{ID: "s1", Name: "Subarea 1"},
			{ID: "s2", Name: "Subarea 2"},
		},
		selectedSubareaIndex: 1,
	}

	m.NavigateUp(FocusSubareas)
	if m.selectedSubareaIndex != 0 {
		t.Errorf("expected 0, got %d", m.selectedSubareaIndex)
	}
}

func TestNavigateDownDispatch(t *testing.T) {
	m := Model{
		tasks: []domain.Task{
			{ID: "t1", Title: "Task 1"},
			{ID: "t2", Title: "Task 2"},
		},
		selectedTaskIndex: 0,
	}

	m.NavigateDown(FocusTasks)
	if m.selectedTaskIndex != 1 {
		t.Errorf("expected 1, got %d", m.selectedTaskIndex)
	}
}

func TestSingleItemNoWrap(t *testing.T) {
	m := Model{
		subareas: []domain.Subarea{
			{ID: "s1", Name: "Only One"},
		},
		selectedSubareaIndex: 0,
	}

	m.navigateSubareasDown()
	if m.selectedSubareaIndex != 0 {
		t.Errorf("single item should stay at 0, got %d", m.selectedSubareaIndex)
	}

	m.navigateSubareasUp()
	if m.selectedSubareaIndex != 0 {
		t.Errorf("single item should stay at 0, got %d", m.selectedSubareaIndex)
	}
}
