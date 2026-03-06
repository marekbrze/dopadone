package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/example/dopadone/internal/tui/tree"
)

func (m *Model) RenderSubareas() string {
	if m.isLoadingSubareas {
		return m.spinner.View() + " " + LoadingMessageSubareas
	}

	if len(m.subareas) == 0 {
		return EmptyStateNoSubareas
	}

	var lines []string
	for i, subarea := range m.subareas {
		if i == m.selectedSubareaIndex {
			lines = append(lines, "  "+m.renderSelectedLine(subarea.Name))
		} else {
			lines = append(lines, "  "+subarea.Name)
		}
	}
	return joinLines(lines)
}

func (m *Model) RenderProjects() string {
	if m.isLoadingProjects {
		return m.spinner.View() + " " + LoadingMessageProjects
	}

	if m.projectTree == nil {
		return EmptyStateNoProjects
	}

	visibleNodes := tree.GetAllVisibleNodes(m.projectTree)
	if len(visibleNodes) == 0 {
		return EmptyStateNoProjects
	}

	renderer := tree.NewRenderer()
	return renderer.Render(m.projectTree, m.selectedProjectID)
}

func (m *Model) RenderTasks() string {
	if m.isLoadingTasks {
		return m.spinner.View() + " " + LoadingMessageTasks
	}

	if len(m.tasks) == 0 {
		return EmptyStateNoTasks
	}

	var lines []string
	for i, task := range m.tasks {
		var prefix string
		var text string
		var style lipgloss.Style

		if task.IsCompleted() {
			prefix = "✓ "
			text = task.Title
			style = lipgloss.NewStyle().
				Strikethrough(true).
				Foreground(m.theme.Muted)
		} else {
			prefix = "  "
			text = task.Title
			style = lipgloss.NewStyle().
				Foreground(m.theme.Foreground)
		}

		if i == m.selectedTaskIndex {
			style = style.Bold(true).Reverse(true)
		}

		lines = append(lines, "  "+style.Render(prefix+text))
	}
	return joinLines(lines)
}

func (m *Model) renderSelectedLine(text string) string {
	return lipgloss.NewStyle().
		Bold(true).
		Reverse(true).
		Render(text)
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func overlay(background, foreground string, width, height int) string {
	return lipgloss.JoinVertical(lipgloss.Center, foreground)
}
