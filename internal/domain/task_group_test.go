package domain

import (
	"testing"
)

func TestNewGroupedTasks(t *testing.T) {
	tests := []struct {
		name            string
		tasks           []Task
		parentProjectID string
		projectNames    map[string]string
		wantDirectCount int
		wantGroupCount  int
		wantTotalCount  int
		validateFunc    func(t *testing.T, g *GroupedTasks)
	}{
		{
			name:            "empty tasks",
			tasks:           nil,
			parentProjectID: "proj-1",
			wantDirectCount: 0,
			wantGroupCount:  0,
			wantTotalCount:  0,
		},
		{
			name: "direct tasks only",
			tasks: []Task{
				{ID: "t1", ProjectID: "proj-1", Title: "Task 1"},
				{ID: "t2", ProjectID: "proj-1", Title: "Task 2"},
			},
			parentProjectID: "proj-1",
			wantDirectCount: 2,
			wantGroupCount:  0,
			wantTotalCount:  2,
		},
		{
			name: "subproject tasks only",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-1", Title: "Sub Task 1"},
				{ID: "t2", ProjectID: "sub-2", Title: "Sub Task 2"},
			},
			parentProjectID: "proj-1",
			projectNames:    map[string]string{"sub-1": "Subproject 1", "sub-2": "Subproject 2"},
			wantDirectCount: 0,
			wantGroupCount:  2,
			wantTotalCount:  2,
		},
		{
			name: "mixed direct and subproject",
			tasks: []Task{
				{ID: "t1", ProjectID: "proj-1", Title: "Direct 1"},
				{ID: "t2", ProjectID: "sub-1", Title: "Sub 1"},
				{ID: "t3", ProjectID: "proj-1", Title: "Direct 2"},
				{ID: "t4", ProjectID: "sub-2", Title: "Sub 2"},
			},
			parentProjectID: "proj-1",
			projectNames:    map[string]string{"sub-1": "Sub 1", "sub-2": "Sub 2"},
			wantDirectCount: 2,
			wantGroupCount:  2,
			wantTotalCount:  4,
		},
		{
			name: "task order preservation within groups",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-1", Title: "First"},
				{ID: "t2", ProjectID: "sub-1", Title: "Second"},
				{ID: "t3", ProjectID: "sub-1", Title: "Third"},
			},
			parentProjectID: "proj-1",
			projectNames:    map[string]string{"sub-1": "Sub 1"},
			wantGroupCount:  1,
			wantTotalCount:  3,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups) != 1 {
					t.Fatalf("expected 1 group, got %d", len(g.Groups))
				}
				tasks := g.Groups[0].Tasks
				if tasks[0].Title != "First" || tasks[1].Title != "Second" || tasks[2].Title != "Third" {
					t.Error("task order not preserved")
				}
			},
		},
		{
			name: "group order preservation",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-A", Title: "A"},
				{ID: "t2", ProjectID: "sub-B", Title: "B"},
				{ID: "t3", ProjectID: "sub-C", Title: "C"},
			},
			parentProjectID: "proj-1",
			projectNames:    map[string]string{"sub-A": "A", "sub-B": "B", "sub-C": "C"},
			wantGroupCount:  3,
			wantTotalCount:  3,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if g.Groups[0].ProjectID != "sub-A" ||
					g.Groups[1].ProjectID != "sub-B" ||
					g.Groups[2].ProjectID != "sub-C" {
					t.Error("group order not preserved")
				}
			},
		},
		{
			name: "default expanded state",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-1", Title: "Task"},
			},
			parentProjectID: "proj-1",
			projectNames:    map[string]string{"sub-1": "Sub 1"},
			wantGroupCount:  1,
			wantTotalCount:  1,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				for _, group := range g.Groups {
					if !group.IsExpanded {
						t.Error("expected IsExpanded to be true by default")
					}
				}
			},
		},
		{
			name: "missing project name",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-1", Title: "Task"},
			},
			parentProjectID: "proj-1",
			projectNames:    nil,
			wantGroupCount:  1,
			wantTotalCount:  1,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups) != 1 {
					t.Fatalf("expected 1 group, got %d", len(g.Groups))
				}
				if g.Groups[0].ProjectName != "Unknown Project" {
					t.Errorf("expected 'Unknown Project', got %s", g.Groups[0].ProjectName)
				}
			},
		},
		{
			name: "empty parent project ID",
			tasks: []Task{
				{ID: "t1", ProjectID: "sub-1", Title: "Task 1"},
				{ID: "t2", ProjectID: "sub-2", Title: "Task 2"},
			},
			parentProjectID: "",
			wantDirectCount: 0,
			wantGroupCount:  2,
			wantTotalCount:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGroupedTasks(tt.tasks, tt.parentProjectID, tt.projectNames)

			if len(g.DirectTasks) != tt.wantDirectCount {
				t.Errorf("DirectTasks count = %d, want %d", len(g.DirectTasks), tt.wantDirectCount)
			}

			if len(g.Groups) != tt.wantGroupCount {
				t.Errorf("Groups count = %d, want %d", len(g.Groups), tt.wantGroupCount)
			}

			if g.TotalCount != tt.wantTotalCount {
				t.Errorf("TotalCount = %d, want %d", g.TotalCount, tt.wantTotalCount)
			}

			if tt.validateFunc != nil {
				tt.validateFunc(t, g)
			}
		})
	}
}

