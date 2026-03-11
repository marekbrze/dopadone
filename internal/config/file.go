package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type FileConfig struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Path         string      `yaml:"path"`
	Mode         string      `yaml:"mode"`
	SyncInterval string      `yaml:"sync_interval"`
	Turso        TursoConfig `yaml:"turso"`
}

type TursoConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}

func (f *FileConfig) SyncIntervalDuration() (time.Duration, error) {
	if f.Database.SyncInterval == "" {
		return 0, nil
	}
	return time.ParseDuration(f.Database.SyncInterval)
}

func ParseFile(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if len(data) == 0 {
		return &FileConfig{}, nil
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *FileConfig) error {
	if cfg.Database.Mode != "" {
		validModes := map[string]bool{
			"local":   true,
			"remote":  true,
			"replica": true,
			"auto":    true,
		}
		if !validModes[cfg.Database.Mode] {
			return fmt.Errorf("invalid database mode: %s (valid: local, remote, replica, auto)", cfg.Database.Mode)
		}
	}

	if cfg.Database.SyncInterval != "" {
		if _, err := time.ParseDuration(cfg.Database.SyncInterval); err != nil {
			return fmt.Errorf("invalid sync_interval: %w", err)
		}
	}

	return nil
}

var (
	ErrNoConfigFile = errors.New("no config file found")
)

func DiscoverConfig() (string, error) {
	paths := getDiscoveryPaths()

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", ErrNoConfigFile
}

func DiscoverConfigWithExplicit(explicitPath string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); err != nil {
			return "", fmt.Errorf("config file not found: %s", explicitPath)
		}
		return explicitPath, nil
	}

	return DiscoverConfig()
}

func getDiscoveryPaths() []string {
	var paths []string

	cwd, err := os.Getwd()
	if err == nil {
		paths = append(paths, filepath.Join(cwd, "dopadone.yaml"))
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		paths = append(paths, filepath.Join(xdgConfigHome, "dopadone", "config.yaml"))
	}

	if homeDir, err := os.UserHomeDir(); err == nil {
		if xdgConfigHome == "" {
			paths = append(paths, filepath.Join(homeDir, ".config", "dopadone", "config.yaml"))
		}
		paths = append(paths, filepath.Join(homeDir, ".dopadone.yaml"))
	}

	return paths
}
