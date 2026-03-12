package configwizard

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/cli"
	"github.com/marekbrze/dopadone/internal/config"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/tui/internal/constants"
)

type Wizard struct {
	step wizardStep

	mode        DatabaseMode
	localPath   textinput.Model
	tursoURL    textinput.Model
	tursoToken  textinput.Model
	replicaPath textinput.Model

	selectedModeIndex int

	errorMsg    string
	spinner     spinner.Model
	isVerifying bool

	width  int
	height int

	configPath string
}

func New() *Wizard {
	defaultPath := cli.DefaultDBPath()

	localPath := textinput.New()
	localPath.Placeholder = "Path to local database"
	localPath.SetValue(defaultPath)
	localPath.CharLimit = 256
	localPath.Width = 50

	tursoURL := textinput.New()
	tursoURL.Placeholder = "libsql://your-db.turso.io"
	tursoURL.CharLimit = 256
	tursoURL.Width = 50

	tursoToken := textinput.New()
	tursoToken.Placeholder = "Your Turso auth token"
	tursoToken.CharLimit = 256
	tursoToken.Width = 50
	tursoToken.EchoMode = textinput.EchoPassword

	replicaPath := textinput.New()
	replicaPath.Placeholder = "Path to local replica"
	replicaPath.SetValue(defaultPath)
	replicaPath.CharLimit = 256
	replicaPath.Width = 50

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return &Wizard{
		step:              stepWelcome,
		mode:              ModeLocal,
		localPath:         localPath,
		tursoURL:          tursoURL,
		tursoToken:        tursoToken,
		replicaPath:       replicaPath,
		selectedModeIndex: 0,
		spinner:           s,
		configPath:        "",
	}
}

func (w *Wizard) Init() tea.Cmd {
	return w.spinner.Tick
}

func (w *Wizard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		return w, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		w.spinner, cmd = w.spinner.Update(msg)
		return w, cmd

	case VerificationSuccessMsg:
		w.isVerifying = false
		w.step = stepSuccess
		return w, nil

	case VerificationErrorMsg:
		w.isVerifying = false
		w.errorMsg = msg.Error.Error()
		w.step = stepConfig
		return w, nil

	case tea.KeyMsg:
		if w.isVerifying {
			return w, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return w, w.createFallbackAndQuit()

		case constants.KeyEsc:
			if w.step == stepWelcome {
				return w, tea.Quit
			}
			return w.handleBack()

		case constants.KeyEnter:
			return w.handleEnter()

		case "up", "k":
			return w.handleUp()

		case "down", "j":
			return w.handleDown()

		case "tab", "shift+tab":
			return w.handleTab()
		}
	}

	var cmd tea.Cmd
	switch w.step {
	case stepConfig:
		switch w.mode {
		case ModeLocal:
			w.localPath, cmd = w.localPath.Update(msg)
		case ModeRemote:
			if w.tursoURL.Focused() {
				w.tursoURL, cmd = w.tursoURL.Update(msg)
			} else if w.tursoToken.Focused() {
				w.tursoToken, cmd = w.tursoToken.Update(msg)
			}
		case ModeReplica:
			if w.replicaPath.Focused() {
				w.replicaPath, cmd = w.replicaPath.Update(msg)
			} else if w.tursoURL.Focused() {
				w.tursoURL, cmd = w.tursoURL.Update(msg)
			} else if w.tursoToken.Focused() {
				w.tursoToken, cmd = w.tursoToken.Update(msg)
			}
		}
	}

	return w, cmd
}

func (w *Wizard) handleEnter() (tea.Model, tea.Cmd) {
	switch w.step {
	case stepWelcome:
		w.step = stepModeSelection
		return w, nil

	case stepModeSelection:
		modes := []DatabaseMode{ModeLocal, ModeRemote, ModeReplica}
		w.mode = modes[w.selectedModeIndex]
		w.step = stepConfig
		w.focusFirstInput()
		return w, nil

	case stepConfig:
		if w.validateConfig() != nil {
			return w, nil
		}
		return w, w.verifyAndSave()

	case stepSuccess:
		return w, tea.Quit
	}

	return w, nil
}

