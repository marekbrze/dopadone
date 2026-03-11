package driver

import (
	"strings"
	"testing"
)

func TestDetectMode(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected DriverType
		reason   string
	}{
		{
			name: "local_sqlite_default",
			config: &Config{
				DatabasePath: "/tmp/test.db",
			},
			expected: DriverSQLite,
			reason:   "local SQLite",
		},
		{
			name: "remote_turso",
			config: &Config{
				TursoURL:   "libsql://test.turso.io",
				TursoToken: "test-token",
			},
			expected: DriverTursoRemote,
			reason:   "remote Turso",
		},
		{
			name: "embedded_replica",
			config: &Config{
				DatabasePath: "/tmp/test.db",
				TursoURL:     "libsql://test.turso.io",
				TursoToken:   "test-token",
			},
			expected: DriverTursoReplica,
			reason:   "embedded replica",
		},
		{
			name:     "empty_config",
			config:   &Config{},
			expected: DriverSQLite,
			reason:   "default fallback",
		},
		{
			name: "turso_url_only_no_token",
			config: &Config{
				TursoURL: "libsql://test.turso.io",
			},
			expected: DriverSQLite,
			reason:   "local SQLite",
		},
		{
			name: "turso_token_only_no_url",
			config: &Config{
				TursoToken: "test-token",
			},
			expected: DriverSQLite,
			reason:   "local SQLite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectMode(tt.config)
			if result.Type != tt.expected {
				t.Errorf("DetectMode() = %v, want %v", result.Type, tt.expected)
			}
			if !strings.Contains(result.Reason, tt.reason) {
				t.Errorf("Reason = %v, want to contain %v", result.Reason, tt.reason)
			}
		})
	}
}

func TestParseExplicitMode(t *testing.T) {
	tests := []struct {
		input    string
		expected DriverType
		wantErr  bool
	}{
		{"", "", false},
		{"auto", "", false},
		{"local", DriverSQLite, false},
		{"sqlite", DriverSQLite, false},
		{"remote", DriverTursoRemote, false},
		{"turso-remote", DriverTursoRemote, false},
		{"replica", DriverTursoReplica, false},
		{"turso-replica", DriverTursoReplica, false},
		{"invalid", "", true},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseExplicitMode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseExplicitMode(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ParseExplicitMode(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDetectOrExplicitMode(t *testing.T) {
	tests := []struct {
		name       string
		config     *Config
		expected   DriverType
		wantErr    bool
		errContain string
	}{
		{
			name:     "auto_detect_local",
			config:   &Config{DatabasePath: "/tmp/test.db", Type: ""},
			expected: DriverSQLite,
			wantErr:  false,
		},
		{
			name:     "auto_detect_remote",
			config:   &Config{TursoURL: "libsql://test.turso.io", TursoToken: "token", Type: ""},
			expected: DriverTursoRemote,
			wantErr:  false,
		},
		{
			name:     "auto_detect_replica",
			config:   &Config{DatabasePath: "/tmp/test.db", TursoURL: "libsql://test.turso.io", TursoToken: "token", Type: ""},
			expected: DriverTursoReplica,
			wantErr:  false,
		},
		{
			name:     "explicit_local",
			config:   &Config{Type: "local"},
			expected: DriverSQLite,
			wantErr:  false,
		},
		{
			name:     "explicit_remote",
			config:   &Config{Type: "remote"},
			expected: DriverTursoRemote,
			wantErr:  false,
		},
		{
			name:     "explicit_replica",
			config:   &Config{Type: "replica"},
			expected: DriverTursoReplica,
			wantErr:  false,
		},
		{
			name:       "invalid_mode",
			config:     &Config{Type: "invalid"},
			expected:   "",
			wantErr:    true,
			errContain: "invalid database mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DetectOrExplicitMode(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectOrExplicitMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errContain != "" {
				if !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("Error = %v, want to contain %v", err, tt.errContain)
				}
				return
			}
			if result.Type != tt.expected {
				t.Errorf("DetectOrExplicitMode() = %v, want %v", result.Type, tt.expected)
			}
		})
	}
}

func TestValidateConfigForMode(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		drvType DriverType
		wantErr bool
	}{
		{
			name:    "valid_sqlite",
			config:  &Config{DatabasePath: "/tmp/test.db"},
			drvType: DriverSQLite,
			wantErr: false,
		},
		{
			name:    "invalid_sqlite_missing_path",
			config:  &Config{DatabasePath: ""},
			drvType: DriverSQLite,
			wantErr: true,
		},
		{
			name:    "valid_remote",
			config:  &Config{TursoURL: "libsql://test.turso.io", TursoToken: "token"},
			drvType: DriverTursoRemote,
			wantErr: false,
		},
		{
			name:    "invalid_remote_missing_url",
			config:  &Config{TursoURL: "", TursoToken: "token"},
			drvType: DriverTursoRemote,
			wantErr: true,
		},
		{
			name:    "invalid_remote_missing_token",
			config:  &Config{TursoURL: "libsql://test.turso.io", TursoToken: ""},
			drvType: DriverTursoRemote,
			wantErr: true,
		},
		{
			name:    "valid_replica",
			config:  &Config{DatabasePath: "/tmp/test.db", TursoURL: "libsql://test.turso.io", TursoToken: "token"},
			drvType: DriverTursoReplica,
			wantErr: false,
		},
		{
			name:    "invalid_replica_missing_path",
			config:  &Config{TursoURL: "libsql://test.turso.io", TursoToken: "token"},
			drvType: DriverTursoReplica,
			wantErr: true,
		},
		{
			name:    "invalid_replica_missing_token",
			config:  &Config{DatabasePath: "/tmp/test.db", TursoURL: "libsql://test.turso.io"},
			drvType: DriverTursoReplica,
			wantErr: true,
		},
		{
			name:    "unknown_driver_type",
			config:  &Config{},
			drvType: DriverType("unknown"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfigForMode(tt.config, tt.drvType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfigForMode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
