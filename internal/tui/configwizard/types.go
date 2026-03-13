package configwizard

import (
	"github.com/marekbrze/dopadone/internal/config"
)

type ConfigResult struct {
	DatabasePath string
	Mode         DatabaseMode
	TursoURL     string
	TursoToken   string
}

func (r *ConfigResult) ToFileConfig() *config.FileConfig {
	cfg := &config.FileConfig{
		Database: config.DatabaseConfig{
			Mode: string(r.Mode),
		},
	}
	switch r.Mode {
	case ModeLocal:
		cfg.Database.Path = r.DatabasePath
	case ModeRemote:
		cfg.Database.Turso = config.TursoConfig{
			URL:   r.TursoURL,
			Token: r.TursoToken,
		}
	case ModeReplica:
		cfg.Database.Path = r.DatabasePath
		cfg.Database.Turso = config.TursoConfig{
			URL:   r.TursoURL,
			Token: r.TursoToken,
		}
	}
	return cfg
}

type SubmitMsg struct {
	Result *ConfigResult
}

type CancelMsg struct{}

type VerificationSuccessMsg struct{}

type VerificationErrorMsg struct {
	Error error
}

type FallbackCreatedMsg struct {
	ConfigPath string
}

type DatabaseMode string

const (
	ModeLocal   DatabaseMode = "local"
	ModeRemote  DatabaseMode = "remote"
	ModeReplica DatabaseMode = "replica"
)

type wizardStep int

const (
	stepWelcome wizardStep = iota
	stepModeSelection
	stepConfig
	stepVerifying
	stepSuccess
)

type WelcomeOption int

const (
	WelcomeOptionQuickStart WelcomeOption = iota
	WelcomeOptionCustomSetup
	WelcomeOptionExit
)
