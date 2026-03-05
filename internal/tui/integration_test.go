package tui

import (
	"database/sql"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/db"
	_ "modernc.org/sqlite"
)

func TestTUIIntegrationWithRealDB(t *testing.T) {
	dbPath := "../../../test-verify.db"
	
	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	// First, test the repo directly
	repo := db.New(database)
	
	t.Log("Testing repo directly...")
	areas, err := repo.ListAreas(nil)
	if err != nil {
		t.Fatalf("Failed to list areas directly: %v", err)
	}
	t.Logf("Direct repo query found %d areas", len(areas))
	for i, area := range areas {
		t.Logf("  Area %d: %s (ID: %s)", i, area.Name, area.ID)
	}
	
	// Now test through the TUI
	model := InitialModel(repo)
	
	// Manually load areas
	cmd := LoadAreasCmd(repo)
	msg := cmd()
	
	// Handle AreasLoadedMsg
	areasMsg, ok := msg.(AreasLoadedMsg)
	if !ok {
		t.Fatalf("Expected AreasLoadedMsg, got %T", msg)
	}
	
	t.Logf("AreasLoadedMsg - Err: %v, Areas count: %d", areasMsg.Err, len(areasMsg.Areas))
	
	if areasMsg.Err != nil {
		t.Fatalf("AreasLoadedMsg has error: %v", areasMsg.Err)
	}
	
	// Simulate window resize
	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = newModel.(Model)
	
	// Update with areas message
	newModel, cmd = model.Update(areasMsg)
	model = newModel.(Model)
	
	// Check areas loaded
	if len(model.areas) == 0 {
		t.Error("No areas loaded in model after AreasLoadedMsg")
	} else {
		t.Logf("✓ Model has %d areas", len(model.areas))
	}
	
	// Check tabs created
	if len(model.tabs) == 0 {
		t.Error("No tabs created from areas")
	} else {
		t.Logf("✓ Model has %d tabs", len(model.tabs))
	}
}
