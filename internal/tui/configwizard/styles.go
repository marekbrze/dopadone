package configwizard

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

var (
	brandStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Primary).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Primary).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginBottom(1)

	modeTitleStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Foreground).
			Bold(true)

	modeDescStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted)

	modeSelectedStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Primary).
				Bold(true)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Foreground).
			MarginBottom(0)

	inputFieldStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Muted).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Success).
			Bold(true)

	hintStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(1, 2).
			Width(60)

	spinnerStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Primary)

	ModalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(1, 2).
			Width(65)
)
