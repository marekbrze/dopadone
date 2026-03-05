package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/db"
)

func New(repo db.Querier) *tea.Program {
	return tea.NewProgram(
		InitialModel(repo),
		tea.WithAltScreen(),
	)
}
