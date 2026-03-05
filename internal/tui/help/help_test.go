package help

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	h := New()
	if h == nil {
		t.Fatal("New() returned nil")
	}
	if h.width != 0 {
		t.Errorf("expected width 0, got %d", h.width)
	}
	if h.height != 0 {
		t.Errorf("expected height 0, got %d", h.height)
	}
}

func TestHelpModalUpdateWindowSize(t *testing.T) {
	h := New()
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	updated, cmd := h.Update(msg)
	if cmd != nil {
		t.Errorf("expected nil cmd, got %v", cmd)
	}
	if updated.width != 100 {
		t.Errorf("expected width 100, got %d", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("expected height 50, got %d", updated.height)
	}
}

func TestHelpModalUpdateCloseKeys(t *testing.T) {
	keys := []string{"?", "esc", "q"}

	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			h := New()
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}

			_, cmd := h.Update(msg)
			if cmd == nil {
				t.Errorf("expected non-nil cmd for key %s", key)
			}

			result := cmd()
			if _, ok := result.(CloseMsg); !ok {
				t.Errorf("expected CloseMsg, got %T", result)
			}
		})
	}
}

func TestHelpModalView(t *testing.T) {
	h := New()
	h.width = 100
	h.height = 50

	view := h.View()
	if view == "" {
		t.Error("View() returned empty string")
	}

	expectedStrings := []string{
		"Keyboard Shortcuts",
		"Navigation",
		"h, ←",
		"l, →",
		"j, ↓",
		"k, ↑",
		"Actions",
		"a",
		"General",
		"?",
	}

	for _, expected := range expectedStrings {
		if !contains(view, expected) {
			t.Errorf("View() missing expected string: %s", expected)
		}
	}
}

func TestHelpModalGetModalWidth(t *testing.T) {
	tests := []struct {
		width       int
		expectedMin int
		expectedMax int
	}{
		{0, 50, 70},
		{40, 50, 70},
		{100, 50, 70},
		{200, 70, 70},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			h := New()
			h.width = tt.width
			got := h.getModalWidth()
			if got < tt.expectedMin || got > tt.expectedMax {
				t.Errorf("getModalWidth() = %d, want between %d and %d", got, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func TestHelpModalGetCategories(t *testing.T) {
	h := New()
	categories := h.getCategories()

	if len(categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(categories))
	}

	expectedCategories := []string{"Navigation", "Actions", "General"}
	for i, expected := range expectedCategories {
		if categories[i].Name != expected {
			t.Errorf("category %d: expected %s, got %s", i, expected, categories[i].Name)
		}
	}

	for _, cat := range categories {
		if len(cat.Shortcuts) == 0 {
			t.Errorf("category %s has no shortcuts", cat.Name)
		}
		for _, shortcut := range cat.Shortcuts {
			if shortcut.Key == "" {
				t.Errorf("category %s has empty key", cat.Name)
			}
			if shortcut.Description == "" {
				t.Errorf("category %s has empty description for key %s", cat.Name, shortcut.Key)
			}
		}
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
