package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/marekbrze/dopadone/internal/db/driver"
	"github.com/marekbrze/dopadone/internal/service"
)

func New(
	areaSvc service.AreaServiceInterface,
	subareaSvc service.SubareaServiceInterface,
	projectSvc service.ProjectServiceInterface,
	taskSvc service.TaskServiceInterface,
	dbDriver driver.DatabaseDriver,
) *tea.Program {
	return tea.NewProgram(
		InitialModel(areaSvc, subareaSvc, projectSvc, taskSvc, dbDriver),
		tea.WithAltScreen(),
	)
}
