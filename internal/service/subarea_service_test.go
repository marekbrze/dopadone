package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/example/dopadone/internal/db"
	"github.com/example/dopadone/internal/domain"
)

type mockSubareaQuerier struct {
	createSubareaFunc          func(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error)
	getSubareaByIDFunc         func(ctx context.Context, id string) (db.Subarea, error)
	listSubareasByAreaFunc     func(ctx context.Context, areaID string) ([]db.Subarea, error)
	listAllSubareasFunc        func(ctx context.Context) ([]db.Subarea, error)
	updateSubareaFunc          func(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error)
	softDeleteSubareaFunc      func(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error)
	hardDeleteSubareaFunc      func(ctx context.Context, id string) error
	countProjectsBySubareaFunc func(ctx context.Context, subareaID sql.NullString) (int64, error)
}

func (m *mockSubareaQuerier) CreateSubarea(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
	if m.createSubareaFunc != nil {
		return m.createSubareaFunc(ctx, arg)
	}
	return db.Subarea{}, nil
}

func (m *mockSubareaQuerier) GetSubareaByID(ctx context.Context, id string) (db.Subarea, error) {
	if m.getSubareaByIDFunc != nil {
		return m.getSubareaByIDFunc(ctx, id)
	}
	return db.Subarea{}, nil
}

func (m *mockSubareaQuerier) ListSubareasByArea(ctx context.Context, areaID string) ([]db.Subarea, error) {
	if m.listSubareasByAreaFunc != nil {
		return m.listSubareasByAreaFunc(ctx, areaID)
	}
	return nil, nil
}

func (m *mockSubareaQuerier) ListAllSubareas(ctx context.Context) ([]db.Subarea, error) {
	if m.listAllSubareasFunc != nil {
		return m.listAllSubareasFunc(ctx)
	}
	return nil, nil
}

func (m *mockSubareaQuerier) UpdateSubarea(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error) {
	if m.updateSubareaFunc != nil {
		return m.updateSubareaFunc(ctx, arg)
	}
	return db.Subarea{}, nil
}

func (m *mockSubareaQuerier) SoftDeleteSubarea(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error) {
	if m.softDeleteSubareaFunc != nil {
		return m.softDeleteSubareaFunc(ctx, arg)
	}
	return db.Subarea{}, nil
}

func (m *mockSubareaQuerier) HardDeleteSubarea(ctx context.Context, id string) error {
	if m.hardDeleteSubareaFunc != nil {
		return m.hardDeleteSubareaFunc(ctx, id)
	}
	return nil
}

func (m *mockSubareaQuerier) CountProjectsBySubarea(ctx context.Context, subareaID sql.NullString) (int64, error) {
	if m.countProjectsBySubareaFunc != nil {
		return m.countProjectsBySubareaFunc(ctx, subareaID)
	}
	return 0, nil
}

func (m *mockSubareaQuerier) CreateArea(ctx context.Context, arg db.CreateAreaParams) (db.CreateAreaRow, error) {
	return db.CreateAreaRow{}, nil
}

func (m *mockSubareaQuerier) GetAreaByID(ctx context.Context, id string) (db.GetAreaByIDRow, error) {
	return db.GetAreaByIDRow{}, nil
}

func (m *mockSubareaQuerier) ListAreas(ctx context.Context) ([]db.ListAreasRow, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) UpdateArea(ctx context.Context, arg db.UpdateAreaParams) (db.UpdateAreaRow, error) {
	return db.UpdateAreaRow{}, nil
}

func (m *mockSubareaQuerier) UpdateAreaSortOrder(ctx context.Context, arg db.UpdateAreaSortOrderParams) error {
	return nil
}

func (m *mockSubareaQuerier) SoftDeleteArea(ctx context.Context, arg db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
	return db.SoftDeleteAreaRow{}, nil
}

func (m *mockSubareaQuerier) HardDeleteArea(ctx context.Context, id string) error {
	return nil
}

func (m *mockSubareaQuerier) CountSubareasByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockSubareaQuerier) CountProjectsByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockSubareaQuerier) CountTasksByArea(ctx context.Context, areaID string) (int64, error) {
	return 0, nil
}

func (m *mockSubareaQuerier) CountProjectsByParent(ctx context.Context, parentID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockSubareaQuerier) CountTasksByProject(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}

