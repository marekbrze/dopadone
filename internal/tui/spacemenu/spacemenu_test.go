package spacemenu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSpaceMenu_New(t *testing.T) {
	sm := New()
	if sm == nil {
		t.Fatal("New() returned nil")
	}
	if sm.State() != StateMain {
		t.Errorf("expected initial state to be StateMain, got %v", sm.State())
	}
}

func TestSpaceMenu_WindowSize(t *testing.T) {
	sm := New()
	updated, cmd := sm.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if cmd != nil {
		t.Errorf("expected no command on window size update, got %v", cmd)
	}
	if updated.width != 100 {
		t.Errorf("expected width 100, got %d", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("expected height 50, got %d", updated.height)
	}
}

func TestSpaceMenu_CloseKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"escape closes", "esc"},
		{"space closes", " "},
		{"q closes from main", "q"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New()
			updated, cmd := sm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})

			if cmd == nil {
				t.Errorf("expected command for key %q, got nil", tt.key)
				return
			}

			msg := cmd()
			if _, ok := msg.(CloseMsg); !ok {
				t.Errorf("expected CloseMsg for key %q, got %T", tt.key, msg)
			}

			_ = updated
		})
	}
}

func TestSpaceMenu_Navigation(t *testing.T) {
	t.Run("c key opens config menu", func(t *testing.T) {
		sm := New()
		updated, _ := sm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})

		if updated.State() != StateConfig {
			t.Errorf("expected state StateConfig, got %v", updated.State())
		}
	})

	t.Run("q key goes back from config to main", func(t *testing.T) {
		sm := New()
		sm.state = StateConfig

		updated, _ := sm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})

		if updated.State() != StateMain {
			t.Errorf("expected state StateMain, got %v", updated.State())
		}
	})

	t.Run("escape goes back from config to main", func(t *testing.T) {
		sm := New()
		sm.state = StateConfig

		updated, _ := sm.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")})

		if updated.State() != StateMain {
			t.Errorf("expected state StateMain, got %v", updated.State())
		}
	})
}

func TestSpaceMenu_View(t *testing.T) {
	t.Run("renders main menu", func(t *testing.T) {
		sm := New()
		sm.width = 80
		sm.height = 24

		view := sm.View()

		if view == "" {
			t.Error("View() returned empty string")
		}

		expectedStrings := []string{"Command Menu", "c: Config", "Area management", "q: Quit", "Exit application"}
		for _, expected := range expectedStrings {
			if !contains(view, expected) {
				t.Errorf("View() missing expected string %q", expected)
			}
		}
	})

	t.Run("renders config menu", func(t *testing.T) {
		sm := New()
		sm.width = 80
		sm.height = 24
		sm.state = StateConfig

		view := sm.View()

		if view == "" {
			t.Error("View() returned empty string")
		}

		expectedStrings := []string{"Config Menu", "n: New Area", "Create a new area", "e: Edit Area", "Edit current area", "d: Delete Area", "Delete current area"}
		for _, expected := range expectedStrings {
			if !contains(view, expected) {
				t.Errorf("View() missing expected string %q", expected)
			}
		}
	})
}

func TestSpaceMenu_GetModalWidth(t *testing.T) {
	tests := []struct {
		name        string
		windowWidth int
		expectedMin int
		expectedMax int
	}{
		{"small window", 50, 45, 45},
		{"medium window", 100, 45, 50},
		{"large window", 200, 45, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := New()
			sm.width = tt.windowWidth

			width := sm.getModalWidth()

			if width < tt.expectedMin || width > tt.expectedMax {
				t.Errorf("getModalWidth() = %d, want between %d and %d", width, tt.expectedMin, tt.expectedMax)
			}
		})
	}
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
