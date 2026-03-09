package confirmmodal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

var (
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Error).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Error)

	MessageStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Foreground)

	HintStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)
)
