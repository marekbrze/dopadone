package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/example/dopadone/internal/domain"
	"github.com/example/dopadone/internal/service"
	"github.com/example/dopadone/internal/tui/areamodal"
	"github.com/example/dopadone/internal/tui/help"
	"github.com/example/dopadone/internal/tui/modal"
	"github.com/example/dopadone/internal/tui/spacemenu"
	"github.com/example/dopadone/internal/tui/theme"
	"github.com/example/dopadone/internal/tui/toast"
	"github.com/example/dopadone/internal/tui/tree"
	"github.com/example/dopadone/internal/tui/views"
)

type Model struct {
	areaSvc     service.AreaServiceInterface
	subareaSvc  service.SubareaServiceInterface
	projectSvc  service.ProjectServiceInterface
	taskSvc     service.TaskServiceInterface
	theme       theme.ColorTheme
	focus       FocusColumn
	width       int
	height      int
	ready       bool
	tabs        []views.Tab
	selectedTab int

	areas    []domain.Area
	subareas []domain.Subarea
	projects []domain.Project
	tasks    []domain.Task

	groupedTasks       *domain.GroupedTasks
	expandedTaskGroups map[string]bool

	selectedAreaIndex    int
	selectedSubareaIndex int
	selectedProjectIndex int
	selectedTaskIndex    int

	areaStates        map[string]*AreaState
	projectTree       *tree.TreeNode
	selectedProjectID string

	isLoadingAreas    bool
	isLoadingSubareas bool
	isLoadingProjects bool
	isLoadingTasks    bool

	spinner     spinner.Model
	modal       *modal.Modal
	isModalOpen bool

	areaModal       *areamodal.Modal
	isAreaModalOpen bool

	helpModal  *help.HelpModal
	isHelpOpen bool

	spaceMenu       *spacemenu.SpaceMenu
	isSpaceMenuOpen bool

	toasts []toast.Toast
}

func InitialModel(
	areaSvc service.AreaServiceInterface,
	subareaSvc service.SubareaServiceInterface,
	projectSvc service.ProjectServiceInterface,
	taskSvc service.TaskServiceInterface,
) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.Default.Primary)

	return Model{
		areaSvc:     areaSvc,
		subareaSvc:  subareaSvc,
		projectSvc:  projectSvc,
		taskSvc:     taskSvc,
		theme:       theme.Default,
		focus:       FocusSubareas,
		width:       0,
		height:      0,
		ready:       false,
		tabs:        []views.Tab{},
		selectedTab: 0,

		areas:    []domain.Area{},
		subareas: []domain.Subarea{},
		projects: []domain.Project{},
		tasks:    []domain.Task{},

		selectedAreaIndex:    0,
		selectedSubareaIndex: 0,
		selectedProjectIndex: 0,
		selectedTaskIndex:    0,

		areaStates: make(map[string]*AreaState),

		isLoadingAreas:    false,
		isLoadingSubareas: false,
		isLoadingProjects: false,
		isLoadingTasks:    false,

		spinner: s,

		helpModal:  nil,
		isHelpOpen: false,

		spaceMenu:       nil,
		isSpaceMenuOpen: false,

		areaModal:       nil,
		isAreaModalOpen: false,

		toasts: []toast.Toast{},
	}
}

