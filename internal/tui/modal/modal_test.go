package modal

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewModal(t *testing.T) {
	tests := []struct {
		name          string
		parentName    string
		entityType    EntityType
		parentID      string
		expectedTitle string
	}{
		{
			name:          "subarea with parent",
			parentName:    "Work",
			entityType:    EntityTypeSubarea,
			parentID:      "area-123",
			expectedTitle: "New Subarea in: Work",
		},
		{
			name:          "project with parent",
			parentName:    "Tasks",
			entityType:    EntityTypeProject,
			parentID:      "",
			expectedTitle: "New Project in: Tasks",
		},
		{
			name:          "task with parent",
			parentName:    "Build Feature",
			entityType:    EntityTypeTask,
			parentID:      "project-456",
			expectedTitle: "New Task in: Build Feature",
		},
		{
			name:          "without parent name",
			parentName:    "",
			entityType:    EntityTypeSubarea,
			parentID:      "area-789",
			expectedTitle: "New Subarea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.parentName, tt.entityType, tt.parentID, nil)
			if !strings.Contains(m.title, tt.expectedTitle) {
				t.Errorf("expected title to contain %q, got %q", tt.expectedTitle, m.title)
			}
			if m.entityType != tt.entityType {
				t.Errorf("expected entity type %v, got %v", tt.entityType, m.entityType)
			}
			if m.parentID != tt.parentID {
				t.Errorf("expected parent ID %q, got %q", tt.parentID, m.parentID)
			}
		})
	}
}

func TestModalView(t *testing.T) {
	m := New("Test Parent", EntityTypeProject, "", nil)
	m.width = 80
	m.height = 24

	view := m.View()

	if !strings.Contains(view, "New Project") {
		t.Error("expected view to contain 'New Project'")
	}
	if !strings.Contains(view, "Test Parent") {
		t.Error("expected view to contain parent name")
	}
	if !strings.Contains(view, "Enter: Create") {
		t.Error("expected view to contain hint text")
	}
}

func TestModalInput(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}}
	updated, _ := m.Update(keyMsg)

	if updated.input.Value() != "Test" {
		t.Errorf("expected input value 'Test', got %q", updated.input.Value())
	}
}

func TestModalSubmit(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil)
	m.input.SetValue("New Item")

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := m.Update(enterMsg)

	if cmd == nil {
		t.Fatal("expected command to be returned")
	}

	msg := cmd()
	submitMsg, ok := msg.(SubmitMsg)
	if !ok {
		t.Fatalf("expected SubmitMsg, got %T", msg)
	}

	if submitMsg.Title != "New Item" {
		t.Errorf("expected title 'New Item', got %q", submitMsg.Title)
	}
	if submitMsg.EntityType != EntityTypeSubarea {
		t.Errorf("expected entity type Subarea, got %v", submitMsg.EntityType)
	}
	if submitMsg.ParentID != "area-1" {
		t.Errorf("expected parent ID 'area-1', got %q", submitMsg.ParentID)
	}
}

func TestModalCancel(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil)
	m.input.SetValue("Some text")

	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	_, cmd := m.Update(escMsg)

	if cmd == nil {
		t.Fatal("expected command to be returned")
	}

	msg := cmd()
	_, ok := msg.(CloseMsg)
	if !ok {
		t.Fatalf("expected CloseMsg, got %T", msg)
	}
}

func TestModalValidation(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		shouldError   bool
		expectedError string
	}{
		{
			name:        "valid input",
			input:       "Valid Title",
			shouldError: false,
		},
		{
			name:          "empty input",
			input:         "",
			shouldError:   true,
			expectedError: ErrTitleEmpty.Error(),
		},
		{
			name:          "whitespace only",
			input:         "   ",
			shouldError:   true,
			expectedError: ErrTitleEmpty.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", EntityTypeSubarea, "area-1", nil)
			m.input.SetValue(tt.input)

			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			updated, cmd := m.Update(enterMsg)

			if tt.shouldError {
				if cmd != nil {
					t.Error("expected no command on validation error")
				}
				if updated.errorMsg == "" {
					t.Error("expected error message to be set")
				}
				if !strings.Contains(updated.errorMsg, tt.expectedError) {
					t.Errorf("expected error to contain %q, got %q", tt.expectedError, updated.errorMsg)
				}
			} else {
				if cmd == nil {
					t.Error("expected command on valid input")
				}
			}
		})
	}
}

func TestModalSetError(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil)

	m.SetError("Test error message")
	if m.errorMsg != "Test error message" {
		t.Errorf("expected error message 'Test error message', got %q", m.errorMsg)
	}

	view := m.View()
	if !strings.Contains(view, "Test error message") {
		t.Error("expected view to contain error message")
	}
}

func TestModalClearError(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil)
	m.errorMsg = "Previous error"

	m.ClearError()
	if m.errorMsg != "" {
		t.Errorf("expected error message to be cleared, got %q", m.errorMsg)
	}
}

func TestModalInputWidth(t *testing.T) {
	tests := []struct {
		name        string
		parentName  string
		expectedMin int
		expectedMax int
	}{
		{
			name:        "short title",
			parentName:  "",
			expectedMin: 30,
			expectedMax: 40,
		},
		{
			name:        "medium title",
			parentName:  "Test Parent",
			expectedMin: 30,
			expectedMax: 50,
		},
		{
			name:        "long title caps at 60",
			parentName:  "This is a very long parent name that should expand the modal",
			expectedMin: 55,
			expectedMax: 60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.parentName, EntityTypeProject, "area-1", nil)
			m.width = 100
			m.height = 24

			_ = m.View()

			if m.input.Width < tt.expectedMin {
				t.Errorf("expected input width >= %d, got %d", tt.expectedMin, m.input.Width)
			}
			if m.input.Width > tt.expectedMax {
				t.Errorf("expected input width <= %d, got %d", tt.expectedMax, m.input.Width)
			}
		})
	}
}
