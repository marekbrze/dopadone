package confirmmodal

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EntityType string

const (
	EntityTypeSubarea EntityType = "Subarea"
	EntityTypeProject EntityType = "Project"
	EntityTypeTask    EntityType = "Task"
)

type Modal struct {
	itemName   string
	entityType EntityType
	entityID   string
	width      int
	height     int
}

type ConfirmMsg struct {
	EntityType EntityType
	EntityID   string
	EntityName string
}

type CancelMsg struct{}

func New(itemName string, entityType EntityType, entityID string) *Modal {
	return &Modal{
		itemName:   itemName,
		entityType: entityType,
		entityID:   entityID,
		width:      0,
		height:     0,
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
		case "y":
			return m, func() tea.Msg {
				return ConfirmMsg{
					EntityType: m.entityType,
					EntityID:   m.entityID,
					EntityName: m.itemName,
				}
			}

		case "n", "esc":
			return m, func() tea.Msg {
				return CancelMsg{}
			}
		}
	}

	return m, nil
}

func (m *Modal) View() string {
	title := TitleStyle.Render(fmt.Sprintf("Delete %s?", m.entityType))
	message := MessageStyle.Render(fmt.Sprintf("%s", m.truncateItemName()))
	hint := HintStyle.Render("y: confirm | n/esc: cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		message,
		"",
		hint,
	)

	box := BorderStyle.Render(content)

	if m.width == 0 || m.height == 0 {
		return box
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (m *Modal) truncateItemName() string {
	maxLen := 40
	if len(m.itemName) <= maxLen {
		return m.itemName
	}
	return m.itemName[:maxLen-3] + "..."
}
