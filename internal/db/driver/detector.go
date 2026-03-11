package driver

import (
	"fmt"
)

type DetectionResult struct {
	Type   DriverType
	Reason string
}

func DetectMode(cfg *Config) DetectionResult {
	hasTurso := cfg.TursoURL != "" && cfg.TursoToken != ""
	hasLocalPath := cfg.DatabasePath != ""

	if hasTurso && hasLocalPath {
		return DetectionResult{
			Type:   DriverTursoReplica,
			Reason: "embedded replica (Turso URL + local path configured)",
		}
	}

	if hasTurso && !hasLocalPath {
		return DetectionResult{
			Type:   DriverTursoRemote,
			Reason: "remote Turso (Turso URL configured without local path)",
		}
	}

	if !hasTurso && hasLocalPath {
		return DetectionResult{
			Type:   DriverSQLite,
			Reason: "local SQLite (no Turso configuration found)",
		}
	}

	return DetectionResult{
		Type:   DriverSQLite,
		Reason: "local SQLite (default fallback)",
	}
}

func ParseExplicitMode(mode string) (DriverType, error) {
	switch mode {
	case "", "auto":
		return "", nil
	case "local", "sqlite":
		return DriverSQLite, nil
	case "remote", "turso-remote":
		return DriverTursoRemote, nil
	case "replica", "turso-replica":
		return DriverTursoReplica, nil
	default:
		return "", fmt.Errorf("invalid database mode: %s (valid: local, remote, replica, auto)", mode)
	}
}

func DetectOrExplicitMode(cfg *Config) (DetectionResult, error) {
	explicitMode, err := ParseExplicitMode(string(cfg.Type))
	if err != nil {
		return DetectionResult{}, err
	}

	if explicitMode != "" {
		return DetectionResult{
			Type:   explicitMode,
			Reason: fmt.Sprintf("explicit mode: %s", cfg.Type),
		}, nil
	}

	return DetectMode(cfg), nil
}

func ValidateConfigForMode(cfg *Config, detectedType DriverType) error {
	switch detectedType {
	case DriverSQLite:
		if cfg.DatabasePath == "" {
			return NewDriverError(detectedType, "validate",
				fmt.Errorf("%w: database path required", ErrInvalidConfig))
		}

	case DriverTursoRemote:
		if cfg.TursoURL == "" || cfg.TursoToken == "" {
			return NewDriverError(detectedType, "validate",
				fmt.Errorf("%w: turso URL and token required", ErrInvalidConfig))
		}

	case DriverTursoReplica:
		if cfg.TursoURL == "" || cfg.TursoToken == "" || cfg.DatabasePath == "" {
			return NewDriverError(detectedType, "validate",
				fmt.Errorf("%w: turso URL, token, and database path required", ErrInvalidConfig))
		}

	default:
		return NewDriverError(detectedType, "validate", ErrDriverNotRegistered)
	}

	return nil
}
