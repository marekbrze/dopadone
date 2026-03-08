package tree

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/marekbrze/dopadone/internal/tui/theme"
)

// Renderer produces visual string representations of tree structures.
// It uses lipgloss for styling and supports selected node highlighting.
type Renderer struct {
	style          TreeStyle
	selectedStyle  lipgloss.Style
	expandedStyle  lipgloss.Style
	collapsedStyle lipgloss.Style
}

// NewRenderer creates a new Renderer with default lipgloss styling.
// The default selected style uses bold reverse video.
func NewRenderer() *Renderer {
	return &Renderer{
		style: DefaultStyle,
		selectedStyle: lipgloss.NewStyle().
			Bold(true).
			Reverse(true),
		expandedStyle: lipgloss.NewStyle().
			Foreground(theme.Default.Secondary),
		collapsedStyle: lipgloss.NewStyle().
			Foreground(theme.Default.Muted),
	}
}

// Render produces a string representation of the tree with visual indicators.
// The output includes:
//   - Simple indentation for depth levels
//   - Arrow indicators (▸/▾) for collapsed/expanded non-leaf nodes
//   - Selected node highlighting via lipgloss
//
// The selectedID parameter identifies which node should be highlighted.
// Pass empty string for no selection.
func (r *Renderer) Render(root *TreeNode, selectedID string) string {
	if root == nil {
		return ""
	}

	if root.Name == "root" && root.ID == "" {
		var lines []string
		for i, child := range root.Children {
			isLast := i == len(root.Children)-1
			lines = append(lines, r.renderNode(child, selectedID, []bool{}, isLast)...)
		}
		return strings.Join(lines, "\n")
	}

	lines := r.renderNode(root, selectedID, []bool{}, true)
	return strings.Join(lines, "\n")
}

// renderNode recursively renders a node and its visible children.
func (r *Renderer) renderNode(node *TreeNode, selectedID string, levels []bool, isLast bool) []string {
	var lines []string

	line := r.buildLine(node, levels, isLast, selectedID)
	lines = append(lines, line)

	if node.HasChildren() && node.IsExpanded {
		newLevels := append(levels, !isLast)
		for i, child := range node.Children {
			childIsLast := i == len(node.Children)-1
			childLines := r.renderNode(child, selectedID, newLevels, childIsLast)
			lines = append(lines, childLines...)
		}
	}

	return lines
}

// buildLine constructs a single line of rendered output for a node.
func (r *Renderer) buildLine(node *TreeNode, levels []bool, isLast bool, selectedID string) string {
	var prefix strings.Builder

	for i := 0; i < len(levels); i++ {
		prefix.WriteString(r.style.Indent)
	}

	if !node.IsRoot() || (node.Name != "root" || node.ID != "") {
		prefix.WriteString(r.style.Indent)
	}

	indicator := ""
	if node.HasChildren() {
		if node.IsExpanded {
			indicator = ExpandedIcon + " "
		} else {
			indicator = CollapsedIcon + " "
		}
	} else {
		indicator = "  "
	}

	name := node.Name
	isSelected := node.ID == selectedID
	if isSelected {
		name = r.selectedStyle.Render(name)
	}

	return prefix.String() + indicator + name
}

// RenderCompact produces a compact string representation of the tree.
// Currently an alias for Render, reserved for future compact formatting.
func (r *Renderer) RenderCompact(root *TreeNode, selectedID string) string {
	return r.Render(root, selectedID)
}

// SetSelectedStyle configures the lipgloss style for selected nodes.
// Returns the renderer for method chaining.
func (r *Renderer) SetSelectedStyle(style lipgloss.Style) *Renderer {
	r.selectedStyle = style
	return r
}

// SetExpandedStyle configures the lipgloss style for expanded nodes.
// Returns the renderer for method chaining.
func (r *Renderer) SetExpandedStyle(style lipgloss.Style) *Renderer {
	r.expandedStyle = style
	return r
}

// SetCollapsedStyle configures the lipgloss style for collapsed nodes.
// Returns the renderer for method chaining.
func (r *Renderer) SetCollapsedStyle(style lipgloss.Style) *Renderer {
	r.collapsedStyle = style
	return r
}
