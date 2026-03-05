package views

import (
	"github.com/charmbracelet/lipgloss"
)

type Column struct {
	Title     string
	Content   string
	IsFocused bool
	Width     int
	Height    int
}

func ColumnView(col Column) string {
	var style lipgloss.Style
	if col.IsFocused {
		style = FocusedColumnStyle
	} else {
		style = UnfocusedColumnStyle
	}

	header := ColumnHeaderStyle.Render(col.Title)

	content := col.Content
	if content == "" {
		content = EmptyContentStyle.Render("No items")
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left, header, content)

	if col.Width > 0 && col.Height > 0 {
		return style.Width(col.Width).Height(col.Height).Render(fullContent)
	}

	return style.Render(fullContent)
}

func Layout(columns []Column, width, height int) string {
	if len(columns) != 3 {
		return ""
	}

	tabsHeight := 2
	footerHeight := 2
	availableHeight := height - tabsHeight - footerHeight - 2
	if availableHeight < 5 {
		availableHeight = 5
	}

	columnWidth := (width - 6) / 3
	if columnWidth < 10 {
		columnWidth = 10
	}

	for i := range columns {
		columns[i].Width = columnWidth
		columns[i].Height = availableHeight
	}

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
