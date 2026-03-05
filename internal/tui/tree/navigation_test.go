package tree

import "testing"

func buildTestTree() *TreeNode {
	root := NewTreeNode("root", "Root", nil)
	child1 := NewTreeNode("c1", "Child 1", nil)
	child2 := NewTreeNode("c2", "Child 2", nil)
	grandchild1 := NewTreeNode("gc1", "Grandchild 1", nil)
	grandchild2 := NewTreeNode("gc2", "Grandchild 2", nil)

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild1)
	child1.AddChild(grandchild2)

	return root
}

func TestGetNextVisibleNode(t *testing.T) {
	root := buildTestTree()
	child1 := root.Children[0]
	child2 := root.Children[1]
	grandchild1 := child1.Children[0]
	grandchild2 := child1.Children[1]

	tests := []struct {
		name     string
		current  *TreeNode
		expected string
	}{
		{"root to first child", root, "c1"},
		{"first child to grandchild", child1, "gc1"},
		{"grandchild to sibling", grandchild1, "gc2"},
		{"last grandchild to next parent", grandchild2, "c2"},
		{"last node returns nil", child2, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next := GetNextVisibleNode(root, tt.current)
			if tt.expected == "" {
				if next != nil {
					t.Errorf("expected nil, got '%s'", next.ID)
				}
			} else {
				if next == nil {
					t.Errorf("expected '%s', got nil", tt.expected)
				} else if next.ID != tt.expected {
					t.Errorf("expected '%s', got '%s'", tt.expected, next.ID)
				}
			}
		})
	}
}

func TestGetPrevVisibleNode(t *testing.T) {
	root := buildTestTree()
	child1 := root.Children[0]
	grandchild2 := child1.Children[1]

	tests := []struct {
		name     string
		current  *TreeNode
		expected string
	}{
		{"grandchild to sibling", grandchild2, "gc1"},
		{"first grandchild to parent", child1.Children[0], "c1"},
		{"first child to root", child1, "root"},
		{"root returns nil", root, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prev := GetPrevVisibleNode(root, tt.current)
			if tt.expected == "" {
				if prev != nil {
					t.Errorf("expected nil, got '%s'", prev.ID)
				}
			} else {
				if prev == nil {
					t.Errorf("expected '%s', got nil", tt.expected)
				} else if prev.ID != tt.expected {
					t.Errorf("expected '%s', got '%s'", tt.expected, prev.ID)
				}
			}
		})
	}
}

func TestNavigationSkipsCollapsed(t *testing.T) {
	root := buildTestTree()
	child1 := root.Children[0]

	child1.IsExpanded = false

	next := GetNextVisibleNode(root, child1)
	if next == nil {
		t.Error("expected non-nil next node")
	} else if next.ID != "c2" {
		t.Errorf("expected next to be 'c2' (skipping collapsed children), got '%s'", next.ID)
	}

	visible := GetAllVisibleNodes(root)
	for _, node := range visible {
		if node.ID == "gc1" || node.ID == "gc2" {
			t.Error("expected collapsed children to not be visible")
		}
	}
}

func TestNavigationBoundaries(t *testing.T) {
	root := buildTestTree()

	next := GetNextVisibleNode(root, root.Children[1])
	if next != nil {
		t.Errorf("expected nil at end, got '%s'", next.ID)
	}

	prev := GetPrevVisibleNode(root, root)
	if prev != nil {
		t.Errorf("expected nil at start, got '%s'", prev.ID)
	}

	last := GetLastVisibleNode(root)
	if last == nil {
		t.Error("expected non-nil last node")
	} else if last.ID != "c2" {
		t.Errorf("expected last to be 'c2', got '%s'", last.ID)
	}
}

func TestGetAllVisibleNodes(t *testing.T) {
	root := buildTestTree()

	visible := GetAllVisibleNodes(root)

	if len(visible) != 5 {
		t.Errorf("expected 5 visible nodes, got %d", len(visible))
	}

	expectedOrder := []string{"root", "c1", "gc1", "gc2", "c2"}
	for i, expected := range expectedOrder {
		if visible[i].ID != expected {
			t.Errorf("expected visible[%d] to be '%s', got '%s'", i, expected, visible[i].ID)
		}
	}
}

