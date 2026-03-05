package modal

import "github.com/charmbracelet/lipgloss"

var (
	ModalBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	ErrorText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	HintText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("62"))

	OverlayStyle = lipgloss.NewStyle()
)
