package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type ColorTheme struct {
	Primary    lipgloss.AdaptiveColor
	Secondary  lipgloss.AdaptiveColor
	Success    lipgloss.AdaptiveColor
	Error      lipgloss.AdaptiveColor
	Warning    lipgloss.AdaptiveColor
	Muted      lipgloss.AdaptiveColor
	Dimmed     lipgloss.AdaptiveColor
	Background lipgloss.AdaptiveColor
	Foreground lipgloss.AdaptiveColor
}

var Default = ColorTheme{
	Primary: lipgloss.AdaptiveColor{
		Light: "#0066CC",
		Dark:  "#4D9FFF",
	},
	Secondary: lipgloss.AdaptiveColor{
		Light: "#6B7280",
		Dark:  "#9CA3AF",
	},
	Success: lipgloss.AdaptiveColor{
		Light: "#059669",
		Dark:  "#10B981",
	},
	Error: lipgloss.AdaptiveColor{
		Light: "#DC2626",
		Dark:  "#EF4444",
	},
	Warning: lipgloss.AdaptiveColor{
		Light: "#D97706",
		Dark:  "#F59E0B",
	},
	Muted: lipgloss.AdaptiveColor{
		Light: "#9CA3AF",
		Dark:  "#6B7280",
	},
	Dimmed: lipgloss.AdaptiveColor{
		Light: "#D1D5DB",
		Dark:  "#374151",
	},
	Background: lipgloss.AdaptiveColor{
		Light: "#FFFFFF",
		Dark:  "#1F2937",
	},
	Foreground: lipgloss.AdaptiveColor{
		Light: "#1F2937",
		Dark:  "#F9FAFB",
	},
}

func (t ColorTheme) TabActiveBackground() lipgloss.AdaptiveColor {
	return t.Primary
}

func (t ColorTheme) TabActiveForeground() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: "#FFFFFF",
		Dark:  "#FFFFFF",
	}
}

func (t ColorTheme) TabInactiveBackground() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: "#E5E7EB",
		Dark:  "#374151",
	}
}

func (t ColorTheme) TabInactiveForeground() lipgloss.AdaptiveColor {
	return lipgloss.AdaptiveColor{
		Light: "#6B7280",
		Dark:  "#9CA3AF",
	}
}

func (t ColorTheme) ColumnFocusedBorder() lipgloss.AdaptiveColor {
	return t.Primary
}

func (t ColorTheme) ColumnUnfocusedBorder() lipgloss.AdaptiveColor {
	return t.Dimmed
}

func (t ColorTheme) ColumnHeader() lipgloss.AdaptiveColor {
	return t.Success
}

func (t ColorTheme) EmptyText() lipgloss.AdaptiveColor {
	return t.Muted
}

func (t ColorTheme) FooterForeground() lipgloss.AdaptiveColor {
	return t.Secondary
}

func (t ColorTheme) FooterBackground() lipgloss.AdaptiveColor {
	return t.Dimmed
}