func TestAddTask(t *testing.T) {
	tests := []struct {
		name         string
		initial      *GroupedTasks
		addTask      Task
		validateFunc func(t *testing.T, g *GroupedTasks)
	}{
		{
			name:    "add to direct tasks",
			initial: &GroupedTasks{DirectTasks: []Task{}, Groups: []TaskGroup{}, parentProjectID: "proj-1"},
			addTask: Task{ID: "t1", ProjectID: "proj-1", Title: "Direct Task"},
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.DirectTasks) != 1 {
					t.Error("task not added to DirectTasks")
				}
				if g.TotalCount != 1 {
					t.Error("TotalCount not incremented")
				}
			},
		},
		{
			name: "add to existing group",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", Tasks: []Task{{ID: "t1"}}},
				},
				parentProjectID: "proj-1",
			},
			addTask: Task{ID: "t2", ProjectID: "sub-1", Title: "Group Task"},
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups[0].Tasks) != 2 {
					t.Error("task not added to existing group")
				}
			},
		},
		{
			name:    "create new group",
			initial: &GroupedTasks{Groups: []TaskGroup{}, parentProjectID: "proj-1"},
			addTask: Task{ID: "t1", ProjectID: "new-sub", Title: "New Group Task"},
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups) != 1 {
					t.Error("new group not created")
				}
				if g.Groups[0].ProjectID != "new-sub" {
					t.Error("wrong project ID for new group")
				}
				if !g.Groups[0].IsExpanded {
					t.Error("new group should be expanded by default")
				}
			},
		},
		{
			name:    "add task with empty project id",
			initial: &GroupedTasks{DirectTasks: []Task{}, Groups: []TaskGroup{}, parentProjectID: "proj-1"},
			addTask: Task{ID: "t1", ProjectID: "", Title: "Empty Project Task"},
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.DirectTasks) != 1 {
					t.Error("task with empty project ID should be added to DirectTasks")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.AddTask(tt.addTask)
			if tt.validateFunc != nil {
				tt.validateFunc(t, tt.initial)
			}
		})
	}
}

