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

func TestTUICompleteFlow(t *testing.T) {
	wd, _ := os.Getwd()
	dbPath := filepath.Join(wd, "..", "..", "test-with-tasks.db")

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

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	areaMsg := LoadAreasCmd(areaSvc)().(AreasLoadedMsg)
	t.Logf("Loaded %d areas", len(areaMsg.Areas))

	newModel, cmd := model.Update(areaMsg)
	model = *newModel.(*Model)

	if len(model.tabs) != len(model.areas) {
		t.Errorf("Tabs mismatch: %d tabs, %d areas", len(model.tabs), len(model.areas))
	}
	t.Logf("Created %d tabs", len(model.tabs))

	if cmd != nil {
		msg := cmd().(SubareasLoadedMsg)
		newModel, cmd = model.Update(msg)
		model = *newModel.(*Model)
		t.Logf("Loaded %d subareas", len(model.subareas))

		if cmd != nil {
			msg := cmd().(ProjectsLoadedMsg)
			newModel, cmd = model.Update(msg)
			model = *newModel.(*Model)
			t.Logf("Loaded %d projects", len(model.projects))

			if cmd != nil {
				msg := cmd().(TasksLoadedMsg)
				newModel, _ = model.Update(msg)
				model = *newModel.(*Model)
				t.Logf("Loaded %d tasks", len(model.tasks))
			}
		}
	}

	if len(model.areas) == 0 {
		t.Error("No areas loaded")
	}
	if len(model.tabs) == 0 {
		t.Error("No tabs created")
	}
	if len(model.subareas) == 0 {
		t.Error("No subareas loaded")
	}
	if len(model.projects) == 0 {
		t.Error("No projects loaded")
	}
	if len(model.tasks) == 0 {
		t.Error("No tasks loaded")
	}

	if len(model.areas) > 0 && len(model.tabs) > 0 && len(model.subareas) > 0 &&
		len(model.projects) > 0 && len(model.tasks) > 0 {
		t.Log("\n✅ SUCCESS: TUI displays all seeded data (areas, subareas, projects, tasks)")
	} else {
		t.Error("\n❌ FAILED: TUI does not display all seeded data")
	}
}
