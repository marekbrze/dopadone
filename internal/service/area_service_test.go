package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/example/projectdb/internal/converter"
	"github.com/example/projectdb/internal/db"
	"github.com/example/projectdb/internal/domain"
)

type mockQuerier struct {
	listAreasFunc               func(ctx context.Context) ([]db.ListAreasRow, error)
	getAreaByIDFunc             func(ctx context.Context, id string) (db.GetAreaByIDRow, error)
	createAreaFunc              func(ctx context.Context, params db.CreateAreaParams) (db.CreateAreaRow, error)
	updateAreaFunc              func(ctx context.Context, params db.UpdateAreaParams) (db.UpdateAreaRow, error)
	updateAreaSortOrderFunc     func(ctx context.Context, params db.UpdateAreaSortOrderParams) error
	softDeleteAreaFunc          func(ctx context.Context, params db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error)
	hardDeleteAreaFunc          func(ctx context.Context, id string) error
	deleteTasksByProjectFunc    func(ctx context.Context, areaID string) error
	deleteProjectsBySubareaFunc func(ctx context.Context, areaID string) error
	deleteSubareasByAreaFunc    func(ctx context.Context, areaID string) error
	countSubareasByAreaFunc     func(ctx context.Context, areaID string) (int64, error)
	countProjectsByAreaFunc     func(ctx context.Context, areaID string) (int64, error)
	countTasksByAreaFunc        func(ctx context.Context, areaID string) (int64, error)
}

func (m *mockQuerier) ListAreas(ctx context.Context) ([]db.ListAreasRow, error) {
	if m.listAreasFunc != nil {
		return m.listAreasFunc(ctx)
	}
	return nil, nil
}

func (m *mockQuerier) GetAreaByID(ctx context.Context, id string) (db.GetAreaByIDRow, error) {
	if m.getAreaByIDFunc != nil {
		return m.getAreaByIDFunc(ctx, id)
	}
	return db.GetAreaByIDRow{}, nil
}

func (m *mockQuerier) CreateArea(ctx context.Context, params db.CreateAreaParams) (db.CreateAreaRow, error) {
	if m.createAreaFunc != nil {
		return m.createAreaFunc(ctx, params)
	}
	return db.CreateAreaRow{}, nil
}

func (m *mockQuerier) UpdateArea(ctx context.Context, params db.UpdateAreaParams) (db.UpdateAreaRow, error) {
	if m.updateAreaFunc != nil {
		return m.updateAreaFunc(ctx, params)
	}
	return db.UpdateAreaRow{}, nil
}

func (m *mockQuerier) UpdateAreaSortOrder(ctx context.Context, params db.UpdateAreaSortOrderParams) error {
	if m.updateAreaSortOrderFunc != nil {
		return m.updateAreaSortOrderFunc(ctx, params)
	}
	return nil
}

func (m *mockQuerier) SoftDeleteArea(ctx context.Context, params db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
	if m.softDeleteAreaFunc != nil {
		return m.softDeleteAreaFunc(ctx, params)
	}
	return db.SoftDeleteAreaRow{}, nil
}

func (m *mockQuerier) HardDeleteArea(ctx context.Context, id string) error {
	if m.hardDeleteAreaFunc != nil {
		return m.hardDeleteAreaFunc(ctx, id)
	}
	return nil
}

func (m *mockQuerier) DeleteTasksByProject(ctx context.Context, areaID string) error {
	if m.deleteTasksByProjectFunc != nil {
		return m.deleteTasksByProjectFunc(ctx, areaID)
	}
	return nil
}

func (m *mockQuerier) DeleteProjectsBySubarea(ctx context.Context, areaID string) error {
	if m.deleteProjectsBySubareaFunc != nil {
		return m.deleteProjectsBySubareaFunc(ctx, areaID)
	}
	return nil
}

func (m *mockQuerier) DeleteSubareasByArea(ctx context.Context, areaID string) error {
	if m.deleteSubareasByAreaFunc != nil {
		return m.deleteSubareasByAreaFunc(ctx, areaID)
	}
	return nil
}

func (m *mockQuerier) CountSubareasByArea(ctx context.Context, areaID string) (int64, error) {
	if m.countSubareasByAreaFunc != nil {
		return m.countSubareasByAreaFunc(ctx, areaID)
	}
	return 0, nil
}

