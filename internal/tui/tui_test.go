package tui

import "testing"

func TestFocusColumnString(t *testing.T) {
	tests := []struct {
		column   FocusColumn
		expected string
	}{
		{FocusSubareas, "Subareas"},
		{FocusProjects, "Projects"},
		{FocusTasks, "Tasks"},
		{FocusColumn(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.column.String(); got != tt.expected {
				t.Errorf("FocusColumn.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFocusColumnPrev(t *testing.T) {
	tests := []struct {
		name     string
		column   FocusColumn
		expected FocusColumn
	}{
		{"Subareas wraps to Tasks", FocusSubareas, FocusTasks},
		{"Projects to Subareas", FocusProjects, FocusSubareas},
		{"Tasks to Projects", FocusTasks, FocusProjects},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.column.Prev(); got != tt.expected {
				t.Errorf("FocusColumn.Prev() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFocusColumnNext(t *testing.T) {
	tests := []struct {
		name     string
		column   FocusColumn
		expected FocusColumn
	}{
		{"Subareas to Projects", FocusSubareas, FocusProjects},
		{"Projects to Tasks", FocusProjects, FocusTasks},
		{"Tasks wraps to Subareas", FocusTasks, FocusSubareas},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.column.Next(); got != tt.expected {
				t.Errorf("FocusColumn.Next() = %v, want %v", got, tt.expected)
			}
		})
	}
}
