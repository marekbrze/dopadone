package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func TestMigrationUp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	}()

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../../migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("failed to run goose up: %v", err)
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Logf("failed to close rows: %v", err)
		}
	}()

	tables := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		tables[name] = true
	}

	expectedTables := []string{"areas", "subareas", "projects", "tasks"}
	for _, table := range expectedTables {
		if !tables[table] {
			t.Errorf("expected table %s to exist, but it was not found", table)
		}
	}

	if !tables["goose_db_version"] {
		t.Error("expected goose_db_version table to exist for migration tracking")
	}

	fkRows, err := db.Query("PRAGMA foreign_key_list(subareas)")
	if err != nil {
		t.Fatalf("failed to query foreign keys for subareas: %v", err)
	}
	defer func() {
		if err := fkRows.Close(); err != nil {
			t.Logf("failed to close fkRows: %v", err)
		}
	}()

	var hasAreaFK bool
	for fkRows.Next() {
		var id, seq int
		var table, from, to string
		var onUpdate, onDelete, match string
		if err := fkRows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("failed to scan foreign key row: %v", err)
		}
		if table == "areas" && from == "area_id" && to == "id" {
			hasAreaFK = true
		}
	}
	if !hasAreaFK {
		t.Error("expected foreign key from subareas.area_id to areas.id")
	}

	fkRows, err = db.Query("PRAGMA foreign_key_list(projects)")
	if err != nil {
		t.Fatalf("failed to query foreign keys for projects: %v", err)
	}
	defer func() {
		if err := fkRows.Close(); err != nil {
			t.Logf("failed to close fkRows: %v", err)
		}
	}()

	var hasParentFK, hasSubareaFK bool
	for fkRows.Next() {
		var id, seq int
		var table, from, to string
		var onUpdate, onDelete, match string
		if err := fkRows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("failed to scan foreign key row: %v", err)
		}
		if table == "projects" && from == "parent_id" && to == "id" {
			hasParentFK = true
		}
		if table == "subareas" && from == "subarea_id" && to == "id" {
			hasSubareaFK = true
		}
	}
	if !hasParentFK {
		t.Error("expected foreign key from projects.parent_id to projects.id")
	}
	if !hasSubareaFK {
		t.Error("expected foreign key from projects.subarea_id to subareas.id")
	}

	fkRows, err = db.Query("PRAGMA foreign_key_list(tasks)")
	if err != nil {
		t.Fatalf("failed to query foreign keys for tasks: %v", err)
	}
	defer func() {
		if err := fkRows.Close(); err != nil {
			t.Logf("failed to close fkRows: %v", err)
		}
	}()

	var hasProjectFK bool
	for fkRows.Next() {
		var id, seq int
		var table, from, to string
		var onUpdate, onDelete, match string
		if err := fkRows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			t.Fatalf("failed to scan foreign key row: %v", err)
		}
		if table == "projects" && from == "project_id" && to == "id" {
			hasProjectFK = true
		}
	}
	if !hasProjectFK {
		t.Error("expected foreign key from tasks.project_id to projects.id")
	}

	indexRows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name LIKE 'idx_%' ORDER BY name")
	if err != nil {
		t.Fatalf("failed to query indexes: %v", err)
	}
	defer func() {
		if err := indexRows.Close(); err != nil {
			t.Logf("failed to close indexRows: %v", err)
		}
	}()

	indexes := make(map[string]bool)
	for indexRows.Next() {
		var name string
		if err := indexRows.Scan(&name); err != nil {
			t.Fatalf("failed to scan index name: %v", err)
		}
		indexes[name] = true
	}

	expectedIndexes := []string{
		"idx_projects_deadline",
		"idx_projects_status_priority",
		"idx_projects_parent_id",
		"idx_projects_subarea_id",
		"idx_subareas_area_id",
		"idx_tasks_project_id",
		"idx_tasks_status",
		"idx_tasks_deadline",
		"idx_tasks_is_next",
		"idx_tasks_priority",
	}
	for _, idx := range expectedIndexes {
		if !indexes[idx] {
			t.Errorf("expected index %s to exist", idx)
		}
	}

	_, err = db.Exec("INSERT INTO projects (id, name, status, priority, subarea_id) VALUES ('test1', 'Test', 'invalid_status', 'medium', NULL)")
	if err == nil {
		t.Error("expected error for invalid status value")
	}

	_, err = db.Exec("INSERT INTO projects (id, name, status, priority, subarea_id) VALUES ('test2', 'Test', 'active', 'invalid_priority', NULL)")
	if err == nil {
		t.Error("expected error for invalid priority value")
	}

	for _, table := range expectedTables {
		colRows, err := db.Query("PRAGMA table_info(" + table + ")")
		if err != nil {
			t.Fatalf("failed to query columns for table %s: %v", table, err)
		}

		var hasDeletedAt, hasCreatedAt, hasUpdatedAt bool
		for colRows.Next() {
			var cid int
			var name, colType string
			var notNull, pk int
			var dfltValue interface{}
			if err := colRows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
				t.Fatalf("failed to scan column row: %v", err)
			}
			switch name {
			case "deleted_at":
				hasDeletedAt = true
			case "created_at":
				hasCreatedAt = true
			case "updated_at":
				hasUpdatedAt = true
			}
		}
		if err := colRows.Close(); err != nil {
			t.Logf("failed to close colRows: %v", err)
		}

		if !hasDeletedAt {
			t.Errorf("expected deleted_at column in table %s for soft delete support", table)
		}
		if !hasCreatedAt {
			t.Errorf("expected created_at column in table %s", table)
		}
		if !hasUpdatedAt {
			t.Errorf("expected updated_at column in table %s", table)
		}
	}

	colRows, err := db.Query("PRAGMA table_info(projects)")
	if err != nil {
		t.Fatalf("failed to query columns for projects table: %v", err)
	}
	var hasCompletedAt bool
	for colRows.Next() {
		var cid int
		var name, colType string
		var notNull, pk int
		var dfltValue interface{}
		if err := colRows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			t.Fatalf("failed to scan column row: %v", err)
		}
		if name == "completed_at" {
			hasCompletedAt = true
			break
		}
	}
	if err := colRows.Close(); err != nil {
		t.Logf("failed to close colRows: %v", err)
	}

	if !hasCompletedAt {
		t.Error("expected completed_at column in projects table")
	}
}

func TestMigrationDown(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "db_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("failed to remove temp dir: %v", err)
		}
	}()

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close database: %v", err)
		}
	}()

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../../migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("failed to run goose up: %v", err)
	}

	if err := goose.DownTo(db, migrationsDir, 0); err != nil {
		t.Fatalf("failed to run goose down: %v", err)
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' AND name != 'goose_db_version'")
	if err != nil {
		t.Fatalf("failed to query tables: %v", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			t.Logf("failed to close rows: %v", err)
		}
	}()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("failed to scan table name: %v", err)
		}
		tables = append(tables, name)
	}

	if len(tables) > 0 {
		t.Errorf("expected all application tables to be dropped, but found: %v", tables)
	}
}
