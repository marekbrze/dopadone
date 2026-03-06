package views

import (
	"strings"
	"testing"
)

func TestCalculateColumnWidths(t *testing.T) {
	tests := []struct {
		name         string
		totalWidth   int
		wantSubareas int
		wantProjects int
		wantTasks    int
	}{
		{
			name:         "standard 120 cols",
			totalWidth:   120,
			wantSubareas: 28,
			wantProjects: 28,
			wantTasks:    58,
		},
		{
			name:         "wide 160 cols",
			totalWidth:   160,
			wantSubareas: 38,
			wantProjects: 38,
			wantTasks:    78,
		},
		{
			name:         "exact minimum 80 cols",
			totalWidth:   80,
			wantSubareas: 20,
			wantProjects: 20,
			wantTasks:    40,
		},
		{
			name:         "below minimum 70 cols",
			totalWidth:   70,
			wantSubareas: 20,
			wantProjects: 20,
			wantTasks:    40,
		},
		{
			name:         "narrow 90 cols",
			totalWidth:   90,
			wantSubareas: 21,
			wantProjects: 21,
			wantTasks:    42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubareas, gotProjects, gotTasks := calculateColumnWidths(tt.totalWidth)

			if gotSubareas != tt.wantSubareas {
				t.Errorf("subareas = %d, want %d", gotSubareas, tt.wantSubareas)
			}
			if gotProjects != tt.wantProjects {
				t.Errorf("projects = %d, want %d", gotProjects, tt.wantProjects)
			}
			if gotTasks != tt.wantTasks {
				t.Errorf("tasks = %d, want %d", gotTasks, tt.wantTasks)
			}
		})
	}
}

func TestColumnViewTruncation(t *testing.T) {
	tests := []struct {
		name                  string
		col                   Column
		shouldContainEllipsis bool
	}{
		{
			name: "long title truncated",
			col: Column{
				Title:     "This is a very long title that should be truncated",
				Content:   "Short",
				Width:     30,
				Height:    10,
				IsFocused: false,
			},
			shouldContainEllipsis: true,
		},
		{
			name: "long content lines truncated",
			col: Column{
				Title:     "Title",
				Content:   "Line 1 is very long and should be truncated\nLine 2 is also very long",
				Width:     30,
				Height:    10,
				IsFocused: false,
			},
			shouldContainEllipsis: true,
		},
		{
			name: "short text not truncated",
			col: Column{
				Title:     "Short",
				Content:   "Short content",
				Width:     30,
				Height:    10,
				IsFocused: false,
			},
			shouldContainEllipsis: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ColumnView(tt.col)

			if result == "" {
				t.Error("ColumnView returned empty string")
			}

			if tt.shouldContainEllipsis && !strings.Contains(result, "…") {
				t.Error("Expected output to contain ellipsis for truncated text")
			}

			if !tt.shouldContainEllipsis && strings.Contains(result, "…") {
				t.Error("Did not expect ellipsis for short text")
			}
		})
	}
}

func TestShouldUseStackedLayout(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		expected bool
	}{
		{"narrow 119 cols", 119, true},
		{"exactly 120 cols", 120, false},
		{"wide 121 cols", 121, false},
		{"very narrow 80 cols", 80, true},
		{"medium 100 cols", 100, true},
		{"wide 160 cols", 160, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUseStackedLayout(tt.width)
			if got != tt.expected {
				t.Errorf("shouldUseStackedLayout(%d) = %v, want %v", tt.width, got, tt.expected)
			}
		})
	}
}

func TestCalculateStackedLayoutWidths(t *testing.T) {
	tests := []struct {
		name           string
		totalWidth     int
		wantLeftWidth  int
		wantTasksWidth int
	}{
		{
			name:           "narrow 80 cols",
			totalWidth:     80,
			wantLeftWidth:  19,
			wantTasksWidth: 59,
		},
		{
			name:           "medium 100 cols",
			totalWidth:     100,
			wantLeftWidth:  24,
			wantTasksWidth: 74,
		},
		{
			name:           "exactly 119 cols",
			totalWidth:     119,
			wantLeftWidth:  29,
			wantTasksWidth: 88,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeftWidth, gotTasksWidth := calculateStackedLayoutWidths(tt.totalWidth)

			if gotLeftWidth != tt.wantLeftWidth {
				t.Errorf("leftWidth = %d, want %d", gotLeftWidth, tt.wantLeftWidth)
			}
			if gotTasksWidth != tt.wantTasksWidth {
				t.Errorf("tasksWidth = %d, want %d", gotTasksWidth, tt.wantTasksWidth)
			}
		})
	}
}

func TestCalculateStackedLayoutHeights(t *testing.T) {
	tests := []struct {
		name         string
		totalHeight  int
		wantSubareas int
		wantProjects int
	}{
		{
			name:         "standard 30 lines",
			totalHeight:  30,
			wantSubareas: 14,
			wantProjects: 14,
		},
		{
			name:         "short 20 lines",
			totalHeight:  20,
			wantSubareas: 9,
			wantProjects: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSubareas, gotProjects := calculateStackedLayoutHeights(tt.totalHeight)

			if gotSubareas != tt.wantSubareas {
				t.Errorf("subareasHeight = %d, want %d", gotSubareas, tt.wantSubareas)
			}
			if gotProjects != tt.wantProjects {
				t.Errorf("projectsHeight = %d, want %d", gotProjects, tt.wantProjects)
			}
		})
	}
}

func TestLayoutStacked(t *testing.T) {
	columns := []Column{
		{Title: "Subareas", Content: "Item 1\nItem 2", IsFocused: false},
		{Title: "Projects", Content: "Project A\nProject B", IsFocused: true},
		{Title: "Tasks", Content: "Task 1\nTask 2\nTask 3", IsFocused: false},
	}

	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"narrow 80 cols", 80, 30},
		{"medium 100 cols", 100, 30},
		{"exactly 119 cols", 119, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LayoutStacked(columns, tt.width, tt.height)

			if result == "" {
				t.Error("LayoutStacked returned empty string")
			}

			for _, col := range columns {
				if !strings.Contains(result, col.Title) {
					t.Errorf("LayoutStacked missing title: %s", col.Title)
				}
			}
		})
	}
}

func TestLayoutModeSwitching(t *testing.T) {
	columns := []Column{
		{Title: "Subareas", Content: "Item 1", IsFocused: false},
		{Title: "Projects", Content: "Project A", IsFocused: false},
		{Title: "Tasks", Content: "Task 1", IsFocused: true},
	}

	tests := []struct {
		name          string
		width         int
		expectStacked bool
	}{
		{"stacked mode at 119", 119, true},
		{"side-by-side at 120", 120, false},
		{"side-by-side at 121", 121, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Layout(columns, tt.width, 30)

			if result == "" {
				t.Error("Layout returned empty string")
			}

			for _, col := range columns {
				if !strings.Contains(result, col.Title) {
					t.Errorf("Layout missing title at width %d: %s", tt.width, col.Title)
				}
			}
		})
	}
}
