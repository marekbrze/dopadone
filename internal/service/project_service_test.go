package service

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/marekbrze/dopadone/internal/db"
	"github.com/marekbrze/dopadone/internal/domain"
)

type mockProjectQuerier struct {
	createProjectFunc                  func(ctx context.Context, arg db.CreateProjectParams) (db.Project, error)
	getProjectByIDFunc                 func(ctx context.Context, id string) (db.Project, error)
	listProjectsBySubareaFunc          func(ctx context.Context, subareaID sql.NullString) ([]db.Project, error)
	listProjectsBySubareaRecursiveFunc func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error)
	listProjectsByParentFunc           func(ctx context.Context, parentID sql.NullString) ([]db.Project, error)
	listAllProjectsFunc                func(ctx context.Context) ([]db.Project, error)
	getProjectsByStatusFunc            func(ctx context.Context, status string) ([]db.Project, error)
	listProjectsByPriorityFunc         func(ctx context.Context, priority string) ([]db.Project, error)
	updateProjectFunc                  func(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error)
	softDeleteProjectFunc              func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error)
	softDeleteTasksByProjectFunc       func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error
	hardDeleteProjectFunc              func(ctx context.Context, id string) error
	countTasksByProjectFunc            func(ctx context.Context, projectID string) (int64, error)
	countProjectsByParentFunc          func(ctx context.Context, parentID sql.NullString) (int64, error)
}

func (m *mockProjectQuerier) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	if m.createProjectFunc != nil {
		return m.createProjectFunc(ctx, arg)
	}
	return db.Project{}, nil
}

func (m *mockProjectQuerier) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	if m.getProjectByIDFunc != nil {
		return m.getProjectByIDFunc(ctx, id)
	}
	return db.Project{}, nil
}

func (m *mockProjectQuerier) ListProjectsBySubarea(ctx context.Context, subareaID sql.NullString) ([]db.Project, error) {
	if m.listProjectsBySubareaFunc != nil {
		return m.listProjectsBySubareaFunc(ctx, subareaID)
	}
	return nil, nil
}

func (m *mockProjectQuerier) ListProjectsBySubareaRecursive(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
	if m.listProjectsBySubareaRecursiveFunc != nil {
		return m.listProjectsBySubareaRecursiveFunc(ctx, subareaID)
	}
	return nil, nil
}

func (m *mockProjectQuerier) ListProjectsByParent(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
	if m.listProjectsByParentFunc != nil {
		return m.listProjectsByParentFunc(ctx, parentID)
	}
	return nil, nil
}

func (m *mockProjectQuerier) ListAllProjects(ctx context.Context) ([]db.Project, error) {
	if m.listAllProjectsFunc != nil {
		return m.listAllProjectsFunc(ctx)
	}
	return nil, nil
}

func (m *mockProjectQuerier) GetProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	if m.getProjectsByStatusFunc != nil {
		return m.getProjectsByStatusFunc(ctx, status)
	}
	return nil, nil
}

func (m *mockProjectQuerier) ListProjectsByPriority(ctx context.Context, priority string) ([]db.Project, error) {
	if m.listProjectsByPriorityFunc != nil {
		return m.listProjectsByPriorityFunc(ctx, priority)
	}
	return nil, nil
}

func (m *mockProjectQuerier) ListProjectsByIDs(ctx context.Context, ids []string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockProjectQuerier) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	if m.updateProjectFunc != nil {
		return m.updateProjectFunc(ctx, arg)
	}
	return db.Project{}, nil
}

func (m *mockProjectQuerier) SoftDeleteProject(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
	if m.softDeleteProjectFunc != nil {
		return m.softDeleteProjectFunc(ctx, arg)
	}
	return db.Project{}, nil
}

func (m *mockProjectQuerier) SoftDeleteTasksByProject(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
	if m.softDeleteTasksByProjectFunc != nil {
		return m.softDeleteTasksByProjectFunc(ctx, arg)
	}
	return nil
}

func (m *mockProjectQuerier) HardDeleteProject(ctx context.Context, id string) error {
	if m.hardDeleteProjectFunc != nil {
		return m.hardDeleteProjectFunc(ctx, id)
	}
	return nil
}

