package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreateRootProject(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	projectID := uuid.New().String()
	project, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          projectID,
		Name:        "Q1 Marketing Campaign",
		Description: sql.NullString{String: "Launch new marketing campaign for Q1", Valid: true},
		Goal:        sql.NullString{String: "Increase brand awareness by 25%", Valid: true},
		Status:      "active",
		Priority:    "high",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{String: "#FF5500", Valid: true},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if project.ID != projectID {
		t.Errorf("expected ID %s, got %s", projectID, project.ID)
	}
	if project.Name != "Q1 Marketing Campaign" {
		t.Errorf("expected Name 'Q1 Marketing Campaign', got %s", project.Name)
	}
	if !project.Description.Valid || project.Description.String != "Launch new marketing campaign for Q1" {
		t.Errorf("expected Description 'Launch new marketing campaign for Q1', got %v", project.Description)
	}
	if !project.Goal.Valid || project.Goal.String != "Increase brand awareness by 25%" {
		t.Errorf("expected Goal 'Increase brand awareness by 25%%', got %v", project.Goal)
	}
	if project.Status != "active" {
		t.Errorf("expected Status 'active', got %s", project.Status)
	}
	if project.Priority != "high" {
		t.Errorf("expected Priority 'high', got %s", project.Priority)
	}
	if project.Progress != 0 {
		t.Errorf("expected Progress 0, got %d", project.Progress)
	}
	if !project.SubareaID.Valid || project.SubareaID.String != subareaID {
		t.Errorf("expected SubareaID %s, got %v", subareaID, project.SubareaID)
	}
	if project.ParentID.Valid {
		t.Errorf("expected ParentID to be NULL for root project, got %v", project.ParentID)
	}

	fetched, err := queries.GetProjectByID(ctx, projectID)
	if err != nil {
		t.Fatalf("GetProjectByID failed: %v", err)
	}
	if fetched.ID != project.ID {
		t.Errorf("expected fetched ID %s, got %s", project.ID, fetched.ID)
	}
}

func TestCreateNestedProject(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	parentID := uuid.New().String()
	parent, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          parentID,
		Name:        "Q1 Marketing Campaign",
		Description: sql.NullString{String: "Launch new marketing campaign for Q1", Valid: true},
		Goal:        sql.NullString{String: "Increase brand awareness by 25%", Valid: true},
		Status:      "active",
		Priority:    "high",
		Progress:    25,
		Deadline:    nil,
		Color:       sql.NullString{String: "#FF5500", Valid: true},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (parent) failed: %v", err)
	}

	childID := uuid.New().String()
	child, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          childID,
		Name:        "Social Media Ads",
		Description: sql.NullString{String: "Create social media ad campaign", Valid: true},
		Goal:        sql.NullString{String: "Reach 100k impressions", Valid: true},
		Status:      "active",
		Priority:    "medium",
		Progress:    50,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{String: parentID, Valid: true},
		SubareaID:   sql.NullString{Valid: false},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (child) failed: %v", err)
	}

	if !child.ParentID.Valid || child.ParentID.String != parentID {
		t.Errorf("expected ParentID %s, got %v", parentID, child.ParentID)
	}
	if child.SubareaID.Valid {
		t.Errorf("expected SubareaID to be NULL for nested project, got %v", child.SubareaID)
	}

	fetched, err := queries.GetProjectByID(ctx, childID)
	if err != nil {
		t.Fatalf("GetProjectByID failed: %v", err)
	}
	if !fetched.ParentID.Valid || fetched.ParentID.String != parent.ID {
		t.Errorf("expected child ParentID to match parent ID %s, got %v", parent.ID, fetched.ParentID)
	}
}