func (w *Wizard) handleBack() (tea.Model, tea.Cmd) {
	switch w.step {
	case stepModeSelection:
		w.step = stepWelcome
		return w, nil

	case stepConfig:
		w.step = stepModeSelection
		w.errorMsg = ""
		return w, nil

	case stepSuccess:
		w.step = stepConfig
		return w, nil
	}

	return w, nil
}

func (w *Wizard) handleUp() (tea.Model, tea.Cmd) {
	if w.step == stepModeSelection {
		modes := []DatabaseMode{ModeLocal, ModeRemote, ModeReplica}
		if w.selectedModeIndex > 0 {
			w.selectedModeIndex--
		} else {
			w.selectedModeIndex = len(modes) - 1
		}
	}
	return w, nil
}

func (w *Wizard) handleDown() (tea.Model, tea.Cmd) {
	if w.step == stepModeSelection {
		modes := []DatabaseMode{ModeLocal, ModeRemote, ModeReplica}
		if w.selectedModeIndex < len(modes)-1 {
			w.selectedModeIndex++
		} else {
			w.selectedModeIndex = 0
		}
	}
	return w, nil
}

func (w *Wizard) handleTab() (tea.Model, tea.Cmd) {
	if w.step != stepConfig {
		return w, nil
	}

	switch w.mode {
	case ModeLocal:
	// Only one input, nothing to tab through
	case ModeRemote:
		if w.tursoURL.Focused() {
			w.tursoURL.Blur()
			w.tursoToken.Focus()
		} else {
			w.tursoToken.Blur()
			w.tursoURL.Focus()
		}
	case ModeReplica:
		if w.replicaPath.Focused() {
			w.replicaPath.Blur()
			w.tursoURL.Focus()
		} else if w.tursoURL.Focused() {
			w.tursoURL.Blur()
			w.tursoToken.Focus()
		} else {
			w.tursoToken.Blur()
			w.replicaPath.Focus()
		}
	}

	return w, nil
}

func (w *Wizard) focusFirstInput() {
	switch w.mode {
	case ModeLocal:
		w.localPath.Focus()
	case ModeRemote:
		w.tursoURL.Focus()
	case ModeReplica:
		w.replicaPath.Focus()
	}
}

func (w *Wizard) validateConfig() error {
	switch w.mode {
	case ModeLocal:
		if strings.TrimSpace(w.localPath.Value()) == "" {
			w.errorMsg = "Database path cannot be empty"
			return fmt.Errorf("empty path")
		}
	case ModeRemote:
		if strings.TrimSpace(w.tursoURL.Value()) == "" {
			w.errorMsg = "Turso URL cannot be empty"
			return fmt.Errorf("empty URL")
		}
		if strings.TrimSpace(w.tursoToken.Value()) == "" {
			w.errorMsg = "Turso auth token cannot be empty"
			return fmt.Errorf("empty token")
		}
	case ModeReplica:
		if strings.TrimSpace(w.replicaPath.Value()) == "" {
			w.errorMsg = "Local replica path cannot be empty"
			return fmt.Errorf("empty replica path")
		}
		if strings.TrimSpace(w.tursoURL.Value()) == "" {
			w.errorMsg = "Turso URL cannot be empty"
			return fmt.Errorf("empty URL")
		}
		if strings.TrimSpace(w.tursoToken.Value()) == "" {
			w.errorMsg = "Turso auth token cannot be empty"
			return fmt.Errorf("empty token")
		}
	}

	w.errorMsg = ""
	return nil
}

