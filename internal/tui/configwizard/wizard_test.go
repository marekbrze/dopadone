package configwizard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/config"
)

func TestNewWizard(t *testing.T) {
	w := New()
	if w.step != stepWelcome {
		t.Errorf("Expected initial step to be stepWelcome, got %d", w.step)
	}
	if w.mode != ModeLocal {
		t.Errorf("Expected initial mode to be ModeLocal, got %s", w.mode)
	}
	if w.selectedModeIndex != 0 {
		t.Errorf("Expected initial selectedModeIndex to be 0, got %d", w.selectedModeIndex)
	}
}

func TestNewWizard_DefaultPath(t *testing.T) {
	defaultPath := cli.DefaultDBPath()
	if defaultPath == "" {
		t.Error("Default path should not be empty")
	}

	w := New()
	if w.localPath.Value() != defaultPath {
		t.Errorf("Expected local path to be %s, got %s", defaultPath, w.localPath.Value())
	}
}

func TestWizardModeSelection_NavigateDown(t *testing.T) {
	w := New()
	w.step = stepModeSelection

	model, _ := w.Update(tea.KeyMsg{Type: tea.KeyDown})
	updated := model.(*Wizard)
	if updated.selectedModeIndex != 1 {
		t.Errorf("Expected selectedModeIndex to be 1 after down, got %d", updated.selectedModeIndex)
	}
}

func TestWizardModeSelection_NavigateUp(t *testing.T) {
	w := New()
	w.step = stepModeSelection
	w.selectedModeIndex = 0

	model, _ := w.Update(tea.KeyMsg{Type: tea.KeyUp})
	updated := model.(*Wizard)
	if updated.selectedModeIndex != 2 {
		t.Errorf("Expected selectedModeIndex to wrap to 2 after up from 0, got %d", updated.selectedModeIndex)
	}
}

func TestWizardModeSelection_SelectMode(t *testing.T) {
	tests := []struct {
		name          string
		selectedIndex int
		expectedMode  DatabaseMode
	}{
		{"select local", 0, ModeLocal},
		{"select remote", 1, ModeRemote},
		{"select replica", 2, ModeReplica},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.step = stepModeSelection
			w.selectedModeIndex = tt.selectedIndex

			model, _ := w.Update(tea.KeyMsg{Type: tea.KeyEnter})
			updated := model.(*Wizard)
			if updated.mode != tt.expectedMode {
				t.Errorf("Expected mode %s, got %s", tt.expectedMode, updated.mode)
			}
			if updated.step != stepConfig {
				t.Errorf("Expected step to be stepConfig, got %d", updated.step)
			}
		})
	}
}

func TestWizardCancelOnWelcome(t *testing.T) {
	w := New()
	_, cmd := w.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Error("Expected cmd to be non-nil on ctrl+c at welcome step")
	}
}

func TestWizardBackNavigation(t *testing.T) {
	w := New()
	w.step = stepModeSelection

	model, _ := w.Update(tea.KeyMsg{Type: tea.KeyEsc})
	updated := model.(*Wizard)
	if updated.step != stepWelcome {
		t.Errorf("Expected step to be stepWelcome after ESC, got %d", updated.step)
	}
}

