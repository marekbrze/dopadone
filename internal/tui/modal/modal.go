package modal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EntityType string

const (
	EntityTypeSubarea    EntityType = "Subarea"
	EntityTypeProject    EntityType = "Project"
	EntityTypeSubproject EntityType = "Subproject"
	EntityTypeTask       EntityType = "Task"
)

type Modal struct {
	title      string
	input      textinput.Model
	errorMsg   string
	parentName string
	entityType EntityType
	width      int
	height     int
	parentID   string
	subareaID  *string
}

type SubmitMsg struct {
	Title      string
	EntityType EntityType
	ParentID   string
	SubareaID  *string
}

type CloseMsg struct{}

func New(parentName string, entityType EntityType, parentID string, subareaID *string) *Modal {
	ti := textinput.New()
	ti.Placeholder = "Enter title..."
	ti.Focus()
	ti.CharLimit = MaxTitleLength
	ti.Width = 40

	title := fmt.Sprintf("New %s", entityType)
	if parentName != "" {
		title = fmt.Sprintf("New %s in: %s", entityType, parentName)
	}

	return &Modal{
		title:      title,
		input:      ti,
		errorMsg:   "",
		parentName: parentName,
		entityType: entityType,
		width:      0,
		height:     0,
		parentID:   parentID,
		subareaID:  subareaID,
	}
}

func (m *Modal) Init() tea.Cmd {
	return nil
}

func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			title := strings.TrimSpace(m.input.Value())
			if err := ValidateTitle(title); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			return m, func() tea.Msg {
				return SubmitMsg{
					Title:      title,
					EntityType: m.entityType,
					ParentID:   m.parentID,
					SubareaID:  m.subareaID,
				}
			}

		case "esc":
			return m, func() tea.Msg {
				return CloseMsg{}
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Modal) View() string {
	titleWidth := lipgloss.Width(TitleStyle.Render(m.title))
	inputWidth := titleWidth
	if inputWidth < 30 {
		inputWidth = 30
	}
	if inputWidth > 60 {
		inputWidth = 60
	}
	m.input.Width = inputWidth

	inputView := m.input.View()

	var content strings.Builder
	content.WriteString(TitleStyle.Render(m.title))
	content.WriteString("\n\n")

	if m.errorMsg != "" {
		content.WriteString(ErrorText.Render("✗ " + m.errorMsg))
		content.WriteString("\n\n")
	}

	content.WriteString(InputField.Render(inputView))
	content.WriteString("\n\n")
	content.WriteString(HintText.Render("Enter: Create • Esc: Cancel"))

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
