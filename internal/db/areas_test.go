package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to open database: %v", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		db.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../../migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		db.Close()
		os.RemoveAll(tmpDir)
		t.Fatalf("failed to run goose up: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	return db, cleanup
}

func TestCreateArea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New().String()
	color := sql.NullString{String: "#FF0000", Valid: true}

	area, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        id,
		Name:      "Work",
		Color:     color,
		SortOrder: 1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	if area.ID != id {
		t.Errorf("expected ID %s, got %s", id, area.ID)
	}
	if area.Name != "Work" {
		t.Errorf("expected Name 'Work', got %s", area.Name)
	}
	if !area.Color.Valid || area.Color.String != "#FF0000" {
		t.Errorf("expected Color '#FF0000', got %v", area.Color)
	}
	if area.CreatedAt.Unix() != now.Unix() {
		t.Errorf("expected CreatedAt %v, got %v", now, area.CreatedAt)
	}
	if area.UpdatedAt.Unix() != now.Unix() {
		t.Errorf("expected UpdatedAt %v, got %v", now, area.UpdatedAt)
	}
	if area.DeletedAt != nil {
		t.Errorf("expected DeletedAt to be nil, got %v", area.DeletedAt)
	}

	fetched, err := queries.GetAreaByID(ctx, id)
	if err != nil {
		t.Fatalf("GetAreaByID failed: %v", err)
	}

	if fetched.ID != area.ID {
		t.Errorf("expected fetched ID %s, got %s", area.ID, fetched.ID)
	}
	if fetched.Name != area.Name {
		t.Errorf("expected fetched Name %s, got %s", area.Name, fetched.Name)
	}
	if fetched.Color.String != area.Color.String {
		t.Errorf("expected fetched Color %s, got %s", area.Color.String, fetched.Color.String)
	}
}

func TestListAreas(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaNames := []string{"Zulu", "Alpha", "Beta", "Gamma"}
	sortOrders := []int64{4, 1, 2, 3}
	for i, name := range areaNames {
		_, err := queries.CreateArea(ctx, CreateAreaParams{
			ID:        uuid.New().String(),
			Name:      name,
			Color:     sql.NullString{String: "#123456", Valid: true},
			SortOrder: sortOrders[i],
			CreatedAt: now,
			UpdatedAt: now,
			DeletedAt: nil,
		})
		if err != nil {
			t.Fatalf("CreateArea failed for %s: %v", name, err)
		}
	}

	areas, err := queries.ListAreas(ctx)
	if err != nil {
		t.Fatalf("ListAreas failed: %v", err)
	}

	if len(areas) != 4 {
		t.Errorf("expected 4 areas, got %d", len(areas))
	}

	expectedOrder := []string{"Alpha", "Beta", "Gamma", "Zulu"}
	for i, expected := range expectedOrder {
		if i >= len(areas) {
			t.Errorf("missing area at index %d", i)
			continue
		}
		if areas[i].Name != expected {
			t.Errorf("expected area %d to be '%s', got '%s'", i, expected, areas[i].Name)
		}
	}
}

func TestUpdateArea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New().String()

	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        id,
		Name:      "Original Name",
		Color:     sql.NullString{String: "#000000", Valid: true},
		SortOrder: 1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	updatedAt := time.Now().UTC().Truncate(time.Second)

	updated, err := queries.UpdateArea(ctx, UpdateAreaParams{
		ID:        id,
		Name:      "Updated Name",
		Color:     sql.NullString{String: "#FFFFFF", Valid: true},
		UpdatedAt: updatedAt,
	})
	if err != nil {
		t.Fatalf("UpdateArea failed: %v", err)
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

func TestSoftDeleteArea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	id := uuid.New().String()

	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        id,
		Name:      "To Be Deleted",
		Color:     sql.NullString{String: "#123456", Valid: true},
		SortOrder: 1,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	deletedAt := time.Now().UTC().Truncate(time.Second)
	deleted, err := queries.SoftDeleteArea(ctx, SoftDeleteAreaParams{
		ID:        id,
		DeletedAt: &deletedAt,
	})
	if err != nil {
		t.Fatalf("SoftDeleteArea failed: %v", err)
	}

	if deleted.DeletedAt == nil {
		t.Error("expected DeletedAt to be set after soft delete")
	}

	_, err = queries.GetAreaByID(ctx, id)
	if err == nil {
		t.Error("expected GetAreaByID to fail for soft-deleted area")
	}

	areas, err := queries.ListAreas(ctx)
	if err != nil {
		t.Fatalf("ListAreas failed: %v", err)
	}

	for _, area := range areas {
		if area.ID == id {
			t.Error("soft-deleted area should not appear in ListAreas")
		}
	}
}