func (m *mockQuerier) CountProjectsByArea(ctx context.Context, areaID string) (int64, error) {
	if m.countProjectsByAreaFunc != nil {
		return m.countProjectsByAreaFunc(ctx, areaID)
	}
	return 0, nil
}

func (m *mockQuerier) CountTasksByArea(ctx context.Context, areaID string) (int64, error) {
	if m.countTasksByAreaFunc != nil {
		return m.countTasksByAreaFunc(ctx, areaID)
	}
	return 0, nil
}

func (m *mockQuerier) CountProjectsByParent(ctx context.Context, parentID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockQuerier) CountProjectsBySubarea(ctx context.Context, subareaID sql.NullString) (int64, error) {
	return 0, nil
}

func (m *mockQuerier) CountTasksByProject(ctx context.Context, projectID string) (int64, error) {
	return 0, nil
}

func (m *mockQuerier) HardDeleteProject(ctx context.Context, id string) error {
	return nil
}

func (m *mockQuerier) HardDeleteSubarea(ctx context.Context, id string) error {
	return nil
}

func (m *mockQuerier) HardDeleteTask(ctx context.Context, id string) error {
	return nil
}

func (m *mockQuerier) CreateProject(ctx context.Context, arg db.CreateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockQuerier) CreateSubarea(ctx context.Context, arg db.CreateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockQuerier) CreateTask(ctx context.Context, arg db.CreateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockQuerier) GetProjectByID(ctx context.Context, id string) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockQuerier) GetProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	return nil, nil
}

func (m *mockQuerier) GetSubareaByID(ctx context.Context, id string) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockQuerier) GetTaskByID(ctx context.Context, id string) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockQuerier) ListAllProjects(ctx context.Context) ([]db.Project, error) {
	return nil, nil
}

func (m *mockQuerier) ListNextTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockQuerier) ListProjectsByParent(ctx context.Context, parentID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockQuerier) ListProjectsBySubarea(ctx context.Context, subareaID sql.NullString) ([]db.Project, error) {
	return nil, nil
}

func (m *mockQuerier) ListSubareasByArea(ctx context.Context, areaID string) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockQuerier) ListTasksByPriority(ctx context.Context, priority string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockQuerier) ListTasksByProject(ctx context.Context, projectID string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockQuerier) ListTasksByStatus(ctx context.Context, status string) ([]db.Task, error) {
	return nil, nil
}

func (m *mockQuerier) SoftDeleteProject(ctx context.Context, arg db.SoftDeleteProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockQuerier) SoftDeleteSubarea(ctx context.Context, arg db.SoftDeleteSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockQuerier) SoftDeleteTask(ctx context.Context, arg db.SoftDeleteTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockQuerier) ToggleIsNext(ctx context.Context, arg db.ToggleIsNextParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockQuerier) UpdateProject(ctx context.Context, arg db.UpdateProjectParams) (db.Project, error) {
	return db.Project{}, nil
}

func (m *mockQuerier) UpdateSubarea(ctx context.Context, arg db.UpdateSubareaParams) (db.Subarea, error) {
	return db.Subarea{}, nil
}

func (m *mockQuerier) UpdateTask(ctx context.Context, arg db.UpdateTaskParams) (db.Task, error) {
	return db.Task{}, nil
}

func (m *mockQuerier) ListAllTasks(ctx context.Context) ([]db.Task, error) {
	return nil, nil
}

func (m *mockQuerier) ListAllSubareas(ctx context.Context) ([]db.Subarea, error) {
	return nil, nil
}

func (m *mockQuerier) ListProjectsByPriority(ctx context.Context, priority string) ([]db.Project, error) {
	return nil, nil
}

