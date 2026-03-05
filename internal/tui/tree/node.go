package tree

// TreeNode represents a single node in the hierarchical tree structure.
// It supports unlimited nesting depth and expand/collapse state management.
type TreeNode struct {
	// ID is the unique identifier for this node.
	ID string
	// Name is the display text shown for this node.
	Name string
	// Depth indicates the nesting level (0 for root, increments for children).
	Depth int
	// IsExpanded controls whether children are visible in the rendered output.
	IsExpanded bool
	// Children contains the child nodes of this node.
	Children []*TreeNode
	// Parent is a reference to the parent node (nil for root).
	Parent *TreeNode
	// Data holds arbitrary data associated with this node (e.g., *domain.Project).
	Data interface{}
}

// NewTreeNode creates a new tree node with the given ID, name, and associated data.
// The node is created with IsExpanded=true and no children.
func NewTreeNode(id, name string, data interface{}) *TreeNode {
	return &TreeNode{
		ID:         id,
		Name:       name,
		Depth:      0,
		IsExpanded: true,
		Children:   make([]*TreeNode, 0),
		Parent:     nil,
		Data:       data,
	}
}

// IsLeaf returns true if this node has no children.
func (n *TreeNode) IsLeaf() bool {
	return len(n.Children) == 0
}

// HasChildren returns true if this node has one or more children.
func (n *TreeNode) HasChildren() bool {
	return len(n.Children) > 0
}

// ToggleExpanded switches the expanded state between true and false.
// When collapsed, children are not shown in the rendered output.
func (n *TreeNode) ToggleExpanded() {
	n.IsExpanded = !n.IsExpanded
}

// AddChild appends a child node to this node, setting its depth and parent reference.
func (n *TreeNode) AddChild(child *TreeNode) {
	child.Depth = n.Depth + 1
	child.Parent = n
	n.Children = append(n.Children, child)
}

// IsRoot returns true if this node has no parent.
func (n *TreeNode) IsRoot() bool {
	return n.Parent == nil
}
