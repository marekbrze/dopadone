package cli

import (
	"testing"
)

func TestParseProjectStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"active", "active", false},
		{"completed", "completed", false},
		{"on_hold", "on_hold", false},
		{"archived", "archived", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := ParseProjectStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !IsValidationError(err) && err != nil {
				t.Error("ParseProjectStatus() should return ValidationError")
			}
			if !tt.wantErr && string(status) != tt.input {
				t.Errorf("ParseProjectStatus() = %v, want %v", status, tt.input)
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"low", "low", false},
		{"medium", "medium", false},
		{"high", "high", false},
		{"urgent", "urgent", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority, err := ParsePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(priority) != tt.input {
				t.Errorf("ParsePriority() = %v, want %v", priority, tt.input)
			}
		})
	}
}

func TestParseProgress(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		wantErr bool
	}{
		{"zero", 0, false},
		{"fifty", 50, false},
		{"hundred", 100, false},
		{"negative", -1, true},
		{"over hundred", 101, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			progress, err := ParseProgress(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && int(progress) != tt.input {
				t.Errorf("ParseProgress() = %v, want %v", progress, tt.input)
			}
		})
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid red", "#FF0000", false},
		{"valid green", "#00FF00", false},
		{"valid blue", "#0000FF", false},
		{"valid lowercase", "#ff0000", false},
		{"empty", "", false},
		{"invalid no hash", "FF0000", true},
		{"invalid short", "#FFF", true},
		{"invalid chars", "#ZZZZZZ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color, err := ParseColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseColor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && string(color) != tt.input {
				t.Errorf("ParseColor() = %v, want %v", color, tt.input)
			}
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name        string
		startDate   string
		deadline    string
		wantErr     bool
		errContains string
	}{
		{"both empty", "", "", false, ""},
		{"valid start only", "2024-01-01", "", false, ""},
		{"valid both", "2024-01-01", "2024-12-31", false, ""},
		{"deadline without start", "", "2024-12-31", true, "deadline"},
		{"start after deadline", "2024-12-31", "2024-01-01", true, "date_range"},
		{"invalid start format", "01-01-2024", "", true, "start_date"},
		{"invalid deadline format", "", "31-12-2024", true, "deadline"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, deadline, err := ParseDate(tt.startDate, tt.deadline)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.startDate == "" && start != nil {
					t.Error("ParseDate() start should be nil")
				}
				if tt.deadline == "" && deadline != nil {
					t.Error("ParseDate() deadline should be nil")
				}
			}
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		err := ValidateProjectName("My Project")
		if err != nil {
			t.Errorf("ValidateProjectName() error = %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		err := ValidateProjectName("")
		if err == nil {
			t.Error("ValidateProjectName() should return error for empty name")
		}
		if !IsValidationError(err) {
			t.Error("ValidateProjectName() should return ValidationError")
		}
	})
}
