package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/tui/areamodal"
	"github.com/example/projectdb/internal/tui/help"
	"github.com/example/projectdb/internal/tui/modal"
	"github.com/example/projectdb/internal/tui/toast"
	"github.com/example/projectdb/internal/tui/tree"
)

func (m *Model) handleEnterOrSpace() {
	if m.focus != FocusProjects {
		return
	}
	m.toggleTreeExpand()
}

func (m *Model) toggleTreeExpand() {
	if m.selectedProjectID == "" || m.projectTree == nil {
		return
	}
	node := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
	if node == nil || node.IsLeaf() {
		return
	}
	node.ToggleExpanded()
	if len(m.areas) > 0 && m.selectedAreaIndex < len(m.areas) {
		areaID := m.areas[m.selectedAreaIndex].ID
		state := m.GetAreaState(areaID)
		state.ExpandedProjects[m.selectedProjectID] = node.IsExpanded
	}
}

func (m *Model) handleQuickAdd() (tea.Model, tea.Cmd) {
	parentName, entityType, parentID, subareaID := m.getParentContext()
	if parentName == "" {
		return m, nil
	}

	m.modal = modal.New(parentName, entityType, parentID, subareaID)
	m.isModalOpen = true

	return m, nil
}

func (m *Model) getParentContext() (string, modal.EntityType, string, *string) {
	switch m.focus {
	case FocusSubareas:
		if len(m.areas) == 0 || m.selectedAreaIndex >= len(m.areas) {
			return "", "", "", nil
		}
		area := m.areas[m.selectedAreaIndex]
		return area.Name, modal.EntityTypeSubarea, area.ID, nil

	case FocusProjects:
		if m.selectedProjectID != "" {
			projectName := m.getProjectNameByID(m.selectedProjectID)
			if projectName != "" {
				return projectName, modal.EntityTypeSubproject, m.selectedProjectID, nil
			}
		}
		if len(m.subareas) == 0 || m.selectedSubareaIndex >= len(m.subareas) {
			return "", "", "", nil
		}
		subarea := m.subareas[m.selectedSubareaIndex]
		return subarea.Name, modal.EntityTypeProject, "", &subarea.ID

	case FocusTasks:
		if m.selectedProjectID == "" {
			return "", "", "", nil
		}
		projectName := m.getProjectNameByID(m.selectedProjectID)
		if projectName == "" {
			return "", "", "", nil
		}
		return projectName, modal.EntityTypeTask, m.selectedProjectID, nil
	}

	return "", "", "", nil
}

func (m *Model) getProjectNameByID(id string) string {
	if m.projectTree == nil {
		return ""
	}
	node := tree.FindNodeByID(m.projectTree, id)
	if node == nil {
		return ""
	}
	return node.Name
}

func (m *Model) handleModalSubmit(msg modal.SubmitMsg) (tea.Model, tea.Cmd) {
	switch msg.EntityType {
	case modal.EntityTypeSubarea:
		return m, CreateSubareaCmd(m.subareaSvc, msg.Title, msg.ParentID)

	case modal.EntityTypeProject:
		return m, CreateProjectCmd(m.projectSvc, msg.Title, nil, msg.SubareaID)

	case modal.EntityTypeSubproject:
		return m, CreateProjectCmd(m.projectSvc, msg.Title, &msg.ParentID, nil)

	case modal.EntityTypeTask:
		return m, CreateTaskCmd(m.taskSvc, msg.Title, msg.ParentID)
	}

	m.isModalOpen = false
	m.modal = nil
	return m, nil
}

func (m *Model) handleHelp() *Model {
	m.helpModal = help.New()
	m.helpModal.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	m.isHelpOpen = true
	return m
}

func (m *Model) handleOpenAreaModal() (tea.Model, tea.Cmd) {
	areas := make([]areamodal.Area, len(m.areas))
	for i, a := range m.areas {
		areas[i] = areamodal.Area{
			ID:        a.ID,
			Name:      a.Name,
			Color:     a.Color,
			SortOrder: a.SortOrder,
		}
	}
	m.areaModal = areamodal.New(areas)
	m.isAreaModalOpen = true
	return m, nil
}

func (m *Model) addToast(t toast.Toast) {
	m.toasts = append(m.toasts, t)
	if len(m.toasts) > toast.MaxToasts {
		m.toasts = m.toasts[len(m.toasts)-toast.MaxToasts:]
	}
}

func (m *Model) removeExpiredToasts() {
	var validToasts []toast.Toast
	for _, t := range m.toasts {
		if !t.IsExpired() {
			validToasts = append(validToasts, t)
		}
	}
	m.toasts = validToasts
}
