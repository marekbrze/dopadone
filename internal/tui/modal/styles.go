package modal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/example/projectdb/internal/tui/theme"
)

var (
	ModalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(1, 2)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(0, 1)

	ErrorText = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Bold(true)

	HintText = lipgloss.NewStyle().
			Foreground(theme.Default.Muted)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Primary)

	OverlayStyle = lipgloss.NewStyle()

	CheckboxStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted)

	CheckboxFocusedStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Primary).
				Bold(true)

	CheckboxCheckedStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Success)

	CheckboxFocusedCheckedStyle = lipgloss.NewStyle().
					Foreground(theme.Default.Success).
					Bold(true)
)
