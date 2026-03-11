package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/statusindicator"
	"github.com/marekbrze/dopadone/internal/tui/theme"
	"github.com/marekbrze/dopadone/internal/tui/toast"
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
		Foreground(theme.Default.FooterForeground()).
		Background(theme.Default.FooterBackground()).
		Padding(0, 1)

	indicator := statusindicator.New(
		m.connectionStatus.Mode,
		m.connectionStatus.Status,
		m.connectionStatus.SyncStatus,
	)

	statusPart := indicator.Render()
	shortcuts := "h/l: columns | j/k: navigate | a: add | d: delete | ?: help | q: quit"

	footer := fmt.Sprintf("%s | %s", statusPart, shortcuts)
	return footerStyle.Render(footer)
}
