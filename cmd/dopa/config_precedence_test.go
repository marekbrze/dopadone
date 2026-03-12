package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig_CLI_Overrides_File(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	configContent := `database:
  path: ./file-path.db
  mode: remote
  sync_interval: 120s
  turso:
    url: libsql://file.turso.io
    token: file-token
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "./cli-path.db",
		TursoURL:     "libsql://cli.turso.io",
		TursoToken:   "cli-token",
		DBMode:       "local",
		SyncInterval: 30 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./cli-path.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "./cli-path.db")
	}
	if cfg.TursoURL != "libsql://cli.turso.io" {
		t.Errorf("TursoURL = %q, want %q", cfg.TursoURL, "libsql://cli.turso.io")
	}
	if cfg.TursoToken != "cli-token" {
		t.Errorf("TursoToken = %q, want %q", cfg.TursoToken, "cli-token")
	}
	if cfg.DBMode != "local" {
		t.Errorf("DBMode = %q, want %q", cfg.DBMode, "local")
	}
	if cfg.SyncInterval != 30*time.Second {
		t.Errorf("SyncInterval = %v, want %v", cfg.SyncInterval, 30*time.Second)
	}
}

func TestLoadConfig_File_Overrides_Env(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	configContent := `database:
  path: ./file-path.db
  mode: remote
  sync_interval: 120s
  turso:
    url: libsql://file.turso.io
    token: file-token
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	os.Setenv("DOPA_DB_PATH", "./env-path.db")
	os.Setenv("TURSO_DATABASE_URL", "libsql://env.turso.io")
	os.Setenv("TURSO_AUTH_TOKEN", "env-token")
	os.Setenv("DOPA_DB_MODE", "local")
	defer func() {
		_ = os.Unsetenv("DOPA_DB_PATH")
		_ = os.Unsetenv("TURSO_DATABASE_URL")
		_ = os.Unsetenv("TURSO_AUTH_TOKEN")
		_ = os.Unsetenv("DOPA_DB_MODE")
	}()

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./env-path.db" {
		t.Errorf("DatabasePath = %q, want %q (env should win over file)", cfg.DatabasePath, "./env-path.db")
	}
	if cfg.TursoURL != "libsql://env.turso.io" {
		t.Errorf("TursoURL = %q, want %q (env should win over file)", cfg.TursoURL, "libsql://env.turso.io")
	}
	if cfg.TursoToken != "env-token" {
		t.Errorf("TursoToken = %q, want %q (env should win over file)", cfg.TursoToken, "env-token")
	}
	if cfg.DBMode != "local" {
		t.Errorf("DBMode = %q, want %q (env should win over file)", cfg.DBMode, "local")
	}
}

func TestLoadConfig_CLI_Overrides_Env(t *testing.T) {
	os.Setenv("DOPA_DB_PATH", "./env-path.db")
	os.Setenv("TURSO_DATABASE_URL", "libsql://env.turso.io")
	os.Setenv("TURSO_AUTH_TOKEN", "env-token")
	os.Setenv("DOPA_DB_MODE", "local")
	defer func() {
		_ = os.Unsetenv("DOPA_DB_PATH")
		_ = os.Unsetenv("TURSO_DATABASE_URL")
		_ = os.Unsetenv("TURSO_AUTH_TOKEN")
		_ = os.Unsetenv("DOPA_DB_MODE")
	}()

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "./cli-path.db",
		TursoURL:     "libsql://cli.turso.io",
		TursoToken:   "cli-token",
		DBMode:       "remote",
		SyncInterval: 30 * time.Second,
		ConfigPath:   "",
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./cli-path.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "./cli-path.db")
	}
	if cfg.TursoURL != "libsql://cli.turso.io" {
		t.Errorf("TursoURL = %q, want %q", cfg.TursoURL, "libsql://cli.turso.io")
	}
	if cfg.TursoToken != "cli-token" {
		t.Errorf("TursoToken = %q, want %q", cfg.TursoToken, "cli-token")
	}
	if cfg.DBMode != "remote" {
		t.Errorf("DBMode = %q, want %q", cfg.DBMode, "remote")
	}
}