func (m *mockProjectQuerier) CountTasksByProject(ctx context.Context, projectID string) (int64, error) {
	if m.countTasksByProjectFunc != nil {
		return m.countTasksByProjectFunc(ctx, projectID)
	}
	return 0, nil
}

func (m *mockProjectQuerier) CountProjectsByParent(ctx context.Context, parentID sql.NullString) (int64, error) {
	if m.countProjectsByParentFunc != nil {
		return m.countProjectsByParentFunc(ctx, parentID)
	}
	return 0, nil
}

func (m *mockProjectQuerier) CreateArea(ctx context.Context, arg db.CreateAreaParams) (db.CreateAreaRow, error) {
	return db.CreateAreaRow{}, nil
}

func (m *mockProjectQuerier) GetAreaByID(ctx context.Context, id string) (db.GetAreaByIDRow, error) {
	return db.GetAreaByIDRow{}, nil
}

func (m *mockProjectQuerier) ListAreas(ctx context.Context) ([]db.ListAreasRow, error) {
	return nil, nil
}

func (m *mockProjectQuerier) UpdateArea(ctx context.Context, arg db.UpdateAreaParams) (db.UpdateAreaRow, error) {
	return db.UpdateAreaRow{}, nil
}

func (m *mockProjectQuerier) UpdateAreaSortOrder(ctx context.Context, arg db.UpdateAreaSortOrderParams) error {
	return nil
}

func (m *mockProjectQuerier) SoftDeleteArea(ctx context.Context, arg db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
	return db.SoftDeleteAreaRow{}, nil
}

func (m *mockProjectQuerier) HardDeleteArea(ctx context.Context, id string) error {
	return nil
}

func (m *mockProjectQuerier) CountSubareasByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockProjectQuerier) CountProjectsByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockProjectQuerier) CountTasksByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockProjectQuerier) CountProjectsBySubarea(ctx context.Context, subareaID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockProjectQuerier) DeleteTasksByProject(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockProjectQuerier) DeleteProjectsBySubarea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockProjectQuerier) DeleteSubareasByArea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockProjectQuerier) CreateSubarea(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockProjectQuerier) GetSubareaByID(ctx context.Context, id string) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockProjectQuerier) ListSubareasByArea(ctx context.Context, areaID string) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockProjectQuerier) UpdateSubarea(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockProjectQuerier) SoftDeleteSubarea(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockProjectQuerier) HardDeleteSubarea(ctx context.Context, id string) error {
	return nil
}

func (m *mockProjectQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockProjectQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockProjectQuerier) ListTasksByProject(ctx context.Context, projectID string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) ListNextTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) ListTasksByStatus(ctx context.Context, status string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) ListTasksByPriority(ctx context.Context, priority string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockProjectQuerier) SoftDeleteTask(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockProjectQuerier) HardDeleteTask(ctx context.Context, id string) error {
	return nil
}

func (m *mockProjectQuerier) ToggleIsNext(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockProjectQuerier) ListAllTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) ListTasksByProjectRecursive(ctx context.Context, projectID sql.NullString) ([]db.Task, error) {
	return nil, nil
}

func (m *mockProjectQuerier) ListAllSubareas(ctx context.Context) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockProjectQuerier) DeleteProjectsByParentID(ctx context.Context, parentID sql.NullString) error {
	return nil
}

func (m *mockProjectQuerier) DeleteProjectsBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockProjectQuerier) DeleteTasksBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockProjectQuerier) DeleteTasksByProjectID(ctx context.Context, projectID string) error {
	return nil
}

