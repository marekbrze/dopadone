package domain

// TaskGroup represents a collection of tasks belonging to a subproject.
// It tracks the subproject metadata and expansion state for TUI rendering.
type TaskGroup struct {
	ProjectID   string
	ProjectName string
	Tasks       []Task
	IsExpanded  bool
}

// GroupedTasks organizes tasks by subproject with group metadata.
// Direct tasks belong to the selected project, while groups contain
// tasks from nested subprojects.
type GroupedTasks struct {
	DirectTasks []Task
	Groups      []TaskGroup
	TotalCount  int

	parentProjectID string
}

// NewGroupedTasks creates a GroupedTasks instance from a flat list of tasks.
// Tasks are grouped by their ProjectID, with direct tasks (matching parentProjectID)
// separated into DirectTasks field.
//
// Parameters:
//   - tasks: Flat list of tasks to group (can be nil/empty)
//   - parentProjectID: ID of the selected project for direct tasks
//   - projectNames: Map of ProjectID → ProjectName for display (can be nil)
//
// Returns a pointer to GroupedTasks with:
//   - Direct tasks in DirectTasks field
//   - Subproject tasks in Groups field
//   - TotalCount reflecting total number of tasks
//   - All groups default to IsExpanded=true
//
// Task order is preserved within groups (append order).
// Group order is preserved by discovery order from tasks slice.
// Missing project names default to "Unknown Project".
func NewGroupedTasks(tasks []Task, parentProjectID string, projectNames map[string]string) *GroupedTasks {
	if tasks == nil {
		tasks = []Task{}
	}
	if projectNames == nil {
		projectNames = make(map[string]string)
	}

	grouped := &GroupedTasks{
		DirectTasks:     []Task{},
		Groups:          []TaskGroup{},
		TotalCount:      0,
		parentProjectID: parentProjectID,
	}

	tasksByProject := make(map[string][]Task)
	groupOrder := []string{}

	for _, task := range tasks {
		if _, exists := tasksByProject[task.ProjectID]; !exists {
			groupOrder = append(groupOrder, task.ProjectID)
		}
		tasksByProject[task.ProjectID] = append(tasksByProject[task.ProjectID], task)
	}

	if directTasks, exists := tasksByProject[parentProjectID]; exists {
		grouped.DirectTasks = directTasks
		delete(tasksByProject, parentProjectID)
	}

	for _, projectID := range groupOrder {
		if projectID == parentProjectID {
			continue
		}

		projectTasks, exists := tasksByProject[projectID]
		if !exists {
			continue
		}

		groupName := projectNames[projectID]
		if groupName == "" {
			groupName = "Unknown Project"
		}

		grouped.Groups = append(grouped.Groups, TaskGroup{
			ProjectID:   projectID,
			ProjectName: groupName,
			Tasks:       projectTasks,
			IsExpanded:  true,
		})
	}

	grouped.TotalCount = len(grouped.DirectTasks)
	for _, g := range grouped.Groups {
		grouped.TotalCount += len(g.Tasks)
	}

	return grouped
}

// AddTask adds a task to the appropriate group or DirectTasks.
// If task.ProjectID matches parent project → DirectTasks
// If group exists for task.ProjectID → existing group
// Otherwise → creates new group with IsExpanded=true
func (g *GroupedTasks) AddTask(task Task) {
	if task.ProjectID == "" || task.ProjectID == g.parentProjectID {
		g.DirectTasks = append(g.DirectTasks, task)
	} else {
		found := false
		for i := range g.Groups {
			if g.Groups[i].ProjectID == task.ProjectID {
				g.Groups[i].Tasks = append(g.Groups[i].Tasks, task)
				found = true
				break
			}
		}
		if !found {
			g.Groups = append(g.Groups, TaskGroup{
				ProjectID:   task.ProjectID,
				ProjectName: "Unknown Project",
				Tasks:       []Task{task},
				IsExpanded:  true,
			})
		}
	}
	g.TotalCount++
}

// RemoveTask removes a task by ID from any group or DirectTasks.
// Returns true if found and removed, false if not found.
func (g *GroupedTasks) RemoveTask(taskID string) bool {
	for i, task := range g.DirectTasks {
		if task.ID == taskID {
			g.DirectTasks = append(g.DirectTasks[:i], g.DirectTasks[i+1:]...)
			g.TotalCount--
			return true
		}
	}

	for i := range g.Groups {
		for j, task := range g.Groups[i].Tasks {
			if task.ID == taskID {
				g.Groups[i].Tasks = append(g.Groups[i].Tasks[:j], g.Groups[i].Tasks[j+1:]...)
				g.TotalCount--
				return true
			}
		}
	}

	return false
}

// ToggleGroup toggles the expansion state of a group by projectID.
// Returns true if group found and toggled, false if not found.
func (g *GroupedTasks) ToggleGroup(projectID string) bool {
	for i := range g.Groups {
		if g.Groups[i].ProjectID == projectID {
			g.Groups[i].IsExpanded = !g.Groups[i].IsExpanded
			return true
		}
	}
	return false
}

// Clear resets DirectTasks, Groups, and TotalCount to empty/zero.
func (g *GroupedTasks) Clear() {
	g.DirectTasks = []Task{}
	g.Groups = []TaskGroup{}
	g.TotalCount = 0
}

func (g *GroupedTasks) isParentProject(projectID string) bool {
	return projectID == g.parentProjectID
}

// Flattened returns all tasks as a flat slice (for backward compatibility).
func (g *GroupedTasks) Flattened() []Task {
	tasks := make([]Task, 0, g.TotalCount)
	tasks = append(tasks, g.DirectTasks...)
	for _, group := range g.Groups {
		tasks = append(tasks, group.Tasks...)
	}
	return tasks
}
