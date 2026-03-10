package tui

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/service"
	_ "modernc.org/sqlite"
)

func TestDatabaseConnectionAndData(t *testing.T) {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "..", "..", "test-verify.db")
	t.Logf("Database path: %s", dbPath)

	if info, err := os.Stat(dbPath); os.IsNotExist(err) || info.Size() == 0 {
		t.Skip("test database file not found or empty, skipping integration test")
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
	t.Log("✓ Database connection successful")

	repo := db.New(database)

	areas, err := repo.ListAreas(context.Background())
	if err != nil {
		t.Fatalf("Failed to list areas: %v", err)
	}

	if len(areas) == 0 {
		t.Error("No areas found in database")
	} else {
		t.Logf("✓ Found %d areas", len(areas))
		for i, area := range areas {
			t.Logf("  Area %d: ID=%s, Name=%s", i, area.ID, area.Name)
		}
	}

	if len(areas) > 0 {
		subareas, err := repo.ListSubareasByArea(context.Background(), areas[0].ID)
		if err != nil {
			t.Fatalf("Failed to list subareas: %v", err)
		}
		t.Logf("✓ Found %d subareas for area %s", len(subareas), areas[0].Name)
	}
}

func TestTUILoadAreasFromDB(t *testing.T) {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "..", "..", "test-verify.db")

	if info, err := os.Stat(dbPath); os.IsNotExist(err) || info.Size() == 0 {
		t.Skip("test database file not found or empty, skipping integration test")
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	repo := db.New(database)

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

	if areasMsg.Err != nil {
		t.Fatalf("AreasLoadedMsg has error: %v", areasMsg.Err)
	}

	if len(areasMsg.Areas) == 0 {
		t.Error("No areas in AreasLoadedMsg")
	} else {
		t.Logf("✓ LoadAreasCmd found %d areas", len(areasMsg.Areas))
		for i, area := range areasMsg.Areas {
			t.Logf("  Area %d: %s (ID: %s)", i, area.Name, area.ID)
		}
	}

	newModel, _ := model.Update(areasMsg)
	model = newModel.(Model)

	if len(model.areas) == 0 {
		t.Error("Model has no areas after update")
	} else {
		t.Logf("✓ Model has %d areas after update", len(model.areas))
	}

	if len(model.tabs) == 0 {
		t.Error("Model has no tabs after update")
	} else {
		t.Logf("✓ Model has %d tabs after update", len(model.tabs))
		for i, tab := range model.tabs {
			t.Logf("  Tab %d: %s (active: %v)", i, tab.Name, tab.IsActive)
		}
	}
}

func TestFullTUIFlowFromDB(t *testing.T) {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "..", "..", "test-verify.db")

	if info, err := os.Stat(dbPath); os.IsNotExist(err) || info.Size() == 0 {
		t.Skip("test database file not found or empty, skipping integration test")
	}

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()

	repo := db.New(database)

	areaSvc := service.NewAreaService(repo, nil)
	subareaSvc := service.NewSubareaService(repo, nil)
	projectSvc := service.NewProjectService(repo, nil)
	taskSvc := service.NewTaskService(repo, nil, nil)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	areaCmd := LoadAreasCmd(areaSvc)
	areaMsg := areaCmd()
	areasMsg := areaMsg.(AreasLoadedMsg)

	if len(areasMsg.Areas) == 0 {
		t.Fatal("No areas found")
	}
	t.Logf("✓ Loaded %d areas", len(areasMsg.Areas))

	newModel, subareaCmd := model.Update(areasMsg)
	model = newModel.(Model)

	if len(model.tabs) != len(model.areas) {
		t.Errorf("Tabs count mismatch: tabs=%d, areas=%d", len(model.tabs), len(model.areas))
	}

	if subareaCmd != nil {
		subareaMsg := subareaCmd()
		if subareasMsg, ok := subareaMsg.(SubareasLoadedMsg); ok {
			newModel, projCmd := model.Update(subareasMsg)
			model = newModel.(Model)
			t.Logf("✓ Loaded %d subareas", len(model.subareas))

			if projCmd != nil {
				projMsg := projCmd()
				if projectsMsg, ok := projMsg.(ProjectsLoadedMsg); ok {
					newModel, taskCmd := model.Update(projectsMsg)
					model = newModel.(Model)
					t.Logf("✓ Loaded %d projects", len(model.projects))

					if taskCmd != nil {
						taskMsg := taskCmd()
						if tasksMsg, ok := taskMsg.(TasksLoadedMsg); ok {
							newModel, _ := model.Update(tasksMsg)
							model = newModel.(Model)
							t.Logf("✓ Loaded %d tasks", len(model.tasks))
						}
					}
				}
			}
		}
	}

	t.Logf("\n=== Final State ===")
	t.Logf("Areas: %d (Tabs: %d)", len(model.areas), len(model.tabs))
	t.Logf("Subareas: %d", len(model.subareas))
	t.Logf("Projects: %d", len(model.projects))
	t.Logf("Tasks: %d", len(model.tasks))

	if len(model.areas) == 0 || len(model.tabs) == 0 {
		t.Error("FAILED: TUI does not display seeded data - areas or tabs are empty")
	} else {
		t.Log("✓ SUCCESS: TUI can load and display seeded data")
	}
}
