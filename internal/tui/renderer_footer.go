package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/example/projectdb/internal/tui/toast"
)

func (m *Model) RenderToasts() string {
	if len(m.toasts) == 0 {
		return ""
	}

	var lines []string
	for _, t := range m.toasts {
		var line string
		switch t.Type {
		case toast.TypeError:
			line = toast.ErrorStyle.Render("✗ " + t.Message)
		case toast.TypeSuccess:
			line = toast.SuccessStyle.Render("✓ " + t.Message)
		case toast.TypeInfo:
			line = toast.InfoStyle.Render("ℹ " + t.Message)
		}
		lines = append(lines, line)
	}

	return joinLines(lines)
}

func (m *Model) RenderFooter() string {
	if !m.ready {
		return ""
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Background(lipgloss.Color("235")).
		Padding(0, 1)

	shortcuts := "h/l: columns | j/k: navigate | a: add | ?: help | q: quit"
	return footerStyle.Render(shortcuts)
}
