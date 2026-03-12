package tui

import (
	"context"
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/tui/mocks"
	"github.com/marekbrze/dopadone/internal/tui/welcome"
)

func TestEmptyAreasShowsWelcomeModal(t *testing.T) {
	areaSvc := &mocks.MockAreaService{
		ListFunc: func(ctx context.Context) ([]domain.Area, error) {
			return []domain.Area{}, nil
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = *newModel.(*Model)

	msg := AreasLoadedMsg{Areas: []domain.Area{}}
	newModel, _ = model.Update(msg)
	model = *newModel.(*Model)

	if !model.isWelcomeOpen {
		t.Error("Expected welcome modal to be open when areas list is empty")
	}
	if model.welcomeModal == nil {
		t.Error("Expected welcome modal to be initialized")
	}
}

func TestWelcomeModalNotShownWhenAreasExist(t *testing.T) {
	areaSvc := &mocks.MockAreaService{
		ListFunc: func(ctx context.Context) ([]domain.Area, error) {
			return []domain.Area{
				{ID: "area-1", Name: "Test Area", Color: "#3B82F6"},
			}, nil
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)

	newModel, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	model = *newModel.(*Model)

	msg := AreasLoadedMsg{Areas: []domain.Area{
		{ID: "area-1", Name: "Test Area", Color: "#3B82F6"},
	}}
	newModel, _ = model.Update(msg)
	model = *newModel.(*Model)

	if model.isWelcomeOpen {
		t.Error("Expected welcome modal to NOT be open when areas exist")
	}
	if model.welcomeModal != nil {
		t.Error("Expected welcome modal to be nil when areas exist")
	}
}

func TestWelcomeModalSubmitCreatesArea(t *testing.T) {
	areaCreated := false
	areaSvc := &mocks.MockAreaService{
		CreateFunc: func(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
			areaCreated = true
			return &domain.Area{ID: "area-1", Name: name, Color: color}, nil
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isWelcomeOpen = true
	model.welcomeModal = welcome.New()

	msg := welcome.SubmitMsg{Name: "New Area", Color: "#3B82F6"}
	newModel, cmd := model.Update(msg)
	model = *newModel.(*Model)

	if model.isWelcomeOpen {
		t.Error("Expected welcome modal to be closed after submit")
	}
	if !model.isFromWelcomeFlow {
		t.Error("Expected isFromWelcomeFlow to be true after submit")
	}
	if cmd == nil {
		t.Error("Expected command to be returned for area creation")
	}

	cmd()
	if !areaCreated {
		t.Error("Expected area to be created via CreateAreaCmd")
	}
}

func TestWelcomeModalExitQuitsApp(t *testing.T) {
	areaSvc := &mocks.MockAreaService{}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isWelcomeOpen = true
	model.welcomeModal = welcome.New()

	msg := welcome.ExitMsg{}
	newModel, cmd := model.Update(msg)
	model = *newModel.(*Model)

	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}

	isQuit := false
	if cmd != nil {
		_, isQuit = cmd().(tea.QuitMsg)
	}
	if !isQuit {
		t.Error("Expected command to return tea.Quit")
	}
}

func TestAutoSelectionAfterFirstAreaCreated(t *testing.T) {
	areaSvc := &mocks.MockAreaService{
		ListFunc: func(ctx context.Context) ([]domain.Area, error) {
			return []domain.Area{
				{ID: "area-1", Name: "New Area", Color: "#3B82F6"},
			}, nil
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isFromWelcomeFlow = true

	msg := AreasLoadedMsg{Areas: []domain.Area{
		{ID: "area-1", Name: "New Area", Color: "#3B82F6"},
	}}
	newModel, cmd := model.Update(msg)
	model = *newModel.(*Model)

	if model.selectedAreaIndex != 0 {
		t.Errorf("Expected selectedAreaIndex to be 0, got %d", model.selectedAreaIndex)
	}
	if model.isFromWelcomeFlow {
		t.Error("Expected isFromWelcomeFlow to be false after areas loaded")
	}
	if cmd == nil {
		t.Error("Expected LoadSubareasCmd to be returned for auto-selection")
	}
}

func TestSpaceMenuAreaCreationStillWorks(t *testing.T) {
	areaSvc := &mocks.MockAreaService{
		CreateFunc: func(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
			return &domain.Area{ID: "area-2", Name: name, Color: color}, nil
		},
		ListFunc: func(ctx context.Context) ([]domain.Area, error) {
			return []domain.Area{
				{ID: "area-1", Name: "Existing Area", Color: "#3B82F6"},
				{ID: "area-2", Name: "New Area", Color: "#10B981"},
			}, nil
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.areas = []domain.Area{
		{ID: "area-1", Name: "Existing Area", Color: "#3B82F6"},
	}

	msg := AreasLoadedMsg{Areas: []domain.Area{
		{ID: "area-1", Name: "Existing Area", Color: "#3B82F6"},
		{ID: "area-2", Name: "New Area", Color: "#10B981"},
	}}
	newModel, _ := model.Update(msg)
	model = *newModel.(*Model)

	if model.isWelcomeOpen {
		t.Error("Expected welcome modal to NOT be open when areas exist via Space menu")
	}
}

func TestKeyboardRoutingToWelcomeModal(t *testing.T) {
	areaSvc := &mocks.MockAreaService{}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isWelcomeOpen = true
	model.welcomeModal = welcome.New()
	model.ready = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ := model.Update(msg)
	model = *newModel.(*Model)

	if !model.isWelcomeOpen {
		t.Error("Expected welcome modal to still be open after key press")
	}
}

func TestWelcomeModalSubmitWithError(t *testing.T) {
	areaSvc := &mocks.MockAreaService{
		CreateFunc: func(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
			return nil, errors.New("database error")
		},
	}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isWelcomeOpen = true
	model.welcomeModal = welcome.New()

	msg := welcome.SubmitMsg{Name: "New Area", Color: "#3B82F6"}
	newModel, cmd := model.Update(msg)
	model = *newModel.(*Model)

	if model.isWelcomeOpen {
		t.Error("Expected welcome modal to be closed after submit")
	}
	if !model.isFromWelcomeFlow {
		t.Error("Expected isFromWelcomeFlow to be true after submit")
	}

	if cmd != nil {
		createdMsg := cmd()
		if areaCreatedMsg, ok := createdMsg.(AreaCreatedMsg); ok {
			if areaCreatedMsg.Err == nil {
				t.Error("Expected error in AreaCreatedMsg")
			}
		}
	}
}

func TestWelcomeAreasLoadedWithError(t *testing.T) {
	areaSvc := &mocks.MockAreaService{}
	subareaSvc := &mocks.MockSubareaService{}
	projectSvc := &mocks.MockProjectService{}
	taskSvc := &mocks.MockTaskService{}

	model := InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, nil)
	model.isFromWelcomeFlow = true

	msg := AreasLoadedMsg{Err: errors.New("database error")}
	newModel, _ := model.Update(msg)
	model = *newModel.(*Model)

	if model.isFromWelcomeFlow {
		t.Error("Expected isFromWelcomeFlow to be false after error")
	}
	if len(model.toasts) == 0 {
		t.Error("Expected error toast to be added")
	}
}
