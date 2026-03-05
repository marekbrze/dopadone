package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MockHandler struct {
	HandleMessageFunc func(msg tea.Msg) (tea.Model, tea.Cmd)
}

func (m *MockHandler) HandleMessage(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.HandleMessageFunc != nil {
		return m.HandleMessageFunc(msg)
	}
	return nil, nil
}
