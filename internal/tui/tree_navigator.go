package tui

import "github.com/example/projectdb/internal/tui/tree"

func (m *Model) navigateTreeUp() {
	if m.projectTree == nil {
		return
	}
	current := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
	if current == nil {
		first := tree.GetFirstVisibleNode(m.projectTree)
		if first != nil {
			m.selectedProjectID = first.ID
			m.syncTreeSelectionToIndex()
		}
		return
	}
	prev := tree.GetPrevVisibleNode(m.projectTree, current)
	if prev != nil {
		m.selectedProjectID = prev.ID
		m.syncTreeSelectionToIndex()
	} else {
		last := tree.GetLastVisibleNode(m.projectTree)
		if last != nil {
			m.selectedProjectID = last.ID
			m.syncTreeSelectionToIndex()
		}
	}
}

func (m *Model) navigateTreeDown() {
	if m.projectTree == nil {
		return
	}
	current := tree.FindNodeByID(m.projectTree, m.selectedProjectID)
	if current == nil {
		first := tree.GetFirstVisibleNode(m.projectTree)
		if first != nil {
			m.selectedProjectID = first.ID
			m.syncTreeSelectionToIndex()
		}
		return
	}
	next := tree.GetNextVisibleNode(m.projectTree, current)
	if next != nil {
		m.selectedProjectID = next.ID
		m.syncTreeSelectionToIndex()
	} else {
		first := tree.GetFirstVisibleNode(m.projectTree)
		if first != nil {
			m.selectedProjectID = first.ID
			m.syncTreeSelectionToIndex()
		}
	}
}

func (m *Model) syncTreeSelectionToIndex() {
	if m.projectTree == nil {
		return
	}
	visibleNodes := tree.GetAllVisibleNodes(m.projectTree)
	for i, node := range visibleNodes {
		if node.ID == m.selectedProjectID {
			m.selectedProjectIndex = i
			return
		}
	}
	m.selectedProjectIndex = 0
}
