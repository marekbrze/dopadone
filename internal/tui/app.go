package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/service"
	"github.com/marekbrze/dopadone/internal/tui/areamodal"
	"github.com/marekbrze/dopadone/internal/tui/confirmmodal"
	"github.com/marekbrze/dopadone/internal/tui/help"
	"github.com/marekbrze/dopadone/internal/tui/modal"
	"github.com/marekbrze/dopadone/internal/tui/spacemenu"
	"github.com/marekbrze/dopadone/internal/tui/theme"
	"github.com/marekbrze/dopadone/internal/tui/toast"
	"github.com/marekbrze/dopadone/internal/tui/tree"
	"github.com/marekbrze/dopadone/internal/tui/views"
)

type Model struct {
	areaSvc     service.AreaServiceInterface
	subareaSvc  service.SubareaServiceInterface
	projectSvc  service.ProjectServiceInterface
	taskSvc     service.TaskServiceInterface
	dbDriver    driver.DatabaseDriver
	theme       theme.ColorTheme
	focus       FocusColumn
	width       int
	height      int
	ready       bool
	tabs        []views.Tab
	selectedTab int

	connectionStatus ConnectionStatusView

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

	areaLoadError    error
	subareaLoadError error
	projectLoadError error
	taskLoadError    error

	spinner     spinner.Model
	modal       *modal.Modal
	isModalOpen bool

	areaModal       *areamodal.Modal
	isAreaModalOpen bool

	helpModal  *help.HelpModal
	isHelpOpen bool

	spaceMenu       *spacemenu.SpaceMenu
	isSpaceMenuOpen bool

	confirmModal       *confirmmodal.Modal
	isConfirmModalOpen bool

	toasts []toast.Toast
}

func InitialModel(
	areaSvc service.AreaServiceInterface,
	subareaSvc service.SubareaServiceInterface,
	projectSvc service.ProjectServiceInterface,
	taskSvc service.TaskServiceInterface,
	dbDriver driver.DatabaseDriver,
) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.Default.Primary)

	return Model{
		areaSvc:     areaSvc,
		subareaSvc:  subareaSvc,
		projectSvc:  projectSvc,
		taskSvc:     taskSvc,
		dbDriver:    dbDriver,
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

		areaLoadError:    nil,
		subareaLoadError: nil,
		projectLoadError: nil,
		taskLoadError:    nil,

		connectionStatus: NewConnectionStatusView(dbDriver),
	}
}

func (m Model) Init() tea.Cmd {
	if m.areaSvc == nil {
		return nil
	}
	cmds := []tea.Cmd{
		m.spinner.Tick,
		LoadAreasCmd(m.areaSvc),
		tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return ToastTickMsg{}
		}),
	}
	if m.dbDriver != nil {
		cmds = append(cmds, PollConnectionStatusCmd(m.dbDriver))
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if model, cmd, handled := m.handleAreaMessages(msg); handled {
		return model, cmd
	}

	if model, cmd, handled := m.handleConnectionMessages(msg); handled {
		return model, cmd
	}

	switch msg := msg.(type) {
	case spinner.TickMsg:
		return m.handleSpinnerTick(msg)

	case AreasLoadedMsg:
		return m.handleAreasLoaded(msg)

	case SubareasLoadedMsg:
		return m.handleSubareasLoaded(msg)

	case ProjectsLoadedMsg:
		return m.handleProjectsLoaded(msg)

	case TasksLoadedMsg:
		return m.handleTasksLoaded(msg)

	case modal.SubmitMsg:
		return m.handleModalSubmit(msg)

	case modal.CloseMsg:
		return m.handleModalClose()

	case SubareaCreatedMsg:
		return m.handleSubareaCreated(msg)

	case ProjectCreatedMsg:
		return m.handleProjectCreated(msg)

	case TaskCreatedMsg:
		return m.handleTaskCreated(msg)

	case TaskStatusToggledMsg:
		return m.handleTaskStatusToggled(msg)

	case ToastTickMsg:
		return m.handleToastTick()

	case help.CloseMsg:
		return m.handleHelpClose()

	case spacemenu.CloseMsg:
		return m.handleSpaceMenuClose()

	case spacemenu.ActionMsg:
		return m.handleSpaceMenuAction(msg)

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case confirmmodal.ConfirmMsg:
		return m.handleConfirmModalConfirm(msg)
	case confirmmodal.CancelMsg:
		m.handleConfirmModalCancel()
		return m, nil
	case DeleteSuccessMsg:
		return m.handleDeleteSuccess(msg)
	case DeleteErrorMsg:
		m.handleDeleteError(msg)
		return m, nil

	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	}

	return m, nil
}

func (m *Model) toggleTaskCompletion() tea.Cmd {
	task := m.getTaskAtLine(m.selectedTaskIndex)
	if task == nil {
		return nil
	}

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

func (m *Model) handleSpaceMenuClose() (tea.Model, tea.Cmd) {
	m.isSpaceMenuOpen = false
	m.spaceMenu = nil
	return m, nil
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

	if m.isConfirmModalOpen && m.confirmModal != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			overlay(baseView, m.confirmModal.View(), m.width, m.height),
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
