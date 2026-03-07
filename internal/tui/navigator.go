package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/dopadone/internal/domain"
)

func (m *Model) NavigateUp(column FocusColumn) {
	switch column {
	case FocusSubareas:
		m.navigateSubareasUp()
	case FocusProjects:
		m.navigateTreeUp()
	case FocusTasks:
		m.navigateTasksUp()
	}
}

func (m *Model) NavigateDown(column FocusColumn) {
	switch column {
	case FocusSubareas:
		m.navigateSubareasDown()
	case FocusProjects:
		m.navigateTreeDown()
	case FocusTasks:
		m.navigateTasksDown()
	}
}

func (m *Model) NavigateDownWithLoad(column FocusColumn) (tea.Model, tea.Cmd) {
	if column == FocusSubareas {
		prevIndex := m.selectedSubareaIndex
		m.navigateSubareasDown()
		if m.selectedSubareaIndex != prevIndex && len(m.subareas) > 0 {
			m.isLoadingProjects = true
			m.projects = nil
			m.tasks = nil
			m.projectTree = nil
			m.selectedProjectID = ""
			return m, LoadProjectsCmd(m.projectSvc, &m.subareas[m.selectedSubareaIndex].ID)
		}
		return m, nil
	}
	if column == FocusProjects {
		prevID := m.selectedProjectID
		m.navigateTreeDown()
		if m.selectedProjectID != prevID && m.selectedProjectID != "" {
			m.isLoadingTasks = true
			m.tasks = nil
			m.selectedTaskIndex = 0
			return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
		}
		return m, nil
	}
	m.NavigateDown(column)
	return m, nil
}

func (m *Model) NavigateUpWithLoad(column FocusColumn) (tea.Model, tea.Cmd) {
	if column == FocusSubareas {
		prevIndex := m.selectedSubareaIndex
		m.navigateSubareasUp()
		if m.selectedSubareaIndex != prevIndex && len(m.subareas) > 0 {
			m.isLoadingProjects = true
			m.projects = nil
			m.tasks = nil
			m.projectTree = nil
			m.selectedProjectID = ""
			return m, LoadProjectsCmd(m.projectSvc, &m.subareas[m.selectedSubareaIndex].ID)
		}
		return m, nil
	}
	if column == FocusProjects {
		prevID := m.selectedProjectID
		m.navigateTreeUp()
		if m.selectedProjectID != prevID && m.selectedProjectID != "" {
			m.isLoadingTasks = true
			m.tasks = nil
			m.selectedTaskIndex = 0
			return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
		}
		return m, nil
	}
	m.NavigateUp(column)
	return m, nil
}

func (m *Model) navigateSubareasUp() {
	if len(m.subareas) == 0 {
		return
	}
	if m.selectedSubareaIndex == 0 {
		m.selectedSubareaIndex = len(m.subareas) - 1
	} else {
		m.selectedSubareaIndex--
	}
}

func (m *Model) navigateSubareasDown() {
	if len(m.subareas) == 0 {
		return
	}
	if m.selectedSubareaIndex >= len(m.subareas)-1 {
		m.selectedSubareaIndex = 0
	} else {
		m.selectedSubareaIndex++
	}
}

func (m *Model) navigateTasksUp() {
	totalLines := m.getTotalTaskLines()
	if totalLines == 0 {
		return
	}
	if m.selectedTaskIndex == 0 {
		m.selectedTaskIndex = totalLines - 1
	} else {
		m.selectedTaskIndex--
	}
}

func (m *Model) navigateTasksDown() {
	totalLines := m.getTotalTaskLines()
	if totalLines == 0 {
		return
	}
	if m.selectedTaskIndex >= totalLines-1 {
		m.selectedTaskIndex = 0
	} else {
		m.selectedTaskIndex++
	}
}

func (m *Model) getTotalTaskLines() int {
	if m.groupedTasks == nil {
		return len(m.tasks)
	}

	count := 0

	count += len(m.groupedTasks.DirectTasks)

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		count++
	}

	for _, group := range m.groupedTasks.Groups {
		count++
		if m.expandedTaskGroups != nil && m.expandedTaskGroups[group.ProjectID] {
			count += len(group.Tasks)
		}
	}

	return count
}

