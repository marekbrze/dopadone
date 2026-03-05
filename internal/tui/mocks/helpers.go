package mocks

import (
	"context"

	"github.com/example/projectdb/internal/domain"
	"github.com/example/projectdb/internal/service"
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
