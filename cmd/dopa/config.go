package main

import (
	"fmt"
	"os"
	"time"

	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/config"
	"github.com/marekbrze/dopadone/internal/db/driver"
)

type Config struct {
	DatabasePath string
	TursoURL     string
	TursoToken   string
	DBMode       string
	SyncInterval time.Duration
}

type LoadConfigParams struct {
	DBPath       string
	TursoURL     string
	TursoToken   string
	DBMode       string
	SyncInterval time.Duration
	ConfigPath   string
}

func LoadConfig(params LoadConfigParams) (*Config, error) {
	fileCfg, err := loadFileConfig(params.ConfigPath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DatabasePath: resolveDBPath(params.DBPath, fileCfg),
		TursoURL:     resolveTursoURL(params.TursoURL, fileCfg),
		TursoToken:   resolveTursoToken(params.TursoToken, fileCfg),
		DBMode:       resolveDBMode(params.DBMode, fileCfg),
		SyncInterval: resolveSyncInterval(params.SyncInterval, fileCfg),
	}

	return cfg, nil
}

func loadFileConfig(configPath string) (*config.FileConfig, error) {
	path, err := config.DiscoverConfigWithExplicit(configPath)
	if err != nil {
		if err == config.ErrNoConfigFile {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to discover config file: %w", err)
	}

	fileCfg, err := config.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	return fileCfg, nil
}

func resolveDBPath(cliValue string, fileCfg *config.FileConfig) string {
	if cliValue != "" {
		return cliValue
	}
	if env := os.Getenv("DOPA_DB_PATH"); env != "" {
		return env
	}
	if fileCfg != nil && fileCfg.Database.Path != "" {
		return fileCfg.Database.Path
	}
	return cli.DefaultDBPath()
}

func resolveTursoURL(cliValue string, fileCfg *config.FileConfig) string {
	if cliValue != "" {
		return cliValue
	}
	if env := os.Getenv("TURSO_DATABASE_URL"); env != "" {
		return env
	}
	if fileCfg != nil && fileCfg.Database.Turso.URL != "" {
		return fileCfg.Database.Turso.URL
	}
	return ""
}

func resolveTursoToken(cliValue string, fileCfg *config.FileConfig) string {
	if cliValue != "" {
		return cliValue
	}
	if env := os.Getenv("TURSO_AUTH_TOKEN"); env != "" {
		return env
	}
	if fileCfg != nil && fileCfg.Database.Turso.Token != "" {
		return fileCfg.Database.Turso.Token
	}
	return ""
}

func resolveDBMode(cliValue string, fileCfg *config.FileConfig) string {
	if cliValue != "" {
		return cliValue
	}
	if env := os.Getenv("DOPA_DB_MODE"); env != "" {
		return env
	}
	if fileCfg != nil && fileCfg.Database.Mode != "" {
		return fileCfg.Database.Mode
	}
	return ""
}

func resolveSyncInterval(cliValue time.Duration, fileCfg *config.FileConfig) time.Duration {
	if cliValue != 0 && cliValue != 60*time.Second {
		return cliValue
	}
	if fileCfg != nil && fileCfg.Database.SyncInterval != "" {
		if dur, err := fileCfg.SyncIntervalDuration(); err == nil && dur != 0 {
			return dur
		}
	}
	return cliValue
}

func (c *Config) ToDriverConfig() *driver.Config {
	return &driver.Config{
		Type:         driver.DriverType(c.DBMode),
		DatabasePath: c.DatabasePath,
		TursoURL:     c.TursoURL,
		TursoToken:   c.TursoToken,
		SyncInterval: c.SyncInterval,
	}
}
