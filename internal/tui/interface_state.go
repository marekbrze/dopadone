package tui

type StateManager interface {
	GetAreaState(areaID string) *AreaState
	SaveCurrentAreaState()
	RestoreAreaState(areaID string)
	IsEmpty(column FocusColumn) bool
}