func (w *Wizard) verifyAndSave() tea.Cmd {
	return func() tea.Msg {
		result := w.getConfigResult()

		var dbConn *sql.DB
		var closeFunc func()

		switch w.mode {
		case ModeLocal:
			db, err := cli.Connect(result.DatabasePath)
			if err != nil {
				return VerificationErrorMsg{Error: fmt.Errorf("failed to connect to database: %w", err)}
			}
			dbConn = db
			closeFunc = func() { _ = db.Close() }

		case ModeRemote, ModeReplica:
			driverType := driver.DriverTursoRemote
			if w.mode == ModeReplica {
				driverType = driver.DriverTursoReplica
			}

			drv, err := cli.ConnectWithDriver(
				driver.WithDriverType(driverType),
				driver.WithDatabasePath(result.DatabasePath),
				driver.WithTurso(result.TursoURL, result.TursoToken),
				driver.WithSyncInterval(60*time.Second),
			)
			if err != nil {
				return VerificationErrorMsg{Error: fmt.Errorf("failed to create database driver: %w", err)}
			}

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err := drv.Connect(ctx); err != nil {
				_ = drv.Close()
				return VerificationErrorMsg{Error: fmt.Errorf("failed to connect to Turso: %w", err)}
			}

			if err := drv.Ping(ctx); err != nil {
				_ = drv.Close()
				return VerificationErrorMsg{Error: fmt.Errorf("failed to ping Turso database: %w", err)}
			}

			dbConn = drv.GetDB()
			closeFunc = func() { _ = drv.Close() }
		}

		if closeFunc != nil {
			defer closeFunc()
		}

		if err := cli.EnsureMigrations(dbConn); err != nil {
			return VerificationErrorMsg{Error: fmt.Errorf("failed to run migrations: %w", err)}
		}

		cfg := result.ToFileConfig()
		configPath, err := config.GetDefaultConfigPath()
		if err != nil {
			return VerificationErrorMsg{Error: fmt.Errorf("failed to get config path: %w", err)}
		}

		if err := config.SaveConfig(cfg, configPath); err != nil {
			return VerificationErrorMsg{Error: fmt.Errorf("failed to save config: %w", err)}
		}

		w.configPath = configPath
		return VerificationSuccessMsg{}
	}
}

func (w *Wizard) getConfigResult() *ConfigResult {
	result := &ConfigResult{
		Mode: w.mode,
	}

	switch w.mode {
	case ModeLocal:
		result.DatabasePath = strings.TrimSpace(w.localPath.Value())
	case ModeRemote:
		result.TursoURL = strings.TrimSpace(w.tursoURL.Value())
		result.TursoToken = strings.TrimSpace(w.tursoToken.Value())
	case ModeReplica:
		result.DatabasePath = strings.TrimSpace(w.replicaPath.Value())
		result.TursoURL = strings.TrimSpace(w.tursoURL.Value())
		result.TursoToken = strings.TrimSpace(w.tursoToken.Value())
	}

	return result
}