func (m *Model) isLineGroupHeader(lineIndex int) bool {
	if m.groupedTasks == nil {
		return false
	}

	currentLine := 0

	currentLine += len(m.groupedTasks.DirectTasks)
	if lineIndex < currentLine {
		return false
	}

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		if lineIndex == currentLine {
			return false
		}
		currentLine++
	}

	for _, group := range m.groupedTasks.Groups {
		if lineIndex == currentLine {
			return true
		}
		currentLine++

		if m.expandedTaskGroups != nil && m.expandedTaskGroups[group.ProjectID] {
			currentLine += len(group.Tasks)
			if lineIndex < currentLine {
				return false
			}
		}
	}

	return false
}

func (m *Model) getGroupAtLine(lineIndex int) *domain.TaskGroup {
	if m.groupedTasks == nil {
		return nil
	}

	currentLine := 0

	currentLine += len(m.groupedTasks.DirectTasks)
	if lineIndex < currentLine {
		return nil
	}

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		if lineIndex == currentLine {
			return nil
		}
		currentLine++
	}

	for i := range m.groupedTasks.Groups {
		if lineIndex == currentLine {
			return &m.groupedTasks.Groups[i]
		}
		currentLine++

		if m.expandedTaskGroups != nil && m.expandedTaskGroups[m.groupedTasks.Groups[i].ProjectID] {
			currentLine += len(m.groupedTasks.Groups[i].Tasks)
			if lineIndex < currentLine {
				return nil
			}
		}
	}

	return nil
}

func (m *Model) getTaskAtLine(lineIndex int) *domain.Task {
	if m.groupedTasks == nil {
		if lineIndex >= 0 && lineIndex < len(m.tasks) {
			return &m.tasks[lineIndex]
		}
		return nil
	}

	currentLine := 0

	for i := range m.groupedTasks.DirectTasks {
		if lineIndex == currentLine {
			return &m.groupedTasks.DirectTasks[i]
		}
		currentLine++
	}

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		if lineIndex == currentLine {
			return nil
		}
		currentLine++
	}

	for gi := range m.groupedTasks.Groups {
		if lineIndex == currentLine {
			return nil
		}
		currentLine++

		if m.expandedTaskGroups != nil && m.expandedTaskGroups[m.groupedTasks.Groups[gi].ProjectID] {
			for ti := range m.groupedTasks.Groups[gi].Tasks {
				if lineIndex == currentLine {
					return &m.groupedTasks.Groups[gi].Tasks[ti]
				}
				currentLine++
			}
		}
	}

	return nil
}

func (m *Model) getGroupHeaderLineForGroup(groupID string) int {
	if m.groupedTasks == nil {
		return -1
	}

	currentLine := 0

	currentLine += len(m.groupedTasks.DirectTasks)

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		currentLine++
	}

	for _, group := range m.groupedTasks.Groups {
		if group.ProjectID == groupID {
			return currentLine
		}
		currentLine++

		if m.expandedTaskGroups != nil && m.expandedTaskGroups[group.ProjectID] {
			currentLine += len(group.Tasks)
		}
	}

	return -1
}

func (m *Model) SwitchToPreviousArea() tea.Cmd {
	if len(m.areas) == 0 {
		return nil
	}
	m.SaveCurrentAreaState()
	if m.selectedAreaIndex == 0 {
		m.selectedAreaIndex = len(m.areas) - 1
	} else {
		m.selectedAreaIndex--
	}
	areaID := m.areas[m.selectedAreaIndex].ID
	m.RestoreAreaState(areaID)
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)
	m.selectedTab = m.selectedAreaIndex
	return m.loadAreaData(areaID)
}

func (m *Model) SwitchToNextArea() tea.Cmd {
	if len(m.areas) == 0 {
		return nil
	}
	m.SaveCurrentAreaState()
	if m.selectedAreaIndex >= len(m.areas)-1 {
		m.selectedAreaIndex = 0
	} else {
		m.selectedAreaIndex++
	}
	areaID := m.areas[m.selectedAreaIndex].ID
	m.RestoreAreaState(areaID)
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)
	m.selectedTab = m.selectedAreaIndex
	return m.loadAreaData(areaID)
}

func (m *Model) loadAreaData(areaID string) tea.Cmd {
	m.isLoadingSubareas = true
	m.subareas = nil
	m.projects = nil
	m.tasks = nil
	m.projectTree = nil
	return LoadSubareasCmd(m.subareaSvc, areaID)
}