func TestProjectService_Create(t *testing.T) {
	now := time.Now()
	subareaID := "subarea-1"

	tests := []struct {
		name    string
		params  CreateProjectParams
		mock    func() *mockProjectQuerier
		wantErr bool
	}{
		{
			name: "creates project successfully",
			params: CreateProjectParams{
				Name:        "Test Project",
				Description: "Test description",
				Status:      domain.ProjectStatusActive,
				Priority:    domain.PriorityHigh,
				Progress:    domain.Progress(0),
				SubareaID:   &subareaID,
				Position:    0,
			},
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					createProjectFunc: func(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
						return db.Project{
							ID:          arg.ID,
							Name:        arg.Name,
							Description: arg.Description,
							Status:      arg.Status,
							Priority:    arg.Priority,
							Progress:    arg.Progress,
							SubareaID:   arg.SubareaID,
							Position:    arg.Position,
							CreatedAt:   now,
							UpdatedAt:   now,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name: "rejects project without parent",
			params: CreateProjectParams{
				Name:     "Test Project",
				Status:   domain.ProjectStatusActive,
				Priority: domain.PriorityHigh,
				Progress: domain.Progress(0),
			},
			mock:    func() *mockProjectQuerier { return &mockProjectQuerier{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.Create(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ProjectService.Create() returned nil project")
			}
		})
	}
}

func TestProjectService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockProjectQuerier
		wantErr bool
	}{
		{
			name: "retrieves project by ID",
			id:   "project-1",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
						if id == "project-1" {
							return db.Project{
								ID:        "project-1",
								Name:      "Test Project",
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							}, nil
						}
						return db.Project{}, sql.ErrNoRows
					},
				}
			},
			wantErr: false,
		},
		{
			name: "returns error for non-existent project",
			id:   "nonexistent",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
						return db.Project{}, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.GetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ProjectService.GetByID() returned nil project")
			}
		})
	}
}

func TestProjectService_ListAll(t *testing.T) {
	tests := []struct {
		name    string
		mock    func() *mockProjectQuerier
		want    int
		wantErr bool
	}{
		{
			name: "lists all projects",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
						return []db.Project{
							{ID: "project-1", Name: "Project 1", Status: "active", Priority: "high", Progress: 0},
							{ID: "project-2", Name: "Project 2", Status: "active", Priority: "medium", Progress: 50},
						}, nil
					},
				}
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "returns empty list",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listAllProjectsFunc: func(ctx context.Context) ([]db.Project, error) {
						return []db.Project{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.ListAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ProjectService.ListAll() returned %d projects, want %d", len(got), tt.want)
			}
		})
	}
}

func TestProjectService_GetStats(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mock      func() *mockProjectQuerier
		wantStats *ProjectStats
		wantErr   bool
	}{
		{
			name: "returns correct stats",
			id:   "project-1",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					countTasksByProjectFunc: func(ctx context.Context, projectID string) (int64, error) {
						return 5, nil
					},
					countProjectsByParentFunc: func(ctx context.Context, parentID sql.NullString) (int64, error) {
						return 2, nil
					},
				}
			},
			wantStats: &ProjectStats{TaskCount: 5, ProjectCount: 2},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.GetStats(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.GetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.TaskCount != tt.wantStats.TaskCount {
					t.Errorf("TaskCount = %d, want %d", got.TaskCount, tt.wantStats.TaskCount)
				}
				if got.ProjectCount != tt.wantStats.ProjectCount {
					t.Errorf("ProjectCount = %d, want %d", got.ProjectCount, tt.wantStats.ProjectCount)
				}
			}
		})
	}
}

