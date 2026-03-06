package help

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Shortcut struct {
	Key         string
	Description string
}

type Category struct {
	Name      string
	Shortcuts []Shortcut
}

type HelpModal struct {
	width  int
	height int
}

type CloseMsg struct{}

func New() *HelpModal {
	return &HelpModal{
		width:  0,
		height: 0,
	}
}

func (h *HelpModal) Init() tea.Cmd {
	return nil
}

func (h *HelpModal) Update(msg tea.Msg) (*HelpModal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = msg.Width
		h.height = msg.Height
		return h, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "?", "esc", "q":
			return h, func() tea.Msg { return CloseMsg{} }
		}
	}

	return h, nil
}

func (h *HelpModal) View() string {
	categories := h.getCategories()

	var sections []string

	title := TitleStyle.Render("Keyboard Shortcuts")
	sections = append(sections, title)

	for _, cat := range categories {
		catTitle := CategoryStyle.Render(cat.Name)
		sections = append(sections, catTitle)

		for _, shortcut := range cat.Shortcuts {
			line := fmt.Sprintf("  %s  %s",
				KeyStyle.Render(fmt.Sprintf("%-15s", shortcut.Key)),
				DescriptionStyle.Render(shortcut.Description),
			)
			sections = append(sections, line)
		}
	}

	hint := HintStyle.Render("Press ? or Esc to close")
	sections = append(sections, hint)

	content := strings.Join(sections, "\n")

	modalWidth := h.getModalWidth()

	box := HelpBorder.
		Width(modalWidth).
		Render(content)

	if h.width == 0 || h.height == 0 {
		return box
	}

	return lipgloss.Place(
		h.width,
		h.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (h *HelpModal) getModalWidth() int {
	minWidth := 50
	maxWidth := 70

	percentageWidth := h.width * 60 / 100

	if percentageWidth < minWidth {
		return minWidth
	}
	if percentageWidth > maxWidth {
		return maxWidth
	}

	return percentageWidth
}

func (h *HelpModal) getCategories() []Category {
	return []Category{
		{
			Name: "Navigation",
			Shortcuts: []Shortcut{
				{Key: "h, ←", Description: "Move focus left (wrap)"},
				{Key: "l, →", Description: "Move focus right (wrap)"},
				{Key: "Tab", Description: "Cycle through columns"},
				{Key: "j, ↓", Description: "Navigate down (wrap)"},
				{Key: "k, ↑", Description: "Navigate up (wrap)"},
				{Key: "[", Description: "Previous area (wrap)"},
				{Key: "]", Description: "Next area (wrap)"},
			},
		},
		{
			Name: "Actions",
			Shortcuts: []Shortcut{
				{Key: "a", Description: "Quick-add item (context-aware)"},
				{Key: "Enter, Space", Description: "Toggle expand/collapse"},
				{Key: "x", Description: "Toggle task completion (Tasks column)"},
			},
		},
		{
			Name: "General",
			Shortcuts: []Shortcut{
				{Key: "Space", Description: "Open command menu"},
				{Key: "?", Description: "Show this help"},
				{Key: "q, Ctrl+C", Description: "Quit"},
			},
		},
	}
}
