package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAreaNameEmpty = errors.New("area name cannot be empty")
)

type Area struct {
	ID        string     // Unique identifier (UUID)
	Name      string     // Display name of the area
	Color     Color      // Optional color for UI display (hex format: #RRGGBB)
	SortOrder int        // Order for display sorting (0-indexed)
	CreatedAt time.Time  // Timestamp when the area was created
	UpdatedAt time.Time  // Timestamp when the area was last updated
	DeletedAt *time.Time // Timestamp when the area was soft-deleted (nil if not deleted)
}

func NewArea(name string, color Color, sortOrder int) (*Area, error) {
	if name == "" {
		return nil, ErrAreaNameEmpty
	}

	if !color.IsValid() {
		return nil, ErrInvalidColorFormat
	}

	now := time.Now()
	return &Area{
		ID:        uuid.New().String(),
		Name:      name,
		Color:     color,
		SortOrder: sortOrder,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}, nil
}

func (a *Area) IsDeleted() bool {
	return a.DeletedAt != nil
}
