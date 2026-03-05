package tui

import (
	"context"
	"errors"
	"testing"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/service"
	"github.com/example/projectdb/internal/tui/mocks"
)

func TestLoadAreasCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		mockAreaSvc, _, _, _ := mocks.NewMockServices()
		expectedAreas := []domain.Area{
			{ID: "1", Name: "Area 1"},
			{ID: "2", Name: "Area 2"},
		}
		mocks.SetupMockAreaSuccess(mockAreaSvc, expectedAreas)

		cmd := LoadAreasCmd(mockAreaSvc)
		msg := cmd()

		loaded, ok := msg.(AreasLoadedMsg)
		if !ok {
			t.Fatal("Expected AreasLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Areas) != 2 {
			t.Errorf("Expected 2 areas, got %d", len(loaded.Areas))
		}
	})
}

func TestLoadSubareasCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		_, mockSubareaSvc, _, _ := mocks.NewMockServices()
		expectedSubareas := []domain.Subarea{
			{ID: "1", Name: "Subarea 1", AreaID: "area-1"},
		}
		mocks.SetupMockSubareaSuccess(mockSubareaSvc, expectedSubareas)

		cmd := LoadSubareasCmd(mockSubareaSvc, "area-1")
		msg := cmd()

		loaded, ok := msg.(SubareasLoadedMsg)
		if !ok {
			t.Fatal("Expected SubareasLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Subareas) != 1 {
			t.Errorf("Expected 1 subarea, got %d", len(loaded.Subareas))
		}
	})
}

func TestLoadProjectsCmd(t *testing.T) {
	tests := []struct {
		name           string
		subareaID      *string
		setupMock      func(*mocks.MockProjectService)
		wantCount      int
		wantErr        bool
		wantProjectIDs []string
	}{
		{
			name:      "recursive load - direct members only",
			subareaID: ptrToString("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				m.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
					return []domain.Project{
						{ID: "proj-1", Name: "Direct Project", SubareaID: ptrToString("subarea-1")},
					}, nil
				}
			},
			wantCount:      1,
			wantProjectIDs: []string{"proj-1"},
		},
		{
			name:      "recursive load - nested projects included",
			subareaID: ptrToString("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				m.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
					return []domain.Project{
						{ID: "root-1", Name: "Root", SubareaID: ptrToString("subarea-1")},
						{ID: "child-1", Name: "Child", ParentID: ptrToString("root-1")},
						{ID: "grandchild-1", Name: "Grandchild", ParentID: ptrToString("child-1")},
					}, nil
				}
			},
			wantCount:      3,
			wantProjectIDs: []string{"root-1", "child-1", "grandchild-1"},
		},
		{
			name:      "load all projects when subareaID is nil",
			subareaID: nil,
			setupMock: func(m *mocks.MockProjectService) {
				m.ListAllFunc = func(ctx context.Context) ([]domain.Project, error) {
					return []domain.Project{
						{ID: "proj-1", Name: "Project 1"},
						{ID: "proj-2", Name: "Project 2"},
					}, nil
				}
			},
			wantCount: 2,
		},
		{
			name:      "empty result - no projects in subarea",
			subareaID: ptrToString("empty-subarea"),
			setupMock: func(m *mocks.MockProjectService) {
				m.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
					return []domain.Project{}, nil
				}
			},
			wantCount: 0,
		},
		{
			name:      "service error - database failure",
			subareaID: ptrToString("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				m.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
					return nil, errors.New("database connection failed")
				}
			},
			wantErr: true,
		},
		{
			name:      "service error - context cancelled",
			subareaID: ptrToString("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				m.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
					return nil, context.Canceled
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, mockProjectSvc, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockProjectSvc)
			}

			cmd := LoadProjectsCmd(mockProjectSvc, tt.subareaID)
			msg := cmd()

			loaded, ok := msg.(ProjectsLoadedMsg)
			if !ok {
				t.Fatal("Expected ProjectsLoadedMsg")
			}

			if tt.wantErr {
				if loaded.Err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if loaded.Err != nil {
				t.Fatalf("unexpected error: %v", loaded.Err)
			}

			if len(loaded.Projects) != tt.wantCount {
				t.Errorf("got %d projects, want %d", len(loaded.Projects), tt.wantCount)
			}

			if tt.wantProjectIDs != nil {
				gotIDs := make([]string, len(loaded.Projects))
				for i, p := range loaded.Projects {
					gotIDs[i] = p.ID
				}
				for _, wantID := range tt.wantProjectIDs {
					found := false
					for _, gotID := range gotIDs {
						if gotID == wantID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected project ID %s not found in results", wantID)
					}
				}
			}
		})
	}
}

