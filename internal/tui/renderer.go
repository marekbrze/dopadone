package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/example/dopadone/internal/domain"
	"github.com/example/dopadone/internal/tui/tree"
	"github.com/example/dopadone/internal/tui/views"
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

	if m.groupedTasks == nil || m.groupedTasks.TotalCount == 0 {
		return EmptyStateNoTasks
	}

	var lines []string
	taskIndex := 0

	for _, task := range m.groupedTasks.DirectTasks {
		lines = append(lines, m.renderTaskLine(task, taskIndex, 0))
		taskIndex++
	}

	if len(m.groupedTasks.DirectTasks) > 0 && len(m.groupedTasks.Groups) > 0 {
		lines = append(lines, "")
	}

	for _, group := range m.groupedTasks.Groups {
		lines = append(lines, m.renderGroupHeader(group))

		if group.IsExpanded {
			for _, task := range group.Tasks {
				lines = append(lines, m.renderTaskLine(task, taskIndex, 1))
				taskIndex++
			}
		}
	}

	return joinLines(lines)
}

func (m *Model) renderGroupHeader(group domain.TaskGroup) string {
	icon := "▾"
	if !group.IsExpanded {
		icon = "▸"
	}

	taskCount := len(group.Tasks)
	taskLabel := "task"
	if taskCount != 1 {
		taskLabel = "tasks"
	}

	header := fmt.Sprintf("%s %s (%d %s)", icon, group.ProjectName, taskCount, taskLabel)

	style := lipgloss.NewStyle().
		Foreground(m.theme.Secondary).
		PaddingLeft(1)

	return style.Render(header)
}

func (m *Model) renderTaskLine(task domain.Task, index int, indentLevel int) string {
	indent := strings.Repeat("  ", indentLevel+1)

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

	if index == m.selectedTaskIndex {
		style = style.Bold(true).Reverse(true)
	}

	maxTextWidth := m.taskColumnWidth() - len(indent) - len(prefix) - 4
	if maxTextWidth > 0 && len(text) > maxTextWidth {
		text = text[:maxTextWidth-1] + "…"
	}

	return indent + style.Render(prefix+text)
}

func (m *Model) taskColumnWidth() int {
	subareasWeight := 1
	projectsWeight := 1
	tasksWeight := 2
	totalWeight := subareasWeight + projectsWeight + tasksWeight

	availableWidth := m.width - (views.ColumnGap * 3)

	subareasWidth := (availableWidth * subareasWeight) / totalWeight
	projectsWidth := (availableWidth * projectsWeight) / totalWeight
	tasksWidth := availableWidth - subareasWidth - projectsWidth

	if tasksWidth < views.MinTasksWidth {
		tasksWidth = views.MinTasksWidth
	}

	return tasksWidth
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
