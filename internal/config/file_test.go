package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile_ValidConfig(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected *FileConfig
	}{
		{
			name: "full config",
			yaml: `database:
  path: ./custom.db
  mode: remote
  sync_interval: 120s
  turso:
    url: libsql://example.turso.io
    token: secret-token
`,
			expected: &FileConfig{
				Database: DatabaseConfig{
					Path:         "./custom.db",
					Mode:         "remote",
					SyncInterval: "120s",
					Turso: TursoConfig{
						URL:   "libsql://example.turso.io",
						Token: "secret-token",
					},
				},
			},
		},
		{
			name: "partial config - only path",
			yaml: `database:
  path: ./mydb.db
`,
			expected: &FileConfig{
				Database: DatabaseConfig{
					Path: "./mydb.db",
				},
			},
		},
		{
			name: "partial config - only turso",
			yaml: `database:
  turso:
    url: libsql://test.turso.io
    token: mytoken
`,
			expected: &FileConfig{
				Database: DatabaseConfig{
					Turso: TursoConfig{
						URL:   "libsql://test.turso.io",
						Token: "mytoken",
					},
				},
			},
		},
		{
			name: "local mode",
			yaml: `database:
  path: ./local.db
  mode: local
`,
			expected: &FileConfig{
				Database: DatabaseConfig{
					Path: "./local.db",
					Mode: "local",
				},
			},
		},
		{
			name: "replica mode with sync",
			yaml: `database:
  mode: replica
  sync_interval: 30s
  turso:
    url: libsql://replica.turso.io
    token: replica-token
`,
			expected: &FileConfig{
				Database: DatabaseConfig{
					Mode:         "replica",
					SyncInterval: "30s",
					Turso: TursoConfig{
						URL:   "libsql://replica.turso.io",
						Token: "replica-token",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "dopadone.yaml")
			if err := os.WriteFile(configPath, []byte(tt.yaml), 0644); err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			cfg, err := ParseFile(configPath)
			if err != nil {
				t.Fatalf("ParseFile() error = %v", err)
			}

			if cfg == nil {
				t.Fatal("ParseFile() returned nil")
			}

			if cfg.Database.Path != tt.expected.Database.Path {
				t.Errorf("Database.Path = %q, want %q", cfg.Database.Path, tt.expected.Database.Path)
			}
			if cfg.Database.Mode != tt.expected.Database.Mode {
				t.Errorf("Database.Mode = %q, want %q", cfg.Database.Mode, tt.expected.Database.Mode)
			}
			if cfg.Database.SyncInterval != tt.expected.Database.SyncInterval {
				t.Errorf("Database.SyncInterval = %q, want %q", cfg.Database.SyncInterval, tt.expected.Database.SyncInterval)
			}
			if cfg.Database.Turso.URL != tt.expected.Database.Turso.URL {
				t.Errorf("Database.Turso.URL = %q, want %q", cfg.Database.Turso.URL, tt.expected.Database.Turso.URL)
			}
			if cfg.Database.Turso.Token != tt.expected.Database.Turso.Token {
				t.Errorf("Database.Turso.Token = %q, want %q", cfg.Database.Turso.Token, tt.expected.Database.Turso.Token)
			}
		})
	}
}

func TestParseFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	if err := os.WriteFile(configPath, []byte{}, 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := ParseFile(configPath)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("ParseFile() returned nil for empty file")
	}
}

func TestParseFile_NonExistent(t *testing.T) {
	cfg, err := ParseFile("/nonexistent/path/dopadone.yaml")
	if err != nil {
		t.Fatalf("ParseFile() error = %v, want nil for non-existent file", err)
	}
	if cfg != nil {
		t.Error("ParseFile() should return nil for non-existent file")
	}
}

func TestParseFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	invalidYAML := `database:
  path: [this is not valid
`
	if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := ParseFile(configPath)
	if err == nil {
		t.Error("ParseFile() should return error for invalid YAML")
	}
	if cfg != nil {
		t.Error("ParseFile() should return nil config on error")
	}
}

func TestParseFile_InvalidMode(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	yaml := `database:
  mode: invalid_mode
`
	if err := os.WriteFile(configPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := ParseFile(configPath)
	if err == nil {
		t.Error("ParseFile() should return error for invalid mode")
	}
	if cfg != nil {
		t.Error("ParseFile() should return nil config on validation error")
	}
}

func TestParseFile_InvalidSyncInterval(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dopadone.yaml")
	yaml := `database:
  sync_interval: not-a-duration
`
	if err := os.WriteFile(configPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := ParseFile(configPath)
	if err == nil {
		t.Error("ParseFile() should return error for invalid sync_interval")
	}
	if cfg != nil {
		t.Error("ParseFile() should return nil config on validation error")
	}
}

func TestFileConfig_SyncIntervalDuration(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *FileConfig
		expected    string
		expectError bool
	}{
		{
			name:        "empty interval",
			cfg:         &FileConfig{Database: DatabaseConfig{SyncInterval: ""}},
			expected:    "0s",
			expectError: false,
		},
		{
			name:        "valid interval",
			cfg:         &FileConfig{Database: DatabaseConfig{SyncInterval: "120s"}},
			expected:    "2m0s",
			expectError: false,
		},
		{
			name:        "valid complex interval",
			cfg:         &FileConfig{Database: DatabaseConfig{SyncInterval: "1h30m"}},
			expected:    "1h30m0s",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dur, err := tt.cfg.SyncIntervalDuration()
			if tt.expectError {
				if err == nil {
					t.Error("SyncIntervalDuration() should return error")
				}
			} else {
				if err != nil {
					t.Fatalf("SyncIntervalDuration() error = %v", err)
				}
				if dur.String() != tt.expected {
					t.Errorf("SyncIntervalDuration() = %v, want %v", dur, tt.expected)
				}
			}
		})
	}
}