func TestLoadTasksCmd(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		_, _, _, mockTaskSvc := mocks.NewMockServices()
		expectedTasks := []domain.Task{
			{ID: "1", Title: "Task 1", ProjectID: "project-1"},
		}
		mocks.SetupMockTaskSuccess(mockTaskSvc, expectedTasks)

		cmd := LoadTasksCmd(mockTaskSvc, "project-1")
		msg := cmd()

		loaded, ok := msg.(TasksLoadedMsg)
		if !ok {
			t.Fatal("Expected TasksLoadedMsg")
		}
		if loaded.Err != nil {
			t.Errorf("Unexpected error: %v", loaded.Err)
		}
		if len(loaded.Tasks) != 1 {
			t.Errorf("Expected 1 task, got %d", len(loaded.Tasks))
		}
	})
}

func ptrToString(s string) *string {
	return &s
}

func strPtr(s string) *string {
	return &s
}

func TestCreateSubareaCmd(t *testing.T) {
	tests := []struct {
		name        string
		subareaName string
		areaID      string
		setupMock   func(*mocks.MockSubareaService)
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful creation",
			subareaName: "Test Subarea",
			areaID:      "area-1",
			setupMock: func(m *mocks.MockSubareaService) {
				expected := &domain.Subarea{ID: "subarea-1", Name: "Test Subarea", AreaID: "area-1"}
				mocks.SetupMockSubareaCreate(m, expected)
			},
		},
		{
			name:        "creation error - database failure",
			subareaName: "Test Subarea",
			areaID:      "area-1",
			setupMock: func(m *mocks.MockSubareaService) {
				mocks.SetupMockSubareaCreateError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:        "creation error - validation failure",
			subareaName: "",
			areaID:      "area-1",
			setupMock: func(m *mocks.MockSubareaService) {
				mocks.SetupMockSubareaCreateError(m, errors.New("name cannot be empty"))
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name:        "creation error - context cancelled",
			subareaName: "Test Subarea",
			areaID:      "area-1",
			setupMock: func(m *mocks.MockSubareaService) {
				mocks.SetupMockSubareaCreateError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, mockSubareaSvc, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockSubareaSvc)
			}

			cmd := CreateSubareaCmd(mockSubareaSvc, tt.subareaName, tt.areaID)
			msg := cmd()

			created, ok := msg.(SubareaCreatedMsg)
			if !ok {
				t.Fatal("Expected SubareaCreatedMsg")
			}

			if tt.wantErr {
				if created.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && created.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, created.Err.Error())
				}
				return
			}

			if created.Err != nil {
				t.Fatalf("unexpected error: %v", created.Err)
			}

			if created.Subarea.Name != tt.subareaName {
				t.Errorf("expected name %q, got %q", tt.subareaName, created.Subarea.Name)
			}

			if created.Subarea.AreaID != tt.areaID {
				t.Errorf("expected areaID %q, got %q", tt.areaID, created.Subarea.AreaID)
			}
		})
	}
}

func TestCreateProjectCmd(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		parentID    *string
		subareaID   *string
		setupMock   func(*mocks.MockProjectService)
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "successful creation with subarea",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   strPtr("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				expected := &domain.Project{ID: "proj-1", Name: "Test Project", SubareaID: strPtr("subarea-1")}
				mocks.SetupMockProjectCreate(m, expected)
			},
		},
		{
			name:        "successful creation with parent",
			projectName: "Nested Project",
			parentID:    strPtr("parent-1"),
			subareaID:   nil,
			setupMock: func(m *mocks.MockProjectService) {
				expected := &domain.Project{ID: "proj-1", Name: "Nested Project", ParentID: strPtr("parent-1")}
				mocks.SetupMockProjectCreate(m, expected)
			},
		},
		{
			name:        "creation error - database failure",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   strPtr("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				mocks.SetupMockProjectCreateError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:        "creation error - validation failure",
			projectName: "",
			parentID:    nil,
			subareaID:   strPtr("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				mocks.SetupMockProjectCreateError(m, errors.New("name cannot be empty"))
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name:        "creation error - context cancelled",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   strPtr("subarea-1"),
			setupMock: func(m *mocks.MockProjectService) {
				mocks.SetupMockProjectCreateError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, mockProjectSvc, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockProjectSvc)
			}

			cmd := CreateProjectCmd(mockProjectSvc, tt.projectName, tt.parentID, tt.subareaID)
			msg := cmd()

			created, ok := msg.(ProjectCreatedMsg)
			if !ok {
				t.Fatal("Expected ProjectCreatedMsg")
			}

			if tt.wantErr {
				if created.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && created.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, created.Err.Error())
				}
				return
			}

			if created.Err != nil {
				t.Fatalf("unexpected error: %v", created.Err)
			}

			if created.Project.Name != tt.projectName {
				t.Errorf("expected name %q, got %q", tt.projectName, created.Project.Name)
			}
		})
	}
}

