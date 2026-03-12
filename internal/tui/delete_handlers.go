package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/tui/confirmmodal"
	"github.com/marekbrze/dopadone/internal/tui/toast"
	"github.com/marekbrze/dopadone/internal/tui/tree"
)

const (
	entityTypeSubarea = "Subarea"
	entityTypeProject = "Project"
	entityTypeTask    = "Task"
)

func (m *Model) handleDeleteKey() tea.Cmd {
	switch m.focus {
	case FocusSubareas:
		if len(m.subareas) == 0 || m.selectedSubareaIndex >= len(m.subareas) {
			return nil
		}
		subarea := m.subareas[m.selectedSubareaIndex]
		m.confirmModal = confirmmodal.New(
			subarea.Name,
			confirmmodal.EntityTypeSubarea,
			subarea.ID,
		)
		m.isConfirmModalOpen = true
		return nil

	case FocusProjects:
		if m.projectTree == nil {
			return nil
		}
		project := m.getSelectedProject()
		if project == nil {
			return nil
		}
		m.confirmModal = confirmmodal.New(
			project.Name,
			confirmmodal.EntityTypeProject,
			project.ID,
		)
		m.isConfirmModalOpen = true
		return nil
	case FocusTasks:
		task := m.getTaskAtLine(m.selectedTaskIndex)
		if task == nil {
			return nil
		}
		m.confirmModal = confirmmodal.New(
			task.Title,
			confirmmodal.EntityTypeTask,
			task.ID,
		)
		m.isConfirmModalOpen = true
		return nil
	}
	return nil
}

func (m *Model) getSelectedProject() *domain.Project {
	if m.projectTree == nil || m.selectedProjectID == "" {
		return nil
	}
	node := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
	if node == nil || node.Data == nil {
		return nil
	}
	project, ok := node.Data.(*domain.Project)
	if !ok {
		return nil
	}
	return project
}

func (m *Model) handleConfirmModalConfirm(msg confirmmodal.ConfirmMsg) (tea.Model, tea.Cmd) {
	m.isConfirmModalOpen = false
	m.confirmModal = nil
	switch msg.EntityType {
	case confirmmodal.EntityTypeSubarea:
		return m, DeleteSubareaCmd(m.subareaSvc, msg.EntityID, msg.EntityName)
	case confirmmodal.EntityTypeProject:
		return m, DeleteProjectCmd(m.projectSvc, msg.EntityID, msg.EntityName)
	case confirmmodal.EntityTypeTask:
		return m, DeleteTaskCmd(m.taskSvc, msg.EntityID, msg.EntityName)
	}
	return m, nil
}

func (m *Model) handleConfirmModalCancel() {
	m.isConfirmModalOpen = false
	m.confirmModal = nil
}

func (m *Model) handleDeleteSuccess(msg DeleteSuccessMsg) (tea.Model, tea.Cmd) {
	m.addToast(toast.NewSuccess(fmt.Sprintf("%s '%s' deleted successfully", msg.EntityType, msg.EntityName)))
	switch msg.EntityType {
	case entityTypeSubarea:
		m.subareas = nil
		m.projects = nil
		m.tasks = nil
		m.projectTree = nil
		if len(m.areas) > 0 {
			m.isLoadingSubareas = true
			return m, LoadSubareasCmd(m.subareaSvc, m.areas[m.selectedTab].ID)
		}
	case entityTypeProject:
		m.projects = nil
		m.tasks = nil
		m.projectTree = nil
		if len(m.subareas) > 0 && m.selectedSubareaIndex < len(m.subareas) {
			m.isLoadingProjects = true
			return m, LoadProjectsCmd(m.projectSvc, &m.subareas[m.selectedSubareaIndex].ID)
		}
	case entityTypeTask:
		m.tasks = nil
		if m.selectedProjectID != "" {
			m.isLoadingTasks = true
			return m, LoadTasksCmd(m.taskSvc, m.selectedProjectID)
		}
	}
	return m, nil
}

func (m *Model) handleDeleteError(msg DeleteErrorMsg) {
	m.addToast(toast.NewError(fmt.Sprintf("Failed to delete %s '%s': %v", msg.EntityType, msg.EntityName, msg.Err)))
}
