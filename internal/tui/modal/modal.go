package modal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/internal/constants"
)

type EntityType string

const (
	EntityTypeSubarea    EntityType = "Subarea"
	EntityTypeProject    EntityType = "Project"
	EntityTypeSubproject EntityType = "Subproject"
	EntityTypeTask       EntityType = "Task"
)

type focusedElement int

const (
	focusedInput focusedElement = iota
	focusedCheckbox
)

type Modal struct {
	title           string
	input           textinput.Model
	errorMsg        string
	parentName      string
	entityType      EntityType
	width           int
	height          int
	parentID        string
	subareaID       *string
	showCheckbox    bool
	checkboxChecked bool
	focusedElement  focusedElement
}

type SubmitMsg struct {
	Title        string
	EntityType   EntityType
	ParentID     string
	SubareaID    *string
	AsSubproject bool
}

type CloseMsg struct{}

func New(parentName string, entityType EntityType, parentID string, subareaID *string, showCheckbox bool) *Modal {
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
		title:           title,
		input:           ti,
		errorMsg:        "",
		parentName:      parentName,
		entityType:      entityType,
		width:           0,
		height:          0,
		parentID:        parentID,
		subareaID:       subareaID,
		showCheckbox:    showCheckbox,
		checkboxChecked: false,
		focusedElement:  focusedInput,
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
		case constants.KeyEnter:
			title := strings.TrimSpace(m.input.Value())
			if err := ValidateTitle(title); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			asSubproject := m.showCheckbox && m.checkboxChecked
			return m, func() tea.Msg {
				return SubmitMsg{
					Title:        title,
					EntityType:   m.entityType,
					ParentID:     m.parentID,
					SubareaID:    m.subareaID,
					AsSubproject: asSubproject,
				}
			}

		case constants.KeyEsc:
			return m, func() tea.Msg {
				return CloseMsg{}
			}

		case "tab":
			if m.showCheckbox {
				if m.focusedElement == focusedInput {
					m.focusedElement = focusedCheckbox
					m.input.Blur()
				} else {
					m.focusedElement = focusedInput
					m.input.Focus()
				}
			}
			return m, nil

		case "shift+tab":
			if m.showCheckbox {
				if m.focusedElement == focusedInput {
					m.focusedElement = focusedCheckbox
					m.input.Blur()
				} else {
					m.focusedElement = focusedInput
					m.input.Focus()
				}
			}
			return m, nil

		case " ":
			if m.showCheckbox && m.focusedElement == focusedCheckbox {
				m.checkboxChecked = !m.checkboxChecked
				return m, nil
			}
		}
	}

	if m.focusedElement == focusedInput {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	}

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
	content.WriteString("\n")

	if m.showCheckbox {
		content.WriteString("\n")
		checkboxText := fmt.Sprintf("[ ] Add as subproject of %s", m.parentName)
		if m.checkboxChecked {
			checkboxText = fmt.Sprintf("[✓] Add as subproject of %s", m.parentName)
		}

		var style lipgloss.Style
		if m.focusedElement == focusedCheckbox {
			if m.checkboxChecked {
				style = CheckboxFocusedCheckedStyle
			} else {
				style = CheckboxFocusedStyle
			}
		} else {
			if m.checkboxChecked {
				style = CheckboxCheckedStyle
			} else {
				style = CheckboxStyle
			}
		}

		content.WriteString(style.Render(checkboxText))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	if m.showCheckbox {
		content.WriteString(HintText.Render("Tab: Switch • Space: Toggle • Enter: Create • Esc: Cancel"))
	} else {
		content.WriteString(HintText.Render("Enter: Create • Esc: Cancel"))
	}

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