func TestCreateTaskCmd(t *testing.T) {
	tests := []struct {
		name      string
		taskTitle string
		projectID string
		setupMock func(*mocks.MockTaskService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "successful creation",
			taskTitle: "Test Task",
			projectID: "project-1",
			setupMock: func(m *mocks.MockTaskService) {
				expected := &domain.Task{ID: "task-1", Title: "Test Task", ProjectID: "project-1"}
				mocks.SetupMockTaskCreate(m, expected)
			},
		},
		{
			name:      "creation error - database failure",
			taskTitle: "Test Task",
			projectID: "project-1",
			setupMock: func(m *mocks.MockTaskService) {
				mocks.SetupMockTaskCreateError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:      "creation error - validation failure",
			taskTitle: "",
			projectID: "project-1",
			setupMock: func(m *mocks.MockTaskService) {
				mocks.SetupMockTaskCreateError(m, errors.New("title cannot be empty"))
			},
			wantErr: true,
			errMsg:  "title cannot be empty",
		},
		{
			name:      "creation error - context cancelled",
			taskTitle: "Test Task",
			projectID: "project-1",
			setupMock: func(m *mocks.MockTaskService) {
				mocks.SetupMockTaskCreateError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, mockTaskSvc := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockTaskSvc)
			}

			cmd := CreateTaskCmd(mockTaskSvc, tt.taskTitle, tt.projectID)
			msg := cmd()

			created, ok := msg.(TaskCreatedMsg)
			if !ok {
				t.Fatal("Expected TaskCreatedMsg")
			}

			if tt.wantErr {
				if created.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && created.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, created.Err.Error())
				}
				return
			}

			if created.Err != nil {
				t.Fatalf("unexpected error: %v", created.Err)
			}

			if created.Task.Title != tt.taskTitle {
				t.Errorf("expected title %q, got %q", tt.taskTitle, created.Task.Title)
			}
		})
	}
}

func TestCreateAreaCmd(t *testing.T) {
	tests := []struct {
		name      string
		areaName  string
		color     domain.Color
		setupMock func(*mocks.MockAreaService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful creation",
			areaName: "Test Area",
			color:    "#0000FF",
			setupMock: func(m *mocks.MockAreaService) {
				expected := &domain.Area{ID: "area-1", Name: "Test Area", Color: "#0000FF"}
				mocks.SetupMockAreaCreate(m, expected)
			},
		},
		{
			name:     "creation error - database failure",
			areaName: "Test Area",
			color:    "#0000FF",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaCreateError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:     "creation error - validation failure",
			areaName: "",
			color:    "#0000FF",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaCreateError(m, errors.New("name cannot be empty"))
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAreaSvc, _, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockAreaSvc)
			}

			cmd := CreateAreaCmd(mockAreaSvc, tt.areaName, tt.color)
			msg := cmd()

			created, ok := msg.(AreaCreatedMsg)
			if !ok {
				t.Fatal("Expected AreaCreatedMsg")
			}

			if tt.wantErr {
				if created.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && created.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, created.Err.Error())
				}
				return
			}

			if created.Err != nil {
				t.Fatalf("unexpected error: %v", created.Err)
			}

			if created.Area.Name != tt.areaName {
				t.Errorf("expected name %q, got %q", tt.areaName, created.Area.Name)
			}

			if created.Area.Color != tt.color {
				t.Errorf("expected color %v, got %v", tt.color, created.Area.Color)
			}
		})
	}
}