func (w *Wizard) View() string {
	var content strings.Builder

	switch w.step {
	case stepWelcome:
		content.WriteString(w.renderWelcome())
	case stepModeSelection:
		content.WriteString(w.renderModeSelection())
	case stepConfig:
		content.WriteString(w.renderConfig())
	case stepSuccess:
		content.WriteString(w.renderSuccess())
	}

	box := boxStyle.Render(content.String())

	return lipgloss.Place(
		w.width, w.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func (w *Wizard) renderWelcome() string {
	var content strings.Builder

	content.WriteString(brandStyle.Render("╔═══════════════════════════════════════════╗"))
	content.WriteString("\n")
	content.WriteString(brandStyle.Render("║                                           ║"))
	content.WriteString("\n")
	content.WriteString(brandStyle.Render("║        Welcome to Dopadone!               ║"))
	content.WriteString("\n")
	content.WriteString(brandStyle.Render("║                                           ║"))
	content.WriteString("\n")
	content.WriteString(brandStyle.Render("╚═══════════════════════════════════════════╝"))
	content.WriteString("\n\n")

	content.WriteString(subtitleStyle.Render("Your ADHD-friendly project management companion"))
	content.WriteString("\n\n")

	content.WriteString("Let's set up your database configuration.")
	content.WriteString("\n")
	content.WriteString("This will only take a moment.")
	content.WriteString("\n\n")

	content.WriteString(hintStyle.Render("Press Enter to continue • ESC to exit"))

	return content.String()
}

func (w *Wizard) renderModeSelection() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Choose Database Mode"))
	content.WriteString("\n\n")

	modes := []struct {
		name        string
		description string
	}{
		{
			name:        "Local SQLite",
			description: "Single device, offline-first (recommended)",
		},
		{
			name:        "Turso Remote",
			description: "Cloud-only, requires internet connection",
		},
		{
			name:        "Turso Replica",
			description: "Hybrid: local storage + cloud sync",
		},
	}

	for i, mode := range modes {
		prefix := "  "
		if i == w.selectedModeIndex {
			prefix = modeSelectedStyle.Render("▸ ")
		}

		line := fmt.Sprintf("%s%s", prefix, mode.name)
		if i == w.selectedModeIndex {
			line = modeSelectedStyle.Render(line)
		} else {
			line = modeTitleStyle.Render(line)
		}
		content.WriteString(line)
		content.WriteString("\n")

		desc := fmt.Sprintf("%s%s", strings.Repeat(" ", 4), mode.description)
		content.WriteString(modeDescStyle.Render(desc))
		content.WriteString("\n\n")
	}

	content.WriteString(hintStyle.Render("↑/↓: Navigate • Enter: Select • ESC: Back"))

	return content.String()
}

func (w *Wizard) renderConfig() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Configure " + string(w.mode) + " Mode"))
	content.WriteString("\n\n")

	switch w.mode {
	case ModeLocal:
		content.WriteString(inputLabelStyle.Render("Database Path:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.localPath.View()))
		content.WriteString("\n")
		content.WriteString(modeDescStyle.Render("The database will be stored locally on this device."))
		content.WriteString("\n\n")

	case ModeRemote:
		content.WriteString(inputLabelStyle.Render("Turso Database URL:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.tursoURL.View()))
		content.WriteString("\n\n")

		content.WriteString(inputLabelStyle.Render("Auth Token:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.tursoToken.View()))
		content.WriteString("\n")
		content.WriteString(modeDescStyle.Render("Get your credentials from turso.io"))
		content.WriteString("\n\n")

	case ModeReplica:
		content.WriteString(inputLabelStyle.Render("Local Replica Path:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.replicaPath.View()))
		content.WriteString("\n\n")

		content.WriteString(inputLabelStyle.Render("Turso Database URL:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.tursoURL.View()))
		content.WriteString("\n\n")

		content.WriteString(inputLabelStyle.Render("Auth Token:"))
		content.WriteString("\n")
		content.WriteString(inputFieldStyle.Render(w.tursoToken.View()))
		content.WriteString("\n\n")
	}

	if w.errorMsg != "" {
		content.WriteString(errorStyle.Render("✗ " + w.errorMsg))
		content.WriteString("\n")
	}

	if w.isVerifying {
		content.WriteString(spinnerStyle.Render(w.spinner.View() + " Saving configuration..."))
		content.WriteString("\n")
	}

	content.WriteString(hintStyle.Render("Enter: Continue • Tab: Next field • ESC: Back"))

	return content.String()
}

func (w *Wizard) renderSuccess() string {
	var content strings.Builder

	content.WriteString(successStyle.Render("✓ Configuration Complete!"))
	content.WriteString("\n\n")

	content.WriteString("Your database has been configured.")
	content.WriteString("\n\n")

	if w.configPath != "" {
		content.WriteString(modeDescStyle.Render(fmt.Sprintf("Config saved to: %s", w.configPath)))
		content.WriteString("\n\n")
	}

	switch w.mode {
	case ModeLocal:
		content.WriteString(modeDescStyle.Render(fmt.Sprintf("Database: %s", w.localPath.Value())))
	case ModeRemote:
		content.WriteString(modeDescStyle.Render("Connected to Turso cloud"))
	case ModeReplica:
		content.WriteString(modeDescStyle.Render(fmt.Sprintf("Local: %s", w.replicaPath.Value())))
		content.WriteString("\n")
		content.WriteString(modeDescStyle.Render("Syncing with Turso cloud"))
	}

	content.WriteString("\n\n")
	content.WriteString(hintStyle.Render("Press Enter to start Dopadone"))

	return content.String()
}

func (w *Wizard) SetError(err string) {
	w.errorMsg = err
}

func (w *Wizard) ClearError() {
	w.errorMsg = ""
}

func (w *Wizard) IsComplete() bool {
	return w.step == stepSuccess
}

func (w *Wizard) GetConfigPath() string {
	return w.configPath
}

func (w *Wizard) createFallbackAndQuit() tea.Cmd {
	return func() tea.Msg {
		fallbackResult := &ConfigResult{
			Mode:         ModeLocal,
			DatabasePath: cli.DefaultDBPath(),
		}
		cfg := fallbackResult.ToFileConfig()

		configPath, err := config.GetDefaultConfigPath()
		if err == nil {
			_ = config.SaveConfig(cfg, configPath)
			w.configPath = configPath
		}

		return tea.Quit()
	}
}
