package tui

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/service"
	_ "modernc.org/sqlite"
)

func TestTUIDisplaysSeededData(t *testing.T) {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "..", "..", "test-seed-final.db")

	if info, err := os.Stat(dbPath); os.IsNotExist(err) || info.Size() == 0 {
		t.Skip("test database file not found or empty, skipping integration test")
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

	areaSvc := service.NewAreaService(repo, nil)
	subareaSvc := service.NewSubareaService(repo, nil)
	projectSvc := service.NewProjectService(repo, nil)
	taskSvc := service.NewTaskService(repo, nil, nil)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)

	areaCmd := LoadAreasCmd(areaSvc)
	areaMsg := areaCmd()
	areasMsg := areaMsg.(AreasLoadedMsg)

	if len(areasMsg.Areas) != 3 {
		t.Fatalf("Expected 3 areas, got %d", len(areasMsg.Areas))
	}
	t.Logf("✓ Loaded %d areas: %s, %s, %s",
		len(areasMsg.Areas),
		areasMsg.Areas[0].Name,
		areasMsg.Areas[1].Name,
		areasMsg.Areas[2].Name)

	newModel, subareaCmd := model.Update(areasMsg)
	model = *newModel.(*Model)

	if len(model.tabs) != 3 {
		t.Errorf("Expected 3 tabs, got %d", len(model.tabs))
	} else {
		t.Logf("✓ Created %d tabs: %s (active: %v), %s (active: %v), %s (active: %v)",
			len(model.tabs),
			model.tabs[0].Name, model.tabs[0].IsActive,
			model.tabs[1].Name, model.tabs[1].IsActive,
			model.tabs[2].Name, model.tabs[2].IsActive)
	}

	if subareaCmd != nil {
		subareaMsg := subareaCmd()
		if subareasMsg, ok := subareaMsg.(SubareasLoadedMsg); ok {
			newModel, projCmd := model.Update(subareasMsg)
			model = *newModel.(*Model)

			if len(model.subareas) == 0 {
				t.Error("No subareas loaded")
			} else {
				t.Logf("✓ Loaded %d subareas for area '%s'", len(model.subareas), model.areas[0].Name)
			}

			if projCmd != nil {
				projMsg := projCmd()
				if projectsMsg, ok := projMsg.(ProjectsLoadedMsg); ok {
					newModel, taskCmd := model.Update(projectsMsg)
					model = *newModel.(*Model)

					if len(model.projects) == 0 {
						t.Error("No projects loaded")
					} else {
						t.Logf("✓ Loaded %d projects for subarea '%s'", len(model.projects), model.subareas[0].Name)
					}

					if taskCmd != nil {
						taskMsg := taskCmd()
						if tasksMsg, ok := taskMsg.(TasksLoadedMsg); ok {
							newModel, _ := model.Update(tasksMsg)
							model = *newModel.(*Model)

							if len(model.tasks) == 0 {
								t.Log("No tasks loaded for first project (this is OK if project has no tasks)")
							} else {
								t.Logf("✓ Loaded %d tasks for first project", len(model.tasks))
							}
						}
					}
				}
			}
		}
	}

	t.Logf("\n=== Final TUI State ===")
	t.Logf("Areas: %d (Tabs: %d)", len(model.areas), len(model.tabs))
	t.Logf("Subareas: %d", len(model.subareas))
	t.Logf("Projects: %d", len(model.projects))
	t.Logf("Tasks: %d", len(model.tasks))

	success := len(model.areas) > 0 &&
		len(model.tabs) == len(model.areas) &&
		len(model.subareas) > 0 &&
		len(model.projects) > 0

	if success {
		t.Log("\n✅ SUCCESS: TUI displays all seeded data correctly")
	} else {
		t.Error("\n❌ FAILED: TUI does not display all seeded data")
	}
}
