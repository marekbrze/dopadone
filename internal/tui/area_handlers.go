package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/tui/areamodal"
	"github.com/marekbrze/dopadone/internal/tui/toast"
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

func (m *Model) handleAreaMessages(msg interface{}) (tea.Model, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case areamodal.SubmitMsg:
		model, cmd := m.handleAreaModalSubmit(msg)
		return model, cmd, true
	case areamodal.UpdateMsg:
		model, cmd := m.handleAreaModalUpdate(msg)
		return model, cmd, true
	case areamodal.DeleteMsg:
		model, cmd := m.handleAreaModalDelete(msg)
		return model, cmd, true
	case areamodal.ReorderMsg:
		model, cmd := m.handleAreaModalReorder(msg)
		return model, cmd, true
	case areamodal.CloseMsg:
		model, cmd := m.handleAreaModalClose()
		return model, cmd, true
	case areamodal.LoadStatsMsg:
		return m, LoadAreaStatsCmd(m.areaSvc, msg.AreaID), true
	case AreaCreatedMsg:
		model, cmd := m.handleAreaCreated(msg)
		return model, cmd, true
	case AreaUpdatedMsg:
		model, cmd := m.handleAreaUpdated(msg)
		return model, cmd, true
	case AreaDeletedMsg:
		model, cmd := m.handleAreaDeleted(msg)
		return model, cmd, true
	case AreasReorderedMsg:
		model, cmd := m.handleAreasReordered(msg)
		return model, cmd, true
	case AreaStatsLoadedMsg:
		model, cmd := m.handleAreaStatsLoaded(msg)
		return model, cmd, true
	case LoadAreaStatsMsg:
		model, cmd := m.handleLoadAreaStats(msg)
		return model, cmd, true
	}
	return m, nil, false
}
