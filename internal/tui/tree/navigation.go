package tree

// GetNextVisibleNode returns the next visible node after the current node.
// It respects the collapsed state of nodes, skipping hidden children.
// Returns nil if current is the last visible node or if inputs are nil.
func GetNextVisibleNode(root, current *TreeNode) *TreeNode {
	if root == nil || current == nil {
		return nil
	}

	visibleNodes := GetAllVisibleNodes(root)
	for i, node := range visibleNodes {
		if node.ID == current.ID && i < len(visibleNodes)-1 {
			return visibleNodes[i+1]
		}
	}

	return nil
}

// GetPrevVisibleNode returns the previous visible node before the current node.
// It respects the collapsed state of nodes, skipping hidden children.
// Returns nil if current is the first visible node or if inputs are nil.
func GetPrevVisibleNode(root, current *TreeNode) *TreeNode {
	if root == nil || current == nil {
		return nil
	}

	visibleNodes := GetAllVisibleNodes(root)
	for i, node := range visibleNodes {
		if node.ID == current.ID && i > 0 {
			return visibleNodes[i-1]
		}
	}

	return nil
}

// GetAllVisibleNodes returns a flat list of all visible nodes in depth-first order.
// Nodes with collapsed parents are excluded from the result.
// A dummy root (empty ID, name "root") has its children listed directly.
func GetAllVisibleNodes(root *TreeNode) []*TreeNode {
	if root == nil {
		return nil
	}

	var nodes []*TreeNode
	collectVisibleNodes(root, &nodes, root.Name == "root" && root.ID == "")
	return nodes
}

// collectVisibleNodes recursively collects visible nodes into the provided slice.
func collectVisibleNodes(node *TreeNode, nodes *[]*TreeNode, skipRoot bool) {
	if !skipRoot {
		*nodes = append(*nodes, node)
	}

	if node.HasChildren() && node.IsExpanded {
		for _, child := range node.Children {
			collectVisibleNodes(child, nodes, false)
		}
	}
}

// FindNodeByID searches the tree for a node with the given ID.
// Returns nil if not found or if root is nil.
func FindNodeByID(root *TreeNode, id string) *TreeNode {
	if root == nil {
		return nil
	}

	if root.ID == id {
		return root
	}

	for _, child := range root.Children {
		if found := FindNodeByID(child, id); found != nil {
			return found
		}
	}

	return nil
}

// GetFirstVisibleNode returns the first visible node in the tree.
// For a tree with a single root, this is the root itself.
// For a tree with multiple roots, it's the first child of the dummy root.
func GetFirstVisibleNode(root *TreeNode) *TreeNode {
	visibleNodes := GetAllVisibleNodes(root)
	if len(visibleNodes) == 0 {
		return nil
	}
	return visibleNodes[0]
}

// GetLastVisibleNode returns the last visible node in the tree.
// Useful for boundary checks during navigation.
func GetLastVisibleNode(root *TreeNode) *TreeNode {
	visibleNodes := GetAllVisibleNodes(root)
	if len(visibleNodes) == 0 {
		return nil
	}
	return visibleNodes[len(visibleNodes)-1]
}

// GetVisibleNodeCount returns the number of currently visible nodes.
// This count respects the expanded/collapsed state of nodes.
func GetVisibleNodeCount(root *TreeNode) int {
	return len(GetAllVisibleNodes(root))
}

// ExpandAll sets IsExpanded to true for all nodes in the tree.
// This makes all descendants visible.
func ExpandAll(root *TreeNode) {
	if root == nil {
		return
	}

	root.IsExpanded = true
	for _, child := range root.Children {
		ExpandAll(child)
	}
}

// CollapseAll sets IsExpanded to false for all non-root nodes.
// The root node remains expanded to allow navigation back into the tree.
func CollapseAll(root *TreeNode) {
	if root == nil {
		return
	}

	if !root.IsRoot() {
		root.IsExpanded = false
	}
	for _, child := range root.Children {
		CollapseAll(child)
	}
}

// ExpandToNode expands all nodes on the path from root to target.
// This ensures the target node is visible in the rendered output.
// Returns true if the target was found and the path expanded, false otherwise.
func ExpandToNode(root, target *TreeNode) bool {
	if root == nil || target == nil {
		return false
	}

	path := findPath(root, target.ID)
	if path == nil {
		return false
	}

	for _, node := range path[:len(path)-1] {
		node.IsExpanded = true
	}

	return true
}

// findPath returns the path from root to the node with the given ID.
// Returns nil if no such node exists.
func findPath(root *TreeNode, targetID string) []*TreeNode {
	if root == nil {
		return nil
	}

	if root.ID == targetID {
		return []*TreeNode{root}
	}

	for _, child := range root.Children {
		if path := findPath(child, targetID); path != nil {
			return append([]*TreeNode{root}, path...)
		}
	}

	return nil
}
