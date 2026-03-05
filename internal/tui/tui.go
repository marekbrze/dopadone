package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/example/projectdb/internal/service"
)

func New(
	areaSvc service.AreaServiceInterface,
	subareaSvc service.SubareaServiceInterface,
	projectSvc service.ProjectServiceInterface,
	taskSvc service.TaskServiceInterface,
) *tea.Program {
	return tea.NewProgram(
		InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc),
		tea.WithAltScreen(),
	)
}
