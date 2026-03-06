package spacemenu

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CloseMsg struct{}

type ActionMsg struct {
	Action MenuAction
}

type SpaceMenu struct {
	state  MenuState
	width  int
	height int
}

func New() *SpaceMenu {
	return &SpaceMenu{
		state:  StateMain,
		width:  0,
		height: 0,
	}
}

func (sm *SpaceMenu) Init() tea.Cmd {
	return nil
}

func (sm *SpaceMenu) Update(msg tea.Msg) (*SpaceMenu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		sm.width = msg.Width
		sm.height = msg.Height
		return sm, nil

	case tea.KeyMsg:
		switch msg.String() {
		case " ", "esc", "q":
			if sm.state == StateMain {
				return sm, func() tea.Msg { return CloseMsg{} }
			}
			sm.state = StateMain
			return sm, nil
		case "c":
			if sm.state == StateMain {
				sm.state = StateConfig
				return sm, nil
			}
		}
	}

	return sm, nil
}

func (sm *SpaceMenu) View() string {
	var content strings.Builder

	switch sm.state {
	case StateMain:
		content.WriteString(sm.renderMainMenu())
	case StateConfig:
		content.WriteString(sm.renderConfigMenu())
	}

	box := MenuBorder.
		Width(sm.getModalWidth()).
		Render(content.String())

	if sm.width == 0 || sm.height == 0 {
		return box
	}

	return lipgloss.Place(
		sm.width,
		sm.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

func (sm *SpaceMenu) renderMainMenu() string {
	var sections []string

	title := TitleStyle.Render("Command Menu")
	sections = append(sections, title)

	commands := sm.getMainCommands()
	for _, cmd := range commands {
		keyLabel := fmt.Sprintf("%s: %s", cmd.Key, cmd.Label)
		line := fmt.Sprintf("  %s  %s",
			KeyStyle.Render(fmt.Sprintf("%-15s", keyLabel)),
			DescriptionStyle.Render(cmd.Description),
		)
		sections = append(sections, line)
	}

	hint := HintStyle.Render("Press Space, Esc, or q to close")
	sections = append(sections, hint)

	return strings.Join(sections, "\n")
}

func (sm *SpaceMenu) renderConfigMenu() string {
	var sections []string

	title := TitleStyle.Render("Config Menu")
	sections = append(sections, title)

	commands := sm.getConfigCommands()
	for _, cmd := range commands {
		keyLabel := fmt.Sprintf("%s: %s", cmd.Key, cmd.Label)
		line := fmt.Sprintf("  %s  %s",
			KeyStyle.Render(fmt.Sprintf("%-15s", keyLabel)),
			DescriptionStyle.Render(cmd.Description),
		)
		sections = append(sections, line)
	}

	hint := HintStyle.Render("Press Space, Esc, or q to go back")
	sections = append(sections, hint)

	return strings.Join(sections, "\n")
}

func (sm *SpaceMenu) getMainCommands() []Command {
	return []Command{
		{
			Key:         "c",
			Label:       "Config",
			Description: "Area management",
			Action:      ActionConfig,
		},
		{
			Key:         "q",
			Label:       "Quit",
			Description: "Exit application",
			Action:      ActionQuit,
		},
	}
}

func (sm *SpaceMenu) getConfigCommands() []Command {
	return []Command{
		{
			Key:         "n",
			Label:       "New Area",
			Description: "Create a new area",
			Action:      ActionCreateArea,
		},
		{
			Key:         "e",
			Label:       "Edit Area",
			Description: "Edit current area",
			Action:      ActionEditArea,
		},
		{
			Key:         "d",
			Label:       "Delete Area",
			Description: "Delete current area",
			Action:      ActionDeleteArea,
		},
	}
}

func (sm *SpaceMenu) getModalWidth() int {
	minWidth := 45
	maxWidth := 60

	percentageWidth := sm.width * 50 / 100

	if percentageWidth < minWidth {
		return minWidth
	}
	if percentageWidth > maxWidth {
		return maxWidth
	}

	return percentageWidth
}

func (sm *SpaceMenu) State() MenuState {
	return sm.state
}
