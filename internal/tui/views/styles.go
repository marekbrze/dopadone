package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/example/dopadone/internal/tui/theme"
)

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Background(theme.Default.TabActiveBackground()).
			Foreground(theme.Default.TabActiveForeground()).
			Bold(true).
			Underline(true)

	InactiveTabStyle = lipgloss.NewStyle().
				Padding(0, 3).
				Background(theme.Default.TabInactiveBackground()).
				Foreground(theme.Default.TabInactiveForeground())

	TabSeparator = lipgloss.NewStyle().
			Foreground(theme.Default.Dimmed).
			SetString("│")

	FocusedColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(theme.Default.ColumnFocusedBorder()).
				Padding(0, 1)

	UnfocusedColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(theme.Default.ColumnUnfocusedBorder()).
				Padding(0, 1)

	ColumnHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(theme.Default.ColumnHeader()).
				MarginBottom(1)

	EmptyContentStyle = lipgloss.NewStyle().
				Foreground(theme.Default.EmptyText()).
				Italic(true)
)
