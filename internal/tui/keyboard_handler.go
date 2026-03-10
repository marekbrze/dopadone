package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/spacemenu"
)

const (
	KeySpace        = " "
	KeyH            = "h"
	KeyL            = "l"
	KeyLeft         = "left"
	KeyRight        = "right"
	KeyTab          = "tab"
	KeyJ            = "j"
	KeyK            = "k"
	KeyDown         = "down"
	KeyUp           = "up"
	KeyX            = "x"
	KeyA            = "a"
	KeyD            = "d"
	KeyHelp         = "?"
	KeyCtrlA        = "ctrl+a"
	KeyBracketOpen  = "["
	KeyBracketClose = "]"
	KeyQ            = "q"
	KeyCtrlC        = "ctrl+c"
	KeyEnter        = "enter"
)

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.isHelpOpen {
		return m.handleHelpModalKeyPress(msg)
	}

	if m.isModalOpen {
		return m.handleGenericModalKeyPress(msg)
	}

	if m.isAreaModalOpen && m.areaModal != nil {
		return m.handleAreaModalKeyPress(msg)
	}

	if m.isSpaceMenuOpen && m.spaceMenu != nil {
		return m.handleSpaceMenuKeyPress(msg)
	}

	if m.isConfirmModalOpen && m.confirmModal != nil {
		return m.handleConfirmModalKeyPress(msg)
	}

	return m.handleMainKeyPress(msg)
}

func (m *Model) handleHelpModalKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.helpModal, cmd = m.helpModal.Update(msg)
	return m, cmd
}

func (m *Model) handleGenericModalKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.modal, cmd = m.modal.Update(msg)
	return m, cmd
}

func (m *Model) handleAreaModalKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.areaModal, cmd = m.areaModal.Update(msg)
	return m, cmd
}

func (m *Model) handleSpaceMenuKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		if m.spaceMenu != nil && m.spaceMenu.State() == 0 {
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.spaceMenu, cmd = m.spaceMenu.Update(msg)
	return m, cmd
}

func (m *Model) handleConfirmModalKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	}
	var cmd tea.Cmd
	m.confirmModal, cmd = m.confirmModal.Update(msg)
	return m, cmd
}

func (m *Model) handleMainKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case KeyQ, KeyCtrlC:
		return m, tea.Quit
	case KeySpace:
		if !m.isModalOpen && !m.isAreaModalOpen && !m.isHelpOpen {
			m.isSpaceMenuOpen = true
			m.spaceMenu = spacemenu.New()
			m.spaceMenu, _ = m.spaceMenu.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m, nil
		}
		m.handleEnterOrSpace()
	case KeyH, KeyLeft:
		m.focus = m.focus.Prev()
	case KeyL, KeyRight:
		m.focus = m.focus.Next()
	case KeyTab:
		m.focus = m.focus.Next()
	case KeyJ, KeyDown:
		if !m.IsEmpty(m.focus) {
			return m.NavigateDownWithLoad(m.focus)
		}
	case KeyK, KeyUp:
		if !m.IsEmpty(m.focus) {
			return m.NavigateUpWithLoad(m.focus)
		}
	case KeyBracketOpen:
		return m, m.SwitchToPreviousArea()
	case KeyBracketClose:
		return m, m.SwitchToNextArea()
	case KeyEnter:
		m.handleEnterOrSpace()
	case KeyX:
		if m.focus == FocusTasks && (len(m.tasks) > 0 || (m.groupedTasks != nil && m.groupedTasks.TotalCount > 0)) {
			return m, m.toggleTaskCompletion()
		}
	case KeyA:
		return m.handleQuickAdd()
	case KeyD:
		return m, m.handleDeleteKey()
	case KeyHelp:
		return m.handleHelp(), nil
	case KeyCtrlA:
		return m.handleOpenAreaModal()
	}
	return m, nil
}
