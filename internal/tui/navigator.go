package tui

import (
	tea "github.com/charmbracelet/bubbletea"
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
	if len(m.tasks) == 0 {
		return
	}
	if m.selectedTaskIndex == 0 {
		m.selectedTaskIndex = len(m.tasks) - 1
	} else {
		m.selectedTaskIndex--
	}
}

func (m *Model) navigateTasksDown() {
	if len(m.tasks) == 0 {
		return
	}
	if m.selectedTaskIndex >= len(m.tasks)-1 {
		m.selectedTaskIndex = 0
	} else {
		m.selectedTaskIndex++
	}
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