func TestProjectService_ValidateParentHierarchy(t *testing.T) {
	tests := []struct {
		name      string
		parentID  string
		projectID string
		mock      func() *mockProjectQuerier
		wantErr   bool
	}{
		{
			name:      "accepts valid parent",
			parentID:  "parent-1",
			projectID: "project-1",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
						return db.Project{
							ID:        "parent-1",
							Name:      "Parent",
							ParentID:  sql.NullString{Valid: false},
							Status:    "active",
							Priority:  "high",
							Progress:  0,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:      "rejects self-reference",
			parentID:  "project-1",
			projectID: "project-1",
			mock:      func() *mockProjectQuerier { return &mockProjectQuerier{} },
			wantErr:   true,
		},
		{
			name:      "rejects circular reference",
			parentID:  "parent-1",
			projectID: "ancestor-1",
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
						if id == "parent-1" {
							return db.Project{
								ID:       "parent-1",
								ParentID: sql.NullString{String: "ancestor-1", Valid: true},
								Status:   "active",
								Priority: "high",
								Progress: 0,
							}, nil
						}
						return db.Project{}, errors.New("not found")
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			err := svc.ValidateParentHierarchy(context.Background(), tt.parentID, tt.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.ValidateParentHierarchy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProjectService_ListByPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority domain.Priority
		mock     func() *mockProjectQuerier
		want     int
		wantErr  bool
	}{
		{
			name:     "lists projects by high priority",
			priority: domain.PriorityHigh,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsByPriorityFunc: func(ctx context.Context, priority string) ([]db.Project, error) {
						return []db.Project{
							{ID: "project-1", Name: "Project 1", Status: "active", Priority: "high", Progress: 0},
							{ID: "project-2", Name: "Project 2", Status: "active", Priority: "high", Progress: 50},
						}, nil
					},
				}
			},
			want:    2,
			wantErr: false,
		},
		{
			name:     "returns empty list for priority with no projects",
			priority: domain.PriorityLow,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsByPriorityFunc: func(ctx context.Context, priority string) ([]db.Project, error) {
						return []db.Project{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.ListByPriority(context.Background(), tt.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProjectService.ListByPriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("ProjectService.ListByPriority() returned %d projects, want %d", len(got), tt.want)
			}
		})
	}
}

func TestProjectService_ListBySubareaRecursive(t *testing.T) {
	now := time.Now()
	subareaA := "subarea-a"

	tests := []struct {
		name      string
		subareaID string
		mock      func() *mockProjectQuerier
		wantCount int
		wantErr   bool
		wantIDs   []string
	}{
		{
			name:      "empty subareaID returns empty slice",
			subareaID: "",
			mock:      func() *mockProjectQuerier { return &mockProjectQuerier{} },
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "no projects in database",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{}, nil
					},
				}
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "direct membership only",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 1,
			wantIDs:   []string{"proj-root-a"},
		},
		{
			name:      "nested project via parent",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "proj-child-a1",
								ParentID:  sql.NullString{String: "proj-root-a", Valid: true},
								SubareaID: sql.NullString{Valid: false},
								Status:    "active",
								Priority:  "medium",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 2,
			wantIDs:   []string{"proj-root-a", "proj-child-a1"},
		},
		{
			name:      "deep nesting (3 levels)",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "proj-child-a1",
								ParentID:  sql.NullString{String: "proj-root-a", Valid: true},
								Status:    "active",
								Priority:  "medium",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "proj-grandchild-a1",
								ParentID:  sql.NullString{String: "proj-child-a1", Valid: true},
								Status:    "active",
								Priority:  "low",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 3,
		},
		{
			name:      "excludes projects in other subareas",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 1,
			wantIDs:   []string{"proj-root-a"},
		},
		{
			name:      "excludes soft-deleted projects",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-active",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 1,
		},
		{
			name:      "mixed direct and nested",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a1",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "proj-root-a2",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "medium",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "proj-child-a1",
								ParentID:  sql.NullString{String: "proj-root-a1", Valid: true},
								Status:    "active",
								Priority:  "low",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 3,
		},
		{
			name:      "orphaned project (parent doesn't exist)",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{}, nil
					},
				}
			},
			wantCount: 0,
		},
		{
			name:      "root project with no parent",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "proj-root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								ParentID:  sql.NullString{Valid: false},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 1,
		},
		{
			name:      "database error",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return nil, errors.New("database connection failed")
					},
				}
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name:      "complex hierarchy",
			subareaID: subareaA,
			mock: func() *mockProjectQuerier {
				return &mockProjectQuerier{
					listProjectsBySubareaRecursiveFunc: func(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
						return []db.ListProjectsBySubareaRecursiveRow{
							projectToRow(db.Project{
								ID:        "root-a",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "child-a",
								ParentID:  sql.NullString{String: "root-a", Valid: true},
								Status:    "active",
								Priority:  "medium",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "grandchild-a",
								ParentID:  sql.NullString{String: "child-a", Valid: true},
								Status:    "active",
								Priority:  "low",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
							projectToRow(db.Project{
								ID:        "root-a2",
								SubareaID: sql.NullString{String: subareaA, Valid: true},
								Status:    "active",
								Priority:  "high",
								Progress:  0,
								CreatedAt: now,
								UpdatedAt: now,
							}),
						}, nil
					},
				}
			},
			wantCount: 4,
			wantIDs:   []string{"root-a", "child-a", "grandchild-a", "root-a2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewProjectService(tt.mock(), nil)
			got, err := svc.ListBySubareaRecursive(context.Background(), tt.subareaID)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListBySubareaRecursive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.wantCount {
				t.Errorf("ListBySubareaRecursive() returned %d projects, want %d", len(got), tt.wantCount)
				t.Logf("Returned IDs: %v", getProjectIDs(got))
			}

			if tt.wantIDs != nil && !tt.wantErr {
				gotIDs := getProjectIDs(got)
				for _, wantID := range tt.wantIDs {
					found := false
					for _, gotID := range gotIDs {
						if gotID == wantID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected project ID %s not found in results", wantID)
					}
				}
			}
		})
	}
}

