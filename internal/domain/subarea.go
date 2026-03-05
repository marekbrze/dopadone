package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSubareaNameEmpty   = errors.New("subarea name cannot be empty")
	ErrSubareaAreaIDEmpty = errors.New("subarea area_id cannot be empty")
)

type Subarea struct {
	ID        string     // Unique identifier (UUID)
	Name      string     // Display name of the subarea
	AreaID    string     // Foreign key to parent Area
	Color     Color      // Optional color for UI display (hex format: #RRGGBB)
	CreatedAt time.Time  // Timestamp when the subarea was created
	UpdatedAt time.Time  // Timestamp when the subarea was last updated
	DeletedAt *time.Time // Timestamp when the subarea was soft-deleted (nil if not deleted)
}

func NewSubarea(name string, areaID string, color Color) (*Subarea, error) {
	if name == "" {
		return nil, ErrSubareaNameEmpty
	}
	if areaID == "" {
		return nil, ErrSubareaAreaIDEmpty
	}

	if !color.IsValid() {
		return nil, ErrInvalidColorFormat
	}

	now := time.Now()
	return &Subarea{
		ID:        uuid.New().String(),
		Name:      name,
		AreaID:    areaID,
		Color:     color,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}, nil
}

func (s *Subarea) IsDeleted() bool {
	return s.DeletedAt != nil
}

func (s *Subarea) GetEffectiveColor(parentArea *Area) Color {
	if s.Color != "" {
		return s.Color
	}
	if parentArea != nil {
		return parentArea.Color
	}
	return Color("")
}
