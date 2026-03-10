package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleToastTick() (tea.Model, tea.Cmd) {
	m.removeExpiredToasts()
	return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return ToastTickMsg{}
	})
}
