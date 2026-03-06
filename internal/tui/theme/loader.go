package theme

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme string `yaml:"theme"`
}

type ThemeMode string

const (
	ThemeAuto  ThemeMode = "auto"
	ThemeLight ThemeMode = "light"
	ThemeDark  ThemeMode = "dark"
)

func LoadTheme(configPath string) (ColorTheme, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return Default, nil
	}

	return GetTheme(ThemeMode(config.Theme))
}

func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func GetTheme(mode ThemeMode) (ColorTheme, error) {
	switch mode {
	case ThemeAuto, "":
		return Default, nil
	case ThemeLight:
		return getLightTheme(), nil
	case ThemeDark:
		return getDarkTheme(), nil
	default:
		return Default, fmt.Errorf("unknown theme mode: %s, using default", mode)
	}
}

func getLightTheme() ColorTheme {
	return ColorTheme{
		Primary:    lipgloss.AdaptiveColor{Light: "#0066CC", Dark: "#0066CC"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"},
		Success:    lipgloss.AdaptiveColor{Light: "#059669", Dark: "#059669"},
		Error:      lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#DC2626"},
		Warning:    lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#D97706"},
		Muted:      lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#9CA3AF"},
		Dimmed:     lipgloss.AdaptiveColor{Light: "#D1D5DB", Dark: "#D1D5DB"},
		Background: lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#FFFFFF"},
		Foreground: lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#1F2937"},
	}
}

func getDarkTheme() ColorTheme {
	return ColorTheme{
		Primary:    lipgloss.AdaptiveColor{Light: "#4D9FFF", Dark: "#4D9FFF"},
		Secondary:  lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#9CA3AF"},
		Success:    lipgloss.AdaptiveColor{Light: "#10B981", Dark: "#10B981"},
		Error:      lipgloss.AdaptiveColor{Light: "#EF4444", Dark: "#EF4444"},
		Warning:    lipgloss.AdaptiveColor{Light: "#F59E0B", Dark: "#F59E0B"},
		Muted:      lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#6B7280"},
		Dimmed:     lipgloss.AdaptiveColor{Light: "#374151", Dark: "#374151"},
		Background: lipgloss.AdaptiveColor{Light: "#1F2937", Dark: "#1F2937"},
		Foreground: lipgloss.AdaptiveColor{Light: "#F9FAFB", Dark: "#F9FAFB"},
	}
}
