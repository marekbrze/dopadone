package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleConnectionStatusUpdated(msg ConnectionStatusUpdatedMsg) (tea.Model, tea.Cmd) {
	m.connectionStatus = ConnectionStatusView{
		Mode:       msg.DriverType,
		Status:     msg.Status,
		SyncStatus: msg.SyncInfo.Status,
		LastSyncAt: msg.SyncInfo.LastSyncAt,
	}

	if msg.SyncInfo.LastError != nil {
		m.connectionStatus.ErrorMessage = msg.SyncInfo.LastError.Error()
	}

	if m.dbDriver != nil {
		return m, PollConnectionStatusCmd(m.dbDriver)
	}

	return m, nil
}

func (m *Model) handleConnectionMessages(msg interface{}) (tea.Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case ConnectionStatusUpdatedMsg:
		model, cmd := m.handleConnectionStatusUpdated(msg)
		return model, cmd, true
	}
	return m, nil, false
}