func TestLoadConfig_PartialMerge(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	configContent := `database:
  path: ./file-path.db
  mode: replica
  turso:
    url: libsql://file.turso.io
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	os.Setenv("TURSO_AUTH_TOKEN", "env-token")
	defer func() { _ = os.Unsetenv("TURSO_AUTH_TOKEN") }()

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./file-path.db" {
		t.Errorf("DatabasePath = %q, want %q (from file)", cfg.DatabasePath, "./file-path.db")
	}
	if cfg.TursoURL != "libsql://file.turso.io" {
		t.Errorf("TursoURL = %q, want %q (from file)", cfg.TursoURL, "libsql://file.turso.io")
	}
	if cfg.TursoToken != "env-token" {
		t.Errorf("TursoToken = %q, want %q (from env)", cfg.TursoToken, "env-token")
	}
	if cfg.DBMode != "replica" {
		t.Errorf("DBMode = %q, want %q (from file)", cfg.DBMode, "replica")
	}
}

func TestLoadConfig_ExplicitConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-config.yaml")
	configContent := `database:
  path: ./custom-path.db
  mode: remote
  turso:
    url: libsql://custom.turso.io
    token: custom-token
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./custom-path.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "./custom-path.db")
	}
	if cfg.TursoURL != "libsql://custom.turso.io" {
		t.Errorf("TursoURL = %q, want %q", cfg.TursoURL, "libsql://custom.turso.io")
	}
	if cfg.TursoToken != "custom-token" {
		t.Errorf("TursoToken = %q, want %q", cfg.TursoToken, "custom-token")
	}
	if cfg.DBMode != "remote" {
		t.Errorf("DBMode = %q, want %q", cfg.DBMode, "remote")
	}
}

func TestLoadConfig_NoConfigFile(t *testing.T) {
	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "./test.db",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   "",
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./test.db" {
		t.Errorf("DatabasePath = %q, want %q", cfg.DatabasePath, "./test.db")
	}
}

func TestLoadConfig_SyncIntervalFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	configContent := `database:
  sync_interval: 90s
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.SyncInterval != 90*time.Second {
		t.Errorf("SyncInterval = %v, want %v", cfg.SyncInterval, 90*time.Second)
	}
}

func TestLoadConfig_InvalidConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	invalidContent := `database:
  mode: [invalid
`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err == nil {
		t.Error("LoadConfig() should return error for invalid config file")
	}
}

func TestLoadConfig_DefaultSyncInterval(t *testing.T) {
	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   "",
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.SyncInterval != 60*time.Second {
		t.Errorf("SyncInterval = %v, want %v", cfg.SyncInterval, 60*time.Second)
	}
}

func TestPrecedence_CLI_Over_File_Over_Env(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	configContent := `database:
  path: ./file-path.db
  turso:
    url: libsql://file.turso.io
    token: file-token
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	os.Setenv("DOPA_DB_PATH", "./env-path.db")
	os.Setenv("TURSO_DATABASE_URL", "libsql://env.turso.io")
	os.Setenv("TURSO_AUTH_TOKEN", "env-token")
	defer func() {
		_ = os.Unsetenv("DOPA_DB_PATH")
		_ = os.Unsetenv("TURSO_DATABASE_URL")
		_ = os.Unsetenv("TURSO_AUTH_TOKEN")
	}()

	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "./cli-path.db",
		TursoURL:     "",
		TursoToken:   "",
		DBMode:       "",
		SyncInterval: 60 * time.Second,
		ConfigPath:   configPath,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./cli-path.db" {
		t.Errorf("DatabasePath = %q, want %q (CLI wins)", cfg.DatabasePath, "./cli-path.db")
	}
	if cfg.TursoURL != "libsql://env.turso.io" {
		t.Errorf("TursoURL = %q, want %q (env wins, CLI empty)", cfg.TursoURL, "libsql://env.turso.io")
	}
	if cfg.TursoToken != "env-token" {
		t.Errorf("TursoToken = %q, want %q (env wins, CLI empty)", cfg.TursoToken, "env-token")
	}
}