func TestRemoveTask(t *testing.T) {
	tests := []struct {
		name         string
		initial      *GroupedTasks
		removeID     string
		wantFound    bool
		validateFunc func(t *testing.T, g *GroupedTasks)
	}{
		{
			name: "remove from direct tasks",
			initial: &GroupedTasks{
				DirectTasks: []Task{{ID: "t1"}, {ID: "t2"}},
				TotalCount:  2,
			},
			removeID:  "t1",
			wantFound: true,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.DirectTasks) != 1 {
					t.Error("task not removed")
				}
				if g.TotalCount != 1 {
					t.Error("TotalCount not decremented")
				}
			},
		},
		{
			name: "remove from group",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", Tasks: []Task{{ID: "t1"}, {ID: "t2"}}},
				},
				TotalCount: 2,
			},
			removeID:  "t1",
			wantFound: true,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups[0].Tasks) != 1 {
					t.Error("task not removed from group")
				}
			},
		},
		{
			name: "task not found",
			initial: &GroupedTasks{
				DirectTasks: []Task{{ID: "t1"}},
				TotalCount:  1,
			},
			removeID:  "nonexistent",
			wantFound: false,
		},
		{
			name: "remove last task from group",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", Tasks: []Task{{ID: "t1"}}},
				},
				TotalCount: 1,
			},
			removeID:  "t1",
			wantFound: true,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if len(g.Groups[0].Tasks) != 0 {
					t.Error("task not removed")
				}
				if len(g.Groups) != 1 {
					t.Error("group should remain even when empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.initial.RemoveTask(tt.removeID)
			if found != tt.wantFound {
				t.Errorf("RemoveTask() = %v, want %v", found, tt.wantFound)
			}
			if tt.validateFunc != nil {
				tt.validateFunc(t, tt.initial)
			}
		})
	}
}

func TestToggleGroup(t *testing.T) {
	tests := []struct {
		name         string
		initial      *GroupedTasks
		toggleID     string
		wantFound    bool
		validateFunc func(t *testing.T, g *GroupedTasks)
	}{
		{
			name: "toggle expanded to collapsed",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", IsExpanded: true},
				},
			},
			toggleID:  "sub-1",
			wantFound: true,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if g.Groups[0].IsExpanded {
					t.Error("group should be collapsed")
				}
			},
		},
		{
			name: "toggle collapsed to expanded",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", IsExpanded: false},
				},
			},
			toggleID:  "sub-1",
			wantFound: true,
			validateFunc: func(t *testing.T, g *GroupedTasks) {
				if !g.Groups[0].IsExpanded {
					t.Error("group should be expanded")
				}
			},
		},
		{
			name: "group not found",
			initial: &GroupedTasks{
				Groups: []TaskGroup{},
			},
			toggleID:  "nonexistent",
			wantFound: false,
		},
		{
			name: "empty project id",
			initial: &GroupedTasks{
				Groups: []TaskGroup{
					{ProjectID: "sub-1", IsExpanded: true},
				},
			},
			toggleID:  "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.initial.ToggleGroup(tt.toggleID)
			if found != tt.wantFound {
				t.Errorf("ToggleGroup() = %v, want %v", found, tt.wantFound)
			}
			if tt.validateFunc != nil {
				tt.validateFunc(t, tt.initial)
			}
		})
	}
}

func TestClear(t *testing.T) {
	tests := []struct {
		name    string
		initial *GroupedTasks
	}{
		{
			name: "clear populated tasks",
			initial: &GroupedTasks{
				DirectTasks: []Task{{ID: "t1"}},
				Groups: []TaskGroup{
					{ProjectID: "sub-1", Tasks: []Task{{ID: "t2"}}},
				},
				TotalCount: 2,
			},
		},
		{
			name:    "clear empty tasks",
			initial: &GroupedTasks{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.Clear()

			if len(tt.initial.DirectTasks) != 0 {
				t.Error("DirectTasks should be empty after clear")
			}

			if len(tt.initial.Groups) != 0 {
				t.Error("Groups should be empty after clear")
			}

			if tt.initial.TotalCount != 0 {
				t.Error("TotalCount should be 0 after clear")
			}
		})
	}
}

func TestIsParentProject(t *testing.T) {
	g := &GroupedTasks{parentProjectID: "proj-1"}

	if !g.isParentProject("proj-1") {
		t.Error("should return true for parent project")
	}

	if g.isParentProject("sub-1") {
		t.Error("should return false for non-parent project")
	}
}
