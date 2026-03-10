package tui

import (
	"context"
	"database/sql"
	"os"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/service"
)

func TestTUIIntegrationWithRealDB(t *testing.T) {
	dbPath := "../../../test-verify.db"

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skip("test database file not found, skipping integration test")
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	}()

	repo := db.New(database)

	t.Log("Testing repo directly...")
	areas, err := repo.ListAreas(context.TODO())
	if err != nil {
		t.Fatalf("Failed to list areas directly: %v", err)
	}
	t.Logf("Direct repo query found %d areas", len(areas))
	for i, area := range areas {
		t.Logf("  Area %d: %s (ID: %s)", i, area.Name, area.ID)
	}

	areaSvc := service.NewAreaService(repo, nil)
	subareaSvc := service.NewSubareaService(repo, nil)
	projectSvc := service.NewProjectService(repo, nil)
	taskSvc := service.NewTaskService(repo, nil, nil)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	cmd := LoadAreasCmd(areaSvc)
	msg := cmd()

	areasMsg, ok := msg.(AreasLoadedMsg)
	if !ok {
		t.Fatalf("Expected AreasLoadedMsg, got %T", msg)
	}

	t.Logf("AreasLoadedMsg - Err: %v, Areas count: %d", areasMsg.Err, len(areasMsg.Areas))

	if areasMsg.Err != nil {
		t.Fatalf("AreasLoadedMsg has error: %v", areasMsg.Err)
	}

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = newModel.(Model)

	newModel, _ = model.Update(areasMsg)
	model = newModel.(Model)

	if len(model.areas) == 0 {
		t.Error("No areas loaded in model after AreasLoadedMsg")
	} else {
		t.Logf("✓ Model has %d areas", len(model.areas))
	}

	if len(model.tabs) == 0 {
		t.Error("No tabs created from areas")
	} else {
		t.Logf("✓ Model has %d tabs", len(model.tabs))
	}
}
