package views

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Tab struct {
	Name     string
	ID       string
	IsActive bool
}

func TabsView(tabs []Tab, selectedIndex int) string {
	if len(tabs) == 0 {
		return ""
	}

	var renderedTabs []string
	for i, tab := range tabs {
		var style lipgloss.Style
		if i == selectedIndex {
			style = ActiveTabStyle
		} else {
			style = InactiveTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(tab.Name))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	return tabRow
}

func TabsWithSeparator(tabs []Tab, selectedIndex int) string {
	if len(tabs) == 0 {
		return ""
	}

	var parts []string
	for i, tab := range tabs {
		var style lipgloss.Style
		if i == selectedIndex {
			style = ActiveTabStyle
		} else {
			style = InactiveTabStyle
		}
		parts = append(parts, style.Render(tab.Name))
		if i < len(tabs)-1 {
			parts = append(parts, TabSeparator.String())
		}
	}

	return strings.Join(parts, " ")
}
