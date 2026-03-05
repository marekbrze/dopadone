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
		showCheckbox  bool
		expectedTitle string
	}{
		{
			name:          "subarea with parent",
			parentName:    "Work",
			entityType:    EntityTypeSubarea,
			parentID:      "area-123",
			showCheckbox:  false,
			expectedTitle: "New Subarea in: Work",
		},
		{
			name:          "project with parent and checkbox",
			parentName:    "Tasks",
			entityType:    EntityTypeProject,
			parentID:      "project-123",
			showCheckbox:  true,
			expectedTitle: "New Project in: Tasks",
		},
		{
			name:          "task with parent",
			parentName:    "Build Feature",
			entityType:    EntityTypeTask,
			parentID:      "project-456",
			showCheckbox:  false,
			expectedTitle: "New Task in: Build Feature",
		},
		{
			name:          "without parent name",
			parentName:    "",
			entityType:    EntityTypeSubarea,
			parentID:      "area-789",
			showCheckbox:  false,
			expectedTitle: "New Subarea",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.parentName, tt.entityType, tt.parentID, nil, tt.showCheckbox)
			if !strings.Contains(m.title, tt.expectedTitle) {
				t.Errorf("expected title to contain %q, got %q", tt.expectedTitle, m.title)
			}
			if m.entityType != tt.entityType {
				t.Errorf("expected entity type %v, got %v", tt.entityType, m.entityType)
			}
			if m.parentID != tt.parentID {
				t.Errorf("expected parent ID %q, got %q", tt.parentID, m.parentID)
			}
			if m.showCheckbox != tt.showCheckbox {
				t.Errorf("expected showCheckbox %v, got %v", tt.showCheckbox, m.showCheckbox)
			}
			if tt.showCheckbox && m.checkboxChecked {
				t.Error("expected checkbox to be unchecked initially")
			}
			if m.focusedElement != focusedInput {
				t.Errorf("expected initial focus on input, got %v", m.focusedElement)
			}
		})
	}
}

func TestModalView(t *testing.T) {
	m := New("Test Parent", EntityTypeProject, "", nil, false)
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
	m := New("Test", EntityTypeSubarea, "area-1", nil, false)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'T', 'e', 's', 't'}}
	updated, _ := m.Update(keyMsg)

	if updated.input.Value() != "Test" {
		t.Errorf("expected input value 'Test', got %q", updated.input.Value())
	}
}

func TestModalSubmit(t *testing.T) {
	m := New("Test", EntityTypeSubarea, "area-1", nil, false)
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
	m := New("Test", EntityTypeSubarea, "area-1", nil, false)
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
			m := New("Test", EntityTypeSubarea, "area-1", nil, false)
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
	m := New("Test", EntityTypeSubarea, "area-1", nil, false)

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
	m := New("Test", EntityTypeSubarea, "area-1", nil, false)
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
			m := New(tt.parentName, EntityTypeProject, "area-1", nil, false)
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

func TestCheckboxVisibility(t *testing.T) {
	tests := []struct {
		name         string
		showCheckbox bool
		shouldShow   bool
	}{
		{
			name:         "checkbox shown when showCheckbox is true",
			showCheckbox: true,
			shouldShow:   true,
		},
		{
			name:         "checkbox hidden when showCheckbox is false",
			showCheckbox: false,
			shouldShow:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test Project", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.width = 80
			m.height = 24

			view := m.View()

			if tt.shouldShow {
				if !strings.Contains(view, "Add as subproject") {
					t.Error("expected view to contain checkbox text")
				}
				if !strings.Contains(view, "Tab: Switch") {
					t.Error("expected hint text to mention Tab navigation")
				}
			} else {
				if strings.Contains(view, "Add as subproject") {
					t.Error("expected checkbox text to be hidden")
				}
				if strings.Contains(view, "Tab: Switch") {
					t.Error("expected hint text not to mention Tab navigation")
				}
			}
		})
	}
}

func TestCheckboxInitialState(t *testing.T) {
	m := New("Test Project", EntityTypeProject, "proj-1", nil, true)

	if m.checkboxChecked {
		t.Error("expected checkbox to be unchecked initially")
	}
	if m.focusedElement != focusedInput {
		t.Error("expected initial focus on input field")
	}
}

func TestTabNavigation(t *testing.T) {
	tests := []struct {
		name              string
		showCheckbox      bool
		initialFocus      focusedElement
		expectedFocus     focusedElement
		shouldToggleFocus bool
	}{
		{
			name:              "Tab from input to checkbox when checkbox visible",
			showCheckbox:      true,
			initialFocus:      focusedInput,
			expectedFocus:     focusedCheckbox,
			shouldToggleFocus: true,
		},
		{
			name:              "Tab from checkbox back to input",
			showCheckbox:      true,
			initialFocus:      focusedCheckbox,
			expectedFocus:     focusedInput,
			shouldToggleFocus: true,
		},
		{
			name:              "Tab does nothing when checkbox hidden",
			showCheckbox:      false,
			initialFocus:      focusedInput,
			expectedFocus:     focusedInput,
			shouldToggleFocus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.focusedElement = tt.initialFocus

			tabMsg := tea.KeyMsg{Type: tea.KeyTab}
			updated, _ := m.Update(tabMsg)

			if updated.focusedElement != tt.expectedFocus {
				t.Errorf("expected focus %v, got %v", tt.expectedFocus, updated.focusedElement)
			}
		})
	}
}

