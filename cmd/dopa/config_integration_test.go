package main

import (
	"os"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "./test.db",
		SyncInterval: 60 * time.Second,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "./test.db" {
		t.Errorf("DatabasePath = %v, want ./test.db", cfg.DatabasePath)
	}
	if cfg.TursoURL != "" {
		t.Errorf("TursoURL = %v, want empty", cfg.TursoURL)
	}
	if cfg.TursoToken != "" {
		t.Errorf("TursoToken = %v, want empty", cfg.TursoToken)
	}
	if cfg.DBMode != "" {
		t.Errorf("DBMode = %v, want empty", cfg.DBMode)
	}
	if cfg.SyncInterval != 60*time.Second {
		t.Errorf("SyncInterval = %v, want 60s", cfg.SyncInterval)
	}
}

func TestLoadConfig_AllFlagsProvided(t *testing.T) {
	cfg, err := LoadConfig(LoadConfigParams{
		DBPath:       "/custom/path.db",
		TursoURL:     "libsql://test.turso.io",
		TursoToken:   "test-token",
		DBMode:       "replica",
		SyncInterval: 30 * time.Second,
	})
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DatabasePath != "/custom/path.db" {
		t.Errorf("DatabasePath = %v, want /custom/path.db", cfg.DatabasePath)
	}
	if cfg.TursoURL != "libsql://test.turso.io" {
		t.Errorf("TursoURL = %v, want libsql://test.turso.io", cfg.TursoURL)
	}
	if cfg.TursoToken != "test-token" {
		t.Errorf("TursoToken = %v, want test-token", cfg.TursoToken)
	}
	if cfg.DBMode != "replica" {
		t.Errorf("DBMode = %v, want replica", cfg.DBMode)
	}
	if cfg.SyncInterval != 30*time.Second {
		t.Errorf("SyncInterval = %v, want 30s", cfg.SyncInterval)
	}
}

func TestResolveDBPath_FlagPrecedence(t *testing.T) {
	if err := os.Setenv("DOPA_DB_PATH", "/env/path.db"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("DOPA_DB_PATH") }()

	path := resolveDBPath("/flag/path.db", nil)

	if path != "/flag/path.db" {
		t.Errorf("resolveDBPath = %v, want /flag/path.db (flag should override env)", path)
	}
}

func TestResolveDBPath_EnvFallback(t *testing.T) {
	if err := os.Setenv("DOPA_DB_PATH", "/env/path.db"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("DOPA_DB_PATH") }()

	path := resolveDBPath("./dopadone.db", nil)

	if path != "/env/path.db" {
		t.Errorf("resolveDBPath = %v, want /env/path.db (env should be used when flag is default)", path)
	}
}

func TestResolveDBPath_Default(t *testing.T) {
	if err := os.Unsetenv("DOPA_DB_PATH"); err != nil {
		t.Fatalf("Unsetenv failed: %v", err)
	}

	path := resolveDBPath("./dopadone.db", nil)

	if path != "./dopadone.db" {
		t.Errorf("resolveDBPath = %v, want ./dopadone.db (default)", path)
	}
}

func TestResolveTursoURL_FlagPrecedence(t *testing.T) {
	if err := os.Setenv("TURSO_DATABASE_URL", "libsql://env.turso.io"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("TURSO_DATABASE_URL") }()

	url := resolveTursoURL("libsql://flag.turso.io", nil)

	if url != "libsql://flag.turso.io" {
		t.Errorf("resolveTursoURL = %v, want libsql://flag.turso.io (flag should override env)", url)
	}
}

func TestResolveTursoURL_EnvFallback(t *testing.T) {
	if err := os.Setenv("TURSO_DATABASE_URL", "libsql://env.turso.io"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("TURSO_DATABASE_URL") }()

	url := resolveTursoURL("", nil)

	if url != "libsql://env.turso.io" {
		t.Errorf("resolveTursoURL = %v, want libsql://env.turso.io (from env)", url)
	}
}

func TestResolveTursoToken_FlagPrecedence(t *testing.T) {
	if err := os.Setenv("TURSO_AUTH_TOKEN", "env-token"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("TURSO_AUTH_TOKEN") }()

	token := resolveTursoToken("flag-token", nil)

	if token != "flag-token" {
		t.Errorf("resolveTursoToken = %v, want flag-token (flag should override env)", token)
	}
}

func TestResolveTursoToken_EnvFallback(t *testing.T) {
	if err := os.Setenv("TURSO_AUTH_TOKEN", "env-token"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("TURSO_AUTH_TOKEN") }()

	token := resolveTursoToken("", nil)

	if token != "env-token" {
		t.Errorf("resolveTursoToken = %v, want env-token (from env)", token)
	}
}

func TestResolveDBMode_FlagPrecedence(t *testing.T) {
	if err := os.Setenv("DOPA_DB_MODE", "remote"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("DOPA_DB_MODE") }()

	mode := resolveDBMode("replica", nil)

	if mode != "replica" {
		t.Errorf("resolveDBMode = %v, want replica (flag should override env)", mode)
	}
}

func TestResolveDBMode_EnvFallback(t *testing.T) {
	if err := os.Setenv("DOPA_DB_MODE", "remote"); err != nil {
		t.Fatalf("Setenv failed: %v", err)
	}
	defer func() { _ = os.Unsetenv("DOPA_DB_MODE") }()

	mode := resolveDBMode("", nil)

	if mode != "remote" {
		t.Errorf("resolveDBMode = %v, want remote (from env)", mode)
	}
}

