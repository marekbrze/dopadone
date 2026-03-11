package main

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/service"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func TestSQLite_ComprehensiveCRUD(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("Area CRUD operations", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Test Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := services.Areas.GetByID(ctx, area.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Name != "Test Area" {
			t.Errorf("Name = %v, want %v", got.Name, "Test Area")
		}

		_, err = services.Areas.Update(ctx, area.ID, "Updated Area", domain.Color("#ff0000"))
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		err = services.Areas.SoftDelete(ctx, area.ID)
		if err != nil {
			t.Fatalf("SoftDelete() error = %v", err)
		}

		_, err = services.Areas.GetByID(ctx, area.ID)
		if err == nil {
			t.Error("GetByID() should fail for soft-deleted area")
		}
	})

	t.Run("Subarea CRUD operations", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Parent Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Test Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := services.Subareas.GetByID(ctx, subarea.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Name != "Test Subarea" {
			t.Errorf("Name = %v, want %v", got.Name, "Test Subarea")
		}

		err = services.Subareas.SoftDelete(ctx, subarea.ID)
		if err != nil {
			t.Fatalf("SoftDelete() error = %v", err)
		}
	})

	t.Run("Project CRUD operations", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Parent Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Parent Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		project, err := services.Projects.Create(ctx, service.CreateProjectParams{
			Name:      "Test Project",
			SubareaID: &subarea.ID,
			Status:    domain.ProjectStatusActive,
			Priority:  domain.PriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := services.Projects.GetByID(ctx, project.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Name != "Test Project" {
			t.Errorf("Name = %v, want %v", got.Name, "Test Project")
		}

		err = services.Projects.SoftDelete(ctx, project.ID)
		if err != nil {
			t.Fatalf("SoftDelete() error = %v", err)
		}
	})

	t.Run("Task CRUD operations", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Parent Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Parent Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		project, err := services.Projects.Create(ctx, service.CreateProjectParams{
			Name:      "Test Project",
			SubareaID: &subarea.ID,
			Status:    domain.ProjectStatusActive,
			Priority:  domain.PriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
			ProjectID: project.ID,
			Title:     "Test Task",
			Status:    domain.TaskStatusTodo,
			Priority:  domain.TaskPriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := services.Tasks.GetByID(ctx, task.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Title != "Test Task" {
			t.Errorf("Title = %v, want %v", got.Title, "Test Task")
		}

		_, err = services.Tasks.SetStatus(ctx, task.ID, domain.TaskStatusDone)
		if err != nil {
			t.Fatalf("SetStatus() error = %v", err)
		}

		err = services.Tasks.SoftDelete(ctx, task.ID)
		if err != nil {
			t.Fatalf("SoftDelete() error = %v", err)
		}
	})

	t.Run("Full hierarchy CRUD", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Full Hierarchy Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Full Hierarchy Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		project, err := services.Projects.Create(ctx, service.CreateProjectParams{
			Name:      "Full Hierarchy Project",
			SubareaID: &subarea.ID,
			Status:    domain.ProjectStatusActive,
			Priority:  domain.PriorityHigh,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
			ProjectID: project.ID,
			Title:     "Full Hierarchy Task",
			Status:    domain.TaskStatusTodo,
			Priority:  domain.TaskPriorityHigh,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		got, err := services.Tasks.GetByID(ctx, task.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}

		if got.Title != "Full Hierarchy Task" {
			t.Errorf("Title = %v, want %v", got.Title, "Full Hierarchy Task")
		}
	})
}

func TestSQLite_TransactionHandling(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	area, err := services.Areas.Create(ctx, "Transaction Test Area", domain.Color(""))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("successful transaction", func(t *testing.T) {
		err := services.TM.WithTransaction(ctx, func(txCtx context.Context, _ db.Querier) error {
			_, err := services.Subareas.Create(txCtx, "In Transaction", area.ID, domain.Color(""))
			return err
		})
		if err != nil {
			t.Errorf("WithTransaction() error = %v", err)
		}
	})
}

func TestSQLite_ConcurrentAccess(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	area, err := services.Areas.Create(ctx, "Concurrent Test Area", domain.Color(""))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	t.Run("concurrent creates with retry", func(t *testing.T) {
		const numGoroutines = 10
		var wg sync.WaitGroup
		errCh := make(chan error, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				var lastErr error
				for retries := 0; retries < 5; retries++ {
					_, err := services.Subareas.Create(ctx, "Concurrent Subarea", area.ID, domain.Color(""))
					if err == nil {
						errCh <- nil
						return
					}
					lastErr = err
					time.Sleep(time.Duration(retries+1) * 10 * time.Millisecond)
				}
				errCh <- lastErr
			}(i)
		}

		wg.Wait()
		close(errCh)

		successCount := 0
		for err := range errCh {
			if err == nil {
				successCount++
			}
		}

		if successCount < numGoroutines/3 {
			t.Errorf("Expected at least %d successful creates, got %d", numGoroutines/3, successCount)
		} else {
			t.Logf("Successful creates: %d/%d", successCount, numGoroutines)
		}
	})

	t.Run("concurrent reads", func(t *testing.T) {
		const numReaders = 20
		var wg sync.WaitGroup
		errCh := make(chan error, numReaders)

		for i := 0; i < numReaders; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := services.Areas.List(ctx)
				errCh <- err
			}()
		}

		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				t.Errorf("Concurrent read error: %v", err)
			}
		}
	})
}