func getProjectIDs(projects []domain.Project) []string {
	ids := make([]string, len(projects))
	for i, p := range projects {
		ids[i] = p.ID
	}
	return ids
}

func projectToRow(p db.Project) db.ListProjectsBySubareaRecursiveRow {
	return db.ListProjectsBySubareaRecursiveRow{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Goal:        p.Goal,
		Status:      p.Status,
		Priority:    p.Priority,
		Progress:    p.Progress,
		Deadline:    p.Deadline,
		Color:       p.Color,
		ParentID:    p.ParentID,
		SubareaID:   p.SubareaID,
		Position:    p.Position,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		CompletedAt: p.CompletedAt,
		DeletedAt:   p.DeletedAt,
	}
}

func projectsToRows(projects []db.Project) []db.ListProjectsBySubareaRecursiveRow {
	rows := make([]db.ListProjectsBySubareaRecursiveRow, len(projects))
	for i, p := range projects {
		rows[i] = projectToRow(p)
	}
	return rows
}

func TestProjectService_SoftDeleteWithCascade(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name        string
		projectID   string
		setupMock   func(*mockProjectQuerier)
		wantErr     error
		verifyCalls func(t *testing.T, mock *mockProjectQuerier)
	}{
		{
			name:      "single project with no children",
			projectID: "proj-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{
						ID:        "proj-1",
						Name:      "Test Project",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				}
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					return []db.Project{}, nil
				}
				m.softDeleteTasksByProjectFunc = func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
					return nil
				}
				m.softDeleteProjectFunc = func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
					return db.Project{ID: "proj-1", DeletedAt: arg.DeletedAt}, nil
				}
			},
			wantErr: nil,
		},
		{
			name:      "project with direct children",
			projectID: "parent-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{
						ID:        "parent-1",
						Name:      "Parent",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				}
				callCount := 0
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					callCount++
					if parentID.String == "parent-1" && callCount == 1 {
						return []db.Project{
							{ID: "child-1", Name: "Child 1"},
							{ID: "child-2", Name: "Child 2"},
						}, nil
					}
					return []db.Project{}, nil
				}
				m.softDeleteTasksByProjectFunc = func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
					return nil
				}
				m.softDeleteProjectFunc = func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
					return db.Project{ID: arg.ID, DeletedAt: arg.DeletedAt}, nil
				}
			},
			wantErr: nil,
		},
		{
			name:      "deeply nested hierarchy",
			projectID: "level-0",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{
						ID:        "level-0",
						Name:      "Level 0",
						CreatedAt: now,
						UpdatedAt: now,
					}, nil
				}
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					switch parentID.String {
					case "level-0":
						return []db.Project{{ID: "level-1", Name: "Level 1"}}, nil
					case "level-1":
						return []db.Project{{ID: "level-2", Name: "Level 2"}}, nil
					default:
						return []db.Project{}, nil
					}
				}
				m.softDeleteTasksByProjectFunc = func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
					return nil
				}
				m.softDeleteProjectFunc = func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
					return db.Project{ID: arg.ID, DeletedAt: arg.DeletedAt}, nil
				}
			},
			wantErr: nil,
		},
		{
			name:      "already deleted project",
			projectID: "proj-deleted",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{
						ID:        "proj-deleted",
						Name:      "Deleted Project",
						DeletedAt: &now,
					}, nil
				}
			},
			wantErr: nil,
		},
		{
			name:      "non-existent project",
			projectID: "nonexistent",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{}, sql.ErrNoRows
				}
			},
			wantErr: ErrProjectNotFound,
		},
		{
			name:      "database error on get project",
			projectID: "proj-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{}, errors.New("database connection failed")
				}
			},
			wantErr: errors.New("get project proj-1"),
		},
		{
			name:      "database error on list children",
			projectID: "proj-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{ID: "proj-1", Name: "Test"}, nil
				}
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr: errors.New("list child projects"),
		},
		{
			name:      "database error on delete tasks",
			projectID: "proj-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{ID: "proj-1", Name: "Test"}, nil
				}
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					return []db.Project{}, nil
				}
				m.softDeleteTasksByProjectFunc = func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
					return errors.New("delete tasks failed")
				}
			},
			wantErr: errors.New("soft delete tasks"),
		},
		{
			name:      "database error on delete project",
			projectID: "proj-1",
			setupMock: func(m *mockProjectQuerier) {
				m.getProjectByIDFunc = func(ctx context.Context, id string) (db.Project, error) {
					return db.Project{ID: "proj-1", Name: "Test"}, nil
				}
				m.listProjectsByParentFunc = func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
					return []db.Project{}, nil
				}
				m.softDeleteTasksByProjectFunc = func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
					return nil
				}
				m.softDeleteProjectFunc = func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
					return db.Project{}, errors.New("delete project failed")
				}
			},
			wantErr: errors.New("soft delete project"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockProjectQuerier{}
			tt.setupMock(mock)

			service := NewProjectService(mock, nil)

			err := service.SoftDeleteWithCascade(ctx, tt.projectID)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr.Error())
					return
				}
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {
					t.Errorf("error = %v, want error containing %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.verifyCalls != nil {
				tt.verifyCalls(t, mock)
			}
		})
	}
}