func (m *mockSubareaQuerier) DeleteTasksByProject(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockSubareaQuerier) DeleteProjectsBySubarea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockSubareaQuerier) DeleteSubareasByArea(ctx context.Context, areaID string) error {
	return nil
}

func (m *mockSubareaQuerier) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockSubareaQuerier) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockSubareaQuerier) ListProjectsBySubarea(ctx context.Context, subareaID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListProjectsByParent(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListAllProjects(ctx context.Context) ([]db.Project, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListProjectsBySubareaRecursive(ctx context.Context, subareaID sql.NullString) ([]db.ListProjectsBySubareaRecursiveRow, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) GetProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockSubareaQuerier) SoftDeleteProject(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockSubareaQuerier) HardDeleteProject(ctx context.Context, id string) error {
	return nil
}

func (m *mockSubareaQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockSubareaQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockSubareaQuerier) ListTasksByProject(ctx context.Context, projectID string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListNextTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListTasksByStatus(ctx context.Context, status string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListTasksByPriority(ctx context.Context, priority string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockSubareaQuerier) SoftDeleteTask(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockSubareaQuerier) HardDeleteTask(ctx context.Context, id string) error {
	return nil
}

func (m *mockSubareaQuerier) ToggleIsNext(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockSubareaQuerier) ListAllTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListTasksByProjectRecursive(ctx context.Context, projectID sql.NullString) ([]db.Task, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) ListProjectsByPriority(ctx context.Context, priority string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockSubareaQuerier) DeleteProjectsByParentID(ctx context.Context, parentID sql.NullString) error {
	return nil
}

func (m *mockSubareaQuerier) DeleteProjectsBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockSubareaQuerier) DeleteTasksBySubareaID(ctx context.Context, subareaID sql.NullString) error {
	return nil
}

func (m *mockSubareaQuerier) DeleteTasksByProjectID(ctx context.Context, projectID string) error {
	return nil
}

func TestSubareaService_Create(t *testing.T) {
	tests := []struct {
		name    string
		areaID  string
		name_   string
		color   domain.Color
		mock    func() *mockSubareaQuerier
		wantErr bool
	}{
		{
			name:   "creates subarea successfully",
			areaID: "area-1",
			name_:  "Test Subarea",
			color:  domain.Color("#3B82F6"),
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					createSubareaFunc: func(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
						return db.Subarea{
							ID:        arg.ID,
							Name:      arg.Name,
							AreaID:    arg.AreaID,
							Color:     arg.Color,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:    "rejects empty name",
			areaID:  "area-1",
			name_:   "",
			color:   domain.Color("#3B82F6"),
			mock:    func() *mockSubareaQuerier { return &mockSubareaQuerier{} },
			wantErr: true,
		},
		{
			name:    "rejects empty area ID",
			areaID:  "",
			name_:   "Test",
			color:   domain.Color("#3B82F6"),
			mock:    func() *mockSubareaQuerier { return &mockSubareaQuerier{} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(tt.mock(), nil)
			got, err := svc.Create(context.Background(), tt.name_, tt.areaID, tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubareaService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("SubareaService.Create() returned nil subarea")
			}
		})
	}
}

func TestSubareaService_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockSubareaQuerier
		wantErr bool
	}{
		{
			name: "retrieves subarea by ID",
			id:   "subarea-1",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					getSubareaByIDFunc: func(ctx context.Context, id string) (db.Subarea, error) {
						if id == "subarea-1" {
							return db.Subarea{
								ID:        "subarea-1",
								Name:      "Test Subarea",
								AreaID:    "area-1",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							}, nil
						}
						return db.Subarea{}, sql.ErrNoRows
					},
				}
			},
			wantErr: false,
		},
		{
			name: "returns error for non-existent subarea",
			id:   "nonexistent",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					getSubareaByIDFunc: func(ctx context.Context, id string) (db.Subarea, error) {
						return db.Subarea{}, sql.ErrNoRows
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(tt.mock(), nil)
			got, err := svc.GetByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubareaService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("SubareaService.GetByID() returned nil subarea")
			}
		})
	}
}

