package mocks

import (
	"context"

	"github.com/example/dopadone/internal/domain"
	"github.com/example/dopadone/internal/service"
)

type MockAreaService struct {
	ListFunc            func(ctx context.Context) ([]domain.Area, error)
	GetByIDFunc         func(ctx context.Context, id string) (*domain.Area, error)
	CreateFunc          func(ctx context.Context, name string, color domain.Color) (*domain.Area, error)
	UpdateFunc          func(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error)
	UpdateSortOrderFunc func(ctx context.Context, id string, sortOrder int) error
	ReorderAllFunc      func(ctx context.Context, areaIDs []string) error
	SoftDeleteFunc      func(ctx context.Context, id string) error
	HardDeleteFunc      func(ctx context.Context, id string) error
	GetStatsFunc        func(ctx context.Context, id string) (*service.AreaStats, error)
}

func (m *MockAreaService) List(ctx context.Context) ([]domain.Area, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []domain.Area{}, nil
}

func (m *MockAreaService) GetByID(ctx context.Context, id string) (*domain.Area, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockAreaService) Create(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, color)
	}
	return nil, nil
}

func (m *MockAreaService) Update(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, name, color)
	}
	return nil, nil
}

func (m *MockAreaService) UpdateSortOrder(ctx context.Context, id string, sortOrder int) error {
	if m.UpdateSortOrderFunc != nil {
		return m.UpdateSortOrderFunc(ctx, id, sortOrder)
	}
	return nil
}

func (m *MockAreaService) ReorderAll(ctx context.Context, areaIDs []string) error {
	if m.ReorderAllFunc != nil {
		return m.ReorderAllFunc(ctx, areaIDs)
	}
	return nil
}

func (m *MockAreaService) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFunc != nil {
		return m.SoftDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockAreaService) HardDelete(ctx context.Context, id string) error {
	if m.HardDeleteFunc != nil {
		return m.HardDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockAreaService) GetStats(ctx context.Context, id string) (*service.AreaStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx, id)
	}
	return &service.AreaStats{}, nil
}

type MockSubareaService struct {
	CreateFunc            func(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error)
	GetByIDFunc           func(ctx context.Context, id string) (*domain.Subarea, error)
	ListByAreaFunc        func(ctx context.Context, areaID string) ([]domain.Subarea, error)
	UpdateFunc            func(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error)
	SoftDeleteFunc        func(ctx context.Context, id string) error
	HardDeleteFunc        func(ctx context.Context, id string) error
	GetStatsFunc          func(ctx context.Context, id string) (*service.SubareaStats, error)
	GetEffectiveColorFunc func(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color
	ListAllFunc           func(ctx context.Context) ([]domain.Subarea, error)
}

func (m *MockSubareaService) Create(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, name, areaID, color)
	}
	return nil, nil
}

func (m *MockSubareaService) GetByID(ctx context.Context, id string) (*domain.Subarea, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockSubareaService) ListByArea(ctx context.Context, areaID string) ([]domain.Subarea, error) {
	if m.ListByAreaFunc != nil {
		return m.ListByAreaFunc(ctx, areaID)
	}
	return []domain.Subarea{}, nil
}

func (m *MockSubareaService) Update(ctx context.Context, id string, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, name, areaID, color)
	}
	return nil, nil
}

func (m *MockSubareaService) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFunc != nil {
		return m.SoftDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockSubareaService) HardDelete(ctx context.Context, id string) error {
	if m.HardDeleteFunc != nil {
		return m.HardDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockSubareaService) GetStats(ctx context.Context, id string) (*service.SubareaStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx, id)
	}
	return &service.SubareaStats{}, nil
}

func (m *MockSubareaService) GetEffectiveColor(ctx context.Context, subarea *domain.Subarea, parentArea *domain.Area) domain.Color {
	if m.GetEffectiveColorFunc != nil {
		return m.GetEffectiveColorFunc(ctx, subarea, parentArea)
	}
	return ""
}

