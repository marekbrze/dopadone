package spacemenu

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

var (
	MenuBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Primary).
			MarginBottom(1)

	KeyStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Warning).
			Bold(true)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Foreground)

	HintStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)
)
