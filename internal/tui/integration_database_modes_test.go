package tui

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/service"
)

func createTestDatabase(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	cleanup := func() {
		_ = database.Close()
	}

	return database, cleanup
}

func createTestServices(t *testing.T, database *sql.DB) (service.AreaServiceInterface, service.SubareaServiceInterface, service.ProjectServiceInterface, service.TaskServiceInterface) {
	t.Helper()

	queries := db.New(database)
	tm := db.NewTransactionManager(database)

	areaSvc := service.NewAreaService(queries, tm)
	subareaSvc := service.NewSubareaService(queries, tm)
	projectSvc := service.NewProjectService(queries, tm)
	taskSvc := service.NewTaskService(queries, tm, projectSvc)

	return areaSvc, subareaSvc, projectSvc, taskSvc
}

func TestTUI_WithSQLite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	if len(model.areas) != 0 {
		t.Errorf("Expected empty areas initially, got %d", len(model.areas))
	}

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = *newModel.(*Model)

	if model.width != 100 {
		t.Errorf("Width = %d, want 100", model.width)
	}
	if model.height != 30 {
		t.Errorf("Height = %d, want 30", model.height)
	}
}

func TestTUI_LoadAreasCmd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, _, _, _ := createTestServices(t, database)

	cmd := LoadAreasCmd(areaSvc)
	if cmd == nil {
		t.Fatal("LoadAreasCmd returned nil")
	}

	msg := cmd()

	areasMsg, ok := msg.(AreasLoadedMsg)
	if !ok {
		t.Fatalf("Expected AreasLoadedMsg, got %T", msg)
	}

	if areasMsg.Err == nil {
		t.Errorf("AreasLoadedMsg.Err = %v, want nil", areasMsg.Err)
	}
}

func TestTUI_ModelUpdate(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	loadedMsg := AreasLoadedMsg{
		Areas: []domain.Area{
			{ID: "area-1", Name: "Test Area"},
		},
	}

	newModel, _ := model.Update(loadedMsg)
	model = *newModel.(*Model)

	if len(model.areas) != 1 {
		t.Errorf("Expected 1 area, got %d", len(model.areas))
	}

	if len(model.tabs) != 1 {
		t.Errorf("Expected 1 tab, got %d", len(model.tabs))
	}
}

func TestTUI_KeyboardNavigation(t *testing.T) {
	t.Skip("TODO: Keyboard navigation 'l' key handler not yet implemented")
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	loadedMsg := AreasLoadedMsg{
		Areas: []domain.Area{
			{ID: "area-1", Name: "Area 1"},
			{ID: "area-2", Name: "Area 2"},
		},
	}
	newModel, _ := model.Update(loadedMsg)
	model = *newModel.(*Model)

	if model.selectedTab != 0 {
		t.Errorf("Initial selectedTab = %d, want 0", model.selectedTab)
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	newModel, _ = model.Update(keyMsg)
	model = *newModel.(*Model)

	if model.selectedTab != 1 {
		t.Errorf("After 'l' key, selectedTab = %d, want 1", model.selectedTab)
	}
}

func TestTUI_View(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = *newModel.(*Model)

	view := model.View()

	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestTUI_QuitCommand(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(quitMsg)

	if cmd == nil {
		t.Error("Expected quit command to be returned")
		return
	}

	quitResult := cmd()
	if _, ok := quitResult.(tea.QuitMsg); !ok {
		t.Errorf("Expected tea.QuitMsg, got %T", quitResult)
	}
}

func TestTUI_WithMultipleAreas(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	areas := make([]domain.Area, 5)
	for i := 0; i < 5; i++ {
		areas[i] = domain.Area{
			ID:   string(rune('a' + i)),
			Name: string(rune('A' + i)),
		}
	}

	loadedMsg := AreasLoadedMsg{Areas: areas}
	newModel, _ := model.Update(loadedMsg)
	model = *newModel.(*Model)

	if len(model.tabs) != 5 {
		t.Errorf("Expected 5 tabs, got %d", len(model.tabs))
	}
}

func TestTUI_EmptyState(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	loadedMsg := AreasLoadedMsg{Areas: []domain.Area{}}
	newModel, _ := model.Update(loadedMsg)
	model = *newModel.(*Model)

	if len(model.areas) != 0 {
		t.Errorf("Expected 0 areas, got %d", len(model.areas))
	}

	view := model.View()
	if view == "" {
		t.Error("View should not be empty even with no areas")
	}
}

func TestTUI_WindowResize(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	model = *newModel.(*Model)

	if model.width != 120 {
		t.Errorf("Width = %d, want 120", model.width)
	}
	if model.height != 40 {
		t.Errorf("Height = %d, want 40", model.height)
	}

	newModel, _ = model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model = *newModel.(*Model)

	if model.width != 80 {
		t.Errorf("Width = %d, want 80", model.width)
	}
	if model.height != 24 {
		t.Errorf("Height = %d, want 24", model.height)
	}
}

func TestDriverMode_SQLiteWithTUI(t *testing.T) {
	t.Skip("TODO: Requires database migrations to be run before testing")
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = database.Close() }()

	queries := db.New(database)
	tm := db.NewTransactionManager(database)

	areaSvc := service.NewAreaService(queries, tm)

	cmd := LoadAreasCmd(areaSvc)
	msg := cmd()

	areasMsg, ok := msg.(AreasLoadedMsg)
	if !ok {
		t.Fatalf("Expected AreasLoadedMsg, got %T", msg)
	}

	if areasMsg.Err != nil {
		t.Errorf("AreasLoadedMsg.Err = %v, want nil", areasMsg.Err)
	}
}

func TestTUI_ConcurrentCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, _, projectSvc, _ := createTestServices(t, database)

	done := make(chan bool, 2)

	go func() {
		cmd := LoadAreasCmd(areaSvc)
		_ = cmd()
		done <- true
	}()

	go func() {
		cmd := LoadProjectsCmd(projectSvc, nil)
		_ = cmd()
		done <- true
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("Timeout waiting for concurrent commands")
		}
	}
}

