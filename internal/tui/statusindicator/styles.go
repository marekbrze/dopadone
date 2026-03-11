package statusindicator

import "github.com/charmbracelet/lipgloss"

const (
	SymbolConnected = "●"
	SymbolSyncing   = "◐"
	SymbolOffline   = "○"
	SymbolLocalOnly = "■"
)

var (
	ColorGreen  = lipgloss.Color("#10B981")
	ColorYellow = lipgloss.Color("#F59E0B")
	ColorRed    = lipgloss.Color("#EF4444")
	ColorGray   = lipgloss.Color("#6B7280")
)
