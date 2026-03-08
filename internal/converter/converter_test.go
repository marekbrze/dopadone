package converter

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/domain"
)

func TestDbAreaToDomain(t *testing.T) {
	now := time.Now()
	dbArea := db.Area{
		ID:        "area-1",
		Name:      "Test Area",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	result := DbAreaToDomain(dbArea)

	if result.ID != "area-1" {
		t.Errorf("Expected ID 'area-1', got %s", result.ID)
	}
	if result.Name != "Test Area" {
		t.Errorf("Expected Name 'Test Area', got %s", result.Name)
	}
	if result.Color != "#FF0000" {
		t.Errorf("Expected Color '#FF0000', got %s", result.Color)
	}
	if result.DeletedAt != nil {
		t.Error("Expected DeletedAt to be nil")
	}
}

func TestDbSubareaToDomain(t *testing.T) {
	now := time.Now()
	dbSubarea := db.Subarea{
		ID:        "subarea-1",
		Name:      "Test Subarea",
		AreaID:    "area-1",
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := DbSubareaToDomain(dbSubarea)

	if result.ID != "subarea-1" {
		t.Errorf("Expected ID 'subarea-1', got %s", result.ID)
	}
	if result.AreaID != "area-1" {
		t.Errorf("Expected AreaID 'area-1', got %s", result.AreaID)
	}
}

func TestDbProjectToDomain(t *testing.T) {
	now := time.Now()
	dbProject := db.Project{
		ID:          "project-1",
		Name:        "Test Project",
		Description: sql.NullString{String: "Description", Valid: true},
		Status:      "active",
		Priority:    "high",
		Progress:    50,
		SubareaID:   sql.NullString{String: "subarea-1", Valid: true},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := DbProjectToDomain(dbProject)

	if result.ID != "project-1" {
		t.Errorf("Expected ID 'project-1', got %s", result.ID)
	}
	if result.Name != "Test Project" {
		t.Errorf("Expected Name 'Test Project', got %s", result.Name)
	}
	if result.Status != domain.ProjectStatus("active") {
		t.Errorf("Expected Status 'active', got %s", result.Status)
	}
	if result.SubareaID == nil || *result.SubareaID != "subarea-1" {
		t.Error("Expected SubareaID 'subarea-1'")
	}
}

func TestDbTaskToDomain(t *testing.T) {
	now := time.Now()
	dbTask := db.Task{
		ID:        "task-1",
		ProjectID: "project-1",
		Title:     "Test Task",
		Status:    "todo",
		Priority:  "medium",
		IsNext:    1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := DbTaskToDomain(dbTask)

	if result.ID != "task-1" {
		t.Errorf("Expected ID 'task-1', got %s", result.ID)
	}
	if result.Title != "Test Task" {
		t.Errorf("Expected Title 'Test Task', got %s", result.Title)
	}
	if result.Status != domain.TaskStatus("todo") {
		t.Errorf("Expected Status 'todo', got %s", result.Status)
	}
	if !result.IsNext {
		t.Error("Expected IsNext to be true")
	}
}

func TestNullStringToString(t *testing.T) {
	tests := []struct {
		name     string
		input    sql.NullString
		expected string
	}{
		{"valid string", sql.NullString{String: "test", Valid: true}, "test"},
		{"null string", sql.NullString{String: "", Valid: false}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := nullStringToString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDbAreaToDomainWithNullColor(t *testing.T) {
	dbArea := db.Area{
		ID:    "area-1",
		Name:  "Test Area",
		Color: sql.NullString{Valid: false},
	}

	result := DbAreaToDomain(dbArea)

	if result.Color != "" {
		t.Errorf("Expected empty color, got %s", result.Color)
	}
}

func TestDbSubareaToDomainWithNullColor(t *testing.T) {
	dbSubarea := db.Subarea{
		ID:     "subarea-1",
		Name:   "Test Subarea",
		AreaID: "area-1",
		Color:  sql.NullString{Valid: false},
	}

	result := DbSubareaToDomain(dbSubarea)

	if result.Color != "" {
		t.Errorf("Expected empty color, got %s", result.Color)
	}
}

func TestDbProjectToDomainWithNullFields(t *testing.T) {
	dbProject := db.Project{
		ID:          "project-1",
		Name:        "Test Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{Valid: false},
	}

	result := DbProjectToDomain(dbProject)

	if result.Description != "" {
		t.Errorf("Expected empty description, got %s", result.Description)
	}
	if result.Goal != "" {
		t.Errorf("Expected empty goal, got %s", result.Goal)
	}
	if result.ParentID != nil {
		t.Error("Expected ParentID to be nil")
	}
	if result.SubareaID != nil {
		t.Error("Expected SubareaID to be nil")
	}
}

func TestDbTaskToDomainWithNullFields(t *testing.T) {
	dbTask := db.Task{
		ID:                "task-1",
		ProjectID:         "project-1",
		Title:             "Test Task",
		Description:       sql.NullString{Valid: false},
		Context:           sql.NullString{Valid: false},
		EstimatedDuration: sql.NullInt64{Valid: false},
		IsNext:            0,
	}

	result := DbTaskToDomain(dbTask)

	if result.Description != "" {
		t.Errorf("Expected empty description, got %s", result.Description)
	}
	if result.Context != "" {
		t.Errorf("Expected empty context, got %s", result.Context)
	}
	if result.IsNext {
		t.Error("Expected IsNext to be false")
	}
}