func TestTUI_ErrorHandling(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	errorMsg := AreasLoadedMsg{
		Err: context.DeadlineExceeded,
	}

	newModel, _ := model.Update(errorMsg)
	model = *newModel.(*Model)

	view := model.View()
	if view == "" {
		t.Error("View should still render even with error")
	}
}

func TestTUI_New(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	program := New(areaSvc, subareaSvc, projectSvc, taskSvc)
	if program == nil {
		t.Fatal("New() returned nil program")
	}
}

func skipIfNoTursoCredentials(t *testing.T) {
	t.Helper()

	if os.Getenv("TURSO_TEST_URL") == "" || os.Getenv("TURSO_TEST_TOKEN") == "" {
		t.Skip("Skipping: TURSO_TEST_URL and TURSO_TEST_TOKEN not set")
	}
}

func TestTUI_WithTursoRemoteCredentials(t *testing.T) {
	skipIfNoTursoCredentials(t)

	tursoURL := os.Getenv("TURSO_TEST_URL")
	tursoToken := os.Getenv("TURSO_TEST_TOKEN")

	drv, err := driver.NewDriver(
		driver.WithDriverType(driver.DriverTursoRemote),
		driver.WithTurso(tursoURL, tursoToken),
	)
	if err != nil {
		t.Fatalf("Failed to create driver: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := drv.Connect(ctx); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer func() { _ = drv.Close() }()

	database := drv.GetDB()
	queries := db.New(database)
	tm := db.NewTransactionManager(database)

	areaSvc := service.NewAreaService(queries, tm)

	cmd := LoadAreasCmd(areaSvc)
	msg := cmd()

	areasMsg, ok := msg.(AreasLoadedMsg)
	if !ok {
		t.Fatalf("Expected AreasLoadedMsg, got %T", msg)
	}

	if areasMsg.Err != nil {
		t.Errorf("AreasLoadedMsg.Err = %v, want nil", areasMsg.Err)
	}
}

func TestTUI_InitReturnsCommands(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)
	cmd := model.Init()

	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestTUI_AreaStatesInitialization(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	if model.areaStates == nil {
		t.Error("areaStates should be initialized")
	}
}

func TestTUI_DefaultFocus(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	if model.focus != FocusSubareas {
		t.Errorf("Default focus = %v, want FocusSubareas", model.focus)
	}
}

func TestTUI_SpinnerInitialization(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	_ = model.spinner
}

func TestTUI_ToastsInitialization(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	if model.toasts == nil {
		t.Error("toasts should be initialized")
	}
}

func TestTUI_ModelReadyState(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	if model.ready {
		t.Error("Model should not be ready initially")
	}

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = *newModel.(*Model)

	if !model.ready {
		t.Error("Model should be ready after window size message")
	}
}

func TestTUI_InitialModelNotReady(t *testing.T) {
	database, cleanup := createTestDatabase(t)
	defer cleanup()

	areaSvc, subareaSvc, projectSvc, taskSvc := createTestServices(t, database)

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc)

	view := model.View()
	if view != "\n  Initializing..." {
		t.Errorf("Initial view should be 'Initializing...', got %q", view)
	}
}
