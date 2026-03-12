package welcome

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/tui/internal/constants"
)

var PredefinedColors = []domain.Color{
	"#3B82F6", "#10B981", "#F59E0B", "#EF4444",
	"#8B5CF6", "#EC4899", "#F97316", "#14B8A6",
	"#6366F1", "#06B6D4", "#6B7280", "#92400E",
}

type Modal struct {
	input      textinput.Model
	colorIndex int
	errorMsg   string
	width      int
	height     int
}

func New() *Modal {
	ti := textinput.New()
	ti.Placeholder = "Enter area name..."
	ti.Focus()
	ti.CharLimit = MaxNameLength
	ti.Width = 40

	return &Modal{
		input:      ti,
		colorIndex: 0,
		errorMsg:   "",
	}
}

func (m *Modal) Init() tea.Cmd {
	return nil
}

func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.colorIndex = (m.colorIndex + 1) % len(PredefinedColors)
			return m, nil

		case "shift+tab":
			m.colorIndex = (m.colorIndex - 1 + len(PredefinedColors)) % len(PredefinedColors)
			return m, nil

		case constants.KeyEnter:
			name := strings.TrimSpace(m.input.Value())
			if err := ValidateName(name); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			color := PredefinedColors[m.colorIndex]
			return m, func() tea.Msg {
				return SubmitMsg{Name: name, Color: color}
			}

		case constants.KeyEsc:
			return m, func() tea.Msg {
				return ExitMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != "" && m.errorMsg != "" {
		m.errorMsg = ""
	}
	return m, cmd
}

func (m *Modal) View() string {
	var content strings.Builder

	content.WriteString(TitleStyle.Render("Welcome to Dopadone"))
	content.WriteString("\n")
	content.WriteString(SubtitleStyle.Render("Your project management companion"))
	content.WriteString("\n\n")
	content.WriteString(GuidanceStyle.Render("Create your first area to get started"))
	content.WriteString("\n\n")

	content.WriteString(InputLabelStyle.Render("Area Name:"))
	content.WriteString("\n")
	content.WriteString(InputField.Render(m.input.View()))
	content.WriteString("\n\n")

	content.WriteString(InputLabelStyle.Render("Color (Tab to change):"))
	content.WriteString("\n")
	colorDisplay := fmt.Sprintf("%s %s",
		ColorPreviewStyle.Foreground(lipgloss.Color(string(PredefinedColors[m.colorIndex]))).Render("■"),
		string(PredefinedColors[m.colorIndex]))
	content.WriteString(InputField.Render(colorDisplay))
	content.WriteString("\n")

	if m.errorMsg != "" {
		content.WriteString("\n")
		content.WriteString(ErrorText.Render("✗ " + m.errorMsg))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	content.WriteString(HintText.Render("Enter: Create • Tab: Color • ESC: Exit"))

	box := ModalBorder.Render(content.String())

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func (m *Modal) SetError(err string) {
	m.errorMsg = err
}

func (m *Modal) ClearError() {
	m.errorMsg = ""
}
