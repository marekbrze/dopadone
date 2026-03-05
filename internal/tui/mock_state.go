package tui

type MockStateManager struct {
	GetAreaStateFunc         func(areaID string) *AreaState
	SaveCurrentAreaStateFunc func()
	RestoreAreaStateFunc     func(areaID string)
	IsEmptyFunc              func(column FocusColumn) bool
}

func (m *MockStateManager) GetAreaState(areaID string) *AreaState {
	if m.GetAreaStateFunc != nil {
		return m.GetAreaStateFunc(areaID)
	}
	return nil
}

func (m *MockStateManager) SaveCurrentAreaState() {
	if m.SaveCurrentAreaStateFunc != nil {
		m.SaveCurrentAreaStateFunc()
	}
}

func (m *MockStateManager) RestoreAreaState(areaID string) {
	if m.RestoreAreaStateFunc != nil {
		m.RestoreAreaStateFunc(areaID)
	}
}

func (m *MockStateManager) IsEmpty(column FocusColumn) bool {
	if m.IsEmptyFunc != nil {
		return m.IsEmptyFunc(column)
	}
	return false
}
