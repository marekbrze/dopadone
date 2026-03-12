package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/toast"
	"github.com/marekbrze/dopadone/internal/tui/tree"
	"github.com/marekbrze/dopadone/internal/tui/welcome"
)

func (m *Model) handleAreasLoaded(msg AreasLoadedMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.isLoadingAreas = false
	if msg.Err != nil {
		m.areaLoadError = msg.Err
		m.addToast(toast.NewError("Failed to load areas: " + msg.Err.Error()))
		return m, nil
	}

	if len(msg.Areas) == 0 {
		m.welcomeModal = welcome.New()
		m.isWelcomeOpen = true
		m.welcomeModal, _ = m.welcomeModal.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		return m, nil
	}

	m.areaLoadError = nil
	m.areas = msg.Areas
	m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)
	m.selectedTab = m.selectedAreaIndex
	if len(m.areas) > 0 && m.selectedAreaIndex == 0 {
		m.selectedAreaIndex = 0
		m.isLoadingSubareas = true
		cmds = append(cmds, LoadSubareasCmd(m.subareaSvc, m.areas[0].ID))
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleSubareasLoaded(msg SubareasLoadedMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.isLoadingSubareas = false
	if msg.Err != nil {
		m.subareaLoadError = msg.Err
		m.addToast(toast.NewError("Failed to load subareas: " + msg.Err.Error()))
		return m, nil
	}
	m.subareaLoadError = nil
	m.subareas = msg.Subareas
	if len(m.subareas) > 0 && m.selectedSubareaIndex == 0 {
		m.selectedSubareaIndex = 0
		m.isLoadingProjects = true
		cmds = append(cmds, LoadProjectsCmd(m.projectSvc, &m.subareas[0].ID))
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleProjectsLoaded(msg ProjectsLoadedMsg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.isLoadingProjects = false
	if msg.Err != nil {
		m.projectLoadError = msg.Err
		m.addToast(toast.NewError("Failed to load projects: " + msg.Err.Error()))
		return m, nil
	}
	m.projectLoadError = nil
	m.projects = msg.Projects

	builder := tree.NewBuilder()
	m.projectTree = builder.BuildFromProjects(m.projects)

	if m.projectTree != nil {
		firstNode := tree.GetFirstVisibleNode(m.projectTree)
		if firstNode != nil {
			m.selectedProjectID = firstNode.ID
			m.selectedProjectIndex = 0
			m.isLoadingTasks = true
			m.tasks = nil
			m.selectedTaskIndex = 0
			cmds = append(cmds, LoadTasksCmd(m.taskSvc, firstNode.ID))
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) handleTasksLoaded(msg TasksLoadedMsg) (tea.Model, tea.Cmd) {
	m.isLoadingTasks = false
	if msg.Err != nil {
		m.taskLoadError = msg.Err
		m.addToast(toast.NewError("Failed to load tasks: " + msg.Err.Error()))
		return m, nil
	}
	m.taskLoadError = nil
	m.tasks = msg.Tasks
	m.groupedTasks = msg.GroupedTasks

	if m.expandedTaskGroups == nil {
		m.expandedTaskGroups = make(map[string]bool)
	}

	if m.groupedTasks != nil {
		for i := range m.groupedTasks.Groups {
			groupID := m.groupedTasks.Groups[i].ProjectID
			if _, exists := m.expandedTaskGroups[groupID]; !exists {
				m.expandedTaskGroups[groupID] = true
				m.groupedTasks.Groups[i].IsExpanded = true
			} else {
				m.groupedTasks.Groups[i].IsExpanded = m.expandedTaskGroups[groupID]
			}
		}
	}

	if m.selectedTaskIndex >= len(m.tasks) {
		m.selectedTaskIndex = 0
	}
	return m, nil
}
