package views

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	ActiveTabStyle = lipgloss.NewStyle().
			Padding(0, 3).
			Background(lipgloss.Color("39")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Underline(true)

	InactiveTabStyle = lipgloss.NewStyle().
				Padding(0, 3).
				Background(lipgloss.Color("238")).
				Foreground(lipgloss.Color("252"))

	TabSeparator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")).
			SetString("│")

	FocusedColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(lipgloss.Color("62")).
				Padding(0, 1)

	UnfocusedColumnStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("238")).
				Padding(0, 1)

	ColumnHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("86")).
				MarginBottom(1)

	EmptyContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("241")).
				Italic(true)
)
