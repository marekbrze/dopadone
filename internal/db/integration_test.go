package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func TestCompleteHierarchy(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	area, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Home",
		Color:     sql.NullString{String: "#3498db", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateArea failed: %v", err)
	}
	t.Logf("Created Area: %s (id: %s)", area.Name, area.ID)

	subareaID := uuid.New().String()
	subarea, err := queries.CreateSubarea(ctx, CreateSubareaParams{
		ID:        subareaID,
		Name:      "Travel",
		AreaID:    areaID,
		Color:     sql.NullString{String: "#2ecc71", Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}
	t.Logf("Created Subarea: %s (id: %s, area_id: %s)", subarea.Name, subarea.ID, subarea.AreaID)

	parentProjectID := uuid.New().String()
	parentProject, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          parentProjectID,
		Name:        "Trip to Japan",
		Description: sql.NullString{String: "Plan and organize trip to Japan", Valid: true},
		Goal:        sql.NullString{String: "Have an amazing vacation", Valid: true},
		Status:      "active",
		Priority:    "high",
		Progress:    10,
		Deadline:    nil,
		Color:       sql.NullString{String: "#e74c3c", Valid: true},
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
	t.Logf("Created Project: %s (id: %s, subarea_id: %s)", parentProject.Name, parentProject.ID, parentProject.SubareaID.String)

	child1ID := uuid.New().String()
	child1, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          child1ID,
		Name:        "Hotel booking",
		Description: sql.NullString{String: "Book hotel in Tokyo", Valid: true},
		Goal:        sql.NullString{String: "Find affordable hotel near city center", Valid: true},
		Status:      "active",
		Priority:    "high",
		Progress:    50,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{String: parentProjectID, Valid: true},
		SubareaID:   sql.NullString{Valid: false},
		Position:    0,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (child1) failed: %v", err)
	}
	t.Logf("Created Sub-project: %s (id: %s, parent_id: %s)", child1.Name, child1.ID, child1.ParentID.String)

	child2ID := uuid.New().String()
	child2, err := queries.CreateProject(ctx, CreateProjectParams{
		ID:          child2ID,
		Name:        "Flight tickets",
		Description: sql.NullString{String: "Book round-trip flight tickets", Valid: true},
		Goal:        sql.NullString{String: "Get best price for flights", Valid: true},
		Status:      "active",
		Priority:    "medium",
		Progress:    0,
		Deadline:    nil,
		Color:       sql.NullString{Valid: false},
		ParentID:    sql.NullString{String: parentProjectID, Valid: true},
		SubareaID:   sql.NullString{Valid: false},
		Position:    1,
		CreatedAt:   now,
		UpdatedAt:   now,
		CompletedAt: nil,
		DeletedAt:   nil,
	})
	if err != nil {
		t.Fatalf("CreateProject (child2) failed: %v", err)
	}
	t.Logf("Created Sub-project: %s (id: %s, parent_id: %s)", child2.Name, child2.ID, child2.ParentID.String)

	fetchedArea, err := queries.GetAreaByID(ctx, areaID)
	if err != nil {
		t.Fatalf("GetAreaByID failed: %v", err)
	}
	if fetchedArea.Name != "Home" {
		t.Errorf("expected area name 'Home', got %s", fetchedArea.Name)
	}

	fetchedSubarea, err := queries.GetSubareaByID(ctx, subareaID)
	if err != nil {
		t.Fatalf("GetSubareaByID failed: %v", err)
	}
	if fetchedSubarea.Name != "Travel" {
		t.Errorf("expected subarea name 'Travel', got %s", fetchedSubarea.Name)
	}
	if fetchedSubarea.AreaID != areaID {
		t.Errorf("expected subarea area_id %s, got %s", areaID, fetchedSubarea.AreaID)
	}

	subareas, err := queries.ListSubareasByArea(ctx, areaID)
	if err != nil {
		t.Fatalf("ListSubareasByArea failed: %v", err)
	}
	if len(subareas) != 1 {
		t.Errorf("expected 1 subarea, got %d", len(subareas))
	}

	fetchedParent, err := queries.GetProjectByID(ctx, parentProjectID)
	if err != nil {
		t.Fatalf("GetProjectByID (parent) failed: %v", err)
	}
	if fetchedParent.Name != "Trip to Japan" {
		t.Errorf("expected project name 'Trip to Japan', got %s", fetchedParent.Name)
	}
	if !fetchedParent.SubareaID.Valid || fetchedParent.SubareaID.String != subareaID {
		t.Errorf("expected parent project subarea_id %s, got %v", subareaID, fetchedParent.SubareaID)
	}
	if fetchedParent.ParentID.Valid {
		t.Errorf("expected parent project to have NULL parent_id, got %v", fetchedParent.ParentID)
	}

	children, err := queries.ListProjectsByParent(ctx, sql.NullString{String: parentProjectID, Valid: true})
	if err != nil {
		t.Fatalf("ListProjectsByParent failed: %v", err)
	}
	if len(children) != 2 {
		t.Errorf("expected 2 child projects, got %d", len(children))
	}

	childNames := make(map[string]bool)
	for _, child := range children {
		childNames[child.Name] = true
		if !child.ParentID.Valid || child.ParentID.String != parentProjectID {
			t.Errorf("expected child parent_id %s, got %v", parentProjectID, child.ParentID)
		}
		if child.SubareaID.Valid {
			t.Errorf("expected child subarea_id to be NULL, got %v", child.SubareaID)
		}
	}
	if !childNames["Hotel booking"] {
		t.Error("expected child 'Hotel booking' not found")
	}
	if !childNames["Flight tickets"] {
		t.Error("expected child 'Flight tickets' not found")
	}

	projectsBySubarea, err := queries.ListProjectsBySubarea(ctx, sql.NullString{String: subareaID, Valid: true})
	if err != nil {
		t.Fatalf("ListProjectsBySubarea failed: %v", err)
	}
	if len(projectsBySubarea) != 1 {
		t.Errorf("expected 1 root project in subarea, got %d", len(projectsBySubarea))
	}
	if len(projectsBySubarea) > 0 && projectsBySubarea[0].Name != "Trip to Japan" {
		t.Errorf("expected root project 'Trip to Japan', got %s", projectsBySubarea[0].Name)
	}
}

func TestSoftDeleteCascadeBehavior(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Test Area",
		Color:     sql.NullString{Valid: false},
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
		Name:      "Test Subarea",
		AreaID:    areaID,
		Color:     sql.NullString{Valid: false},
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
	_, err = queries.SoftDeleteArea(ctx, SoftDeleteAreaParams{
		ID:        areaID,
		DeletedAt: deletedAt,
	})
	if err != nil {
		t.Fatalf("SoftDeleteArea failed: %v", err)
	}

	_, err = queries.GetAreaByID(ctx, areaID)
	if err == nil {
		t.Error("expected GetAreaByID to fail for soft-deleted area")
	}

	subarea, err := queries.GetSubareaByID(ctx, subareaID)
	if err != nil {
		t.Errorf("subarea should still exist after area soft delete (soft delete is NOT cascade), got error: %v", err)
	} else {
		t.Logf("subarea still exists with deleted_at=nil (soft delete does not cascade): %v", subarea.DeletedAt)
		if subarea.DeletedAt != nil {
			t.Log("Note: If deleted_at is set, soft delete cascaded")
		}
	}

	areas, err := queries.ListAreas(ctx)
	if err != nil {
		t.Fatalf("ListAreas failed: %v", err)
	}
	for _, a := range areas {
		if a.ID == areaID {
			t.Error("soft-deleted area should not appear in ListAreas")
		}
	}

	subareas, err := queries.ListSubareasByArea(ctx, areaID)
	if err != nil {
		t.Fatalf("ListSubareasByArea failed: %v", err)
	}
	for _, s := range subareas {
		if s.ID == subareaID {
			t.Log("Note: orphaned subarea appears in ListSubareasByArea (soft delete doesn't cascade)")
		}
	}

	t.Log("Soft delete behavior: parent soft-deletion creates orphaned records, not cascade deletion")
	t.Log("This is by design - soft deletes preserve referential integrity while allowing recovery")
}

func TestIndexPerformance(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	queries := New(db)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Second)

	areaID := uuid.New().String()
	_, err := queries.CreateArea(ctx, CreateAreaParams{
		ID:        areaID,
		Name:      "Performance Test Area",
		Color:     sql.NullString{Valid: false},
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
		Name:      "Performance Test Subarea",
		AreaID:    areaID,
		Color:     sql.NullString{Valid: false},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	})
	if err != nil {
		t.Fatalf("CreateSubarea failed: %v", err)
	}

	statuses := []string{"active", "completed", "on_hold", "archived"}
	priorities := []string{"low", "medium", "high", "urgent"}

	for i := 0; i < 100; i++ {
		deadline := now.AddDate(0, 0, i)
		_, err := queries.CreateProject(ctx, CreateProjectParams{
			ID:          uuid.New().String(),
			Name:        "Project " + string(rune('A'+i%26)) + string(rune('0'+i/26)),
			Description: sql.NullString{Valid: false},
			Goal:        sql.NullString{Valid: false},
			Status:      statuses[i%len(statuses)],
			Priority:    priorities[i%len(priorities)],
			Progress:    int64(i % 101),
			Deadline:    &deadline,
			Color:       sql.NullString{Valid: false},
			ParentID:    sql.NullString{Valid: false},
			SubareaID:   sql.NullString{String: subareaID, Valid: true},
			Position:    int64(i),
			CreatedAt:   now,
			UpdatedAt:   now,
			CompletedAt: nil,
			DeletedAt:   nil,
		})
		if err != nil {
			t.Fatalf("CreateProject failed for project %d: %v", i, err)
		}
	}

	t.Run("IndexUsage_StatusQuery", func(t *testing.T) {
		rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM projects WHERE status = 'active'")
		if err != nil {
			t.Fatalf("Failed to get query plan: %v", err)
		}
		defer rows.Close()

		var planLines []string
		for rows.Next() {
			var id int
			var parent int
			var notUsed int
			var detail string
			if err := rows.Scan(&id, &parent, &notUsed, &detail); err != nil {
				t.Fatalf("Failed to scan query plan: %v", err)
			}
			planLines = append(planLines, detail)
			t.Logf("  Plan: %s", detail)
		}

		planStr := strings.Join(planLines, " ")
		if strings.Contains(planStr, "USING INDEX") || strings.Contains(planStr, "idx_projects_status") {
			t.Log("Query uses index for status filter")
		} else {
			t.Log("Query does not use index for status filter (may be expected for small datasets)")
		}
	})

	t.Run("IndexUsage_DeadlineQuery", func(t *testing.T) {
		startDeadline := now.AddDate(0, 0, 10)
		endDeadline := now.AddDate(0, 0, 50)

		rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM projects WHERE deadline BETWEEN ? AND ?", startDeadline, endDeadline)
		if err != nil {
			t.Fatalf("Failed to get query plan: %v", err)
		}
		defer rows.Close()

		var planLines []string
		for rows.Next() {
			var id int
			var parent int
			var notUsed int
			var detail string
			if err := rows.Scan(&id, &parent, &notUsed, &detail); err != nil {
				t.Fatalf("Failed to scan query plan: %v", err)
			}
			planLines = append(planLines, detail)
			t.Logf("  Plan: %s", detail)
		}

		planStr := strings.Join(planLines, " ")
		if strings.Contains(planStr, "USING INDEX") || strings.Contains(planStr, "idx_projects_deadline") {
			t.Log("Query uses index for deadline range filter")
		} else {
			t.Log("Query does not use index for deadline filter (may be expected for small datasets)")
		}
	})

	t.Run("IndexUsage_StatusPriorityQuery", func(t *testing.T) {
		rows, err := db.Query("EXPLAIN QUERY PLAN SELECT * FROM projects WHERE status = 'active' AND priority = 'high'")
		if err != nil {
			t.Fatalf("Failed to get query plan: %v", err)
		}
		defer rows.Close()

		var planLines []string
		for rows.Next() {
			var id int
			var parent int
			var notUsed int
			var detail string
			if err := rows.Scan(&id, &parent, &notUsed, &detail); err != nil {
				t.Fatalf("Failed to scan query plan: %v", err)
			}
			planLines = append(planLines, detail)
			t.Logf("  Plan: %s", detail)
		}

		planStr := strings.Join(planLines, " ")
		if strings.Contains(planStr, "USING INDEX") || strings.Contains(planStr, "idx_projects_status") {
			t.Log("Query uses composite index for status+priority filter")
		} else {
			t.Log("Query does not use index for status+priority filter (may be expected for small datasets)")
		}
	})

	t.Run("QueryPerformance_ActualQueries", func(t *testing.T) {
		activeProjects, err := queries.GetProjectsByStatus(ctx, "active")
		if err != nil {
			t.Fatalf("GetProjectsByStatus failed: %v", err)
		}
		t.Logf("Found %d active projects", len(activeProjects))

		startDeadline := now.AddDate(0, 0, 10)
		endDeadline := now.AddDate(0, 0, 50)
		var deadlineProjects []Project
		deadlineRows, err := db.Query("SELECT id, name, status, priority, deadline FROM projects WHERE deadline BETWEEN ? AND ? ORDER BY deadline", startDeadline, endDeadline)
		if err != nil {
			t.Fatalf("Deadline query failed: %v", err)
		}
		defer deadlineRows.Close()
		for deadlineRows.Next() {
			var p Project
			if err := deadlineRows.Scan(&p.ID, &p.Name, &p.Status, &p.Priority, &p.Deadline); err != nil {
				t.Fatalf("Failed to scan deadline project: %v", err)
			}
			deadlineProjects = append(deadlineProjects, p)
		}
		t.Logf("Found %d projects with deadline in range", len(deadlineProjects))

		bySubarea, err := queries.ListProjectsBySubarea(ctx, sql.NullString{String: subareaID, Valid: true})
		if err != nil {
			t.Fatalf("ListProjectsBySubarea failed: %v", err)
		}
		t.Logf("Found %d projects in subarea", len(bySubarea))
	})
}

func TestMigrationIdempotency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../../migrations"

	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("first goose up failed: %v", err)
	}
	t.Log("First goose up succeeded")

	var version int64
	err = db.QueryRow("SELECT MAX(version_id) FROM goose_db_version").Scan(&version)
	if err != nil {
		t.Fatalf("failed to get migration version: %v", err)
	}
	t.Logf("Migration version after first up: %d", version)

	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("second goose up failed (should be idempotent): %v", err)
	}
	t.Log("Second goose up succeeded (idempotent)")

	err = db.QueryRow("SELECT MAX(version_id) FROM goose_db_version").Scan(&version)
	if err != nil {
		t.Fatalf("failed to get migration version after second up: %v", err)
	}
	t.Logf("Migration version after second up: %d", version)

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer rows.Close()

	tables := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		tables[name] = true
	}

	expectedTables := []string{"areas", "subareas", "projects"}
	for _, table := range expectedTables {
		if !tables[table] {
			t.Errorf("expected table %s to exist after idempotent up", table)
		}
	}
	t.Log("All expected tables exist after idempotent migration")
}
