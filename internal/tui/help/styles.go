package help

import "github.com/charmbracelet/lipgloss"

var (
	HelpBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("86")).
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(1)

	CategoryStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginTop(1).
			MarginBottom(0)

	ShortcutStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86"))

	DescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252"))

	KeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Bold(true)

	HintStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)
