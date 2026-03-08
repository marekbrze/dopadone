package tree

import (
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/domain"
)

func createTestProject(id, name, parentID, subareaID string, position int) domain.Project {
	var pid *string
	var sid *string
	if parentID != "" {
		pid = &parentID
	}
	if subareaID != "" {
		sid = &subareaID
	}

	return domain.Project{
		ID:        id,
		Name:      name,
		ParentID:  pid,
		SubareaID: sid,
		Position:  position,
		Status:    domain.ProjectStatusActive,
		Priority:  domain.PriorityMedium,
		Progress:  domain.Progress(0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestBuildFromProjectsEmpty(t *testing.T) {
	builder := NewBuilder()
	root := builder.BuildFromProjects([]domain.Project{})

	if root != nil {
		t.Errorf("expected nil root for empty projects, got %+v", root)
	}
}

func TestBuildFromProjectsNil(t *testing.T) {
	builder := NewBuilder()
	root := builder.BuildFromProjects(nil)

	if root != nil {
		t.Errorf("expected nil root for nil projects, got %+v", root)
	}
}

func TestBuildFromProjectsSingleNode(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p1", "Project 1", "", "sa1", 1),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if root.ID != "p1" {
		t.Errorf("expected root ID 'p1', got '%s'", root.ID)
	}
	if root.Name != "Project 1" {
		t.Errorf("expected root Name 'Project 1', got '%s'", root.Name)
	}
	if root.HasChildren() {
		t.Error("expected root to have no children")
	}
}

func TestBuildFromProjectsFlat(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p1", "Project 1", "", "sa1", 1),
		createTestProject("p2", "Project 2", "", "sa1", 2),
		createTestProject("p3", "Project 3", "", "sa1", 3),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if root.Name != "root" {
		t.Errorf("expected dummy root for multiple top-level projects")
	}
	if len(root.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(root.Children))
	}

	expectedOrder := []string{"p1", "p2", "p3"}
	for i, expectedID := range expectedOrder {
		if root.Children[i].ID != expectedID {
			t.Errorf("expected child %d to have ID '%s', got '%s'", i, expectedID, root.Children[i].ID)
		}
	}
}

func TestBuildFromProjectsNested(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p1", "Parent", "", "sa1", 1),
		createTestProject("p2", "Child 1", "p1", "", 1),
		createTestProject("p3", "Child 2", "p1", "", 2),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if root.ID != "p1" {
		t.Errorf("expected root ID 'p1', got '%s'", root.ID)
	}
	if len(root.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(root.Children))
	}
	if root.Children[0].Parent != root {
		t.Error("expected child's Parent to be set")
	}
	if root.Children[0].Depth != 1 {
		t.Errorf("expected child depth 1, got %d", root.Children[0].Depth)
	}
}

func TestBuildFromProjectsDeep(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p1", "Level 1", "", "sa1", 1),
		createTestProject("p2", "Level 2", "p1", "", 1),
		createTestProject("p3", "Level 3", "p2", "", 1),
		createTestProject("p4", "Level 4", "p3", "", 1),
		createTestProject("p5", "Level 5", "p4", "", 1),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}

	current := root
	for i := 0; i < 5; i++ {
		if current.Depth != i {
			t.Errorf("expected depth %d, got %d", i, current.Depth)
		}
		if i < 4 {
			if len(current.Children) != 1 {
				t.Fatalf("expected 1 child at level %d, got %d", i, len(current.Children))
			}
			current = current.Children[0]
		}
	}
}

func TestBuildFromProjectsPositionOrdering(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p3", "Project 3", "", "sa1", 3),
		createTestProject("p1", "Project 1", "", "sa1", 1),
		createTestProject("p2", "Project 2", "", "sa1", 2),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}

	expectedOrder := []string{"p1", "p2", "p3"}
	for i, expectedID := range expectedOrder {
		if root.Children[i].ID != expectedID {
			t.Errorf("expected child %d to have ID '%s' (ordered by position), got '%s'", i, expectedID, root.Children[i].ID)
		}
	}
}

func TestBuildFromProjectsNestedPositionOrdering(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("parent", "Parent", "", "sa1", 1),
		createTestProject("c3", "Child 3", "parent", "", 3),
		createTestProject("c1", "Child 1", "parent", "", 1),
		createTestProject("c2", "Child 2", "parent", "", 2),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if len(root.Children) != 3 {
		t.Errorf("expected 3 children, got %d", len(root.Children))
	}

	expectedOrder := []string{"c1", "c2", "c3"}
	for i, expectedID := range expectedOrder {
		if root.Children[i].ID != expectedID {
			t.Errorf("expected child %d to have ID '%s', got '%s'", i, expectedID, root.Children[i].ID)
		}
	}
}

func TestBuildFromProjectsMultipleRoots(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("r1", "Root 1", "", "sa1", 1),
		createTestProject("r2", "Root 2", "", "sa1", 2),
		createTestProject("c1", "Child of Root 1", "r1", "", 1),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if root.Name != "root" {
		t.Error("expected dummy root for multiple root projects")
	}
	if len(root.Children) != 2 {
		t.Errorf("expected 2 root children, got %d", len(root.Children))
	}
	if len(root.Children[0].Children) != 1 {
		t.Errorf("expected first root to have 1 child")
	}
	if len(root.Children[1].Children) != 0 {
		t.Errorf("expected second root to have no children")
	}
}

func TestBuildFromProjectsComplex(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("root1", "Root 1", "", "sa1", 1),
		createTestProject("root2", "Root 2", "", "sa1", 2),
		createTestProject("child1-1", "Child 1-1", "root1", "", 1),
		createTestProject("child1-2", "Child 1-2", "root1", "", 2),
		createTestProject("child2-1", "Child 2-1", "root2", "", 1),
		createTestProject("grandchild1-1-1", "Grandchild 1-1-1", "child1-1", "", 1),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if len(root.Children) != 2 {
		t.Errorf("expected 2 root children, got %d", len(root.Children))
	}
	if len(root.Children[0].Children) != 2 {
		t.Errorf("expected Root 1 to have 2 children, got %d", len(root.Children[0].Children))
	}
	if len(root.Children[0].Children[0].Children) != 1 {
		t.Errorf("expected Child 1-1 to have 1 grandchild")
	}
}

func TestBuildFromProjectsOrphans(t *testing.T) {
	builder := NewBuilder()
	projects := []domain.Project{
		createTestProject("p1", "Project 1", "", "sa1", 1),
		createTestProject("p2", "Orphan", "nonexistent", "", 1),
	}

	root := builder.BuildFromProjects(projects)

	if root == nil {
		t.Fatal("expected non-nil root")
	}
	if root.ID != "p1" {
		t.Errorf("expected root ID 'p1', got '%s'", root.ID)
	}
	if root.HasChildren() {
		t.Error("expected root to have no children (orphan should not be attached)")
	}
}