func TestValidateConfig_LocalMode(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"empty path", "", true},
		{"valid path", "/path/to/db.sqlite", false},
		{"whitespace only", "   ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.mode = ModeLocal
			w.localPath.SetValue(tt.path)
			err := w.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfig_RemoteMode(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		token   string
		wantErr bool
	}{
		{"empty url", "", "token", true},
		{"empty token", "libsql://test.turso.io", "", true},
		{"both empty", "", "", true},
		{"both valid", "libsql://test.turso.io", "token", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.mode = ModeRemote
			w.tursoURL.SetValue(tt.url)
			w.tursoToken.SetValue(tt.token)
			err := w.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfig_ReplicaMode(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		url     string
		token   string
		wantErr bool
	}{
		{"empty path", "", "libsql://test.turso.io", "token", true},
		{"empty url", "/path/to/replica.db", "", "token", true},
		{"empty token", "/path/to/replica.db", "libsql://test.turso.io", "", true},
		{"all valid", "/path/to/replica.db", "libsql://test.turso.io", "token", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := New()
			w.mode = ModeReplica
			w.replicaPath.SetValue(tt.path)
			w.tursoURL.SetValue(tt.url)
			w.tursoToken.SetValue(tt.token)
			err := w.validateConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigResultConversion_LocalMode(t *testing.T) {
	result := &ConfigResult{
		Mode:         ModeLocal,
		DatabasePath: "/path/to/db.sqlite",
	}
	cfg := result.ToFileConfig()
	if cfg.Database.Path != "/path/to/db.sqlite" {
		t.Errorf("Expected path /path/to/db.sqlite, got %s", cfg.Database.Path)
	}
	if cfg.Database.Mode != "local" {
		t.Errorf("Expected mode local, got %s", cfg.Database.Mode)
	}
}

func TestConfigResultConversion_RemoteMode(t *testing.T) {
	result := &ConfigResult{
		Mode:       ModeRemote,
		TursoURL:   "libsql://test.turso.io",
		TursoToken: "secret-token",
	}
	cfg := result.ToFileConfig()
	if cfg.Database.Mode != "remote" {
		t.Errorf("Expected mode remote, got %s", cfg.Database.Mode)
	}
	if cfg.Database.Turso.URL != "libsql://test.turso.io" {
		t.Errorf("Expected Turso URL libsql://test.turso.io, got %s", cfg.Database.Turso.URL)
	}
	if cfg.Database.Turso.Token != "secret-token" {
		t.Errorf("Expected Turso token secret-token, got %s", cfg.Database.Turso.Token)
	}
}

func TestConfigResultConversion_ReplicaMode(t *testing.T) {
	result := &ConfigResult{
		Mode:         ModeReplica,
		DatabasePath: "/path/to/replica.db",
		TursoURL:     "libsql://test.turso.io",
		TursoToken:   "secret-token",
	}
	cfg := result.ToFileConfig()
	if cfg.Database.Path != "/path/to/replica.db" {
		t.Errorf("Expected path /path/to/replica.db, got %s", cfg.Database.Path)
	}
	if cfg.Database.Mode != "replica" {
		t.Errorf("Expected mode replica, got %s", cfg.Database.Mode)
	}
	if cfg.Database.Turso.URL != "libsql://test.turso.io" {
		t.Errorf("Expected Turso URL libsql://test.turso.io, got %s", cfg.Database.Turso.URL)
	}
	if cfg.Database.Turso.Token != "secret-token" {
		t.Errorf("Expected Turso token secret-token, got %s", cfg.Database.Turso.Token)
	}
}

func TestGetConfigResult(t *testing.T) {
	t.Run("local mode", func(t *testing.T) {
		w := New()
		w.mode = ModeLocal
		w.localPath.SetValue("/custom/path.db")
		result := w.getConfigResult()
		if result.Mode != ModeLocal {
			t.Errorf("Expected mode ModeLocal, got %s", result.Mode)
		}
		if result.DatabasePath != "/custom/path.db" {
			t.Errorf("Expected database path /custom/path.db, got %s", result.DatabasePath)
		}
	})

	t.Run("remote mode", func(t *testing.T) {
		w := New()
		w.mode = ModeRemote
		w.tursoURL.SetValue("libsql://remote.turso.io")
		w.tursoToken.SetValue("my-token")
		result := w.getConfigResult()
		if result.Mode != ModeRemote {
			t.Errorf("Expected mode ModeRemote, got %s", result.Mode)
		}
		if result.TursoURL != "libsql://remote.turso.io" {
			t.Errorf("Expected Turso URL libsql://remote.turso.io, got %s", result.TursoURL)
		}
		if result.TursoToken != "my-token" {
			t.Errorf("Expected Turso token my-token, got %s", result.TursoToken)
		}
	})

	t.Run("replica mode", func(t *testing.T) {
		w := New()
		w.mode = ModeReplica
		w.replicaPath.SetValue("/replica/path.db")
		w.tursoURL.SetValue("libsql://replica.turso.io")
		w.tursoToken.SetValue("replica-token")
		result := w.getConfigResult()
		if result.Mode != ModeReplica {
			t.Errorf("Expected mode ModeReplica, got %s", result.Mode)
		}
		if result.DatabasePath != "/replica/path.db" {
			t.Errorf("Expected database path /replica/path.db, got %s", result.DatabasePath)
		}
		if result.TursoURL != "libsql://replica.turso.io" {
			t.Errorf("Expected Turso URL libsql://replica.turso.io, got %s", result.TursoURL)
		}
		if result.TursoToken != "replica-token" {
			t.Errorf("Expected Turso token replica-token, got %s", result.TursoToken)
		}
	})
}

func TestIsComplete(t *testing.T) {
	w := New()
	if w.IsComplete() {
		t.Error("Wizard should not be complete initially")
	}
	w.step = stepSuccess
	if !w.IsComplete() {
		t.Error("Wizard should be complete at stepSuccess")
	}
}

func TestSetError(t *testing.T) {
	w := New()
	w.SetError("test error")
	if w.errorMsg != "test error" {
		t.Errorf("Expected error message 'test error', got %s", w.errorMsg)
	}
	w.ClearError()
	if w.errorMsg != "" {
		t.Errorf("Expected error message to be cleared, got %s", w.errorMsg)
	}
}

func TestFirstRunDetection(t *testing.T) {
	// Skip test if a config file already exists in the system
	// This test verifies the function works but can't guarantee isolation
	// in environments where a config file may already be present
	_, err := config.DiscoverConfig()
	if err == nil {
		t.Skip("Skipping: config file already exists in system")
	}

	if !config.IsFirstRun() {
		t.Error("IsFirstRun should return true when no config exists")
	}
}

func TestGetDefaultConfigPath(t *testing.T) {
	path, err := config.GetDefaultConfigPath()
	if err != nil {
		t.Errorf("GetDefaultConfigPath should not error: %v", err)
	}
	if path == "" {
		t.Error("GetDefaultConfigPath should return non-empty path")
	}
}
