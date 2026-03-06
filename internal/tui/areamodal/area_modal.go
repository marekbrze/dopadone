package areamodal

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/tui/theme"
)

type Mode int

const (
	ModeList Mode = iota
	ModeCreate
	ModeEdit
	ModeDeleteConfirm
	ModeReorder
)

type DeleteChoice int

const (
	DeleteChoiceNone DeleteChoice = iota
	DeleteChoiceSoft
	DeleteChoiceHard
)

var PredefinedColors = []domain.Color{
	"#3B82F6", "#10B981", "#F59E0B", "#EF4444",
	"#8B5CF6", "#EC4899", "#F97316", "#14B8A6",
	"#6366F1", "#06B6D4", "#6B7280", "#92400E",
}

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

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(theme.Default.Warning).
				Bold(true)

	NormalItemStyle = lipgloss.NewStyle()

	ColorPreviewStyle = lipgloss.NewStyle().
				Padding(0, 1)

	InputField = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Default.Primary).
			Padding(0, 1).
			MarginTop(1)

	ErrorText = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Bold(true)

	HintText = lipgloss.NewStyle().
			Foreground(theme.Default.Muted).
			MarginTop(1)

	StatsStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Warning)

	WarningStyle = lipgloss.NewStyle().
			Foreground(theme.Default.Error).
			Bold(true)
)

type Area struct {
	ID        string
	Name      string
	Color     domain.Color
	SortOrder int
}

type Stats struct {
	Subareas int64
	Projects int64
	Tasks    int64
}

type Modal struct {
	mode           Mode
	areas          []Area
	selectedIndex  int
	input          textinput.Model
	colorIndex     int
	errorMsg       string
	width          int
	height         int
	stats          Stats
	editAreaID     string
	reorderChanged bool
	deleteChoice   DeleteChoice
	statsLoaded    bool
}

type SubmitMsg struct {
	Name  string
	Color domain.Color
}

type UpdateMsg struct {
	ID    string
	Name  string
	Color domain.Color
}

type DeleteMsg struct {
	ID   string
	Hard bool
}

type ReorderMsg struct {
	AreaIDs []string
}

type CloseMsg struct{}

type LoadStatsMsg struct {
	AreaID string
}

type StatsLoadedMsg struct {
	Stats Stats
}

func New(areas []Area) *Modal {
	ti := textinput.New()
	ti.Placeholder = "Enter area name..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	return &Modal{
		mode:          ModeList,
		areas:         areas,
		selectedIndex: 0,
		input:         ti,
		colorIndex:    0,
		errorMsg:      "",
		stats:         Stats{},
		deleteChoice:  DeleteChoiceNone,
	}
}

func (m *Modal) Init() tea.Cmd {
	return nil
}

func (m *Modal) SetStats(stats Stats) {
	m.stats = stats
	m.statsLoaded = true
}

func (m *Modal) UpdateAreas(areas []Area) {
	m.areas = areas
	if m.selectedIndex >= len(areas) && len(areas) > 0 {
		m.selectedIndex = len(areas) - 1
	}
}