func (m Model) Init() tea.Cmd {
	if m.areaSvc == nil {
		return nil
	}
	m.isLoadingAreas = true
	return tea.Batch(
		m.spinner.Tick,
		LoadAreasCmd(m.areaSvc),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return ToastTickMsg{}
		}),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case AreasLoadedMsg:
		m.isLoadingAreas = false
		if msg.Err != nil {
			m.addToast(toast.NewError("Failed to load areas: " + msg.Err.Error()))
			return m, nil
		}
		m.areas = msg.Areas
		m.tabs = updateTabsFromAreas(m.areas, m.selectedAreaIndex)
		m.selectedTab = m.selectedAreaIndex
		if len(m.areas) > 0 && m.selectedAreaIndex == 0 {
			m.selectedAreaIndex = 0
			m.isLoadingSubareas = true
			cmds = append(cmds, LoadSubareasCmd(m.subareaSvc, m.areas[0].ID))
		}
		return m, tea.Batch(cmds...)

	case SubareasLoadedMsg:
		m.isLoadingSubareas = false
		if msg.Err != nil {
			m.addToast(toast.NewError("Failed to load subareas: " + msg.Err.Error()))
			return m, nil
		}
		m.subareas = msg.Subareas
		if len(m.subareas) > 0 && m.selectedSubareaIndex == 0 {
			m.selectedSubareaIndex = 0
			m.isLoadingProjects = true
			cmds = append(cmds, LoadProjectsCmd(m.projectSvc, &m.subareas[0].ID))
		}
		return m, tea.Batch(cmds...)

	case ProjectsLoadedMsg:
		m.isLoadingProjects = false
		if msg.Err != nil {
			m.addToast(toast.NewError("Failed to load projects: " + msg.Err.Error()))
			return m, nil
		}
		m.projects = msg.Projects

		builder := tree.NewBuilder()
		m.projectTree = builder.BuildFromProjects(m.projects)

		if m.projectTree != nil {
			firstNode := tree.GetFirstVisibleNode(m.projectTree)
			if firstNode != nil {
				m.selectedProjectID = firstNode.ID
				m.selectedProjectIndex = 0
				m.isLoadingTasks = true
				m.tasks = nil
				m.selectedTaskIndex = 0
				cmds = append(cmds, LoadTasksCmd(m.taskSvc, firstNode.ID))
			}
		}
		return m, tea.Batch(cmds...)

	case TasksLoadedMsg:
		m.isLoadingTasks = false
		if msg.Err != nil {
			m.addToast(toast.NewError("Failed to load tasks: " + msg.Err.Error()))
			return m, nil
		}
		m.tasks = msg.Tasks

		m.groupedTasks = msg.GroupedTasks

		if m.expandedTaskGroups == nil {
			m.expandedTaskGroups = make(map[string]bool)
		}

		if m.groupedTasks != nil {
			for i := range m.groupedTasks.Groups {
				groupID := m.groupedTasks.Groups[i].ProjectID

				if _, exists := m.expandedTaskGroups[groupID]; !exists {
					m.expandedTaskGroups[groupID] = true
					m.groupedTasks.Groups[i].IsExpanded = true
				} else {
					m.groupedTasks.Groups[i].IsExpanded = m.expandedTaskGroups[groupID]
				}
			}
		}

		if m.selectedTaskIndex >= len(m.tasks) {
			m.selectedTaskIndex = 0
		}
		return m, nil

	case modal.SubmitMsg:
		return m.handleModalSubmit(msg)

	case modal.CloseMsg:
		m.isModalOpen = false
		m.modal = nil
		return m, nil

	case SubareaCreatedMsg:
		return m.handleSubareaCreated(msg)

	case ProjectCreatedMsg:
		return m.handleProjectCreated(msg)

	case TaskCreatedMsg:
		return m.handleTaskCreated(msg)

	case TaskStatusToggledMsg:
		if msg.Err != nil {
			if msg.TaskIndex < len(m.tasks) {
				m.tasks[msg.TaskIndex].Status = msg.OriginalStatus
			}

			m.addToast(toast.NewError("Failed to update task status: " + msg.Err.Error()))
			return m, nil
		}

		return m, nil

	case areamodal.SubmitMsg:
		return m.handleAreaModalSubmit(msg)
	case areamodal.UpdateMsg:
		return m.handleAreaModalUpdate(msg)
	case areamodal.DeleteMsg:
		return m.handleAreaModalDelete(msg)
	case areamodal.ReorderMsg:
		return m.handleAreaModalReorder(msg)
	case areamodal.CloseMsg:
		m.isAreaModalOpen = false
		m.areaModal = nil
		return m, nil
	case areamodal.LoadStatsMsg:
		return m, LoadAreaStatsCmd(m.areaSvc, msg.AreaID)

	case AreaCreatedMsg:
		return m.handleAreaCreated(msg)
	case AreaUpdatedMsg:
		return m.handleAreaUpdated(msg)
	case AreaDeletedMsg:
		return m.handleAreaDeleted(msg)
	case AreasReorderedMsg:
		return m.handleAreasReordered(msg)
	case AreaStatsLoadedMsg:
		return m.handleAreaStatsLoaded(msg)
	case LoadAreaStatsMsg:
		return m.handleLoadAreaStats(msg)

	case ToastTickMsg:
		m.removeExpiredToasts()
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return ToastTickMsg{}
		})

	case help.CloseMsg:
		m.isHelpOpen = false
		m.helpModal = nil
		return m, nil

	case spacemenu.CloseMsg:
		m.isSpaceMenuOpen = false
		m.spaceMenu = nil
		return m, nil

	case spacemenu.ActionMsg:
		return m.handleSpaceMenuAction(msg)

	case tea.KeyMsg:
		if m.isHelpOpen {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.helpModal, cmd = m.helpModal.Update(msg)
			return m, cmd
		}

		if m.isModalOpen {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.modal, cmd = m.modal.Update(msg)
			return m, cmd
		}

		if m.isAreaModalOpen && m.areaModal != nil {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.areaModal, cmd = m.areaModal.Update(msg)
			return m, cmd
		}

		if m.isSpaceMenuOpen && m.spaceMenu != nil {
			switch msg.String() {
			case "q", "ctrl+c":
				if m.spaceMenu != nil && m.spaceMenu.State() == spacemenu.StateMain {
					return m, tea.Quit
				}
			}
			var cmd tea.Cmd
			m.spaceMenu, cmd = m.spaceMenu.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			if !m.isModalOpen && !m.isAreaModalOpen && !m.isHelpOpen {
				m.isSpaceMenuOpen = true
				m.spaceMenu = spacemenu.New()
				m.spaceMenu, _ = m.spaceMenu.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
				return m, nil
			}
			m.handleEnterOrSpace()
		case "h", "left":
			m.focus = m.focus.Prev()
		case "l", "right":
			m.focus = m.focus.Next()
		case "tab":
			m.focus = m.focus.Next()
		case "j", "down":
			if !m.IsEmpty(m.focus) {
				return m.NavigateDownWithLoad(m.focus)
			}
		case "k", "up":
			if !m.IsEmpty(m.focus) {
				return m.NavigateUpWithLoad(m.focus)
			}
		case "[":
			return m, m.SwitchToPreviousArea()
		case "]":
			return m, m.SwitchToNextArea()
		case "enter":
			m.handleEnterOrSpace()
		case "x":
			if m.focus == FocusTasks && len(m.tasks) > 0 {
				return m, m.toggleTaskCompletion()
			}
		case "a":
			return m.handleQuickAdd()
		case "?":
			return m.handleHelp(), nil
		case "ctrl+a":
			return m.handleOpenAreaModal()
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		if m.helpModal != nil {
			m.helpModal, _ = m.helpModal.Update(msg)
		}
		if m.areaModal != nil {
			m.areaModal, _ = m.areaModal.Update(msg)
		}
		if m.spaceMenu != nil {
			m.spaceMenu, _ = m.spaceMenu.Update(msg)
		}
	}

	return m, nil
}

