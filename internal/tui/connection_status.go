package tui

import (
	"time"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

const (
	SymbolConnected = "●"
	SymbolSyncing   = "◐"
	SymbolOffline   = "○"
	SymbolLocalOnly = "■"
)

type ConnectionStatusView struct {
	Mode         driver.DriverType
	Status       driver.ConnectionStatus
	SyncStatus   driver.SyncStatus
	LastSyncAt   time.Time
	ErrorMessage string
}

func NewConnectionStatusView(drv driver.DatabaseDriver) ConnectionStatusView {
	if drv == nil {
		return ConnectionStatusView{
			Mode:   driver.DriverSQLite,
			Status: driver.StatusConnected,
		}
	}

	view := ConnectionStatusView{
		Mode:   drv.Type(),
		Status: drv.Status(),
	}

	if replica, ok := drv.(interface{ SyncInfo() driver.SyncInfo }); ok {
		info := replica.SyncInfo()
		view.SyncStatus = info.Status
		view.LastSyncAt = info.LastSyncAt
		if info.LastError != nil {
			view.ErrorMessage = info.LastError.Error()
		}
	}

	return view
}

func (v ConnectionStatusView) Symbol() string {
	switch v.Mode {
	case driver.DriverSQLite:
		return SymbolLocalOnly
	case driver.DriverTursoRemote, driver.DriverTursoReplica:
		switch v.Status {
		case driver.StatusConnected:
			if v.SyncStatus == driver.SyncStatusSyncing {
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

func (v ConnectionStatusView) Tooltip() string {
	switch v.Mode {
	case driver.DriverSQLite:
		return "Local database (no sync)"
	case driver.DriverTursoRemote:
		switch v.Status {
		case driver.StatusConnected:
			return "Connected to Turso"
		case driver.StatusConnecting:
			return "Connecting to Turso..."
		case driver.StatusDisconnected:
			return "Disconnected from Turso"
		case driver.StatusError:
			if v.ErrorMessage != "" {
				return "Error: " + v.ErrorMessage
			}
			return "Connection error"
		}
	case driver.DriverTursoReplica:
		switch v.Status {
		case driver.StatusConnected:
			if v.SyncStatus == driver.SyncStatusSyncing {
				return "Syncing with Turso..."
			}
			if v.LastSyncAt.IsZero() {
				return "Connected (not yet synced)"
			}
			elapsed := time.Since(v.LastSyncAt)
			return "Connected (last sync: " + formatDuration(elapsed) + " ago)"
		case driver.StatusConnecting:
			return "Connecting to Turso..."
		case driver.StatusDisconnected:
			return "Offline - changes will sync when connected"
		case driver.StatusError:
			if v.ErrorMessage != "" {
				return "Error: " + v.ErrorMessage
			}
			return "Connection error"
		}
	}
	return "Unknown status"
}

func (v ConnectionStatusView) ModeLabel() string {
	switch v.Mode {
	case driver.DriverSQLite:
		return "local"
	case driver.DriverTursoRemote:
		return "remote"
	case driver.DriverTursoReplica:
		return "replica"
	}
	return "unknown"
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "<1m"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		return string(rune(mins/10+'0')) + string(rune(mins%10+'0')) + "m"
	}
	hours := int(d.Hours())
	if hours < 24 {
		return string(rune(hours/10+'0')) + string(rune(hours%10+'0')) + "h"
	}
	days := hours / 24
	return string(rune(days/10+'0')) + string(rune(days%10+'0')) + "d"
}
