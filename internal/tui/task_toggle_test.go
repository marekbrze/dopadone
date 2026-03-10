package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/domain"
)

func TestToggleTaskCompletion(t *testing.T) {
	tests := []struct {
		name           string
		currentStatus  domain.TaskStatus
		expectedStatus domain.TaskStatus
	}{
		{"todo to done", domain.TaskStatusTodo, domain.TaskStatusDone},
		{"in_progress to done", domain.TaskStatusInProgress, domain.TaskStatusDone},
		{"waiting to done", domain.TaskStatusWaiting, domain.TaskStatusDone},
		{"done to todo", domain.TaskStatusDone, domain.TaskStatusTodo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Model{
				tasks: []domain.Task{
					{ID: "task-1", Title: "Test Task", Status: tt.currentStatus},
				},
				selectedTaskIndex: 0,
				focus:             FocusTasks,
			}

			cmd := m.toggleTaskCompletion()
			if cmd == nil {
				t.Error("expected command to be returned")
			}

			if m.tasks[0].Status != tt.expectedStatus {
				t.Errorf("status = %v, want %v", m.tasks[0].Status, tt.expectedStatus)
			}
		})
	}
}

func TestToggleTaskCompletion_KeyboardBinding(t *testing.T) {
	m := Model{
		tasks: []domain.Task{
			{ID: "task-1", Title: "Test Task", Status: domain.TaskStatusTodo},
		},
		selectedTaskIndex: 0,
		focus:             FocusTasks,
		ready:             true,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
	model, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("x key should trigger toggle command when focus is on Tasks")
	}

	m = *model.(*Model)
	if m.tasks[0].Status != domain.TaskStatusDone {
		t.Errorf("task should be toggled to done, got %v", m.tasks[0].Status)
	}

	m.focus = FocusProjects
	_, cmd = m.Update(msg)
	if cmd != nil {
		t.Error("x key should not trigger when focus is not on Tasks")
	}
}

func TestToggleTaskCompletion_EmptyTaskList(t *testing.T) {
	m := Model{
		tasks:             []domain.Task{},
		selectedTaskIndex: 0,
		focus:             FocusTasks,
		ready:             true,
	}

	cmd := m.toggleTaskCompletion()
	if cmd != nil {
		t.Error("expected no command for empty task list")
	}
}

func TestToggleTaskCompletion_InvalidIndex(t *testing.T) {
	m := Model{
		tasks: []domain.Task{
			{ID: "task-1", Title: "Test Task", Status: domain.TaskStatusTodo},
		},
		selectedTaskIndex: 5,
		focus:             FocusTasks,
		ready:             true,
	}

	cmd := m.toggleTaskCompletion()
	if cmd != nil {
		t.Error("expected no command for invalid index")
	}
}
