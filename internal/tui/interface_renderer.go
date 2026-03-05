package tui

type Renderer interface {
	RenderSubareas() string
	RenderProjects() string
	RenderTasks() string
	RenderFooter() string
	RenderToasts() string
}
