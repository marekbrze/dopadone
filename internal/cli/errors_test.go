package cli

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		message  string
		expected string
	}{
		{
			name:     "with field",
			field:    "name",
			message:  "cannot be empty",
			expected: "name: cannot be empty",
		},
		{
			name:     "without field",
			field:    "",
			message:  "invalid input",
			expected: "invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.message)
			if err.Error() != tt.expected {
				t.Errorf("Error() = %q, want %q", err.Error(), tt.expected)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	t.Run("wrap error", func(t *testing.T) {
		original := errors.New("original error")
		wrapped := WrapError(original, "context")

		if wrapped == nil {
			t.Error("WrapError() returned nil")
		}

		if wrapped.Error() != "context: original error" {
			t.Errorf("WrapError() = %q, want %q", wrapped.Error(), "context: original error")
		}
	})

	t.Run("nil error", func(t *testing.T) {
		wrapped := WrapError(nil, "context")
		if wrapped != nil {
			t.Errorf("WrapError(nil) = %v, want nil", wrapped)
		}
	})
}

func TestIsValidationError(t *testing.T) {
	t.Run("is validation error", func(t *testing.T) {
		err := NewValidationError("field", "message")
		if !IsValidationError(err) {
			t.Error("IsValidationError() = false, want true")
		}
	})

	t.Run("is not validation error", func(t *testing.T) {
		err := errors.New("regular error")
		if IsValidationError(err) {
			t.Error("IsValidationError() = true, want false")
		}
	})
}
