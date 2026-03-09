package confirmmodal

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		itemName   string
		entityType EntityType
		entityID   string
	}{
		{
			name:       "subarea",
			itemName:   "Test Subarea",
			entityType: EntityTypeSubarea,
			entityID:   "subarea-123",
		},
		{
			name:       "project",
			itemName:   "Test Project",
			entityType: EntityTypeProject,
			entityID:   "project-456",
		},
		{
			name:       "task",
			itemName:   "Test Task",
			entityType: EntityTypeTask,
			entityID:   "task-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.itemName, tt.entityType, tt.entityID)
			if m == nil {
				t.Fatal("New() returned nil")
			}
			if m.itemName != tt.itemName {
				t.Errorf("expected itemName %q, got %q", tt.itemName, m.itemName)
			}
			if m.entityType != tt.entityType {
				t.Errorf("expected entityType %v, got %v", tt.entityType, m.entityType)
			}
			if m.entityID != tt.entityID {
				t.Errorf("expected entityID %q, got %q", tt.entityID, m.entityID)
			}
		})
	}
}

func TestModalInit(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}

func TestModalUpdateWindowSize(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}

	updated, cmd := m.Update(msg)
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

func TestModalConfirmKey(t *testing.T) {
	tests := []struct {
		name       string
		entityType EntityType
		entityID   string
		itemName   string
	}{
		{
			name:       "confirm subarea",
			entityType: EntityTypeSubarea,
			entityID:   "subarea-1",
			itemName:   "My Subarea",
		},
		{
			name:       "confirm project",
			entityType: EntityTypeProject,
			entityID:   "project-1",
			itemName:   "My Project",
		},
		{
			name:       "confirm task",
			entityType: EntityTypeTask,
			entityID:   "task-1",
			itemName:   "My Task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.itemName, tt.entityType, tt.entityID)
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}

			_, cmd := m.Update(msg)
			if cmd == nil {
				t.Fatal("expected command to be returned")
			}

			result := cmd()
			confirmMsg, ok := result.(ConfirmMsg)
			if !ok {
				t.Fatalf("expected ConfirmMsg, got %T", result)
			}
			if confirmMsg.EntityType != tt.entityType {
				t.Errorf("expected entityType %v, got %v", tt.entityType, confirmMsg.EntityType)
			}
			if confirmMsg.EntityID != tt.entityID {
				t.Errorf("expected entityID %q, got %q", tt.entityID, confirmMsg.EntityID)
			}
			if confirmMsg.EntityName != tt.itemName {
				t.Errorf("expected entityName %q, got %q", tt.itemName, confirmMsg.EntityName)
			}
		})
	}
}

func TestModalCancelKeys(t *testing.T) {
	keys := []string{"n", "esc"}

	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			m := New("Test Item", EntityTypeTask, "task-1")

			var msg tea.Msg
			if key == "esc" {
				msg = tea.KeyMsg{Type: tea.KeyEsc}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
			}

			_, cmd := m.Update(msg)
			if cmd == nil {
				t.Fatalf("expected command for key %s", key)
			}

			result := cmd()
			if _, ok := result.(CancelMsg); !ok {
				t.Errorf("expected CancelMsg, got %T", result)
			}
		})
	}
}

func TestModalView(t *testing.T) {
	tests := []struct {
		name       string
		itemName   string
		entityType EntityType
		expected   []string
	}{
		{
			name:       "subarea deletion",
			itemName:   "Test Subarea",
			entityType: EntityTypeSubarea,
			expected:   []string{"Delete Subarea?", "Test Subarea", "y: confirm", "n/esc: cancel"},
		},
		{
			name:       "project deletion",
			itemName:   "Test Project",
			entityType: EntityTypeProject,
			expected:   []string{"Delete Project?", "Test Project", "y: confirm", "n/esc: cancel"},
		},
		{
			name:       "task deletion",
			itemName:   "Test Task",
			entityType: EntityTypeTask,
			expected:   []string{"Delete Task?", "Test Task", "y: confirm", "n/esc: cancel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.itemName, tt.entityType, "id-1")
			m.width = 80
			m.height = 24

			view := m.View()

			for _, exp := range tt.expected {
				if !strings.Contains(view, exp) {
					t.Errorf("expected view to contain %q", exp)
				}
			}
		})
	}
}

func TestModalViewWithoutSize(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")

	view := m.View()
	if view == "" {
		t.Error("expected non-empty view even without size")
	}
}

func TestModalTruncateItemName(t *testing.T) {
	tests := []struct {
		name     string
		itemName string
		expected string
	}{
		{
			name:     "short name",
			itemName: "Short",
			expected: "Short",
		},
		{
			name:     "exactly max length",
			itemName: strings.Repeat("a", 40),
			expected: strings.Repeat("a", 40),
		},
		{
			name:     "long name",
			itemName: strings.Repeat("a", 50),
			expected: strings.Repeat("a", 37) + "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.itemName, EntityTypeTask, "task-1")
			got := m.truncateItemName()
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestModalUnknownKey(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	updated, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("expected nil cmd for unknown key")
	}
	if updated.itemName != "Test" {
		t.Error("expected modal state to be unchanged")
	}
}

func TestModalMultipleKeyPresses(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")

	unknownMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	updated, cmd := m.Update(unknownMsg)
	if cmd != nil {
		t.Error("expected nil cmd for unknown key")
	}

	confirmMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	_, cmd = updated.Update(confirmMsg)
	if cmd == nil {
		t.Fatal("expected command for confirm key")
	}

	result := cmd()
	if _, ok := result.(ConfirmMsg); !ok {
		t.Errorf("expected ConfirmMsg, got %T", result)
	}
}

func TestModalSpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		itemName string
	}{
		{
			name:     "with quotes",
			itemName: `Test "quoted" item`,
		},
		{
			name:     "with unicode",
			itemName: "Test 你好 🎉",
		},
		{
			name:     "with newlines",
			itemName: "Test\nItem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.itemName, EntityTypeTask, "task-1")
			m.width = 80
			m.height = 24

			view := m.View()
			if view == "" {
				t.Error("expected non-empty view")
			}
		})
	}
}

func TestModalEmptyItemName(t *testing.T) {
	m := New("", EntityTypeTask, "task-1")
	m.width = 80
	m.height = 24

	view := m.View()
	if !strings.Contains(view, "Delete Task?") {
		t.Error("expected view to contain entity type even with empty name")
	}
}

func TestConfirmMsgFields(t *testing.T) {
	m := New("My Item", EntityTypeProject, "proj-123")
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}

	_, cmd := m.Update(msg)
	result := cmd()
	confirmMsg, ok := result.(ConfirmMsg)
	if !ok {
		t.Fatalf("expected ConfirmMsg, got %T", result)
	}

	if confirmMsg.EntityName != "My Item" {
		t.Errorf("expected EntityName 'My Item', got %q", confirmMsg.EntityName)
	}
	if confirmMsg.EntityType != EntityTypeProject {
		t.Errorf("expected EntityType Project, got %v", confirmMsg.EntityType)
	}
	if confirmMsg.EntityID != "proj-123" {
		t.Errorf("expected EntityID 'proj-123', got %q", confirmMsg.EntityID)
	}
}

func TestModalCentering(t *testing.T) {
	m := New("Test", EntityTypeTask, "task-1")
	m.width = 100
	m.height = 50

	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}
