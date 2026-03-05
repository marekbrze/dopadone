package tui

import tea "github.com/charmbracelet/bubbletea"

type Navigator interface {
	NavigateUp(column FocusColumn)
	NavigateDown(column FocusColumn)
	NavigateUpWithLoad(column FocusColumn) (tea.Model, tea.Cmd)
	NavigateDownWithLoad(column FocusColumn) (tea.Model, tea.Cmd)
	SwitchToPreviousArea() tea.Cmd
	SwitchToNextArea() tea.Cmd
}