func TestGetAllVisibleNodesCollapsed(t *testing.T) {
	root := buildTestTree()
	root.Children[0].IsExpanded = false

	visible := GetAllVisibleNodes(root)

	if len(visible) != 3 {
		t.Errorf("expected 3 visible nodes (collapsed), got %d", len(visible))
	}

	for _, node := range visible {
		if node.ID == "gc1" || node.ID == "gc2" {
			t.Error("expected collapsed children to not be visible")
		}
	}
}

func TestGetAllVisibleNodesNil(t *testing.T) {
	visible := GetAllVisibleNodes(nil)
	if visible != nil {
		t.Errorf("expected nil for nil root, got %v", visible)
	}
}

func TestFindNodeByID(t *testing.T) {
	root := buildTestTree()

	found := FindNodeByID(root, "gc1")
	if found == nil {
		t.Error("expected to find node 'gc1'")
	} else if found.ID != "gc1" {
		t.Errorf("expected ID 'gc1', got '%s'", found.ID)
	}

	notFound := FindNodeByID(root, "nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent ID")
	}

	nilFound := FindNodeByID(nil, "any")
	if nilFound != nil {
		t.Error("expected nil for nil root")
	}
}

func TestGetFirstVisibleNode(t *testing.T) {
	root := buildTestTree()

	first := GetFirstVisibleNode(root)
	if first == nil {
		t.Error("expected non-nil first node")
	} else if first.ID != "root" {
		t.Errorf("expected first to be 'root', got '%s'", first.ID)
	}

	nilFirst := GetFirstVisibleNode(nil)
	if nilFirst != nil {
		t.Error("expected nil for nil root")
	}
}

func TestGetLastVisibleNode(t *testing.T) {
	root := buildTestTree()

	last := GetLastVisibleNode(root)
	if last == nil {
		t.Error("expected non-nil last node")
	} else if last.ID != "c2" {
		t.Errorf("expected last to be 'c2', got '%s'", last.ID)
	}

	nilLast := GetLastVisibleNode(nil)
	if nilLast != nil {
		t.Error("expected nil for nil root")
	}
}

func TestGetVisibleNodeCount(t *testing.T) {
	root := buildTestTree()

	count := GetVisibleNodeCount(root)
	if count != 5 {
		t.Errorf("expected 5 visible nodes, got %d", count)
	}

	root.Children[0].IsExpanded = false
	collapsedCount := GetVisibleNodeCount(root)
	if collapsedCount != 3 {
		t.Errorf("expected 3 visible nodes (collapsed), got %d", collapsedCount)
	}
}

func TestExpandAll(t *testing.T) {
	root := buildTestTree()
	root.Children[0].IsExpanded = false

	ExpandAll(root)

	if !root.IsExpanded {
		t.Error("expected root to be expanded")
	}
	if !root.Children[0].IsExpanded {
		t.Error("expected child to be expanded")
	}
	if !root.Children[0].Children[0].IsExpanded {
		t.Error("expected grandchild to be expanded")
	}
}

func TestCollapseAll(t *testing.T) {
	root := buildTestTree()

	CollapseAll(root)

	if !root.IsExpanded {
		t.Error("expected root to remain expanded")
	}
	if root.Children[0].IsExpanded {
		t.Error("expected child to be collapsed")
	}
	if root.Children[0].Children[0].IsExpanded {
		t.Error("expected grandchild to be collapsed")
	}
}

func TestExpandToNode(t *testing.T) {
	root := buildTestTree()
	CollapseAll(root)

	grandchild := root.Children[0].Children[0]

	result := ExpandToNode(root, grandchild)

	if !result {
		t.Error("expected ExpandToNode to return true")
	}
	if !root.IsExpanded {
		t.Error("expected root to be expanded")
	}
	if !root.Children[0].IsExpanded {
		t.Error("expected parent to be expanded")
	}
}

func TestExpandToNodeNotFound(t *testing.T) {
	root := buildTestTree()

	nonexistent := NewTreeNode("notfound", "Not Found", nil)

	result := ExpandToNode(root, nonexistent)

	if result {
		t.Error("expected ExpandToNode to return false for nonexistent node")
	}
}
