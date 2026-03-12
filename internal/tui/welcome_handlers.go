package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/toast"
	"github.com/marekbrze/dopadone/internal/tui/welcome"
)

func (m *Model) handleWelcomeSubmit(msg welcome.SubmitMsg) (tea.Model, tea.Cmd) {
	m.isWelcomeOpen = false
	m.welcomeModal = nil
	m.isFromWelcomeFlow = true
	return m, CreateAreaCmd(m.areaSvc, msg.Name, msg.Color)
}

func (m *Model) handleWelcomeExit(msg welcome.ExitMsg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m *Model) handleWelcomeAreasLoaded(msg AreasLoadedMsg) (tea.Model, tea.Cmd) {
	m.isLoadingAreas = false
	if msg.Err != nil {
		m.areaLoadError = msg.Err
		m.addToast(toast.NewError("Failed to load areas: " + msg.Err.Error()))
		m.isFromWelcomeFlow = false
		return m, nil
	}

	m.areaLoadError = nil
	m.areas = msg.Areas
	m.tabs = updateTabsFromAreas(m.areas, 0)
	m.selectedTab = 0
	m.selectedAreaIndex = 0

	if len(m.areas) > 0 {
		m.isLoadingSubareas = true
		m.isFromWelcomeFlow = false
		return m, LoadSubareasCmd(m.subareaSvc, m.areas[0].ID)
	}

	m.isFromWelcomeFlow = false
	return m, nil
}

func (m *Model) handleWelcomeMessages(msg interface{}) (tea.Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case welcome.SubmitMsg:
		if m.isWelcomeOpen {
			model, cmd := m.handleWelcomeSubmit(msg)
			return model, cmd, true
		}
	case welcome.ExitMsg:
		if m.isWelcomeOpen {
			model, cmd := m.handleWelcomeExit(msg)
			return model, cmd, true
		}
	case AreaCreatedMsg:
		if m.isFromWelcomeFlow {
			if msg.Err != nil {
				m.addToast(toast.NewError("Failed to create area: " + msg.Err.Error()))
				m.isFromWelcomeFlow = false
				return m, nil, true
			}
			m.addToast(toast.NewSuccess("Area created successfully"))
			m.isLoadingAreas = true
			return m, LoadAreasCmd(m.areaSvc), true
		}
	case AreasLoadedMsg:
		if m.isFromWelcomeFlow {
			model, cmd := m.handleWelcomeAreasLoaded(msg)
			return model, cmd, true
		}
	}

	return m, nil, false
}

func (m *Model) handleWelcomeKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.welcomeModal, cmd = m.welcomeModal.Update(msg)
	return m, cmd
}
