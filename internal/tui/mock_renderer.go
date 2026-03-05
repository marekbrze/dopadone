package tui

type MockRenderer struct {
	RenderSubareasFunc func() string
	RenderProjectsFunc func() string
	RenderTasksFunc    func() string
	RenderFooterFunc   func() string
	RenderToastsFunc   func() string
}

func (m *MockRenderer) RenderSubareas() string {
	if m.RenderSubareasFunc != nil {
		return m.RenderSubareasFunc()
	}
	return ""
}

func (m *MockRenderer) RenderProjects() string {
	if m.RenderProjectsFunc != nil {
		return m.RenderProjectsFunc()
	}
	return ""
}

func (m *MockRenderer) RenderTasks() string {
	if m.RenderTasksFunc != nil {
		return m.RenderTasksFunc()
	}
	return ""
}

func (m *MockRenderer) RenderFooter() string {
	if m.RenderFooterFunc != nil {
		return m.RenderFooterFunc()
	}
	return ""
}

func (m *MockRenderer) RenderToasts() string {
	if m.RenderToastsFunc != nil {
		return m.RenderToastsFunc()
	}
	return ""
}
