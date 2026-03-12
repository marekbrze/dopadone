package welcome

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/internal/constants"
)

func TestNew(t *testing.T) {
	m := New()

	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.input.Value() != "" {
		t.Errorf("expected empty input value, got %q", m.input.Value())
	}
	if m.colorIndex != 0 {
		t.Errorf("expected colorIndex 0, got %d", m.colorIndex)
	}
	if m.errorMsg != "" {
		t.Errorf("expected empty error message, got %q", m.errorMsg)
	}
}

func TestViewInit(t *testing.T) {
	m := New()
	cmd := m.Init()
	if cmd != nil {
		t.Errorf("expected nil Init command, got %v", cmd)
	}
}

func TestViewContainsBranding(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	view := m.View()

	expectedStrings := []string{
		"Welcome to Dopadone",
		"Your project management companion",
		"Create your first area to get started",
		"Area Name:",
		"Color",
		"Tab",
		"Enter: Create",
		"ESC: Exit",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("view should contain %q", expected)
		}
	}
}

func TestViewContainsInputField(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24
	view := m.View()

	if !strings.Contains(view, "Enter area name...") {
		t.Error("view should contain input placeholder")
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty name",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "valid name",
			input:   "My Area",
			wantErr: false,
		},
		{
			name:    "name with leading/trailing whitespace",
			input:   "  My Area  ",
			wantErr: false,
		},
		{
			name:    "name at max length",
			input:   strings.Repeat("a", 100),
			wantErr: false,
		},
		{
			name:    "name exceeding max length",
			input:   strings.Repeat("a", 101),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestEnterWithEmptyInputShowsError(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(msg)

	if cmd != nil {
		t.Errorf("expected nil command for empty input, got %v", cmd)
	}
	if updated.errorMsg == "" {
		t.Error("expected error message for empty input")
	}

	view := updated.View()
	if !strings.Contains(view, "✗") {
		t.Error("view should contain error indicator")
	}
}

func TestColorCycling(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	initialIndex := m.colorIndex

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if updated.colorIndex != (initialIndex+1)%len(PredefinedColors) {
		t.Errorf("expected colorIndex %d, got %d", (initialIndex+1)%len(PredefinedColors), updated.colorIndex)
	}

	for i := 0; i < len(PredefinedColors)-1; i++ {
		updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	if updated.colorIndex != initialIndex {
		t.Errorf("expected colorIndex to wrap around to %d, got %d", initialIndex, updated.colorIndex)
	}
}

func TestShiftTabColorCycling(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	m.colorIndex = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})

	expectedIndex := len(PredefinedColors) - 1
	if updated.colorIndex != expectedIndex {
		t.Errorf("expected colorIndex %d, got %d", expectedIndex, updated.colorIndex)
	}
}

func TestESCBehavior(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if cmd == nil {
		t.Fatal("expected command for ESC key")
	}

	msg := cmd()
	if _, ok := msg.(ExitMsg); !ok {
		t.Errorf("expected ExitMsg, got %T", msg)
	}
}

func TestEnterWithValidInputReturnsSubmitMsg(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	m.input.SetValue("My First Area")

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected command for valid input")
	}

	msg := cmd()
	submitMsg, ok := msg.(SubmitMsg)
	if !ok {
		t.Fatalf("expected SubmitMsg, got %T", msg)
	}

	if submitMsg.Name != "My First Area" {
		t.Errorf("expected name %q, got %q", "My First Area", submitMsg.Name)
	}
	if submitMsg.Color != PredefinedColors[0] {
		t.Errorf("expected color %s, got %s", PredefinedColors[0], submitMsg.Color)
	}
}

func TestEnterPreservesSelectedColor(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	for i := 0; i < 3; i++ {
		m.Update(tea.KeyMsg{Type: tea.KeyTab})
	}
	expectedColor := PredefinedColors[m.colorIndex]

	m.input.SetValue("Test Area")

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	msg := cmd()
	submitMsg := msg.(SubmitMsg)

	if submitMsg.Color != expectedColor {
		t.Errorf("expected color %s, got %s", expectedColor, submitMsg.Color)
	}
}

func TestErrorClearsOnInput(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.errorMsg == "" {
		t.Fatal("expected error message after empty submit")
	}

	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	if m.errorMsg != "" {
		t.Errorf("expected error to be cleared on input, got %q", m.errorMsg)
	}
}

func TestWindowSizeMsg(t *testing.T) {
	m := New()

	updated, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if cmd != nil {
		t.Errorf("expected nil command for WindowSizeMsg, got %v", cmd)
	}
	if updated.width != 100 {
		t.Errorf("expected width 100, got %d", updated.width)
	}
	if updated.height != 50 {
		t.Errorf("expected height 50, got %d", updated.height)
	}
}

func TestSetError(t *testing.T) {
	m := New()
	m.SetError("test error")

	if m.errorMsg != "test error" {
		t.Errorf("expected error %q, got %q", "test error", m.errorMsg)
	}
}

func TestClearError(t *testing.T) {
	m := New()
	m.errorMsg = "existing error"
	m.ClearError()

	if m.errorMsg != "" {
		t.Errorf("expected empty error, got %q", m.errorMsg)
	}
}

func TestPredefinedColorsNotEmpty(t *testing.T) {
	if len(PredefinedColors) == 0 {
		t.Error("PredefinedColors should not be empty")
	}
}

func TestEscKeyString(t *testing.T) {
	m := New()
	m.width = 80
	m.height = 24

	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	updated, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Error("expected command for ESC key")
	}

	msg := cmd()
	if _, ok := msg.(ExitMsg); !ok {
		t.Errorf("expected ExitMsg, got %T", msg)
	}

	if updated.width != 80 {
		t.Errorf("model state should be preserved")
	}
}

func TestEnterKeyString(t *testing.T) {
	m := New()
	m.input.SetValue("Valid Name")

	enterKey := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(enterKey)

	if cmd == nil {
		t.Fatal("expected command for Enter with valid input")
	}

	msg := cmd()
	if _, ok := msg.(SubmitMsg); !ok {
		t.Errorf("expected SubmitMsg, got %T", msg)
	}
}

func TestConstantsMatch(t *testing.T) {
	if constants.KeyEnter != "enter" {
		t.Errorf("KeyEnter constant should be 'enter'")
	}
	if constants.KeyEsc != "esc" {
		t.Errorf("KeyEsc constant should be 'esc'")
	}
}
