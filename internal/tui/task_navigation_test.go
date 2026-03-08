package tui

import (
	"testing"

	"github.com/marekbrze/dopadone/internal/domain"
)

func TestGetTotalTaskLinesEmpty(t *testing.T) {
	m := Model{
		groupedTasks: nil,
	}

	if total := m.getTotalTaskLines(); total != 0 {
		t.Errorf("expected 0, got %d", total)
	}
}

func TestGetTotalTaskLinesDirectTasks(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
				{ID: "t2", Title: "Task 2"},
			},
		},
	}

	if total := m.getTotalTaskLines(); total != 2 {
		t.Errorf("expected 2, got %d", total)
	}
}

func TestGetTotalTaskLinesWithGroups(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
						{ID: "gt2", Title: "Group Task 2"},
					},
					IsExpanded: true,
				},
				{
					ProjectID:   "proj-2",
					ProjectName: "Project 2",
					Tasks: []domain.Task{
						{ID: "gt3", Title: "Group Task 3"},
						{ID: "gt4", Title: "Group Task 4"},
					},
					IsExpanded: false,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
			"proj-2": false,
		},
	}

	if total := m.getTotalTaskLines(); total != 6 {
		t.Errorf("expected 6, got %d", total)
	}
}

func TestGetGroupAtLine(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
					},
					IsExpanded: true,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
		},
	}

	tests := []struct {
		name      string
		lineIndex int
		expected  string
	}{
		{"Direct task 0", 0, ""},
		{"Separator", 1, ""},
		{"Group header", 2, "proj-1"},
		{"Group task", 3, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := m.getGroupAtLine(tt.lineIndex)
			if tt.expected == "" {
				if group != nil {
					t.Fatalf("expected nil for line %d, got group", tt.lineIndex)
				}
			} else {
				if group == nil {
					t.Fatalf("expected group, got nil")
				}
				if group.ProjectID != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, group.ProjectID)
				}
			}
		})
	}
}

func TestGetTaskAtLine(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
					},
					IsExpanded: true,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
		},
	}

	tests := []struct {
		name      string
		lineIndex int
		expected  string
	}{
		{"Direct task", 0, "t1"},
		{"Separator", 1, ""},
		{"Group header", 2, ""},
		{"Group task", 3, "gt1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := m.getTaskAtLine(tt.lineIndex)
			if tt.expected == "" {
				if task != nil {
					t.Fatalf("expected nil for line %d, got task", tt.lineIndex)
				}
			} else {
				if task == nil {
					t.Fatalf("expected task, got nil")
				}
				if task.ID != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, task.ID)
				}
			}
		})
	}
}

func TestGetTaskAtLineWithCollapsedGroup(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
					},
					IsExpanded: false,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": false,
		},
	}

	task := m.getTaskAtLine(2)
	if task != nil {
		t.Errorf("expected nil for group header line, got task")
	}
}

func TestGetTaskAtLineWithEmptyGroup(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks:       []domain.Task{},
					IsExpanded:  true,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
		},
	}

	task := m.getTaskAtLine(0)
	if task != nil {
		t.Errorf("expected nil for empty group header, got task")
	}
}

func TestNavigationWithGroupHeaders(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
					},
					IsExpanded: true,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
		},
		selectedTaskIndex: 0,
	}

	m.navigateTasksDown()
	if m.selectedTaskIndex != 1 {
		t.Errorf("expected line 1 (header) after task 0, got %d", m.selectedTaskIndex)
	}

	m.navigateTasksDown()
	if m.selectedTaskIndex != 2 {
		t.Errorf("expected line 2 (group task) after header, got %d", m.selectedTaskIndex)
	}
}

func TestSelectionAdjustmentOnCollapse(t *testing.T) {
	m := Model{
		groupedTasks: &domain.GroupedTasks{
			DirectTasks: []domain.Task{
				{ID: "t1", Title: "Task 1"},
			},
			Groups: []domain.TaskGroup{
				{
					ProjectID:   "proj-1",
					ProjectName: "Project 1",
					Tasks: []domain.Task{
						{ID: "gt1", Title: "Group Task 1"},
						{ID: "gt2", Title: "Group Task 2"},
					},
					IsExpanded: true,
				},
			},
		},
		expandedTaskGroups: map[string]bool{
			"proj-1": true,
		},
		selectedTaskIndex: 3,
	}

	if total := m.getTotalTaskLines(); total != 5 {
		t.Errorf("expected 5 lines, got %d", total)
	}

	m.expandedTaskGroups["proj-1"] = false
	m.groupedTasks.Groups[0].IsExpanded = false

	if total := m.getTotalTaskLines(); total != 3 {
		t.Errorf("expected 3 lines after collapse, got %d", total)
	}
}