func TestDiscoverConfig_Order(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	defer func() {
		_ = os.Chdir(originalWd)
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	cwdConfig := filepath.Join(tmpDir, "dopadone.yaml")
	if err := os.WriteFile(cwdConfig, []byte("database:\n  path: ./cwd.db"), 0644); err != nil {
		t.Fatalf("failed to write cwd config: %v", err)
	}

	path, err := DiscoverConfig()
	if err != nil {
		t.Fatalf("DiscoverConfig() error = %v", err)
	}

	expected := "dopadone.yaml"
	if filepath.Base(path) != expected {
		t.Errorf("DiscoverConfig() = %q, want base to be %q (cwd should take precedence)", path, expected)
	}
}

func TestDiscoverConfig_XDGPath(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	defer func() {
		_ = os.Chdir(originalWd)
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}
	if err := os.Chdir(emptyDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	xdgDir := filepath.Join(tmpDir, "xdg-config")
	xdgConfigDir := filepath.Join(xdgDir, "dopadone")
	if err := os.MkdirAll(xdgConfigDir, 0755); err != nil {
		t.Fatalf("failed to create xdg config dir: %v", err)
	}
	xdgConfigPath := filepath.Join(xdgConfigDir, "config.yaml")
	if err := os.WriteFile(xdgConfigPath, []byte("database:\n  path: ./xdg.db"), 0644); err != nil {
		t.Fatalf("failed to write xdg config: %v", err)
	}

	if err := os.Setenv("XDG_CONFIG_HOME", xdgDir); err != nil {
		t.Fatalf("failed to set XDG_CONFIG_HOME: %v", err)
	}

	path, err := DiscoverConfig()
	if err != nil {
		t.Fatalf("DiscoverConfig() error = %v", err)
	}
	if path != xdgConfigPath {
		t.Errorf("DiscoverConfig() = %q, want %q", path, xdgConfigPath)
	}
}

func TestDiscoverConfig_NoFile(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	defer func() {
		_ = os.Chdir(originalWd)
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}
	if err := os.Chdir(emptyDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	_ = os.Unsetenv("XDG_CONFIG_HOME")

	_, err := DiscoverConfig()
	if err != ErrNoConfigFile {
		t.Errorf("DiscoverConfig() error = %v, want ErrNoConfigFile", err)
	}
}

func TestDiscoverConfigWithExplicit_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "custom-config.yaml")
	if err := os.WriteFile(configPath, []byte("database:\n  path: ./test.db"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	path, err := DiscoverConfigWithExplicit(configPath)
	if err != nil {
		t.Fatalf("DiscoverConfigWithExplicit() error = %v", err)
	}
	if path != configPath {
		t.Errorf("DiscoverConfigWithExplicit() = %q, want %q", path, configPath)
	}
}

func TestDiscoverConfigWithExplicit_NonExistent(t *testing.T) {
	_, err := DiscoverConfigWithExplicit("/nonexistent/config.yaml")
	if err == nil {
		t.Error("DiscoverConfigWithExplicit() should return error for non-existent file")
	}
}

func TestDiscoverConfigWithExplicit_EmptyString(t *testing.T) {
	originalWd, _ := os.Getwd()
	tmpDir := t.TempDir()

	defer func() {
		_ = os.Chdir(originalWd)
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	_, err := DiscoverConfigWithExplicit("")
	if err != ErrNoConfigFile {
		t.Errorf("DiscoverConfigWithExplicit() error = %v, want ErrNoConfigFile", err)
	}
}

func TestGetDiscoveryPaths(t *testing.T) {
	originalWd, _ := os.Getwd()
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	tmpDir := t.TempDir()

	defer func() {
		_ = os.Chdir(originalWd)
		if originalXDG == "" {
			_ = os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			_ = os.Setenv("XDG_CONFIG_HOME", originalXDG)
		}
	}()

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}
	_ = os.Unsetenv("XDG_CONFIG_HOME")

	paths := getDiscoveryPaths()

	if len(paths) < 2 {
		t.Errorf("getDiscoveryPaths() returned only %d paths, expected at least 2", len(paths))
	}

	found := make(map[string]bool)
	for _, p := range paths {
		found[filepath.Base(p)] = true
	}

	if !found["dopadone.yaml"] {
		t.Error("getDiscoveryPaths() should include dopadone.yaml in cwd")
	}
	if !found["config.yaml"] {
		t.Error("getDiscoveryPaths() should include config.yaml in XDG path")
	}
	if !found[".dopadone.yaml"] {
		t.Error("getDiscoveryPaths() should include .dopadone.yaml in home dir")
	}
}
