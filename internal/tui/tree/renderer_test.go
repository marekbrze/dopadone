package tree

import (
	"strings"
	"testing"
)

func TestRenderEmpty(t *testing.T) {
	renderer := NewRenderer()
	result := renderer.Render(nil, "")

	if result != "" {
		t.Errorf("expected empty string for nil root, got '%s'", result)
	}
}

func TestRenderSingleNode(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("p1", "Project 1", nil)

	result := renderer.Render(root, "")

	if !strings.Contains(result, "Project 1") {
		t.Errorf("expected result to contain 'Project 1', got '%s'", result)
	}
}

func TestRenderMultiLevel(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child1 := NewTreeNode("c1", "Child 1", nil)
	child2 := NewTreeNode("c2", "Child 2", nil)
	grandchild := NewTreeNode("gc1", "Grandchild", nil)

	root.AddChild(child1)
	root.AddChild(child2)
	child1.AddChild(grandchild)

	result := renderer.Render(root, "")
	lines := strings.Split(result, "\n")

	if len(lines) < 3 {
		t.Errorf("expected at least 3 lines, got %d: %v", len(lines), lines)
	}

	if !strings.Contains(lines[0], "Root") {
		t.Errorf("expected first line to contain 'Root', got '%s'", lines[0])
	}

	if !strings.Contains(result, "Child 1") || !strings.Contains(result, "Child 2") || !strings.Contains(result, "Grandchild") {
		t.Errorf("expected result to contain all node names, got '%s'", result)
	}
}

func TestRenderWithCollapsedNodes(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child := NewTreeNode("c1", "Child", nil)
	grandchild := NewTreeNode("gc1", "Grandchild", nil)

	root.AddChild(child)
	child.AddChild(grandchild)

	child.IsExpanded = false
	result := renderer.Render(root, "")

	if strings.Contains(result, "Grandchild") {
		t.Errorf("expected collapsed node's children to not be rendered, got '%s'", result)
	}

	if !strings.Contains(result, "▸") {
		t.Errorf("expected collapsed indicator '▸', got '%s'", result)
	}
}

func TestRenderWithExpandedNodes(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child := NewTreeNode("c1", "Child", nil)
	grandchild := NewTreeNode("gc1", "Grandchild", nil)

	root.AddChild(child)
	child.AddChild(grandchild)

	child.IsExpanded = true
	result := renderer.Render(root, "")

	if !strings.Contains(result, "Grandchild") {
		t.Errorf("expected expanded node's children to be rendered, got '%s'", result)
	}

	if !strings.Contains(result, "▾") {
		t.Errorf("expected expanded indicator '▾', got '%s'", result)
	}
}

func TestRenderSelectedNode(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child := NewTreeNode("c1", "Child", nil)

	root.AddChild(child)

	result := renderer.Render(root, "c1")

	if !strings.Contains(result, "Child") {
		t.Errorf("expected result to contain 'Child', got '%s'", result)
	}
}

func TestRenderTreeIndicators(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child1 := NewTreeNode("c1", "Child 1", nil)
	child2 := NewTreeNode("c2", "Child 2", nil)

	root.AddChild(child1)
	root.AddChild(child2)

	result := renderer.Render(root, "")

	if strings.Contains(result, "├─") || strings.Contains(result, "└─") || strings.Contains(result, "│") {
		t.Errorf("expected no box-drawing characters in output, got '%s'", result)
	}

	if strings.Contains(result, "  ") {
		t.Logf("output uses simple indentation: %s", result)
	}
}

func TestRenderLeafNodeNoIndicator(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("p1", "Leaf Project", nil)

	result := renderer.Render(root, "")

	if strings.Contains(result, "▸") || strings.Contains(result, "▾") {
		t.Errorf("expected leaf node to have no arrow indicator (▸/▾), got '%s'", result)
	}
}

func TestRenderMultipleRoots(t *testing.T) {
	renderer := NewRenderer()
	dummyRoot := NewTreeNode("", "root", nil)
	root1 := NewTreeNode("r1", "Root 1", nil)
	root2 := NewTreeNode("r2", "Root 2", nil)

	dummyRoot.AddChild(root1)
	dummyRoot.AddChild(root2)

	result := renderer.Render(dummyRoot, "")

	if !strings.Contains(result, "Root 1") || !strings.Contains(result, "Root 2") {
		t.Errorf("expected result to contain both roots, got '%s'", result)
	}

	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines for 2 roots, got %d", len(lines))
	}
}

func TestRenderCompact(t *testing.T) {
	renderer := NewRenderer()
	root := NewTreeNode("root", "Root", nil)
	child := NewTreeNode("c1", "Child", nil)

	root.AddChild(child)

	result := renderer.RenderCompact(root, "")

	if result == "" {
		t.Error("expected non-empty result from RenderCompact")
	}
}

func TestRendererSetStyles(t *testing.T) {
	renderer := NewRenderer()

	newStyle := renderer.SetSelectedStyle(renderer.selectedStyle).
		SetExpandedStyle(renderer.expandedStyle).
		SetCollapsedStyle(renderer.collapsedStyle)

	if newStyle != renderer {
		t.Error("expected Set*Style to return same renderer instance")
	}
}
