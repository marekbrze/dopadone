package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	MinSubareasWidth = 20
	MinProjectsWidth = 20
	MinTasksWidth    = 40

	ColumnGap = 2

	StackedLayoutThreshold = 120
)

type Column struct {
	Title     string
	Content   string
	IsFocused bool
	Width     int
	Height    int
}

func calculateColumnWidths(totalWidth int) (int, int, int) {
	subareasWeight := 1
	projectsWeight := 1
	tasksWeight := 2
	totalWeight := subareasWeight + projectsWeight + tasksWeight

	availableWidth := totalWidth - (ColumnGap * 3)

	subareasWidth := (availableWidth * subareasWeight) / totalWeight
	projectsWidth := (availableWidth * projectsWeight) / totalWeight
	tasksWidth := availableWidth - subareasWidth - projectsWidth

	if subareasWidth < MinSubareasWidth {
		subareasWidth = MinSubareasWidth
	}
	if projectsWidth < MinProjectsWidth {
		projectsWidth = MinProjectsWidth
	}
	if tasksWidth < MinTasksWidth {
		tasksWidth = MinTasksWidth
	}

	return subareasWidth, projectsWidth, tasksWidth
}

func shouldUseStackedLayout(width int) bool {
	return width < StackedLayoutThreshold
}

func calculateStackedLayoutWidths(totalWidth int) (int, int) {
	availableWidth := totalWidth - ColumnGap

	leftWidth := availableWidth / 4
	tasksWidth := availableWidth - leftWidth

	return leftWidth, tasksWidth
}

func calculateStackedLayoutHeights(totalHeight int) (int, int) {
	availableHeight := totalHeight - 2

	subareasHeight := availableHeight / 2
	projectsHeight := availableHeight - subareasHeight

	return subareasHeight, projectsHeight
}

func stripANSI(s string) string {
	var result []rune
	var inEscape bool
	var escapeLen int

	for _, r := range s {
		if r == '\x1b' && !inEscape {
			inEscape = true
			escapeLen = 0
			continue
		}

		if inEscape {
			escapeLen++
			if escapeLen > 1 && (r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
				inEscape = false
			} else if escapeLen > 2 && (r >= '0' && r <= '9') {
				inEscape = false
			}
			continue
		}

		result = append(result, r)
	}

	return string(result)
}

func truncateString(s string, maxLen int) string {
	if s == "" {
		return ""
	}
	if maxLen <= 1 {
		return "…"
	}

	stripped := stripANSI(s)
	if len(stripped) <= maxLen {
		return s
	}

	if len(stripped) > maxLen {
		return "…"
	}

	return s[:maxLen-1] + "…"
}

func ColumnView(col Column) string {
	var style lipgloss.Style
	if col.IsFocused {
		style = FocusedColumnStyle
	} else {
		style = UnfocusedColumnStyle
	}

	maxTextWidth := col.Width - 4
	if maxTextWidth < 1 {
		maxTextWidth = 1
	}

	header := ColumnHeaderStyle.Render(truncateString(col.Title, maxTextWidth))

	content := col.Content
	if content == "" {
		content = EmptyContentStyle.Render("No items")
	} else {
		lines := strings.Split(content, "\n")
		truncatedLines := make([]string, len(lines))
		for i, line := range lines {
			truncatedLines[i] = truncateString(line, maxTextWidth)
		}
		content = strings.Join(truncatedLines, "\n")
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, header, content)

	if col.Width > 0 && col.Height > 0 {
		return style.Width(col.Width).Height(col.Height).Render(fullContent)
	}

	return style.Render(fullContent)
}

func LayoutStacked(columns []Column, width, height int) string {
	if len(columns) != 3 {
		return ""
	}

	tabsHeight := 2
	footerHeight := 2
	availableHeight := height - tabsHeight - footerHeight - 2
	if availableHeight < 5 {
		availableHeight = 5
	}

	leftWidth, tasksWidth := calculateStackedLayoutWidths(width)

	subareasHeight, projectsHeight := calculateStackedLayoutHeights(availableHeight)

	columns[0].Width = leftWidth
	columns[0].Height = subareasHeight
	columns[1].Width = leftWidth
	columns[1].Height = projectsHeight
	columns[2].Width = tasksWidth
	columns[2].Height = availableHeight

	stackedLeft := lipgloss.JoinVertical(
		lipgloss.Left,
		ColumnView(columns[0]),
		ColumnView(columns[1]),
	)

	tasksColumn := ColumnView(columns[2])

	return lipgloss.JoinHorizontal(lipgloss.Top, stackedLeft, tasksColumn)
}

func Layout(columns []Column, width, height int) string {
	if len(columns) != 3 {
		return ""
	}

	if shouldUseStackedLayout(width) {
		return LayoutStacked(columns, width, height)
	}

	tabsHeight := 2
	footerHeight := 2
	availableHeight := height - tabsHeight - footerHeight - 2
	if availableHeight < 5 {
		availableHeight = 5
	}

	subareasWidth, projectsWidth, tasksWidth := calculateColumnWidths(width)

	columns[0].Width = subareasWidth
	columns[0].Height = availableHeight
	columns[1].Width = projectsWidth
	columns[1].Height = availableHeight
	columns[2].Width = tasksWidth
	columns[2].Height = availableHeight

	renderedColumns := make([]string, 3)
	for i, col := range columns {
		renderedColumns[i] = ColumnView(col)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...)
}

func LayoutWithTabs(tabs string, columns []Column, width, height int) string {
	columnLayout := Layout(columns, width, height)

	return lipgloss.JoinVertical(lipgloss.Left, tabs, columnLayout)
}
