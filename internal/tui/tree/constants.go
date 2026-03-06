// Package tree provides hierarchical tree rendering for TUI applications.
// It supports unlimited nesting depth, expand/collapse functionality, and navigation helpers.
// The package is designed to be a pure logic component with no database dependencies.
package tree

// Tree character constants for visual rendering.
// These avoid magic strings throughout the codebase.
const (
	// TreeIndent is the indentation per depth level (2 spaces).
	TreeIndent = "  "
	// TreeBranch is the indentation for non-last child nodes (simple indent).
	TreeBranch = "  "
	// TreeLast is the indentation for the last child node at a level (simple indent).
	TreeLast = "  "
	// TreeVertical is the vertical continuation character (simple indent).
	TreeVertical = "  "
	// ExpandedIcon indicates a node is expanded and its children are visible.
	ExpandedIcon = "▾"
	// CollapsedIcon indicates a node is collapsed and its children are hidden.
	CollapsedIcon = "▸"
	// SelectedIcon marks the currently selected node (reserved for future use).
	SelectedIcon = "► "
	// DefaultSelected is the default selection marker (empty).
	DefaultSelected = ""
)

// TreeStyle defines the visual characters used for tree rendering.
// This allows customization of the tree appearance.
type TreeStyle struct {
	// Branch is the character prefix for non-last child nodes.
	Branch string
	// Last is the character prefix for the last child at a level.
	Last string
	// Vertical is the vertical continuation character for indented levels.
	Vertical string
	// Indent is the spacing per depth level.
	Indent string
}

// DefaultStyle provides the standard tree rendering characters.
var DefaultStyle = TreeStyle{
	Branch:   TreeBranch,
	Last:     TreeLast,
	Vertical: TreeVertical,
	Indent:   TreeIndent,
}
