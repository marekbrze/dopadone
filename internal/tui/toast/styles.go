package toast

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

var (
	ErrorStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Background(theme.Default.Error).
			Bold(true).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Success).
			Background(theme.Default.Success).
			Bold(true).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	InfoStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Warning).
			Background(theme.Default.Warning).
			Padding(0, 1).
			Margin(0, 0, 1, 0)
)

const (
	TypeError   = "error"
	TypeSuccess = "success"
	TypeInfo    = "info"
)

const (
	ToastDuration = 3000
	MaxToasts     = 3
)
