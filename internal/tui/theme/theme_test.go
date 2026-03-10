package theme

import (
	"testing"
)

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name    string
		mode    ThemeMode
		wantErr bool
	}{
		{
			name:    "auto theme",
			mode:    ThemeAuto,
			wantErr: false,
		},
		{
			name:    "light theme",
			mode:    ThemeLight,
			wantErr: false,
		},
		{
			name:    "dark theme",
			mode:    ThemeDark,
			wantErr: false,
		},
		{
			name:    "empty string defaults to auto",
			mode:    "",
			wantErr: false,
		},
		{
			name:    "invalid theme returns default with error",
			mode:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := GetTheme(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTheme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if theme.Primary.Light == "" || theme.Primary.Dark == "" {
				t.Error("GetTheme() returned theme with empty primary colors")
			}
		})
	}
}

func TestThemeMethods(t *testing.T) {
	theme := Default

	tests := []struct {
		name     string
		getColor func() interface{}
	}{
		{
			name:     "TabActiveBackground",
			getColor: func() interface{} { return theme.TabActiveBackground() },
		},
		{
			name:     "TabActiveForeground",
			getColor: func() interface{} { return theme.TabActiveForeground() },
		},
		{
			name:     "TabInactiveBackground",
			getColor: func() interface{} { return theme.TabInactiveBackground() },
		},
		{
			name:     "TabInactiveForeground",
			getColor: func() interface{} { return theme.TabInactiveForeground() },
		},
		{
			name:     "ColumnFocusedBorder",
			getColor: func() interface{} { return theme.ColumnFocusedBorder() },
		},
		{
			name:     "ColumnUnfocusedBorder",
			getColor: func() interface{} { return theme.ColumnUnfocusedBorder() },
		},
		{
			name:     "ColumnHeader",
			getColor: func() interface{} { return theme.ColumnHeader() },
		},
		{
			name:     "EmptyText",
			getColor: func() interface{} { return theme.EmptyText() },
		},
		{
			name:     "FooterForeground",
			getColor: func() interface{} { return theme.FooterForeground() },
		},
		{
			name:     "FooterBackground",
			getColor: func() interface{} { return theme.FooterBackground() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color := tt.getColor()
			if color == nil {
				t.Errorf("%s() returned nil", tt.name)
			}
		})
	}
}

func TestThemeColorValues(t *testing.T) {
	theme := Default

	if theme.Primary.Light == "" {
		t.Error("Primary.Light color is empty")
	}
	if theme.Primary.Dark == "" {
		t.Error("Primary.Dark color is empty")
	}
	if theme.Success.Light == "" {
		t.Error("Success.Light color is empty")
	}
	if theme.Error.Light == "" {
		t.Error("Error.Light color is empty")
	}
	if theme.Warning.Light == "" {
		t.Error("Warning.Light color is empty")
	}
}

func TestLightTheme(t *testing.T) {
	theme, err := GetTheme(ThemeLight)
	if err != nil {
		t.Fatalf("GetTheme(ThemeLight) error = %v", err)
	}

	if theme.Primary.Light == "" {
		t.Error("Light theme primary color is empty")
	}
}

func TestDarkTheme(t *testing.T) {
	theme, err := GetTheme(ThemeDark)
	if err != nil {
		t.Fatalf("GetTheme(ThemeDark) error = %v", err)
	}

	if theme.Primary.Dark == "" {
		t.Error("Dark theme primary color is empty")
	}
}
