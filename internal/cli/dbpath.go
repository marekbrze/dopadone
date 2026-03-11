package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	AppName = "dopadone"
	DBName  = "dopadone.db"
)

func DefaultDBPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "./" + DBName
	}
	return filepath.Join(configDir, AppName, DBName)
}

func GetDBPathWithFallback() (path string, fallback bool) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "./" + DBName, true
	}
	return filepath.Join(configDir, AppName, DBName), false
}

func EnsureDirExists(dbPath string) error {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return WrapError(err, "failed to resolve database path")
	}

	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return WrapError(err, fmt.Sprintf("failed to create database directory: %s", dir))
	}

	return nil
}

func MigrateFromOldPath(oldPath, newPath string) error {
	if _, err := os.Stat(newPath); err == nil {
		return nil
	}

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil
	}

	if err := EnsureDirExists(newPath); err != nil {
		return WrapError(err, "failed to ensure new database directory exists")
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return WrapError(err, fmt.Sprintf("failed to migrate database from %s to %s", oldPath, newPath))
	}

	return nil
}