func (m *Modal) Update(msg tea.Msg) (*Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case ModeList:
			return m.handleListKeys(msg)
		case ModeCreate, ModeEdit:
			return m.handleFormKeys(msg)
		case ModeDeleteConfirm:
			return m.handleDeleteConfirmKeys(msg)
		case ModeReorder:
			return m.handleReorderKeys(msg)
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Modal) handleListKeys(msg tea.KeyMsg) (*Modal, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedIndex > 0 {
			m.selectedIndex--
		}
	case "down", "j":
		if m.selectedIndex < len(m.areas)-1 {
			m.selectedIndex++
		}
	case "a":
		if len(m.areas) == 0 {
			m.mode = ModeCreate
			m.input.SetValue("")
			m.colorIndex = 0
			return m, nil
		}
		m.mode = ModeCreate
		m.input.SetValue("")
		m.input.Focus()
		m.colorIndex = 0
		return m, nil
	case "e":
		if len(m.areas) == 0 {
			return m, nil
		}
		area := m.areas[m.selectedIndex]
		m.mode = ModeEdit
		m.editAreaID = area.ID
		m.input.SetValue(area.Name)
		m.input.Focus()
		for i, c := range PredefinedColors {
			if c == area.Color {
				m.colorIndex = i
				break
			}
		}
		return m, nil
	case "d":
		if len(m.areas) == 0 {
			return m, nil
		}
		m.mode = ModeDeleteConfirm
		m.deleteChoice = DeleteChoiceNone
		m.statsLoaded = false
		areaID := m.areas[m.selectedIndex].ID
		return m, func() tea.Msg {
			return LoadStatsMsg{AreaID: areaID}
		}
	case "r":
		if len(m.areas) < 2 {
			return m, nil
		}
		m.mode = ModeReorder
		m.reorderChanged = false
		return m, nil
	case "enter":
		if len(m.areas) > 0 {
			return m, func() tea.Msg {
				return CloseMsg{}
			}
		}
	case "esc":
		return m, func() tea.Msg {
			return CloseMsg{}
		}
	}
	return m, nil
}

func (m *Modal) handleFormKeys(msg tea.KeyMsg) (*Modal, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.colorIndex = (m.colorIndex + 1) % len(PredefinedColors)
		return m, nil
	case "shift+tab":
		m.colorIndex = (m.colorIndex - 1 + len(PredefinedColors)) % len(PredefinedColors)
		return m, nil
	case "enter":
		name := strings.TrimSpace(m.input.Value())
		if name == "" {
			m.errorMsg = "Name is required"
			return m, nil
		}
		color := PredefinedColors[m.colorIndex]
		if m.mode == ModeCreate {
			return m, func() tea.Msg {
				return SubmitMsg{Name: name, Color: color}
			}
		} else {
			return m, func() tea.Msg {
				return UpdateMsg{ID: m.editAreaID, Name: name, Color: color}
			}
		}
	case "esc":
		m.mode = ModeList
		m.errorMsg = ""
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Modal) handleDeleteConfirmKeys(msg tea.KeyMsg) (*Modal, tea.Cmd) {
	switch msg.String() {
	case "s":
		return m, func() tea.Msg {
			return DeleteMsg{ID: m.areas[m.selectedIndex].ID, Hard: false}
		}
	case "p":
		return m, func() tea.Msg {
			return DeleteMsg{ID: m.areas[m.selectedIndex].ID, Hard: true}
		}
	case "esc":
		m.mode = ModeList
		return m, nil
	}
	return m, nil
}

func (m *Modal) handleReorderKeys(msg tea.KeyMsg) (*Modal, tea.Cmd) {
	switch msg.String() {
	case "u", "up":
		if m.selectedIndex > 0 {
			m.areas[m.selectedIndex], m.areas[m.selectedIndex-1] = m.areas[m.selectedIndex-1], m.areas[m.selectedIndex]
			m.selectedIndex--
			m.reorderChanged = true
		}
	case "d", "down":
		if m.selectedIndex < len(m.areas)-1 {
			m.areas[m.selectedIndex], m.areas[m.selectedIndex+1] = m.areas[m.selectedIndex+1], m.areas[m.selectedIndex]
			m.selectedIndex++
			m.reorderChanged = true
		}
	case "enter":
		if m.reorderChanged {
			ids := make([]string, len(m.areas))
			for i, a := range m.areas {
				ids[i] = a.ID
			}
			return m, func() tea.Msg {
				return ReorderMsg{AreaIDs: ids}
			}
		}
		m.mode = ModeList
		return m, nil
	case "esc":
		m.mode = ModeList
		return m, nil
	}
	return m, nil
}

func (m *Modal) View() string {
	switch m.mode {
	case ModeList:
		return m.viewList()
	case ModeCreate:
		return m.viewForm("Create New Area")
	case ModeEdit:
		return m.viewForm("Edit Area")
	case ModeDeleteConfirm:
		return m.viewDeleteConfirm()
	case ModeReorder:
		return m.viewReorder()
	}
	return ""
}

