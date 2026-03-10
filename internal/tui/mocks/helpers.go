package mocks

import (
	"context"

	"github.com/marekbrze/dopadone/internal/domain"
	"github.com/marekbrze/dopadone/internal/service"
)

func NewMockServices() (*MockAreaService, *MockSubareaService, *MockProjectService, *MockTaskService) {
	return &MockAreaService{}, &MockSubareaService{}, &MockProjectService{}, &MockTaskService{}
}

func SetupMockAreaSuccess(areaSvc *MockAreaService, areas []domain.Area) {
	areaSvc.ListFunc = func(ctx context.Context) ([]domain.Area, error) {
		return areas, nil
	}
}

func SetupMockAreaError(areaSvc *MockAreaService, err error) {
	areaSvc.ListFunc = func(ctx context.Context) ([]domain.Area, error) {
		return nil, err
	}
}

func SetupMockSubareaSuccess(subareaSvc *MockSubareaService, subareas []domain.Subarea) {
	subareaSvc.ListByAreaFunc = func(ctx context.Context, areaID string) ([]domain.Subarea, error) {
		return subareas, nil
	}
}

func SetupMockSubareaError(subareaSvc *MockSubareaService, err error) {
	subareaSvc.ListByAreaFunc = func(ctx context.Context, areaID string) ([]domain.Subarea, error) {
		return nil, err
	}
}

func SetupMockProjectSuccess(projectSvc *MockProjectService, projects []domain.Project) {
	projectSvc.ListBySubareaFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
		return projects, nil
	}
	projectSvc.ListBySubareaRecursiveFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
		return projects, nil
	}
	projectSvc.ListAllFunc = func(ctx context.Context) ([]domain.Project, error) {
		return projects, nil
	}
}

func SetupMockProjectError(projectSvc *MockProjectService, err error) {
	projectSvc.ListBySubareaFunc = func(ctx context.Context, subareaID string) ([]domain.Project, error) {
		return nil, err
	}
}

func SetupMockTaskSuccess(taskSvc *MockTaskService, tasks []domain.Task) {
	taskSvc.ListByProjectFunc = func(ctx context.Context, projectID string) ([]domain.Task, error) {
		return tasks, nil
	}
}

func SetupMockTaskError(taskSvc *MockTaskService, err error) {
	taskSvc.ListByProjectFunc = func(ctx context.Context, projectID string) ([]domain.Task, error) {
		return nil, err
	}
}

func SetupMockAreaCreate(areaSvc *MockAreaService, area *domain.Area) {
	areaSvc.CreateFunc = func(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
		return area, nil
	}
}

func SetupMockSubareaCreate(subareaSvc *MockSubareaService, subarea *domain.Subarea) {
	subareaSvc.CreateFunc = func(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
		return subarea, nil
	}
}

func SetupMockProjectCreate(projectSvc *MockProjectService, project *domain.Project) {
	projectSvc.CreateFunc = func(ctx context.Context, params service.CreateProjectParams) (*domain.Project, error) {
		return project, nil
	}
}

func SetupMockTaskCreate(taskSvc *MockTaskService, task *domain.Task) {
	taskSvc.CreateFunc = func(ctx context.Context, params service.CreateTaskParams) (*domain.Task, error) {
		return task, nil
	}
}

func SetupMockAreaCreateError(areaSvc *MockAreaService, err error) {
	areaSvc.CreateFunc = func(ctx context.Context, name string, color domain.Color) (*domain.Area, error) {
		return nil, err
	}
}

func SetupMockSubareaCreateError(subareaSvc *MockSubareaService, err error) {
	subareaSvc.CreateFunc = func(ctx context.Context, name string, areaID string, color domain.Color) (*domain.Subarea, error) {
		return nil, err
	}
}

func SetupMockProjectCreateError(projectSvc *MockProjectService, err error) {
	projectSvc.CreateFunc = func(ctx context.Context, params service.CreateProjectParams) (*domain.Project, error) {
		return nil, err
	}
}

func SetupMockTaskCreateError(taskSvc *MockTaskService, err error) {
	taskSvc.CreateFunc = func(ctx context.Context, params service.CreateTaskParams) (*domain.Task, error) {
		return nil, err
	}
}

func SetupMockAreaUpdate(areaSvc *MockAreaService, area *domain.Area) {
	areaSvc.UpdateFunc = func(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error) {
		return area, nil
	}
}

func SetupMockAreaUpdateError(areaSvc *MockAreaService, err error) {
	areaSvc.UpdateFunc = func(ctx context.Context, id string, name string, color domain.Color) (*domain.Area, error) {
		return nil, err
	}
}

func SetupMockAreaDelete(areaSvc *MockAreaService) {
	areaSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
}

func SetupMockAreaDeleteError(areaSvc *MockAreaService, err error) {
	areaSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return err
	}
}

func SetupMockAreaHardDelete(areaSvc *MockAreaService) {
	areaSvc.HardDeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
}

func SetupMockAreaHardDeleteError(areaSvc *MockAreaService, err error) {
	areaSvc.HardDeleteFunc = func(ctx context.Context, id string) error {
		return err
	}
}

func SetupMockAreaReorder(areaSvc *MockAreaService) {
	areaSvc.ReorderAllFunc = func(ctx context.Context, areaIDs []string) error {
		return nil
	}
}

func SetupMockAreaReorderError(areaSvc *MockAreaService, err error) {
	areaSvc.ReorderAllFunc = func(ctx context.Context, areaIDs []string) error {
		return err
	}
}

func SetupMockAreaStats(areaSvc *MockAreaService, stats *service.AreaStats) {
	areaSvc.GetStatsFunc = func(ctx context.Context, id string) (*service.AreaStats, error) {
		return stats, nil
	}
}

func SetupMockAreaStatsError(areaSvc *MockAreaService, err error) {
	areaSvc.GetStatsFunc = func(ctx context.Context, id string) (*service.AreaStats, error) {
		return nil, err
	}
}

func SetupMockSubareaDelete(subareaSvc *MockSubareaService) {
	subareaSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
}

func SetupMockSubareaDeleteError(subareaSvc *MockSubareaService, err error) {
	subareaSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return err
	}
}

func SetupMockProjectDeleteWithCascade(projectSvc *MockProjectService) {
	projectSvc.SoftDeleteWithCascadeFunc = func(ctx context.Context, id string) error {
		return nil
	}
}

func SetupMockProjectDeleteWithCascadeError(projectSvc *MockProjectService, err error) {
	projectSvc.SoftDeleteWithCascadeFunc = func(ctx context.Context, id string) error {
		return err
	}
}

func SetupMockTaskDelete(taskSvc *MockTaskService) {
	taskSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return nil
	}
}

func SetupMockTaskDeleteError(taskSvc *MockTaskService, err error) {
	taskSvc.SoftDeleteFunc = func(ctx context.Context, id string) error {
		return err
	}
}