func TestSQLite_BackwardCompatibility(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("default behavior without flags", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Default Mode Test", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if area.ID == "" {
			t.Error("Area should be created successfully in default SQLite mode")
		}
	})

	t.Run("all services work with SQLite", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Compatibility Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Compatibility Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		project, err := services.Projects.Create(ctx, service.CreateProjectParams{
			Name:      "Compatibility Project",
			SubareaID: &subarea.ID,
			Status:    domain.ProjectStatusActive,
			Priority:  domain.PriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
			ProjectID: project.ID,
			Title:     "Compatibility Task",
			Status:    domain.TaskStatusTodo,
			Priority:  domain.TaskPriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if area.ID == "" || subarea.ID == "" || project.ID == "" || task.ID == "" {
			t.Error("All entities should be created successfully")
		}
	})
}

func TestSQLite_CascadeOperations(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	t.Run("soft delete project cascades to tasks", func(t *testing.T) {
		area, err := services.Areas.Create(ctx, "Cascade Area", domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		subarea, err := services.Subareas.Create(ctx, "Cascade Subarea", area.ID, domain.Color(""))
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		project, err := services.Projects.Create(ctx, service.CreateProjectParams{
			Name:      "Cascade Project",
			SubareaID: &subarea.ID,
			Status:    domain.ProjectStatusActive,
			Priority:  domain.PriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		task, err := services.Tasks.Create(ctx, service.CreateTaskParams{
			ProjectID: project.ID,
			Title:     "Cascade Task",
			Status:    domain.TaskStatusTodo,
			Priority:  domain.TaskPriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		err = services.Projects.SoftDeleteWithCascade(ctx, project.ID)
		if err != nil {
			t.Fatalf("SoftDeleteWithCascade() error = %v", err)
		}

		_, err = services.Tasks.GetByID(ctx, task.ID)
		if err == nil {
			t.Error("Task should be soft-deleted when project is cascade deleted")
		}
	})
}

func TestSQLite_QueryPerformance(t *testing.T) {
	services, cleanup := setupTestServices(t)
	defer cleanup()
	ctx := context.Background()

	area, err := services.Areas.Create(ctx, "Perf Area", domain.Color(""))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	subarea, err := services.Subareas.Create(ctx, "Perf Subarea", area.ID, domain.Color(""))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	project, err := services.Projects.Create(ctx, service.CreateProjectParams{
		Name:      "Perf Project",
		SubareaID: &subarea.ID,
		Status:    domain.ProjectStatusActive,
		Priority:  domain.PriorityMedium,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	for i := 0; i < 50; i++ {
		_, err := services.Tasks.Create(ctx, service.CreateTaskParams{
			ProjectID: project.ID,
			Title:     fmt.Sprintf("Perf Task %d", i),
			Status:    domain.TaskStatusTodo,
			Priority:  domain.TaskPriorityMedium,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
	}

	t.Run("list tasks performance", func(t *testing.T) {
		start := time.Now()
		tasks, err := services.Tasks.ListByProject(ctx, project.ID)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("ListByProject() error = %v", err)
		}
		if len(tasks) != 50 {
			t.Errorf("Expected 50 tasks, got %d", len(tasks))
		}
		if elapsed > 100*time.Millisecond {
			t.Logf("Warning: ListByProject took %v", elapsed)
		}
	})

	t.Run("list all areas performance", func(t *testing.T) {
		start := time.Now()
		areas, err := services.Areas.List(ctx)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("List() error = %v", err)
		}
		if elapsed > 50*time.Millisecond {
			t.Logf("Warning: List took %v", elapsed)
		}
		_ = areas
	})
}

type TestServices struct {
	Areas    service.AreaServiceInterface
	Subareas service.SubareaServiceInterface
	Projects service.ProjectServiceInterface
	Tasks    service.TaskServiceInterface
	TM       *db.TransactionManager
	DB       *sql.DB
}

func setupTestServices(t *testing.T) (*TestServices, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("Failed to set goose dialect: %v", err)
	}

	migrationsDir := "../../migrations"
	if err := goose.Up(database, migrationsDir); err != nil {
		_ = database.Close()
		t.Fatalf("Failed to run migrations: %v", err)
	}

	queries := db.New(database)
	tm := db.NewTransactionManager(database)

	projectSvc := service.NewProjectService(queries, tm)
	taskSvc := service.NewTaskService(queries, tm, projectSvc)

	services := &TestServices{
		Areas:    service.NewAreaService(queries, tm),
		Subareas: service.NewSubareaService(queries, tm),
		Projects: projectSvc,
		Tasks:    taskSvc,
		TM:       tm,
		DB:       database,
	}

	cleanup := func() {
		if err := database.Close(); err != nil {
			t.Logf("Failed to close database: %v", err)
		}
	}

	return services, cleanup
}

func TestDriverDetection_AllModes(t *testing.T) {
	tests := []struct {
		name         string
		dbPath       string
		tursoURL     string
		tursoToken   string
		dbMode       string
		expectedType driver.DriverType
	}{
		{
			name:         "auto_detect_local",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "",
			expectedType: driver.DriverSQLite,
		},
		{
			name:         "explicit_local",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "local",
			expectedType: driver.DriverSQLite,
		},
		{
			name:         "explicit_sqlite_alias",
			dbPath:       "/tmp/test.db",
			tursoURL:     "",
			tursoToken:   "",
			dbMode:       "sqlite",
			expectedType: driver.DriverSQLite,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := LoadConfig(tt.dbPath, tt.tursoURL, tt.tursoToken, tt.dbMode, 60*time.Second)
			driverCfg := cfg.ToDriverConfig()

			result, err := driver.DetectOrExplicitMode(driverCfg)
			if err != nil {
				t.Fatalf("DetectOrExplicitMode() error = %v", err)
			}

			if result.Type != tt.expectedType {
				t.Errorf("Type = %v, want %v", result.Type, tt.expectedType)
			}
		})
	}
}
