package spacemenu

type MenuState int

const (
	StateMain MenuState = iota
	StateConfig
)

type MenuAction int

const (
	ActionNone MenuAction = iota
	ActionQuit
	ActionConfig
	ActionCreateArea
	ActionEditArea
	ActionDeleteArea
)

type Command struct {
	Key         string
	Label       string
	Description string
	Action      MenuAction
}
