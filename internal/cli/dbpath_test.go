package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultDBPath(t *testing.T) {
	t.Run("returns valid path with user config dir", func(t *testing.T) {
		path := DefaultDBPath()

		if path == "" {
			t.Error("DefaultDBPath() returned empty string")
		}

		configDir, err := os.UserConfigDir()
		if err == nil {
			expectedSuffix := filepath.Join(AppName, DBName)
			expected := filepath.Join(configDir, expectedSuffix)
			if path != expected {
				t.Errorf("DefaultDBPath() = %q, want %q", path, expected)
			}
		}
	})

	t.Run("path contains app name and db name", func(t *testing.T) {
		path := DefaultDBPath()

		if path == "./"+DBName {
			return
		}

		if !filepath.IsAbs(path) {
			t.Errorf("DefaultDBPath() returned relative path %q, expected absolute", path)
		}

		if filepath.Base(path) != DBName {
			t.Errorf("DefaultDBPath() filename = %q, want %q", filepath.Base(path), DBName)
		}

		dir := filepath.Dir(path)
		if filepath.Base(dir) != AppName {
			t.Errorf("DefaultDBPath() directory name = %q, want %q", filepath.Base(dir), AppName)
		}
	})
}

func TestGetDBPathWithFallback(t *testing.T) {
	t.Run("returns path and fallback indicator", func(t *testing.T) {
		path, fallback := GetDBPathWithFallback()

		if path == "" {
			t.Error("GetDBPathWithFallback() returned empty path")
		}

		configDir, err := os.UserConfigDir()
		if err != nil {
			if !fallback {
				t.Error("Expected fallback=true when UserConfigDir fails")
			}
			if path != "./"+DBName {
				t.Errorf("Fallback path = %q, want %q", path, "./"+DBName)
			}
		} else {
			if fallback {
				t.Error("Expected fallback=false when UserConfigDir succeeds")
			}
			expected := filepath.Join(configDir, AppName, DBName)
			if path != expected {
				t.Errorf("Path = %q, want %q", path, expected)
			}
		}
	})
}

func TestEnsureDirExists(t *testing.T) {
	t.Run("creates directory if not exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "subdir", "anotherdir", "test.db")

		if err := EnsureDirExists(dbPath); err != nil {
			t.Fatalf("EnsureDirExists() error = %v", err)
		}

		dir := filepath.Dir(dbPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory %q was not created", dir)
		}
	})

	t.Run("no error if directory exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		dbPath := filepath.Join(tmpDir, "test.db")

		if err := EnsureDirExists(dbPath); err != nil {
			t.Fatalf("First EnsureDirExists() error = %v", err)
		}

		if err := EnsureDirExists(dbPath); err != nil {
			t.Fatalf("Second EnsureDirExists() error = %v", err)
		}
	})

	t.Run("works with absolute path", func(t *testing.T) {
		tmpDir := t.TempDir()
		absTmpDir, err := filepath.Abs(tmpDir)
		if err != nil {
			t.Fatalf("Failed to get absolute path: %v", err)
		}
		dbPath := filepath.Join(absTmpDir, "subdir", "test.db")

		if err := EnsureDirExists(dbPath); err != nil {
			t.Fatalf("EnsureDirExists() error = %v", err)
		}
	})
}

func TestMigrateFromOldPath(t *testing.T) {
	t.Run("no migration if old path does not exist", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldPath := filepath.Join(tmpDir, "nonexistent.db")
		newPath := filepath.Join(tmpDir, "new", "new.db")

		if err := MigrateFromOldPath(oldPath, newPath); err != nil {
			t.Fatalf("MigrateFromOldPath() error = %v", err)
		}

		if _, err := os.Stat(newPath); !os.IsNotExist(err) {
			t.Error("New file should not exist when old file doesn't exist")
		}
	})

	t.Run("no migration if new path already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldPath := filepath.Join(tmpDir, "old.db")
		newPath := filepath.Join(tmpDir, "new.db")

		oldContent := []byte("old content")
		if err := os.WriteFile(oldPath, oldContent, 0644); err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		newContent := []byte("new content")
		if err := os.WriteFile(newPath, newContent, 0644); err != nil {
			t.Fatalf("Failed to create new file: %v", err)
		}

		if err := MigrateFromOldPath(oldPath, newPath); err != nil {
			t.Fatalf("MigrateFromOldPath() error = %v", err)
		}

		data, err := os.ReadFile(newPath)
		if err != nil {
			t.Fatalf("Failed to read new file: %v", err)
		}
		if string(data) != "new content" {
			t.Error("New file should not be overwritten")
		}
	})

	t.Run("migrates file from old to new location", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldPath := filepath.Join(tmpDir, "old.db")
		newDir := filepath.Join(tmpDir, "subdir")
		newPath := filepath.Join(newDir, "new.db")

		content := []byte("test database content")
		if err := os.WriteFile(oldPath, content, 0644); err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		if err := MigrateFromOldPath(oldPath, newPath); err != nil {
			t.Fatalf("MigrateFromOldPath() error = %v", err)
		}

		data, err := os.ReadFile(newPath)
		if err != nil {
			t.Fatalf("Failed to read migrated file: %v", err)
		}
		if string(data) != string(content) {
			t.Error("Migrated file content mismatch")
		}

		if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
			t.Error("Old file should be removed after migration")
		}
	})

	t.Run("creates new directory if needed", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldPath := filepath.Join(tmpDir, "old.db")
		newPath := filepath.Join(tmpDir, "deeply", "nested", "dir", "new.db")

		content := []byte("test")
		if err := os.WriteFile(oldPath, content, 0644); err != nil {
			t.Fatalf("Failed to create old file: %v", err)
		}

		if err := MigrateFromOldPath(oldPath, newPath); err != nil {
			t.Fatalf("MigrateFromOldPath() error = %v", err)
		}

		if _, err := os.Stat(newPath); err != nil {
			t.Errorf("New file should exist: %v", err)
		}
	})
}
