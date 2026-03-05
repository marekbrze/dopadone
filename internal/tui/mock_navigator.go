package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MockNavigator struct {
	NavigateUpFunc           func(column FocusColumn)
	NavigateDownFunc         func(column FocusColumn)
	NavigateUpWithLoadFunc   func(column FocusColumn) (tea.Model, tea.Cmd)
	NavigateDownWithLoadFunc func(column FocusColumn) (tea.Model, tea.Cmd)
	SwitchToPreviousAreaFunc func() tea.Cmd
	SwitchToNextAreaFunc     func() tea.Cmd
}

func (m *MockNavigator) NavigateUp(column FocusColumn) {
	if m.NavigateUpFunc != nil {
		m.NavigateUpFunc(column)
	}
}

func (m *MockNavigator) NavigateDown(column FocusColumn) {
	if m.NavigateDownFunc != nil {
		m.NavigateDownFunc(column)
	}
}

func (m *MockNavigator) NavigateUpWithLoad(column FocusColumn) (tea.Model, tea.Cmd) {
	if m.NavigateUpWithLoadFunc != nil {
		return m.NavigateUpWithLoadFunc(column)
	}
	return nil, nil
}

func (m *MockNavigator) NavigateDownWithLoad(column FocusColumn) (tea.Model, tea.Cmd) {
	if m.NavigateDownWithLoadFunc != nil {
		return m.NavigateDownWithLoadFunc(column)
	}
	return nil, nil
}

func (m *MockNavigator) SwitchToPreviousArea() tea.Cmd {
	if m.SwitchToPreviousAreaFunc != nil {
		return m.SwitchToPreviousAreaFunc()
	}
	return nil
}

func (m *MockNavigator) SwitchToNextArea() tea.Cmd {
	if m.SwitchToNextAreaFunc != nil {
		return m.SwitchToNextAreaFunc()
	}
	return nil
}
