package tui

import "github.com/marekbrze/dopadone/internal/tui/tree"

func (m *Model) GetAreaState(areaID string) *AreaState {
	if m.areaStates[areaID] == nil {
		m.areaStates[areaID] = NewAreaState()
	}
	return m.areaStates[areaID]
}

func (m *Model) SaveCurrentAreaState() {
	if len(m.areas) == 0 || m.selectedAreaIndex >= len(m.areas) {
		return
	}
	areaID := m.areas[m.selectedAreaIndex].ID
	state := m.GetAreaState(areaID)
	state.SelectedSubareaIndex = m.selectedSubareaIndex
	state.SelectedProjectIndex = m.selectedProjectIndex
	state.SelectedTaskIndex = m.selectedTaskIndex
	m.saveTreeExpandState(state)
	m.saveTaskGroupExpandState(state)
}

func (m *Model) RestoreAreaState(areaID string) {
	state := m.GetAreaState(areaID)
	m.selectedSubareaIndex = state.SelectedSubareaIndex
	m.selectedProjectIndex = state.SelectedProjectIndex
	m.selectedTaskIndex = state.SelectedTaskIndex
	m.restoreTreeExpandState(state)
	m.restoreTaskGroupExpandState(state)
}

func (m *Model) saveTreeExpandState(state *AreaState) {
	if m.projectTree == nil {
		return
	}
	visibleNodes := tree.GetAllVisibleNodes(m.projectTree)
	for _, node := range visibleNodes {
		state.ExpandedProjects[node.ID] = node.IsExpanded
	}
}

func (m *Model) restoreTreeExpandState(state *AreaState) {
	if m.projectTree == nil {
		return
	}
	for _, node := range tree.GetAllVisibleNodes(m.projectTree) {
		if expanded, ok := state.ExpandedProjects[node.ID]; ok {
			node.IsExpanded = expanded
		}
	}
}

func (m *Model) saveTaskGroupExpandState(state *AreaState) {
	if m.expandedTaskGroups == nil {
		return
	}
	for groupID, expanded := range m.expandedTaskGroups {
		state.ExpandedTaskGroups[groupID] = expanded
	}
}

func (m *Model) restoreTaskGroupExpandState(state *AreaState) {
	if m.expandedTaskGroups == nil {
		m.expandedTaskGroups = make(map[string]bool)
	}
	for groupID, expanded := range state.ExpandedTaskGroups {
		m.expandedTaskGroups[groupID] = expanded
	}

	if m.groupedTasks != nil {
		for i := range m.groupedTasks.Groups {
			groupID := m.groupedTasks.Groups[i].ProjectID
			if expanded, ok := m.expandedTaskGroups[groupID]; ok {
				m.groupedTasks.Groups[i].IsExpanded = expanded
			}
		}
	}
}

func (m *Model) IsEmpty(column FocusColumn) bool {
	switch column {
	case FocusSubareas:
		return len(m.subareas) == 0
	case FocusProjects:
		return m.projectTree == nil || tree.GetVisibleNodeCount(m.projectTree) == 0
	case FocusTasks:
		return len(m.tasks) == 0
	}
	return true
}
