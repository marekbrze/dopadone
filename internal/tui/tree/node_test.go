package tree

import "testing"

func TestNewTreeNode(t *testing.T) {
	data := struct{ Value int }{Value: 42}
	node := NewTreeNode("id-1", "Test Node", data)

	if node.ID != "id-1" {
		t.Errorf("expected ID 'id-1', got '%s'", node.ID)
	}
	if node.Name != "Test Node" {
		t.Errorf("expected Name 'Test Node', got '%s'", node.Name)
	}
	if node.Depth != 0 {
		t.Errorf("expected Depth 0, got %d", node.Depth)
	}
	if !node.IsExpanded {
		t.Errorf("expected IsExpanded to be true")
	}
	if len(node.Children) != 0 {
		t.Errorf("expected empty Children slice")
	}
	if node.Parent != nil {
		t.Errorf("expected nil Parent")
	}
	if node.Data != data {
		t.Errorf("expected Data to match input data")
	}
}

func TestTreeNodeIsLeaf(t *testing.T) {
	node := NewTreeNode("id", "name", nil)
	if !node.IsLeaf() {
		t.Error("expected node without children to be leaf")
	}

	child := NewTreeNode("child", "child", nil)
	node.AddChild(child)
	if node.IsLeaf() {
		t.Error("expected node with children to not be leaf")
	}
}

func TestTreeNodeHasChildren(t *testing.T) {
	node := NewTreeNode("id", "name", nil)
	if node.HasChildren() {
		t.Error("expected node without children to return false")
	}

	child := NewTreeNode("child", "child", nil)
	node.AddChild(child)
	if !node.HasChildren() {
		t.Error("expected node with children to return true")
	}
}

func TestTreeNodeToggleExpanded(t *testing.T) {
	node := NewTreeNode("id", "name", nil)
	if !node.IsExpanded {
		t.Error("expected initial IsExpanded to be true")
	}

	node.ToggleExpanded()
	if node.IsExpanded {
		t.Error("expected IsExpanded to be false after toggle")
	}

	node.ToggleExpanded()
	if !node.IsExpanded {
		t.Error("expected IsExpanded to be true after second toggle")
	}
}

func TestTreeNodeAddChild(t *testing.T) {
	parent := NewTreeNode("parent", "Parent", nil)
	child := NewTreeNode("child", "Child", nil)

	parent.AddChild(child)

	if len(parent.Children) != 1 {
		t.Errorf("expected 1 child, got %d", len(parent.Children))
	}
	if parent.Children[0] != child {
		t.Error("expected child to be in parent's Children slice")
	}
	if child.Depth != 1 {
		t.Errorf("expected child depth 1, got %d", child.Depth)
	}
	if child.Parent != parent {
		t.Error("expected child's Parent to be set")
	}
}

func TestTreeNodeIsRoot(t *testing.T) {
	node := NewTreeNode("id", "name", nil)
	if !node.IsRoot() {
		t.Error("expected node without parent to be root")
	}

	parent := NewTreeNode("parent", "Parent", nil)
	parent.AddChild(node)
	if node.IsRoot() {
		t.Error("expected node with parent to not be root")
	}
}

func TestTreeNodeNestedChildren(t *testing.T) {
	root := NewTreeNode("root", "Root", nil)
	child1 := NewTreeNode("child1", "Child 1", nil)
	child2 := NewTreeNode("child2", "Child 2", nil)
	grandchild := NewTreeNode("grandchild", "Grandchild", nil)

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild)

	if child1.Depth != 1 {
		t.Errorf("expected child1 depth 1, got %d", child1.Depth)
	}
	if grandchild.Depth != 2 {
		t.Errorf("expected grandchild depth 2, got %d", grandchild.Depth)
	}
	if len(root.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(root.Children))
	}
	if len(child1.Children) != 1 {
		t.Errorf("expected 1 grandchild, got %d", len(child1.Children))
	}
}