func TestConfig_ToDriverConfig(t *testing.T) {
	cfg := &Config{
		DatabasePath: "/path/to/db.db",
		TursoURL:     "libsql://test.turso.io",
		TursoToken:   "test-token",
		DBMode:       "replica",
		SyncInterval: 45 * time.Second,
	}

	driverCfg := cfg.ToDriverConfig()

	if driverCfg.DatabasePath != "/path/to/db.db" {
		t.Errorf("DatabasePath = %v, want /path/to/db.db", driverCfg.DatabasePath)
	}
	if driverCfg.TursoURL != "libsql://test.turso.io" {
		t.Errorf("TursoURL = %v, want libsql://test.turso.io", driverCfg.TursoURL)
	}
	if driverCfg.TursoToken != "test-token" {
		t.Errorf("TursoToken = %v, want test-token", driverCfg.TursoToken)
	}
	if driverCfg.Type != driver.DriverType("replica") {
		t.Errorf("Type = %v, want replica", driverCfg.Type)
	}
	if driverCfg.SyncInterval != 45*time.Second {
		t.Errorf("SyncInterval = %v, want 45s", driverCfg.SyncInterval)
	}
}

func TestConfigPrecedence_FullIntegration(t *testing.T) {
	tests := []struct {
		name           string
		flagDBPath     string
		flagTursoURL   string
		flagTursoToken string
		flagDBMode     string
		envDBPath      string
		envTursoURL    string
		envTursoToken  string
		envDBMode      string
		wantDBPath     string
		wantTursoURL   string
		wantTursoToken string
		wantDBMode     string
	}{
		{
			name:           "flags_override_all_env",
			flagDBPath:     "/flag/db.db",
			flagTursoURL:   "libsql://flag.turso.io",
			flagTursoToken: "flag-token",
			flagDBMode:     "replica",
			envDBPath:      "/env/db.db",
			envTursoURL:    "libsql://env.turso.io",
			envTursoToken:  "env-token",
			envDBMode:      "remote",
			wantDBPath:     "/flag/db.db",
			wantTursoURL:   "libsql://flag.turso.io",
			wantTursoToken: "flag-token",
			wantDBMode:     "replica",
		},
		{
			name:           "env_used_when_flag_empty",
			flagDBPath:     "./dopadone.db",
			flagTursoURL:   "",
			flagTursoToken: "",
			flagDBMode:     "",
			envDBPath:      "/env/db.db",
			envTursoURL:    "libsql://env.turso.io",
			envTursoToken:  "env-token",
			envDBMode:      "remote",
			wantDBPath:     "/env/db.db",
			wantTursoURL:   "libsql://env.turso.io",
			wantTursoToken: "env-token",
			wantDBMode:     "remote",
		},
		{
			name:           "defaults_when_nothing_set",
			flagDBPath:     "./dopadone.db",
			flagTursoURL:   "",
			flagTursoToken: "",
			flagDBMode:     "",
			envDBPath:      "",
			envTursoURL:    "",
			envTursoToken:  "",
			envDBMode:      "",
			wantDBPath:     "./dopadone.db",
			wantTursoURL:   "",
			wantTursoToken: "",
			wantDBMode:     "",
		},
		{
			name:           "partial_override_turso_only",
			flagDBPath:     "./dopadone.db",
			flagTursoURL:   "libsql://flag.turso.io",
			flagTursoToken: "",
			flagDBMode:     "",
			envDBPath:      "",
			envTursoURL:    "",
			envTursoToken:  "env-token",
			envDBMode:      "",
			wantDBPath:     "./dopadone.db",
			wantTursoURL:   "libsql://flag.turso.io",
			wantTursoToken: "env-token",
			wantDBMode:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setOrUnsetEnv := func(key, value string) {
				t.Helper()
				if value != "" {
					if err := os.Setenv(key, value); err != nil {
						t.Fatalf("Setenv(%s) failed: %v", key, err)
					}
				} else {
					_ = os.Unsetenv(key)
				}
			}

			setOrUnsetEnv("DOPA_DB_PATH", tt.envDBPath)
			setOrUnsetEnv("TURSO_DATABASE_URL", tt.envTursoURL)
			setOrUnsetEnv("TURSO_AUTH_TOKEN", tt.envTursoToken)
			setOrUnsetEnv("DOPA_DB_MODE", tt.envDBMode)

			defer func() {
				_ = os.Unsetenv("DOPA_DB_PATH")
				_ = os.Unsetenv("TURSO_DATABASE_URL")
				_ = os.Unsetenv("TURSO_AUTH_TOKEN")
				_ = os.Unsetenv("DOPA_DB_MODE")
			}()

			cfg, err := LoadConfig(LoadConfigParams{
				DBPath:       tt.flagDBPath,
				TursoURL:     tt.flagTursoURL,
				TursoToken:   tt.flagTursoToken,
				DBMode:       tt.flagDBMode,
				SyncInterval: 60 * time.Second,
			})
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}

			if cfg.DatabasePath != tt.wantDBPath {
				t.Errorf("DatabasePath = %v, want %v", cfg.DatabasePath, tt.wantDBPath)
			}
			if cfg.TursoURL != tt.wantTursoURL {
				t.Errorf("TursoURL = %v, want %v", cfg.TursoURL, tt.wantTursoURL)
			}
			if cfg.TursoToken != tt.wantTursoToken {
				t.Errorf("TursoToken = %v, want %v", cfg.TursoToken, tt.wantTursoToken)
			}
			if cfg.DBMode != tt.wantDBMode {
				t.Errorf("DBMode = %v, want %v", cfg.DBMode, tt.wantDBMode)
			}
		})
	}
}