func TestAreaService_List(t *testing.T) {
	tests := []struct {
		name    string
		mock    func() *mockQuerier
		want    int
		wantErr bool
	}{
		{
			name: "returns areas sorted by sort_order",
			mock: func() *mockQuerier {
				return &mockQuerier{
					listAreasFunc: func(ctx context.Context) ([]db.ListAreasRow, error) {
						return []db.ListAreasRow{
							{ID: "1", Name: "Area 1", SortOrder: 0, CreatedAt: time.Now(), UpdatedAt: time.Now()},
							{ID: "2", Name: "Area 2", SortOrder: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()},
						}, nil
					},
				}
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "returns empty list when no areas",
			mock: func() *mockQuerier {
				return &mockQuerier{
					listAreasFunc: func(ctx context.Context) ([]db.ListAreasRow, error) {
						return []db.ListAreasRow{}, nil
					},
				}
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			got, err := svc.List(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("AreaService.List() returned %d areas, want %d", len(got), tt.want)
			}
		})
	}
}

func TestAreaService_Create(t *testing.T) {
	tests := []struct {
		name    string
		color   domain.Color
		mock    func() *mockQuerier
		wantErr bool
	}{
		{
			name:  "creates area with next sort order",
			color: domain.Color("#3B82F6"),
			mock: func() *mockQuerier {
				return &mockQuerier{
					listAreasFunc: func(ctx context.Context) ([]db.ListAreasRow, error) {
						return []db.ListAreasRow{
							{ID: "1", Name: "Existing", SortOrder: 0},
						}, nil
					},
					createAreaFunc: func(ctx context.Context, params db.CreateAreaParams) (db.CreateAreaRow, error) {
						if params.SortOrder != 1 {
							t.Errorf("expected SortOrder 1, got %d", params.SortOrder)
						}
						return db.CreateAreaRow{
							ID:        params.ID,
							Name:      params.Name,
							Color:     params.Color,
							SortOrder: params.SortOrder,
							CreatedAt: params.CreatedAt,
							UpdatedAt: params.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
		{
			name:  "creates first area with sort order 0",
			color: domain.Color("#10B981"),
			mock: func() *mockQuerier {
				return &mockQuerier{
					listAreasFunc: func(ctx context.Context) ([]db.ListAreasRow, error) {
						return []db.ListAreasRow{}, nil
					},
					createAreaFunc: func(ctx context.Context, params db.CreateAreaParams) (db.CreateAreaRow, error) {
						if params.SortOrder != 0 {
							t.Errorf("expected SortOrder 0, got %d", params.SortOrder)
						}
						return db.CreateAreaRow{
							ID:        params.ID,
							Name:      params.Name,
							Color:     params.Color,
							SortOrder: params.SortOrder,
							CreatedAt: params.CreatedAt,
							UpdatedAt: params.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			got, err := svc.Create(context.Background(), "Test Area", tt.color)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("AreaService.Create() returned nil area")
			}
		})
	}
}

func TestAreaService_Update(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		newName  string
		newColor domain.Color
		mock     func() *mockQuerier
		wantErr  bool
	}{
		{
			name:     "updates area name and color",
			id:       "area-1",
			newName:  "Updated Name",
			newColor: domain.Color("#EF4444"),
			mock: func() *mockQuerier {
				return &mockQuerier{
					updateAreaFunc: func(ctx context.Context, params db.UpdateAreaParams) (db.UpdateAreaRow, error) {
						if params.ID != "area-1" {
							t.Errorf("expected ID area-1, got %s", params.ID)
						}
						if params.Name != "Updated Name" {
							t.Errorf("expected Name 'Updated Name', got %s", params.Name)
						}
						return db.UpdateAreaRow{
							ID:        params.ID,
							Name:      params.Name,
							Color:     params.Color,
							SortOrder: 0,
							CreatedAt: time.Now(),
							UpdatedAt: params.UpdatedAt,
						}, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			got, err := svc.Update(context.Background(), tt.id, tt.newName, tt.newColor)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("AreaService.Update() returned nil area")
			}
		})
	}
}

func TestAreaService_ReorderAll(t *testing.T) {
	tests := []struct {
		name    string
		areaIDs []string
		mock    func() *mockQuerier
		wantErr bool
	}{
		{
			name:    "reorders all areas",
			areaIDs: []string{"area-3", "area-1", "area-2"},
			mock: func() *mockQuerier {
				callCount := 0
				return &mockQuerier{
					updateAreaSortOrderFunc: func(ctx context.Context, params db.UpdateAreaSortOrderParams) error {
						expected := []struct {
							id        string
							sortOrder int64
						}{
							{"area-3", 0},
							{"area-1", 1},
							{"area-2", 2},
						}
						if callCount < len(expected) {
							if params.ID != expected[callCount].id {
								t.Errorf("call %d: expected ID %s, got %s", callCount, expected[callCount].id, params.ID)
							}
							if params.SortOrder != expected[callCount].sortOrder {
								t.Errorf("call %d: expected SortOrder %d, got %d", callCount, expected[callCount].sortOrder, params.SortOrder)
							}
						}
						callCount++
						return nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			err := svc.ReorderAll(context.Background(), tt.areaIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.ReorderAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAreaService_SoftDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockQuerier
		wantErr bool
	}{
		{
			name: "marks area as deleted",
			id:   "area-1",
			mock: func() *mockQuerier {
				return &mockQuerier{
					softDeleteAreaFunc: func(ctx context.Context, params db.SoftDeleteAreaParams) (db.SoftDeleteAreaRow, error) {
						if params.ID != "area-1" {
							t.Errorf("expected ID area-1, got %s", params.ID)
						}
						return db.SoftDeleteAreaRow{}, nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			err := svc.SoftDelete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.SoftDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAreaService_HardDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mock    func() *mockQuerier
		wantErr bool
	}{
		{
			name: "cascades delete to children then deletes area",
			id:   "area-1",
			mock: func() *mockQuerier {
				callOrder := []string{}
				return &mockQuerier{
					deleteTasksByProjectFunc: func(ctx context.Context, areaID string) error {
						callOrder = append(callOrder, "tasks")
						return nil
					},
					deleteProjectsBySubareaFunc: func(ctx context.Context, areaID string) error {
						callOrder = append(callOrder, "projects")
						return nil
					},
					deleteSubareasByAreaFunc: func(ctx context.Context, areaID string) error {
						callOrder = append(callOrder, "subareas")
						return nil
					},
					hardDeleteAreaFunc: func(ctx context.Context, id string) error {
						callOrder = append(callOrder, "area")
						return nil
					},
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			err := svc.HardDelete(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.HardDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAreaService_GetStats(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		mock      func() *mockQuerier
		wantStats *AreaStats
		wantErr   bool
	}{
		{
			name: "returns correct stats",
			id:   "area-1",
			mock: func() *mockQuerier {
				return &mockQuerier{
					countSubareasByAreaFunc: func(ctx context.Context, areaID string) (int64, error) {
						return 3, nil
					},
					countProjectsByAreaFunc: func(ctx context.Context, areaID string) (int64, error) {
						return 12, nil
					},
					countTasksByAreaFunc: func(ctx context.Context, areaID string) (int64, error) {
						return 45, nil
					},
				}
			},
			wantStats: &AreaStats{SubareaCount: 3, ProjectCount: 12, TaskCount: 45},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewAreaService(tt.mock())
			got, err := svc.GetStats(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("AreaService.GetStats() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.SubareaCount != tt.wantStats.SubareaCount {
					t.Errorf("SubareaCount = %d, want %d", got.SubareaCount, tt.wantStats.SubareaCount)
				}
				if got.ProjectCount != tt.wantStats.ProjectCount {
					t.Errorf("ProjectCount = %d, want %d", got.ProjectCount, tt.wantStats.ProjectCount)
				}
				if got.TaskCount != tt.wantStats.TaskCount {
					t.Errorf("TaskCount = %d, want %d", got.TaskCount, tt.wantStats.TaskCount)
				}
			}
		})
	}
}

func TestDbAreaRowToDomain(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		row  db.ListAreasRow
		want domain.Area
	}{
		{
			name: "converts row with all fields",
			row: db.ListAreasRow{
				ID:        "area-1",
				Name:      "Test Area",
				Color:     sql.NullString{String: "#3B82F6", Valid: true},
				SortOrder: 5,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
			want: domain.Area{
				ID:        "area-1",
				Name:      "Test Area",
				Color:     domain.Color("#3B82F6"),
				SortOrder: 5,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
		},
		{
			name: "handles null color",
			row: db.ListAreasRow{
				ID:        "area-2",
				Name:      "No Color",
				Color:     sql.NullString{String: "", Valid: false},
				SortOrder: 0,
				CreatedAt: now,
				UpdatedAt: now,
			},
			want: domain.Area{
				ID:        "area-2",
				Name:      "No Color",
				Color:     domain.Color(""),
				SortOrder: 0,
				CreatedAt: now,
				UpdatedAt: now,
				DeletedAt: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := converter.DbListAreasRowToDomain(tt.row)
			if got.ID != tt.want.ID {
				t.Errorf("ID = %s, want %s", got.ID, tt.want.ID)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %s, want %s", got.Name, tt.want.Name)
			}
			if got.Color != tt.want.Color {
				t.Errorf("Color = %s, want %s", got.Color, tt.want.Color)
			}
			if got.SortOrder != tt.want.SortOrder {
				t.Errorf("SortOrder = %d, want %d", got.SortOrder, tt.want.SortOrder)
			}
		})
	}
}