func TestShiftTabNavigation(t *testing.T) {
	tests := []struct {
		name          string
		showCheckbox  bool
		initialFocus  focusedElement
		expectedFocus focusedElement
	}{
		{
			name:          "Shift+Tab from input to checkbox",
			showCheckbox:  true,
			initialFocus:  focusedInput,
			expectedFocus: focusedCheckbox,
		},
		{
			name:          "Shift+Tab from checkbox back to input",
			showCheckbox:  true,
			initialFocus:  focusedCheckbox,
			expectedFocus: focusedInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.focusedElement = tt.initialFocus

			shiftTabMsg := tea.KeyMsg{Type: tea.KeyShiftTab}
			updated, _ := m.Update(shiftTabMsg)

			if updated.focusedElement != tt.expectedFocus {
				t.Errorf("expected focus %v, got %v", tt.expectedFocus, updated.focusedElement)
			}
		})
	}
}

func TestCheckboxToggle(t *testing.T) {
	tests := []struct {
		name            string
		showCheckbox    bool
		focusedElement  focusedElement
		initialChecked  bool
		expectedChecked bool
		shouldToggle    bool
	}{
		{
			name:            "Space toggles checkbox on when focused",
			showCheckbox:    true,
			focusedElement:  focusedCheckbox,
			initialChecked:  false,
			expectedChecked: true,
			shouldToggle:    true,
		},
		{
			name:            "Space toggles checkbox off when focused",
			showCheckbox:    true,
			focusedElement:  focusedCheckbox,
			initialChecked:  true,
			expectedChecked: false,
			shouldToggle:    true,
		},
		{
			name:            "Space does not toggle when input focused",
			showCheckbox:    true,
			focusedElement:  focusedInput,
			initialChecked:  false,
			expectedChecked: false,
			shouldToggle:    false,
		},
		{
			name:            "Space does nothing when checkbox hidden",
			showCheckbox:    false,
			focusedElement:  focusedInput,
			initialChecked:  false,
			expectedChecked: false,
			shouldToggle:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.focusedElement = tt.focusedElement
			m.checkboxChecked = tt.initialChecked

			spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
			updated, _ := m.Update(spaceMsg)

			if updated.checkboxChecked != tt.expectedChecked {
				t.Errorf("expected checkbox checked %v, got %v", tt.expectedChecked, updated.checkboxChecked)
			}
		})
	}
}

func TestCheckboxSubmitMsg(t *testing.T) {
	tests := []struct {
		name                 string
		showCheckbox         bool
		checkboxChecked      bool
		expectedAsSubproject bool
	}{
		{
			name:                 "submit with checkbox checked",
			showCheckbox:         true,
			checkboxChecked:      true,
			expectedAsSubproject: true,
		},
		{
			name:                 "submit with checkbox unchecked",
			showCheckbox:         true,
			checkboxChecked:      false,
			expectedAsSubproject: false,
		},
		{
			name:                 "submit without checkbox",
			showCheckbox:         false,
			checkboxChecked:      false,
			expectedAsSubproject: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.checkboxChecked = tt.checkboxChecked
			m.input.SetValue("New Project")

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

			if submitMsg.AsSubproject != tt.expectedAsSubproject {
				t.Errorf("expected AsSubproject %v, got %v", tt.expectedAsSubproject, submitMsg.AsSubproject)
			}
		})
	}
}

func TestCheckboxRendering(t *testing.T) {
	tests := []struct {
		name            string
		showCheckbox    bool
		checkboxChecked bool
		focusedElement  focusedElement
		expectedText    string
	}{
		{
			name:            "unchecked checkbox",
			showCheckbox:    true,
			checkboxChecked: false,
			focusedElement:  focusedInput,
			expectedText:    "[ ]",
		},
		{
			name:            "checked checkbox",
			showCheckbox:    true,
			checkboxChecked: true,
			focusedElement:  focusedInput,
			expectedText:    "[✓]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New("Test Project", EntityTypeProject, "proj-1", nil, tt.showCheckbox)
			m.checkboxChecked = tt.checkboxChecked
			m.focusedElement = tt.focusedElement
			m.width = 80
			m.height = 24

			view := m.View()

			if !strings.Contains(view, tt.expectedText) {
				t.Errorf("expected view to contain %q", tt.expectedText)
			}
		})
	}
}

func TestInputBlurOnCheckboxFocus(t *testing.T) {
	m := New("Test", EntityTypeProject, "proj-1", nil, true)

	if !m.input.Focused() {
		t.Error("expected input to be focused initially")
	}

	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	updated, _ := m.Update(tabMsg)

	if updated.focusedElement != focusedCheckbox {
		t.Error("expected focus to be on checkbox")
	}
	if updated.input.Focused() {
		t.Error("expected input to be blurred when checkbox is focused")
	}

	tabMsg2 := tea.KeyMsg{Type: tea.KeyTab}
	updated2, _ := updated.Update(tabMsg2)

	if updated2.focusedElement != focusedInput {
		t.Error("expected focus to be back on input")
	}
	if !updated2.input.Focused() {
		t.Error("expected input to be focused again")
	}
}

func TestMultipleToggleCycles(t *testing.T) {
	m := New("Test", EntityTypeProject, "proj-1", nil, true)
	m.focusedElement = focusedCheckbox

	for i := 0; i < 5; i++ {
		spaceMsg := tea.KeyMsg{Type: tea.KeySpace}
		updated, _ := m.Update(spaceMsg)

		expectedChecked := (i%2 == 0)
		if updated.checkboxChecked != expectedChecked {
			t.Errorf("cycle %d: expected checked %v, got %v", i, expectedChecked, updated.checkboxChecked)
		}

		m = updated
	}
}
