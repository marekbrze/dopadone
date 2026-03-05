package toast

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Background(lipgloss.Color("52")).
			Bold(true).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Background(lipgloss.Color("22")).
			Bold(true).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("58")).
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
