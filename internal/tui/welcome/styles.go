package welcome

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

var (
	ModalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(1, 2).
			Width(60)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Default.Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Secondary).
			MarginBottom(1)

	GuidanceStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginBottom(1)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(0, 1).
			MarginTop(1)

	InputLabelStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Foreground)

	ColorPreviewStyle = lipgloss.NewStyle().
				Padding(0, 1)

	ErrorText = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Bold(true)

	HintText = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)
)
