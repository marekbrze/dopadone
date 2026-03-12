package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.ready = true
	if m.helpModal != nil {
		m.helpModal, _ = m.helpModal.Update(msg)
	}
	if m.areaModal != nil {
		m.areaModal, _ = m.areaModal.Update(msg)
	}
	if m.spaceMenu != nil {
		m.spaceMenu, _ = m.spaceMenu.Update(msg)
	}
	if m.confirmModal != nil {
		m.confirmModal, _ = m.confirmModal.Update(msg)
	}
	if m.welcomeModal != nil {
		m.welcomeModal, _ = m.welcomeModal.Update(msg)
	}
	return m, nil
}
