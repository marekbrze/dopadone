package domain

import (
	"testing"
	"time"
)

func TestColorValidation(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		isValid bool
	}{
		{
			name:    "valid 6-digit hex color uppercase",
			color:   "#FF0000",
			isValid: true,
		},
		{
			name:    "valid 6-digit hex color lowercase",
			color:   "#ff0000",
			isValid: true,
		},
		{
			name:    "valid 6-digit hex color mixed case",
			color:   "#Ff00Aa",
			isValid: true,
		},
		{
			name:    "valid color with all digits",
			color:   "#123456",
			isValid: true,
		},
		{
			name:    "empty string is valid (optional color)",
			color:   "",
			isValid: true,
		},
		{
			name:    "invalid - missing hash prefix",
			color:   "FF0000",
			isValid: false,
		},
		{
			name:    "invalid - 3-digit hex (not supported)",
			color:   "#F00",
			isValid: false,
		},
		{
			name:    "invalid - 8-digit hex with alpha",
			color:   "#FF0000FF",
			isValid: false,
		},
		{
			name:    "invalid - 5 digits",
			color:   "#FF000",
			isValid: false,
		},
		{
			name:    "invalid - 7 digits",
			color:   "#FF00000",
			isValid: false,
		},
		{
			name:    "invalid - invalid characters",
			color:   "#GG0000",
			isValid: false,
		},
		{
			name:    "invalid - special characters",
			color:   "#FF-000",
			isValid: false,
		},
		{
			name:    "invalid - spaces",
			color:   "#FF 000",
			isValid: false,
		},
		{
			name:    "invalid - lowercase g",
			color:   "#gg0000",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.color.IsValid(); got != tt.isValid {
				t.Errorf("Color(%q).IsValid() = %v, want %v", tt.color, got, tt.isValid)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Color
		wantErr bool
	}{
		{
			name:    "valid hex color uppercase",
			input:   "#FF0000",
			want:    Color("#FF0000"),
			wantErr: false,
		},
		{
			name:    "valid hex color lowercase",
			input:   "#00ff00",
			want:    Color("#00ff00"),
			wantErr: false,
		},
		{
			name:    "empty string returns empty color",
			input:   "",
			want:    Color(""),
			wantErr: false,
		},
		{
			name:    "invalid - missing hash",
			input:   "FF0000",
			want:    Color(""),
			wantErr: true,
		},
		{
			name:    "invalid - wrong length",
			input:   "#FFF",
			want:    Color(""),
			wantErr: true,
		},
		{
			name:    "invalid - invalid characters",
			input:   "#GGHHII",
			want:    Color(""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestColorString(t *testing.T) {
	color := Color("#FF0000")
	if got := color.String(); got != "#FF0000" {
		t.Errorf("Color.String() = %q, want %q", got, "#FF0000")
	}
}

func TestProjectStatusValidation(t *testing.T) {
	tests := []struct {
		name    string
		status  ProjectStatus
		isValid bool
	}{
		{"active", ProjectStatusActive, true},
		{"completed", ProjectStatusCompleted, true},
		{"on_hold", ProjectStatusOnHold, true},
		{"archived", ProjectStatusArchived, true},
		{"invalid", ProjectStatus("invalid"), false},
		{"empty", ProjectStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.isValid {
				t.Errorf("ProjectStatus(%q).IsValid() = %v, want %v", tt.status, got, tt.isValid)
			}
		})
	}
}

func TestParseProjectStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ProjectStatus
		wantErr bool
	}{
		{"active", "active", ProjectStatusActive, false},
		{"completed", "completed", ProjectStatusCompleted, false},
		{"on_hold", "on_hold", ProjectStatusOnHold, false},
		{"archived", "archived", ProjectStatusArchived, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProjectStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseProjectStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestProjectStatusString(t *testing.T) {
	status := ProjectStatusActive
	if got := status.String(); got != "active" {
		t.Errorf("ProjectStatus.String() = %q, want %q", got, "active")
	}
}

func TestPriorityValidation(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		isValid  bool
	}{
		{"low", PriorityLow, true},
		{"medium", PriorityMedium, true},
		{"high", PriorityHigh, true},
		{"urgent", PriorityUrgent, true},
		{"invalid", Priority("invalid"), false},
		{"empty", Priority(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.priority.IsValid(); got != tt.isValid {
				t.Errorf("Priority(%q).IsValid() = %v, want %v", tt.priority, got, tt.isValid)
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Priority
		wantErr bool
	}{
		{"low", "low", PriorityLow, false},
		{"medium", "medium", PriorityMedium, false},
		{"high", "high", PriorityHigh, false},
		{"urgent", "urgent", PriorityUrgent, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePriority(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePriority(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestPriorityString(t *testing.T) {
	priority := PriorityHigh
	if got := priority.String(); got != "high" {
		t.Errorf("Priority.String() = %q, want %q", got, "high")
	}
}

func TestProgressValidation(t *testing.T) {
	tests := []struct {
		name    string
		value   Progress
		isValid bool
	}{
		{"zero", Progress(0), true},
		{"50", Progress(50), true},
		{"100", Progress(100), true},
		{"negative", Progress(-1), false},
		{"101", Progress(101), false},
		{"very negative", Progress(-100), false},
		{"very high", Progress(200), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.value.IsValid(); got != tt.isValid {
				t.Errorf("Progress(%d).IsValid() = %v, want %v", tt.value, got, tt.isValid)
			}
		})
	}
}

func TestParseProgress(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    Progress
		wantErr bool
	}{
		{"zero", 0, Progress(0), false},
		{"50", 50, Progress(50), false},
		{"100", 100, Progress(100), false},
		{"negative", -1, Progress(0), true},
		{"101", 101, Progress(0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProgress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProgress(%d) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseProgress(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestProgressInt(t *testing.T) {
	progress := Progress(75)
	if got := progress.Int(); got != 75 {
		t.Errorf("Progress.Int() = %d, want %d", got, 75)
	}
}

func TestDateRangeValidation(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)

	tests := []struct {
		name      string
		startDate *time.Time
		deadline  *time.Time
		isValid   bool
	}{
		{"both nil", nil, nil, true},
		{"only start date", &now, nil, true},
		{"start before deadline", &now, &future, true},
		{"start equal deadline", &now, &now, false},
		{"start after deadline", &future, &now, false},
		{"deadline without start", nil, &future, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dr := DateRange{StartDate: tt.startDate, Deadline: tt.deadline}
			if got := dr.IsValid(); got != tt.isValid {
				t.Errorf("DateRange.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestNewDateRange(t *testing.T) {
	now := time.Now()
	future := now.Add(24 * time.Hour)

	t.Run("valid date range", func(t *testing.T) {
		dr, err := NewDateRange(&now, &future)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if dr.StartDate == nil || !dr.StartDate.Equal(now) {
			t.Errorf("StartDate mismatch")
		}
		if dr.Deadline == nil || !dr.Deadline.Equal(future) {
			t.Errorf("Deadline mismatch")
		}
	})

	t.Run("invalid date range - deadline before start", func(t *testing.T) {
		_, err := NewDateRange(&future, &now)
		if err != ErrInvalidDateRange {
			t.Errorf("expected ErrInvalidDateRange, got %v", err)
		}
	})

	t.Run("invalid date range - deadline without start", func(t *testing.T) {
		_, err := NewDateRange(nil, &future)
		if err != ErrInvalidDateRange {
			t.Errorf("expected ErrInvalidDateRange, got %v", err)
		}
	})
}
