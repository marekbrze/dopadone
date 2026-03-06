package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/dopadone/internal/tui/areamodal"
	"github.com/example/dopadone/internal/tui/toast"
)

func (m *Model) handleAreaModalSubmit(msg areamodal.SubmitMsg) (tea.Model, tea.Cmd) {
	m.isAreaModalOpen = false
	m.areaModal = nil
	return m, CreateAreaCmd(m.areaSvc, msg.Name, msg.Color)
}

func (m *Model) handleAreaModalUpdate(msg areamodal.UpdateMsg) (tea.Model, tea.Cmd) {
	m.isAreaModalOpen = false
	m.areaModal = nil
	return m, UpdateAreaCmd(m.areaSvc, msg.ID, msg.Name, msg.Color)
}

func (m *Model) handleAreaModalDelete(msg areamodal.DeleteMsg) (tea.Model, tea.Cmd) {
	m.isAreaModalOpen = false
	m.areaModal = nil
	return m, DeleteAreaCmd(m.areaSvc, msg.ID, msg.Hard)
}

func (m *Model) handleAreaModalReorder(msg areamodal.ReorderMsg) (tea.Model, tea.Cmd) {
	m.isAreaModalOpen = false
	m.areaModal = nil
	return m, ReorderAreasCmd(m.areaSvc, msg.AreaIDs)
}

func (m *Model) handleLoadAreaStats(msg LoadAreaStatsMsg) (tea.Model, tea.Cmd) {
	return m, LoadAreaStatsCmd(m.areaSvc, msg.AreaID)
}

func (m *Model) handleAreaStatsLoaded(msg AreaStatsLoadedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.addToast(toast.NewError("Failed to load area stats: " + msg.Err.Error()))
		return m, nil
	}
	if m.areaModal != nil && m.isAreaModalOpen {
		m.areaModal.SetStats(areamodal.Stats{
			Subareas: msg.Stats.Subareas,
			Projects: msg.Stats.Projects,
			Tasks:    msg.Stats.Tasks,
		})
	}
	return m, nil
}

func (m *Model) handleAreaCreated(msg AreaCreatedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.addToast(toast.NewError("Failed to create area: " + msg.Err.Error()))
		return m, nil
	}

	m.addToast(toast.NewSuccess("Area created successfully"))
	m.isLoadingAreas = true
	return m, LoadAreasCmd(m.areaSvc)
}

func (m *Model) handleAreaUpdated(msg AreaUpdatedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.addToast(toast.NewError("Failed to update area: " + msg.Err.Error()))
		return m, nil
	}

	m.addToast(toast.NewSuccess("Area updated successfully"))
	m.isLoadingAreas = true
	return m, LoadAreasCmd(m.areaSvc)
}

func (m *Model) handleAreaDeleted(msg AreaDeletedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.addToast(toast.NewError("Failed to delete area: " + msg.Err.Error()))
		return m, nil
	}

	deleteType := "permanently"
	if !msg.Hard {
		deleteType = "soft"
	}
	m.addToast(toast.NewSuccess(fmt.Sprintf("Area %s deleted", deleteType)))
	m.isLoadingAreas = true
	return m, LoadAreasCmd(m.areaSvc)
}

func (m *Model) handleAreasReordered(msg AreasReorderedMsg) (tea.Model, tea.Cmd) {
	if msg.Err != nil {
		m.addToast(toast.NewError("Failed to reorder areas: " + msg.Err.Error()))
		return m, nil
	}

	m.addToast(toast.NewSuccess("Areas reordered successfully"))
	m.isLoadingAreas = true
	return m, LoadAreasCmd(m.areaSvc)
}
