package tui

type FocusColumn int

const (
	FocusSubareas FocusColumn = iota
	FocusProjects
	FocusTasks
)

func (f FocusColumn) String() string {
	switch f {
	case FocusSubareas:
		return "Subareas"
	case FocusProjects:
		return "Projects"
	case FocusTasks:
		return "Tasks"
	default:
		return "Unknown"
	}
}

func (f FocusColumn) Prev() FocusColumn {
	switch f {
	case FocusSubareas:
		return FocusTasks
	case FocusProjects:
		return FocusSubareas
	case FocusTasks:
		return FocusProjects
	default:
		return FocusSubareas
	}
}

func (f FocusColumn) Next() FocusColumn {
	switch f {
	case FocusSubareas:
		return FocusProjects
	case FocusProjects:
		return FocusTasks
	case FocusTasks:
		return FocusSubareas
	default:
		return FocusSubareas
	}
}

type AreaState struct {
	SelectedSubareaIndex int
	SelectedProjectIndex int
	SelectedTaskIndex    int
	ExpandedProjects     map[string]bool
	ExpandedTaskGroups   map[string]bool
}

func NewAreaState() *AreaState {
	return &AreaState{
		SelectedSubareaIndex: 0,
		SelectedProjectIndex: 0,
		SelectedTaskIndex:    0,
		ExpandedProjects:     make(map[string]bool),
		ExpandedTaskGroups:   make(map[string]bool),
	}
}