func (m *Modal) viewList() string {
	var content strings.Builder

	content.WriteString(TitleStyle.Render("Area Management"))
	content.WriteString("\n")

	if len(m.areas) == 0 {
		content.WriteString("\n")
		content.WriteString(HintText.Render("No areas yet."))
		content.WriteString("\n\n")
		content.WriteString(HintText.Render("Press 'a' to create your first area"))
		content.WriteString("\n")
	} else {
		for i, area := range m.areas {
			var line string
			prefix := "  "
			if i == m.selectedIndex {
				prefix = "> "
				line = fmt.Sprintf("%s%d. %s", prefix, i+1, SelectedItemStyle.Render(area.Name))
			} else {
				line = fmt.Sprintf("%s%d. %s", prefix, i+1, NormalItemStyle.Render(area.Name))
			}

			colorBlock := ""
			if area.Color != "" {
				colorBlock = ColorPreviewStyle.Foreground(lipgloss.Color(string(area.Color))).Render("■")
			}
			content.WriteString(fmt.Sprintf("%s %s\n", line, colorBlock))
		}
	}

	content.WriteString("\n")
	content.WriteString(HintText.Render("a: New • e: Edit • d: Delete • r: Reorder • Enter: Select • Esc: Close"))

	box := ModalBorder.Render(content.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m *Modal) viewForm(title string) string {
	var content strings.Builder

	content.WriteString(TitleStyle.Render(title))
	content.WriteString("\n\n")

	content.WriteString("Name:\n")
	content.WriteString(InputField.Render(m.input.View()))
	content.WriteString("\n\n")

	content.WriteString("Color (Tab to change):\n")
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
	content.WriteString(HintText.Render("Enter: Save • Tab: Change Color • Esc: Cancel"))

	box := ModalBorder.Render(content.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m *Modal) viewDeleteConfirm() string {
	var content strings.Builder

	area := m.areas[m.selectedIndex]
	content.WriteString(TitleStyle.Render("Delete Area"))
	content.WriteString("\n\n")

	content.WriteString(fmt.Sprintf("Area: %s\n\n", SelectedItemStyle.Render(area.Name)))

	statsText := fmt.Sprintf("%d subareas, %d projects, %d tasks will be affected",
		m.stats.Subareas, m.stats.Projects, m.stats.Tasks)
	content.WriteString(StatsStyle.Render(statsText))
	content.WriteString("\n\n")

	content.WriteString(WarningStyle.Render("⚠ This action cannot be undone!"))
	content.WriteString("\n\n")

	content.WriteString("Delete options:\n")
	content.WriteString("  s: Soft delete (children become orphaned)\n")
	content.WriteString("  p: Permanent delete (cascades to all children)\n")
	content.WriteString("\n")
	content.WriteString(HintText.Render("s: Soft • p: Permanent • Esc: Cancel"))

	box := ModalBorder.Render(content.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func (m *Modal) viewReorder() string {
	var content strings.Builder

	content.WriteString(TitleStyle.Render("Reorder Areas"))
	content.WriteString("\n\n")

	for i, area := range m.areas {
		var line string
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
			line = fmt.Sprintf("%s%d. %s", prefix, i+1, SelectedItemStyle.Render(area.Name))
		} else {
			line = fmt.Sprintf("%s%d. %s", prefix, i+1, NormalItemStyle.Render(area.Name))
		}

		colorBlock := ""
		if area.Color != "" {
			colorBlock = ColorPreviewStyle.Foreground(lipgloss.Color(string(area.Color))).Render("■")
		}
		content.WriteString(fmt.Sprintf("%s %s\n", line, colorBlock))
	}

	content.WriteString("\n")
	content.WriteString(HintText.Render("u/d: Move Up/Down • Enter: Save • Esc: Cancel"))

	box := ModalBorder.Render(content.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}
