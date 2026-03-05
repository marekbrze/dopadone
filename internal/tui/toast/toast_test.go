package toast

import (
	"testing"
	"time"
)

func TestNewError(t *testing.T) {
	msg := "Database connection failed"
	toast := NewError(msg)

	if toast.Type != TypeError {
		t.Errorf("expected type %s, got %s", TypeError, toast.Type)
	}
	if toast.Message != msg {
		t.Errorf("expected message %s, got %s", msg, toast.Message)
	}
	if toast.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestNewSuccess(t *testing.T) {
	msg := "Item created successfully"
	toast := NewSuccess(msg)

	if toast.Type != TypeSuccess {
		t.Errorf("expected type %s, got %s", TypeSuccess, toast.Type)
	}
	if toast.Message != msg {
		t.Errorf("expected message %s, got %s", msg, toast.Message)
	}
	if toast.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestNewInfo(t *testing.T) {
	msg := "Loading data..."
	toast := NewInfo(msg)

	if toast.Type != TypeInfo {
		t.Errorf("expected type %s, got %s", TypeInfo, toast.Type)
	}
	if toast.Message != msg {
		t.Errorf("expected message %s, got %s", msg, toast.Message)
	}
	if toast.CreatedAt.IsZero() {
		t.Error("CreatedAt is zero")
	}
}

func TestIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		createdAt time.Time
		expected  bool
	}{
		{
			name:      "just created",
			createdAt: time.Now(),
			expected:  false,
		},
		{
			name:      "created 2 seconds ago",
			createdAt: time.Now().Add(-2 * time.Second),
			expected:  false,
		},
		{
			name:      "created 2.9 seconds ago",
			createdAt: time.Now().Add(-2900 * time.Millisecond),
			expected:  false,
		},
		{
			name:      "created 3.1 seconds ago",
			createdAt: time.Now().Add(-3100 * time.Millisecond),
			expected:  true,
		},
		{
			name:      "created 4 seconds ago",
			createdAt: time.Now().Add(-4 * time.Second),
			expected:  true,
		},
		{
			name:      "created 10 seconds ago",
			createdAt: time.Now().Add(-10 * time.Second),
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toast := Toast{CreatedAt: tt.createdAt}
			result := toast.IsExpired()
			if result != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if TypeError != "error" {
		t.Errorf("TypeError = %s, want error", TypeError)
	}
	if TypeSuccess != "success" {
		t.Errorf("TypeSuccess = %s, want success", TypeSuccess)
	}
	if TypeInfo != "info" {
		t.Errorf("TypeInfo = %s, want info", TypeInfo)
	}
	if ToastDuration != 3000 {
		t.Errorf("ToastDuration = %d, want 3000", ToastDuration)
	}
	if MaxToasts != 3 {
		t.Errorf("MaxToasts = %d, want 3", MaxToasts)
	}
}