func TestProjectService_SoftDeleteWithCascade_WithTransaction(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("with transaction manager", func(t *testing.T) {
		mock := &mockProjectQuerier{
			getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
				return db.Project{
					ID:        "proj-1",
					Name:      "Test Project",
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			listProjectsByParentFunc: func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
				return []db.Project{}, nil
			},
			softDeleteTasksByProjectFunc: func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
				return nil
			},
			softDeleteProjectFunc: func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
				return db.Project{ID: arg.ID, DeletedAt: arg.DeletedAt}, nil
			},
		}

		service := NewProjectService(mock, nil)

		err := service.SoftDeleteWithCascade(ctx, "proj-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("without transaction manager", func(t *testing.T) {
		mock := &mockProjectQuerier{
			getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
				return db.Project{
					ID:        "proj-1",
					Name:      "Test Project",
					CreatedAt: now,
					UpdatedAt: now,
				}, nil
			},
			listProjectsByParentFunc: func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
				return []db.Project{}, nil
			},
			softDeleteTasksByProjectFunc: func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
				return nil
			},
			softDeleteProjectFunc: func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
				return db.Project{ID: arg.ID, DeletedAt: arg.DeletedAt}, nil
			},
		}

		service := NewProjectService(mock, nil)

		err := service.SoftDeleteWithCascade(ctx, "proj-1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestProjectService_SoftDeleteWithCascade_VerifyOrder(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	var deleteOrder []string
	mock := &mockProjectQuerier{
		getProjectByIDFunc: func(ctx context.Context, id string) (db.Project, error) {
			return db.Project{
				ID:        "parent",
				Name:      "Parent",
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		},
		listProjectsByParentFunc: func(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
			if parentID.String == "parent" {
				return []db.Project{
					{ID: "child-1", Name: "Child 1"},
					{ID: "child-2", Name: "Child 2"},
				}, nil
			}
			return []db.Project{}, nil
		},
		softDeleteTasksByProjectFunc: func(ctx context.Context, arg db.SoftDeleteTasksByProjectParams) error {
			deleteOrder = append(deleteOrder, "tasks-"+arg.ProjectID)
			return nil
		},
		softDeleteProjectFunc: func(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
			deleteOrder = append(deleteOrder, "project-"+arg.ID)
			return db.Project{ID: arg.ID, DeletedAt: arg.DeletedAt}, nil
		},
	}

	service := NewProjectService(mock, nil)

	err := service.SoftDeleteWithCascade(ctx, "parent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOrder := []string{
		"tasks-child-1",
		"project-child-1",
		"tasks-child-2",
		"project-child-2",
		"tasks-parent",
		"project-parent",
	}

	if !reflect.DeepEqual(deleteOrder, expectedOrder) {
		t.Errorf("delete order = %v, want %v", deleteOrder, expectedOrder)
	}
}
