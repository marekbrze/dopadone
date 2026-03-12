package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/areamodal"
	"github.com/marekbrze/dopadone/internal/tui/help"
	"github.com/marekbrze/dopadone/internal/tui/modal"
	"github.com/marekbrze/dopadone/internal/tui/toast"
	"github.com/marekbrze/dopadone/internal/tui/tree"
)

func (m *Model) handleEnterOrSpace() {
	switch m.focus {
	case FocusProjects:
		m.toggleTreeExpand()
	case FocusTasks:
		m.toggleTaskGroup()
	}
}

func (m *Model) toggleTaskGroup() {
	if !m.isLineGroupHeader(m.selectedTaskIndex) {
		return
	}

	group := m.getGroupAtLine(m.selectedTaskIndex)
	if group == nil {
		return
	}

	wasExpanded := m.expandedTaskGroups[group.ProjectID]
	m.expandedTaskGroups[group.ProjectID] = !wasExpanded

	for i := range m.groupedTasks.Groups {
		if m.groupedTasks.Groups[i].ProjectID == group.ProjectID {
			m.groupedTasks.Groups[i].IsExpanded = !wasExpanded
			break
		}
	}

	if wasExpanded {
		headerLine := m.getGroupHeaderLineForGroup(group.ProjectID)
		m.selectedTaskIndex = headerLine
	}
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
	parentName, entityType, parentID, subareaID, showCheckbox := m.getParentContext()
	if parentName == "" {
		return m, nil
	}

	m.modal = modal.New(parentName, entityType, parentID, subareaID, showCheckbox)
	m.isModalOpen = true

	return m, nil
}

func (m *Model) getParentContext() (string, modal.EntityType, string, *string, bool) {
	switch m.focus {
	case FocusSubareas:
		if len(m.areas) == 0 || m.selectedAreaIndex >= len(m.areas) {
			return "", "", "", nil, false
		}
		area := m.areas[m.selectedAreaIndex]
		return area.Name, modal.EntityTypeSubarea, area.ID, nil, false

	case FocusProjects:
		if m.selectedProjectID != "" {
			projectName := m.getProjectNameByID(m.selectedProjectID)
			if projectName != "" && len(m.subareas) > 0 && m.selectedSubareaIndex < len(m.subareas) {
				subareaID := m.subareas[m.selectedSubareaIndex].ID
				return projectName, modal.EntityTypeProject, m.selectedProjectID, &subareaID, true
			}
		}
		if len(m.subareas) == 0 || m.selectedSubareaIndex >= len(m.subareas) {
			return "", "", "", nil, false
		}
		subarea := m.subareas[m.selectedSubareaIndex]
		return subarea.Name, modal.EntityTypeProject, "", &subarea.ID, false

	case FocusTasks:
		if m.selectedProjectID == "" {
			return "", "", "", nil, false
		}
		projectName := m.getProjectNameByID(m.selectedProjectID)
		if projectName == "" {
			return "", "", "", nil, false
		}
		return projectName, modal.EntityTypeTask, m.selectedProjectID, nil, false
	}

	return "", "", "", nil, false
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
		if msg.AsSubproject {
			return m, CreateProjectCmd(m.projectSvc, msg.Title, &msg.ParentID, nil)
		}
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
	return m.handleOpenAreaModalWithMode(areamodal.ModeList)
}

func (m *Model) handleOpenAreaModalWithMode(mode areamodal.Mode) (tea.Model, tea.Cmd) {
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

	if len(m.areas) > 0 && m.selectedAreaIndex < len(m.areas) {
		m.areaModal.SetSelectedIndex(m.selectedAreaIndex)
	}

	switch mode {
	case areamodal.ModeCreate:
		m.areaModal.SetupForCreate()
	case areamodal.ModeEdit:
		m.areaModal.SetupForEdit()
	case areamodal.ModeDeleteConfirm:
		m.areaModal.SetupForDelete()
	}

	m.areaModal, _ = m.areaModal.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
	m.isAreaModalOpen = true

	if mode == areamodal.ModeDeleteConfirm && len(m.areas) > 0 {
		areaID := m.areas[m.selectedAreaIndex].ID
		return m, LoadAreaStatsCmd(m.areaSvc, areaID)
	}

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

func (m *Model) handleModalClose() (tea.Model, tea.Cmd) {
	m.isModalOpen = false
	m.modal = nil
	return m, nil
}

func (m *Model) handleAreaModalClose() (tea.Model, tea.Cmd) {
	m.isAreaModalOpen = false
	m.areaModal = nil
	return m, nil
}

func (m *Model) handleHelpClose() (tea.Model, tea.Cmd) {
	m.isHelpOpen = false
	m.helpModal = nil
	return m, nil
}
