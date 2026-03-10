package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/help"
	"github.com/marekbrze/dopadone/internal/tui/mocks"
	"github.com/marekbrze/dopadone/internal/tui/spacemenu"
)

func TestSpaceMenuIntegration(t *testing.T) {
	t.Run("Space key opens spacemenu when no modal is open", func(t *testing.T) {
		areaSvc, subareaSvc, projectSvc, taskSvc := mocks.NewMockServices()
		model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		model = *newModel.(*Model)

		if model.isSpaceMenuOpen {
			t.Error("spacemenu should not be open initially")
		}

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
		model = *newModel.(*Model)

		if !model.isSpaceMenuOpen {
			t.Error("spacemenu should be open after pressing Space")
		}

		if model.spaceMenu == nil {
			t.Error("spacemenu component should be initialized")
		}
	})

	t.Run("Space key does not open spacemenu when modal is open", func(t *testing.T) {
		areaSvc, subareaSvc, projectSvc, taskSvc := mocks.NewMockServices()
		model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		model = *newModel.(*Model)

		model.isHelpOpen = true
		model.helpModal = help.New()

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
		model = *newModel.(*Model)

		if model.isSpaceMenuOpen {
			t.Error("spacemenu should not open when help modal is open")
		}
	})

	t.Run("spacemenu closes on Space key", func(t *testing.T) {
		areaSvc, subareaSvc, projectSvc, taskSvc := mocks.NewMockServices()
		model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		model = *newModel.(*Model)
		model.isSpaceMenuOpen = true
		model.spaceMenu = spacemenu.New()

		newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" ")})
		model = *newModel.(*Model)

		if cmd == nil {
			t.Fatal("Expected command to be returned when Space is pressed")
		}

		msg := cmd()
		if _, ok := msg.(spacemenu.CloseMsg); !ok {
			t.Errorf("Expected spacemenu.CloseMsg, got %T", msg)
		}

		newModel, _ = model.Update(msg)
		model = *newModel.(*Model)

		if model.isSpaceMenuOpen {
			t.Error("spacemenu should close on Space key")
		}

		if model.spaceMenu != nil {
			t.Error("spacemenu component should be nil after closing")
		}
	})

	t.Run("spacemenu closes on Escape key", func(t *testing.T) {
		areaSvc, subareaSvc, projectSvc, taskSvc := mocks.NewMockServices()
		model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		model = *newModel.(*Model)
		model.isSpaceMenuOpen = true
		model.spaceMenu = spacemenu.New()

		newModel, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		model = *newModel.(*Model)

		if cmd == nil {
			t.Fatal("Expected command to be returned when Escape is pressed")
		}

		msg := cmd()
		if _, ok := msg.(spacemenu.CloseMsg); !ok {
			t.Errorf("Expected spacemenu.CloseMsg, got %T", msg)
		}

		newModel, _ = model.Update(msg)
		model = *newModel.(*Model)

		if model.isSpaceMenuOpen {
			t.Error("spacemenu should close on Escape key")
		}
	})

	t.Run("spacemenu navigates to config on c key", func(t *testing.T) {
		areaSvc, subareaSvc, projectSvc, taskSvc := mocks.NewMockServices()
		model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

		newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
		model = *newModel.(*Model)
		model.isSpaceMenuOpen = true
		model.spaceMenu = spacemenu.New()

		newModel, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
		model = *newModel.(*Model)

		if model.spaceMenu == nil || model.spaceMenu.State() != spacemenu.StateConfig {
			t.Error("spacemenu should be in Config state after pressing 'c'")
		}
	})
}