func TestUpdateAreaCmd(t *testing.T) {
	tests := []struct {
		name      string
		areaID    string
		areaName  string
		color     domain.Color
		setupMock func(*mocks.MockAreaService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:     "successful update",
			areaID:   "area-1",
			areaName: "Updated Area",
			color:    "#FF0000",
			setupMock: func(m *mocks.MockAreaService) {
				expected := &domain.Area{ID: "area-1", Name: "Updated Area", Color: "#FF0000"}
				mocks.SetupMockAreaUpdate(m, expected)
			},
		},
		{
			name:     "update error - database failure",
			areaID:   "area-1",
			areaName: "Updated Area",
			color:    "#FF0000",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaUpdateError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:     "update error - not found",
			areaID:   "nonexistent",
			areaName: "Updated Area",
			color:    "#FF0000",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaUpdateError(m, errors.New("area not found"))
			},
			wantErr: true,
			errMsg:  "area not found",
		},
		{
			name:     "update error - context cancelled",
			areaID:   "area-1",
			areaName: "Updated Area",
			color:    "#FF0000",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaUpdateError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAreaSvc, _, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockAreaSvc)
			}

			cmd := UpdateAreaCmd(mockAreaSvc, tt.areaID, tt.areaName, tt.color)
			msg := cmd()

			updated, ok := msg.(AreaUpdatedMsg)
			if !ok {
				t.Fatal("Expected AreaUpdatedMsg")
			}

			if tt.wantErr {
				if updated.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && updated.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, updated.Err.Error())
				}
				return
			}

			if updated.Err != nil {
				t.Fatalf("unexpected error: %v", updated.Err)
			}

			if updated.Area.ID != tt.areaID {
				t.Errorf("expected ID %q, got %q", tt.areaID, updated.Area.ID)
			}

			if updated.Area.Name != tt.areaName {
				t.Errorf("expected name %q, got %q", tt.areaName, updated.Area.Name)
			}

			if updated.Area.Color != tt.color {
				t.Errorf("expected color %v, got %v", tt.color, updated.Area.Color)
			}
		})
	}
}

func TestDeleteAreaCmd(t *testing.T) {
	tests := []struct {
		name      string
		areaID    string
		hard      bool
		setupMock func(*mocks.MockAreaService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:   "successful soft delete",
			areaID: "area-1",
			hard:   false,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaDelete(m)
			},
		},
		{
			name:   "successful hard delete",
			areaID: "area-1",
			hard:   true,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaHardDelete(m)
			},
		},
		{
			name:   "soft delete error - database failure",
			areaID: "area-1",
			hard:   false,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaDeleteError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:   "hard delete error - has children",
			areaID: "area-1",
			hard:   true,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaHardDeleteError(m, errors.New("area has children"))
			},
			wantErr: true,
			errMsg:  "area has children",
		},
		{
			name:   "delete error - not found",
			areaID: "nonexistent",
			hard:   false,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaDeleteError(m, errors.New("area not found"))
			},
			wantErr: true,
			errMsg:  "area not found",
		},
		{
			name:   "delete error - context cancelled",
			areaID: "area-1",
			hard:   false,
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaDeleteError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAreaSvc, _, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockAreaSvc)
			}

			cmd := DeleteAreaCmd(mockAreaSvc, tt.areaID, tt.hard)
			msg := cmd()

			deleted, ok := msg.(AreaDeletedMsg)
			if !ok {
				t.Fatal("Expected AreaDeletedMsg")
			}

			if tt.wantErr {
				if deleted.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && deleted.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, deleted.Err.Error())
				}
				return
			}

			if deleted.Err != nil {
				t.Fatalf("unexpected error: %v", deleted.Err)
			}

			if deleted.AreaID != tt.areaID {
				t.Errorf("expected areaID %q, got %q", tt.areaID, deleted.AreaID)
			}

			if deleted.Hard != tt.hard {
				t.Errorf("expected hard %v, got %v", tt.hard, deleted.Hard)
			}
		})
	}
}

