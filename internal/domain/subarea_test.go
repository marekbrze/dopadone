package domain

import (
	"testing"
)

func TestSubareaColorInheritance(t *testing.T) {
	parentArea, err := NewArea("Work", Color("#FF0000"), 0)
	if err != nil {
		t.Fatalf("failed to create parent area: %v", err)
	}

	t.Run("uses own color when set", func(t *testing.T) {
		subarea, err := NewSubarea("Marketing", parentArea.ID, Color("#00FF00"))
		if err != nil {
			t.Fatalf("failed to create subarea: %v", err)
		}

		effectiveColor := subarea.GetEffectiveColor(parentArea)
		if effectiveColor != "#00FF00" {
			t.Errorf("expected effective color '#00FF00', got '%s'", effectiveColor)
		}
	})

	t.Run("inherits color from parent when subarea color is empty", func(t *testing.T) {
		subarea, err := NewSubarea("Marketing", parentArea.ID, Color(""))
		if err != nil {
			t.Fatalf("failed to create subarea: %v", err)
		}

		effectiveColor := subarea.GetEffectiveColor(parentArea)
		if effectiveColor != "#FF0000" {
			t.Errorf("expected effective color '#FF0000' (inherited from parent), got '%s'", effectiveColor)
		}
	})

	t.Run("returns empty when both subarea and parent have no color", func(t *testing.T) {
		colorlessArea, err := NewArea("Personal", Color(""), 1)
		if err != nil {
			t.Fatalf("failed to create colorless area: %v", err)
		}

		subarea, err := NewSubarea("Hobbies", colorlessArea.ID, Color(""))
		if err != nil {
			t.Fatalf("failed to create subarea: %v", err)
		}

		effectiveColor := subarea.GetEffectiveColor(colorlessArea)
		if effectiveColor != "" {
			t.Errorf("expected empty effective color, got '%s'", effectiveColor)
		}
	})

	t.Run("returns empty when parent area is nil", func(t *testing.T) {
		subarea, err := NewSubarea("Marketing", "some-area-id", Color(""))
		if err != nil {
			t.Fatalf("failed to create subarea: %v", err)
		}

		effectiveColor := subarea.GetEffectiveColor(nil)
		if effectiveColor != "" {
			t.Errorf("expected empty effective color when parent is nil, got '%s'", effectiveColor)
		}
	})
}

func TestNewSubareaValidation(t *testing.T) {
	t.Run("returns error when name is empty", func(t *testing.T) {
		_, err := NewSubarea("", "area-id", Color("#FF0000"))
		if err != ErrSubareaNameEmpty {
			t.Errorf("expected ErrSubareaNameEmpty, got %v", err)
		}
	})

	t.Run("returns error when area_id is empty", func(t *testing.T) {
		_, err := NewSubarea("Marketing", "", Color("#FF0000"))
		if err != ErrSubareaAreaIDEmpty {
			t.Errorf("expected ErrSubareaAreaIDEmpty, got %v", err)
		}
	})

	t.Run("returns error when color is invalid", func(t *testing.T) {
		_, err := NewSubarea("Marketing", "area-id", Color("invalid"))
		if err != ErrInvalidColorFormat {
			t.Errorf("expected ErrInvalidColorFormat, got %v", err)
		}
	})
}
