package tui

const (
	SpinnerIDMain = "main"
)

const (
	EmptyStateNoAreas    = "No areas - press a to add"
	EmptyStateNoSubareas = "No subareas - press s to add"
	EmptyStateNoProjects = "No projects - press p to add"
	EmptyStateNoTasks    = "No tasks - press t to add"
)

const (
	LoadingMessageAreas    = "Loading areas..."
	LoadingMessageSubareas = "Loading subareas..."
	LoadingMessageProjects = "Loading projects..."
	LoadingMessageTasks    = "Loading tasks..."
)

const (
	ErrMsgDatabase  = "Unable to load data. Please restart the application."
	ErrMsgTimeout   = "Loading took too long. Please try again."
	ErrMsgCancelled = "Operation cancelled"
	ErrMsgNotFound  = "Resource not found"
)