func (m *MockSubareaService) ListAll(ctx context.Context) ([]domain.Subarea, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc(ctx)
	}
	return []domain.Subarea{}, nil
}

type MockProjectService struct {
	CreateFunc                  func(ctx context.Context, params service.CreateProjectParams) (*domain.Project, error)
	GetByIDFunc                 func(ctx context.Context, id string) (*domain.Project, error)
	ListBySubareaFunc           func(ctx context.Context, subareaID string) ([]domain.Project, error)
	ListByParentFunc            func(ctx context.Context, parentID string) ([]domain.Project, error)
	ListAllFunc                 func(ctx context.Context) ([]domain.Project, error)
	ListByStatusFunc            func(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error)
	ListByPriorityFunc          func(ctx context.Context, priority domain.Priority) ([]domain.Project, error)
	ListBySubareaRecursiveFunc  func(ctx context.Context, subareaID string) ([]domain.Project, error)
	UpdateFunc                  func(ctx context.Context, params service.UpdateProjectParams) (*domain.Project, error)
	SoftDeleteFunc              func(ctx context.Context, id string) error
	HardDeleteFunc              func(ctx context.Context, id string) error
	GetStatsFunc                func(ctx context.Context, id string) (*service.ProjectStats, error)
	ValidateParentHierarchyFunc func(ctx context.Context, parentID string, projectID string) error
}

