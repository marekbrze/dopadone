package domain

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{"not found", ErrNotFound, "resource not found"},
		{"invalid input", ErrInvalidInput, "invalid input provided"},
		{"database error", ErrDatabaseError, "database operation failed"},
		{"empty id", ErrEmptyID, "id cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.message {
				t.Errorf("got %q, want %q", tt.err.Error(), tt.message)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("title", "cannot be empty")

	expected := "validation failed on title: cannot be empty"
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}

	if !errors.Is(err, ErrInvalidInput) {
		t.Error("ValidationError should wrap ErrInvalidInput")
	}
}

func TestDatabaseError(t *testing.T) {
	originalErr := errors.New("connection failed")
	err := NewDatabaseError("GetTaskByID", originalErr)

	expected := "database error during GetTaskByID: connection failed"
	if err.Error() != expected {
		t.Errorf("got %q, want %q", err.Error(), expected)
	}

	if !errors.Is(err, originalErr) {
		t.Error("DatabaseError should wrap original error")
	}

	var dbErr *DatabaseError
	if !errors.As(err, &dbErr) {
		t.Error("should be able to extract DatabaseError")
	}
	if dbErr.Operation != "GetTaskByID" {
		t.Errorf("got operation %q, want %q", dbErr.Operation, "GetTaskByID")
	}
}

func TestNotFoundError(t *testing.T) {
	t.Run("with ID", func(t *testing.T) {
		err := NewNotFoundError("task", "task-123")

		expected := "task not found: task-123"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}

		if !errors.Is(err, ErrNotFound) {
			t.Error("NotFoundError should wrap ErrNotFound")
		}
	})

	t.Run("without ID", func(t *testing.T) {
		err := NewNotFoundError("project", "")

		expected := "project not found"
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err.Error(), expected)
		}
	})
}

func TestErrorHelpers(t *testing.T) {
	t.Run("IsNotFound", func(t *testing.T) {
		if !IsNotFound(NewNotFoundError("task", "1")) {
			t.Error("IsNotFound should return true for NotFoundError")
		}
		if IsNotFound(errors.New("other error")) {
			t.Error("IsNotFound should return false for other errors")
		}
	})

	t.Run("IsDatabaseError", func(t *testing.T) {
		dbErr := NewDatabaseError("op", errors.New("err"))
		if !IsDatabaseError(dbErr) {
			t.Error("IsDatabaseError should return true for DatabaseError")
		}
		if IsDatabaseError(errors.New("other error")) {
			t.Error("IsDatabaseError should return false for other errors")
		}
	})

	t.Run("IsValidationError", func(t *testing.T) {
		valErr := NewValidationError("field", "message")
		if !IsValidationError(valErr) {
			t.Error("IsValidationError should return true for ValidationError")
		}
		if IsValidationError(errors.New("other error")) {
			t.Error("IsValidationError should return false for other errors")
		}
	})
}
