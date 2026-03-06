package help

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/example/projectdb/internal/tui/theme"
)

var (
	HelpBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Success).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Success).
			MarginBottom(1)

	CategoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Warning).
			MarginTop(1).
			MarginBottom(0)

	ShortcutStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Success)

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Foreground)

	KeyStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Warning).
			Bold(true)

	HintStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)
)
