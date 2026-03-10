package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/toast"
)

func (m *Model) handleTaskStatusToggled(msg TaskStatusToggledMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		if msg.TaskIndex < len(m.tasks) {
			m.tasks[msg.TaskIndex].Status = msg.OriginalStatus
		}

		m.addToast(toast.NewError("Failed to update task status: " + msg.Err.Error()))
		return m, nil
	}

	return m, nil
}