func TestListProjectsBySubarea(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subarea1ID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subarea1ID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea (subarea1) failed: %v", err)
	}

	subarea2ID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subarea2ID,
		Name:      "Sales",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#0000FF", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea (subarea2) failed: %v", err)
	}

	projectNames := []string{"Zebra Project", "Alpha Project", "Beta Project"}
	for _, name := range projectNames {
		_, err := queries.CreateProject(ctx, CreateProjectParams{
			ID:          uuid.New().String(),
			Name:        name,
			Description: sql.NullString{Valid: false},
			Goal:        sql.NullString{Valid: false},
			Status:      "active",
			Priority:    "medium",
			Progress:    0,
			Deadline:    nil,
			Color:       sql.NullString{Valid: false},
			ParentID:    sql.NullString{Valid: false},
			SubareaID:   sql.NullString{String: subarea1ID, Valid: true},
			Position:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
			CompletedAt: nil,
			DeletedAt:   nil,
		})
		if err != nil {
			t.Fatalf("CreateProject failed for %s: %v", name, err)
		}
	}

	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          uuid.New().String(),
		Name:        "Sales Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "medium",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subarea2ID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (subarea2) failed: %v", err)
	}

	projects, err := queries.ListProjectsBySubarea(ctx, sql.NullString{String: subarea1ID, Valid: true})
	if err != nil {
		t.Fatalf("ListProjectsBySubarea failed: %v", err)
	}

	if len(projects) != 3 {
		t.Errorf("expected 3 projects for subarea1, got %d", len(projects))
	}

	expectedOrder := []string{"Alpha Project", "Beta Project", "Zebra Project"}
	for i, expected := range expectedOrder {
		if i >= len(projects) {
			t.Errorf("missing project at index %d", i)
			continue
		}
		if projects[i].Name != expected {
			t.Errorf("expected project %d to be '%s', got '%s'", i, expected, projects[i].Name)
		}
	}

	for _, project := range projects {
		if !project.SubareaID.Valid || project.SubareaID.String != subarea1ID {
			t.Errorf("expected all projects to have SubareaID %s, got %v", subarea1ID, project.SubareaID)
		}
	}
}

func TestListProjectsByParent(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	parentID := uuid.New().String()
	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          parentID,
		Name:        "Parent Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "high",
		Progress:    10,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (parent) failed: %v", err)
	}

	childNames := []string{"Child Z", "Child A", "Child B"}
	for _, name := range childNames {
		_, err := queries.CreateProject(ctx, CreateProjectParams{
			ID:          uuid.New().String(),
			Name:        name,
			Description: sql.NullString{Valid: false},
			Goal:        sql.NullString{Valid: false},
			Status:      "active",
			Priority:    "medium",
			Progress:    0,
			Deadline:    nil,
			Color:       sql.NullString{Valid: false},
			ParentID:    sql.NullString{String: parentID, Valid: true},
			SubareaID:   sql.NullString{Valid: false},
			Position:    0,
			CreatedAt:   now,
			UpdatedAt:   now,
			CompletedAt: nil,
			DeletedAt:   nil,
		})
		if err != nil {
			t.Fatalf("CreateProject (child %s) failed: %v", name, err)
		}
	}

	otherParentID := uuid.New().String()
	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          otherParentID,
		Name:        "Other Parent",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "low",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (other parent) failed: %v", err)
	}

	children, err := queries.ListProjectsByParent(ctx, sql.NullString{String: parentID, Valid: true})
	if err != nil {
		t.Fatalf("ListProjectsByParent failed: %v", err)
	}

	if len(children) != 3 {
		t.Errorf("expected 3 children for parent, got %d", len(children))
	}

	expectedOrder := []string{"Child A", "Child B", "Child Z"}
	for i, expected := range expectedOrder {
		if i >= len(children) {
			t.Errorf("missing child at index %d", i)
			continue
		}
		if children[i].Name != expected {
			t.Errorf("expected child %d to be '%s', got '%s'", i, expected, children[i].Name)
		}
	}

	for _, child := range children {
		if !child.ParentID.Valid || child.ParentID.String != parentID {
			t.Errorf("expected all children to have ParentID %s, got %v", parentID, child.ParentID)
		}
	}
}

func TestUpdateProjectStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	projectID := uuid.New().String()
	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "medium",
		Progress:    50,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	updatedAt := time.Now().UTC().Truncate(time.Second)
	completedAt := updatedAt

	updated, err := queries.UpdateProject(ctx, UpdateProjectParams{
		ID:          projectID,
		Name:        "Test Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "completed",
		Priority:    "medium",
		Progress:    100,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		UpdatedAt:   updatedAt,
		CompletedAt: completedAt,
	})
	if err != nil {
		t.Fatalf("UpdateProject failed: %v", err)
	}

	if updated.Status != "completed" {
		t.Errorf("expected Status 'completed', got %s", updated.Status)
	}
	if updated.Progress != 100 {
		t.Errorf("expected Progress 100, got %d", updated.Progress)
	}
	if updated.CompletedAt == nil {
		t.Error("expected CompletedAt to be set when status is 'completed'")
	}
}

func TestGetProjectsByStatus(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          uuid.New().String(),
		Name:        "Active Project 1",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "high",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          uuid.New().String(),
		Name:        "Completed Project 1",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "completed",
		Priority:    "medium",
		Progress:    100,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: now,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          uuid.New().String(),
		Name:        "On Hold Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "on_hold",
		Priority:    "low",
		Progress:    25,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    2,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	activeProjects, err := queries.GetProjectsByStatus(ctx, "active")
	if err != nil {
		t.Fatalf("GetProjectsByStatus (active) failed: %v", err)
	}
	if len(activeProjects) != 1 {
		t.Errorf("expected 1 active project, got %d", len(activeProjects))
	}

	completedProjects, err := queries.GetProjectsByStatus(ctx, "completed")
	if err != nil {
		t.Fatalf("GetProjectsByStatus (completed) failed: %v", err)
	}
	if len(completedProjects) != 1 {
		t.Errorf("expected 1 completed project, got %d", len(completedProjects))
	}
	if completedProjects[0].CompletedAt == nil {
		t.Error("expected completed project to have CompletedAt set")
	}

	onHoldProjects, err := queries.GetProjectsByStatus(ctx, "on_hold")
	if err != nil {
		t.Fatalf("GetProjectsByStatus (on_hold) failed: %v", err)
	}
	if len(onHoldProjects) != 1 {
		t.Errorf("expected 1 on_hold project, got %d", len(onHoldProjects))
	}
}

func TestSoftDeleteProject(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Work",
		Color:     sql.NullString{String: "#FF0000", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}

	subareaID := uuid.New().String()
	_, err = queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Marketing",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#00FF00", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	projectID := uuid.New().String()
	_, err = queries.CreateProject(ctx, CreateProjectParams{
		ID:          projectID,
		Name:        "To Be Deleted",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "medium",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{String: subareaID, Valid: true},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	deletedAt := time.Now().UTC().Truncate(time.Second)
	deleted, err := queries.SoftDeleteProject(ctx, SoftDeleteProjectParams{
		ID:        projectID,
		DeletedAt: deletedAt,
	})
	if err != nil {
		t.Fatalf("SoftDeleteProject failed: %v", err)
	}

	if deleted.DeletedAt == nil {
		t.Error("expected DeletedAt to be set after soft delete")
	}

	_, err = queries.GetProjectByID(ctx, projectID)
	if err == nil {
		t.Error("expected GetProjectByID to fail for soft-deleted project")
	}

	projects, err := queries.ListProjectsBySubarea(ctx, sql.NullString{String: subareaID, Valid: true})
	if err != nil {
		t.Fatalf("ListProjectsBySubarea failed: %v", err)
	}

	for _, project := range projects {
		if project.ID == projectID {
			t.Error("soft-deleted project should not appear in ListProjectsBySubarea")
		}
	}
}

func TestConstraintViolationBothParentAndSubareaNull(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	projectID := uuid.New().String()
	_, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          projectID,
		Name:        "Invalid Project",
		Description: sql.NullString{Valid: false},
		Goal:        sql.NullString{Valid: false},
		Status:      "active",
		Priority:    "medium",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{Valid: false},
		SubareaID:   sql.NullString{Valid: false},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err == nil {
		t.Error("expected constraint violation when both parent_id and subarea_id are NULL")
	}
}