func (m *Model) toggleTaskCompletion() tea.Cmd {
	if len(m.tasks) == 0 || m.selectedTaskIndex >= len(m.tasks) {
		return nil
	}

	task := &m.tasks[m.selectedTaskIndex]

	var newStatus domain.TaskStatus
	if task.IsCompleted() {
		newStatus = domain.TaskStatusTodo
	} else {
		newStatus = domain.TaskStatusDone
	}

	originalStatus := task.Status

	task.Status = newStatus

	return ToggleTaskStatusCmd(m.taskSvc, task.ID, newStatus, originalStatus, m.selectedTaskIndex)
}

func (m Model) handleSpaceMenuAction(msg spacemenu.ActionMsg) (tea.Model, tea.Cmd) {
	m.isSpaceMenuOpen = false
	m.spaceMenu = nil

	switch msg.Action {
	case spacemenu.ActionQuit:
		return m, tea.Quit
	case spacemenu.ActionConfig:
		return m.handleOpenAreaModal()
	case spacemenu.ActionCreateArea, spacemenu.ActionEditArea, spacemenu.ActionDeleteArea:
		return m.handleOpenAreaModal()
	}

	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	tabs := views.TabsView(m.tabs, m.selectedTab)
	tabsRow := lipgloss.NewStyle().MarginBottom(1).Render(tabs)

	subareasContent := m.RenderSubareas()
	projectsContent := m.RenderProjects()
	tasksContent := m.RenderTasks()

	columns := []views.Column{
		{
			Title:     "Subareas",
			Content:   subareasContent,
			IsFocused: m.focus == FocusSubareas,
		},
		{
			Title:     "Projects",
			Content:   projectsContent,
			IsFocused: m.focus == FocusProjects,
		},
		{
			Title:     "Tasks",
			Content:   tasksContent,
			IsFocused: m.focus == FocusTasks,
		},
	}

	baseView := views.LayoutWithTabs(tabsRow, columns, m.width, m.height)

	toastView := m.RenderToasts()
	if toastView != "" {
		baseView = toastView + "\n" + baseView
	}

	footer := m.RenderFooter()
	if footer != "" {
		baseView = baseView + "\n" + footer
	}

	if m.isHelpOpen && m.helpModal != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			overlay(baseView, m.helpModal.View(), m.width, m.height),
		)
	}

	if m.isModalOpen && m.modal != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			overlay(baseView, m.modal.View(), m.width, m.height),
		)
	}

	if m.isAreaModalOpen && m.areaModal != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			overlay(baseView, m.areaModal.View(), m.width, m.height),
		)
	}

	if m.isSpaceMenuOpen && m.spaceMenu != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			overlay(baseView, m.spaceMenu.View(), m.width, m.height),
		)
	}

	return baseView
}

func updateTabsFromAreas(areas []domain.Area, selectedIndex int) []views.Tab {
	tabs := make([]views.Tab, len(areas))
	for i, area := range areas {
		tabs[i] = views.Tab{
			Name:     area.Name,
			ID:       area.ID,
			IsActive: i == selectedIndex,
		}
	}
	return tabs
}
