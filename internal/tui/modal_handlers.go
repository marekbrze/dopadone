package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/dopadone/internal/tui/toast"
)

func (m *Model) handleSubareaCreated(msg SubareaCreatedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.isModalOpen = false
		m.modal = nil
		m.addToast(toast.NewError("Failed to create subarea: " + msg.Err.Error()))
		return m, nil
	}

	m.isModalOpen = false
	m.modal = nil
	m.addToast(toast.NewSuccess("Subarea created successfully"))

	if len(m.areas) == 0 {
		return m, nil
	}
	areaID := m.areas[m.selectedAreaIndex].ID
	m.isLoadingSubareas = true

	return m, LoadSubareasCmd(m.subareaSvc, areaID)
}

func (m *Model) handleProjectCreated(msg ProjectCreatedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.isModalOpen = false
		m.modal = nil
		m.addToast(toast.NewError("Failed to create project: " + msg.Err.Error()))
		return m, nil
	}

	m.isModalOpen = false
	m.modal = nil
	m.addToast(toast.NewSuccess("Project created successfully"))

	if len(m.subareas) == 0 {
		return m, nil
	}
	subareaID := m.subareas[m.selectedSubareaIndex].ID
	m.isLoadingProjects = true

	return m, LoadProjectsCmd(m.projectSvc, &subareaID)
}

func (m *Model) handleTaskCreated(msg TaskCreatedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.isModalOpen = false
		m.modal = nil
		m.addToast(toast.NewError("Failed to create task: " + msg.Err.Error()))
		return m, nil
	}

	m.isModalOpen = false
	m.modal = nil
	m.addToast(toast.NewSuccess("Task created successfully"))

	if m.selectedProjectID == "" {
		return m, nil
	}
	m.isLoadingTasks = true

	return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
}
