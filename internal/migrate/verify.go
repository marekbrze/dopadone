package migrate

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

type TableInfo struct {
	Name    string
	Columns int
	Indexes int
}

type SchemaVerification struct {
	LocalVersion   int64
	RemoteVersion  int64
	Consistent     bool
	Tables         []TableInfo
	ExpectedTables []string
	Errors         []string
	Warnings       []string
}

var expectedTables = []string{
	"areas",
	"subareas",
	"projects",
	"tasks",
	"goose_db_version",
}

func VerifySchema(d driver.DatabaseDriver) (*SchemaVerification, error) {
	db := d.GetDB()
	if db == nil {
		return nil, fmt.Errorf("driver not connected")
	}

	verification := &SchemaVerification{
		ExpectedTables: expectedTables,
		Consistent:     true,
	}

	localVersion, err := getGooseVersion(db)
	if err != nil {
		verification.Warnings = append(verification.Warnings,
			fmt.Sprintf("Could not determine local migration version: %v", err))
	} else {
		verification.LocalVersion = localVersion
	}

	tables, err := getTableInfo(db)
	if err != nil {
		verification.Errors = append(verification.Errors,
			fmt.Sprintf("Failed to get table info: %v", err))
		verification.Consistent = false
	} else {
		verification.Tables = tables
	}

	verification.Consistent = len(verification.Errors) == 0

	return verification, nil
}

func VerifyConsistency(d driver.DatabaseDriver) (*SchemaVerification, error) {
	verification, err := VerifySchema(d)
	if err != nil {
		return nil, err
	}

	db := d.GetDB()

	missingTables := findMissingTables(db, verification.ExpectedTables)
	if len(missingTables) > 0 {
		verification.Errors = append(verification.Errors,
			fmt.Sprintf("Missing expected tables: %s", strings.Join(missingTables, ", ")))
		verification.Consistent = false
	}

	return verification, nil
}

func getGooseVersion(db *sql.DB) (int64, error) {
	var version int64
	var dirty bool

	err := db.QueryRow(`
		SELECT MAX(version_id), false 
		FROM goose_db_version 
		WHERE is_applied = 1
	`).Scan(&version, &dirty)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to query goose_db_version: %w", err)
	}

	return version, nil
}

func getTableInfo(db *sql.DB) ([]TableInfo, error) {
	rows, err := db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
		ORDER BY name
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query tables: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var tables []TableInfo
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}

		columnCount, err := getColumnCount(db, name)
		if err != nil {
			columnCount = -1
		}

		indexCount, err := getIndexCount(db, name)
		if err != nil {
			indexCount = -1
		}

		tables = append(tables, TableInfo{
			Name:    name,
			Columns: columnCount,
			Indexes: indexCount,
		})
	}

	return tables, nil
}

func getColumnCount(db *sql.DB, tableName string) (int, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()

	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

func getIndexCount(db *sql.DB, tableName string) (int, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA index_list(%s)", tableName))
	if err != nil {
		return 0, err
	}
	defer func() { _ = rows.Close() }()

	count := 0
	for rows.Next() {
		count++
	}
	return count, nil
}

func findMissingTables(db *sql.DB, expected []string) []string {
	existingTables := make(map[string]bool)

	rows, err := db.Query(`
		SELECT name FROM sqlite_master 
		WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
	`)
	if err != nil {
		return expected
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		existingTables[name] = true
	}

	var missing []string
	for _, table := range expected {
		if !existingTables[table] {
			missing = append(missing, table)
		}
	}

	return missing
}

func (v *SchemaVerification) String() string {
	var sb strings.Builder

	sb.WriteString("Migration Status:\n")
	if v.LocalVersion > 0 {
		fmt.Fprintf(&sb, "  Local Version:  %d\n", v.LocalVersion)
	} else {
		sb.WriteString("  Local Version:  (none)\n")
	}

	if v.RemoteVersion > 0 {
		fmt.Fprintf(&sb, "  Remote Version: %d\n", v.RemoteVersion)
	}

	if v.Consistent {
		sb.WriteString("  Status:         Consistent\n")
	} else {
		sb.WriteString("  Status:         Issues detected\n")
	}

	if len(v.Tables) > 0 {
		sb.WriteString("\nTables:\n")
		for _, t := range v.Tables {
			check := "✓"
			if t.Columns < 0 {
				check = "?"
			}
			fmt.Fprintf(&sb, "  %s %s (%d columns, %d indexes)\n",
				check, t.Name, t.Columns, t.Indexes)
		}
	}

	if len(v.Warnings) > 0 {
		sb.WriteString("\nWarnings:\n")
		for _, w := range v.Warnings {
			fmt.Fprintf(&sb, "  ⚠ %s\n", w)
		}
	}

	if len(v.Errors) > 0 {
		sb.WriteString("\nErrors:\n")
		for _, e := range v.Errors {
			fmt.Fprintf(&sb, "  ✗ %s\n", e)
		}
	}

	return sb.String()
}