func TestSubareaService_ListByArea(t *testing.T) {
	tests := []struct {
		name    string
		areaID  string
		mock    func() *mockSubareaQuerier
		want    int
		wantErr bool
	}{
		{
			name:   "lists subareas by area",
			areaID: "area-1",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					listSubareasByAreaFunc: func(ctx context.Context, areaID string) ([]db.Subarea, error) {
						return []db.Subarea{
							{ID: "subarea-1", Name: "Subarea 1", AreaID: "area-1"},
							{ID: "subarea-2", Name: "Subarea 2", AreaID: "area-1"},
						}, nil
					},
				}
			},
			want:    2,
			wantErr: false,
		},
		{
			name:   "returns empty list",
			areaID: "area-2",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					listSubareasByAreaFunc: func(ctx context.Context, areaID string) ([]db.Subarea, error) {
						return []db.Subarea{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(tt.mock(), nil)
			got, err := svc.ListByArea(context.Background(), tt.areaID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubareaService.ListByArea() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("SubareaService.ListByArea() returned %d subareas, want %d", len(got), tt.want)
			}
		})
	}
}

func TestSubareaService_GetStats(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mock      func() *mockSubareaQuerier
		wantStats *SubareaStats
		wantErr   bool
	}{
		{
			name: "returns correct stats",
			id:   "subarea-1",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					countProjectsBySubareaFunc: func(ctx context.Context, subareaID sql.NullString) (int64, error) {
						return 5, nil
					},
				}
			},
			wantStats: &SubareaStats{ProjectCount: 5},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(tt.mock(), nil)
			got, err := svc.GetStats(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubareaService.GetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.ProjectCount != tt.wantStats.ProjectCount {
					t.Errorf("ProjectCount = %d, want %d", got.ProjectCount, tt.wantStats.ProjectCount)
				}
			}
		})
	}
}

func TestSubareaService_GetEffectiveColor(t *testing.T) {
	tests := []struct {
		name       string
		subarea    *domain.Subarea
		parentArea *domain.Area
		wantColor  domain.Color
	}{
		{
			name: "returns subarea color when set",
			subarea: &domain.Subarea{
				ID:     "subarea-1",
				Name:   "Test",
				AreaID: "area-1",
				Color:  domain.Color("#3B82F6"),
			},
			parentArea: &domain.Area{
				ID:    "area-1",
				Name:  "Parent",
				Color: domain.Color("#10B981"),
			},
			wantColor: domain.Color("#3B82F6"),
		},
		{
			name: "inherits from parent when color not set",
			subarea: &domain.Subarea{
				ID:     "subarea-1",
				Name:   "Test",
				AreaID: "area-1",
				Color:  domain.Color(""),
			},
			parentArea: &domain.Area{
				ID:    "area-1",
				Name:  "Parent",
				Color: domain.Color("#10B981"),
			},
			wantColor: domain.Color("#10B981"),
		},
		{
			name: "returns empty when no color available",
			subarea: &domain.Subarea{
				ID:     "subarea-1",
				Name:   "Test",
				AreaID: "area-1",
				Color:  domain.Color(""),
			},
			parentArea: nil,
			wantColor:  domain.Color(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(&mockSubareaQuerier{}, nil)
			got := svc.GetEffectiveColor(context.Background(), tt.subarea, tt.parentArea)
			if got != tt.wantColor {
				t.Errorf("GetEffectiveColor() = %v, want %v", got, tt.wantColor)
			}
		})
	}
}

func TestSubareaService_ListAll(t *testing.T) {
	tests := []struct {
		name    string
		mock    func() *mockSubareaQuerier
		want    int
		wantErr bool
	}{
		{
			name: "lists all subareas",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					listAllSubareasFunc: func(ctx context.Context) ([]db.Subarea, error) {
						return []db.Subarea{
							{ID: "subarea-1", Name: "Subarea 1", AreaID: "area-1"},
							{ID: "subarea-2", Name: "Subarea 2", AreaID: "area-1"},
							{ID: "subarea-3", Name: "Subarea 3", AreaID: "area-2"},
						}, nil
					},
				}
			},
			want:    3,
			wantErr: false,
		},
		{
			name: "returns empty list",
			mock: func() *mockSubareaQuerier {
				return &mockSubareaQuerier{
					listAllSubareasFunc: func(ctx context.Context) ([]db.Subarea, error) {
						return []db.Subarea{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSubareaService(tt.mock(), nil)
			got, err := svc.ListAll(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("SubareaService.ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("SubareaService.ListAll() returned %d subareas, want %d", len(got), tt.want)
			}
		})
	}
}
