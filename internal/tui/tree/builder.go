package tree

import (
	"sort"

	"github.com/marekbrze/dopadone/internal/domain"
)

// Builder constructs tree structures from domain data.
// It handles the transformation of flat project lists into hierarchical trees.
type Builder struct {
	style TreeStyle
}

// NewBuilder creates a new Builder with default styling.
func NewBuilder() *Builder {
	return &Builder{
		style: DefaultStyle,
	}
}

// BuildFromProjects transforms a flat list of projects into a hierarchical tree.
// Projects are organized by ParentID relationships, with roots being projects
// that have a SubareaID but no ParentID.
//
// The function handles:
//   - Empty lists: returns nil
//   - Single root: returns that node directly
//   - Multiple roots: returns a dummy root node containing all roots
//   - Orphans (children referencing non-existent parents): skipped
//   - Position ordering: siblings are sorted by their Position field
func (b *Builder) BuildFromProjects(projects []domain.Project) *TreeNode {
	if len(projects) == 0 {
		return nil
	}

	childrenMap := make(map[string][]*TreeNode)
	var roots []*TreeNode

	for i := range projects {
		node := NewTreeNode(projects[i].ID, projects[i].Name, &projects[i])

		if projects[i].ParentID != nil {
			childrenMap[*projects[i].ParentID] = append(childrenMap[*projects[i].ParentID], node)
		} else {
			roots = append(roots, node)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return getProjectPosition(roots[i]) < getProjectPosition(roots[j])
	})

	for _, children := range childrenMap {
		sort.Slice(children, func(i, j int) bool {
			return getProjectPosition(children[i]) < getProjectPosition(children[j])
		})
	}

	dummyRoot := NewTreeNode("", RootNodeName, nil)
	for _, root := range roots {
		b.attachChildren(root, childrenMap)
		dummyRoot.AddChild(root)
	}

	if len(dummyRoot.Children) == 1 {
		singleRoot := dummyRoot.Children[0]
		singleRoot.Parent = nil
		singleRoot.Depth = 0
		return singleRoot
	}

	return dummyRoot
}

// attachChildren recursively attaches children to a node based on the children map.
func (b *Builder) attachChildren(node *TreeNode, childrenMap map[string][]*TreeNode) {
	children, exists := childrenMap[node.ID]
	if !exists {
		return
	}

	for _, child := range children {
		node.AddChild(child)
		b.attachChildren(child, childrenMap)
	}
}

// getProjectPosition extracts the Position field from a node's Data if it contains a Project.
func getProjectPosition(node *TreeNode) int {
	if node.Data == nil {
		return 0
	}
	if proj, ok := node.Data.(*domain.Project); ok {
		return proj.Position
	}
	return 0
}
