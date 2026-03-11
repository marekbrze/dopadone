package statusindicator

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/db/driver"
)

type StatusIndicator struct {
	mode       driver.DriverType
	status     driver.ConnectionStatus
	syncStatus driver.SyncStatus
}

func New(mode driver.DriverType, status driver.ConnectionStatus, syncStatus driver.SyncStatus) *StatusIndicator {
	return &StatusIndicator{
		mode:       mode,
		status:     status,
		syncStatus: syncStatus,
	}
}

func (s *StatusIndicator) Render() string {
	symbol := s.getSymbol()
	color := s.getColor()

	styledSymbol := lipgloss.NewStyle().
		Foreground(color).
		Render(symbol)

	modeLabel := s.getModeLabel()

	return fmt.Sprintf("%s %s", styledSymbol, modeLabel)
}

func (s *StatusIndicator) getSymbol() string {
	switch s.mode {
	case driver.DriverSQLite:
		return SymbolLocalOnly
	case driver.DriverTursoRemote, driver.DriverTursoReplica:
		switch s.status {
		case driver.StatusConnected:
			if s.syncStatus == driver.SyncStatusSyncing {
				return SymbolSyncing
			}
			return SymbolConnected
		case driver.StatusConnecting:
			return SymbolSyncing
		case driver.StatusDisconnected, driver.StatusError:
			return SymbolOffline
		}
	}
	return SymbolLocalOnly
}

func (s *StatusIndicator) getColor() lipgloss.Color {
	switch s.mode {
	case driver.DriverSQLite:
		return ColorGray
	case driver.DriverTursoRemote, driver.DriverTursoReplica:
		switch s.status {
		case driver.StatusConnected:
			if s.syncStatus == driver.SyncStatusSyncing {
				return ColorYellow
			}
			return ColorGreen
		case driver.StatusConnecting:
			return ColorYellow
		case driver.StatusDisconnected, driver.StatusError:
			return ColorRed
		}
	}
	return ColorGray
}

func (s *StatusIndicator) getModeLabel() string {
	switch s.mode {
	case driver.DriverSQLite:
		return "local"
	case driver.DriverTursoRemote:
		return "remote"
	case driver.DriverTursoReplica:
		return "replica"
	}
	return "unknown"
}
