package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/service"
)

func TestCreateSubareaCmd(t *testing.T) {
	tests := []struct {
		name        string
		subareaName string
		areaID      string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful creation",
			subareaName: "Test Subarea",
			areaID:      "area-123",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			subareaName: "Test Subarea",
			areaID:      "area-123",
			mockError:   errors.New("database error"),
			expectError: true,
		},
		{
			name:        "empty name validation error",
			subareaName: "",
			areaID:      "area-123",
			mockError:   nil,
			expectError: true,
		},
		{
			name:        "empty area ID validation error",
			subareaName: "Test Subarea",
			areaID:      "",
			mockError:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			queries := db.New(sqlDB)
			subareaSvc := service.NewSubareaService(queries)

			if tt.subareaName != "" && tt.areaID != "" && tt.mockError == nil {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "name", "area_id", "color", "created_at", "updated_at", "deleted_at",
				}).AddRow(
					"test-id", tt.subareaName, tt.areaID, nil, now, now, nil,
				)
				mock.ExpectQuery("INSERT INTO subareas").
					WithArgs(sqlmock.AnyArg(), tt.subareaName, tt.areaID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(rows)
			} else if tt.mockError != nil {
				mock.ExpectQuery("INSERT INTO subareas").
					WithArgs(sqlmock.AnyArg(), tt.subareaName, tt.areaID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(tt.mockError)
			}

			cmd := CreateSubareaCmd(subareaSvc, tt.subareaName, tt.areaID)
			msg := cmd()

			createdMsg, ok := msg.(SubareaCreatedMsg)
			if !ok {
				t.Fatalf("expected SubareaCreatedMsg, got %T", msg)
			}

			if tt.expectError {
				if createdMsg.Err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if createdMsg.Err != nil {
					t.Errorf("expected no error but got %v", createdMsg.Err)
				}
				if createdMsg.Subarea.Name != tt.subareaName {
					t.Errorf("expected subarea name %q, got %q", tt.subareaName, createdMsg.Subarea.Name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
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
		mockError   error
		expectError bool
	}{
		{
			name:        "successful creation with subarea",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   strPtr("subarea-123"),
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "successful creation with parent project",
			projectName: "Nested Project",
			parentID:    strPtr("project-456"),
			subareaID:   nil,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   strPtr("subarea-123"),
			mockError:   errors.New("database error"),
			expectError: true,
		},
		{
			name:        "empty name validation error",
			projectName: "",
			parentID:    nil,
			subareaID:   strPtr("subarea-123"),
			mockError:   nil,
			expectError: true,
		},
		{
			name:        "no parent validation error",
			projectName: "Test Project",
			parentID:    nil,
			subareaID:   nil,
			mockError:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			queries := db.New(sqlDB)
			projectSvc := service.NewProjectService(queries)

			if tt.projectName != "" && (tt.parentID != nil || tt.subareaID != nil) && tt.mockError == nil {
				now := time.Now()
				var parentID interface{}
				var subareaID interface{}
				if tt.parentID != nil {
					parentID = *tt.parentID
				}
				if tt.subareaID != nil {
					subareaID = *tt.subareaID
				}
				rows := sqlmock.NewRows([]string{
					"id", "name", "description", "goal", "status", "priority", "progress",
					"deadline", "color", "parent_id", "subarea_id", "position",
					"created_at", "updated_at", "completed_at", "deleted_at",
				}).AddRow(
					"test-id", tt.projectName, nil, nil, "active", "medium", 0,
					nil, nil, parentID, subareaID, 0,
					now, now, nil, nil,
				)
				mock.ExpectQuery("INSERT INTO projects").
					WillReturnRows(rows)
			} else if tt.mockError != nil {
				mock.ExpectQuery("INSERT INTO projects").
					WillReturnError(tt.mockError)
			}

			cmd := CreateProjectCmd(projectSvc, tt.projectName, tt.parentID, tt.subareaID)
			msg := cmd()

			createdMsg, ok := msg.(ProjectCreatedMsg)
			if !ok {
				t.Fatalf("expected ProjectCreatedMsg, got %T", msg)
			}

			if tt.expectError {
				if createdMsg.Err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if createdMsg.Err != nil {
					t.Errorf("expected no error but got %v", createdMsg.Err)
				}
				if createdMsg.Project.Name != tt.projectName {
					t.Errorf("expected project name %q, got %q", tt.projectName, createdMsg.Project.Name)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestCreateTaskCmd(t *testing.T) {
	tests := []struct {
		name        string
		taskTitle   string
		projectID   string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful creation",
			taskTitle:   "Test Task",
			projectID:   "project-123",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "database error",
			taskTitle:   "Test Task",
			projectID:   "project-123",
			mockError:   errors.New("database error"),
			expectError: true,
		},
		{
			name:        "empty title validation error",
			taskTitle:   "",
			projectID:   "project-123",
			mockError:   nil,
			expectError: true,
		},
		{
			name:        "empty project ID validation error",
			taskTitle:   "Test Task",
			projectID:   "",
			mockError:   nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %v", err)
			}
			defer sqlDB.Close()

			queries := db.New(sqlDB)
			taskSvc := service.NewTaskService(queries)

			if tt.taskTitle != "" && tt.projectID != "" && tt.mockError == nil {
				now := time.Now()
				rows := sqlmock.NewRows([]string{
					"id", "project_id", "title", "description", "start_date", "deadline",
					"priority", "context", "estimated_duration", "status", "is_next",
					"created_at", "updated_at", "deleted_at",
				}).AddRow(
					"test-id", tt.projectID, tt.taskTitle, nil, nil, nil,
					"medium", nil, nil, "todo", 0,
					now, now, nil,
				)
				mock.ExpectQuery("INSERT INTO tasks").
					WillReturnRows(rows)
			} else if tt.mockError != nil {
				mock.ExpectQuery("INSERT INTO tasks").
					WillReturnError(tt.mockError)
			}

			cmd := CreateTaskCmd(taskSvc, tt.taskTitle, tt.projectID)
			msg := cmd()

			createdMsg, ok := msg.(TaskCreatedMsg)
			if !ok {
				t.Fatalf("expected TaskCreatedMsg, got %T", msg)
			}

			if tt.expectError {
				if createdMsg.Err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if createdMsg.Err != nil {
					t.Errorf("expected no error but got %v", createdMsg.Err)
				}
				if createdMsg.Task.Title != tt.taskTitle {
					t.Errorf("expected task title %q, got %q", tt.taskTitle, createdMsg.Task.Title)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