func TestReorderAreasCmd(t *testing.T) {
	tests := []struct {
		name      string
		areaIDs   []string
		setupMock func(*mocks.MockAreaService)
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "successful reorder",
			areaIDs: []string{"area-1", "area-2", "area-3"},
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaReorder(m)
			},
		},
		{
			name:    "reorder error - database failure",
			areaIDs: []string{"area-1", "area-2"},
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaReorderError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:    "reorder error - empty list",
			areaIDs: []string{},
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaReorderError(m, errors.New("area IDs cannot be empty"))
			},
			wantErr: true,
			errMsg:  "area IDs cannot be empty",
		},
		{
			name:    "reorder error - invalid area ID",
			areaIDs: []string{"area-1", "nonexistent"},
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaReorderError(m, errors.New("invalid area ID"))
			},
			wantErr: true,
			errMsg:  "invalid area ID",
		},
		{
			name:    "reorder error - context cancelled",
			areaIDs: []string{"area-1", "area-2"},
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaReorderError(m, context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAreaSvc, _, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockAreaSvc)
			}

			cmd := ReorderAreasCmd(mockAreaSvc, tt.areaIDs)
			msg := cmd()

			reordered, ok := msg.(AreasReorderedMsg)
			if !ok {
				t.Fatal("Expected AreasReorderedMsg")
			}

			if tt.wantErr {
				if reordered.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && reordered.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, reordered.Err.Error())
				}
				return
			}

			if reordered.Err != nil {
				t.Fatalf("unexpected error: %v", reordered.Err)
			}
		})
	}
}

func TestLoadAreaStatsCmd(t *testing.T) {
	tests := []struct {
		name      string
		areaID    string
		setupMock func(*mocks.MockAreaService)
		wantErr   bool
		errMsg    string
		wantStats *service.AreaStats
	}{
		{
			name:   "successful load stats",
			areaID: "area-1",
			setupMock: func(m *mocks.MockAreaService) {
				stats := &service.AreaStats{
					SubareaCount: 2,
					ProjectCount: 5,
					TaskCount:    10,
				}
				mocks.SetupMockAreaStats(m, stats)
			},
			wantStats: &service.AreaStats{
				SubareaCount: 2,
				ProjectCount: 5,
				TaskCount:    10,
			},
		},
		{
			name:   "load stats error - database failure",
			areaID: "area-1",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaStatsError(m, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "database error",
		},
		{
			name:   "load stats error - not found",
			areaID: "nonexistent",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaStatsError(m, errors.New("area not found"))
			},
			wantErr: true,
			errMsg:  "area not found",
		},
		{
			name:   "load stats error - context cancelled",
			areaID: "area-1",
			setupMock: func(m *mocks.MockAreaService) {
				mocks.SetupMockAreaStatsError(m, context.Canceled)
			},
			wantErr: true,
		},
		{
			name:   "successful load stats - zero counts",
			areaID: "area-1",
			setupMock: func(m *mocks.MockAreaService) {
				stats := &service.AreaStats{
					SubareaCount: 0,
					ProjectCount: 0,
					TaskCount:    0,
				}
				mocks.SetupMockAreaStats(m, stats)
			},
			wantStats: &service.AreaStats{
				SubareaCount: 0,
				ProjectCount: 0,
				TaskCount:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAreaSvc, _, _, _ := mocks.NewMockServices()
			if tt.setupMock != nil {
				tt.setupMock(mockAreaSvc)
			}

			cmd := LoadAreaStatsCmd(mockAreaSvc, tt.areaID)
			msg := cmd()

			loaded, ok := msg.(AreaStatsLoadedMsg)
			if !ok {
				t.Fatal("Expected AreaStatsLoadedMsg")
			}

			if tt.wantErr {
				if loaded.Err == nil {
					t.Error("expected error, got nil")
				} else if tt.errMsg != "" && loaded.Err.Error() != tt.errMsg {
					t.Errorf("expected error message %q, got %q", tt.errMsg, loaded.Err.Error())
				}
				return
			}

			if loaded.Err != nil {
				t.Fatalf("unexpected error: %v", loaded.Err)
			}

			if tt.wantStats != nil {
				if loaded.Stats.Subareas != tt.wantStats.SubareaCount {
					t.Errorf("expected subarea count %d, got %d", tt.wantStats.SubareaCount, loaded.Stats.Subareas)
				}
				if loaded.Stats.Projects != tt.wantStats.ProjectCount {
					t.Errorf("expected project count %d, got %d", tt.wantStats.ProjectCount, loaded.Stats.Projects)
				}
				if loaded.Stats.Tasks != tt.wantStats.TaskCount {
					t.Errorf("expected task count %d, got %d", tt.wantStats.TaskCount, loaded.Stats.Tasks)
				}
			}
		})
	}
}
