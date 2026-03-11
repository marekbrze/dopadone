package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db/driver"
)

func TestConnectionStatusView_Symbol(t *testing.T) {
	tests := []struct {
		name       string
		view       ConnectionStatusView
		wantSymbol string
	}{
		{
			name: "local_sqlite",
			view: ConnectionStatusView{
				Mode:   driver.DriverSQLite,
				Status: driver.StatusConnected,
			},
			wantSymbol: SymbolLocalOnly,
		},
		{
			name: "remote_connected",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusConnected,
			},
			wantSymbol: SymbolConnected,
		},
		{
			name: "remote_connecting",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusConnecting,
			},
			wantSymbol: SymbolSyncing,
		},
		{
			name: "remote_disconnected",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusDisconnected,
			},
			wantSymbol: SymbolOffline,
		},
		{
			name: "remote_error",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusError,
			},
			wantSymbol: SymbolOffline,
		},
		{
			name: "replica_connected_idle",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusConnected,
				SyncStatus: driver.SyncStatusIdle,
			},
			wantSymbol: SymbolConnected,
		},
		{
			name: "replica_connected_syncing",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusConnected,
				SyncStatus: driver.SyncStatusSyncing,
			},
			wantSymbol: SymbolSyncing,
		},
		{
			name: "replica_offline",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusDisconnected,
				SyncStatus: driver.SyncStatusOffline,
			},
			wantSymbol: SymbolOffline,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.view.Symbol()
			if got != tt.wantSymbol {
				t.Errorf("Symbol() = %q, want %q", got, tt.wantSymbol)
			}
		})
	}
}

func TestConnectionStatusView_Tooltip(t *testing.T) {
	tests := []struct {
		name        string
		view        ConnectionStatusView
		wantContain string
	}{
		{
			name: "local_sqlite",
			view: ConnectionStatusView{
				Mode:   driver.DriverSQLite,
				Status: driver.StatusConnected,
			},
			wantContain: "Local database",
		},
		{
			name: "remote_connected",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusConnected,
			},
			wantContain: "Connected to Turso",
		},
		{
			name: "remote_connecting",
			view: ConnectionStatusView{
				Mode:   driver.DriverTursoRemote,
				Status: driver.StatusConnecting,
			},
			wantContain: "Connecting",
		},
		{
			name: "replica_connected_with_sync",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusConnected,
				SyncStatus: driver.SyncStatusSyncing,
			},
			wantContain: "Syncing",
		},
		{
			name: "replica_connected_recent_sync",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusConnected,
				SyncStatus: driver.SyncStatusIdle,
				LastSyncAt: time.Now().Add(-30 * time.Second),
			},
			wantContain: "sync",
		},
		{
			name: "replica_offline",
			view: ConnectionStatusView{
				Mode:       driver.DriverTursoReplica,
				Status:     driver.StatusDisconnected,
				SyncStatus: driver.SyncStatusOffline,
			},
			wantContain: "Offline",
		},
		{
			name: "error_with_message",
			view: ConnectionStatusView{
				Mode:         driver.DriverTursoRemote,
				Status:       driver.StatusError,
				ErrorMessage: "connection refused",
			},
			wantContain: "connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.view.Tooltip()
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("Tooltip() = %q, want to contain %q", got, tt.wantContain)
			}
		})
	}
}

func TestConnectionStatusView_ModeLabel(t *testing.T) {
	tests := []struct {
		name      string
		view      ConnectionStatusView
		wantLabel string
	}{
		{
			name:      "sqlite",
			view:      ConnectionStatusView{Mode: driver.DriverSQLite},
			wantLabel: "local",
		},
		{
			name:      "turso_remote",
			view:      ConnectionStatusView{Mode: driver.DriverTursoRemote},
			wantLabel: "remote",
		},
		{
			name:      "turso_replica",
			view:      ConnectionStatusView{Mode: driver.DriverTursoReplica},
			wantLabel: "replica",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.view.ModeLabel()
			if got != tt.wantLabel {
				t.Errorf("ModeLabel() = %q, want %q", got, tt.wantLabel)
			}
		})
	}
}

func TestNewConnectionStatusView_NilDriver(t *testing.T) {
	view := NewConnectionStatusView(nil)
	if view.Mode != driver.DriverSQLite {
		t.Errorf("NewConnectionStatusView(nil).Mode = %v, want %v", view.Mode, driver.DriverSQLite)
	}
	if view.Status != driver.StatusConnected {
		t.Errorf("NewConnectionStatusView(nil).Status = %v, want %v", view.Status, driver.StatusConnected)
	}
}
