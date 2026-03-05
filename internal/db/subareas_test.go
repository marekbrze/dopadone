package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateSubarea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	areaID := uuid.New().String()

	area, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	subareaColor := sql.NullString{String: "#00FF00", Valid: true}

	subarea, err := queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     subareaColor,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	if subarea.ID != subareaID {
		t.Errorf("expected ID %s, got %s", subareaID, subarea.ID)
	}
	if subarea.Name != "Marketing" {
		t.Errorf("expected Name 'Marketing', got %s", subarea.Name)
	}
	if subarea.AreaID != areaID {
		t.Errorf("expected AreaID %s, got %s", areaID, subarea.AreaID)
	}
	if !subarea.Color.Valid || subarea.Color.String != "#00FF00" {
		t.Errorf("expected Color '#00FF00', got %v", subarea.Color)
	}

	fetched, err := queries.GetSubareaByID(ctx, subareaID)
	if err != nil {
		t.Fatalf("GetSubareaByID failed: %v", err)
	}

	if fetched.ID != subarea.ID {
		t.Errorf("expected fetched ID %s, got %s", subarea.ID, fetched.ID)
	}
	if fetched.AreaID != area.ID {
		t.Errorf("expected subarea AreaID to match parent area ID %s, got %s", area.ID, fetched.AreaID)
	}
}

func TestListSubareasByArea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	area1ID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        area1ID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed for area1: %v", err)
	}

	area2ID := uuid.New().String()
	_, err = queries.CreateArea(ctx, CreateAreaParams{
		ID:        area2ID,
		Name:      "Personal",
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed for area2: %v", err)
	}

	subareaNames := []string{"Zebra", "Alpha", "Beta", "Gamma"}
	for _, name := range subareaNames {
		_, err := queries.CreateSubarea(ctx, CreateSubareaParams{
			ID:        uuid.New().String(),
			Name:      name,
			AreaID:    area1ID,
			Color:     sql.NullString{String: "#123456", Valid: true},
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		})
		if err != nil {
			t.Fatalf("CreateSubarea failed for %s: %v", name, err)
		}
	}

	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        uuid.New().String(),
		Name:      "OtherAreaSubarea",
		AreaID:    area2ID,
		Color:     sql.NullString{String: "#654321", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed for OtherAreaSubarea: %v", err)
	}

	subareas, err := queries.ListSubareasByArea(ctx, area1ID)
	if err != nil {
		t.Fatalf("ListSubareasByArea failed: %v", err)
	}

	if len(subareas) != 4 {
		t.Errorf("expected 4 subareas for area1, got %d", len(subareas))
	}

	expectedOrder := []string{"Alpha", "Beta", "Gamma", "Zebra"}
	for i, expected := range expectedOrder {
		if i >= len(subareas) {
			t.Errorf("missing subarea at index %d", i)
			continue
		}
		if subareas[i].Name != expected {
			t.Errorf("expected subarea %d to be '%s', got '%s'", i, expected, subareas[i].Name)
		}
	}

	for _, subarea := range subareas {
		if subarea.AreaID != area1ID {
			t.Errorf("expected all subareas to have AreaID %s, got %s", area1ID, subarea.AreaID)
		}
	}
}

func TestUpdateSubarea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Original Name",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#000000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	updatedAt := time.Now().UTC().Truncate(time.Second)

	updated, err := queries.UpdateSubarea(ctx, UpdateSubareaParams{
		ID:        subareaID,
		Name:      "Updated Name",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#FFFFFF", Valid: true},
		UpdatedAt: updatedAt,
	})
	if err != nil {
		t.Fatalf("UpdateSubarea failed: %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("expected Name 'Updated Name', got %s", updated.Name)
	}
	if !updated.Color.Valid || updated.Color.String != "#FFFFFF" {
		t.Errorf("expected Color '#FFFFFF', got %v", updated.Color)
	}
	if !updated.UpdatedAt.After(now) && updated.UpdatedAt.Unix() != now.Unix() {
		t.Errorf("expected UpdatedAt to change, got %v", updated.UpdatedAt)
	}
	if updated.CreatedAt.Unix() != now.Unix() {
		t.Errorf("expected CreatedAt to remain unchanged, got %v", updated.CreatedAt)
	}
}

func TestSoftDeleteSubarea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "To Be Deleted",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#123456", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	deletedAt := time.Now().UTC().Truncate(time.Second)
	deleted, err := queries.SoftDeleteSubarea(ctx, SoftDeleteSubareaParams{
		ID:        subareaID,
		DeletedAt: &deletedAt,
	})
	if err != nil {
		t.Fatalf("SoftDeleteSubarea failed: %v", err)
	}

	if deleted.DeletedAt == nil {
		t.Error("expected DeletedAt to be set after soft delete")
	}

	_, err = queries.GetSubareaByID(ctx, subareaID)
	if err == nil {
		t.Error("expected GetSubareaByID to fail for soft-deleted subarea")
	}

	subareas, err := queries.ListSubareasByArea(ctx, areaID)
	if err != nil {
		t.Fatalf("ListSubareasByArea failed: %v", err)
	}

	for _, subarea := range subareas {
		if subarea.ID == subareaID {
			t.Error("soft-deleted subarea should not appear in ListSubareasByArea")
		}
	}
}