func (m *MockProjectService) Create(ctx context.Context, params service.CreateProjectParams) (*domain.Project, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) GetByID(ctx context.Context, id string) (*domain.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectService) ListBySubarea(ctx context.Context, subareaID string) ([]domain.Project, error) {
	if m.ListBySubareaFunc != nil {
		return m.ListBySubareaFunc(ctx, subareaID)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) ListByParent(ctx context.Context, parentID string) ([]domain.Project, error) {
	if m.ListByParentFunc != nil {
		return m.ListByParentFunc(ctx, parentID)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) ListAll(ctx context.Context) ([]domain.Project, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc(ctx)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) ListByStatus(ctx context.Context, status domain.ProjectStatus) ([]domain.Project, error) {
	if m.ListByStatusFunc != nil {
		return m.ListByStatusFunc(ctx, status)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) ListByPriority(ctx context.Context, priority domain.Priority) ([]domain.Project, error) {
	if m.ListByPriorityFunc != nil {
		return m.ListByPriorityFunc(ctx, priority)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) ListBySubareaRecursive(ctx context.Context, subareaID string) ([]domain.Project, error) {
	if m.ListBySubareaRecursiveFunc != nil {
		return m.ListBySubareaRecursiveFunc(ctx, subareaID)
	}
	return []domain.Project{}, nil
}

func (m *MockProjectService) Update(ctx context.Context, params service.UpdateProjectParams) (*domain.Project, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockProjectService) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFunc != nil {
		return m.SoftDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectService) HardDelete(ctx context.Context, id string) error {
	if m.HardDeleteFunc != nil {
		return m.HardDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockProjectService) GetStats(ctx context.Context, id string) (*service.ProjectStats, error) {
	if m.GetStatsFunc != nil {
		return m.GetStatsFunc(ctx, id)
	}
	return &service.ProjectStats{}, nil
}

func (m *MockProjectService) ValidateParentHierarchy(ctx context.Context, parentID string, projectID string) error {
	if m.ValidateParentHierarchyFunc != nil {
		return m.ValidateParentHierarchyFunc(ctx, parentID, projectID)
	}
	return nil
}

type MockTaskService struct {
	CreateFunc                 func(ctx context.Context, params service.CreateTaskParams) (*domain.Task, error)
	GetByIDFunc                func(ctx context.Context, id string) (*domain.Task, error)
	ListByProjectFunc          func(ctx context.Context, projectID string) ([]domain.Task, error)
	ListByProjectRecursiveFunc func(ctx context.Context, projectID string) ([]domain.Task, error)
	GetGroupedTasksFunc        func(ctx context.Context, projectID string) (*domain.GroupedTasks, error)
	ListByStatusFunc           func(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error)
	ListByPriorityFunc         func(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error)
	ListNextFunc               func(ctx context.Context) ([]domain.Task, error)
	ListAllFunc                func(ctx context.Context) ([]domain.Task, error)
	UpdateFunc                 func(ctx context.Context, params service.UpdateTaskParams) (*domain.Task, error)
	SoftDeleteFunc             func(ctx context.Context, id string) error
	HardDeleteFunc             func(ctx context.Context, id string) error
	SetStatusFunc              func(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error)
	MarkCompletedFunc          func(ctx context.Context, id string) (*domain.Task, error)
	SetPriorityFunc            func(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error)
	ToggleIsNextFunc           func(ctx context.Context, id string) (*domain.Task, error)
}

func (m *MockTaskService) Create(ctx context.Context, params service.CreateTaskParams) (*domain.Task, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockTaskService) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskService) ListByProject(ctx context.Context, projectID string) ([]domain.Task, error) {
	if m.ListByProjectFunc != nil {
		return m.ListByProjectFunc(ctx, projectID)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) ListByProjectRecursive(ctx context.Context, projectID string) ([]domain.Task, error) {
	if m.ListByProjectRecursiveFunc != nil {
		return m.ListByProjectRecursiveFunc(ctx, projectID)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) GetGroupedTasks(ctx context.Context, projectID string) (*domain.GroupedTasks, error) {
	if m.GetGroupedTasksFunc != nil {
		return m.GetGroupedTasksFunc(ctx, projectID)
	}
	return &domain.GroupedTasks{}, nil
}

func (m *MockTaskService) ListByStatus(ctx context.Context, status domain.TaskStatus) ([]domain.Task, error) {
	if m.ListByStatusFunc != nil {
		return m.ListByStatusFunc(ctx, status)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) ListByPriority(ctx context.Context, priority domain.TaskPriority) ([]domain.Task, error) {
	if m.ListByPriorityFunc != nil {
		return m.ListByPriorityFunc(ctx, priority)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) ListNext(ctx context.Context) ([]domain.Task, error) {
	if m.ListNextFunc != nil {
		return m.ListNextFunc(ctx)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) ListAll(ctx context.Context) ([]domain.Task, error) {
	if m.ListAllFunc != nil {
		return m.ListAllFunc(ctx)
	}
	return []domain.Task{}, nil
}

func (m *MockTaskService) Update(ctx context.Context, params service.UpdateTaskParams) (*domain.Task, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, params)
	}
	return nil, nil
}

func (m *MockTaskService) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFunc != nil {
		return m.SoftDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTaskService) HardDelete(ctx context.Context, id string) error {
	if m.HardDeleteFunc != nil {
		return m.HardDeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTaskService) SetStatus(ctx context.Context, id string, status domain.TaskStatus) (*domain.Task, error) {
	if m.SetStatusFunc != nil {
		return m.SetStatusFunc(ctx, id, status)
	}
	return nil, nil
}

func (m *MockTaskService) MarkCompleted(ctx context.Context, id string) (*domain.Task, error) {
	if m.MarkCompletedFunc != nil {
		return m.MarkCompletedFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskService) SetPriority(ctx context.Context, id string, priority domain.TaskPriority) (*domain.Task, error) {
	if m.SetPriorityFunc != nil {
		return m.SetPriorityFunc(ctx, id, priority)
	}
	return nil, nil
}

func (m *MockTaskService) ToggleIsNext(ctx context.Context, id string) (*domain.Task, error) {
	if m.ToggleIsNextFunc != nil {
		return m.ToggleIsNextFunc(ctx, id)
	}
	return nil, nil
}
